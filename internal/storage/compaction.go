package storage

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// CompactionManager manages the compaction of segments
type CompactionManager struct {
	mu             sync.RWMutex
	segmentsDir    string
	levels         []*CompactionLevel
	segmentReader  SegmentReaderInterface
	segmentWriter  SegmentWriterInterface
	compactionChan chan compactionTask
	stopChan       chan struct{}
	running        bool
}

// compactionTask represents a compaction job
type compactionTask struct {
	Level     int
	Segments  []*Segment
	Priority  int
	CreatedAt time.Time
}

// CompactionConfig holds configuration for compaction
type CompactionConfig struct {
	SegmentsDir         string
	MaxLevels           int
	MaxSegmentsPerLevel int
	MaxSegmentSize      int64
	CompactionInterval  time.Duration
	MaxConcurrent       int
}

// NewCompactionManager creates a new compaction manager
func NewCompactionManager(config CompactionConfig, reader SegmentReaderInterface, writer SegmentWriterInterface) *CompactionManager {
	// Initialize compaction levels
	levels := make([]*CompactionLevel, config.MaxLevels)
	for i := 0; i < config.MaxLevels; i++ {
		levels[i] = &CompactionLevel{
			Level:    i,
			Segments: make([]*Segment, 0),
			MaxSize:  config.MaxSegmentSize * int64(1<<uint(i)), // Exponential growth
			MaxFiles: config.MaxSegmentsPerLevel,
		}
	}

	return &CompactionManager{
		segmentsDir:    config.SegmentsDir,
		levels:         levels,
		segmentReader:  reader,
		segmentWriter:  writer,
		compactionChan: make(chan compactionTask, 100),
		stopChan:       make(chan struct{}),
		running:        false,
	}
}

// Start starts the compaction manager
func (cm *CompactionManager) Start() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.running {
		return fmt.Errorf("compaction manager already running")
	}

	cm.running = true

	// Start compaction workers
	go cm.compactionWorker()

	// Start periodic compaction scheduler
	go cm.scheduleCompactions()

	return nil
}

// Stop stops the compaction manager
func (cm *CompactionManager) Stop() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if !cm.running {
		return nil
	}

	cm.running = false
	close(cm.stopChan)

	return nil
}

// AddSegment adds a segment to the appropriate level
func (cm *CompactionManager) AddSegment(segment *Segment) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Find appropriate level for the segment
	level := cm.findLevelForSegment(segment)

	// Add segment to level
	cm.levels[level].Segments = append(cm.levels[level].Segments, segment)

	// Sort segments by creation time (oldest first)
	sort.Slice(cm.levels[level].Segments, func(i, j int) bool {
		return cm.levels[level].Segments[i].CreatedAt.Before(cm.levels[level].Segments[j].CreatedAt)
	})

	// Check if compaction is needed
	if cm.shouldCompactLevel(level) {
		cm.scheduleLevelCompaction(level)
	}

	return nil
}

// findLevelForSegment determines which level a segment should be placed in
func (cm *CompactionManager) findLevelForSegment(segment *Segment) int {
	for i, level := range cm.levels {
		if segment.Size <= level.MaxSize {
			return i
		}
	}
	return len(cm.levels) - 1
}

// shouldCompactLevel checks if a level needs compaction
func (cm *CompactionManager) shouldCompactLevel(level int) bool {
	if level >= len(cm.levels) {
		return false
	}

	levelData := cm.levels[level]
	return len(levelData.Segments) > levelData.MaxFiles
}

// scheduleLevelCompaction schedules compaction for a specific level
func (cm *CompactionManager) scheduleLevelCompaction(level int) {
	task := compactionTask{
		Level:     level,
		Segments:  cm.levels[level].Segments,
		Priority:  level, // Higher levels have higher priority
		CreatedAt: time.Now(),
	}

	select {
	case cm.compactionChan <- task:
		// Task scheduled successfully
	default:
		// Channel full, log warning
		fmt.Printf("Compaction task queue full, level %d compaction delayed\n", level)
	}
}

// scheduleCompactions periodically schedules compaction tasks
func (cm *CompactionManager) scheduleCompactions() {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.checkAndScheduleCompactions()
		case <-cm.stopChan:
			return
		}
	}
}

// checkAndScheduleCompactions checks all levels and schedules compaction if needed
func (cm *CompactionManager) checkAndScheduleCompactions() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for i := range cm.levels {
		if cm.shouldCompactLevel(i) {
			cm.scheduleLevelCompaction(i)
		}
	}
}

// compactionWorker processes compaction tasks
func (cm *CompactionManager) compactionWorker() {
	for {
		select {
		case task := <-cm.compactionChan:
			if err := cm.processCompactionTask(task); err != nil {
				fmt.Printf("Compaction task failed: %v\n", err)
			}
		case <-cm.stopChan:
			return
		}
	}
}

