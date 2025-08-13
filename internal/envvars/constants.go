package envvars

// Environment variable keys for server configuration
const (
	// Server Configuration
	Port         = "PORT"
	ReadTimeout  = "READ_TIMEOUT"
	WriteTimeout = "WRITE_TIMEOUT"
	IdleTimeout  = "IDLE_TIMEOUT"
)

// Environment variable keys for storage configuration
const (
	// Storage Configuration
	DataFile    = "DATA_FILE"
	DataDir     = "DATA_DIR"
	MaxFileSize = "MAX_FILE_SIZE"
	BackupDir   = "BACKUP_DIR"
	Compression = "COMPRESSION"
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

// Environment variable keys for database configuration
const (
	// Database Configuration
	MaxConnections = "MAX_CONNECTIONS"
	ConnectionTTL  = "CONNECTION_TTL"
	QueryTimeout   = "QUERY_TIMEOUT"
)
