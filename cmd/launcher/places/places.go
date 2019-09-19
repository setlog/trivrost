package places

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/setlog/trivrost/cmd/launcher/resources"
	"github.com/setlog/trivrost/pkg/system"
)

func MakePlaces() {
	err := os.MkdirAll(GetAppCacheFolderPath(), 0700)
	if err != nil {
		panic(fmt.Sprintf("Could not create app cache folder \"%s\": %v", GetAppCacheFolderPath(), err))
	}
	err = os.MkdirAll(GetAppDataFolderPath(), 0700)
	if err != nil {
		panic(fmt.Sprintf("Could not create app cache folder \"%s\": %v", GetAppDataFolderPath(), err))
	}
	err = os.MkdirAll(GetAppLocalDataFolderPath(), 0700)
	if err != nil {
		panic(fmt.Sprintf("Could not create app cache folder \"%s\": %v", GetAppLocalDataFolderPath(), err))
	}
}

func DetectPlaces(useRoamingOnly bool) {
	detectPlaces(useRoamingOnly)
}

func ReportResults() {
	reportResults()
}

func GetAppCacheFolderPath() string {
	return filepath.Join(localCacheFolder, resources.LauncherConfig.VendorName, resources.LauncherConfig.ProductName)
}

func GetAppDataFolderPath() string {
	return filepath.Join(globalSettingFolder, resources.LauncherConfig.VendorName, resources.LauncherConfig.ProductName)
}

func GetAppLocalDataFolderPath() string {
	return filepath.Join(localSettingFolder, resources.LauncherConfig.VendorName, resources.LauncherConfig.ProductName)
}

func GetAppLogFolderPath() string {
	return filepath.Join(GetAppCacheFolderPath(), "log")
}

func GetLauncherTargetDirectoryPath() string {
	return GetAppDataFolderPath()
}

func GetLauncherIconPath() string {
	return filepath.Join(GetAppLocalDataFolderPath(), "icon.png")
}

func GetSystemWideBundleFolderPath() string {
	return filepath.Join(filepath.Dir(system.GetProgramPath()), "systembundles")
}

func GetBundleFolderPath() string {
	return filepath.Join(GetAppLocalDataFolderPath(), "bundles")
}

func GetLaunchDesktopShortcutPath() string {
	return getLaunchDesktopShortcutPath()
}

func GetLaunchStartMenuShortcutPath() string {
	return getLaunchStartMenuShortcutPath()
}

func GetUninstallStartMenuShortcutPath() string {
	return getUninstallStartMenuShortcutPath()
}

func GetTimestampsFilePath() string {
	return filepath.Join(GetAppLocalDataFolderPath(), "timestamps.json")
}
