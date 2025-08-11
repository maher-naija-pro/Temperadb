# Test Directory

This directory contains all the test files for the TimeSeriesDB project, organized in a clean and maintainable structure.

## Directory Structure

```
test/
├── README.md           # This file
├── config.go           # Test configuration and environment setup
├── helpers.go          # Test helper utilities and functions
├── suite.go            # Base test suite with common functionality
├── main_test.go        # Main test suite for HTTP endpoints
├── benchmark_test.go   # Performance benchmarks
└── run.go              # Test runner and main entry point
```

## Files Overview

### `config.go`
- **TestConfig**: Configuration struct for test settings
- **SetupTestEnvironment**: Initialize test environment
- **CleanupTestEnvironment**: Clean up after tests
- **CI Environment Detection**: Automatic timeout adjustment for CI

### `helpers.go`
- **TestHelper**: HTTP request creation and response validation
- **BenchmarkHelper**: Benchmarking utilities
- **Line Protocol Generation**: Helper functions for creating test data

### `suite.go`
- **BaseTestSuite**: Common functionality for all test suites
- **Storage Management**: Automatic storage setup and cleanup
- **Server Management**: Test server creation and management
- **Environment Management**: Automatic environment setup/teardown

### `main_test.go`
- **MainTestSuite**: Tests for HTTP endpoints
- **Route Testing**: Comprehensive endpoint validation
- **Integration Tests**: End-to-end functionality testing
- **Edge Case Testing**: Boundary conditions and error scenarios

### `benchmark_test.go`
- **Performance Tests**: HTTP endpoint benchmarks
- **Bulk Operations**: Large dataset performance testing
- **Request Creation**: HTTP request creation benchmarks

### `run.go`
- **TestMain**: Entry point for test execution
- **Test Orchestration**: Custom test suite execution

## Usage

### Running Tests

```bash
# From project root
make test                    # Run all tests
make test-verbose           # Run with verbose output
make test-race             # Run with race detection
make test-coverage         # Generate coverage report
make benchmark             # Run all benchmarks
make benchmark-write       # Run specific benchmark
```

### From Go Commands

```bash
# Run all tests in test directory
go test ./test/...

# Run specific test file
go test ./test/main_test.go

# Run specific test function
go test -run TestWriteEndpoint ./test/...

# Run benchmarks
go test -bench=. ./test/...
```

## Test Organization Principles

1. **Separation of Concerns**: Each file has a specific responsibility
2. **Reusability**: Common functionality in base classes and helpers
3. **Maintainability**: Clear structure and documentation
4. **Performance**: Efficient test execution and cleanup
5. **CI Friendly**: Automatic environment detection and configuration

## Adding New Tests

### 1. Create Test Suite

```go
type MyFeatureTestSuite struct {
    BaseTestSuite
}

func (suite *MyFeatureTestSuite) TestMyFeature() {
    // Test implementation
}

func TestMyFeatureSuite(t *testing.T) {
    suite.Run(t, new(MyFeatureTestSuite))
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
testData := suite.GenerateTestData("cpu", tags, fields, timestamp)
```

## Configuration

Test configuration is centralized in `config.go`:

- **Port**: Test server port (default: 8080)
- **Data File**: Test storage file (default: test_data.tsv)
- **Temp Directory**: Test files location (default: test/)
- **Timeouts**: Configurable test timeouts
- **CI Detection**: Automatic CI environment detection

## Best Practices

1. **Use BaseTestSuite**: Inherit from BaseTestSuite for common functionality
2. **Cleanup Resources**: Always clean up in TearDownSuite
3. **Use Test Helpers**: Leverage helper functions for common operations
4. **Meaningful Names**: Use descriptive test and function names
5. **Documentation**: Document complex test scenarios
6. **Error Handling**: Provide clear error messages in assertions

## Troubleshooting

### Common Issues

1. **Storage Errors**: Check if test directory exists and is writable
2. **Port Conflicts**: Tests use random ports via httptest
3. **File Permissions**: Ensure test directory has proper permissions
4. **Environment Variables**: Check if required env vars are set

### Debug Mode

```bash
# Run tests with verbose output
go test -v -count=1 ./test/...

# Run specific test with verbose output
go test -v -run TestWriteEndpoint ./test/...
```

## Contributing

When adding new tests:

1. Follow the existing structure and patterns
2. Use the base test suite for common functionality
3. Leverage test helpers for HTTP operations
4. Add appropriate documentation
5. Ensure tests pass in CI environment
6. Maintain or improve test coverage
