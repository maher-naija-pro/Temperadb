package storage

import (
	"fmt"
	"testing"
	"time"

	"timeseriesdb/internal/logger"
)

func init() {
	logger.Init()
}

// createTestSegmentReader creates a test segment reader
func createTestSegmentReader(tempDir string) *SegmentReader {
	return NewSegmentReader(tempDir)
}

// createTestSegmentWriter creates a test segment writer
func createTestSegmentWriter(tempDir string) (*SegmentWriter, error) {
	return NewSegmentWriter(SegmentWriterConfig{
		SegmentsDir: tempDir,
		Compression: false,
		BufferSize:  64 * 1024,
	})
}

func TestNewCompactionManager(t *testing.T) {
	t.Run("create new compaction manager", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := createTestSegmentReader(tempDir)
		writer, err := createTestSegmentWriter(tempDir)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024 * 1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       2,
		}

		manager := NewCompactionManager(config, reader, writer, nil)

		if manager == nil {
			t.Fatal("CompactionManager should not be nil")
		}
		if manager.segmentsDir != tempDir {
			t.Errorf("Expected segmentsDir %s, got %s", tempDir, manager.segmentsDir)
		}
		if len(manager.levels) != 3 {
			t.Errorf("Expected 3 levels, got %d", len(manager.levels))
		}
		if manager.segmentReader != reader {
			t.Error("SegmentReader should be set correctly")
		}
		if manager.segmentWriter != writer {
			t.Error("SegmentWriter should be set correctly")
		}
		if manager.running {
			t.Error("Manager should not be running initially")
		}
	})

	t.Run("compaction level initialization", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := createTestSegmentReader(tempDir)
		writer, err := createTestSegmentWriter(tempDir)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           4,
			MaxSegmentsPerLevel: 10,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, reader, writer, nil)

		// Check level configuration
		for i, level := range manager.levels {
			if level.Level != i {
				t.Errorf("Level %d should have Level field %d, got %d", i, i, level.Level)
			}
			if level.MaxSize != int64(1024*(1<<uint(i))) {
				t.Errorf("Level %d should have MaxSize %d, got %d", i, 1024*(1<<uint(i)), level.MaxSize)
			}
			if level.MaxFiles != 10 {
				t.Errorf("Level %d should have MaxFiles 10, got %d", i, level.MaxFiles)
			}
			if len(level.Segments) != 0 {
				t.Errorf("Level %d should start with 0 segments, got %d", i, len(level.Segments))
			}
		}
	})
}

func TestCompactionManagerStartStop(t *testing.T) {
	t.Run("start and stop compaction manager", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := createTestSegmentReader(tempDir)
		writer, err := createTestSegmentWriter(tempDir)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  100 * time.Millisecond, // Fast for testing
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, reader, writer, nil)

		// Start manager
		if err := manager.Start(); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}

		if !manager.running {
			t.Error("Manager should be running after Start()")
		}

		// Try to start again
		if err := manager.Start(); err == nil {
			t.Error("Expected error when starting already running manager")
		}

		// Stop manager
		if err := manager.Stop(); err != nil {
			t.Fatalf("Failed to stop manager: %v", err)
		}

		if manager.running {
			t.Error("Manager should not be running after Stop()")
		}

		// Stop again should not error
		if err := manager.Stop(); err != nil {
			t.Errorf("Second stop should not error: %v", err)
		}
	})
}

