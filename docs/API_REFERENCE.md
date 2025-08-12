# API Reference

Complete API documentation for TimeSeriesDB.

## Base URL

```
http://localhost:8080
```

## Endpoints

### POST /write

Writes time series data in InfluxDB line protocol format.

#### Request

- **Method**: `POST`
- **Content-Type**: `text/plain`
- **Body**: Line protocol data

#### Line Protocol Format

```
measurement[,tag_key=tag_value...] field_key=field_value[,field_key=field_value...] [timestamp]
```

#### Examples

```bash
# Basic measurement
curl -X POST http://localhost:8080/write \
  -d "cpu value=0.64"

# With tags
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01,region=us-west value=0.64"

# With timestamp
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01 value=0.64 1434055562000000000"

# Multiple lines
curl -X POST http://localhost:8080/write \
  -d "cpu,host=server01 value=0.64 1434055562000000000
cpu,host=server01 value=0.65 1434055563000000000"
```

#### Response

**Success (200 OK):**
```
OK
```

**Error (400 Bad Request):**
```json
{
  "error": "Invalid line protocol format"
}
```

### GET /health

Health check endpoint.

#### Example

```bash
curl http://localhost:8080/health
```

#### Response

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "uptime": "2h30m15s"
}
```

### GET /metrics

Server metrics and statistics.

#### Example

```bash
curl http://localhost:8080/metrics
```

#### Response

```json
{
  "server": {
    "uptime": "2h30m15s",
    "requests_total": 1250
  },
  "storage": {
    "data_points_total": 45678,
    "disk_usage_bytes": 1048576
  }
}
```

## Data Types

### Supported Types

- **Float**: `value=0.64`
- **Integer**: `count=42`
- **String**: `status="running"`
- **Boolean**: `active=true`

### Rules

- Tags are indexed for fast querying
- Field values can be of different types
- Timestamps are in Unix nanoseconds
- Maximum 256 tags and fields per measurement

## Error Handling

### Common Errors

- **400**: Invalid line protocol format
- **405**: Method not allowed
- **500**: Internal server error

### Error Response

```json
{
  "error": "Error description",
  "details": "Additional details"
}
```

## Performance Tips

1. **Batch writes** - Send multiple points per request
2. **Use tags sparingly** - Tags are indexed
3. **Include timestamps** - Avoid server-generated timestamps
4. **Monitor metrics** - Use `/metrics` endpoint



## Testing

### Using cURL

```bash
# Test health
curl http://localhost:8080/health

# Test write
curl -X POST http://localhost:8080/write \
  -d "test,host=localhost value=1.0"

# Verbose output
curl -v -X POST http://localhost:8080/write \
  -d "test,host=localhost value=1.0"
```

## Next Steps

- **[Installation Guide](INSTALLATION.md)** - Set up TimeSeriesDB
- **[Performance Guide](PERFORMANCE.md)** - Optimize performance
- **[CI/CD Guide](CI_CD.md)** - Set up automation
