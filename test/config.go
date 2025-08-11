package test

import (
	"os"
	"path/filepath"
	"time"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	Port           string
	DataFile       string
	TempDir        string
	TestTimeout    time.Duration
	BenchmarkCount int
}

// DefaultTestConfig returns the default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		Port:           "8080",
		DataFile:       "test_data.tsv",
		TempDir:        "test",
		TestTimeout:    30 * time.Second,
		BenchmarkCount: 1000,
	}
}

// SetupTestEnvironment sets up the test environment
func SetupTestEnvironment(config *TestConfig) error {
	if config == nil {
		config = DefaultTestConfig()
	}

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(config.TempDir, 0755); err != nil {
		return err
	}

	// Set environment variables
	os.Setenv("PORT", config.Port)
	os.Setenv("DATA_FILE", filepath.Join(config.TempDir, config.DataFile))
	os.Setenv("LOG_LEVEL", "error") // Reduce log noise during tests

	return nil
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment(config *TestConfig) error {
	if config == nil {
		config = DefaultTestConfig()
	}

	// Remove test files
	testFile := filepath.Join(config.TempDir, config.DataFile)
	if err := os.Remove(testFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Remove benchmark files
	benchmarkFile := filepath.Join(config.TempDir, "benchmark_data.tsv")
	if err := os.Remove(benchmarkFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// GetTestDataPath returns the full path for test data files
func GetTestDataPath(filename string) string {
	config := DefaultTestConfig()
	return filepath.Join(config.TempDir, filename)
}

// IsCIEnvironment checks if running in a CI environment
func IsCIEnvironment() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" ||
		os.Getenv("TRAVIS") != ""
}

// GetTestTimeout returns the appropriate timeout for tests
func GetTestTimeout() time.Duration {
	config := DefaultTestConfig()

	// Use shorter timeouts in CI environments
	if IsCIEnvironment() {
		return 10 * time.Second
	}

	return config.TestTimeout
}