func TestAddSegment(t *testing.T) {
	t.Run("add segment to appropriate level", func(t *testing.T) {
		tempDir := t.TempDir()
		reader := createTestSegmentReader(tempDir)
		writer, err := createTestSegmentWriter(tempDir)
		if err != nil {
			t.Fatalf("Failed to create segment writer: %v", err)
		}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  100 * time.Millisecond,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, reader, writer, nil)

		// Create segments of different sizes
		smallSegment := &Segment{
			ID:        1,
			Path:      testPath("small.seg"),
			Size:      512, // Fits in level 0
			MinTime:   time.Now(),
			MaxTime:   time.Now().Add(time.Hour),
			SeriesIDs: []string{"series1"},
			CreatedAt: time.Now(),
		}

		mediumSegment := &Segment{
			ID:        2,
			Path:      testPath("medium.seg"),
			Size:      2048, // Fits in level 1
			MinTime:   time.Now(),
			MaxTime:   time.Now().Add(time.Hour),
			SeriesIDs: []string{"series2"},
			CreatedAt: time.Now(),
		}

		largeSegment := &Segment{
			ID:        3,
			Path:      testPath("large.seg"),
			Size:      4096, // Fits in level 2
			MinTime:   time.Now(),
			MaxTime:   time.Now().Add(time.Hour),
			SeriesIDs: []string{"series3"},
			CreatedAt: time.Now(),
		}

		// Add segments
		if err := manager.AddSegment(smallSegment); err != nil {
			t.Fatalf("Failed to add small segment: %v", err)
		}
		if err := manager.AddSegment(mediumSegment); err != nil {
			t.Fatalf("Failed to add medium segment: %v", err)
		}
		if err := manager.AddSegment(largeSegment); err != nil {
			t.Fatalf("Failed to add large segment: %v", err)
		}

		// Check segment placement
		if len(manager.levels[0].Segments) != 1 {
			t.Errorf("Level 0 should have 1 segment, got %d", len(manager.levels[0].Segments))
		}
		if len(manager.levels[1].Segments) != 1 {
			t.Errorf("Level 1 should have 1 segment, got %d", len(manager.levels[1].Segments))
		}
		if len(manager.levels[2].Segments) != 1 {
			t.Errorf("Level 2 should have 1 segment, got %d", len(manager.levels[2].Segments))
		}

		// Check if segments are sorted by creation time
		if manager.levels[0].Segments[0].ID != smallSegment.ID {
			t.Error("Level 0 should contain small segment")
		}
		if manager.levels[1].Segments[0].ID != mediumSegment.ID {
			t.Error("Level 1 should contain medium segment")
		}
		if manager.levels[2].Segments[0].ID != largeSegment.ID {
			t.Error("Level 2 should contain large segment")
		}
	})

	t.Run("add segment that triggers compaction", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 2, // Small limit to trigger compaction
			MaxSegmentSize:      1024,
			CompactionInterval:  100 * time.Millisecond,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Start manager to enable compaction
		if err := manager.Start(); err != nil {
			t.Fatalf("Failed to start manager: %v", err)
		}
		defer manager.Stop()

		// Add segments until compaction is triggered
		for i := 0; i < 3; i++ {
			segment := &Segment{
				ID:        uint64(i + 1),
				Path:      fmt.Sprintf("./testdata/segment_%d.seg", i),
				Size:      512,
				MinTime:   time.Now(),
				MaxTime:   time.Now().Add(time.Hour),
				SeriesIDs: []string{fmt.Sprintf("series%d", i)},
				CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
			}

			if err := manager.AddSegment(segment); err != nil {
				t.Fatalf("Failed to add segment %d: %v", i, err)
			}
		}

		// Wait a bit for compaction to be scheduled
		time.Sleep(200 * time.Millisecond)

		// Check if compaction was scheduled (level should have fewer segments)
		if len(manager.levels[0].Segments) >= 3 {
			t.Error("Compaction should have reduced segment count")
		}
	})
}

func TestFindLevelForSegment(t *testing.T) {
	t.Run("find appropriate level for segment", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           3,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Test different segment sizes
		testCases := []struct {
			size  int64
			level int
		}{
			{512, 0},  // Fits in level 0
			{1024, 0}, // Exactly fits in level 0
			{2048, 1}, // Fits in level 1
			{4096, 2}, // Fits in level 2
			{8192, 2}, // Too large, goes to highest level
		}

		for _, tc := range testCases {
			segment := &Segment{
				ID:        1,
				Path:      "./testdata/test.seg",
				Size:      tc.size,
				MinTime:   time.Now(),
				MaxTime:   time.Now().Add(time.Hour),
				SeriesIDs: []string{"test"},
				CreatedAt: time.Now(),
			}

			level := manager.findLevelForSegment(segment)
			if level != tc.level {
				t.Errorf("Segment size %d should go to level %d, got %d", tc.size, tc.level, level)
			}
		}
	})
}

func TestShouldCompactLevel(t *testing.T) {
	t.Run("check if level needs compaction", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 3,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Level 0 should not need compaction initially
		if manager.shouldCompactLevel(0) {
			t.Error("Level 0 should not need compaction initially")
		}

		// Add segments to trigger compaction
		for i := 0; i < 4; i++ {
			segment := &Segment{
				ID:        uint64(i + 1),
				Path:      fmt.Sprintf("./testdata/segment_%d.seg", i),
				Size:      512,
				MinTime:   time.Now(),
				MaxTime:   time.Now().Add(time.Hour),
				SeriesIDs: []string{fmt.Sprintf("series%d", i)},
				CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
			}

			manager.levels[0].Segments = append(manager.levels[0].Segments, segment)
		}

		// Now level 0 should need compaction
		if !manager.shouldCompactLevel(0) {
			t.Error("Level 0 should need compaction after adding 4 segments")
		}

		// Invalid level should not need compaction
		if manager.shouldCompactLevel(5) {
			t.Error("Invalid level should not need compaction")
		}
	})
}

