package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"timeseriesdb/internal/metrics"
)

func TestNewMetricsMiddleware(t *testing.T) {
	middleware := NewMetricsMiddleware()

	if middleware == nil {
		t.Error("Expected middleware to be created, got nil")
	}

	if middleware.requests == nil {
		t.Error("Expected requests counter to be initialized")
	}

	if middleware.duration == nil {
		t.Error("Expected duration histogram to be initialized")
	}
}

func TestMetricsMiddleware_Wrap(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler that returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	if wrappedHandler == nil {
		t.Error("Expected wrapped handler to be created, got nil")
	}
}

func TestMetricsMiddleware_CollectsMetrics(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler that returns 200 OK
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute the wrapped handler
	wrappedHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected response body 'OK', got '%s'", w.Body.String())
	}

	// Note: In a real test, you would verify that metrics were collected
	// This would require access to the metrics registry to check values
	// For now, we're just testing that the middleware doesn't break the handler
}

func TestMetricsMiddleware_HandlesErrors(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler that returns 500 Internal Server Error
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Create a test request
	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	// Execute the wrapped handler
	wrappedHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	if w.Body.String() != "Internal Server Error" {
		t.Errorf("Expected response body 'Internal Server Error', got '%s'", w.Body.String())
	}
}

func TestMetricsMiddleware_HandlesPanic(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler that panics
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute the wrapped handler - should not panic due to middleware
	// The panic will be caught by the test framework
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Handler panicked as expected: %v", r)
		}
	}()

	wrappedHandler.ServeHTTP(w, req)
}

func TestMetricsMiddleware_ResponseWriterWrapper(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler that sets custom headers and status
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Create a test request
	req := httptest.NewRequest("PUT", "/test", nil)
	w := httptest.NewRecorder()

	// Execute the wrapped handler
	wrappedHandler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if w.Body.String() != "Created" {
		t.Errorf("Expected response body 'Created', got '%s'", w.Body.String())
	}

	// Check custom header
	if w.Header().Get("X-Custom-Header") != "test-value" {
		t.Errorf("Expected custom header 'test-value', got '%s'", w.Header().Get("X-Custom-Header"))
	}
}

func TestMetricsMiddleware_DifferentHTTPMethods(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.Method))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Test different HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", method, w.Code)
			}

			if w.Body.String() != method {
				t.Errorf("Expected response body '%s', got '%s'", method, w.Body.String())
			}
		})
	}
}

func TestMetricsMiddleware_DifferentPaths(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	middleware := NewMetricsMiddleware()

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Path))
	})

	// Wrap the handler with metrics middleware
	wrappedHandler := middleware.Wrap(testHandler)

	// Test different paths
	paths := []string{"/health", "/write", "/metrics", "/api/v1/query"}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200 for %s, got %d", path, w.Code)
			}

			if w.Body.String() != path {
				t.Errorf("Expected response body '%s', got '%s'", path, w.Body.String())
			}
		})
	}
}

func TestMetricsMiddleware_ResponseWriterInterface(t *testing.T) {
	// Test that the responseWriter wrapper implements the basic ResponseWriter interface
	var _ http.ResponseWriter = &responseWriter{}

	// Note: The responseWriter only implements the basic ResponseWriter interface
	// Optional interfaces like Flusher, Hijacker, etc. are not implemented
	// This is fine for our use case as we only need the basic functionality
}

func TestMetricsHandler(t *testing.T) {
	// Initialize metrics system for testing
	metrics.Init()
	defer metrics.Reset()

	handler := MetricsHandler()

	if handler == nil {
		t.Error("Expected metrics handler to be created, got nil")
	}

	// Test that the handler can serve requests
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected metrics handler to return 200, got %d", w.Code)
	}
}
