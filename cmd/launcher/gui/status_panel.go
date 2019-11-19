package gui

import (
	"github.com/andlabs/ui"
)

type DownloadStatusPanel struct {
	*ui.Box

	labelStage *ui.Label

	barTotalProgress *ui.ProgressBar
	labelStatus      *ui.Label

	progressPrevious, progressCurrent, progressTarget uint64 // Whether these refer to amount of bytes downloaded or something else depends on the current GUI stage.
	currentProblemMessage                             string
	stage                                             Stage
}

func newDownloadStatusPanel() *DownloadStatusPanel {
	panel := &DownloadStatusPanel{Box: ui.NewVerticalBox()}
	panel.SetPadded(true)

	panel.labelStage = ui.NewLabel("Initializing...")
	panel.barTotalProgress = ui.NewProgressBar()
	panel.barTotalProgress.SetValue(-1)
	panel.labelStatus = ui.NewLabel("")

	panel.Box.Append(panel.labelStage, false)
	panel.Box.Append(panel.barTotalProgress, false)

	hBox := ui.NewHorizontalBox()
	hBox.Append(panel.labelStatus, false)
	hBox.Append(newLinkLabel("Show logs...", showLogFolder), true)
	panel.Box.Append(hBox, false)

	return panel
}
