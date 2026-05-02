package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/launcher/hashing"
	"github.com/sirupsen/logrus"
)

const timeFormat = "2006-01-02 15:04:05"
const bundleFileName = "bundleinfo.json"

type discardingWriter struct{}

func (dw discardingWriter) Write(p []byte) (n int, err error) { return len(p), nil }

func main() {
	log.SetFlags(0)
	logrus.SetOutput(discardingWriter{})
	var verifyOnly, overwrite bool
	var absentFilePathsString string
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite existing "+bundleFileName+" if present.")
	flag.BoolVar(&verifyOnly, "verify", false, "Verify only")
	flag.StringVar(&absentFilePathsString, "absent", "", "Comma-separated list of disk files to treat as absent. Escape with ',,'")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Println("Hasher expects exactly two parameters.")
		log.Println("The first parameter is the unique bundle name.")
		log.Println("The second parameter is the path to the directory to hash.")
		log.Println("Additional flags:")
		log.Println("  -overwrite                Overwrite existing " + bundleFileName + " if present.")
		log.Println("  -verify                   Verify that state recorded in bundle info file matches state on disk and unique bundle name matches.")
		log.Println("  -absent file1,file2,...   Comma-separated list of disk files to treat as absent. Escape with ',,'.")
		log.Fatal("Wrong number of arguments for hasher. Stopping.")
	}
	absentFilePaths := stringSplitDoubleSepEscapable(absentFilePathsString, ',')

	uniqueBundleName := flag.Arg(0)
	pathToHash := flag.Arg(1)
	bundleFilePath := filepath.Join(pathToHash, bundleFileName)
	if verifyOnly {
		mustVerifyDirectory(uniqueBundleName, pathToHash, bundleFilePath, absentFilePaths)
	} else {
		mustHashDirectory(uniqueBundleName, pathToHash, bundleFilePath, overwrite)
	}
}

func stringSplitDoubleSepEscapable(str string, sep rune) (res []string) {
	var ss strings.Builder
	checkTerminate := false
	for _, r := range str {
		if r == sep {
			checkTerminate = !checkTerminate
			if checkTerminate {
				continue
			}
		} else if checkTerminate {
			if ss.Len() > 0 {
				res = append(res, ss.String())
			}
			ss.Reset()
			checkTerminate = false
		}
		ss.WriteRune(r)
	}
	if ss.Len() > 0 {
		res = append(res, ss.String())
	}
	return res
}

func mustVerifyDirectory(uniqueBundleName, pathToHash, bundleFilePath string, absentFilePaths []string) {
	pathInfo, err := os.Stat(pathToHash)
	if err != nil {
		log.Fatalf("Verification of %#q failed: cannot stat %#q: %v.\n", bundleFilePath, pathToHash, err)
	}
	if !pathInfo.IsDir() {
		log.Fatalf("Verification of %#q failed: \"%s\" must be a directory.\n", bundleFilePath, pathToHash)
	}
	fileInfoDisk := hashing.MustHash(context.Background(), pathToHash).WithForwardSlashes()
	for _, absentFilePath := range absentFilePaths {
		delete(fileInfoDisk, filepath.ToSlash(absentFilePath))
	}
	delete(fileInfoDisk, bundleFileName)

	bundleInfoRecorded := config.ReadInfo(bundleFilePath)
	fileInfoRecorded := bundleInfoRecorded.BundleFiles

	if bundleInfoRecorded.UniqueBundleName != uniqueBundleName {
		log.Fatalf("Verification of %#q failed: unique bundle name of %#q presents as %#q. Expected  %#q.\n", bundleFilePath, bundleFilePath, bundleInfoRecorded.UniqueBundleName, uniqueBundleName)
	}

	diffToDisk := config.MakeDiffFileInfoMap(fileInfoRecorded, fileInfoDisk)
	if diffToDisk.HasChanges() {
		log.Fatalf("Verification of %#q failed: files on disk have %d changes:\n", bundleFilePath, len(diffToDisk))
		for diffFilePath, diffFileInfo := range diffToDisk {
			if diffFileInfo.SHA256 == "" {
				log.Printf("%#q: %s on disk but absent in %#q.\n", diffFilePath, diffFileInfo.SHA256, bundleFilePath)
			} else {
				log.Printf("%#q: %s on disk but %s in %#q.\n", diffFilePath, diffFileInfo.SHA256, fileInfoRecorded[diffFilePath].GetSHA256(), bundleFilePath)
			}
		}
	}
}

func mustHashDirectory(uniqueBundleName, pathToHash, bundleFilePath string, overwrite bool) {
	log.Printf("Hashing directory %#q for bundle %#q.\n", pathToHash, uniqueBundleName)
	pathInfo, err := os.Stat(pathToHash)
	if err != nil {
		log.Fatalf("Cannot hash \"%s\". %s\n", pathToHash, err)
	}
	if !pathInfo.IsDir() {
		log.Fatalf("\"%s\" must be a directory.\n", pathToHash)
	}
	bundleInfo := &config.BundleInfo{
		BundleFiles:      hashing.MustHash(context.Background(), pathToHash).WithForwardSlashes(),
		Timestamp:        time.Now().UTC().Format(timeFormat),
		UniqueBundleName: uniqueBundleName,
	}
	if bundleInfo.BundleFiles[bundleFileName] != nil && !overwrite {
		log.Fatalf("Found existing \"%s\". Aborting.\n", bundleFilePath)
	}
	delete(bundleInfo.BundleFiles, bundleFileName)
	if len(bundleInfo.BundleFiles) == 0 {
		log.Fatalf("No files to hash at %v\n", pathToHash)
	}
	log.Printf("Writing %#q.\n", bundleFilePath)
	config.WriteInfo(bundleInfo, bundleFilePath)
}
