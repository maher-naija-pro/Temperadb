package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/errors"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger early
	logger.InitWithConfig(cfg.Logging)
	logger.Info("Starting TimeSeriesDB...")

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		// Check if it's our custom error type
		var appErr *errors.AppError
		if errors.As(err, &appErr) {
			logger.Fatalf("Failed to create server: %s (Type: %s)", appErr.Message, appErr.Type)
		} else {
			logger.Fatalf("Failed to create server: %v", err)
		}
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			// Check if it's our custom error type
			var appErr *errors.AppError
			if errors.As(err, &appErr) {
				logger.Errorf("Server error: %s (Type: %s)", appErr.Message, appErr.Type)
			} else {
				logger.Errorf("Server error: %v", err)
			}
			cancel()
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-sigChan:
		logger.Infof("Received signal %v, shutting down gracefully...", sig)
	case <-ctx.Done():
		logger.Info("Server context cancelled, shutting down...")
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Check if it's our custom error type
		var appErr *errors.AppError
		if errors.As(err, &appErr) {
			logger.Errorf("Error during shutdown: %s (Type: %s)", appErr.Message, appErr.Type)
		} else {
			logger.Errorf("Error during shutdown: %v", err)
		}
		os.Exit(1)
	}

	logger.Info("TimeSeriesDB shutdown complete")
}
