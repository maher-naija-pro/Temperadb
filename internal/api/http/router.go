package http

import (
	"net/http"
	handlers "timeseriesdb/internal/api/http/http_handlers"
	"timeseriesdb/internal/storage"
)

// Router manages all API routes and handlers
type Router struct {
	writeHandler  *handlers.WriteHandler
	healthHandler *handlers.HealthHandler
}

// NewRouter creates a new router instance with all handlers
func NewRouter(storage *storage.Storage) *Router {
	return &Router{
		writeHandler:  handlers.NewWriteHandler(storage),
		healthHandler: handlers.NewHealthHandler(),
	}
}

// RegisterRoutes registers all API routes with the default HTTP mux
func (r *Router) RegisterRoutes() {
	http.HandleFunc("/write", r.writeHandler.Handle)
	http.HandleFunc("/health", r.healthHandler.Handle)
}

// GetMux returns a new HTTP mux with all routes registered
func (r *Router) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/write", r.writeHandler.Handle)
	mux.HandleFunc("/health", r.healthHandler.Handle)
	return mux
}
