package http

import (
	"net/http"
	handlers "timeseriesdb/internal/api/http/http_handlers"
	"timeseriesdb/internal/api/middleware"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/internal/storage"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router manages all API routes and handlers
type Router struct {
	writeHandler      *handlers.WriteHandler
	healthHandler     *handlers.HealthHandler
	metricsMiddleware *middleware.MetricsMiddleware
}

// NewRouter creates a new router instance with all handlers
func NewRouter(storage *storage.Storage) *Router {
	return &Router{
		writeHandler:      handlers.NewWriteHandler(storage),
		healthHandler:     handlers.NewHealthHandler(),
		metricsMiddleware: middleware.NewMetricsMiddleware(),
	}
}

// RegisterRoutes registers all API routes with the default HTTP mux
func (r *Router) RegisterRoutes() {
	// Wrap handlers with metrics middleware
	http.Handle("/write", r.metricsMiddleware.Wrap(http.HandlerFunc(r.writeHandler.Handle)))
	http.Handle("/health", r.metricsMiddleware.Wrap(http.HandlerFunc(r.healthHandler.Handle)))
	// Expose Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.HandlerFor(metrics.GetRegistry(), promhttp.HandlerOpts{}))
}

// GetMux returns a new HTTP mux with all routes registered
func (r *Router) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	// Wrap handlers with metrics middleware
	mux.Handle("/write", r.metricsMiddleware.Wrap(http.HandlerFunc(r.writeHandler.Handle)))
	mux.Handle("/health", r.metricsMiddleware.Wrap(http.HandlerFunc(r.healthHandler.Handle)))
	// Expose Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.GetRegistry(), promhttp.HandlerOpts{}))
	return mux
}
