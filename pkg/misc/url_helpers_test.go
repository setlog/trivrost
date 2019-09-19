package misc_test

import (
	"github.com/setlog/trivrost/pkg/misc"
	"testing"
)

func TestMustJoinURL(t *testing.T) {
	tests := []struct {
		base, reference, expected string
	}{
		{"http://example.com", "exemplary", "http://example.com/exemplary"},
		{"http://example", "exemplary", "http://example/exemplary"},
		{"https://example.com", "exemplary", "https://example.com/exemplary"},
		{"https://example", "exemplary", "https://example/exemplary"},
		{"http://example.com/", "/exemplary", "http://example.com/exemplary"},
		{"http://example/", "/exemplary", "http://example/exemplary"},
		{"http://example.com/", "/exemplary/thing", "http://example.com/exemplary/thing"},
		{"http://example/", "/exemplary/thing", "http://example/exemplary/thing"},
		{"http://example.com/", "/exemplary/thing/", "http://example.com/exemplary/thing"},
		{"http://example/", "/exemplary/thing/", "http://example/exemplary/thing"},
	}
	for i, test := range tests {
		result := misc.MustJoinURL(test.base, test.reference)
		if result != test.expected {
			t.Errorf("Test #%d failed: MustJoinURL(\"%s\", \"%s\") yielded %s. Expected %s.", i+1, test.base, test.reference, result, test.expected)
		}
	}
}

func TestMustStripLastURLPathElement(t *testing.T) {
	tests := []struct {
		original, expected string
	}{
		{"http://example.com/exemplary", "http://example.com/"},
		{"http://example/exemplary", "http://example/"},
		{"https://example.com/exemplary", "https://example.com/"},
		{"https://example/exemplary", "https://example/"},
		{"http://example.com//exemplary", "http://example.com/"},
		{"http://example//exemplary", "http://example/"},
		{"http://example.com//exemplary/thing", "http://example.com/exemplary"},
		{"http://example/exemplary/thing", "http://example/exemplary"},
		{"http://example.com/exemplary/thing/", "http://example.com/exemplary/thing"},
		{"http://example/exemplary/thing/", "http://example/exemplary/thing"},
		{"http://example/exemplary/bundleinfo.json", "http://example/exemplary"},
	}
	for i, test := range tests {
		result := misc.MustStripLastURLPathElement(test.original)
		if result != test.expected {
			t.Errorf("Test #%d failed: MustStripLastURLPathElement(\"%s\") yielded %s. Expected %s.", i+1, test.original, result, test.expected)
		}
	}
}

func TestEllipsisURL(t *testing.T) {
	tests := []struct {
		original, expected string
	}{
		{"http://example.com/exemplary", "http://example.com/exemplary"},
		{"http://example/exemplary", "http://example/exemplary"},
		{"https://example.com/exemplary", "https://example.com/exemplary"},
		{"https://example/exemplary", "https://example/exemplary"},
		{"http://example.com//exemplary", "http://example.com//exemplary"},
		{"http://example//exemplary", "http://example//exemplary"},
		{"http://example.com//exemplary/thing", "http://example.com/.../thing"},
		{"http://example/exemplary/thing", "http://example/.../thing"},
		{"http://example/exa/thing", "http://example/exa/thing"},
		{"http://example/exam/thing", "http://example/.../thing"},
		{"http://example.com/exemplary/thing/", "http://example.com/.../thing/"},
		{"http://example.com/exe/thing/", "http://example.com/exe/thing/"},
		{"http://example.com/exem/thing/", "http://example.com/.../thing/"},
		{"http://example/exemplary/thing/", "http://example/.../thing/"},
		{"http://example/exemplary/bundleinfo.json", "http://example/.../bundleinfo.json"},
		{"http://example", "http://example"},
		{"http://example/", "http://example/"},
		{"http://example//", "http://example//"},
		{"http://example///", "http://example///"},
	}
	for i, test := range tests {
		result := misc.EllipsisURL(test.original, 0)
		if result != test.expected {
			t.Errorf("Test #%d failed: EllipsisURL(\"%s\") yielded %s. Expected %s.", i+1, test.original, result, test.expected)
		}
	}
}
