package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/setlog/trivrost/pkg/misc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

func determineProgramPath() (string, error) {
	return binaryPath, nil
}

// Returns the path to the renamed binary, or its original path if it could not be renamed.
func undeployProgram(undeployProgramPath string) (string, error) {
	binaryOldPath := undeployProgramPath
	randomHex := misc.MustGetRandomHexString(8)
	binaryNewName := "~" + filepath.Base(binaryOldPath) + ".delete." + randomHex
	binaryNewPath := filepath.Join(filepath.Dir(binaryOldPath), binaryNewName)

	if err := os.Rename(binaryOldPath, binaryNewPath); err != nil {
		return binaryOldPath, fmt.Errorf("Failed renaming binary \"%s\" to \"%s\": %v", binaryOldPath, binaryNewPath, err)
	}
	return binaryNewPath, nil
}

func deleteProgram(deleteAppPath string) error {
	return delayDeleteFile(deleteAppPath, 3)
}

func delayDeleteFile(fileToDelete string, waitSeconds int) error {
	log.WithFields(log.Fields{"fileToDelete": fileToDelete}).Infof("Try to remove file in %d seconds.", waitSeconds)

	files := make([]*os.File, 3)
	waitSecondsArg := strconv.Itoa(waitSeconds + 1) // Ping waits on intervals, not pings, so add 1.
	runArgs := []string{"cmd.exe", "/C", "ping", "127.0.0.1", "-n", waitSecondsArg, "&", "del", fileToDelete}

	pathOfCmd, err := exec.LookPath("cmd.exe")
	if err != nil {
		return fmt.Errorf("Could not find cmd.exe: %v", err)
	}

	procAttr := &windows.SysProcAttr{}
	procAttr.HideWindow = true
	proc, err := os.StartProcess("cmd.exe", runArgs, &os.ProcAttr{
		Dir:   filepath.Dir(pathOfCmd),
		Env:   os.Environ(),
		Files: files,
		Sys:   procAttr,
	})
	if err != nil {
		return fmt.Errorf("Could not start process: %v", err)
	}

	err = proc.Release()
	if err != nil {
		return fmt.Errorf("Could not release process: %v", err)
	}
	return nil
}
