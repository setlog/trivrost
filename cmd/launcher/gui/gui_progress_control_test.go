package gui

import (
	"testing"
	"time"
)

func TestPanelDownloadStatusAverageConstantDelta(t *testing.T) {
	var mockedProgress uint64
	progressFunc := func() uint64 {
		return mockedProgress
	}
	const sampleCount = 42
	const sampleInterval = time.Millisecond * 200
	movingAverage := NewMovingAverage(sampleCount, sampleInterval, progressFunc)

	average := movingAverage.GetAverageDelta()
	if average != 0 {
		t.Fatalf("average was %f. Expected 0.", average)
	}

	for i := 0; i < 100; i++ {
		mockedProgress += 1000
		movingAverage.TakeSample()
		average = movingAverage.GetAverageDelta()
		if int(average+0.5) != 5000 {
			t.Fatalf("average was %d. Expected 5000. i = %d", int(average+0.5), i)
		}
	}
}

func TestPanelDownloadStatusAverageChaotic(t *testing.T) {
	var mockedProgress uint64
	progressFunc := func() uint64 {
		return mockedProgress
	}
	const sampleCount = 8
	const sampleInterval = time.Millisecond * 200
	movingAverage := NewMovingAverage(sampleCount, sampleInterval, progressFunc)

	average := movingAverage.GetAverageDelta()
	if average != 0 {
		t.Fatalf("average was %f. Expected 0.", average)
	}

	tests := []struct {
		increment uint64
		expected  int
	}{
		{1000, 5000},
		{2000, 7500},
		{5000, 13333},
		{5000, 16250},
		{5000, 18000},
		{5000, 19167},
		{5000, 20000},
		{5000, 20625},
		{5000, 23125},
		{5000, 25000},
		{5000, 25000},
	}

	for i, test := range tests {
		mockedProgress += test.increment
		movingAverage.TakeSample()
		average = movingAverage.GetAverageDelta()
		if int(average+0.5) != test.expected {
			t.Fatalf("average was %d. Expected %d. i = %d", int(average+0.5), test.expected, i)
		}
	}
}
