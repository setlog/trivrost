package launcher

import (
	"math"
	"time"
)

type progressFaker struct {
	progressStartTime time.Time
	progressRate      float64
}

func newProgressFaker(progressRate float64) *progressFaker {
	return &progressFaker{progressRate: progressRate}
}

func (pf *progressFaker) getProgress() uint64 {
	t := time.Time{}
	if pf.progressStartTime == t {
		pf.progressStartTime = time.Now()
		return 0
	}
	diff := time.Since(pf.progressStartTime)
	return uint64(math.Round(diff.Seconds() * pf.progressRate))
}
