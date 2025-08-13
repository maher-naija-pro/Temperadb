package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer handles serving Prometheus metrics
type MetricsServer struct {
	addr     string
	registry *prometheus.Registry
}

// NewMetricsServer creates a new metrics server
func NewMetricsServer(addr string) *MetricsServer {
	return &MetricsServer{
		addr:     addr,
		registry: Registry,
	}
}

// Start starts the metrics server
func (s *MetricsServer) Start() error {
	mux := http.NewServeMux()

	// Register metrics endpoint
	mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))

	// Register health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register readiness endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Printf("Starting metrics server on %s\n", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

// GetMetricsEndpoint returns the metrics endpoint URL
func (s *MetricsServer) GetMetricsEndpoint() string {
	return fmt.Sprintf("http://%s/metrics", s.addr)
}

// GetHealthEndpoint returns the health check endpoint URL
func (s *MetricsServer) GetHealthEndpoint() string {
	return fmt.Sprintf("http://%s/health", s.addr)
}

// GetReadyEndpoint returns the readiness endpoint URL
func (s *MetricsServer) GetReadyEndpoint() string {
	return fmt.Sprintf("http://%s/ready", s.addr)
}
