package ingestion

import (
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	t.Run("NewMetrics", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Verify all metrics were initialized
		if metrics.IngestedPoints == nil {
			t.Error("Expected IngestedPoints to be initialized")
		}

		if metrics.IngestedBatches == nil {
			t.Error("Expected IngestedBatches to be initialized")
		}

		if metrics.WriteErrors == nil {
			t.Error("Expected WriteErrors to be initialized")
		}

		if metrics.IngestionLatency == nil {
			t.Error("Expected IngestionLatency to be initialized")
		}

		if metrics.BatchQueueWaitTime == nil {
			t.Error("Expected BatchQueueWaitTime to be initialized")
		}

		if metrics.WALAppendLatency == nil {
			t.Error("Expected WALAppendLatency to be initialized")
		}
	})
}

func TestRecordIngestion(t *testing.T) {
	t.Run("RecordIngestion", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording ingestion metrics
		points := 100
		duration := 50 * time.Millisecond

		// Record ingestion
		metrics.RecordIngestion(points, duration)

		// Verify metrics were recorded
		// Note: We can't easily verify the actual values without exposing internal metric state
		// but we can verify the operations don't panic
		t.Log("Ingestion metrics recorded successfully")
	})

	t.Run("RecordIngestionWithZeroPoints", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with zero points
		points := 0
		duration := time.Duration(0)

		// Record ingestion
		metrics.RecordIngestion(points, duration)

		t.Log("Ingestion metrics recorded successfully with zero values")
	})

	t.Run("RecordIngestionWithLargeValues", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with large values
		points := 1000000
		duration := 10 * time.Second

		// Record ingestion
		metrics.RecordIngestion(points, duration)

		t.Log("Ingestion metrics recorded successfully with large values")
	})
}

func TestRecordBatchIngestion(t *testing.T) {
	t.Run("RecordBatchIngestion", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording batch ingestion metrics
		batchSize := 500
		queueWaitTime := 25 * time.Millisecond
		ingestionTime := 75 * time.Millisecond

		// Record batch ingestion
		metrics.RecordBatchIngestion(batchSize, queueWaitTime, ingestionTime)

		// Verify metrics were recorded
		t.Log("Batch ingestion metrics recorded successfully")
	})

	t.Run("RecordBatchIngestionWithZeroValues", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with zero values
		batchSize := 0
		queueWaitTime := time.Duration(0)
		ingestionTime := time.Duration(0)

		// Record batch ingestion
		metrics.RecordBatchIngestion(batchSize, queueWaitTime, ingestionTime)

		t.Log("Batch ingestion metrics recorded successfully with zero values")
	})

	t.Run("RecordBatchIngestionWithLargeValues", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with large values
		batchSize := 10000
		queueWaitTime := 5 * time.Second
		ingestionTime := 15 * time.Second

		// Record batch ingestion
		metrics.RecordBatchIngestion(batchSize, queueWaitTime, ingestionTime)

		t.Log("Batch ingestion metrics recorded successfully with large values")
	})
}

func TestRecordWALAppend(t *testing.T) {
	t.Run("RecordWALAppend", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording WAL append metrics
		duration := 30 * time.Millisecond

		// Record WAL append
		metrics.RecordWALAppend(duration)

		// Verify metrics were recorded
		t.Log("WAL append metrics recorded successfully")
	})

	t.Run("RecordWALAppendWithZeroDuration", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with zero duration
		duration := time.Duration(0)

		// Record WAL append
		metrics.RecordWALAppend(duration)

		t.Log("WAL append metrics recorded successfully with zero duration")
	})

	t.Run("RecordWALAppendWithLongDuration", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with long duration
		duration := 2 * time.Second

		// Record WAL append
		metrics.RecordWALAppend(duration)

		t.Log("WAL append metrics recorded successfully with long duration")
	})
}

func TestRecordWriteError(t *testing.T) {
	t.Run("RecordWriteError", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording write error
		metrics.RecordWriteError()

		// Verify metrics were recorded
		t.Log("Write error metrics recorded successfully")
	})

	t.Run("RecordMultipleWriteErrors", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording multiple write errors
		for i := 0; i < 5; i++ {
			metrics.RecordWriteError()
		}

		t.Log("Multiple write error metrics recorded successfully")
	})
}

func TestRecordBatchQueueWait(t *testing.T) {
	t.Run("RecordBatchQueueWait", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test recording batch queue wait metrics
		duration := 45 * time.Millisecond

		// Record batch queue wait
		metrics.RecordBatchQueueWait(duration)

		// Verify metrics were recorded
		t.Log("Batch queue wait metrics recorded successfully")
	})

	t.Run("RecordBatchQueueWaitWithZeroDuration", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with zero duration
		duration := time.Duration(0)

		// Record batch queue wait
		metrics.RecordBatchQueueWait(duration)

		t.Log("Batch queue wait metrics recorded successfully with zero duration")
	})

	t.Run("RecordBatchQueueWaitWithLongDuration", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test with long duration
		duration := 3 * time.Second

		// Record batch queue wait
		metrics.RecordBatchQueueWait(duration)

		t.Log("Batch queue wait metrics recorded successfully with long duration")
	})
}

func TestMetricsIntegration(t *testing.T) {
	t.Run("MetricsIntegration", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test all metrics functions together
		startTime := time.Now()

		// Record various metrics
		metrics.RecordIngestion(100, 50*time.Millisecond)
		metrics.RecordBatchIngestion(500, 25*time.Millisecond, 75*time.Millisecond)
		metrics.RecordWALAppend(30 * time.Millisecond)
		metrics.RecordWriteError()
		metrics.RecordBatchQueueWait(45 * time.Millisecond)

		// Verify all operations completed without error
		totalTime := time.Since(startTime)
		t.Logf("All metrics operations completed successfully in %v", totalTime)
	})
}

func TestMetricsEdgeCases(t *testing.T) {
	t.Run("MetricsEdgeCases", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test edge cases
		// Very small durations
		metrics.RecordIngestion(1, time.Nanosecond)
		metrics.RecordWALAppend(time.Nanosecond)
		metrics.RecordBatchQueueWait(time.Nanosecond)

		// Very large durations
		metrics.RecordIngestion(1, time.Hour)
		metrics.RecordWALAppend(time.Hour)
		metrics.RecordBatchQueueWait(time.Hour)

		// Large point counts
		metrics.RecordIngestion(1000000, time.Millisecond)
		metrics.RecordBatchIngestion(1000000, time.Millisecond, time.Millisecond)

		t.Log("Edge case metrics recorded successfully")
	})
}

func TestMetricsConcurrency(t *testing.T) {
	t.Run("MetricsConcurrency", func(t *testing.T) {
		metrics := NewMetrics()
		if metrics == nil {
			t.Fatal("Expected metrics to be created")
		}

		// Test concurrent access to metrics
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				defer func() { done <- true }()

				// Record various metrics concurrently
				metrics.RecordIngestion(id*10, time.Duration(id)*time.Millisecond)
				metrics.RecordBatchIngestion(id*5, time.Duration(id)*time.Millisecond, time.Duration(id)*time.Millisecond)
				metrics.RecordWALAppend(time.Duration(id) * time.Millisecond)
				metrics.RecordWriteError()
				metrics.RecordBatchQueueWait(time.Duration(id) * time.Millisecond)
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		t.Log("Concurrent metrics operations completed successfully")
	})
}
