package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testPath creates a test path for testing purposes
func testPath(filename string) string {
	return fmt.Sprintf("./testdata/%s", filename)
}

func TestSegmentCreation(t *testing.T) {
	// Create test data directory
	err := os.MkdirAll("./testdata", 0755)
	assert.NoError(t, err)
	defer os.RemoveAll("./testdata")

	// Test data
	timestamp := time.Now().UnixNano()

	// Create segment
	segment := &Segment{
		ID:        12345,
		Path:      testPath("segment_12345.seg"),
		Size:      1024,
		MinTime:   time.Unix(0, timestamp),
		MaxTime:   time.Unix(0, timestamp),
		SeriesIDs: []string{"series1"},
		CreatedAt: time.Now(),
	}

	// Verify segment properties
	assert.Equal(t, uint64(12345), segment.ID)
	assert.Equal(t, testPath("segment_12345.seg"), segment.Path)
	assert.Equal(t, int64(1024), segment.Size)
	assert.Equal(t, time.Unix(0, timestamp), segment.MinTime)
	assert.Equal(t, time.Unix(0, timestamp), segment.MaxTime)
	assert.Equal(t, []string{"series1"}, segment.SeriesIDs)
	assert.NotZero(t, segment.CreatedAt)
}
