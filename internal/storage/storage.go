package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/types"
)

// Storage persists time-series data in TSV format
type Storage struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex
}

// NewStorage creates or opens a storage file
func NewStorage(path string) *Storage {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Error opening storage file: %v", err)
	}
	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	return &Storage{file: f, writer: writer}
}

// WritePoint writes a time-series point to file
func (s *Storage) WritePoint(p types.Point) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range p.Fields {
		row := []string{
			p.Measurement,
			formatTags(p.Tags),
			k,
			formatFloat(v),
			p.Timestamp.Format(time.RFC3339Nano),
		}
		err := s.writer.Write(row)
		if err != nil {
			return err
		}
	}
	s.writer.Flush()
	return nil
}

func formatTags(tags map[string]string) string {
	var parts []string
	for k, v := range tags {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ",")
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Close closes the file
func (s *Storage) Close() {
	s.file.Close()
}

// Clear truncates the file to clear all data (useful for testing)
func (s *Storage) Clear() error {
	// Close current file
	s.file.Close()

	// Truncate file to 0 bytes
	err := os.Truncate(s.file.Name(), 0)
	if err != nil {
		return err
	}

	// Reopen file for writing
	f, err := os.OpenFile(s.file.Name(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	s.file = f
	s.writer = csv.NewWriter(f)
	s.writer.Comma = '\t'
	return nil
}
