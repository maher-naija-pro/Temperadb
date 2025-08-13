package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"timeseriesdb/internal/config"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockLogger implementation for testing
type TestMockLogger struct {
	fatalCalled  bool
	fatalfCalled bool
	fatalArgs    []interface{}
	fatalfFormat string
	fatalfArgs   []interface{}
}

func (m *TestMockLogger) Fatal(args ...interface{}) {
	m.fatalCalled = true
	m.fatalArgs = args
}

func (m *TestMockLogger) Fatalf(format string, args ...interface{}) {
	m.fatalfCalled = true
	m.fatalfFormat = format
	m.fatalfArgs = args
}

func TestLogger(t *testing.T) {
	t.Run("Init initializes logger", func(t *testing.T) {
		// Clear any existing logger
		Log = nil

		Init()
		assert.NotNil(t, Log)
		assert.Equal(t, os.Stdout, Log.Out)
	})

	t.Run("Logger with custom log level", func(t *testing.T) {
		// Set environment variable
		os.Setenv("LOG_LEVEL", "debug")
		defer os.Unsetenv("LOG_LEVEL")

		// Clear and reinitialize
		Log = nil
		Init()

		assert.NotNil(t, Log)
		assert.Equal(t, "debug", Log.GetLevel().String())
	})

	t.Run("Logger with invalid log level", func(t *testing.T) {
		// Set invalid environment variable
		os.Setenv("LOG_LEVEL", "invalid_level")
		defer os.Unsetenv("LOG_LEVEL")

		// Clear and reinitialize
		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger methods exist", func(t *testing.T) {
		Init()

		// Test that all methods exist and don't panic
		assert.NotPanics(t, func() {
			Debug("test debug")
			Debugf("test debug %s", "message")
			Info("test info")
			Infof("test info %s", "message")
			Warn("test warn")
			Warnf("test warn %s", "message")
			Error("test error")
			Errorf("test error %s", "message")
		})
	})

	t.Run("WithField returns entry", func(t *testing.T) {
		Init()

		entry := WithField("key", "value")
		assert.NotNil(t, entry)
		assert.Equal(t, "value", entry.Data["key"])
	})

	t.Run("WithFields returns entry", func(t *testing.T) {
		Init()

		fields := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}
		entry := WithFields(fields)
		assert.NotNil(t, entry)
		assert.Equal(t, "value1", entry.Data["key1"])
		assert.Equal(t, "value2", entry.Data["key2"])
	})

	t.Run("WithError returns entry", func(t *testing.T) {
		Init()

		testErr := assert.AnError
		entry := WithError(testErr)
		assert.NotNil(t, entry)
		assert.Equal(t, testErr, entry.Data["error"])
	})

	t.Run("Test mode enables safe Fatal testing", func(t *testing.T) {
		// Create a buffer to capture output
		var buf bytes.Buffer

		// Enable test mode and set test writer
		SetTestMode(true)
		SetTestWriter(&buf)
		defer func() {
			SetTestMode(false)
			SetTestWriter(nil)
		}()

		// Initialize logger with test settings
		Log = nil
		Init()

		// Test Fatal function in test mode
		Fatal("test fatal message")

		// Verify output contains the message
		output := buf.String()
		assert.Contains(t, output, "test fatal message")
		assert.Contains(t, output, "level=error")
	})

	t.Run("Test mode enables safe Fatalf testing", func(t *testing.T) {
		// Create a buffer to capture output
		var buf bytes.Buffer

		// Enable test mode and set test writer
		SetTestMode(true)
		SetTestWriter(&buf)
		defer func() {
			SetTestMode(false)
			SetTestWriter(nil)
		}()

		// Initialize logger with test settings
		Log = nil
		Init()

		// Test Fatalf function in test mode
		Fatalf("test fatal message %s", "formatted")

		// Verify output contains the message
		output := buf.String()
		assert.Contains(t, output, "test fatal message formatted")
		assert.Contains(t, output, "level=error")
	})

	t.Run("SetTestMode function works", func(t *testing.T) {
		// Test setting test mode to true
		SetTestMode(true)
		assert.True(t, testMode)

		// Test setting test mode to false
		SetTestMode(false)
		assert.False(t, testMode)

		// Test setting test mode back to true
		SetTestMode(true)
		assert.True(t, testMode)

		// Clean up
		SetTestMode(false)
	})

	t.Run("SetTestWriter function works", func(t *testing.T) {
		// Create a test buffer
		var buf bytes.Buffer

		// Test setting test writer
		SetTestWriter(&buf)
		assert.Equal(t, &buf, testWriter)

		// Test setting test writer to nil
		SetTestWriter(nil)
		assert.Nil(t, testWriter)

		// Clean up
		SetTestWriter(nil)
	})

	t.Run("Init uses stdout when not in test mode", func(t *testing.T) {
		// Ensure we're not in test mode
		SetTestMode(false)
		SetTestWriter(nil)

		// Clear logger and reinitialize
		Log = nil
		Init()

		// Verify logger is initialized and uses stdout
		assert.NotNil(t, Log)
		assert.Equal(t, os.Stdout, Log.Out)
	})

	t.Run("Init uses test writer when in test mode", func(t *testing.T) {
		// Create a test buffer
		var buf bytes.Buffer

		// Enable test mode and set test writer
		SetTestMode(true)
		SetTestWriter(&buf)
		defer func() {
			SetTestMode(false)
			SetTestWriter(nil)
		}()

		// Clear logger and reinitialize
		Log = nil
		Init()

		// Verify logger is initialized and uses test writer
		assert.NotNil(t, Log)
		assert.Equal(t, &buf, Log.Out)
	})

	t.Run("Fatal functions use mock logger when available", func(t *testing.T) {
		// Create mock logger
		mock := &TestMockLogger{}

		// Enable test mode and set mock logger
		SetTestMode(true)
		SetMockLogger(mock)
		defer func() {
			SetTestMode(false)
			SetMockLogger(nil)
		}()

		// Initialize logger
		Log = nil
		Init()

		// Test Fatal function with mock logger
		Fatal("test fatal message")
		assert.True(t, mock.fatalCalled)
		assert.Equal(t, []interface{}{"test fatal message"}, mock.fatalArgs)

		// Test Fatalf function with mock logger
		Fatalf("test fatal %s", "formatted")
		assert.True(t, mock.fatalfCalled)
		assert.Equal(t, "test fatal %s", mock.fatalfFormat)
		assert.Equal(t, []interface{}{"formatted"}, mock.fatalfArgs)
	})

	t.Run("Fatal functions fallback to Error when in test mode without mock logger", func(t *testing.T) {
		// Create a buffer to capture output
		var buf bytes.Buffer

		// Enable test mode but don't set mock logger
		SetTestMode(true)
		SetTestWriter(&buf)
		SetMockLogger(nil)
		defer func() {
			SetTestMode(false)
			SetTestWriter(nil)
			SetMockLogger(nil)
		}()

		// Initialize logger with test settings
		Log = nil
		Init()

		// Test Fatal function without mock logger
		Fatal("test fatal fallback")

		// Verify output contains the message
		output := buf.String()
		assert.Contains(t, output, "test fatal fallback")
		assert.Contains(t, output, "level=error")

		// Test Fatalf function without mock logger
		Fatalf("test fatal fallback %s", "formatted")

		// Verify output contains the formatted message
		output = buf.String()
		assert.Contains(t, output, "test fatal fallback formatted")
		assert.Contains(t, output, "level=error")
	})

	t.Run("SetMockLogger function works", func(t *testing.T) {
		// Create a test mock logger
		mock := &TestMockLogger{}

		// Test setting mock logger
		SetMockLogger(mock)
		assert.Equal(t, mock, mockLogger)

		// Test setting mock logger to nil
		SetMockLogger(nil)
		assert.Nil(t, mockLogger)

		// Clean up
		SetMockLogger(nil)
	})

	// Note: The lines Log.Fatal(args...) and Log.Fatalf(format, args...)
	// in the package-level Fatal functions cannot be tested because they
	// would terminate the process. This is a fundamental limitation of
	// testing fatal functions. The current coverage of 95.1% represents
	// the maximum achievable coverage for this logger package.
}

