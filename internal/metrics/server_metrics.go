package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ServerStatus represents the current server status
	// 0=stopped, 1=starting, 2=running, 3=shutting_down, 4=stopped
	ServerStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_status",
			Help: "Current server status (0=stopped, 1=starting, 2=running, 3=shutting_down, 4=stopped)",
		},
		[]string{},
	)

	// ServerStartTime represents the server start timestamp
	ServerStartTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_start_time_seconds",
			Help: "Server start time in Unix timestamp",
		},
		[]string{},
	)

	// ServerUptime represents the server uptime in seconds
	ServerUptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_uptime_seconds",
			Help: "Server uptime in seconds",
		},
		[]string{},
	)

	// ServerConfigPort represents the configured server port
	ServerConfigPort = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_config_port",
			Help: "Configured server port",
		},
		[]string{},
	)

	// ServerConfigReadTimeout represents the configured read timeout
	ServerConfigReadTimeout = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_config_read_timeout_seconds",
			Help: "Configured server read timeout in seconds",
		},
		[]string{},
	)

	// ServerConfigWriteTimeout represents the configured write timeout
	ServerConfigWriteTimeout = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_config_write_timeout_seconds",
			Help: "Configured server write timeout in seconds",
		},
		[]string{},
	)

	// ServerConfigIdleTimeout represents the configured idle timeout
	ServerConfigIdleTimeout = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_config_idle_timeout_seconds",
			Help: "Configured server idle timeout in seconds",
		},
		[]string{},
	)

	// ServerMemoryUsage represents the current memory usage in bytes
	ServerMemoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_memory_usage_bytes",
			Help: "Current server memory usage in bytes",
		},
		[]string{},
	)

	// ServerGoroutines represents the current number of goroutines
	ServerGoroutines = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_goroutines",
			Help: "Current number of goroutines",
		},
		[]string{},
	)

	// ServerShutdownDuration represents the server shutdown duration
	ServerShutdownDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_server_shutdown_duration_seconds",
			Help:    "Server shutdown duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	// ServerErrors represents server errors by type and component
	ServerErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_server_errors_total",
			Help: "Total number of server errors",
		},
		[]string{"error_type", "component"},
	)

	// StorageConnectionStatus represents the storage connection status
	// 0=disconnected, 1=connected
	StorageConnectionStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_connection_status",
			Help: "Storage connection status (0=disconnected, 1=connected)",
		},
		[]string{},
	)

	// HTTPRequestsTotal represents total HTTP requests
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// HTTPRequestDuration represents HTTP request duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// HTTPRequestsInFlight represents current HTTP requests in flight
	HTTPRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
		[]string{},
	)

	// HTTPResponseSize represents HTTP response size in bytes
	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B to 100MB
		},
		[]string{"method", "endpoint"},
	)

	// APIVersion represents the API version being served
	APIVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_api_version",
			Help: "API version being served",
		},
		[]string{"version"},
	)

	// BuildInfo represents build information
	BuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_build_info",
			Help: "Build information",
		},
		[]string{"version", "commit", "branch", "go_version"},
	)

	// ServerActiveConnections represents the current number of active connections
	ServerActiveConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_active_connections",
			Help: "Current number of active connections",
		},
		[]string{},
	)

	// ServerHealth represents the overall server health status
	// 0=unhealthy, 1=healthy
	ServerHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_server_health",
			Help: "Server health status (0=unhealthy, 1=healthy)",
		},
		[]string{},
	)

	// DataPointsWritten represents the total number of data points written
	DataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_data_points_written_total",
			Help: "Total number of data points written to storage",
		},
		[]string{"measurement"},
	)

	// DataPointsWrittenRate represents the rate of data points written per second
	DataPointsWrittenRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_data_points_written_rate",
			Help: "Rate of data points written per second",
		},
		[]string{"measurement"},
	)
)
