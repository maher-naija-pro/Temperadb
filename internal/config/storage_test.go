package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageConfig(t *testing.T) {
	t.Run("NewStorageConfig with defaults", func(t *testing.T) {
		cfg := NewStorageConfig()
		assert.Equal(t, "/tmp/data.tsv", cfg.DataFile)
		assert.Equal(t, int64(1073741824), cfg.MaxFileSize)
		assert.Equal(t, "/tmp/backups", cfg.BackupDir)
		assert.False(t, cfg.Compression)
	})

	t.Run("NewStorageConfig with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("DATA_FILE", "custom.tsv")
		os.Setenv("MAX_FILE_SIZE", "2048")
		os.Setenv("BACKUP_DIR", "custom_backups")
		os.Setenv("COMPRESSION", "true")
		defer func() {
			os.Unsetenv("DATA_FILE")
			os.Unsetenv("MAX_FILE_SIZE")
			os.Unsetenv("BACKUP_DIR")
			os.Unsetenv("COMPRESSION")
		}()

		cfg := NewStorageConfig()
		assert.Equal(t, "custom.tsv", cfg.DataFile)
		assert.Equal(t, int64(2048), cfg.MaxFileSize)
		assert.Equal(t, "custom_backups", cfg.BackupDir)
		assert.True(t, cfg.Compression)
	})

	t.Run("NewStorageConfig with invalid environment variables", func(t *testing.T) {
		// Set invalid environment variables
		os.Setenv("MAX_FILE_SIZE", "invalid")
		os.Setenv("COMPRESSION", "invalid")
		defer func() {
			os.Unsetenv("MAX_FILE_SIZE")
			os.Unsetenv("COMPRESSION")
		}()

		cfg := NewStorageConfig()
		// Should use defaults for invalid values
		assert.Equal(t, int64(1073741824), cfg.MaxFileSize)
		assert.False(t, cfg.Compression)
	})

	t.Run("NewStorageConfig with edge case values", func(t *testing.T) {
		// Test with very large file size
		os.Setenv("MAX_FILE_SIZE", "9223372036854775807") // Max int64
		defer os.Unsetenv("MAX_FILE_SIZE")

		cfg := NewStorageConfig()
		assert.Equal(t, int64(9223372036854775807), cfg.MaxFileSize)
	})

	t.Run("NewStorageConfig with zero and negative values", func(t *testing.T) {
		// Test with zero file size
		os.Setenv("MAX_FILE_SIZE", "0")
		defer os.Unsetenv("MAX_FILE_SIZE")

		cfg := NewStorageConfig()
		assert.Equal(t, int64(0), cfg.MaxFileSize)
	})

	t.Run("NewStorageConfig with empty string values", func(t *testing.T) {
		// Test with empty string values (should use defaults)
		os.Setenv("DATA_FILE", "")
		os.Setenv("BACKUP_DIR", "")
		defer func() {
			os.Unsetenv("DATA_FILE")
			os.Unsetenv("BACKUP_DIR")
		}()

		cfg := NewStorageConfig()
		assert.Equal(t, "/tmp/data.tsv", cfg.DataFile)
		assert.Equal(t, "/tmp/backups", cfg.BackupDir)
	})

	t.Run("NewStorageConfig with whitespace values", func(t *testing.T) {
		// Test with whitespace values (should trim whitespace and use defaults)
		os.Setenv("DATA_FILE", "   ")
		os.Setenv("BACKUP_DIR", "  ")
		defer func() {
			os.Unsetenv("DATA_FILE")
			os.Unsetenv("BACKUP_DIR")
		}()

		cfg := NewStorageConfig()
		// Whitespace is trimmed, so empty strings fall back to defaults
		assert.Equal(t, "/tmp/data.tsv", cfg.DataFile)
		assert.Equal(t, "/tmp/backups", cfg.BackupDir)
	})
}
