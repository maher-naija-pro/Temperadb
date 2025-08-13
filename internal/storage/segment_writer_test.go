package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSegmentWriter(t *testing.T) {
	t.Run("create new segment writer", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		if writer == nil {
			t.Fatal("SegmentWriter should not be nil")
		}
		if writer.segmentsDir != tempDir {
			t.Errorf("Expected segmentsDir %s, got %s", tempDir, writer.segmentsDir)
		}
		if writer.nextID == 0 {
			t.Error("nextID should be initialized")
		}
	})

	t.Run("create segment writer with non-existent directory", func(t *testing.T) {
		tempDir := t.TempDir()
		segmentsDir := filepath.Join(tempDir, "nonexistent", "segments")

		config := SegmentWriterConfig{
			SegmentsDir: segmentsDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Directory should be created automatically
		if _, err := os.Stat(segmentsDir); os.IsNotExist(err) {
			t.Error("Directory should have been created")
		}

		if writer.segmentsDir != segmentsDir {
			t.Errorf("Expected segmentsDir %s, got %s", segmentsDir, writer.segmentsDir)
		}
	})
}

func TestWriteMemTable(t *testing.T) {
	t.Run("write simple memtable", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create a simple memtable
		now := time.Now()
		memTable := &MemTable{
			ID: 1,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
					{Timestamp: now.Add(time.Second), Value: 2.0},
				},
			},
			Size:      128,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		if segment == nil {
			t.Fatal("Segment should not be nil")
		}
		if segment.ID == 0 {
			t.Error("Segment ID should be set")
		}
		if segment.Path == "" {
			t.Error("Segment path should be set")
		}
		if segment.Size <= 0 {
			t.Error("Segment size should be positive")
		}
		if segment.MinTime.IsZero() {
			t.Error("Segment MinTime should be set")
		}
		if segment.MaxTime.IsZero() {
			t.Error("Segment MaxTime should be set")
		}
		if len(segment.SeriesIDs) != 1 {
			t.Errorf("Expected 1 series ID, got %d", len(segment.SeriesIDs))
		}
		if segment.SeriesIDs[0] != "series1" {
			t.Errorf("Expected series ID 'series1', got %s", segment.SeriesIDs[0])
		}

		// Check if file was created
		if _, err := os.Stat(segment.Path); os.IsNotExist(err) {
			t.Error("Segment file should have been created")
		}
	})

	t.Run("write memtable with multiple series", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		now := time.Now()
		memTable := &MemTable{
			ID: 2,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
				},
				"series2": {
					{Timestamp: now.Add(time.Second), Value: 2.0},
				},
				"series3": {
					{Timestamp: now.Add(2 * time.Second), Value: 3.0},
				},
			},
			Size:      256,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		if len(segment.SeriesIDs) != 3 {
			t.Errorf("Expected 3 series IDs, got %d", len(segment.SeriesIDs))
		}

		// Check if series IDs are sorted
		expectedSeries := []string{"series1", "series2", "series3"}
		for i, expected := range expectedSeries {
			if segment.SeriesIDs[i] != expected {
				t.Errorf("Expected series ID %s at position %d, got %s", expected, i, segment.SeriesIDs[i])
			}
		}
	})

	t.Run("write empty memtable", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		now := time.Now()
		memTable := &MemTable{
			ID:        3,
			Data:      make(map[string][]DataPoint),
			Size:      0,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write empty memtable: %v", err)
		}

		if len(segment.SeriesIDs) != 0 {
			t.Errorf("Expected 0 series IDs, got %d", len(segment.SeriesIDs))
		}
	})
}

