package types

import (
	"testing"
	"time"
)

func TestPoint(t *testing.T) {
	t.Run("PointCreation", func(t *testing.T) {
		// Test creating a new Point
		now := time.Now()
		point := Point{
			Measurement: "cpu",
			Tags: map[string]string{
				"host":   "server1",
				"region": "us-west",
			},
			Fields: map[string]float64{
				"usage": 75.5,
				"idle":  24.5,
			},
			Timestamp: now,
		}

		// Verify Point fields
		if point.Measurement != "cpu" {
			t.Errorf("Expected measurement 'cpu', got '%s'", point.Measurement)
		}

		if len(point.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(point.Tags))
		}

		if point.Tags["host"] != "server1" {
			t.Errorf("Expected host tag 'server1', got '%s'", point.Tags["host"])
		}

		if point.Tags["region"] != "us-west" {
			t.Errorf("Expected region tag 'us-west', got '%s'", point.Tags["region"])
		}

		if len(point.Fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(point.Fields))
		}

		if point.Fields["usage"] != 75.5 {
			t.Errorf("Expected usage field 75.5, got %f", point.Fields["usage"])
		}

		if point.Fields["idle"] != 24.5 {
			t.Errorf("Expected idle field 24.5, got %f", point.Fields["idle"])
		}

		if !point.Timestamp.Equal(now) {
			t.Errorf("Expected timestamp %v, got %v", now, point.Timestamp)
		}
	})

	t.Run("PointWithEmptyTags", func(t *testing.T) {
		// Test creating a Point with empty tags
		point := Point{
			Measurement: "memory",
			Tags:        map[string]string{},
			Fields: map[string]float64{
				"used": 1024.0,
			},
			Timestamp: time.Now(),
		}

		if point.Measurement != "memory" {
			t.Errorf("Expected measurement 'memory', got '%s'", point.Measurement)
		}

		if len(point.Tags) != 0 {
			t.Errorf("Expected 0 tags, got %d", len(point.Tags))
		}

		if point.Fields["used"] != 1024.0 {
			t.Errorf("Expected used field 1024.0, got %f", point.Fields["used"])
		}
	})

	t.Run("PointWithEmptyFields", func(t *testing.T) {
		// Test creating a Point with empty fields
		point := Point{
			Measurement: "network",
			Tags: map[string]string{
				"interface": "eth0",
			},
			Fields:    map[string]float64{},
			Timestamp: time.Now(),
		}

		if point.Measurement != "network" {
			t.Errorf("Expected measurement 'network', got '%s'", point.Measurement)
		}

		if len(point.Tags) != 1 {
			t.Errorf("Expected 1 tag, got %d", len(point.Tags))
		}

		if len(point.Fields) != 0 {
			t.Errorf("Expected 0 fields, got %d", len(point.Fields))
		}
	})

	t.Run("PointWithNilTags", func(t *testing.T) {
		// Test creating a Point with nil tags
		point := Point{
			Measurement: "disk",
			Tags:        nil,
			Fields: map[string]float64{
				"free": 2048.0,
			},
			Timestamp: time.Now(),
		}

		if point.Measurement != "disk" {
			t.Errorf("Expected measurement 'disk', got '%s'", point.Measurement)
		}

		if point.Tags != nil {
			t.Errorf("Expected nil tags, got %v", point.Tags)
		}

		if point.Fields["free"] != 2048.0 {
			t.Errorf("Expected free field 2048.0, got %f", point.Fields["free"])
		}
	})

	t.Run("PointWithNilFields", func(t *testing.T) {
		// Test creating a Point with nil fields
		point := Point{
			Measurement: "process",
			Tags: map[string]string{
				"name": "nginx",
			},
			Fields:    nil,
			Timestamp: time.Now(),
		}

		if point.Measurement != "process" {
			t.Errorf("Expected measurement 'process', got '%s'", point.Measurement)
		}

		if point.Fields != nil {
			t.Errorf("Expected nil fields, got %v", point.Fields)
		}
	})

	t.Run("PointWithZeroTimestamp", func(t *testing.T) {
		// Test creating a Point with zero timestamp
		zeroTime := time.Time{}
		point := Point{
			Measurement: "temperature",
			Tags: map[string]string{
				"sensor": "temp1",
			},
			Fields: map[string]float64{
				"value": 23.5,
			},
			Timestamp: zeroTime,
		}

		if point.Measurement != "temperature" {
			t.Errorf("Expected measurement 'temperature', got '%s'", point.Measurement)
		}

		if !point.Timestamp.Equal(zeroTime) {
			t.Errorf("Expected zero timestamp, got %v", point.Timestamp)
		}
	})

	t.Run("PointWithSpecialCharacters", func(t *testing.T) {
		// Test creating a Point with special characters in measurement and tags
		point := Point{
			Measurement: "cpu-usage_metric",
			Tags: map[string]string{
				"host-name":   "server-01",
				"data_center": "us-west-2",
			},
			Fields: map[string]float64{
				"cpu_usage": 85.7,
			},
			Timestamp: time.Now(),
		}

		if point.Measurement != "cpu-usage_metric" {
			t.Errorf("Expected measurement 'cpu-usage_metric', got '%s'", point.Measurement)
		}

		if point.Tags["host-name"] != "server-01" {
			t.Errorf("Expected host-name tag 'server-01', got '%s'", point.Tags["host-name"])
		}

		if point.Tags["data_center"] != "us-west-2" {
			t.Errorf("Expected data_center tag 'us-west-2', got '%s'", point.Tags["data_center"])
		}
	})

	t.Run("PointWithLargeNumbers", func(t *testing.T) {
		// Test creating a Point with large numbers
		point := Point{
			Measurement: "memory",
			Tags: map[string]string{
				"host": "bigserver",
			},
			Fields: map[string]float64{
				"total": 1.073741824e+12, // 1TB
				"used":  8.589934592e+11, // 800GB
			},
			Timestamp: time.Now(),
		}

		if point.Fields["total"] != 1.073741824e+12 {
			t.Errorf("Expected total field 1.073741824e+12, got %f", point.Fields["total"])
		}

		if point.Fields["used"] != 8.589934592e+11 {
			t.Errorf("Expected used field 8.589934592e+11, got %f", point.Fields["used"])
		}
	})

	t.Run("PointWithNegativeValues", func(t *testing.T) {
		// Test creating a Point with negative values
		point := Point{
			Measurement: "temperature",
			Tags: map[string]string{
				"location": "freezer",
			},
			Fields: map[string]float64{
				"value": -18.5,
			},
			Timestamp: time.Now(),
		}

		if point.Fields["value"] != -18.5 {
			t.Errorf("Expected value field -18.5, got %f", point.Fields["value"])
		}
	})

	t.Run("PointWithDecimalPrecision", func(t *testing.T) {
		// Test creating a Point with high decimal precision
		point := Point{
			Measurement: "precision",
			Tags: map[string]string{
				"type": "sensor",
			},
			Fields: map[string]float64{
				"value": 3.14159265359,
			},
			Timestamp: time.Now(),
		}

		if point.Fields["value"] != 3.14159265359 {
			t.Errorf("Expected value field 3.14159265359, got %f", point.Fields["value"])
		}
	})
}

