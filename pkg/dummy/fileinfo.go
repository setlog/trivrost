package dummy

import (
	"io/fs"
	"os"
	"time"
)

func NewFileInfo(name string, isDir bool) fs.FileInfo {
	return &FileInfo{name: name, isDir: isDir}
}

type FileInfo struct {
	name  string
	isDir bool
}

// Satisfy fs.FileInfo interface requirements.
func (dfi *FileInfo) Name() string {
	return dfi.name
}

func (dfi *FileInfo) Size() int64 {
	return 42
}

func (dfi *FileInfo) Mode() os.FileMode {
	return 0755
}

func (dfi *FileInfo) ModTime() time.Time {
	return time.Now()
}

func (dfi *FileInfo) IsDir() bool {
	return dfi.isDir
}

func (dfi *FileInfo) Sys() any {
	return nil
}
