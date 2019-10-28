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

type platformChecks []platformCheck

type platformCheck struct {
	reason   checkReason
	platform config.Platform
}

func (pChecks platformChecks) String() string {
	if len(pChecks) > 1 {
		if len(pChecks) > 2 {
			return fmt.Sprintf("%s on platform %v and %d others", pChecks[0].reason, pChecks[0].platform, len(pChecks)-1)
		}
		return fmt.Sprintf("%s on platform %v and one other", pChecks[0].reason, pChecks[0].platform)
	}
	return fmt.Sprintf("%s on platform %v", pChecks[0].reason, pChecks[0].platform)
}

type checkReason int

const (
	reasonUpdate checkReason = iota
	reasonBundle
	reasonCommand
	reasonJar
)

func (cr checkReason) String() string {
	switch cr {
	case reasonUpdate:
		return "URL required for self-update"
	case reasonBundle:
		return "URL required for bundle-update"
	case reasonCommand:
		return "URL required for command binary"
	case reasonJar:
		return "URL required for Java application .jar"
	}
	panic(fmt.Sprintf("Unknown checkReason %d", cr))
}

func main() {
	filePath, skipurlcheck, skipJarCheck := parseFlags()
	data := system.MustReadFile(filePath)

	// validate json schema
	err := config.ValidateDeploymentConfig(string(data))
	if err != nil {
		fatalf("Could not validate deployment config \"%s\": %v", filePath, err)
	}

	if !skipurlcheck {
		urlMap, err := collectURLs(data, skipJarCheck)
		if err != nil {
			fatalf(err.Error())
		}
		err = checkURLs(urlMap, filePath, skipJarCheck)
		if err != nil {
			fatalf(err.Error())
		}
	}
}

func checkURLs(urlMap map[string]platformChecks, filePath string, skipJarCheck bool) error {
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(len(urlMap))
	var gotCertError bool

	// check all url in parallel
	var errorCount int32
	for url, pChecks := range urlMap {
		go func(url string, pChecks platformChecks) {
			defer waitgroup.Done()
			code, err := getUrlHeadResult(url)
			if err != nil {
				_, isUnknownAuthorityError := err.(x509.UnknownAuthorityError)
				gotCertError = gotCertError || isUnknownAuthorityError
				fmt.Printf("\033[0;91mHTTP HEAD request to URL '%s' failed: %v. (Check reason: %v)\033[0m\n", url, err, pChecks[0])
				atomic.AddInt32(&errorCount, 1)
			} else if code != http.StatusOK {
				fmt.Printf("\033[0;91mHTTP HEAD request to URL '%s' yielded bad response code %d. (Check reason: %v).\033[0m\n", url, code, pChecks[0])
				atomic.AddInt32(&errorCount, 1)
			} else {
				fmt.Printf("OK: Resource %s is available. (Reason for check: %v)\n", url, pChecks[0])
			}
		}(url, pChecks)
	}
	waitgroup.Wait()
	if errorCount > 0 {
		if gotCertError {
			fmt.Printf("\033[0;91mThere was at least one certificate-related error. The system's certificate pool may be out of date.\033[0m\n")
		}
		return fmt.Errorf("%d out of %d tested URLs from \"%s\" do not point to valid resources", errorCount, len(urlMap), filePath)
	}
	return nil
}

func collectURLs(data []byte, skipJarCheck bool) (urlMap map[string]platformChecks, err error) {
	urlMap = make(map[string]platformChecks)
	failCount := 0
	for _, os := range []string{"windows", "darwin", "linux"} {
		for _, arch := range []string{"386", "amd64"} {
			deploymentConfig := config.ParseDeploymentConfig(strings.NewReader(string(data)), os, arch)
			for _, update := range deploymentConfig.LauncherUpdate {
				addUrlWithDetails(urlMap, update.BundleInfoURL, platformCheck{reasonUpdate, config.Platform{OS: os, Arch: arch}})
			}
			for _, update := range deploymentConfig.Bundles {
				addUrlWithDetails(urlMap, update.BundleInfoURL, platformCheck{reasonBundle, config.Platform{OS: os, Arch: arch}})
			}
			for _, command := range deploymentConfig.Execution.Commands {
				if !collectCommandURLs(urlMap, deploymentConfig, os, arch, command, skipJarCheck) {
					failCount++
				}
			}
		}
	}
	if failCount > 0 {
		err = fmt.Errorf("%d commands were malformed", failCount)
	}
	return urlMap, err
}

func addUrlWithDetails(urlMap map[string]platformChecks, url string, pCheck platformCheck) {
	presentChecks, ok := urlMap[url]
	if ok {
		urlMap[url] = append(presentChecks, pCheck)
	} else {
		urlMap[url] = platformChecks{pCheck}
	}
}

func isAbsForOS(filePath string, os string) bool {
	if os == system.OsWindows {
		return filePath != "" && ((filePath[0] >= 'A' && filePath[0] <= 'Z') || (filePath[0] >= 'a' && filePath[0] <= 'z')) &&
			(len(filePath) == 1 || (filePath[1] == ':' && (len(filePath) == 2 || filePath[2] == '\\' || filePath[2] == '/')))
	}
	return path.IsAbs(filePath)
}

func collectCommandURLs(urlMap map[string]platformChecks, deploymentConfig *config.DeploymentConfig, os, arch string, command config.Command, skipJarCheck bool) (success bool) {
	commandNameUnix := strings.ReplaceAll(command.Name, `\`, "/")
	bundleName := misc.FirstElementOfPath(commandNameUnix)
	if isAbsForOS(command.Name, os) {
		fmt.Printf("\033[0;96mNote: cannot validate absolute command path for bundle \"%s\". (Command \"%s\" for platform %s-%s is absolute path).\033[0m\n", bundleName, command.Name, os, arch)
		return true
	} else if !strings.Contains(commandNameUnix, "/") {
		fmt.Printf("\033[0;91mCould not get bundle URL for bundle \"%s\": relative command \"%s\" for platform %s-%s does not descend into a bundle directory.\033[0m\n", bundleName, command.Name, os, arch)
		return false
	}
	bundleURL := getBundleURL(bundleName, deploymentConfig)
	if bundleURL == "" {
		fmt.Printf("\033[0;91mNo BaseURL configured or inferable for bundle \"%s\": not set. (Required for command \"%s\" on platform %s-%s).\033[0m\n", bundleName, command.Name, os, arch)
		return false
	}
	binaryURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(commandNameUnix))
	if os == system.OsWindows && !strings.HasSuffix(binaryURL, ".exe") {
		binaryURL += ".exe"
	}
	addUrlWithDetails(urlMap, binaryURL, platformCheck{reasonCommand, config.Platform{OS: os, Arch: arch}})
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

func collectJarURL(urlMap map[string]platformChecks, deploymentConfig *config.DeploymentConfig, command config.Command, os, arch string) {
	check := false
	for _, arg := range command.Arguments {
		if check {
			jarPath := strings.ReplaceAll(arg, `\`, "/")
			bundleName := misc.FirstElementOfPath(jarPath)
			bundleURL := getBundleURL(bundleName, deploymentConfig)
			jarURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(jarPath))
			addUrlWithDetails(urlMap, jarURL, platformCheck{reasonJar, config.Platform{OS: os, Arch: arch}})
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
