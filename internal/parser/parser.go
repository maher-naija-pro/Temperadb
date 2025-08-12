package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"timeseriesdb/internal/types"
)

// ParseLineProtocol parses InfluxDB line protocol into []types.Point
func ParseLineProtocol(input string) ([]types.Point, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var points []types.Point

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid line format: expected 3 parts, got %d", len(parts))
		}

		// Parse measurement and tags
		measurementAndTags := strings.Split(parts[0], ",")
		measurement := measurementAndTags[0]
		if measurement == "" {
			return nil, fmt.Errorf("missing measurement name")
		}

		tags := map[string]string{}
		for _, tag := range measurementAndTags[1:] {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("malformed tag: %s", tag)
			}
			if kv[0] == "" || kv[1] == "" {
				return nil, fmt.Errorf("invalid tag key or value: %s", tag)
			}
			tags[kv[0]] = kv[1]
		}

		// Parse fields
		fields := map[string]float64{}
		fieldPairs := strings.Split(parts[1], ",")
		if len(fieldPairs) == 0 {
			return nil, fmt.Errorf("no fields provided")
		}

		for _, fieldPair := range fieldPairs {
			kv := strings.SplitN(fieldPair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("malformed field: %s", fieldPair)
			}
			if kv[0] == "" {
				return nil, fmt.Errorf("empty field name")
			}

			val, err := strconv.ParseFloat(strings.TrimSuffix(kv[1], "i"), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid field value '%s': %v", kv[1], err)
			}
			fields[kv[0]] = val
		}

		// Parse timestamp
		tsInt, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %v", err)
		}
		timestamp := time.Unix(0, tsInt)

		points = append(points, types.Point{
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
			Timestamp:   timestamp,
		})
	}

	return points, nil
}
