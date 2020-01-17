package config_test

import (
	"os"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/system"
)

func TestDeploymentConfigWithTemplates(t *testing.T) {
	f, err := os.Open("test/dpc_template.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	configString, err := config.ReadDeploymentConfig(f, "AmigaOS", "68000x16")
	if err != nil {
		t.Fatal(err)
	}
	expectedConfigString := string(system.MustReadFile("test/dpc_template.expected.json"))
	if configString != expectedConfigString {
		t.Fatalf("Files do not match. Expected: \n%s\nActual:\n%s\n", expectedConfigString, configString)
	}
}
