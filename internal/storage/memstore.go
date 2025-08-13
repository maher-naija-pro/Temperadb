package storage

import (
	"sync"
	"time"
)

// MemStore represents the in-memory storage layer
type MemStore struct {
	mu       sync.RWMutex
	memTable *MemTable
	maxSize  int64
	wal      WALInterface
	onFlush  func(*MemTable) error
}

// WALInterface defines the interface for WAL operations
type WALInterface interface {
	Write(entry WALEntry) error
	Flush() error
}

// NewMemStore creates a new memory store
func NewMemStore(maxSize int64, wal WALInterface, onFlush func(*MemTable) error) *MemStore {
	return &MemStore{
		memTable: &MemTable{
			ID:        uint64(time.Now().UnixNano()),
			Data:      make(map[string][]DataPoint),
			Size:      0,
			MaxSize:   maxSize,
			CreatedAt: time.Now(),
			IsFlushed: false,
		},
		maxSize: maxSize,
		wal:     wal,
		onFlush: onFlush,
	}
}

// Write writes data points to the memory store
func (ms *MemStore) Write(seriesID string, points []DataPoint) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Add points to memtable
	if ms.memTable.Data[seriesID] == nil {
		ms.memTable.Data[seriesID] = make([]DataPoint, 0)
	}

	ms.memTable.Data[seriesID] = append(ms.memTable.Data[seriesID], points...)

	// Update size estimate (rough calculation)
	ms.memTable.Size += int64(len(points) * 64) // Approximate size per point

	// Check if we need to flush the current memtable after updating size
	if ms.memTable.Size >= ms.maxSize {
		if err := ms.flushMemTable(); err != nil {
			return err
		}
	}

	// Write to WAL for durability
	for _, point := range points {
		entry := WALEntry{
			ID:        uint64(time.Now().UnixNano()),
			Timestamp: time.Now(),
			SeriesID:  seriesID,
			Points:    []DataPoint{point},
			Checksum:  calculateChecksum(seriesID, point),
		}

		if err := ms.wal.Write(entry); err != nil {
			return err
		}
	}

	return nil
}

// Read reads data points from the memory store
func (ms *MemStore) Read(seriesID string, start, end time.Time) ([]DataPoint, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if ms.memTable.Data[seriesID] == nil {
		return []DataPoint{}, nil
	}

	var result []DataPoint
	for _, point := range ms.memTable.Data[seriesID] {
		if (point.Timestamp.Equal(start) || point.Timestamp.After(start)) &&
			(point.Timestamp.Equal(end) || point.Timestamp.Before(end)) {
			result = append(result, point)
		}
	}

	return result, nil
}

// GetMemTable returns the current memtable
func (ms *MemStore) GetMemTable() *MemTable {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.memTable
}

// flushMemTable flushes the current memtable and creates a new one
func (ms *MemStore) flushMemTable() error {
	if ms.onFlush != nil {
		if err := ms.onFlush(ms.memTable); err != nil {
			return err
		}
	}

	// Mark current memtable as flushed
	ms.memTable.IsFlushed = true

	// Create new memtable
	ms.memTable = &MemTable{
		ID:        uint64(time.Now().UnixNano()),
		Data:      make(map[string][]DataPoint),
		Size:      0,
		MaxSize:   ms.maxSize,
		CreatedAt: time.Now(),
		IsFlushed: false,
	}

	return nil
}

// ForceFlush forces a flush of the current memtable
func (ms *MemStore) ForceFlush() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return ms.flushMemTable()
}

// GetSize returns the current size of the memtable
func (ms *MemStore) GetSize() int64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.memTable.Size
}

// calculateChecksum calculates a simple checksum for data integrity
func calculateChecksum(seriesID string, point DataPoint) uint32 {
	// Simple checksum implementation
	var sum uint32
	for _, b := range []byte(seriesID) {
		sum += uint32(b)
	}
	sum += uint32(point.Timestamp.Unix())
	sum += uint32(point.Value * 1000) // Convert float to int for checksum
	return sum
}