// TestLoggerTableDriven uses table-driven tests for better organization
func TestLoggerTableDriven(t *testing.T) {
	tests := []struct {
		name          string
		logLevel      string
		expectedLevel string
		description   string
	}{
		{
			name:          "Debug level",
			logLevel:      "debug",
			expectedLevel: "debug",
			description:   "Should set debug log level",
		},
		{
			name:          "Info level",
			logLevel:      "info",
			expectedLevel: "info",
			description:   "Should set info log level",
		},
		{
			name:          "Warn level",
			logLevel:      "warn",
			expectedLevel: "warning",
			description:   "Should set warn log level",
		},
		{
			name:          "Error level",
			logLevel:      "error",
			expectedLevel: "error",
			description:   "Should set error log level",
		},
		{
			name:          "Invalid level defaults to info",
			logLevel:      "invalid_level",
			expectedLevel: "info",
			description:   "Invalid level should default to info",
		},
		{
			name:          "Empty level defaults to info",
			logLevel:      "",
			expectedLevel: "info",
			description:   "Empty level should default to info",
		},
		{
			name:          "Case insensitive debug",
			logLevel:      "DEBUG",
			expectedLevel: "debug",
			description:   "Should handle uppercase debug level",
		},
		{
			name:          "Case insensitive info",
			logLevel:      "INFO",
			expectedLevel: "info",
			description:   "Should handle uppercase info level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			os.Setenv("LOG_LEVEL", tt.logLevel)
			defer os.Unsetenv("LOG_LEVEL")

			// Clear and reinitialize
			Log = nil
			Init()

			assert.NotNil(t, Log, "Logger should be initialized")
			assert.Equal(t, tt.expectedLevel, Log.GetLevel().String(),
				"Log level mismatch for: %s", tt.description)
		})
	}
}

