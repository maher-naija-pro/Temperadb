package storage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SegmentReader reads data from immutable on-disk segments
type SegmentReader struct {
	mu          sync.Mutex
	segmentsDir string
}

// SegmentReadResult contains the result of reading from a segment
type SegmentReadResult struct {
	SeriesID string
	Points   []DataPoint
	Error    error
}

// NewSegmentReader creates a new segment reader
func NewSegmentReader(segmentsDir string) *SegmentReader {
	return &SegmentReader{
		segmentsDir: segmentsDir,
	}
}

// ReadSegment reads all data from a specific segment
func (sr *SegmentReader) ReadSegment(segmentPath string) (*Segment, []SegmentReadResult, error) {
	file, err := os.Open(segmentPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open segment file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read segment header
	header, err := sr.readHeader(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read segment header: %w", err)
	}

	// Create segment object
	segment := &Segment{
		ID:        header.ID,
		Path:      segmentPath,
		Size:      0, // Will be updated below
		MinTime:   header.MinTime,
		MaxTime:   header.MaxTime,
		SeriesIDs: make([]string, 0),
		CreatedAt: header.CreatedAt,
	}

	// Get file size
	if stat, err := file.Stat(); err == nil {
		segment.Size = stat.Size()
	}

	// Read series data
	results, err := sr.readSeriesData(reader, header)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read series data: %w", err)
	}

	// Extract series IDs from results
	for _, result := range results {
		if result.Error == nil {
			segment.SeriesIDs = append(segment.SeriesIDs, result.SeriesID)
		}
	}

	return segment, results, nil
}

// ReadSegmentRange reads data from a segment within a specific time range
func (sr *SegmentReader) ReadSegmentRange(segmentPath string, start, end time.Time) ([]SegmentReadResult, error) {
	file, err := os.Open(segmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open segment file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read header to get metadata
	header, err := sr.readHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read segment header: %w", err)
	}

	// Check if segment overlaps with time range
	if header.MaxTime.Before(start) || header.MinTime.After(end) {
		return []SegmentReadResult{}, nil // No overlap
	}

	// Read series data with time filtering
	return sr.readSeriesDataFiltered(reader, header, start, end)
}

// readHeader reads the segment header from the file
func (sr *SegmentReader) readHeader(reader *bufio.Reader) (*SegmentHeader, error) {
	// Read header length
	headerLenBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerLenBytes); err != nil {
		return nil, fmt.Errorf("failed to read header length: %w", err)
	}

	headerLen := binary.LittleEndian.Uint32(headerLenBytes)

	// Read header data
	headerData := make([]byte, headerLen)
	if _, err := io.ReadFull(reader, headerData); err != nil {
		return nil, fmt.Errorf("failed to read header data: %w", err)
	}

	// Deserialize header
	var header SegmentHeader
	if err := json.Unmarshal(headerData, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %w", err)
	}

	return &header, nil
}

// readSeriesData reads all series data from the segment
func (sr *SegmentReader) readSeriesData(reader *bufio.Reader, header *SegmentHeader) ([]SegmentReadResult, error) {
	var results []SegmentReadResult

	for i := 0; i < header.SeriesCount; i++ {
		result, err := sr.readSeries(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Error = err
		}
		results = append(results, result)
	}

	return results, nil
}

// readSeriesDataFiltered reads series data with time filtering
func (sr *SegmentReader) readSeriesDataFiltered(reader *bufio.Reader, header *SegmentHeader, start, end time.Time) ([]SegmentReadResult, error) {
	var results []SegmentReadResult

	for i := 0; i < header.SeriesCount; i++ {
		result, err := sr.readSeriesFiltered(reader, start, end)
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Error = err
		}

		// Only include results with points
		if len(result.Points) > 0 {
			results = append(results, result)
		}
	}

	return results, nil
}

// readSeries reads a single series from the segment
func (sr *SegmentReader) readSeries(reader *bufio.Reader) (SegmentReadResult, error) {
	result := SegmentReadResult{}

	// Read series header
	seriesHeader, err := sr.readSeriesHeader(reader)
	if err != nil {
		return result, err
	}

	result.SeriesID = seriesHeader.SeriesID

	// Read points
	points, err := sr.readPoints(reader, seriesHeader.PointCount)
	if err != nil {
		return result, err
	}

	result.Points = points
	return result, nil
}

