package storage

import (
	"time"
)

// DataPoint represents a single time series data point
type DataPoint struct {
	Timestamp time.Time
	Value     float64
	Labels    map[string]string
}

// Series represents a time series with metadata
type Series struct {
	ID     string
	Name   string
	Labels map[string]string
	Points []DataPoint
}

// WriteRequest represents a write operation
type WriteRequest struct {
	SeriesID string
	Points   []DataPoint
}

// ReadRequest represents a read operation
type ReadRequest struct {
	SeriesID string
	Start    time.Time
	End      time.Time
	Limit    int
}

// StorageEngine defines the interface for the storage engine
type StorageEngine interface {
	Write(req WriteRequest) error
	Read(req ReadRequest) ([]DataPoint, error)
	Close() error
}

// MemTable represents an in-memory table for buffering writes
type MemTable struct {
	ID        uint64
	Data      map[string][]DataPoint
	Size      int64
	MaxSize   int64
	CreatedAt time.Time
	IsFlushed bool
}

// Segment represents an immutable on-disk segment
type Segment struct {
	ID        uint64
	Path      string
	Size      int64
	MinTime   time.Time
	MaxTime   time.Time
	SeriesIDs []string
	CreatedAt time.Time
}

// WALEntry represents a write-ahead log entry
type WALEntry struct {
	ID        uint64
	Timestamp time.Time
	SeriesID  string
	Points    []DataPoint
	Checksum  uint32
}

// CompactionLevel represents a level in the LSM tree
type CompactionLevel struct {
	Level    int
	Segments []*Segment
	MaxSize  int64
	MaxFiles int
}
