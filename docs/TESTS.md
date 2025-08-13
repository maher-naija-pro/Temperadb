# Test Architecture

This directory contains a well-organized test suite for the TimeSeriesDB project with improved architecture and maintainability.

## Directory Structure

```
test/
├── suite_test.go               # Test suite setup and teardown
├── integration/                 # Integration tests
│   ├── metrics_endpoint_test.go    # Metrics endpoint tests
│   └── write_endpoint_test.go      # Write endpoint tests
├── utils/                      # Shared test utilities
│   ├── assertions.go           # Custom assertions and helpers
│   ├── data_test.go            # Test data factories
│   └── server_test.go          # HTTP test server factory
├── helpers/                    # Test-specific helper functions
│   ├── test_helpers.go         # Common test helpers
│   ├── config_helpers.go       # Configuration helpers
│   └── validation_helpers.go   # Validation helpers
├── benchmark/                  # Benchmark tests
│   ├── parser_benchmark_test.go    # Parser performance benchmarks
│   ├── storage_benchmark_test.go   # Storage performance benchmarks
│   └── api_benchmark_test.go       # API performance benchmarks
```

## Architecture Benefits

### 1. **Test Suite (`suite_test.go`)**
- **Centralized Setup**: Common test configuration and storage setup
- **Isolation**: Each test gets its own temporary directory
- **Cleanup**: Automatic cleanup of test files and resources
- **Reusability**: Consistent test environment across all tests

### 2. **Test Server Factory (`server_test.go`)**
- **HTTP Testing**: Reusable HTTP test server for API testing
- **Storage Integration**: Direct access to storage instance for verification
- **Cleanup**: Proper server cleanup after tests

### 3. **Test Data Factory (`data_test.go`)**
- **Consistent Data**: Standardized test data across all tests
- **Line Protocol**: Easy conversion between Points and line protocol format
- **Scalability**: Generate datasets of any size for testing

### 4. **Test Utilities (`utils/`)**
- **Custom Assertions**: Domain-specific assertion functions
- **Reusable Helpers**: Common test operations
- **Consistent Error Messages**: Standardized test failure output

### 5. **Organized Benchmarks (`benchmark/`)**
- **Separated Concerns**: Parser, storage, and API benchmarks in separate files
- **No Duplication**: Shared helper functions and data
- **Focused Testing**: Each benchmark file focuses on specific functionality

## Usage Examples

### Using the Test Suite

```go
func TestMyFunction(t *testing.T) {
    suite := NewTestSuite(t)
    defer suite.Cleanup()
    
    // Use suite.Storage for testing
    // Use suite.Config for configuration
    // Use suite.TempDir for temporary files
}
```

### Using the Test Server

```go
import "timeseriesdb/test/utils"

func TestAPIEndpoint(t *testing.T) {
    suite := NewTestSuite(t)
    defer suite.Cleanup()
    
    server := utils.NewTestServer(suite.Storage)
    defer server.Close()
    
    // Test HTTP endpoints using server.URL
}
```

### Using Test Data Factory

```go
import "timeseriesdb/test/utils"

func TestDataProcessing(t *testing.T) {
    // Get standard test data
    point := utils.DataFactory.SimplePoint()
    
    // Convert to line protocol
    lineProtocol := utils.DataFactory.LineProtocol(point)
    
    // Generate large datasets
    largeData := utils.DataFactory.GenerateLargeDataset(1000)
}
```

### Using Custom Assertions

```go
import "timeseriesdb/test/utils"

func TestPointEquality(t *testing.T) {
    expected := utils.DataFactory.SimplePoint()
    actual := processPoint(expected)
    
    utils.AssertPointEqual(t, expected, actual)
}
```

### Using Test Helpers

```go
import "timeseriesdb/test/helpers"

func TestPointCreation(t *testing.T) {
    // Create test points with helpers
    point := helpers.Helpers.CreateTestPoint("cpu", nil, nil)
    
    // Validate point structure
    helpers.Validation.ValidatePointStructure(t, point, "cpu")
    
    // Create test configuration
    cfg := helpers.Config.CreateTestConfig(t)
    
    // Create test data
    data := helpers.Helpers.GenerateTestData(100, "memory")
}
```

## Running Tests

### Run All Tests
```bash
go test ./test/...
```

### Run Specific Test Categories
```bash
# Run only unit tests
go test ./test/ -run TestUnit

# Run only integration tests
go test ./test/ -run TestIntegration

# Run only benchmarks
go test ./test/benchmark/... -bench=.
```

### Run Benchmarks with Specific Parameters
```bash
# Run parser benchmarks
go test ./test/benchmark/ -run=^$ -bench=BenchmarkParse

# Run storage benchmarks
go test ./test/benchmark/ -run=^$ -bench=BenchmarkWrite

# Run with memory profiling
go test ./test/benchmark/ -run=^$ -bench=BenchmarkMemoryUsage -benchmem
```

## Best Practices

1. **Always use the test suite** for tests that need storage or configuration
2. **Use the test server factory** for HTTP API testing
3. **Leverage the data factory** for consistent test data
4. **Use custom assertions** for domain-specific validations
5. **Keep benchmarks focused** on specific functionality
6. **Clean up resources** using defer statements

## Migration from Old Tests

The existing test files (`metrics_endpoint_test.go`, `write_endpoint_test.go`) can be gradually migrated to use the new architecture:

1. Replace manual setup with `NewTestSuite(t)`
2. Use `NewTestServer()` instead of manual HTTP server creation
3. Replace hardcoded test data with `DataFactory` methods
4. Use custom assertions from `utils/` package

This architecture provides better maintainability, reduces code duplication, and makes tests more reliable and consistent.
