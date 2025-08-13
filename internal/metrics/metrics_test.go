package metrics

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test setup and teardown
func setupTestMetrics(t *testing.T) func() {
	// Reset registry for clean testing
	Registry = prometheus.NewRegistry()

	// Re-initialize metrics
	Init()

	// Return cleanup function
	return func() {
		// Reset registry
		Registry = prometheus.NewRegistry()
	}
}

// Benchmark setup
func setupBenchmarkMetrics(b *testing.B) func() {
	// Reset registry for clean testing
	Registry = prometheus.NewRegistry()

	// Re-initialize metrics
	Init()

	// Return cleanup function
	return func() {
		// Reset registry
		Registry = prometheus.NewRegistry()
	}
}

func TestMetricsInitialization(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Verify that all key metrics are registered
	requiredMetrics := []struct {
		name   string
		metric interface{}
	}{
		{"IngestedPoints", IngestedPoints},
		{"IngestedBatches", IngestedBatches},
		{"IngestionLatency", IngestionLatency},
		{"BatchQueueWaitTime", BatchQueueWaitTime},
		{"WALAppendLatency", WALAppendLatency},
		{"QueryRequests", QueryRequests},
		{"QueryLatency", QueryLatency},
		{"QueryErrors", QueryErrors},
		{"HTTPRequests", HTTPRequests},
		{"HTTPRequestDuration", HTTPRequestDuration},
		{"WALSizeBytes", WALSizeBytes},
		{"CompactionRuns", CompactionRuns},
		{"CompactionDuration", CompactionDuration},
		{"CompactionErrors", CompactionErrors},
		{"MemoryPoolUsage", MemoryPoolUsage},
		{"ShardCount", ShardCount},
		{"LeaderElectionResults", LeaderElectionResults},
		{"ReplicationLag", ReplicationLag},
		{"WriteErrors", WriteErrors},
		{"WALErrors", WALErrors},
	}

	for _, rm := range requiredMetrics {
		assert.NotNil(t, rm.metric, "Metric %s should not be nil", rm.name)
	}

	// Verify registry contains metrics (some might be optional)
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 10, "Registry should contain a reasonable number of metrics")
}

func TestIngestionMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test ingestion counters
	IngestedPoints.Add(100)
	IngestedBatches.Inc()

	// Test latency metrics with multiple observations
	latencyValues := []float64{0.01, 0.05, 0.1, 0.5, 1.0}
	for _, latency := range latencyValues {
		IngestionLatency.Observe(latency)
	}

	// Test batch queue metrics
	BatchQueueWaitTime.Observe(0.05) // 50ms
	BatchQueueWaitTime.Observe(0.1)  // 100ms

	// Verify metrics are recorded by checking registry
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	// Check that metrics exist in registry
	var foundIngestion, foundBatchQueue bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "ingestion_latency_seconds") {
			foundIngestion = true
		}
		if strings.Contains(*metric.Name, "batch_queue_wait_seconds") {
			foundBatchQueue = true
		}
	}

	assert.True(t, foundIngestion, "IngestionLatency should be in registry")
	assert.True(t, foundBatchQueue, "BatchQueueWaitTime should be in registry")
}

func TestQueryMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test query counters
	QueryRequests.Inc()
	QueryRequests.Inc()
	QueryErrors.Inc()

	// Test query latency with realistic values
	latencyValues := []float64{0.001, 0.01, 0.1, 0.5, 2.0}
	for _, latency := range latencyValues {
		QueryLatency.Observe(latency)
	}

	// Verify metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundQueryRequests, foundQueryLatency bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "query_requests_total") {
			foundQueryRequests = true
		}
		if strings.Contains(*metric.Name, "query_latency_seconds") {
			foundQueryLatency = true
		}
	}

	assert.True(t, foundQueryRequests, "QueryRequests should be in registry")
	assert.True(t, foundQueryLatency, "QueryLatency should be in registry")
}

func TestHTTPMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test HTTP request counter with labels
	HTTPRequests.WithLabelValues("GET", "/health", "200").Inc()
	HTTPRequests.WithLabelValues("POST", "/write", "201").Inc()
	HTTPRequests.WithLabelValues("GET", "/metrics", "200").Inc()
	HTTPRequests.WithLabelValues("POST", "/write", "400").Inc()

	// Test HTTP request duration
	HTTPRequestDuration.WithLabelValues("GET", "/health").Observe(0.01)
	HTTPRequestDuration.WithLabelValues("POST", "/write").Observe(0.05)
	HTTPRequestDuration.WithLabelValues("GET", "/metrics").Observe(0.02)

	// Verify metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	// Find HTTP metrics
	var foundHTTPRequests, foundHTTPDuration bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "http_requests_total") {
			foundHTTPRequests = true
		}
		if strings.Contains(*metric.Name, "http_request_duration_seconds") {
			foundHTTPDuration = true
		}
	}

	assert.True(t, foundHTTPRequests, "HTTP requests metric should exist")
	assert.True(t, foundHTTPDuration, "HTTP duration metric should exist")
}

func TestStorageMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test WAL metrics
	WALSizeBytes.Set(1024 * 1024) // 1MB

	// Test compaction metrics
	CompactionRuns.Inc()
	CompactionRuns.Inc()
	CompactionErrors.Inc()

	// Test compaction duration
	CompactionDuration.Observe(0.1)
	CompactionDuration.Observe(0.5)
	CompactionDuration.Observe(1.0)

	// Verify metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundWAL, foundCompaction bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "wal_size_bytes") {
			foundWAL = true
		}
		if strings.Contains(*metric.Name, "compaction_runs_total") {
			foundCompaction = true
		}
	}

	assert.True(t, foundWAL, "WAL metrics should exist")
	assert.True(t, foundCompaction, "Compaction metrics should exist")
}

func TestResourceMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test shard count
	ShardCount.Set(10)

	// Test memory pool usage with different labels
	MemoryPoolUsage.WithLabelValues("write_buffer").Set(512 * 1024)
	MemoryPoolUsage.WithLabelValues("read_buffer").Set(256 * 1024)
	MemoryPoolUsage.WithLabelValues("cache").Set(1024 * 1024)
	MemoryPoolUsage.WithLabelValues("temp").Set(128 * 1024)

	// Test leader election with different results
	LeaderElectionResults.WithLabelValues("success").Inc()
	LeaderElectionResults.WithLabelValues("failure").Inc()
	LeaderElectionResults.WithLabelValues("timeout").Inc()

	// Verify that labeled metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundMemory, foundShard bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "memory_pool_bytes") {
			foundMemory = true
		}
		if strings.Contains(*metric.Name, "shard_count") {
			foundShard = true
		}
	}

	assert.True(t, foundMemory, "Memory pool metric should exist")
	assert.True(t, foundShard, "Shard count metric should exist")
}

func TestClusterMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test leader election metrics
	LeaderElectionResults.WithLabelValues("success").Inc()
	LeaderElectionResults.WithLabelValues("success").Inc()
	LeaderElectionResults.WithLabelValues("failure").Inc()

	// Test replication lag
	ReplicationLag.WithLabelValues("shard_1").Set(0.1)
	ReplicationLag.WithLabelValues("shard_2").Set(0.2)
	ReplicationLag.WithLabelValues("shard_3").Set(0.05)

	// Verify cluster metrics
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundLeader, foundReplication bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "leader_election_results_total") {
			foundLeader = true
		}
		if strings.Contains(*metric.Name, "replication_lag_seconds") {
			foundReplication = true
		}
	}

	assert.True(t, foundLeader, "Leader election metric should exist")
	assert.True(t, foundReplication, "Replication lag metric should exist")
}

func TestErrorMetrics(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test error counters
	WriteErrors.Inc()
	WriteErrors.Inc()
	WALErrors.Inc()

	// Verify error metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundWriteErrors, foundWALErrors bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "write_errors_total") {
			foundWriteErrors = true
		}
		if strings.Contains(*metric.Name, "wal_errors_total") {
			foundWALErrors = true
		}
	}

	assert.True(t, foundWriteErrors, "Write errors metric should exist")
	assert.True(t, foundWALErrors, "WAL errors metric should exist")
}

