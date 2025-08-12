# Performance Guide

Complete guide for performance testing, benchmarking, and optimization of TimeSeriesDB.

## Overview

TimeSeriesDB includes a comprehensive performance testing and monitoring system that helps maintain consistent performance across code changes.

## Quick Start

### 1. Set Up Performance Baseline
```bash
# Run all benchmarks and set as baseline
make regression-baseline
# Or use the script directly
./scripts/run-benchmarks.sh -b
```

### 2. Run Performance Monitoring
```bash
# Complete monitoring workflow (recommended)
make performance-monitor

# Or run individual steps:
make benchmark-all          # Run benchmarks
make regression-detect      # Detect regressions
make dashboard              # Generate dashboard
```

### 3. View Results
```bash
# Open performance dashboard in browser
make dashboard-open
# View regression reports
ls benchmark-results/regression_report_*.txt
```

## Benchmark Categories & Commands

### Available Commands

**Makefile (Recommended):**
```bash
make benchmark              # Run all benchmarks
make benchmark-parser       # Parser performance
make benchmark-storage      # Storage performance  
make benchmark-http         # HTTP endpoint performance
make benchmark-e2e          # End-to-end workflows
make benchmark-memory       # Memory usage
make benchmark-profile      # With profiling
make benchmark-clean        # Clean artifacts
```

**Scripts:**
```bash
./scripts/run-benchmarks.sh           # All benchmarks
./scripts/run-benchmarks.sh -p        # Parser only
./scripts/run-benchmarks.sh -s        # Storage only
./scripts/run-benchmarks.sh -e        # HTTP endpoints only
./scripts/run-benchmarks.sh -t        # End-to-end workflows only
./scripts/run-benchmarks.sh -m        # Memory usage only
./scripts/run-benchmarks.sh -c        # Compare with baseline
./scripts/run-benchmarks.sh -b        # Set as baseline
```

**Direct Go Commands:**
```bash
go test -bench=. -benchmem ./test/                    # All benchmarks
go test -bench=BenchmarkParse -benchmem ./test/       # Parser tests
go test -bench=BenchmarkWrite -benchmem ./test/       # Storage tests
go test -bench=BenchmarkHTTP -benchmem ./test/        # HTTP tests
go test -bench=. -benchmem -timeout=10m ./test/       # With timeout
```

### Benchmark Types

**Parser Performance:**
- Simple/complex line parsing, multi-line parsing, large datasets (1000+ lines), scalability tests

**Storage Performance:**
- Single/multiple point writes, points with many fields/tags

**HTTP Endpoint Performance:**
- Single/multiple point HTTP writes, large dataset operations

**End-to-End Workflows:**
- Parse → Store → Retrieve cycles, concurrent operations

**Memory Usage:**
- Allocation tracking, garbage collection pressure analysis

## Understanding Results

### Benchmark Output Format
```
BenchmarkName-16         1000        1234567 ns/op        1234 B/op        10 allocs/op
```

**Fields:**
- **BenchmarkName-16**: Name and CPU cores
- **1000**: Number of iterations
- **1234567 ns/op**: Time per operation (nanoseconds)
- **1234 B/op**: Memory allocated per operation (bytes)
- **10 allocs/op**: Number of allocations per operation

### Performance Metrics
- **Throughput**: Operations per second (higher is better)
- **Latency**: Time per operation (lower is better)
- **Memory**: Bytes allocated per operation (lower is better)
- **Allocations**: Number of memory allocations per operation (lower is better)

### Performance Tips
1. **Run multiple times** to account for system variance
2. **Use consistent environment** for comparable results
3. **Monitor system resources** during execution
4. **Profile first** to identify bottlenecks before optimizing

## Performance Regression Detection

### Setting & Detecting Regressions
```bash
# Set baseline
make regression-baseline

# Detect regressions
make regression-detect              # Basic detection
make regression-detect-html         # With HTML report
make regression-detect-json         # With JSON output
make regression-detect-full         # All outputs
```

### Regression Thresholds
- **Warning Threshold**: 5% (default) - Triggers warnings
- **Critical Threshold**: 15% (default) - Triggers critical alerts

**Customize thresholds:**
```bash
./scripts/detect-regressions.sh -t 3.0 --critical-threshold 10.0
```

### Regression Report Example
```
=== Performance Regression Analysis ===
Benchmark Name                    | Baseline (ns/op) | Current (ns/op) | Change (%) | Status
----------------------------------|------------------|-----------------|------------|---------
BenchmarkParseSimpleLine          |             1000 |            1100 |     +10.00% | ⚠️  REGRESSION
BenchmarkParseComplexLine         |             2000 |            1900 |      -5.00% | 🚀 IMPROVEMENT
BenchmarkWritePoint               |             5000 |            6000 |     +20.00% | 🚨 CRITICAL
```

## Profiling & Analysis

