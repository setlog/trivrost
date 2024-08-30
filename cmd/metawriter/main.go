package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

var (
	launcherConfigPath                               string
	versionInfoTemplatePath, versionInfoPath         string
	mainExeManifestTemplatePath, mainExeManifestPath string
	infoPListTemplatePath, infoPListPath             string
)

var (
	launcherConfig               *config.LauncherConfig
	versionSemantic, versionFull string
)

func main() {
	parseFlags()
	determineVariables()
	validateVariables()
	writeMetaFiles()
}

func writeMetaFiles() {
	mustWriteFile(versionInfoPath, []byte(replacePlaceholders(mustReadFile(versionInfoTemplatePath))))
	mustWriteFile(mainExeManifestPath, []byte(replacePlaceholders(mustReadFile(mainExeManifestTemplatePath))))
	mustWriteFile(infoPListPath, []byte(replacePlaceholders(mustReadFile(infoPListTemplatePath))))
}

func replacePlaceholders(text string) string {
	text = strings.Replace(text, "${LAUNCHER_BINARY}", html.EscapeString(os.Getenv("LAUNCHER_BINARY"))), -1)
	text = strings.Replace(text, "${LAUNCHER_BINARY_EXT}", html.EscapeString(os.Getenv("LAUNCHER_BINARY_EXT")), -1)
	text = strings.Replace(text, "${LAUNCHER_VENDOR_NAME}", html.EscapeString(launcherConfig.VendorName), -1)
	text = strings.Replace(text, "${LAUNCHER_BRANDING_NAME}", html.EscapeString(launcherConfig.BrandingName), -1)
	text = strings.Replace(text, "${LAUNCHER_BRANDING_NAME_SHORT}", html.EscapeString(launcherConfig.BrandingNameShort), -1)
	text = strings.Replace(text, "${LAUNCHER_REVERSE_DNS_PRODUCT_ID}", html.EscapeString(launcherConfig.ReverseDnsProductId), -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_MAJOR}", strconv.Itoa(html.EscapeString(launcherConfig.ProductVersion.Major)), -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_MINOR}", strconv.Itoa(html.EscapeString(launcherConfig.ProductVersion.Minor)), -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_PATCH}", strconv.Itoa(html.EscapeString(launcherConfig.ProductVersion.Patch)), -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_BUILD}", strconv.Itoa(html.EscapeString(launcherConfig.ProductVersion.Build)), -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_SEMANTIC}", versionSemantic, -1)
	text = strings.Replace(text, "${LAUNCHER_VERSION_FULL}", versionFull, -1)
	return text
}

func mustReadFile(filePath string) string {
	fmt.Printf("Metawriter: Reading \"%s\".\n", filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fatalf("Could not read \"%s\": %v", filePath, err)
	}
	return string(data)
}

func mustWriteFile(filePath string, data []byte) {
	fmt.Printf("Metawriter: Writing \"%s\".\n", filePath)
	err := ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		fatalf("Could not open file \"%s\" for writing: %v", filePath, err)
	}
}

func determineVariables() {
	launcherConfig = config.ReadLauncherConfigFromReader(mustReaderForFile(launcherConfigPath))
	versionSemantic = fmt.Sprintf("%d.%d.%d", launcherConfig.ProductVersion.Major, launcherConfig.ProductVersion.Minor, launcherConfig.ProductVersion.Patch)
	versionFull = fmt.Sprintf("%s.%d", versionSemantic, launcherConfig.ProductVersion.Build)
}

func validateVariables() {
	if launcherConfig.DeploymentConfigURL == "" {
		fatalf("'DeploymentConfigURL' is not set in the launcher config.")
	}
	if launcherConfig.VendorName == "" {
		fatalf("'VendorName' is not set in the launcher config.")
	}
	if launcherConfig.ProductName == "" {
		fatalf("'ProductName' is not set in the launcher config.")
	}
	if launcherConfig.BrandingName == "" {
		fatalf("'BrandingName' is not set in the launcher config.")
	}
	if launcherConfig.BrandingNameShort == "" {
		fatalf("'BrandingNameShort' is not set in the launcher config.")
	}
	if len(launcherConfig.BrandingNameShort) > 15 {
		fatalf("'BrandingNameShort' in the launcher config is longer than 15 bytes.")
	}
	if launcherConfig.ReverseDnsProductId == "" {
		fatalf("'ReverseDnsProductId' is not set in the launcher config.")
	}
	if launcherConfig.BinaryName == "" {
		fatalf("'BinaryName' is not set in the launcher config.")
	}
	if versionFull == "0.0.0.0" {
		fmt.Println("Warning: 'ProductVersion' is not set or is '0.0.0.0' in the launcher config. This is not fatal but users might think it looks strange.")
	}
}

func mustReaderForFile(filePath string) io.Reader {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fatalf("Could not read file \"%s\": %v", filePath, err)
	}
	return bytes.NewReader(data)
}

func parseFlags() {
	flag.Parse()
	if flag.NArg() != 7 {
		fatalf("Need 7 args: launcherConfigPath versionInfoTemplatePath versionInfoPath mainExeManifestTemplatePath mainExeManifestPath infoPListTemplatePath infoPListPath")
	}
	launcherConfigPath = flag.Arg(0)
	versionInfoTemplatePath, versionInfoPath = flag.Arg(1), flag.Arg(2)
	mainExeManifestTemplatePath, mainExeManifestPath = flag.Arg(3), flag.Arg(4)
	infoPListTemplatePath, infoPListPath = flag.Arg(5), flag.Arg(6)
	if launcherConfigPath == "" {
		fatalf("launcher config path (1st arg) empty")
	}
	if versionInfoTemplatePath == "" {
		fatalf("version info template path (2nd arg) empty")
	}
	if versionInfoPath == "" {
		fatalf("version info path (3rd arg) empty")
	}
	if mainExeManifestTemplatePath == "" {
		fatalf("main exe manifest template path (4th arg) empty")
	}
	if mainExeManifestPath == "" {
		fatalf("main exe manifest path (5th arg) empty")
	}
	if infoPListTemplatePath == "" {
		fatalf("info plist template path (6th arg) empty")
	}
	if infoPListPath == "" {
		fatalf("info plist path (7th arg) empty")
	}
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf("Fatal: "+formatMessage+"\n", args...)
	os.Exit(1)
}
