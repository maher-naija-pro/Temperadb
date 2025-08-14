package storage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WAL represents a Write-Ahead Log for durability
type WAL struct {
	mu          sync.Mutex
	file        *os.File
	writer      *bufio.Writer
	path        string
	maxFileSize int64
	currentSize int64
	sequenceNum uint64
	closed      bool
	metrics     *StorageMetrics
}

// WALConfig holds configuration for WAL
type WALConfig struct {
	Path        string
	MaxFileSize int64
	Metrics     *StorageMetrics
}

// NewWAL creates a new Write-Ahead Log
func NewWAL(config WALConfig) (*WAL, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(config.Path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %w", err)
	}

	// Open or create WAL file
	file, err := os.OpenFile(config.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}

	// Get current file size
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat WAL file: %w", err)
	}

	wal := &WAL{
		file:        file,
		writer:      bufio.NewWriter(file),
		path:        config.Path,
		maxFileSize: config.MaxFileSize,
		currentSize: stat.Size(),
		sequenceNum: uint64(time.Now().UnixNano()),
		closed:      false,
		metrics:     config.Metrics,
	}

	// Update metrics if available
	if wal.metrics != nil {
		wal.metrics.RecordWALSize(wal.currentSize)
		wal.metrics.RecordWALFileCount(1) // Single WAL file initially
	}

	return wal, nil
}

// Write writes a WAL entry to the log
func (w *WAL) Write(entry WALEntry) error {
	startTime := time.Now()

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("WAL is closed")
	}

	// Check if we need to rotate the file
	if w.currentSize >= w.maxFileSize {
		rotationStartTime := time.Now()
		if err := w.rotateFile(); err != nil {
			if w.metrics != nil {
				w.metrics.RecordWALFileRotationComplete(rotationStartTime, err)
			}
			return fmt.Errorf("failed to rotate WAL file: %w", err)
		}
		if w.metrics != nil {
			w.metrics.RecordWALFileRotationComplete(rotationStartTime, nil)
		}
	}

	// Assign sequence number
	entry.ID = w.sequenceNum
	w.sequenceNum++

	// Serialize entry
	data, err := w.serializeEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to serialize WAL entry: %w", err)
	}

	// Write entry length and data
	entryLen := uint32(len(data))
	if err := binary.Write(w.writer, binary.LittleEndian, entryLen); err != nil {
		return fmt.Errorf("failed to write entry length: %w", err)
	}

	if _, err := w.writer.Write(data); err != nil {
		return fmt.Errorf("failed to write entry data: %w", err)
	}

	// Update current size
	w.currentSize += int64(4 + len(data)) // 4 bytes for length + data

	// Update metrics
	if w.metrics != nil {
		w.metrics.RecordWALSize(w.currentSize)
		w.metrics.RecordWALEntriesWritten(1)
		w.metrics.RecordWALEntrySize(len(data))
		w.metrics.RecordWALWriteLatency(time.Since(startTime))
	}

	return nil
}

// Flush flushes the WAL buffer to disk
func (w *WAL) Flush() error {
	startTime := time.Now()

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("WAL is closed")
	}

	if err := w.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush WAL buffer: %w", err)
	}

	// Update metrics
	if w.metrics != nil {
		// Record WAL flush operation and latency
		WALFlushOperations.WithLabelValues("success").Inc()
		WALFlushLatency.WithLabelValues().Observe(time.Since(startTime).Seconds())
	}

	return nil
}

// Close closes the WAL
func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}

	w.closed = true

	// Flush any remaining data
	if err := w.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush WAL: %w", err)
	}

	// Close the file
	if err := w.file.Close(); err != nil {
		return fmt.Errorf("failed to close WAL file: %w", err)
	}

	return nil
}

// rotateFile rotates the current WAL file and creates a new one
func (w *WAL) rotateFile() error {
	// Close current file
	if err := w.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush WAL buffer before rotation: %w", err)
	}

	if err := w.file.Close(); err != nil {
		return fmt.Errorf("failed to close WAL file: %w", err)
	}

	// Rename current file with timestamp
	timestamp := time.Now().Format("20060102-150405.000")
	oldPath := w.path + "." + timestamp
	if err := os.Rename(w.path, oldPath); err != nil {
		return fmt.Errorf("failed to rename old WAL file: %w", err)
	}

	// Create new WAL file
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new WAL file: %w", err)
	}

	w.file = file
	w.writer = bufio.NewWriter(file)
	w.currentSize = 0
	w.sequenceNum = uint64(time.Now().UnixNano())

	// Update metrics
	if w.metrics != nil {
		w.metrics.RecordWALSize(0) // New file starts with 0 size
		// Note: File count is managed externally since we don't track all files here
	}

	return nil
}

// serializeEntry serializes a WAL entry to bytes
func (w *WAL) serializeEntry(entry WALEntry) ([]byte, error) {
	// Create a serializable version of the entry
	serializableEntry := struct {
		ID        uint64      `json:"id"`
		Timestamp time.Time   `json:"timestamp"`
		SeriesID  string      `json:"series_id"`
		Points    []DataPoint `json:"points"`
		Checksum  uint32      `json:"checksum"`
	}{
		ID:        entry.ID,
		Timestamp: entry.Timestamp,
		SeriesID:  entry.SeriesID,
		Points:    entry.Points,
		Checksum:  entry.Checksum,
	}

	return json.Marshal(serializableEntry)
}

// GetPath returns the WAL file path
func (w *WAL) GetPath() string {
	return w.path
}

// GetSize returns the current WAL file size
func (w *WAL) GetSize() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.currentSize
}

// IsClosed returns whether the WAL is closed
func (w *WAL) IsClosed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closed
}
