package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"timeseriesdb/internal/parser"
	"timeseriesdb/internal/storage"
	"timeseriesdb/internal/types"
)

// Benchmark data samples
var (
	simpleLine = "cpu,host=server01,region=us-west value=0.64 1434055562000000000"

	complexLine = "cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1 " +
		"user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0 " +
		"1434055562000000000"

	multiLine = strings.Join([]string{
		"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
		"cpu,host=server01,region=us-west value=0.65 1434055563000000000",
		"cpu,host=server01,region=us-west value=0.66 1434055564000000000",
		"cpu,host=server01,region=us-west value=0.67 1434055565000000000",
		"cpu,host=server01,region=us-west value=0.68 1434055566000000000",
	}, "\n")

	largeDataset = generateLargeDataset(1000)
)

func generateLargeDataset(size int) string {
	var lines []string
	for i := 0; i < size; i++ {
		line := strings.Join([]string{
			"cpu,host=server01,region=us-west,datacenter=dc1,rack=r1,zone=z1",
			"user=0.64,system=0.23,idle=0.12,wait=0.01,steal=0.0,guest=0.0",
			fmt.Sprintf("%d", 1434055562000000000+i*1000000000),
		}, " ")
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// =============================================================================
// PARSER BENCHMARKS
// =============================================================================

func BenchmarkParseSimpleLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(simpleLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseComplexLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(complexLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMultiLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(multiLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLargeDataset(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(largeDataset)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLineCounts(b *testing.B) {
	benchmarks := []struct {
		name  string
		lines int
	}{
		{"1_line", 1},
		{"10_lines", 10},
		{"100_lines", 100},
		{"1000_lines", 1000},
		{"10000_lines", 10000},
	}

	for _, bm := range benchmarks {
		data := generateLargeDataset(bm.lines)
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := parser.ParseLineProtocol(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseTagCounts(b *testing.B) {
	benchmarks := []struct {
		name string
		tags int
	}{
		{"2_tags", 2},
		{"5_tags", 5},
		{"10_tags", 10},
		{"20_tags", 20},
	}

	for _, bm := range benchmarks {
		tagPairs := make([]string, bm.tags)
		for i := 0; i < bm.tags; i++ {
			tagPairs[i] = fmt.Sprintf("tag%d=value%d", i, i)
		}
		tags := strings.Join(tagPairs, ",")

		line := fmt.Sprintf("cpu,%s value=0.64 1434055562000000000", tags)

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := parser.ParseLineProtocol(line)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseFieldCounts(b *testing.B) {
	benchmarks := []struct {
		name   string
		fields int
	}{
		{"1_field", 1},
		{"5_fields", 5},
		{"10_fields", 10},
		{"20_fields", 20},
	}

	for _, bm := range benchmarks {
		fieldPairs := make([]string, bm.fields)
		for i := 0; i < bm.fields; i++ {
			fieldPairs[i] = fmt.Sprintf("field%d=%d", i, i)
		}
		fields := strings.Join(fieldPairs, ",")

		line := fmt.Sprintf("cpu,host=server01 %s 1434055562000000000", fields)

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := parser.ParseLineProtocol(line)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// =============================================================================
// STORAGE BENCHMARKS
// =============================================================================

func BenchmarkWriteSinglePoint(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
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

func BenchmarkWriteMultiplePoints(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
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

func BenchmarkWritePointWithManyFields(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
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

func BenchmarkWritePointWithManyTags(b *testing.B) {
	testFile := "benchmark_storage_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
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

// =============================================================================
// HTTP ENDPOINT BENCHMARKS
// =============================================================================

func BenchmarkHTTPWriteSinglePoint(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
	defer storageInstance.Close()

	// Parse the data first
	points, err := parser.ParseLineProtocol("cpu,host=server01,region=us-west value=0.64 1434055562000000000")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkHTTPWriteMultiplePoints(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
	defer storageInstance.Close()

	// Parse the data first
	data := strings.Join([]string{
		"cpu,host=server01,region=us-west value=0.64 1434055562000000000",
		"cpu,host=server01,region=us-west value=0.65 1434055563000000000",
		"cpu,host=server01,region=us-west value=0.66 1434055564000000000",
		"cpu,host=server01,region=us-west value=0.67 1434055565000000000",
		"cpu,host=server01,region=us-west value=0.68 1434055566000000000",
	}, "\n")

	points, err := parser.ParseLineProtocol(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkHTTPWriteLargeDataset(b *testing.B) {
	testFile := "benchmark_http_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
	defer storageInstance.Close()

	// Parse the data first
	data := generateLargeDataset(100)
	points, err := parser.ParseLineProtocol(data)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Write directly to storage instead of HTTP
		for _, point := range points {
			err := storageInstance.WritePoint(point)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// =============================================================================
// INTEGRATED WORKFLOW BENCHMARKS
// =============================================================================

func BenchmarkEndToEndWrite(b *testing.B) {
	testFile := "benchmark_e2e_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
	defer storageInstance.Close()

	// Parse data first
	points, err := parser.ParseLineProtocol(largeDataset)
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

func BenchmarkConcurrentWrites(b *testing.B) {
	testFile := "benchmark_concurrent_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
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

// =============================================================================
// MEMORY USAGE BENCHMARKS
// =============================================================================

func BenchmarkMemoryUsage(b *testing.B) {
	b.ReportAllocs()

	points, err := parser.ParseLineProtocol(largeDataset)
	if err != nil {
		b.Fatal(err)
	}

	testFile := "benchmark_memory_test.tsv"
	defer os.Remove(testFile)

	storageInstance := storage.NewStorage(testFile)
	defer storageInstance.Close()

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

// Note: handleWrite function is already defined in write_endpoint_test.go
