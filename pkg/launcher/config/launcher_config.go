package config

import (
	"io"

	"github.com/setlog/trivrost/pkg/misc"
)

type LauncherConfig struct {
	DeploymentConfigURL            string         `json:"DeploymentConfigURL"`
	VendorName                     string         `json:"VendorName"`
	ProductName                    string         `json:"ProductName"`
	BrandingName                   string         `json:"BrandingName"`
	BrandingNameShort              string         `json:"BrandingNameShort"`
	ReverseDnsProductId            string         `json:"ReverseDnsProductId"`
	ProductVersion                 VersionData    `json:"ProductVersion"`
	BinaryName                     string         `json:"BinaryName"`
	StatusMessages                 StatusMessages `json:"StatusMessages"`
	IgnoreLauncherBundleInfoHashes []string       `json:"IgnoreLauncherBundleInfoHashes"`
}

type StatusMessages struct {
	AcquireLock                   string `json:"AcquireLock"`
	GetDeploymentConfig           string `json:"GetDeploymentConfig"`
	DetermineLocalLauncherVersion string `json:"DetermineLocalLauncherVersion"`
	RetrieveRemoteLauncherVersion string `json:"RetrieveRemoteLauncherVersion"`
	SelfUpdate                    string `json:"SelfUpdate"`
	DetermineLocalBundleVersions  string `json:"DetermineLocalBundleVersions"`
	RetrieveRemoteBundleVersions  string `json:"RetrieveRemoteBundleVersions"`
	AwaitApplicationsTerminated   string `json:"AwaitApplicationsTerminated"`
	DownloadBundleUpdates         string `json:"DownloadBundleUpdates"`
	LaunchApplication             string `json:"LaunchApplication"`
}

type VersionData struct {
	Major int `json:"Major"`
	Minor int `json:"Minor"`
	Patch int `json:"Patch"`
	Build int `json:"Build"`
}

func ReadLauncherConfigFromReader(reader io.Reader) (launcherConfig *LauncherConfig) {
	misc.MustUnmarshalJSON(misc.MustReadAll(reader), &launcherConfig)
	return launcherConfig
}
