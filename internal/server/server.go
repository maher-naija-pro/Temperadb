package server

import (
	"context"
	"net/http"
	"runtime"
	"strconv"
	"time"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/errors"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"
)

// Server represents the TimeSeriesDB server
type Server struct {
	httpServer *http.Server
	storage    *storage.Storage
	config     *config.Config
	startTime  time.Time
	status     int   // 0=stopped, 1=starting, 2=running, 3=shutting_down, 4=stopped
	connCount  int64 // active connection count
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Check for nil config
	if cfg == nil {
		return nil, errors.NewValidationError("config cannot be nil")
	}

	// Initialize logger with configuration
	logger.InitWithConfig(cfg.Logging)

	// Initialize metrics system
	metrics.Init()

	// Initialize storage with configuration
	storageInstance := storage.NewStorage(cfg.Storage)

	// Initialize API router
	router := aphttp.NewRouter(storageInstance)

	// Use custom mux for testing isolation
	mux := router.GetMux()

	// Create HTTP server with configuration
	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      mux,
	}

	// Create server instance
	server := &Server{
		httpServer: httpServer,
		storage:    storageInstance,
		config:     cfg,
		startTime:  time.Now(),
		status:     1, // starting
	}

	// Initialize server metrics
	server.initializeMetrics()

	// Set initial health status
	metrics.ServerHealth.WithLabelValues().Set(1) // healthy

	return server, nil
}

// Start starts the server
func (s *Server) Start() error {
	logger.Infof("Starting TimeSeriesDB on port %s...", s.config.Server.Port)
	logger.Infof("Configuration: %s", s.config.String())
	logger.Infof("Metrics available at: http://localhost:%s/metrics", s.config.Server.Port)

	// Update server status to running
	s.status = 2 // running
	metrics.ServerStatus.WithLabelValues().Set(float64(s.status))

	// Start metrics collection goroutine
	go s.collectMetrics()

	return s.httpServer.ListenAndServe()
}

// Close closes the server and cleans up resources
func (s *Server) Close() error {
	// Update server status to stopped
	s.status = 4 // stopped
	metrics.ServerStatus.WithLabelValues().Set(float64(s.status))

	if s.storage != nil {
		s.storage.Close()
		// Update storage connection status
		metrics.StorageConnectionStatus.WithLabelValues().Set(0) // disconnected
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server gracefully...")

	// Update server status to shutting down
	s.status = 3 // shutting_down
	metrics.ServerStatus.WithLabelValues().Set(float64(s.status))

	shutdownStart := time.Now()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
		metrics.ServerErrors.WithLabelValues("shutdown_error", "http_server").Inc()
	}

	// Close storage
	if err := s.Close(); err != nil {
		logger.Errorf("Storage close error: %v", err)
		metrics.ServerErrors.WithLabelValues("close_error", "storage").Inc()
	}

	// Record shutdown duration
	shutdownDuration := time.Since(shutdownStart).Seconds()
	metrics.ServerShutdownDuration.WithLabelValues().Observe(shutdownDuration)

	logger.Info("Server shutdown complete")
	return nil
}

// IncrementConnection increments the active connection count
func (s *Server) IncrementConnection() {
	s.connCount++
	metrics.ServerActiveConnections.WithLabelValues().Set(float64(s.connCount))
}

// DecrementConnection decrements the active connection count
func (s *Server) DecrementConnection() {
	s.connCount--
	if s.connCount < 0 {
		s.connCount = 0
	}
	metrics.ServerActiveConnections.WithLabelValues().Set(float64(s.connCount))
}

// GetActiveConnections returns the current active connection count
func (s *Server) GetActiveConnections() int64 {
	return s.connCount
}

// SetHealth sets the server health status
func (s *Server) SetHealth(healthy bool) {
	if healthy {
		metrics.ServerHealth.WithLabelValues().Set(1)
	} else {
		metrics.ServerHealth.WithLabelValues().Set(0)
	}
}

// initializeMetrics initializes server-specific metrics
func (s *Server) initializeMetrics() {
	// Set configuration metrics
	port, _ := strconv.Atoi(s.config.Server.Port)
	metrics.ServerConfigPort.WithLabelValues().Set(float64(port))
	metrics.ServerConfigReadTimeout.WithLabelValues().Set(s.config.Server.ReadTimeout.Seconds())
	metrics.ServerConfigWriteTimeout.WithLabelValues().Set(s.config.Server.WriteTimeout.Seconds())
	metrics.ServerConfigIdleTimeout.WithLabelValues().Set(s.config.Server.IdleTimeout.Seconds())

	// Set initial server status
	metrics.ServerStatus.WithLabelValues().Set(float64(s.status))
	metrics.ServerStartTime.WithLabelValues().Set(float64(s.startTime.Unix()))

	// Set storage connection status
	metrics.StorageConnectionStatus.WithLabelValues().Set(1) // connected
}

// collectMetrics continuously collects server metrics
func (s *Server) collectMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if server is still running
			if s.status == 4 { // stopped
				return
			}

			// Update uptime
			uptime := time.Since(s.startTime).Seconds()
			metrics.ServerUptime.WithLabelValues().Set(uptime)

			// Update resource metrics
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			metrics.ServerMemoryUsage.WithLabelValues().Set(float64(m.Alloc))
			metrics.ServerGoroutines.WithLabelValues().Set(float64(runtime.NumGoroutine()))
		}
	}
}

// GetMetrics returns current server metrics
func (s *Server) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"status":            s.status,
		"uptime_seconds":    time.Since(s.startTime).Seconds(),
		"start_time":        s.startTime,
		"port":              s.config.Server.Port,
		"goroutines":        runtime.NumGoroutine(),
		"storage_connected": s.storage != nil,
	}
}
