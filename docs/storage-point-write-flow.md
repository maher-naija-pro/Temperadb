# Storage Point Write Flow Documentation

## Overview

This document explains how time-series data points flow through the storage system from initial write request to persistent storage on disk. The system uses an LSM (Log-Structured Merge) tree architecture to provide both high write performance and data durability.

## Architecture Components

The storage system consists of several key components that work together:

- **Storage Engine**: The main coordinator that routes writes to appropriate shards
- **Shards**: Individual storage units that handle data for specific partitions
- **MemStore**: In-memory buffer for recent writes
- **WAL (Write-Ahead Log)**: Durability mechanism that logs all writes
- **Segments**: Immutable on-disk files containing compacted data
- **Compaction Manager**: Background process that merges and optimizes segments

## Point Write Journey

```
Client Write Request
         │
         ▼
   Storage Engine
         │
         ▼
    Shard Router
         │
         ▼
      MemStore ──────────┐
         │               │
         ▼               │
   Write-Ahead Log       │
         │               │
         ▼               │
   (Durability)          │
                         │
         ┌───────────────┘
         │
         ▼
   Size Check
         │
         ▼
   Flush Trigger
         │
         ▼
   Segment Writer
         │
         ▼
   Immutable Segment
         │
         ▼
   Compaction Manager
         │
         ▼
   Optimized Storage
```

### 1. Entry Point - Storage Layer

When a client writes a time-series point, the request first reaches the main Storage engine. The Storage engine acts as a coordinator, determining which shard should handle the data based on the measurement name and current shard distribution.

The Storage engine converts the incoming point format into the internal storage format, creating a unique series identifier that combines the measurement name, tags, and field name. This series ID ensures that related data points are grouped together efficiently.

### 2. Shard Routing

Each shard is responsible for a subset of the data, allowing the system to distribute load and scale horizontally. The Storage engine either routes the write to an existing shard or creates a new one if needed.

Shards are identified by unique IDs and contain their own complete storage stack, including memory buffers, write-ahead logs, and segment files. This isolation ensures that operations on one shard don't interfere with others.

### 3. Memory Storage - MemStore

Once the write reaches the appropriate shard, it's immediately stored in the MemStore, which is an in-memory table that buffers recent writes. The MemStore provides extremely fast write performance since it only involves memory operations.

The MemStore maintains a size estimate for the current data and automatically triggers a flush operation when the configured size limit is reached. This ensures that memory usage remains bounded and that data eventually makes its way to persistent storage.

### 4. Durability - Write-Ahead Log

Simultaneously with writing to memory, every point is also written to the Write-Ahead Log (WAL). The WAL is a sequential log file that records all write operations before they're acknowledged to the client.

This WAL provides durability guarantees - even if the system crashes immediately after a write, the data can be recovered by replaying the WAL entries. The WAL files are rotated when they reach a certain size to prevent them from growing indefinitely.

### 5. Flush to Segments

When the MemStore reaches its size limit, it triggers a flush operation. During the flush, all the buffered data is written to an immutable segment file on disk. The segment file contains a header with metadata and the actual data organized by series.

Segment files are immutable once written, which simplifies concurrent access and provides a stable foundation for the compaction process. Each segment file is identified by a unique ID and contains data for a specific time range.

### 6. Segment File Structure

Segment files follow a structured binary format that optimizes both read and write performance. The file begins with a header containing metadata such as the segment ID, creation timestamp, number of series, and time range.

Following the header, the file contains series data organized by series ID. Each series section includes a header with the series identifier and point count, followed by the actual data points. The binary format uses length-prefixed encoding to enable efficient parsing and random access.

### 7. Background Compaction

After segments are written, they become part of the LSM tree structure. The Compaction Manager runs in the background, continuously merging smaller segments into larger ones and organizing them into levels based on size and age.

Compaction serves multiple purposes: it reduces the number of files that need to be read during queries, it removes duplicate or obsolete data, and it optimizes the storage layout for better read performance. The compaction process is designed to be non-blocking, allowing writes to continue while background optimization occurs.

## Performance Characteristics

The write path is optimized for maximum throughput:

- **Immediate Durability**: WAL ensures data is safe even with system crashes
- **Batched Writes**: Multiple points are grouped together in memory before flushing
- **Sequential I/O**: Both WAL and segment writes use sequential access patterns
- **Non-blocking Compaction**: Background optimization doesn't interfere with writes
- **Memory Buffering**: Recent writes are served from fast memory storage

## Recovery and Durability

The system provides strong durability guarantees through the combination of WAL and segment storage. If the system crashes, it can recover by:

1. Reading the latest segment files to restore the persistent data
2. Replaying the WAL to restore any writes that were in memory but not yet flushed
3. Reconstructing the MemStore state from the recovered data

This recovery process ensures that no acknowledged writes are lost, even in the event of unexpected system failures.

## Summary

The point write flow demonstrates the LSM tree architecture's ability to provide both high write performance and strong durability guarantees. By buffering writes in memory, logging them to a WAL, and periodically flushing to immutable segments, the system achieves optimal performance while maintaining data safety.

The background compaction process ensures that the storage remains efficient over time, automatically organizing data for optimal read performance without impacting write throughput. This design makes the system suitable for high-throughput time-series workloads where both write performance and data reliability are critical.
