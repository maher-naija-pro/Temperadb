package ingestion

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"timeseriesdb/internal/types"
)

func TestParseLineProtocol(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []types.Point
		expectError bool
		errorMsg    string
	}{
		{
			name:  "Single valid point",
			input: "cpu,host=server01,region=us-west value=0.64 1434055562000000000",
			expected: []types.Point{
				{
					Measurement: "cpu",
					Tags: map[string]string{
						"host":   "server01",
						"region": "us-west",
					},
					Fields: map[string]float64{
						"value": 0.64,
					},
					Timestamp: time.Unix(0, 1434055562000000000),
				},
			},
			expectError: false,
		},
		{
			name: "Multiple valid points",
			input: `cpu,host=server01,region=us-west value=0.64 1434055562000000000
cpu,host=server01,region=us-west value=0.65 1434055563000000000`,
			expected: []types.Point{
				{
					Measurement: "cpu",
					Tags: map[string]string{
						"host":   "server01",
						"region": "us-west",
					},
					Fields: map[string]float64{
						"value": 0.64,
					},
					Timestamp: time.Unix(0, 1434055562000000000),
				},
				{
					Measurement: "cpu",
					Tags: map[string]string{
						"host":   "server01",
						"region": "us-west",
					},
					Fields: map[string]float64{
						"value": 0.65,
					},
					Timestamp: time.Unix(0, 1434055563000000000),
				},
			},
			expectError: false,
		},
		{
			name:  "Point with multiple fields",
			input: "cpu,host=server01 value=0.64,load=0.85 1434055562000000000",
			expected: []types.Point{
				{
					Measurement: "cpu",
					Tags: map[string]string{
						"host": "server01",
					},
					Fields: map[string]float64{
						"value": 0.64,
						"load":  0.85,
					},
					Timestamp: time.Unix(0, 1434055562000000000),
				},
			},
			expectError: false,
		},
		{
			name:  "Point with integer field value",
			input: "cpu,host=server01 value=64i 1434055562000000000",
			expected: []types.Point{
				{
					Measurement: "cpu",
					Tags: map[string]string{
						"host": "server01",
					},
					Fields: map[string]float64{
						"value": 64.0,
					},
					Timestamp: time.Unix(0, 1434055562000000000),
				},
			},
			expectError: false,
		},
		{
			name:        "Empty input",
			input:       "",
			expected:    []types.Point{},
			expectError: false,
		},
		{
			name:        "Whitespace only",
			input:       "   \n  \t  ",
			expected:    []types.Point{},
			expectError: false,
		},
		{
			name:        "Single measurement only",
			input:       "cpu",
			expectError: true,
			errorMsg:    "invalid line format: expected 3 parts, got 1",
		},
		{
			name:        "Missing fields",
			input:       "cpu,host=server01 1434055562000000000",
			expectError: true,
			errorMsg:    "invalid line format: expected 3 parts, got 2",
		},
		{
			name:        "Empty measurement",
			input:       ",host=server01 value=0.64 1434055562000000000",
			expectError: true,
			errorMsg:    "missing measurement name",
		},
		{
			name:        "Malformed tag",
			input:       "cpu,host=server01, value=0.64 1434055562000000000",
			expectError: true,
			errorMsg:    "malformed tag: ",
		},
		{
			name:        "Empty tag key",
			input:       "cpu,=server01 value=0.64 1434055562000000000",
			expectError: true,
			errorMsg:    "invalid tag key or value: =server01",
		},
		{
			name:        "Empty tag value",
			input:       "cpu,host= value=0.64 1434055562000000000",
			expectError: true,
			errorMsg:    "invalid tag key or value: host=",
		},
		{
			name:        "No fields",
			input:       "cpu,host=server01  1434055562000000000",
			expectError: true,
			errorMsg:    "malformed field: ",
		},
		{
			name:        "Empty field name",
			input:       "cpu,host=server01 =0.64 1434055562000000000",
			expectError: true,
			errorMsg:    "empty field name",
		},
		{
			name:        "Invalid field value",
			input:       "cpu,host=server01 value=abc 1434055562000000000",
			expectError: true,
			errorMsg:    "invalid field value 'abc': strconv.ParseFloat: parsing \"abc\": invalid syntax",
		},
		{
			name:        "Invalid timestamp length",
			input:       "cpu,host=server01 value=0.64 143405556200000000",
			expectError: true,
			errorMsg:    "invalid timestamp length, expected 19 digits for nanoseconds",
		},
		{
			name:        "Invalid timestamp format",
			input:       "cpu,host=server01 value=0.64 143405556200000000a",
			expectError: true,
			errorMsg:    "invalid timestamp: strconv.ParseInt: parsing \"143405556200000000a\": invalid syntax",
		},
		{
			name: "Mixed valid and invalid lines",
			input: `cpu,host=server01 value=0.64 1434055562000000000
invalid_line
cpu,host=server02 value=0.65 1434055563000000000`,
			expectError: true,
			errorMsg:    "invalid line format: expected 3 parts, got 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseLineProtocol(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d points, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if i >= len(result) {
					t.Errorf("Missing point at index %d", i)
					continue
				}

				actual := result[i]
				if actual.Measurement != expected.Measurement {
					t.Errorf("Point %d: expected measurement '%s', got '%s'", i, expected.Measurement, actual.Measurement)
				}

				if len(actual.Tags) != len(expected.Tags) {
					t.Errorf("Point %d: expected %d tags, got %d", i, len(expected.Tags), len(actual.Tags))
				} else {
					for k, v := range expected.Tags {
						if actual.Tags[k] != v {
							t.Errorf("Point %d: expected tag '%s'='%s', got '%s'", i, k, v, actual.Tags[k])
						}
					}
				}

				if len(actual.Fields) != len(expected.Fields) {
					t.Errorf("Point %d: expected %d fields, got %d", i, len(expected.Fields), len(actual.Fields))
				} else {
					for k, v := range expected.Fields {
						if actual.Fields[k] != v {
							t.Errorf("Point %d: expected field '%s'=%f, got %f", i, k, v, actual.Fields[k])
						}
					}
				}

				if !actual.Timestamp.Equal(expected.Timestamp) {
					t.Errorf("Point %d: expected timestamp %v, got %v", i, expected.Timestamp, actual.Timestamp)
				}
			}
		})
	}
}

