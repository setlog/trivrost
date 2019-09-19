package bundle

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/setlog/trivrost/pkg/launcher/timestamps"
	"github.com/setlog/trivrost/pkg/system"

	"github.com/MMulthaupt/chronometry"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/hashing"
)

func (u *Updater) DetermineBundleUpdateRequired(userBundlesFolderPath, systemBundlesFolderPath string) bool {
	u.userBundlesFolderPath, u.systemBundlesFolderPath = userBundlesFolderPath, systemBundlesFolderPath
	u.determineLocalBundleVersions()
	u.removeUnknownBundles()
	u.determineBundleChanges()
	isUpdatedRequired := u.isAtLeastOneChangeRequired()
	if isUpdatedRequired {
		u.assertUpdatePossible()
	}
	return isUpdatedRequired
}

func (u *Updater) InstallBundleUpdates() {
	log.Infof("Downloading bundle updates.")
	u.installBundleUpdates()
}

func (u *Updater) determineLocalBundleVersions() {
	u.announceStatus(DetermineLocalBundleVersions, 0)
	for _, bundleConfig := range u.deploymentConfig.Bundles {
		var bundleUpdateInfo *BundleUpdateInfo
		if u.haveSystemBundleWithName(bundleConfig.LocalDirectory) {
			bundleUpdateInfo = u.makeBundleUpdateConfigFromBundle(bundleConfig, u.systemBundlesFolderPath)
			bundleUpdateInfo.IsSystemBundle = true
			log.Debugf("Identified bundle \"%s\" as system bundle.", bundleConfig.LocalDirectory)
		} else {
			bundleUpdateInfo = u.makeBundleUpdateConfigFromBundle(bundleConfig, u.userBundlesFolderPath)
			log.Debugf("Identified bundle \"%s\" as user bundle.", bundleConfig.LocalDirectory)
		}
		u.bundleUpdateInfos = append(u.bundleUpdateInfos, bundleUpdateInfo)
	}
}

func (u *Updater) makeBundleUpdateConfigFromBundle(bundleConfig config.BundleConfig, bundleFolderPath string) *BundleUpdateInfo {
	bundleUpdateConfig := BundleUpdateInfo{BundleConfig: bundleConfig}
	sw := chronometry.NewStartedStopwatch()
	bundleUpdateConfig.PresentState = hashing.MustHash(filepath.Join(bundleFolderPath, bundleConfig.LocalDirectory))
	log.Infof("Hashing directory of bundle \"%s\" took %v.", bundleConfig.LocalDirectory, sw.TakeLapTime())
	return &bundleUpdateConfig
}

func (u *Updater) removeUnknownBundles() {
	fileInfos, err := ioutil.ReadDir(u.userBundlesFolderPath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() && !u.wantBundleWithName(fileInfo.Name()) {
			removePath := filepath.Join(u.userBundlesFolderPath, fileInfo.Name())
			log.Infof("Remove unknown bundle folder \"%s\".", removePath)
			err = os.RemoveAll(removePath)
			if err != nil {
				panic(fmt.Sprintf("Failed removing unknown bundle: %v", err))
			}
		}
	}
}

func (u *Updater) determineBundleChanges() {
	u.announceStatus(RetrieveRemoteBundleVersions, 0)
	urls := make([]string, 0, len(u.bundleUpdateInfos))
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		urls = append(urls, bundleUpdateInfo.BundleInfoURL)
	}
	log.Infof("Downloading bundle information for bundles from these URLs: %v.", urls)
	bundleInfos, err := u.RetrieveBundleInfos(urls, u.publicKeys)
	if err != nil {
		panic(err)
	}
	for _, bundleUpdateInfo := range u.bundleUpdateInfos {
		bundleUpdateInfo.RemoteState = bundleInfos[bundleUpdateInfo.BundleInfoURL].GetFileHashes()
		bundleUpdateInfo.WantedState = config.MakeDiffFileInfoMap(bundleUpdateInfo.PresentState, bundleUpdateInfo.RemoteState)
	}
}

func (u *Updater) RetrieveBundleInfos(urls []string, publicKeys []*rsa.PublicKey) (bundleInfos map[string]*config.BundleInfo, err error) {
	bundleInfosData, err := u.downloader.DownloadSignedResources(urls, u.publicKeys)
	if err != nil {
		return nil, err
	}
	bundleInfos = make(map[string]*config.BundleInfo)
	for _, url := range urls {
		bundleInfos[url] = config.ReadInfoFromByteSlice(bundleInfosData[url])
	}
	return bundleInfos, err
}

func (u *Updater) RetrieveBundleInfo(fromURL string, publicKeys []*rsa.PublicKey) (info *config.BundleInfo, sha string) {
	bundleInfosData, err := u.downloader.DownloadSignedResources([]string{fromURL}, publicKeys)
	if err != nil {
		panic(err)
	}
	info = config.ReadInfoFromReader(strings.NewReader(string(bundleInfosData[fromURL])))
	if u.timestampFilePath != "" {
		timestamps.VerifyBundleInfoTimestamp(info.UniqueBundleName, info.Timestamp, u.timestampFilePath)
	}
	shaBytes := sha256.Sum256(bundleInfosData[fromURL])
	return info, hex.EncodeToString(shaBytes[:])
}

func (u *Updater) installBundleUpdates() {
	u.announceStatus(DownloadBundleFiles, countUpdatesBytes(u.bundleUpdateInfos))
	for _, bundleUpdateConfig := range u.bundleUpdateInfos {
		log.Infof("Downloading %d files for bundle \"%s\".", bundleUpdateConfig.WantedState.UpdateFileCount(), bundleUpdateConfig.LocalDirectory)
		bundleDirectory := filepath.Join(u.userBundlesFolderPath, bundleUpdateConfig.LocalDirectory)
		deleteChangedFiles(bundleUpdateConfig.WantedState, bundleDirectory)
		u.downloader.MustDownloadToDirectory(bundleUpdateConfig.BaseURL, bundleUpdateConfig.WantedState, bundleDirectory)
		system.MustRecursivelyRemoveEmptyFolders(bundleDirectory)
	}
}

func applyBundleUpdate(fileMap config.FileInfoMap, fromPath, toPath string) {
	sw := chronometry.NewStartedStopwatch()
	deleteChangedFiles(fileMap, toPath)
	log.Infof("Moving %d files from \"%s\" to \"%s\".", fileMap.UpdateFileCount(), fromPath, toPath)
	system.MustMoveFiles(fromPath, toPath)
	system.MustRecursivelyRemoveEmptyFolders(toPath)
	log.Infof("Applying bundle update to folder \"%s\" took %v.", toPath, sw.TakeLapTime())
}

func deleteChangedFiles(fileMap config.FileInfoMap, localDirPath string) {
	log.Infof("Deleting %d files, %d of which will be updated in \"%s\".", len(fileMap), uint64(len(fileMap))-fileMap.DeleteFileCount(), localDirPath)
	for filePath := range fileMap {
		system.MustRemoveFile(filepath.Join(localDirPath, filePath))
	}
}

func (u *Updater) DownloadBundle(fromURL string, bundleInfoURL string, publicKeys []*rsa.PublicKey, toFolder string) {
	bundleInfo, _ := u.RetrieveBundleInfo(bundleInfoURL, publicKeys)
	u.downloader.MustDownloadToDirectory(fromURL, bundleInfo.BundleFiles, toFolder)
}
