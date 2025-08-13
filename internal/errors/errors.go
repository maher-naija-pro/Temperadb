package errors

import (
	"fmt"
	"runtime"
	"time"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeValidation represents validation errors
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypeNotFound represents not found errors
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeDatabase represents database errors
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypeStorage represents storage errors
	ErrorTypeStorage ErrorType = "storage"
	// ErrorTypeNetwork represents network errors
	ErrorTypeNetwork ErrorType = "network"
	// ErrorTypeInternal represents internal server errors
	ErrorTypeInternal ErrorType = "internal"
	// ErrorTypeTimeout represents timeout errors
	ErrorTypeTimeout ErrorType = "timeout"
)

// AppError represents an application-specific error
type AppError struct {
	Type      ErrorType
	Message   string
	Err       error
	Code      int
	Context   map[string]interface{}
	Stack     []Frame
	Timestamp int64
}

// Frame represents a stack frame
type Frame struct {
	Function string
	File     string
	Line     int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is checks if the error is of a specific type
func (e *AppError) Is(target error) bool {
	if target == nil {
		return false
	}

	if appErr, ok := target.(*AppError); ok {
		return e.Type == appErr.Type
	}

	return false
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// New creates a new AppError
func New(errType ErrorType, message string) *AppError {
	return &AppError{
		Type:      errType,
		Message:   message,
		Stack:     captureStack(),
		Timestamp: getTimestamp(),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	appErr, ok := err.(*AppError)
	if ok {
		appErr.Message = message + ": " + appErr.Message
		return appErr
	}

	return &AppError{
		Type:      ErrorTypeInternal,
		Message:   message,
		Err:       err,
		Stack:     captureStack(),
		Timestamp: getTimestamp(),
	}
}

// WrapWithType wraps an error with a specific type
func WrapWithType(err error, errType ErrorType, message string) *AppError {
	if err == nil {
		return nil
	}

	appErr := Wrap(err, message)
	appErr.Type = errType
	return appErr
}

// IsType checks if an error is of a specific type
func IsType(err error, errType ErrorType) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if ok := As(err, &appErr); ok {
		return appErr.Type == errType
	}

	return false
}

// As attempts to convert an error to an AppError
func As(err error, target **AppError) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		*target = appErr
		return true
	}

	return false
}

// captureStack captures the current stack trace
func captureStack() []Frame {
	var frames []Frame
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		frames = append(frames, Frame{
			Function: fn.Name(),
			File:     file,
			Line:     line,
		})
	}
	return frames
}

// getTimestamp returns the current timestamp
func getTimestamp() int64 {
	return time.Now().UnixNano()
}

// Common error constructors
func NewValidationError(message string) *AppError {
	return New(ErrorTypeValidation, message)
}

func NewNotFoundError(message string) *AppError {
	return New(ErrorTypeNotFound, message)
}

func NewDatabaseError(message string) *AppError {
	return New(ErrorTypeDatabase, message)
}

func NewStorageError(message string) *AppError {
	return New(ErrorTypeStorage, message)
}

func NewNetworkError(message string) *AppError {
	return New(ErrorTypeNetwork, message)
}

func NewInternalError(message string) *AppError {
	return New(ErrorTypeInternal, message)
}

func NewTimeoutError(message string) *AppError {
	return New(ErrorTypeTimeout, message)
}
