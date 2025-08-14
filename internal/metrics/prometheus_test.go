package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRegistry(t *testing.T) {
	// Test that GetRegistry returns the same registry instance
	registry1 := GetRegistry()
	registry2 := GetRegistry()

	assert.NotNil(t, registry1)
	assert.Equal(t, registry1, registry2)
	assert.Equal(t, Registry, registry1)
}

func TestReset(t *testing.T) {
	// Initialize metrics first
	Init()

	// Verify metrics are registered
	registry := GetRegistry()
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Should have metrics registered after Init")

	// Reset metrics
	Reset()

	// Verify metrics are unregistered
	metrics, err = registry.Gather()
	require.NoError(t, err)
	assert.Equal(t, 0, len(metrics), "Should have no metrics after Reset")
}

func TestInit(t *testing.T) {
	// Reset first to ensure clean state
	Reset()

	// Test successful initialization
	assert.NotPanics(t, func() {
		Init()
	})

	// Verify metrics are registered
	registry := GetRegistry()
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Should have metrics registered after Init")
}

func TestInitAfterReset(t *testing.T) {
	// Initialize, reset, then initialize again
	Init()
	Reset()

	// Should not panic on second Init
	assert.NotPanics(t, func() {
		Init()
	})

	// Verify metrics are registered again
	registry := GetRegistry()
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Should have metrics registered after second Init")
}

func TestMultipleInitCalls(t *testing.T) {
	// Reset first
	Reset()

	// Multiple Init calls should not panic
	assert.NotPanics(t, func() {
		Init()
		Init()
		Init()
	})

	// Verify metrics are registered
	registry := GetRegistry()
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Should have metrics registered after multiple Init calls")
}

func TestRegistryConsistency(t *testing.T) {
	// Test that the global Registry is consistent
	originalRegistry := Registry

	// Initialize should not change the registry instance
	Init()
	assert.Equal(t, originalRegistry, Registry)

	// Reset should not change the registry instance
	Reset()
	assert.Equal(t, originalRegistry, Registry)

	// GetRegistry should return the same instance
	assert.Equal(t, originalRegistry, GetRegistry())
}

func TestResetWithoutInit(t *testing.T) {
	// Reset without Init should not panic
	assert.NotPanics(t, func() {
		Reset()
	})

	// Registry should still be valid
	registry := GetRegistry()
	assert.NotNil(t, registry)
}

func TestRegistryType(t *testing.T) {
	// Verify Registry is of correct type
	assert.IsType(t, &prometheus.Registry{}, Registry)

	// Verify GetRegistry returns correct type
	registry := GetRegistry()
	assert.IsType(t, &prometheus.Registry{}, registry)
}

func TestMetricsRegistrationFlow(t *testing.T) {
	// Test complete flow: Reset -> Init -> Reset
	Reset()

	// Initial state should have no metrics
	registry := GetRegistry()
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Equal(t, 0, len(metrics), "Initial state should have no metrics")

	// Init should register metrics
	Init()
	metrics, err = registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Init should register metrics")

	// Reset should unregister all metrics
	Reset()
	metrics, err = registry.Gather()
	require.NoError(t, err)
	assert.Equal(t, 0, len(metrics), "Reset should unregister all metrics")
}

func TestConcurrentAccess(t *testing.T) {
	// Test that concurrent access to Registry doesn't cause issues
	done := make(chan bool)

	// Start multiple goroutines accessing Registry
	for i := 0; i < 10; i++ {
		go func() {
			registry := GetRegistry()
			assert.NotNil(t, registry)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRegistryMetrics(t *testing.T) {
	// Test that the registry can actually collect metrics
	Init()
	defer Reset()

	registry := GetRegistry()

	// Test that we can gather metrics
	metrics, err := registry.Gather()
	require.NoError(t, err)
	assert.Greater(t, len(metrics), 0, "Registry should contain metrics")

	// Test that metrics have names and help text
	for _, metric := range metrics {
		assert.NotEmpty(t, metric.GetName(), "Metric should have a name")
		// Note: Some metrics might not have help text, so we don't assert on that
	}
}
