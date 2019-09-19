package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func determineProgramPath() (string, error) {
	binaryDirPath := filepath.Dir(binaryPath)
	binaryDirName := filepath.Base(binaryDirPath)
	if binaryDirName != "MacOS" {
		return binaryPath, fmt.Errorf(`Binary "%s" not contained in a folder called "MacOS", but "%s"`, binaryPath, binaryDirName)
	}
	contentsDirPath := filepath.Dir(binaryDirPath)
	contentsDirName := filepath.Base(contentsDirPath)
	if contentsDirName != "Contents" {
		binaryName := filepath.Base(binaryPath)
		return binaryPath, fmt.Errorf(`Binary "%s" not contained in a proper application bundle. Expected "Contents/MacOS/%s" `+
			`but found "%s/%s/%s"`, binaryPath, binaryName, contentsDirName, binaryDirName, binaryName)
	}
	appDirPath := filepath.Dir(contentsDirPath)
	appDirName := filepath.Base(appDirPath)
	if !strings.HasSuffix(appDirName, ".app") {
		return binaryPath, fmt.Errorf(`Binary "%s" not contained in an application bundle ending on ".app". Found "%s"`,
			binaryPath, appDirPath)
	}
	return appDirPath, nil
}

func undeployProgram(undeployProgramPath string) (string, error) {
	return undeployProgramPath, fmt.Errorf("Program undeployment not required on Darwin")
}

func deleteProgram(deleteAppPath string) error {
	err := os.RemoveAll(deleteAppPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Could not delete application bundle at \"%s\": %v", deleteAppPath, err)
	}
	return nil
}
