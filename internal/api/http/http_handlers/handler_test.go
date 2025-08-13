package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseHandler_WriteJSON(t *testing.T) {
	handler := &BaseHandler{}
	w := httptest.NewRecorder()

	// Test WriteJSON method
	handler.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})

	// Check response headers
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got '%s'", w.Header().Get("Content-Type"))
	}

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBaseHandler_WriteError(t *testing.T) {
	handler := &BaseHandler{}
	w := httptest.NewRecorder()

	// Test WriteError method
	handler.WriteError(w, http.StatusBadRequest, "Bad request")

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check response body contains error message
	body := w.Body.String()
	if body != "Bad request\n" {
		t.Errorf("Expected response body 'Bad request\\n', got '%s'", body)
	}
}

func TestBaseHandler_MethodNotAllowed(t *testing.T) {
	handler := &BaseHandler{}
	w := httptest.NewRecorder()

	// Test MethodNotAllowed method
	handler.MethodNotAllowed(w, http.MethodGet, http.MethodPost)

	// Check Allow header
	allowHeader := w.Header().Get("Allow")
	expectedAllow := "GET, POST"
	if allowHeader != expectedAllow {
		t.Errorf("Expected Allow header '%s', got '%s'", expectedAllow, allowHeader)
	}

	// Check status code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	// Check response body
	body := w.Body.String()
	if body != "Method not allowed\n" {
		t.Errorf("Expected response body 'Method not allowed\\n', got '%s'", body)
	}
}

func TestBaseHandler_MethodNotAllowed_DefaultMethods(t *testing.T) {
	handler := &BaseHandler{}
	w := httptest.NewRecorder()

	// Test MethodNotAllowed method without specifying methods
	handler.MethodNotAllowed(w)

	// Check Allow header (should default to GET, POST)
	allowHeader := w.Header().Get("Allow")
	expectedAllow := "GET, POST"
	if allowHeader != expectedAllow {
		t.Errorf("Expected Allow header '%s', got '%s'", expectedAllow, allowHeader)
	}

	// Check status code
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
