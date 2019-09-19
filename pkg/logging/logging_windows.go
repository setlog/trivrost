package logging

import (
	"fmt"

	"golang.org/x/sys/windows"
)

// Enables parsing of ANSI escape sequences (for color and such) in Windows 10 cmd.exe.
const enableVirtualTerminalProcessingFlag uint32 = 0x0004

// var stdOutputHandleIdentifier int = -11
var stdErrorHandleIdentifier int = -12

func enableVirtualTerminalProcessing() error {
	v := windows.RtlGetVersion()
	if v.MajorVersion < 10 {
		return fmt.Errorf("Windows Version %d.%d does not support virtual terminal processing. Version 10.0 is needed", v.MajorVersion, v.MinorVersion)
	}

	handle, err := windows.GetStdHandle(uint32(stdErrorHandleIdentifier))
	if err != nil {
		return fmt.Errorf("windows.GetStdHandle() failed: %v", err)
	}
	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return fmt.Errorf("GetConsoleMode() failed: %v", err)
	}
	newMode := mode | enableVirtualTerminalProcessingFlag
	err = windows.SetConsoleMode(handle, newMode)
	if err != nil {
		errCode, ok := err.(windows.Errno)
		if !ok || errCode != 0 {
			return fmt.Errorf("SetConsoleMode() failed: %v. Code: %d; Old mode: %d; New Mode: %d", err, errCode, mode, newMode)
		}
	}
	return nil
}
