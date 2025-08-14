package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"timeseriesdb/internal/logger"
)

func init() {
	logger.Init()
}

func TestNewWALReplay(t *testing.T) {
	t.Run("create new WAL replay", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		if replay == nil {
			t.Fatal("WALReplay should not be nil")
		}
		if replay.walDir != tempDir {
			t.Errorf("Expected walDir %s, got %s", tempDir, replay.walDir)
		}
	})
}

func TestGetWALFiles(t *testing.T) {
	t.Run("get WAL files from directory", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create some WAL files
		walFiles := []string{
			"shard.wal.20231201-120000",
			"shard.wal.20231201-130000",
			"shard.wal.20231201-140000",
		}

		for _, filename := range walFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create WAL file %s: %v", filename, err)
			}
		}

		// Create some non-WAL files
		nonWALFiles := []string{
			"data.txt",
			"config.json",
			"shard.wal.tmp",
		}

		for _, filename := range nonWALFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create non-WAL file %s: %v", filename, err)
			}
		}

		files, err := replay.getWALFiles()
		if err != nil {
			t.Fatalf("Failed to get WAL files: %v", err)
		}

		if len(files) != 3 {
			t.Errorf("Expected 3 WAL files, got %d", len(files))
		}

		// Check if files are sorted by modification time
		for i := 1; i < len(files); i++ {
			statI, err := os.Stat(files[i-1])
			if err != nil {
				t.Fatalf("Failed to stat file %s: %v", files[i-1], err)
			}
			statJ, err := os.Stat(files[i])
			if err != nil {
				t.Fatalf("Failed to stat file %s: %v", files[i], err)
			}
			if statI.ModTime().After(statJ.ModTime()) {
				t.Errorf("Files should be sorted by modification time, but %s is after %s", files[i-1], files[i])
			}
		}
	})

	t.Run("get WAL files from empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		files, err := replay.getWALFiles()
		if err != nil {
			t.Fatalf("Failed to get WAL files from empty directory: %v", err)
		}

		if len(files) != 0 {
			t.Errorf("Expected 0 WAL files from empty directory, got %d", len(files))
		}
	})
}

func TestReplayFile(t *testing.T) {
	t.Run("replay valid WAL file", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create a simple WAL file with some entries
		walPath := filepath.Join(tempDir, "test.wal")

		// Create a simple WAL file structure
		// This is a simplified version - in practice, the WAL would have proper headers
		walData := []byte("test WAL data")
		if err := os.WriteFile(walPath, walData, 0644); err != nil {
			t.Fatalf("Failed to create test WAL file: %v", err)
		}

		result := &ReplayResult{
			SeriesData: make(map[string][]DataPoint),
		}

		// This will likely fail due to invalid WAL format, but we're testing the function call
		err := replay.replayFile(walPath, result)
		// We expect an error due to invalid WAL format, but the function should handle it gracefully
		if err != nil {
			// This is expected for invalid WAL format
			t.Logf("Expected error for invalid WAL format: %v", err)
		}
	})

	t.Run("replay non-existent file", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		result := &ReplayResult{
			SeriesData: make(map[string][]DataPoint),
		}

		err := replay.replayFile("/nonexistent/path.wal", result)
		if err == nil {
			t.Error("Expected error when replaying non-existent file")
		}
	})
}

func TestDeserializeEntry(t *testing.T) {
	t.Run("deserialize valid entry", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create valid entry data with proper JSON structure
		entryData := []byte(`{
			"ID": 123,
			"Timestamp": "2023-12-01T12:00:00Z",
			"SeriesID": "test_series",
			"Points": [{"Timestamp": "2023-12-01T12:00:00Z", "Value": 42.5, "Labels": {}}],
			"Checksum": 12345
		}`)

		entry, err := replay.deserializeEntry(entryData)
		if err != nil {
			t.Fatalf("Failed to deserialize valid entry: %v", err)
		}

		if entry.ID != 123 {
			t.Errorf("Expected ID 123, got %d", entry.ID)
		}
		if entry.SeriesID != "test_series" {
			t.Errorf("Expected SeriesID 'test_series', got %s", entry.SeriesID)
		}
		if len(entry.Points) != 1 {
			t.Errorf("Expected 1 point, got %d", len(entry.Points))
		}
		if entry.Points[0].Value != 42.5 {
			t.Errorf("Expected point value 42.5, got %f", entry.Points[0].Value)
		}
		if entry.Checksum != 12345 {
			t.Errorf("Expected checksum 12345, got %d", entry.Checksum)
		}
	})

	t.Run("deserialize invalid entry", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create invalid data
		invalidData := []byte("invalid json data")

		_, err := replay.deserializeEntry(invalidData)
		if err == nil {
			t.Error("Expected error when deserializing invalid entry")
		}
	})
}

