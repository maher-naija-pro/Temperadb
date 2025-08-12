package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler handles the /health endpoint for health checks
type HealthHandler struct {
	BaseHandler
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Handle processes health check requests
func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.MethodNotAllowed(w, http.MethodGet)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"healthy","service":"TimeSeriesDB"}`)
}
