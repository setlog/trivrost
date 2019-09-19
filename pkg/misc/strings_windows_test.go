package misc_test

import (
	"github.com/setlog/trivrost/pkg/misc"
	"testing"
)

func TestExtensionlessFileName(t *testing.T) {
	tests := []struct {
		filePath, expected string
	}{
		{`C:\Folder\SomeFile.txt.backup`, "SomeFile.txt"},
		{`Folder\SomeFile.txt.backup`, "SomeFile.txt"},
		{`SomeFile.txt.backup`, "SomeFile.txt"},
		{`C:\Folder\SomeFile.txt.`, "SomeFile.txt"},
		{`Folder\SomeFile.txt.`, "SomeFile.txt"},
		{`SomeFile.txt.`, "SomeFile.txt"},
		{`C:\Folder\SomeFile`, "SomeFile"},
		{`Folder\SomeFile`, "SomeFile"},
		{`SomeFile`, "SomeFile"},
	}
	for i, test := range tests {
		result := misc.ExtensionlessFileName(test.filePath)
		if result != test.expected {
			t.Errorf("Test #%d failed: ExtensionlessFileName(\"%s\") yielded %s. Expected %s.", i+1, test.filePath, result, test.expected)
		}
	}
}
