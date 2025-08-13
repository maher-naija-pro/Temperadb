package config

import (
	"time"

	"timeseriesdb/internal/envvars"

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
	parser := envvars.NewParser()

	return LoggingConfig{
		Level:      parser.String(envvars.LogLevel, envvars.DefaultLogLevel),
		Format:     parser.String(envvars.LogFormat, envvars.DefaultLogFormat),
		Output:     parser.String(envvars.LogOutput, envvars.DefaultLogOutput),
		MaxSize:    parser.Int(envvars.LogMaxSize, envvars.DefaultLogMaxSize),
		MaxBackups: parser.Int(envvars.LogMaxBackups, envvars.DefaultLogMaxBackups),
		MaxAge:     parser.Int(envvars.LogMaxAge, envvars.DefaultLogMaxAge),
		Compress:   parser.Bool(envvars.LogCompress, envvars.DefaultLogCompress),
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
