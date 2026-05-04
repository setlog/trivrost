package hashing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/setlog/trivrost/pkg/dummy"
	"github.com/setlog/trivrost/pkg/launcher/config"
)

var infoForContent = config.FileInfoMap{
	"abc": &config.FileInfo{SHA256: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", Size: 3},
	"def": &config.FileInfo{SHA256: "cb8379ac2098aa165029e3938a51da0bcecfc008fd6795f401178647f96c5b34", Size: 3},
	"ghi": &config.FileInfo{SHA256: "50ae61e841fac4e8f9e40baf2ad36ec868922ea48368c18f9535e47db56dd7fb", Size: 3},
	"jkl": &config.FileInfo{SHA256: "268f277c6d766d31334fda0f7a5533a185598d269e61c76a805870244828a5f1", Size: 3},
	"mno": &config.FileInfo{SHA256: "cf63b8eb216845d24edd4b249b146957b42199cd12759647df90cb57525b4e90", Size: 3},
	"pqr": &config.FileInfo{SHA256: "d24bd97b5fb24761112354dec329c70a5c6e2dedcc9a6df160eefd1d671efe56", Size: 3},
}

func TestMustHashRelatively(t *testing.T) {
	fileMap := mustHashRelatively(context.Background(), dummyListDirectory, dummyReadFile, dummyStatFile, "x")
	expected := config.FileInfoMap{"foo": infoForContent["abc"], filepath.FromSlash("foo/bar"): infoForContent["def"], filepath.FromSlash("fuu/baaar"): infoForContent["ghi"], filepath.FromSlash("fuu/moo/meow/bla"): infoForContent["jkl"]}
	if !reflect.DeepEqual(fileMap, expected) {
		t.Errorf("Mismatch!\nGot:\n%v\nExpected:\n%v\n", fileMap, expected)
	}
}

func dummyListDirectory(dirPath string) ([]fs.DirEntry, error) {
	switch dirPath {
	case "x":
		return []fs.DirEntry{newDirEntry("foo", true), newDirEntry("fuu", true), newDirEntry("foo", false)}, nil
	case filepath.FromSlash("x/foo"):
		return []fs.DirEntry{newDirEntry("bar", false)}, nil
	case filepath.FromSlash("x/fuu"):
		return []fs.DirEntry{newDirEntry("moo", true), newDirEntry("baaar", false)}, nil
	case filepath.FromSlash("x/fuu/moo"):
		return []fs.DirEntry{newDirEntry("meow", true)}, nil
	case filepath.FromSlash("x/fuu/moo/meow"):
		return []fs.DirEntry{newDirEntry("bla", false)}, nil
	}
	return []fs.DirEntry{}, fmt.Errorf("Could not find the specified directory \"%s\"", dirPath)
}

func newDirEntry(name string, isDir bool) fs.DirEntry {
	return dirEntry{name: name, isDir: isDir}
}

type dirEntry struct {
	name  string
	isDir bool
}

func (de dirEntry) Name() string {
	return de.name
}

func (de dirEntry) IsDir() bool {
	return de.isDir
}

func (de dirEntry) Type() fs.FileMode {
	if de.isDir {
		return fs.ModeDir
	}
	return 0
}

func (de dirEntry) Info() (fs.FileInfo, error) {
	return dummy.NewFileInfo(de.name, de.isDir), nil
}

func dummyReadFile(filePath string) (io.ReadCloser, error) {
	switch filePath {
	case filepath.FromSlash("x/foo"):
		return &dummy.ByteReadCloser{Buffer: bytes.NewBuffer([]byte("abc"))}, nil
	case filepath.FromSlash("x/foo/bar"):
		return &dummy.ByteReadCloser{Buffer: bytes.NewBuffer([]byte("def"))}, nil
	case filepath.FromSlash("x/fuu/baaar"):
		return &dummy.ByteReadCloser{Buffer: bytes.NewBuffer([]byte("ghi"))}, nil
	case filepath.FromSlash("x/fuu/moo/meow/bla"):
		return &dummy.ByteReadCloser{Buffer: bytes.NewBuffer([]byte("jkl"))}, nil
	}
	return nil, fmt.Errorf("File not found")
}

func dummyStatFile(filePath string) (fs.FileInfo, error) {
	switch filePath {
	case filepath.FromSlash("x"):
		return dummy.NewFileInfo("x", true), nil
	}
	return nil, fmt.Errorf("File not found")
}
