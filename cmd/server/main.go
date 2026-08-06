package main

import (
	"fmt"

	"kvdb/internal/storage"
)

func main() {
	memtable := storage.NewMemTable()

	memtable.Set("zebra", "1")
	memtable.Set("apple", "2")
	memtable.Set("mango", "3")
	memtable.Set("banana", "4")

	entries := memtable.Entries()

	fmt.Println("Before sorting:")
	for _, e := range entries {
		fmt.Printf("Key: %s, Value: %s\n", e.Key, e.Value)
	}

	storage.SortEntries(entries)

	fmt.Println("After sorting:")
	for _, e := range entries {
		fmt.Printf("Key: %s, Value: %s\n", e.Key, e.Value)
	}
}
