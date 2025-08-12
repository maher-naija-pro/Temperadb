package main

import (
	"log"
	"net/http"
	"os"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/storage"

	"github.com/joho/godotenv"
)

var (
	storageInstance *storage.Storage
)

func main() {
	// Load .env config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.tsv"
	}

	// Initialize logger
	logger.Init()

	// Initialize storage
	storageInstance = storage.NewStorage(dataFile)
	defer storageInstance.Close()

	// Initialize API router
	router := aphttp.NewRouter(storageInstance)
	router.RegisterRoutes()

	logger.Infof("Starting TimeSeriesDB on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
