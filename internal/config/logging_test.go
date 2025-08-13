package config

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggingConfig(t *testing.T) {
	t.Run("NewLoggingConfig with defaults", func(t *testing.T) {
		cfg := NewLoggingConfig()
		assert.Equal(t, "info", cfg.Level)
		assert.Equal(t, "text", cfg.Format)
		assert.Equal(t, "stdout", cfg.Output)
		assert.Equal(t, 100, cfg.MaxSize)
		assert.Equal(t, 3, cfg.MaxBackups)
		assert.Equal(t, 28, cfg.MaxAge)
		assert.True(t, cfg.Compress)
	})

	t.Run("NewLoggingConfig with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("LOG_FORMAT", "json")
		os.Setenv("LOG_OUTPUT", "file.log")
		os.Setenv("LOG_MAX_SIZE", "200")
		os.Setenv("LOG_MAX_BACKUPS", "5")
		os.Setenv("LOG_MAX_AGE", "30")
		os.Setenv("LOG_COMPRESS", "false")
		defer func() {
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_FORMAT")
			os.Unsetenv("LOG_OUTPUT")
			os.Unsetenv("LOG_MAX_SIZE")
			os.Unsetenv("LOG_MAX_BACKUPS")
			os.Unsetenv("LOG_MAX_AGE")
			os.Unsetenv("LOG_COMPRESS")
		}()

		cfg := NewLoggingConfig()
		assert.Equal(t, "debug", cfg.Level)
		assert.Equal(t, "json", cfg.Format)
		assert.Equal(t, "file.log", cfg.Output)
		assert.Equal(t, 200, cfg.MaxSize)
		assert.Equal(t, 5, cfg.MaxBackups)
		assert.Equal(t, 30, cfg.MaxAge)
		assert.False(t, cfg.Compress)
	})

	t.Run("NewLoggingConfig with invalid environment variables", func(t *testing.T) {
		// Set invalid environment variables
		os.Setenv("LOG_LEVEL", "invalid_level")
		os.Setenv("LOG_MAX_SIZE", "invalid")
		os.Setenv("LOG_MAX_BACKUPS", "invalid")
		os.Setenv("LOG_MAX_AGE", "invalid")
		os.Setenv("LOG_COMPRESS", "invalid")
		defer func() {
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_MAX_SIZE")
			os.Unsetenv("LOG_MAX_BACKUPS")
			os.Unsetenv("LOG_MAX_AGE")
			os.Unsetenv("LOG_COMPRESS")
		}()

		cfg := NewLoggingConfig()
		// Should use defaults for invalid values
		assert.Equal(t, "invalid_level", cfg.Level) // Level field contains invalid value
		assert.Equal(t, 100, cfg.MaxSize)
		assert.Equal(t, 3, cfg.MaxBackups)
		assert.Equal(t, 28, cfg.MaxAge)
		assert.True(t, cfg.Compress)
	})

	t.Run("GetLogLevel with valid level", func(t *testing.T) {
		cfg := LoggingConfig{Level: "debug"}
		level := cfg.GetLogLevel()
		assert.Equal(t, "debug", level.String())
	})

	t.Run("GetLogLevel with invalid level", func(t *testing.T) {
		cfg := LoggingConfig{Level: "invalid"}
		level := cfg.GetLogLevel()
		assert.Equal(t, "info", level.String()) // Should default to info
	})

	t.Run("GetLogFormat with json", func(t *testing.T) {
		cfg := LoggingConfig{Format: "json"}
		formatter := cfg.GetLogFormat()
		assert.IsType(t, &logrus.JSONFormatter{}, formatter)
	})

	t.Run("GetLogFormat with text", func(t *testing.T) {
		cfg := LoggingConfig{Format: "text"}
		formatter := cfg.GetLogFormat()
		assert.IsType(t, &logrus.TextFormatter{}, formatter)
	})

	t.Run("GetLogFormat with invalid", func(t *testing.T) {
		cfg := LoggingConfig{Format: "invalid"}
		formatter := cfg.GetLogFormat()
		assert.IsType(t, &logrus.TextFormatter{}, formatter) // Should default to text
	})
}