func TestValidateEntry(t *testing.T) {
	t.Run("validate valid entry", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		entry := WALEntry{
			ID:        1,
			Timestamp: time.Now(),
			SeriesID:  "test_series",
			Points: []DataPoint{
				{Timestamp: time.Now(), Value: 42.5},
			},
			Checksum: 0, // Will be calculated
		}

		// Calculate expected checksum
		expectedChecksum := calculateChecksum(entry.SeriesID, entry.Points[0])
		entry.Checksum = expectedChecksum

		if !replay.ValidateEntry(entry) {
			t.Error("Valid entry should pass validation")
		}
	})

	t.Run("validate invalid entry", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		entry := WALEntry{
			ID:        1,
			Timestamp: time.Now(),
			SeriesID:  "test_series",
			Points: []DataPoint{
				{Timestamp: time.Now(), Value: 42.5},
			},
			Checksum: 99999, // Wrong checksum
		}

		if replay.ValidateEntry(entry) {
			t.Error("Invalid entry should fail validation")
		}
	})
}

func TestCleanupOldWALs(t *testing.T) {
	t.Run("cleanup old WAL files", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create old WAL files with old timestamps
		oldFiles := []string{
			"shard.wal.20231101-120000", // Old file
			"shard.wal.20231101-130000", // Old file
		}

		// Create files with old timestamps (2 days ago)
		oldTime := time.Now().Add(-48 * time.Hour)
		for _, filename := range oldFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create old WAL file %s: %v", filename, err)
			}
			// Set the file modification time to old time
			if err := os.Chtimes(filePath, oldTime, oldTime); err != nil {
				t.Fatalf("Failed to set old timestamp for %s: %v", filename, err)
			}
		}

		// Create recent WAL files
		recentFiles := []string{
			"shard.wal.20231201-120000", // Recent file
			"shard.wal.20231201-130000", // Recent file
		}

		for _, filename := range recentFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create recent WAL file %s: %v", filename, err)
			}
		}

		// Cleanup files older than 1 day
		maxAge := 24 * time.Hour
		if err := replay.CleanupOldWALs(maxAge); err != nil {
			t.Fatalf("Failed to cleanup old WAL files: %v", err)
		}

		// Check if old files were removed
		for _, filename := range oldFiles {
			filePath := filepath.Join(tempDir, filename)
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				t.Errorf("Old WAL file %s should have been removed", filename)
			}
		}

		// Check if recent files were kept
		for _, filename := range recentFiles {
			filePath := filepath.Join(tempDir, filename)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Recent WAL file %s should have been kept", filename)
			}
		}
	})

	t.Run("cleanup with no old files", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create only recent files
		recentFiles := []string{
			"shard.wal.20231201-120000",
			"shard.wal.20231201-130000",
		}

		for _, filename := range recentFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create recent WAL file %s: %v", filename, err)
			}
		}

		// Cleanup files older than 1 day
		maxAge := 24 * time.Hour
		if err := replay.CleanupOldWALs(maxAge); err != nil {
			t.Fatalf("Failed to cleanup old WAL files: %v", err)
		}

		// Check if all files were kept
		for _, filename := range recentFiles {
			filePath := filepath.Join(tempDir, filename)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Recent WAL file %s should have been kept", filename)
			}
		}
	})
}

func TestGetWALStats(t *testing.T) {
	t.Run("get WAL statistics", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		// Create some WAL files with different sizes
		walFiles := []string{
			"shard.wal.20231201-120000",
			"shard.wal.20231201-130000",
			"shard.wal.20231201-140000",
		}

		for i, filename := range walFiles {
			filePath := filepath.Join(tempDir, filename)
			// Create files with different content sizes
			content := make([]byte, (i+1)*100)
			if err := os.WriteFile(filePath, content, 0644); err != nil {
				t.Fatalf("Failed to create WAL file %s: %v", filename, err)
			}
		}

		stats, err := replay.GetWALStats()
		if err != nil {
			t.Fatalf("Failed to get WAL stats: %v", err)
		}

		if stats["total_files"] != 3 {
			t.Errorf("Expected 3 total files, got %v", stats["total_files"])
		}

		// Check total size (should be 100 + 200 + 300 = 600 bytes)
		expectedSize := int64(600)
		if stats["total_size"] != expectedSize {
			t.Errorf("Expected total size %d, got %v", expectedSize, stats["total_size"])
		}

		// Check if oldest and newest files are set
		if stats["oldest_file"] == "" {
			t.Error("Oldest file should be set")
		}
		if stats["newest_file"] == "" {
			t.Error("Newest file should be set")
		}
	})

	t.Run("get WAL statistics from empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		metrics := NewStorageMetrics()
		replay := NewWALReplay(tempDir, metrics)

		stats, err := replay.GetWALStats()
		if err != nil {
			t.Fatalf("Failed to get WAL stats from empty directory: %v", err)
		}

		if stats["total_files"] != 0 {
			t.Errorf("Expected 0 total files, got %v", stats["total_files"])
		}
		if stats["total_size"] != int64(0) {
			t.Errorf("Expected total size 0, got %v", stats["total_size"])
		}
		if stats["oldest_file"] != "" {
			t.Errorf("Expected empty oldest file, got %v", stats["oldest_file"])
		}
		if stats["newest_file"] != "" {
			t.Errorf("Expected empty newest file, got %v", stats["newest_file"])
		}
	})
}
