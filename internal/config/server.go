package config

import "time"

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewServerConfig creates a new ServerConfig with default values
func NewServerConfig() ServerConfig {
	return ServerConfig{
		Port:         getEnvWithDefault("PORT", defaults["PORT"]),
		ReadTimeout:  getDurationEnvWithDefault("READ_TIMEOUT", defaults["READ_TIMEOUT"]),
		WriteTimeout: getDurationEnvWithDefault("WRITE_TIMEOUT", defaults["WRITE_TIMEOUT"]),
		IdleTimeout:  getDurationEnvWithDefault("IDLE_TIMEOUT", defaults["IDLE_TIMEOUT"]),
	}
}
