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

func TestNewSegmentReader(t *testing.T) {
	t.Run("create new segment reader", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := NewSegmentReader(tempDir)

		if reader == nil {
			t.Fatal("SegmentReader should not be nil")
		}
		if reader.segmentsDir != tempDir {
			t.Errorf("Expected segmentsDir %s, got %s", tempDir, reader.segmentsDir)
		}
	})
}

func TestReadSegment(t *testing.T) {
	t.Run("read segment with data", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a segment writer to generate test data
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create and write a memtable
		now := time.Now()
		memTable := &MemTable{
			ID: 1,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
					{Timestamp: now.Add(time.Second), Value: 2.0},
				},
				"series2": {
					{Timestamp: now.Add(2 * time.Second), Value: 3.0},
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

		// Create segment reader and read the segment
		reader := NewSegmentReader(tempDir)
		readSegment, results, err := reader.ReadSegment(segment.Path)
		if err != nil {
			t.Fatalf("Failed to read segment: %v", err)
		}

		if readSegment == nil {
			t.Fatal("Read segment should not be nil")
		}
		if readSegment.ID != segment.ID {
			t.Errorf("Expected segment ID %d, got %d", segment.ID, readSegment.ID)
		}
		if readSegment.Path != segment.Path {
			t.Errorf("Expected segment path %s, got %s", segment.Path, readSegment.Path)
		}
		if readSegment.Size <= 0 {
			t.Error("Segment size should be positive")
		}
		if readSegment.MinTime.IsZero() {
			t.Error("Segment MinTime should be set")
		}
		if readSegment.MaxTime.IsZero() {
			t.Error("Segment MaxTime should be set")
		}
		if len(readSegment.SeriesIDs) != 2 {
			t.Errorf("Expected 2 series IDs, got %d", len(readSegment.SeriesIDs))
		}

		// Check results
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		// Check first series
		if results[0].SeriesID != "series1" {
			t.Errorf("Expected first series ID 'series1', got %s", results[0].SeriesID)
		}
		if len(results[0].Points) != 2 {
			t.Errorf("Expected 2 points in first series, got %d", len(results[0].Points))
		}
		if results[0].Error != nil {
			t.Errorf("Expected no error for first series, got %v", results[0].Error)
		}

		// Check second series
		if results[1].SeriesID != "series2" {
			t.Errorf("Expected second series ID 'series2', got %s", results[1].SeriesID)
		}
		if len(results[1].Points) != 1 {
			t.Errorf("Expected 1 point in second series, got %d", len(results[1].Points))
		}
		if results[1].Error != nil {
			t.Errorf("Expected no error for second series, got %v", results[1].Error)
		}
	})

	t.Run("read non-existent segment", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := NewSegmentReader(tempDir)

		_, _, err := reader.ReadSegment("/nonexistent/path.seg")
		if err == nil {
			t.Error("Expected error when reading non-existent segment")
		}
	})
}

