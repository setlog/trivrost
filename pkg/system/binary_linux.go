package system

import (
	"fmt"
	"os"
)

func determineProgramPath() (string, error) {
	return binaryPath, nil
}

func undeployProgram(undeployProgramPath string) (string, error) {
	return undeployProgramPath, fmt.Errorf("Program undeployment not required on Linux")
}

func deleteProgram(deleteAppPath string) error {
	err := os.Remove(deleteAppPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Could not delete binary at \"%s\": %v", deleteAppPath, err)
	}
	return nil
}
