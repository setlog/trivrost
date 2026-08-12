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
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/logging"

	"github.com/setlog/trivrost/pkg/system"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/launcher"
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
	runtime.LockOSThread() // This call must be in init().

	postInit()
}

func main() {
	defer misc.LogPanic()
	launcherFlags, envErr := initializeEnvironment()
	ctx, cancelFunc := context.WithCancel(context.Background())

	go runLauncher(ctx, envErr, launcherFlags)
	runGUI(ctx, cancelFunc, launcherFlags, envErr == nil)

	log.Exit(0)
}

func initializeEnvironment() (*flags.LauncherFlags, error) {
	registerSignalOverrides()

	if resources.LauncherConfig == nil {
		panic("LauncherConfig was not found / compiled into the binary.")
	}
	launcherFlags, argumentError, flagError, pathError, evalError, placesError := parseEnvironment()
	launcherFlags.SetNextLogIndex(logging.Initialize(places.GetAppLogFolderPath(), resources.LauncherConfig.ProductName,
		launcherFlags.LogIndexCounter, launcherFlags.LogInstanceCounter))
	logState(argumentError, flagError, pathError, evalError)

	printProxySettings()
	setGuiStatusMessages(resources.LauncherConfig.StatusMessages)
	return launcherFlags, misc.NewUserErrorFromErrors(argumentError, flagError, pathError, placesError)
}

func runLauncher(ctx context.Context, fatalError error, launcherFlags *flags.LauncherFlags) {
	gui.WaitUntilReady()
	defer gui.Quit()
	if fatalError != nil {
		gui.PanicInformatively(fatalError, launcherFlags)
	}
	defer gui.HandlePanic(launcherFlags)
	launcher.LauncherMain(ctx, launcherFlags)
}

func runGUI(ctx context.Context, cancelFunc context.CancelFunc, launcherFlags *flags.LauncherFlags, showMainWindow bool) {
	err := gui.Main(ctx, cancelFunc, resources.LauncherConfig.BrandingName, !launcherFlags.Uninstall && showMainWindow)
	if err != nil {
		log.Fatalf("gui.Main() failed: %v\n", err)
	}
	log.Info("runGUI() terminated.")
}

func parseEnvironment() (launcherFlags *flags.LauncherFlags, argumentError, flagError, pathError, evalError, placesError error) {
	launcherFlags = &flags.LauncherFlags{}
	if len(os.Args) == 0 {
		argumentError = misc.UserErrorf(fmt.Errorf("len(os.Args) was 0"), "Your system launched the application without any arguments, but there must be at least one.")
	} else {
		launcherFlags, flagError = processFlags(os.Args)
		pathError, evalError = system.FindPaths()
		placesError = places.DetectPlaces(launcherFlags.Roaming)
	}
	return launcherFlags, argumentError, flagError, pathError, evalError, placesError
}

func registerSignalOverrides() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGHUP)
	go handleSignals(sigChan)
}

func handleSignals(sigChan chan os.Signal) {
	for {
		s, ok := <-sigChan
		if ok {
			log.Errorf("Received signal \"%v\". Printing stacks and quitting.", s)
		} else {
			log.Errorf("Signal channel has been closed unexpectedly. Printing stacks and quitting.")
		}
		log.Info("\n")
		pprof.Lookup("goroutine").WriteTo(log.StandardLogger().Out, 2)
		log.Exit(1)
	}
}

func processFlags(args []string) (launcherFlags *flags.LauncherFlags, err error) {
	launcherFlags, err = flags.Setup(args)
	if launcherFlags.PrintBuildTime {
		fmt.Print(launcher.BuildTime())
		os.Exit(0)
	}
	if launcherFlags.Debug {
		log.SetLevel(log.TraceLevel)
	}
	if launcherFlags.DeploymentConfig != "" {
		resources.LauncherConfig.DeploymentConfigURL = launcherFlags.DeploymentConfig
	}
	return launcherFlags, err
}

func logState(argumentError, flagError, pathError, evalError error) {
	log.Infof("Git description of this build: %s; Commit hash: %s; Branch: %s; Built with %v", gitDescription, gitHash, gitBranch, runtime.Version())

	launcher.ReportProgramNameDivergence()

	if argumentError != nil {
		log.Errorf("Fatal: Parsing arguments failed: %v", argumentError)
	}

	if flagError != nil {
		log.Errorf("Fatal: Parsing flags failed: %v", flagError)
	}

	if pathError != nil {
		log.Errorf("Fatal: Determining binary path failed: %v", pathError)
	}

	if evalError != nil {
		log.Warnf("Evaluating binary path failed: %v", evalError)
	}

	places.ReportResults()
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
