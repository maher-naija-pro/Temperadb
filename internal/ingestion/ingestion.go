package ingestion

import (
	"strconv"
	"strings"
	"time"
	"timeseriesdb/internal/errors"
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
			return nil, errors.NewValidationError("invalid line format: expected 3 parts, got " + strconv.Itoa(len(parts)))
		}

		// Parse measurement and tags
		measurementAndTags := strings.Split(parts[0], ",")
		measurement := measurementAndTags[0]
		if measurement == "" {
			return nil, errors.NewValidationError("missing measurement name")
		}

		tags := map[string]string{}
		for _, tag := range measurementAndTags[1:] {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) != 2 {
				return nil, errors.NewValidationError("malformed tag: " + tag)
			}
			if kv[0] == "" || kv[1] == "" {
				return nil, errors.NewValidationError("invalid tag key or value: " + tag)
			}
			tags[kv[0]] = kv[1]
		}

		// Parse fields
		fields := map[string]float64{}
		fieldPairs := strings.Split(parts[1], ",")
		if len(fieldPairs) == 0 {
			return nil, errors.NewValidationError("no fields provided")
		}

		for _, fieldPair := range fieldPairs {
			kv := strings.SplitN(fieldPair, "=", 2)
			if len(kv) != 2 {
				return nil, errors.NewValidationError("malformed field: " + fieldPair)
			}
			if kv[0] == "" {
				return nil, errors.NewValidationError("empty field name")
			}

			val, err := strconv.ParseFloat(strings.TrimSuffix(kv[1], "i"), 64)
			if err != nil {
				return nil, errors.WrapWithType(err, errors.ErrorTypeValidation, "invalid field value '"+kv[1]+"'")
			}
			fields[kv[0]] = val
		}

		// Parse timestamp
		tsInt, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, errors.WrapWithType(err, errors.ErrorTypeValidation, "invalid timestamp")
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
