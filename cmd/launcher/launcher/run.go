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

	"github.com/setlog/trivrost/pkg/launcher/bundle"
	"github.com/setlog/trivrost/pkg/launcher/config"
)

func Run(ctx context.Context) {
	doHousekeeping()

	handler := gui.NewGuiDownloadProgressHandler(fetching.MaxConcurrentDownloads)
	updater := bundle.NewUpdater(ctx, handler, resources.PublicRsaKeys)
	updater.EnableTimestampVerification(places.GetTimestampsFilePath())
	updater.SetStatusCallback(func(status bundle.UpdaterStatus, expectedProgressUnits uint64) {
		handler.ResetProgress()
		handleStatusChange(status, expectedProgressUnits)
	})

	gui.ProgressFunc = func(s gui.Stage) uint64 {
		if s.IsDownloadStage() {
			return handler.GetProgress()
		}
		return 0
	}
	gui.SetStage(gui.StageGetDeploymentConfig, 0)

	updater.RetrieveDeploymentConfig(resources.LauncherConfig.DeploymentConfigURL)

	updater.SetIgnoredSelfUpdateBundleInfoSHAs(resources.LauncherConfig.IgnoreLauncherBundleInfoHashes)
	if !(*flags.SkipSelfUpdate || IsInstanceInstalledSystemWide()) {
		if updater.UpdateSelf() {
			locking.Restart(true)
		}
	}
	if updater.DetermineBundleUpdateRequired(places.GetBundleFolderPath(), places.GetSystemWideBundleFolderPath()) {
		locking.AwaitApplicationsTerminated(ctx)
		updater.InstallBundleUpdates()
	}

	launch(&updater.GetDeploymentConfig().Execution)
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

func launch(execution *config.ExecutionConfig) {
	executeCommands(execution.Commands)
	lingerTimeMilliseconds = execution.LingerTimeMilliseconds
}
