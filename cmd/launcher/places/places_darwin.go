package places

import "os"
import "path/filepath"
import "github.com/setlog/trivrost/cmd/launcher/resources"

var globalSettingFolder = os.Getenv("HOME") + "/Library/Application Support"
var localSettingFolder = globalSettingFolder
var localCacheFolder = os.Getenv("HOME") + "/Library/Caches"

func detectPlaces(useRoamingOnly bool) {
}

func reportResults() {
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
