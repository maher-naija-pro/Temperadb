package config

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	DataFile    string
	MaxFileSize int64
	BackupDir   string
	Compression bool
}

// NewStorageConfig creates a new StorageConfig with default values
func NewStorageConfig() StorageConfig {
	return StorageConfig{
		DataFile:    getEnvWithDefault("DATA_FILE", defaults["DATA_FILE"]),
		MaxFileSize: getInt64EnvWithDefault("MAX_FILE_SIZE", defaults["MAX_FILE_SIZE"]),
		BackupDir:   getEnvWithDefault("BACKUP_DIR", defaults["BACKUP_DIR"]),
		Compression: getBoolEnvWithDefault("COMPRESSION", defaults["COMPRESSION"]),
	}
}
