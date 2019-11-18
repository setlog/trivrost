package system

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/misc"
)

type FileSystemError struct {
	message      string
	causingError error
}

func (fse *FileSystemError) Error() string {
	var causingErrorMessage string
	if fse.causingError == nil {
		causingErrorMessage = "<nil>"
	} else {
		causingErrorMessage = fse.causingError.Error()
	}
	if fse == nil {
		return "<nil>: " + causingErrorMessage
	}
	return fse.message + ": " + causingErrorMessage
}

func (fse *FileSystemError) Unwrap() error {
	return fse.causingError
}

func NewFileSystemError(message string, cause error) *FileSystemError {
	return &FileSystemError{message: message, causingError: cause}
}

func MustMakeTempDirectory(forFolderAtPath string) string {
	randomHex := misc.MustGetRandomHexString(8)
	tempDirPath := filepath.Join(filepath.Dir(forFolderAtPath), "~"+filepath.Base(forFolderAtPath)+".dl."+randomHex)
	err := os.MkdirAll(tempDirPath, 0700)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not create temporary directory \"%s\" for directory \"%s\"", tempDirPath, forFolderAtPath), err})
	}
	return tempDirPath
}

func TryRemoveDirectory(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Warnf("Could not remove directory \"%s\": %v", path, err)
	}
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(&FileSystemError{fmt.Sprintf("Could not stat \"%s\"", path), err})
	}
	return info.IsDir()
}

// Recursively moves all content of fromDirectory into toDirectory, overwriting existing files when encountered.
func MustMoveFiles(fromDirectory, toDirectory string) {
	infos, err := ioutil.ReadDir(fromDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			panic(&FileSystemError{fmt.Sprintf("Cannot move content of directory \"%s\" into \"%s\", because the former does not exist", fromDirectory, toDirectory), err})
		}
		panic(&FileSystemError{fmt.Sprintf("Could not read content of directory \"%s\" to move it into \"%s\"", fromDirectory, toDirectory), err})
	}
	for _, info := range infos {
		sourcePath := filepath.Join(fromDirectory, info.Name())
		destPath := filepath.Join(toDirectory, info.Name())
		if info.IsDir() {
			MustMoveFiles(sourcePath, destPath)
		} else {
			err = os.MkdirAll(toDirectory, 0700)
			if err != nil {
				panic(&FileSystemError{fmt.Sprintf("Could not create nested directory structure \"%s\" to move file \"%s\" to \"%s\" in \"%s\"",
					toDirectory, sourcePath, destPath, fromDirectory), err})
			}
			err = os.Rename(sourcePath, destPath)
			if err != nil {
				panic(&FileSystemError{fmt.Sprintf("Could not move file \"%s\" to \"%s\" in attempt to move content of \"%s\" into \"%s\"",
					sourcePath, destPath, fromDirectory, toDirectory), err})
			}
		}
	}
}

func MustMakeDir(dirPath string) {
	err := os.MkdirAll(dirPath, 0744)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not create nested directory structure \"%s\"", dirPath), err})
	}
}

func MustPutFile(localFilePath string, bytes []byte) {
	os.Remove(localFilePath)
	dir := filepath.Dir(localFilePath)
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not create nested directory structure \"%s\" to put file \"%s\"", dir, localFilePath), err})
	}
	file, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not open file \"%s\" for writing", localFilePath), err})
	}
	defer file.Close()

	log.WithFields(log.Fields{
		"localFilePath": localFilePath, "dir": dir, "length": len(bytes)}).Debug("Writing file.")

	_, err = file.Write(bytes)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not write file \"%s\"", localFilePath), err})
	}
}

func MustReadFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not read file \"%s\"", filePath), err})
	}
	return data
}

func MustCopyFile(from, to string) {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not read file \"%s\"", from), err})
	}
	info, err := os.Stat(from)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not stat file \"%s\"", from), err})
	}
	log.Debugf(`Copying "%s" to "%s" with mode %s.`, from, to, strconv.FormatInt(int64(info.Mode()), 8))
	err = ioutil.WriteFile(to, data, info.Mode())
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not write file \"%s\"", to), err})
	}
}

// Move the file or folder at src to dst. If dst is taken by an existing file or folder, it will be removed beforehand.
func MustMoveAll(src, dst string) {
	srcInfo, _ := mustPrepareFileSystemOperation(src, dst)
	mustMoveFile(src, dst, srcInfo.Mode())
}

func mustMoveFile(from, to string, mode os.FileMode) {
	err := os.MkdirAll(filepath.Dir(to), 0700|mode)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf(`Could not move "%s" to "%s": MkdirAll() failed`, from, to), err})
	}
	err = os.Rename(from, to)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf(`Could not move "%s" to "%s": Rename() failed`, from, to), err})
	}
}

