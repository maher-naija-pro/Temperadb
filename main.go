package main

import (
	"log"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Create and start server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("Error creating server:", err)
	}
	defer srv.Close()

	// Start the server
	log.Fatal(srv.Start())
}
