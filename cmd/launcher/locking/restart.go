package locking

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/pkg/launcher/hashing"
	"github.com/setlog/trivrost/pkg/system"
)

func Restart(forwardLauncherLockOwnership bool) {
	RestartWithBinary(forwardLauncherLockOwnership, system.GetBinaryPath())
}

// Starts a new instance of the calling executable, writes the new process signature into the launcher signature file and quits the current instance.
func RestartWithBinary(forwardLauncherLockOwnership bool, binaryPath string) {
	hash, _, hashErr := hashing.CalculateSha256(binaryPath)
	log.WithFields(log.Fields{"forwardLauncherLockOwnership": forwardLauncherLockOwnership, "binaryPath": binaryPath, "hash": hash, "hashErr": hashErr}).Info("Restarting.")
	absoluteBinaryPath := system.MustGetAbsolutePath(binaryPath)
	workingDirectory := filepath.Dir(absoluteBinaryPath)
	_, procSig := system.MustStartProcess(absoluteBinaryPath, workingDirectory, flags.GetTransmittingFlags(), nil, true)
	if forwardLauncherLockOwnership {
		mustWriteProcessSignatureListFile(launcherSignatureFilePath(), []system.ProcessSignature{*procSig})
	}
	ReleaseLock()
	log.Info("Restart appears to have worked. Exiting now.")
	log.Exit(0)
}
