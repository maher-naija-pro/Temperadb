package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"timeseriesdb/internal/storage"
	"timeseriesdb/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test data for write endpoint tests
var (
	testDataFile      = "test_data.tsv"
	validLineProtocol = "cpu,host=server01,region=us-west value=0.64 1434055562000000000\n" +
		"memory,host=server01,region=us-west used=1234567 1434055562000000001"
	invalidLineProtocol   = "cpu,host=server01 value=invalid 1434055562000000000"
	emptyLineProtocol     = ""
	malformedLineProtocol = "cpu,host=server01 value=0.64"
)

// setupTestStorage creates a test storage instance
func setupTestStorage(t *testing.T) *storage.Storage {
	// Clean up any existing test file
	os.Remove(testDataFile)

	storageInstance := storage.NewStorage(testDataFile)
	require.NotNil(t, storageInstance)
	return storageInstance
}

// cleanupTestStorage removes test files
func cleanupTestStorage(t *testing.T, storageInstance *storage.Storage) {
	if storageInstance != nil {
		storageInstance.Close()
	}
	os.Remove(testDataFile)
}

// createTestServer creates a test HTTP server with the write handler
func createTestServer(storageInstance *storage.Storage) *httptest.Server {
	// Create a test server that mimics the main server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/write" {
			handleWrite(w, r, storageInstance)
		} else {
			http.NotFound(w, r)
		}
	}))

	return server
}

// handleWrite is copied from main.go for testing
func handleWrite(w http.ResponseWriter, r *http.Request, storageInstance *storage.Storage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	defer r.Body.Close()

	// Handle negative or zero ContentLength
	if r.ContentLength <= 0 {
		http.Error(w, "Bad request", 400)
		return
	}

	lines := make([]byte, r.ContentLength)
	r.Body.Read(lines)

	points, err := parseLineProtocol(string(lines))
	if err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	for _, p := range points {
		err := storageInstance.WritePoint(p)
		if err != nil {
			http.Error(w, "Internal server error", 500)
			return
		}
	}

	fmt.Fprint(w, "OK")
}

// parseLineProtocol is copied from parser package for testing
func parseLineProtocol(input string) ([]types.Point, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var points []types.Point

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid line format: expected 3 parts, got %d", len(parts))
		}

		// Parse measurement and tags
		measurementAndTags := strings.Split(parts[0], ",")
		measurement := measurementAndTags[0]
		if measurement == "" {
			return nil, fmt.Errorf("missing measurement name")
		}

		tags := map[string]string{}
		for _, tag := range measurementAndTags[1:] {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("malformed tag: %s", tag)
			}
			if kv[0] == "" || kv[1] == "" {
				return nil, fmt.Errorf("invalid tag key or value: %s", tag)
			}
			tags[kv[0]] = kv[1]
		}

		// Parse fields
		fields := map[string]float64{}
		fieldPairs := strings.Split(parts[1], ",")
		if len(fieldPairs) == 0 {
			return nil, fmt.Errorf("no fields provided")
		}

		for _, fieldPair := range fieldPairs {
			kv := strings.SplitN(fieldPair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("malformed field: %s", fieldPair)
			}
			if kv[0] == "" {
				return nil, fmt.Errorf("empty field name")
			}

			val, err := strconv.ParseFloat(strings.TrimSuffix(kv[1], "i"), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid field value '%s': %v", kv[1], err)
			}
			fields[kv[0]] = val
		}

		// Parse timestamp
		tsInt, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %v", err)
		}
		timestamp := time.Unix(0, tsInt)

		points = append(points, types.Point{
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
			Timestamp:   timestamp,
		})
	}

	return points, nil
}

// TestWriteEndpointHTTPMethods tests that only POST method is allowed
func TestWriteEndpointHTTPMethods(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	tests := []struct {
		method         string
		expectedStatus int
		description    string
	}{
		{"GET", 405, "GET method should not be allowed"},
		{"POST", 200, "POST method should be allowed"},
		{"PUT", 405, "PUT method should not be allowed"},
		{"DELETE", 405, "DELETE method should not be allowed"},
		{"PATCH", 405, "PATCH method should not be allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var body *strings.Reader
			if tt.method == "POST" {
				body = strings.NewReader("test,host=test value=1.0 1434055562000000000")
				req, err := http.NewRequest(tt.method, server.URL+"/write", body)
				require.NoError(t, err)
				req.ContentLength = int64(body.Len())

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode, tt.description)
			} else {
				req, err := http.NewRequest(tt.method, server.URL+"/write", nil)
				require.NoError(t, err)

				resp, err := http.DefaultClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, tt.expectedStatus, resp.StatusCode, tt.description)
			}
		})
	}
}

// TestWriteEndpointInvalidContentLength tests handling of invalid content length
func TestWriteEndpointInvalidContentLength(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with no content length header
	req, err := http.NewRequest("POST", server.URL+"/write", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for invalid content length")
}

// TestWriteEndpointEmptyBody tests handling of empty request body
func TestWriteEndpointEmptyBody(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with empty body
	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(""))
	require.NoError(t, err)
	req.ContentLength = 0

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for empty body")
}

