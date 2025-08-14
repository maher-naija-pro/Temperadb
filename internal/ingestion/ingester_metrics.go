package ingestion

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics wraps all ingestion-related metrics
type Metrics struct {
	// Ingestion counters
	IngestedPoints  prometheus.Counter
	IngestedBatches prometheus.Counter
	WriteErrors     prometheus.Counter

	// Latency histograms
	IngestionLatency   prometheus.Histogram
	BatchQueueWaitTime prometheus.Histogram
	WALAppendLatency   prometheus.Histogram
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		IngestedPoints: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ingested_points_total",
			Help: "Total number of ingested data points",
		}),
		IngestedBatches: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "ingested_batches_total",
			Help: "Total number of ingested batches",
		}),
		WriteErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "write_errors_total",
			Help: "Total number of write errors",
		}),
		IngestionLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "ingestion_latency_seconds",
			Help:    "Time spent ingesting data",
			Buckets: prometheus.DefBuckets,
		}),
		BatchQueueWaitTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "batch_queue_wait_time_seconds",
			Help:    "Time spent waiting in batch queue",
			Buckets: prometheus.DefBuckets,
		}),
		WALAppendLatency: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "wal_append_latency_seconds",
			Help:    "Time spent appending to WAL",
			Buckets: prometheus.DefBuckets,
		}),
	}
}

// RecordIngestion records metrics for a successful ingestion
func (m *Metrics) RecordIngestion(points int, duration time.Duration) {
	m.IngestedPoints.Add(float64(points))
	m.IngestionLatency.Observe(duration.Seconds())
}

// RecordBatchIngestion records metrics for batch ingestion
func (m *Metrics) RecordBatchIngestion(batchSize int, queueWaitTime, ingestionTime time.Duration) {
	m.IngestedBatches.Inc()
	m.IngestedPoints.Add(float64(batchSize))
	m.BatchQueueWaitTime.Observe(queueWaitTime.Seconds())
	m.IngestionLatency.Observe(ingestionTime.Seconds())
}

// RecordWALAppend records WAL append latency
func (m *Metrics) RecordWALAppend(duration time.Duration) {
	m.WALAppendLatency.Observe(duration.Seconds())
}

// RecordWriteError records a write error
func (m *Metrics) RecordWriteError() {
	m.WriteErrors.Inc()
}

// RecordBatchQueueWait records time spent waiting in batch queue
func (m *Metrics) RecordBatchQueueWait(duration time.Duration) {
	m.BatchQueueWaitTime.Observe(duration.Seconds())
}
