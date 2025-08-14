package storage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// SegmentWriter writes memtables to immutable on-disk segments
type SegmentWriter struct {
	mu          sync.Mutex
	segmentsDir string
	nextID      uint64
}

// SegmentHeader contains metadata about a segment
type SegmentHeader struct {
	ID          uint64                 `json:"id"`
	CreatedAt   time.Time              `json:"created_at"`
	SeriesCount int                    `json:"series_count"`
	PointCount  int                    `json:"point_count"`
	MinTime     time.Time              `json:"min_time"`
	MaxTime     time.Time              `json:"max_time"`
	Checksum    uint32                 `json:"checksum"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// SegmentWriterConfig holds configuration for segment writing
type SegmentWriterConfig struct {
	SegmentsDir string
	Compression bool
	BufferSize  int
}

// NewSegmentWriter creates a new segment writer
func NewSegmentWriter(config SegmentWriterConfig) (*SegmentWriter, error) {
	// Ensure segments directory exists
	if err := os.MkdirAll(config.SegmentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create segments directory: %w", err)
	}

	return &SegmentWriter{
		segmentsDir: config.SegmentsDir,
		nextID:      uint64(time.Now().UnixNano()),
	}, nil
}

// WriteMemTable writes a memtable to an immutable segment
func (sw *SegmentWriter) WriteMemTable(memTable *MemTable) (*Segment, error) {
	if memTable == nil {
		return nil, fmt.Errorf("memtable cannot be nil")
	}

	if memTable.Data == nil {
		return nil, fmt.Errorf("memtable data cannot be nil")
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	segmentID := sw.nextID
	sw.nextID++

	// Create segment file path
	segmentPath := filepath.Join(sw.segmentsDir, fmt.Sprintf("segment_%d.seg", segmentID))

	// Open segment file for writing
	file, err := os.Create(segmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create segment file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Calculate segment metadata
	header, err := sw.calculateHeader(memTable, segmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate segment header: %w", err)
	}

	// Write header
	if err := sw.writeHeader(writer, header); err != nil {
		return nil, fmt.Errorf("failed to write segment header: %w", err)
	}

	// Write series data
	if err := sw.writeSeriesData(writer, memTable); err != nil {
		return nil, fmt.Errorf("failed to write series data: %w", err)
	}

	// Flush writer
	if err := writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush segment writer: %w", err)
	}

	// Get final file size
	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat segment file: %w", err)
	}

	// Create segment object
	segment := &Segment{
		ID:        segmentID,
		Path:      segmentPath,
		Size:      stat.Size(),
		MinTime:   header.MinTime,
		MaxTime:   header.MaxTime,
		SeriesIDs: sw.getSeriesIDs(memTable),
		CreatedAt: header.CreatedAt,
	}

	return segment, nil
}

// calculateHeader calculates the segment header from memtable data
func (sw *SegmentWriter) calculateHeader(memTable *MemTable, segmentID uint64) (*SegmentHeader, error) {
	header := &SegmentHeader{
		ID:          segmentID,
		CreatedAt:   time.Now(),
		SeriesCount: len(memTable.Data),
		PointCount:  0,
		MinTime:     time.Time{},
		MaxTime:     time.Time{},
		Metadata:    make(map[string]interface{}),
	}

	// Calculate point count and time range
	for _, points := range memTable.Data {
		header.PointCount += len(points)

		for _, point := range points {
			if header.MinTime.IsZero() || point.Timestamp.Before(header.MinTime) {
				header.MinTime = point.Timestamp
			}
			if header.MaxTime.IsZero() || point.Timestamp.After(header.MaxTime) {
				header.MaxTime = point.Timestamp
			}
		}
	}

	// Calculate checksum
	header.Checksum = sw.calculateSegmentChecksum(memTable)

	return header, nil
}

// writeHeader writes the segment header to the file
func (sw *SegmentWriter) writeHeader(writer *bufio.Writer, header *SegmentHeader) error {
	// Serialize header
	headerData, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal segment header: %w", err)
	}

	// Write header length
	headerLen := uint32(len(headerData))
	if err := binary.Write(writer, binary.LittleEndian, headerLen); err != nil {
		return fmt.Errorf("failed to write header length: %w", err)
	}

	// Write header data
	if _, err := writer.Write(headerData); err != nil {
		return fmt.Errorf("failed to write header data: %w", err)
	}

	return nil
}

// writeSeriesData writes the series data to the segment
func (sw *SegmentWriter) writeSeriesData(writer *bufio.Writer, memTable *MemTable) error {
	// Sort series IDs for consistent ordering
	seriesIDs := make([]string, 0, len(memTable.Data))
	for seriesID := range memTable.Data {
		seriesIDs = append(seriesIDs, seriesID)
	}
	sort.Strings(seriesIDs)

	// Write each series
	for _, seriesID := range seriesIDs {
		points := memTable.Data[seriesID]

		// Write series header
		seriesHeader := struct {
			SeriesID   string `json:"series_id"`
			PointCount int    `json:"point_count"`
		}{
			SeriesID:   seriesID,
			PointCount: len(points),
		}

		seriesHeaderData, err := json.Marshal(seriesHeader)
		if err != nil {
			return fmt.Errorf("failed to marshal series header: %w", err)
		}

		// Write series header length and data
		seriesHeaderLen := uint32(len(seriesHeaderData))
		if err := binary.Write(writer, binary.LittleEndian, seriesHeaderLen); err != nil {
			return fmt.Errorf("failed to write series header length: %w", err)
		}

		if _, err := writer.Write(seriesHeaderData); err != nil {
			return fmt.Errorf("failed to write series header data: %w", err)
		}

		// Write points
		for _, point := range points {
			if err := sw.writePoint(writer, point); err != nil {
				return fmt.Errorf("failed to write point: %w", err)
			}
		}
	}

	return nil
}

// writePoint writes a single data point to the segment
func (sw *SegmentWriter) writePoint(writer *bufio.Writer, point DataPoint) error {
	// Serialize point
	pointData, err := json.Marshal(point)
	if err != nil {
		return fmt.Errorf("failed to marshal point: %w", err)
	}

	// Write point length
	pointLen := uint32(len(pointData))
	if err := binary.Write(writer, binary.LittleEndian, pointLen); err != nil {
		return fmt.Errorf("failed to write point length: %w", err)
	}

	// Write point data
	if _, err := writer.Write(pointData); err != nil {
		return fmt.Errorf("failed to write point data: %w", err)
	}

	return nil
}

// calculateSegmentChecksum calculates a checksum for the entire segment
func (sw *SegmentWriter) calculateSegmentChecksum(memTable *MemTable) uint32 {
	var checksum uint32

	for seriesID, points := range memTable.Data {
		for _, b := range []byte(seriesID) {
			checksum += uint32(b)
		}

		for _, point := range points {
			checksum += uint32(point.Timestamp.Unix())
			checksum += uint32(point.Value * 1000)
		}
	}

	return checksum
}

// getSeriesIDs extracts series IDs from memtable
func (sw *SegmentWriter) getSeriesIDs(memTable *MemTable) []string {
	seriesIDs := make([]string, 0, len(memTable.Data))
	for seriesID := range memTable.Data {
		seriesIDs = append(seriesIDs, seriesID)
	}
	sort.Strings(seriesIDs)
	return seriesIDs
}

// GetSegmentsDir returns the segments directory path
func (sw *SegmentWriter) GetSegmentsDir() string {
	return sw.segmentsDir
}

// GetNextID returns the next available segment ID
func (sw *SegmentWriter) GetNextID() uint64 {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.nextID
}
