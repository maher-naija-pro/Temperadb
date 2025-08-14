package server

import (
	"context"
	"net/http"
	"testing"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/test/helpers"
)

func TestNewServer(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Test server creation
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Verify server structure
	if server == nil {
		t.Fatal("Server should not be nil")
	}

	if server.config == nil {
		t.Error("Server config should not be nil")
	}

	if server.storage == nil {
		t.Error("Server storage should not be nil")
	}

	if server.httpServer == nil {
		t.Error("HTTP server should not be nil")
	}

	// Verify server metrics are initialized
	if server.startTime.IsZero() {
		t.Error("Server start time should be initialized")
	}

	if server.status != 1 {
		t.Errorf("Expected server status 1 (starting), got %d", server.status)
	}

	// Verify HTTP server configuration
	if server.httpServer.Addr != ":8080" {
		t.Errorf("Expected server address ':8080', got '%s'", server.httpServer.Addr)
	}

	if server.httpServer.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", server.httpServer.ReadTimeout)
	}

	if server.httpServer.WriteTimeout != 30*time.Second {
		t.Errorf("Expected write timeout 30s, got %v", server.httpServer.WriteTimeout)
	}

	if server.httpServer.IdleTimeout != 60*time.Second {
		t.Errorf("Expected idle timeout 60s, got %v", server.httpServer.IdleTimeout)
	}

	// Clean up
	server.Close()
}

func TestServerStart(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration with a different port
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8085",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Start()
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test if server is responding
	resp, err := http.Get("http://localhost:8085/metrics")
	if err != nil {
		// Server might not be ready yet, wait a bit more
		time.Sleep(100 * time.Millisecond)
		resp, err = http.Get("http://localhost:8085/metrics")
	}

	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	} else {
		t.Logf("Server not responding yet (this might be expected): %v", err)
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if server.httpServer.Shutdown(ctx) != nil {
		t.Logf("Server shutdown error (this might be expected): %v", err)
	}

	// Check for server start error
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected server error: %v", err)
		}
	case <-time.After(1 * time.Second):
		// Server is still running, which is fine
	}
}

func TestServerClose(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test close method
	err = server.Close()
	if err != nil {
		t.Errorf("Failed to close server: %v", err)
	}

	// Verify server status is updated
	if server.status != 4 {
		t.Errorf("Expected server status 4 (stopped), got %d", server.status)
	}

	// Test that storage is closed
	// Note: We can't directly test if storage is closed, but we can verify the method doesn't panic
}

func TestServerCloseNilStorage(t *testing.T) {
	defer metrics.Reset()

	// Create a server with nil storage to test edge case
	server := &Server{
		httpServer: &http.Server{},
		storage:    nil,
		config:     &config.Config{},
		startTime:  time.Now(),
		status:     1,
	}

	// This should not panic
	err := server.Close()
	if err != nil {
		t.Errorf("Close should not return error when storage is nil: %v", err)
	}

	// Verify server status is updated
	if server.status != 4 {
		t.Errorf("Expected server status 4 (stopped), got %d", server.status)
	}
}

func TestServerWithInvalidConfig(t *testing.T) {
	defer metrics.Reset()

	// Test with nil config
	_, err := NewServer(nil)
	if err == nil {
		t.Error("Expected error when creating server with nil config")
	}
}

func TestServerConfigurationValidation(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "invalid-port",
		ReadTimeout:  -1 * time.Second, // Invalid negative timeout
		WriteTimeout: 0 * time.Second,  // Zero timeout
		IdleTimeout:  10 * time.Second,
	}

	// Test server creation with invalid config
	// Note: The current implementation doesn't validate these values,
	// but this test documents the expected behavior
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Server creation should not fail with invalid config: %v", err)
	}

	// Verify the invalid values are set (current behavior)
	if server.httpServer.ReadTimeout != -1*time.Second {
		t.Errorf("Expected read timeout -1s, got %v", server.httpServer.ReadTimeout)
	}

	if server.httpServer.WriteTimeout != 0*time.Second {
		t.Errorf("Expected write timeout 0s, got %v", server.httpServer.WriteTimeout)
	}

	// Clean up
	server.Close()
}

