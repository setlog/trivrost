package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func main() {
	filePath, skipurlcheck, skipJarCheck := parseFlags()
	data := system.MustReadFile(filePath)

	// validate json schema
	err := config.ValidateDeploymentConfig(string(data))
	if err != nil {
		fatalf("Could not validate deployment config \"%s\": %v", filePath, err)
	}

	if !skipurlcheck {
		checkURLs(data, filePath, skipJarCheck)
	}
}

func checkURLs(data []byte, filePath string, skipJarCheck bool) {
	urlMap, success := collectURLs(data, skipJarCheck)

	waitgroup := sync.WaitGroup{}
	waitgroup.Add(len(urlMap))
	var gotCertError bool

	// check all url in parallel
	var errorCount int32
	for url, details := range urlMap {
		go func(url string, details checkDetails) {
			defer waitgroup.Done()
			code, err := getUrlHeadResult(url)
			if err != nil {
				_, isUnknownAuthorityError := err.(x509.UnknownAuthorityError)
				gotCertError = gotCertError || isUnknownAuthorityError
				fmt.Printf("\033[0;91mHTTP HEAD request to URL '%s' failed: %v. (Check reason: %v)\033[0m\n", url, err, details)
				atomic.AddInt32(&errorCount, 1)
			} else if code != http.StatusOK {
				fmt.Printf("\033[0;91mHTTP HEAD request to URL '%s' yielded bad response code %d. (Check reason: %v)\033[0m\n", url, code, details)
				atomic.AddInt32(&errorCount, 1)
			} else {
				fmt.Printf("OK: Resource %s is available. (Reason for check: %v)\n", url, details)
			}
		}(url, details)
	}
	waitgroup.Wait()
	if errorCount > 0 {
		if gotCertError {
			fmt.Printf("\033[0;91mThere was at least one certificate-related error. The system's certificate pool may be out of date.\033[0m\n")
		}
		fatalf("%d out of %d tested URLs from \"%s\" do not point to valid resources.", errorCount, len(urlMap), filePath)
		success = false
	}
	if !success {
		os.Exit(1)
	}
}

func collectURLs(data []byte, skipJarCheck bool) (urlMap map[string]checkDetails, success bool) {
	urlMap = make(map[string]checkDetails)
	success = true
	for _, operatingsystem := range []string{"windows", "darwin", "linux"} {
		for _, arch := range []string{"386", "amd64"} {
			deploymentConfig := config.ParseDeploymentConfig(strings.NewReader(string(data)), operatingsystem, arch)
			for _, update := range deploymentConfig.LauncherUpdate {
				addUrlWithDetails(urlMap, update.BundleInfoURL, checkDetails{reasonUpdate, operatingsystem, arch, 0})
			}
			for _, update := range deploymentConfig.Bundles {
				addUrlWithDetails(urlMap, update.BundleInfoURL, checkDetails{reasonBundle, operatingsystem, arch, 0})
			}
			for _, command := range deploymentConfig.Execution.Commands {
				success = collectCommandURLs(urlMap, deploymentConfig, operatingsystem, arch, command, skipJarCheck) && success
			}
		}
	}
	return urlMap, success
}

func addUrlWithDetails(urlMap map[string]checkDetails, url string, details checkDetails) {
	presentDetails, ok := urlMap[url]
	if ok {
		presentDetails.othersCount++
		urlMap[url] = presentDetails
	} else {
		urlMap[url] = details
	}
}

func collectCommandURLs(urlMap map[string]checkDetails, deploymentConfig *config.DeploymentConfig, os, arch string, command config.Command, skipJarCheck bool) (success bool) {
	commandNameUnix := strings.ReplaceAll(command.Name, `\`, "/")
	if path.IsAbs(commandNameUnix) || !strings.Contains(commandNameUnix, "/") {
		return
	}
	bundleName := misc.FirstElementOfPath(commandNameUnix)
	bundleURL := getBundleURL(bundleName, deploymentConfig)
	if bundleURL == "" {
		fmt.Printf("\033[0;91mCould not get bundle URL for bundle \"%s\" for platform %s-%s. (Required for command \"%s\")\033[0m\n", bundleName, os, arch, command.Name)
		return false
	}
	binaryURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(commandNameUnix))
	if os == system.OsWindows && !strings.HasSuffix(binaryURL, ".exe") {
		binaryURL += ".exe"
	}
	addUrlWithDetails(urlMap, binaryURL, checkDetails{reasonCommand, os, arch, 0})
	if !skipJarCheck {
		if strings.HasSuffix(binaryURL, "/java.exe") || strings.HasSuffix(binaryURL, "/javaw.exe") ||
			strings.HasSuffix(binaryURL, "/java") {
			collectJarURL(urlMap, deploymentConfig, command, os, arch)
		}
	}
	return true
}

func stripFirstPathElement(s string) string {
	s = path.Clean(s)
	parts := strings.Split(s, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts[1:], "/")
}

func collectJarURL(urlMap map[string]checkDetails, deploymentConfig *config.DeploymentConfig, command config.Command, os, arch string) {
	check := false
	for _, arg := range command.Arguments {
		if check {
			jarPath := strings.ReplaceAll(arg, `\`, "/")
			bundleName := misc.FirstElementOfPath(jarPath)
			bundleURL := getBundleURL(bundleName, deploymentConfig)
			jarURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(jarPath))
			addUrlWithDetails(urlMap, jarURL, checkDetails{reasonJar, os, arch, 0})
			break
		}
		if arg == "-jar" {
			check = true
		}
	}
}

func getBundleURL(bundleName string, deploymentConfig *config.DeploymentConfig) string {
	for _, bundle := range deploymentConfig.Bundles {
		if bundle.LocalDirectory == bundleName {
			return bundle.BaseURL
		}
	}
	return ""
}

func getUrlHeadResult(url string) (responseCode int, err error) {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	var response *http.Response
	response, err = client.Head(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	return response.StatusCode, err
}

func parseFlags() (string, bool, bool) {
	skipurlcheck := flag.Bool("skipurlcheck", false, "Disable checking of availability of all URLs in the config.")
	skipJarCheck := flag.Bool("skipjarcheck", false, "Disable checking of availability of .jar files given to java with the -jar argument.")
	flag.Parse()

	if flag.NArg() != 1 {
		fatalf("Need at least one arg: deploymentConfigPath")
	}
	deploymentConfigPath := flag.Arg(0)
	if deploymentConfigPath == "" {
		fatalf("deploymentConfigPath not set")
	}

	return deploymentConfigPath, *skipurlcheck, *skipJarCheck
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf("\033[0;91mFatal: "+formatMessage+"\033[0m\n", args...)
	os.Exit(1)
}
