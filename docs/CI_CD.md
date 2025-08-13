# TimeSeriesDB CI/CD & Benchmark Systems - Integrated Guide

This document provides a comprehensive overview of the integrated CI/CD and Benchmark systems for TimeSeriesDB, combining automated build processes with performance testing and monitoring.

## Overview

The TimeSeriesDB development pipeline integrates three core systems:

1. **CI/CD Pipeline** - Automated building, testing, and deployment
2. **Benchmark CI System** - Performance testing and regression detection
3. **Auto PR Creator** - Automatic pull request creation for test branches

Together, these systems ensure code quality, performance consistency, and reliable releases across multiple platforms and Go versions.

## Integrated Workflows

### 1. Build and Package Pipeline (`.github/workflows/build-packages.yml`)

**Triggers:**
- Push to `main` or `master` branch
- Push of a tag (e.g., `v1.0.0`)
- Pull request to `main` or `master`
- Manual workflow dispatch

**Jobs:**

#### Build Job (Matrix Strategy)
- **Multi-platform builds**:
  - Linux AMD64 & ARM64
  - Windows AMD64
  - macOS AMD64 & ARM64 (Apple Silicon)
- **Artifacts**: Platform-specific binaries and archives
- **Docker**: Linux platform images

#### Release Job
- **Condition**: Only runs on tag pushes
- **Action**: Creates GitHub release with all built artifacts
- **Files**: `.tar.gz` and `.zip` files for each platform

#### Publish Docker Job
- **Condition**: Only runs on tag pushes
- **Action**: Publishes to GitHub Container Registry
- **Tags**: Versioned and `latest` tags

#### Test Binaries Job
- **Condition**: Only runs on tag pushes
- **Action**: Verifies built binaries functionality

#### Security Scan Job
- **Condition**: Only runs on tag pushes
- **Action**: Trivy vulnerability scanner
- **Output**: GitHub Security tab integration

### 2. Benchmark Pipeline (`.github/workflows/benchmark.yml`)

**Triggers:**
- Push/Pull Request to main/master
- Daily scheduled runs (2 AM UTC)
- Manual workflow dispatch

**Jobs:**

#### Benchmark Job (Matrix Strategy)
- **Go version compatibility**: 1.20, 1.21, 1.22
- **Performance consistency**: Ensures compatibility across versions
- **Artifacts**: CPU/memory profiles, benchmark results

#### Benchmark Performance Job
- **Post-benchmark analysis**: Detailed performance metrics
- **PR integration**: Automatic comments with benchmark summaries
- **Regression detection**: Performance change analysis
- **Long-term storage**: 90-day artifact retention

### 3. Auto PR Creator Pipeline

**Workflows:**
- **`auto-pr-creator.yml`** - Full-featured workflow with detailed logging
- **`auto-pr-creator-fast.yml`** - Optimized for speed with minimal steps

**Triggers:**
- **Automatic**: Pushes to branches matching patterns:
  - `test*` (e.g., `test`, `test-feature`, `test-bugfix`)
  - `test`
  - `feature*` (e.g., `feature/new-query`, `feature/optimization`)
  - `dev*` (e.g., `dev`, `dev-experimental`)
- **Manual**: Can be triggered manually via GitHub Actions

**Process:**
1. **Check**: Verifies if a PR already exists for the branch
2. **Create**: If no PR exists, creates a draft PR with:
   - Title: "Auto PR: {branch} → main"
   - Draft status (requires manual review)
   - Base: `main`
   - Head: current branch

**Features:**
- ✅ **Smart**: Only creates PRs when they don't exist
- ✅ **Fast**: Minimal checkout depth and efficient checks
- ✅ **Safe**: Creates draft PRs requiring manual review
- ✅ **Flexible**: Works with any test/feature/dev branch pattern

## Local Development Integration

### Prerequisites
- Go 1.20 or later
- Docker (optional, for Docker builds)
- Make

### Integrated Development Commands

