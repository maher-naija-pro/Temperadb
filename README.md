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

## Quick Start

### Installation

```bash
git clone <your-repo-url>
cd tsdb
go mod download
go build -o timeseriesdb
```

### Configuration

Create a `.env` file:
```env
PORT=8080
DATA_FILE=data.tsv
```

### Usage

```bash
# Start the server
./timeseriesdb

# Write data
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
```

## Documentation

- **[Installation & Setup](docs/INSTALLATION.md)** - Detailed installation instructions and configuration
- **[API Reference](docs/API_REFERENCE.md)** - Complete API documentation and examples
- **[Development Guide](docs/DEVELOPMENT.md)** - Development setup, testing, and contribution guidelines
- **[Performance & Benchmarks](docs/PERFORMANCE.md)** - Performance testing, benchmarks, and optimization
- **[CI/CD Pipeline](docs/CI_CD.md)** - Continuous integration and deployment information

## Project Structure

```
tsdb/
├── main.go                    # Main application and HTTP server
├── internal/                  # Internal packages
│   ├── storage/              # Data storage implementation
│   ├── ingestion/            # Line protocol ingestion
│   ├── types/                # Data type definitions
│   └── logger/               # Logging utilities
├── test/                     # Test files and benchmarks
├── docs/                     # Documentation
├── scripts/                  # Utility scripts
├── go.mod                    # Go module dependencies
└── Makefile                  # Build and test automation
```

## Testing

```bash
# Run all tests
make test

# Run benchmarks
make benchmark

# Run with coverage
make test-coverage
```

## CI/CD and Package Building

This project includes comprehensive CI/CD workflows that automatically build and publish packages for multiple platforms.

### Automated Builds

The CI/CD system automatically builds packages when:
- You push a tag (e.g., `git tag v1.0.0 && git push origin v1.0.0`)
- You create a pull request
- You manually trigger the workflow

### Supported Platforms

- **Linux**: AMD64, ARM64
- **Windows**: AMD64
- **macOS**: AMD64, ARM64 (Apple Silicon)

### Building Locally

```bash
# Build for current platform
make build

# Build for specific platform
make build-linux
make build-windows
make build-darwin

# Build for all platforms
make build-all

# Build Docker image
make build-docker

# Clean build artifacts
make clean-build

# Show build help
make build-help
```

### GitHub Packages

When you create a release tag, the CI/CD system automatically:
1. Builds binaries for all supported platforms
2. Creates a GitHub release with downloadable artifacts
3. Publishes Docker images to GitHub Container Registry
4. Runs security scans and vulnerability checks

### Docker Images

Docker images are available at:
```
ghcr.io/yourusername/timeseriesdb:latest
ghcr.io/yourusername/timeseriesdb:v1.0.0
```

### Manual Workflow Trigger

You can manually trigger the build workflow:
1. Go to Actions → Build and Publish Packages
2. Click "Run workflow"
3. Optionally specify a version
4. Click "Run workflow"

## License

MIT License - see [LICENSE](LICENSE) file for details. 


