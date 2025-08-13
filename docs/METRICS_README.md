# TimeSeriesDB Metrics System

This document describes the comprehensive metrics system integrated into TimeSeriesDB, which exposes Prometheus-compatible metrics for monitoring and observability.

## Overview

The TimeSeriesDB metrics system provides:
- **HTTP metrics**: Request counts, durations, and status codes
- **Ingestion metrics**: Points and batches processed, latency, errors
- **Storage metrics**: WAL operations, compaction, shard management
- **Resource metrics**: Memory usage, system performance
- **Cluster metrics**: Leader election, replication lag

## Quick Start

### 1. Start the Server

The metrics endpoint is automatically available when you start TimeSeriesDB:

```bash
go run .  # or go build && ./timeseriesdb
```

### 2. Access Metrics

- **Metrics endpoint**: `http://localhost:8080/metrics`
- **Health endpoint**: `http://localhost:8080/health`
- **Write endpoint**: `http://localhost:8080/write`

### 3. Demo Script

Run the included demo script to see the metrics in action:

```bash
./scripts/demo_metrics.sh
```

## Available Metrics

### HTTP API Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_http_requests_total` | Counter | Total HTTP requests by method, path, and status code |
| `tsdb_http_request_duration_seconds` | Histogram | HTTP request duration by method and path |

### Ingestion Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_ingestion_points_total` | Counter | Total number of data points ingested |
| `tsdb_ingestion_batches_total` | Counter | Total number of batches ingested |
| `tsdb_ingestion_latency_seconds` | Histogram | Time taken to ingest points |
| `tsdb_batch_queue_wait_seconds` | Histogram | Time spent waiting in batch queue |
| `tsdb_wal_append_latency_seconds` | Histogram | Time taken to append to WAL |

### Storage Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_wal_size_bytes` | Gauge | Current size of WAL in bytes |
| `tsdb_wal_errors_total` | Counter | Total number of WAL errors |
| `tsdb_compaction_runs_total` | Counter | Total number of compaction runs |
| `tsdb_compaction_duration_seconds` | Histogram | Time taken for compaction operations |
| `tsdb_compaction_errors_total` | Counter | Total number of compaction errors |
| `tsdb_shard_count` | Gauge | Current number of shards |

### Query Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_query_requests_total` | Counter | Total number of query requests |
| `tsdb_query_latency_seconds` | Histogram | Time taken to execute queries |
| `tsdb_query_errors_total` | Counter | Total number of query errors |

### Resource Usage Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_memory_pool_bytes` | Gauge | Memory pool usage by pool name |
| `tsdb_write_errors_total` | Counter | Total number of write errors |

### Cluster Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `tsdb_leader_election_results_total` | Counter | Leader election results |
| `tsdb_replication_lag_seconds` | Gauge | Replication lag by shard ID |

## Architecture

### Components

1. **Metrics Registry** (`internal/metrics/prometheus.go`)
   - Central Prometheus registry
   - Metric definitions and initialization

2. **Metrics Server** (`internal/metrics/server.go`)
   - Standalone metrics server (optional)
   - Health and readiness endpoints

3. **HTTP Middleware** (`internal/api/middleware/metrics.go`)
   - Automatic HTTP request metrics collection
   - Request duration and status code tracking

4. **Router Integration** (`internal/api/http/router.go`)
   - Metrics endpoint exposure
   - Middleware integration

### Data Flow

```
HTTP Request → Metrics Middleware → Handler → Response
     ↓
HTTP Metrics (count, duration, status)
     ↓
Prometheus Registry
     ↓
/metrics endpoint
```

## Configuration

### Environment Variables

The metrics system respects the following configuration:

```bash
# Server configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s

# Logging configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

### Customization

To add custom metrics:

1. Define the metric in `internal/metrics/prometheus.go`:
```go
var MyCustomMetric = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "tsdb_my_custom_metric_total",
        Help: "Description of my custom metric",
    },
)
```

2. Register it in the `Init()` function:
```go
func Init() {
    Registry.MustRegister(
        // ... existing metrics
        MyCustomMetric,
    )
}
```

## Monitoring with Prometheus

### Prometheus Configuration

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'timeseriesdb'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboards

Create dashboards using the available metrics:

- **System Overview**: HTTP requests, errors, latency
- **Ingestion Performance**: Points per second, batch processing
- **Storage Health**: WAL size, compaction status
- **Resource Usage**: Memory pools, shard count

## Testing

### Run Tests

```bash
# Run all tests
go test ./test -v

# Run metrics tests only
go test ./test -v -run TestMetrics

# Run specific test
go test ./test -v -run TestMetricsEndpoint
```

### Manual Testing

1. Start the server:
```bash
go run .
```

2. Test endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Metrics
curl http://localhost:8080/metrics

# Write (will generate metrics)
curl -X POST http://localhost:8080/write \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'
```

## Troubleshooting

### Common Issues

1. **Empty metrics response**
   - Ensure `metrics.Init()` is called in main
   - Check that metrics are registered with the correct registry

2. **HTTP metrics not appearing**
   - Verify middleware is properly integrated
   - Check that requests are actually being made

3. **Build errors**
   - Ensure all dependencies are installed: `go mod tidy`
   - Check Go version compatibility

### Debug Mode

Enable debug logging to troubleshoot metrics issues:

```bash
export LOG_LEVEL=debug
go run .
```

## Performance Considerations

- **Memory**: Each metric consumes a small amount of memory
- **CPU**: Histogram operations have minimal CPU overhead
- **Network**: Metrics endpoint response size grows with metric count
- **Storage**: Metrics are in-memory only (no persistence)

## Best Practices

1. **Metric Naming**: Use consistent naming conventions (`tsdb_*`)
2. **Labels**: Use labels sparingly to avoid cardinality explosion
3. **Documentation**: Always provide helpful descriptions for metrics
4. **Testing**: Test metrics collection in your integration tests
5. **Monitoring**: Set up alerts for critical metrics

## Integration Examples

### Custom Application Metrics

```go
import "timeseriesdb/internal/metrics"

// Increment a counter
metrics.IngestedPoints.Inc()

// Record a duration
timer := prometheus.NewTimer(metrics.IngestionLatency)
defer timer.ObserveDuration()

// Set a gauge value
metrics.ShardCount.Set(float64(shardCount))
```

### External Monitoring

```bash
# Check metrics from command line
curl -s http://localhost:8080/metrics | grep tsdb_ingestion_points_total

# Monitor specific metric
watch -n 1 'curl -s http://localhost:8080/metrics | grep tsdb_http_requests_total'
```

## Support

For issues or questions about the metrics system:

1. Check the logs for error messages
2. Verify configuration and environment variables
3. Run the test suite to identify issues
4. Check Prometheus client library documentation

## References

- [Prometheus Client Library](https://github.com/prometheus/client_golang)
- [Prometheus Metrics Best Practices](https://prometheus.io/docs/practices/naming/)
- [Go HTTP Middleware Patterns](https://golang.org/pkg/net/http/#Handler)
