package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Segment Writer metrics
	SegmentWriterWriteOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_write_operations_total",
			Help: "Total number of segment write operations",
		},
		[]string{"status"},
	)

	SegmentWriterWriteLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_segment_writer_write_latency_seconds",
			Help:    "Segment write operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	SegmentWriterWriteErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_write_errors_total",
			Help: "Total number of segment write errors",
		},
		[]string{"error_type"},
	)

	SegmentWriterSegmentsCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_segments_created_total",
			Help: "Total number of segments created",
		},
		[]string{"status"},
	)

	SegmentWriterDataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_data_points_written_total",
			Help: "Total number of data points written to segments",
		},
		[]string{"status"},
	)

	SegmentWriterSeriesWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_series_written_total",
			Help: "Total number of series written to segments",
		},
		[]string{"status"},
	)

	SegmentWriterBytesWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_bytes_written_total",
			Help: "Total number of bytes written to segments",
		},
		[]string{"status"},
	)

	SegmentWriterSegmentSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_segment_writer_segment_size_bytes",
			Help:    "Size of created segments in bytes",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 20), // 1KB to 1GB
		},
		[]string{},
	)

	SegmentWriterMemTableFlushes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_writer_memtable_flushes_total",
			Help: "Total number of memtable flushes to segments",
		},
		[]string{"status"},
	)

	SegmentWriterCompressionRatio = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_segment_writer_compression_ratio",
			Help:    "Compression ratio achieved when writing segments",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10), // 0.1 to 1.0
		},
		[]string{},
	)
)
