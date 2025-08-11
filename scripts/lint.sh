#!/bin/bash

# TimeSeriesDB Linting Script
# This script runs comprehensive linting checks on the codebase

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install missing tools
install_tool() {
    local tool=$1
    local install_cmd=$2
    
    if ! command_exists "$tool"; then
        print_warning "$tool not found. Installing..."
        eval "$install_cmd"
        if ! command_exists "$tool"; then
            print_error "Failed to install $tool"
            return 1
        fi
        print_success "$tool installed successfully"
    fi
}

# Check and install required tools
print_status "Checking required tools..."

# Install golangci-lint if not present
if ! command_exists golangci-lint; then
    print_warning "golangci-lint not found. Installing..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Install gosec if not present
if ! command_exists gosec; then
    print_warning "gosec not found. Installing..."
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
fi

# Install gocyclo if not present
if ! command_exists gocyclo; then
    print_warning "gocyclo not found. Installing..."
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
fi

# Install misspell if not present
if ! command_exists misspell; then
    print_warning "misspell not found. Installing..."
    go install github.com/client9/misspell/cmd/misspell@latest
fi

print_success "All required tools are available"

# Function to run linting checks
run_lint_check() {
    local name=$1
    local command=$2
    local description=$3
    
    print_status "Running $name: $description"
    
    if eval "$command"; then
        print_success "$name passed"
        return 0
    else
        print_error "$name failed"
        return 1
    fi
}

# Initialize error counter
ERRORS=0

# Go formatting check
print_status "=== Go Code Formatting ==="
if ! run_lint_check "go fmt" "go fmt ./..." "Check code formatting"; then
    ((ERRORS++))
fi

# Go imports check
if ! run_lint_check "go imports" "goimports -w ." "Fix import formatting"; then
    ((ERRORS++))
fi

# Go vet check
if ! run_lint_check "go vet" "go vet ./..." "Check for common Go mistakes"; then
    ((ERRORS++))
fi

# Go mod tidy
if ! run_lint_check "go mod tidy" "go mod tidy" "Clean up go.mod and go.sum"; then
    ((ERRORS++))
fi

# Go mod verify
if ! run_lint_check "go mod verify" "go mod verify" "Verify module dependencies"; then
    ((ERRORS++))
fi

# golangci-lint
print_status "=== Static Analysis ==="
if ! run_lint_check "golangci-lint" "golangci-lint run ./..." "Run comprehensive static analysis"; then
    ((ERRORS++))
fi

# Security scan
print_status "=== Security Analysis ==="
if ! run_lint_check "gosec" "gosec ./..." "Run security vulnerability scan"; then
    ((ERRORS++))
fi

# Code complexity check
print_status "=== Code Complexity ==="
COMPLEX_FUNCTIONS=$(find . -name "*.go" -not -path "./test/*" -not -path "./vendor/*" -not -path "./.git/*" | xargs gocyclo -over 15 2>/dev/null || true)

if [ -n "$COMPLEX_FUNCTIONS" ]; then
    print_warning "Found functions with complexity > 15:"
    echo "$COMPLEX_FUNCTIONS"
    ((ERRORS++))
else
    print_success "No overly complex functions found"
fi

# Duplicate code check
print_status "=== Duplicate Code ==="
DUPLICATE_CODE=$(find . -name "*.go" -not -path "./test/*" -not -path "./vendor/*" -not -path "./.git/*" | xargs gocyclo -over 100 2>/dev/null || true)

if [ -n "$DUPLICATE_CODE" ]; then
    print_warning "Found potential duplicate code:"
    echo "$DUPLICATE_CODE"
    ((ERRORS++))
else
    print_success "No duplicate code detected"
fi

# Spell check
print_status "=== Spell Check ==="
if ! run_lint_check "misspell" "misspell -error ." "Check for spelling mistakes"; then
    ((ERRORS++))
fi

# Test coverage check
print_status "=== Test Coverage ==="
if [ -d "./test" ]; then
    COVERAGE_OUTPUT=$(go test -coverprofile=coverage.out ./test/... 2>&1 || true)
    if echo "$COVERAGE_OUTPUT" | grep -q "FAIL"; then
        print_error "Tests failed during coverage check"
        echo "$COVERAGE_OUTPUT"
        ((ERRORS++))
    else
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        print_status "Test coverage: $COVERAGE%"
        
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            print_error "Test coverage is below 80%"
            ((ERRORS++))
        else
            print_success "Test coverage is above 80%"
        fi
        
        # Clean up coverage file
        rm -f coverage.out
    fi
else
    print_warning "No test directory found, skipping coverage check"
fi

# Documentation check
print_status "=== Documentation Check ==="
REQUIRED_DOCS=("README.md" "TESTING.md" "CI.md")
for doc in "${REQUIRED_DOCS[@]}"; do
    if [ -f "$doc" ]; then
        print_success "$doc exists"
    else
        print_error "$doc is missing"
        ((ERRORS++))
    fi
done

# Makefile check
print_status "=== Makefile Check ==="
if [ -f "Makefile" ]; then
    TARGET_COUNT=$(grep -E "^[a-zA-Z_-]+:" Makefile | wc -l)
    print_status "Found $TARGET_COUNT Makefile targets"
    
    if [ "$TARGET_COUNT" -lt 5 ]; then
        print_warning "Makefile has fewer than 5 targets"
        ((ERRORS++))
    else
        print_success "Makefile has sufficient targets"
    fi
else
    print_error "Makefile not found"
    ((ERRORS++))
fi

# Pre-commit hooks check
print_status "=== Pre-commit Hooks Check ==="
if [ -f ".pre-commit-config.yaml" ]; then
    print_success "Pre-commit configuration found"
    
    if command_exists pre-commit; then
        print_status "Installing pre-commit hooks..."
        if pre-commit install --install-hooks; then
            print_success "Pre-commit hooks installed"
        else
            print_error "Failed to install pre-commit hooks"
            ((ERRORS++))
        fi
    else
        print_warning "pre-commit not installed. Install with: pip install pre-commit"
    fi
else
    print_warning "Pre-commit configuration not found"
fi

# Summary
print_status "=== Linting Summary ==="
if [ $ERRORS -eq 0 ]; then
    print_success "All linting checks passed! 🎉"
    exit 0
else
    print_error "$ERRORS linting check(s) failed ❌"
    print_status "Please fix the issues above and run the script again"
    exit 1
fi
