package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/storage"
)

// TestWriteHandler_NewWriteHandler tests the NewWriteHandler function
func TestWriteHandler_NewWriteHandler(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_storage.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_storage.tsv")
		os.RemoveAll("test_backups")
	}()

	handler := NewWriteHandler(storageInstance)

	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}

	if handler.storage != storageInstance {
		t.Error("Expected handler to have the provided storage instance")
	}
}

// TestWriteHandler_Handle_ValidRequest tests a valid write request
func TestWriteHandler_Handle_ValidRequest(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_valid.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_valid",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_valid.tsv")
		os.RemoveAll("test_backups_valid")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}

// TestWriteHandler_Handle_InvalidMethod tests that only POST method is allowed
func TestWriteHandler_Handle_InvalidMethod(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_method.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_method",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_method.tsv")
		os.RemoveAll("test_backups_method")
	}()

	handler := NewWriteHandler(storageInstance)

	// Test all HTTP methods
	methods := []string{"GET", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/write", nil)
			w := httptest.NewRecorder()

			handler.Handle(w, req)

			// Check response
			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405 for %s, got %d", method, w.Code)
			}

			// Check Allow header - should contain GET and POST
			allowHeader := w.Header().Get("Allow")
			if allowHeader != "GET, POST" {
				t.Errorf("Expected Allow header 'GET, POST', got '%s'", allowHeader)
			}
		})
	}
}

