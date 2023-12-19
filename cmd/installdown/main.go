package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
	"github.com/sirupsen/logrus"
)

const (
	launcherConfigFlag    = "launcher-config"
	componentGroupDirFlag = "componentgroupdir"
	componentGroupIdFlag  = "componentgroupid"
	msiOutputFileFlag     = "msioutputfile"
	launcherVersionFlag   = "launcherversion"
	wxsTemplateFlag       = "wxstemplate"
	archFlag              = "arch"
	outDirFlag            = "out"
)

const archWin64 = "x64"
const archGo64 = "amd64"
const archWin32 = "x86"
const archGo32 = "386"

var validVersionRegex = regexp.MustCompile(`v?([0-9]+\.[0-9]+\.[0-9]+).*`)

type wxsConfig struct {
	VendorName        string
	ProductName       string
	BrandingName      string
	BinaryName        string
	ComponentGroupDir string
	ComponentGroupId  string
	GuidUpgradeCode   string
	MsiOutputFile     string
	LauncherVersion   string
	WxsTemplatePath   string
	Arch              string
	OutDir            string
}

func main() {
	cfg := configure()

	logrus.Info("Building WIX script")
	createWXSFile(cfg)

	logrus.Info("Harvesting components for " + cfg.Arch)
	generateComponentGroupsFile(cfg)

	logrus.Info("Generating MSI for " + cfg.Arch)
	compileMsi(cfg)

	logrus.Info("Done")
}

func compileMsi(cfg *wxsConfig) {
	// -arch run for specified arch
	logrus.Info("Running candle to build component wixobj")
	mustRunCommand("candle",
		filepath.Join(cfg.OutDir, cfg.ComponentGroupId+cfg.Arch+".wxs"),
		"-o", filepath.Join(cfg.OutDir, cfg.ComponentGroupId+cfg.Arch+".wixobj"))
	logrus.Info("Running candle to build launcher.wixobj")
	mustRunCommand("candle",
		filepath.Join(cfg.OutDir, "launcher.wxs"),
		"-o", filepath.Join(cfg.OutDir, "launcher wixobj"),
		"-arch", cfg.Arch)

	// -sice:ICE61 suppress ICE61, the warning about same-version-upgrade which we need to allow for updating bundled bundles
	//  without updating the launcher
	// -sacl suppress ACL warning
	logrus.Info("Running light to compile final msi")
	mustRunCommand("light",
		"-sice:ICE61",
		"-sacl",
		"-sval",
		"-ext WixUIExtension",
		"-b", cfg.ComponentGroupDir,
		"-out", filepath.Join(cfg.OutDir, cfg.MsiOutputFile),
		filepath.Join(cfg.OutDir, cfg.ComponentGroupId+cfg.Arch+".wixobj"),
		filepath.Join(cfg.OutDir, "launcher.wixobj"))
}

func generateComponentGroupsFile(cfg *wxsConfig) {
	// need to exclude the binary because it is part of the wxs file already. Moving it to the side temporarily
	// which is an utterly stupid way to do it. FIXME
	dir, err := ioutil.TempDir(cfg.OutDir, "launcher")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir) // clean up

	// move main binary away
	src, temp := filepath.Join(cfg.ComponentGroupDir, cfg.BinaryName+".exe"), filepath.Join(dir, cfg.BinaryName+".exe")
	err = os.Rename(src, temp)
	if err != nil {
		panic(err)
	}

	// generate the fragment for the specified bundle
	//  -cg: the name of the componentgroup = bundle name
	//  -srd: suppress harvesting the root directory, do not generate a component for the root directory
	//  -dr: directory reference, the directory under which these components should go
	//  -gg: generate GUIDs for the components immediately. Same input generate same GUIDs
	//  -platform: sets the target platform, but apparently not used when generating componentgroups
	if cfg.Arch == archWin64 {
		mustRunCommand("heat",
			"dir", cfg.ComponentGroupDir,
			"-cg", cfg.ComponentGroupId,
			"-srd",
			"-dr", "APPLICATIONROOTDIRECTORY",
			"-t", "build/HeatTransform.xslt",
			"-gg",
			"-out", filepath.Join(cfg.OutDir, cfg.ComponentGroupId+cfg.Arch+".wxs"))
	} else {
		mustRunCommand("heat",
			"dir", cfg.ComponentGroupDir,
			"-cg", cfg.ComponentGroupId,
			"-srd",
			"-dr", "APPLICATIONROOTDIRECTORY",
			"-gg",
			"-out", filepath.Join(cfg.OutDir, cfg.ComponentGroupId+cfg.Arch+".wxs"))
	}

	// move main binary back
	err = os.Rename(temp, src)
	if err != nil {
		panic(err)
	}
}

