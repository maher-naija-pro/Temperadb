package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/storage"

	"github.com/stretchr/testify/suite"
)

// BaseTestSuite provides common functionality for all test suites
type BaseTestSuite struct {
	suite.Suite
	server      *httptest.Server
	storage     *storage.Storage
	config      *TestConfig
	originalEnv map[string]string
}

// SetupSuite runs once before all tests
func (suite *BaseTestSuite) SetupSuite() {
	// Initialize configuration
	suite.config = DefaultTestConfig()

	// Setup test environment
	err := SetupTestEnvironment(suite.config)
	suite.Require().NoError(err, "Failed to setup test environment")

	// Backup original environment variables
	suite.originalEnv = make(map[string]string)
	for _, key := range []string{"PORT", "DATA_FILE", "LOG_LEVEL"} {
		if val := os.Getenv(key); val != "" {
			suite.originalEnv[key] = val
		}
	}

	// Initialize logger for tests
	logger.Init()

	// Initialize test storage
	suite.storage = storage.NewStorage(GetTestDataPath(suite.config.DataFile))
}

// TearDownSuite runs once after all tests
func (suite *BaseTestSuite) TearDownSuite() {
	// Cleanup
	if suite.server != nil {
		suite.server.Close()
	}
	if suite.storage != nil {
		suite.storage.Close()
	}

	// Cleanup test environment
	err := CleanupTestEnvironment(suite.config)
	suite.Require().NoError(err, "Failed to cleanup test environment")

	// Restore original environment
	for key, val := range suite.originalEnv {
		os.Setenv(key, val)
	}
	for key := range suite.originalEnv {
		if _, exists := suite.originalEnv[key]; !exists {
			os.Unsetenv(key)
		}
	}
}

// SetupTest runs before each test
func (suite *BaseTestSuite) SetupTest() {
	// Clear storage before each test
	err := suite.storage.Clear()
	suite.Require().NoError(err, "Failed to clear storage for test")
}

// GetStorage returns the test storage instance
func (suite *BaseTestSuite) GetStorage() *storage.Storage {
	return suite.storage
}

// GetConfig returns the test configuration
func (suite *BaseTestSuite) GetConfig() *TestConfig {
	return suite.config
}

// SetServer sets the test server
func (suite *BaseTestSuite) SetServer(server *httptest.Server) {
	suite.server = server
}

// GetServer returns the test server
func (suite *BaseTestSuite) GetServer() *httptest.Server {
	return suite.server
}

// CreateTestServer creates a new test server with the given handler
func (suite *BaseTestSuite) CreateTestServer(handler http.Handler) *httptest.Server {
	server := httptest.NewServer(handler)
	suite.SetServer(server)
	return server
}

// AssertHTTPResponse is a helper method to assert HTTP responses
func (suite *BaseTestSuite) AssertHTTPResponse(t *testing.T, resp *http.Response, expectedStatus int, expectedBody string) {
	helper := NewTestHelper()
	helper.AssertResponseStatus(t, resp, expectedStatus, "Status code mismatch")
	helper.AssertResponseBody(t, resp, expectedBody, "Response body mismatch")
}

// GenerateTestData generates test data using the test helper
func (suite *BaseTestSuite) GenerateTestData(measurement string, tags map[string]string, fields map[string]float64, timestamp int64) string {
	helper := NewTestHelper()
	return helper.GenerateLineProtocolData(measurement, tags, fields, timestamp)
}

// GenerateBulkTestData generates bulk test data
func (suite *BaseTestSuite) GenerateBulkTestData(count int, baseTime int64) string {
	helper := NewTestHelper()
	return helper.GenerateBulkLineProtocolData(count, baseTime)
}
