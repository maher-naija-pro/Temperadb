package config

import (
	"os"
	"strconv"
	"time"
)

// Helper functions for environment variable handling
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvWithDefault(key, defaultValue string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.Atoi(defaultValue); err == nil {
		return intValue
	}
	return 0
}

func getInt64EnvWithDefault(key, defaultValue string) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
		return intValue
	}
	return 0
}

func getBoolEnvWithDefault(key, defaultValue string) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	if boolValue, err := strconv.ParseBool(defaultValue); err == nil {
		return boolValue
	}
	return false
}

func getDurationEnvWithDefault(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value + "s"); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue + "s"); err == nil {
		return duration
	}
	return 30 * time.Second
}
