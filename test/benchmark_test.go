package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// BenchmarkWriteEndpoint provides performance benchmarks
func BenchmarkWriteEndpoint(b *testing.B) {
	// Setup
	os.Setenv("DATA_FILE", "benchmark_data.tsv")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	// Test data
	testData := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, server.URL+"/write", strings.NewReader(testData))
		if err != nil {
			b.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}

// BenchmarkBulkWriteEndpoint benchmarks bulk write operations
func BenchmarkBulkWriteEndpoint(b *testing.B) {
	// Setup
	os.Setenv("DATA_FILE", "benchmark_bulk_data.tsv")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	// Generate bulk test data
	helper := NewTestHelper()
	testData := helper.GenerateBulkLineProtocolData(100, 1434055562000000000)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, server.URL+"/write", strings.NewReader(testData))
		if err != nil {
			b.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
	}
}

// BenchmarkHTTPRequestCreation benchmarks HTTP request creation
func BenchmarkHTTPRequestCreation(b *testing.B) {
	helper := NewTestHelper()
	testData := "cpu,host=server01 value=0.64 1434055562000000000"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := helper.CreateTestRequest(http.MethodPost, "http://localhost:8080/write", testData, nil)
		if req == nil {
			b.Fatal("Failed to create request")
		}
	}
}
