# Makefile for TimeSeriesDB - Build, Test and Coverage

# Go command and version
GOCMD=go
GOVERSION=$(shell go version | awk '{print $$3}')

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(VERSION)

# Directories
BIN_DIR=bin
DIST_DIR=dist
COVERAGE_DIR=coverage
BENCHMARK_DIR=benchmark-results
DASHBOARD_DIR=performance-dashboard

# Ensure directories exist
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(COVERAGE_DIR):
	mkdir -p $(COVERAGE_DIR)

# Build targets
.PHONY: build
build: $(BIN_DIR)
	@echo "Building TimeSeriesDB version $(VERSION)..."
	@echo "Go version: $(GOVERSION)"
	$(GOCMD) build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/timeseriesdb .
	@echo "Build completed: $(BIN_DIR)/timeseriesdb"

.PHONY: build-linux
build-linux: $(DIST_DIR)
	@echo "Building TimeSeriesDB for Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOCMD) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/timeseriesdb-linux-amd64 .
	@echo "Linux build completed: $(DIST_DIR)/timeseriesdb-linux-amd64"


.PHONY: build-all
build-all: build-linux build-windows build-darwin build-darwin-arm64
	@echo "All builds completed!"
	@echo "Build artifacts in $(DIST_DIR)/"

.PHONY: build-docker
build-docker:
	@echo "Building Docker image..."
	docker build --build-arg VERSION=$(VERSION) -t timeseriesdb:$(VERSION) .
	@echo "Docker image built: timeseriesdb:$(VERSION)"



