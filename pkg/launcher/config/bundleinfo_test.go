package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

func TestReadBundleInfoPreservesForwardSlashPaths(t *testing.T) {
	reader := strings.NewReader(`{
		"Timestamp": "2019-02-07 14:53:17",
		"UniqueBundleName": "bundle",
		"BundleFiles": {
			"app/bin/run": { "SHA256": "abc", "Size": 1 },
			"top.txt": { "SHA256": "def", "Size": 2 }
		}
	}`)

	info := config.ReadInfoFromReader(reader)

	if _, ok := info.BundleFiles["app/bin/run"]; !ok {
		t.Fatalf("expected nested clean path to be kept")
	}
	if _, ok := info.BundleFiles["top.txt"]; !ok {
		t.Fatalf("expected top-level clean path to be kept")
	}
}

func TestGetFileHashesReturnsOSPaths(t *testing.T) {
	info := &config.BundleInfo{
		BundleFiles: config.FileInfoMap{
			"app/bin/run": {SHA256: "abc", Size: 1},
		},
	}

	if _, ok := info.GetFileHashes()[filepath.FromSlash("app/bin/run")]; !ok {
		t.Fatalf("expected file hashes to use OS-specific separators")
	}
}

func TestReadBundleInfoRejectsUnsafePaths(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
	}{
		{name: "empty", filePath: ""},
		{name: "absolute", filePath: "/etc/passwd"},
		{name: "parentTraversal", filePath: "../outside"},
		{name: "embeddedTraversal", filePath: "app/../outside"},
		{name: "dotPrefix", filePath: "./app"},
		{name: "doubleSlash", filePath: "app//file"},
		{name: "backslash", filePath: `app\\file`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := strings.NewReader(`{
				"Timestamp": "2019-02-07 14:53:17",
				"UniqueBundleName": "bundle",
				"BundleFiles": {
					"` + test.filePath + `": { "SHA256": "abc", "Size": 1 }
				}
			}`)

			defer func() {
				if recover() == nil {
					t.Fatalf("expected panic for unsafe bundle file path %q", test.filePath)
				}
			}()

			config.ReadInfoFromReader(reader)
		})
	}
}