func TestParseLineProtocol_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "Very long measurement name",
			input:       strings.Repeat("a", 1000) + ",host=server01 value=0.64 1434055562000000000",
			description: "Should handle very long measurement names",
		},
		{
			name:        "Very long tag values",
			input:       "cpu,host=" + strings.Repeat("a", 1000) + " value=0.64 1434055562000000000",
			description: "Should handle very long tag values",
		},
		{
			name:        "Very long field names",
			input:       "cpu,host=server01 " + strings.Repeat("a", 1000) + "=0.64 1434055562000000000",
			description: "Should handle very long field names",
		},
		{
			name:        "Special characters in measurement",
			input:       "cpu-test_123,host=server01 value=0.64 1434055562000000000",
			description: "Should handle special characters in measurement names",
		},
		{
			name:        "Unicode characters",
			input:       "cpu世界,host=server01 value=0.64 1434055562000000000",
			description: "Should handle unicode characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLineProtocol(tt.input)
			if err != nil {
				t.Errorf("Failed to parse line with %s: %v", tt.description, err)
			}
		})
	}
}

func TestParseLineProtocol_Performance(t *testing.T) {
	// Create a large input with many lines
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = fmt.Sprintf("cpu,host=server%03d value=%d 1434055562000000000", i, i)
	}
	input := strings.Join(lines, "\n")

	// Test parsing performance
	start := time.Now()
	_, err := ParseLineProtocol(input)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Ensure parsing completes in reasonable time (should be under 100ms for 1000 lines)
	if duration > 100*time.Millisecond {
		t.Errorf("Parsing 1000 lines took too long: %v", duration)
	}
}

func BenchmarkParseLineProtocol_SinglePoint(b *testing.B) {
	input := "cpu,host=server01,region=us-west value=0.64 1434055562000000000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseLineProtocol(input)
		if err != nil {
			b.Fatalf("Failed to parse: %v", err)
		}
	}
}

func BenchmarkParseLineProtocol_MultiplePoints(b *testing.B) {
	input := `cpu,host=server01,region=us-west value=0.64 1434055562000000000
cpu,host=server01,region=us-west value=0.65 1434055563000000000
cpu,host=server01,region=us-west value=0.66 1434055564000000000
cpu,host=server01,region=us-west value=0.67 1434055565000000000
cpu,host=server01,region=us-west value=0.68 1434055566000000000`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseLineProtocol(input)
		if err != nil {
			b.Fatalf("Failed to parse: %v", err)
		}
	}
}
