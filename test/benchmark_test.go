package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
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

	b.Run("Single Point", func(b *testing.B) {
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
	})

	b.Run("Multiple Points", func(b *testing.B) {
		// Generate multiple points
		var lines []string
		for i := 0; i < 10; i++ {
			line := fmt.Sprintf("cpu,host=server%02d,region=us-west value=%d.0 1434055562000000000", i, i)
			lines = append(lines, line)
		}
		multipleData := strings.Join(lines, "\n")

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req, err := http.NewRequest(http.MethodPost, server.URL+"/write", strings.NewReader(multipleData))
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
	})

	b.Run("Large Dataset", func(b *testing.B) {
		// Generate 100 points
		var lines []string
		for i := 0; i < 100; i++ {
			line := fmt.Sprintf("cpu,host=server%03d,region=us-west value=%d.0,load=%d.5 1434055562000000000", i, i, i%10)
			lines = append(lines, line)
		}
		largeData := strings.Join(lines, "\n")

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req, err := http.NewRequest(http.MethodPost, server.URL+"/write", strings.NewReader(largeData))
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
	})
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

	b.Run("100 Points", func(b *testing.B) {
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
	})

	b.Run("1000 Points", func(b *testing.B) {
		testData := helper.GenerateBulkLineProtocolData(1000, 1434055562000000000)

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
	})

	b.Run("10000 Points", func(b *testing.B) {
		testData := helper.GenerateBulkLineProtocolData(10000, 1434055562000000000)

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
	})
}

// BenchmarkHTTPRequestCreation benchmarks HTTP request creation
func BenchmarkHTTPRequestCreation(b *testing.B) {
	helper := NewTestHelper()
	testData := "cpu,host=server01 value=0.64 1434055562000000000"

	b.Run("Basic Request", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, "http://localhost:8080/write", testData, nil)
			if req == nil {
				b.Fatal("Failed to create request")
			}
		}
	})

	b.Run("Request with Headers", func(b *testing.B) {
		headers := map[string]string{
			"Content-Type": "text/plain",
			"User-Agent":   "BenchmarkTest/1.0",
			"Accept":       "*/*",
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, "http://localhost:8080/write", testData, headers)
			if req == nil {
				b.Fatal("Failed to create request")
			}
		}
	})

	b.Run("Request with Large Body", func(b *testing.B) {
		largeData := strings.Repeat("cpu,host=server01 value=0.64 1434055562000000000\n", 1000)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, "http://localhost:8080/write", largeData, nil)
			if req == nil {
				b.Fatal("Failed to create request")
			}
		}
	})
}

// BenchmarkDataGeneration benchmarks data generation functions
func BenchmarkDataGeneration(b *testing.B) {
	helper := NewTestHelper()
	baseTime := time.Now().UnixNano()

	b.Run("Single Line Protocol", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateLineProtocolData(
				"cpu",
				map[string]string{"host": "server01", "region": "us-west"},
				map[string]float64{"value": 0.64, "load": 0.8},
				baseTime,
			)
		}
	})

	b.Run("Bulk Data 100", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateBulkLineProtocolData(100, baseTime)
		}
	})

	b.Run("Bulk Data 1000", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateBulkLineProtocolData(1000, baseTime)
		}
	})

	b.Run("Random Data 100", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomLineProtocolData(100, baseTime)
		}
	})

	b.Run("Random Data 1000", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomLineProtocolData(1000, baseTime)
		}
	})

	b.Run("Stress Test Data 100", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateStressTestData(100, baseTime)
		}
	})

	b.Run("Unicode Test Data 100", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateUnicodeTestData(100, baseTime)
		}
	})
}

// BenchmarkConcurrentRequests benchmarks concurrent HTTP requests
func BenchmarkConcurrentRequests(b *testing.B) {
	// Setup
	os.Setenv("DATA_FILE", "benchmark_concurrent_data.tsv")

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	testData := "cpu,host=server01 value=0.64 1434055562000000000"

	b.Run("Concurrent 10", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
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
		})
	})

	b.Run("Concurrent 100", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
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
		})
	})
}

// BenchmarkMemoryUsage benchmarks memory usage patterns
func BenchmarkMemoryUsage(b *testing.B) {
	helper := NewTestHelper()
	baseTime := time.Now().UnixNano()

	b.Run("Memory Allocation - Small Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			data := helper.GenerateLineProtocolData(
				"cpu",
				map[string]string{"host": "server01"},
				map[string]float64{"value": 0.64},
				baseTime,
			)
			_ = data
		}
	})

	b.Run("Memory Allocation - Large Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			data := helper.GenerateBulkLineProtocolData(1000, baseTime)
			_ = data
		}
	})

	b.Run("Memory Allocation - Stress Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			data := helper.GenerateStressTestData(100, baseTime)
			_ = data
		}
	})
}

// BenchmarkStringOperations benchmarks string manipulation operations
func BenchmarkStringOperations(b *testing.B) {
	helper := NewTestHelper()

	b.Run("Random String Generation", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomString(100)
		}
	})

	b.Run("Random Tags Generation", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomTags(10)
		}
	})

	b.Run("Random Fields Generation", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomFields(10)
		}
	})
}

// BenchmarkHTTPClientOperations benchmarks HTTP client operations
func BenchmarkHTTPClientOperations(b *testing.B) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	}))
	defer server.Close()

	helper := NewTestHelper()
	testData := "cpu,host=server01 value=0.64 1434055562000000000"

	b.Run("Request Creation and Execution", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, server.URL+"/write", testData, nil)
			resp, err := helper.ExecuteRequest(req)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("Request with Retry Logic", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, server.URL+"/write", testData, nil)
			resp, err := helper.ExecuteRequestWithRetry(req, 3, 10*time.Millisecond)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})

	b.Run("Request with Custom Timeout", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			req := helper.CreateTestRequest(http.MethodPost, server.URL+"/write", testData, nil)
			resp, err := helper.ExecuteRequestWithTimeout(req, 5*time.Second)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}

// BenchmarkHelperFunctions benchmarks helper utility functions
func BenchmarkHelperFunctions(b *testing.B) {
	helper := NewTestHelper()
	baseTime := time.Now().UnixNano()

	b.Run("Generate Line Protocol Data", func(b *testing.B) {
		tags := map[string]string{
			"host":   "server01",
			"region": "us-west",
			"env":    "production",
		}
		fields := map[string]float64{
			"value": 0.64,
			"load":  0.8,
			"count": 42.0,
		}

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateLineProtocolData("cpu", tags, fields, baseTime)
		}
	})

	b.Run("Generate Bulk Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateBulkLineProtocolData(100, baseTime)
		}
	})

	b.Run("Generate Random Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateRandomLineProtocolData(100, baseTime)
		}
	})

	b.Run("Generate Stress Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateStressTestData(100, baseTime)
		}
	})

	b.Run("Generate Unicode Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateUnicodeTestData(100, baseTime)
		}
	})

	b.Run("Generate Malformed Data", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = helper.GenerateMalformedData(100)
		}
	})
}
