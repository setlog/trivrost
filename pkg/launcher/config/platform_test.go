package config_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func TestFilterApplicableBundles(t *testing.T) {
	testos := runtime.GOOS
	testarch := system.GetOSArch()
	tests := []struct {
		bundles []config.BundleConfig
		want    []config.BundleConfig
	}{
		{bundles(bundle(testos)), bundles(bundle(testos))},
		{bundles(bundle(testarch)), bundles(bundle(testarch))},
		{bundles(bundle(testos + "-" + testarch)), bundles(bundle(testos + "-" + testarch))},
		{bundles(bundle("foo")), bundles()},
		{bundles(bundle("foo-" + testarch)), bundles()},
		{bundles(bundle(testos + "-foo")), bundles()},
		{bundles(bundle(testos, "foo")), bundles(bundle(testos, "foo"))},
		{bundles(bundle()), bundles(bundle())},
	}
	for i, test := range tests {
		result := config.FilterBundlesByPlatform(test.bundles, testos, testarch)
		// assert both are empty or equal
		if !((len(result) == 0 && len(test.want) == 0) || reflect.DeepEqual(result, test.want)) {
			t.Errorf("Test #%d failed: Got %v. Expected %v.", i+1, result, test.want)
		}
	}
}

func bundles(bundles ...config.BundleConfig) []config.BundleConfig {
	return append([]config.BundleConfig{}, bundles...)
}

func bundle(platforms ...string) config.BundleConfig {
	return config.BundleConfig{
		HashDataConfig: config.HashDataConfig{
			BundleInfoURL: "https://example.com/windows/testapp/bundleinfo.json",
			BaseURL:       "https://example.com/windows/testapp"},
		LocalDirectory:  "app",
		TargetPlatforms: append([]string{}, platforms...)}
}

func TestFilterApplicableCommands(t *testing.T) {
	testos := runtime.GOOS
	testarch := system.GetOSArch()
	tests := []struct {
		commands []config.Command
		want     []config.Command
	}{
		{commands(command(testos)), commands(command(testos))},
		{commands(command(testarch)), commands(command(testarch))},
		{commands(command(testos + "-" + testarch)), commands(command(testos + "-" + testarch))},
		{commands(command("foo")), commands()},
		{commands(command("foo-" + testarch)), commands()},
		{commands(command(testos + "-foo")), commands()},
		{commands(command(testos, "foo")), commands(command(testos, "foo"))},
		{commands(command()), commands(command())},
	}
	for i, test := range tests {
		result := config.FilterCommandsByPlatform(test.commands, testos, testarch)
		if !((len(result) == 0 && len(test.want) == 0) || reflect.DeepEqual(result, test.want)) {
			t.Errorf("Test #%d failed: Got %v. Expected %v.", i+1, result, test.want)
		}
	}
}

func commands(commands ...config.Command) []config.Command {
	return append([]config.Command{}, commands...)
}

func command(platforms ...string) config.Command {
	return config.Command{
		Name:            "java/bin/java",
		Arguments:       []string{"java", "-jar"},
		TargetPlatforms: append([]string{}, platforms...)}
}
