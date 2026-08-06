package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const maxMemTableSize = 1000

type DB struct {
	dir        string
	walPath    string
	memtable   *MemTable
	wal        *WAL
	sstables   []*SSTable
	generation int
}

func New(path string) (*DB, error) {
	dir := filepath.Dir(path)

	sstables, maxGen, err := loadSSTables(dir)
	if err != nil {
		return nil, err
	}

	wal, err := NewWAL(path)
	if err != nil {
		return nil, err
	}

	data, err := wal.Replay()
	if err != nil {
		return nil, err
	}

	memtable := NewMemTable()
	for key, value := range data {
		memtable.Set(key, value)
	}

	return &DB{
		dir:        dir,
		walPath:    path,
		memtable:   memtable,
		wal:        wal,
		sstables:   sstables,
		generation: maxGen,
	}, nil
}

// loadSSTables scans dir for existing sst_NNNN.dat files, opens each one,
// and returns them ordered oldest -> newest, along with the highest
// generation number found (0 if none exist).
func loadSSTables(dir string) ([]*SSTable, int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}

	type found struct {
		gen  int
		path string
	}

	var matches []found

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "sst_") || !strings.HasSuffix(name, ".dat") {
			continue
		}

		genStr := strings.TrimSuffix(strings.TrimPrefix(name, "sst_"), ".dat")
		gen, err := strconv.Atoi(genStr)
		if err != nil {
			continue
		}

		matches = append(matches, found{gen: gen, path: filepath.Join(dir, name)})
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].gen < matches[j].gen
	})

	var sstables []*SSTable
	maxGen := 0
	for _, m := range matches {
		sstables = append(sstables, OpenSSTable(m.path))
		if m.gen > maxGen {
			maxGen = m.gen
		}
	}

	return sstables, maxGen, nil
}

func (db *DB) Set(key, value string) error {
	if err := db.wal.WriteSet(key, value); err != nil {
		return err
	}

	db.memtable.Set(key, value)

	if db.memtable.Size() >= maxMemTableSize {
		return db.Flush()
	}

	return nil
}

func (db *DB) Delete(key string) error {
	if err := db.wal.WriteDelete(key); err != nil {
		return err
	}

	db.memtable.Delete(key)
	return nil
}

func (db *DB) Get(key string) (string, bool, error) {
	if value, ok := db.memtable.Get(key); ok {
		return value, true, nil
	}

	for i := len(db.sstables) - 1; i >= 0; i-- {
		value, ok, err := db.sstables[i].Get(key)
		if err != nil {
			return "", false, err
		}
		if ok {
			return value, true, nil
		}
	}

	return "", false, nil
}

// Flush writes the memtable to a new immutable SSTable, then rotates the
// WAL, and only then swaps in a fresh memtable and registers the SSTable.
// In-memory state is never cleared until the durable state has fully
// transitioned, so a crash mid-flush leaves the database in a state that
// recovery can still reconstruct correctly.
func (db *DB) Flush() error {
	entries := db.memtable.Entries()
	SortEntries(entries)

	nextGen := db.generation + 1
	sstPath := filepath.Join(db.dir, fmt.Sprintf("sst_%04d.dat", nextGen))

	if err := WriteSSTable(entries, sstPath); err != nil {
		return err
	}

	if err := db.rotateWAL(); err != nil {
		return err
	}

	db.generation = nextGen
	db.memtable = NewMemTable()
	db.sstables = append(db.sstables, OpenSSTable(sstPath))

	return nil
}

func (db *DB) rotateWAL() error {
	if err := db.wal.Close(); err != nil {
		return err
	}

	if err := os.Remove(db.walPath); err != nil {
		return err
	}

	wal, err := NewWAL(db.walPath)
	if err != nil {
		return err
	}

	db.wal = wal
	return nil
}
