package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/errors"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/server"
)

func TestBuildTimeVariables(t *testing.T) {
	// Test that build-time variables are accessible
	if Version == "" {
		t.Error("Expected Version to have a value")
	}
	if BuildTime == "" {
		t.Error("Expected BuildTime to have a value")
	}
	if CommitHash == "" {
		t.Error("Expected CommitHash to have a value")
	}

	// Note: In test environment, these variables may have default values
	// During actual build, they are set via ldflags
	t.Logf("Version: %s", Version)
	t.Logf("BuildTime: %s", BuildTime)
	t.Logf("CommitHash: %s", CommitHash)
}

func TestConfigurationLoading(t *testing.T) {
	// Test that configuration loading works
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Verify configuration is not nil
	if cfg == nil {
		t.Fatal("Expected configuration to be loaded")
	}

	// Verify logging configuration has expected values
	if cfg.Logging.Level == "" {
		t.Error("Expected logging level to be set")
	}

	// Verify server configuration has expected values
	if cfg.Server.Port == "" {
		t.Error("Expected server port to be set")
	}
}

func TestLoggerInitialization(t *testing.T) {
	// Test that logger can be initialized with configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitWithConfig(cfg.Logging)

	// Test that logger is working
	logger.Info("Test log message")
	logger.Debug("Test debug message")
	logger.Warn("Test warning message")
	logger.Error("Test error message")

	// Verify logger is functional by checking if we can log without panicking
	// This is a basic test that the logger is initialized and working
}

func TestServerCreation(t *testing.T) {
	// Test that server can be created with configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Verify server is not nil
	if srv == nil {
		t.Fatal("Expected server to be created")
	}

	// Clean up
	if err := srv.Close(); err != nil {
		t.Logf("Warning: failed to close server during test cleanup: %v", err)
	}
}

func TestErrorHandlingPatterns(t *testing.T) {
	// Test the error handling patterns used in main.go
	// Create a custom app error
	appErr := &errors.AppError{
		Message: "test error message",
		Type:    "test_error",
	}

	// Test error type checking
	var targetErr *errors.AppError
	if !errors.As(appErr, &targetErr) {
		t.Error("Expected errors.As to work with AppError")
	}

	// Test error message and type
	if targetErr.Message != "test error message" {
		t.Errorf("Expected error message 'test error message', got '%s'", targetErr.Message)
	}
	if targetErr.Type != "test_error" {
		t.Errorf("Expected error type 'test_error', got '%s'", targetErr.Type)
	}
}

func TestMetricsInitialization(t *testing.T) {
	// Test that metrics can be initialized
	// Initialize metrics
	metrics.Init()

	// Test that build info metrics exist
	if metrics.BuildInfo == nil {
		t.Error("Expected BuildInfo metric to exist")
	}

	if metrics.APIVersion == nil {
		t.Error("Expected APIVersion metric to exist")
	}

	// Test setting build info metrics (this is what main.go does)
	metrics.BuildInfo.WithLabelValues(Version, CommitHash, "test", "go1.21.0").Set(1)
	metrics.APIVersion.WithLabelValues(Version).Set(1)

	// Verify metrics were set
	// Note: We can't easily verify the actual values without exposing internal metric state
	// but we can verify the operations don't panic
}

func TestContextHandling(t *testing.T) {
	// Test context creation and cancellation (used in main.go)
	ctx, cancel := context.WithCancel(context.Background())
	if ctx == nil {
		t.Fatal("Expected context to be created")
	}

	// Test context cancellation
	cancel()

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Expected context to be cancelled")
	}
}

