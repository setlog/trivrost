package locking

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/pkg/launcher/hashing"
	"github.com/setlog/trivrost/pkg/system"
)

func Restart(forwardLauncherLockOwnership bool, launcherFlags *flags.LauncherFlags) {
	RestartWithBinary(forwardLauncherLockOwnership, system.GetBinaryPath(), launcherFlags)
}

// Starts a new instance of the calling executable, writes the new process signature into the launcher signature file and quits the current instance.
func RestartWithBinary(forwardLauncherLockOwnership bool, binaryPath string, launcherFlags *flags.LauncherFlags) {
	hash, _, hashErr := hashing.CalculateSha256(binaryPath)
	log.WithFields(log.Fields{"forwardLauncherLockOwnership": forwardLauncherLockOwnership, "binaryPath": binaryPath, "hash": hash, "hashErr": hashErr}).Info("Restarting.")
	absoluteBinaryPath := system.MustGetAbsolutePath(binaryPath)
	workingDirectory := filepath.Dir(absoluteBinaryPath)
	_, procSig, err := system.StartProcess(absoluteBinaryPath, workingDirectory, launcherFlags.GetTransmittingFlags(), nil, true)
	if err != nil {
		panic(err)
	}
	if forwardLauncherLockOwnership {
		mustWriteProcessSignatureListFile(launcherSignatureFilePath(), []system.ProcessSignature{*procSig})
	}
	ReleaseLock()
	log.Info("Restart appears to have worked. Exiting now.")
	log.Exit(0)
}
