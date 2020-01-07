package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func main() {
	flags := parseFlags()
	if flags.ActAsService {
		actAsService(flags)
	} else {
		log.SetFlags(0)
		if len(validateDeploymentConfig(flags.DeploymentConfigUrl, flags.SkipUrlCheck, flags.SkipJarChek)) > 0 {
			os.Exit(1)
		}
	}
}

func validateDeploymentConfig(url string, skipUrlCheck bool, skipJarCheck bool) []error {
	log.Printf("Validating deployment-config at %s...\n", url)
	expandedDeploymentConfig, err := getFile(url)
	if err != nil {
		log.Printf("\033[0;91mCould not validate deployment-config at URL %s: %v\033[0m\n", url, err)
		return []error{err}
	}

	err = config.ValidateDeploymentConfig(string(expandedDeploymentConfig))
	if err != nil {
		log.Printf("\033[0;91mCould not validate deployment-config at URL %s: %v\033[0m\n", url, err)
		return []error{err}
	}

	if !skipUrlCheck {
		return checkURLs(expandedDeploymentConfig, skipJarCheck)
	}
	return nil
}

func checkURLs(expandedDeploymentConfig []byte, skipJarCheck bool) []error {
	urlMap, errs := collectURLs(expandedDeploymentConfig, skipJarCheck)

	waitgroup := sync.WaitGroup{}
	waitgroup.Add(len(urlMap))
	var gotCertError bool

	// check all url in parallel
	var errorCount int32
	for url, details := range urlMap {
		go func(url string, details checkDetails) {
			defer waitgroup.Done()
			code, err := getHttpHeadResult(url)
			if err != nil {
				_, isUnknownAuthorityError := err.(x509.UnknownAuthorityError)
				gotCertError = gotCertError || isUnknownAuthorityError
				log.Printf("\033[0;91mHTTP HEAD request to URL '%s' failed: %v. (Check reason: %v)\033[0m\n", url, err, details)
				atomic.AddInt32(&errorCount, 1)
			} else if code != http.StatusOK {
				log.Printf("\033[0;91mHTTP HEAD request to URL '%s' yielded bad response code %d. (Check reason: %v)\033[0m\n", url, code, details)
				atomic.AddInt32(&errorCount, 1)
			} else {
				log.Printf("OK: Resource %s is available. (Reason for check: %v)\n", url, details)
			}
		}(url, details)
	}
	waitgroup.Wait()
	if errorCount > 0 {
		if gotCertError {
			log.Printf("\033[0;91mThere was at least one certificate-related error. The system's certificate pool may be out of date.\033[0m\n")
			errs = append(errs, fmt.Errorf("there was at least one certificate-related error. The system's certificate pool may be out of date"))
		}
		errs = append(errs, fmt.Errorf("%d out of %d tested URLs do not point to valid resources", errorCount, len(urlMap)))
	}
	return errs
}

func collectURLs(data []byte, skipJarCheck bool) (urlMap map[string]checkDetails, errs []error) {
	urlMap = make(map[string]checkDetails)
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
				err := collectCommandURLs(urlMap, deploymentConfig, operatingsystem, arch, command, skipJarCheck)
				if err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	return urlMap, errs
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

func collectCommandURLs(urlMap map[string]checkDetails, deploymentConfig *config.DeploymentConfig, os, arch string, command config.Command, skipJarCheck bool) error {
	commandNameUnix := strings.ReplaceAll(command.Name, `\`, "/")
	if path.IsAbs(commandNameUnix) || !strings.Contains(commandNameUnix, "/") {
		return fmt.Errorf("%s is not a relative path which descends into at least one folder", commandNameUnix)
	}
	bundleName := misc.FirstElementOfPath(commandNameUnix)
	bundleURL := getBundleURL(bundleName, deploymentConfig)
	if bundleURL == "" {
		log.Printf("\033[0;91mCould not get bundle URL for bundle \"%s\" for platform %s-%s. (Required for command \"%s\")\033[0m\n", bundleName, os, arch, command.Name)
		return fmt.Errorf("could not get bundle URL for bundle \"%s\" for platform %s-%s (Required for command \"%s\")", bundleName, os, arch, command.Name)
	}
	binaryURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(commandNameUnix))
	if os == system.OsWindows && !strings.HasSuffix(binaryURL, ".exe") {
		binaryURL += ".exe"
	}
	addUrlWithDetails(urlMap, binaryURL, checkDetails{reasonCommand, os, arch, 0})
	if !skipJarCheck {
		if strings.HasSuffix(binaryURL, "/java.exe") || strings.HasSuffix(binaryURL, "/javaw.exe") ||
			strings.HasSuffix(binaryURL, "/java") {
			err := collectJarURL(urlMap, deploymentConfig, command, os, arch)
			if err != nil {
				log.Printf("\033[0;91mCould not get JAR URL for bundle \"%s\" for platform %s-%s (Required for command \"%s\"): %v\033[0m\n", bundleName, os, arch, command.Name, err)
				return fmt.Errorf("could not get JAR URL for bundle \"%s\" for platform %s-%s (Required for command \"%s\"): %w", bundleName, os, arch, command.Name, err)
			}
		}
	}
	return nil
}

func stripFirstPathElement(s string) string {
	s = path.Clean(s)
	parts := strings.Split(s, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts[1:], "/")
}

func collectJarURL(urlMap map[string]checkDetails, deploymentConfig *config.DeploymentConfig, command config.Command, os, arch string) error {
	check := false
	for _, arg := range command.Arguments {
		if check {
			jarPath := strings.ReplaceAll(arg, `\`, "/")
			bundleName := misc.FirstElementOfPath(jarPath)
			bundleURL := getBundleURL(bundleName, deploymentConfig)
			if bundleURL == "" {
				return fmt.Errorf("jar path '%s' does not descend into a bundle directory", arg)
			}
			jarURL := misc.MustJoinURL(bundleURL, stripFirstPathElement(jarPath))
			addUrlWithDetails(urlMap, jarURL, checkDetails{reasonJar, os, arch, 0})
			break
		}
		if arg == "-jar" {
			check = true
		}
	}
	return nil
}

func getBundleURL(bundleName string, deploymentConfig *config.DeploymentConfig) string {
	for _, bundle := range deploymentConfig.Bundles {
		if bundle.LocalDirectory == bundleName {
			return bundle.BaseURL
		}
	}
	return ""
}

func fatalf(formatMessage string, args ...interface{}) {
	fmt.Printf("\033[0;91mFatal: "+formatMessage+"\033[0m\n", args...)
	os.Exit(1)
}
