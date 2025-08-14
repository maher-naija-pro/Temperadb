package metrics

import (
	"time"
	"timeseriesdb/internal/storage"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Registry holds the Prometheus metrics registry
	Registry = prometheus.NewRegistry()
)

// Init initializes the metrics system
func Init() {
	// Check if metrics are already registered
	metrics, err := Registry.Gather()
	if err == nil && len(metrics) > 0 {
		// Metrics are already registered, don't re-register
		return
	}

	// Register all storage metrics
	registerStorageMetrics()

	// Register all server metrics
	registerServerMetrics()

	// Initialize metrics with default values so they can be collected
	initializeMetricsWithDefaults()
}

// Reset resets the metrics system for testing purposes
func Reset() {
	// Unregister metrics from the default registry first
	unregisterFromDefaultRegistry()

	// Create a new registry but keep the same instance reference
	*Registry = *prometheus.NewRegistry()
}

// unregisterFromDefaultRegistry unregisters all our metrics from the default Prometheus registry
func unregisterFromDefaultRegistry() {
	// List of all storage metrics that might be registered with the default registry
	storageMetrics := []prometheus.Collector{
		storage.CompactionOperations,
		storage.MemStoreSize,
		storage.SegmentReaderReadOperations,
		storage.SegmentWriterWriteOperations,
		storage.ShardWriteOperations,
		storage.StorageShardCount,
		storage.WALWriteOperations,
		storage.WALReplayOperations,
	}

	// Unregister each metric from the default registry
	for _, metric := range storageMetrics {
		if metric != nil {
			prometheus.Unregister(metric)
		}
	}
}

// GetRegistry returns the Prometheus registry
func GetRegistry() *prometheus.Registry {
	return Registry
}

// registerStorageMetrics registers all storage-related metrics
func registerStorageMetrics() {
	// Storage metrics that are defined but may not be registered yet
	storageMetrics := []prometheus.Collector{
		storage.CompactionOperations,
		storage.MemStoreSize,
		storage.SegmentReaderReadOperations,
		storage.SegmentWriterWriteOperations,
		storage.ShardWriteOperations,
		storage.StorageShardCount,
		storage.WALWriteOperations,
		storage.WALReplayOperations,
	}

	for _, metric := range storageMetrics {
		if metric == nil {
			continue
		}

		// Try to register the metric, but don't fail if it's already registered
		err := Registry.Register(metric)
		if err != nil {
			// If the metric is already registered, try to unregister it first
			if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
				// Try to unregister from default registry first
				prometheus.Unregister(metric)
				// Then register with our registry
				err = Registry.Register(metric)
				if err != nil {
					// Log error but don't panic
				}
			} else {
				// Log error but don't panic
			}
		}
	}
}

// registerServerMetrics registers all server-related metrics
func registerServerMetrics() {
	// Manually register server metrics with our custom registry
	metrics := []prometheus.Collector{
		ServerStatus,
		ServerStartTime,
		ServerUptime,
		ServerConfigPort,
		ServerConfigReadTimeout,
		ServerConfigWriteTimeout,
		ServerConfigIdleTimeout,
		ServerMemoryUsage,
		ServerGoroutines,
		ServerShutdownDuration,
		ServerErrors,
		StorageConnectionStatus,
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
		HTTPResponseSize,
		APIVersion,
		BuildInfo,
		ServerActiveConnections,
		ServerHealth,
		DataPointsWritten,
		DataPointsWrittenRate,
	}

	for _, metric := range metrics {
		if metric != nil {
			err := Registry.Register(metric)
			if err != nil {
				// Log error but don't panic
			}
		}
	}
}

// initializeMetricsWithDefaults initializes metrics with default values so they can be collected
func initializeMetricsWithDefaults() {
	// Initialize storage metrics with default labels
	if storage.CompactionOperations != nil {
		storage.CompactionOperations.WithLabelValues("0", "success").Add(0)
	}
	if storage.MemStoreSize != nil {
		storage.MemStoreSize.WithLabelValues("0").Set(0)
	}
	if storage.SegmentReaderReadOperations != nil {
		storage.SegmentReaderReadOperations.WithLabelValues("read", "success").Add(0)
	}
	if storage.SegmentWriterWriteOperations != nil {
		storage.SegmentWriterWriteOperations.WithLabelValues("success").Add(0)
	}
	if storage.ShardWriteOperations != nil {
		storage.ShardWriteOperations.WithLabelValues("0", "success").Add(0)
	}
	if storage.StorageShardCount != nil {
		storage.StorageShardCount.WithLabelValues().Set(0)
	}
	if storage.WALWriteOperations != nil {
		storage.WALWriteOperations.WithLabelValues("success").Add(0)
	}
	if storage.WALReplayOperations != nil {
		storage.WALReplayOperations.WithLabelValues("success").Add(0)
	}

	// Initialize server metrics with default values
	if ServerStatus != nil {
		ServerStatus.WithLabelValues().Set(1)
	}
	if ServerStartTime != nil {
		ServerStartTime.WithLabelValues().Set(float64(time.Now().Unix()))
	}
	if ServerUptime != nil {
		ServerUptime.WithLabelValues().Set(0)
	}
	if ServerConfigPort != nil {
		ServerConfigPort.WithLabelValues().Set(8080)
	}
	if ServerConfigReadTimeout != nil {
		ServerConfigReadTimeout.WithLabelValues().Set(30)
	}
	if ServerConfigWriteTimeout != nil {
		ServerConfigWriteTimeout.WithLabelValues().Set(30)
	}
	if ServerConfigIdleTimeout != nil {
		ServerConfigIdleTimeout.WithLabelValues().Set(60)
	}
	if ServerMemoryUsage != nil {
		ServerMemoryUsage.WithLabelValues().Set(0)
	}
	if ServerGoroutines != nil {
		ServerGoroutines.WithLabelValues().Set(0)
	}
	if ServerShutdownDuration != nil {
		ServerShutdownDuration.WithLabelValues().Observe(0)
	}
	if ServerErrors != nil {
		ServerErrors.WithLabelValues("general", "none").Add(0)
	}
	if StorageConnectionStatus != nil {
		StorageConnectionStatus.WithLabelValues().Set(1)
	}
	if HTTPRequestsTotal != nil {
		HTTPRequestsTotal.WithLabelValues("GET", "/metrics", "200").Add(0)
	}
	if HTTPRequestDuration != nil {
		HTTPRequestDuration.WithLabelValues("GET", "/metrics").Observe(0)
	}
	if HTTPRequestsInFlight != nil {
		HTTPRequestsInFlight.WithLabelValues().Set(0)
	}
	if HTTPResponseSize != nil {
		HTTPResponseSize.WithLabelValues("GET", "/metrics").Observe(0)
	}
	if APIVersion != nil {
		APIVersion.WithLabelValues("v1").Set(1)
	}
	if BuildInfo != nil {
		BuildInfo.WithLabelValues("1.0.0", "unknown", "main", "1.21").Set(1)
	}
	if ServerActiveConnections != nil {
		ServerActiveConnections.WithLabelValues().Set(0)
	}
	if ServerHealth != nil {
		ServerHealth.WithLabelValues().Set(1)
	}
	if DataPointsWritten != nil {
		DataPointsWritten.WithLabelValues("default").Add(0)
	}
	if DataPointsWrittenRate != nil {
		DataPointsWrittenRate.WithLabelValues("default").Set(0)
	}
}

// For direct access to metrics, import the storage package:
// import "timeseriesdb/internal/storage"
// Then use: storage.CompactionOperations, storage.MemStoreSize, etc.
