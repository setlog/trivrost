package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
)

//#include <windows.h>
import "C"

func mustDetectArchitecture() {
	if runtime.GOARCH == Arch64 {
		is64BitOS = true
	} else {
		handle, err := windows.GetCurrentProcess()
		if err != nil {
			panic(fmt.Sprintf("Could not get current process handle: %v", err))
		}
		err = windows.IsWow64Process(handle, &is64BitOS)
		if err != nil {
			panic(fmt.Sprintf("Could not detect architecture: %v", err))
		}
	}
}

func removeEnv(envs []string, name string) []string {
	for i := 0; i < len(envs); i++ {
		env := envs[i]
		kv := strings.SplitN(env, "=", 2)
		if len(kv) == 2 && strings.EqualFold(kv[0], name) {
			envs = append(envs[:i], envs[i+1:]...)
			i--
		}
	}
	return envs
}

func showLocalFileInFileManager(path string) error {
	cmd := exec.Command("explorer", "/select,", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %v (this is often a false positive)", string(output), err)
	}
	return nil
}

func isProcessRunning(p *os.Process) bool {
	handle := C.OpenProcess(C.PROCESS_QUERY_INFORMATION, C.FALSE, C.DWORD(p.Pid))
	if handle == C.HANDLE(C.NULL) {
		return false
	}
	defer C.CloseHandle(handle)
	var lpExitCode C.DWORD
	result := C.GetExitCodeProcess(handle, &lpExitCode)
	return (result != 0) && (lpExitCode == C.STILL_ACTIVE)
}
