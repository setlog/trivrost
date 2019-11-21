package launcher

import (
	"context"
	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"

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

	updater := createUpdater(ctx, wireHandler(gui.NewGuiDownloadProgressHandler(fetching.MaxConcurrentDownloads)))

	gui.SetStage(gui.StageGetDeploymentConfig, 0)
	isSelfUpdateMandatory := updater.Prepare(resources.LauncherConfig.DeploymentConfigURL)

	errSelfUpdate := updateSelf(updater, launcherFlags)
	if isSelfUpdateMandatory && system.IsPermission(errSelfUpdate) {
		handleInsufficientPrivileges(ctx, true)
	}
	updateBundles(ctx, updater)

	gui.SetStage(gui.StageLaunchApplication, 0)
	handleUpdateOmissions(ctx, updater, system.IsPermission(errSelfUpdate))
	launch(ctx, updater.GetDeploymentConfig().Execution, launcherFlags)
}

func doHousekeeping() {
	logging.DeleteOldLogFiles()
	locking.MinimizeApplicationSignaturesList()
	deleteLeftoverBinaries()
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

func createUpdater(ctx context.Context, handler *gui.GuiDownloadProgressHandler) *bundle.Updater {
	updater := bundle.NewUpdater(ctx, handler, resources.PublicRsaKeys)
	updater.EnableTimestampVerification(places.GetTimestampsFilePath())
	updater.SetStatusCallback(func(status bundle.UpdaterStatus, expectedProgressUnits uint64) {
		handler.ResetProgress()
		handleStatusChange(status, expectedProgressUnits)
	})
	return updater
}

func updateSelf(updater *bundle.Updater, launcherFlags *flags.LauncherFlags) (err error) {
	updater.SetIgnoredSelfUpdateBundleInfoSHAs(resources.LauncherConfig.IgnoreLauncherBundleInfoHashes)
	if !(launcherFlags.SkipSelfUpdate || IsInstanceInstalledSystemWide()) {
		defer permissionPanicToError(&err)
		if updater.UpdateSelf() {
			runPostBinaryUpdateProvisioning()
			locking.Restart(true, launcherFlags)
		}
	}
	return nil
}

func permissionPanicToError(errPtr *error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok && system.IsPermission(err) {
			*errPtr = err
		} else {
			panic(r)
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

func handleUpdateOmissions(ctx context.Context, updater *bundle.Updater, wasAtLeastOneOptionalUpdateOmittedDueToInsufficientPrivileges bool) {
	if updater.HasChangesToSystemBundles(true) {
		handleInsufficientPrivileges(ctx, true)
	} else if updater.HasChangesToSystemBundles(false) || wasAtLeastOneOptionalUpdateOmittedDueToInsufficientPrivileges {
		handleInsufficientPrivileges(ctx, false)
	}
}

func handleInsufficientPrivileges(ctx context.Context, wasAtLeastOneMandatoryUpdateOmitted bool) {
	const howTo = "To bring the application up to date, its latest release needs to be installed with administrative privileges.\n" +
		"You may click \"Continue\" to ignore this for the time being."
	if wasAtLeastOneMandatoryUpdateOmitted {
		panic(misc.NewNestedError("A mandatory update was not applied because it needs to write files in protected system folders. "+howTo, nil))
	} else {
		gui.Pause(ctx, "Some optional updates were not applied because they need to write files in protected system folders. "+howTo)
	}
}

func launch(ctx context.Context, execution config.ExecutionConfig, launcherFlags *flags.LauncherFlags) {
	executeCommands(ctx, execution.Commands, launcherFlags)
	lingerTimeMilliseconds = execution.LingerTimeMilliseconds
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
