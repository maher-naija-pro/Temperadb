package storage

import (
	"testing"
	"time"
	"timeseriesdb/internal/logger"
)

func init() {
	logger.Init()
}

func TestDataPoint(t *testing.T) {
	t.Run("DataPoint creation", func(t *testing.T) {
		now := time.Now()
		point := DataPoint{
			Timestamp: now,
			Value:     42.5,
			Labels: map[string]string{
				"host": "server1",
				"cpu":  "0",
			},
		}

		if point.Timestamp != now {
			t.Errorf("Expected timestamp %v, got %v", now, point.Timestamp)
		}
		if point.Value != 42.5 {
			t.Errorf("Expected value 42.5, got %f", point.Value)
		}
		if len(point.Labels) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(point.Labels))
		}
	})

	t.Run("DataPoint with empty labels", func(t *testing.T) {
		point := DataPoint{
			Timestamp: time.Now(),
			Value:     100.0,
			Labels:    make(map[string]string),
		}

		if point.Labels == nil {
			t.Error("Labels should not be nil")
		}
		if len(point.Labels) != 0 {
			t.Errorf("Expected 0 labels, got %d", len(point.Labels))
		}
	})
}

func TestWriteRequest(t *testing.T) {
	t.Run("WriteRequest creation", func(t *testing.T) {
		points := []DataPoint{
			{Timestamp: time.Now(), Value: 10.5},
		}

		req := WriteRequest{
			SeriesID: "test_series",
			Points:   points,
		}

		if req.SeriesID != "test_series" {
			t.Errorf("Expected SeriesID 'test_series', got %s", req.SeriesID)
		}
		if len(req.Points) != 1 {
			t.Errorf("Expected 1 point, got %d", len(req.Points))
		}
	})

	t.Run("WriteRequest with empty points", func(t *testing.T) {
		req := WriteRequest{
			SeriesID: "empty_series",
			Points:   []DataPoint{},
		}

		if req.SeriesID != "empty_series" {
			t.Errorf("Expected SeriesID 'empty_series', got %s", req.SeriesID)
		}
		if len(req.Points) != 0 {
			t.Errorf("Expected 0 points, got %d", len(req.Points))
		}
	})
}

func TestReadRequest(t *testing.T) {
	t.Run("ReadRequest creation", func(t *testing.T) {
		start := time.Now()
		end := start.Add(time.Hour)

		req := ReadRequest{
			SeriesID: "test_series",
			Start:    start,
			End:      end,
			Limit:    100,
		}

		if req.SeriesID != "test_series" {
			t.Errorf("Expected SeriesID 'test_series', got %s", req.SeriesID)
		}
		if req.Start != start {
			t.Errorf("Expected Start %v, got %v", start, req.Start)
		}
		if req.End != end {
			t.Errorf("Expected End %v, got %v", end, req.End)
		}
		if req.Limit != 100 {
			t.Errorf("Expected Limit 100, got %d", req.Limit)
		}
	})

	t.Run("ReadRequest with zero limit", func(t *testing.T) {
		req := ReadRequest{
			SeriesID: "test_series",
			Start:    time.Now(),
			End:      time.Now().Add(time.Hour),
			Limit:    0,
		}

		if req.Limit != 0 {
			t.Errorf("Expected Limit 0, got %d", req.Limit)
		}
	})
}

