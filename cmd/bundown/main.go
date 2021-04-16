package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/setlog/trivrost/pkg/system"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/resources"

	"github.com/setlog/trivrost/pkg/fetching"
	"github.com/setlog/trivrost/pkg/launcher/bundle"
	"github.com/setlog/trivrost/pkg/launcher/config"
)

const (
	deploymentConfigFlag = "deployment-config"
	osFlag               = "os"
	archFlag             = "arch"
	outDirPathFlag       = "out"
	publicKeyPathFlag    = "pub"
	tagsFlag             = "tags"
)

const (
	tagUntagged = "untagged"
	tagAll      = "all"
)

func main() {
	deploymentConfigPath, os, arch, outDirPath, publicKeyPath, tags := parseFlags()
	if publicKeyPath != "" {
		resources.PublicRsaKeys = resources.ReadPublicRsaKeysAsset(string(system.MustReadFile(publicKeyPath)))
	}
	deploymentConfig := config.ParseDeploymentConfig(mustReaderForFile(deploymentConfigPath), os, arch)
	downloadBundles(deploymentConfig, outDirPath, tags)
}

func downloadBundles(deploymentConfig *config.DeploymentConfig, outDirPath string, tags []string) {
	outDirPath, err := filepath.Abs(outDirPath)
	if err != nil {
		fatalf("%v", err)
	}
	updater := bundle.NewUpdaterWithDeploymentConfig(context.Background(), deploymentConfig, &fetching.ConsoleDownloadProgressHandler{}, resources.PublicRsaKeys)
	for _, bundle := range deploymentConfig.Bundles {
		if shouldDownloadBundle(bundle.Tags, tags) {
			log.Infof("Starting download of bundle %s", bundle.BaseURL)
			bundleDirectory := filepath.Join(outDirPath, bundle.LocalDirectory)
			updater.DownloadBundle(bundle.BaseURL, bundle.BundleInfoURL, resources.PublicRsaKeys, bundleDirectory)
		} else {
			log.Infof("Not downloading bundle %s: none of its tags %v match the supplied tags %v.", bundle.BaseURL, bundle.Tags, tags)
		}
	}
}

func shouldDownloadBundle(bundleTags []string, allowedTags []string) bool {
	for _, allowedTag := range allowedTags {
		for _, bundleTag := range bundleTags {
			if bundleTag == allowedTag {
				return true
			}
		}
		if (len(bundleTags) == 0) && (allowedTag == tagUntagged) {
			return true
		}
		if allowedTag == tagAll {
			return true
		}
	}

	return false
}

func parseFlags() (string, string, string, string, string, []string) {
	deploymentConfig := flag.String(deploymentConfigFlag, "trivrost/deployment-config.json",
		"Path to a deployment-config to download bundles for.")
	os := flag.String(osFlag, "", "GOOS-style name of the operating system to download bundles for.")
	arch := flag.String(archFlag, "", "GOARCH-style name of the architecture to download bundles for.")
	outDirPath := flag.String(outDirPathFlag, "bundles", "Path to the directory to download files to. Will be created if missing.")
	publicKeyPath := flag.String(publicKeyPathFlag, "", "Path to a custom public key file to verify signatures of downloaded bundleinfo.json files. (optional)")
	tags := flag.String(tagsFlag, tagUntagged, "Only download bundles with one of these comma-separated tags. "+
		"The special tag '"+tagUntagged+"' implicitly exists on all bundles without tags. The special tag '"+tagAll+"' "+
		"will instruct bundown to download all bundles regardless of tags.")
	flag.Parse()

	if *deploymentConfig == "" {
		fatalf("Parameter --%s cannot be empty.", deploymentConfigFlag)
	}
	if *os == "" {
		fatalf("Parameter --%s cannot be empty.", osFlag)
	}
	if *arch == "" {
		fatalf("Parameter --%s cannot be empty.", archFlag)
	}
	if *outDirPath == "" {
		fatalf("Parameter --%s cannot be empty.", outDirPathFlag)
	}
	if *tags == "" {
		fatalf("Parameter --%s cannot be empty.", tagsFlag)
	}

	return *deploymentConfig, *os, *arch, *outDirPath, *publicKeyPath, strings.Split(*tags, ",")
}

func mustReaderForFile(filePath string) io.Reader {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fatalf("Could not read file \"%s\": %v", filePath, err)
	}
	return bytes.NewReader(data)
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf(formatMessage+"\n", args...)
	os.Exit(1)
}
