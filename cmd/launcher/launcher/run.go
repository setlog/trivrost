package launcher

import (
	"context"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/cmd/launcher/resources"

	"github.com/setlog/trivrost/pkg/fetching"
	"github.com/setlog/trivrost/pkg/logging"
	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/pkg/launcher/bundle"
)

func Run(ctx context.Context, launcherFlags *flags.LauncherFlags) {
	doHousekeeping()

	handler := wireHandler(gui.NewGuiDownloadProgressHandler(fetching.MaxConcurrentDownloads))
	updater := wireUpdater(bundle.NewUpdater(ctx, handler, resources.PublicRsaKeys), handler)

	gui.SetStage(gui.StageGetDeploymentConfig, 0)
	updater.RetrieveDeploymentConfig(resources.LauncherConfig.DeploymentConfigURL)

	updateSelf(updater, launcherFlags)
	updateBundles(ctx, updater)

	launch(ctx, updater, launcherFlags)
}

func wireHandler(handler *gui.GuiDownloadProgressHandler) *gui.GuiDownloadProgressHandler {
	hashLauncherProgress, hashBundlesProgress := newProgressFaker(10), newProgressFaker(10)
	gui.ProgressFunc = func(s gui.Stage) uint64 {
		if s.IsDownloadStage() {
			return handler.GetProgress()
		} else if s == gui.StageDetermineLocalLauncherVersion {
			return hashLauncherProgress.getProgress()
		} else if s == gui.StageDetermineLocalBundleVersions {
			return hashBundlesProgress.getProgress()
		}
		return 0
	}
	return handler
}

func wireUpdater(updater *bundle.Updater, handler *gui.GuiDownloadProgressHandler) *bundle.Updater {
	updater.EnableTimestampVerification(places.GetTimestampsFilePath())
	updater.SetStatusCallback(func(status bundle.UpdaterStatus, expectedProgressUnits uint64) {
		handler.ResetProgress()
		handleStatusChange(status, expectedProgressUnits)
	})
	return updater
}

func updateSelf(updater *bundle.Updater, launcherFlags *flags.LauncherFlags) {
	updater.SetIgnoredSelfUpdateBundleInfoSHAs(resources.LauncherConfig.IgnoreLauncherBundleInfoHashes)
	if !(launcherFlags.SkipSelfUpdate || IsInstanceInstalledSystemWide()) {
		if updater.UpdateSelf() {
			runPostBinaryUpdateProvisioning()
			locking.Restart(true, launcherFlags)
		}
	}
}

func updateBundles(ctx context.Context, updater *bundle.Updater) {
	updater.DetermineBundleRequirements(places.GetBundleFolderPath(), places.GetSystemWideBundleFolderPath())
	if updater.HasChangesToUserBundles() {
		locking.AwaitApplicationsTerminated(ctx)
		updater.InstallBundleUpdates()
	}
}

func handleSystemBundleChanges(ctx context.Context, updater *bundle.Updater) {
	const howTo = "To bring the application up to date, its latest release needs to be installed with administrative privileges."
	if updater.HasChangesToSystemBundles(true) {
		panic(misc.NewNestedError("A mandatory update was not applied because it needs to write files in protected system folders. "+howTo, nil))
	} else if updater.HasChangesToSystemBundles(false) {
		gui.Pause(ctx, "Some optional updates were not applied because they need to write files in protected system folders. "+howTo)
	}
}

func handleStatusChange(status bundle.UpdaterStatus, expectedProgressUnits uint64) {
	switch status {
	case bundle.DetermineLocalLauncherVersion:
		gui.SetStage(gui.StageDetermineLocalLauncherVersion, expectedProgressUnits)
	case bundle.RetrieveRemoteLauncherVersion:
		gui.SetStage(gui.StageRetrieveRemoteLauncherVersion, expectedProgressUnits)
	case bundle.DownloadLauncherFiles:
		gui.SetStage(gui.StageSelfUpdate, expectedProgressUnits)
	case bundle.DetermineLocalBundleVersions:
		gui.SetStage(gui.StageDetermineLocalBundleVersions, expectedProgressUnits)
	case bundle.RetrieveRemoteBundleVersions:
		gui.SetStage(gui.StageRetrieveRemoteBundleVersions, expectedProgressUnits)
	case bundle.DownloadBundleFiles:
		gui.SetStage(gui.StageDownloadBundleUpdates, expectedProgressUnits)
	}
}

func doHousekeeping() {
	logging.DeleteOldLogFiles()
	locking.MinimizeApplicationSignaturesList()
	deleteLeftoverBinaries()
}

func launch(ctx context.Context, updater *bundle.Updater, launcherFlags *flags.LauncherFlags) {
	gui.SetStage(gui.StageLaunchApplication, 0)
	handleSystemBundleChanges(ctx, updater)
	execution := updater.GetDeploymentConfig().Execution
	executeCommands(ctx, execution.Commands, launcherFlags)
	lingerTimeMilliseconds = execution.LingerTimeMilliseconds
}
