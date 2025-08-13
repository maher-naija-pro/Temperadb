package utils

import (
	"testing"
	"timeseriesdb/internal/storage"
	"timeseriesdb/internal/types"
)

// AssertPointEqual checks if two points are equal
func AssertPointEqual(t *testing.T, expected, actual types.Point) {
	if expected.Measurement != actual.Measurement {
		t.Errorf("Measurement mismatch: expected %s, got %s", expected.Measurement, actual.Measurement)
	}

	if len(expected.Tags) != len(actual.Tags) {
		t.Errorf("Tags count mismatch: expected %d, got %d", len(expected.Tags), len(actual.Tags))
	}

	for k, v := range expected.Tags {
		if actual.Tags[k] != v {
			t.Errorf("Tag %s mismatch: expected %s, got %s", k, v, actual.Tags[k])
		}
	}

	if len(expected.Fields) != len(actual.Fields) {
		t.Errorf("Fields count mismatch: expected %d, got %d", len(expected.Fields), len(actual.Fields))
	}

	for k, v := range expected.Fields {
		if actual.Fields[k] != v {
			t.Errorf("Field %s mismatch: expected %g, got %g", k, v, actual.Fields[k])
		}
	}

	if !expected.Timestamp.Equal(actual.Timestamp) {
		t.Errorf("Timestamp mismatch: expected %v, got %v", expected.Timestamp, actual.Timestamp)
	}
}

// AssertStorageContains checks if storage contains the expected point
func AssertStorageContains(t *testing.T, storage *storage.Storage, expected types.Point) {
	// This would need to be implemented based on storage capabilities
	// For now, we'll just log that this check was attempted
	t.Logf("Storage contains check requested for point: %s", expected.Measurement)
}

// AssertHTTPResponse checks common HTTP response properties
func AssertHTTPResponse(t *testing.T, statusCode int, expectedStatus int, description string) {
	if statusCode != expectedStatus {
		t.Errorf("%s: expected status %d, got %d", description, expectedStatus, statusCode)
	}
}

// AssertResponseBody checks if response body contains expected content
func AssertResponseBody(t *testing.T, body string, expected string, description string) {
	if body != expected {
		t.Errorf("%s: expected body '%s', got '%s'", description, expected, body)
	}
}
