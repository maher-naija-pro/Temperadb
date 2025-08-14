package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricsServer(t *testing.T) {
	t.Run("creates new metrics server", func(t *testing.T) {
		addr := ":8080"
		server := NewMetricsServer(addr)

		assert.NotNil(t, server)
		assert.Equal(t, addr, server.addr)
		assert.NotNil(t, server.registry)
	})
}

func TestMetricsServerEndpoints(t *testing.T) {
	t.Run("returns correct endpoint URLs", func(t *testing.T) {
		addr := ":8080"
		server := NewMetricsServer(addr)

		metricsEndpoint := server.GetMetricsEndpoint()
		healthEndpoint := server.GetHealthEndpoint()
		readyEndpoint := server.GetReadyEndpoint()

		assert.Equal(t, "http://:8080/metrics", metricsEndpoint)
		assert.Equal(t, "http://:8080/health", healthEndpoint)
		assert.Equal(t, "http://:8080/ready", readyEndpoint)
	})
}

func TestMetricsServerConfiguration(t *testing.T) {
	t.Run("server has correct configuration", func(t *testing.T) {
		addr := ":9090"
		server := NewMetricsServer(addr)

		assert.Equal(t, addr, server.addr)
		assert.NotNil(t, server.registry)
	})
}

func TestMetricsServerRegistryAccess(t *testing.T) {
	t.Run("server registry is accessible", func(t *testing.T) {
		server := NewMetricsServer(":8080")

		// The registry should be the global Registry from prometheus.go
		assert.Equal(t, Registry, server.registry)
	})
}
