package config

import "timeseriesdb/internal/envvars"

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	DataFile    string
	DataDir     string
	MaxFileSize int64
	BackupDir   string
	Compression bool
}

// NewStorageConfig creates a new StorageConfig with default values
func NewStorageConfig() StorageConfig {
	parser := envvars.NewParser()

	return StorageConfig{
		DataFile:    parser.String(envvars.DataFile, envvars.DefaultDataFile),
		DataDir:     parser.String(envvars.DataDir, envvars.DefaultDataDir),
		MaxFileSize: parser.FileSize(envvars.MaxFileSize, envvars.DefaultMaxFileSize),
		BackupDir:   parser.String(envvars.BackupDir, envvars.DefaultBackupDir),
		Compression: parser.Bool(envvars.Compression, envvars.DefaultCompression),
	}
}