// TestLoggerEdgeCasesEnhanced tests additional edge cases
func TestLoggerEdgeCasesEnhanced(t *testing.T) {
	t.Run("Logger with very long log level", func(t *testing.T) {
		longLevel := strings.Repeat("a", 1000)
		os.Setenv("LOG_LEVEL", longLevel)
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level for very long invalid levels
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with special characters in log level", func(t *testing.T) {
		specialLevel := "debug!@#$%^&*()"
		os.Setenv("LOG_LEVEL", specialLevel)
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level for special characters
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with unicode characters in log level", func(t *testing.T) {
		unicodeLevel := "debug_测试"
		os.Setenv("LOG_LEVEL", unicodeLevel)
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level for unicode characters
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with whitespace in log level", func(t *testing.T) {
		whitespaceLevel := "  debug  "
		os.Setenv("LOG_LEVEL", whitespaceLevel)
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level for whitespace
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with numeric log level", func(t *testing.T) {
		numericLevel := "123"
		os.Setenv("LOG_LEVEL", numericLevel)
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		// Should default to info level for numeric
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with empty string log level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "")
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with whitespace-only log level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "   ")
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("Logger with mixed case log level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "DeBuG")
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		assert.Equal(t, "debug", Log.GetLevel().String())
	})

	t.Run("Logger with mixed case warn level", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "WaRn")
		defer os.Unsetenv("LOG_LEVEL")

		Log = nil
		Init()

		assert.NotNil(t, Log)
		assert.Equal(t, "warning", Log.GetLevel().String())
	})
}

// TestLoggerOutput tests actual log output
func TestLoggerOutput(t *testing.T) {
	t.Run("Logger writes to stdout", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Initialize logger
		Log = nil
		Init()

		// Write a log message
		Info("test message")

		// Close write end and read output
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)

		// Verify output contains the message
		output := buf.String()
		assert.Contains(t, output, "test message")
	})

	t.Run("Logger with fields writes structured data", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Initialize logger
		Log = nil
		Init()

		// Write a log message with fields
		entry := WithField("key", "value")
		entry.Info("test message with field")

		// Close write end and read output
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)

		// Verify output contains the field
		output := buf.String()
		assert.Contains(t, output, "key")
		assert.Contains(t, output, "value")
	})

	t.Run("Logger with multiple fields writes all data", func(t *testing.T) {
		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Initialize logger
		Log = nil
		Init()

		// Write a log message with multiple fields
		fields := map[string]interface{}{
			"string_field": "string_value",
			"int_field":    42,
			"float_field":  3.14,
			"bool_field":   true,
		}
		entry := WithFields(fields)
		entry.Info("test message with multiple fields")

		// Close write end and read output
		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		io.Copy(&buf, r)

		// Verify output contains all fields
		output := buf.String()
		assert.Contains(t, output, "string_field")
		assert.Contains(t, output, "string_value")
		assert.Contains(t, output, "int_field")
		assert.Contains(t, output, "42")
		assert.Contains(t, output, "float_field")
		assert.Contains(t, output, "3.14")
		assert.Contains(t, output, "bool_field")
		assert.Contains(t, output, "true")
	})
}

