package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

var (
	launcherConfigPath string
	keyOfValue         string
)

func main() {
	parseFlags()
	launcherConfig := config.ReadLauncherConfigFromReader(mustReaderForFile(launcherConfigPath))
	if keyOfValue == "BinaryName" {
		fmt.Print(launcherConfig.BinaryName)
	} else if keyOfValue == "BrandingName" {
		fmt.Print(launcherConfig.BrandingName)
	} else {
		fatalf("Unknown launcher-config key \"%s\".", keyOfValue)
	}
}

func parseFlags() {
	flag.Parse()
	if flag.NArg() != 2 {
		fatalf("Need 2 args: launcherConfigPath keyOfValue")
	}
	launcherConfigPath = flag.Arg(0)
	keyOfValue = flag.Arg(1)
}

func mustReaderForFile(filePath string) io.Reader {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fatalf("Could not read file \"%s\": %v", filePath, err)
	}
	return bytes.NewReader(data)
}

func fatalf(formatMessage string, args ...any) {
	fmt.Printf(formatMessage+"\n", args...)
	os.Exit(1)
}
