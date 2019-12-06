package launcher

import (
	"time"
)

var lingerTimeMilliseconds int

func Linger() {
	time.Sleep(time.Duration(lingerTimeMilliseconds) * time.Millisecond)
}
