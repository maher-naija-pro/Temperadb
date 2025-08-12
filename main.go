package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/parser"
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

	// HTTP handler for line protocol writes
	http.HandleFunc("/write", handleWrite)

	logger.Infof("Starting TimeSeriesDB on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Accepts InfluxDB line protocol via POST body
func handleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	defer r.Body.Close()

	// Handle negative or zero ContentLength
	if r.ContentLength <= 0 {
		http.Error(w, "Bad request", 400)
		return
	}

	lines := make([]byte, r.ContentLength)
	r.Body.Read(lines)

	points, err := parser.ParseLineProtocol(string(lines))
	if err != nil {
		logger.Errorf("Failed to parse line protocol: %v", err)
		http.Error(w, "Bad request", 400)
		return
	}

	for _, p := range points {
		err := storageInstance.WritePoint(p)
		if err != nil {
			logger.Errorf("Failed to write point: %v", err)
		}
	}

	logger.Infof("Wrote %d points", len(points))
	fmt.Fprint(w, "OK")
}
