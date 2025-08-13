package benchmark

import (
	"fmt"
	"os"
	"testing"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/storage"
	"timeseriesdb/internal/types"
)

// BenchmarkWriteSinglePoint benchmarks writing a single point
func BenchmarkWriteSinglePoint(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	point := types.Point{
		Measurement: "cpu",
		Tags:        map[string]string{"host": "server01", "region": "us-west"},
		Fields:      map[string]float64{"value": 0.64},
		Timestamp:   time.Unix(0, 1434055562000000000),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := storageInstance.WritePoint(point)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWriteMultiplePoints benchmarks writing multiple points
func BenchmarkWriteMultiplePoints(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	points := make([]types.Point, 100)
	for i := 0; i < 100; i++ {
		points[i] = types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01", "region": "us-west"},
			Fields:      map[string]float64{"value": float64(i)},
			Timestamp:   time.Unix(0, 1434055562000000000+int64(i)),
		}
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

// BenchmarkWritePointWithManyFields benchmarks writing points with many fields
func BenchmarkWritePointWithManyFields(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	fields := make(map[string]float64, 50)
	for i := 0; i < 50; i++ {
		fields[fmt.Sprintf("field%d", i)] = float64(i)
	}

	point := types.Point{
		Measurement: "cpu",
		Tags:        map[string]string{"host": "server01", "region": "us-west"},
		Fields:      fields,
		Timestamp:   time.Unix(0, 1434055562000000000),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := storageInstance.WritePoint(point)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkWritePointWithManyTags benchmarks writing points with many tags
func BenchmarkWritePointWithManyTags(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	tags := make(map[string]string, 50)
	for i := 0; i < 50; i++ {
		tags[fmt.Sprintf("tag%d", i)] = fmt.Sprintf("value%d", i)
	}

	point := types.Point{
		Measurement: "cpu",
		Tags:        tags,
		Fields:      map[string]float64{"value": 0.64},
		Timestamp:   time.Unix(0, 1434055562000000000),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := storageInstance.WritePoint(point)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConcurrentWrites benchmarks concurrent write operations
func BenchmarkConcurrentWrites(b *testing.B) {
	testFile := "benchmark_concurrent_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	points := make([]types.Point, 1000)
	for i := 0; i < 1000; i++ {
		points[i] = types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01", "region": "us-west"},
			Fields:      map[string]float64{"value": float64(i)},
			Timestamp:   time.Unix(0, 1434055562000000000+int64(i)),
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			point := points[i%len(points)]
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// BenchmarkMemoryUsage benchmarks memory allocation during writes
func BenchmarkMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	testFile := "benchmark_memory_test.tsv"
	defer os.Remove(testFile)

	storageConfig := config.StorageConfig{
		DataFile:    testFile,
		MaxFileSize: 1073741824, // 1GB
		BackupDir:   "backups",
		Compression: false,
	}
	storageInstance := storage.NewStorage(storageConfig)
	defer storageInstance.Close()

	points := make([]types.Point, 100)
	for i := 0; i < 100; i++ {
		points[i] = types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01", "region": "us-west"},
			Fields:      map[string]float64{"value": float64(i)},
			Timestamp:   time.Unix(0, 1434055562000000000+int64(i)),
		}
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