func TestMemTable(t *testing.T) {
	t.Run("MemTable creation", func(t *testing.T) {
		now := time.Now()
		memTable := MemTable{
			ID:        12345,
			Data:      make(map[string][]DataPoint),
			Size:      1024,
			MaxSize:   1048576,
			CreatedAt: now,
			IsFlushed: false,
		}

		if memTable.ID != 12345 {
			t.Errorf("Expected ID 12345, got %d", memTable.ID)
		}
		if memTable.Size != 1024 {
			t.Errorf("Expected Size 1024, got %d", memTable.Size)
		}
		if memTable.MaxSize != 1048576 {
			t.Errorf("Expected MaxSize 1048576, got %d", memTable.MaxSize)
		}
		if memTable.CreatedAt != now {
			t.Errorf("Expected CreatedAt %v, got %v", now, memTable.CreatedAt)
		}
		if memTable.IsFlushed != false {
			t.Errorf("Expected IsFlushed false, got %v", memTable.IsFlushed)
		}
	})

	t.Run("MemTable with data", func(t *testing.T) {
		data := map[string][]DataPoint{
			"series1": {{Timestamp: time.Now(), Value: 1.0}},
			"series2": {{Timestamp: time.Now(), Value: 2.0}},
		}

		memTable := MemTable{
			ID:        1,
			Data:      data,
			Size:      2048,
			MaxSize:   1048576,
			CreatedAt: time.Now(),
			IsFlushed: false,
		}

		if len(memTable.Data) != 2 {
			t.Errorf("Expected 2 series, got %d", len(memTable.Data))
		}
		if len(memTable.Data["series1"]) != 1 {
			t.Errorf("Expected 1 point in series1, got %d", len(memTable.Data["series1"]))
		}
	})
}

func TestSegment(t *testing.T) {
	t.Run("Segment creation", func(t *testing.T) {
		now := time.Now()
		segment := Segment{
			ID:        67890,
			Path:      "/tmp/segment_67890.seg",
			Size:      2048,
			MinTime:   now,
			MaxTime:   now.Add(time.Hour),
			SeriesIDs: []string{"series1", "series2"},
			CreatedAt: now,
		}

		if segment.ID != 67890 {
			t.Errorf("Expected ID 67890, got %d", segment.ID)
		}
		if segment.Path != "/tmp/segment_67890.seg" {
			t.Errorf("Expected Path '/tmp/segment_67890.seg', got %s", segment.Path)
		}
		if segment.Size != 2048 {
			t.Errorf("Expected Size 2048, got %d", segment.Size)
		}
		if len(segment.SeriesIDs) != 2 {
			t.Errorf("Expected 2 series IDs, got %d", len(segment.SeriesIDs))
		}
	})

	t.Run("Segment with empty series", func(t *testing.T) {
		segment := Segment{
			ID:        1,
			Path:      "/tmp/empty.seg",
			Size:      0,
			MinTime:   time.Time{},
			MaxTime:   time.Time{},
			SeriesIDs: []string{},
			CreatedAt: time.Now(),
		}

		if len(segment.SeriesIDs) != 0 {
			t.Errorf("Expected 0 series IDs, got %d", len(segment.SeriesIDs))
		}
		if segment.Size != 0 {
			t.Errorf("Expected Size 0, got %d", segment.Size)
		}
	})
}

func TestWALEntry(t *testing.T) {
	t.Run("WALEntry creation", func(t *testing.T) {
		now := time.Now()
		points := []DataPoint{
			{Timestamp: now, Value: 42.0},
		}

		entry := WALEntry{
			ID:        11111,
			Timestamp: now,
			SeriesID:  "test_series",
			Points:    points,
			Checksum:  12345,
		}

		if entry.ID != 11111 {
			t.Errorf("Expected ID 11111, got %d", entry.ID)
		}
		if entry.SeriesID != "test_series" {
			t.Errorf("Expected SeriesID 'test_series', got %s", entry.SeriesID)
		}
		if len(entry.Points) != 1 {
			t.Errorf("Expected 1 point, got %d", len(entry.Points))
		}
		if entry.Checksum != 12345 {
			t.Errorf("Expected Checksum 12345, got %d", entry.Checksum)
		}
	})
}

func TestCompactionLevel(t *testing.T) {
	t.Run("CompactionLevel creation", func(t *testing.T) {
		level := CompactionLevel{
			Level:    1,
			Segments: []*Segment{},
			MaxSize:  1048576,
			MaxFiles: 10,
		}

		if level.Level != 1 {
			t.Errorf("Expected Level 1, got %d", level.Level)
		}
		if level.MaxSize != 1048576 {
			t.Errorf("Expected MaxSize 1048576, got %d", level.MaxSize)
		}
		if level.MaxFiles != 10 {
			t.Errorf("Expected MaxFiles 10, got %d", level.MaxFiles)
		}
		if len(level.Segments) != 0 {
			t.Errorf("Expected 0 segments, got %d", len(level.Segments))
		}
	})
}
