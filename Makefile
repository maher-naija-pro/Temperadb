# Makefile for TimeSeriesDB - Test and Coverage only

# Go command
GOCMD=go

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

.PHONY: benchmark-parser
benchmark-parser:
	@echo "Running parser benchmarks..."
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
	@echo "1. Parser Performance..."
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

.PHONY: benchmark-help
benchmark-help:
	@echo "Available benchmark targets:"
	@echo "  benchmark         - Run all benchmarks"
	@echo "  benchmark-parser  - Parser performance only"
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
	@echo "  go test -bench=BenchmarkParse -benchmem ./test/       # Parser only"
	@echo "  go test -bench=BenchmarkWrite -benchmem ./test/       # Storage only"
	@echo "  go test -bench=BenchmarkHTTP -benchmem ./test/        # HTTP only"

# Default target
.DEFAULT_GOAL := test
