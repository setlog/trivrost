package launcher

import (
	"github.com/setlog/trivrost/pkg/misc"
)

var lingerTimeMilliseconds int

func Linger() {
	misc.SleepMilliseconds(lingerTimeMilliseconds)
}
