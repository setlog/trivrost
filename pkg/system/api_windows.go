package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

// #cgo LDFLAGS: -lMpr
//#include <windows.h>
//#include <winnetwk.h>
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

func universalPathName(p string) (string, error) {
	s, lpBufferSize, err := universalPathNameWithBufferSize(p, 1000)
	if err != nil && err.(*universalNameRetrievalError).ErrorType() == errorMoreData {
		s, _, err = universalPathNameWithBufferSize(s, lpBufferSize)
	}
	if err != nil {
		return p, err
	}
	return s, err
}

func universalPathNameWithBufferSize(p string, lpBufferSizeUse C.DWORD) (universalPath string, lpBufferSize C.DWORD, err error) {
	cp := C.LPCWSTR(StringToUTF16UnmanagedString(p))
	defer C.free(unsafe.Pointer(cp))

	// The possible data written to infoStruct (we request a UNIVERSAL_NAME_INFO below) not only consists of the struct, but also of the data (strings)
	// pointed to by pointer-members within the struct. That's why this allocation needs to be much larger than just large enough to hold the struct itself.
	infoStruct := C.LPVOID(C.calloc(C.size_t(lpBufferSizeUse), 1))
	defer C.free(unsafe.Pointer(infoStruct))

	lpBufferSize = lpBufferSizeUse
	errorCode := C.WNetGetUniversalNameW(cp, C.UNIVERSAL_NAME_INFO_LEVEL, infoStruct, &lpBufferSize)
	err = getErrorOfWNetGetUniversalNameW(errorCode)
	if err == nil {
		lpUniversalName := unsafe.Pointer(*(*C.LPWSTR)(infoStruct))
		universalPath = UTF16StringToString(lpUniversalName)
	}
	return universalPath, lpBufferSize, err
}

func getErrorOfWNetGetUniversalNameW(returnCode C.DWORD) error {
	if returnCode == C.NO_ERROR {
		return nil
	}
	if returnCode == C.ERROR_BAD_DEVICE {
		return &universalNameRetrievalError{errorType: errorBadDevice,
			message: `the string pointed to by the lpLocalPath parameter is invalid`}
	}
	if returnCode == C.ERROR_CONNECTION_UNAVAIL {
		return &universalNameRetrievalError{errorType: errorConnectionUnavailable,
			message: `there is no current connection to the remote device, but there is a remembered (persistent) connection to it`}
	}
	if returnCode == C.ERROR_EXTENDED_ERROR {
		errorMessage, providerName, err := getLastWNetError()
		if err != nil {
			return &universalNameRetrievalError{errorType: errorExtendedError,
				message: `a network-specific error occurred; getting extended error information failed: ` + err.Error()}
		}
		return &universalNameRetrievalError{errorType: errorExtendedError,
			message: `a network-specific error occurred; Network provider "` + providerName + `" reports: ` + errorMessage}
	}
	if returnCode == C.ERROR_MORE_DATA {
		return &universalNameRetrievalError{errorType: errorMoreData,
			message: `despite trying to query with the requested buffer size, the buffer pointed to by the lpBuffer parameter was too small`}
	}
	if returnCode == C.ERROR_NOT_SUPPORTED {
		return &universalNameRetrievalError{errorType: errorNotSupported,
			message: `the dwInfoLevel parameter is set to UNIVERSAL_NAME_INFO_LEVEL, but the network provider does not support UNC names. (None of the network providers support this function)`}
	}
	if returnCode == C.ERROR_NO_NET_OR_BAD_PATH {
		return &universalNameRetrievalError{errorType: errorNoNetOrBadPath,
			message: `none of the network providers recognize the local name as having a connection. However, the network is not available for at least one provider to whom the connection may belong`}
	}
	if returnCode == C.ERROR_NO_NETWORK {
		return &universalNameRetrievalError{errorType: errorNoNetwork,
			message: `the network is unavailable`}
	}
	if returnCode == C.ERROR_NOT_CONNECTED {
		return &universalNameRetrievalError{errorType: errorNotConnected,
			message: `the device specified by the path is not redirected`}
	}
	return &universalNameRetrievalError{errorType: errorUndocumented,
		message: fmt.Sprintf(`undocumented error code %d`, returnCode)}
}

func getLastWNetError() (errorMessage, providerName string, err error) {
	var lpError C.DWORD

	const errorBufferSize = 5000
	const nErrorBufSize C.DWORD = errorBufferSize
	lpErrorBuf := (C.LPWSTR)(C.calloc(C.size_t(errorBufferSize+1), C.size_t(unsafe.Sizeof(uint16(0)))))
	defer C.free(unsafe.Pointer(lpErrorBuf))

	const nameBufferSize = 1000
	const nNameBufSize C.DWORD = nameBufferSize
	lpNameBuf := (C.LPWSTR)(C.calloc(C.size_t(nameBufferSize+1), C.size_t(unsafe.Sizeof(uint16(0)))))
	defer C.free(unsafe.Pointer(lpNameBuf))

	returnCode := C.WNetGetLastErrorW(&lpError, lpErrorBuf, nErrorBufSize, lpNameBuf, nNameBufSize)
	if returnCode == C.NO_ERROR {
		return UTF16StringToString(unsafe.Pointer(lpErrorBuf)), UTF16StringToString(unsafe.Pointer(lpNameBuf)), nil
	}
	if returnCode == C.ERROR_INVALID_ADDRESS {
		return "", "", fmt.Errorf("could not get last WNet error: ERROR_INVALID_ADDRESS")
	}
	return "", "", fmt.Errorf("could not get last WNet error: undocumented extended error code %d", returnCode)
}

// StringToUTF16UnmanagedString returns an unmanaged, null-terminated UTF16 string for given string.
// The caller is responsible for freeing the returned pointer.
func StringToUTF16UnmanagedString(s string) unsafe.Pointer {
	utf16String := utf16.Encode([]rune(s))
	utf16StringPointer := (*uint16)(C.calloc(C.size_t(len(utf16String)+1), C.size_t(unsafe.Sizeof(uint16(0)))))
	currentCharPointer := utf16StringPointer
	for _, c := range utf16String {
		*currentCharPointer = c
		currentCharPointer = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(currentCharPointer)) + unsafe.Sizeof(uint16(0))))
	}
	return unsafe.Pointer(utf16StringPointer)
}

// UTF16StringToString returns a string for a given null-terminated UTF16 string.
// This function does not call free on the parameter.
func UTF16StringToString(lpwString unsafe.Pointer) string {
	ptr := (*uint16)(lpwString)
	data := make([]uint16, 0, 0)
	for {
		if *ptr == 0 {
			break
		}
		data = append(data, *ptr)
		ptr = (*uint16)(unsafe.Pointer(((uintptr)(unsafe.Pointer(ptr))) + unsafe.Sizeof(uint16(0))))
	}
	s := utf16.Decode(data)
	return string(s)
}
