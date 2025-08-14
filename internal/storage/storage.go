// Package storage provides functionality for persisting time-series data using LSM tree architecture
package storage

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/errors"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/types"
)

// Storage represents the main storage engine using LSM tree architecture
// This struct coordinates multiple shards, compaction, and provides a unified interface
type Storage struct {
	mu            sync.RWMutex
	config        config.StorageConfig
	shards        map[string]*Shard
	compactionMgr *CompactionManager
	metrics       *StorageMetrics
	closed        bool
}

// NewStorage creates a new storage engine with LSM tree architecture
// Takes a StorageConfig and returns a configured Storage instance
func NewStorage(cfg config.StorageConfig) *Storage {
	// Initialize storage metrics
	metrics := NewStorageMetrics()

	// Create storage instance
	storage := &Storage{
		config:  cfg,
		shards:  make(map[string]*Shard),
		metrics: metrics,
		closed:  false,
	}

	// Initialize default shard
	if err := storage.createShard("default"); err != nil {
		logger.Fatalf("Error creating default shard: %v", err)
	}

	return storage
}

// WritePoint writes a time-series point to the appropriate shard
func (s *Storage) WritePoint(p types.Point) error {
	startTime := time.Now()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return errors.WrapWithType(fmt.Errorf("storage is closed"), errors.ErrorTypeStorage, "write operation on closed storage")
	}

	// Determine which shard to use (for now, use default or create based on measurement)
	shardID := s.determineShardID(p.Measurement)

	shard, exists := s.shards[shardID]
	if !exists {
		// Need to upgrade to write lock to create shard
		s.mu.RUnlock()
		s.mu.Lock()

		// Double-check that shard still doesn't exist after acquiring write lock
		shard, exists = s.shards[shardID]
		if !exists {
			// Create new shard if it doesn't exist
			if err := s.createShard(shardID); err != nil {
				s.mu.Unlock()
				return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to create shard")
			}
			shard = s.shards[shardID]
		}

		s.mu.Unlock()
		s.mu.RLock() // Re-acquire read lock for the rest of the function
	}

	// Convert types.Point to storage.DataPoint
	dataPoints := make([]DataPoint, 0, len(p.Fields))
	for fieldName, fieldValue := range p.Fields {
		// Create a unique series ID combining measurement, tags, and field
		seriesID := s.createSeriesID(p.Measurement, p.Tags, fieldName)

		// Convert field value to float64
		value, err := s.convertToFloat64(fieldValue)
		if err != nil {
			logger.Warnf("Skipping field %s with invalid value %v: %v", fieldName, fieldValue, err)
			continue
		}

		dataPoint := DataPoint{
			Timestamp: p.Timestamp,
			Value:     value,
			Labels:    p.Tags,
		}
		dataPoints = append(dataPoints, dataPoint)

		// Write to the specific series
		writeReq := WriteRequest{
			SeriesID: seriesID,
			Points:   []DataPoint{dataPoint},
		}

		if err := shard.Write(writeReq); err != nil {
			return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to write point to shard")
		}
	}

	// Update metrics
	if s.metrics != nil {
		s.metrics.RecordShardCount(len(s.shards))
		s.metrics.RecordStorageWriteOperation("storage", "write_point")
		s.metrics.RecordDataPointsWritten("storage", len(dataPoints))
		s.metrics.RecordStorageWriteLatency("storage", "write_point", time.Since(startTime))
	}

	return nil
}

// ReadPoints reads time-series points from the storage engine
func (s *Storage) ReadPoints(measurement string, tags map[string]string, field string, start, end time.Time, limit int) ([]types.Point, error) {
	startTime := time.Now()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, errors.WrapWithType(fmt.Errorf("storage is closed"), errors.ErrorTypeStorage, "read operation on closed storage")
	}

	// Create series ID for the query
	seriesID := s.createSeriesID(measurement, tags, field)

	// Try to find the series in any shard
	var allPoints []DataPoint
	for _, shard := range s.shards {
		readReq := ReadRequest{
			SeriesID: seriesID,
			Start:    start,
			End:      end,
			Limit:    limit,
		}

		points, err := shard.Read(readReq)
		if err != nil {
			logger.Warnf("Failed to read from shard %s: %v", shard.GetID(), err)
			continue
		}

		allPoints = append(allPoints, points...)
	}

	// Convert DataPoints back to types.Point
	result := make([]types.Point, 0, len(allPoints))
	for _, dp := range allPoints {
		point := types.Point{
			Measurement: measurement,
			Tags:        dp.Labels,
			Fields:      map[string]float64{field: dp.Value},
			Timestamp:   dp.Timestamp,
		}
		result = append(result, point)
	}

	// Update metrics
	if s.metrics != nil {
		s.metrics.RecordStorageReadOperation("storage", "read_points")
		s.metrics.RecordDataPointsRead("storage", len(result))
		s.metrics.RecordStorageReadLatency("storage", "read_points", time.Since(startTime))
	}

	return result, nil
}

