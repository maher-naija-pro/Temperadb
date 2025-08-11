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

# Linting and code quality
lint:
	@echo "Running comprehensive linting checks..."
	@./scripts/lint.sh

lint-fast:
	@echo "Running fast linting checks..."
	golangci-lint run --fast ./...

lint-fix:
	@echo "Fixing code formatting issues..."
	go fmt ./...
	goimports -w .
	pre-commit run go-fmt --all-files || true
	pre-commit run go-imports --all-files || true

# Pre-commit hooks
install-hooks:
	@echo "Installing pre-commit hooks..."
	pre-commit install --install-hooks

run-hooks:
	@echo "Running pre-commit hooks on all files..."
	pre-commit run --all-files

run-hooks-staged:
	@echo "Running pre-commit hooks on staged files..."
	pre-commit run

update-hooks:
	@echo "Updating pre-commit hooks..."
	pre-commit autoupdate

# Security checks
security:
	@echo "Running security vulnerability scan..."
	gosec ./...

security-html:
	@echo "Generating HTML security report..."
	gosec -fmt=html -out=security-report.html ./...

# Code complexity analysis
complexity:
	@echo "Analyzing code complexity..."
	find . -name "*.go" -not -path "./test/*" -not -path "./vendor/*" | xargs gocyclo -over 15

# Spell checking
spell:
	@echo "Checking for spelling mistakes..."
	misspell -error .

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install golang.org/x/tools/cmd/goimports@latest

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
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Run comprehensive linting checks"
	@echo "  lint-fast      - Run fast linting checks"
	@echo "  lint-fix       - Fix code formatting issues"
	@echo "  security       - Run security vulnerability scan"
	@echo "  complexity     - Analyze code complexity"
	@echo "  spell          - Check for spelling mistakes"
	@echo ""
	@echo "Pre-commit Hooks:"
	@echo "  install-hooks  - Install pre-commit hooks"
	@echo "  run-hooks      - Run pre-commit hooks on all files"
	@echo "  run-hooks-staged- Run pre-commit hooks on staged files"
	@echo "  update-hooks   - Update pre-commit hooks"
	@echo ""
	@echo "Development:"
	@echo "  install-tools  - Install development tools"
	@echo "  help           - Show this help message"
