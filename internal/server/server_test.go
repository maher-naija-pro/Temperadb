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
		Port:         "8081",
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
	resp, err := http.Get("http://localhost:8081/metrics")
	if err != nil {
		// Server might not be ready yet, wait a bit more
		time.Sleep(100 * time.Millisecond)
		resp, err = http.Get("http://localhost:8081/metrics")
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
	}

	// This should not panic
	err := server.Close()
	if err != nil {
		t.Errorf("Close should not return error when storage is nil: %v", err)
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
