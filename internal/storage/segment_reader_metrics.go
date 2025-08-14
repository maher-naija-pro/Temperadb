package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Segment Reader metrics
	SegmentReaderReadOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_read_operations_total",
			Help: "Total number of segment read operations",
		},
		[]string{"operation", "status"},
	)

	SegmentReaderReadLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_segment_reader_read_latency_seconds",
			Help:    "Segment read operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	SegmentReaderReadErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_read_errors_total",
			Help: "Total number of segment read errors",
		},
		[]string{"operation", "error_type"},
	)

	SegmentReaderSegmentsRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_segments_read_total",
			Help: "Total number of segments read",
		},
		[]string{"status"},
	)

	SegmentReaderDataPointsRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_data_points_read_total",
			Help: "Total number of data points read from segments",
		},
		[]string{"status"},
	)

	SegmentReaderSeriesRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_series_read_total",
			Help: "Total number of series read from segments",
		},
		[]string{"status"},
	)

	SegmentReaderBytesRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_bytes_read_total",
			Help: "Total number of bytes read from segments",
		},
		[]string{"status"},
	)

	SegmentReaderTimeRangeQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_time_range_queries_total",
			Help: "Total number of time range queries on segments",
		},
		[]string{"status"},
	)

	SegmentReaderCorruptedSegments = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_corrupted_segments_total",
			Help: "Total number of corrupted segments encountered",
		},
		[]string{},
	)

	SegmentReaderListOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_segment_reader_list_operations_total",
			Help: "Total number of segment listing operations",
		},
		[]string{"status"},
	)
)
