package fetching

import (
	"testing"
	"time"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		requestedRange                     string
		expectedStart, expectedEnd, endMax int64
		expectError                        bool
	}{
		{"bytez=1-5", 0, 0, 10, true},
		{"bytes=1-5", 1, 5, 10, false},
		{"bytes=1-", 1, 10, 10, false},
		{"bytes=  1   -  5 ", 1, 5, 10, false},
	}
	for i, test := range tests {
		rangeStart, rangeEnd, err := ParseRange(test.requestedRange, test.endMax)
		if rangeStart != test.expectedStart {
			t.Errorf("Test %d: rangeStart != expectedStart; %d != %d", i+1, rangeStart, test.expectedStart)
		}
		if rangeEnd != test.expectedEnd {
			t.Errorf("Test %d: rangeEnd != expectedEnd; %d != %d", i+1, rangeEnd, test.expectedEnd)
		}
		if (err != nil) != test.expectError {
			t.Errorf("Test %d: error: %v; expected: %v", i+1, err, test.expectError)
		}
	}
}

func TestParseTotalLengthFromContentRangeHeader(t *testing.T) {
	tests := []struct {
		contentRange   string
		expectedLength int64
	}{
		{"3-5/1000", 1000},
		{"/42", 42},
		{"/", -1},
		{"", -1},
	}
	for i, test := range tests {
		length := parseTotalLengthFromContentRangeHeader(test.contentRange)
		if length != test.expectedLength {
			t.Errorf("Test #%d for \"%s\" returned %d. Expected %d.", i+1, test.contentRange, length, test.expectedLength)
		}
	}
}

func TestBitRateToByteDuration(t *testing.T) {
	tests := []struct {
		rate     int64
		expected time.Duration
	}{
		{8, time.Second},
		{80, time.Second / 10},
		{80000, time.Millisecond / 10},
		{1000, time.Millisecond * 8},
	}
	for _, test := range tests {
		result := BitRateToByteDuration(test.rate)
		if result != test.expected {
			t.Errorf("byteDuration was %v. Expected %v.", result, test.expected)
		}
	}
}
