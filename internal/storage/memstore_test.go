package storage

import (
	"fmt"
	"testing"
	"time"
)

// MockWAL implements WALInterface for testing
type MockWAL struct {
	entries []WALEntry
	errors  []error
	index   int
}

func (m *MockWAL) Write(entry WALEntry) error {
	if m.index < len(m.errors) && m.errors[m.index] != nil {
		err := m.errors[m.index]
		m.index++
		return err
	}
	m.entries = append(m.entries, entry)
	m.index++
	return nil
}

func (m *MockWAL) Flush() error {
	return nil
}

func TestNewMemStore(t *testing.T) {
	t.Run("create new memstore", func(t *testing.T) {
		mockWAL := &MockWAL{}
		maxSize := int64(1024 * 1024) // 1MB

		memStore := NewMemStore(maxSize, mockWAL, nil, nil, "test_shard")

		if memStore == nil {
			t.Fatal("MemStore should not be nil")
		}
		if memStore.maxSize != maxSize {
			t.Errorf("Expected maxSize %d, got %d", maxSize, memStore.maxSize)
		}
		if memStore.wal != mockWAL {
			t.Error("WAL should be set correctly")
		}
		if memStore.onFlush != nil {
			t.Error("onFlush should be nil initially")
		}
		if memStore.memTable == nil {
			t.Fatal("MemTable should be initialized")
		}
		if memStore.memTable.MaxSize != maxSize {
			t.Errorf("Expected MemTable MaxSize %d, got %d", maxSize, memStore.memTable.MaxSize)
		}
	})

	t.Run("create memstore with flush callback", func(t *testing.T) {
		mockWAL := &MockWAL{}
		maxSize := int64(1024)
		flushCalled := false

		onFlush := func(memTable *MemTable) error {
			flushCalled = true
			return nil
		}

		memStore := NewMemStore(maxSize, mockWAL, onFlush, nil, "test_shard")

		if memStore.onFlush == nil {
			t.Error("onFlush should be set")
		}
		// The flush callback is not called during initialization, only when flushing
		if flushCalled {
			t.Error("onFlush callback should not be called during initialization")
		}
	})
}

func TestMemStoreWrite(t *testing.T) {
	t.Run("write single point", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		point := DataPoint{
			Timestamp: time.Now(),
			Value:     42.5,
		}

		err := memStore.Write("test_series", []DataPoint{point})
		if err != nil {
			t.Fatalf("Failed to write point: %v", err)
		}

		// Check if point was added to memtable
		memTable := memStore.GetMemTable()
		if len(memTable.Data["test_series"]) != 1 {
			t.Errorf("Expected 1 point in series, got %d", len(memTable.Data["test_series"]))
		}

		// Check if WAL was written to
		if len(mockWAL.entries) != 1 {
			t.Errorf("Expected 1 WAL entry, got %d", len(mockWAL.entries))
		}

		// Check WAL entry content
		entry := mockWAL.entries[0]
		if entry.SeriesID != "test_series" {
			t.Errorf("Expected SeriesID 'test_series', got %s", entry.SeriesID)
		}
		if len(entry.Points) != 1 {
			t.Errorf("Expected 1 point in WAL entry, got %d", len(entry.Points))
		}
		if entry.Points[0].Value != 42.5 {
			t.Errorf("Expected point value 42.5, got %f", entry.Points[0].Value)
		}
	})

	t.Run("write multiple points", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		points := []DataPoint{
			{Timestamp: time.Now(), Value: 1.0},
			{Timestamp: time.Now().Add(time.Second), Value: 2.0},
			{Timestamp: time.Now().Add(2 * time.Second), Value: 3.0},
		}

		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		// Check if all points were added
		memTable := memStore.GetMemTable()
		if len(memTable.Data["test_series"]) != 3 {
			t.Errorf("Expected 3 points in series, got %d", len(memTable.Data["test_series"]))
		}

		// Check if WAL was written to for each point
		if len(mockWAL.entries) != 3 {
			t.Errorf("Expected 3 WAL entries, got %d", len(mockWAL.entries))
		}
	})

	t.Run("write to multiple series", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		// Write to first series
		err := memStore.Write("series1", []DataPoint{{Timestamp: time.Now(), Value: 1.0}})
		if err != nil {
			t.Fatalf("Failed to write to series1: %v", err)
		}

		// Write to second series
		err = memStore.Write("series2", []DataPoint{{Timestamp: time.Now(), Value: 2.0}})
		if err != nil {
			t.Fatalf("Failed to write to series2: %v", err)
		}

		memTable := memStore.GetMemTable()
		if len(memTable.Data) != 2 {
			t.Errorf("Expected 2 series, got %d", len(memTable.Data))
		}
		if len(memTable.Data["series1"]) != 1 {
			t.Errorf("Expected 1 point in series1, got %d", len(memTable.Data["series1"]))
		}
		if len(memTable.Data["series2"]) != 1 {
			t.Errorf("Expected 1 point in series2, got %d", len(memTable.Data["series2"]))
		}
	})

	t.Run("write with WAL error", func(t *testing.T) {
		mockWAL := &MockWAL{
			errors: []error{fmt.Errorf("WAL write failed")},
		}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		point := DataPoint{
			Timestamp: time.Now(),
			Value:     42.5,
		}

		err := memStore.Write("test_series", []DataPoint{point})
		if err == nil {
			t.Error("Expected error when WAL write fails")
		}
		if err.Error() != "WAL write failed" {
			t.Errorf("Expected error 'WAL write failed', got %v", err)
		}
	})
}