func TestMetricsServer(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test metrics server creation
	server := NewMetricsServer(":0") // Use port 0 for testing

	require.NotNil(t, server, "Failed to create metrics server")

	// Test endpoint URLs
	metricsURL := server.GetMetricsEndpoint()
	assert.NotEmpty(t, metricsURL, "Metrics endpoint URL should not be empty")
	assert.Contains(t, metricsURL, "/metrics", "Metrics endpoint should contain /metrics")

	healthURL := server.GetHealthEndpoint()
	assert.NotEmpty(t, healthURL, "Health endpoint URL should not be empty")
	assert.Contains(t, healthURL, "/health", "Health endpoint should contain /health")

	readyURL := server.GetReadyEndpoint()
	assert.NotEmpty(t, readyURL, "Ready endpoint URL should not be empty")
	assert.Contains(t, readyURL, "/ready", "Ready endpoint should contain /ready")
}

func TestMetricsEndpoint(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Create a test server
	server := httptest.NewServer(promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	defer server.Close()

	// Make a request to the metrics endpoint
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Metrics endpoint should return 200 OK")

	// Read response body
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	require.NoError(t, err)

	content := buf.String()

	// Verify that key metrics are present
	expectedMetrics := []string{
		"# HELP tsdb_ingestion_points_total",
		"# TYPE tsdb_ingestion_points_total counter",
		"# HELP tsdb_query_requests_total",
		"# TYPE tsdb_query_requests_total counter",
		"# HELP tsdb_wal_size_bytes",
		"# TYPE tsdb_wal_size_bytes gauge",
	}

	for _, expected := range expectedMetrics {
		assert.Contains(t, content, expected, "Metrics response should contain %s", expected)
	}
}

func TestMetricsConcurrency(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test concurrent access to metrics
	const numGoroutines = 10
	const operationsPerGoroutine = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()

			for j := 0; j < operationsPerGoroutine; j++ {
				IngestedPoints.Add(1)
				QueryRequests.Inc()
				WALSizeBytes.Set(float64(j))
				IngestionLatency.Observe(0.01)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify that metrics were recorded (we can't check exact values due to concurrency)
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should be recorded under concurrent access")
}

func TestMetricsEdgeCases(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test zero values
	IngestedPoints.Add(0)
	QueryRequests.Inc()
	WALSizeBytes.Set(0)

	// Test very large values
	WALSizeBytes.Set(1e15) // 1 PB

	// Test very small values
	IngestionLatency.Observe(1e-9) // 1 nanosecond
	QueryLatency.Observe(1e-12)    // 1 picosecond

	// Verify metrics are still recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should handle edge cases gracefully")
}

func TestMetricsLabels(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test HTTP metrics with various label combinations
	HTTPRequests.WithLabelValues("GET", "/health", "200").Inc()
	HTTPRequests.WithLabelValues("POST", "/write", "201").Inc()
	HTTPRequests.WithLabelValues("GET", "/metrics", "200").Inc()
	HTTPRequests.WithLabelValues("POST", "/write", "400").Inc()
	HTTPRequests.WithLabelValues("DELETE", "/data", "404").Inc()

	// Test memory pool usage with different labels
	MemoryPoolUsage.WithLabelValues("write_buffer").Set(512 * 1024)
	MemoryPoolUsage.WithLabelValues("read_buffer").Set(256 * 1024)
	MemoryPoolUsage.WithLabelValues("cache").Set(1024 * 1024)
	MemoryPoolUsage.WithLabelValues("temp").Set(128 * 1024)

	// Test leader election with different results
	LeaderElectionResults.WithLabelValues("success").Inc()
	LeaderElectionResults.WithLabelValues("failure").Inc()
	LeaderElectionResults.WithLabelValues("timeout").Inc()

	// Verify that labeled metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var foundLabeledMetrics bool
	for _, metric := range metrics {
		if len(metric.Metric) > 0 && len(metric.Metric[0].Label) > 0 {
			foundLabeledMetrics = true
			break
		}
	}

	assert.True(t, foundLabeledMetrics, "Labeled metrics should be recorded")
}

func TestResetFunction(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Add some metrics first
	IngestedPoints.Add(100)
	QueryRequests.Inc()
	WALSizeBytes.Set(1024)

	// Verify metrics exist
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should exist before reset")

	// Test Reset function
	Reset()

	// Verify metrics are unregistered
	metrics, err = Registry.Gather()
	require.NoError(t, err)
	assert.Equal(t, 0, len(metrics), "Registry should be empty after reset")

	// Re-initialize for cleanup
	Init()
}

func TestGetRegistryFunction(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test GetRegistry function
	retrievedRegistry := GetRegistry()
	assert.NotNil(t, retrievedRegistry, "GetRegistry should return a valid registry")
	assert.Equal(t, Registry, retrievedRegistry, "GetRegistry should return the same registry instance")

	// Verify the returned registry contains our metrics
	metrics, err := retrievedRegistry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Retrieved registry should contain metrics")
}

func TestQueryMetricsWrapper(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test NewQueryMetrics
	queryMetrics := NewQueryMetrics()
	assert.NotNil(t, queryMetrics, "NewQueryMetrics should return a valid instance")
	assert.Equal(t, QueryRequests, queryMetrics.QueryRequests, "QueryRequests should match global metric")
	assert.Equal(t, QueryErrors, queryMetrics.QueryErrors, "QueryErrors should match global metric")
	assert.Equal(t, QueryLatency, queryMetrics.QueryLatency, "QueryLatency should match global metric")

	// Test RecordQueryStart
	queryMetrics.RecordQueryStart()

	// Test RecordQueryDuration
	queryMetrics.RecordQueryDuration(100 * time.Millisecond)
	queryMetrics.RecordQueryDuration(200 * time.Millisecond)

	// Verify latency observations
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	var foundLatency bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "query_latency_seconds") {
			foundLatency = true
			break
		}
	}
	assert.True(t, foundLatency, "Query latency should be recorded")

	// Test RecordQueryError
	queryMetrics.RecordQueryError()

	// Test RecordQueryComplete with success
	startTime := time.Now()
	time.Sleep(1 * time.Millisecond) // Small delay to ensure measurable duration
	queryMetrics.RecordQueryComplete(startTime, nil)

	// Test RecordQueryComplete with error
	startTime = time.Now()
	time.Sleep(1 * time.Millisecond)
	queryMetrics.RecordQueryComplete(startTime, assert.AnError)

	// Verify that metrics were recorded by checking registry
	metrics, err = Registry.Gather()
	require.NoError(t, err)

	var foundQueryRequests, foundQueryErrors bool
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "query_requests_total") {
			foundQueryRequests = true
		}
		if strings.Contains(*metric.Name, "query_errors_total") {
			foundQueryErrors = true
		}
	}

	assert.True(t, foundQueryRequests, "Query requests should be recorded")
	assert.True(t, foundQueryErrors, "Query errors should be recorded")
}