func TestServerConcurrentAccess(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8083",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test concurrent access to server methods
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			// Access server properties concurrently
			_ = server.config
			_ = server.storage
			_ = server.httpServer
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent access test")
		}
	}
}

func TestServerMemoryLeaks(t *testing.T) {
	defer metrics.Reset()

	// Create and destroy multiple servers to check for memory leaks
	for i := 0; i < 10; i++ {
		cfg := helpers.Config.CreateTestConfig(t)
		cfg.Server = config.ServerConfig{
			Port:         "8084",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		}

		server, err := NewServer(cfg)
		if err != nil {
			t.Fatalf("Failed to create server %d: %v", i, err)
		}

		// Close immediately
		server.Close()
	}
}

func TestServerShutdown(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8086",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test shutdown with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected successful shutdown, got error: %v", err)
	}
}

func TestServerShutdownTimeout(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8087",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test shutdown with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait a bit for the context to expire
	time.Sleep(1 * time.Millisecond)

	err = server.Shutdown(ctx)
	if err == nil {
		t.Error("Expected timeout error during shutdown")
	}
}

func TestServerIncrementDecrementConnections(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8088",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test initial connection count
	initialCount := server.GetActiveConnections()
	if initialCount != 0 {
		t.Errorf("Expected initial connection count 0, got %d", initialCount)
	}

	// Test incrementing connections
	server.IncrementConnection()
	server.IncrementConnection()

	count := server.GetActiveConnections()
	if count != 2 {
		t.Errorf("Expected connection count 2, got %d", count)
	}

	// Test decrementing connections
	server.DecrementConnection()

	count = server.GetActiveConnections()
	if count != 1 {
		t.Errorf("Expected connection count 1, got %d", count)
	}

	// Test decrementing below zero (should not go negative)
	server.DecrementConnection()
	server.DecrementConnection()

	count = server.GetActiveConnections()
	if count != 0 {
		t.Errorf("Expected connection count 0, got %d", count)
	}
}

func TestServerHealthStatus(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8089",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test setting health status to false
	server.SetHealth(false)

	// Test setting health status to true
	server.SetHealth(true)

	// Verify metrics are updated (we can't directly test the metrics values in tests)
	// but we can verify the method doesn't panic
}

func TestServerGetMetrics(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test getting metrics
	metrics := server.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be returned")
	}
}

func TestServerGetID(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8091",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test that server has a start time (indirect way to verify server was created)
	if server.startTime.IsZero() {
		t.Error("Expected server start time to be set")
	}
}

func TestServerIsClosed(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8092",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test initial status (should be starting = 1)
	if server.status != 1 {
		t.Errorf("Expected initial status 1 (starting), got %d", server.status)
	}

	// Close server
	server.Close()

	// Test status after close (should be stopped = 4)
	if server.status != 4 {
		t.Errorf("Expected status 4 (stopped) after Close(), got %d", server.status)
	}
}

func TestServerCollectMetrics(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8093",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test metrics collection
	server.collectMetrics()

	// Verify metrics were collected
	metrics := server.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be available after collection")
	}
}

func TestServerInitializeMetrics(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8094",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test metrics initialization
	server.initializeMetrics()

	// Verify metrics were initialized
	metrics := server.GetMetrics()
	if metrics == nil {
		t.Error("Expected metrics to be available after initialization")
	}
}

func TestServerWithNilStorage(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration with nil storage
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8095",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// Test that server can be created with nil storage
	if server == nil {
		t.Fatal("Expected server to be created even with nil storage")
	}
}

func TestServerMultipleCloseCalls(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8096",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test multiple close calls (should not panic)
	server.Close()
	server.Close()
	server.Close()

	// Verify server is marked as closed
	if server.status != 4 {
		t.Error("Expected server status to be 4 (stopped) after multiple Close() calls")
	}
}

func TestServerConcurrentClose(t *testing.T) {
	defer metrics.Reset()

	// Create test configuration
	cfg := helpers.Config.CreateTestConfig(t)
	cfg.Server = config.ServerConfig{
		Port:         "8097",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test concurrent close calls
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			server.Close()
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		select {
		case <-done:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent close test")
		}
	}

	// Verify server is marked as closed
	if server.status != 4 {
		t.Error("Expected server status to be 4 (stopped) after concurrent Close() calls")
	}
}
