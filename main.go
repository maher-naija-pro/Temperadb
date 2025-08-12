package main

import (
	"log"
	"net/http"
	"timeseriesdb/config"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/storage"
)

var (
	storageInstance *storage.Storage
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Initialize logger with configuration
	logger.InitWithConfig(cfg.Logging)

	// Initialize storage with configuration
	storageInstance = storage.NewStorage(cfg.Storage)
	defer storageInstance.Close()

	// Initialize API router
	router := aphttp.NewRouter(storageInstance)
	router.RegisterRoutes()

	// Create HTTP server with configuration
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.Infof("Starting TimeSeriesDB on port %s...", cfg.Server.Port)
	logger.Infof("Configuration: %s", cfg.String())
	log.Fatal(server.ListenAndServe())
}
