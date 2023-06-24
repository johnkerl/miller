// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"time"

	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------------------------------------
type tDataForSec2GMT struct {
	epochSeconds     float64
	numDecimalPlaces int
	expectedOutput   string
}

var dataForSec2GMT = []tDataForSec2GMT{
	{0.0, 0, "1970-01-01T00:00:00Z"},
	{0.0, 6, "1970-01-01T00:00:00.000000Z"},
	{1.0, 6, "1970-01-01T00:00:01.000000Z"},
	{123456789.25, 3, "1973-11-29T21:33:09.250Z"},
}

func TestSec2GMT(t *testing.T) {
	for _, entry := range dataForSec2GMT {
		assert.Equal(t, entry.expectedOutput, Sec2GMT(entry.epochSeconds, entry.numDecimalPlaces))
	}
}

// ----------------------------------------------------------------
type tDataForNsec2GMT struct {
	epochNanoseconds int64
	numDecimalPlaces int
	expectedOutput   string
}

var dataForNsec2GMT = []tDataForNsec2GMT{
	{0, 0, "1970-01-01T00:00:00Z"},
	{0, 6, "1970-01-01T00:00:00.000000Z"},
	{946684800123456789, 0, "2000-01-01T00:00:00Z"},
	{946684800123456789, 1, "2000-01-01T00:00:00.1Z"},
	{946684800123456789, 2, "2000-01-01T00:00:00.12Z"},
	{946684800123456789, 3, "2000-01-01T00:00:00.123Z"},
	{946684800123456789, 4, "2000-01-01T00:00:00.1234Z"},
	{946684800123456789, 5, "2000-01-01T00:00:00.12345Z"},
	{946684800123456789, 6, "2000-01-01T00:00:00.123456Z"},
	{946684800123456789, 7, "2000-01-01T00:00:00.1234567Z"},
	{946684800123456789, 8, "2000-01-01T00:00:00.12345678Z"},
	{946684800123456789, 9, "2000-01-01T00:00:00.123456789Z"},
}

func TestNsec2GMT(t *testing.T) {
	for _, entry := range dataForNsec2GMT {
		actualOutput := Nsec2GMT(entry.epochNanoseconds, entry.numDecimalPlaces)
		assert.Equal(t, entry.expectedOutput, actualOutput)
	}
}

// ----------------------------------------------------------------
type tDataForEpochSecondsToGMT struct {
	epochSeconds   float64
	expectedOutput time.Time
}

var dataForEpochSecondsToGMT = []tDataForEpochSecondsToGMT{
	{0.0, time.Unix(0, 0).UTC()},
	{1.25, time.Unix(1, 250000000).UTC()},
	{123456789.25, time.Unix(123456789, 250000000).UTC()},
}

func TestEpochSecondsToGMT(t *testing.T) {
	for _, entry := range dataForEpochSecondsToGMT {
		assert.Equal(t, entry.expectedOutput, EpochSecondsToGMT(entry.epochSeconds))
	}
}

// ----------------------------------------------------------------
type tDataForEpochNanosecondsToGMT struct {
	epochNanoseconds int64
	expectedOutput   time.Time
}

var dataForEpochNanosecondsToGMT = []tDataForEpochNanosecondsToGMT{
	{0, time.Unix(0, 0).UTC()},
	{1000000000, time.Unix(1, 0).UTC()},
	{1200000000, time.Unix(1, 200000000).UTC()},
	{-1000000000, time.Unix(-1, 0).UTC()},
	{-1200000000, time.Unix(-1, -200000000).UTC()},
	{123456789250000047, time.Unix(123456789, 250000047).UTC()},
}

func TestEpochNanosecondsToGMT(t *testing.T) {
	for _, entry := range dataForEpochNanosecondsToGMT {
		assert.Equal(t, entry.expectedOutput, EpochNanosecondsToGMT(entry.epochNanoseconds))
	}
}