### CPU & Memory Profiling
```bash
# Generate profiles
go test -bench=BenchmarkParseLargeDataset -cpuprofile=cpu.prof -benchmem ./test/
go test -bench=BenchmarkMemoryUsage -memprofile=memory.prof -benchmem ./test/

# Analyze profiles
go tool pprof cpu.prof
go tool pprof memory.prof

# Web interface
go tool pprof -http=:8080 cpu.prof
go tool pprof -http=:8080 memory.prof

# Using Makefile
make benchmark-profile  # Generates cpu_profile.prof and memory_profile.prof
```

## Performance Dashboard

### Generating Dashboard
```bash
# Complete dashboard
make dashboard
make dashboard-trends      # Trends analysis
make dashboard-summary     # Summary report only
make dashboard-open        # Open in browser

# Using script
./scripts/performance-dashboard.sh -g              # Generate complete
./scripts/performance-dashboard.sh -t              # Trends analysis
./scripts/performance-dashboard.sh -s              # Summary only
./scripts/performance-dashboard.sh -g -o my-dashboard  # Custom output
./scripts/performance-dashboard.sh -g --days 7     # Last 7 days
```

## CI/CD Integration

### Automated Workflows
- **Automated Testing**: Runs on every PR and push
- **Performance Monitoring**: Tracks benchmarks over time
- **Regression Detection**: Automatically flags performance issues
- **PR Comments**: Provides performance summaries on pull requests

### Workflow Files
- `.github/workflows/benchmark.yml` - Basic benchmarking
- `.github/workflows/performance-regression.yml` - Regression detection

### Automated Checks
- ✅ **Pass** when no regressions are detected
- ⚠️ **Warn** when performance regressions are detected
- 🚨 **Fail** when critical regressions are detected

## Best Practices

### For Developers
1. **Run Benchmarks Locally**: Always test performance changes locally before pushing
2. **Set Baselines**: Establish performance baselines for your development environment
3. **Monitor Trends**: Watch for performance regressions over time
4. **Profile Issues**: Use profiling tools to identify bottlenecks

### For CI/CD
1. **Consistent Environment**: Use the same Go versions and OS for reliable results
2. **Artifact Retention**: Keep benchmark results for trend analysis
3. **Failure Handling**: Ensure benchmarks don't block critical deployments

### Regular Monitoring
```bash
# Update baseline weekly or after significant changes
make regression-baseline

# Generate trends analysis weekly
make dashboard-trends

# Monitor for regressions
make regression-detect
```

## Performance Optimization

### Code-Level Optimizations
1. **Reduce allocations**: Reuse objects when possible
2. **Use sync.Pool**: For frequently allocated/deallocated objects
3. **Avoid string concatenation**: Use strings.Builder for multiple concatenations
4. **Profile first**: Identify bottlenecks before optimizing

### Storage Optimizations
1. **Batch writes**: Group multiple points into single operations
2. **Efficient data structures**: Use appropriate data types
3. **Memory mapping**: For large datasets
4. **Compression**: For historical data

### HTTP Optimizations
1. **Connection pooling**: Reuse HTTP connections
2. **Request batching**: Send multiple points per request
3. **Compression**: Enable gzip compression
4. **Async processing**: Process requests asynchronously

## Integration & Advanced Usage

### Pre-commit Hooks
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running performance regression check..."
make regression-detect

if [ $? -eq 1 ]; then
    echo "⚠️  Performance regressions detected. Consider reviewing changes."
elif [ $? -eq 2 ]; then
    echo "🚨 Critical performance regressions detected. Please fix before committing."
    exit 1
fi
```

### IDE Integration
Many IDEs support running Make targets:
- **VS Code**: Use the "Make" extension
- **GoLand**: Configure external tools
- **Vim/Emacs**: Use terminal integration

### Custom Benchmark Suites
```go
// test/custom_benchmark_test.go
func BenchmarkCustomOperation(b *testing.B) {
    // Your custom benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Operation to benchmark
    }
}
```

### Continuous Monitoring & Alerts
```bash
# Cron job for daily monitoring
0 2 * * * cd /path/to/tsdb && make performance-monitor

# Check for critical regressions and send alerts
if ./scripts/detect-regressions.sh --critical-threshold 10.0; then
    # Send alert (email, Slack, etc.)
    echo "Critical performance regression detected!" | mail -s "Performance Alert" admin@example.com
fi
```

## Troubleshooting

### Common Issues
**No baseline found:**
```bash
make regression-baseline
```

**Benchmark failures:**
```bash
go version              # Check Go version compatibility
go mod tidy             # Verify dependencies
go mod verify
```

**Script permissions:**
```bash
chmod +x scripts/*.sh
```

### Debug Mode
```bash
# Verbose output for debugging
./scripts/detect-regressions.sh -v
./scripts/performance-dashboard.sh -v
```

## Next Steps

- **[Installation Guide](INSTALLATION.md)** - Set up TimeSeriesDB
- **[API Reference](API_REFERENCE.md)** - Understand the API
- **[Development Guide](DEVELOPMENT.md)** - Contribute to the project
- **[CI/CD Guide](CI_CD.md)** - Set up automation

For performance questions, check the [GitHub Issues](https://github.com/yourusername/timeseriesdb/issues) or create a new one.
