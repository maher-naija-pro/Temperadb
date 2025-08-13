package storage

import (
	"time"

	"timeseriesdb/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

// StorageMetrics wraps all storage-related metrics
type StorageMetrics struct {
	// WAL metrics
	WALSizeBytes prometheus.Gauge
	WALErrors    prometheus.Counter

	// Compaction metrics
	CompactionRuns     prometheus.Counter
	CompactionDuration prometheus.Histogram
	CompactionErrors   prometheus.Counter

	// Resource metrics
	ShardCount prometheus.Gauge
}

// NewStorageMetrics creates a new StorageMetrics instance
func NewStorageMetrics() *StorageMetrics {
	return &StorageMetrics{
		WALSizeBytes:       metrics.WALSizeBytes,
		WALErrors:          metrics.WALErrors,
		CompactionRuns:     metrics.CompactionRuns,
		CompactionDuration: metrics.CompactionDuration,
		CompactionErrors:   metrics.CompactionErrors,
		ShardCount:         metrics.ShardCount,
	}
}

// RecordWALSize records the current WAL size in bytes
func (m *StorageMetrics) RecordWALSize(sizeBytes int64) {
	m.WALSizeBytes.Set(float64(sizeBytes))
}

// RecordWALError records a WAL error
func (m *StorageMetrics) RecordWALError() {
	m.WALErrors.Inc()
}

// RecordCompactionStart records the start of a compaction operation
func (m *StorageMetrics) RecordCompactionStart() {
	m.CompactionRuns.Inc()
}

// RecordCompactionDuration records the duration of a compaction operation
func (m *StorageMetrics) RecordCompactionDuration(duration time.Duration) {
	m.CompactionDuration.Observe(duration.Seconds())
}

// RecordCompactionError records a compaction error
func (m *StorageMetrics) RecordCompactionError() {
	m.CompactionErrors.Inc()
}

// RecordShardCount records the current number of shards
func (m *StorageMetrics) RecordShardCount(count int) {
	m.ShardCount.Set(float64(count))
}

// RecordCompactionComplete records a completed compaction with timing
func (m *StorageMetrics) RecordCompactionComplete(startTime time.Time, err error) {
	duration := time.Since(startTime)
	m.RecordCompactionDuration(duration)

	if err != nil {
		m.RecordCompactionError()
	}
}
