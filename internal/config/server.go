package config

import (
	"time"

	"timeseriesdb/internal/envvars"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// NewServerConfig creates a new ServerConfig with default values
func NewServerConfig() ServerConfig {
	parser := envvars.NewParser()

	return ServerConfig{
		Port:            parser.String(envvars.Port, envvars.DefaultPort),
		ReadTimeout:     parser.Duration(envvars.ReadTimeout, envvars.DefaultReadTimeout),
		WriteTimeout:    parser.Duration(envvars.WriteTimeout, envvars.DefaultWriteTimeout),
		IdleTimeout:     parser.Duration(envvars.IdleTimeout, envvars.DefaultIdleTimeout),
		ShutdownTimeout: parser.Duration(envvars.ShutdownTimeout, envvars.DefaultShutdownTimeout),
	}
}
