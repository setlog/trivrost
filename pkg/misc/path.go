package misc

import (
	"path/filepath"
	"strings"
)

func MakePathAbsolute(path string, absoluteReferencePath string) string {
	path = filepath.FromSlash(path)
	if !filepath.IsAbs(path) {
		return filepath.Join(filepath.FromSlash(absoluteReferencePath), path)
	}
	return filepath.Clean(path)
}

func FirstElementOfPath(path string) string {
	elements := strings.Split(filepath.ToSlash(path), "/")
	if len(elements) == 0 {
		return ""
	}
	if len(elements[0]) > 0 {
		return elements[0]
	}
	if len(elements) > 1 {
		return elements[1]
	}
	return ""
}

func StripFirstPathElement(filePath string, sep rune) (newPath, strippedElement string) {
	elements := strings.Split(filePath, string(sep))
	return strings.Join(elements[1:], string(sep)), elements[0]
}
