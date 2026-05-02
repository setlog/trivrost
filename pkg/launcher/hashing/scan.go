package hashing

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/setlog/trivrost/pkg/launcher/config"
	"github.com/setlog/trivrost/pkg/misc"
	log "github.com/sirupsen/logrus"
)

func fopen(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

func stat(filePath string) (fs.FileInfo, error) {
	return os.Stat(filePath)
}

func MustHash(ctx context.Context, hashFilePath string) config.FileInfoMap {
	log.Infof("Hash \"%s\".", hashFilePath)
	return mustHashRelatively(ctx, readDir, fopen, stat, hashFilePath)
}

type readDirFunc func(dirPath string) ([]fs.FileInfo, error)
type readFileFunc func(filePath string) (io.ReadCloser, error)
type statFunc func(filePath string) (fs.FileInfo, error)

func readDir(dirPath string) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func mustHashRelatively(ctx context.Context, readDir readDirFunc, readFile readFileFunc, stat statFunc, hashFilePath string) config.FileInfoMap {
	info, err := stat(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(fmt.Errorf("Failed hashing file or folder \"%s\": %w", hashFilePath, err))
	}
	if !info.IsDir() {
		fileInfo := &config.FileInfo{}
		fileInfo.SHA256, fileInfo.Size, err = calculateSha256(ctx, hashFilePath, readFile)
		if err != nil {
			panic(fmt.Errorf("failed hashing file \"%s\": %w", hashFilePath, err))
		}
		return config.FileInfoMap{"": fileInfo}
	}
	fileMap := mustHashDir(ctx, readDir, readFile, stat, hashFilePath)
	fileMapR := make(config.FileInfoMap)
	for k, v := range fileMap {
		rel, err := filepath.Rel(hashFilePath, k)
		if err != nil {
			panic(fmt.Errorf("Could not create relative path for \"%s\" in \"%s\": %w", k, hashFilePath, err))
		}
		fileMapR[rel] = v
	}
	return fileMapR
}

func mustHashDir(ctx context.Context, readDir readDirFunc, readFile readFileFunc, stat statFunc, hashFilePath string) config.FileInfoMap {
	fm := make(config.FileInfoMap)
	for _, info := range mustReadDir(readDir, hashFilePath) {
		if info.IsDir() {
			fm.Join(mustHashDir(ctx, readDir, readFile, stat, filepath.Join(hashFilePath, info.Name())))
		} else {
			filePath := filepath.Join(hashFilePath, info.Name())
			sha, size, err := calculateSha256(ctx, filePath, readFile)
			if err != nil {
				panic(fmt.Errorf("failed hashing file \"%s\": %w", hashFilePath, err))
			}
			fm[filePath] = &config.FileInfo{SHA256: sha, Size: size}
		}
	}
	return fm
}

func mustReadDir(readDir readDirFunc, directoryPath string) []fs.FileInfo {
	infos, err := readDir(directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(fmt.Errorf("Could not list directory \"%s\": %w", directoryPath, err))
	}
	return infos
}

func CalculateSha256(ctx context.Context, filePath string) (sha string, n int64, err error) {
	return calculateSha256(ctx, filePath, fopen)
}

func calculateSha256(ctx context.Context, filePath string, readFile readFileFunc) (sha string, n int64, err error) {
	file, err := readFile(filePath)
	if err != nil {
		return "", n, fmt.Errorf("could not open file \"%s\": %w", filePath, err)
	}
	defer file.Close()
	hash := sha256.New()
	if n, err = misc.IOCopyWithContext(ctx, hash, file); err != nil {
		return "", n, fmt.Errorf("could not read file \"%s\": %w", filePath, err)
	}
	shaSlice := hash.Sum(nil)
	return hex.EncodeToString(shaSlice), n, nil
}
