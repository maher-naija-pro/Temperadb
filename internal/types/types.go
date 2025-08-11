package types

import "time"

// Point represents a time-series data point
type Point struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]float64
	Timestamp   time.Time
}
