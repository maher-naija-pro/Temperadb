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

	// Logging Configuration Defaults
	DefaultLogLevel      = "info"
	DefaultLogFormat     = "text"
	DefaultLogOutput     = "stdout"
	DefaultLogMaxSize    = 100
	DefaultLogMaxBackups = 3
	DefaultLogMaxAge     = 28
	DefaultLogCompress   = true

	// Database Configuration Defaults
	DefaultMaxConnections = 100
	DefaultConnectionTTL  = 300 * time.Second // 5 minutes
	DefaultQueryTimeout   = 30 * time.Second  // 30 seconds
)
