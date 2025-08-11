# Testing Documentation

This document describes the testing framework and practices used in the TimeSeriesDB project.

## Overview

The project uses Go's built-in testing framework with the following key components:

- **Go testing package**: Standard Go testing framework
- **httptest**: For HTTP route testing
- **testify**: For assertions and test suites
- **Table-driven tests**: Go's preferred testing pattern
- **Test coverage**: Comprehensive coverage reporting

## Test Structure

### Test Directory Organization

All tests are now organized in a dedicated `test/` directory for better project structure:

```
test/
├── README.md           # Test directory documentation
├── config.go           # Test configuration and environment setup
├── helpers.go          # Test helper utilities and functions
├── suite.go            # Base test suite with common functionality
├── main_test.go        # Main test suite for HTTP endpoints
├── benchmark_test.go   # Performance benchmarks
└── run.go              # Test runner and main entry point
```

### Main Test Suite (`test/main_test.go`)

The main test suite (`MainTestSuite`) provides comprehensive testing of the HTTP endpoints:

- **TestWriteEndpoint**: Tests various scenarios for the `/write` endpoint
- **TestWriteEndpointIntegration**: Integration tests
- **TestWriteEndpointPerformance**: Performance testing
- **TestWriteEndpointEdgeCases**: Edge case testing

### Test Helpers (`test/helpers.go`)

Utility functions to make testing easier and more maintainable:

- **TestHelper**: HTTP request creation and response assertion utilities
- **BenchmarkHelper**: Benchmarking utilities
- **Line Protocol Generation**: Helper functions for creating test data

### Base Test Suite (`test/suite.go`)

Common functionality for all test suites:

- **BaseTestSuite**: Provides storage, server, and environment management
- **Automatic Setup/Teardown**: Handles test environment lifecycle
- **Resource Management**: Automatic cleanup and resource management

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./test/...

# Run tests with verbose output
go test -v ./test/...

# Run tests with race detection
go test -race ./test/...

# Run specific test file
go test ./test/main_test.go

# Run specific test function
go test -run TestWriteEndpoint ./test/...
```

### Using Makefile

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with race detection
make test-race

# Generate coverage report
make test-coverage

# Run benchmarks
make benchmark

# Clean up test artifacts
make clean

# Show available commands
make help
```

## Test Coverage

Generate a coverage report:

```bash
make test-coverage
```

This will:
1. Run tests with coverage profiling
2. Generate an HTML coverage report (`coverage.html`)
3. Show function-level coverage in the terminal

## Benchmarking

Run performance benchmarks:

```bash
# Run all benchmarks
make benchmark

# Run specific benchmark
make benchmark-write
```

Benchmark results show:
- Operations per second
- Memory allocation per operation
- Number of allocations per operation

## Test Data

The tests use:
- **Temporary files**: `test_data.tsv` for test storage
- **Mock data**: Generated InfluxDB line protocol data
- **Cleanup**: Automatic cleanup after tests complete

## Testing Best Practices

### 1. Table-Driven Tests

Use table-driven tests for multiple test scenarios:

```go
tests := []struct {
    name           string
    method         string
    body           string
    expectedStatus int
    expectedBody   string
    description    string
}{
    // Test cases...
}

for _, tt := range tests {
    suite.Run(tt.name, func() {
        // Test implementation
    })
}
```

### 2. Test Setup and Teardown

Use the test suite lifecycle methods:

```go
func (suite *MainTestSuite) SetupSuite() {
    // One-time setup
}

func (suite *MainTestSuite) TearDownSuite() {
    // One-time cleanup
}

func (suite *MainTestSuite) SetupTest() {
    // Setup before each test
}
```

### 3. Assertions

Use testify assertions for clear test failures:

```go
assert.Equal(suite.T(), expectedStatus, resp.StatusCode)
require.NoError(suite.T(), err, "Failed to create request")
assert.Less(suite.T(), duration, 5*time.Second)
```

### 4. Error Handling

Always check for errors and provide meaningful messages:

```go
resp, err := http.DefaultClient.Do(req)
require.NoError(suite.T(), err, "Failed to execute request")
```

## Test Categories

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Fast execution

### Integration Tests
- Test component interactions
- Use real storage and HTTP handlers
- Verify end-to-end functionality

### Performance Tests
- Benchmark critical paths
- Measure response times
- Identify bottlenecks

### Edge Case Tests
- Test boundary conditions
- Test error scenarios
- Test invalid inputs

## Adding New Tests

### 1. Create Test Function

```go
func (suite *MainTestSuite) TestNewFeature() {
    // Test implementation
}
```

### 2. Use Test Helpers

```go
helper := NewTestHelper()
req := helper.CreateTestRequest(http.MethodPost, "/endpoint", "data", nil)
resp, err := helper.ExecuteRequest(req)
require.NoError(suite.T(), err)

helper.AssertResponseStatus(suite.T(), resp, http.StatusOK, "Should return success")
```

### 3. Generate Test Data

```go
helper := NewTestHelper()
testData := helper.GenerateLineProtocolData("cpu", tags, fields, timestamp)
```

## Troubleshooting

### Common Issues

1. **Logger not initialized**: Ensure `logger.Init()` is called in test setup
2. **File already closed**: Check storage cleanup in test teardown
3. **Port conflicts**: Tests use random ports via `httptest.NewServer`

### Debug Mode

Run tests with verbose output to see detailed execution:

```bash
go test -v -count=1 ./...
```

The `-count=1` flag ensures tests run in isolation.

## Continuous Integration

Tests should pass in CI environments:

```yaml
# Example GitHub Actions
- name: Run Tests
  run: make test

- name: Generate Coverage
  run: make test-coverage

- name: Run Benchmarks
  run: make benchmark
```

## Performance Considerations

- Tests run in parallel where possible
- Use `httptest.NewServer` for HTTP testing
- Clean up resources after each test
- Avoid long-running operations in unit tests

## Contributing

When adding new features:
1. Write tests first (TDD approach)
2. Ensure all tests pass
3. Maintain or improve test coverage
4. Follow existing test patterns
5. Document complex test scenarios
