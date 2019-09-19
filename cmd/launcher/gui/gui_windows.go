package gui

import (
	"unicode/utf16"
	"unsafe"

	"github.com/setlog/trivrost/pkg/system"

	log "github.com/sirupsen/logrus"

	// #include "gui_windows.h"
	"C"
)

var didLoadIcons bool

func centerWindow(handle uintptr) {
	result := C.centerWindow(C.ULONG_PTR(handle))
	if result != 0 {
		if result == 1 {
			log.Errorf("Could not center window: getting monitor from window failed.")
		} else if result == 2 {
			log.Errorf("Could not center window: getting monitor info failed.")
		} else if result == 3 {
			log.Errorf("Could not center window: getting window rect failed.")
		}
	}
}

func applyIconToWindow(handle uintptr) {
	if !didLoadIcons {
		loadIcons()
	}
	C.applyIconToWindow(C.ULONG_PTR(handle))
}

func applyWindowStyle(handle uintptr) {
	C.applyWindowStyle(C.ULONG_PTR(handle))
}

func loadIcons() {
	binaryPath := goStringToConstantUTF16WinApiString(system.GetBinaryPath())
	extractedIconCount := C.loadIcons(binaryPath)
	didLoadIcons = true
	C.free(unsafe.Pointer(binaryPath))
	if extractedIconCount == 0 {
		log.Errorf("Extracted no icons. Expected 2.")
	} else if extractedIconCount == 1 {
		log.Warnf("Extracted only one icon. Expected 2.")
	} else if extractedIconCount == 2 {
		log.Debugf("Extracted 2 icons.")
	} else {
		log.Errorf("Extracted %d icons. Expected 2.", extractedIconCount)
	}
}

func goStringToConstantUTF16WinApiString(s string) C.LPCWSTR {
	utf16String := utf16.Encode([]rune(s))
	utf16StringPointer := (*uint16)(C.calloc(C.size_t(len(utf16String)+1), C.size_t(unsafe.Sizeof(uint16(0)))))
	currentCharPointer := utf16StringPointer
	for _, c := range utf16String {
		*currentCharPointer = c
		currentCharPointer = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(currentCharPointer)) + unsafe.Sizeof(uint16(0))))
	}
	return (C.LPCWSTR)(unsafe.Pointer(utf16StringPointer))
}
