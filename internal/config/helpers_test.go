package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHelperFunctions(t *testing.T) {
	t.Run("getEnvWithDefault", func(t *testing.T) {
		// Test with environment variable set
		os.Setenv("TEST_KEY", "test_value")
		defer os.Unsetenv("TEST_KEY")

		result := getEnvWithDefault("TEST_KEY", "default_value")
		assert.Equal(t, "test_value", result)

		// Test with environment variable not set
		result = getEnvWithDefault("NONEXISTENT_KEY", "default_value")
		assert.Equal(t, "default_value", result)
	})

	t.Run("getIntEnvWithDefault", func(t *testing.T) {
		// Test with valid integer
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		result := getIntEnvWithDefault("TEST_INT", "100")
		assert.Equal(t, 42, result)

		// Test with invalid integer
		os.Setenv("TEST_INVALID", "invalid")
		defer os.Unsetenv("TEST_INVALID")

		result = getIntEnvWithDefault("TEST_INVALID", "100")
		assert.Equal(t, 100, result)

		// Test with environment variable not set
		result = getIntEnvWithDefault("NONEXISTENT_INT", "200")
		assert.Equal(t, 200, result)

		// Test with invalid default value (should fallback to 0)
		result = getIntEnvWithDefault("NONEXISTENT_INT", "invalid_default")
		assert.Equal(t, 0, result)
	})

	t.Run("getInt64EnvWithDefault", func(t *testing.T) {
		// Test with valid int64
		os.Setenv("TEST_INT64", "9223372036854775807")
		defer os.Unsetenv("TEST_INT64")

		result := getInt64EnvWithDefault("TEST_INT64", "100")
		assert.Equal(t, int64(9223372036854775807), result)

		// Test with invalid int64
		os.Setenv("TEST_INVALID64", "invalid")
		defer os.Unsetenv("TEST_INVALID64")

		result = getInt64EnvWithDefault("TEST_INVALID64", "100")
		assert.Equal(t, int64(100), result)

		// Test with environment variable not set
		result = getInt64EnvWithDefault("NONEXISTENT_INT64", "200")
		assert.Equal(t, int64(200), result)

		// Test with invalid default value (should fallback to 0)
		result = getInt64EnvWithDefault("NONEXISTENT_INT64", "invalid_default")
		assert.Equal(t, int64(0), result)
	})

	t.Run("getBoolEnvWithDefault", func(t *testing.T) {
		// Test with valid boolean
		os.Setenv("TEST_BOOL", "true")
		defer os.Unsetenv("TEST_BOOL")

		result := getBoolEnvWithDefault("TEST_BOOL", "false")
		assert.True(t, result)

		// Test with invalid boolean
		os.Setenv("TEST_INVALID_BOOL", "invalid")
		defer os.Unsetenv("TEST_INVALID_BOOL")

		result = getBoolEnvWithDefault("TEST_INVALID_BOOL", "true")
		assert.True(t, result)

		// Test with environment variable not set
		result = getBoolEnvWithDefault("NONEXISTENT_BOOL", "false")
		assert.False(t, result)

		// Test with invalid default value (should fallback to false)
		result = getBoolEnvWithDefault("NONEXISTENT_BOOL", "invalid_default")
		assert.False(t, result)
	})

	t.Run("getDurationEnvWithDefault", func(t *testing.T) {
		// Test with valid duration
		os.Setenv("TEST_DURATION", "60")
		defer os.Unsetenv("TEST_DURATION")

		result := getDurationEnvWithDefault("TEST_DURATION", "30")
		assert.Equal(t, 60*time.Second, result)

		// Test with invalid duration
		os.Setenv("TEST_INVALID_DURATION", "invalid")
		defer os.Unsetenv("TEST_INVALID_DURATION")

		result = getDurationEnvWithDefault("TEST_INVALID_DURATION", "45")
		assert.Equal(t, 45*time.Second, result)

		// Test with environment variable not set
		result = getDurationEnvWithDefault("NONEXISTENT_DURATION", "30")
		assert.Equal(t, 30*time.Second, result)

		// Test with invalid default value (should fallback to 30s)
		result = getDurationEnvWithDefault("NONEXISTENT_DURATION", "invalid_default")
		assert.Equal(t, 30*time.Second, result)
	})
}
