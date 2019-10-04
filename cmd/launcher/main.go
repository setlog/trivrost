//go:generate -command metawriter go run ../metawriter/main.go
//go:generate metawriter resources/launcher-config.json resources/versioninfo.json.template resources/versioninfo.json resources/main.exe.template.manifest resources/main.exe.manifest resources/Info.template.plist resources/Info.plist

//go:generate goversioninfo -manifest=resources/main.exe.manifest -platform-specific=true resources/versioninfo.json

//go:generate -command asset go run asset.go
//go:generate asset -var=LauncherConfig -wrap=readLauncherConfigAsset resources/launcher-config.json
//go:generate asset -var=PublicRsaKeys -wrap=ReadPublicRsaKeysAsset resources/public-rsa-keys.pem
//go:generate asset -var=LauncherIcon -wrap=readIconAsset -ignore-missing resources/icon.png

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/logging"

	"github.com/setlog/trivrost/pkg/system"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/launcher"
	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/cmd/launcher/resources"

	"github.com/mattn/go-ieproxy"
	"golang.org/x/net/http/httpproxy"
)

// These are overwritten via ldflags.
var gitDescription string
var gitHash string
var gitBranch string

func init() {
	// On MacOS, only the first thread created by the OS is allowed to be the main GUI thread.
	// Also, on Windows, OLE code needs to run on the main thread, which we rely on when creating shortcuts.
	runtime.LockOSThread()

	system.MustFindPaths()

	flagError := flags.Setup()
	if *flags.PrintBuildTime {
		fmt.Print(launcher.BuildTime())
		os.Exit(0)
	}
	if *flags.Debug {
		log.SetLevel(log.TraceLevel)
	}
	if *flags.DeploymentConfig != "" {
		resources.LauncherConfig.DeploymentConfigURL = *flags.DeploymentConfig
	}
	places.DetectPlaces(*flags.Roaming)

	flags.SetNextLogIndex(logging.Initialize(places.GetAppLogFolderPath(), resources.LauncherConfig.ProductName, *flags.LogIndexCounter, *flags.LogInstanceCounter))
	logState(flagError)
	setGuiStatusMessages(resources.LauncherConfig.StatusMessages)

	printProxySettings()
}

func main() {
	defer misc.LogPanic()
	ctx, cancelFunc := context.WithCancel(context.Background())
	go runLauncher(ctx)
	err := gui.Main(ctx, cancelFunc, resources.LauncherConfig.BrandingName, !*flags.Uninstall)
	if err != nil {
		log.Fatalf("gui.Main() failed: %v\n", err)
	}

	log.Info("End of main().")
	log.Exit(0)
}

func runLauncher(ctx context.Context) {
	gui.WaitUntilReady()
	defer gui.Quit()
	defer handlePanic()

	places.MakePlaces()
	defer launcher.Linger()
	locking.AcquireLock(ctx)
	defer locking.ReleaseLock()

	if *flags.Uninstall {
		log.Info("Goal of this launcher instance: Uninstall.")
		launcher.UninstallPrompt()
	} else if !launcher.IsInstanceInstalled() {
		if launcher.HasInstallation() {
			if launcher.IsInstallationOutdated() {
				log.Info("Goal of this launcher instance: Reinstall.")
				launcher.Install()
			} else {
				log.Info("Goal of this launcher instance: Act as shortcut.")
				if !launcher.RestartWithInstalledBinary() {
					log.Info("Trying to reinstall instead.")
					launcher.Install()
				}
			}
		} else {
			log.Info("Goal of this launcher instance: Install.")
			launcher.Install()
		}
	} else {
		log.Info("Goal of this launcher instance: Run.")
		launcher.Run(ctx)
	}

	log.Info("End of runLauncher().")
}

func handlePanic() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok && errors.Is(err, context.Canceled) {
			log.Infof("Quitting: %v", err)
		} else {
			defer presentError(getPanicMessage(r))
			// The stack printed when panic() is not recover()ed bypasses our file-logging, so log it explicitly here.
			log.Panicf("Exiting due to unrecoverable panic: %v\n%v", r, misc.TryRemoveLines(string(debug.Stack()), 1, 3))
		}
	}
}

