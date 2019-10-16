package system

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

const (
	OsWindows = "windows"
	OsMac     = "darwin"
	OsLinux   = "linux"

	Arch64 = "amd64"
	Arch32 = "386"
)

var is64BitOS bool

type ProcessSignature struct {
	Pid        int   `json:"Pid"`
	CreateTime int64 `json:"CreateTime"`
}

func init() {
	mustDetectArchitecture()
}

func GetCurrentProcessSignature() *ProcessSignature {
	return GetPidProcessSignature(os.Getpid())
}

func GetPidProcessSignature(pid int) *ProcessSignature {
	createTime, _ := GetPidCreateTime(pid)
	return &ProcessSignature{Pid: pid, CreateTime: createTime}
}

func GetPidCreateTime(pid int) (time int64, err error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, fmt.Errorf("Could not open process: %v", err)
	}
	return proc.CreateTime()
}

// Reports the time at which the given process launched in some format, with at least millisecond precision.
func GetProcessCreateTime(proc *os.Process) (time int64, err error) {
	return GetPidCreateTime(proc.Pid)
}

func IsProcessSignatureRunning(procSig *ProcessSignature) bool {
	if IsPidRunning(procSig.Pid) {
		// Process with proc.Pid is running. Is it also the one we are looking for?
		proc, err := process.NewProcess(int32(procSig.Pid))
		if err != nil {
			return false
		}
		createTime, err := proc.CreateTime()
		if (err != nil) || (createTime != procSig.CreateTime) {
			return false
		}
		return true
	}
	return false
}

func IsPidRunning(pid int) bool {
	p, err := os.FindProcess(pid) // Only can return an error on Windows, indicating that no process with given pid is running.
	if err != nil {
		return false
	}
	defer func() {
		err := p.Release()
		if err != nil {
			log.Warn(err)
		}
	}()
	return isProcessRunning(p)
}

// Returns "amd64" if the underlying OS is 64 bit. Returns "386" otherwise. This is compliant with the GOARCH naming scheme.
// Do not confuse the result of this function with runtime.GOARCH, which describes the architecture this binary was built for instead.
func GetOSArch() string {
	if Is64BitOS() {
		return Arch64
	}
	return Arch32
}

func Is64BitOS() bool {
	return is64BitOS
}

func MatchesPlatform(platform string, os string, arch string) bool {
	return (platform == os) || (platform == os+"-"+arch) || (platform == arch)
}

func StartProcess(binaryPath string, workingDirectoryPath string, args []string, extraEnvironmentVariables map[string]*string, passStreams bool) (*exec.Cmd, *ProcessSignature, error) {
	command := exec.Command(binaryPath, args...)
	command.Dir = MustGetAbsolutePath(workingDirectoryPath)
	command.Env = buildEnvironmentVariables(extraEnvironmentVariables)
	if passStreams {
		// Reason for this condition: relaying standard streams causes a crash on Windows if the launcher was built
		// as a console app but tries to start a GUI app (such as javaw.exe).
		command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
	}

	err := command.Start()
	if err != nil {
		return command, nil, fmt.Errorf("Could not start process \"%s\" with working directory \"%s\": %w", binaryPath, workingDirectoryPath, err)
	}

	procSig := &ProcessSignature{Pid: command.Process.Pid}
	procSig.CreateTime, err = GetProcessCreateTime(command.Process)
	if err != nil {
		log.WithFields(log.Fields{"pid": procSig.Pid}).Warnf("Could not get creation time of created process: %v", err)
	}

	return command, procSig, nil
}

func buildEnvironmentVariables(extraEnvs map[string]*string) []string {
	osEnvs := os.Environ()
	patchedEnvs := make([]string, len(osEnvs))
	copy(patchedEnvs, osEnvs)

	for newEnvKey, newEnvValue := range extraEnvs {
		if newEnvValue == nil {
			patchedEnvs = removeEnv(patchedEnvs, newEnvKey)
		} else {
			newEnv := newEnvKey + "=" + *newEnvValue
			patchedEnvs = removeEnv(patchedEnvs, newEnvKey)
			patchedEnvs = append(patchedEnvs, newEnv)
		}
	}

	return patchedEnvs
}

func ShowLocalFileInFileManager(path string) error {
	return showLocalFileInFileManager(path)
}
