package config

import "time"

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	MaxConnections int
	ConnectionTTL  time.Duration
	QueryTimeout   time.Duration
}

// NewDatabaseConfig creates a new DatabaseConfig with default values
func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MaxConnections: getIntEnvWithDefault("MAX_CONNECTIONS", defaults["MAX_CONNECTIONS"]),
		ConnectionTTL:  getDurationEnvWithDefault("CONNECTION_TTL", defaults["CONNECTION_TTL"]),
		QueryTimeout:   getDurationEnvWithDefault("QUERY_TIMEOUT", defaults["QUERY_TIMEOUT"]),
	}
}
