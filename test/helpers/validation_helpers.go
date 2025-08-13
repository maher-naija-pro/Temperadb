package helpers

import (
	"reflect"
	"strings"
	"testing"
	"timeseriesdb/internal/types"
)

// ValidationHelpers provides common validation functions for tests
type ValidationHelpers struct{}

// ValidatePointStructure validates that a point has the expected structure
func (h *ValidationHelpers) ValidatePointStructure(t *testing.T, point types.Point, expectedMeasurement string) {
	if point.Measurement != expectedMeasurement {
		t.Errorf("Expected measurement '%s', got '%s'", expectedMeasurement, point.Measurement)
	}

	if point.Tags == nil {
		t.Error("Point tags should not be nil")
	}

	if point.Fields == nil {
		t.Error("Point fields should not be nil")
	}

	if point.Timestamp.IsZero() {
		t.Error("Point timestamp should not be zero")
	}
}

// ValidateTags validates that tags contain expected key-value pairs
func (h *ValidationHelpers) ValidateTags(t *testing.T, tags map[string]string, expected map[string]string) {
	if len(tags) != len(expected) {
		t.Errorf("Expected %d tags, got %d", len(expected), len(tags))
		return
	}

	for key, expectedValue := range expected {
		if actualValue, exists := tags[key]; !exists {
			t.Errorf("Expected tag '%s' not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected tag '%s' to be '%s', got '%s'", key, expectedValue, actualValue)
		}
	}
}

// ValidateFields validates that fields contain expected key-value pairs
func (h *ValidationHelpers) ValidateFields(t *testing.T, fields map[string]float64, expected map[string]float64) {
	if len(fields) != len(expected) {
		t.Errorf("Expected %d fields, got %d", len(expected), len(fields))
		return
	}

	for key, expectedValue := range expected {
		if actualValue, exists := fields[key]; !exists {
			t.Errorf("Expected field '%s' not found", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected field '%s' to be %g, got %g", key, expectedValue, actualValue)
		}
	}
}

// ValidateSliceLength validates that a slice has the expected length
func (h *ValidationHelpers) ValidateSliceLength(t *testing.T, slice interface{}, expectedLength int) {
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		t.Errorf("Expected slice, got %T", slice)
		return
	}

	if sliceValue.Len() != expectedLength {
		t.Errorf("Expected slice length %d, got %d", expectedLength, sliceValue.Len())
	}
}

// ValidateMapLength validates that a map has the expected length
func (h *ValidationHelpers) ValidateMapLength(t *testing.T, m interface{}, expectedLength int) {
	mapValue := reflect.ValueOf(m)
	if mapValue.Kind() != reflect.Map {
		t.Errorf("Expected map, got %T", m)
		return
	}

	if mapValue.Len() != expectedLength {
		t.Errorf("Expected map length %d, got %d", expectedLength, mapValue.Len())
	}
}

// ValidateStringContains validates that a string contains expected substring
func (h *ValidationHelpers) ValidateStringContains(t *testing.T, str, expected string) {
	if !strings.Contains(str, expected) {
		t.Errorf("Expected string to contain '%s', got '%s'", expected, str)
	}
}

// ValidateStringEquals validates that a string equals expected value
func (h *ValidationHelpers) ValidateStringEquals(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, actual)
	}
}

// ValidateIntEquals validates that an integer equals expected value
func (h *ValidationHelpers) ValidateIntEquals(t *testing.T, actual, expected int) {
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

// ValidateFloatEquals validates that a float equals expected value with tolerance
func (h *ValidationHelpers) ValidateFloatEquals(t *testing.T, actual, expected, tolerance float64) {
	diff := actual - expected
	if diff < 0 {
		diff = -diff
	}

	if diff > tolerance {
		t.Errorf("Expected %g ± %g, got %g", expected, tolerance, actual)
	}
}

// Global instance for easy access
var Validation = &ValidationHelpers{}
