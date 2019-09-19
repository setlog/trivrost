// +build !windows

package logging

func enableVirtualTerminalProcessing() error {
	// Nothing to do on Unix
	return nil
}