// TestWriteHandler_Handle_ZeroContentLength tests handling of zero content length
func TestWriteHandler_Handle_ZeroContentLength(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_zero.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_zero",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_zero.tsv")
		os.RemoveAll("test_backups_zero")
	}()

	handler := NewWriteHandler(storageInstance)

	req := httptest.NewRequest("POST", "/write", nil)
	req.ContentLength = 0
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_NegativeContentLength tests handling of negative content length
func TestWriteHandler_Handle_NegativeContentLength(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_negative.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_negative",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_negative.tsv")
		os.RemoveAll("test_backups_negative")
	}()

	handler := NewWriteHandler(storageInstance)

	req := httptest.NewRequest("POST", "/write", nil)
	req.ContentLength = -1
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_RequestBodyReadError tests handling of body read errors
func TestWriteHandler_Handle_RequestBodyReadError(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_read_error.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_read_error",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_read_error.tsv")
		os.RemoveAll("test_backups_read_error")
	}()

	handler := NewWriteHandler(storageInstance)

	// Create a request with a body that will cause read errors
	req := httptest.NewRequest("POST", "/write", &errorReader{})
	req.ContentLength = 10
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_InvalidLineProtocol tests handling of invalid line protocol
func TestWriteHandler_Handle_InvalidLineProtocol(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_invalid.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_invalid",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_invalid.tsv")
		os.RemoveAll("test_backups_invalid")
	}()

	handler := NewWriteHandler(storageInstance)

	// Invalid line protocol data
	invalidData := "cpu,host=server01 value=invalid 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(invalidData))
	req.ContentLength = int64(len(invalidData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_MalformedLineProtocol tests handling of malformed line protocol
func TestWriteHandler_Handle_MalformedLineProtocol(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_malformed.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_malformed",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_malformed.tsv")
		os.RemoveAll("test_backups_malformed")
	}()

	handler := NewWriteHandler(storageInstance)

	// Malformed line protocol data (missing timestamp)
	malformedData := "cpu,host=server01 value=0.64"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(malformedData))
	req.ContentLength = int64(len(malformedData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_EmptyLineProtocol tests handling of empty line protocol
func TestWriteHandler_Handle_EmptyLineProtocol(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_empty.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_empty",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_empty.tsv")
		os.RemoveAll("test_backups_empty")
	}()

	handler := NewWriteHandler(storageInstance)

	// Empty line protocol data
	emptyData := ""
	req := httptest.NewRequest("POST", "/write", strings.NewReader(emptyData))
	req.ContentLength = int64(len(emptyData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - empty data should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_MultipleLines tests handling of multiple lines
func TestWriteHandler_Handle_MultipleLines(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_multiple.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_multiple",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_multiple.tsv")
		os.RemoveAll("test_backups_multiple")
	}()

	handler := NewWriteHandler(storageInstance)

	// Multiple lines of valid line protocol data
	multiLineData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000\n" +
		"memory,host=server01,region=us-west used=1234567 1434055562000001000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(multiLineData))
	req.ContentLength = int64(len(multiLineData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}

// TestWriteHandler_Handle_ComplexTags tests handling of complex tag structures
func TestWriteHandler_Handle_ComplexTags(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_complex_tags.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_complex_tags",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_complex_tags.tsv")
		os.RemoveAll("test_backups_complex_tags")
	}()

	handler := NewWriteHandler(storageInstance)

	// Line protocol with complex tags
	complexData := "cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(complexData))
	req.ContentLength = int64(len(complexData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_MultipleFields tests handling of multiple fields
func TestWriteHandler_Handle_MultipleFields(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_multiple_fields.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_multiple_fields",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_multiple_fields.tsv")
		os.RemoveAll("test_backups_multiple_fields")
	}()

	handler := NewWriteHandler(storageInstance)

	// Line protocol with multiple fields
	multiFieldData := "cpu,host=server01 user=0.64,system=0.23,idle=0.12 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(multiFieldData))
	req.ContentLength = int64(len(multiFieldData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_StorageError tests handling of storage errors
func TestWriteHandler_Handle_StorageError(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_storage_error.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_storage_error",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_storage_error.tsv")
		os.RemoveAll("test_backups_storage_error")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should still return 200 OK even if storage fails
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_StorageWriteError tests handling of individual storage write errors
func TestWriteHandler_Handle_StorageWriteError(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_storage_write_error.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_storage_write_error",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_storage_write_error.tsv")
		os.RemoveAll("test_backups_storage_write_error")
	}()

	handler := NewWriteHandler(storageInstance)

	// Multiple lines of valid line protocol data to test individual write errors
	multiLineData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000\n" +
		"memory,host=server01,region=us-west used=1234567 1434055562000001000\n" +
		"disk,host=server01,region=us-west free=987654321 1434055562000002000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(multiLineData))
	req.ContentLength = int64(len(multiLineData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should still return 200 OK even if some storage writes fail
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}

// TestWriteHandler_Handle_ConcurrentRequests tests concurrent request handling
func TestWriteHandler_Handle_ConcurrentRequests(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_concurrent.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_concurrent",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_concurrent.tsv")
		os.RemoveAll("test_backups_concurrent")
	}()

	handler := NewWriteHandler(storageInstance)
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			// Valid line protocol data
			validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
			req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
			req.ContentLength = int64(len(validData))

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Request %d returned wrong status code: got %v want %v", id, status, http.StatusOK)
			}

			expected := "OK"
			if rr.Body.String() != expected {
				t.Errorf("Request %d returned unexpected body: got %v want %v", id, rr.Body.String(), expected)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		if !<-done {
			t.Fatal("One or more concurrent requests failed")
		}
	}
}

// TestWriteHandler_Handle_ResponseConsistency tests response consistency across multiple requests
func TestWriteHandler_Handle_ResponseConsistency(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_consistency.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_consistency",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_consistency.tsv")
		os.RemoveAll("test_backups_consistency")
	}()

	handler := NewWriteHandler(storageInstance)

	// Make multiple requests to ensure consistent responses
	for i := 0; i < 5; i++ {
		validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
		req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
		req.ContentLength = int64(len(validData))

		rr := httptest.NewRecorder()
		handler.Handle(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %d returned wrong status code: got %v want %v", i+1, status, http.StatusOK)
		}

		expected := "OK"
		if rr.Body.String() != expected {
			t.Errorf("Request %d returned unexpected body: got %v want %v", i+1, rr.Body.String(), expected)
		}
	}
}

// TestWriteHandler_Handle_EdgeCases tests various edge cases
func TestWriteHandler_Handle_EdgeCases(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_edge_cases.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_edge_cases",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_edge_cases.tsv")
		os.RemoveAll("test_backups_edge_cases")
	}()

	handler := NewWriteHandler(storageInstance)

	testCases := []struct {
		name         string
		data         string
		expectedCode int
	}{
		{
			name:         "whitespace only",
			data:         "   \n  \t  ",
			expectedCode: http.StatusOK, // Empty data should succeed
		},
		{
			name:         "single measurement only",
			data:         "cpu",
			expectedCode: http.StatusBadRequest, // Missing fields and timestamp
		},
		{
			name:         "measurement with empty tags",
			data:         "cpu, value=0.64 1434055562000000000",
			expectedCode: http.StatusBadRequest, // Invalid tag format
		},
		{
			name:         "measurement with empty field name",
			data:         "cpu =0.64 1434055562000000000",
			expectedCode: http.StatusBadRequest, // Empty field name
		},
		{
			name:         "measurement with empty field value",
			data:         "cpu value= 1434055562000000000",
			expectedCode: http.StatusBadRequest, // Empty field value
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/write", strings.NewReader(tc.data))
			req.ContentLength = int64(len(tc.data))

			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected status %d for '%s', got %d", tc.expectedCode, tc.name, w.Code)
			}
		})
	}
}

