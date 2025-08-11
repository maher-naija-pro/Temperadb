package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	storage *Storage
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
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Initialize storage
	storage = NewStorage(dataFile)
	defer storage.Close()

	// HTTP handler for line protocol writes
	http.HandleFunc("/write", handleWrite)

	logrus.Infof("Starting TimeSeriesDB on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Accepts InfluxDB line protocol via POST body
func handleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	defer r.Body.Close()

	lines := make([]byte, r.ContentLength)
	r.Body.Read(lines)

	points, err := ParseLineProtocol(string(lines))
	if err != nil {
		logrus.Error("Failed to parse line protocol: ", err)
		http.Error(w, "Bad request", 400)
		return
	}

	for _, p := range points {
		err := storage.WritePoint(p)
		if err != nil {
			logrus.Error("Failed to write point: ", err)
		}
	}

	logrus.Infof("Wrote %d points", len(points))
	fmt.Fprint(w, "OK")
}

