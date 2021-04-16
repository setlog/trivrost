package logging

import (
	"fmt"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

type LogLimiter struct {
	limit        int64
	seenMessages *sync.Map
}

func NewLogLimiter(limit int) *LogLimiter {
	return &LogLimiter{limit: int64(limit), seenMessages: &sync.Map{}}
}

func (l *LogLimiter) Logf(level log.Level, format string, args ...interface{}) {
	l.Log(level, fmt.Sprintf(format, args...))
}

func (l *LogLimiter) Log(level log.Level, message string) {
	var seenCount int64 = 0
	seenCountP, _ := l.seenMessages.LoadOrStore(message, &seenCount)
	seenCount = atomic.AddInt64(seenCountP.(*int64), 1)
	if seenCount <= l.limit {
		if seenCount == l.limit {
			log.StandardLogger().Logf(level, "%v (Last log of this message due to excess)", message)
		} else {
			log.StandardLogger().Log(level, message)
		}
	}
}

func (l *LogLimiter) Infof(format string, args ...interface{}) {
	l.Logf(log.InfoLevel, format, args...)
}

func (l *LogLimiter) Info(message string) {
	l.Log(log.InfoLevel, message)
}

func (l *LogLimiter) Warnf(format string, args ...interface{}) {
	l.Logf(log.WarnLevel, format, args...)
}

func (l *LogLimiter) Warn(message string) {
	l.Log(log.WarnLevel, message)
}

func (l *LogLimiter) Errorf(format string, args ...interface{}) {
	l.Logf(log.ErrorLevel, format, args...)
}

func (l *LogLimiter) Error(message string) {
	l.Log(log.ErrorLevel, message)
}
