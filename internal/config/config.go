package config

import (
	"strconv"
	"time"

	"timeseriesdb/internal/envvars"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the TimeSeriesDB application
type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Logging  LoggingConfig
	Database DatabaseConfig
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	MaxConnections int
	ConnectionTTL  time.Duration
	QueryTimeout   time.Duration
}

// NewDatabaseConfig creates a new DatabaseConfig with default values
func NewDatabaseConfig() DatabaseConfig {
	parser := envvars.NewParser()

	return DatabaseConfig{
		MaxConnections: parser.Int(envvars.MaxConnections, envvars.DefaultMaxConnections),
		ConnectionTTL:  parser.Duration(envvars.ConnectionTTL, envvars.DefaultConnectionTTL),
		QueryTimeout:   parser.Duration(envvars.QueryTimeout, envvars.DefaultQueryTimeout),
	}
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server:   NewServerConfig(),
		Storage:  NewStorageConfig(),
		Logging:  NewLoggingConfig(),
		Database: NewDatabaseConfig(),
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	return "TimeSeriesDB Configuration:\n" +
		"Server:\n" +
		"  Port: " + c.Server.Port + "\n" +
		"  ReadTimeout: " + c.Server.ReadTimeout.String() + "\n" +
		"  WriteTimeout: " + c.Server.WriteTimeout.String() + "\n" +
		"  IdleTimeout: " + c.Server.IdleTimeout.String() + "\n" +
		"Storage:\n" +
		"  DataFile: " + c.Storage.DataFile + "\n" +
		"  MaxFileSize: " + strconv.FormatInt(c.Storage.MaxFileSize, 10) + "\n" +
		"  BackupDir: " + c.Storage.BackupDir + "\n" +
		"  Compression: " + strconv.FormatBool(c.Storage.Compression) + "\n" +
		"Database:\n" +
		"  MaxConnections: " + strconv.Itoa(c.Database.MaxConnections) + "\n" +
		"  ConnectionTTL: " + c.Database.ConnectionTTL.String() + "\n" +
		"  QueryTimeout: " + c.Database.QueryTimeout.String() + "\n" +
		"Logging:\n" +
		"  Level: " + c.Logging.Level + "\n" +
		"  Format: " + c.Logging.Format + "\n" +
		"  Output: " + c.Logging.Output + "\n"

}
