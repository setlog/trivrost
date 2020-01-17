package launcher

import (
	"context"

	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/cmd/launcher/places"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
)

func LauncherMain(ctx context.Context, launcherFlags *flags.LauncherFlags) {
	places.MakePlaces()
	defer Linger()
	locking.AcquireLock(ctx, launcherFlags)
	defer locking.ReleaseLock()

	Branch(ctx, launcherFlags)

	log.Info("End of LauncherMain().")
}

func Branch(ctx context.Context, launcherFlags *flags.LauncherFlags) {
	if launcherFlags.Uninstall {
		log.Info("Goal of this launcher instance: Uninstall.")
		UninstallPrompt(launcherFlags)
	} else if !IsInstanceInstalled() {
		if HasInstallation() {
			if IsInstallationOutdated() {
				log.Info("Goal of this launcher instance: Reinstall.")
				Install(launcherFlags)
			} else {
				log.Info("Goal of this launcher instance: Act as shortcut.")
				if !RestartWithInstalledBinary(launcherFlags) {
					log.Info("Trying to reinstall instead.")
					Install(launcherFlags)
				}
			}
		} else {
			log.Info("Goal of this launcher instance: Install.")
			Install(launcherFlags)
		}
	} else {
		log.Info("Goal of this launcher instance: Run.")
		Run(ctx, launcherFlags)
	}
}
