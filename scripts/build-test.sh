#!/bin/bash

# Build Test Script for TimeSeriesDB
# This script tests the build process locally to verify CI/CD setup

set -e

echo "========================================="
echo "  TimeSeriesDB Build Test Script        "
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
    print_error "This script must be run from the project root directory"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_status "Go version: $GO_VERSION"

# Clean previous builds
print_status "Cleaning previous build artifacts..."
make clean-build 2>/dev/null || true

# Test basic build
print_status "Testing basic build..."
make build
if [ -f "timeseriesdb" ]; then
    print_status "Basic build successful"
    rm -f timeseriesdb
else
    print_error "Basic build failed"
    exit 1
fi

# Test platform-specific builds
print_status "Testing platform-specific builds..."

# Linux build
print_status "Building for Linux AMD64..."
make build-linux
if [ -f "timeseriesdb-linux-amd64" ]; then
    print_status "Linux build successful"
    rm -f timeseriesdb-linux-amd64
else
    print_error "Linux build failed"
    exit 1
fi

# Windows build
print_status "Building for Windows AMD64..."
make build-windows
if [ -f "timeseriesdb-windows-amd64.exe" ]; then
    print_status "Windows build successful"
    rm -f timeseriesdb-windows-amd64.exe
else
    print_error "Windows build failed"
    exit 1
fi

# macOS build
print_status "Building for macOS AMD64..."
make build-darwin
if [ -f "timeseriesdb-darwin-amd64" ]; then
    print_status "macOS build successful"
    rm -f timeseriesdb-darwin-amd64
else
    print_error "macOS build failed"
    exit 1
fi

# Test Docker build (if Docker is available)
if command -v docker &> /dev/null; then
    print_status "Testing Docker build..."
    make build-docker
    print_status "Docker build successful"
    
    # Clean up Docker image
    docker rmi timeseriesdb:dev 2>/dev/null || true
else
    print_warning "Docker not available, skipping Docker build test"
fi

# Test with custom version
print_status "Testing build with custom version..."
make build VERSION=v1.0.0-test
if [ -f "timeseriesdb" ]; then
    print_status "Custom version build successful"
    rm -f timeseriesdb
else
    print_error "Custom version build failed"
    exit 1
fi

# Clean up
print_status "Cleaning up..."
make clean-build

echo ""
echo "========================================="
echo "  All build tests passed!               "
echo "========================================="
echo ""
print_status "Your build system is working correctly"
print_status "You can now use the CI/CD workflows:"
echo "  - Push a tag to trigger automated builds"
echo "  - Use 'make build-help' for build options"
echo "  - Check .github/workflows/ for CI/CD details"
