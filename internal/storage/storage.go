// Package storage provides functionality for persisting time-series data to disk
package storage

import (
	"encoding/csv"                 // Import CSV package for writing tab-separated values
	"os"                           // Import OS package for file operations
	"strconv"                      // Import strconv package for string conversions
	"strings"                      // Import strings package for string manipulation
	"sync"                         // Import sync package for thread-safe operations
	"time"                         // Import time package for timestamp handling
	"timeseriesdb/internal/config" // Import internal config package
	"timeseriesdb/internal/errors" // Import internal errors package
	"timeseriesdb/internal/logger" // Import internal logger package
	"timeseriesdb/internal/types"  // Import internal types package
)

// Storage persists time-series data in TSV format
// This struct holds all the necessary components for file-based storage
type Storage struct {
	file        *os.File             // File handle for the data file
	writer      *csv.Writer          // CSV writer configured for TSV format
	mu          sync.Mutex           // Mutex for thread-safe file operations
	config      config.StorageConfig // Configuration settings for storage
	maxFileSize int64                // Maximum file size before rotation (in bytes)
}

// NewStorage creates or opens a storage file
// Takes a StorageConfig and returns a configured Storage instance
func NewStorage(cfg config.StorageConfig) *Storage {
	// Open file in append mode, create if doesn't exist, write-only permissions
	f, err := os.OpenFile(cfg.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// If file cannot be opened, log fatal error and exit
		logger.Fatalf("Error opening storage file: %v", err)
	}
	// Create a new CSV writer for the file
	writer := csv.NewWriter(f)
	// Set the delimiter to tab character for TSV format
	writer.Comma = '\t'
	// Return a new Storage instance with all fields initialized
	return &Storage{
		file:        f,               // Store the file handle
		writer:      writer,          // Store the CSV writer
		config:      cfg,             // Store the configuration
		maxFileSize: cfg.MaxFileSize, // Store the max file size setting
	}
}

// WritePoint writes a time-series point to file
// Takes a Point struct and writes all its fields to the storage file
func (s *Storage) WritePoint(p types.Point) error {
	// Lock the mutex to ensure thread-safe file access
	s.mu.Lock()
	// Ensure mutex is unlocked when function returns
	defer s.mu.Unlock()

	// Check if we need to rotate the file due to size limits
	if err := s.checkAndRotateFile(); err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "file rotation failed")
	}

	// Iterate through all fields in the point
	for k, v := range p.Fields {
		// Create a row with: measurement, tags, field name, field value, timestamp
		row := []string{
			p.Measurement,                        // Measurement name (e.g., "cpu_usage")
			formatTags(p.Tags),                   // Formatted tags string
			k,                                    // Field name (e.g., "value")
			formatFloat(v),                       // Formatted field value
			p.Timestamp.Format(time.RFC3339Nano), // ISO timestamp with nanoseconds
		}
		// Write the row to the TSV file
		err := s.writer.Write(row)
		if err != nil {
			return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to write point to storage")
		}
	}
	// Flush the writer to ensure data is written to disk
	s.writer.Flush()
	return nil // Return success
}

// checkAndRotateFile checks if the current file exceeds max size and rotates if needed
// This method ensures files don't grow indefinitely large
func (s *Storage) checkAndRotateFile() error {
	// If maxFileSize is 0 or negative, no size limit is enforced
	if s.maxFileSize <= 0 {
		return nil // No size limit, so no rotation needed
	}

	// Get file information to check current size
	fileInfo, err := s.file.Stat()
	if err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to get file info")
	}

	// Check if current file size exceeds the maximum allowed size
	if fileInfo.Size() >= s.maxFileSize {
		// If size limit exceeded, rotate the file
		return s.rotateFile()
	}
	return nil // File size is within limits
}

// rotateFile rotates the current file and creates a new one
// This creates a backup of the current file and starts a fresh one
func (s *Storage) rotateFile() error {
	// Close the current file handle
	s.file.Close()

	// Create backup filename with current timestamp
	timestamp := time.Now().Format("20060102-150405") // Format: YYYYMMDD-HHMMSS
	backupPath := s.config.BackupDir + "/" + s.config.DataFile + "." + timestamp

	// Ensure the backup directory exists, create it if it doesn't
	if err := os.MkdirAll(s.config.BackupDir, 0755); err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to create backup directory")
	}

	// Rename the current data file to the backup location
	if err := os.Rename(s.config.DataFile, backupPath); err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to rename file for rotation")
	}

	// Open a new file with the same name for continued writing
	f, err := os.OpenFile(s.config.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to open new file after rotation")
	}

	// Update the storage struct with the new file and writer
	s.file = f
	s.writer = csv.NewWriter(f)
	s.writer.Comma = '\t' // Set tab as delimiter for TSV format

	// Log the successful rotation for monitoring purposes
	logger.Infof("Storage file rotated: %s -> %s", s.config.DataFile, backupPath)
	return nil // Return success
}

// formatTags converts a map of tags to a comma-separated string
// Example: {"host": "server1", "region": "us-west"} -> "host=server1,region=us-west"
func formatTags(tags map[string]string) string {
	var parts []string // Slice to hold formatted tag parts
	// Iterate through each tag key-value pair
	for k, v := range tags {
		// Format each tag as "key=value" and add to parts slice
		parts = append(parts, k+"="+v)
	}
	// Join all parts with commas and return the result
	return strings.Join(parts, ",")
}

// formatFloat converts a float64 to a string representation
// Uses scientific notation if needed, otherwise decimal notation
func formatFloat(f float64) string {
	// Convert float to string with full precision
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Close closes the file and releases system resources
// This should be called when the storage is no longer needed
func (s *Storage) Close() {
	s.file.Close() // Close the underlying file handle
}

// Clear truncates the file to clear all data (useful for testing)
// This method removes all data from the storage file
func (s *Storage) Clear() error {
	// Close the current file handle
	s.file.Close()

	// Truncate the file to 0 bytes, effectively clearing all data
	err := os.Truncate(s.config.DataFile, 0)
	if err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to truncate file")
	}

	// Reopen the file for writing (it will be empty now)
	f, err := os.OpenFile(s.config.DataFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WrapWithType(err, errors.ErrorTypeStorage, "failed to reopen file after clearing")
	}

	// Update the storage struct with the new empty file and writer
	s.file = f
	s.writer = csv.NewWriter(f)
	s.writer.Comma = '\t' // Set tab as delimiter for TSV format
	return nil            // Return success
}
