// Package flags manages the launcher's flags. For any changes, you should check if they
// are consistent with locking.Restart().
package flags

import (
	"flag"
	"fmt"
	"os"
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

func Setup() error {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.BoolVar(Uninstall, UninstallFlag, false, "Flag to uninstall the launcher and its bundles on the local machine.")
	flagSet.BoolVar(Debug, DebugFlag, false, "Enable debug log level.")
	flagSet.BoolVar(SkipSelfUpdate, SkipSelfUpdateFlag, false, "Never perform a self-update.")
	flagSet.BoolVar(NoStreamPassing, NoStreamPassingFlag, false, "Do not relay standard streams to executed commands.")
	flagSet.BoolVar(Roaming, RoamingFlag, false, "Put all files which would go under %LOCALAPPDATA% on Windows to %APPDATA% instead.")
	flagSet.BoolVar(PrintBuildTime, PrintBuildTimeFlag, false, "Print the output of 'date -u \"+%Y-%m-%d %H:%M:%S UTC\"' from the time the binary "+
		"was built to standard out and exit immediately.")
	flagSet.StringVar(DeploymentConfig, DeploymentConfigFlag, "", "Override the embedded URL of the deployment-config.")

	flagSet.BoolVar(AcceptInstall, AcceptInstallFlag, false, fmt.Sprintf("Accept install prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flagSet.BoolVar(AcceptUninstall, AcceptUninstallFlag, false, fmt.Sprintf("Accept uninstall prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flagSet.BoolVar(DismissGuiPrompts, DismissGuiPromptsFlag, false, "Automatically dismiss GUI prompts.")
	flagSet.IntVar(LogIndexCounter, LogIndexCounterFlag, -1, "Number to increment when restarting.")
	flagSet.IntVar(LogInstanceCounter, LogInstanceCounterFlag, 0, "Number to increment when started by user.")

	setDeprecatedFlags(flagSet)

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if !*DismissGuiPrompts && *AcceptInstall {
		return fmt.Errorf("-%s was set when -%s was not.", AcceptInstallFlag, DismissGuiPromptsFlag)
	}

	if !*DismissGuiPrompts && *AcceptUninstall {
		return fmt.Errorf("-%s was set when -%s was not.", AcceptUninstallFlag, DismissGuiPromptsFlag)
	}

	return nil
}

// GetTransmittingFlags returns those flags which the launcher should hand to itself when restarting.
func GetTransmittingFlags() (transmittingFlags []string) {
	transmittingFlags = append(transmittingFlags, "-"+LogIndexCounterFlag, strconv.Itoa(nextLogIndex))
	transmittingFlags = append(transmittingFlags, "-"+LogInstanceCounterFlag, strconv.Itoa(*LogInstanceCounter+1))
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

	return transmittingFlags
}

func SetNextLogIndex(index int) {
	nextLogIndex = index
}

func setDeprecatedFlags(flagSet *flag.FlagSet) {
	flagSet.String("remove", "", "DEPRECATED: Name of binary to remove upon launch.")
}
