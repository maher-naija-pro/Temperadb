package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	cfg := NewServerConfig()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, 30*time.Second, cfg.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.WriteTimeout)
	assert.Equal(t, 120*time.Second, cfg.IdleTimeout)
	assert.Equal(t, 30*time.Second, cfg.ShutdownTimeout)
}

func TestServerConfig(t *testing.T) {
	t.Run("NewServerConfig with defaults", func(t *testing.T) {
		cfg := NewServerConfig()
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, 30*time.Second, cfg.ReadTimeout)
		assert.Equal(t, 30*time.Second, cfg.WriteTimeout)
		assert.Equal(t, 120*time.Second, cfg.IdleTimeout)
	})

	t.Run("NewServerConfig with environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("PORT", "9090")
		os.Setenv("READ_TIMEOUT", "60")
		os.Setenv("WRITE_TIMEOUT", "45")
		os.Setenv("IDLE_TIMEOUT", "180")
		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("READ_TIMEOUT")
			os.Unsetenv("WRITE_TIMEOUT")
			os.Unsetenv("IDLE_TIMEOUT")
		}()

		cfg := NewServerConfig()
		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, 60*time.Second, cfg.ReadTimeout)
		assert.Equal(t, 45*time.Second, cfg.WriteTimeout)
		assert.Equal(t, 180*time.Second, cfg.IdleTimeout)
	})

	t.Run("NewServerConfig with invalid environment variables", func(t *testing.T) {
		// Set invalid environment variables
		os.Setenv("READ_TIMEOUT", "invalid")
		os.Setenv("WRITE_TIMEOUT", "invalid")
		os.Setenv("IDLE_TIMEOUT", "invalid")
		defer func() {
			os.Unsetenv("READ_TIMEOUT")
			os.Unsetenv("WRITE_TIMEOUT")
			os.Unsetenv("IDLE_TIMEOUT")
		}()

		cfg := NewServerConfig()
		// Should use defaults for invalid values
		assert.Equal(t, 30*time.Second, cfg.ReadTimeout)
		assert.Equal(t, 30*time.Second, cfg.WriteTimeout)
		assert.Equal(t, 120*time.Second, cfg.IdleTimeout)
	})
}
