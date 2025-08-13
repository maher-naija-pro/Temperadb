package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Shard represents a storage shard with LSM tree architecture
type Shard struct {
	mu          sync.RWMutex
	id          string
	dataDir     string
	walDir      string
	segmentsDir string

	// Core components
	memStore      *MemStore
	wal           *WAL
	walReplay     *WALReplay
	segmentWriter *SegmentWriter
	segmentReader *SegmentReader
	compactionMgr *CompactionManager

	// Configuration
	config ShardConfig

	// State
	closed     bool
	recovering bool
}

// ShardConfig holds configuration for a shard
type ShardConfig struct {
	ID                  string
	DataDir             string
	MaxMemTableSize     int64
	MaxWALSize          int64
	MaxLevels           int
	MaxSegmentsPerLevel int
	MaxSegmentSize      int64
	CompactionInterval  time.Duration
}

// NewShard creates a new storage shard
func NewShard(config ShardConfig) (*Shard, error) {
	// Ensure directories exist
	walDir := filepath.Join(config.DataDir, "wal")
	segmentsDir := filepath.Join(config.DataDir, "segments")

	if err := os.MkdirAll(walDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %w", err)
	}

	if err := os.MkdirAll(segmentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create segments directory: %w", err)
	}

	// Create WAL
	wal, err := NewWAL(WALConfig{
		Path:        filepath.Join(walDir, "shard.wal"),
		MaxFileSize: config.MaxWALSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL: %w", err)
	}

	// Create WAL replay handler
	walReplay := NewWALReplay(walDir)

	// Create segment writer
	segmentWriter, err := NewSegmentWriter(SegmentWriterConfig{
		SegmentsDir: segmentsDir,
		Compression: false,
		BufferSize:  64 * 1024,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create segment writer: %w", err)
	}

	// Create segment reader
	segmentReader := NewSegmentReader(segmentsDir)

	// Create compaction manager
	compactionMgr := NewCompactionManager(CompactionConfig{
		SegmentsDir:         segmentsDir,
		MaxLevels:           config.MaxLevels,
		MaxSegmentsPerLevel: config.MaxSegmentsPerLevel,
		MaxSegmentSize:      config.MaxSegmentSize,
		CompactionInterval:  config.CompactionInterval,
		MaxConcurrent:       1,
	}, segmentReader, segmentWriter)

	// Create memstore with flush callback
	memStore := NewMemStore(config.MaxMemTableSize, wal, func(memTable *MemTable) error {
		// Flush callback: write memtable to segment and add to compaction manager
		segment, err := segmentWriter.WriteMemTable(memTable)
		if err != nil {
			return fmt.Errorf("failed to write memtable to segment: %w", err)
		}

		// Add segment to compaction manager
		if err := compactionMgr.AddSegment(segment); err != nil {
			return fmt.Errorf("failed to add segment to compaction manager: %w", err)
		}

		return nil
	})

	shard := &Shard{
		id:            config.ID,
		dataDir:       config.DataDir,
		walDir:        walDir,
		segmentsDir:   segmentsDir,
		memStore:      memStore,
		wal:           wal,
		walReplay:     walReplay,
		segmentWriter: segmentWriter,
		segmentReader: segmentReader,
		compactionMgr: compactionMgr,
		config:        config,
		closed:        false,
		recovering:    false,
	}

	return shard, nil
}

// Open opens the shard and performs recovery if needed
func (s *Shard) Open() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("shard is closed")
	}

	// Start compaction manager
	if err := s.compactionMgr.Start(); err != nil {
		return fmt.Errorf("failed to start compaction manager: %w", err)
	}

	// Perform WAL recovery
	if err := s.performRecovery(); err != nil {
		return fmt.Errorf("failed to perform recovery: %w", err)
	}

	return nil
}

// Close closes the shard
func (s *Shard) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	// Stop compaction manager
	if err := s.compactionMgr.Stop(); err != nil {
		return fmt.Errorf("failed to stop compaction manager: %w", err)
	}

	// Flush memstore
	if err := s.memStore.ForceFlush(); err != nil {
		return fmt.Errorf("failed to flush memstore: %w", err)
	}

	// Close WAL
	if err := s.wal.Close(); err != nil {
		return fmt.Errorf("failed to close WAL: %w", err)
	}

	return nil
}

