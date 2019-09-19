package locking

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/system"
)

const applicationSignaturesFileName = ".execution-lock"

// Adds the supplied ProcessSignature to the execution lock file.
func AddApplicationSignature(sig *system.ProcessSignature) {
	sigs := readApplicationSignatures()
	sigs = append(sigs, *sig)
	mustSetApplicationSignatures(sigs)
}

// Waits for the processes from all previous executions to stop running and then removes the execution lock file.
func AwaitApplicationsTerminated(ctx context.Context) {
	sigs := readApplicationSignatures()
	gui.SetStage(gui.StageAwaitApplicationsTerminated, 0)
	for _, sig := range sigs {
		waitSig := sig
		WaitForProcessSignatureToStopRunning(ctx, &waitSig)
	}
	mustRemoveApplicationSignatures()
}

// MinimizeApplicationSignaturesList removes any process signatures no longer running from the execution lock file if it exists.
// If no process signature is running, the execution lock file is removed and the function returns true.
func MinimizeApplicationSignaturesList() bool {
	sigs := readApplicationSignatures()
	runningSigs := make([]system.ProcessSignature, 0)
	for _, sig := range sigs {
		checkSig := sig
		if system.IsProcessSignatureRunning(&checkSig) {
			runningSigs = append(runningSigs, checkSig)
		}
	}
	return !mustSetApplicationSignatures(runningSigs)
}

func readApplicationSignatures() (sigs []system.ProcessSignature) {
	return readProcessSignatureListFile(applicationSignaturesFilePath())
}

func mustSetApplicationSignatures(sigs []system.ProcessSignature) bool {
	if len(sigs) > 0 {
		mustWriteProcessSignatureListFile(applicationSignaturesFilePath(), sigs)
		return true
	}
	mustRemoveApplicationSignatures()
	return false
}

func mustRemoveApplicationSignatures() {
	err := os.Remove(applicationSignaturesFilePath())
	if err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf(`Could not remove Execution Lock "%s": %v`, applicationSignaturesFilePath(), err))
	}
}

func applicationSignaturesFilePath() string {
	return filepath.Join(places.GetAppLocalDataFolderPath(), applicationSignaturesFileName)
}
