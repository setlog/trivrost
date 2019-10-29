package system

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var binaryPath string
var programPath string

func FindPaths() (err error) {
	binaryPath, err = os.Executable()
	if err != nil {
		return fmt.Errorf("Could not get path to binary: %v", err)
	}
	evaluatedBinaryPath, err := filepath.EvalSymlinks(binaryPath)
	if err != nil {
		log.Printf("Could not evaluate path to binary \"%s\": %v", binaryPath, err)
	} else {
		binaryPath = evaluatedBinaryPath
	}
	programPath, err = determineProgramPath()
	if err != nil {
		return fmt.Errorf("Could not determine path of application bundle for binary path \"%s\": %v", binaryPath, err)
	}
	return nil
}

// GetProgramPath returns the path where the calling program lies.
// This will be identical to the result of GetBinaryPath() unless the OS is MacOS,
// where GetProgramPath() will return the path of the program's application bundle folder.
func GetProgramPath() string {
	return programPath
}

// GetBinaryPath returns the path where the binary of the calling program lies.
func GetBinaryPath() string {
	return binaryPath
}

// UndeployProgram renames the file at undeployProgramPath on operating systems
// where running binaries cannot be deleted (i.e. Windows) and returns the path to the
// renamed file. Otherwise, the provided string is returned with a non-nil error.
func UndeployProgram(undeployProgramPath string) (string, error) {
	return undeployProgram(undeployProgramPath)
}

// DeleteProgram immediately deletes the binary or MacOS application bundle at deleteProgramPath,
// unless the operating system is Windows; in that case the file at deleteProgramPath will be deleted after
// 3 seconds.
func DeleteProgram(deleteProgramPath string) error {
	return deleteProgram(deleteProgramPath)
}