// Write writes data points to the shard
func (s *Shard) Write(req WriteRequest) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("shard is closed")
	}

	if s.recovering {
		return fmt.Errorf("shard is recovering")
	}

	return s.memStore.Write(req.SeriesID, req.Points)
}

// Read reads data points from the shard
func (s *Shard) Read(req ReadRequest) ([]DataPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, fmt.Errorf("shard is closed")
	}

	var allPoints []DataPoint

	// Read from memstore
	memPoints, err := s.memStore.Read(req.SeriesID, req.Start, req.End)
	if err != nil {
		return nil, fmt.Errorf("failed to read from memstore: %w", err)
	}
	allPoints = append(allPoints, memPoints...)

	// Read from segments
	segments, err := s.segmentReader.ListSegments()
	if err != nil {
		return nil, fmt.Errorf("failed to list segments: %w", err)
	}

	for _, segment := range segments {
		// Check if segment contains the series
		containsSeries := false
		for _, seriesID := range segment.SeriesIDs {
			if seriesID == req.SeriesID {
				containsSeries = true
				break
			}
		}

		if !containsSeries {
			continue
		}

		// Check if segment overlaps with time range
		if segment.MaxTime.Before(req.Start) || segment.MinTime.After(req.End) {
			continue
		}

		// Read from segment
		results, err := s.segmentReader.ReadSegmentRange(segment.Path, req.Start, req.End)
		if err != nil {
			continue // Skip corrupted segments
		}

		// Find the right series
		for _, result := range results {
			if result.SeriesID == req.SeriesID && result.Error == nil {
				allPoints = append(allPoints, result.Points...)
				break
			}
		}
	}

	// Sort by timestamp
	sort.Slice(allPoints, func(i, j int) bool {
		return allPoints[i].Timestamp.Before(allPoints[j].Timestamp)
	})

	// Apply limit if specified
	if req.Limit > 0 && len(allPoints) > req.Limit {
		allPoints = allPoints[:req.Limit]
	}

	return allPoints, nil
}

// performRecovery performs WAL recovery on startup
func (s *Shard) performRecovery() error {
	s.recovering = true
	defer func() { s.recovering = false }()

	// Replay WAL
	result, err := s.walReplay.Replay()
	if err != nil {
		return fmt.Errorf("failed to replay WAL: %w", err)
	}

	if result.TotalCount == 0 {
		return nil // No WAL files to recover
	}

	fmt.Printf("Recovering shard %s: %d entries, %d errors\n",
		s.id, result.TotalCount, result.ErrorCount)

	// Reconstruct memstore from recovered data
	for seriesID, points := range result.SeriesData {
		if err := s.memStore.Write(seriesID, points); err != nil {
			return fmt.Errorf("failed to write recovered data for series %s: %w", seriesID, err)
		}
	}

	// Clean up old WAL files
	if err := s.walReplay.CleanupOldWALs(24 * time.Hour); err != nil {
		fmt.Printf("Warning: failed to cleanup old WAL files: %v\n", err)
	}

	return nil
}

// GetStats returns statistics about the shard
func (s *Shard) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"id":               s.id,
		"closed":           s.closed,
		"recovering":       s.recovering,
		"memstore_size":    s.memStore.GetSize(),
		"wal_size":         s.wal.GetSize(),
		"compaction_stats": s.compactionMgr.GetLevelStats(),
	}

	// Get segment stats
	if segments, err := s.segmentReader.ListSegments(); err == nil {
		stats["segment_count"] = len(segments)
		totalSize := int64(0)
		for _, segment := range segments {
			totalSize += segment.Size
		}
		stats["total_segment_size"] = totalSize
	}

	return stats
}

// ForceFlush forces a flush of the memstore
func (s *Shard) ForceFlush() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("shard is closed")
	}

	return s.memStore.ForceFlush()
}

// ForceCompaction forces compaction of a specific level
func (s *Shard) ForceCompaction(level int) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return fmt.Errorf("shard is closed")
	}

	return s.compactionMgr.ForceCompaction(level)
}

// GetID returns the shard ID
func (s *Shard) GetID() string {
	return s.id
}

// IsClosed returns whether the shard is closed
func (s *Shard) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}
