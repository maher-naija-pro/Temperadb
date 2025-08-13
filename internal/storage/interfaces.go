package storage

import "time"

// SegmentReaderInterface defines the interface for reading segments
type SegmentReaderInterface interface {
	ListSegments() ([]*Segment, error)
	ReadSegment(segmentPath string) (*Segment, []SegmentReadResult, error)
	ReadSegmentRange(segmentPath string, start, end time.Time) ([]SegmentReadResult, error)
	GetSegmentsDir() string
}

// SegmentWriterInterface defines the interface for writing segments
type SegmentWriterInterface interface {
	WriteMemTable(memTable *MemTable) (*Segment, error)
	GetSegmentsDir() string
	GetNextID() uint64
}