// createShard creates a new storage shard
func (s *Storage) createShard(shardID string) error {
	// Use existing config fields and provide sensible defaults for LSM tree
	shardConfig := ShardConfig{
		ID:                  shardID,
		DataDir:             fmt.Sprintf("%s/shard_%s", s.config.DataDir, shardID),
		MaxMemTableSize:     64 * 1024 * 1024, // 64MB default
		MaxWALSize:          s.config.MaxFileSize,
		MaxLevels:           7, // Standard LSM tree levels
		MaxSegmentsPerLevel: 10,
		MaxSegmentSize:      256 * 1024 * 1024, // 256MB default
		CompactionInterval:  30 * time.Second,
	}

	shard, err := NewShard(shardConfig, s.metrics)
	if err != nil {
		return fmt.Errorf("failed to create shard %s: %w", shardID, err)
	}

	// Open the shard
	if err := shard.Open(); err != nil {
		return fmt.Errorf("failed to open shard %s: %w", shardID, err)
	}

	s.shards[shardID] = shard

	// Update metrics
	if s.metrics != nil {
		s.metrics.RecordShardCount(len(s.shards))
	}

	logger.Infof("Created and opened shard: %s", shardID)

	return nil
}

// determineShardID determines which shard should store the data
// This is a simple hash-based approach that can be enhanced with more sophisticated routing
func (s *Storage) determineShardID(measurement string) string {
	// For now, use a simple approach - could be enhanced with consistent hashing
	if len(s.shards) == 0 {
		return "default"
	}

	// Simple hash-based shard selection
	hash := 0
	for _, char := range measurement {
		hash = (hash*31 + int(char)) % len(s.shards)
	}

	// Get shard IDs and return the selected one
	shardIDs := make([]string, 0, len(s.shards))
	for id := range s.shards {
		shardIDs = append(shardIDs, id)
	}

	if len(shardIDs) > 0 {
		return shardIDs[hash%len(shardIDs)]
	}

	return "default"
}

// createSeriesID creates a unique series identifier from measurement, tags, and field
func (s *Storage) createSeriesID(measurement string, tags map[string]string, field string) string {
	// Simple concatenation approach - could be enhanced with more sophisticated encoding
	seriesID := measurement + ":" + field

	// Add tags if present
	if len(tags) > 0 {
		// Sort tags for consistent ordering
		tagKeys := make([]string, 0, len(tags))
		for k := range tags {
			tagKeys = append(tagKeys, k)
		}

		for _, key := range tagKeys {
			seriesID += ":" + key + "=" + tags[key]
		}
	}

	return seriesID
}

// convertToFloat64 converts various field values to float64
func (s *Storage) convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case string:
		// Try to parse as float
		if f, err := parseFloat(v); err == nil {
			return f, nil
		}
		// If not a number, hash the string to a float
		return float64(hashString(v)), nil
	default:
		return 0, fmt.Errorf("unsupported value type: %T", value)
	}
}

// parseFloat attempts to parse a string as a float
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// hashString creates a simple hash of a string
func hashString(s string) uint32 {
	var hash uint32
	for _, char := range s {
		hash = hash*31 + uint32(char)
	}
	return hash
}

// GetStats returns comprehensive statistics about the storage engine
func (s *Storage) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"shard_count": len(s.shards),
		"closed":      s.closed,
		"config": map[string]interface{}{
			"data_file":     s.config.DataFile,
			"max_file_size": s.config.MaxFileSize,
			"backup_dir":    s.config.BackupDir,
			"compression":   s.config.Compression,
		},
		"shards": make(map[string]interface{}),
	}

	// Collect stats from each shard
	for shardID, shard := range s.shards {
		stats["shards"].(map[string]interface{})[shardID] = shard.GetStats()
	}

	// Add metrics
	if s.metrics != nil {
		// Note: In a real implementation, you'd want to expose Prometheus metrics
		// For now, we'll just indicate that metrics are available
		stats["metrics_available"] = true
	}

	return stats
}

// ForceCompaction forces compaction on all shards
func (s *Storage) ForceCompaction() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return errors.WrapWithType(fmt.Errorf("storage is closed"), errors.ErrorTypeStorage, "compaction on closed storage")
	}

	var lastError error
	for shardID, shard := range s.shards {
		// Force compaction on level 0 (most common)
		if err := shard.ForceCompaction(0); err != nil {
			logger.Warnf("Failed to force compaction on shard %s: %v", shardID, err)
			lastError = err
		}
	}

	return lastError
}

// Close closes the storage engine and all its shards
func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	s.closed = true

	// Close all shards
	var lastError error
	for shardID, shard := range s.shards {
		if err := shard.Close(); err != nil {
			logger.Warnf("Failed to close shard %s: %v", shardID, err)
			lastError = err
		}
	}

	// Clear shards map
	s.shards = make(map[string]*Shard)

	logger.Infof("Storage engine closed")
	return lastError
}

// Clear truncates all data from the storage engine (useful for testing)
func (s *Storage) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return errors.WrapWithType(fmt.Errorf("storage is closed"), errors.ErrorTypeStorage, "clear operation on closed storage")
	}

	// Close and remove all shards
	for shardID, shard := range s.shards {
		if err := shard.Close(); err != nil {
			logger.Warnf("Failed to close shard %s during clear: %v", shardID, err)
		}
	}

	// Clear shards map
	s.shards = make(map[string]*Shard)

	// Create default shard
	if err := s.createShard("default"); err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to recreate default shard after clear")
	}

	logger.Infof("Storage engine cleared")
	return nil
}
