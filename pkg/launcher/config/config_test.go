package config_test

import (
	"strings"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

const launcherConfigContent string = `{
    "DeploymentConfigURL": "https://example.com/deployment-config.json",
    "VendorName": "example_company",
    "ProductName": "example_product"
}
`

const deploymentConfigContent string = `{
    "Timestamp": "2019-02-07 14:53:17",
    "LauncherUpdate": [
        {
            "BundleInfoURL": "https://example.com/windows/launcher/bundleinfo.json",
            "BaseURL": "https://example.com/windows/launcher",
            "TargetPlatforms": [ "windows" ]
        }
    ],
    "Bundles": [
        {
            "BundleInfoURL": "https://example.com/windows/testapp/bundleinfo.json",
            "BaseURL": "https://example.com/windows/testapp",
            "LocalDirectory": "app"
        },
        {
            "BundleInfoURL": "https://example.com/windows/java/bundleinfo.json",
            "BaseURL": "https://example.com/windows/java",
            "LocalDirectory": "java",
            "TargetPlatforms": [ "windows" ]
        }
    ],
    "Execution": {
        "Commands": [
            {
                "Name": "java/bin/java",
                "Arguments": [ "-jar", "foo", "-Xm1024M" ],
                "Env": {
                    "NEW_ENV": "New env variable.",
                    "OMIT_ENV": null
                },
                "TargetPlatforms": [ "windows", "linux" ]
            }
        ]
    }
}`

const deploymentConfigWithoutTargetPlatformsContent string = `{
    "Timestamp": "2019-02-07 14:53:17",
    "Bundles": [
        {
            "BundleInfoURL": "https://example.com/windows/testapp/bundleinfo.json",
            "BaseURL": "https://example.com/windows/testapp",
            "LocalDirectory": "app"
        },
        {
            "BundleInfoURL": "https://example.com/windows/java/bundleinfo.json",
            "BaseURL": "https://example.com/windows/java",
            "LocalDirectory": "java"
        }
    ],
    "Execution": {
        "Commands": [
            {
                "Name": "java/bin/java",
                "Arguments": [ "-jar", "foo", "-Xm1024M" ],
                "Env": {
                    "NEW_ENV": "New env variable.",
                    "OMIT_ENV": null
                }
            }
        ]
    }
}`

const deploymentConfigWithoutEnvSectionContent string = `{
    "Timestamp": "2019-02-07 14:53:17",
    "Bundles": [
        {
            "BundleInfoURL": "https://example.com/windows/testapp/bundleinfo.json",
            "BaseURL": "https://example.com/windows/testapp",
            "LocalDirectory": "app"
        },
        {
            "BundleInfoURL": "https://example.com/windows/java/bundleinfo.json",
            "BaseURL": "https://example.com/windows/java",
            "LocalDirectory": "java",
            "TargetPlatforms": [ "windows" ]
        },
        {
            "BundleInfoURL": "https://example.com/linux/java/bundleinfo.json",
            "BaseURL": "https://example.com/linux/java",
            "LocalDirectory": "java",
            "TargetPlatforms": [ "linux" ]
        }
    ],
    "Execution": {
        "Commands": [
            {
                "Name": "java/bin/java",
                "Arguments": [ "-jar", "foo", "-Xm1024M" ],
                "TargetPlatforms": [ "windows", "linux" ]
            }
        ]
    }
}`

func TestReadLauncherConfig(t *testing.T) {
	reader := strings.NewReader(launcherConfigContent)
	cfg := config.ReadLauncherConfigFromReader(reader)
	tests := []struct {
		value, valueName, expected string
	}{
		{cfg.DeploymentConfigURL, "Launcher.DeploymentConfigURL", "https://example.com/deployment-config.json"},
		{cfg.VendorName, "Launcher.VendorName", "example_company"},
		{cfg.ProductName, "Launcher.ProductName", "example_product"},
	}
	for _, test := range tests {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Got: %s; Expected: %s.", test.valueName, test.value, test.expected)
		}
	}
}

