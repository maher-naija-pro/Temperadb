package envvars

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Parser provides methods to parse environment variables with proper type conversion
type Parser struct{}

// NewParser creates a new environment variable parser
func NewParser() *Parser {
	return &Parser{}
}

// String parses a string environment variable with a default value
func (p *Parser) String(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return defaultValue
}

// Int parses an integer environment variable with a default value
func (p *Parser) Int(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Int64 parses an int64 environment variable with a default value
func (p *Parser) Int64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Bool parses a boolean environment variable with a default value
func (p *Parser) Bool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.ToLower(strings.TrimSpace(value))

		// Handle common boolean representations
		switch value {
		case "true", "1", "yes", "on", "enabled":
			return true
		case "false", "0", "no", "off", "disabled":
			return false
		}

		// Try standard boolean parsing as fallback
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// Duration parses a duration environment variable with a default value
// Expects values in seconds (e.g., "30" for 30 seconds)
func (p *Parser) Duration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(strings.TrimSpace(value) + "s"); err == nil {
			return duration
		}
	}
	return defaultValue
}

// FileSize parses a file size environment variable with a default value
// Supports human-readable formats like "1GB", "100MB", etc.
func (p *Parser) FileSize(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if size, err := parseFileSize(strings.TrimSpace(value)); err == nil {
			return size
		}
	}
	return defaultValue
}

// parseFileSize converts human-readable file size strings to bytes
func parseFileSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	// Handle numeric values (assumed to be in bytes)
	if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
		return size, nil
	}

	// Handle human-readable formats
	var multiplier int64 = 1
	var numericPart string

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		numericPart = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		numericPart = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		numericPart = strings.TrimSuffix(sizeStr, "GB")
	} else if strings.HasSuffix(sizeStr, "TB") {
		multiplier = 1024 * 1024 * 1024 * 1024
		numericPart = strings.TrimSuffix(sizeStr, "TB")
	} else {
		return 0, strconv.ErrSyntax
	}

	if size, err := strconv.ParseInt(numericPart, 10, 64); err == nil {
		return size * multiplier, nil
	}

	return 0, strconv.ErrSyntax
}

// RequiredString returns an environment variable value or panics if not set
func (p *Parser) RequiredString(key string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	panic("required environment variable " + key + " is not set")
}

// RequiredInt returns an environment variable value as int or panics if not set
func (p *Parser) RequiredInt(key string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return intValue
		}
		panic("environment variable " + key + " must be a valid integer")
	}
	panic("required environment variable " + key + " is not set")
}

// Has checks if an environment variable is set
func (p *Parser) Has(key string) bool {
	return os.Getenv(key) != ""
}

// IsSet checks if an environment variable is set and not empty
func (p *Parser) IsSet(key string) bool {
	value := os.Getenv(key)
	return value != "" && strings.TrimSpace(value) != ""
}
