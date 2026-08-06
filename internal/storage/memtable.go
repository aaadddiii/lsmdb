package storage

type MemTable struct {
	data map[string]Entry
}

func NewMemTable() *MemTable {
	return &MemTable{
		data: make(map[string]Entry),
	}
}

func (m *MemTable) Set(key, value string) {
	m.data[key] = Entry{Key: key, Value: value}
}

func (m *MemTable) Delete(key string) {
	m.data[key] = Entry{Key: key, Deleted: true}
}

func (m *MemTable) Get(key string) (string, bool) {
	entry, ok := m.data[key]
	if !ok || entry.Deleted {
		return "", false
	}
	return entry.Value, true
}

// Lookup returns the raw entry for key exactly as stored, including
// tombstones, and reports whether the memtable has any record of key
// at all — set or deleted. DB.Get needs this to know when to stop
// searching older tiers, which plain Get can't tell it.
func (m *MemTable) Lookup(key string) (Entry, bool) {
	entry, ok := m.data[key]
	return entry, ok
}

func (m *MemTable) Size() int {
	return len(m.data)
}

func (m *MemTable) Entries() []Entry {
	entries := make([]Entry, 0, len(m.data))
	for _, entry := range m.data {
		entries = append(entries, entry)
	}
	return entries
}