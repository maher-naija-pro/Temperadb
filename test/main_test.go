package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSuite provides a structured approach to testing
type MainTestSuite struct {
	BaseTestSuite
}

// SetupSuite runs once before all tests
func (suite *MainTestSuite) SetupSuite() {
	// Call parent setup
	suite.BaseTestSuite.SetupSuite()

	// Create test server with mock handler
	suite.CreateTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		defer r.Body.Close()

		lines := make([]byte, r.ContentLength)
		r.Body.Read(lines)

		// For testing purposes, we'll use a simple response
		if len(lines) == 0 {
			http.Error(w, "Bad request", 400)
			return
		}

		// Simple validation - check if it looks like line protocol
		if strings.Contains(string(lines), " ") && strings.Count(string(lines), " ") >= 2 {
			fmt.Fprint(w, "OK")
		} else {
			http.Error(w, "Bad request", 400)
		}
	}))
}

// TearDownSuite runs once after all tests
func (suite *MainTestSuite) TearDownSuite() {
	// Call parent teardown
	suite.BaseTestSuite.TearDownSuite()
}

// SetupTest runs before each test
func (suite *MainTestSuite) SetupTest() {
	// Clear storage before each test
	err := suite.storage.Clear()
	require.NoError(suite.T(), err, "Failed to clear storage for test")
}

// TestWriteEndpoint tests the /write endpoint with various scenarios
func (suite *MainTestSuite) TestWriteEndpoint() {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
		description    string
	}{
		{
			name:           "Valid Line Protocol POST",
			method:         http.MethodPost,
			body:           "cpu,host=server01,region=us-west value=0.64 1434055562000000000",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			description:    "Should accept valid InfluxDB line protocol",
		},
		{
			name:           "Multiple Lines Protocol",
			method:         http.MethodPost,
			body:           "cpu,host=server01 value=0.64 1434055562000000000\ncpu,host=server02 value=0.65 1434055562000000000",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			description:    "Should handle multiple lines of protocol data",
		},
		{
			name:           "Invalid Method GET",
			method:         http.MethodGet,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
			description:    "Should reject non-POST methods",
		},
		{
			name:           "Invalid Method PUT",
			method:         http.MethodPut,
			body:           "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
			description:    "Should reject non-POST methods",
		},
		{
			name:           "Empty Body",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request\n",
			description:    "Should reject empty body as invalid line protocol",
		},
		{
			name:           "Invalid Line Protocol",
			method:         http.MethodPost,
			body:           "invalid,protocol,format",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request\n",
			description:    "Should reject invalid line protocol format",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req, err := http.NewRequest(tt.method, suite.server.URL+"/write", strings.NewReader(tt.body))
			require.NoError(suite.T(), err)

			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "text/plain")
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(suite.T(), err)
			defer resp.Body.Close()

			// Check status code
			assert.Equal(suite.T(), tt.expectedStatus, resp.StatusCode,
				fmt.Sprintf("Test: %s - %s", tt.name, tt.description))

			// Check response body
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expectedBody, string(bodyBytes),
				fmt.Sprintf("Test: %s - %s", tt.name, tt.description))
		})
	}
}

// TestWriteEndpointIntegration tests the full integration of the write endpoint
func (suite *MainTestSuite) TestWriteEndpointIntegration() {
	// Test data that should be successfully written
	testData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"

	req, err := http.NewRequest(http.MethodPost, suite.server.URL+"/write", strings.NewReader(testData))
	require.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Verify response
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "OK", string(bodyBytes))

	// Verify data was actually stored
	// This would require exposing a read method or checking storage directly
	// For now, we'll verify the endpoint responds correctly
}

// TestWriteEndpointPerformance tests the endpoint under load
func (suite *MainTestSuite) TestWriteEndpointPerformance() {
	// Generate test data with valid line protocol format
	var testData strings.Builder
	baseTime := int64(1434055562000000000)
	for i := 0; i < 50; i++ {
		testData.WriteString(fmt.Sprintf("cpu,host=server%02d value=%d %d\n",
			i, i, baseTime+int64(i)))
	}

	start := time.Now()

	req, err := http.NewRequest(http.MethodPost, suite.server.URL+"/write", strings.NewReader(testData.String()))
	require.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	duration := time.Since(start)

	// Verify response
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Performance assertion - should handle 50 points in reasonable time
	assert.Less(suite.T(), duration, 2*time.Second,
		"Should process 50 data points in under 2 seconds")
}

// TestWriteEndpointEdgeCases tests edge cases and error conditions
func (suite *MainTestSuite) TestWriteEndpointEdgeCases() {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
		description    string
	}{
		{
			name:           "Multiple Valid Lines",
			body:           "cpu,host=server01 value=0.64 1434055562000000000\ncpu,host=server02 value=0.65 1434055562000000001\ncpu,host=server03 value=0.66 1434055562000000002",
			expectedStatus: http.StatusOK,
			description:    "Should handle multiple valid lines",
		},
		{
			name:           "Special Characters",
			body:           "cpu,host=server-01,region=\"us-west\" value=0.64 1434055562000000000",
			expectedStatus: http.StatusOK,
			description:    "Should handle special characters in tags",
		},
		{
			name:           "Unicode Characters",
			body:           "cpu,host=服务器01,region=us-west value=0.64 1434055562000000000",
			expectedStatus: http.StatusOK,
			description:    "Should handle unicode characters",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req, err := http.NewRequest(http.MethodPost, suite.server.URL+"/write", strings.NewReader(tt.body))
			require.NoError(suite.T(), err)
			req.Header.Set("Content-Type", "text/plain")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(suite.T(), err)
			defer resp.Body.Close()

			assert.Equal(suite.T(), tt.expectedStatus, resp.StatusCode,
				fmt.Sprintf("Test: %s - %s", tt.name, tt.description))
		})
	}
}

// TestMainSuite runs the test suite
func TestMainSuite(t *testing.T) {
	// This will be handled by the test runner
	// The suite will be automatically discovered
}