.PHONY: clean-build
clean-build:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)/* $(DIST_DIR)/* *.tar.gz *.zip
	@echo "Build artifacts cleaned."

# Dependency management
.PHONY: deps
deps:
	@echo "Downloading Go dependencies..."
	$(GOCMD) mod download
	@echo "Dependencies downloaded"

.PHONY: deps-update
deps-update:
	@echo "Updating Go dependencies..."
	$(GOCMD) get -u ./...
	$(GOCMD) mod tidy
	@echo "Dependencies updated and tidied"

.PHONY: deps-check
deps-check:
	@echo "Checking Go dependencies..."
	$(GOCMD) mod verify
	$(GOCMD) mod why -m all
	@echo "Dependency check completed"

# Test targets
.PHONY: test
test: $(BIN_DIR)
	@echo "Running tests for all modules..."
	@echo "========================================="
	@echo "  TimeSeriesDB Module Testing           "
	@echo "========================================="
	@echo ""
	@echo "1. Testing internal/errors..."
	$(GOCMD) test -timeout=30s -v ./internal/errors/...
	@echo ""
	@echo "2. Testing internal/ingestion..."
	$(GOCMD) test -timeout=30s -v ./internal/ingestion/...
	@echo ""
	@echo "3. Testing internal/metrics..."
	$(GOCMD) test -timeout=30s -v ./internal/metrics/...
	@echo ""
	@echo "4. Testing internal/server..."
	$(GOCMD) test -timeout=30s -v ./internal/server/...
	@echo ""
	@echo "5. Testing internal/storage..."
	$(GOCMD) test -timeout=30s -v ./internal/storage/...
	@echo ""
	@echo "6. Testing internal/config..."
	$(GOCMD) test -timeout=30s -v ./internal/config/...
	@echo ""
	@echo "7. Testing internal/envvars..."
	$(GOCMD) test -timeout=30s -v ./internal/envvars/...
	@echo ""
	@echo "8. Testing internal/api..."
	$(GOCMD) test -timeout=30s -v ./internal/api/...
	@echo ""
	@echo "9. Testing internal/types..."
	$(GOCMD) test -timeout=30s -v ./internal/types/...
	@echo ""
	@echo "10. Testing internal/logger..."
	$(GOCMD) test -timeout=30s -v ./internal/logger/...
	@echo ""
	@echo "11. Testing test/helpers..."
	$(GOCMD) test -timeout=30s -v ./test/helpers/...
	@echo ""
	@echo "12. Testing test/utils..."
	$(GOCMD) test -timeout=30s -v ./test/utils/...
	@echo ""
	@echo "13. Testing test/integration..."
	$(GOCMD) test -timeout=30s -v ./test/integration/...
	@echo ""
	@echo "14. Testing test/benchmark..."
	$(GOCMD) test -timeout=30s -v ./test/benchmark/...
	@echo ""
	@echo "15. Testing root package..."
	$(GOCMD) test -timeout=30s -v ./
	@echo ""
	@echo "========================================="
	@echo "  All module tests completed!           "
	@echo "========================================="

# Test with coverage
.PHONY: coverage
coverage: $(COVERAGE_DIR)
	@echo "Running tests with coverage..."
	$(GOCMD) test -v -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"



# Test with memory profiling
.PHONY: test-mem
test-mem: $(COVERAGE_DIR)
	@echo "Running tests with memory profiling..."
	$(GOCMD) test -memprofile=$(COVERAGE_DIR)/mem.prof -v ./...
	@echo "Memory profile generated: $(COVERAGE_DIR)/mem.prof"



.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOCMD) test -v -tags=integration ./test/...

# Test individual modules
.PHONY: test-errors
test-errors:
	@echo "Testing internal/errors module..."
	$(GOCMD) test -v ./internal/errors/...

.PHONY: test-ingestion
test-ingestion:
	@echo "Testing internal/ingestion module..."
	$(GOCMD) test -v ./internal/ingestion/...

.PHONY: test-metrics
test-metrics:
	@echo "Testing internal/metrics module..."
	$(GOCMD) test -v ./internal/metrics/...

.PHONY: test-server
test-server:
	@echo "Testing internal/server module..."
	$(GOCMD) test -v ./internal/server/...

.PHONY: test-storage
test-storage:
	@echo "Testing internal/storage module..."
	$(GOCMD) test -v ./internal/storage/...

.PHONY: test-config
test-config:
	@echo "Testing internal/config module..."
	$(GOCMD) test -v ./internal/config/...

.PHONY: test-envvars
test-envvars:
	@echo "Testing internal/envvars module..."
	$(GOCMD) test -v ./internal/envvars/...

.PHONY: test-api
test-api:
	@echo "Testing internal/api module..."
	$(GOCMD) test -v ./internal/api/...
	@echo "Cleaning up test data..."
	rm -rf internal/api/http/test_data_register/
	@echo "Test data cleanup completed"

.PHONY: test-types
test-types:
	@echo "Testing internal/types module..."
	$(GOCMD) test -v ./internal/types/...

.PHONY: test-logger
test-logger:
	@echo "Testing internal/logger module..."
	$(GOCMD) test -v ./internal/logger/...

.PHONY: test-helpers
test-helpers:
	@echo "Testing test/helpers module..."
	$(GOCMD) test -v ./test/helpers/...

.PHONY: test-utils
test-utils:
	@echo "Testing test/utils module..."
	$(GOCMD) test -v ./test/utils/...

.PHONY: test-benchmark
test-benchmark:
	@echo "Testing test/benchmark module..."
	$(GOCMD) test -v ./test/benchmark/...



# Clean coverage files
.PHONY: clean
clean:
	@echo "Cleaning test artifacts..."
	rm -rf $(COVERAGE_DIR)/*
	@echo "Test artifacts cleaned"

# Linting and code quality
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	$(GOCMD) fmt ./...
	@echo "Code formatting completed"

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "Go vet completed"

.PHONY: tidy
tidy:
	@echo "Tidying Go modules..."
	$(GOCMD) mod tidy
	@echo "Go modules tidied"

# Benchmark targets
.PHONY: benchmark
benchmark: $(BENCHMARK_DIR)
	@echo "Running all benchmarks..."
	$(GOCMD) test -bench=. -benchmem -timeout=5m ./test/

.PHONY: benchmark-ingestion
benchmark-ingestion: $(BENCHMARK_DIR)
	@echo "Running ingestion benchmarks..."
	$(GOCMD) test -bench=BenchmarkParse -benchmem -timeout=5m ./test/

.PHONY: benchmark-storage
benchmark-storage: $(BENCHMARK_DIR)
	@echo "Running storage benchmarks..."
	$(GOCMD) test -bench=BenchmarkWrite -benchmem -timeout=5m ./test/

.PHONY: benchmark-http
benchmark-http: $(BENCHMARK_DIR)
	@echo "Running HTTP endpoint benchmarks..."
	$(GOCMD) test -bench=BenchmarkHTTP -benchmem -timeout=5m ./test/

.PHONY: benchmark-e2e
benchmark-e2e: $(BENCHMARK_DIR)
	@echo "Running end-to-end workflow benchmarks..."
	$(GOCMD) test -bench="BenchmarkEndToEnd|BenchmarkConcurrent" -benchmem -timeout=5m ./test/

.PHONY: benchmark-memory
benchmark-memory: $(BENCHMARK_DIR)
	@echo "Running memory usage benchmarks..."
	$(GOCMD) test -bench=BenchmarkMemory -benchmem -timeout=5m ./test/

.PHONY: benchmark-all
benchmark-all: $(BENCHMARK_DIR)
	@echo "========================================="
	@echo "  TimeSeriesDB Performance Benchmarks   "
	@echo "========================================="
	@echo ""
	@echo "1. Ingestion Performance..."
	$(GOCMD) test -bench=BenchmarkParse -benchmem -timeout=5m ./test/
	@echo ""
	@echo "2. Storage Performance..."
	$(GOCMD) test -bench=BenchmarkWrite -benchmem -timeout=5m ./test/
	@echo ""
	@echo "3. HTTP Endpoint Performance..."
	$(GOCMD) test -bench=BenchmarkHTTP -benchmem -timeout=5m ./test/
	@echo ""
	@echo "4. End-to-End Workflow Performance..."
	$(GOCMD) test -bench="BenchmarkEndToEnd|BenchmarkConcurrent" -benchmem -timeout=5m ./test/
	@echo ""
	@echo "5. Memory Usage Performance..."
	$(GOCMD) test -bench=BenchmarkMemory -benchmem -timeout=5m ./test/
	@echo ""
	@echo "========================================="
	@echo "  All benchmarks completed!             "
	@echo "========================================="

.PHONY: benchmark-profile
benchmark-profile: $(BENCHMARK_DIR)
	@echo "Running benchmarks with CPU and memory profiling..."
	@echo "CPU profiling..."
	$(GOCMD) test -bench=BenchmarkParseLargeDataset -cpuprofile=$(BENCHMARK_DIR)/cpu_profile.prof -benchmem -timeout=5m ./test/
	@echo "Memory profiling..."
	$(GOCMD) test -bench=BenchmarkMemoryUsage -memprofile=$(BENCHMARK_DIR)/memory_profile.prof -benchmem -timeout=5m ./test/
	@echo ""
	@echo "Profiles generated:"
	@echo "  CPU: $(BENCHMARK_DIR)/cpu_profile.prof"
	@echo "  Memory: $(BENCHMARK_DIR)/memory_profile.prof"
	@echo ""
	@echo "To analyze profiles:"
	@echo "  go tool pprof $(BENCHMARK_DIR)/cpu_profile.prof"
	@echo "  go tool pprof $(BENCHMARK_DIR)/memory_profile.prof"

.PHONY: benchmark-clean
benchmark-clean:
	@echo "Cleaning benchmark artifacts..."
	rm -f $(BENCHMARK_DIR)/cpu_profile.prof $(BENCHMARK_DIR)/memory_profile.prof $(BENCHMARK_DIR)/block_profile.prof
	@echo "Benchmark artifacts cleaned."

# Performance Regression Detection targets
.PHONY: regression-detect
regression-detect: $(BENCHMARK_DIR)
	@echo "Running performance regression detection..."
	@if [ -f "./scripts/detect-regressions.sh" ]; then \
		./scripts/detect-regressions.sh; \
	else \
		echo "Script not found: ./scripts/detect-regressions.sh"; \
	fi


.PHONY: regression-detect-json
regression-detect-json: $(BENCHMARK_DIR)
	@echo "Running performance regression detection with JSON output..."
	@if [ -f "./scripts/detect-regressions.sh" ]; then \
		./scripts/detect-regressions.sh -j; \
	else \
		echo "Script not found: ./scripts/detect-regressions.sh"; \
	fi

.PHONY: regression-detect-full
regression-detect-full: $(BENCHMARK_DIR)
	@echo "Running performance regression detection with all outputs..."
	@if [ -f "./scripts/detect-regressions.sh" ]; then \
		./scripts/detect-regressions.sh -H -j; \
	else \
		echo "Script not found: ./scripts/detect-regressions.sh"; \
	fi

.PHONY: regression-baseline
regression-baseline: $(BENCHMARK_DIR)
	@echo "Setting current benchmark results as baseline..."
	@if [ -f "./scripts/run-benchmarks.sh" ]; then \
		./scripts/run-benchmarks.sh -b; \
	else \
		echo "Script not found: ./scripts/run-benchmarks.sh"; \
	fi

.PHONY: regression-compare
regression-compare: $(BENCHMARK_DIR)
	@echo "Comparing current results with baseline..."
	@if [ -f "./scripts/run-benchmarks.sh" ]; then \
		./scripts/run-benchmarks.sh -c; \
	else \
		echo "Script not found: ./scripts/run-benchmarks.sh"; \
	fi

# Performance Dashboard targets
.PHONY: dashboard
dashboard: $(DASHBOARD_DIR)
	@echo "Generating performance dashboard..."
	@if [ -f "./scripts/performance-dashboard.sh" ]; then \
		./scripts/performance-dashboard.sh -g; \
	else \
		echo "Script not found: ./scripts/performance-dashboard.sh"; \
	fi

.PHONY: dashboard-trends
dashboard-trends: $(DASHBOARD_DIR)
	@echo "Generating performance trends analysis..."
	@if [ -f "./scripts/performance-dashboard.sh" ]; then \
		./scripts/performance-dashboard.sh -t; \
	else \
		echo "Script not found: ./scripts/performance-dashboard.sh"; \
	fi

.PHONY: dashboard-summary
dashboard-summary: $(DASHBOARD_DIR)
	@echo "Generating performance summary report..."
	@if [ -f "./scripts/performance-dashboard.sh" ]; then \
		./scripts/performance-dashboard.sh -s; \
	else \
		echo "Script not found: ./scripts/performance-dashboard.sh"; \
	fi

.PHONY: dashboard-open
dashboard-open: $(DASHBOARD_DIR)
	@echo "Opening performance dashboard in browser..."
	@if [ -f "$(DASHBOARD_DIR)/index.html" ]; then \
		xdg-open $(DASHBOARD_DIR)/index.html 2>/dev/null || \
		open $(DASHBOARD_DIR)/index.html 2>/dev/null || \
		echo "Please open $(DASHBOARD_DIR)/index.html manually in your browser"; \
	else \
		echo "Dashboard not found. Generate it first with: make dashboard"; \
	fi

# Performance Monitoring Workflow
.PHONY: performance-monitor
performance-monitor: $(BENCHMARK_DIR) $(DASHBOARD_DIR)
	@echo "========================================="
	@echo "  TimeSeriesDB Performance Monitoring   "
	@echo "========================================="
	@echo ""
	@echo "1. Running benchmarks..."
	@if [ -f "./scripts/run-benchmarks.sh" ]; then \
		./scripts/run-benchmarks.sh -a; \
	else \
		echo "Script not found: ./scripts/run-benchmarks.sh"; \
	fi
	@echo ""
	@echo "2. Detecting performance regressions..."
	@if [ -f "./scripts/detect-regressions.sh" ]; then \
		./scripts/detect-regressions.sh -H -j; \
	else \
		echo "Script not found: ./scripts/detect-regressions.sh"; \
	fi
	@echo ""
	@echo "3. Generating performance dashboard..."
	@if [ -f "./scripts/performance-dashboard.sh" ]; then \
		./scripts/performance-dashboard.sh -g; \
	else \
		echo "Script not found: ./scripts/performance-dashboard.sh"; \
	fi
	@echo ""
	@echo "========================================="
	@echo "  Performance monitoring completed!     "
	@echo "========================================="
	@echo ""
	@echo "Next steps:"
	@echo "  - View regression reports: $(BENCHMARK_DIR)/regression_report_*.txt"
	@echo "  - Open dashboard: make dashboard-open"
	@echo "  - Set baseline: make regression-baseline"

.PHONY: performance-clean
performance-clean:
	@echo "Cleaning performance monitoring artifacts..."
	rm -rf $(DASHBOARD_DIR)/
	rm -f $(BENCHMARK_DIR)/regression_report_*.txt
	rm -f $(BENCHMARK_DIR)/regression_report_*.html
	rm -f $(BENCHMARK_DIR)/regression_report_*.json
	@echo "Performance monitoring artifacts cleaned."



.PHONY: pre-commit
pre-commit: fmt vet lint test
	@echo "Pre-commit checks completed"

.PHONY: ci
ci: deps-check lint test coverage
	@echo "CI checks completed"



# Release targets
.PHONY: release
release: clean build-all
	@echo "Creating release packages..."
	@cd $(DIST_DIR) && \
	for file in timeseriesdb-*; do \
		if [[ "$$file" == *".exe" ]]; then \
			zip "$${file%.exe}.zip" "$$file"; \
		else \
			tar -czf "$$file.tar.gz" "$$file"; \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"

.PHONY: clean-all
clean-all: clean clean-build benchmark-clean performance-clean
	@echo "All artifacts cleaned"

# Help targets
.PHONY: help
help:
	@echo "TimeSeriesDB Makefile - Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-linux        - Build for Linux AMD64"
	@echo "  build-windows      - Build for Windows AMD64"
	@echo "  build-darwin       - Build for macOS AMD64"
	@echo "  build-darwin-arm64 - Build for macOS ARM64"
	@echo "  build-all          - Build for all platforms"
	@echo "  build-docker       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-stop        - Stop Docker container"
	@echo "  docker-test        - Test Docker image"
	@echo "  docker-push        - Show push commands for registry"
	@echo "  install            - Install binary"
	@echo "  clean-build        - Clean build artifacts"
	@echo ""
	@echo "Dependency management:"
	@echo "  deps               - Download dependencies"
	@echo "  deps-update        - Update dependencies"
	@echo "  deps-check         - Check dependencies"
	@echo ""
	@echo "Test targets:"
	@echo "  test               - Run all tests for all modules"
	@echo "  test-race          - Run tests with race detection"
	@echo "  test-mem           - Run tests with memory profiling"
	@echo "  test-unit          - Run unit tests only"
	@echo "  test-integration   - Run integration tests only"
	@echo "  coverage           - Run tests with coverage report"
	@echo "  clean              - Clean test artifacts"
	@echo ""
	@echo "Individual module tests:"
	@echo "  test-errors        - Test internal/errors module"
	@echo "  test-ingestion     - Test internal/ingestion module"
	@echo "  test-metrics       - Test internal/metrics module"
	@echo "  test-server        - Test internal/server module"
	@echo "  test-storage       - Test internal/storage module"
	@echo "  test-config        - Test internal/config module"
	@echo "  test-envvars       - Test internal/envvars module"
	@echo "  test-api           - Test internal/api module"
	@echo "  test-types         - Test internal/types module"
	@echo "  test-logger        - Test internal/logger module"
	@echo "  test-helpers       - Test test/helpers module"
	@echo "  test-utils         - Test test/utils module"
	@echo "  test-benchmark     - Test test/benchmark module"
	@echo "  test-integration-module - Test test/integration module"
	@echo "  test-root          - Test root package"
	@echo ""
	@echo "Code quality:"
	@echo "  lint               - Run linter"
	@echo "  fmt                - Format code"
	@echo "  vet                - Run go vet"
	@echo "  tidy               - Tidy Go modules"
	@echo ""
	@echo "Benchmark targets:"
	@echo "  benchmark          - Run all benchmarks"
	@echo "  benchmark-all      - Run all benchmarks with progress"
	@echo "  benchmark-profile  - Run with profiling"
	@echo "  benchmark-clean    - Clean benchmark artifacts"
	@echo ""
	@echo "Performance monitoring:"
	@echo "  performance-monitor - Complete monitoring workflow"
	@echo "  dashboard          - Generate performance dashboard"
	@echo "  regression-detect  - Detect performance regressions"
	@echo ""
	@echo "Development:"
	@echo "  dev-setup          - Setup development environment"
	@echo "  pre-commit         - Run pre-commit checks"
	@echo "  ci                 - Run CI checks"
	@echo ""
	@echo "Other:"
	@echo "  docs               - Generate documentation"
	@echo "  release            - Create release packages"
	@echo "  clean-all          - Clean all artifacts"
	@echo "  benchmark-help     - Show benchmark help message"
	@echo "  test-help          - Show testing help message"
	@echo "  performance-help   - Show performance help message"
	@echo "  build-help         - Show build help message"
	@echo "  help               - Show this help message"

.PHONY: benchmark-help
benchmark-help:
	@echo "Available benchmark targets:"
	@echo "  benchmark         - Run all benchmarks"
	@echo "  benchmark-ingestion  - Ingestion performance only"
	@echo "  benchmark-storage - Storage performance only"
	@echo "  benchmark-http    - HTTP endpoint performance only"
	@echo "  benchmark-e2e     - End-to-end workflow performance only"
	@echo "  benchmark-memory  - Memory usage performance only"
	@echo "  benchmark-all     - Run all benchmarks with progress display"
	@echo "  benchmark-profile - Run with CPU and memory profiling"
	@echo "  benchmark-clean   - Clean up profile files"
	@echo "  benchmark-help    - Show this help message"
	@echo ""
	@echo "Quick commands:"
	@echo "  go test -bench=. -benchmem ./test/                    # All benchmarks"
	@echo "  go test -bench=BenchmarkParse -benchmem ./test/       # Ingestion only"
	@echo "  go test -bench=BenchmarkWrite -benchmem ./test/       # Storage only"
	@echo "  go test -bench=BenchmarkHTTP -benchmem ./test/        # HTTP only"

.PHONY: test-help
test-help:
	@echo "Available testing targets:"
	@echo ""
	@echo "Complete Testing:"
	@echo "  test                  - Run all tests for all modules"
	@echo "  test-mem              - Run tests with memory profiling"
	@echo "  test-unit             - Run unit tests only"
	@echo "  test-integration      - Run integration tests only"
	@echo "  coverage              - Run tests with coverage report"
	@echo ""
	@echo "Individual Module Tests:"
	@echo "  test-errors           - Test internal/errors module"
	@echo "  test-ingestion        - Test internal/ingestion module"
	@echo "  test-metrics          - Test internal/metrics module"
	@echo "  test-server           - Test internal/server module"
	@echo "  test-storage          - Test internal/storage module"
	@echo "  test-config           - Test internal/config module"
	@echo "  test-envvars          - Test internal/envvars module"
	@echo "  test-api              - Test internal/api module"
	@echo "  test-types            - Test internal/types module"
	@echo "  test-logger           - Test internal/logger module"
	@echo "  test-helpers          - Test test/helpers module"
	@echo "  test-utils            - Test test/utils module"
	@echo "  test-benchmark        - Test test/benchmark module"
	@echo "  test-integration-module - Test test/integration module"
	@echo "  test-root             - Test root package"
	@echo ""
	@echo "Utility:"
	@echo "  clean                 - Clean test artifacts"
	@echo "  test-help             - Show this help message"
	@echo ""
	@echo "Quick commands:"
	@echo "  make test             # Test all modules"
	@echo "  make test-storage     # Test storage module only"
	@echo "  make test-api         # Test API module only"
	@echo "  make coverage         # Test with coverage report"

.PHONY: performance-help
performance-help:
	@echo "Available performance monitoring targets:"
	@echo ""
	@echo "Performance Regression Detection:"
	@echo "  regression-detect      - Detect performance regressions"
	@echo "  regression-detect-html - Detect regressions with HTML report"
	@echo "  regression-detect-json - Detect regressions with JSON output"
	@echo "  regression-detect-full - Detect regressions with all outputs"
	@echo "  regression-baseline    - Set current results as baseline"
	@echo "  regression-compare     - Compare with baseline"
	@echo ""
	@echo "Performance Dashboard:"
	@echo "  dashboard              - Generate performance dashboard"
	@echo "  dashboard-trends       - Generate trends analysis"
	@echo "  dashboard-summary      - Generate summary report"
	@echo "  dashboard-open         - Open dashboard in browser"
	@echo ""
	@echo "Complete Workflow:"
	@echo "  performance-monitor    - Run complete monitoring workflow"
	@echo "  performance-clean      - Clean monitoring artifacts"
	@echo "  performance-help       - Show this help message"
	@echo ""
	@echo "Quick workflow:"
	@echo "  make performance-monitor  # Complete monitoring workflow"
	@echo "  make dashboard-open       # View results"

.PHONY: build-help
build-help:
	@echo "Available build targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-linux   - Build for Linux AMD64"
	@echo "  build-windows - Build for Windows AMD64"
	@echo "  build-darwin  - Build for macOS AMD64"
	@echo "  build-darwin-arm64 - Build for macOS ARM64"
	@echo "  build-all     - Build for all platforms"
	@echo "  build-docker  - Build Docker image"
	@echo "  clean-build   - Clean build artifacts"
	@echo ""
	@echo "Build variables:"
	@echo "  VERSION       - Version to build (default: auto-detected)"
	@echo "  BUILD_TIME    - Build timestamp (default: auto-detected)"
	@echo "  COMMIT_HASH   - Git commit hash (default: auto-detected)"
	@echo ""
	@echo "Examples:"
	@echo "  make build VERSION=v1.0.0"
	@echo "  make build-all"
	@echo "  make build-docker VERSION=v1.0.0"

# Default target
.DEFAULT_GOAL := help
