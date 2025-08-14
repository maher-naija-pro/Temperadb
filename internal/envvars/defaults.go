package envvars

import "time"

// Default values for environment variables
var (
	// Server Configuration Defaults
	DefaultPort         = "8080"
	DefaultReadTimeout  = 30 * time.Second
	DefaultWriteTimeout = 30 * time.Second
	DefaultIdleTimeout  = 120 * time.Second

	// Storage Configuration Defaults
	DefaultDataFile    = "/tmp/data.tsv"
	DefaultDataDir     = "/tmp"
	DefaultMaxFileSize = int64(1073741824) // 1GB
	DefaultBackupDir   = "/tmp/backups"
	DefaultCompression = false

	// LSM Tree Configuration Defaults
	DefaultMaxMemTableSize          = int64(64 * 1024 * 1024)  // 64MB
	DefaultMaxWALSize               = int64(64 * 1024 * 1024)  // 64MB
	DefaultMaxLevels                = 7                        // Standard LSM tree levels
	DefaultMaxSegmentsPerLevel      = 10                       // 10 segments per level
	DefaultMaxSegmentSize           = int64(256 * 1024 * 1024) // 256MB
	DefaultCompactionInterval       = 30 * time.Second
	DefaultMaxConcurrentCompactions = 2

	// Sharding Configuration Defaults
	DefaultShardCount     = 1
	DefaultShardStrategy  = "hash"
	DefaultShardKeyFields = []string{"measurement"}

	// Performance Tuning Defaults
	DefaultBufferSize          = 64 * 1024 // 64KB buffer
	DefaultFlushThreshold      = 0.8       // 80% memory usage
	DefaultCompactionThreshold = 0.9       // 90% level fullness

	// Durability Settings Defaults
	DefaultWALEnabled       = true
	DefaultWALFlushInterval = 100 * time.Millisecond
	DefaultSyncOnWrite      = false          // Performance over durability by default
	DefaultBackupInterval   = 24 * time.Hour // Daily backups

	// Monitoring and Metrics Defaults
	DefaultMetricsEnabled = true
	DefaultStatsInterval  = 60 * time.Second // Stats every minute

	// Logging Configuration Defaults
	DefaultLogLevel      = "info"
	DefaultLogFormat     = "text"
	DefaultLogOutput     = "stdout"
	DefaultLogMaxSize    = 100
	DefaultLogMaxBackups = 3
	DefaultLogMaxAge     = 28
	DefaultLogCompress   = true
)