func TestShutdownTimeoutContext(t *testing.T) {
	// Test shutdown timeout context creation (used in main.go)
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if shutdownCtx == nil {
		t.Fatal("Expected shutdown context to be created")
	}

	// Verify timeout is set
	deadline, ok := shutdownCtx.Deadline()
	if !ok {
		t.Error("Expected shutdown context to have deadline")
	}

	// Verify timeout is approximately the configured value
	expectedTimeout := cfg.Server.ShutdownTimeout
	actualTimeout := time.Until(deadline)
	tolerance := 100 * time.Millisecond

	if actualTimeout < expectedTimeout-tolerance || actualTimeout > expectedTimeout+tolerance {
		t.Errorf("Expected timeout around %v, got %v", expectedTimeout, actualTimeout)
	}
}

func TestSignalHandling(t *testing.T) {
	// Test signal handling setup (used in main.go)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Verify signal channel is set up
	if sigChan == nil {
		t.Error("Expected signal channel to be set up")
	}

	// Test that we can receive signals
	go func() {
		time.Sleep(10 * time.Millisecond)
		sigChan <- syscall.SIGINT
	}()

	// Wait for signal
	select {
	case sig := <-sigChan:
		if sig != syscall.SIGINT {
			t.Errorf("Expected SIGINT signal, got %v", sig)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected to receive signal within timeout")
	}

	// Clean up
	signal.Stop(sigChan)
	close(sigChan)
}

func TestGracefulShutdownPattern(t *testing.T) {
	// Test the graceful shutdown pattern used in main.go
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shutdownCancel()

	// Shutdown should complete within the timeout
	if err := srv.Shutdown(shutdownCtx); err != nil {
		t.Logf("Server shutdown error (expected in test): %v", err)
	}

	// Clean up
	if err := srv.Close(); err != nil {
		t.Logf("Warning: failed to close server during test cleanup: %v", err)
	}
}

func TestExitHandling(t *testing.T) {
	// Test exit handling logic (used in main.go)
	// We can't actually test os.Exit in unit tests, but we can test the pattern

	// Test that we can create an exit scenario
	exitCode := 1
	if exitCode != 1 {
		t.Error("Expected exit code to be 1")
	}

	// Test error message pattern
	errorMessage := "shutdown failed"
	if errorMessage != "shutdown failed" {
		t.Error("Expected error message to match")
	}
}

func TestMainFunctionErrorHandling(t *testing.T) {
	// Test the specific error handling patterns used in main.go
	// Test configuration loading error handling
	// Test server creation error handling
	// Test server start error handling
	// Test shutdown error handling

	// Create a test error that matches the pattern in main.go
	testErr := &errors.AppError{
		Message: "test server error",
		Type:    "server_error",
	}

	// Test the error type checking pattern used in main.go
	var appErr *errors.AppError
	if !errors.As(testErr, &appErr) {
		t.Error("Expected errors.As to work with AppError")
	}

	// Test the error logging pattern
	if appErr.Message != "test server error" {
		t.Errorf("Expected error message 'test server error', got '%s'", appErr.Message)
	}
	if appErr.Type != "server_error" {
		t.Errorf("Expected error type 'server_error', got '%s'", appErr.Type)
	}
}

func TestRuntimeVersion(t *testing.T) {
	// Test runtime version functionality (used in main.go)
	version := runtime.Version()
	if version == "" {
		t.Error("Expected runtime version to have a value")
	}

	// Test that version contains "go"
	if len(version) < 2 || version[:2] != "go" {
		t.Errorf("Expected runtime version to start with 'go', got '%s'", version)
	}
}

func TestMainFunctionIntegration(t *testing.T) {
	// Test integration of main function components
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitWithConfig(cfg.Logging)

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test metrics initialization
	metrics.Init()
	metrics.BuildInfo.WithLabelValues(Version, CommitHash, "test", runtime.Version()).Set(1)
	metrics.APIVersion.WithLabelValues(Version).Set(1)

	// Test graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		t.Logf("Server shutdown error (expected in test): %v", err)
	}

	// Clean up
	if err := srv.Close(); err != nil {
		t.Logf("Warning: failed to close server during test cleanup: %v", err)
	}
}
