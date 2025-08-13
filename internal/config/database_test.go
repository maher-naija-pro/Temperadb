package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig(t *testing.T) {
	t.Run("NewDatabaseConfig with defaults", func(t *testing.T) {
		cfg := NewDatabaseConfig()
		assert.Equal(t, 100, cfg.MaxConnections)
		assert.Equal(t, 5*time.Minute, cfg.ConnectionTTL)
		assert.Equal(t, 30*time.Second, cfg.QueryTimeout)
	})

	t.Run("NewDatabaseConfig with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("MAX_CONNECTIONS", "200")
		os.Setenv("CONNECTION_TTL", "600")
		os.Setenv("QUERY_TIMEOUT", "60")
		defer func() {
			os.Unsetenv("MAX_CONNECTIONS")
			os.Unsetenv("CONNECTION_TTL")
			os.Unsetenv("QUERY_TIMEOUT")
		}()

		cfg := NewDatabaseConfig()
		assert.Equal(t, 200, cfg.MaxConnections)
		assert.Equal(t, 600*time.Second, cfg.ConnectionTTL)
		assert.Equal(t, 60*time.Second, cfg.QueryTimeout)
	})

	t.Run("NewDatabaseConfig with invalid environment variables", func(t *testing.T) {
		// Set invalid environment variables
		os.Setenv("MAX_CONNECTIONS", "invalid")
		os.Setenv("CONNECTION_TTL", "invalid")
		os.Setenv("QUERY_TIMEOUT", "invalid")
		defer func() {
			os.Unsetenv("MAX_CONNECTIONS")
			os.Unsetenv("CONNECTION_TTL")
			os.Unsetenv("QUERY_TIMEOUT")
		}()

		cfg := NewDatabaseConfig()
		// Should use defaults for invalid values
		assert.Equal(t, 100, cfg.MaxConnections)
		assert.Equal(t, 5*time.Minute, cfg.ConnectionTTL)
		assert.Equal(t, 30*time.Second, cfg.QueryTimeout)
	})
}
