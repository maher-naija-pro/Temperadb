# Integration Tests

This directory contains integration tests for the TimeSeriesDB API that test the complete workflow from HTTP requests to storage.

## Test Files

### `api_integration_test.go`
Basic integration tests that test the API components in isolation:
- **TestWriteEndpointIntegration**: Tests write endpoint with various data formats
- **TestWriteEndpointHTTPMethods**: Tests all HTTP methods on write endpoint
- **TestHealthEndpointIntegration**: Tests health endpoint functionality
- **TestEndToEndWorkflow**: Tests complete data flow
- **TestConcurrentWrites**: Tests concurrent write operations
- **TestErrorHandling**: Tests various error scenarios

### `api_full_integration_test.go`
Full integration tests that use a real HTTP server:
- **TestWriteEndpointFullIntegration**: Tests write endpoint with real HTTP requests
- **TestWriteEndpointHTTPMethodsFull**: Tests HTTP methods with real server
- **TestHealthEndpointFullIntegration**: Tests health endpoint with real HTTP
- **TestMetricsEndpointFullIntegration**: Tests metrics endpoint with real HTTP
- **TestEndToEndWorkflowFull**: Tests complete workflow with real server
- **TestConcurrentWritesFull**: Tests concurrent writes with real HTTP
- **TestErrorHandlingFull**: Tests error scenarios with real HTTP

## Running the Tests

### Run all integration tests:
```bash
go test ./test/integration/...
```

### Run specific test file:
```bash
go test ./test/integration/api_integration_test.go
go test ./test/integration/api_full_integration_test.go
```

### Run with verbose output:
```bash
go test -v ./test/integration/...
```

### Run with race detection:
```bash
go test -race ./test/integration/...
```

## Test Coverage

The integration tests cover:

### API Endpoints
- **POST /write**: Data ingestion in InfluxDB line protocol format
- **GET /health**: Health check endpoint
- **GET /metrics**: Prometheus metrics endpoint

### Test Scenarios
- Valid data ingestion (single point, multiple points, complex points)
- Invalid data handling (malformed line protocol, empty bodies)
- HTTP method validation (only POST allowed for write)
- Concurrent write operations
- Error handling and edge cases
- End-to-end data flow validation

### Data Formats
- Simple measurements with basic tags and fields
- Complex measurements with many tags and fields
- Multiple data points in single request
- Various timestamp formats
- Edge cases and error conditions

## Test Architecture

### TestSuite Structure
Each test suite provides:
- Isolated storage instance with temporary files
- Cleanup after each test
- Configuration management
- Helper methods for common operations

### HTTP Testing
- Uses `httptest.Server` for real HTTP testing
- Tests actual HTTP requests and responses
- Validates response status codes, headers, and body content
- Tests error conditions and edge cases

### Storage Integration
- Tests actual storage operations
- Validates data persistence
- Tests concurrent access patterns
- Verifies data integrity

## Adding New Tests

### For Basic Integration Tests:
1. Add test function to `api_integration_test.go`
2. Use `NewTestSuite(t)` for setup
3. Test components in isolation
4. Use helper functions for common operations

### For Full Integration Tests:
1. Add test function to `api_full_integration_test.go`
2. Use `NewFullAPITestSuite(t)` for setup
3. Make real HTTP requests to test server
4. Validate complete request/response cycle

### Test Naming Convention:
- Use descriptive test names that explain what is being tested
- Group related tests using `t.Run()` for subtests
- Use table-driven tests for multiple test cases
- Include both positive and negative test cases

## Dependencies

The integration tests depend on:
- `timeseriesdb/internal/api/http` - HTTP router and handlers
- `timeseriesdb/internal/storage` - Storage implementation
- `timeseriesdb/internal/config` - Configuration management
- `timeseriesdb/test/helpers` - Test helper functions
- Standard Go testing package
- `net/http/httptest` for HTTP testing

## Best Practices

1. **Isolation**: Each test should be independent and not affect others
2. **Cleanup**: Always clean up resources after tests
3. **Realistic Data**: Use realistic test data that matches production scenarios
4. **Error Testing**: Test both success and failure cases
5. **Concurrent Testing**: Test concurrent access patterns where applicable
6. **Performance**: Keep tests fast and efficient
7. **Documentation**: Document complex test scenarios and expected behavior

## Troubleshooting

### Common Issues:
- **Port Conflicts**: Tests use random ports, but conflicts can occur
- **File Permissions**: Ensure test directories can be created and cleaned up
- **Storage Errors**: Check that storage configuration is valid
- **HTTP Timeouts**: Some tests may need longer timeouts for large data

### Debug Mode:
Run tests with verbose output and logging:
```bash
go test -v -log.level=debug ./test/integration/...
```

### Memory Leaks:
Use Go's built-in race detector:
```bash
go test -race ./test/integration/...
```
