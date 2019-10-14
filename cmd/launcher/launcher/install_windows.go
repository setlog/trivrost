package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/places"

	log "github.com/sirupsen/logrus"
)

func prepareShortcutInstallation() {
}

func createLaunchDesktopShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetLaunchDesktopShortcutPath()
	createShortcutWindows(shortcutLocation, destination, getArgs(nil, launcherFlags))
}

func createLaunchStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetLaunchStartMenuShortcutPath()
	createShortcutWindows(shortcutLocation, destination, getArgs(nil, launcherFlags))
}

func createUninstallStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetUninstallStartMenuShortcutPath()
	createShortcutWindows(shortcutLocation, destination, getArgs([]string{"-" + flags.UninstallFlag}, launcherFlags))
}

func getArgs(baseArgs []string, launcherFlags *flags.LauncherFlags) string {
	if launcherFlags.Roaming {
		baseArgs = append(baseArgs, "-"+flags.RoamingFlag)
	}
	return strings.Join(baseArgs, " ")
}

func createShortcutWindows(location, destination string, arguments string) {
	err := os.MkdirAll(filepath.Dir(location), 0700)
	if err != nil {
		panic(fmt.Sprintf("Could not create directory \"%s\": %v", filepath.Dir(location), err))
	}

	err = os.Remove(location) // The OLE code below cannot overwrite the shortcut, so we remove it here if it exists.
	if err != nil {
		if !os.IsNotExist(err) {
			log.Errorf("Could not remove shortcut \"%s\": %v", location, err)
		}
	}

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		panic(fmt.Sprintf("Could not create OLE shell object: %v", err))
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(fmt.Sprintf("Could not query interface: %v", err))
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", location)
	if err != nil {
		panic(fmt.Sprintf("Could not call method: %v", err))
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", destination)
	if arguments != "" {
		oleutil.PutProperty(idispatch, "Arguments", arguments)
	}
	oleutil.CallMethod(idispatch, "Save")
	log.Infof("Installed shortcut \"%s\" which links to \"%s\".\n", location, destination)
}
