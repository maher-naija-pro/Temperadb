package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	t.Run("Load with defaults", func(t *testing.T) {
		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// Check server defaults
		assert.Equal(t, "8080", cfg.Server.Port)
		assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
		assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
		assert.Equal(t, 120*time.Second, cfg.Server.IdleTimeout)

		// Check storage defaults
		assert.Equal(t, "./data/data.tsv", cfg.Storage.DataFile)
		assert.Equal(t, int64(1073741824), cfg.Storage.MaxFileSize)
		assert.Equal(t, "./data/backups", cfg.Storage.BackupDir)
		assert.False(t, cfg.Storage.Compression)

		// Check logging defaults
		assert.Equal(t, "info", cfg.Logging.Level)
		assert.Equal(t, "text", cfg.Logging.Format)
		assert.Equal(t, "stdout", cfg.Logging.Output)

		// Check database defaults
		assert.Equal(t, 100, cfg.Database.MaxConnections)
		assert.Equal(t, 5*time.Minute, cfg.Database.ConnectionTTL)
		assert.Equal(t, 30*time.Second, cfg.Database.QueryTimeout)
	})

	t.Run("Load with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("PORT", "9090")
		os.Setenv("DATA_FILE", "custom.tsv")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("MAX_FILE_SIZE", "2048")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("DATA_FILE")
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("MAX_FILE_SIZE")
		}()

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "9090", cfg.Server.Port)
		assert.Equal(t, "custom.tsv", cfg.Storage.DataFile)
		assert.Equal(t, "debug", cfg.Logging.Level)
		assert.Equal(t, int64(2048), cfg.Storage.MaxFileSize)
	})

	t.Run("Load with invalid environment variables", func(t *testing.T) {
		// Set invalid environment variables
		os.Setenv("MAX_FILE_SIZE", "invalid")
		os.Setenv("LOG_LEVEL", "invalid_level")
		defer func() {
			os.Unsetenv("MAX_FILE_SIZE")
			os.Unsetenv("LOG_LEVEL")
		}()

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// Should use defaults for invalid values
		assert.Equal(t, int64(1073741824), cfg.Storage.MaxFileSize)
		// The Level field still contains the invalid value, but GetLogLevel() should return info
		assert.Equal(t, "invalid_level", cfg.Logging.Level)
		assert.Equal(t, "info", cfg.Logging.GetLogLevel().String())
	})
}

func TestConfigString(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		Storage: StorageConfig{
			DataFile:    "data.tsv",
			MaxFileSize: 1073741824,
			BackupDir:   "backups",
			Compression: false,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Database: DatabaseConfig{
			MaxConnections: 100,
			ConnectionTTL:  5 * time.Minute,
			QueryTimeout:   30 * time.Second,
		},
	}

	configStr := cfg.String()
	assert.Contains(t, configStr, "Port: 8080")
	assert.Contains(t, configStr, "DataFile: data.tsv")
	assert.Contains(t, configStr, "Level: info")
	assert.Contains(t, configStr, "MaxConnections: 100")
}

func TestConfigValidation(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	assert.NoError(t, err) // Currently no validation logic
}
