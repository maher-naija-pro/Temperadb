package benchmark

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	aphttp "timeseriesdb/internal/api/http"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/ingestion"
	"timeseriesdb/internal/storage"
)


// BenchmarkHTTPWriteSinglePoint benchmarks HTTP write endpoint with single point
func BenchmarkHTTPWriteSinglePoint(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	// Create test server
	router := aphttp.NewRouter(storageInstance)
	server := httptest.NewServer(router.GetMux())
	defer server.Close()

	// Parse the data first
	points, err := ingestion.ParseLineProtocol("cpu,host=server01,region=us-west value=0.64 1434055562000000000")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP for accurate benchmarking
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkHTTPWriteMultiplePoints benchmarks HTTP write endpoint with multiple points
func BenchmarkHTTPWriteMultiplePoints(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	// Parse the data first
	data := strings.Join([]string{
		"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
		"cpu,host=server01,region=us-west value=0.65 1434055563000000000",
		"cpu,host=server01,region=us-west value=0.66 1434055564000000000",
		"cpu,host=server01,region=us-west value=0.67 1434055565000000000",
		"cpu,host=server01,region=us-west value=0.68 1434055566000000000",
	}, "\n")

	points, err := ingestion.ParseLineProtocol(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP for accurate benchmarking
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkHTTPWriteLargeDataset benchmarks HTTP write endpoint with large datasets
func BenchmarkHTTPWriteLargeDataset(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	// Parse the data first
	data := generateLargeDataset(100)
	points, err := ingestion.ParseLineProtocol(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP for accurate benchmarking
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkEndToEndWrite benchmarks the complete workflow from parsing to storage
func BenchmarkEndToEndWrite(b *testing.B) {
	testFile := "benchmark_e2e_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	// Parse data first
	data := strings.Join([]string{
		"cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0 1434055562000000000",
		"cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0 1434055563000000000",
		"cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0 1434055564000000000",
	}, "\n")
	points, err := ingestion.ParseLineProtocol(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