func TestGetLevelStats(t *testing.T) {
	t.Run("get level statistics", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Add some segments to level 0
		for i := 0; i < 3; i++ {
			segment := &Segment{
				ID:        uint64(i + 1),
				Path:      fmt.Sprintf("./testdata/segment_%d.seg", i),
				Size:      512,
				MinTime:   time.Now(),
				MaxTime:   time.Now().Add(time.Hour),
				SeriesIDs: []string{fmt.Sprintf("series%d", i)},
				CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
			}

			manager.levels[0].Segments = append(manager.levels[0].Segments, segment)
		}

		stats := manager.GetLevelStats()

		// Check level 0 stats
		level0Stats, exists := stats["level_0"]
		if !exists {
			t.Fatal("Level 0 stats should exist")
		}

		level0Map := level0Stats.(map[string]interface{})
		if level0Map["segment_count"] != 3 {
			t.Errorf("Expected 3 segments in level 0, got %v", level0Map["segment_count"])
		}
		if level0Map["max_size"] != int64(1024) {
			t.Errorf("Expected max size 1024 in level 0, got %v", level0Map["max_size"])
		}
		if level0Map["max_files"] != 5 {
			t.Errorf("Expected max files 5 in level 0, got %v", level0Map["max_files"])
		}

		// Check total size calculation
		expectedTotalSize := int64(3 * 512)
		if level0Map["total_size"] != expectedTotalSize {
			t.Errorf("Expected total size %d in level 0, got %v", expectedTotalSize, level0Map["total_size"])
		}

		// Check level 1 stats (should be empty)
		level1Stats, exists := stats["level_1"]
		if !exists {
			t.Fatal("Level 1 stats should exist")
		}

		level1Map := level1Stats.(map[string]interface{})
		if level1Map["segment_count"] != 0 {
			t.Errorf("Expected 0 segments in level 1, got %v", level1Map["segment_count"])
		}
	})
}

func TestForceCompaction(t *testing.T) {
	t.Run("force compaction of specific level", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Add segments to level 0
		for i := 0; i < 3; i++ {
			segment := &Segment{
				ID:        uint64(i + 1),
				Path:      fmt.Sprintf("./testdata/segment_%d.seg", i),
				Size:      512,
				MinTime:   time.Now(),
				MaxTime:   time.Now().Add(time.Hour),
				SeriesIDs: []string{fmt.Sprintf("series%d", i)},
				CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
			}

			manager.levels[0].Segments = append(manager.levels[0].Segments, segment)
		}

		// Force compaction
		if err := manager.ForceCompaction(0); err != nil {
			t.Fatalf("Failed to force compaction: %v", err)
		}

		// Check if compaction task was scheduled
		select {
		case task := <-manager.compactionChan:
			if task.Level != 0 {
				t.Errorf("Expected compaction task for level 0, got level %d", task.Level)
			}
			if len(task.Segments) != 3 {
				t.Errorf("Expected 3 segments in compaction task, got %d", len(task.Segments))
			}
			if task.Priority != 1000 {
				t.Errorf("Expected priority 1000, got %d", task.Priority)
			}
		default:
			t.Error("Compaction task should have been scheduled")
		}
	})

	t.Run("force compaction of invalid level", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Try to force compaction of invalid level
		if err := manager.ForceCompaction(-1); err == nil {
			t.Error("Expected error when forcing compaction of negative level")
		}

		if err := manager.ForceCompaction(5); err == nil {
			t.Error("Expected error when forcing compaction of level beyond max levels")
		}
	})

	t.Run("force compaction of empty level", func(t *testing.T) {
		tempDir := t.TempDir()
		mockReader := &MockSegmentReader{}
		mockWriter := &MockSegmentWriter{}

		config := CompactionConfig{
			SegmentsDir:         tempDir,
			MaxLevels:           2,
			MaxSegmentsPerLevel: 5,
			MaxSegmentSize:      1024,
			CompactionInterval:  30 * time.Second,
			MaxConcurrent:       1,
		}

		manager := NewCompactionManager(config, mockReader, mockWriter, nil)

		// Force compaction of empty level
		if err := manager.ForceCompaction(0); err != nil {
			t.Fatalf("Force compaction of empty level should not error: %v", err)
		}
	})
}
