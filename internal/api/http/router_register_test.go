package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"
)

func TestRouter_RegisterRoutes(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	// Initialize logger for testing
	logger.Init()

	// Reset the default mux to avoid conflicts
	http.DefaultServeMux = http.NewServeMux()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataFile:    "test_router_register_storage.tsv",
		DataDir:     "test_data_register",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_register",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_register_storage.tsv")
		os.RemoveAll("test_backups_register")
	}()

	router := NewRouter(storageInstance)
	router.RegisterRoutes()

	// Test that routes are registered by making requests
	server := httptest.NewServer(http.DefaultServeMux)
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

func TestRouter_RegisterRoutes_WriteEndpointPOST(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	// Initialize logger for testing
	logger.Init()

	// Reset the default mux to avoid conflicts
	http.DefaultServeMux = http.NewServeMux()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_register_write_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_register_write",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_register_write_storage.tsv")
		os.RemoveAll("test_backups_register_write")
	}()

	router := NewRouter(storageInstance)
	router.RegisterRoutes()

	// Test that routes are registered by making requests
	server := httptest.NewServer(http.DefaultServeMux)
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

func TestRouter_RegisterRoutes_HealthEndpointPOST(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	// Reset the default mux to avoid conflicts
	http.DefaultServeMux = http.NewServeMux()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_router_register_health_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_register_health",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_router_register_health_storage.tsv")
		os.RemoveAll("test_backups_register_health")
	}()

	router := NewRouter(storageInstance)
	router.RegisterRoutes()

	// Test that routes are registered by making requests
	server := httptest.NewServer(http.DefaultServeMux)
	defer server.Close()

	// Test POST to health endpoint (should return 405)
	resp, err := http.Post(server.URL+"/health", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to make POST request to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected health endpoint to return 405 for POST, got %d", resp.StatusCode)
	}
}
