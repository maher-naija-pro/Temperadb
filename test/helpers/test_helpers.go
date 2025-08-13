package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"timeseriesdb/internal/types"
)

// TestHelpers provides common test helper functions
type TestHelpers struct{}

// CreateTestPoint creates a test point with default values
func (h *TestHelpers) CreateTestPoint(measurement string, tags map[string]string, fields map[string]float64) types.Point {
	if tags == nil {
		tags = map[string]string{"host": "server01", "region": "us-west"}
	}
	if fields == nil {
		fields = map[string]float64{"value": 0.64}
	}

	return types.Point{
		Measurement: measurement,
		Tags:        tags,
		Fields:      fields,
		Timestamp:   time.Unix(0, 1434055562000000000),
	}
}

// CreateTestPoints creates multiple test points
func (h *TestHelpers) CreateTestPoints(count int, measurement string, tags map[string]string, fields map[string]float64) []types.Point {
	points := make([]types.Point, count)
	for i := 0; i < count; i++ {
		points[i] = h.CreateTestPoint(measurement, tags, fields)
		// Add some variation to timestamps
		points[i].Timestamp = time.Unix(0, 1434055562000000000+int64(i)*1000000000)
	}
	return points
}

// CreateTestRequest creates an HTTP request for testing
func (h *TestHelpers) CreateTestRequest(method, path, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, bytes.NewBufferString(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	return req
}

// CreateTestResponse creates a response recorder for testing
func (h *TestHelpers) CreateTestResponse() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// AssertJSONResponse checks if a response contains expected JSON
func (h *TestHelpers) AssertJSONResponse(t *testing.T, response *httptest.ResponseRecorder, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
		return
	}

	// Simple equality check - in real tests you might want to use a more sophisticated comparison
	if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
		t.Errorf("Expected response %v, got %v", expected, actual)
	}
}

// AssertHTTPStatus checks if a response has the expected status code
func (h *TestHelpers) AssertHTTPStatus(t *testing.T, response *httptest.ResponseRecorder, expectedStatus int) {
	if response.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, response.Code)
	}
}

// AssertResponseBodyContains checks if response body contains expected text
func (h *TestHelpers) AssertResponseBodyContains(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	body := response.Body.String()
	if !bytes.Contains(response.Body.Bytes(), []byte(expected)) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expected, body)
	}
}

// GenerateTestData generates test data of specified size
func (h *TestHelpers) GenerateTestData(size int, measurement string) string {
	var lines []string
	for i := 0; i < size; i++ {
		line := fmt.Sprintf("%s,host=server01,region=us-west value=%d %d",
			measurement, i, 1434055562000000000+i*1000000000)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Global instance for easy access
var Helpers = &TestHelpers{}
