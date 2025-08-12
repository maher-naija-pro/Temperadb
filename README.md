# TimeSeriesDB

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=flat-square)](https://github.com/yourusername/timeseriesdb/actions)
[![Coverage](https://img.shields.io/badge/coverage-80%25-green?style=flat-square)](https://github.com/yourusername/timeseriesdb/actions)
[![Go Version](https://img.shields.io/badge/go-1.20+-blue?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

A lightweight, high-performance time series database written in Go that accepts InfluxDB line protocol for data ingestion.

## Features

- **HTTP API**: RESTful endpoint for writing time series data
- **InfluxDB Line Protocol**: Compatible with InfluxDB line protocol format
- **TSV Storage**: Efficient tab-separated value storage backend
- **Environment Configuration**: Configurable via environment variables
- **Structured Logging**: Comprehensive logging with logrus
- **High Performance**: Optimized for high-throughput time series data ingestion

## Architecture

The project consists of several key components:

- **`main.go`**: HTTP server and main application logic
- **`storage.go`**: Data storage and persistence layer
- **`parser.go`**: InfluxDB line protocol parser
- **`data.tsv`**: Time series data storage file

## Prerequisites

- Go 1.20 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd tsdb
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o timeseriesdb
```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
PORT=8080
DATA_FILE=data.tsv
```

- `PORT`: HTTP server port (default: 8080)
- `DATA_FILE`: Path to the TSV data file (default: data.tsv)

## Usage

### Starting the Server

```bash
./timeseriesdb
```

The server will start on the configured port and begin accepting HTTP requests.

### Writing Data

Send POST requests to `/write` endpoint with InfluxDB line protocol data:

```bash
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
```

### Line Protocol Format

The server accepts standard InfluxDB line protocol:

```
measurement[,tag_key=tag_value...] field_key=field_value[,field_key=field_value...] [timestamp]
```

Example:
```
cpu,host=server01,region=us-west value=0.64,user=23 1434055562000000000
```

## API Endpoints

### POST /write

Accepts time series data in InfluxDB line protocol format.

**Request:**
- Method: POST
- Content-Type: text/plain
- Body: Line protocol data (one or more lines)

## 🧪 Testing & Quality

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Generate coverage report
make test-coverage

# Run benchmarks
make benchmark

# Run specific benchmark categories
make benchmark-parser      # Parser performance only
make benchmark-storage     # Storage performance only
make benchmark-http        # HTTP endpoint performance only
make benchmark-e2e         # End-to-end workflow performance only
make benchmark-memory      # Memory usage performance only

# Run all benchmarks with progress display
make benchmark-all

# Run with CPU and memory profiling
make benchmark-profile

# Get benchmark help
make benchmark-help
```

### Test Coverage

The project maintains high test coverage with comprehensive testing of:
- HTTP endpoints and routing
- Data parsing and validation
- Storage operations
- Error handling and edge cases
- Performance benchmarks

### Performance Benchmarks

The project includes a comprehensive benchmark suite covering all major components:

#### Benchmark Categories
- **Parser Performance**: Line protocol parsing with various data sizes and complexities
- **Storage Performance**: Point writing with different tag/field counts
- **HTTP Endpoint Performance**: End-to-end HTTP write operations
- **End-to-End Workflows**: Complete data processing pipelines
- **Memory Usage**: Allocation tracking and memory profiling
- **Concurrent Operations**: Parallel write performance testing

#### Benchmark Features
- **Realistic test data** with CPU metrics, tags, and fields
- **Scalable datasets** from 1 to 10,000 lines
- **Performance metrics** for throughput, latency, memory, and allocations
- **Profiling support** for CPU and memory analysis
- **Timeout protection** for long-running benchmarks

#### Quick Benchmark Commands
```bash
# Run all benchmarks with progress display
make benchmark-all

# Run specific component benchmarks
make benchmark-parser
make benchmark-storage
make benchmark-http

# Run with profiling
make benchmark-profile

# Get help
make benchmark-help
```

#### Manual Benchmark Execution
```bash
# All benchmarks
go test -bench=. -benchmem ./test/

# Specific patterns
go test -bench=BenchmarkParse -benchmem ./test/
go test -bench=BenchmarkWrite -benchmem ./test/
go test -bench=BenchmarkHTTP -benchmem ./test/

# With profiling
go test -bench=BenchmarkParseLargeDataset -cpuprofile=cpu.prof -benchmem ./test/
go test -bench=BenchmarkMemoryUsage -memprofile=memory.prof -benchmem ./test/
```

#### Understanding Benchmark Results

Benchmark output format:
```
BenchmarkName-16         1000        1234567 ns/op        1234 B/op        10 allocs/op
```

- **BenchmarkName-16**: Name and CPU cores
- **1000**: Number of iterations
- **1234567 ns/op**: Time per operation (nanoseconds)
- **1234 B/op**: Memory allocated per operation (bytes)
- **10 allocs/op**: Number of allocations per operation

#### Performance Metrics
- **Throughput**: Operations per second (higher is better)
- **Latency**: Time per operation (lower is better)
- **Memory**: Bytes allocated per operation (lower is better)
- **Allocations**: Number of memory allocations per operation (lower is better)

#### Profiling and Analysis

After running `make benchmark-profile`, you'll get:
- `cpu_profile.prof` - CPU usage analysis
- `memory_profile.prof` - Memory usage analysis

Analyze profiles with:
```bash
go tool pprof cpu_profile.prof
go tool pprof memory_profile.prof

# Web interface
go tool pprof -http=:8080 cpu_profile.prof
go tool pprof -http=:8080 memory_profile.prof
```

#### Performance Tips
1. **Run multiple times** to account for system variance
2. **Use consistent environment** for comparable results
3. **Monitor system resources** during execution
4. **Profile first** to identify bottlenecks before optimizing

### Code Quality

- **Linting**: golangci-lint with comprehensive rules
- **Formatting**: Automatic code formatting with `go fmt`
- **Security**: Regular security scans with gosec
- **Coverage**: Minimum 80% test coverage enforced

## 🚀 CI/CD Pipeline

The project uses GitHub Actions for continuous integration:

- **Automated Testing**: Runs on every push and PR
- **Multi-Platform**: Tests on Ubuntu, Windows, and macOS
- **Go Versions**: Supports Go 1.20, 1.21, and 1.22
- **Quality Gates**: Enforces code quality and security standards
- **Performance Monitoring**: Tracks benchmarks and performance metrics
- **Benchmark Regression Detection**: Monitors performance changes over time

For detailed CI/CD information, see [CI.md](CI.md).

**Response:**
- Success: `200 OK` with "OK" message
- Error: `400 Bad Request` for invalid line protocol
- Error: `405 Method Not Allowed` for non-POST requests

## Development

### Project Structure

```
tsdb/
├── main.go                    # Main application and HTTP server
├── internal/                  # Internal packages
│   ├── storage/              # Data storage implementation
│   ├── parser/               # Line protocol parser
│   ├── types/                # Data type definitions
│   └── logger/               # Logging utilities
├── test/                     # Test files
│   ├── benchmark_test.go     # Comprehensive performance benchmarks
│   └── write_endpoint_test.go # HTTP endpoint tests
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
├── .gitignore               # Git ignore patterns
├── Makefile                 # Build and test automation
├── data.tsv                 # Time series data storage
└── README.md                # This file
```

### Dependencies

- `github.com/joho/godotenv`: Environment variable loading
- `github.com/sirupsen/logrus`: Structured logging

### Building

```bash
# Development build
go build

# Production build with optimizations
go build -ldflags="-s -w" -o timeseriesdb

# Cross-compilation for different platforms
GOOS=linux GOARCH=amd64 go build -o timeseriesdb-linux
GOOS=darwin GOARCH=amd64 go build -o timeseriesdb-macos
GOOS=windows GOARCH=amd64 go build -o timeseriesdb-windows.exe
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./test/

# Run specific test patterns
go test -run TestWrite ./test/
go test -bench=BenchmarkParse ./test/
```

## License

MIT 