func TestSegmentHeaderCalculation(t *testing.T) {
	t.Run("header calculation with data", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		now := time.Now()
		memTable := &MemTable{
			ID: 4,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
					{Timestamp: now.Add(time.Hour), Value: 2.0},
				},
				"series2": {
					{Timestamp: now.Add(30 * time.Minute), Value: 1.5},
				},
			},
			Size:      192,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		// Check header calculations
		if segment.MinTime != now {
			t.Errorf("Expected MinTime %v, got %v", now, segment.MinTime)
		}
		if segment.MaxTime != now.Add(time.Hour) {
			t.Errorf("Expected MaxTime %v, got %v", now.Add(time.Hour), segment.MaxTime)
		}
		if len(segment.SeriesIDs) != 2 {
			t.Errorf("Expected 2 series, got %d", len(segment.SeriesIDs))
		}
	})

	t.Run("header calculation with single point", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		now := time.Now()
		memTable := &MemTable{
			ID: 5,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 42.5},
				},
			},
			Size:      64,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		// For single point, MinTime and MaxTime should be the same
		if !segment.MinTime.Equal(segment.MaxTime) {
			t.Errorf("MinTime and MaxTime should be equal for single point, got %v and %v", segment.MinTime, segment.MaxTime)
		}
		if segment.MinTime != now {
			t.Errorf("Expected MinTime %v, got %v", now, segment.MinTime)
		}
	})
}

func TestSegmentFileCreation(t *testing.T) {
	t.Run("segment file naming", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Get initial nextID
		initialID := writer.GetNextID()

		now := time.Now()
		memTable := &MemTable{
			ID: 6,
			Data: map[string][]DataPoint{
				"series1": {{Timestamp: now, Value: 1.0}},
			},
			Size:      64,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		segment, err := writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		// Check if file name follows expected pattern
		expectedFileName := fmt.Sprintf("segment_%d.seg", initialID)
		actualFileName := filepath.Base(segment.Path)
		if actualFileName != expectedFileName {
			t.Errorf("Expected file name %s, got %s", expectedFileName, actualFileName)
		}

		// Check if nextID was incremented
		newID := writer.GetNextID()
		if newID <= initialID {
			t.Errorf("Expected nextID to increase, got %d (was %d)", newID, initialID)
		}
	})

	t.Run("multiple segment files", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create multiple segments
		for i := 0; i < 3; i++ {
			now := time.Now()
			memTable := &MemTable{
				ID: uint64(7 + i),
				Data: map[string][]DataPoint{
					"series1": {{Timestamp: now, Value: float64(i)}},
				},
				Size:      64,
				MaxSize:   1024,
				CreatedAt: now,
				IsFlushed: false,
			}

			segment, err := writer.WriteMemTable(memTable)
			if err != nil {
				t.Fatalf("Failed to write memtable %d: %v", i, err)
			}

			// Check if file exists
			if _, err := os.Stat(segment.Path); os.IsNotExist(err) {
				t.Errorf("Segment file %d should exist", i)
			}
		}

		// Check total number of files
		files, err := os.ReadDir(tempDir)
		if err != nil {
			t.Fatalf("Failed to read temp dir: %v", err)
		}

		segmentFiles := 0
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".seg" {
				segmentFiles++
			}
		}

		if segmentFiles != 3 {
			t.Errorf("Expected 3 segment files, got %d", segmentFiles)
		}
	})
}

func TestSegmentWriterGetters(t *testing.T) {
	t.Run("get segments directory", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		if writer.GetSegmentsDir() != tempDir {
			t.Errorf("Expected segments directory %s, got %s", tempDir, writer.GetSegmentsDir())
		}
	})

	t.Run("get next ID", func(t *testing.T) {
		tempDir := t.TempDir()

		config := SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		}

		writer, err := NewSegmentWriter(config)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		initialID := writer.GetNextID()
		if initialID == 0 {
			t.Error("Next ID should not be 0")
		}

		// Write a segment to increment the ID
		now := time.Now()
		memTable := &MemTable{
			ID: 10,
			Data: map[string][]DataPoint{
				"series1": {{Timestamp: now, Value: 1.0}},
			},
			Size:      64,
			MaxSize:   1024,
			CreatedAt: now,
			IsFlushed: false,
		}

		_, err = writer.WriteMemTable(memTable)
		if err != nil {
			t.Fatalf("Failed to write memtable: %v", err)
		}

		newID := writer.GetNextID()
		if newID <= initialID {
			t.Errorf("Expected next ID to increase, got %d (was %d)", newID, initialID)
		}
	})
}
