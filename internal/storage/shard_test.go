package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewShard(t *testing.T) {
	t.Run("create new shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		if shard == nil {
			t.Fatal("Shard should not be nil")
		}
		if shard.id != "test_shard" {
			t.Errorf("Expected ID 'test_shard', got %s", shard.id)
		}
		if shard.dataDir != tempDir {
			t.Errorf("Expected dataDir %s, got %s", tempDir, shard.dataDir)
		}
		if shard.closed {
			t.Error("Shard should not be closed initially")
		}
		if shard.recovering {
			t.Error("Shard should not be recovering initially")
		}

		// Check if directories were created
		walDir := filepath.Join(tempDir, "wal")
		segmentsDir := filepath.Join(tempDir, "segments")

		if _, err := os.Stat(walDir); os.IsNotExist(err) {
			t.Error("WAL directory should have been created")
		}
		if _, err := os.Stat(segmentsDir); os.IsNotExist(err) {
			t.Error("Segments directory should have been created")
		}
	})

	t.Run("create shard with non-existent parent directory", func(t *testing.T) {
		tempDir := t.TempDir()
		dataDir := filepath.Join(tempDir, "nonexistent", "shard")

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             dataDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		_, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Directory should be created automatically
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			t.Error("Data directory should have been created")
		}
	})
}

func TestShardOpenClose(t *testing.T) {
	t.Run("open and close shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  100 * time.Millisecond, // Fast for testing
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Open shard
		if err := shard.Open(); err != nil {
			t.Fatalf("Failed to open shard: %v", err)
		}

		// Check if compaction manager is running
		if !shard.compactionMgr.running {
			t.Error("Compaction manager should be running after open")
		}

		// Close shard
		if err := shard.Close(); err != nil {
			t.Fatalf("Failed to close shard: %v", err)
		}

		if !shard.closed {
			t.Error("Shard should be marked as closed")
		}

		// Try to open closed shard
		if err := shard.Open(); err == nil {
			t.Error("Expected error when opening closed shard")
		}
	})

	t.Run("close already closed shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Close shard
		if err := shard.Close(); err != nil {
			t.Fatalf("Failed to close shard: %v", err)
		}

		// Close again should not error
		if err := shard.Close(); err != nil {
			t.Errorf("Second close should not error: %v", err)
		}
	})
}

func TestShardWrite(t *testing.T) {
	t.Run("write data points", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  100 * time.Millisecond,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		// Open shard
		if err := shard.Open(); err != nil {
			t.Fatalf("Failed to open shard: %v", err)
		}

		// Write data points
		req := WriteRequest{
			SeriesID: "test_series",
			Points: []DataPoint{
				{Timestamp: time.Now(), Value: 42.5},
				{Timestamp: time.Now().Add(time.Second), Value: 43.0},
			},
		}

		if err := shard.Write(req); err != nil {
			t.Fatalf("Failed to write data: %v", err)
		}

		// Check if data was written to memstore
		memTable := shard.memStore.GetMemTable()
		if len(memTable.Data["test_series"]) != 2 {
			t.Errorf("Expected 2 points in memstore, got %d", len(memTable.Data["test_series"]))
		}
	})

	t.Run("write to closed shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Close shard
		shard.Close()

		req := WriteRequest{
			SeriesID: "test_series",
			Points:   []DataPoint{{Timestamp: time.Now(), Value: 42.5}},
		}

		if err := shard.Write(req); err == nil {
			t.Error("Expected error when writing to closed shard")
		}
	})

	t.Run("write to recovering shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		// Manually set recovering state
		shard.recovering = true

		req := WriteRequest{
			SeriesID: "test_series",
			Points:   []DataPoint{{Timestamp: time.Now(), Value: 42.5}},
		}

		if err := shard.Write(req); err == nil {
			t.Error("Expected error when writing to recovering shard")
		}
	})
}

