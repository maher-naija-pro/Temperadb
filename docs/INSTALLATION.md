# Installation & Setup

This guide covers the complete installation and setup process for TimeSeriesDB.

## Prerequisites

- **Go 1.20 or higher** - [Download Go](https://golang.org/dl/)
- **Git** - For cloning the repository
- **Make** - For build automation (optional but recommended)

### Verify Prerequisites

```bash
# Check Go version
go version

# Check Git
git --version

# Check Make (optional)
make --version
```

## Installation

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd tsdb
```

### 2. Install Dependencies

```bash
go mod download
go mod verify
```

### 3. Build the Application

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

### 4. Verify Installation

```bash
# Check if binary was created
ls -la timeseriesdb

# Test the binary
./timeseriesdb --help
```

## Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Server Configuration
PORT=8080
HOST=0.0.0.0

# Data Storage
DATA_FILE=data.tsv
DATA_DIR=./data

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Performance
MAX_CONNECTIONS=1000
READ_TIMEOUT=30s
WRITE_TIMEOUT=30s
```

### Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `HOST` | `0.0.0.0` | HTTP server host |
| `DATA_FILE` | `data.tsv` | Path to TSV data file |
| `DATA_DIR` | `./data` | Directory for data files |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `LOG_FORMAT` | `json` | Log format (json, text) |
| `MAX_CONNECTIONS` | `1000` | Maximum concurrent connections |
| `READ_TIMEOUT` | `30s` | HTTP read timeout |
| `WRITE_TIMEOUT` | `30s` | HTTP write timeout |
| `SHUTDOWN_TIMEOUT` | `30s` | Application shutdown timeout |

### Configuration File

Alternatively, you can use a configuration file:

```yaml
# config.yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: "30s"
  write_timeout: "30s"
  shutdown_timeout: "30s"
  max_connections: 1000

storage:
  data_file: "data.tsv"
  data_dir: "./data"

logging:
  level: "info"
  format: "json"
```

## Usage

### Starting the Server

```bash
# Basic start
./timeseriesdb

# With custom config
./timeseriesdb -config config.yaml

# With environment variables
PORT=9090 DATA_FILE=my_data.tsv ./timeseriesdb

# In background
nohup ./timeseriesdb > timeseriesdb.log 2>&1 &

# As a service (systemd)
sudo systemctl start timeseriesdb
```

### Server Management

```bash
# Check server status
curl http://localhost:8080/health

# Stop server gracefully
pkill -TERM timeseriesdb

# Force stop
pkill -KILL timeseriesdb
```

### Data Writing

```bash
# Single data point
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01,region=us-west value=0.64 1434055562000000000"

# Multiple data points
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01,region=us-west value=0.64 1434055562000000000
cpu,host=server01,region=us-west value=0.65 1434055563000000000"

# From file
curl -X POST http://localhost:8080/write \
  --data-binary @data_points.txt
```

## Line Protocol Format

TimeSeriesDB accepts standard InfluxDB line protocol:

```
measurement[,tag_key=tag_value...] field_key=field_value[,field_key=field_value...] [timestamp]
```

### Examples

```bash
# Basic measurement
cpu value=0.64

# With tags
cpu,host=server01,region=us-west value=0.64

# With multiple fields
cpu,host=server01 value=0.64,user=23,system=41

# With timestamp
cpu,host=server01 value=0.64 1434055562000000000

# Multiple lines
cpu,host=server01 value=0.64 1434055562000000000
memory,host=server01 used=1024,free=2048 1434055562000000000
```

## Data Storage

### TSV Format

Data is stored in tab-separated value format:

```
timestamp	measurement	tag_key	tag_value	field_key	field_value
1434055562000000000	cpu	host	server01	value	0.64
1434055562000000000	cpu	region	us-west	value	0.64
```

### Storage Locations

- **Default**: `./data.tsv`
- **Custom**: Set via `DATA_FILE` environment variable
- **Directory**: Set via `DATA_DIR` environment variable

### Data Persistence

- Data is automatically persisted to disk
- No additional configuration required
- Data survives server restarts

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Check what's using the port
sudo netstat -tlnp | grep :8080

# Kill the process
sudo kill -9 <PID>

# Or use a different port
PORT=8081 ./timeseriesdb
```

**Permission denied:**
```bash
# Make binary executable
chmod +x timeseriesdb

# Check file permissions
ls -la timeseriesdb
```

**Data file not writable:**
```bash
# Check directory permissions
ls -la data/

# Fix permissions
chmod 755 data/
chmod 644 data.tsv
```

**Go version too old:**
```bash
# Update Go
go version
# If < 1.20, download from https://golang.org/dl/
```

### Logs and Debugging

```bash
# Enable debug logging
LOG_LEVEL=debug ./timeseriesdb

# Check logs
tail -f timeseriesdb.log

# Check system resources
top -p $(pgrep timeseriesdb)
```

### Performance Issues

```bash
# Check memory usage
ps aux | grep timeseriesdb

# Monitor disk I/O
iostat -x 1

# Check network connections
netstat -an | grep :8080
```

## Production Deployment

### System Requirements

- **CPU**: 2+ cores recommended
- **Memory**: 4GB+ RAM recommended
- **Disk**: SSD recommended for high I/O
- **Network**: 1Gbps+ recommended

### Security Considerations

```bash
# Run as non-root user
sudo useradd -r -s /bin/false timeseriesdb
sudo chown timeseriesdb:timeseriesdb timeseriesdb
sudo -u timeseriesdb ./timeseriesdb

# Firewall configuration
sudo ufw allow 8080/tcp

# Reverse proxy (nginx)
# See nginx configuration examples
```

### Monitoring

```bash
# Health check endpoint
curl http://localhost:8080/health

# Metrics endpoint (if implemented)
curl http://localhost:8080/metrics

# Process monitoring
ps aux | grep timeseriesdb
```

## Next Steps

After installation, you may want to:

1. **[Read the API Reference](API_REFERENCE.md)** - Learn about all available endpoints
2. **[Set up Development Environment](DEVELOPMENT.md)** - Configure your development workflow
3. **[Run Performance Tests](PERFORMANCE.md)** - Verify system performance
4. **[Configure CI/CD](CI_CD.md)** - Set up automated testing and deployment

For additional help, check the [GitHub Issues](https://github.com/yourusername/timeseriesdb/issues) or create a new one.