func TestReadDeploymentConfig(t *testing.T) {
	reader := strings.NewReader(deploymentConfigContent)
	cfg := config.ParseDeploymentConfig(reader, "windows", "amd64")

	testsValues := []struct {
		value, valueName, expected string
	}{
		{cfg.Timestamp, "Timestamp", "2019-02-07 14:53:17"},

		{cfg.LauncherUpdate[0].BundleInfoURL, "LauncherUpdate[0].BundleInfoURL", "https://example.com/windows/launcher/bundleinfo.json"},
		{cfg.LauncherUpdate[0].BaseURL, "LauncherUpdate[0].BaseURL", "https://example.com/windows/launcher"},
		{cfg.LauncherUpdate[0].TargetPlatforms[0], "LauncherUpdate[0].TargetPlatforms[0]", "windows"},

		{cfg.Bundles[0].BundleInfoURL, "Bundles[0].BundleInfoURL", "https://example.com/windows/testapp/bundleinfo.json"},
		{cfg.Bundles[0].BaseURL, "Bundles[0].BaseURL", "https://example.com/windows/testapp"},
		{cfg.Bundles[0].LocalDirectory, "Bundles[0].LocalDirectory", "app"},

		{cfg.Bundles[1].BundleInfoURL, "Bundles[1].BundleInfoURL", "https://example.com/windows/java/bundleinfo.json"},
		{cfg.Bundles[1].BaseURL, "Bundles[1].BaseURL", "https://example.com/windows/java"},
		{cfg.Bundles[1].LocalDirectory, "Bundles[1].LocalDirectory", "java"},
		{cfg.Bundles[1].TargetPlatforms[0], "Bundles[1].TargetPlatforms[0]", "windows"},

		{cfg.Execution.Commands[0].Name, "Execution.Commands[0].Name", "java/bin/java"},
		{cfg.Execution.Commands[0].Arguments[0], "Execution.Commands[0].Arguments[0]", "-jar"},
		{cfg.Execution.Commands[0].Arguments[1], "Execution.Commands[0].Arguments[1]", "foo"},
		{cfg.Execution.Commands[0].Arguments[2], "Execution.Commands[0].Arguments[2]", "-Xm1024M"},
		{cfg.Execution.Commands[0].TargetPlatforms[0], "Execution.Commands[0].TargetPlatforms[0]", "windows"},
		{cfg.Execution.Commands[0].TargetPlatforms[1], "Execution.Commands[0].TargetPlatforms[1]", "linux"},
	}
	for _, test := range testsValues {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Got: %s; Expected: %s.", test.valueName, test.value, test.expected)
		}
	}

	if *cfg.Execution.Commands[0].Env["NEW_ENV"] != "New env variable." {
		t.Errorf("config.ReadConfig() test failed for Execution.Commands[0].Env[\"NEW_ENV\"]. Got: %v; Expected: %s.",
			cfg.Execution.Commands[0].Env["NEW_ENV"], "New env variable.")
	}
	if cfg.Execution.Commands[0].Env["OMIT_ENV"] != nil {
		t.Errorf("config.ReadConfig() test failed for Execution.Commands[0].Env[\"OMIT_ENV\"]. Got: %v; Expected null",
			cfg.Execution.Commands[0].Env["OMIT_ENV"])
	}

	testsLengths := []struct {
		value     int
		valueName string
		expected  int
	}{
		{len(cfg.Bundles[0].TargetPlatforms), "Bundles[0].TargetPlatforms", 0},
		{len(cfg.Bundles[1].TargetPlatforms), "Bundles[1].TargetPlatforms", 1},
		{len(cfg.Execution.Commands[0].TargetPlatforms), "Execution.Commands[0].TargetPlatforms", 2},
	}
	for _, test := range testsLengths {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Actual length: %d; Expected length: %d.", test.valueName, test.value, test.expected)
		}
	}
}

