package metrics

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testPort = ":8080"
	altPort  = ":9090"
)

func TestNewMetricsServer(t *testing.T) {
	addr := testPort
	server := NewMetricsServer(addr)

	assert.NotNil(t, server)
	assert.Equal(t, addr, server.addr)
	assert.NotNil(t, server.registry)
}

func TestMetricsServerEndpoints(t *testing.T) {
	addr := testPort
	server := NewMetricsServer(addr)

	// Test metrics endpoint
	metricsEndpoint := server.GetMetricsEndpoint()
	healthEndpoint := server.GetHealthEndpoint()
	readyEndpoint := server.GetReadyEndpoint()

	assert.Equal(t, "http://"+addr+"/metrics", metricsEndpoint)
	assert.Equal(t, "http://"+addr+"/health", healthEndpoint)
	assert.Equal(t, "http://"+addr+"/ready", readyEndpoint)
}

func TestMetricsServerStartStop(t *testing.T) {
	addr := altPort
	server := NewMetricsServer(addr)

	// Start server in goroutine
	go func() {
		err := server.Start()
		assert.NoError(t, err)
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://localhost" + addr + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Note: Server will continue running, no Stop method available
	// In a real test, you might want to use a context with timeout
}
