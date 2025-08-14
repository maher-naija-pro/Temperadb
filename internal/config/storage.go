package config

import (
	"time"
	"timeseriesdb/internal/envvars"
)

// StorageConfig holds storage-related configuration for LSM tree architecture
type StorageConfig struct {
	// Basic storage settings
	DataFile    string
	DataDir     string
	MaxFileSize int64
	BackupDir   string
	Compression bool

	// LSM Tree specific settings
	MaxMemTableSize          int64         // Maximum size of memtable before flush
	MaxWALSize               int64         // Maximum WAL file size before rotation
	MaxLevels                int           // Number of LSM tree levels
	MaxSegmentsPerLevel      int           // Maximum segments per level before compaction
	MaxSegmentSize           int64         // Maximum segment size
	CompactionInterval       time.Duration // Interval between compaction checks
	MaxConcurrentCompactions int           // Maximum concurrent compaction operations

	// Sharding configuration
	ShardCount     int      // Number of storage shards
	ShardStrategy  string   // Sharding strategy (hash, range, etc.)
	ShardKeyFields []string // Fields to use for shard key calculation

	// Performance tuning
	BufferSize          int     // Buffer size for file operations
	FlushThreshold      float64 // Memory usage threshold to trigger flush
	CompactionThreshold float64 // Level fullness threshold to trigger compaction

	// Durability settings
	WALEnabled       bool          // Whether to enable WAL
	WALFlushInterval time.Duration // WAL flush interval
	SyncOnWrite      bool          // Whether to sync on every write
	BackupInterval   time.Duration // Interval for automatic backups

	// Monitoring and metrics
	MetricsEnabled bool          // Whether to enable metrics collection
	StatsInterval  time.Duration // Interval for statistics collection
}

// NewStorageConfig creates a new StorageConfig with default values
func NewStorageConfig() StorageConfig {
	parser := envvars.NewParser()

	return StorageConfig{
		// Basic storage settings
		DataFile:    parser.String(envvars.DataFile, envvars.DefaultDataFile),
		DataDir:     parser.String(envvars.DataDir, envvars.DefaultDataDir),
		MaxFileSize: parser.FileSize(envvars.MaxFileSize, envvars.DefaultMaxFileSize),
		BackupDir:   parser.String(envvars.BackupDir, envvars.DefaultBackupDir),
		Compression: parser.Bool(envvars.Compression, envvars.DefaultCompression),

		// LSM Tree specific settings with sensible defaults
		MaxMemTableSize:          parser.Int64(envvars.MaxMemTableSize, envvars.DefaultMaxMemTableSize),
		MaxWALSize:               parser.Int64(envvars.MaxWALSize, envvars.DefaultMaxWALSize),
		MaxLevels:                parser.Int(envvars.MaxLevels, envvars.DefaultMaxLevels),
		MaxSegmentsPerLevel:      parser.Int(envvars.MaxSegmentsPerLevel, envvars.DefaultMaxSegmentsPerLevel),
		MaxSegmentSize:           parser.Int64(envvars.MaxSegmentSize, envvars.DefaultMaxSegmentSize),
		CompactionInterval:       parser.Duration(envvars.CompactionInterval, envvars.DefaultCompactionInterval),
		MaxConcurrentCompactions: parser.Int(envvars.MaxConcurrentCompactions, envvars.DefaultMaxConcurrentCompactions),

		// Sharding configuration
		ShardCount:     parser.Int(envvars.ShardCount, envvars.DefaultShardCount),
		ShardStrategy:  parser.String(envvars.ShardStrategy, envvars.DefaultShardStrategy),
		ShardKeyFields: []string{envvars.DefaultShardKeyFields[0]}, // For now, use first default field

		// Performance tuning
		BufferSize:          parser.Int(envvars.BufferSize, envvars.DefaultBufferSize),
		FlushThreshold:      parser.Float64(envvars.FlushThreshold, envvars.DefaultFlushThreshold),
		CompactionThreshold: parser.Float64(envvars.CompactionThreshold, envvars.DefaultCompactionThreshold),

		// Durability settings
		WALEnabled:       parser.Bool(envvars.WALEnabled, envvars.DefaultWALEnabled),
		WALFlushInterval: parser.Duration(envvars.WALFlushInterval, envvars.DefaultWALFlushInterval),
		SyncOnWrite:      parser.Bool(envvars.SyncOnWrite, envvars.DefaultSyncOnWrite),
		BackupInterval:   parser.Duration(envvars.BackupInterval, envvars.DefaultBackupInterval),

		// Monitoring and metrics
		MetricsEnabled: parser.Bool(envvars.MetricsEnabled, envvars.DefaultMetricsEnabled),
		StatsInterval:  parser.Duration(envvars.StatsInterval, envvars.DefaultStatsInterval),
	}
}

// GetShardConfig returns a shard configuration based on the storage config
func (sc *StorageConfig) GetShardConfig(shardID string) map[string]interface{} {
	return map[string]interface{}{
		"id":                     shardID,
		"max_memtable_size":      sc.MaxMemTableSize,
		"max_wal_size":           sc.MaxWALSize,
		"max_levels":             sc.MaxLevels,
		"max_segments_per_level": sc.MaxSegmentsPerLevel,
		"max_segment_size":       sc.MaxSegmentSize,
		"compaction_interval":    sc.CompactionInterval,
		"buffer_size":            sc.BufferSize,
		"wal_enabled":            sc.WALEnabled,
		"wal_flush_interval":     sc.WALFlushInterval,
		"sync_on_write":          sc.SyncOnWrite,
	}
}

// Validate validates the storage configuration
func (sc *StorageConfig) Validate() error {
	// Add validation logic here if needed
	// For now, return nil as all defaults are valid
	return nil
}
