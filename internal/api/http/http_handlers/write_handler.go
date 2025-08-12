package handlers

import (
	"fmt"
	"net/http"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/parser"
	"timeseriesdb/internal/storage"
)

// WriteHandler handles the /write endpoint for InfluxDB line protocol
type WriteHandler struct {
	BaseHandler
	storage *storage.Storage
}

// NewWriteHandler creates a new write handler instance
func NewWriteHandler(storage *storage.Storage) *WriteHandler {
	return &WriteHandler{
		storage: storage,
	}
}

// Handle processes write requests
func (h *WriteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.MethodNotAllowed(w, http.MethodPost)
		return
	}

	defer r.Body.Close()

	// Handle negative or zero ContentLength
	if r.ContentLength <= 0 {
		h.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	lines := make([]byte, r.ContentLength)
	_, err := r.Body.Read(lines)
	if err != nil {
		logger.Errorf("Failed to read request body: %v", err)
		h.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	points, err := parser.ParseLineProtocol(string(lines))
	if err != nil {
		logger.Errorf("Failed to parse line protocol: %v", err)
		h.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	// Write points to storage
	successCount := 0
	for _, p := range points {
		err := h.storage.WritePoint(p)
		if err != nil {
			logger.Errorf("Failed to write point: %v", err)
		} else {
			successCount++
		}
	}

	logger.Infof("Wrote %d points successfully", successCount)
	fmt.Fprint(w, "OK")
}
