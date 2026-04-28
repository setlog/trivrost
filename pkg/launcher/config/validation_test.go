package config_test

import (
	"strings"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func TestAcceptCorrectDeploymentConfig(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_complex.json")))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestAllowUnknownFields(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_pseudo_future.json")))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDetectBadPlatform(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_bad_platform.json")))
	if err == nil || !strings.Contains(err.Error(), "LauncherUpdate.0.TargetPlatforms.0: Does not match pattern") ||
		!strings.Contains(err.Error(), "Bundles.0.TargetPlatforms.0: Does not match pattern") ||
		!strings.Contains(err.Error(), "Bundles.1.TargetPlatforms.0: Does not match pattern") ||
		!strings.Contains(err.Error(), "Execution.Commands.1.TargetPlatforms.1: Does not match pattern") {
		t.Fatalf("%v", err)
	}
}

func TestDetectBadEnvVarValue(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_bad_env_var.json")))
	if err == nil || !strings.Contains(err.Error(), "Invalid type. Expected: string, given: integer") {
		t.Fatalf("%v", err)
	}
}

func TestNoSelfUpdateAllowed(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_no_self_update.json")))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestEmptySelfUpdateAllowed(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_empty_self_update.json")))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestDetectMissingExecution(t *testing.T) {
	err := config.ValidateDeploymentConfig(string(system.MustReadFile("test/dpc_exec_missing.json")))
	if err == nil || !strings.Contains(err.Error(), "Execution is required") {
		t.Fatalf("%v", err)
	}
}

func TestDetectInvalidBaseURLAndIsUpdateMandatoryTypes(t *testing.T) {
	err := config.ValidateDeploymentConfig(`{
		"Timestamp": "2019-02-07 14:53:17",
		"Bundles": [
			{
				"BundleInfoURL": "https://example.com/testapp/bundleinfo.json",
				"BaseURL": 42,
				"LocalDirectory": "app",
				"IsUpdateMandatory": "yes"
			}
		],
		"Execution": {
			"Commands": [
				{
					"Name": "java"
				}
			]
		}
	}`)
	if err == nil ||
		!strings.Contains(err.Error(), "Bundles.0.BaseURL: Invalid type. Expected: string, given: integer") ||
		!strings.Contains(err.Error(), "Bundles.0.IsUpdateMandatory: Invalid type. Expected: boolean, given: string") {
		t.Fatalf("%v", err)
	}
}
