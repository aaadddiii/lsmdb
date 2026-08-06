package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


type SSTable struct {
	path string
}

// WriteSSTable writes entries to disk
// The caller must ensure entries are already sorted by key
func WriteSSTable(entries []Entry, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, entry := range entries {
		if _, err := fmt.Fprintf(writer, "%s %s\n", entry.Key, entry.Value); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return file.Sync()
}


func OpenSSTable(path string) *SSTable {
	return &SSTable{path: path}
}

func (s *SSTable) Get(key string) (string, bool, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return "", false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		if parts[0] == key {
			return parts[1], true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", false, err
	}

	return "", false, nil
}