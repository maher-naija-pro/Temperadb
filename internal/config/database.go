package config

import (
	"time"

	"timeseriesdb/internal/envvars"
)

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	MaxConnections int
	ConnectionTTL  time.Duration
	QueryTimeout   time.Duration
}

// NewDatabaseConfig creates a new DatabaseConfig with default values
func NewDatabaseConfig() DatabaseConfig {
	parser := envvars.NewParser()

	return DatabaseConfig{
		MaxConnections: parser.Int(envvars.MaxConnections, envvars.DefaultMaxConnections),
		ConnectionTTL:  parser.Duration(envvars.ConnectionTTL, envvars.DefaultConnectionTTL),
		QueryTimeout:   parser.Duration(envvars.QueryTimeout, envvars.DefaultQueryTimeout),
	}
}
