package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// WAL Replay metrics
	WALReplayOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_operations_total",
			Help: "Total number of WAL replay operations",
		},
		[]string{"status"},
	)

	WALReplayLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_replay_latency_seconds",
			Help:    "WAL replay operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	WALReplayErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_errors_total",
			Help: "Total number of WAL replay errors",
		},
		[]string{"error_type"},
	)

	WALReplayFilesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_files_processed_total",
			Help: "Total number of WAL files processed during replay",
		},
		[]string{"status"},
	)

	WALReplayEntriesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_entries_processed_total",
			Help: "Total number of WAL entries processed during replay",
		},
		[]string{"status"},
	)

	WALReplayDataPointsRecovered = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_data_points_recovered_total",
			Help: "Total number of data points recovered during WAL replay",
		},
		[]string{"status"},
	)

	WALReplaySeriesRecovered = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_series_recovered_total",
			Help: "Total number of series recovered during WAL replay",
		},
		[]string{"status"},
	)

	WALReplayBytesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_bytes_processed_total",
			Help: "Total number of bytes processed during WAL replay",
		},
		[]string{"status"},
	)

	WALReplayValidationErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_validation_errors_total",
			Help: "Total number of WAL entry validation errors during replay",
		},
		[]string{"error_type"},
	)

	WALReplayCleanupOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_cleanup_operations_total",
			Help: "Total number of WAL cleanup operations",
		},
		[]string{"status"},
	)

	WALReplayCleanupLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_replay_cleanup_latency_seconds",
			Help:    "WAL cleanup operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	WALReplayOldFilesRemoved = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_replay_old_files_removed_total",
			Help: "Total number of old WAL files removed during cleanup",
		},
		[]string{},
	)

	WALReplayRecoveryTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_replay_recovery_time_seconds",
			Help:    "Total time taken for WAL recovery in seconds",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 15), // 0.1s to 1638.4s
		},
		[]string{},
	)
)