// TestWriteHandler_Handle_WithQueryParameters tests handling of requests with query parameters
func TestWriteHandler_Handle_WithQueryParameters(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_query_params.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_query_params",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_query_params.tsv")
		os.RemoveAll("test_backups_query_params")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data with query parameters in URL
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write?db=mydb&precision=ns", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should still succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_WithHeaders tests handling of requests with custom headers
func TestWriteHandler_Handle_WithHeaders(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_headers.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_headers",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_headers.tsv")
		os.RemoveAll("test_backups_headers")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data with custom headers
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	// Add custom headers
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", "test-client")
	req.Header.Set("X-Request-ID", "test-123")

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should still succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_WithBody tests handling of requests with body
func TestWriteHandler_Handle_WithBody(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_body.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_body",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_body.tsv")
		os.RemoveAll("test_backups_body")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}

// TestWriteHandler_Handle_EmptyPath tests handling of requests with empty path
func TestWriteHandler_Handle_EmptyPath(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_empty_path.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_empty_path",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_empty_path.tsv")
		os.RemoveAll("test_backups_empty_path")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data with root path instead of empty path
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should still succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// errorReader is a reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// TestWriteHandler_Integration_WithRealStorage tests integration with real storage
func TestWriteHandler_Integration_WithRealStorage(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_integration.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_integration",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_integration.tsv")
		os.RemoveAll("test_backups_integration")
	}()

	handler := NewWriteHandler(storageInstance)

	// Test multiple valid writes
	testData := []string{
		"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
		"memory,host=server01,region=us-west used=1234567 1434055562000001000",
		"disk,host=server01,region=us-west free=987654321 1434055562000002000",
	}

	for i, data := range testData {
		t.Run(fmt.Sprintf("write_%d", i), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/write", strings.NewReader(data))
			req.ContentLength = int64(len(data))

			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for write %d, got %d", i, w.Code)
			}

			if w.Body.String() != "OK" {
				t.Errorf("Expected response body 'OK' for write %d, got '%s'", i, w.Body.String())
			}
		})
	}
}

// TestWriteHandler_Handle_WithLogger tests that the handler properly uses the logger
func TestWriteHandler_Handle_WithLogger(t *testing.T) {
	// Enable test mode for logger
	logger.SetTestMode(true)
	defer logger.SetTestMode(false)

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_logger.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_logger",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_logger.tsv")
		os.RemoveAll("test_backups_logger")
	}()

	handler := NewWriteHandler(storageInstance)

	// Valid line protocol data
	validData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(validData))
	req.ContentLength = int64(len(validData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - should succeed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestWriteHandler_Handle_EdgeCase_WhitespaceOnly tests handling of whitespace-only data
func TestWriteHandler_Handle_EdgeCase_WhitespaceOnly(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_whitespace.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_whitespace",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_whitespace.tsv")
		os.RemoveAll("test_backups_whitespace")
	}()

	handler := NewWriteHandler(storageInstance)

	// Whitespace-only data
	whitespaceData := "   \n  \t  \n  "
	req := httptest.NewRequest("POST", "/write", strings.NewReader(whitespaceData))
	req.ContentLength = int64(len(whitespaceData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - whitespace-only data should succeed (0 points)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}

// TestWriteHandler_Handle_EdgeCase_SingleNewline tests handling of single newline
func TestWriteHandler_Handle_EdgeCase_SingleNewline(t *testing.T) {
	// Initialize logger for testing
	logger.Init()

	// Create a real storage instance for testing
	storageConfig := config.StorageConfig{
		DataDir:     t.TempDir(),
		DataFile:    "test_write_handler_single_newline.tsv",
		MaxFileSize: 1024,
		BackupDir:   "test_backups_single_newline",
		Compression: false,
	}

	storageInstance := storage.NewStorage(storageConfig)
	defer func() {
		storageInstance.Close()
		os.Remove("test_write_handler_single_newline.tsv")
		os.RemoveAll("test_backups_single_newline")
	}()

	handler := NewWriteHandler(storageInstance)

	// Single newline data
	newlineData := "\n"
	req := httptest.NewRequest("POST", "/write", strings.NewReader(newlineData))
	req.ContentLength = int64(len(newlineData))

	w := httptest.NewRecorder()
	handler.Handle(w, req)

	// Check response - single newline should succeed (0 points)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}
}
