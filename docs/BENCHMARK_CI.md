
# Benchmark CI System


This document describes the comprehensive benchmark CI system for TimeSeriesDB, which automatically runs performance tests and provides detailed analysis.

## Overview

The benchmark CI system consists of:

1. **GitHub Actions Workflow** (`.github/workflows/benchmark.yml`) - Automated benchmark execution
2. **Local Benchmark Script** (`scripts/run-benchmarks.sh`) - Developer tools for local testing
3. **Makefile Targets** - Convenient commands for running benchmarks
4. **Artifact Storage** - Benchmark results and profiling data

## GitHub Actions Workflow

### Trigger Events

- **Push/Pull Request**: Runs on main/master branch changes
- **Manual**: Can be triggered manually via GitHub Actions UI
- **Scheduled**: Runs daily at 2 AM UTC for continuous monitoring

### Jobs

#### 1. Benchmark Job (Matrix Strategy)

Runs benchmarks across multiple Go versions (1.20, 1.21, 1.22) to ensure compatibility and performance consistency.


#### 2. Benchmark Performance Job

Runs after the main benchmark job completes, providing detailed performance analysis and PR comments.

**Features:**
- Detailed benchmark output
- Results comparison
- Automatic PR comments with benchmark summaries
- Long-term artifact storage (90 days)

### Artifacts

- **Benchmark Profiles**: CPU and memory profiling data
- **Benchmark Results**: Detailed performance metrics
- **Coverage Reports**: Test coverage information
- **Comparison Reports**: Performance regression analysis

## Local Development

### Using the Benchmark Script

The `scripts/run-benchmarks.sh` script provides a convenient way to run benchmarks locally:

```bash
# Run all benchmarks
./scripts/run-benchmarks.sh

# Run specific benchmark categories
./scripts/run-benchmarks.sh -p  # Parser only
./scripts/run-benchmarks.sh -s  # Storage only
./scripts/run-benchmarks.sh -e  # HTTP endpoints only
./scripts/run-benchmarks.sh -t  # End-to-end workflows only
./scripts/run-benchmarks.sh -m  # Memory usage only

# Compare with baseline
./scripts/run-benchmarks.sh -c

# Set current results as baseline
./scripts/run-benchmarks.sh -b

# Custom output file
./scripts/run-benchmarks.sh -o my_results.txt

# Custom timeout
./scripts/run-benchmarks.sh --timeout 5m
```

### Using Makefile Targets

The Makefile provides convenient targets for running benchmarks:

```bash
# Run all benchmarks
make benchmark

# Run specific categories
make benchmark-parser
make benchmark-storage
make benchmark-http
make benchmark-e2e
make benchmark-memory

# Run with profiling
make benchmark-profile

# Clean up artifacts
make benchmark-clean

# Show help
make benchmark-help
```

### Direct Go Commands

For more control, you can run benchmarks directly:

```bash
# All benchmarks
go test -bench=. -benchmem ./test/

# Specific patterns
go test -bench=BenchmarkParse -benchmem ./test/
go test -bench=BenchmarkWrite -benchmem ./test/

# With profiling
go test -bench=BenchmarkParseLargeDataset -cpuprofile=cpu.prof -benchmem ./test/
go test -bench=BenchmarkMemoryUsage -memprofile=mem.prof -benchmem ./test/

# With timeout
go test -bench=. -benchmem -timeout=10m ./test/
```

## Benchmark Categories

### 1. Parser Benchmarks

- **Simple Line Parsing**: Basic line protocol parsing
- **Complex Line Parsing**: Lines with many tags and fields
- **Multi-line Parsing**: Batch processing of multiple lines
- **Large Dataset Parsing**: Performance with 1000+ lines
- **Scalability Tests**: Various line counts (1, 10, 100, 1000, 10000)
- **Tag Count Tests**: Performance with different numbers of tags
- **Field Count Tests**: Performance with different numbers of fields

### 2. Storage Benchmarks

- **Single Point Write**: Individual point insertion
- **Multiple Points Write**: Batch point insertion
- **Many Fields**: Points with numerous fields
- **Many Tags**: Points with numerous tags

### 3. HTTP Endpoint Benchmarks

