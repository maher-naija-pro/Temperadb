package logger

import (
	"io"
	"os"
	"timeseriesdb/config"

	"github.com/sirupsen/logrus"
)

// MockLogger interface for testing fatal functions
type MockLogger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// CustomLogger wraps logrus.Logger to allow testing fatal functions
type CustomLogger struct {
	*logrus.Logger
	testMode bool
}

// Fatal overrides logrus.Logger.Fatal to allow testing
func (l *CustomLogger) Fatal(args ...interface{}) {
	if l.testMode {
		l.Error(args...)
		return
	}
	l.Logger.Fatal(args...)
}

// Fatalf overrides logrus.Logger.Fatalf to allow testing
func (l *CustomLogger) Fatalf(format string, args ...interface{}) {
	if l.testMode {
		l.Errorf(format, args...)
		return
	}
	l.Logger.Fatalf(format, args...)
}

var (
	// Log is the global logger instance
	Log *CustomLogger
	// testMode indicates if the logger is running in test mode
	testMode bool
	// testWriter is used to capture output during testing
	testWriter io.Writer
	// mockLogger is used for testing fatal functions
	mockLogger MockLogger
)

// SetTestMode enables or disables test mode for the logger
func SetTestMode(enabled bool) {
	testMode = enabled
}

// SetTestWriter sets a custom writer for testing
func SetTestWriter(writer io.Writer) {
	testWriter = writer
}

// SetMockLogger sets a mock logger for testing fatal functions
func SetMockLogger(logger MockLogger) {
	mockLogger = logger
}

// Init initializes the global logger with configuration
func Init() {
	baseLogger := logrus.New()

	// Create custom logger
	Log = &CustomLogger{
		Logger:   baseLogger,
		testMode: testMode,
	}

	// Set output based on mode
	if testMode && testWriter != nil {
		Log.SetOutput(testWriter)
	} else {
		// Set output to stdout
		Log.SetOutput(os.Stdout)
	}

	// Set formatter with full timestamp
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level (can be configured via environment variable)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err == nil {
			Log.SetLevel(level)
		}
	}
}

// InitWithConfig initializes the global logger with the provided configuration
func InitWithConfig(cfg config.LoggingConfig) {
	baseLogger := logrus.New()

	// Create custom logger
	Log = &CustomLogger{
		Logger:   baseLogger,
		testMode: testMode,
	}

	// Set output based on mode
	if testMode && testWriter != nil {
		Log.SetOutput(testWriter)
	} else {
		// Set output based on configuration
		switch cfg.Output {
		case "stdout":
			Log.SetOutput(os.Stdout)
		case "stderr":
			Log.SetOutput(os.Stderr)
		default:
			Log.SetOutput(os.Stdout)
		}
	}

	// Set formatter based on configuration
	Log.SetFormatter(cfg.GetLogFormat())

	// Set log level from configuration
	Log.SetLevel(cfg.GetLogLevel())
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal logs a fatal message and exits (unless in test mode)
func Fatal(args ...interface{}) {
	if testMode && mockLogger != nil {
		mockLogger.Fatal(args...)
		return
	}
	Log.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits (unless in test mode)
func Fatalf(format string, args ...interface{}) {
	if testMode && mockLogger != nil {
		mockLogger.Fatalf(format, args...)
		return
	}
	Log.Fatalf(format, args...)
}

// WithField returns a logger with a field
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields returns a logger with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}

// WithError returns a logger with an error
func WithError(err error) *logrus.Entry {
	return Log.WithError(err)
}
