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
	deploymentConfigFlag   = "deployment-config"
	osFlag                 = "os"
	archFlag               = "arch"
	outDirPathFlag         = "out"
	publicKeyPathFlag      = "pub"
	tagsFlag               = "tags"
	skipPresentBundlesFlag = "skip-present-bundles"
)

const (
	tagUntagged = "untagged"
	tagAll      = "all"
)

type Flags struct {
	deploymentConfigPath string
	os                   string
	arch                 string
	outDirPath           string
	publicKeyPath        string
	tags                 []string
	skipPresentBundles   bool
}

func main() {
	flags := parseFlags()
	if flags.publicKeyPath != "" {
		resources.PublicRsaKeys = resources.ReadPublicRsaKeysAsset(string(system.MustReadFile(flags.publicKeyPath)))
	}
	deploymentConfig := config.ParseDeploymentConfig(mustReaderForFile(flags.deploymentConfigPath), flags.os, flags.arch)
	downloadBundles(deploymentConfig, flags.outDirPath, flags.tags, flags.skipPresentBundles)
}

func downloadBundles(deploymentConfig *config.DeploymentConfig, outDirPath string, tags []string, skipPresentBundles bool) {
	outDirPathAbs, err := filepath.Abs(outDirPath)
	if err != nil {
		fatalf("%v", err)
	}
	updater := bundle.NewUpdaterWithDeploymentConfig(context.Background(), deploymentConfig, &fetching.ConsoleDownloadProgressHandler{}, resources.PublicRsaKeys)
	for _, bundle := range deploymentConfig.Bundles {
		if shouldDownloadBundle(bundle.Tags, tags) {
			if skipPresentBundles && isFolder(filepath.Join(outDirPathAbs, bundle.LocalDirectory)) {
				log.Infof("Not downloading bundle %s: excluded by --%s because \"%s\" already exists.", bundle.BaseURL, skipPresentBundlesFlag, filepath.Join(outDirPathAbs, bundle.LocalDirectory))
			} else {
				log.Infof("Starting download of bundle %s", bundle.BaseURL)
				bundleDirectory := filepath.Join(outDirPathAbs, bundle.LocalDirectory)
				updater.DownloadBundle(bundle.BaseURL, bundle.BundleInfoURL, resources.PublicRsaKeys, bundleDirectory)
			}
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

func parseFlags() *Flags {
	deploymentConfig := flag.String(deploymentConfigFlag, "trivrost/deployment-config.json",
		"Path to a deployment-config to download bundles for.")
	os := flag.String(osFlag, "", "GOOS-style name of the operating system to download bundles for.")
	arch := flag.String(archFlag, "", "GOARCH-style name of the architecture to download bundles for.")
	outDirPath := flag.String(outDirPathFlag, "bundles", "Path to the directory to download files to. Will be created if missing.")
	publicKeyPath := flag.String(publicKeyPathFlag, "", "Path to a custom public key file to verify signatures of downloaded bundleinfo.json files. (optional)")
	tags := flag.String(tagsFlag, tagUntagged, "Only download bundles with one of these comma-separated tags. "+
		"The special tag '"+tagUntagged+"' implicitly exists on all bundles without tags. The special tag '"+tagAll+"' "+
		"will instruct bundown to download all bundles regardless of tags.")
	skipPresentBundles := flag.Bool(skipPresentBundlesFlag, false, "If set, skip download of bundles the corresponding directory of which already exist under --out.")
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

	return &Flags{
		deploymentConfigPath: *deploymentConfig,
		os:                   *os,
		arch:                 *arch,
		outDirPath:           *outDirPath,
		publicKeyPath:        *publicKeyPath,
		tags:                 strings.Split(*tags, ","),
		skipPresentBundles:   *skipPresentBundles,
	}
}

func isFolder(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		fatalf("Could not check if \"%s\" is folder: %v", filePath, err)
	}
	return info.IsDir()
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