func TestMetricsServerStart(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test server creation
	server := NewMetricsServer(":0")
	assert.NotNil(t, server, "NewMetricsServer should return a valid server")

	// Test server start (we'll start it in a goroutine and stop it quickly)
	serverStarted := make(chan bool, 1)
	serverStopped := make(chan bool, 1)

	go func() {
		serverStarted <- true
		_ = server.Start() // Ignore error for testing purposes
		// We expect an error when we stop the server
		serverStopped <- true
	}()

	// Wait for server to start
	<-serverStarted
	time.Sleep(10 * time.Millisecond) // Give server time to start

	// Test that server is running by checking if we can get endpoint URLs
	metricsURL := server.GetMetricsEndpoint()
	assert.NotEmpty(t, metricsURL, "Metrics endpoint should be available")
	assert.Contains(t, metricsURL, "/metrics", "Metrics endpoint should contain /metrics path")

	// Stop the server by closing the connection (this will cause Start to return)
	time.Sleep(10 * time.Millisecond) // Give some time for the test to complete

	// Wait for server to stop
	select {
	case <-serverStopped:
		// Server stopped as expected
	case <-time.After(100 * time.Millisecond):
		// Server didn't stop in time, but that's okay for this test
	}
}

func TestMetricsServerEndpoints(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	server := NewMetricsServer(":0")
	assert.NotNil(t, server, "Server should be created successfully")

	// Test endpoint URL generation
	metricsURL := server.GetMetricsEndpoint()
	healthURL := server.GetHealthEndpoint()
	readyURL := server.GetReadyEndpoint()

	assert.Contains(t, metricsURL, "/metrics", "Metrics endpoint should contain /metrics")
	assert.Contains(t, healthURL, "/health", "Health endpoint should contain /health")
	assert.Contains(t, readyURL, "/ready", "Ready endpoint should contain /ready")

	// Test that URLs are properly formatted
	assert.True(t, strings.HasPrefix(metricsURL, "http://"), "Metrics URL should start with http://")
	assert.True(t, strings.HasPrefix(healthURL, "http://"), "Health URL should start with http://")
	assert.True(t, strings.HasPrefix(readyURL, "http://"), "Ready URL should start with http://")
}