func getPanicMessage(r interface{}) string {
	message := "Something went wrong. The program will now close."

	userError, ok := r.(misc.IUserError)
	if ok && !misc.IsNil(userError) {
		return userError.UserError()
	}

	fileSystemError, ok := r.(*system.FileSystemError)
	if ok && fileSystemError != nil {
		if os.IsPermission(fileSystemError.CausingError) {
			return "Error: Insufficient permissions to write files in your own user directory. " +
				"Please contact your system administrator and verify that you have full access to your user directory."
		} else {
			return "Error: Your machine's file system denied a required operation. The error received was: " + fileSystemError.CausingError.Error()
		}
	}

	if !strings.HasSuffix(message, ".") && !strings.HasSuffix(message, "!") && !strings.HasSuffix(message, "?") {
		message += "."
	}

	return message
}

func presentError(message string) {
	if gui.BlockingDialog("Error", fmt.Sprintf("%s\n\nYou can find technical information in the log files under\n%s\n",
		message, places.GetAppLogFolderPath()), []string{"Open log folder and close", "Close"}, 1) == 0 {
		log.Infof("Showing file \"%s\" in file manager.", logging.GetLogFilePath())
		err := system.ShowLocalFileInFileManager(logging.GetLogFilePath())
		if err != nil {
			log.Errorf("Error showing file \"%s\" in file manager: %v", logging.GetLogFilePath(), err)
		}
	}
}

func printProxySettings() {
	envcfg := httpproxy.FromEnvironment()
	log.Infof("Environment proxy: HTTPProxy: \"%s\"; HTTPSProxy: \"%s\".", envcfg.HTTPProxy, envcfg.HTTPSProxy)
	conf := ieproxy.GetConf()
	log.Infof("Automatic proxy: %v; Preconfigured URL: \"%s\".", conf.Automatic.Active, conf.Automatic.PreConfiguredURL)
	log.Infof("Static proxy: %v, Protocols: %v, No proxy: \"%s\".", conf.Static.Active, conf.Static.Protocols, conf.Static.NoProxy)
}

func setGuiStatusMessages(statusMessages config.StatusMessages) {
	setGuiStatusMessage(gui.StageAcquireLock, statusMessages.AcquireLock)
	setGuiStatusMessage(gui.StageGetDeploymentConfig, statusMessages.GetDeploymentConfig)
	setGuiStatusMessage(gui.StageDetermineLocalLauncherVersion, statusMessages.DetermineLocalLauncherVersion)
	setGuiStatusMessage(gui.StageRetrieveRemoteLauncherVersion, statusMessages.RetrieveRemoteLauncherVersion)
	setGuiStatusMessage(gui.StageSelfUpdate, statusMessages.SelfUpdate)
	setGuiStatusMessage(gui.StageDetermineLocalBundleVersions, statusMessages.DetermineLocalBundleVersions)
	setGuiStatusMessage(gui.StageRetrieveRemoteBundleVersions, statusMessages.RetrieveRemoteBundleVersions)
	setGuiStatusMessage(gui.StageAwaitApplicationsTerminated, statusMessages.AwaitApplicationsTerminated)
	setGuiStatusMessage(gui.StageDownloadBundleUpdates, statusMessages.DownloadBundleUpdates)
	setGuiStatusMessage(gui.StageLaunchApplication, statusMessages.LaunchApplication)
}

func setGuiStatusMessage(s gui.Stage, text string) {
	if text != "" {
		gui.SetStageText(s, text)
	}
}

func logState(flagError error) {
	log.Infof("Git commit of this build: Tag: %s; Hash: %s; Branch: %s", gitDescription, gitHash, gitBranch)

	if filepath.Base(system.GetProgramPath()) != resources.LauncherConfig.BinaryName {
		log.Warnf("Program name on disk (\"%s\") has diverged from configured program name (\"%s\").",
			filepath.Base(system.GetProgramPath()), resources.LauncherConfig.BinaryName)
	}

	places.ReportResults()

	if flagError != nil {
		log.Fatalf("Parsing flags failed: %v", flagError)
	}
}
