# Performance Regression Detection

TimeSeriesDB includes a comprehensive performance monitoring and regression detection system that helps maintain consistent performance across code changes.

## Overview

The performance regression detection system consists of:

1. **Automated Benchmarking** - Runs performance tests on your code
2. **Regression Detection** - Compares current performance against baseline
3. **Performance Dashboard** - Visualizes performance trends and metrics
4. **CI/CD Integration** - Automated checks in GitHub Actions
5. **Alerting** - Notifications when performance regressions are detected

## Quick Start

### 1. Set Up Performance Baseline

First, establish a performance baseline by running benchmarks:

```bash
# Run all benchmarks and set as baseline
make regression-baseline

# Or use the script directly
./scripts/run-benchmarks.sh -b
```

### 2. Run Performance Monitoring

After making code changes, run the complete performance monitoring workflow:

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

## Available Make Targets

### Performance Regression Detection

```bash
make regression-detect          # Basic regression detection
make regression-detect-html     # With HTML report
make regression-detect-json     # With JSON output
make regression-detect-full     # With all outputs
make regression-baseline        # Set current results as baseline
make regression-compare         # Compare with baseline
```

### Performance Dashboard

```bash
make dashboard                  # Generate complete dashboard
make dashboard-trends           # Generate trends analysis
make dashboard-summary          # Generate summary report
make dashboard-open             # Open dashboard in browser
```

### Complete Workflow

```bash
make performance-monitor        # Run complete monitoring workflow
make performance-clean          # Clean monitoring artifacts
make performance-help           # Show help
```

## Script Usage

### Regression Detection Script

```bash
# Basic usage
./scripts/detect-regressions.sh

# With specific file
./scripts/detect-regressions.sh -c benchmark-results/my_results.txt

# With custom threshold (default: 5%)
./scripts/detect-regressions.sh -t 10.0

# Generate HTML and JSON reports
./scripts/detect-regressions.sh -H -j

# Help
./scripts/detect-regressions.sh -h
```

**Options:**
- `-c, --current FILE` - Current results file to analyze
- `-b, --baseline FILE` - Baseline file (default: benchmark-results/baseline.txt)
- `-t, --threshold PERCENT` - Performance regression threshold (default: 5%)
- `--critical-threshold PERCENT` - Critical regression threshold (default: 15%)
- `-H, --html` - Generate HTML report
- `-j, --json` - Generate JSON output
- `-o, --output FILE` - Custom output file

### Performance Dashboard Script

```bash
# Generate complete dashboard
./scripts/performance-dashboard.sh -g

# Generate trends analysis
./scripts/performance-dashboard.sh -t

# Generate summary report only
./scripts/performance-dashboard.sh -s

# Custom output directory
./scripts/performance-dashboard.sh -g -o my-dashboard

# Analyze last 7 days
./scripts/performance-dashboard.sh -g --days 7
```

## Configuration

### Thresholds

The system uses configurable thresholds to determine performance regressions:

- **Warning Threshold**: 5% (default) - Triggers warnings for performance regressions
- **Critical Threshold**: 15% (default) - Triggers critical alerts and can fail CI/CD

You can customize these thresholds:

```bash
./scripts/detect-regressions.sh -t 3.0 --critical-threshold 10.0
```

### Baseline Management

The baseline represents your "known good" performance state. Update it when:

- Performance improvements are merged
- After major refactoring
- When switching to new Go versions
- After dependency updates

```bash
# Set new baseline after performance improvements
make regression-baseline

# Or manually
./scripts/run-benchmarks.sh -b
```

## CI/CD Integration

### GitHub Actions

The system includes GitHub Actions workflows that automatically:

1. **Run benchmarks** on every PR and push
2. **Detect regressions** against the baseline
3. **Comment on PRs** with performance summary
4. **Fail workflows** on critical regressions
5. **Generate reports** and upload artifacts

### Workflow Files

- `.github/workflows/benchmark.yml` - Basic benchmarking
- `.github/workflows/performance-regression.yml` - Regression detection

### Automated Checks

The CI/CD system will:

- ✅ **Pass** when no regressions are detected
- ⚠️ **Warn** when performance regressions are detected
- 🚨 **Fail** when critical regressions are detected

## Understanding Results

### Regression Report Format

