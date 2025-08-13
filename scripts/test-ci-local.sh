#!/bin/bash

# Local CI Test Script - Replicates .github/workflows/tests-pipeline.yml

set -e

echo "🚀 Testing CI locally..."
echo "=========================="

# Step 1: Checkout code (already done locally)
echo "✅ Code checkout (local)"

# Step 2: Set up Go (check version)
echo "🔍 Checking Go version..."
go version
echo "✅ Go setup complete"

# Step 3: Download dependencies
echo "📦 Downloading dependencies..."
go mod tidy
echo "✅ Dependencies downloaded"

# Step 4: Run tests
echo "🧪 Running tests..."
go test -v ./...
echo "✅ Tests completed"

# Step 5: Run tests with coverage
echo "📊 Running tests with coverage..."
go test -v -coverprofile=coverage.out ./...
echo "✅ Coverage tests completed"

# Step 6: Show coverage
echo "📈 Coverage summary..."
go tool cover -func=coverage.out
echo "✅ Coverage displayed"

# Step 7: Quick benchmark check
echo "⚡ Running quick benchmark check..."
go test -bench=. -benchmem -timeout=2m ./test/ | head -20
echo "✅ Benchmark check completed"

echo ""
echo "🎉 Local CI test completed successfully!"
echo "📁 Coverage file: coverage.out"
echo "🌐 Open coverage.html in browser: go tool cover -html=coverage.out -o coverage.html"
