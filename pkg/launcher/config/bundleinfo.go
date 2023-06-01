package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/setlog/trivrost/pkg/misc"
	"github.com/setlog/trivrost/pkg/system"
)

type BundleInfo struct {
	Timestamp        string `json:"Timestamp"`
	UniqueBundleName string `json:"UniqueBundleName"`

	BundleFiles FileInfoMap `json:"BundleFiles"` // Within BundleInfo, keys are filepaths with forward slashes.
}

type FileInfoMap map[string]*FileInfo

type FileInfo struct {
	SHA256 string `json:"SHA256"`
	Size   int64  `json:"Size"`
}

// GetFileHashes returns a FileInfoMap for the info's BundleFiles using filepath.Separator in the place of forward slashes.
func (info *BundleInfo) GetFileHashes() FileInfoMap {
	fm := NewFileInfoMap()
	for filePath, fileInfo := range info.BundleFiles {
		genericFileInfo := *fileInfo
		fm[filepath.FromSlash(filePath)] = &genericFileInfo
	}
	return fm
}

// ForOS returns a copy of the FileInfoMap where all keys' forward slashes are replaced by the separator
// character of the underling operating system.
func (bundleFiles FileInfoMap) ForOS() FileInfoMap {
	osBundleFiles := make(FileInfoMap)
	for filePath, fileInfo := range bundleFiles {
		osFileInfo := *fileInfo
		osBundleFiles[filepath.FromSlash(filePath)] = &osFileInfo
	}
	return osBundleFiles
}

func (bundleFiles FileInfoMap) OmitEntriesWithMissingSha() FileInfoMap {
	newBundleFiles := make(FileInfoMap)
	for filePath, fileInfo := range bundleFiles {
		if fileInfo.SHA256 != "" {
			osFileInfo := *fileInfo
			newBundleFiles[filepath.ToSlash(filePath)] = &osFileInfo
		}
	}
	return newBundleFiles
}

// WithForwardSlashes returns a copy of the FileInfoMap where all occurrences of os.PathSeparator
// in the keys are replaced with forward slashes if os.PathSeparator is not already the forward slash.
func (bundleFiles FileInfoMap) WithForwardSlashes() FileInfoMap {
	forwardBundleFiles := make(FileInfoMap)
	for filePath, fileInfo := range bundleFiles {
		osFileInfo := *fileInfo
		forwardBundleFiles[filepath.ToSlash(filePath)] = &osFileInfo
	}
	return forwardBundleFiles
}

func (infoMap FileInfoMap) StripFirstPathElement(sep rune) FileInfoMap {
	newInfoMap := make(FileInfoMap)
	var stripElement string
	for filePath, fileInfo := range infoMap {
		newFilePath, strippedElement := misc.StripFirstPathElement(filePath, sep)
		if stripElement != "" && stripElement != strippedElement {
			panic(fmt.Sprintf("Mixed first path elements: \"%s\" vs \"%s\"", stripElement, strippedElement))
		} else if strippedElement == "" {
			panic(fmt.Sprintf("Nothing to strip from \"%s\"", filePath))
		}
		stripElement = strippedElement
		fileInfoNew := *fileInfo
		newInfoMap[newFilePath] = &fileInfoNew
	}
	return newInfoMap
}

func (infoMap FileInfoMap) FirstPathElement(sep rune) (firstElement string) {
	for filePath := range infoMap {
		_, strippedElement := misc.StripFirstPathElement(filePath, sep)
		if firstElement != "" && firstElement != strippedElement {
			panic(fmt.Sprintf("Mixed first path elements: \"%s\" vs \"%s\"", firstElement, strippedElement))
		}
		firstElement = strippedElement
	}
	return firstElement
}

func (bundleFiles FileInfoMap) Prepend(pathElement string, sep rune) FileInfoMap {
	newBundleFiles := make(FileInfoMap)
	for filePath, fileInfo := range bundleFiles {
		newFileInfo := *fileInfo
		if filePath == "" {
			newBundleFiles[pathElement] = &newFileInfo
		} else if pathElement == "" {
			newBundleFiles[filePath] = &newFileInfo
		} else {
			newBundleFiles[strings.TrimRight(pathElement, string(sep))+string(sep)+strings.TrimLeft(filePath, string(sep))] = &newFileInfo
		}
	}
	return newBundleFiles
}

func (bundleFiles FileInfoMap) FilePaths() []string {
	paths := make([]string, 0, len(bundleFiles))
	for filePath := range bundleFiles {
		paths = append(paths, filePath)
	}
	return paths
}

func ReadInfo(filePath string) *BundleInfo {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return ReadInfoFromByteSlice(data)
}

func ReadInfoFromReader(reader *strings.Reader) *BundleInfo {
	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return ReadInfoFromByteSlice(data)
}

func ReadInfoFromByteSlice(data []byte) *BundleInfo {
	info := BundleInfo{}
	err := json.Unmarshal(data, &info)
	if err != nil {
		panic(err)
	}
	return &info
}

func WriteInfo(info *BundleInfo, filePath string) {
	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	system.MustPutFile(filePath, data)
}