func TestMemStoreRead(t *testing.T) {
	t.Run("read from empty memstore", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		points, err := memStore.Read("nonexistent_series", time.Now(), time.Now().Add(time.Hour))
		if err != nil {
			t.Fatalf("Read should not error for empty series: %v", err)
		}
		if len(points) != 0 {
			t.Errorf("Expected 0 points, got %d", len(points))
		}
	})

	t.Run("read from populated memstore", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		now := time.Now()
		points := []DataPoint{
			{Timestamp: now, Value: 1.0},
			{Timestamp: now.Add(time.Second), Value: 2.0},
			{Timestamp: now.Add(2 * time.Second), Value: 3.0},
		}

		// Write points
		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		// Read all points
		readPoints, err := memStore.Read("test_series", now, now.Add(3*time.Second))
		if err != nil {
			t.Fatalf("Failed to read points: %v", err)
		}
		if len(readPoints) != 3 {
			t.Errorf("Expected 3 points, got %d", len(readPoints))
		}

		// Read with time range
		readPoints, err = memStore.Read("test_series", now.Add(time.Second), now.Add(2*time.Second))
		if err != nil {
			t.Fatalf("Failed to read points with time range: %v", err)
		}
		if len(readPoints) != 2 {
			t.Errorf("Expected 2 points in time range, got %d", len(readPoints))
		}
		if readPoints[0].Value != 2.0 {
			t.Errorf("Expected first point value 2.0, got %f", readPoints[0].Value)
		}
		if readPoints[1].Value != 3.0 {
			t.Errorf("Expected second point value 3.0, got %f", readPoints[1].Value)
		}
	})

	t.Run("read with exact time boundaries", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		now := time.Now()
		points := []DataPoint{
			{Timestamp: now, Value: 1.0},
			{Timestamp: now.Add(time.Second), Value: 2.0},
		}

		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		// Read with exact start and end times
		readPoints, err := memStore.Read("test_series", now, now.Add(time.Second))
		if err != nil {
			t.Fatalf("Failed to read points: %v", err)
		}
		if len(readPoints) != 2 {
			t.Errorf("Expected 2 points, got %d", len(readPoints))
		}
	})
}

