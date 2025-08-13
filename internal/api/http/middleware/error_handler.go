package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"timeseriesdb/internal/errors"
	"timeseriesdb/internal/logger"
)

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Type    string                 `json:"type,omitempty"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// ErrorHandler middleware handles errors and provides consistent error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v\nStack trace: %s", err, debug.Stack())

				// Return 500 Internal Server Error for panics
				writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error", "internal", nil)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// HandleError handles application errors and writes appropriate HTTP responses
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var appErr *errors.AppError
	if ok := errors.As(err, &appErr); ok {
		// Handle application-specific errors
		statusCode := getStatusCodeForErrorType(appErr.Type)
		writeErrorResponse(w, statusCode, appErr.Message, string(appErr.Type), appErr.Context)
		return
	}

	// Handle generic errors
	logger.Errorf("Unhandled error: %v", err)
	writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error", "internal", nil)
}

// writeErrorResponse writes a JSON error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, message, errorType string, context map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   errorType,
		Type:    errorType,
		Code:    statusCode,
		Message: message,
		Context: context,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode error response: %v", err)
		// Fallback to plain text
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}
}

// getStatusCodeForErrorType maps error types to HTTP status codes
func getStatusCodeForErrorType(errorType errors.ErrorType) int {
	switch errorType {
	case errors.ErrorTypeValidation:
		return http.StatusBadRequest
	case errors.ErrorTypeNotFound:
		return http.StatusNotFound
	case errors.ErrorTypeDatabase, errors.ErrorTypeStorage:
		return http.StatusServiceUnavailable
	case errors.ErrorTypeNetwork:
		return http.StatusBadGateway
	case errors.ErrorTypeTimeout:
		return http.StatusRequestTimeout
	case errors.ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// LogError logs an error with appropriate context
func LogError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	var appErr *errors.AppError
	if ok := errors.As(err, &appErr); ok {
		// Log application error with context
		logger.WithFields(context).Errorf("Application error [%s]: %s", appErr.Type, appErr.Message)
		if appErr.Err != nil {
			logger.WithFields(context).Errorf("Wrapped error: %v", appErr.Err)
		}
	} else {
		// Log generic error
		logger.WithFields(context).Errorf("Generic error: %v", err)
	}
}
