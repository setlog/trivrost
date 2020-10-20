package stats

import "time"

// MovingAverage models a series of arbitrary uint64 values placed equidistantly in the time domain.
type MovingAverage struct {
	totals         []uint64
	maxTotals      int
	sampleInterval time.Duration // TODO: Don't terminologically restrict the moving average to the time domain.
	sampleFunc     func() uint64
}

// NewMovingAverage constructs a new MovingAverage with given sample count limit, given assumed sample interval
// and given sample function.
func NewMovingAverage(maxSampleCount int, sampleInterval time.Duration, sampleFunc func() uint64) *MovingAverage {
	return &MovingAverage{maxTotals: maxSampleCount + 1, sampleInterval: sampleInterval, sampleFunc: sampleFunc}
}

// TakeSample calls the sample function the MovingAverage was constructed with using NewMovingAverage
// and stores the returned value internally, discarding any old values
func (ma *MovingAverage) TakeSample() {
	currentTotal := ma.sampleFunc()
	if len(ma.totals) >= ma.maxTotals {
		ma.totals = append(ma.totals[1:], currentTotal)
	} else {
		ma.totals = append(ma.totals, currentTotal)
	}
}

// AveragePerSecondDelta returns the average per-second change from sample to sample using all available samples.
func (ma *MovingAverage) AveragePerSecondDelta() float64 {
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

// TODO: Principally, something like AverageTotal() would be nice to have as well, but isn't needed for trivrost.

// Total returns the most recent value sampled by TakeSample(), or 0 if it has not been called yet.
func (ma *MovingAverage) Total() uint64 {
	if ma.totals == nil {
		return 0
	}
	return ma.totals[len(ma.totals)-1]
}

// Reset returns this moving average to its initial state, as if it had just been returned from NewMovingAverage().
func (ma *MovingAverage) Reset() {
	ma.totals = nil
}
