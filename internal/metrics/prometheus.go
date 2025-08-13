package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Registry holds the Prometheus metrics registry
	Registry = prometheus.NewRegistry()

	// Ingestion metrics
	IngestedPoints = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_ingestion_points_total",
			Help: "Total number of points ingested",
		},
	)

	IngestedBatches = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_ingestion_batches_total",
			Help: "Total number of batches ingested",
		},
	)

	IngestionLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "tsdb_ingestion_latency_seconds",
			Help:    "Time taken to ingest points",
			Buckets: prometheus.DefBuckets,
		},
	)

	BatchQueueWaitTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "tsdb_batch_queue_wait_seconds",
			Help:    "Time spent waiting in batch queue",
			Buckets: prometheus.DefBuckets,
		},
	)

	WALAppendLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_append_latency_seconds",
			Help:    "Time taken to append to WAL",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Query metrics
	QueryRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_query_requests_total",
			Help: "Total number of query requests",
		},
	)

	QueryLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "tsdb_query_latency_seconds",
			Help:    "Time taken to execute queries",
			Buckets: prometheus.DefBuckets,
		},
	)

	QueryErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_query_errors_total",
			Help: "Total number of query errors",
		},
	)

	// HTTP API metrics
	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Storage metrics
	WALSizeBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tsdb_wal_size_bytes",
			Help: "Current size of WAL in bytes",
		},
	)

	CompactionRuns = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_runs_total",
			Help: "Total number of compaction runs",
		},
	)

	CompactionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "tsdb_compaction_duration_seconds",
			Help:    "Time taken for compaction operations",
			Buckets: prometheus.DefBuckets,
		},
	)

	CompactionErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_errors_total",
			Help: "Total number of compaction errors",
		},
	)

	// Resource usage metrics
	MemoryPoolUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_memory_pool_bytes",
			Help: "Memory pool usage in bytes",
		},
		[]string{"pool_name"},
	)

	ShardCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tsdb_shard_count",
			Help: "Current number of shards",
		},
	)

	// Cluster metrics
	LeaderElectionResults = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_leader_election_results_total",
			Help: "Leader election results",
		},
		[]string{"result"},
	)

	ReplicationLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_replication_lag_seconds",
			Help: "Replication lag in seconds",
		},
		[]string{"shard_id"},
	)

	// Error metrics
	WriteErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_write_errors_total",
			Help: "Total number of write errors",
		},
	)

	WALErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tsdb_wal_errors_total",
			Help: "Total number of WAL errors",
		},
	)
)

// Init initializes the metrics system
func Init() {
	// Reset any existing metrics to avoid duplicate registration
	Reset()

	// Register all metrics with our custom registry
	Registry.MustRegister(
		IngestedPoints,
		IngestedBatches,
		IngestionLatency,
		BatchQueueWaitTime,
		WALAppendLatency,
		QueryRequests,
		QueryLatency,
		QueryErrors,
		HTTPRequests,
		HTTPRequestDuration,
		WALSizeBytes,
		CompactionRuns,
		CompactionDuration,
		CompactionErrors,
		MemoryPoolUsage,
		ShardCount,
		LeaderElectionResults,
		ReplicationLag,
		WriteErrors,
		WALErrors,
	)
}

// Reset resets the metrics system for testing purposes
func Reset() {
	// Unregister all metrics from the registry
	Registry.Unregister(IngestedPoints)
	Registry.Unregister(IngestedBatches)
	Registry.Unregister(IngestionLatency)
	Registry.Unregister(BatchQueueWaitTime)
	Registry.Unregister(WALAppendLatency)
	Registry.Unregister(QueryRequests)
	Registry.Unregister(QueryLatency)
	Registry.Unregister(QueryErrors)
	Registry.Unregister(HTTPRequests)
	Registry.Unregister(HTTPRequestDuration)
	Registry.Unregister(WALSizeBytes)
	Registry.Unregister(CompactionRuns)
	Registry.Unregister(CompactionDuration)
	Registry.Unregister(CompactionErrors)
	Registry.Unregister(MemoryPoolUsage)
	Registry.Unregister(ShardCount)
	Registry.Unregister(LeaderElectionResults)
	Registry.Unregister(ReplicationLag)
	Registry.Unregister(WriteErrors)
	Registry.Unregister(WALErrors)
}

// GetRegistry returns the Prometheus registry
func GetRegistry() *prometheus.Registry {
	return Registry
}
