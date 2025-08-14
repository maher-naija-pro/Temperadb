package storage

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// StorageMetrics provides a wrapper around Prometheus metrics for the storage system
type StorageMetrics struct{}

// NewStorageMetrics creates a new StorageMetrics instance
func NewStorageMetrics() *StorageMetrics {
	return &StorageMetrics{}
}

var (
	// Main Storage metrics
	StorageShardCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_shard_count",
			Help: "Total number of storage shards",
		},
		[]string{},
	)

	StorageWriteOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_write_operations_total",
			Help: "Total number of storage write operations",
		},
		[]string{"operation", "status"},
	)

	StorageReadOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_read_operations_total",
			Help: "Total number of storage read operations",
		},
		[]string{"operation", "status"},
	)

	StorageDataPointsWritten = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_data_points_written_total",
			Help: "Total number of data points written to storage",
		},
		[]string{"shard_id"},
	)

	StorageDataPointsRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_data_points_read_total",
			Help: "Total number of data points read from storage",
		},
		[]string{"shard_id"},
	)

	StorageWriteLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_storage_write_latency_seconds",
			Help:    "Storage write operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	StorageReadLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_storage_read_latency_seconds",
			Help:    "Storage read operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	StorageWriteErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_write_errors_total",
			Help: "Total number of storage write errors",
		},
		[]string{"operation", "error_type"},
	)

	StorageReadErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_read_errors_total",
			Help: "Total number of storage read errors",
		},
		[]string{"operation", "error_type"},
	)

	StorageCompactionOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_compaction_operations_total",
			Help: "Total number of storage compaction operations",
		},
		[]string{"status"},
	)

	StorageCompactionLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_storage_compaction_latency_seconds",
			Help:    "Storage compaction operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	StorageSeriesCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_series_count",
			Help: "Total number of series across all shards",
		},
		[]string{},
	)

	StorageTotalSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_total_size_bytes",
			Help: "Total size of all data in storage in bytes",
		},
		[]string{},
	)

	StorageShardCreationOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_shard_creation_operations_total",
			Help: "Total number of shard creation operations",
		},
		[]string{"status"},
	)

	// WAL metrics
	StorageWALSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_wal_size_bytes",
			Help: "Current size of WAL files in bytes",
		},
		[]string{},
	)

	StorageWALFileCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_wal_file_count",
			Help: "Current number of WAL files",
		},
		[]string{},
	)

	StorageWALErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_wal_errors_total",
			Help: "Total number of WAL errors",
		},
		[]string{},
	)

	StorageWALCorruptionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_wal_corruption_errors_total",
			Help: "Total number of WAL corruption errors",
		},
		[]string{},
	)

	StorageWALRecoveryOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_wal_recovery_operations_total",
			Help: "Total number of WAL recovery operations",
		},
		[]string{"status"},
	)

	StorageWALEntriesRead = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_wal_entries_read_total",
			Help: "Total number of WAL entries read",
		},
		[]string{},
	)

	// MemTable metrics
	StorageMemTableSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_memtable_size_bytes",
			Help: "Current size of memtable in bytes",
		},
		[]string{},
	)

	StorageMemTableFlushOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tsdb_storage_memtable_flush_operations_total",
			Help: "Total number of memtable flush operations",
		},
		[]string{"status"},
	)

	StorageMemTableFlushLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tsdb_storage_memtable_flush_latency_seconds",
			Help:    "Memtable flush operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{},
	)

	// Segment metrics
	StorageSegmentCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_segment_count",
			Help: "Number of segments in each compaction level",
		},
		[]string{"level"},
	)

	StorageSegmentSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tsdb_storage_segment_size_bytes",
			Help: "Total size of segments in each compaction level in bytes",
		},
		[]string{"level"},
	)
)

// RecordShardCount records the current number of shards
func (m *StorageMetrics) RecordShardCount(count int) {
	StorageShardCount.WithLabelValues().Set(float64(count))
}

// RecordStorageWriteOperation records a storage write operation
func (m *StorageMetrics) RecordStorageWriteOperation(shardID, operation string) {
	StorageWriteOperations.WithLabelValues(operation, "success").Inc()
}

// RecordStorageWriteError records a storage write error
func (m *StorageMetrics) RecordStorageWriteError(shardID, operation string) {
	StorageWriteErrors.WithLabelValues(operation, "error").Inc()
}

// RecordStorageReadOperation records a storage read operation
func (m *StorageMetrics) RecordStorageReadOperation(shardID, operation string) {
	StorageReadOperations.WithLabelValues(operation, "success").Inc()
}

// RecordStorageReadError records a storage read error
func (m *StorageMetrics) RecordStorageReadError(shardID, operation string) {
	StorageReadErrors.WithLabelValues(operation, "error").Inc()
}

// RecordDataPointsWritten records the number of data points written
func (m *StorageMetrics) RecordDataPointsWritten(shardID string, count int) {
	StorageDataPointsWritten.WithLabelValues(shardID).Add(float64(count))
}

