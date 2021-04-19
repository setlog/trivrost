// +build linux darwin

package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

func mustDetectArchitecture() {
	cmd := exec.Command("uname", "-m")
	data, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("Could not determine system architecture with \"uname -m\": %v", err))
	}
	is64BitOS = runtime.GOARCH == "amd64" || strings.Contains(string(data), "x86_64")
}

func removeEnv(envs []string, name string) []string {
	for i := 0; i < len(envs); i++ {
		env := envs[i]
		kv := strings.SplitN(env, "=", 2)
		if len(kv) == 2 && kv[0] == name {
			envs = append(envs[:i], envs[i+1:]...)
			break
		}
	}
	return envs
}

func showLocalFileInFileManager(path string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == OsMac {
		cmd = exec.Command("open", filepath.Dir(path))
	} else {
		cmd = exec.Command("xdg-open", filepath.Dir(path))
	}
	return cmd.Run()
}

func isProcessRunning(p *os.Process) bool {
	return p.Signal(unix.Signal(0)) == nil
}

func universalPathName(p string) (string, error) {
	return p, nil
}