func TestShardRead(t *testing.T) {
	t.Run("read data points", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  100 * time.Millisecond,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		// Open shard
		if err := shard.Open(); err != nil {
			t.Fatalf("Failed to open shard: %v", err)
		}

		// Write some data
		now := time.Now()
		req := WriteRequest{
			SeriesID: "test_series",
			Points: []DataPoint{
				{Timestamp: now, Value: 42.5},
				{Timestamp: now.Add(time.Second), Value: 43.0},
				{Timestamp: now.Add(2 * time.Second), Value: 43.5},
			},
		}

		if err := shard.Write(req); err != nil {
			t.Fatalf("Failed to write data: %v", err)
		}

		// Read data
		readReq := ReadRequest{
			SeriesID: "test_series",
			Start:    now,
			End:      now.Add(3 * time.Second),
			Limit:    10,
		}

		points, err := shard.Read(readReq)
		if err != nil {
			t.Fatalf("Failed to read data: %v", err)
		}

		if len(points) != 3 {
			t.Errorf("Expected 3 points, got %d", len(points))
		}

		// Check point values
		if points[0].Value != 42.5 {
			t.Errorf("Expected first point value 42.5, got %f", points[0].Value)
		}
		if points[1].Value != 43.0 {
			t.Errorf("Expected second point value 43.0, got %f", points[1].Value)
		}
		if points[2].Value != 43.5 {
			t.Errorf("Expected third point value 43.5, got %f", points[2].Value)
		}
	})

	t.Run("read with time range", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  100 * time.Millisecond,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		// Open shard
		if err := shard.Open(); err != nil {
			t.Fatalf("Failed to open shard: %v", err)
		}

		// Write data
		now := time.Now()
		req := WriteRequest{
			SeriesID: "test_series",
			Points: []DataPoint{
				{Timestamp: now, Value: 42.5},
				{Timestamp: now.Add(time.Hour), Value: 43.0},
				{Timestamp: now.Add(2 * time.Hour), Value: 43.5},
			},
		}

		if err := shard.Write(req); err != nil {
			t.Fatalf("Failed to write data: %v", err)
		}

		// Read with time range
		readReq := ReadRequest{
			SeriesID: "test_series",
			Start:    now.Add(30 * time.Minute),
			End:      now.Add(90 * time.Minute),
			Limit:    10,
		}

		points, err := shard.Read(readReq)
		if err != nil {
			t.Fatalf("Failed to read data: %v", err)
		}

		// Should only get the point at 1 hour (within range)
		if len(points) != 1 {
			t.Errorf("Expected 1 point in time range, got %d", len(points))
		}
		if points[0].Value != 43.0 {
			t.Errorf("Expected point value 43.0, got %f", points[0].Value)
		}
	})

	t.Run("read from closed shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Close shard
		shard.Close()

		readReq := ReadRequest{
			SeriesID: "test_series",
			Start:    time.Now(),
			End:      time.Now().Add(time.Hour),
			Limit:    10,
		}

		_, err = shard.Read(readReq)
		if err == nil {
			t.Error("Expected error when reading from closed shard")
		}
	})
}

func TestShardStats(t *testing.T) {
	t.Run("get shard statistics", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  100 * time.Millisecond,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		// Open shard
		if err := shard.Open(); err != nil {
			t.Fatalf("Failed to open shard: %v", err)
		}

		stats := shard.GetStats()

		// Check basic stats
		if stats["id"] != "test_shard" {
			t.Errorf("Expected ID 'test_shard', got %v", stats["id"])
		}
		if stats["closed"] != false {
			t.Errorf("Expected closed false, got %v", stats["closed"])
		}
		if stats["recovering"] != false {
			t.Errorf("Expected recovering false, got %v", stats["recovering"])
		}

		// Check memstore size
		memstoreSize, exists := stats["memstore_size"]
		if !exists {
			t.Error("memstore_size should exist in stats")
		}
		if memstoreSize.(int64) < 0 {
			t.Error("memstore_size should be non-negative")
		}

		// Check WAL size
		walSize, exists := stats["wal_size"]
		if !exists {
			t.Error("wal_size should exist in stats")
		}
		if walSize.(int64) < 0 {
			t.Error("wal_size should be non-negative")
		}

		// Check compaction stats
		compactionStats, exists := stats["compaction_stats"]
		if !exists {
			t.Error("compaction_stats should exist in stats")
		}
		if compactionStats == nil {
			t.Error("compaction_stats should not be nil")
		}
	})

	t.Run("get stats from closed shard", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		// Close shard
		shard.Close()

		stats := shard.GetStats()

		if stats["closed"] != true {
			t.Errorf("Expected closed true, got %v", stats["closed"])
		}
	})
}

func TestShardGetters(t *testing.T) {
	t.Run("get shard ID", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}
		defer shard.Close()

		if shard.GetID() != "test_shard" {
			t.Errorf("Expected ID 'test_shard', got %s", shard.GetID())
		}
	})

	t.Run("check if shard is closed", func(t *testing.T) {
		tempDir := t.TempDir()

		config := ShardConfig{
			ID:                  "test_shard",
			DataDir:             tempDir,
			MaxMemTableSize:     1024 * 1024,
			MaxWALSize:          64 * 1024,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
		}

		shard, err := NewShard(config)
		if err != nil {
			t.Fatalf("Failed to create shard: %v", err)
		}

		if shard.IsClosed() {
			t.Error("Shard should not be closed initially")
		}

		shard.Close()

		if !shard.IsClosed() {
			t.Error("Shard should be closed after Close()")
		}
	})
}
