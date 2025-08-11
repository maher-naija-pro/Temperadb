.PHONY: test test-verbose test-coverage test-race benchmark clean help

# Default target
all: test

# Run all tests
test:
	go test ./test/...

# Run tests with verbose output
test-verbose:
	go test -v ./test/...

# Run tests with race detection
test-race:
	go test -race ./test/...

# Run tests and generate coverage report
test-coverage:
	go test -coverprofile=coverage.out ./test/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./test/...

# Run specific benchmark
benchmark-write:
	go test -bench=BenchmarkWriteEndpoint -benchmem ./test/...

# Clean up test artifacts
clean:
	rm -f coverage.out coverage.html
	rm -f test_data.tsv benchmark_data.tsv
	go clean -testcache

# Install test dependencies
deps:
	go mod tidy
	go get github.com/stretchr/testify

# Show help
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-coverage  - Run tests and generate coverage report"
	@echo "  benchmark      - Run all benchmarks"
	@echo "  benchmark-write- Run write endpoint benchmark"
	@echo "  clean          - Clean up test artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  help           - Show this help message"
