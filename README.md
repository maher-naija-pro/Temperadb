# TimeSeriesDB

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

**Response:**
- Success: `200 OK` with "OK" message
- Error: `400 Bad Request` for invalid line protocol
- Error: `405 Method Not Allowed` for non-POST requests

## Development

### Project Structure

```
tsdb/
├── main.go          # Main application and HTTP server
├── storage.go       # Data storage implementation
├── parser.go        # Line protocol parser
├── go.mod          # Go module dependencies
├── go.sum          # Dependency checksums
├── .gitignore      # Git ignore patterns
├── data.tsv        # Time series data storage
└── README.md       # This file
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
go test ./...
```

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

## Support

[Add support information here]
