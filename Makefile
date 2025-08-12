# Makefile for TimeSeriesDB - Build, Test and Coverage

# Go command
GOCMD=go

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build targets
.PHONY: build
build:
	@echo "Building TimeSeriesDB version $(VERSION)..."
	$(GOCMD) build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)" -o timeseriesdb .

.PHONY: build-linux
build-linux:
	@echo "Building TimeSeriesDB for Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)" -o timeseriesdb-linux-amd64 .

.PHONY: build-windows
build-windows:
	@echo "Building TimeSeriesDB for Windows AMD64..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)" -o timeseriesdb-windows-amd64.exe .

.PHONY: build-darwin
build-darwin:
	@echo "Building TimeSeriesDB for macOS AMD64..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)" -o timeseriesdb-darwin-amd64 .

.PHONY: build-all
build-all: build-linux build-windows build-darwin
	@echo "All builds completed!"

.PHONY: build-docker
build-docker:
	@echo "Building Docker image..."
	docker build --build-arg VERSION=$(VERSION) -t timeseriesdb:$(VERSION) .

.PHONY: clean-build
clean-build:
	@echo "Cleaning build artifacts..."
	rm -f timeseriesdb timeseriesdb-* *.tar.gz *.zip
	@echo "Build artifacts cleaned."

# Test targets
.PHONY: test
test:
	$(GOCMD) test -v ./...

# Test with coverage
.PHONY: coverage
coverage:
	$(GOCMD) test -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean coverage files
.PHONY: clean
clean:
	rm -f coverage.out coverage.html

# Benchmark targets
.PHONY: benchmark
benchmark:
	@echo "Running all benchmarks..."
	$(GOCMD) test -bench=. -benchmem -timeout=5m ./test/

.PHONY: benchmark-ingestion
benchmark-ingestion:
	@echo "Running ingestion benchmarks..."
	$(GOCMD) test -bench=BenchmarkParse -benchmem -timeout=5m ./test/

.PHONY: benchmark-storage
benchmark-storage:
	@echo "Running storage benchmarks..."
	$(GOCMD) test -bench=BenchmarkWrite -benchmem -timeout=5m ./test/

.PHONY: benchmark-http
benchmark-http:
	@echo "Running HTTP endpoint benchmarks..."
	$(GOCMD) test -bench=BenchmarkHTTP -benchmem -timeout=5m ./test/

.PHONY: benchmark-e2e
benchmark-e2e:
	@echo "Running end-to-end workflow benchmarks..."
	$(GOCMD) test -bench="BenchmarkEndToEnd|BenchmarkConcurrent" -benchmem -timeout=5m ./test/

.PHONY: benchmark-memory
benchmark-memory:
	@echo "Running memory usage benchmarks..."
	$(GOCMD) test -bench=BenchmarkMemory -benchmem -timeout=5m ./test/

.PHONY: benchmark-all
benchmark-all:
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
benchmark-profile:
	@echo "Running benchmarks with CPU and memory profiling..."
	@echo "CPU profiling..."
	$(GOCMD) test -bench=BenchmarkParseLargeDataset -cpuprofile=cpu_profile.prof -benchmem -timeout=5m ./test/
	@echo "Memory profiling..."
	$(GOCMD) test -bench=BenchmarkMemoryUsage -memprofile=memory_profile.prof -benchmem -timeout=5m ./test/
	@echo ""
	@echo "Profiles generated:"
	@echo "  CPU: cpu_profile.prof"
	@echo "  Memory: memory_profile.prof"
	@echo ""
	@echo "To analyze profiles:"
	@echo "  go tool pprof cpu_profile.prof"
	@echo "  go tool pprof memory_profile.prof"

.PHONY: benchmark-clean
benchmark-clean:
	@echo "Cleaning benchmark artifacts..."
	rm -f cpu_profile.prof memory_profile.prof block_profile.prof
	@echo "Benchmark artifacts cleaned."

# Performance Regression Detection targets
.PHONY: regression-detect
regression-detect:
	@echo "Running performance regression detection..."
	@echo "Make sure you have a baseline set first:"
	@echo "  ./scripts/run-benchmarks.sh -b"
	@echo ""
	./scripts/detect-regressions.sh

.PHONY: regression-detect-html
regression-detect-html:
	@echo "Running performance regression detection with HTML report..."
	./scripts/detect-regressions.sh -H

.PHONY: regression-detect-json
regression-detect-json:
	@echo "Running performance regression detection with JSON output..."
	./scripts/detect-regressions.sh -j

.PHONY: regression-detect-full
regression-detect-full:
	@echo "Running performance regression detection with all outputs..."
	./scripts/detect-regressions.sh -H -j

.PHONY: regression-baseline
regression-baseline:
	@echo "Setting current benchmark results as baseline..."
	./scripts/run-benchmarks.sh -b

.PHONY: regression-compare
regression-compare:
	@echo "Comparing current results with baseline..."
	./scripts/run-benchmarks.sh -c

# Performance Dashboard targets
.PHONY: dashboard
dashboard:
	@echo "Generating performance dashboard..."
	./scripts/performance-dashboard.sh -g

.PHONY: dashboard-trends
dashboard-trends:
	@echo "Generating performance trends analysis..."
	./scripts/performance-dashboard.sh -t

.PHONY: dashboard-summary
dashboard-summary:
	@echo "Generating performance summary report..."
	./scripts/performance-dashboard.sh -s

.PHONY: dashboard-open
dashboard-open:
	@echo "Opening performance dashboard in browser..."
	@if [ -f "performance-dashboard/index.html" ]; then \
		xdg-open performance-dashboard/index.html 2>/dev/null || \
		open performance-dashboard/index.html 2>/dev/null || \
		echo "Please open performance-dashboard/index.html manually in your browser"; \
	else \
		echo "Dashboard not found. Generate it first with: make dashboard"; \
	fi

# Performance Monitoring Workflow
.PHONY: performance-monitor
performance-monitor:
	@echo "========================================="
	@echo "  TimeSeriesDB Performance Monitoring   "
	@echo "========================================="
	@echo ""
	@echo "1. Running benchmarks..."
	./scripts/run-benchmarks.sh -a
	@echo ""
	@echo "2. Detecting performance regressions..."
	./scripts/detect-regressions.sh -H -j
	@echo ""
	@echo "3. Generating performance dashboard..."
	./scripts/performance-dashboard.sh -g
	@echo ""
	@echo "========================================="
	@echo "  Performance monitoring completed!     "
	@echo "========================================="
	@echo ""
	@echo "Next steps:"
	@echo "  - View regression reports: benchmark-results/regression_report_*.txt"
	@echo "  - Open dashboard: make dashboard-open"
	@echo "  - Set baseline: make regression-baseline"

.PHONY: performance-clean
performance-clean:
	@echo "Cleaning performance monitoring artifacts..."
	rm -rf performance-dashboard/
	rm -f benchmark-results/regression_report_*.txt
	rm -f benchmark-results/regression_report_*.html
	rm -f benchmark-results/regression_report_*.json
	@echo "Performance monitoring artifacts cleaned."

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

# Build help
.PHONY: build-help
build-help:
	@echo "Available build targets:"
	@echo "  build         - Build for current platform"
	@echo "  build-linux   - Build for Linux AMD64"
	@echo "  build-windows - Build for Windows AMD64"
	@echo "  build-darwin  - Build for macOS AMD64"
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
.DEFAULT_GOAL := test
