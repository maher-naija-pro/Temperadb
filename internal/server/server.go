package server

import (
	"fmt"
	"net/http"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"
)

// Server represents the TimeSeriesDB server
type Server struct {
	httpServer *http.Server
	storage    *storage.Storage
	config     *config.Config
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Check for nil config
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
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

	return &Server{
		httpServer: httpServer,
		storage:    storageInstance,
		config:     cfg,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	logger.Infof("Starting TimeSeriesDB on port %s...", s.config.Server.Port)
	logger.Infof("Configuration: %s", s.config.String())
	logger.Infof("Metrics available at: http://localhost:%s/metrics", s.config.Server.Port)

	return s.httpServer.ListenAndServe()
}

// Close closes the server and cleans up resources
func (s *Server) Close() error {
	if s.storage != nil {
		s.storage.Close()
	}
	return nil
}