// readSeriesFiltered reads a single series with time filtering
func (sr *SegmentReader) readSeriesFiltered(reader *bufio.Reader, start, end time.Time) (SegmentReadResult, error) {
	result := SegmentReadResult{}

	// Read series header
	seriesHeader, err := sr.readSeriesHeader(reader)
	if err != nil {
		return result, err
	}

	result.SeriesID = seriesHeader.SeriesID

	// Read points with filtering
	points, err := sr.readPointsFiltered(reader, seriesHeader.PointCount, start, end)
	if err != nil {
		return result, err
	}

	result.Points = points
	return result, nil
}

// readSeriesHeader reads a series header from the segment
func (sr *SegmentReader) readSeriesHeader(reader *bufio.Reader) (struct {
	SeriesID   string `json:"series_id"`
	PointCount int    `json:"point_count"`
}, error) {
	var seriesHeader struct {
		SeriesID   string `json:"series_id"`
		PointCount int    `json:"point_count"`
	}

	// Read series header length
	headerLenBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, headerLenBytes); err != nil {
		return seriesHeader, fmt.Errorf("failed to read series header length: %w", err)
	}

	headerLen := binary.LittleEndian.Uint32(headerLenBytes)

	// Read series header data
	headerData := make([]byte, headerLen)
	if _, err := io.ReadFull(reader, headerData); err != nil {
		return seriesHeader, fmt.Errorf("failed to read series header data: %w", err)
	}

	// Deserialize series header
	if err := json.Unmarshal(headerData, &seriesHeader); err != nil {
		return seriesHeader, fmt.Errorf("failed to unmarshal series header: %w", err)
	}

	return seriesHeader, nil
}

// readPoints reads a specified number of points from the segment
func (sr *SegmentReader) readPoints(reader *bufio.Reader, count int) ([]DataPoint, error) {
	var points []DataPoint

	for i := 0; i < count; i++ {
		point, err := sr.readPoint(reader)
		if err != nil {
			return points, err
		}
		points = append(points, point)
	}

	return points, nil
}

// readPointsFiltered reads points with time filtering
func (sr *SegmentReader) readPointsFiltered(reader *bufio.Reader, count int, start, end time.Time) ([]DataPoint, error) {
	var points []DataPoint

	for i := 0; i < count; i++ {
		point, err := sr.readPoint(reader)
		if err != nil {
			return points, err
		}

		// Apply time filter
		if (point.Timestamp.Equal(start) || point.Timestamp.After(start)) &&
			(point.Timestamp.Equal(end) || point.Timestamp.Before(end)) {
			points = append(points, point)
		}
	}

	return points, nil
}

// readPoint reads a single data point from the segment
func (sr *SegmentReader) readPoint(reader *bufio.Reader) (DataPoint, error) {
	var point DataPoint

	// Read point length
	pointLenBytes := make([]byte, 4)
	if _, err := io.ReadFull(reader, pointLenBytes); err != nil {
		return point, fmt.Errorf("failed to read point length: %w", err)
	}

	pointLen := binary.LittleEndian.Uint32(pointLenBytes)

	// Read point data
	pointData := make([]byte, pointLen)
	if _, err := io.ReadFull(reader, pointData); err != nil {
		return point, fmt.Errorf("failed to read point data: %w", err)
	}

	// Deserialize point
	if err := json.Unmarshal(pointData, &point); err != nil {
		return point, fmt.Errorf("failed to unmarshal point: %w", err)
	}

	return point, nil
}

// ListSegments returns a list of all available segments
func (sr *SegmentReader) ListSegments() ([]*Segment, error) {
	files, err := os.ReadDir(sr.segmentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read segments directory: %w", err)
	}

	var segments []*Segment

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".seg") {
			continue
		}

		segmentPath := filepath.Join(sr.segmentsDir, file.Name())

		// Try to read segment header to get metadata
		if segment, _, err := sr.ReadSegment(segmentPath); err == nil {
			segments = append(segments, segment)
		}
	}

	return segments, nil
}

// GetSegmentsDir returns the segments directory path
func (sr *SegmentReader) GetSegmentsDir() string {
	return sr.segmentsDir
}
