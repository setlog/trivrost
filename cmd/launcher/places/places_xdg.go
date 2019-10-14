// +build !windows,!darwin

package places

import (
	"fmt"
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
var warning error

func detectPlaces(useRoamingOnly bool) error {
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
		warning = fmt.Errorf("Could not run xdg-user-dir to locate DESKTOP folder: %v. Falling back to $HOME/Desktop", err)
	} else {
		desktopFolder = strings.Trim(string(output), "\n\"")
	}
	return nil
}

func reportResults() {
	if warning != nil {
		log.Warn(warning)
	}
	log.Infof("globalSettingFolder: %v", globalSettingFolder)
	log.Infof("localSettingFolder: %v", localSettingFolder)
	log.Infof("localCacheFolder: %v", localCacheFolder)
	log.Infof("desktopFolder: %v", desktopFolder)
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
