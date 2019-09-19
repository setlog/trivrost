package misc

import (
	"context"
	"time"
)

func MustWaitForContext(ctx context.Context, timeout time.Duration) {
	if err := WaitForContext(ctx, timeout); err != nil {
		panic(err)
	}
}

func WaitForContext(ctx context.Context, timeout time.Duration) error {
	if WaitForSignal(ctx.Done(), timeout) {
		return ctx.Err()
	}
	return nil
}

func WaitForSignal(c <-chan struct{}, timeout time.Duration) bool {
	select {
	case <-c:
		return true
	case <-time.NewTimer(timeout).C:
		return false
	}
}
