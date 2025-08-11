package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Point represents a time-series data point
type Point struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]float64
	Timestamp   time.Time
}

// ParseLineProtocol parses InfluxDB line protocol into []Point
func ParseLineProtocol(input string) ([]Point, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var points []Point

	for _, line := range lines {
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid line format")
		}

		// Parse measurement and tags
		measurementAndTags := strings.Split(parts[0], ",")
		measurement := measurementAndTags[0]
		tags := map[string]string{}
		for _, tag := range measurementAndTags[1:] {
			kv := strings.SplitN(tag, "=", 2)
			if len(kv) == 2 {
				tags[kv[0]] = kv[1]
			}
		}

		// Parse fields
		fields := map[string]float64{}
		for _, fieldPair := range strings.Split(parts[1], ",") {
			kv := strings.SplitN(fieldPair, "=", 2)
			if len(kv) != 2 {
				continue
			}
			val, err := strconv.ParseFloat(strings.TrimSuffix(kv[1], "i"), 64)
			if err != nil {
				continue
			}
			fields[kv[0]] = val
		}

		// Parse timestamp
		tsInt, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp")
		}
		timestamp := time.Unix(0, tsInt)

		points = append(points, Point{
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
			Timestamp:   timestamp,
		})
	}

	return points, nil
}

