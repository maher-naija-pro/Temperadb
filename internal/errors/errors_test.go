package errors

import (
	"testing"
)

func TestNewAppError(t *testing.T) {
	tests := []struct {
		name        string
		errorType   ErrorType
		message     string
		expectedMsg string
	}{
		{
			name:        "Validation error",
			errorType:   ErrorTypeValidation,
			message:     "invalid input",
			expectedMsg: "invalid input",
		},
		{
			name:        "Database error",
			errorType:   ErrorTypeDatabase,
			message:     "connection failed",
			expectedMsg: "connection failed",
		},
		{
			name:        "Storage error",
			errorType:   ErrorTypeStorage,
			message:     "disk full",
			expectedMsg: "disk full",
		},
		{
			name:        "Network error",
			errorType:   ErrorTypeNetwork,
			message:     "timeout",
			expectedMsg: "timeout",
		},
		{
			name:        "Internal error",
			errorType:   ErrorTypeInternal,
			message:     "unexpected state",
			expectedMsg: "unexpected state",
		},
		{
			name:        "Timeout error",
			errorType:   ErrorTypeTimeout,
			message:     "operation timed out",
			expectedMsg: "operation timed out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.errorType, tt.message)

			if err == nil {
				t.Fatal("Expected error but got nil")
			}

			// err is already *AppError, no need for type assertion
			if err.Type != tt.errorType {
				t.Errorf("Expected error type %s, got %s", tt.errorType, err.Type)
			}

			if err.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, err.Message)
			}

			if err.Timestamp == 0 {
				t.Error("Expected timestamp to be set")
			}

			if len(err.Stack) == 0 {
				t.Error("Expected stack trace to be captured")
			}
		})
	}
}

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name        string
		errorType   ErrorType
		message     string
		expectedMsg string
	}{
		{
			name:        "Simple error",
			errorType:   ErrorTypeValidation,
			message:     "invalid input",
			expectedMsg: "invalid input",
		},
		{
			name:        "Error with wrapped error",
			errorType:   ErrorTypeDatabase,
			message:     "connection failed",
			expectedMsg: "connection failed: network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err *AppError
			if tt.errorType == ErrorTypeDatabase {
				err = &AppError{
					Type:    tt.errorType,
					Message: tt.message,
					Err:     &AppError{Message: "network error"},
				}
			} else {
				err = New(tt.errorType, tt.message)
			}

			msg := err.Error()

			if msg != tt.expectedMsg {
				t.Errorf("Expected error message '%s', got '%s'", tt.expectedMsg, msg)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := &AppError{Message: "original error"}
	err := &AppError{
		Type: ErrorTypeDatabase,
		Err:  cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Expected unwrapped error %v, got %v", cause, unwrapped)
	}
}

func TestAppError_Is(t *testing.T) {
	err1 := New(ErrorTypeValidation, "error 1")
	err2 := New(ErrorTypeValidation, "error 2")
	err3 := New(ErrorTypeDatabase, "error 3")

	if !err1.Is(err2) {
		t.Error("Expected errors of same type to be equal")
	}

	if err1.Is(err3) {
		t.Error("Expected errors of different types to not be equal")
	}
}

func TestAppError_As(t *testing.T) {
	err := New(ErrorTypeInternal, "internal error")

	var appErr *AppError
	if !As(err, &appErr) {
		t.Error("Expected As to return true for *AppError")
	}

	if appErr.Type != ErrorTypeInternal {
		t.Errorf("Expected error type %s, got %s", ErrorTypeInternal, appErr.Type)
	}
}

func TestAppError_WithContext(t *testing.T) {
	err := New(ErrorTypeValidation, "invalid input")

	withCtx := err.WithContext("field", "username")

	// WithContext returns the same error instance
	if withCtx != err {
		t.Error("Expected WithContext to return the same error instance")
	}

	if withCtx.Context["field"] != "username" {
		t.Errorf("Expected context field 'username', got %v", withCtx.Context["field"])
	}
}

func TestWrap(t *testing.T) {
	original := &AppError{
		Type:    ErrorTypeValidation,
		Message: "original error",
	}
	wrapped := Wrap(original, "wrapped message")

	// When wrapping an AppError, Wrap modifies the original and returns it
	if wrapped != original {
		t.Error("Expected Wrap to return the modified original error when wrapping AppError")
	}

	// Wrap doesn't change the type when wrapping an AppError
	if wrapped.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, wrapped.Type)
	}

	if wrapped.Message != "wrapped message: original error" {
		t.Errorf("Expected message 'wrapped message: original error', got '%s'", wrapped.Message)
	}

	// When wrapping an AppError, the Err field is not set
	if wrapped.Err != nil {
		t.Errorf("Expected Err to be nil when wrapping AppError, got %v", wrapped.Err)
	}
}