```bash
# Build and test cycle
make build              # Build current platform
make test              # Run unit tests
make benchmark         # Run performance tests
make clean-build       # Clean artifacts

# Platform-specific builds
make build-linux       # Linux AMD64
make build-windows     # Windows AMD64
make build-darwin      # macOS AMD64
make build-all         # All platforms

# Docker integration
make build-docker      # Docker image build
```

### Benchmark Script Integration

```bash
# Run all benchmarks
./scripts/run-benchmarks.sh

# Category-specific testing
./scripts/run-benchmarks.sh -p  # Ingestion only
./scripts/run-benchmarks.sh -s  # Storage only
./scripts/run-benchmarks.sh -e  # HTTP endpoints only
./scripts/run-benchmarks.sh -t  # End-to-end workflows only
./scripts/run-benchmarks.sh -m  # Memory usage only

# Performance comparison
./scripts/run-benchmarks.sh -c  # Compare with baseline
./scripts/run-benchmarks.sh -b  # Set current as baseline
```

## Benchmark Categories

### 1. Ingestion Performance
- Simple line parsing
- Complex line parsing (many tags/fields)
- Multi-line batch processing
- Large dataset scalability (1-10,000 lines)
- Tag and field count variations

### 2. Storage Performance
- Single point writes
- Batch point operations
- Many fields/tags scenarios

### 3. HTTP Endpoint Performance
- Single and batch writes
- Large dataset operations

### 4. End-to-End Workflows
- Complete parse → store → retrieve cycles
- Concurrent operation performance

### 5. Memory Usage Analysis
- Allocation patterns
- Garbage collection pressure

## Quality Assurance Flow

```
Code Change → Local Development → Local Testing → Local Benchmarks
                ↓
        Push to Repository → CI/CD Pipeline Activation
                ↓
        Parallel Execution:
        ├── Build & Package Pipeline
        │   ├── Multi-platform builds
        │   ├── Docker image creation
        │   ├── Security scanning
        │   └── Binary verification
        ├── Benchmark Pipeline
        │   ├── Multi-version testing
        │   ├── Performance analysis
        │   ├── Regression detection
        │   └── PR integration
        └── Auto PR Creator
            ├── Branch pattern matching
            ├── PR existence checking
            └── Draft PR creation
                ↓
        Quality Gates → Release (if tag) → Artifact Storage
                ↓
        Continuous Monitoring → Trend Analysis → Alerting
```

## Performance Analysis Integration

### Benchmark Output Format
```
BenchmarkParseSimpleLine-8         1000000              1234 ns/op             256 B/op          8 allocs/op
BenchmarkParseComplexLine-8          500000              2468 ns/op             512 B/op         16 allocs/op
```

**Metrics:**
- **Time per Operation**: Nanoseconds per operation
- **Memory per Operation**: Bytes allocated per operation
- **Allocations per Operation**: Number of allocations per operation

### Profiling Integration
- **CPU Profiles**: Execution time analysis
- **Memory Profiles**: Allocation pattern analysis
- **Web Interface**: `go tool pprof -http=:8080`

## Continuous Monitoring

### Daily Performance Tracking
- Scheduled benchmark execution
- Long-term trend analysis
- Performance regression detection
- Cross-version compatibility monitoring

### PR Integration
- Automatic benchmark execution
- Performance comparison with baseline
- Regression flagging
- Performance summary comments

## Security and Compliance

### Automated Security
- Trivy vulnerability scanning
- GitHub Security tab integration
- Release artifact verification
- Multi-platform binary validation

### Quality Gates
- All tests must pass
- No significant performance regressions
- Security scans must pass
- Binaries must be functional

## Artifact Management

### Storage Strategy
- **Benchmark Results**: 90-day retention
- **Build Artifacts**: Platform-specific binaries
- **Docker Images**: Versioned and latest tags
- **Security Reports**: Vulnerability scan results