```
=== Performance Regression Analysis ===
Generated: 2024-01-15 10:30:00
Baseline: benchmark-results/baseline.txt
Current: benchmark-results/benchmark_20240115_103000.txt
Threshold: 5.0%
Critical threshold: 15.0%

Benchmark Name                    | Baseline (ns/op) | Current (ns/op) | Change (%) | Status
----------------------------------|------------------|-----------------|------------|---------
BenchmarkParseSimpleLine          |             1000 |            1100 |     +10.00% | ⚠️  REGRESSION
BenchmarkParseComplexLine         |             2000 |            1900 |      -5.00% | 🚀 IMPROVEMENT
BenchmarkWritePoint               |             5000 |            6000 |     +20.00% | 🚨 CRITICAL

=== Summary ===
Total benchmarks analyzed: 3
Performance regressions: 1
Critical regressions: 1
Performance improvements: 1
```

### Status Indicators

- ✅ **No Change** - Performance within threshold
- ⚠️ **REGRESSION** - Performance degraded above warning threshold
- 🚨 **CRITICAL** - Performance degraded above critical threshold
- 🚀 **IMPROVEMENT** - Performance improved above threshold

## Best Practices

### 1. Regular Baseline Updates

```bash
# Update baseline weekly or after significant changes
make regression-baseline
```

### 2. Monitor Trends

```bash
# Generate trends analysis weekly
make dashboard-trends
```

### 3. Investigate Regressions

When regressions are detected:

1. **Review the code changes** that caused the regression
2. **Check for obvious issues** (inefficient algorithms, memory leaks)
3. **Use profiling tools** to identify bottlenecks
4. **Consider reverting** if the regression is significant
5. **Document the regression** and track resolution

### 4. Performance Budgets

Set performance budgets for critical operations:

```bash
# Use stricter thresholds for critical operations
./scripts/detect-regressions.sh -t 2.0 --critical-threshold 5.0
```

## Troubleshooting

### Common Issues

**No baseline found:**
```bash
# Set a baseline first
make regression-baseline
```

**Benchmark failures:**
```bash
# Check Go version compatibility
go version

# Verify dependencies
go mod tidy
go mod verify
```

**Script permissions:**
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

### Debug Mode

Enable verbose output for debugging:

```bash
# Verbose regression detection
./scripts/detect-regressions.sh -v

# Verbose dashboard generation
./scripts/performance-dashboard.sh -v
```

## Performance Metrics

The system tracks several performance metrics:

### Time-based Metrics
- **ns/op** - Nanoseconds per operation (primary metric)
- **us/op** - Microseconds per operation
- **ms/op** - Milliseconds per operation

### Memory Metrics
- **B/op** - Bytes allocated per operation
- **allocs/op** - Memory allocations per operation

### Throughput Metrics
- **op/s** - Operations per second
- **MB/s** - Megabytes processed per second

## Integration with Development Workflow

### Pre-commit Hooks

Consider adding performance checks to your pre-commit workflow:

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

## Advanced Usage

### Custom Benchmark Suites

Create custom benchmark suites for specific use cases:

```go
// test/custom_benchmark_test.go
func BenchmarkCustomOperation(b *testing.B) {
    // Your custom benchmark
}
```

### Continuous Monitoring

Set up continuous performance monitoring:

```bash
# Cron job for daily monitoring
0 2 * * * cd /path/to/tsdb && make performance-monitor
```

### Performance Alerts

Integrate with alerting systems:

```bash
# Check for critical regressions and send alerts
if ./scripts/detect-regressions.sh --critical-threshold 10.0; then
    # Send alert (email, Slack, etc.)
    echo "Critical performance regression detected!" | mail -s "Performance Alert" admin@example.com
fi
```

## Support and Contributing

### Getting Help

- Check the help commands: `make performance-help`
- Review this documentation
- Check GitHub Issues for known problems

### Contributing

To improve the performance monitoring system:

1. **Enhance scripts** with new features
2. **Add new benchmark types** for uncovered areas
3. **Improve visualization** in the dashboard
4. **Add new CI/CD integrations**

### Reporting Issues

When reporting performance monitoring issues:

1. Include the command that failed
2. Attach relevant log files
3. Specify your Go version and OS
4. Include benchmark results if available

---

For more information, see the main [README.md](../README.md) or run `make performance-help`.
