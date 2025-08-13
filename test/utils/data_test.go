package utils

import (
	"fmt"
	"strings"
	"time"
	"timeseriesdb/internal/types"
)

// TestDataFactory provides reusable test data for tests
type TestDataFactory struct{}

// SimplePoint returns a basic test point
func (f *TestDataFactory) SimplePoint() types.Point {
	return types.Point{
		Measurement: "cpu",
		Tags:        map[string]string{"host": "server01", "region": "us-west"},
		Fields:      map[string]float64{"value": 0.64},
		Timestamp:   time.Unix(0, 1434055562000000000),
	}
}

// ComplexPoint returns a point with many tags and fields
func (f *TestDataFactory) ComplexPoint() types.Point {
	tags := make(map[string]string, 5)
	tags["host"] = "server01"
	tags["region"] = "us-west"
	tags["datacenter"] = "dc1"
	tags["rack"] = "r1"
	tags["zone"] = "z1"

	fields := make(map[string]float64, 6)
	fields["user"] = 0.64
	fields["system"] = 0.23
	fields["idle"] = 0.12
	fields["wait"] = 0.01
	fields["steal"] = 0.0
	fields["guest"] = 0.0

	return types.Point{
		Measurement: "cpu",
		Tags:        tags,
		Fields:      fields,
		Timestamp:   time.Unix(0, 1434055562000000000),
	}
}

// MultiPoint returns multiple test points
func (f *TestDataFactory) MultiPoint(count int) []types.Point {
	points := make([]types.Point, count)
	for i := 0; i < count; i++ {
		points[i] = types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01", "region": "us-west"},
			Fields:      map[string]float64{"value": float64(i)},
			Timestamp:   time.Unix(0, 1434055562000000000+int64(i)),
		}
	}
	return points
}

// LineProtocol converts points to line protocol format
func (f *TestDataFactory) LineProtocol(points ...types.Point) string {
	var lines []string
	for _, p := range points {
		lines = append(lines, f.pointToLineProtocol(p))
	}
	return strings.Join(lines, "\n")
}

// pointToLineProtocol converts a single point to line protocol format
func (f *TestDataFactory) pointToLineProtocol(p types.Point) string {
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

// GenerateLargeDataset creates a dataset with the specified number of lines
func (f *TestDataFactory) GenerateLargeDataset(size int) string {
	var lines []string
	for i := 0; i < size; i++ {
		point := types.Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server01", "region": "us-west", "datacenter": "dc1", "rack": "r1", "zone": "z1"},
			Fields:      map[string]float64{"user": 0.64, "system": 0.23, "idle": 0.12, "wait": 0.01, "steal": 0.0, "guest": 0.0},
			Timestamp:   time.Unix(0, 1434055562000000000+int64(i)*1000000000),
		}
		lines = append(lines, f.pointToLineProtocol(point))
	}
	return strings.Join(lines, "\n")
}

// Global instance for easy access
var DataFactory = &TestDataFactory{}
