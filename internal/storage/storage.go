package storage

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"timeseriesdb/internal/config"
	"timeseriesdb/internal/logger"
	"timeseriesdb/internal/types"
)

// Storage persists time-series data in TSV format
type Storage struct {
	file        *os.File
	writer      *csv.Writer
	mu          sync.Mutex
	config      config.StorageConfig
	maxFileSize int64
}

// NewStorage creates or opens a storage file
func NewStorage(cfg config.StorageConfig) *Storage {
	f, err := os.OpenFile(cfg.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("Error opening storage file: %v", err)
	}
	writer := csv.NewWriter(f)
	writer.Comma = '\t'
	return &Storage{
		file:        f,
		writer:      writer,
		config:      cfg,
		maxFileSize: cfg.MaxFileSize,
	}
}

// WritePoint writes a time-series point to file
func (s *Storage) WritePoint(p types.Point) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we need to rotate the file
	if err := s.checkAndRotateFile(); err != nil {
		return err
	}

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

// checkAndRotateFile checks if the current file exceeds max size and rotates if needed
func (s *Storage) checkAndRotateFile() error {
	if s.maxFileSize <= 0 {
		return nil // No size limit
	}

	fileInfo, err := s.file.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() >= s.maxFileSize {
		return s.rotateFile()
	}
	return nil
}

// rotateFile rotates the current file and creates a new one
func (s *Storage) rotateFile() error {
	// Close current file
	s.file.Close()

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := s.config.BackupDir + "/" + s.config.DataFile + "." + timestamp

	// Ensure backup directory exists
	if err := os.MkdirAll(s.config.BackupDir, 0755); err != nil {
		return err
	}

	// Rename current file to backup
	if err := os.Rename(s.config.DataFile, backupPath); err != nil {
		return err
	}

	// Open new file
	f, err := os.OpenFile(s.config.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	s.file = f
	s.writer = csv.NewWriter(f)
	s.writer.Comma = '\t'

	logger.Infof("Storage file rotated: %s -> %s", s.config.DataFile, backupPath)
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
	err := os.Truncate(s.config.DataFile, 0)
	if err != nil {
		return err
	}

	// Reopen file for writing
	f, err := os.OpenFile(s.config.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	s.file = f
	s.writer = csv.NewWriter(f)
	s.writer.Comma = '\t'
	return nil
}
