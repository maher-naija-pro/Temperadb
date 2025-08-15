package envvars

// Environment variable keys for server configuration
const (
	// Server Configuration
	Port            = "PORT"
	ReadTimeout     = "READ_TIMEOUT"
	WriteTimeout    = "WRITE_TIMEOUT"
	IdleTimeout     = "IDLE_TIMEOUT"
	ShutdownTimeout = "SHUTDOWN_TIMEOUT"
)

// Environment variable keys for storage configuration
const (
	// Storage Configuration
	DataFile    = "DATA_FILE"
	DataDir     = "DATA_DIR"
	MaxFileSize = "MAX_FILE_SIZE"
	BackupDir   = "BACKUP_DIR"
	Compression = "COMPRESSION"

	// LSM Tree Storage Configuration
	MaxMemTableSize          = "MAX_MEMTABLE_SIZE"
	MaxWALSize               = "MAX_WAL_SIZE"
	MaxLevels                = "MAX_LEVELS"
	MaxSegmentsPerLevel      = "MAX_SEGMENTS_PER_LEVEL"
	MaxSegmentSize           = "MAX_SEGMENT_SIZE"
	CompactionInterval       = "COMPACTION_INTERVAL"
	MaxConcurrentCompactions = "MAX_CONCURRENT_COMPACTIONS"

	// Sharding Configuration
	ShardCount     = "SHARD_COUNT"
	ShardStrategy  = "SHARD_STRATEGY"
	ShardKeyFields = "SHARD_KEY_FIELDS"

	// Performance Tuning
	BufferSize          = "BUFFER_SIZE"
	FlushThreshold      = "FLUSH_THRESHOLD"
	CompactionThreshold = "COMPACTION_THRESHOLD"

	// Durability Settings
	WALEnabled       = "WAL_ENABLED"
	WALFlushInterval = "WAL_FLUSH_INTERVAL"
	SyncOnWrite      = "SYNC_ON_WRITE"
	BackupInterval   = "BACKUP_INTERVAL"

	// Monitoring and Metrics
	MetricsEnabled = "METRICS_ENABLED"
	StatsInterval  = "STATS_INTERVAL"
)

// Environment variable keys for logging configuration
const (
	// Logging Configuration
	LogLevel      = "LOG_LEVEL"
	LogFormat     = "LOG_FORMAT"
	LogOutput     = "LOG_OUTPUT"
	LogMaxSize    = "LOG_MAX_SIZE"
	LogMaxBackups = "LOG_MAX_BACKUPS"
	LogMaxAge     = "LOG_MAX_AGE"
	LogCompress   = "LOG_COMPRESS"
)