// RecordDataPointsRead records the number of data points read
func (m *StorageMetrics) RecordDataPointsRead(shardID string, count int) {
	StorageDataPointsRead.WithLabelValues(shardID).Add(float64(count))
}

// RecordStorageWriteLatency records write operation latency
func (m *StorageMetrics) RecordStorageWriteLatency(shardID, operation string, duration time.Duration) {
	StorageWriteLatency.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordStorageReadLatency records read operation latency
func (m *StorageMetrics) RecordStorageReadLatency(shardID, operation string, duration time.Duration) {
	StorageReadLatency.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordCompactionStart records the start of a compaction operation
func (m *StorageMetrics) RecordCompactionStart() {
	StorageCompactionOperations.WithLabelValues("started").Inc()
}

// RecordCompactionComplete records the completion of a compaction operation
func (m *StorageMetrics) RecordCompactionComplete(startTime time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	StorageCompactionOperations.WithLabelValues(status).Inc()

	if err == nil {
		duration := time.Since(startTime)
		StorageCompactionLatency.WithLabelValues().Observe(duration.Seconds())
	}
}

// RecordCompactionError records a compaction error
func (m *StorageMetrics) RecordCompactionError() {
	StorageCompactionOperations.WithLabelValues("error").Inc()
}

// RecordWALSize records the current WAL size
func (m *StorageMetrics) RecordWALSize(size int64) {
	StorageWALSize.WithLabelValues().Set(float64(size))
}

// RecordWALFileCount records the number of WAL files
func (m *StorageMetrics) RecordWALFileCount(count int) {
	StorageWALFileCount.WithLabelValues().Set(float64(count))
}

// RecordWALError records a WAL error
func (m *StorageMetrics) RecordWALError() {
	StorageWALErrors.WithLabelValues().Inc()
}

// RecordMemTableSize records the current memtable size
func (m *StorageMetrics) RecordMemTableSize(size int64) {
	StorageMemTableSize.WithLabelValues().Set(float64(size))
}

// RecordMemTableFlushStart records the start of a memtable flush
func (m *StorageMetrics) RecordMemTableFlushStart() {
	StorageMemTableFlushOperations.WithLabelValues("started").Inc()
}

// RecordMemTableFlushComplete records the completion of a memtable flush
func (m *StorageMetrics) RecordMemTableFlushComplete(startTime time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	StorageMemTableFlushOperations.WithLabelValues(status).Inc()

	if err == nil {
		duration := time.Since(startTime)
		StorageMemTableFlushLatency.WithLabelValues().Observe(duration.Seconds())
	}
}

// RecordWALCorruptionError records a WAL corruption error
func (m *StorageMetrics) RecordWALCorruptionError() {
	StorageWALCorruptionErrors.WithLabelValues().Inc()
}

// RecordWALRecoveryComplete records the completion of WAL recovery
func (m *StorageMetrics) RecordWALRecoveryComplete(startTime time.Time, totalCount int, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	StorageWALRecoveryOperations.WithLabelValues(status).Inc()

	if err == nil {
		_ = time.Since(startTime) // Duration calculated but not used yet
		// This would need a WAL recovery latency metric to be added
	}
}

// RecordWALEntriesRead records the number of WAL entries read
func (m *StorageMetrics) RecordWALEntriesRead(count int) {
	StorageWALEntriesRead.WithLabelValues().Add(float64(count))
}

// RecordSegmentCount records the number of segments in a compaction level
func (m *StorageMetrics) RecordSegmentCount(level int, count int) {
	StorageSegmentCount.WithLabelValues(strconv.Itoa(level)).Set(float64(count))
}

// RecordSegmentSize records the total size of segments in a compaction level
func (m *StorageMetrics) RecordSegmentSize(level int, size int64) {
	StorageSegmentSize.WithLabelValues(strconv.Itoa(level)).Set(float64(size))
}

// RecordRecoveryStart records the start of a recovery operation
func (m *StorageMetrics) RecordRecoveryStart() {
	StorageWALRecoveryOperations.WithLabelValues("started").Inc()
}

// RecordRecoveryComplete records the completion of a recovery operation
func (m *StorageMetrics) RecordRecoveryComplete(startTime time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	StorageWALRecoveryOperations.WithLabelValues(status).Inc()
}

// RecordWALFileRotationComplete records the completion of a WAL file rotation
func (m *StorageMetrics) RecordWALFileRotationComplete(startTime time.Time, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	WALFileRotations.WithLabelValues(status).Inc()
}

// RecordWALEntriesWritten records the number of WAL entries written
func (m *StorageMetrics) RecordWALEntriesWritten(count int) {
	WALEntriesWritten.WithLabelValues("success").Add(float64(count))
}

// RecordWALEntrySize records the size of a WAL entry
func (m *StorageMetrics) RecordWALEntrySize(size int) {
	WALBytesWritten.WithLabelValues("success").Add(float64(size))
}

// RecordWALWriteLatency records WAL write operation latency
func (m *StorageMetrics) RecordWALWriteLatency(duration time.Duration) {
	WALWriteLatency.WithLabelValues().Observe(duration.Seconds())
}
