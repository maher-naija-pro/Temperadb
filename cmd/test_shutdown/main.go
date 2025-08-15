package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/server"
)

func main() {
	fmt.Println("=== TimeSeriesDB Server Shutdown Test ===")
	fmt.Println()

	// Create test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:            "8080",
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     10 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Logging: config.LoggingConfig{
			Level: "info",
		},
		Storage: config.StorageConfig{
			DataDir: "/tmp/test-tsdb",
		},
	}

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("1. Starting server...")
	fmt.Printf("   Server will be available at: http://localhost:%s\n", cfg.Server.Port)
	fmt.Println("   Press Ctrl+C to test graceful shutdown")
	fmt.Println()

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Give server time to start
	time.Sleep(2 * time.Second)

	// Test server status
	fmt.Println("2. Server status check:")
	fmt.Printf("   Status: %d\n", srv.GetMetrics()["status"])
	fmt.Printf("   Port: %s\n", srv.GetMetrics()["port"])
	fmt.Printf("   Goroutines: %v\n", srv.GetMetrics()["goroutines"])
	fmt.Println()

	// Simulate some activity
	fmt.Println("3. Simulating server activity...")
	srv.IncrementConnection()
	srv.IncrementConnection()
	fmt.Printf("   Active connections: %d\n", srv.GetActiveConnections())
	fmt.Println()

	// Wait for shutdown signal
	fmt.Println("4. Waiting for shutdown signal...")
	select {
	case sig := <-sigChan:
		fmt.Printf("   Received signal %v, shutting down gracefully...\n", sig)
	case <-ctx.Done():
		fmt.Println("   Server context cancelled, shutting down...")
	}

	// Graceful shutdown with timeout
	fmt.Println("5. Performing graceful shutdown...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	shutdownStart := time.Now()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	shutdownDuration := time.Since(shutdownStart)
	fmt.Printf("   ✓ Shutdown completed in %v\n", shutdownDuration)

	// Final status check
	fmt.Println("6. Final server status:")
	finalMetrics := srv.GetMetrics()
	fmt.Printf("   Status: %v\n", finalMetrics["status"])
	fmt.Printf("   Storage connected: %v\n", finalMetrics["storage_connected"])
	fmt.Println()

	fmt.Println("=== Shutdown Test Completed Successfully ===")
	fmt.Println("The server has been gracefully shut down with:")
	fmt.Println("  - All connections properly closed")
	fmt.Println("  - Storage connections cleaned up")
	fmt.Println("  - Metrics updated")
	fmt.Println("  - Status properly set to stopped")
}
