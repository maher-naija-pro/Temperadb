package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestMain is the entry point for running tests
func TestMain(m *testing.M) {
	// Setup test environment
	config := DefaultTestConfig()
	err := SetupTestEnvironment(config)
	if err != nil {
		panic("Failed to setup test environment: " + err.Error())
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup test environment
	err = CleanupTestEnvironment(config)
	if err != nil {
		// Log cleanup error but don't fail tests
		println("Warning: Failed to cleanup test environment: " + err.Error())
	}

	os.Exit(exitCode)
}

// RunTestSuite runs a specific test suite
func RunTestSuite(t *testing.T, testSuite suite.TestingSuite) {
	// This is a placeholder for custom test suite execution
	// The actual suite.Run is handled by the testify framework
}

// RunAllTests runs all test suites
func RunAllTests(t *testing.T) {
	// This function can be used to run multiple test suites
	// or for custom test orchestration
}
