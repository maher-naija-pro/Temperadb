package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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

// CreateTestRequestWithRandomBody creates a test request with a random body for load testing
func (h *TestHelper) CreateTestRequestWithRandomBody(method, url string, bodySize int, headers map[string]string) *http.Request {
	randomBody := h.GenerateRandomString(bodySize)
	return h.CreateTestRequest(method, url, randomBody, headers)
}

// CreateTestRequestWithLargeBody creates a test request with a large body for stress testing
func (h *TestHelper) CreateTestRequestWithLargeBody(method, url string, bodySizeMB int, headers map[string]string) *http.Request {
	largeBody := strings.Repeat("a", bodySizeMB*1024*1024)
	return h.CreateTestRequest(method, url, largeBody, headers)
}

// ExecuteRequest executes an HTTP request and returns the response
func (h *TestHelper) ExecuteRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return client.Do(req)
}

// ExecuteRequestWithTimeout executes an HTTP request with a custom timeout
func (h *TestHelper) ExecuteRequestWithTimeout(req *http.Request, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Do(req)
}

// ExecuteRequestWithRetry executes an HTTP request with retry logic
func (h *TestHelper) ExecuteRequestWithRetry(req *http.Request, maxRetries int, backoff time.Duration) (*http.Response, error) {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		resp, err := h.ExecuteRequest(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if i < maxRetries {
			time.Sleep(backoff * time.Duration(i+1))
		}
	}
	return nil, fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
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

// AssertResponseContains asserts that the response body contains the expected substring
func (h *TestHelper) AssertResponseContains(t *testing.T, resp *http.Response, expectedSubstring string, message string) {
	bodyBytes, err := readResponseBody(resp)
	require.NoError(t, err, "Failed to read response body")
	assert.Contains(t, string(bodyBytes), expectedSubstring, message)
}

// AssertResponseNotContains asserts that the response body does not contain the unexpected substring
func (h *TestHelper) AssertResponseNotContains(t *testing.T, resp *http.Response, unexpectedSubstring string, message string) {
	bodyBytes, err := readResponseBody(resp)
	require.NoError(t, err, "Failed to read response body")
	assert.NotContains(t, string(bodyBytes), unexpectedSubstring, message)
}

// AssertResponseHeaders asserts that the response has the expected headers
func (h *TestHelper) AssertResponseHeaders(t *testing.T, resp *http.Response, expectedHeaders map[string]string, message string) {
	for key, expectedValue := range expectedHeaders {
		actualValue := resp.Header.Get(key)
		assert.Equal(t, expectedValue, actualValue, "Header %s mismatch: %s", key, message)
	}
}

// AssertResponseTime asserts that the response time is within acceptable limits
func (h *TestHelper) AssertResponseTime(t *testing.T, startTime time.Time, maxDuration time.Duration, message string) {
	duration := time.Since(startTime)
	assert.Less(t, duration, maxDuration, "Response time %v exceeded limit %v: %s", duration, maxDuration, message)
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

// GenerateRandomLineProtocolData generates random line protocol data for load testing
func (h *TestHelper) GenerateRandomLineProtocolData(count int, baseTime int64) string {
	var lines []string
	measurements := []string{"cpu", "memory", "disk", "network", "database"}
	regions := []string{"us-west", "us-east", "eu-west", "asia-pacific"}
	environments := []string{"production", "staging", "development"}

	for i := 0; i < count; i++ {
		measurement := measurements[rand.Intn(len(measurements))]
		host := fmt.Sprintf("server%03d", rand.Intn(1000))
		region := regions[rand.Intn(len(regions))]
		env := environments[rand.Intn(len(environments))]

		tags := map[string]string{
			"host":   host,
			"region": region,
			"env":    env,
		}

		fields := map[string]float64{
			"value": rand.Float64() * 100,
			"load":  rand.Float64(),
			"count": float64(rand.Intn(1000)),
		}

		line := h.GenerateLineProtocolData(measurement, tags, fields, baseTime+int64(i))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GenerateStressTestData generates data designed to stress test the system
func (h *TestHelper) GenerateStressTestData(count int, baseTime int64) string {
	var lines []string

	// Generate very long measurement names
	longMeasurement := strings.Repeat("a", 1000)

	// Generate very long tag values
	longTagValue := strings.Repeat("b", 1000)

	// Generate special numeric values
	specialValues := []float64{
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		math.Pi,
		math.E,
		math.Sqrt2,
		math.Ln2,
	}

	for i := 0; i < count; i++ {
		measurement := longMeasurement + fmt.Sprintf("_%d", i)
		tags := map[string]string{
			"long_tag":   longTagValue + fmt.Sprintf("_%d", i),
			"normal_tag": fmt.Sprintf("value_%d", i),
		}

		fields := map[string]float64{
			"normal_field":  float64(i),
			"special_field": specialValues[i%len(specialValues)],
		}

		line := h.GenerateLineProtocolData(measurement, tags, fields, baseTime+int64(i))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GenerateUnicodeTestData generates data with unicode characters for internationalization testing
func (h *TestHelper) GenerateUnicodeTestData(count int, baseTime int64) string {
	var lines []string
	unicodeStrings := []string{"测试", "🚀", "café", "naïve", "über", "résumé", "🎉", "🌟", "💻", "🔥"}

	for i := 0; i < count; i++ {
		measurement := fmt.Sprintf("metric_%s_%d", unicodeStrings[i%len(unicodeStrings)], i)
		tags := map[string]string{
			"host":   fmt.Sprintf("server_%s_%d", unicodeStrings[i%len(unicodeStrings)], i),
			"region": fmt.Sprintf("region_%s_%d", unicodeStrings[i%len(unicodeStrings)], i),
		}

		fields := map[string]float64{
			"value": float64(i),
		}

		line := h.GenerateLineProtocolData(measurement, tags, fields, baseTime+int64(i))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// GenerateMalformedData generates intentionally malformed data for error handling testing
func (h *TestHelper) GenerateMalformedData(count int) string {
	var lines []string
	malformedPatterns := []string{
		"cpu,host=server01 value=0.64",                                     // Missing timestamp
		"cpu,host=server01 1434055562000000000",                            // Missing fields
		"cpu,host=server01,malformed_tag value=0.64 1434055562000000000",   // Malformed tag
		"cpu,host=server01 value=0.64,malformed_field 1434055562000000000", // Malformed field
		"cpu,host=server01 value=not_a_number 1434055562000000000",         // Invalid field value
		"cpu,host=server01 value=0.64 invalid_timestamp",                   // Invalid timestamp
		"",          // Empty line
		"   \n\t  ", // Whitespace only
	}

	for i := 0; i < count; i++ {
		pattern := malformedPatterns[i%len(malformedPatterns)]
		if pattern != "" {
			lines = append(lines, pattern)
		}
	}
	return strings.Join(lines, "\n")
}

// GenerateRandomString generates a random string of specified length
func (h *TestHelper) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomTags generates random tags for testing
func (h *TestHelper) GenerateRandomTags(count int) map[string]string {
	tags := make(map[string]string)
	keyPrefixes := []string{"host", "region", "env", "dc", "instance", "service", "version"}
	valuePrefixes := []string{"server", "us-west", "production", "dc1", "i-123", "api", "v1"}

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%s_%d", keyPrefixes[i%len(keyPrefixes)], i)
		value := fmt.Sprintf("%s_%d", valuePrefixes[i%len(valuePrefixes)], i)
		tags[key] = value
	}
	return tags
}

// GenerateRandomFields generates random fields for testing
func (h *TestHelper) GenerateRandomFields(count int) map[string]float64 {
	fields := make(map[string]float64)
	keyPrefixes := []string{"value", "load", "count", "rate", "latency", "throughput", "error_rate"}

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%s_%d", keyPrefixes[i%len(keyPrefixes)], i)
		value := rand.Float64() * 1000
		fields[key] = value
	}
	return fields
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

// BenchmarkHTTPEndpointWithTimeout benchmarks an HTTP endpoint with custom timeout
func (h *BenchmarkHelper) BenchmarkHTTPEndpointWithTimeout(b *testing.B, method, url, body string, headers map[string]string, timeout time.Duration) {
	helper := NewTestHelper()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := helper.CreateTestRequest(method, url, body, headers)
		resp, err := helper.ExecuteRequestWithTimeout(req, timeout)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}

// BenchmarkDataGeneration benchmarks data generation functions
func (h *BenchmarkHelper) BenchmarkDataGeneration(b *testing.B, count int) {
	helper := NewTestHelper()
	baseTime := time.Now().UnixNano()

	b.Run("LineProtocolData", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateLineProtocolData(
				"cpu",
				map[string]string{"host": "server01"},
				map[string]float64{"value": 0.64},
				baseTime,
			)
		}
	})

	b.Run("BulkLineProtocolData", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateBulkLineProtocolData(count, baseTime)
		}
	})

	b.Run("RandomLineProtocolData", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomLineProtocolData(count, baseTime)
		}
	})
}

// BenchmarkConcurrentRequests benchmarks concurrent HTTP requests
func (h *BenchmarkHelper) BenchmarkConcurrentRequests(b *testing.B, method, url, body string, headers map[string]string, concurrency int) {
	helper := NewTestHelper()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
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
	})
}
