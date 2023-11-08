package launcher

import (
	"github.com/setlog/systemuri"
	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/resources"
	"github.com/setlog/trivrost/pkg/launcher/config"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

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

// IsInstanceInstalledInSystemMode returns true if we are in system mode.
func IsInstanceInstalledInSystemMode() bool {
	return system.FolderExists(places.GetSystemWideBundleFolderPath())
}

// IsInstanceInstalledForCurrentUser returns true if the launcher's desired path under user files is occupied by the program running this code.
func IsInstanceInstalledForCurrentUser() bool {
	programPath := system.GetProgramPath()
	targetProgramPath := getTargetProgramPath()
	if programPath == targetProgramPath {
		return true
	}
	if runtime.GOOS == "windows" {
		if strings.EqualFold(programPath, targetProgramPath) || strings.EqualFold(filepath.Clean(programPath), filepath.Clean(targetProgramPath)) {
			return true
		}
	}
	return false
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

	// Registering the scheme handlers had to happen in run.go, since we do not have the deployment config here, yet.

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

// RegisterSchemeHandlers registers the schemes defined in the deployment config
func RegisterSchemeHandlers(launcherFlags *flags.LauncherFlags, schemeHandlers []config.SchemeHandler) {
	transmittingFlags := launcherFlags.GetTransmittingFlags()
	for _, schemeHandler := range schemeHandlers {
		// TODO: Create flag "-lockagnostic" (or such) to allow to eliminate all self-restarting behavior (when used in combination with "-skipselfupdate") to reduce UI flickering?
		// TODO: Create and then always add flag "-skipschemehandlerregistry" (or such) here to prevent -extra-env from being added?
		// TODO: systemuri does not presently implement %%-escapes according to deployment-config.md.
		binaryPath := system.GetBinaryPath()
		// We want to pass a few flags from the current execution (transmitting flags) as well, but only if they are not set in the passed arguments.
		finalArguments := []string{}
		finalArguments = removeFromList(transmittingFlags, extractArguments(schemeHandler.Arguments))
		finalArguments = filterWhitelistArguments(finalArguments)

		transmittingFlagsFiltered := removeFromList(transmittingFlags, finalArguments)
		transmittingFlagsFiltered = []string{} // TODO: We actually want to create a whitelist here of flags that are okay to retain

		arguments := strings.Join(transmittingFlagsFiltered, " ") + schemeHandler.Arguments
		err := systemuri.RegisterURLHandler(resources.LauncherConfig.BrandingName, schemeHandler.Scheme, binaryPath, arguments)
		if err != nil {
			log.Warnf("Registering the scheme \"%s\" failed: %v", schemeHandler.Scheme, err)
		}
	}
}

// TODO: Make case insensitive
func filterWhitelistArguments(arguments []string) []string {
	// slice of entries to allow in the arguments list
	whitelist := []string{
		"debug",
		"skipselfupdate",
		"roaming",
		"deployment-config",
		"extra-env",
	}

	var filteredArguments []string

	sort.Strings(whitelist)
	for _, argument := range arguments {
		if found := sort.SearchStrings(whitelist, argument); found < len(whitelist) && whitelist[found] == argument {
			continue
		}
		filteredArguments = append(filteredArguments, argument)
	}

	return filteredArguments
}

// TODO: Make case insensitive
func removeFromList(sourceList []string, itemsToRemove []string) []string {
	argSet := make(map[string]bool)
	for _, item := range itemsToRemove {
		argSet[item] = true
	}

	var filteredSourceList []string
	for _, item := range sourceList {
		if _, ok := argSet[item]; !ok {
			filteredSourceList = append(filteredSourceList, item)
		}
	}
	return filteredSourceList
}

func extractArguments(input string) []string {
	words := strings.Fields(input)
	var args []string
	for _, w := range words {
		if strings.HasPrefix(w, "-") || strings.HasPrefix(w, "--") {
			arg := strings.SplitN(w, "=", 2)[0]
			args = append(args, arg)
		}
	}
	return args
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
