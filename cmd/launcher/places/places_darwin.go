package places

import (
	"os"
	"path/filepath"

	"github.com/setlog/trivrost/cmd/launcher/resources"
	log "github.com/sirupsen/logrus"
)

var globalSettingFolder = os.Getenv("HOME") + "/Library/Application Support"
var localSettingFolder = globalSettingFolder
var localCacheFolder = os.Getenv("HOME") + "/Library/Caches"

func detectPlaces(useRoamingOnly bool) error {
	return nil
}

func reportResults() {
	log.Infof("globalSettingFolder: %v", globalSettingFolder)
	log.Infof("localSettingFolder: %v", localSettingFolder)
	log.Infof("localCacheFolder: %v", localCacheFolder)
}

func getLaunchDesktopShortcutPath() string {
	return filepath.Join(os.Getenv("HOME"), "Desktop", resources.LauncherConfig.BrandingName)
}

func getLaunchStartMenuShortcutPath() string {
	return filepath.Join(globalSettingFolder, resources.LauncherConfig.VendorName, resources.LauncherConfig.BrandingName)
}

func getUninstallStartMenuShortcutPath() string {
	return filepath.Join(globalSettingFolder, resources.LauncherConfig.VendorName, "Uninstall", "Uninstall "+resources.LauncherConfig.BrandingName)
}
