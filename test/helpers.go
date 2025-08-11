package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelper provides utility functions for testing
type TestHelper struct{}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

// CreateTestRequest creates a test HTTP request with the given method, body, and headers
func (h *TestHelper) CreateTestRequest(method, url, body string, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		panic(fmt.Sprintf("Failed to create test request: %v", err))
	}

	// Set default headers
	if headers == nil {
		headers = make(map[string]string)
	}

	// Set Content-Type for POST requests if not specified
	if method == http.MethodPost && headers["Content-Type"] == "" {
		headers["Content-Type"] = "text/plain"
	}

	// Apply custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req
}

// ExecuteRequest executes an HTTP request and returns the response
func (h *TestHelper) ExecuteRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return client.Do(req)
}

// AssertResponseStatus asserts that the response has the expected status code
func (h *TestHelper) AssertResponseStatus(t *testing.T, resp *http.Response, expectedStatus int, message string) {
	assert.Equal(t, expectedStatus, resp.StatusCode, message)
}

// AssertResponseBody asserts that the response body matches the expected content
func (h *TestHelper) AssertResponseBody(t *testing.T, resp *http.Response, expectedBody string, message string) {
	bodyBytes, err := readResponseBody(resp)
	require.NoError(t, err, "Failed to read response body")
	assert.Equal(t, expectedBody, string(bodyBytes), message)
}

// AssertResponseJSON asserts that the response body contains valid JSON matching the expected structure
func (h *TestHelper) AssertResponseJSON(t *testing.T, resp *http.Response, expected interface{}, message string) {
	bodyBytes, err := readResponseBody(resp)
	require.NoError(t, err, "Failed to read response body")

	var actual interface{}
	err = json.Unmarshal(bodyBytes, &actual)
	require.NoError(t, err, "Failed to parse response as JSON")

	assert.Equal(t, expected, actual, message)
}

// GenerateLineProtocolData generates InfluxDB line protocol data for testing
func (h *TestHelper) GenerateLineProtocolData(measurement string, tags map[string]string, fields map[string]float64, timestamp int64) string {
	var parts []string

	// Measurement and tags
	measurementPart := measurement
	if len(tags) > 0 {
		var tagParts []string
		for k, v := range tags {
			tagParts = append(tagParts, fmt.Sprintf("%s=%s", k, v))
		}
		measurementPart = fmt.Sprintf("%s,%s", measurement, strings.Join(tagParts, ","))
	}
	parts = append(parts, measurementPart)

	// Fields
	var fieldParts []string
	for k, v := range fields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%f", k, v))
	}
	parts = append(parts, strings.Join(fieldParts, ","))

	// Timestamp
	parts = append(parts, fmt.Sprintf("%d", timestamp))

	return strings.Join(parts, " ")
}

// GenerateBulkLineProtocolData generates multiple lines of InfluxDB line protocol data
func (h *TestHelper) GenerateBulkLineProtocolData(count int, baseTime int64) string {
	var lines []string
	for i := 0; i < count; i++ {
		line := h.GenerateLineProtocolData(
			"cpu",
			map[string]string{"host": fmt.Sprintf("server%02d", i)},
			map[string]float64{"value": float64(i)},
			baseTime+int64(i),
		)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// readResponseBody reads and returns the response body as bytes
func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// BenchmarkHelper provides utilities for benchmarking
type BenchmarkHelper struct{}

// NewBenchmarkHelper creates a new benchmark helper instance
func NewBenchmarkHelper() *BenchmarkHelper {
	return &BenchmarkHelper{}
}

// BenchmarkHTTPEndpoint benchmarks an HTTP endpoint with the given request
func (h *BenchmarkHelper) BenchmarkHTTPEndpoint(b *testing.B, method, url, body string, headers map[string]string) {
	helper := NewTestHelper()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := helper.CreateTestRequest(method, url, body, headers)
		resp, err := helper.ExecuteRequest(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}
