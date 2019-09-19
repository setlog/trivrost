package locking

import (
	"context"
	"time"

	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/system"
)

func WaitForProcessSignatureToStopRunning(ctx context.Context, procSig *system.ProcessSignature) {
	for system.IsProcessSignatureRunning(procSig) {
		misc.MustWaitForContext(ctx, time.Millisecond*300)
	}
}
