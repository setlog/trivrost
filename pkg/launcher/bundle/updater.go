package bundle

import (
	"context"
	"crypto/rsa"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/setlog/trivrost/pkg/fetching"
	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/timestamps"
	"github.com/setlog/trivrost/pkg/system"
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

// BundleUpdateInfo contains information on what files need updating on the user's machine for the bundle specified by the embedded BundleConfig.
type BundleUpdateInfo struct {
	config.BundleConfig
	IsSystemBundle bool
	PresentState   config.FileInfoMap
	RemoteState    config.FileInfoMap
	WantedState    config.FileInfoMap
}

type Updater struct {
	downloader       *fetching.Downloader
	deploymentConfig *config.DeploymentConfig
	publicKeys       []*rsa.PublicKey

	bundleUpdateInfos               []*BundleUpdateInfo
	ignoredSelfUpdateBundleInfoSHAs []string

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

func (u *Updater) IsTimestampVerificationEnabled() bool {
	return u.timestampFilePath == ""
}

func (u *Updater) SetStatusCallback(statusCallback func(UpdaterStatus, uint64)) {
	u.statusCallback = statusCallback
}

func (u *Updater) announceStatus(status UpdaterStatus, progressTarget uint64) {
	if u.statusCallback != nil {
		u.statusCallback(status, progressTarget)
	}
}

func (u *Updater) RetrieveDeploymentConfig(fromURL string) {
	log.Infof("Downloading deployment config from \"%s\".", fromURL)
	data, err := u.downloader.DownloadSignedResource(fromURL, u.publicKeys)
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
