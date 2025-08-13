package config

import (
	"time"

	"github.com/sirupsen/logrus"
)

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

// NewLoggingConfig creates a new LoggingConfig with default values
func NewLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:      getEnvWithDefault("LOG_LEVEL", defaults["LOG_LEVEL"]),
		Format:     getEnvWithDefault("LOG_FORMAT", defaults["LOG_FORMAT"]),
		Output:     getEnvWithDefault("LOG_OUTPUT", defaults["LOG_OUTPUT"]),
		MaxSize:    getIntEnvWithDefault("LOG_MAX_SIZE", "100"),
		MaxBackups: getIntEnvWithDefault("LOG_MAX_BACKUPS", "3"),
		MaxAge:     getIntEnvWithDefault("LOG_MAX_AGE", "28"),
		Compress:   getBoolEnvWithDefault("LOG_COMPRESS", "true"),
	}
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
