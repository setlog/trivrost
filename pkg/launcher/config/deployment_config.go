package config

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/misc"
)

type DeploymentConfig struct {
	Timestamp      string                 `json:"Timestamp,omitempty"`
	LauncherUpdate []LauncherUpdateConfig `json:"LauncherUpdate,omitempty"`
	Bundles        []BundleConfig         `json:"Bundles,omitempty"`
	Execution      ExecutionConfig        `json:"Execution,omitempty"`
}

type HashDataConfig struct {
	BundleInfoURL     string `json:"BundleInfoURL"`
	BaseURL           string `json:"BaseURL,omitempty"`
	IsUpdateMandatory bool   `json:"IsUpdateMandatory,omitempty"`
}

type LauncherUpdateConfig struct {
	HashDataConfig
	TargetPlatforms []string `json:"TargetPlatforms,omitempty"`
}

type BundleConfig struct {
	HashDataConfig
	LocalDirectory  string   `json:"LocalDirectory"`
	TargetPlatforms []string `json:"TargetPlatforms,omitempty"`
	Tags            []string `json:"Tags,omitempty"`
}

type ExecutionConfig struct {
	Commands               []Command `json:"Commands,omitempty"`
	LingerTimeMilliseconds int       `json:"LingerTimeMilliseconds,omitempty"`
}

type Command struct {
	Name                       string             `json:"Name"`
	WorkingDirectoryBundleName string             `json:"WorkingDirectoryBundleName,omitempty"`
	Arguments                  []string           `json:"Arguments,omitempty"`
	Env                        map[string]*string `json:"Env,omitempty"`
	TargetPlatforms            []string           `json:"TargetPlatforms,omitempty"`
}

func (dc *DeploymentConfig) HasLauncherUpdateConfig() bool {
	return len(dc.LauncherUpdate) == 1
}

func (dc *DeploymentConfig) GetLauncherUpdateConfig() *LauncherUpdateConfig {
	if dc.HasLauncherUpdateConfig() {
		return &dc.LauncherUpdate[0]
	}
	return nil
}

func ReadDeploymentConfig(reader io.Reader, os string, arch string) (string, error) {
	return expandPlaceholders(string(misc.MustReadAll(reader)), os, arch)
}

func ParseDeploymentConfig(reader io.Reader, os string, arch string) (deploymentConfig *DeploymentConfig) {
	data, err := ReadDeploymentConfig(reader, os, arch)
	if err != nil {
		panic(err)
	}
	misc.MustUnmarshalJSON([]byte(data), &deploymentConfig)
	deploymentConfig.LauncherUpdate = FilterLauncherUpdatesByPlatform(deploymentConfig.LauncherUpdate, os, arch)
	deploymentConfig.Bundles = FilterBundlesByPlatform(deploymentConfig.Bundles, os, arch)
	deploymentConfig.Execution.Commands = FilterCommandsByPlatform(deploymentConfig.Execution.Commands, os, arch)
	configureLauncherUpdates(deploymentConfig.LauncherUpdate)
	configureBundles(deploymentConfig.Bundles)
	return deploymentConfig
}

func configureLauncherUpdates(launchers []LauncherUpdateConfig) {
	launcherCount := len(launchers)
	for i := 0; i < launcherCount; i++ {
		configureLauncherUpdate(&launchers[i])
	}
}

func configureLauncherUpdate(launcher *LauncherUpdateConfig) {
	if launcher.BaseURL == "" {
		launcher.BaseURL = misc.MustStripLastURLPathElement(launcher.BundleInfoURL)
		log.Debugf("BaseURL for launcher for platforms %v was empty. Deriving it from BundleInfoURL \"%s\": \"%s\".",
			launcher.TargetPlatforms, launcher.BundleInfoURL, launcher.BaseURL)
	}
}

func configureBundles(bundles []BundleConfig) {
	bundleCount := len(bundles)
	for i := 0; i < bundleCount; i++ {
		configureBundle(&bundles[i])
	}
}

func configureBundle(bundle *BundleConfig) {
	if bundle.BaseURL == "" {
		bundle.BaseURL = misc.MustStripLastURLPathElement(bundle.BundleInfoURL)
		log.Debugf("BaseURL for bundle \"%s\" was empty. Deriving it from BundleInfoURL \"%s\": \"%s\".",
			bundle.LocalDirectory, bundle.BundleInfoURL, bundle.BaseURL)
	}

	initialLocalDirectoryValue := bundle.LocalDirectory
	if strings.HasPrefix(bundle.LocalDirectory, `/`) {
		panic(fmt.Sprintf(`Bundle with "LocalDirectory"-value of "%s" is invalid: cannot use an absolute path.`, initialLocalDirectoryValue))
	}
	bundle.LocalDirectory = strings.Trim(bundle.LocalDirectory, `/\`)
	if strings.Contains(bundle.LocalDirectory, `/`) || strings.Contains(bundle.LocalDirectory, `\`) {
		panic(fmt.Sprintf(`Bundle with "LocalDirectory"-value of "%s" is invalid: must not have '/' or '\' within name.`, initialLocalDirectoryValue))
	}
}
