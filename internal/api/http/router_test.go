package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"
)

func TestNewRouter(t *testing.T) {
	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_storage.tsv")
		os.RemoveAll("test_backups")
	}()

	router := NewRouter(storageInstance)

	if router == nil {
		t.Error("Expected router to be created, got nil")
	}

	if router.writeHandler == nil {
		t.Error("Expected write handler to be created")
	}

	if router.healthHandler == nil {
		t.Error("Expected health handler to be created")
	}

	if router.metricsMiddleware == nil {
		t.Error("Expected metrics middleware to be created")
	}
}

func TestRouter_GetMux(t *testing.T) {
	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_mux_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_mux",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_mux_storage.tsv")
		os.RemoveAll("test_backups_mux")
	}()

	router := NewRouter(storageInstance)
	mux := router.GetMux()

	if mux == nil {
		t.Error("Expected mux to be created, got nil")
	}

	// Test that the mux has the expected routes
	server := httptest.NewServer(mux)
	defer server.Close()

	// Test health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected health endpoint to return 200, got %d", resp.StatusCode)
	}

	// Test write endpoint (should return 405 for GET)
	resp, err = http.Get(server.URL + "/write")
	if err != nil {
		t.Fatalf("Failed to make request to write endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected write endpoint to return 405 for GET, got %d", resp.StatusCode)
	}

	// Test metrics endpoint
	resp, err = http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to make request to metrics endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected metrics endpoint to return 200, got %d", resp.StatusCode)
	}
}

func TestRouter_Integration_WriteEndpoint(t *testing.T) {
	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_write_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_write",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_write_storage.tsv")
		os.RemoveAll("test_backups_write")
	}()

	router := NewRouter(storageInstance)
	mux := router.GetMux()

	server := httptest.NewServer(mux)
	defer server.Close()

	// Test POST to write endpoint (should work)
	resp, err := http.Post(server.URL+"/write", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to make POST request to write endpoint: %v", err)
	}
	defer resp.Body.Close()

	// The write endpoint should handle the request (even if it fails due to empty body)
	// We're just testing that the route is properly registered and accessible
	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected write endpoint to handle POST request, got %d", resp.StatusCode)
	}
}

func TestRouter_Integration_HealthEndpoint(t *testing.T) {
	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_health_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_health",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_health_storage.tsv")
		os.RemoveAll("test_backups_health")
	}()

	router := NewRouter(storageInstance)
	mux := router.GetMux()

	server := httptest.NewServer(mux)
	defer server.Close()

	// Test GET to health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make GET request to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected health endpoint to return 200, got %d", resp.StatusCode)
	}

	// Test POST to health endpoint (should return 405)
	resp, err = http.Post(server.URL+"/health", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to make POST request to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected health endpoint to return 405 for POST, got %d", resp.StatusCode)
	}
}

func TestRouter_Integration_MetricsEndpoint(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_metrics_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_metrics",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_metrics_storage.tsv")
		os.RemoveAll("test_backups_metrics")
	}()

	router := NewRouter(storageInstance)
	mux := router.GetMux()

	server := httptest.NewServer(mux)
	defer server.Close()

	// Test GET to metrics endpoint
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to make GET request to metrics endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected metrics endpoint to return 200, got %d", resp.StatusCode)
	}

	// Check that response contains Prometheus metrics
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)
	content := string(body[:n])

	// Should contain some Prometheus metrics
	if len(content) == 0 {
		t.Error("Expected metrics endpoint to return non-empty response")
	}
}
