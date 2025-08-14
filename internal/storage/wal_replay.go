package storage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"timeseriesdb/internal/logger"
)

// WALReplay handles replaying WAL files during recovery
type WALReplay struct {
	walDir  string
	metrics *StorageMetrics
}

// NewWALReplay creates a new WAL replay handler
func NewWALReplay(walDir string, metrics *StorageMetrics) *WALReplay {
	return &WALReplay{
		walDir:  walDir,
		metrics: metrics,
	}
}

// ReplayResult contains the result of WAL replay
type ReplayResult struct {
	Entries    []WALEntry
	SeriesData map[string][]DataPoint
	ErrorCount int
	TotalCount int
}

// Replay replays all WAL files and returns the recovered data
func (wr *WALReplay) Replay() (*ReplayResult, error) {
	startTime := time.Now()

	// Get list of WAL files
	files, err := wr.getWALFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get WAL files: %w", err)
	}

	// Safety check to prevent processing too many files
	if len(files) > 1000 {
		return nil, fmt.Errorf("too many WAL files to replay: %d (max: 1000)", len(files))
	}

	result := &ReplayResult{
		Entries:    make([]WALEntry, 0),
		SeriesData: make(map[string][]DataPoint),
		ErrorCount: 0,
		TotalCount: 0,
	}

	// Replay each file
	for _, filePath := range files {
		if err := wr.replayFile(filePath, result); err != nil {
			// Record corruption error if available
			if wr.metrics != nil {
				wr.metrics.RecordWALCorruptionError()
			}
			result.ErrorCount++
			continue
		}
	}

	// Update metrics
	if wr.metrics != nil {
		wr.metrics.RecordWALRecoveryComplete(startTime, result.TotalCount, nil)
		wr.metrics.RecordWALEntriesRead(result.TotalCount)
	}

	return result, nil
}

// getWALFiles returns all WAL files sorted by creation time
func (wr *WALReplay) getWALFiles() ([]string, error) {
	files, err := os.ReadDir(wr.walDir)
	if err != nil {
		return nil, err
	}

	var walFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.HasSuffix(fileName, ".wal") ||
			(strings.Contains(fileName, ".wal.") && !strings.HasSuffix(fileName, ".tmp")) {
			walFiles = append(walFiles, filepath.Join(wr.walDir, fileName))
		}
	}

	// Sort by modification time (oldest first)
	sort.Slice(walFiles, func(i, j int) bool {
		statI, err := os.Stat(walFiles[i])
		if err != nil {
			return false
		}
		statJ, err := os.Stat(walFiles[j])
		if err != nil {
			return false
		}
		return statI.ModTime().Before(statJ.ModTime())
	})

	return walFiles, nil
}

// replayFile replays a single WAL file
func (wr *WALReplay) replayFile(filePath string, result *ReplayResult) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open WAL file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		// Read entry length
		entryLenBytes := make([]byte, 4)
		if _, err := io.ReadFull(reader, entryLenBytes); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read entry length: %w", err)
		}

		entryLen := binary.LittleEndian.Uint32(entryLenBytes)

		// Read entry data
		entryData := make([]byte, entryLen)
		if _, err := io.ReadFull(reader, entryData); err != nil {
			return fmt.Errorf("failed to read entry data: %w", err)
		}

		// Deserialize entry
		entry, err := wr.deserializeEntry(entryData)
		if err != nil {
			result.ErrorCount++
			continue // Skip corrupted entries
		}

		// Add to result
		result.Entries = append(result.Entries, entry)

		// Reconstruct series data
		for _, point := range entry.Points {
			if result.SeriesData[entry.SeriesID] == nil {
				result.SeriesData[entry.SeriesID] = make([]DataPoint, 0)
			}
			result.SeriesData[entry.SeriesID] = append(result.SeriesData[entry.SeriesID], point)
		}
	}

	return nil
}

// deserializeEntry deserializes a WAL entry from bytes
func (wr *WALReplay) deserializeEntry(data []byte) (WALEntry, error) {
	var entry WALEntry

	// Try to deserialize as JSON first
	if err := json.Unmarshal(data, &entry); err == nil {
		return entry, nil
	}

	// If JSON fails, try legacy format or return error
	return entry, fmt.Errorf("failed to deserialize WAL entry")
}

// ValidateEntry validates a WAL entry for integrity
func (wr *WALReplay) ValidateEntry(entry WALEntry) bool {
	// Check if checksum matches
	expectedChecksum := calculateChecksum(entry.SeriesID, entry.Points[0])
	return entry.Checksum == expectedChecksum
}

// CleanupOldWALs removes WAL files older than the specified duration
func (wr *WALReplay) CleanupOldWALs(maxAge time.Duration) error {
	files, err := wr.getWALFiles()
	if err != nil {
		return fmt.Errorf("failed to get WAL files for cleanup: %w", err)
	}

	cutoffTime := time.Now().Add(-maxAge)
	var lastError error

	for _, filePath := range files {
		// Check file age
		info, err := os.Stat(filePath)
		if err != nil {
			logger.Warnf("Failed to stat WAL file %s: %v", filePath, err)
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			if err := os.Remove(filePath); err != nil {
				logger.Warnf("Failed to remove old WAL file %s: %v", filePath, err)
				lastError = err
			} else {
				logger.Debugf("Removed old WAL file: %s", filePath)
			}
		}
	}

	return lastError
}

// GetWALStats returns statistics about WAL files
func (wr *WALReplay) GetWALStats() (map[string]interface{}, error) {
	files, err := wr.getWALFiles()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_files": len(files),
		"total_size":  int64(0),
		"oldest_file": "",
		"newest_file": "",
	}

	if len(files) == 0 {
		return stats, nil
	}

	var oldestTime, newestTime time.Time

	for i, filePath := range files {
		stat, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		stats["total_size"] = stats["total_size"].(int64) + stat.Size()

		if i == 0 || stat.ModTime().Before(oldestTime) {
			oldestTime = stat.ModTime()
			stats["oldest_file"] = filepath.Base(filePath)
		}

		if i == 0 || stat.ModTime().After(newestTime) {
			newestTime = stat.ModTime()
			stats["newest_file"] = filepath.Base(filePath)
		}
	}

	return stats, nil
}
