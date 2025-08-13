# Environment Variables Package

This package provides a robust and organized way to parse environment variables for the TimeSeriesDB application.

## Features

- **Type-safe parsing**: Automatic conversion to appropriate Go types
- **Default values**: Fallback to sensible defaults when environment variables are not set
- **Whitespace handling**: Automatic trimming of whitespace from values
- **Human-readable file sizes**: Support for KB, MB, GB, TB suffixes
- **Boolean parsing**: Support for various boolean representations (true/false, yes/no, 1/0, on/off, enabled/disabled)
- **Duration parsing**: Automatic conversion of seconds to time.Duration
- **Validation**: Built-in validation and error handling

## Usage

### Basic Usage

```go
import "timeseriesdb/internal/envvars"

// Create a new parser
parser := envvars.NewParser()

// Parse environment variables with defaults
port := parser.String(envvars.Port, envvars.DefaultPort)
maxConnections := parser.Int(envvars.MaxConnections, envvars.DefaultMaxConnections)
compression := parser.Bool(envvars.Compression, envvars.DefaultCompression)
timeout := parser.Duration(envvars.ReadTimeout, envvars.DefaultReadTimeout)
fileSize := parser.FileSize(envvars.MaxFileSize, envvars.DefaultMaxFileSize)
```

### Available Methods

- `String(key, defaultValue string) string` - Parse string values
- `Int(key string, defaultValue int) int` - Parse integer values
- `Int64(key string, defaultValue int64) int64` - Parse int64 values
- `Bool(key string, defaultValue bool) bool` - Parse boolean values
- `Duration(key string, defaultValue time.Duration) time.Duration` - Parse duration values (in seconds)
- `FileSize(key string, defaultValue int64) int64` - Parse file size values (supports KB, MB, GB, TB)
- `Has(key string) bool` - Check if an environment variable is set
- `IsSet(key string) bool` - Check if an environment variable is set and not empty

### Environment Variable Keys

The package provides constants for all environment variable keys:

```go
// Server Configuration
envvars.Port         // "PORT"
envvars.ReadTimeout  // "READ_TIMEOUT"
envvars.WriteTimeout // "WRITE_TIMEOUT"
envvars.IdleTimeout  // "IDLE_TIMEOUT"

// Storage Configuration
envvars.DataFile     // "DATA_FILE"
envvars.MaxFileSize  // "MAX_FILE_SIZE"
envvars.BackupDir    // "BACKUP_DIR"
envvars.Compression  // "COMPRESSION"

// Logging Configuration
envvars.LogLevel      // "LOG_LEVEL"
envvars.LogFormat     // "LOG_FORMAT"
envvars.LogOutput     // "LOG_OUTPUT"
envvars.LogMaxSize    // "LOG_MAX_SIZE"
envvars.LogMaxBackups // "LOG_MAX_BACKUPS"
envvars.LogMaxAge     // "LOG_MAX_AGE"
envvars.LogCompress   // "LOG_COMPRESS"

// Database Configuration
envvars.MaxConnections // "MAX_CONNECTIONS"
envvars.ConnectionTTL  // "CONNECTION_TTL"
envvars.QueryTimeout   // "QUERY_TIMEOUT"
```

### Default Values

The package also provides default values for all configuration options:

```go
// Server Defaults
envvars.DefaultPort         // "8080"
envvars.DefaultReadTimeout  // 30 * time.Second
envvars.DefaultWriteTimeout // 30 * time.Second
envvars.DefaultIdleTimeout  // 120 * time.Second

// Storage Defaults
envvars.DefaultDataFile    // "data.tsv"
envvars.DefaultMaxFileSize // 1073741824 (1GB)
envvars.DefaultBackupDir   // "backups"
envvars.DefaultCompression // false

// Logging Defaults
envvars.DefaultLogLevel      // "info"
envvars.DefaultLogFormat     // "text"
envvars.DefaultLogOutput     // "stdout"
envvars.DefaultLogMaxSize    // 100
envvars.DefaultLogMaxBackups // 3
envvars.DefaultLogMaxAge     // 28
envvars.DefaultLogCompress   // true

// Database Defaults
envvars.DefaultMaxConnections // 100
envvars.DefaultConnectionTTL  // 300 * time.Second (5 minutes)
envvars.DefaultQueryTimeout   // 30 * time.Second
```

### File Size Parsing

The `FileSize` method supports human-readable file size formats:

```go
// Examples of supported formats:
"1024"     // 1024 bytes
"1KB"      // 1024 bytes
"2MB"      // 2,097,152 bytes
"1GB"      // 1,073,741,824 bytes
"1TB"      // 1,099,511,627,776 bytes
```

### Boolean Parsing

The `Bool` method supports various boolean representations:

```go
// True values: "true", "1", "yes", "on", "enabled"
// False values: "false", "0", "no", "off", "disabled"
// Case-insensitive
```

### Duration Parsing

The `Duration` method expects values in seconds and converts them to `time.Duration`:

```go
// Examples:
"30"       // 30 seconds
"60"       // 1 minute
"300"      // 5 minutes
"3600"     // 1 hour
```

## Integration with Configuration

This package is designed to work seamlessly with the main configuration system. The configuration structs use this package to parse environment variables:

```go
// Example from server.go
func NewServerConfig() ServerConfig {
    parser := envvars.NewParser()
    
    return ServerConfig{
        Port:         parser.String(envvars.Port, envvars.DefaultPort),
        ReadTimeout:  parser.Duration(envvars.ReadTimeout, envvars.DefaultReadTimeout),
        WriteTimeout: parser.Duration(envvars.WriteTimeout, envvars.DefaultWriteTimeout),
        IdleTimeout:  parser.Duration(envvars.IdleTimeout, envvars.DefaultIdleTimeout),
    }
}
```

## Benefits

1. **Centralized**: All environment variable parsing logic is in one place
2. **Type-safe**: Automatic type conversion with proper error handling
3. **Maintainable**: Easy to add new environment variables and defaults
4. **Testable**: Comprehensive test coverage for all parsing functions
5. **Flexible**: Support for various input formats and edge cases
6. **Consistent**: Uniform behavior across all configuration types