- **Single Point Write**: HTTP write endpoint simulation
- **Multiple Points Write**: Batch HTTP writes
- **Large Dataset Write**: High-volume HTTP operations

### 4. End-to-End Workflows

- **Complete Workflow**: Parse → Store → Retrieve cycle
- **Concurrent Operations**: Parallel processing performance

### 5. Memory Usage Benchmarks

- **Allocation Tracking**: Memory allocation patterns
- **Garbage Collection**: GC pressure analysis

## Performance Analysis

### Benchmark Output Format

```
BenchmarkParseSimpleLine-8         1000000              1234 ns/op             256 B/op          8 allocs/op
BenchmarkParseComplexLine-8          500000              2468 ns/op             512 B/op         16 allocs/op
```

**Fields:**
- **Benchmark Name**: Function name with CPU count
- **Iterations**: Number of times the benchmark ran
- **Time per Operation**: Nanoseconds per operation
- **Memory per Operation**: Bytes allocated per operation
- **Allocations per Operation**: Number of allocations per operation

### Profiling

The CI system generates both CPU and memory profiles:

- **CPU Profile**: Shows where time is spent during execution
- **Memory Profile**: Shows memory allocation patterns

**Analyzing Profiles:**
```bash
# CPU profile analysis
go tool pprof cpu_profile.prof

# Memory profile analysis
go tool pprof memory_profile.prof

# Web interface
go tool pprof -http=:8080 cpu_profile.prof
go tool pprof -http=:8080 memory_profile.prof
```

## Continuous Monitoring

### Daily Benchmarks

The scheduled workflow runs benchmarks daily to:
- Detect performance regressions
- Monitor long-term trends
- Ensure consistent performance across Go versions

### PR Integration

For pull requests, the CI system:
- Runs comprehensive benchmarks
- Compares against baseline
- Comments on PR with performance summary
- Flags potential regressions

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

### Debugging

1. **Verbose Output**: Use `-v` flag for detailed benchmark information
2. **Profile Analysis**: Generate and analyze profiles for performance issues
3. **Log Analysis**: Check CI logs for error messages and warnings
4. **Local Reproduction**: Reproduce issues locally for detailed investigation

## Future Enhancements

Planned improvements to the benchmark CI system:

1. **Performance Regression Alerts**: Automated notifications for significant changes
2. **Historical Trend Analysis**: Long-term performance tracking and visualization
3. **Benchmark Result Dashboard**: Web interface for viewing results and trends
4. **Performance Budgets**: Automated enforcement of performance constraints
5. **Cross-Platform Testing**: Benchmark execution on multiple operating systems
6. **Hardware Profiling**: CPU and memory usage monitoring during benchmarks

## Integration with Other Systems

### Performance Regression Detection

The benchmark CI system integrates with the performance regression detection system to:
- Automatically detect performance regressions
- Generate regression reports
- Update performance dashboards
- Trigger alerts for critical regressions

### CI/CD Pipeline

The benchmark system is part of the larger CI/CD pipeline:
- Runs alongside other quality checks
- Contributes to overall build status
- Provides performance metrics for releases
- Integrates with deployment automation

## Troubleshooting

### Common CI Issues

**Benchmark timeouts:**
- Increase timeout values in workflow files
- Check for infinite loops in benchmark code
- Verify benchmark data sizes

**Memory issues:**
- Monitor memory usage during benchmarks
- Check for memory leaks in test code
- Adjust memory limits in workflow configuration

**Inconsistent results:**
- Ensure consistent environment across runs
- Check for external dependencies
- Verify benchmark isolation

### Local vs CI Differences

**Environment differences:**
- Go version variations
- OS-specific optimizations
- Hardware differences
- Background processes

**Resolution strategies:**
- Use Docker for consistent environments
- Standardize on specific Go versions
- Document environment requirements
- Use relative performance metrics

## Next Steps

- **[Performance Guide](PERFORMANCE.md)** - Complete performance testing guide
- **[CI/CD Guide](CI_CD.md)** - CI/CD pipeline setup and management
- **[Performance Monitoring](PERFORMANCE_MONITORING.md)** - Regression detection system
- **[Installation Guide](INSTALLATION.md)** - Set up development environment

For benchmark CI questions, check the [GitHub Issues](https://github.com/yourusername/timeseriesdb/issues) or create a new one.