func TestReadSegmentRange(t *testing.T) {
	t.Run("read segment range with time filtering", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a segment writer
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create memtable with points at different times
		now := time.Now()
		memTable := &MemTable{
			ID: 2,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
					{Timestamp: now.Add(time.Hour), Value: 2.0},
					{Timestamp: now.Add(2 * time.Hour), Value: 3.0},
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

		// Read with time range
		reader := NewSegmentReader(tempDir)
		start := now.Add(30 * time.Minute)
		end := now.Add(90 * time.Minute)

		results, err := reader.ReadSegmentRange(segment.Path, start, end)
		if err != nil {
			t.Fatalf("Failed to read segment range: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		// Should only get the point at 1 hour (within range)
		if len(results[0].Points) != 1 {
			t.Errorf("Expected 1 point in time range, got %d", len(results[0].Points))
		}
		if results[0].Points[0].Value != 2.0 {
			t.Errorf("Expected point value 2.0, got %f", results[0].Points[0].Value)
		}
	})

	t.Run("read segment range with no overlap", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a segment writer
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create memtable
		now := time.Now()
		memTable := &MemTable{
			ID: 3,
			Data: map[string][]DataPoint{
				"series1": {
					{Timestamp: now, Value: 1.0},
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

		// Read with time range that doesn't overlap
		reader := NewSegmentReader(tempDir)
		start := now.Add(24 * time.Hour) // 24 hours later
		end := start.Add(time.Hour)

		results, err := reader.ReadSegmentRange(segment.Path, start, end)
		if err != nil {
			t.Fatalf("Failed to read segment range: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 results for non-overlapping time range, got %d", len(results))
		}
	})

	t.Run("read segment range with exact boundaries", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a segment writer
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create memtable
		now := time.Now()
		memTable := &MemTable{
			ID: 4,
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

		// Read with exact start and end times
		reader := NewSegmentReader(tempDir)
		results, err := reader.ReadSegmentRange(segment.Path, now, now.Add(time.Second))
		if err != nil {
			t.Fatalf("Failed to read segment range: %v", err)
		}

		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
		if len(results[0].Points) != 2 {
			t.Errorf("Expected 2 points, got %d", len(results[0].Points))
		}
	})
}

func TestListSegments(t *testing.T) {
	t.Run("list segments from directory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a segment writer
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		// Create multiple segments
		for i := 0; i < 3; i++ {
			now := time.Now()
			memTable := &MemTable{
				ID: uint64(10 + i),
				Data: map[string][]DataPoint{
					"series1": {{Timestamp: now, Value: float64(i)}},
				},
				Size:      64,
				MaxSize:   1024,
				CreatedAt: now,
				IsFlushed: false,
			}

			_, err := writer.WriteMemTable(memTable)
			if err != nil {
				t.Fatalf("Failed to write memtable %d: %v", i, err)
			}
		}

		// List segments
		reader := NewSegmentReader(tempDir)
		segments, err := reader.ListSegments()
		if err != nil {
			t.Fatalf("Failed to list segments: %v", err)
		}

		if len(segments) != 3 {
			t.Errorf("Expected 3 segments, got %d", len(segments))
		}

		// Check if all segments have valid data
		for i, segment := range segments {
			if segment.ID == 0 {
				t.Errorf("Segment %d should have valid ID", i)
			}
			if segment.Path == "" {
				t.Errorf("Segment %d should have valid path", i)
			}
			if segment.Size <= 0 {
				t.Errorf("Segment %d should have positive size", i)
			}
		}
	})

	t.Run("list segments from empty directory", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := NewSegmentReader(tempDir)

		segments, err := reader.ListSegments()
		if err != nil {
			t.Fatalf("Failed to list segments from empty directory: %v", err)
		}

		if len(segments) != 0 {
			t.Errorf("Expected 0 segments from empty directory, got %d", len(segments))
		}
	})

	t.Run("list segments with non-segment files", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create some non-segment files
		nonSegmentFiles := []string{"data.txt", "config.json", "temp.tmp"}
		for _, filename := range nonSegmentFiles {
			filePath := filepath.Join(tempDir, filename)
			if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file %s: %v", filename, err)
			}
		}

		// Create one valid segment
		writer, err := NewSegmentWriter(SegmentWriterConfig{
			SegmentsDir: tempDir,
			Compression: false,
			BufferSize:  64 * 1024,
		})
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		now := time.Now()
		memTable := &MemTable{
			ID: 20,
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

		// List segments - should only find the .seg file
		reader := NewSegmentReader(tempDir)
		segments, err := reader.ListSegments()
		if err != nil {
			t.Fatalf("Failed to list segments: %v", err)
		}

		if len(segments) != 1 {
			t.Errorf("Expected 1 segment, got %d", len(segments))
		}
	})
}

func TestSegmentReaderGetters(t *testing.T) {
	t.Run("get segments directory", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := NewSegmentReader(tempDir)

		if reader.GetSegmentsDir() != tempDir {
			t.Errorf("Expected segments directory %s, got %s", tempDir, reader.GetSegmentsDir())
		}
	})
}

func TestSegmentReaderErrorHandling(t *testing.T) {
	t.Run("read corrupted segment", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a corrupted segment file
		corruptedPath := filepath.Join(tempDir, "corrupted.seg")
		corruptedData := []byte("This is not a valid segment file")
		if err := os.WriteFile(corruptedPath, corruptedData, 0644); err != nil {
			t.Fatalf("Failed to create corrupted file: %v", err)
		}

		reader := NewSegmentReader(tempDir)
		_, _, err := reader.ReadSegment(corruptedPath)
		if err == nil {
			t.Error("Expected error when reading corrupted segment")
		}
	})

	t.Run("read segment with invalid header", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a file with invalid header length
		invalidPath := filepath.Join(tempDir, "invalid.seg")
		invalidData := []byte{0xFF, 0xFF, 0xFF, 0xFF} // Invalid header length
		if err := os.WriteFile(invalidPath, invalidData, 0644); err != nil {
			t.Fatalf("Failed to create invalid file: %v", err)
		}

		reader := NewSegmentReader(tempDir)
		_, _, err := reader.ReadSegment(invalidPath)
		if err == nil {
			t.Error("Expected error when reading segment with invalid header")
		}
	})
}