func createWXSFile(cfg *wxsConfig) {
	wxsTemplate, err := template.New("wxsTempalte").Parse(string(system.MustReadFile(cfg.WxsTemplatePath)))
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.OutDir)
	finalWxsFile, err := os.OpenFile(filepath.Join(cfg.OutDir, "launcher.wxs"),
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0770)
	if err != nil {
		panic(err)
	}
	defer finalWxsFile.Close()

	err = wxsTemplate.Execute(finalWxsFile, cfg)
	if err != nil {
		panic(err)
	}
}

func generateUpgradeCode(vendorName string, productName string) string {
	// GUID v4 (random), variant 8 generation using random part in the middle to get unique GUID for this config:
	//  X = SHA256(VendorName + ProductName)
	//  GUID = X[1-8]-X[9-12]-4X[13-15]-bX[16-18]-X[19-30]

	shasumbytes := sha256.Sum256([]byte(vendorName + productName))
	x := hex.EncodeToString(shasumbytes[:])
	return x[0:8] + "-" + x[8:12] + "-4" + x[12:15] + "-b" + x[15:18] + "-" + x[18:30]
}

func configure() *wxsConfig {
	launcherConfigPath := flag.String(launcherConfigFlag,
		"cmd/launcher/resources/launcher-config.json", "Path to a launcher-config for configuration.")
	componentGroupDir := flag.String(componentGroupDirFlag,
		"", "The directory which holds the componentgroup for which the installer is created.")
	componentGroupId := flag.String(componentGroupIdFlag,
		"_components", "The internal id of the componentgroup for which the installer is created.")
	msiOutputFile := flag.String(msiOutputFileFlag,
		"install.msi", "Name of the resulting MSI file.")
	wxsTemplate := flag.String(wxsTemplateFlag,
		"", "Template used for wix installer.")
	launcherVersion := flag.String(launcherVersionFlag,
		"", "Version to embed in the installer. Leading 'v' char will be stripped automatically.")
	arch := flag.String(archFlag,
		archWin32, "Which arch to build for. Either x86/368(default) or x64/amd64.")
	outDir := flag.String(outDirFlag,
		"out", "Output directory.")
	flag.Parse()

	if *launcherConfigPath == "" {
		fatalf("Parameter --%s cannot be empty.", launcherConfigFlag)
	}
	if *componentGroupDir == "" {
		fatalf("Parameter --%s cannot be empty.", componentGroupDirFlag)
	}
	if *componentGroupId == "" {
		fatalf("Parameter --%s cannot be empty.", componentGroupIdFlag)
	}
	if *wxsTemplate == "" {
		fatalf("Parameter --%s cannot be empty.", wxsTemplateFlag)
	}

	if *arch != archGo32 && *arch != archWin32 && *arch != archGo64 && *arch != archWin64 {
		fatalf("Parameter --%s must be either x86, 386 or x64, amd64.", archFlag)
	}
	if *arch == archGo32 {
		*arch = archWin32
	}
	if *arch == archGo64 {
		*arch = archWin64
	}

	versionMatch := validVersionRegex.FindAllStringSubmatch(*launcherVersion, -1)
	if versionMatch == nil {
		fatalf("Parameter --%s must be a version in the format %s. Found: %s", launcherVersionFlag, validVersionRegex.String(), *launcherVersion)
	}
	*launcherVersion = versionMatch[0][1]
	fmt.Printf("Using version %s for MSI.\n", *launcherVersion)

	if *outDir == "" {
		fatalf("Parameter --%s cannot be empty.", outDirFlag)
	}

	launcherConfig := config.ReadLauncherConfigFromReader(mustReaderForFile(*launcherConfigPath))
	guidUpgradeCode := generateUpgradeCode(launcherConfig.VendorName, launcherConfig.ProductName)

	return &wxsConfig{
		VendorName:        launcherConfig.VendorName,
		ProductName:       launcherConfig.ProductName,
		BinaryName:        launcherConfig.BinaryName,
		BrandingName:      launcherConfig.BrandingName,
		ComponentGroupDir: *componentGroupDir,
		ComponentGroupId:  *componentGroupId,
		GuidUpgradeCode:   guidUpgradeCode,
		MsiOutputFile:     *msiOutputFile,
		LauncherVersion:   *launcherVersion,
		WxsTemplatePath:   *wxsTemplate,
		Arch:              *arch,
		OutDir:            *outDir,
	}
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

func mustRunCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
