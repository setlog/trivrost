package bundle

import (
	"context"
	"crypto/rsa"
	"runtime"
	"strings"

	"github.com/setlog/trivrost/pkg/fetching"
	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/timestamps"
	"github.com/setlog/trivrost/pkg/system"
	log "github.com/sirupsen/logrus"
)

type UpdaterStatus int

const (
	DetermineLocalLauncherVersion UpdaterStatus = iota
	RetrieveRemoteLauncherVersion
	DownloadLauncherFiles
	DetermineLocalBundleVersions
	RetrieveRemoteBundleVersions
	DownloadBundleFiles
)

type Updater struct {
	downloader       *fetching.Downloader
	deploymentConfig *config.DeploymentConfig
	publicKeys       []*rsa.PublicKey

	bundleUpdateInfos                   []*BundleUpdateInfo
	ignoredLauncherUpdateBundleInfoSHAs []string

	userBundlesFolderPath   string
	systemBundlesFolderPath string

	timestampFilePath string

	statusCallback func(UpdaterStatus, uint64)

	ctx context.Context
}

func NewUpdater(ctx context.Context, dlHandler fetching.DownloadProgressHandler, publicKeys []*rsa.PublicKey) *Updater {
	return &Updater{ctx: ctx, downloader: fetching.NewDownloader(ctx, dlHandler), publicKeys: publicKeys}
}

func NewUpdaterWithDeploymentConfig(ctx context.Context, deploymentConfig *config.DeploymentConfig, dlHandler fetching.DownloadProgressHandler, publicKeys []*rsa.PublicKey) *Updater {
	return &Updater{ctx: ctx, downloader: fetching.NewDownloader(ctx, dlHandler), publicKeys: publicKeys, deploymentConfig: deploymentConfig}
}

func (u *Updater) EnableTimestampVerification(filePath string) {
	u.timestampFilePath = filePath
}

func (u *Updater) DisableTimestampVerification() {
	u.timestampFilePath = ""
}

func (u *Updater) SetStatusCallback(statusCallback func(UpdaterStatus, uint64)) {
	u.statusCallback = statusCallback
}

func (u *Updater) announceStatus(status UpdaterStatus, progressTarget uint64) {
	if u.statusCallback != nil {
		u.statusCallback(status, progressTarget)
	}
}

func (u *Updater) ObtainDeploymentConfig(deploymentConfigURL string) {
	log.Infof("Obtaining deployment config from \"%s\".", deploymentConfigURL)
	data, err := u.downloader.DownloadSignedResource(deploymentConfigURL, u.publicKeys)
	if err != nil {
		panic(err)
	}

	deploymentConfig := config.ParseDeploymentConfig(strings.NewReader(string(data)), runtime.GOOS, system.GetOSArch())
	if u.timestampFilePath != "" {
		timestamps.VerifyDeploymentConfigTimestamp(deploymentConfig.Timestamp, u.timestampFilePath)
	}
	u.deploymentConfig = deploymentConfig
}

func (u *Updater) GetDeploymentConfig() *config.DeploymentConfig {
	return u.deploymentConfig
}
