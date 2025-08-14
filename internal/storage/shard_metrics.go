package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Shard metrics
	ShardWriteOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_write_operations_total",
			Help: "Total number of shard write operations",
		},
		[]string{"shard_id", "status"},
	)

	ShardReadOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_read_operations_total",
			Help: "Total number of shard read operations",
		},
		[]string{"shard_id", "status"},
	)

	ShardDataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_data_points_written_total",
			Help: "Total number of data points written to shard",
		},
		[]string{"shard_id"},
	)

	ShardDataPointsRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_data_points_read_total",
			Help: "Total number of data points read from shard",
		},
		[]string{"shard_id"},
	)

	ShardWriteLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_shard_write_latency_seconds",
			Help:    "Shard write operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	ShardReadLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_shard_read_latency_seconds",
			Help:    "Shard read operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	ShardWriteErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_write_errors_total",
			Help: "Total number of shard write errors",
		},
		[]string{"shard_id", "error_type"},
	)

	ShardReadErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_read_errors_total",
			Help: "Total number of shard read errors",
		},
		[]string{"shard_id", "error_type"},
	)

	ShardRecoveryOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_recovery_operations_total",
			Help: "Total number of shard recovery operations",
		},
		[]string{"shard_id", "status"},
	)

	ShardRecoveryLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_shard_recovery_latency_seconds",
			Help:    "Shard recovery operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"shard_id"},
	)

	ShardWALEntriesRecovered = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_shard_wal_entries_recovered_total",
			Help: "Total number of WAL entries recovered during shard recovery",
		},
		[]string{"shard_id"},
	)

	ShardSegmentCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_shard_segment_count",
			Help: "Current number of segments in the shard",
		},
		[]string{"shard_id"},
	)

	ShardTotalSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_shard_total_size_bytes",
			Help: "Total size of all data in the shard in bytes",
		},
		[]string{"shard_id"},
	)

	ShardStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_shard_status",
			Help: "Current status of the shard (0=closed, 1=open, 2=recovering)",
		},
		[]string{"shard_id"},
	)
)
