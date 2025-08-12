package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the TimeSeriesDB application
type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Logging  LoggingConfig
	Database DatabaseConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	DataFile     string
	MaxFileSize  int64
	BackupDir    string
	Compression  bool
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	MaxConnections int
	ConnectionTTL  time.Duration
	QueryTimeout   time.Duration
}

// Default configuration values
var defaults = map[string]string{
	"PORT":            "8080",
	"DATA_FILE":       "data.tsv",
	"LOG_LEVEL":       "info",
	"LOG_FORMAT":      "text",
	"LOG_OUTPUT":      "stdout",
	"MAX_FILE_SIZE":   "1073741824", // 1GB
	"BACKUP_DIR":      "backups",
	"COMPRESSION":     "false",
	"MAX_CONNECTIONS": "100",
	"CONNECTION_TTL":  "300", // 5 minutes
	"QUERY_TIMEOUT":   "30",  // 30 seconds
	"READ_TIMEOUT":    "30",  // 30 seconds
	"WRITE_TIMEOUT":   "30",  // 30 seconds
	"IDLE_TIMEOUT":    "120", // 2 minutes
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:         getEnvWithDefault("PORT", defaults["PORT"]),
			ReadTimeout:  getDurationEnvWithDefault("READ_TIMEOUT", defaults["READ_TIMEOUT"]),
			WriteTimeout: getDurationEnvWithDefault("WRITE_TIMEOUT", defaults["WRITE_TIMEOUT"]),
			IdleTimeout:  getDurationEnvWithDefault("IDLE_TIMEOUT", defaults["IDLE_TIMEOUT"]),
		},
		Storage: StorageConfig{
			DataFile:    getEnvWithDefault("DATA_FILE", defaults["DATA_FILE"]),
			MaxFileSize: getInt64EnvWithDefault("MAX_FILE_SIZE", defaults["MAX_FILE_SIZE"]),
			BackupDir:   getEnvWithDefault("BACKUP_DIR", defaults["BACKUP_DIR"]),
			Compression: getBoolEnvWithDefault("COMPRESSION", defaults["COMPRESSION"]),
		},
		Logging: LoggingConfig{
			Level:      getEnvWithDefault("LOG_LEVEL", defaults["LOG_LEVEL"]),
			Format:     getEnvWithDefault("LOG_FORMAT", defaults["LOG_FORMAT"]),
			Output:     getEnvWithDefault("LOG_OUTPUT", defaults["LOG_OUTPUT"]),
			MaxSize:    getIntEnvWithDefault("LOG_MAX_SIZE", "100"),
			MaxBackups: getIntEnvWithDefault("LOG_MAX_BACKUPS", "3"),
			MaxAge:     getIntEnvWithDefault("LOG_MAX_AGE", "28"),
			Compress:   getBoolEnvWithDefault("LOG_COMPRESS", "true"),
		},
		Database: DatabaseConfig{
			MaxConnections: getIntEnvWithDefault("MAX_CONNECTIONS", defaults["MAX_CONNECTIONS"]),
			ConnectionTTL:  getDurationEnvWithDefault("CONNECTION_TTL", defaults["CONNECTION_TTL"]),
			QueryTimeout:   getDurationEnvWithDefault("QUERY_TIMEOUT", defaults["QUERY_TIMEOUT"]),
		},
	}

	return config, nil
}

// GetLogLevel returns the parsed log level
func (c *LoggingConfig) GetLogLevel() logrus.Level {
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		return logrus.InfoLevel
	}
	return level
}

// GetLogFormat returns the log formatter
func (c *LoggingConfig) GetLogFormat() logrus.Formatter {
	switch c.Format {
	case "json":
		return &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		}
	default:
		return &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}
	}
}

// Helper functions for environment variable handling
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvWithDefault(key, defaultValue string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.Atoi(defaultValue); err == nil {
		return intValue
	}
	return 0
}

func getInt64EnvWithDefault(key, defaultValue string) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
		return intValue
	}
	return 0
}

func getBoolEnvWithDefault(key, defaultValue string) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	if boolValue, err := strconv.ParseBool(defaultValue); err == nil {
		return boolValue
	}
	return false
}

func getDurationEnvWithDefault(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value + "s"); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue + "s"); err == nil {
		return duration
	}
	return 30 * time.Second
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
		"Logging:\n" +
		"  Level: " + c.Logging.Level + "\n" +
		"  Format: " + c.Logging.Format + "\n" +
		"  Output: " + c.Logging.Output + "\n" +
		"Database:\n" +
		"  MaxConnections: " + strconv.Itoa(c.Database.MaxConnections) + "\n" +
		"  ConnectionTTL: " + c.Database.ConnectionTTL.String() + "\n" +
		"  QueryTimeout: " + c.Database.QueryTimeout.String()
}