func TestWrapWithType(t *testing.T) {
	original := &AppError{Message: "original error"}
	wrapped := WrapWithType(original, ErrorTypeStorage, "storage error")

	// When wrapping an AppError, WrapWithType modifies the original and returns it
	if wrapped != original {
		t.Error("Expected WrapWithType to return the modified original error when wrapping AppError")
	}

	if wrapped.Type != ErrorTypeStorage {
		t.Errorf("Expected error type %s, got %s", ErrorTypeStorage, wrapped.Type)
	}

	if wrapped.Message != "storage error: original error" {
		t.Errorf("Expected message 'storage error: original error', got '%s'", wrapped.Message)
	}

	// When wrapping an AppError, the Err field is not set
	if wrapped.Err != nil {
		t.Errorf("Expected Err to be nil when wrapping AppError, got %v", wrapped.Err)
	}
}

func TestIsType(t *testing.T) {
	err := New(ErrorTypeValidation, "invalid input")

	if !IsType(err, ErrorTypeValidation) {
		t.Error("Expected IsType to return true for validation error")
	}

	if IsType(err, ErrorTypeDatabase) {
		t.Error("Expected IsType to return false for database error")
	}
}

func TestAs(t *testing.T) {
	err := New(ErrorTypeInternal, "internal error")

	var appErr *AppError
	if !As(err, &appErr) {
		t.Error("Expected As to return true for *AppError")
	}

	if appErr.Type != ErrorTypeInternal {
		t.Errorf("Expected error type %s, got %s", ErrorTypeInternal, appErr.Type)
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %s, got %s", ErrorTypeValidation, err.Type)
	}

	if err.Message != "invalid input" {
		t.Errorf("Expected message 'invalid input', got '%s'", err.Message)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("resource not found")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeNotFound {
		t.Errorf("Expected error type %s, got %s", ErrorTypeNotFound, err.Type)
	}

	if err.Message != "resource not found" {
		t.Errorf("Expected message 'resource not found', got '%s'", err.Message)
	}
}

func TestNewDatabaseError(t *testing.T) {
	err := NewDatabaseError("connection failed")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeDatabase {
		t.Errorf("Expected error type %s, got %s", ErrorTypeDatabase, err.Type)
	}

	if err.Message != "connection failed" {
		t.Errorf("Expected message 'connection failed', got '%s'", err.Message)
	}
}

func TestNewStorageError(t *testing.T) {
	err := NewStorageError("disk full")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeStorage {
		t.Errorf("Expected error type %s, got %s", ErrorTypeStorage, err.Type)
	}

	if err.Message != "disk full" {
		t.Errorf("Expected message 'disk full', got '%s'", err.Message)
	}
}

func TestNewNetworkError(t *testing.T) {
	err := NewNetworkError("timeout")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %s, got %s", ErrorTypeNetwork, err.Type)
	}

	if err.Message != "timeout" {
		t.Errorf("Expected message 'timeout', got '%s'", err.Message)
	}
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("unexpected state")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeInternal {
		t.Errorf("Expected error type %s, got %s", ErrorTypeInternal, err.Type)
	}

	if err.Message != "unexpected state" {
		t.Errorf("Expected message 'unexpected state', got '%s'", err.Message)
	}
}

func TestNewTimeoutError(t *testing.T) {
	err := NewTimeoutError("operation timed out")

	// err is already *AppError, no need for type assertion
	if err.Type != ErrorTypeTimeout {
		t.Errorf("Expected error type %s, got %s", ErrorTypeTimeout, err.Type)
	}

	if err.Message != "operation timed out" {
		t.Errorf("Expected message 'operation timed out', got '%s'", err.Message)
	}
}

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errorType   ErrorType
		expectedStr string
	}{
		{ErrorTypeValidation, "validation"},
		{ErrorTypeNotFound, "not_found"},
		{ErrorTypeDatabase, "database"},
		{ErrorTypeStorage, "storage"},
		{ErrorTypeNetwork, "network"},
		{ErrorTypeInternal, "internal"},
		{ErrorTypeTimeout, "timeout"},
		{ErrorType("UNKNOWN"), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(string(tt.errorType), func(t *testing.T) {
			if string(tt.errorType) != tt.expectedStr {
				t.Errorf("Expected string '%s', got '%s'", tt.expectedStr, string(tt.errorType))
			}
		})
	}
}

func TestAppError_ContextOperations(t *testing.T) {
	err := New(ErrorTypeValidation, "test error")

	// Test adding context
	err = err.WithContext("key1", "value1")
	err = err.WithContext("key2", 42)

	// Test adding more context
	err = err.WithContext("key3", true)

	// Verify all context keys are present
	expectedKeys := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	for k, v := range expectedKeys {
		if err.Context[k] != v {
			t.Errorf("Expected context key '%s' to have value %v, got %v", k, v, err.Context[k])
		}
	}
}

func TestAppError_StackCapture(t *testing.T) {
	err := New(ErrorTypeInternal, "test error")

	// Verify stack trace is captured
	if len(err.Stack) == 0 {
		t.Error("Expected stack trace to be captured")
	}

	// Verify stack trace contains some function names
	found := false
	for _, frame := range err.Stack {
		if frame.Function != "" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected stack trace to contain function names")
	}

	// Verify stack trace contains file information
	found = false
	for _, frame := range err.Stack {
		if frame.File != "" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected stack trace to contain file information")
	}
}
