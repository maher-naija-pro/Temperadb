package benchmark

import (
	"os"
	"testing"
	"timeseriesdb/internal/logger"
)

// TestMain sets up the test environment for all benchmark tests
func TestMain(m *testing.M) {
	// Initialize logger for tests
	logger.SetTestMode(true)
	logger.Init()

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}
