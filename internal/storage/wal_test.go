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

func TestNewWAL(t *testing.T) {
	t.Run("create new WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024 * 1024, // 1MB
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		if wal == nil {
			t.Fatal("WAL should not be nil")
		}
		if wal.path != walPath {
			t.Errorf("Expected path %s, got %s", walPath, wal.path)
		}
		if wal.maxFileSize != 1024*1024 {
			t.Errorf("Expected max file size %d, got %d", 1024*1024, wal.maxFileSize)
		}
		if wal.closed {
			t.Error("WAL should not be closed initially")
		}
	})

	t.Run("create WAL with non-existent directory", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "nonexistent", "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		// Directory should be created automatically
		if _, err := os.Stat(filepath.Dir(walPath)); os.IsNotExist(err) {
			t.Error("Directory should have been created")
		}
	})
}

func TestWALWrite(t *testing.T) {
	t.Run("write single entry", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024 * 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		entry := WALEntry{
			Timestamp: time.Now(),
			SeriesID:  "test_series",
			Points: []DataPoint{
				{Timestamp: time.Now(), Value: 42.5},
			},
		}

		if err := wal.Write(entry); err != nil {
			t.Fatalf("Failed to write entry: %v", err)
		}

		// Check file size increased
		if wal.currentSize == 0 {
			t.Error("File size should have increased after write")
		}
	})

	t.Run("write multiple entries", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024 * 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		entries := []WALEntry{
			{
				Timestamp: time.Now(),
				SeriesID:  "series1",
				Points:    []DataPoint{{Timestamp: time.Now(), Value: 1.0}},
			},
			{
				Timestamp: time.Now(),
				SeriesID:  "series2",
				Points:    []DataPoint{{Timestamp: time.Now(), Value: 2.0}},
			},
		}

		for _, entry := range entries {
			if err := wal.Write(entry); err != nil {
				t.Fatalf("Failed to write entry: %v", err)
			}
		}

		// Check file size increased
		if wal.currentSize == 0 {
			t.Error("File size should have increased after writes")
		}
	})

	t.Run("write to closed WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}

		wal.Close()

		entry := WALEntry{
			Timestamp: time.Now(),
			SeriesID:  "test_series",
			Points:    []DataPoint{{Timestamp: time.Now(), Value: 42.5}},
		}

		if err := wal.Write(entry); err == nil {
			t.Error("Expected error when writing to closed WAL")
		}
	})
}

func TestWALRotation(t *testing.T) {
	t.Run("rotate file when size limit reached", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		// Small max file size to trigger rotation
		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 100, // Very small to trigger rotation
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		// Write enough data to trigger rotation
		for i := 0; i < 10; i++ {
			entry := WALEntry{
				Timestamp: time.Now(),
				SeriesID:  "series" + string(rune(i)),
				Points: []DataPoint{
					{Timestamp: time.Now(), Value: float64(i)},
				},
			}

			if err := wal.Write(entry); err != nil {
				t.Fatalf("Failed to write entry %d: %v", i, err)
			}
		}

		// Check if rotation occurred
		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read temp dir: %v", err)
		}

		// Should have at least 2 files: current .wal and one rotated
		if len(files) < 2 {
			t.Errorf("Expected at least 2 files after rotation, got %d", len(files))
		}

		// Check for rotated file
		foundRotated := false
		for _, file := range files {
			if file.Name() != "test.wal" && filepath.Ext(file.Name()) != ".tmp" {
				foundRotated = true
				break
			}
		}

		if !foundRotated {
			t.Error("Expected to find rotated WAL file")
		}
	})
}

func TestWALFlush(t *testing.T) {
	t.Run("flush WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		if err := wal.Flush(); err != nil {
			t.Fatalf("Failed to flush WAL: %v", err)
		}
	})

	t.Run("flush closed WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}

		wal.Close()

		if err := wal.Flush(); err == nil {
			t.Error("Expected error when flushing closed WAL")
		}
	})
}

func TestWALClose(t *testing.T) {
	t.Run("close WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}

		if err := wal.Close(); err != nil {
			t.Fatalf("Failed to close WAL: %v", err)
		}

		if !wal.closed {
			t.Error("WAL should be marked as closed")
		}
	})

	t.Run("close already closed WAL", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}

		// Close first time
		if err := wal.Close(); err != nil {
			t.Fatalf("Failed to close WAL: %v", err)
		}

		// Close second time should not error
		if err := wal.Close(); err != nil {
			t.Fatalf("Second close should not error: %v", err)
		}
	})
}

func TestWALGetters(t *testing.T) {
	t.Run("get path", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		if wal.GetPath() != walPath {
			t.Errorf("Expected path %s, got %s", walPath, wal.GetPath())
		}
	})

	t.Run("get size", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		size := wal.GetSize()
		if size < 0 {
			t.Errorf("Size should be non-negative, got %d", size)
		}
	})

	t.Run("check if closed", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}

		if wal.IsClosed() {
			t.Error("WAL should not be closed initially")
		}

		wal.Close()

		if !wal.IsClosed() {
			t.Error("WAL should be closed after Close()")
		}
	})
}

func TestWALSerialization(t *testing.T) {
	t.Run("serialize entry", func(t *testing.T) {
		tempDir := t.TempDir()
		walPath := filepath.Join(tempDir, "test.wal")

		config := WALConfig{
			Path:        walPath,
			MaxFileSize: 1024,
		}

		wal, err := NewWAL(config)
		if err != nil {
			t.Fatalf("Failed to create WAL: %v", err)
		}
		defer wal.Close()

		entry := WALEntry{
			Timestamp: time.Now(),
			SeriesID:  "test_series",
			Points: []DataPoint{
				{Timestamp: time.Now(), Value: 42.5},
			},
		}

		// This tests the internal serializeEntry method indirectly
		if err := wal.Write(entry); err != nil {
			t.Fatalf("Failed to write entry: %v", err)
		}
	})
}
