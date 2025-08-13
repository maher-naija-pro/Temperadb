package test

import (
	"os"
	"testing"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/storage"
)

// TestSuite provides common setup and teardown for all tests
type TestSuite struct {
	Storage *storage.Storage
	Config  *config.Config
	TempDir string
}

// NewTestSuite creates a new test suite with isolated storage and temp directory
func NewTestSuite(t *testing.T) *TestSuite {
	// Create temp directory for each test
	tempDir, err := os.MkdirTemp("", "tsdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cfg := &config.Config{
		Storage: config.StorageConfig{
			DataDir:     tempDir,
			DataFile:    tempDir + "/test_data.tsv",
			MaxFileSize: 1073741824, // 1GB
			BackupDir:   tempDir + "/backups",
			Compression: false,
		},
		Logging: config.LoggingConfig{
			Level:  "error", // Reduce log noise in tests
			Format: "json",
		},
	}

	storageInstance := storage.NewStorage(cfg.Storage)

	return &TestSuite{
		Storage: storageInstance,
		Config:  cfg,
		TempDir: tempDir,
	}
}

// Cleanup removes temporary files and closes storage
func (ts *TestSuite) Cleanup() {
	if ts.Storage != nil {
		ts.Storage.Close()
	}
	if ts.TempDir != "" {
		os.RemoveAll(ts.TempDir)
	}
}

// GetTestFilePath returns a path within the test temp directory
func (ts *TestSuite) GetTestFilePath(filename string) string {
	return ts.TempDir + "/" + filename
}