// Copy the file or folder at src to dst. If dst is taken by an existing file or folder, it will be removed beforehand.
func MustCopyAll(src, dst string) {
	srcInfo, _ := mustPrepareFileSystemOperation(src, dst)
	mustCopyAll(src, dst, srcInfo.IsDir(), srcInfo.Mode())
}

func mustPrepareFileSystemOperation(src, dst string) (srcInfo, dstInfo os.FileInfo) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf(`Could not prepare file system operation from "%s" to "%s": Stat() failed on source`, src, dst), err})
	}
	dstInfo, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		panic(&FileSystemError{fmt.Sprintf(`Could not prepare file system operation from "%s" to "%s": Stat() failed on destination`, src, dst), err})
	}
	if err == nil {
		log.Debugf(`Removing "%s".`, dst)
		err = os.RemoveAll(dst)
		if err != nil {
			panic(&FileSystemError{fmt.Sprintf(`Could not prepare file system operation from "%s" to "%s": RemoveAll() failed on destination`, src, dst), err})
		}
	}
	return srcInfo, dstInfo
}

func mustCopyAll(src, dst string, srcIsDir bool, mode os.FileMode) {
	if srcIsDir {
		err := os.MkdirAll(dst, 0700|mode)
		if err != nil {
			panic(&FileSystemError{fmt.Sprintf(`Could not copy "%s" to "%s": MkdirAll() failed on destination`, src, dst), err})
		}
		mustCopyDir(src, dst)
	} else {
		err := os.MkdirAll(filepath.Dir(dst), 0700|mode)
		if err != nil {
			panic(&FileSystemError{fmt.Sprintf(`Could not copy "%s" to "%s": Mkdir() failed on destination parent folder`, src, dst), err})
		}
		MustCopyFile(src, dst)
	}
}

func mustCopyDir(src, dst string) {
	contentInfos, err := ioutil.ReadDir(src)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf(`Could not copy "%s" to "%s": ReadDir() failed`, src, dst), err})
	}
	for _, contentInfo := range contentInfos {
		mustCopyAll(filepath.Join(src, contentInfo.Name()), filepath.Join(dst, contentInfo.Name()), contentInfo.IsDir(), contentInfo.Mode())
	}
}

func MustRemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		panic(&FileSystemError{fmt.Sprintf("Could not delete file \"%s\"", filePath), err})
	}
}

func MustRecursivelyRemoveEmptyFolders(folder string) bool {
	infos, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not read content of directory \"%s\"", folder), err})
	}
	removeCount := 0
	for _, info := range infos {
		if info.IsDir() {
			if MustRecursivelyRemoveEmptyFolders(filepath.Join(folder, info.Name())) {
				removeCount++
			}
		}
	}
	if removeCount == len(infos) {
		err = os.Remove(folder)
		if err != nil {
			panic(&FileSystemError{fmt.Sprintf("Could not remove empty directory \"%s\"", folder), err})
		}
		return true
	}
	return false
}

func MustGetAbsolutePath(filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(&FileSystemError{fmt.Sprintf("Could not make absolute path from \"%s\"", filePath), err})
	}
	return absolutePath
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(&FileSystemError{fmt.Sprintf("Could not check if file \"%s\" exists", filePath), err})
	}
	return !info.IsDir()
}

func FolderExists(folderPath string) bool {
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(&FileSystemError{fmt.Sprintf("Could not check if folder \"%s\" exists", folderPath), err})
	}
	return info.IsDir()
}

func IsEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		panic(&FileSystemError{fmt.Sprintf("Could not check if file or folder \"%s\" is empty", filePath), err})
	}
	if !info.IsDir() {
		return info.Size() == 0
	}
	infos, err := ioutil.ReadDir(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		panic(&FileSystemError{fmt.Sprintf("Could not check if file or folder \"%s\" is empty", filePath), err})
	}
	return len(infos) == 0
}

func TryRemove(filePath string) {
	log.Debugf("Removing file \"%s\".", filePath)
	err := os.Remove(filePath)
	if err != nil {
		log.Errorf("Could not remove file \"%s\": %v", filePath, err)
	}
}

func TryRemoveEmpty(filePath string) {
	if IsEmpty(filePath) {
		log.Debugf("Removing folder \"%s\".", filePath)
		err := os.Remove(filePath)
		if err != nil {
			log.Errorf("Could not remove folder \"%s\": %v", filePath, err)
		}
	}
}

func CleanUpFileOperation(file *os.File, returnError *error) {
	if file != nil {
		filePath := file.Name()
		file.Close()
		if *returnError != nil {
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				log.Warnf("Could not remove file \"%s\" after error: %v", filePath, err)
			}
		}
	}
}
