// Package flags manages the launcher's flags. For any changes, you should check if they
// are consistent with locking.Restart().
package flags

import (
	"flag"
	"fmt"
	"strconv"
)

type LauncherFlags struct {
	Uninstall        bool
	Debug            bool
	SkipSelfUpdate   bool
	NoStreamPassing  bool
	Roaming          bool
	PrintBuildTime   bool
	DeploymentConfig string

	AcceptInstall      bool
	AcceptUninstall    bool
	DismissGuiPrompts  bool
	LogIndexCounter    int
	LogInstanceCounter int

	nextLogIndex int
}

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

func Setup(args []string) (*LauncherFlags, error) {
	launcherFlags := LauncherFlags{nextLogIndex: -1}
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flagSet.BoolVar(&launcherFlags.Uninstall, UninstallFlag, false, "Remove the launcher and its bundles from the local machine.")
	flagSet.BoolVar(&launcherFlags.Debug, DebugFlag, false, "Write verbose information to the log files.")
	flagSet.BoolVar(&launcherFlags.SkipSelfUpdate, SkipSelfUpdateFlag, false, "Skip any updates to this launcher.")
	flagSet.BoolVar(&launcherFlags.NoStreamPassing, NoStreamPassingFlag, false, "Do not relay standard streams to executed commands.")
	flagSet.BoolVar(&launcherFlags.Roaming, RoamingFlag, false, "Put all files which would go under %LOCALAPPDATA% on Windows to %APPDATA% instead.")
	flagSet.BoolVar(&launcherFlags.PrintBuildTime, PrintBuildTimeFlag, false, "Print the output of 'date -u \"+%Y-%m-%d %H:%M:%S UTC\"' from the time the binary "+
		"was built to standard out and exit immediately.")
	flagSet.StringVar(&launcherFlags.DeploymentConfig, DeploymentConfigFlag, "", "Override the embedded URL of the deployment-config.")

	flagSet.BoolVar(&launcherFlags.AcceptInstall, AcceptInstallFlag, false, fmt.Sprintf("Accept install prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flagSet.BoolVar(&launcherFlags.AcceptUninstall, AcceptUninstallFlag, false, fmt.Sprintf("Accept uninstall prompt when it is dismissed. Use with -%s.", DismissGuiPromptsFlag))
	flagSet.BoolVar(&launcherFlags.DismissGuiPrompts, DismissGuiPromptsFlag, false, "Automatically dismiss GUI prompts.")
	flagSet.IntVar(&launcherFlags.LogIndexCounter, LogIndexCounterFlag, -1, "Number to increment when restarting.")
	flagSet.IntVar(&launcherFlags.LogInstanceCounter, LogInstanceCounterFlag, 0, "Number to increment when started by user.")
	setDeprecatedFlags(flagSet)

	err := flagSet.Parse(args[1:])
	if err != nil {
		return &launcherFlags, withSuggestions(err, flagSet, []string{DebugFlag, RoamingFlag, SkipSelfUpdateFlag, UninstallFlag})
	}

	if !launcherFlags.DismissGuiPrompts && launcherFlags.AcceptInstall {
		return &launcherFlags, fmt.Errorf("-%s was set when -%s was not", AcceptInstallFlag, DismissGuiPromptsFlag)
	}

	if !launcherFlags.DismissGuiPrompts && launcherFlags.AcceptUninstall {
		return &launcherFlags, fmt.Errorf("-%s was set when -%s was not", AcceptUninstallFlag, DismissGuiPromptsFlag)
	}

	return &launcherFlags, nil
}

func withSuggestions(err error, flagSet *flag.FlagSet, suggestFlags []string) error {
	if len(suggestFlags) == 0 {
		return err
	}
	suggestionText := ". Maybe you meant one of these:"
	for _, suggestFlag := range suggestFlags {
		if f := flagSet.Lookup(suggestFlag); f != nil {
			suggestionText += "\n" + "\xc2\xa0\xc2\xa0\xc2\xa0\xc2\xa0-" + f.Name + ": " + f.Usage
		} else {
			suggestionText += "\n" + "\xc2\xa0\xc2\xa0\xc2\xa0\xc2\xa0-" + f.Name
		}
	}
	return fmt.Errorf("%w%s", err, suggestionText)
}

// GetTransmittingFlags returns those flags which the launcher should hand to itself when restarting.
func (launcherFlags *LauncherFlags) GetTransmittingFlags() (transmittingFlags []string) {
	transmittingFlags = append(transmittingFlags, "-"+LogIndexCounterFlag, strconv.Itoa(launcherFlags.nextLogIndex))
	transmittingFlags = append(transmittingFlags, "-"+LogInstanceCounterFlag, strconv.Itoa(launcherFlags.LogInstanceCounter+1))
	if launcherFlags.Uninstall {
		transmittingFlags = append(transmittingFlags, "-"+UninstallFlag)
	}
	if launcherFlags.Debug {
		transmittingFlags = append(transmittingFlags, "-"+DebugFlag)
	}
	if launcherFlags.SkipSelfUpdate {
		transmittingFlags = append(transmittingFlags, "-"+SkipSelfUpdateFlag)
	}
	if launcherFlags.Roaming {
		transmittingFlags = append(transmittingFlags, "-"+RoamingFlag)
	}
	if launcherFlags.DeploymentConfig != "" {
		transmittingFlags = append(transmittingFlags, "-"+DeploymentConfigFlag, launcherFlags.DeploymentConfig)
	}
	if launcherFlags.AcceptInstall {
		transmittingFlags = append(transmittingFlags, "-"+AcceptInstallFlag)
	}
	if launcherFlags.AcceptUninstall {
		transmittingFlags = append(transmittingFlags, "-"+AcceptUninstallFlag)
	}
	if launcherFlags.DismissGuiPrompts {
		transmittingFlags = append(transmittingFlags, "-"+DismissGuiPromptsFlag)
	}
	if launcherFlags.NoStreamPassing {
		transmittingFlags = append(transmittingFlags, "-"+NoStreamPassingFlag)
	}

	return transmittingFlags
}

func (launcherFlags *LauncherFlags) SetNextLogIndex(index int) {
	launcherFlags.nextLogIndex = index
}

func setDeprecatedFlags(flagSet *flag.FlagSet) {
	flagSet.String("remove", "", "DEPRECATED: Name of binary to remove upon launch.")
}
