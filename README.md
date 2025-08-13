# TimeSeriesDB


[![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=flat-square)
[![Coverage](https://img.shields.io/badge/coverage-56.0%25-brightgreen?style=flat-square)
[![Go Version](https://img.shields.io/badge/go-1.24.5-blue?style=flat-square)
[![License](https://img.shields.io/badge/license-AGPL%20v3.0-red?style=flat-square)


A lightweight, high-performance time series database written in Go that accepts InfluxDB line protocol for data ingestion.

## What is TimeSeriesDB?

TimeSeriesDB is a simple yet powerful time series database designed for:
- **High-throughput data ingestion** using InfluxDB line protocol
- **Efficient storage** with TSV (tab-separated values) backend
- **HTTP API** for easy integration
- **Built-in metrics** for monitoring and observability
- **Lightweight deployment** with minimal resource requirements

## Quick Start

### 1. Install & Build

```bash
git clone <your-repo-url>
cd tsdb
go mod download
go build -o timeseriesdb
```

### 2. Configure



**Minimum Requirements:**
- Go 1.24 or higher
- 512MB RAM
- 1GB disk space for data storage

### 3. Run

```bash
# Start the server
./timeseriesdb

# Check if it's running
curl http://localhost:8080/health
```

### 4. Write Data

```bash
# Single data point
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01 value=0.64"

# With timestamp
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01 value=0.64 1434055562000000000"

# Multiple points
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01 value=0.64 1434055562000000000
cpu,host=server01 value=0.65 1434055563000000000"
```

## Data Format

TimeSeriesDB accepts standard InfluxDB line protocol:

```
measurement[,tag_key=tag_value...] field_key=field_value[,field_key=field_value...] [timestamp]
```

**Examples:**
- `cpu value=0.64` - Basic measurement
- `cpu,host=server01 value=0.64` - With tags
- `cpu,host=server01 value=0.64,user=23` - Multiple fields
- `cpu,host=server01 value=0.64 1434055562000000000` - With timestamp

## API Endpoints

- **`POST /write`** - Write time series data
- **`GET /health`** - Health check
- **`GET /metrics`** - Prometheus metrics

## Documentation

For detailed information, see:

- **[📖 Installation Guide](docs/INSTALLATION.md)** - Complete setup and configuration
- **[🔌 API Reference](docs/API_REFERENCE.md)** - Full API documentation and examples
- **[📊 Metrics & Monitoring](docs/METRICS.md)** - Prometheus integration and observability
- **[⚡ Performance Guide](docs/PERFORMANCE.md)** - Benchmarking and optimization
- **[🧪 Testing Guide](docs/TESTS.md)** - Test architecture and guidelines
- **[🚀 CI/CD Guide](docs/CI_CD.md)** - Automated workflows and deployment

## Project Structure

```
tsdb/
├── main.go                    # Main application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
├── internal/                  # Core packages
│   ├── storage/              # Data storage implementation
│   ├── ingestion/            # Line protocol parsing
│   ├── api/                  # HTTP handlers and routing
│   ├── metrics/              # Prometheus metrics collection
│   └── config/               # Configuration management
├── test/                     # Test files and benchmarks
├── docs/                     # Documentation files
├── scripts/                  # Utility and build scripts
└── README.md                 # This file
```

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build Docker image
make build-docker
```

## License

GNU Affero General Public License v3.0 - see [LICENSE](LICENSE) file for details.

---

**Need help?** Check the [documentation](docs/) or [create an issue](https://github.com/yourusername/timeseriesdb/issues). 