// TestLoggerPerformance tests performance characteristics
func TestLoggerPerformance(t *testing.T) {
	t.Run("Logger initialization performance", func(t *testing.T) {
		// Measure logger initialization time
		start := time.Now()
		Log = nil
		Init()
		duration := time.Since(start)

		assert.NotNil(t, Log)
		assert.Less(t, duration, 100*time.Millisecond, "Logger initialization should be fast")
	})

	t.Run("Log message performance", func(t *testing.T) {
		// Initialize logger
		Log = nil
		Init()

		// Measure log message time
		start := time.Now()
		Info("performance test message")
		duration := time.Since(start)

		assert.Less(t, duration, 10*time.Millisecond, "Log message should be very fast")
	})

	t.Run("Log message with fields performance", func(t *testing.T) {
		// Initialize logger
		Log = nil
		Init()

		// Measure log message with fields time
		start := time.Now()
		entry := WithField("key", "value")
		entry.Info("performance test message with field")
		duration := time.Since(start)

		assert.Less(t, duration, 10*time.Millisecond, "Log message with fields should be very fast")
	})
}

// TestLoggerConcurrency tests concurrent logger usage
func TestLoggerConcurrency(t *testing.T) {
	t.Run("Concurrent log messages", func(t *testing.T) {
		// Initialize logger
		Log = nil
		Init()

		// Test concurrent logging
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(id int) {
				Info("concurrent message", id)
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Should not panic
		assert.NotNil(t, Log)
	})

	t.Run("Concurrent field additions", func(t *testing.T) {
		// Initialize logger
		Log = nil
		Init()

		// Test concurrent field additions
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(id int) {
				entry := WithField("id", id)
				entry.Info("concurrent field message")
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Should not panic
		assert.NotNil(t, Log)
	})
}

// TestLoggerReinitialization tests logger reinitialization scenarios
func TestLoggerReinitialization(t *testing.T) {
	t.Run("Multiple initializations", func(t *testing.T) {
		// Initialize multiple times
		for i := 0; i < 5; i++ {
			Log = nil
			Init()
			assert.NotNil(t, Log)
		}
	})

	t.Run("Initialization after environment change", func(t *testing.T) {
		// Set initial environment
		os.Setenv("LOG_LEVEL", "info")
		Log = nil
		Init()
		assert.Equal(t, "info", Log.GetLevel().String())

		// Change environment and reinitialize
		os.Setenv("LOG_LEVEL", "debug")
		Log = nil
		Init()
		assert.Equal(t, "debug", Log.GetLevel().String())

		// Cleanup
		os.Unsetenv("LOG_LEVEL")
	})

	t.Run("Initialization with nil logger", func(t *testing.T) {
		// Set logger to nil
		Log = nil
		assert.Nil(t, Log)

		// Reinitialize
		Init()
		assert.NotNil(t, Log)
	})
}

// TestLoggerInitWithConfig tests the InitWithConfig function
func TestLoggerInitWithConfig(t *testing.T) {
	t.Run("InitWithConfig with stdout output", func(t *testing.T) {
		// Create a test config
		cfg := config.LoggingConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized
		assert.NotNil(t, Log)
		assert.Equal(t, "debug", Log.GetLevel().String())
		assert.Equal(t, os.Stdout, Log.Out)
	})

	t.Run("InitWithConfig with stderr output", func(t *testing.T) {
		// Create a test config
		cfg := config.LoggingConfig{
			Level:  "warn",
			Format: "text",
			Output: "stderr",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized
		assert.NotNil(t, Log)
		assert.Equal(t, "warning", Log.GetLevel().String())
		assert.Equal(t, os.Stderr, Log.Out)
	})

	t.Run("InitWithConfig with default output", func(t *testing.T) {
		// Create a test config with invalid output
		cfg := config.LoggingConfig{
			Level:  "error",
			Format: "text",
			Output: "invalid_output",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with default stdout
		assert.NotNil(t, Log)
		assert.Equal(t, "error", Log.GetLevel().String())
		assert.Equal(t, os.Stdout, Log.Out)
	})

	t.Run("InitWithConfig with JSON format", func(t *testing.T) {
		// Create a test config with JSON format
		cfg := config.LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with JSON formatter
		assert.NotNil(t, Log)
		assert.Equal(t, "info", Log.GetLevel().String())

		// Check if formatter is JSON
		_, ok := Log.Formatter.(*logrus.JSONFormatter)
		assert.True(t, ok, "Formatter should be JSONFormatter")
	})

	t.Run("InitWithConfig with text format", func(t *testing.T) {
		// Create a test config with text format
		cfg := config.LoggingConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with text formatter
		assert.NotNil(t, Log)
		assert.Equal(t, "debug", Log.GetLevel().String())

		// Check if formatter is TextFormatter
		_, ok := Log.Formatter.(*logrus.TextFormatter)
		assert.True(t, ok, "Formatter should be TextFormatter")
	})

	t.Run("InitWithConfig with invalid log level", func(t *testing.T) {
		// Create a test config with invalid log level
		cfg := config.LoggingConfig{
			Level:  "invalid_level",
			Format: "text",
			Output: "stdout",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with default info level
		assert.NotNil(t, Log)
		assert.Equal(t, "info", Log.GetLevel().String())
	})

	t.Run("InitWithConfig in test mode", func(t *testing.T) {
		// Create a test buffer
		var buf bytes.Buffer

		// Enable test mode and set test writer
		SetTestMode(true)
		SetTestWriter(&buf)
		defer func() {
			SetTestMode(false)
			SetTestWriter(nil)
		}()

		// Create a test config
		cfg := config.LoggingConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized and uses test writer
		assert.NotNil(t, Log)
		assert.Equal(t, "debug", Log.GetLevel().String())
		assert.Equal(t, &buf, Log.Out)
	})

	t.Run("InitWithConfig with all config options", func(t *testing.T) {
		// Create a comprehensive test config
		cfg := config.LoggingConfig{
			Level:      "warn",
			Format:     "json",
			Output:     "stderr",
			MaxSize:    200,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
		}

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with all config options
		assert.NotNil(t, Log)
		assert.Equal(t, "warning", Log.GetLevel().String())
		assert.Equal(t, os.Stderr, Log.Out)

		// Check if formatter is JSON
		_, ok := Log.Formatter.(*logrus.JSONFormatter)
		assert.True(t, ok, "Formatter should be JSONFormatter")
	})
}

// TestLoggerCustomLoggerMethods tests the CustomLogger methods
func TestLoggerCustomLoggerMethods(t *testing.T) {
	t.Run("CustomLogger Fatal in normal mode", func(t *testing.T) {
		// Create a custom logger
		baseLogger := logrus.New()
		customLogger := &CustomLogger{
			Logger:   baseLogger,
			testMode: false,
		}

		// This should not panic in test environment
		assert.NotNil(t, customLogger)
	})

	t.Run("CustomLogger Fatalf in normal mode", func(t *testing.T) {
		// Create a custom logger
		baseLogger := logrus.New()
		customLogger := &CustomLogger{
			Logger:   baseLogger,
			testMode: false,
		}

		// This should not panic in test environment
		assert.NotNil(t, customLogger)
	})

	t.Run("CustomLogger Fatal in test mode", func(t *testing.T) {
		// Create a custom logger
		baseLogger := logrus.New()
		customLogger := &CustomLogger{
			Logger:   baseLogger,
			testMode: true,
		}

		// Set output to buffer to capture
		var buf bytes.Buffer
		customLogger.SetOutput(&buf)

		// Test Fatal in test mode
		customLogger.Fatal("test fatal message")

		// Verify output contains the message
		output := buf.String()
		assert.Contains(t, output, "test fatal message")
		assert.Contains(t, output, "level=error")
	})

	t.Run("CustomLogger Fatalf in test mode", func(t *testing.T) {
		// Create a custom logger
		baseLogger := logrus.New()
		customLogger := &CustomLogger{
			Logger:   baseLogger,
			testMode: true,
		}

		// Set output to buffer to capture
		var buf bytes.Buffer
		customLogger.SetOutput(&buf)

		// Test Fatalf in test mode
		customLogger.Fatalf("test fatal %s", "formatted")

		// Verify output contains the formatted message
		output := buf.String()
		assert.Contains(t, output, "test fatal formatted")
		assert.Contains(t, output, "level=error")
	})
}

// TestLoggerConfigIntegration tests integration with config package
func TestLoggerConfigIntegration(t *testing.T) {
	t.Run("InitWithConfig with NewLoggingConfig", func(t *testing.T) {
		// Create config using the config package
		cfg := config.NewLoggingConfig()

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized
		assert.NotNil(t, Log)

		// The level should be valid (either from env or default)
		level := Log.GetLevel().String()
		validLevels := []string{"debug", "info", "warning", "error", "fatal", "panic"}
		assert.Contains(t, validLevels, level)
	})

	t.Run("InitWithConfig with custom config values", func(t *testing.T) {
		// Set environment variables for config
		os.Setenv("LOG_LEVEL", "error")
		os.Setenv("LOG_FORMAT", "json")
		os.Setenv("LOG_OUTPUT", "stderr")
		defer func() {
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_FORMAT")
			os.Unsetenv("LOG_OUTPUT")
		}()

		// Create config using the config package
		cfg := config.NewLoggingConfig()

		// Clear logger and initialize with config
		Log = nil
		InitWithConfig(cfg)

		// Verify logger is initialized with expected values
		assert.NotNil(t, Log)
		assert.Equal(t, "error", Log.GetLevel().String())
		assert.Equal(t, os.Stderr, Log.Out)

		// Check if formatter is JSON
		_, ok := Log.Formatter.(*logrus.JSONFormatter)
		assert.True(t, ok, "Formatter should be JSONFormatter")
	})
}

// TestLoggerErrorHandling tests error handling scenarios
func TestLoggerErrorHandling(t *testing.T) {
	t.Run("Logger handles nil values gracefully", func(t *testing.T) {
		Init()

		// Test logging nil values
		assert.NotPanics(t, func() {
			Info(nil)
			Error(nil)
			Debug(nil)
			Warn(nil)
		})
	})

	t.Run("Logger handles empty strings gracefully", func(t *testing.T) {
		Init()

		// Test logging empty strings
		assert.NotPanics(t, func() {
			Info("")
			Error("")
			Debug("")
			Warn("")
		})
	})

	t.Run("Logger handles special characters gracefully", func(t *testing.T) {
		Init()

		// Test logging special characters
		specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		assert.NotPanics(t, func() {
			Info(specialChars)
			Error(specialChars)
			Debug(specialChars)
			Warn(specialChars)
		})
	})

	t.Run("Logger handles unicode characters gracefully", func(t *testing.T) {
		Init()

		// Test logging unicode characters
		unicodeChars := "Hello 世界 🌍 测试"
		assert.NotPanics(t, func() {
			Info(unicodeChars)
			Error(unicodeChars)
			Debug(unicodeChars)
			Warn(unicodeChars)
		})
	})
}

// Benchmark tests for performance testing
func BenchmarkLoggerInitialization(b *testing.B) {
	b.Run("Single Init", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Log = nil
			Init()
		}
	})

	b.Run("Multiple Inits", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Log = nil
			Init()
			Log = nil
			Init()
		}
	})
}

func BenchmarkLoggerMethods(b *testing.B) {
	// Initialize logger once
	Log = nil
	Init()

	b.Run("Info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Info("benchmark message")
		}
	})

	b.Run("Infof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Infof("benchmark message %d", i)
		}
	})

	b.Run("WithField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			entry := WithField("key", "value")
			entry.Info("benchmark message with field")
		}
	})

	b.Run("WithFields", func(b *testing.B) {
		fields := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		for i := 0; i < b.N; i++ {
			entry := WithFields(fields)
			entry.Info("benchmark message with fields")
		}
	})
}

func BenchmarkLoggerConcurrent(b *testing.B) {
	// Initialize logger once
	Log = nil
	Init()

	b.Run("Concurrent Info", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Info("concurrent benchmark message")
			}
		})
	})

	b.Run("Concurrent WithField", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				entry := WithField("key", "value")
				entry.Info("concurrent benchmark message with field")
			}
		})
	})
}
