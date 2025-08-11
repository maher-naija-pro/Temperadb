# Continuous Integration & Deployment

This document describes the CI/CD pipeline for the TimeSeriesDB project, which automatically runs tests, generates reports, and ensures code quality on every push and pull request.

## 🚀 Overview

The CI pipeline consists of multiple workflows that run in parallel to provide comprehensive validation of your code:

- **Test Suite**: Runs all tests with coverage reporting
- **Benchmarks**: Performance testing and benchmarking
- **Code Quality**: Linting and static analysis
- **Security**: Security vulnerability scanning
- **Build Verification**: Multi-platform build testing
- **Badge Generation**: Automatic status badge updates

## 📋 Workflows

### 1. Test Suite (`test.yml`)

**Triggers**: Push to `main`/`develop`, Pull Requests  
**Runs on**: Ubuntu Latest  
**Go Versions**: 1.20, 1.21, 1.22  

**What it does:**
- Runs all tests with race detection
- Generates coverage reports (HTML + text)
- Enforces 80% minimum coverage threshold
- Uploads coverage artifacts for each Go version

**Outputs:**
- Coverage reports as downloadable artifacts
- Test results in GitHub Actions logs
- Coverage threshold validation

### 2. Benchmarks (`test.yml`)

**Triggers**: Push to `main`/`develop`, Pull Requests  
**Runs on**: Ubuntu Latest  
**Go Version**: 1.22  

**What it does:**
- Runs HTTP endpoint benchmarks
- Tests bulk write performance
- Benchmarks request creation
- Generates comprehensive benchmark reports

**Benchmarks included:**
- `BenchmarkWriteEndpoint`: Single write operations
- `BenchmarkBulkWriteEndpoint`: Bulk write operations
- `BenchmarkHTTPRequestCreation`: Request creation overhead

**Outputs:**
- Benchmark reports as markdown files
- Performance metrics in GitHub Actions logs
- Downloadable benchmark artifacts

### 3. Code Quality (`test.yml`)

**Triggers**: Push to `main`/`develop`, Pull Requests  
**Runs on**: Ubuntu Latest  
**Go Version**: 1.22  

**What it does:**
- Runs golangci-lint with comprehensive rules
- Checks code formatting with `go fmt`
- Runs `go vet` for common issues
- Enforces consistent code style

**Linters enabled:**
- **Style**: gofmt, goimports, whitespace
- **Correctness**: govet, errcheck, staticcheck
- **Performance**: gosimple, ineffassign, prealloc
- **Security**: gosec
- **Complexity**: gocyclo, dupl, goconst

**Outputs:**
- Linting results in GitHub Actions logs
- Code quality validation
- Formatting checks

### 4. Security Scan (`test.yml`)

**Triggers**: Push to `main`/`develop`, Pull Requests  
**Runs on**: Ubuntu Latest  
**Go Version**: 1.22  

**What it does:**
- Installs and runs gosec security scanner
- Scans for common security vulnerabilities
- Generates JSON security reports
- Uploads security findings as artifacts

**Security checks:**
- SQL injection vulnerabilities
- Command injection risks
- Insecure cryptographic practices
- File operation security
- Network security issues

**Outputs:**
- Security scan reports as JSON
- Vulnerability findings in logs
- Downloadable security artifacts

### 5. Build Verification (`test.yml`)

**Triggers**: Push to `main`/`develop`, Pull Requests  
**Runs on**: Ubuntu, Windows, macOS  
**Go Version**: 1.22  

**What it does:**
- Builds the application on multiple platforms
- Verifies build artifacts are created
- Tests cross-platform compatibility
- Uploads build artifacts for each OS

**Platforms tested:**
- Ubuntu Latest
- Windows Latest
- macOS Latest

**Outputs:**
- Build artifacts for each platform
- Build success/failure validation
- Cross-platform compatibility verification

### 6. Badge Generation (`badges.yml`)

**Triggers**: After Test Suite completion  
**Runs on**: Ubuntu Latest  
**Dependencies**: Test Suite workflow  

**What it does:**
- Downloads coverage reports
- Generates dynamic badges
- Updates README.md with current status
- Commits and pushes badge updates

**Badges generated:**
- Test coverage percentage
- Build status (passing/failing)
- Dynamic color coding based on metrics

## 🔧 Configuration

### Go Version Support

The pipeline supports multiple Go versions to ensure compatibility:
- **1.20**: Minimum supported version
- **1.21**: Stable version
- **1.22**: Latest stable version

### Coverage Thresholds

