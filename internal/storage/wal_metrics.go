package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// WAL metrics
	WALWriteOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_write_operations_total",
			Help: "Total number of WAL write operations",
		},
		[]string{"status"},
	)

	WALWriteLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_write_latency_seconds",
			Help:    "WAL write operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	WALWriteErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_write_errors_total",
			Help: "Total number of WAL write errors",
		},
		[]string{"error_type"},
	)

	WALEntriesWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_entries_written_total",
			Help: "Total number of WAL entries written",
		},
		[]string{"status"},
	)

	WALBytesWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_bytes_written_total",
			Help: "Total number of bytes written to WAL",
		},
		[]string{"status"},
	)

	WALFileSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_wal_file_size_bytes",
			Help: "Current size of the WAL file in bytes",
		},
		[]string{},
	)

	WALFileRotations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_file_rotations_total",
			Help: "Total number of WAL file rotations",
		},
		[]string{"status"},
	)

	WALFlushOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_flush_operations_total",
			Help: "Total number of WAL flush operations",
		},
		[]string{"status"},
	)

	WALFlushLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_wal_flush_latency_seconds",
			Help:    "WAL flush operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	WALSequenceNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_wal_sequence_number",
			Help: "Current WAL sequence number",
		},
		[]string{},
	)

	WALDataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_data_points_written_total",
			Help: "Total number of data points written to WAL",
		},
		[]string{"status"},
	)

	WALSeriesWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_wal_series_written_total",
			Help: "Total number of series written to WAL",
		},
		[]string{"status"},
	)
)
