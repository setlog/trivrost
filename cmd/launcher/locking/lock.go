package locking

import (
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/setlog/trivrost/cmd/launcher/places"
)

const lockFileName = ".lock" // See https://github.com/gofrs/flock/issues/42

var fileLock *flock.Flock

func lockFilePath() string {
	return filepath.Join(places.GetAppLocalDataFolderPath(), lockFileName)
}

// Releases the lock.
func ReleaseLock() {
	if fileLock != nil && fileLock.Locked() {
		err := fileLock.Unlock()
		if err != nil {
			panic(err)
		}
	}
}

func mustTryLock() bool {
	if fileLock == nil {
		fileLock = flock.New(lockFilePath())
	}

	didAcquireLock, err := fileLock.TryLock()
	if err != nil {
		panic(err)
	}
	return didAcquireLock
}