- **Minimum**: 80% (enforced)
- **Target**: 90%+ (recommended)
- **Badge Colors**:
  - 🟢 90%+: Bright Green
  - 🟢 80-89%: Green
  - 🟡 70-79%: Yellow
  - 🟠 60-69%: Orange
  - 🔴 <60%: Red

### Linter Configuration

The `.golangci.yml` file configures:
- **Timeout**: 5 minutes per linter run
- **Complexity**: Maximum cyclomatic complexity of 15
- **Duplication**: 100-line threshold for duplicate detection
- **Magic Numbers**: Checks for hardcoded values
- **Test Files**: Relaxed rules for test code

## 📊 Artifacts

### Coverage Reports

- **Text Format**: `coverage.out` (for CI tools)
- **HTML Format**: `coverage.html` (for human review)
- **Badge Updates**: Automatic README updates

### Benchmark Reports

- **Markdown Format**: `benchmark-report.md`
- **Performance Metrics**: Operations per second, memory usage
- **Comparison Data**: For performance regression detection

### Security Reports

- **JSON Format**: `security-report.json`
- **Vulnerability Details**: Severity, location, recommendations
- **Actionable Items**: Specific fixes for security issues

### Build Artifacts

- **Executables**: Platform-specific binaries
- **Build Logs**: Compilation and linking details
- **Cross-Platform**: Verification of multi-OS support

## 🚦 Status Checks

### Required for Merge

The following checks must pass for PRs to be merged:
- ✅ All tests passing
- ✅ Coverage above 80%
- ✅ No linting errors
- ✅ Security scan clean
- ✅ Builds on all platforms

### Optional Checks

These provide additional information but don't block merges:
- 📊 Benchmark performance
- 🔍 Detailed coverage analysis
- 🛡️ Security vulnerability details

## 🛠️ Local Development

### Running CI Checks Locally

```bash
# Run all tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Run linter
golangci-lint run ./...

# Run security scan
gosec ./...

# Check formatting
go fmt ./...
go vet ./...
```

### Pre-commit Hooks

Consider installing pre-commit hooks to catch issues early:

```bash
# Install pre-commit
pip install pre-commit

# Install git hooks
pre-commit install

# Run on all files
pre-commit run --all-files
```

## 📈 Performance Monitoring

### Benchmark Tracking

The pipeline automatically tracks:
- **HTTP Endpoint Performance**: Response times and throughput
- **Bulk Operations**: Large dataset handling
- **Memory Usage**: Allocation patterns and efficiency
- **CPU Usage**: Processing overhead

### Regression Detection

- **Historical Comparison**: Track performance over time
- **Threshold Alerts**: Flag significant performance drops
- **Trend Analysis**: Identify performance patterns

## 🔍 Troubleshooting

### Common Issues

1. **Coverage Below Threshold**
   - Add tests for uncovered code paths
   - Review test data generation
   - Check for conditional branches

2. **Linting Failures**
   - Run `golangci-lint run ./...` locally
   - Fix formatting with `go fmt ./...`
   - Address complexity issues

3. **Security Vulnerabilities**
   - Review gosec findings
   - Update dependencies if needed
   - Implement secure coding practices

4. **Build Failures**
   - Check platform-specific code
   - Verify dependency compatibility
   - Review build constraints

### Debug Mode

Enable verbose logging in GitHub Actions:
```yaml
- name: Run tests
  run: |
    go test -v -race -coverprofile=coverage.out ./test/...
  env:
    GITHUB_ACTIONS: true
    DEBUG: true
```

## 📚 Resources

### Documentation
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Go Testing](https://golang.org/pkg/testing/)
- [golangci-lint](https://golangci-lint.run/)
- [gosec Security Scanner](https://github.com/securecodewarrior/gosec)

### Best Practices
- Write tests before implementing features (TDD)
- Maintain high test coverage
- Use meaningful test names and descriptions
- Keep tests fast and focused
- Mock external dependencies appropriately

### Performance Tips
- Use table-driven tests for multiple scenarios
- Benchmark critical code paths
- Profile memory usage in tests
- Use appropriate test timeouts

## 🤝 Contributing

When contributing to the CI pipeline:

1. **Test Locally**: Run all checks before pushing
2. **Update Documentation**: Keep this file current
3. **Add New Checks**: Extend the pipeline as needed
4. **Monitor Performance**: Ensure CI doesn't slow down
5. **Security First**: Prioritize security scanning
6. **Cross-Platform**: Test on multiple operating systems

The CI pipeline is designed to catch issues early and ensure code quality. Regular monitoring and updates will keep it effective and efficient.