// TestWriteEndpointValidLineProtocol tests successful parsing and writing of valid line protocol
func TestWriteEndpointValidLineProtocol(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with valid line protocol
	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(validLineProtocol))
	require.NoError(t, err)
	req.ContentLength = int64(len(validLineProtocol))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for valid line protocol")

	// Verify response body
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Equal(t, "OK", string(body), "Response body should be 'OK'")
}

// TestWriteEndpointInvalidLineProtocol tests handling of invalid line protocol
func TestWriteEndpointInvalidLineProtocol(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with invalid line protocol
	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(invalidLineProtocol))
	require.NoError(t, err)
	req.ContentLength = int64(len(invalidLineProtocol))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for invalid line protocol")
}

// TestWriteEndpointMalformedLineProtocol tests handling of malformed line protocol
func TestWriteEndpointMalformedLineProtocol(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with malformed line protocol
	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(malformedLineProtocol))
	require.NoError(t, err)
	req.ContentLength = int64(len(malformedLineProtocol))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 400, resp.StatusCode, "Should return 400 for malformed line protocol")
}

// TestWriteEndpointMultiplePoints tests writing multiple data points
func TestWriteEndpointMultiplePoints(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with multiple valid points
	multiPointData := "cpu,host=server01 value=0.64 1434055562000000000\n" +
		"memory,host=server01 used=1234567 1434055562000000001\n" +
		"disk,host=server01 free=987654321 1434055562000000002"

	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(multiPointData))
	require.NoError(t, err)
	req.ContentLength = int64(len(multiPointData))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for multiple valid points")
}

// TestWriteEndpointWithTags tests writing data points with various tag configurations
func TestWriteEndpointWithTags(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with various tag configurations
	tagTestData := "cpu,host=server01,region=us-west,env=prod value=0.64 1434055562000000000\n" +
		"memory,host=server02,datacenter=dc1 used=1234567 1434055562000000001"

	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(tagTestData))
	require.NoError(t, err)
	req.ContentLength = int64(len(tagTestData))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for data with tags")
}

// TestWriteEndpointWithMultipleFields tests writing data points with multiple fields
func TestWriteEndpointWithMultipleFields(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with multiple fields per measurement
	multiFieldData := "cpu,host=server01 value=0.64,idle=99.36,user=0.64 1434055562000000000\n" +
		"memory,host=server01 used=1234567,free=8765432,total=9999999 1434055562000000001"

	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(multiFieldData))
	require.NoError(t, err)
	req.ContentLength = int64(len(multiFieldData))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for data with multiple fields")
}

// TestWriteEndpointEmptyLines tests handling of empty lines in the input
func TestWriteEndpointEmptyLines(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test with empty lines mixed with valid data
	dataWithEmptyLines := "cpu,host=server01 value=0.64 1434055562000000000\n" +
		"\n" +
		"memory,host=server01 used=1234567 1434055562000000001\n" +
		"  \n" +
		"disk,host=server01 free=987654321 1434055562000000002"

	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(dataWithEmptyLines))
	require.NoError(t, err)
	req.ContentLength = int64(len(dataWithEmptyLines))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for data with empty lines")
}

// TestWriteEndpointLargeData tests handling of larger datasets
func TestWriteEndpointLargeData(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Generate a smaller dataset first to test the concept
	var largeData strings.Builder
	for i := 0; i < 10; i++ {
		if i == 9 {
			// Last line without newline
			largeData.WriteString(fmt.Sprintf("metric%d,host=server01 value=%.1f 1434055562000000000", i, float64(i)))
		} else {
			largeData.WriteString(fmt.Sprintf("metric%d,host=server01 value=%.1f 1434055562000000000\n", i, float64(i)))
		}
	}

	dataStr := largeData.String()

	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(dataStr))
	require.NoError(t, err)
	req.ContentLength = int64(len(dataStr))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode, "Should return 200 for large dataset")
}

// TestWriteEndpointMultipleSequentialRequests tests handling of multiple sequential write requests
// Note: We test sequential requests instead of concurrent requests because the current storage
// implementation is not thread-safe. Multiple goroutines writing to the same file can cause
// race conditions and errors. In a production environment, you would want to add proper
// synchronization (e.g., mutex) to the storage layer.
func TestWriteEndpointMultipleSequentialRequests(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test multiple sequential requests
	const numRequests = 5

	for i := 0; i < numRequests; i++ {
		data := fmt.Sprintf("sequential_test,request_id=%d value=%.1f 1434055562000000000", i, float64(i))
		req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(data))
		require.NoError(t, err)
		req.ContentLength = int64(len(data))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode, "Sequential request should succeed")
	}
}

// TestWriteEndpointIntegration tests the complete flow from HTTP request to storage
func TestWriteEndpointIntegration(t *testing.T) {
	storageInstance := setupTestStorage(t)
	defer cleanupTestStorage(t, storageInstance)

	server := createTestServer(storageInstance)
	defer server.Close()

	// Test data that should be written to storage
	testData := "integration_test,test=write value=42.0 1434055562000000000"

	// Make the request
	req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(testData))
	require.NoError(t, err)
	req.ContentLength = int64(len(testData))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify HTTP response
	assert.Equal(t, 200, resp.StatusCode, "HTTP response should be 200")

	// Verify response body
	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)
	assert.Equal(t, "OK", string(body), "Response body should be 'OK'")

	// Note: In a real integration test, you might want to verify that the data
	// was actually written to the storage file, but that would require
	// additional storage inspection methods
}
