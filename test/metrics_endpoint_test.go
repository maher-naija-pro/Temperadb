package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"
)

// Test setup function to reset metrics between tests
func setupTest(t *testing.T) func() {
	// Reset metrics before each test
	metrics.Reset()

	// Return cleanup function
	return func() {
		// Cleanup after test if needed
	}
}

func TestMetricsEndpoint(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Initialize metrics system FIRST
	metrics.Init()

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitWithConfig(cfg.Logging)

	// Initialize storage
	storageInstance := storage.NewStorage(cfg.Storage)
	defer storageInstance.Close()

	// Initialize router AFTER metrics initialization
	router := aphttp.NewRouter(storageInstance)

	// Create test server using the router's mux
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	// Test metrics endpoint
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check if response contains Prometheus metrics
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	content := string(body)

	t.Logf("Response status: %d", resp.StatusCode)
	t.Logf("Response headers: %v", resp.Header)
	t.Logf("Response body length: %d", len(content))
	t.Logf("Response body: %q", content)

	// Check for some expected metrics that should be available
	expectedMetrics := []string{
		"# HELP tsdb_batch_queue_wait_seconds",
		"# TYPE tsdb_batch_queue_wait_seconds histogram",
		"# HELP tsdb_compaction_runs_total",
		"# TYPE tsdb_compaction_runs_total counter",
		"# HELP tsdb_ingestion_batches_total",
		"# TYPE tsdb_ingestion_batches_total counter",
	}

	for _, expected := range expectedMetrics {
		if !strings.Contains(content, expected) {
			t.Errorf("Expected metric not found: %s", expected)
		}
	}

	t.Logf("Metrics endpoint response: %s", content[:min(len(content), 500)])
}

func TestHealthEndpoint(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Load test configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitWithConfig(cfg.Logging)

	// Initialize storage
	storageInstance := storage.NewStorage(cfg.Storage)
	defer storageInstance.Close()

	// Initialize metrics system BEFORE creating the router
	metrics.Init()

	// Initialize router AFTER metrics initialization
	router := aphttp.NewRouter(storageInstance)

	// Create test server using the router's mux
	server := httptest.NewServer(router.GetMux())
	defer server.Close()
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to get health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test write endpoint to generate some metrics
	resp, err = http.Post(server.URL+"/write", "application/json", strings.NewReader(`{"test": "data"}`))
	if err != nil {
		t.Logf("Write endpoint test failed (expected for test data): %v", err)
	} else {
		resp.Body.Close()
	}

	// Now check metrics again to see if request metrics were recorded
	resp, err = http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics after request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	content := string(body)

	// Check if HTTP request metrics are present
	// Note: The metrics middleware should record metrics for all HTTP requests
	// including the health endpoint request we just made
	if !strings.Contains(content, "tsdb_http_requests_total") {
		t.Logf("Full metrics response: %s", content)
		t.Error("HTTP request metrics not found - metrics middleware may not be working properly")
	}

	t.Logf("Metrics after request: %s", content[:min(len(content), 1000)])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
