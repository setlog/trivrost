package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/setlog/trivrost/pkg/launcher/config"
)

func fopen(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

func stat(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}

func MustHash(hashFilePath string) config.FileInfoMap {
	log.Infof("Hash \"%s\".", hashFilePath)
	return mustHashRelatively(ioutil.ReadDir, fopen, stat, hashFilePath)
}

type readDirFunc func(dirPath string) ([]os.FileInfo, error)
type readFileFunc func(filePath string) (io.ReadCloser, error)
type statFunc func(filePath string) (os.FileInfo, error)

func mustHashRelatively(readDir readDirFunc, readFile readFileFunc, stat statFunc, hashFilePath string) config.FileInfoMap {
	info, err := stat(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(fmt.Sprintf("Failed hashing file or folder \"%s\": %v", hashFilePath, err))
	}
	if !info.IsDir() {
		fileInfo := &config.FileInfo{}
		fileInfo.SHA256, fileInfo.Size, err = calculateSha256(hashFilePath, readFile)
		if err != nil {
			panic(fmt.Sprintf("failed hashing file \"%s\": %v", hashFilePath, err))
		}
		return config.FileInfoMap{"": fileInfo}
	}
	fileMap := mustHashDir(readDir, readFile, stat, hashFilePath)
	fileMapR := make(config.FileInfoMap)
	for k, v := range fileMap {
		rel, err := filepath.Rel(hashFilePath, k)
		if err != nil {
			panic(fmt.Sprintf("Could not create relative path for \"%s\" in \"%s\": %v", k, hashFilePath, err))
		}
		fileMapR[rel] = v
	}
	return fileMapR
}

func mustHashDir(readDir readDirFunc, readFile readFileFunc, stat statFunc, hashFilePath string) config.FileInfoMap {
	fm := make(config.FileInfoMap)
	for _, info := range mustReadDir(readDir, hashFilePath) {
		if info.IsDir() {
			fm.Join(mustHashDir(readDir, readFile, stat, filepath.Join(hashFilePath, info.Name())))
		} else {
			filePath := filepath.Join(hashFilePath, info.Name())
			sha, size, err := calculateSha256(filePath, readFile)
			if err != nil {
				panic(fmt.Sprintf("failed hashing file \"%s\": %v", hashFilePath, err))
			}
			fm[filePath] = &config.FileInfo{SHA256: sha, Size: size}
		}
	}
	return fm
}

func mustReadDir(readDir readDirFunc, directoryPath string) []os.FileInfo {
	infos, err := readDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(fmt.Sprintf("Could not list directory \"%s\": %v", directoryPath, err))
	}
	return infos
}

func calculateSha256(filePath string, readFile readFileFunc) (sha string, n int64, err error) {
	file, err := readFile(filePath)
	if err != nil {
		return "", n, fmt.Errorf("could not open file \"%s\": %v", filePath, err)
	}
	defer file.Close()
	hash := sha256.New()
	if n, err = io.Copy(hash, file); err != nil {
		return "", n, fmt.Errorf("could not read file \"%s\": %v", filePath, err)
	}
	shaSlice := hash.Sum(nil)
	return hex.EncodeToString(shaSlice), n, nil
}
