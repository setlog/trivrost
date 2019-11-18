package launcher

import (
	"os/exec"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/places"
	log "github.com/sirupsen/logrus"
)

func runPostBinaryUpdateProvisioning() {
}

func createLaunchDesktopShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetLaunchDesktopShortcutPath()
	createShortcutOSX(shortcutLocation, destination)
}

func createLaunchStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	// Not on OSX
}

func createUninstallStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	// Not on OSX
}

func createShortcutOSX(atPath string, destination string) {
	log.Debugf(`Creating soft link to "%s" at "%s".`, destination, atPath)
	c := exec.Command("ln", "-sfn", destination, atPath)
	output, err := c.CombinedOutput()
	if err != nil {
		log.Errorf(`Could not create shortcut "%s" to "%s": %v: %s`, atPath, destination, err, string(output))
	}
}
