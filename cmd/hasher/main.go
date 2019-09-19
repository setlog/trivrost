package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/hashing"
)

const timeFormat = "2006-01-02 15:04:05"
const filename = "/bundleinfo.json"

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("Hasher expects exactly two parameters.")
		fmt.Println("The first parameter is the unique bundle name.")
		fmt.Println("The second parameter is the path to the directory to hash.")

		log.Info("Wrong number of arguments for hasher. Stopping.")

		os.Exit(1)
	}

	uniqueBundleName := flag.Arg(0)
	pathToHash := flag.Arg(1)
	hashesFile := filepath.Join(pathToHash, filename)
	mustHashDirectory(uniqueBundleName, pathToHash, hashesFile)

	log.Info("Finished hasher.")
}

func mustHashDirectory(uniqueBundleName, pathToHash, hashesFile string) {
	log.WithFields(log.Fields{"uniqueBundleName": uniqueBundleName, "pathToHash": pathToHash, "hashesFile": hashesFile}).Info("Hashing directory.")
	bundleInfo := &config.BundleInfo{
		BundleFiles:      hashing.MustHash(pathToHash),
		Timestamp:        time.Now().UTC().Format(timeFormat),
		UniqueBundleName: uniqueBundleName,
	}
	config.WriteInfo(bundleInfo, hashesFile)
}
