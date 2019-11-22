package gui

import (
	"runtime"

	"github.com/andlabs/ui"
	"github.com/setlog/trivrost/pkg/system"
)

func setBarProgress(bar *ui.ProgressBar, progress int) {
	if progress < 0 {
		bar.SetValue(-1)
	} else {
		// See https://stackoverflow.com/questions/2217688/windows-7-aero-theme-progress-bar-bug for the Windows cases.
		if progress >= 100 {
			if runtime.GOOS == system.OsWindows {
				bar.SetValue(100)
				bar.SetValue(99)
				bar.SetValue(100)
			} else {
				bar.SetValue(100)
			}
		} else {
			if runtime.GOOS == system.OsWindows {
				bar.SetValue(progress + 1)
			}
			bar.SetValue(progress)
		}
	}
}
