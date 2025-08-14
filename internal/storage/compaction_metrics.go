package storage

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Compaction metrics
	CompactionOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_operations_total",
			Help: "Total number of compaction operations",
		},
		[]string{"level", "status"},
	)

	CompactionLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_compaction_latency_seconds",
			Help:    "Compaction operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"level"},
	)

	CompactionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_errors_total",
			Help: "Total number of compaction errors",
		},
		[]string{"level", "error_type"},
	)

	CompactionSegmentsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_segments_processed_total",
			Help: "Total number of segments processed during compaction",
		},
		[]string{"level"},
	)

	CompactionDataPointsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_data_points_processed_total",
			Help: "Total number of data points processed during compaction",
		},
		[]string{"level"},
	)

	CompactionLevelSegments = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_compaction_level_segments",
			Help: "Number of segments in each compaction level",
		},
		[]string{"level"},
	)

	CompactionLevelSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_compaction_level_size_bytes",
			Help: "Total size of segments in each compaction level in bytes",
		},
		[]string{"level"},
	)

	CompactionQueueSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_compaction_queue_size",
			Help: "Current size of the compaction task queue",
		},
		[]string{},
	)

	CompactionPromotions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_compaction_segment_promotions_total",
			Help: "Total number of segments promoted to higher levels",
		},
		[]string{"from_level", "to_level"},
	)
)

func init() {
	// Register all compaction metrics
	prometheus.MustRegister(CompactionOperations)
	prometheus.MustRegister(CompactionLatency)
	prometheus.MustRegister(CompactionErrors)
	prometheus.MustRegister(CompactionSegmentsProcessed)
	prometheus.MustRegister(CompactionDataPointsProcessed)
	prometheus.MustRegister(CompactionLevelSegments)
	prometheus.MustRegister(CompactionLevelSize)
	prometheus.MustRegister(CompactionQueueSize)
	prometheus.MustRegister(CompactionPromotions)
}
