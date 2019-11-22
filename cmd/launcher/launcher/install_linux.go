package launcher

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/setlog/trivrost/pkg/system"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/cmd/launcher/resources"
	"github.com/setlog/trivrost/pkg/misc"
)

const desktopFileTemplate string = `[Desktop Entry]
Name={{.Name}}
Comment={{.Comment}}
GenericName=Application Launcher
Exec="{{.Exec}}"
Icon={{.Icon}}
Type=Application
StartupNotify=false
Categories=Utility;
Actions=uninstall;
Keywords={{.Keywords}};

[Desktop Action uninstall]
Name=Uninstall {{.Name}}
Exec="{{.Exec}}" -uninstall
Icon={{.Icon}}
`

type DesktopFileData struct {
	Name     string
	Comment  string
	Exec     string
	Icon     string
	Keywords string
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("freedesktop").Parse(desktopFileTemplate)
	if err != nil {
		panic(fmt.Sprintf("Could not parse template: %v", err))
	}
}

func runPostBinaryUpdateProvisioning() {
	installLauncherIcon()
}

func installLauncherIcon() {
	if len(resources.LauncherIcon) == 0 {
		log.Errorf("Could not install launcher icon: Size was 0 bytes.")
		return
	}
	system.MustPutFile(places.GetLauncherIconPath(), resources.LauncherIcon)
	resources.LauncherIcon = nil // Icon can be pretty large. No reason to keep it around.
}

func createLaunchDesktopShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetLaunchDesktopShortcutPath()
	desktopFileData := getDesktopFileData(misc.ExtensionlessFileName(shortcutLocation), destination)
	createFreeDesktopStandardShortcut(shortcutLocation, desktopFileData)
}

func createLaunchStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetLaunchStartMenuShortcutPath()
	desktopFileData := getDesktopFileData(misc.ExtensionlessFileName(shortcutLocation), destination)
	createFreeDesktopStandardShortcut(shortcutLocation, desktopFileData)
}

func createUninstallStartMenuShortcut(destination string, launcherFlags *flags.LauncherFlags) {
	shortcutLocation := places.GetUninstallStartMenuShortcutPath()
	desktopFileData := getDesktopFileData(misc.ExtensionlessFileName(shortcutLocation), destination+" -"+flags.UninstallFlag)
	createFreeDesktopStandardShortcut(shortcutLocation, desktopFileData)
}

func createFreeDesktopStandardShortcut(atPath string, desktopFileData DesktopFileData) {
	err := os.MkdirAll(filepath.Dir(atPath), 0744)
	if err != nil {
		panic(fmt.Sprintf("Could not create directory \"%s\": %v", filepath.Dir(atPath), err))
	}
	desktopFile, err := os.OpenFile(atPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0744)
	if err != nil {
		panic(fmt.Sprintf("Could not open file for writing: %v", err))
	}
	defer desktopFile.Close()
	err = tmpl.Execute(desktopFile, desktopFileData)
	if err != nil {
		panic(fmt.Sprintf("Could not execute template: %v", err))
	}
}

func getDesktopFileData(name, command string) DesktopFileData {
	return DesktopFileData{name, "", command, places.GetLauncherIconPath(), ""}
}
