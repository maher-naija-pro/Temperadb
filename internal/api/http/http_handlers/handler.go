package handlers

import (
	"net/http"
)

// Handler defines the interface for all API handlers
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	// Add common fields here if needed in the future
}

// Common response methods
func (h *BaseHandler) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// In a real implementation, you'd use json.Marshal here
}

func (h *BaseHandler) WriteError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func (h *BaseHandler) MethodNotAllowed(w http.ResponseWriter, allowedMethods ...string) {
	w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
