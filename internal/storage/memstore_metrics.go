package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// MemStore metrics
	MemStoreSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_memstore_size_bytes",
			Help: "Current size of the memstore in bytes",
		},
		[]string{"shard_id"},
	)

	MemStoreWriteOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_write_operations_total",
			Help: "Total number of memstore write operations",
		},
		[]string{"shard_id", "status"},
	)

	MemStoreReadOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_read_operations_total",
			Help: "Total number of memstore read operations",
		},
		[]string{"shard_id", "status"},
	)

	MemStoreDataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_data_points_written_total",
			Help: "Total number of data points written to memstore",
		},
		[]string{"shard_id"},
	)

	MemStoreDataPointsRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_data_points_read_total",
			Help: "Total number of data points read from memstore",
		},
		[]string{"shard_id"},
	)

	MemStoreWriteLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_memstore_write_latency_seconds",
			Help:    "Memstore write operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	MemStoreReadLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_memstore_read_latency_seconds",
			Help:    "Memstore read operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	MemStoreFlushOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_flush_operations_total",
			Help: "Total number of memstore flush operations",
		},
		[]string{"shard_id", "status"},
	)

	MemStoreFlushLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_memstore_flush_latency_seconds",
			Help:    "Memstore flush operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	MemStoreSeriesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_memstore_series_count",
			Help: "Current number of series in the memstore",
		},
		[]string{"shard_id"},
	)

	MemStoreWALErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_memstore_wal_errors_total",
			Help: "Total number of WAL write errors in memstore",
		},
		[]string{"shard_id"},
	)
)
