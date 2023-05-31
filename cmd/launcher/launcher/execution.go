package launcher

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/setlog/trivrost/cmd/launcher/flags"
	"github.com/setlog/trivrost/cmd/launcher/gui"
	"github.com/setlog/trivrost/cmd/launcher/locking"
	"github.com/setlog/trivrost/cmd/launcher/places"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/system"
)

func executeCommands(ctx context.Context, commandConfigs []config.Command, launcherFlags *flags.LauncherFlags) {
	log.Infof("Executing %d command(s)...", len(commandConfigs))

	lastIndex := len(commandConfigs) - 1
	for i, commandConfig := range commandConfigs {
		command, procSig := executeCommand(ctx, commandConfig, launcherFlags)
		if i != lastIndex {
			err := command.Wait()
			if err != nil {
				log.Errorf("Could not wait for command \"%s\": %v", command.Path, err)
			}
		} else {
			locking.AddApplicationSignature(procSig)
		}
	}
}

func executeCommand(ctx context.Context, commandConfig config.Command, launcherFlags *flags.LauncherFlags) (*exec.Cmd, *system.ProcessSignature) {
	commandWorkingDirectory := findWorkingDirectoryByBundle(commandConfig.WorkingDirectoryBundleName)
	commandBinaryPath := findMatchingExecutablePath(filepath.FromSlash(commandConfig.Name))
	for {
		finalEnv := mergeMaps(commandConfig.Env, launcherFlags.ExtraEnvs)
		log.Infof("Trying to start binary \"%s\" with working directory \"%s\" and args %v", commandBinaryPath, commandWorkingDirectory, commandConfig.Arguments)
		command, procSig, err := system.StartProcess(commandBinaryPath, commandWorkingDirectory, commandConfig.Arguments, finalEnv, !launcherFlags.NoStreamPassing)
		if err != nil {
			log.Info(err)
			gui.NotifyProblem(fmt.Sprintf("System denies launch of \"%s\"", filepath.Base(command.Path)), true)
			misc.MustWaitForContext(ctx, time.Second*3)
		} else {
			log.Infof("Started binary \"%s\" with working directory \"%s\" and args %v", commandBinaryPath, commandWorkingDirectory, commandConfig.Arguments)
			gui.ClearProblem()
			return command, procSig
		}
	}
}

func mergeMaps(map1 map[string]*string, map2 map[string]string) map[string]*string {
	result := make(map[string]*string)
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		v2 := v
		result[k] = &v2
	}
	return result
}

func findWorkingDirectoryByBundle(bundleName string) string {
	if bundleName == "" {
		return places.GetBundleFolderPath()
	}

	systemWideBundlePath := filepath.Join(places.GetSystemWideBundleFolderPath(), bundleName)
	if system.FolderExists(systemWideBundlePath) {
		return places.GetSystemWideBundleFolderPath()
	}
	return places.GetBundleFolderPath()
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