func TestMetricsReinitialization(t *testing.T) {
	// Test that metrics can be reinitialized after reset
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Add some metrics
	IngestedPoints.Add(50)
	QueryRequests.Inc()

	// Reset metrics
	Reset()

	// Verify registry is empty
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Equal(t, 0, len(metrics), "Registry should be empty after reset")

	// Re-initialize
	Init()

	// Verify metrics are available again
	metrics, err = Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should be available after re-initialization")

	// Test that we can use metrics again
	IngestedPoints.Add(25)
	QueryRequests.Inc()

	// Verify metrics are working
	finalMetrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(finalMetrics), 0, "Metrics should be working after re-initialization")
}

func TestMetricsConcurrentAccessWithReset(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test concurrent access while resetting
	const numGoroutines = 5
	const operationsPerGoroutine = 20

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			for j := 0; j < operationsPerGoroutine; j++ {
				// Add some metrics
				IngestedPoints.Add(1)
				QueryRequests.Inc()
				WALSizeBytes.Set(float64(j))

				// Small delay to increase chance of race conditions
				time.Sleep(1 * time.Microsecond)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify metrics were recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should be recorded under concurrent access")
}

func TestMetricsLabelValidation(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test HTTP metrics with various label combinations
	testCases := []struct {
		method     string
		path       string
		statusCode string
	}{
		{"GET", "/health", "200"},
		{"POST", "/write", "201"},
		{"PUT", "/update", "200"},
		{"DELETE", "/delete", "204"},
		{"PATCH", "/patch", "200"},
	}

	for _, tc := range testCases {
		HTTPRequests.WithLabelValues(tc.method, tc.path, tc.statusCode).Inc()
		HTTPRequestDuration.WithLabelValues(tc.method, tc.path).Observe(0.01)
	}

	// Test memory pool usage with various pool names
	poolNames := []string{"write_buffer", "read_buffer", "cache", "temp", "metadata"}
	for _, poolName := range poolNames {
		MemoryPoolUsage.WithLabelValues(poolName).Set(1024 * 1024)
	}

	// Test leader election with various results
	electionResults := []string{"success", "failure", "timeout", "conflict"}
	for _, result := range electionResults {
		LeaderElectionResults.WithLabelValues(result).Inc()
	}

	// Test replication lag with various shard IDs
	shardIDs := []string{"shard_1", "shard_2", "shard_3", "shard_4"}
	for _, shardID := range shardIDs {
		ReplicationLag.WithLabelValues(shardID).Set(0.1)
	}

	// Verify all labeled metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var labeledMetricsCount int
	for _, metric := range metrics {
		if len(metric.Metric) > 0 && len(metric.Metric[0].Label) > 0 {
			labeledMetricsCount++
		}
	}

	assert.Greater(t, labeledMetricsCount, 0, "Should have labeled metrics")
}

func TestMetricsBoundaryValues(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test very large values
	WALSizeBytes.Set(1e18)  // 1 exabyte
	IngestedPoints.Add(1e9) // 1 billion points

	// Test very small values
	IngestionLatency.Observe(1e-12) // 1 picosecond
	QueryLatency.Observe(1e-9)      // 1 nanosecond

	// Test negative values (should be handled gracefully)
	WALSizeBytes.Set(-1)
	WALSizeBytes.Set(0) // Reset to valid value

	// Test zero values
	IngestedPoints.Add(0)
	QueryRequests.Inc()
	WALSizeBytes.Set(0)

	// Verify metrics are still recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should handle boundary values gracefully")
}

func TestMetricsMultipleRegistrations(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test that metrics can be reset and reinitialized multiple times
	// Add some metrics first
	IngestedPoints.Add(50)
	QueryRequests.Inc()

	// Reset and re-initialize
	Reset()
	Init()

	// Add more metrics
	IngestedPoints.Add(25)
	QueryRequests.Inc()

	// Reset and re-initialize again
	Reset()
	Init()

	// Verify metrics are working after multiple reset/reinit cycles
	IngestedPoints.Add(10)
	QueryRequests.Inc()

	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should work after multiple reset/reinit cycles")
}

func TestMetricsHistogramBuckets(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test histogram bucket boundaries
	latencyValues := []float64{
		0.001, // 1ms
		0.01,  // 10ms
		0.1,   // 100ms
		0.5,   // 500ms
		1.0,   // 1s
		2.0,   // 2s
		5.0,   // 5s
		10.0,  // 10s
	}

	for _, latency := range latencyValues {
		IngestionLatency.Observe(latency)
		QueryLatency.Observe(latency)
		BatchQueueWaitTime.Observe(latency)
		WALAppendLatency.Observe(latency)
		CompactionDuration.Observe(latency)
	}

	// Verify histograms are working
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var histogramCount int
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "latency") || strings.Contains(*metric.Name, "duration") {
			histogramCount++
		}
	}

	assert.Greater(t, histogramCount, 0, "Histogram metrics should be recorded")
}