// processCompactionTask processes a single compaction task
func (cm *CompactionManager) processCompactionTask(task compactionTask) error {
	// Read all segments in the task
	var allPoints map[string][]DataPoint
	allPoints = make(map[string][]DataPoint)

	for _, segment := range task.Segments {
		_, results, err := cm.segmentReader.ReadSegment(segment.Path)
		if err != nil {
			return fmt.Errorf("failed to read segment %s: %w", segment.Path, err)
		}

		// Merge points from all segments
		for _, result := range results {
			if result.Error != nil {
				continue
			}

			if allPoints[result.SeriesID] == nil {
				allPoints[result.SeriesID] = make([]DataPoint, 0)
			}
			allPoints[result.SeriesID] = append(allPoints[result.SeriesID], result.Points...)
		}
	}

	// Sort points by timestamp for each series
	for seriesID, points := range allPoints {
		sort.Slice(points, func(i, j int) bool {
			return points[i].Timestamp.Before(points[j].Timestamp)
		})
		allPoints[seriesID] = points
	}

	// Create new memtable with merged data
	memTable := &MemTable{
		ID:        uint64(time.Now().UnixNano()),
		Data:      allPoints,
		Size:      0,
		MaxSize:   cm.levels[task.Level].MaxSize,
		CreatedAt: time.Now(),
		IsFlushed: false,
	}

	// Calculate size
	for _, points := range allPoints {
		memTable.Size += int64(len(points) * 64)
	}

	// Write new segment
	newSegment, err := cm.segmentWriter.WriteMemTable(memTable)
	if err != nil {
		return fmt.Errorf("failed to write compacted segment: %w", err)
	}

	// Remove old segments and add new one
	if err := cm.replaceSegments(task.Level, task.Segments, newSegment); err != nil {
		return fmt.Errorf("failed to replace segments: %w", err)
	}

	// Try to promote to next level if possible
	if task.Level < len(cm.levels)-1 {
		cm.tryPromoteSegment(task.Level, newSegment)
	}

	return nil
}

// replaceSegments replaces old segments with a new one
func (cm *CompactionManager) replaceSegments(level int, oldSegments []*Segment, newSegment *Segment) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Remove old segments from level
	levelData := cm.levels[level]
	newSegments := make([]*Segment, 0, len(levelData.Segments))

	for _, segment := range levelData.Segments {
		shouldKeep := true
		for _, oldSegment := range oldSegments {
			if segment.ID == oldSegment.ID {
				shouldKeep = false
				break
			}
		}
		if shouldKeep {
			newSegments = append(newSegments, segment)
		}
	}

	// Add new segment
	newSegments = append(newSegments, newSegment)

	// Sort by creation time
	sort.Slice(newSegments, func(i, j int) bool {
		return newSegments[i].CreatedAt.Before(newSegments[j].CreatedAt)
	})

	cm.levels[level].Segments = newSegments

	// Delete old segment files
	for _, segment := range oldSegments {
		if err := os.Remove(segment.Path); err != nil {
			fmt.Printf("Failed to remove old segment file %s: %v\n", segment.Path, err)
		}
	}

	return nil
}

// tryPromoteSegment tries to promote a segment to the next level
func (cm *CompactionManager) tryPromoteSegment(currentLevel int, segment *Segment) {
	if currentLevel >= len(cm.levels)-1 {
		return
	}

	nextLevel := currentLevel + 1
	if segment.Size <= cm.levels[nextLevel].MaxSize {
		// Remove from current level
		cm.mu.Lock()
		levelData := cm.levels[currentLevel]
		newSegments := make([]*Segment, 0, len(levelData.Segments))
		for _, s := range levelData.Segments {
			if s.ID != segment.ID {
				newSegments = append(newSegments, s)
			}
		}
		cm.levels[currentLevel].Segments = newSegments

		// Add to next level
		cm.levels[nextLevel].Segments = append(cm.levels[nextLevel].Segments, segment)
		sort.Slice(cm.levels[nextLevel].Segments, func(i, j int) bool {
			return cm.levels[nextLevel].Segments[i].CreatedAt.Before(cm.levels[nextLevel].Segments[j].CreatedAt)
		})
		cm.mu.Unlock()
	}
}

// GetLevelStats returns statistics about compaction levels
func (cm *CompactionManager) GetLevelStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]interface{})

	for i, level := range cm.levels {
		levelStats := map[string]interface{}{
			"segment_count": len(level.Segments),
			"max_size":      level.MaxSize,
			"max_files":     level.MaxFiles,
		}

		// Calculate total size
		totalSize := int64(0)
		for _, segment := range level.Segments {
			totalSize += segment.Size
		}
		levelStats["total_size"] = totalSize

		stats[fmt.Sprintf("level_%d", i)] = levelStats
	}

	return stats
}

// ForceCompaction forces compaction of a specific level
func (cm *CompactionManager) ForceCompaction(level int) error {
	if level < 0 || level >= len(cm.levels) {
		return fmt.Errorf("invalid level: %d", level)
	}

	cm.mu.RLock()
	segments := make([]*Segment, len(cm.levels[level].Segments))
	copy(segments, cm.levels[level].Segments)
	cm.mu.RUnlock()

	if len(segments) == 0 {
		return nil
	}

	task := compactionTask{
		Level:     level,
		Segments:  segments,
		Priority:  1000, // High priority
		CreatedAt: time.Now(),
	}

	select {
	case cm.compactionChan <- task:
		return nil
	default:
		return fmt.Errorf("compaction task queue full")
	}
}
