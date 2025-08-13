package benchmark

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"timeseriesdb/internal/ingestion"
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

func pointToLineProtocol(p types.Point) string {
	// Build measurement and tags
	parts := []string{p.Measurement}
	for k, v := range p.Tags {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	measurementAndTags := strings.Join(parts, ",")

	// Build fields
	var fieldParts []string
	for k, v := range p.Fields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%g", k, v))
	}
	fields := strings.Join(fieldParts, ",")

	// Build timestamp
	timestamp := fmt.Sprintf("%d", p.Timestamp.UnixNano())

	return fmt.Sprintf("%s %s %s", measurementAndTags, fields, timestamp)
}

// BenchmarkParseSimpleLine benchmarks parsing of simple line protocol
func BenchmarkParseSimpleLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ingestion.ParseLineProtocol(simpleLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseComplexLine benchmarks parsing of complex line protocol
func BenchmarkParseComplexLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ingestion.ParseLineProtocol(complexLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseMultiLine benchmarks parsing of multiple lines
func BenchmarkParseMultiLine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ingestion.ParseLineProtocol(multiLine)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseLargeDataset benchmarks parsing of large datasets
func BenchmarkParseLargeDataset(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ingestion.ParseLineProtocol(largeDataset)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseLineCounts benchmarks parsing with different line counts
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
				_, err := ingestion.ParseLineProtocol(data)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkParseTagCounts benchmarks parsing with different tag counts
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
		tags := make(map[string]string, bm.tags)
		for i := 0; i < bm.tags; i++ {
			tags[fmt.Sprintf("tag%d", i)] = fmt.Sprintf("value%d", i)
		}

		point := types.Point{
			Measurement: "cpu",
			Tags:        tags,
			Fields:      map[string]float64{"value": 0.64},
			Timestamp:   time.Unix(0, 1434055562000000000),
		}

		line := pointToLineProtocol(point)

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ingestion.ParseLineProtocol(line)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkParseFieldCounts benchmarks parsing with different field counts
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
		fields := make(map[string]float64, bm.fields)
		for i := 0; i < bm.fields; i++ {
			fields[fmt.Sprintf("field%d", i)] = float64(i)
		}

		point := types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01"},
			Fields:      fields,
			Timestamp:   time.Unix(0, 1434055562000000000),
		}

		line := pointToLineProtocol(point)

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := ingestion.ParseLineProtocol(line)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