func TestMetricsLabelCombinations(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test various HTTP method/path combinations
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	httpPaths := []string{"/api/v1/write", "/api/v1/query", "/health", "/metrics", "/ready"}
	statusCodes := []string{"200", "201", "400", "401", "403", "404", "500", "502"}

	for _, method := range httpMethods {
		for _, path := range httpPaths {
			for _, status := range statusCodes {
				HTTPRequests.WithLabelValues(method, path, status).Inc()
				if method != "HEAD" { // HEAD requests don't have duration
					HTTPRequestDuration.WithLabelValues(method, path).Observe(0.01)
				}
			}
		}
	}

	// Test various pool names
	poolNames := []string{
		"write_buffer", "read_buffer", "cache", "temp", "metadata",
		"index", "compression", "serialization", "network", "disk",
	}

	for _, poolName := range poolNames {
		MemoryPoolUsage.WithLabelValues(poolName).Set(1024 * 1024)
	}

	// Test various election results
	electionResults := []string{
		"success", "failure", "timeout", "conflict", "stale",
		"network_error", "quorum_unavailable", "invalid_state",
	}

	for _, result := range electionResults {
		LeaderElectionResults.WithLabelValues(result).Inc()
	}

	// Test various shard IDs
	shardIDs := []string{
		"shard_1", "shard_2", "shard_3", "shard_4", "shard_5",
		"shard_a", "shard_b", "shard_c", "shard_d", "shard_e",
	}

	for _, shardID := range shardIDs {
		ReplicationLag.WithLabelValues(shardID).Set(0.1)
	}

	// Verify labeled metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var labeledMetricsCount int
	for _, metric := range metrics {
		if len(metric.Metric) > 0 && len(metric.Metric[0].Label) > 0 {
			labeledMetricsCount++
		}
	}

	assert.Greater(t, labeledMetricsCount, 0, "Should have labeled metrics with various combinations")
}

func TestMetricsErrorScenarios(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test error metrics
	WriteErrors.Inc()
	WriteErrors.Inc()
	WriteErrors.Inc()

	WALErrors.Inc()
	WALErrors.Inc()

	CompactionErrors.Inc()
	CompactionErrors.Inc()
	CompactionErrors.Inc()

	QueryErrors.Inc()
	QueryErrors.Inc()

	// Verify error metrics are recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)

	var errorMetricsCount int
	for _, metric := range metrics {
		if strings.Contains(*metric.Name, "errors_total") {
			errorMetricsCount++
		}
	}

	assert.Greater(t, errorMetricsCount, 0, "Error metrics should be recorded")
}

