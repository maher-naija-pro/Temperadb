package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestNewStorageMetrics(t *testing.T) {
	t.Run("create new storage metrics", func(t *testing.T) {
		metrics := NewStorageMetrics()

		if metrics == nil {
			t.Fatal("StorageMetrics should not be nil")
		}

		// Check if all metrics are initialized
		if metrics.WALSizeBytes == nil {
			t.Error("WALSizeBytes should be initialized")
		}
		if metrics.WALErrors == nil {
			t.Error("WALErrors should be initialized")
		}
		if metrics.CompactionRuns == nil {
			t.Error("CompactionRuns should be initialized")
		}
		if metrics.CompactionDuration == nil {
			t.Error("CompactionDuration should be initialized")
		}
		if metrics.CompactionErrors == nil {
			t.Error("CompactionErrors should be initialized")
		}
		if metrics.ShardCount == nil {
			t.Error("ShardCount should be initialized")
		}
	})
}

func TestRecordWALSize(t *testing.T) {
	t.Run("record WAL size", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording different sizes
		testSizes := []int64{0, 1024, 1024 * 1024, 1024 * 1024 * 1024}

		for _, size := range testSizes {
			metrics.RecordWALSize(size)
			// Note: We can't easily verify the metric value without exposing internal state
			// This test mainly ensures the function doesn't panic
		}
	})

	t.Run("record negative WAL size", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording negative size (should not panic)
		metrics.RecordWALSize(-1024)
	})
}

func TestRecordWALError(t *testing.T) {
	t.Run("record WAL error", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Record multiple errors
		for i := 0; i < 5; i++ {
			metrics.RecordWALError()
		}

		// Note: We can't easily verify the metric value without exposing internal state
		// This test mainly ensures the function doesn't panic
	})
}

func TestRecordCompactionStart(t *testing.T) {
	t.Run("record compaction start", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Record multiple compaction starts
		for i := 0; i < 3; i++ {
			metrics.RecordCompactionStart()
		}

		// Note: We can't easily verify the metric value without exposing internal state
		// This test mainly ensures the function doesn't panic
	})
}

func TestRecordCompactionDuration(t *testing.T) {
	t.Run("record compaction duration", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording different durations
		testDurations := []time.Duration{
			0,
			time.Millisecond,
			time.Second,
			time.Minute,
			time.Hour,
		}

		for _, duration := range testDurations {
			metrics.RecordCompactionDuration(duration)
		}
	})

	t.Run("record negative duration", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording negative duration (should not panic)
		negativeDuration := -time.Second
		metrics.RecordCompactionDuration(negativeDuration)
	})
}

func TestRecordCompactionError(t *testing.T) {
	t.Run("record compaction error", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Record multiple compaction errors
		for i := 0; i < 4; i++ {
			metrics.RecordCompactionError()
		}

		// Note: We can't easily verify the metric value without exposing internal state
		// This test mainly ensures the function doesn't panic
	})
}

func TestRecordShardCount(t *testing.T) {
	t.Run("record shard count", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording different shard counts
		testCounts := []int{0, 1, 5, 10, 100}

		for _, count := range testCounts {
			metrics.RecordShardCount(count)
		}
	})

	t.Run("record negative shard count", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test recording negative count (should not panic)
		metrics.RecordShardCount(-5)
	})
}

func TestRecordCompactionComplete(t *testing.T) {
	t.Run("record successful compaction completion", func(t *testing.T) {
		metrics := NewStorageMetrics()

		startTime := time.Now()
		// Simulate some work
		time.Sleep(1 * time.Millisecond)

		// Record successful completion
		metrics.RecordCompactionComplete(startTime, nil)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic
	})

	t.Run("record failed compaction completion", func(t *testing.T) {
		metrics := NewStorageMetrics()

		startTime := time.Now()
		// Simulate some work
		time.Sleep(1 * time.Millisecond)

		// Record failed completion
		err := fmt.Errorf("compaction failed")
		metrics.RecordCompactionComplete(startTime, err)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic
	})

	t.Run("record compaction completion with zero duration", func(t *testing.T) {
		metrics := NewStorageMetrics()

		startTime := time.Now()
		// Record completion immediately (zero duration)
		metrics.RecordCompactionComplete(startTime, nil)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic
	})

	t.Run("record compaction completion with very long duration", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Use a start time from the past to simulate a very long duration
		startTime := time.Now().Add(-24 * time.Hour) // 24 hours ago

		// Record completion
		metrics.RecordCompactionComplete(startTime, nil)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic
	})
}

func TestMetricsConcurrency(t *testing.T) {
	t.Run("concurrent metric recording", func(t *testing.T) {
		metrics := NewStorageMetrics()
		done := make(chan bool)
		numGoroutines := 10
		iterationsPerGoroutine := 100

		// Start multiple goroutines to record metrics concurrently
		for i := 0; i < numGoroutines; i++ {
			go func() {
				for j := 0; j < iterationsPerGoroutine; j++ {
					metrics.RecordWALSize(int64(j))
					metrics.RecordCompactionStart()
					metrics.RecordCompactionDuration(time.Millisecond)
					metrics.RecordShardCount(j % 10)
				}
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic under concurrent access
	})
}

func TestMetricsEdgeCases(t *testing.T) {
	t.Run("record metrics with extreme values", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test extreme values
		metrics.RecordWALSize(0)
		metrics.RecordWALSize(1 << 62) // Large int64 value
		metrics.RecordShardCount(0)
		metrics.RecordShardCount(1 << 30) // Large int32 value
		metrics.RecordCompactionDuration(0)
		metrics.RecordCompactionDuration(time.Duration(1<<62 - 1)) // Large duration

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic with extreme values
	})

	t.Run("record metrics with zero values", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Test zero values
		metrics.RecordWALSize(0)
		metrics.RecordShardCount(0)
		metrics.RecordCompactionDuration(0)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic with zero values
	})
}

func TestMetricsIntegration(t *testing.T) {
	t.Run("complete compaction workflow", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Simulate a complete compaction workflow
		startTime := time.Now()

		// Start compaction
		metrics.RecordCompactionStart()

		// Simulate some work
		time.Sleep(5 * time.Millisecond)

		// Complete successfully
		metrics.RecordCompactionComplete(startTime, nil)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic during a complete workflow
	})

	t.Run("compaction workflow with error", func(t *testing.T) {
		metrics := NewStorageMetrics()

		// Simulate a complete compaction workflow with error
		startTime := time.Now()

		// Start compaction
		metrics.RecordCompactionStart()

		// Simulate some work
		time.Sleep(5 * time.Millisecond)

		// Complete with error
		err := fmt.Errorf("compaction failed due to disk space")
		metrics.RecordCompactionComplete(startTime, err)

		// Note: We can't easily verify the metric values without exposing internal state
		// This test mainly ensures the function doesn't panic during a complete workflow with error
	})
}
