package storage

import "sort"

type Entry struct {
	Key   string
	Value string
	Deleted bool
}

func SortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	}) 
}