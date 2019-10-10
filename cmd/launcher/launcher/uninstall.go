package launcher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/setlog/trivrost/cmd/launcher/locking"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/resources"

	"github.com/setlog/trivrost/cmd/launcher/gui"

	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/system"
	log "github.com/sirupsen/logrus"
)

func UninstallPrompt(launcherFlags *flags.LauncherFlags) {
	brandingName := resources.LauncherConfig.BrandingName
	if locking.MinimizeApplicationSignaturesList() {
		title, message := "Uninstall "+brandingName, "You are about to uninstall "+brandingName+". Do you want to continue?"
		if gui.BlockingDialog(title, message, []string{"Yes", "No"}, 1, launcherFlags.DismissGuiPrompts) == 0 || launcherFlags.AcceptUninstall {
			uninstall(launcherFlags)
		}
	} else {
		title, message := "Uninstall "+brandingName, "Cannot uninstall "+brandingName+" while it is still running.\nPlease close all instances and try again."
		gui.BlockingDialog(title, message, []string{"Close"}, 0, launcherFlags.DismissGuiPrompts)
	}
}

func uninstall(launcherFlags *flags.LauncherFlags) {
	log.Info("Uninstall the launcher now.")
	brandingName := resources.LauncherConfig.BrandingName
	gui.ShowWaitDialog("Uninstalling "+brandingName, "Please wait as "+brandingName+" is uninstalling.")
	deletePlainFiles()
	deleteBundles()
	defer prepareProgramDeletionWithFinalizerFunc()()
	gui.HideWaitDialog()
	gui.BlockingDialog("Uninstallation complete", brandingName+" has been uninstalled.", []string{"Close"}, 0, launcherFlags.DismissGuiPrompts)
}

func deletePlainFiles() {
	deleteDesktopShortcuts()
	if runtime.GOOS != system.OsMac {
		deleteStartMenuEntries()
	}
	deleteTimestampFile()
	deleteIcon()
}

func deleteDesktopShortcuts() {
	launchShortcut := places.GetLaunchDesktopShortcutPath()
	err := os.Remove(launchShortcut)
	if err != nil {
		log.Errorf("Could not remove desktop shortcut \"%s\": %v", launchShortcut, err)
	}
}

func deleteStartMenuEntries() {
	system.TryRemove(places.GetUninstallStartMenuShortcutPath())
	system.TryRemove(places.GetLaunchStartMenuShortcutPath())
	system.TryRemoveEmpty(filepath.Dir(places.GetUninstallStartMenuShortcutPath())) // Delete this first because it might be nested in the other.
	system.TryRemoveEmpty(filepath.Dir(places.GetLaunchStartMenuShortcutPath()))
}

func deleteBundles() {
	bundleFolderPath := places.GetBundleFolderPath()
	err := os.RemoveAll(bundleFolderPath)
	if err != nil {
		log.Errorf("Could not remove folder \"%s\": %v", bundleFolderPath, err)
	}
}

func deleteTimestampFile() {
	system.MustRemoveFile(places.GetTimestampsFilePath())
}

func deleteIcon() {
	if runtime.GOOS == system.OsLinux {
		system.MustRemoveFile(places.GetLauncherIconPath())
	}
}

func prepareProgramDeletionWithFinalizerFunc() (finalizerFunc func()) {
	temporaryProgramPath, err := system.UndeployProgram(getTargetProgramPath())
	finalizerFunc = func() {
		err := system.DeleteProgram(temporaryProgramPath)
		if err != nil {
			log.Error(err)
		}
	}
	if err != nil {
		log.Info(err)
		// Delete immediately on Linux and MacOS.
		finalizerFunc()
		finalizerFunc = func() {}
	}
	// Let the caller defer the deletion on Windows.
	return finalizerFunc
}

func deleteLeftoverBinaries() {
	if runtime.GOOS != system.OsWindows {
		return
	}
	targetDir := places.GetLauncherTargetDirectoryPath()
	fileList, err := ioutil.ReadDir(targetDir)
	if err != nil {
		log.Errorf("Could not read directory \"%s\": %v", targetDir, err)
		return
	}
	regEx := regexp.MustCompile(`^~.*\.(old|new|delete)\.[a-fA-F0-9]{16}$`)
	for _, fileInfo := range fileList {
		if !fileInfo.IsDir() && regEx.Match([]byte(fileInfo.Name())) {
			filePath := filepath.Join(targetDir, fileInfo.Name())
			log.Infof("Removing leftover binary \"%s\".", filePath)
			err = os.Remove(filePath)
			if err != nil {
				log.Errorf("Could not remove leftover binary \"%s\": %v", filePath, err)
			}
		}
	}
}