func TestReadDeploymentWithoutTargetPlatformsConfig(t *testing.T) {
	reader := strings.NewReader(deploymentConfigWithoutTargetPlatformsContent)
	cfg := config.ParseDeploymentConfig(reader, "windows", "amd64")
	tests := []struct {
		value, valueName, expected string
	}{
		{cfg.Timestamp, "Timestamp", "2019-02-07 14:53:17"},

		{cfg.Bundles[0].BundleInfoURL, "Bundles[0].BundleInfoURL", "https://example.com/windows/testapp/bundleinfo.json"},
		{cfg.Bundles[0].BaseURL, "Bundles[0].BaseURL", "https://example.com/windows/testapp"},
		{cfg.Bundles[0].LocalDirectory, "Bundles[0].LocalDirectory", "app"},

		{cfg.Bundles[1].BundleInfoURL, "Bundles[1].BundleInfoURL", "https://example.com/windows/java/bundleinfo.json"},
		{cfg.Bundles[1].BaseURL, "Bundles[1].BaseURL", "https://example.com/windows/java"},
		{cfg.Bundles[1].LocalDirectory, "Bundles[1].LocalDirectory", "java"},

		{cfg.Execution.Commands[0].Name, "Execution.Commands[0].Name", "java/bin/java"},
		{cfg.Execution.Commands[0].Arguments[0], "Execution.Commands[0].Arguments[0]", "-jar"},
		{cfg.Execution.Commands[0].Arguments[1], "Execution.Commands[0].Arguments[1]", "foo"},
		{cfg.Execution.Commands[0].Arguments[2], "Execution.Commands[0].Arguments[2]", "-Xm1024M"},
	}
	for _, test := range tests {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Got: %s; Expected: %s.", test.valueName, test.value, test.expected)
		}
	}

	if *cfg.Execution.Commands[0].Env["NEW_ENV"] != "New env variable." {
		t.Errorf("config.ReadConfig() test failed for Execution.Commands[0].Env[\"NEW_ENV\"]. Got: %v; Expected: %s.",
			cfg.Execution.Commands[0].Env["NEW_ENV"], "New env variable.")
	}
	if cfg.Execution.Commands[0].Env["OMIT_ENV"] != nil {
		t.Errorf("config.ReadConfig() test failed for Execution.Commands[0].Env[\"OMIT_ENV\"]. Got: %v; Expected null",
			cfg.Execution.Commands[0].Env["OMIT_ENV"])
	}

	testsLengths := []struct {
		value     int
		valueName string
		expected  int
	}{
		{len(cfg.Bundles[0].TargetPlatforms), "Bundles[0].TargetPlatforms", 0},
		{len(cfg.Bundles[1].TargetPlatforms), "Bundles[1].TargetPlatforms", 0},
		{len(cfg.Execution.Commands[0].TargetPlatforms), "Execution.Commands[0].TargetPlatforms", 0},
	}
	for _, test := range testsLengths {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Actual length: %d; Expected length: %d.", test.valueName, test.value, test.expected)
		}
	}
}

func TestReadDeploymentConfigWithoutEnvSection(t *testing.T) {
	reader := strings.NewReader(deploymentConfigWithoutEnvSectionContent)
	cfg := config.ParseDeploymentConfig(reader, "windows", "amd64")
	tests := []struct {
		value, valueName, expected string
	}{
		{cfg.Timestamp, "Timestamp", "2019-02-07 14:53:17"},

		{cfg.Bundles[0].BundleInfoURL, "Bundles[0].BundleInfoURL", "https://example.com/windows/testapp/bundleinfo.json"},
		{cfg.Bundles[0].BaseURL, "Bundles[0].BaseURL", "https://example.com/windows/testapp"},
		{cfg.Bundles[0].LocalDirectory, "Bundles[0].LocalDirectory", "app"},

		{cfg.Bundles[1].BundleInfoURL, "Bundles[1].BundleInfoURL", "https://example.com/windows/java/bundleinfo.json"},
		{cfg.Bundles[1].BaseURL, "Bundles[1].BaseURL", "https://example.com/windows/java"},
		{cfg.Bundles[1].LocalDirectory, "Bundles[1].LocalDirectory", "java"},
		{cfg.Bundles[1].TargetPlatforms[0], "Bundles[1].TargetPlatforms[0]", "windows"},

		{cfg.Execution.Commands[0].Name, "Execution.Commands[0].Name", "java/bin/java"},
		{cfg.Execution.Commands[0].Arguments[0], "Execution.Commands[0].Arguments[0]", "-jar"},
		{cfg.Execution.Commands[0].Arguments[1], "Execution.Commands[0].Arguments[1]", "foo"},
		{cfg.Execution.Commands[0].Arguments[2], "Execution.Commands[0].Arguments[2]", "-Xm1024M"},
		{cfg.Execution.Commands[0].TargetPlatforms[0], "Execution.Commands[0].TargetPlatforms[0]", "windows"},
		{cfg.Execution.Commands[0].TargetPlatforms[1], "Execution.Commands[0].TargetPlatforms[1]", "linux"},
	}
	for _, test := range tests {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Got: %s; Expected: %s.", test.valueName, test.value, test.expected)
		}
	}

	testsLengths := []struct {
		value     int
		valueName string
		expected  int
	}{
		{len(cfg.Bundles[0].TargetPlatforms), "Bundles[0].TargetPlatforms", 0},
		{len(cfg.Bundles[1].TargetPlatforms), "Bundles[1].TargetPlatforms", 1},
		{len(cfg.Execution.Commands[0].TargetPlatforms), "Execution.Commands[0].TargetPlatforms", 2},
		{len(cfg.Execution.Commands[0].Env), "Execution.Commands[0].Env", 0},
	}
	for _, test := range testsLengths {
		if test.value != test.expected {
			t.Errorf("config.ReadConfig() test failed for %s. Actual length: %d; Expected length: %d.", test.valueName, test.value, test.expected)
		}
	}
}