### Access Patterns
- GitHub Releases for binaries
- GitHub Container Registry for Docker images
- GitHub Actions artifacts for CI data
- GitHub Security tab for vulnerabilities

## Auto PR Creator Usage

### Automatic PR Creation
Simply push to any matching branch:
```bash
git checkout -b test-new-feature
git push origin test-new-feature
# PR will be created automatically
```

### Manual Trigger
1. Go to GitHub Actions
2. Select "Auto PR Creator" workflow
3. Click "Run workflow"

### Configuration

#### Branch Patterns
Edit the workflow files to modify which branches trigger PR creation:
```yaml
branches:
  - 'test*'      # All branches starting with 'test'
  - 'feature*'   # All branches starting with 'feature'
  - 'dev*'       # All branches starting with 'dev'
```

#### PR Settings
- **Draft**: All PRs are created as drafts
- **Base**: Always targets `main` branch
- **Title**: Auto-generated with branch name
- **Body**: Simple template with checklist

## Troubleshooting

### Common Issues

**Build Failures:**
- Check Go version compatibility
- Verify platform-specific requirements
- Review Docker configuration

**Benchmark Issues:**
- Monitor timeout values
- Check memory usage patterns
- Verify benchmark isolation

**CI/CD Problems:**
- Review workflow permissions
- Check GitHub repository settings
- Verify environment variables

**Auto PR Creator Issues:**
- Check if branch name matches patterns
- Verify GitHub Actions are enabled
- Check workflow run logs for errors
- If duplicates occur, check branch naming conflicts

### Debug Strategies

**Local Debugging:**
```bash
# Build verification
./scripts/build-test.sh

# Benchmark isolation
./scripts/run-benchmarks.sh --timeout 5m

# Profile analysis
go tool pprof cpu_profile.prof
```

**CI Debugging:**
- Enable `ACTIONS_STEP_DEBUG` secret
- Review workflow logs
- Check artifact generation

**Performance Optimization:**
- Use `auto-pr-creator-fast.yml` for maximum speed
- Minimal checkout depth reduces execution time

## Best Practices

### Development Workflow
1. **Local Testing**: Always test builds and benchmarks locally
2. **Performance Baseline**: Establish and maintain performance baselines
3. **Regular Monitoring**: Watch for performance trends and regressions
4. **Security Updates**: Regularly review and update dependencies
5. **Branch Naming**: Use consistent patterns for test/feature branches

### CI/CD Management
1. **Consistent Environments**: Use standardized Go versions and OS
2. **Artifact Retention**: Maintain historical data for analysis
3. **Failure Handling**: Ensure non-blocking benchmark execution
4. **Documentation**: Keep workflows and processes documented
5. **PR Management**: Review auto-created PRs promptly

## Future Enhancements

### Planned Integrations
1. **Performance Regression Alerts**: Automated notification system
2. **Historical Dashboard**: Web-based performance visualization
3. **Performance Budgets**: Automated constraint enforcement
4. **Cross-Platform Benchmarking**: Multi-OS performance testing
5. **Hardware Profiling**: Real-time resource monitoring

### System Evolution
1. **AI-Powered Analysis**: Machine learning for regression detection
2. **Predictive Monitoring**: Performance trend forecasting
3. **Integration APIs**: External system connectivity
4. **Advanced Metrics**: Custom performance indicators
5. **Smart PR Management**: AI-powered PR review and suggestions

## Support and Resources

### Documentation
- [Performance Guide](PERFORMANCE.md)
- [Performance Monitoring](PERFORMANCE_MONITORING.md)
- [Installation Guide](INSTALLATION.md)

### Issue Reporting
- GitHub Issues for CI/CD problems
- Performance regression reports
- Security vulnerability reports

### Community
- Regular performance reviews
- Benchmark result sharing
- Best practice discussions

---

This integrated system ensures TimeSeriesDB maintains high quality standards across code, performance, and security while providing developers with comprehensive tools for local development and CI/CD integration.
