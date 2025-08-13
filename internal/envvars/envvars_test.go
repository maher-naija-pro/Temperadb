package envvars

import (
	"os"
	"testing"
	"time"
)

func TestParser_String(t *testing.T) {
	parser := NewParser()

	// Test with default value
	result := parser.String("NONEXISTENT_KEY", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Test with environment variable set
	os.Setenv("TEST_STRING", "test_value")
	defer os.Unsetenv("TEST_STRING")

	result = parser.String("TEST_STRING", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with whitespace
	os.Setenv("TEST_STRING_WS", "  test_value  ")
	defer os.Unsetenv("TEST_STRING_WS")

	result = parser.String("TEST_STRING_WS", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestParser_Int(t *testing.T) {
	parser := NewParser()

	// Test with default value
	result := parser.Int("NONEXISTENT_KEY", 42)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with valid environment variable
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	result = parser.Int("TEST_INT", 42)
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}

	// Test with invalid environment variable (should fall back to default)
	os.Setenv("TEST_INT_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_INT_INVALID")

	result = parser.Int("TEST_INT_INVALID", 42)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestParser_Int64(t *testing.T) {
	parser := NewParser()

	// Test with default value
	result := parser.Int64("NONEXISTENT_KEY", 42)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test with valid environment variable
	os.Setenv("TEST_INT64", "9223372036854775807")
	defer os.Unsetenv("TEST_INT64")

	result = parser.Int64("TEST_INT64", 42)
	if result != 9223372036854775807 {
		t.Errorf("Expected 9223372036854775807, got %d", result)
	}
}

func TestParser_Bool(t *testing.T) {
	parser := NewParser()

	// Test with default value
	result := parser.Bool("NONEXISTENT_KEY", true)
	if !result {
		t.Errorf("Expected true, got %t", result)
	}

	// Test with valid true values
	trueValues := []string{"true", "TRUE", "True", "1", "yes", "YES", "Yes"}
	for _, value := range trueValues {
		os.Setenv("TEST_BOOL", value)
		result = parser.Bool("TEST_BOOL", false)
		if !result {
			t.Errorf("Expected true for '%s', got %t", value, result)
		}
	}

	// Test with valid false values
	falseValues := []string{"false", "FALSE", "False", "0", "no", "NO", "No"}
	for _, value := range falseValues {
		os.Setenv("TEST_BOOL", value)
		result = parser.Bool("TEST_BOOL", true)
		if result {
			t.Errorf("Expected false for '%s', got %t", value, result)
		}
	}

	os.Unsetenv("TEST_BOOL")
}

func TestParser_Duration(t *testing.T) {
	parser := NewParser()

	// Test with default value
	defaultDuration := 30 * time.Second
	result := parser.Duration("NONEXISTENT_KEY", defaultDuration)
	if result != defaultDuration {
		t.Errorf("Expected %v, got %v", defaultDuration, result)
	}

	// Test with valid environment variable
	os.Setenv("TEST_DURATION", "60")
	defer os.Unsetenv("TEST_DURATION")

	expectedDuration := 60 * time.Second
	result = parser.Duration("TEST_DURATION", defaultDuration)
	if result != expectedDuration {
		t.Errorf("Expected %v, got %v", expectedDuration, result)
	}

	// Test with invalid environment variable (should fall back to default)
	os.Setenv("TEST_DURATION_INVALID", "not_a_number")
	defer os.Unsetenv("TEST_DURATION_INVALID")

	result = parser.Duration("TEST_DURATION_INVALID", defaultDuration)
	if result != defaultDuration {
		t.Errorf("Expected %v, got %v", defaultDuration, result)
	}
}

func TestParser_FileSize(t *testing.T) {
	parser := NewParser()

	// Test with default value
	result := parser.FileSize("NONEXISTENT_KEY", 1024)
	if result != 1024 {
		t.Errorf("Expected 1024, got %d", result)
	}

	// Test with bytes
	os.Setenv("TEST_SIZE", "2048")
	defer os.Unsetenv("TEST_SIZE")

	result = parser.FileSize("TEST_SIZE", 1024)
	if result != 2048 {
		t.Errorf("Expected 2048, got %d", result)
	}

	// Test with KB
	os.Setenv("TEST_SIZE_KB", "1KB")
	defer os.Unsetenv("TEST_SIZE_KB")

	result = parser.FileSize("TEST_SIZE_KB", 1024)
	if result != 1024 {
		t.Errorf("Expected 1024, got %d", result)
	}

	// Test with MB
	os.Setenv("TEST_SIZE_MB", "2MB")
	defer os.Unsetenv("TEST_SIZE_MB")

	result = parser.FileSize("TEST_SIZE_MB", 1024)
	if result != 2*1024*1024 {
		t.Errorf("Expected %d, got %d", 2*1024*1024, result)
	}

	// Test with GB
	os.Setenv("TEST_SIZE_GB", "1GB")
	defer os.Unsetenv("TEST_SIZE_GB")

	result = parser.FileSize("TEST_SIZE_GB", 1024)
	if result != 1024*1024*1024 {
		t.Errorf("Expected %d, got %d", 1024*1024*1024, result)
	}

	// Test with invalid format (should fall back to default)
	os.Setenv("TEST_SIZE_INVALID", "invalid")
	defer os.Unsetenv("TEST_SIZE_INVALID")

	result = parser.FileSize("TEST_SIZE_INVALID", 1024)
	if result != 1024 {
		t.Errorf("Expected 1024, got %d", result)
	}
}

func TestParser_Has(t *testing.T) {
	parser := NewParser()

	// Test with non-existent key
	if parser.Has("NONEXISTENT_KEY") {
		t.Error("Expected false for non-existent key")
	}

	// Test with existing key
	os.Setenv("TEST_HAS", "value")
	defer os.Unsetenv("TEST_HAS")

	if !parser.Has("TEST_HAS") {
		t.Error("Expected true for existing key")
	}
}

func TestParser_IsSet(t *testing.T) {
	parser := NewParser()

	// Test with non-existent key
	if parser.IsSet("NONEXISTENT_KEY") {
		t.Error("Expected false for non-existent key")
	}

	// Test with empty key
	os.Setenv("TEST_EMPTY", "")
	defer os.Unsetenv("TEST_EMPTY")

	if parser.IsSet("TEST_EMPTY") {
		t.Error("Expected false for empty key")
	}

	// Test with whitespace-only key
	os.Setenv("TEST_WS", "   ")
	defer os.Unsetenv("TEST_WS")

	if parser.IsSet("TEST_WS") {
		t.Error("Expected false for whitespace-only key")
	}

	// Test with valid key
	os.Setenv("TEST_VALID", "value")
	defer os.Unsetenv("TEST_VALID")

	if !parser.IsSet("TEST_VALID") {
		t.Error("Expected true for valid key")
	}
}
