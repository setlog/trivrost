package locking

import (
	"context"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/system"
)

const launcherSignatureFileName = ".launcher-lock"

type LockActionResult int

const (
	LockUnavailable LockActionResult = iota // Lock is claimed by another running process.
	LockClaimed                             // The current process claimed the lock, and now owns it.
	LockOwned                               // The current process owns the lock already.
)

// Blocks until the lock is available, claims it and restarts.
func AcquireLock(ctx context.Context) {
	gui.SetStage(gui.StageAcquireLock, 0)
	for {
		switch result := setSignature(system.GetCurrentProcessSignature()); result {
		case LockOwned:
			log.Info("Owning the Launcher Lock.")
			return
		case LockClaimed:
			Restart(true)
		case LockUnavailable:
			misc.MustWaitForContext(ctx, time.Millisecond*300)
		}
	}
}

func setSignature(processSignature *system.ProcessSignature) LockActionResult {
	if !mustTryLock() {
		return LockUnavailable
	}
	if isSignatureSet(processSignature) {
		return LockOwned
	}

	// Avoid race condition when the launcher restarts: between releasing ".lock" and the new instance claiming it.
	if isProcessWithSetSignatureRunning() {
		ReleaseLock()
		time.Sleep(time.Millisecond * 100) // Give other instance an opportunity to acquire ".lock".
		return LockUnavailable
	}

	mustWriteProcessSignatureListFile(launcherSignatureFilePath(), []system.ProcessSignature{*processSignature})
	log.Info("Claimed Launcher Lock.")
	return LockClaimed
}

func isSignatureSet(processSignature *system.ProcessSignature) bool {
	sigs := readProcessSignatureListFile(launcherSignatureFilePath())
	return len(sigs) != 0 && (sigs[0] == *processSignature)
}

func isProcessWithSetSignatureRunning() bool {
	sigs := readProcessSignatureListFile(launcherSignatureFilePath())
	return len(sigs) != 0 && system.IsProcessSignatureRunning(&sigs[0])
}

func launcherSignatureFilePath() string {
	return filepath.Join(places.GetAppLocalDataFolderPath(), launcherSignatureFileName)
}
