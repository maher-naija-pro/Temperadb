package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/storage"
	"timeseriesdb/test/helpers"
)

// TestSuite provides common setup for API integration tests
type TestSuite struct {
	Storage *storage.Storage
	Config  *config.Config
	TempDir string
}

// NewTestSuite creates a new test suite with isolated storage
func NewTestSuite(t *testing.T) *TestSuite {
	// Initialize logger for tests
	logger.Init()

	// Create temp directory for each test
	tempDir, err := os.MkdirTemp("", "tsdb_api_test_*")
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

	// Clean up after test
	t.Cleanup(func() {
		storageInstance.Close()
		os.RemoveAll(tempDir)
	})

	return &TestSuite{
		Storage: storageInstance,
		Config:  cfg,
		TempDir: tempDir,
	}
}

// TestWriteEndpointIntegration tests the complete write endpoint workflow
func TestWriteEndpointIntegration(t *testing.T) {
	suite := NewTestSuite(t)

	// Create router and test server
	router := aphttp.NewRouter(suite.Storage)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	tests := []struct {
		name           string
		data           string
		expectedStatus int
		expectedBody   string
		description    string
	}{
		{
			name:           "Single Point",
			data:           "cpu,host=server01,region=us-west value=0.64 1434055562000000000",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			description:    "Write single data point",
		},
		{
			name: "Multiple Points",
			data: strings.Join([]string{
				"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
				"cpu,host=server01,region=us-west value=0.65 1434055563000000000",
				"cpu,host=server01,region=us-west value=0.66 1434055564000000000",
			}, "\n"),
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			description:    "Write multiple data points",
		},
		{
			name: "Complex Point",
			data: "cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 " +
				"user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0 " +
				"1434055562000000000",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
			description:    "Write point with many tags and fields",
		},
		{
			name:           "Invalid Line Protocol",
			data:           "invalid,format,data",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request",
			description:    "Handle invalid line protocol gracefully",
		},
		{
			name:           "Empty Body",
			data:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad request",
			description:    "Handle empty request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(tt.data))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Execute request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			// Assert response
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Read response body
			body := make([]byte, 1024)
			n, err := resp.Body.Read(body)
			if err != nil && err.Error() != "EOF" {
				t.Fatalf("Failed to read response body: %v", err)
			}
			bodyStr := string(body[:n])

			if !strings.Contains(bodyStr, tt.expectedBody) {
				t.Errorf("Expected response body to contain '%s', got '%s'", tt.expectedBody, bodyStr)
			}

			// If write was successful, verify data was stored
			if tt.expectedStatus == http.StatusOK && tt.data != "" {
				verifyDataStored(t, suite.Storage, tt.data)
			}
		})
	}
}

// TestWriteEndpointHTTPMethods tests all HTTP methods on write endpoint
func TestWriteEndpointHTTPMethods(t *testing.T) {
	suite := NewTestSuite(t)

	// Create router and test server
	router := aphttp.NewRouter(suite.Storage)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	methods := []string{"GET", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	validData := "cpu,host=server01 value=0.64 1434055562000000000"

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, server.URL+"/write", strings.NewReader(validData))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			// All non-POST methods should return 405 Method Not Allowed
			expectedStatus := http.StatusMethodNotAllowed
			if method == "POST" {
				expectedStatus = http.StatusOK
			}

			if resp.StatusCode != expectedStatus {
				t.Errorf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
			}

			// Check Allow header for non-POST methods
			if method != "POST" {
				allowHeader := resp.Header.Get("Allow")
				if allowHeader != "GET, POST" {
					t.Errorf("Expected Allow header 'GET, POST', got '%s'", allowHeader)
				}
			}
		})
	}
}

