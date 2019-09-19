package misc

import (
	"time"
)

func SleepMilliseconds(timeInMilliseconds int) {
	if timeInMilliseconds > 0 {
		time.Sleep(time.Millisecond * time.Duration(timeInMilliseconds))
	}
}
