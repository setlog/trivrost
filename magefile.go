//+build mage

// This is the mage-based buildfile (magefile, https://magefile.org/) for trivrost.
// Project: https://setlog.github.io/trivrost/
//
// To use this magefile, install 'mage' and just run 'mage' within this directory.
package main

import (
	//"log"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"regexp"
	"runtime"
	"time"
)

const modulePathLauncher = "github.com/setlog/trivrost/cmd/launcher"
const modulePathHasher = "github.com/setlog/trivrost/cmd/hasher"
const modulePathValidator = "github.com/setlog/trivrost/cmd/validator"
const modulePathSigner = "github.com/setlog/trivrost/cmd/signer"
const outDir = "out"
const releaseFilesDir = outDir + "/release_files"
const updateFilesDir = outDir + "/update_files"

const hasherBinary = "hasher"
const validatorBinary = "validator"
const signerBinary = "signer"
// allow custom launcher name
var launcherBinary, _ = sh.Output("go", "run", "cmd/echo_field/main.go", "cmd/launcher/resources/launcher-config.json", "BinaryName")
var launcherMsiBinary = launcherBinary
var launcherBrandingName, _ = sh.Output("go", "run", "cmd/echo_field/main.go", "cmd/launcher/resources/launcher-config.json", "BrandingName")

var launcherVersion = version()
var binaryExt = ext()
var gitDesc

// Version is latest tag corresponding to version regex
var versionPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
func version() string {

	versionPattern.MatchString()
}

func ext() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	} else {
		return ""
	}
}

func

func init() {
	timestamp := time.Now().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev" + timestamp + hash
	}
	// fmt.Sprintf(`-X "github.com/magefile/mage/mage.timestamp=%s" -X "github.com/magefile/mage/mage.commitHash=%s" -X "github.com/magefile/mage/mage.gitTag=%s"`, timestamp, hash, tag)
}



var gocmd = mg.GoCmd()
var Default = Build

// Builds trivrost
func Build() {
	mg.Deps(Generate)
}

// Generates all required files
func Generate() {
	println("a " + launcherProgramName)
	println("b " + launcherProgramExt)
	println("c " + launcherBrandingName)
}

// Cleans up generated files
func Clean() error {
	if len(outDir) == 0 {
		return mg.Fatal(1, "Output directory not set.")
	}
	_ = sh.Rm(outDir)
	_ = sh.Rm("cmd/launcher/resources/*.gen.go")
	_ = sh.Rm("cmd/launcher/*.syso")
	_ = sh.Run(gocmd, "clean", modulePathLauncher)
	_ = sh.Run(gocmd, "clean", modulePathHasher)
	return nil
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}