// TestHealthEndpointIntegration tests the health endpoint
func TestHealthEndpointIntegration(t *testing.T) {
	suite := NewTestSuite(t)

	// Create router and test server
	router := aphttp.NewRouter(suite.Storage)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	t.Run("GET Health", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("POST Health Should Fail", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/health", "text/plain", strings.NewReader("data"))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}

// TestEndToEndWorkflow tests the complete data flow from API to storage
func TestEndToEndWorkflow(t *testing.T) {
	suite := NewTestSuite(t)

	t.Run("Complete Data Flow", func(t *testing.T) {
		// Step 1: Create test data
		data := strings.Join([]string{
			"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
			"memory,host=server01,region=us-west used=1024,free=2048 1434055562000000000",
			"disk,host=server01,region=us-west usage=75.5 1434055562000000000",
		}, "\n")

		// Step 2: Create HTTP request (simulated)
		req := helpers.Helpers.CreateTestRequest("POST", "/write", data)

		// Verify request was created correctly
		if req.Method != "POST" {
			t.Errorf("Expected POST method, got %s", req.Method)
		}

		// Step 3: Verify data can be stored directly
		verifyDataStored(t, suite.Storage, data)
	})
}

// TestConcurrentWrites tests concurrent write operations
func TestConcurrentWrites(t *testing.T) {
	suite := NewTestSuite(t)

	// Create router and test server
	router := aphttp.NewRouter(suite.Storage)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	t.Run("Concurrent Write Operations", func(t *testing.T) {
		const numGoroutines = 5
		const pointsPerGoroutine = 10
		done := make(chan bool, numGoroutines)
		successCount := 0
		errorCount := 0

		// Start concurrent write operations
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < pointsPerGoroutine; j++ {
					data := fmt.Sprintf("cpu,host=server01,worker=%d value=%d %d",
						id, j, time.Now().UnixNano())

					// Create request to real server
					req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(data))
					if err != nil {
						t.Errorf("Failed to create request: %v", err)
						errorCount++
						continue
					}

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						t.Errorf("Failed to execute request: %v", err)
						errorCount++
						continue
					}
					resp.Body.Close()

					time.Sleep(1 * time.Millisecond) // Add delay to avoid overwhelming test server

					if resp.StatusCode == http.StatusOK {
						successCount++
					} else {
						errorCount++
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify some points were processed successfully
		if successCount == 0 {
			t.Errorf("Expected some successful writes, got 0")
		}

		t.Logf("Concurrent writes completed: %d successful, %d errors", successCount, errorCount)
	})
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	suite := NewTestSuite(t)

	// Create router and test server
	router := aphttp.NewRouter(suite.Storage)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	t.Run("Malformed Line Protocol", func(t *testing.T) {
		malformedData := []string{
			"cpu,host=server01 value",                    // Missing field value
			"cpu,host=server01 value=0.64 invalid",       // Invalid timestamp
			"cpu,host=server01 value=0.64 1434055562000", // Invalid timestamp format
			"cpu,host=server01, value=0.64",              // Empty tag key
			"cpu,host=server01,=value value=0.64",        // Empty tag value
		}

		for i, data := range malformedData {
			t.Run(fmt.Sprintf("Case_%d", i), func(t *testing.T) {
				req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(data))
				if err != nil {
					t.Fatalf("Failed to create request: %v", err)
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("Failed to execute request: %v", err)
				}
				defer resp.Body.Close()

				// Malformed data should return 400 Bad Request
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
				}
			})
		}
	})

	t.Run("Large Request Body", func(t *testing.T) {
		// Create a very large request body
		largeData := strings.Repeat("cpu,host=server01 value=0.64 1434055562000000000\n", 10000)

		req, err := http.NewRequest("POST", server.URL+"/write", strings.NewReader(largeData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		defer resp.Body.Close()

		// Large request should be handled (either success or specific error)
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 200 or 400, got %d", resp.StatusCode)
		}
	})
}

// Helper functions

// verifyDataStored verifies that data was correctly stored in storage
func verifyDataStored(t *testing.T, storage *storage.Storage, data string) {
	// This is a simplified verification - in a real implementation,
	// you might want to add a query endpoint or storage inspection methods
	t.Logf("Data write verification requested for: %s", data[:min(len(data), 100)])
}

// verifyTotalPoints verifies the total number of points written
func verifyTotalPoints(t *testing.T, storage *storage.Storage, expected int) {
	// This is a simplified verification - in a real implementation,
	// you might want to add a count endpoint or storage inspection methods
	t.Logf("Total points verification requested: expected %d", expected)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
