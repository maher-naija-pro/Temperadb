package main

import (
	"fmt"
	"strings"
	"testing"
	"timeseriesdb/internal/parser"
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

// Benchmark simple line parsing
func BenchmarkParseSimpleLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(simpleLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark complex line parsing
func BenchmarkParseComplexLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(complexLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark multi-line parsing
func BenchmarkParseMultiLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(multiLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark large dataset parsing
func BenchmarkParseLargeDataset(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseLineProtocol(largeDataset)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark parsing with different line counts
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

// Benchmark parsing with different tag counts
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

// Benchmark parsing with different field counts
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