func TestMemStoreFlush(t *testing.T) {
	t.Run("auto-flush when size limit reached", func(t *testing.T) {
		mockWAL := &MockWAL{}
		flushCalled := false
		var flushedMemTable *MemTable

		onFlush := func(memTable *MemTable) error {
			flushCalled = true
			flushedMemTable = memTable
			return nil
		}

		// Create memstore with small max size to trigger flush
		memStore := NewMemStore(100, mockWAL, onFlush, nil, "test_shard")

		// Write enough data to trigger flush
		points := []DataPoint{
			{Timestamp: time.Now(), Value: 1.0},
			{Timestamp: time.Now().Add(time.Second), Value: 2.0},
		}

		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		if !flushCalled {
			t.Error("Flush callback should have been called")
		}
		if flushedMemTable == nil {
			t.Error("Flushed memtable should not be nil")
		}
		if !flushedMemTable.IsFlushed {
			t.Error("Flushed memtable should be marked as flushed")
		}

		// Check if new memtable was created
		newMemTable := memStore.GetMemTable()
		if newMemTable.ID == flushedMemTable.ID {
			t.Error("New memtable should have different ID")
		}
		if newMemTable.IsFlushed {
			t.Error("New memtable should not be marked as flushed")
		}
	})

	t.Run("flush with error", func(t *testing.T) {
		mockWAL := &MockWAL{}
		onFlush := func(memTable *MemTable) error {
			return fmt.Errorf("flush failed")
		}

		memStore := NewMemStore(100, mockWAL, onFlush, nil, "test_shard")

		// Write enough points to trigger flush (each point is ~64 bytes, so 2 points will exceed 100 bytes)
		points := []DataPoint{
			{Timestamp: time.Now(), Value: 1.0},
			{Timestamp: time.Now().Add(time.Second), Value: 2.0},
		}

		err := memStore.Write("test_series", points)
		if err == nil {
			t.Error("Expected error when flush fails")
		}
		if err.Error() != "flush failed" {
			t.Errorf("Expected error 'flush failed', got %v", err)
		}
	})

	t.Run("force flush", func(t *testing.T) {
		mockWAL := &MockWAL{}
		flushCalled := false

		onFlush := func(memTable *MemTable) error {
			flushCalled = true
			return nil
		}

		memStore := NewMemStore(1024*1024, mockWAL, onFlush, nil, "test_shard")

		// Force flush
		err := memStore.ForceFlush()
		if err != nil {
			t.Fatalf("Force flush failed: %v", err)
		}

		if !flushCalled {
			t.Error("Flush callback should have been called")
		}
	})
}

func TestMemStoreGetters(t *testing.T) {
	t.Run("get memtable", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		memTable := memStore.GetMemTable()
		if memTable == nil {
			t.Fatal("GetMemTable should not return nil")
		}
		if memTable.MaxSize != 1024*1024 {
			t.Errorf("Expected MaxSize %d, got %d", 1024*1024, memTable.MaxSize)
		}
	})

	t.Run("get size", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		// Initial size should be 0
		size := memStore.GetSize()
		if size != 0 {
			t.Errorf("Expected initial size 0, got %d", size)
		}

		// Write some data
		points := []DataPoint{
			{Timestamp: time.Now(), Value: 1.0},
			{Timestamp: time.Now().Add(time.Second), Value: 2.0},
		}

		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		// Size should have increased
		newSize := memStore.GetSize()
		if newSize <= size {
			t.Errorf("Expected size to increase, got %d (was %d)", newSize, size)
		}
	})
}

func TestMemStoreSizeCalculation(t *testing.T) {
	t.Run("size calculation with points", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		// Each point is estimated to be 64 bytes
		points := []DataPoint{
			{Timestamp: time.Now(), Value: 1.0},
			{Timestamp: time.Now().Add(time.Second), Value: 2.0},
		}

		err := memStore.Write("test_series", points)
		if err != nil {
			t.Fatalf("Failed to write points: %v", err)
		}

		expectedSize := int64(len(points) * 64)
		actualSize := memStore.GetSize()

		if actualSize != expectedSize {
			t.Errorf("Expected size %d, got %d", expectedSize, actualSize)
		}
	})

	t.Run("size calculation with multiple series", func(t *testing.T) {
		mockWAL := &MockWAL{}
		memStore := NewMemStore(1024*1024, mockWAL, nil, nil, "test_shard")

		// Write to multiple series
		series1Points := []DataPoint{{Timestamp: time.Now(), Value: 1.0}}
		series2Points := []DataPoint{{Timestamp: time.Now(), Value: 2.0}}

		err := memStore.Write("series1", series1Points)
		if err != nil {
			t.Fatalf("Failed to write to series1: %v", err)
		}

		err = memStore.Write("series2", series2Points)
		if err != nil {
			t.Fatalf("Failed to write to series2: %v", err)
		}

		expectedSize := int64((len(series1Points) + len(series2Points)) * 64)
		actualSize := memStore.GetSize()

		if actualSize != expectedSize {
			t.Errorf("Expected size %d, got %d", expectedSize, actualSize)
		}
	})
}
