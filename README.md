# lsmdb

> A persistent LSM-Tree based Key-Value Storage Engine written in Go.

![Go](https://img.shields.io/badge/Go-1.24-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-success)

## Overview

StrataDB is a lightweight storage engine built from scratch in Go to understand the internals of modern databases such as RocksDB, LevelDB, and Cassandra.

The project implements an LSM-Tree architecture that prioritizes high write throughput while maintaining efficient reads through immutable SSTables and background compaction.

---

## Features

### Current

- Write Ahead Log (WAL)
- MemTable
- SSTables
- Automatic flushing
- Background compaction
- Crash recovery
- Persistent storage
- Simple Go API

### Planned

- Bloom Filters
- Block Cache (LRU)
- Skip List MemTable
- Range Iterators
- Snapshot Support
- Compression
- Concurrent Reads/Writes
- Benchmark Suite

---

## Architecture

```
             PUT / GET

                 │
                 ▼

          Write Ahead Log
                 │
                 ▼
            MemTable (RAM)
                 │
       Flush Threshold Reached
                 ▼
         Immutable SSTable
                 │
                 ▼
        Background Compaction
                 │
                 ▼
        Optimized SSTables
```

---

## Storage Layout

```
data/

    wal.log

    sst_0001.db
    sst_0002.db
    sst_0003.db

    manifest
```

---

## Example

```go
db, _ := storage.New("./data")

db.Put("name", "Adithya")

value, _ := db.Get("name")

fmt.Println(value)

db.Delete("name")
```

---

## Why this project?

Modern databases rarely overwrite data in place.

Instead, they rely on Log Structured Merge Trees (LSM Trees), combining:

- WAL for durability
- MemTables for fast writes
- Immutable SSTables for efficient storage
- Background compaction for long-term optimization

StrataDB is an educational implementation of these concepts.

---

## Tech Stack

- Go
- File-based Storage
- LSM Trees
- Write Ahead Logging
- SSTables

---

## Roadmap

- [x] WAL
- [x] MemTable
- [x] SSTable
- [x] Compaction
- [x] Crash Recovery
- [ ] Bloom Filters
- [ ] Block Cache
- [ ] Skip List
- [ ] Compression
- [ ] Range Scan
- [ ] MVCC Snapshots
- [ ] Benchmarks

---

## Inspiration

- Google Bigtable
- LevelDB
- RocksDB
- PebbleDB
- Cassandra

---

## Learning Goals

This project explores:

- Storage Engine Design
- Database Internals
- File Systems
- Memory Management
- Performance Optimization
- Concurrent Systems

---

## License

MIT
