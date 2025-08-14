package middleware

import (
	"net/http"
	"strconv"
	"time"

	"timeseriesdb/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware wraps HTTP handlers to collect Prometheus metrics
type MetricsMiddleware struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		requests: metrics.HTTPRequestsTotal,
		duration: metrics.HTTPRequestDuration,
	}
}

// Wrap wraps an HTTP handler with metrics collection
func (m *MetricsMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Track requests in flight
		metrics.HTTPRequestsInFlight.WithLabelValues().Inc()
		defer metrics.HTTPRequestsInFlight.WithLabelValues().Dec()

		// Create a response writer wrapper to capture status code and response size
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(wrapped.statusCode)

		m.requests.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
		m.duration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		// Record response size if available
		if wrapped.responseSize > 0 {
			metrics.HTTPResponseSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(wrapped.responseSize))
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.responseSize += len(b)
	return rw.ResponseWriter.Write(b)
}

// MetricsHandler returns a handler that serves Prometheus metrics
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(metrics.GetRegistry(), promhttp.HandlerOpts{})
}
