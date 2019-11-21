package launcher

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/resources"

	"github.com/setlog/trivrost/cmd/launcher/places"

	log "github.com/sirupsen/logrus"

	"github.com/andlabs/ui"
	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/pkg/system"
)

var buildTime string

func BuildTime() string {
	return strings.Trim(buildTime, " \n\t\r")
}

// HasInstallation returns true iff the launcher's desired path under user files is occupied by any file.
func HasInstallation() bool {
	return system.FileExists(getTargetBinaryPath())
}

// IsInstanceInstalled returns true iff the binary running this code is to be considered installed.
func IsInstanceInstalled() bool {
	isInstalled := IsInstanceInstalledInSystemMode() || IsInstanceInstalledForCurrentUser()
	if isInstalled {
		log.Debugf(`Launcher is installed. Application path "%s" matches with target application path.`, system.GetProgramPath())
	} else {
		log.Debugf(`Launcher is not installed. Application path "%s" does not match target application path "%s", nor is there`+
			`a 'systembundles'-folder next to the binary.`, system.GetProgramPath(), getTargetProgramPath())
	}
	return isInstalled
}

// IsInstanceInstalledInSystemMode returns true iff we are in system mode.
func IsInstanceInstalledInSystemMode() bool {
	return system.FolderExists(places.GetSystemWideBundleFolderPath())
}

// IsInstanceInstalledForCurrentUser returns true iff the launcher's desired path under user files is occupied by the program running this code.
func IsInstanceInstalledForCurrentUser() bool {
	return system.GetProgramPath() == getTargetProgramPath()
}

// IsInstallationOutdated returns true if the time the installed launcher binary was built
// predates the time the currently running launcher binary was built.
func IsInstallationOutdated() bool {
	cmd := exec.Command(getTargetBinaryPath(), "-"+flags.PrintBuildTimeFlag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Warnf("Asking installed binary for build time failed: %v. This is expected if it is pre-v1.3.0.", err)
		return true
	}
	installationBuildTime := strings.Trim(string(output), " \n\t\r")

	const buildTimeFormat = "2006-01-02 15:04:05 MST"
	timeOfRunning, err := time.Parse(buildTimeFormat, BuildTime())
	if err != nil {
		log.Errorf("Parsing build time of running binary failed: %v", err)
		return false
	}
	timeOfInstalled, err := time.Parse(buildTimeFormat, installationBuildTime)
	if err != nil {
		log.Errorf("Parsing build time of installed binary failed: %v", err)
		return true
	}

	isInstallationOutdated := timeOfInstalled.Before(timeOfRunning)
	if isInstallationOutdated {
		log.Infof("Installation Build Time \"%s\" predates Instance Build Time \"%s\".", installationBuildTime, BuildTime())
	} else {
		log.Infof("Installation Build Time \"%s\" does not predate Instance Build Time \"%s\".", installationBuildTime, BuildTime())
	}
	return isInstallationOutdated
}

func Install(launcherFlags *flags.LauncherFlags) {
	deletePlainFiles()

	programPath, targetProgramPath := system.GetProgramPath(), getTargetProgramPath()
	if filepath.Dir(programPath) == filepath.Dir(targetProgramPath) {
		log.Infof("Renaming running program at \"%s\" to \"%s\" as per embedded launcher-config.", programPath, filepath.Base(targetProgramPath))
		system.MustMoveAll(programPath, targetProgramPath)
	} else {
		log.Infof("Copying running program at \"%s\" to \"%s\".", programPath, targetProgramPath)
		system.MustCopyAll(programPath, targetProgramPath)
	}

	runPostBinaryUpdateProvisioning()

	InstallShortcuts(targetProgramPath, launcherFlags)

	MustRestartWithInstalledBinary(launcherFlags)
}

func InstallShortcuts(targetProgramPath string, launcherFlags *flags.LauncherFlags) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	ui.QueueMain(func() { // Not UI functionality, but required to run on the main thread to be reliable on all OSes.
		installShortcuts(targetProgramPath, launcherFlags)
		waitGroup.Done()
	})
	waitGroup.Wait()
}

func MustRestartWithInstalledBinary(launcherFlags *flags.LauncherFlags) {
	locking.RestartWithBinary(true, getTargetBinaryPath(), launcherFlags)
}

func RestartWithInstalledBinary(launcherFlags *flags.LauncherFlags) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Restarting with binary \"%s\" failed: %v", getTargetBinaryPath(), r)
		}
	}()
	MustRestartWithInstalledBinary(launcherFlags)
	return true
}

func installShortcuts(targetPath string, launcherFlags *flags.LauncherFlags) {
	log.Info("Installing launcher shortcuts.")
	createLaunchDesktopShortcut(targetPath, launcherFlags)
	createLaunchStartMenuShortcut(targetPath, launcherFlags)
	createUninstallStartMenuShortcut(targetPath, launcherFlags)
}

func getTargetProgramPath() string {
	return filepath.Join(places.GetLauncherTargetDirectoryPath(), getTargetProgramName())
}

func ReportProgramNameDivergence() {
	if filepath.Base(system.GetProgramPath()) != getTargetProgramName() {
		log.Warnf("Program name on disk (\"%s\") has diverged from configured program name (\"%s\").",
			filepath.Base(system.GetProgramPath()), getTargetProgramName())
	}
}

func getTargetProgramName() string {
	if runtime.GOOS == system.OsWindows {
		return resources.LauncherConfig.BinaryName + ".exe"
	} else if runtime.GOOS == system.OsMac {
		return resources.LauncherConfig.BinaryName + ".app"
	}
	return resources.LauncherConfig.BinaryName
}

func getTargetBinaryPath() string {
	targetProgramPath := getTargetProgramPath()
	if runtime.GOOS == system.OsMac {
		return filepath.Join(targetProgramPath, "Contents", "MacOS", "launcher")
	}
	return targetProgramPath
}
