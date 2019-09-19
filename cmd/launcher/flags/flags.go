// Package flags manages the launcher's flags. For any changes, you should check if they
// are consistent with locking.Restart().
package flags

import (
	"flag"
	"fmt"
	"strconv"
)

var (
	Uninstall        = new(bool)
	Debug            = new(bool)
	SkipSelfUpdate   = new(bool)
	NoStreamPassing  = new(bool)
	Roaming          = new(bool)
	PrintBuildTime   = new(bool)
	DeploymentConfig = new(string)

	AcceptInstall      = new(bool)
	AcceptUninstall    = new(bool)
	DismissGuiPrompts  = new(bool)
	LogIndexCounter    = new(int)
	LogInstanceCounter = new(int)
)

var nextLogIndex = -1

const (
	UninstallFlag        = "uninstall"
	DebugFlag            = "debug"
	SkipSelfUpdateFlag   = "skipselfupdate"
	NoStreamPassingFlag  = "nostreampassing"
	RoamingFlag          = "roaming"
	PrintBuildTimeFlag   = "build-time"
	DeploymentConfigFlag = "deployment-config"

	AcceptInstallFlag      = "accept-install"
	AcceptUninstallFlag    = "accept-uninstall"
	DismissGuiPromptsFlag  = "dismiss-gui-prompts"
	LogIndexCounterFlag    = "log-index"
	LogInstanceCounterFlag = "log-instance"
)

func Setup() {
	flag.BoolVar(Uninstall, UninstallFlag, false, "Flag to uninstall the launcher and its bundles on the local machine.")
	flag.BoolVar(Debug, DebugFlag, false, "Enable debug log level.")
	flag.BoolVar(SkipSelfUpdate, SkipSelfUpdateFlag, false, "Never perform a self-update.")
	flag.BoolVar(NoStreamPassing, NoStreamPassingFlag, false, "Do not relay standard streams to executed commands.")
	flag.BoolVar(Roaming, RoamingFlag, false, "Put all files which would go under %LOCALAPPDATA% on Windows to %APPDATA% instead.")
	flag.BoolVar(PrintBuildTime, PrintBuildTimeFlag, false, "Print the output of 'date -u \"+%Y-%m-%d %H:%M:%S UTC\"' from the time the binary "+
		"was built to standard out and exit immediately.")
	flag.StringVar(DeploymentConfig, DeploymentConfigFlag, "", "Override the embedded URL of the deployment-config.")

	flag.BoolVar(AcceptInstall, AcceptInstallFlag, false, fmt.Sprintf("Accept install prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flag.BoolVar(AcceptUninstall, AcceptUninstallFlag, false, fmt.Sprintf("Accept uninstall prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flag.BoolVar(DismissGuiPrompts, DismissGuiPromptsFlag, false, "Automatically dismiss GUI prompts.")
	flag.IntVar(LogIndexCounter, LogIndexCounterFlag, -1, "Number to increment when restarting.")
	flag.IntVar(LogInstanceCounter, LogInstanceCounterFlag, 0, "Number to increment when started by user.")

	setDeprecatedFlags()

	flag.Parse()

	if !*DismissGuiPrompts && *AcceptInstall {
		*AcceptInstall = false
		panic(fmt.Sprintf("-%s was set when -%s was not.", AcceptInstallFlag, DismissGuiPromptsFlag))
	}

	if !*DismissGuiPrompts && *AcceptUninstall {
		*AcceptUninstall = false
		panic(fmt.Sprintf("-%s was set when -%s was not.", AcceptUninstallFlag, DismissGuiPromptsFlag))
	}
}

// GetTransmittingFlags returns those flags which the launcher should hand to itself when restarting.
func GetTransmittingFlags() (transmittingFlags []string) {
	if *Uninstall {
		transmittingFlags = append(transmittingFlags, "-"+UninstallFlag)
	}
	if *Debug {
		transmittingFlags = append(transmittingFlags, "-"+DebugFlag)
	}
	if *SkipSelfUpdate {
		transmittingFlags = append(transmittingFlags, "-"+SkipSelfUpdateFlag)
	}
	if *Roaming {
		transmittingFlags = append(transmittingFlags, "-"+RoamingFlag)
	}
	if *DeploymentConfig != "" {
		transmittingFlags = append(transmittingFlags, "-"+DeploymentConfigFlag, *DeploymentConfig)
	}
	if *AcceptInstall {
		transmittingFlags = append(transmittingFlags, "-"+AcceptInstallFlag)
	}
	if *AcceptUninstall {
		transmittingFlags = append(transmittingFlags, "-"+AcceptUninstallFlag)
	}
	if *DismissGuiPrompts {
		transmittingFlags = append(transmittingFlags, "-"+DismissGuiPromptsFlag)
	}
	if *NoStreamPassing {
		transmittingFlags = append(transmittingFlags, "-"+NoStreamPassingFlag)
	}
	transmittingFlags = append(transmittingFlags, "-"+LogIndexCounterFlag, strconv.Itoa(nextLogIndex))
	transmittingFlags = append(transmittingFlags, "-"+LogInstanceCounterFlag, strconv.Itoa(*LogInstanceCounter+1))

	return transmittingFlags
}

func SetNextLogIndex(index int) {
	nextLogIndex = index
}

func setDeprecatedFlags() {
	flag.String("remove", "", "DEPRECATED: Name of binary to remove upon launch.")
}
