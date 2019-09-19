package bundle

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/hashing"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/system"
)

func (u *Updater) UpdateSelf() (needsRestart bool) {
	updateConfig := u.deploymentConfig.GetLauncherUpdateConfig()
	if updateConfig == nil {
		return false
	}
	programPath := system.GetProgramPath()
	log.Infof("Checking for update of launcher at \"%s\".", programPath)
	return u.updateProgram(programPath)
}

func (u *Updater) updateProgram(programPath string) (needsRestart bool) {
	log.Infof("Calculating local hashes.")
	u.announceStatus(DetermineLocalLauncherVersion, 0)
	presentState := hashing.MustHash(programPath)

	log.Infof("Checking against latest version.")
	u.announceStatus(RetrieveRemoteLauncherVersion, 0)
	updateConfig := u.deploymentConfig.GetLauncherUpdateConfig()
	bundleInfo, bundleInfoSha := u.RetrieveBundleInfo(updateConfig.BundleInfoURL, u.publicKeys)
	if u.IsShaIgnored(bundleInfoSha) {
		log.Warnf("Ignoring launcher bundleinfo with sha \"%s\".", bundleInfoSha)
		return false
	}

	remoteState := bundleInfo.GetFileHashes().ForOS()
	wantedState := config.MakeDiffFileInfoMap(presentState.Prepend(remoteState.FirstPathElement(filepath.Separator), filepath.Separator), remoteState)

	if wantedState.HasChanges() {
		log.WithFields(log.Fields{"updateConfig": fmt.Sprintf("%+v", updateConfig)}).
			Infof(spew.Sprintf("Launcher at \"%s\" is outdated. Updating from state %+v to %+v.", programPath, presentState, wantedState))
		if system.IsDir(programPath) {
			u.updateApplicationFolder(updateConfig, wantedState, programPath)
		} else {
			u.updateApplicationBinary(updateConfig, wantedState, programPath)
		}
		return true
	}
	return false
}

func (u *Updater) updateApplicationFolder(updateConfig *config.LauncherUpdateConfig, wantedState config.FileInfoMap, programPath string) {
	u.announceStatus(DownloadLauncherFiles, wantedState.UpdateByteCount())
	tempPath := u.downloader.MustDownloadToTempDirectory(updateConfig.BaseURL, wantedState, programPath)
	defer system.TryRemoveDirectory(tempPath)
	firstPathElement := wantedState.FirstPathElement(filepath.Separator)
	applyBundleUpdate(wantedState.StripFirstPathElement(filepath.Separator), filepath.Join(tempPath, firstPathElement), programPath)
}

func (u *Updater) updateApplicationBinary(updateConfig *config.LauncherUpdateConfig, wantedState config.FileInfoMap, programPath string) {
	binaryName, newFileInfo := wantedState.MustGetOnly()
	u.swapBinary(programPath, misc.MustJoinURL(updateConfig.BaseURL, binaryName), newFileInfo)
}

func (u *Updater) swapBinary(localBinaryPath string, remoteURL string, newFileInfo *config.FileInfo) {
	u.announceStatus(DownloadLauncherFiles, uint64(newFileInfo.Size))

	randomHex := misc.MustGetRandomHexString(8)
	oldBinaryNewPath := filepath.Join(filepath.Dir(localBinaryPath), "~"+filepath.Base(localBinaryPath)+".old."+randomHex)
	newBinaryTempPath := filepath.Join(filepath.Dir(localBinaryPath), "~"+filepath.Base(localBinaryPath)+".new."+randomHex)
	err := u.downloader.DownloadFile(remoteURL, newFileInfo, newBinaryTempPath)
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == system.OsWindows { // On Windows, you cannot delete a running binary, but you can rename it.
		if err := os.Rename(localBinaryPath, oldBinaryNewPath); err != nil {
			panic(&system.FileSystemError{Message: fmt.Sprintf("Could not rename old binary \"%s\" to \"%s\"", localBinaryPath, oldBinaryNewPath), CausingError: err})
		}
	}

	if err := os.Rename(newBinaryTempPath, localBinaryPath); err != nil {
		if runtime.GOOS == system.OsWindows {
			if err2 := os.Rename(oldBinaryNewPath, localBinaryPath); err2 != nil {
				log.WithFields(log.Fields{"err": err, "localBinaryPath": localBinaryPath, "oldBinaryNewPath": oldBinaryNewPath}).
					Error("Could not revert rename. Installation will be broken.")
			}
		}
		panic(&system.FileSystemError{Message: fmt.Sprintf("Could not rename new binary \"%s\" to \"%s\"", newBinaryTempPath, localBinaryPath), CausingError: err})
	}
}

func (u *Updater) SetIgnoredSelfUpdateBundleInfoSHAs(ignoreShas []string) {
	u.ignoredSelfUpdateBundleInfoSHAs = ignoreShas
}

func (u *Updater) IsShaIgnored(sha string) bool {
	for _, ignoreSha := range u.ignoredSelfUpdateBundleInfoSHAs {
		if ignoreSha == sha {
			return true
		}
	}
	return false
}
