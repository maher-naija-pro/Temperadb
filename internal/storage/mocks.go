package storage

import (
	"time"
)

// MockSegmentReader is a mock implementation of SegmentReader for testing
type MockSegmentReader struct {
	ListSegmentsFunc     func() ([]*Segment, error)
	ReadSegmentFunc      func(segmentPath string) (*Segment, []SegmentReadResult, error)
	ReadSegmentRangeFunc func(segmentPath string, start, end time.Time) ([]SegmentReadResult, error)
	GetSegmentsDirFunc   func() string
}

// MockSegmentWriter is a mock implementation of SegmentWriter for testing
type MockSegmentWriter struct {
	WriteMemTableFunc  func(memTable *MemTable) (*Segment, error)
	GetSegmentsDirFunc func() string
	GetNextIDFunc      func() uint64
}

// Implement the required methods for MockSegmentReader
func (m *MockSegmentReader) ListSegments() ([]*Segment, error) {
	if m.ListSegmentsFunc != nil {
		return m.ListSegmentsFunc()
	}
	return []*Segment{}, nil
}

func (m *MockSegmentReader) ReadSegment(segmentPath string) (*Segment, []SegmentReadResult, error) {
	if m.ReadSegmentFunc != nil {
		return m.ReadSegmentFunc(segmentPath)
	}
	return &Segment{}, []SegmentReadResult{}, nil
}

func (m *MockSegmentReader) ReadSegmentRange(segmentPath string, start, end time.Time) ([]SegmentReadResult, error) {
	if m.ReadSegmentRangeFunc != nil {
		return m.ReadSegmentRangeFunc(segmentPath, start, end)
	}
	return []SegmentReadResult{}, nil
}

func (m *MockSegmentReader) GetSegmentsDir() string {
	if m.GetSegmentsDirFunc != nil {
		return m.GetSegmentsDirFunc()
	}
	return ""
}

// Implement the required methods for MockSegmentWriter
func (m *MockSegmentWriter) WriteMemTable(memTable *MemTable) (*Segment, error) {
	if m.WriteMemTableFunc != nil {
		return m.WriteMemTableFunc(memTable)
	}
	return &Segment{}, nil
}

func (m *MockSegmentWriter) GetSegmentsDir() string {
	if m.GetSegmentsDirFunc != nil {
		return m.GetSegmentsDirFunc()
	}
	return ""
}

func (m *MockSegmentWriter) GetNextID() uint64 {
	if m.GetNextIDFunc != nil {
		return m.GetNextIDFunc()
	}
	return 0
}