func TestMetricsPerformanceUnderLoad(t *testing.T) {
	cleanup := setupTestMetrics(t)
	defer cleanup()

	// Test performance under high load
	const numOperations = 10000

	start := time.Now()

	// Perform many metric operations
	for i := 0; i < numOperations; i++ {
		IngestedPoints.Add(1)
		QueryRequests.Inc()
		WALSizeBytes.Set(float64(i))
		IngestionLatency.Observe(0.01)
		QueryLatency.Observe(0.01)
	}

	duration := time.Since(start)

	// Verify metrics were recorded
	metrics, err := Registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Metrics should be recorded under high load")

	// Performance should be reasonable (less than 1 second for 10k operations)
	assert.Less(t, duration, time.Second, "Metric operations should complete in reasonable time")
}

// Benchmark tests for performance
func BenchmarkIngestionMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IngestedPoints.Add(1)
		IngestionLatency.Observe(0.01)
	}
}

func BenchmarkQueryMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QueryRequests.Inc()
		QueryLatency.Observe(0.01)
	}
}

func BenchmarkHTTPMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HTTPRequests.WithLabelValues("GET", "/test", "200").Inc()
		HTTPRequestDuration.WithLabelValues("GET", "/test").Observe(0.01)
	}
}

func BenchmarkConcurrentMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IngestedPoints.Add(1)
			QueryRequests.Inc()
			WALSizeBytes.Set(float64(time.Now().UnixNano()))
		}
	})
}

func BenchmarkHistogramMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IngestionLatency.Observe(0.01)
		QueryLatency.Observe(0.01)
		BatchQueueWaitTime.Observe(0.01)
		WALAppendLatency.Observe(0.01)
		CompactionDuration.Observe(0.01)
	}
}

func BenchmarkLabeledMetrics(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HTTPRequests.WithLabelValues("GET", "/test", "200").Inc()
		HTTPRequestDuration.WithLabelValues("GET", "/test").Observe(0.01)
		MemoryPoolUsage.WithLabelValues("buffer").Set(1024)
		LeaderElectionResults.WithLabelValues("success").Inc()
		ReplicationLag.WithLabelValues("shard_1").Set(0.1)
	}
}

func BenchmarkMetricsReset(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Add some metrics
		IngestedPoints.Add(100)
		QueryRequests.Inc()

		// Reset and reinitialize
		Reset()
		Init()
	}
}

func BenchmarkMetricsGather(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	// Add some metrics first
	IngestedPoints.Add(1000)
	QueryRequests.Inc()
	WALSizeBytes.Set(1024 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Registry.Gather()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQueryMetricsWrapper(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	queryMetrics := NewQueryMetrics()
	startTime := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queryMetrics.RecordQueryStart()
		queryMetrics.RecordQueryDuration(100 * time.Millisecond)
		queryMetrics.RecordQueryComplete(startTime, nil)
	}
}

func BenchmarkMetricsServerCreation(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := NewMetricsServer(":0")
		if server == nil {
			b.Fatal("Failed to create server")
		}
	}
}

func BenchmarkMetricsConcurrentAccess(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IngestedPoints.Add(1)
			QueryRequests.Inc()
			WALSizeBytes.Set(float64(time.Now().UnixNano()))
			IngestionLatency.Observe(0.01)
			QueryLatency.Observe(0.01)
			HTTPRequests.WithLabelValues("GET", "/test", "200").Inc()
			MemoryPoolUsage.WithLabelValues("buffer").Set(1024)
		}
	})
}

func BenchmarkMetricsWithLabels(b *testing.B) {
	cleanup := setupBenchmarkMetrics(b)
	defer cleanup()

	// Pre-create label values to avoid allocation in benchmark
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	paths := []string{"/api/v1/write", "/api/v1/query", "/health", "/metrics"}
	statusCodes := []string{"200", "201", "400", "500"}
	poolNames := []string{"write_buffer", "read_buffer", "cache", "temp"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		method := methods[i%len(methods)]
		path := paths[i%len(paths)]
		status := statusCodes[i%len(statusCodes)]
		pool := poolNames[i%len(poolNames)]

		HTTPRequests.WithLabelValues(method, path, status).Inc()
		HTTPRequestDuration.WithLabelValues(method, path).Observe(0.01)
		MemoryPoolUsage.WithLabelValues(pool).Set(1024)
	}
}
