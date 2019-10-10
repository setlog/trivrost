package misc

import (
	"log"
	"runtime/debug"
)

// LogPanic temporarily recovers from a panic (if there is one) and then logs it with Panicf() of Go's log library
// so that a writer attached with log.SetOutput() can deal with it (e.g. log it to a file), then lets the panic resume.
func LogPanic() {
	if r := recover(); r != nil {
		LogRecoveredValue(r)
	}
}

func LogRecoveredValue(r interface{}) {
	// The stack printed when panic() is not recover()ed bypasses file-logging, so log it explicitly here.
	log.Panicf("Unrecoverable state: %v\n%v", r, TryRemoveLines(string(debug.Stack()), 1, 3))
}