func TestPointValidation(t *testing.T) {
	t.Run("ValidPoint", func(t *testing.T) {
		// Test that a valid Point can be created
		point := Point{
			Measurement: "test",
			Tags:        map[string]string{"key": "value"},
			Fields:      map[string]float64{"field": 1.0},
			Timestamp:   time.Now(),
		}

		// Basic validation checks
		if point.Measurement == "" {
			t.Error("Point measurement should not be empty")
		}

		if point.Timestamp.IsZero() {
			t.Error("Point timestamp should not be zero")
		}
	})

	t.Run("EmptyMeasurement", func(t *testing.T) {
		// Test Point with empty measurement
		point := Point{
			Measurement: "",
			Tags:        map[string]string{"key": "value"},
			Fields:      map[string]float64{"field": 1.0},
			Timestamp:   time.Now(),
		}

		if point.Measurement != "" {
			t.Error("Point measurement should be empty as set")
		}
	})

	t.Run("ZeroTimestamp", func(t *testing.T) {
		// Test Point with zero timestamp
		point := Point{
			Measurement: "test",
			Tags:        map[string]string{"key": "value"},
			Fields:      map[string]float64{"field": 1.0},
			Timestamp:   time.Time{},
		}

		if !point.Timestamp.IsZero() {
			t.Error("Point timestamp should be zero as set")
		}
	})
}

func TestPointEquality(t *testing.T) {
	t.Run("EqualPoints", func(t *testing.T) {
		// Test that two identical Points are equal
		now := time.Now()
		point1 := Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]float64{"usage": 75.0},
			Timestamp:   now,
		}

		point2 := Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]float64{"usage": 75.0},
			Timestamp:   now,
		}

		// Check individual fields
		if point1.Measurement != point2.Measurement {
			t.Error("Points should have equal measurements")
		}

		if len(point1.Tags) != len(point2.Tags) {
			t.Error("Points should have equal number of tags")
		}

		if len(point1.Fields) != len(point2.Fields) {
			t.Error("Points should have equal number of fields")
		}

		if !point1.Timestamp.Equal(point2.Timestamp) {
			t.Error("Points should have equal timestamps")
		}
	})

	t.Run("DifferentPoints", func(t *testing.T) {
		// Test that different Points are not equal
		now := time.Now()
		point1 := Point{
			Measurement: "cpu",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]float64{"usage": 75.0},
			Timestamp:   now,
		}

		point2 := Point{
			Measurement: "memory",
			Tags:        map[string]string{"host": "server1"},
			Fields:      map[string]float64{"usage": 75.0},
			Timestamp:   now,
		}

		// Check that measurements are different
		if point1.Measurement == point2.Measurement {
			t.Error("Points should have different measurements")
		}
	})
}

func TestPointCopy(t *testing.T) {
	t.Run("PointCopy", func(t *testing.T) {
		// Test creating a copy of a Point
		original := Point{
			Measurement: "cpu",
			Tags: map[string]string{
				"host":   "server1",
				"region": "us-west",
			},
			Fields: map[string]float64{
				"usage": 75.5,
				"idle":  24.5,
			},
			Timestamp: time.Now(),
		}

		// Create a copy
		copied := Point{
			Measurement: original.Measurement,
			Tags:        make(map[string]string),
			Fields:      make(map[string]float64),
			Timestamp:   original.Timestamp,
		}

		// Copy tags
		for k, v := range original.Tags {
			copied.Tags[k] = v
		}

		// Copy fields
		for k, v := range original.Fields {
			copied.Fields[k] = v
		}

		// Verify copy is correct
		if copied.Measurement != original.Measurement {
			t.Error("Copied measurement should match original")
		}

		if len(copied.Tags) != len(original.Tags) {
			t.Error("Copied tags should have same length as original")
		}

		if len(copied.Fields) != len(original.Fields) {
			t.Error("Copied fields should have same length as original")
		}

		if !copied.Timestamp.Equal(original.Timestamp) {
			t.Error("Copied timestamp should match original")
		}
	})
}
