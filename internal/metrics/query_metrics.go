package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// QueryMetrics wraps all query-related metrics
type QueryMetrics struct {
	// Query counters
	QueryRequests prometheus.Counter
	QueryErrors   prometheus.Counter

	// Query latency
	QueryLatency prometheus.Histogram
}

// NewQueryMetrics creates a new QueryMetrics instance
func NewQueryMetrics() *QueryMetrics {
	return &QueryMetrics{
		QueryRequests: QueryRequests,
		QueryErrors:   QueryErrors,
		QueryLatency:  QueryLatency,
	}
}

// RecordQueryStart records the start of a query
func (m *QueryMetrics) RecordQueryStart() {
	m.QueryRequests.Inc()
}

// RecordQueryDuration records the duration of a query
func (m *QueryMetrics) RecordQueryDuration(duration time.Duration) {
	m.QueryLatency.Observe(duration.Seconds())
}

// RecordQueryError records a query error
func (m *QueryMetrics) RecordQueryError() {
	m.QueryErrors.Inc()
}

// RecordQueryComplete records a completed query with timing
func (m *QueryMetrics) RecordQueryComplete(startTime time.Time, err error) {
	duration := time.Since(startTime)
	m.RecordQueryDuration(duration)

	if err != nil {
		m.RecordQueryError()
	}
}
