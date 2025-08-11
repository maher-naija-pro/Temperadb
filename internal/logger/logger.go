package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Log is the global logger instance
	Log *logrus.Logger
)

// Init initializes the global logger with configuration
func Init() {
	Log = logrus.New()

	// Set output to stdout
	Log.SetOutput(os.Stdout)

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

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
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
