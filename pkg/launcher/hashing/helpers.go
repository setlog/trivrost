package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func MustCalculateFileHash(filePath string) (sha string, n int64) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Could not open file \"%s\" for reading: %v", filePath, err))
	}
	defer file.Close()
	return MustCalculateHash(file)
}

func MustCalculateHash(file io.Reader) (sha string, n int64) {
	hash := sha256.New()
	n, err := io.Copy(hash, file)
	if err != nil {
		panic(fmt.Sprintf("Could not read from reader: %v", err))
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), n
}

func MustDecodeSha256String(shaString string) (sha [sha256.Size]byte) {
	shaSlice, err := hex.DecodeString(shaString)
	if err != nil {
		panic(err)
	}
	if len(shaSlice) != sha256.Size {
		panic(fmt.Sprintf("shaString has length %d when %d is required.", len(shaSlice), sha256.Size))
	}
	copy(sha[:], shaSlice)
	return sha
}
