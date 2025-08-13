package server

import (
	"testing"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/metrics"
	"timeseriesdb/test/helpers"
)

func BenchmarkNewServer(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration once
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		server, err := NewServer(cfg)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}
		server.Close()
	}
}

func BenchmarkServerStartStop(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration once
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		server, err := NewServer(cfg)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}

		// Start server in background
		go func() {
			server.Start()
		}()

		// Wait a bit for server to start
		time.Sleep(10 * time.Millisecond)

		// Close server
		server.Close()
	}
}

func BenchmarkServerClose(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration once
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		server, err := NewServer(cfg)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}
		server.Close()
	}
}

func BenchmarkServerConcurrentAccess(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration once
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server, err := NewServer(cfg)
	if err != nil {
		b.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Access server properties concurrently
			_ = server.config
			_ = server.storage
			_ = server.httpServer
		}
	})
}

func BenchmarkServerMemoryAllocation(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration once
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		server, err := NewServer(cfg)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}
		server.Close()
	}
}

func BenchmarkServerWithLargeConfig(b *testing.B) {
	defer metrics.Reset()

	// Create test configuration with larger values
	cfg := helpers.Config.CreateTestConfig(&testing.T{})
	cfg.Server = config.ServerConfig{
		Port:         "8080",
		ReadTimeout:  300 * time.Second, // 5 minutes
		WriteTimeout: 300 * time.Second, // 5 minutes
		IdleTimeout:  600 * time.Second, // 10 minutes
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		server, err := NewServer(cfg)
		if err != nil {
			b.Fatalf("Failed to create server: %v", err)
		}
		server.Close()
	}
}
