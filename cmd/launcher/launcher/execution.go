package launcher

import (
	"path/filepath"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/cmd/launcher/places"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func executeCommands(commands []config.Command, launcherFlags *flags.LauncherFlags) {
	log.Infof("Executing %d command(s)...", len(commands))
	gui.SetStage(gui.StageLaunchApplication, 0)

	lastIndex := len(commands) - 1
	for i, c := range commands {
		wait := i != lastIndex
		commandWorkingDirectory := places.GetBundleFolderPath()
		commandBinaryPath := findMatchingExecutablePath(filepath.FromSlash(c.Name))
		log.Infof("Trying to start binary \"%s\" with working directory \"%s\" and args %v", commandBinaryPath, commandWorkingDirectory, c.Arguments)
		command, procSig := system.MustStartProcess(commandBinaryPath, commandWorkingDirectory, c.Arguments, c.Env, !launcherFlags.NoStreamPassing)
		log.Infof("Started binary \"%s\" with working directory \"%s\" and args %v", commandBinaryPath, commandWorkingDirectory, c.Arguments)
		if wait {
			err := command.Wait()
			if err != nil {
				log.Errorf("Could not wait for command \"%s\": %v", commandBinaryPath, err)
			}
		} else {
			locking.AddApplicationSignature(procSig)
		}
	}
}

func findMatchingExecutablePath(filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	systemWideBinaryPath := filepath.Join(places.GetSystemWideBundleFolderPath(), filePath)
	if system.FileExists(systemWideBinaryPath) {
		return systemWideBinaryPath
	}
	return filePath
}
