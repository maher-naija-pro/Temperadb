package helpers

import (
	"os"
	"path/filepath"
	"testing"
	"timeseriesdb/internal/config"
)

// ConfigHelpers provides test configuration management
type ConfigHelpers struct{}

// CreateTestConfig creates a test configuration with temporary paths
func (h *ConfigHelpers) CreateTestConfig(t *testing.T) *config.Config {
	tempDir := h.CreateTempTestDir(t)

	return &config.Config{
		Storage: config.StorageConfig{
			DataFile:    filepath.Join(tempDir, "test_data.tsv"),
			MaxFileSize: 1024 * 1024, // 1MB for tests
			BackupDir:   filepath.Join(tempDir, "backups"),
			Compression: false,
		},
		Logging: config.LoggingConfig{
			Level:  "error", // Reduce log noise in tests
			Format: "json",
		},
	}
}

// CreateTempTestDir creates a temporary directory for test files
func (h *ConfigHelpers) CreateTempTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "tsdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Clean up after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

// CreateTestDataFile creates a test data file with sample content
func (h *ConfigHelpers) CreateTestDataFile(t *testing.T, dir string, filename string, content string) string {
	filepath := filepath.Join(dir, filename)
	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test data file: %v", err)
	}
	return filepath
}

// CleanupTestFiles removes test files and directories
func (h *ConfigHelpers) CleanupTestFiles(paths ...string) {
	for _, path := range paths {
		if path != "" {
			os.RemoveAll(path)
		}
	}
}

// SetTestEnv sets test environment variables
func (h *ConfigHelpers) SetTestEnv(key, value string) func() {
	original := os.Getenv(key)
	os.Setenv(key, value)

	// Return cleanup function
	return func() {
		if original != "" {
			os.Setenv(key, original)
		} else {
			os.Unsetenv(key)
		}
	}
}

// Global instance for easy access
var Config = &ConfigHelpers{}
