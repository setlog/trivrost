// +build !windows,!darwin

package places

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/setlog/trivrost/cmd/launcher/resources"

	log "github.com/sirupsen/logrus"
)

var globalSettingFolder string
var localSettingFolder string
var localCacheFolder string
var desktopFolder string

func detectPlaces(useRoamingOnly bool) {
	globalSettingFolder = filepath.Join(os.Getenv("HOME"), ".local", "share")
	localSettingFolder = globalSettingFolder
	if os.Getenv("XDG_CACHE_HOME") != "" {
		localCacheFolder = os.Getenv("XDG_CACHE_HOME")
	} else {
		localCacheFolder = filepath.Join(os.Getenv("HOME"), ".cache")
	}

	desktopFolder = filepath.Join(os.Getenv("HOME"), "Desktop")
	xdgCommand := exec.Command("xdg-user-dir", "DESKTOP")
	output, err := xdgCommand.CombinedOutput()
	if err != nil {
		log.Errorf("Could not run xdg-user-dir to locate DESKTOP folder: %v", err)
	} else {
		desktopFolder = strings.Trim(string(output), "\n\"")
	}
}

func reportResults() {
}

func getLaunchDesktopShortcutPath() string {
	return filepath.Join(desktopFolder, resources.LauncherConfig.BrandingName+".desktop")
}

func getLaunchStartMenuShortcutPath() string {
	return filepath.Join(globalSettingFolder, "applications", resources.LauncherConfig.VendorName, resources.LauncherConfig.BrandingName+".desktop")
}

func getUninstallStartMenuShortcutPath() string {
	return filepath.Join(globalSettingFolder, "applications", resources.LauncherConfig.VendorName, "Uninstall", "Uninstall "+resources.LauncherConfig.BrandingName+".desktop")
}
