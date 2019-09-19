package places

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/setlog/trivrost/cmd/launcher/resources"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	globalSettingFolder string
	localSettingFolder  string
	localCacheFolder    string
	desktopFolder       string
	startMenuFolder     string
	appData             string
	localAppData        string
)

func detectPlaces(useRoamingOnly bool) {
	appData = getKnownFolderWithFallback(windows.FOLDERID_RoamingAppData, "APPDATA", "")
	localAppData = getKnownFolderWithFallback(windows.FOLDERID_LocalAppData, "LOCALAPPDATA", "")
	desktopFolder = getKnownFolderWithFallback(windows.FOLDERID_Desktop, "USERPROFILE", "Desktop")
	startMenuFolder = getKnownFolderWithFallback(windows.FOLDERID_StartMenu, "", filepath.Join(appData, "Microsoft", "Windows", "Start Menu"))

	globalSettingFolder = appData

	if useRoamingOnly {
		localSettingFolder = appData
	} else {
		localSettingFolder = localAppData
	}
	localCacheFolder = filepath.Join(localSettingFolder, "Temp")
}

func reportResults() {
	log.Infof("Determined APPDATA folder: \"%s\".", appData)
	log.Infof("Determined LOCALAPPDATA folder: \"%s\".", localAppData)
	log.Infof("Determined Desktop folder: \"%s\".", desktopFolder)
	log.Infof("Determined Start Menu folder: \"%s\".", startMenuFolder)
}

func getKnownFolderWithFallback(guid *windows.KNOWNFOLDERID, envVarName string, envVarValueSuffix string) string {
	folderPath, err := getKnownFolderPath(guid)
	if err != nil {
		log.Errorf("Could not get known folder path for GUID %v: %v", *guid, err)
		if envVarName == "" {
			folderPath = envVarValueSuffix
		} else {
			folderPath = os.Getenv(envVarName)
			if folderPath == "" {
				log.Panicf("Could not fall back to environment variable %%%s%%, because it is empty.", envVarName)
			}
			if strings.Contains(folderPath, "%"+envVarName+"%") {
				log.Panicf("Could not fall back to environment variable %%%s%%, because it is defined infinitely recursively: \"%s\".", envVarName, folderPath)
			}
			folderPath = filepath.Join(folderPath, envVarValueSuffix)
		}
	}
	return folderPath
}

func getKnownFolderPath(guid *windows.KNOWNFOLDERID) (string, error) {
	var flagDoNotVerify uint32 = 0x00004000 // https://docs.microsoft.com/en-us/windows/desktop/api/shlobj_core/ne-shlobj_core-known_folder_flag
	return windows.KnownFolderPath(guid, flagDoNotVerify)
}

func getLaunchDesktopShortcutPath() string {
	return filepath.Join(desktopFolder, resources.LauncherConfig.BrandingName+".lnk")
}

func getLaunchStartMenuShortcutPath() string {
	return filepath.Join(startMenuFolder, resources.LauncherConfig.VendorName, resources.LauncherConfig.BrandingName+".lnk")
}

func getUninstallStartMenuShortcutPath() string {
	return filepath.Join(startMenuFolder, resources.LauncherConfig.VendorName, "Uninstall", "Uninstall "+resources.LauncherConfig.BrandingName+".lnk")
}
