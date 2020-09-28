package stats

import "time"

type MovingAverage struct {
	totals         []uint64
	maxTotals      int
	sampleInterval time.Duration
	sampleFunc     func() uint64
}

func NewMovingAverage(maxSampleCount int, sampleInterval time.Duration, sampleFunc func() uint64) *MovingAverage {
	return &MovingAverage{maxTotals: maxSampleCount + 1, sampleInterval: sampleInterval, sampleFunc: sampleFunc}
}

func (ma *MovingAverage) TakeSample() {
	currentTotal := ma.sampleFunc()
	if len(ma.totals) >= ma.maxTotals {
		ma.totals = append(ma.totals[1:], currentTotal)
	} else {
		ma.totals = append(ma.totals, currentTotal)
	}
}

func (ma *MovingAverage) GetAverageDelta() float64 {
	if len(ma.totals) == 0 {
		return 0
	}
	var availableDeltaCount int
	var samplesDelta float64

	if len(ma.totals) == 1 {
		availableDeltaCount = 1
		samplesDelta = float64(ma.totals[0])
	} else if len(ma.totals) < ma.maxTotals {
		availableDeltaCount = len(ma.totals)
		samplesDelta = float64(ma.totals[len(ma.totals)-1])
	} else {
		availableDeltaCount = len(ma.totals) - 1
		samplesDelta = float64(ma.totals[len(ma.totals)-1] - ma.totals[0])
	}

	samplesDeltaDuration := ma.sampleInterval.Seconds() * float64(availableDeltaCount)
	changePerSecond := samplesDelta / samplesDeltaDuration
	return changePerSecond
}

func (ma *MovingAverage) Total() uint64 {
	if ma.totals == nil {
		return 0
	}
	return ma.totals[len(ma.totals)-1]
}

func (ma *MovingAverage) Reset() {
	ma.totals = nil
}
