#!/bin/bash

echo "🚀 Starting TimeSeriesDB with metrics endpoint..."
echo ""

# Build the application
echo "📦 Building application..."
go build -o timeseriesdb .

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"
echo ""

# Start the application in background
echo "🌐 Starting server on port 8080..."
./timeseriesdb &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo ""
echo "📊 Testing metrics endpoint..."
echo ""

# Test health endpoint
echo "🏥 Health check:"
curl -s http://localhost:8080/health
echo ""
echo ""

# Test metrics endpoint
echo "📈 Metrics endpoint:"
curl -s http://localhost:8080/metrics | head -20
echo ""
echo "..."

# Make a request to generate some metrics
echo "📝 Making a test request to generate metrics..."
curl -s -X POST http://localhost:8080/write \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}' || echo "Expected error for invalid data"
echo ""
echo ""

# Check metrics again
echo "📈 Metrics after request:"
curl -s http://localhost:8080/metrics | grep -E "(tsdb_http_requests_total|tsdb_http_request_duration_seconds)" | head -10
echo ""

echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null

echo "✅ Demo completed!"
echo ""
echo "📋 Summary:"
echo "   - Metrics endpoint: http://localhost:8080/metrics"
echo "   - Health endpoint: http://localhost:8080/health"
echo "   - Write endpoint: http://localhost:8080/write"
echo ""
echo "🔍 The metrics endpoint exposes Prometheus-compatible metrics including:"
echo "   - HTTP request counts and durations"
echo "   - Ingestion metrics"
echo "   - Storage and compaction metrics"
echo "   - Resource usage metrics"
