package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler_Handle(t *testing.T) {
	handler := NewHealthHandler()

	// Test GET request (should succeed)
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"TimeSeriesDB"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Test POST request (should fail)
	req, err = http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestHealthHandler_Handle_AllHTTPMethods(t *testing.T) {
	handler := NewHealthHandler()
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			if method == "GET" {
				// GET should succeed
				if status := rr.Code; status != http.StatusOK {
					t.Errorf("GET request returned wrong status code: got %v want %v", status, http.StatusOK)
				}

				expected := `{"status":"healthy","service":"TimeSeriesDB"}`
				if rr.Body.String() != expected {
					t.Errorf("GET request returned unexpected body: got %v want %v", rr.Body.String(), expected)
				}

				// Check Content-Type header
				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			} else {
				// All other methods should return 405 Method Not Allowed
				if status := rr.Code; status != http.StatusMethodNotAllowed {
					t.Errorf("%s request returned wrong status code: got %v want %v", method, status, http.StatusMethodNotAllowed)
				}

				// Check Allow header - should contain GET and POST
				allowHeader := rr.Header().Get("Allow")
				if allowHeader != "GET, POST" {
					t.Errorf("Expected Allow header 'GET, POST', got '%s'", allowHeader)
				}
			}
		})
	}
}

func TestHealthHandler_Handle_WithQueryParameters(t *testing.T) {
	handler := NewHealthHandler()

	// Test GET request with query parameters (should still succeed)
	req, err := http.NewRequest("GET", "/health?format=json&detailed=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"TimeSeriesDB"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHealthHandler_Handle_WithHeaders(t *testing.T) {
	handler := NewHealthHandler()

	// Test GET request with custom headers (should still succeed)
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add custom headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "test-client")
	req.Header.Set("X-Request-ID", "test-123")

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"TimeSeriesDB"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check Content-Type header
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHealthHandler_Handle_WithBody(t *testing.T) {
	handler := NewHealthHandler()

	// Test GET request with body (should still succeed)
	req, err := http.NewRequest("GET", "/health", strings.NewReader("some body content"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"TimeSeriesDB"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHealthHandler_Handle_EmptyPath(t *testing.T) {
	handler := NewHealthHandler()

	// Test GET request with empty path (should still succeed)
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"healthy","service":"TimeSeriesDB"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestHealthHandler_Handle_ConcurrentRequests(t *testing.T) {
	handler := NewHealthHandler()
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req, err := http.NewRequest("GET", "/health", nil)
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
				done <- false
				return
			}

			rr := httptest.NewRecorder()
			handler.Handle(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			expected := `{"status":"healthy","service":"TimeSeriesDB"}`
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
			}

			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		if !<-done {
			t.Fatal("One or more concurrent requests failed")
		}
	}
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	if handler == nil {
		t.Error("Expected health handler to be created, got nil")
	}

	// Test that the handler can handle requests
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.Handle(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("New handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHealthHandler_ResponseConsistency(t *testing.T) {
	handler := NewHealthHandler()

	// Make multiple requests to ensure consistent responses
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.Handle(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Request %d returned wrong status code: got %v want %v", i+1, status, http.StatusOK)
		}

		expected := `{"status":"healthy","service":"TimeSeriesDB"}`
		if rr.Body.String() != expected {
			t.Errorf("Request %d returned unexpected body: got %v want %v", i+1, rr.Body.String(), expected)
		}

		// Check Content-Type header consistency
		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Request %d returned wrong Content-Type: got '%s' want 'application/json'", i+1, contentType)
		}
	}
}
