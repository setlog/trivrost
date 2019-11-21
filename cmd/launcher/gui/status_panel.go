package gui

import (
	"github.com/andlabs/ui"
)

type DownloadStatusPanel struct {
	*ui.Box

	mainVerticalBox *ui.Box
	inlineStatusBox *ui.Box
	pauseStatusBox  *ui.Box

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
	setBarProgress(panel.barTotalProgress, -1)
	panel.labelStatus = ui.NewLabel("")

	panel.Box.Append(panel.labelStage, false)
	panel.Box.Append(panel.barTotalProgress, false)

	panel.pauseStatusBox = ui.NewVerticalBox()
	panel.pauseStatusBox.Hide()
	panel.Box.Append(panel.pauseStatusBox, false)

	panel.inlineStatusBox = ui.NewHorizontalBox()
	panel.inlineStatusBox.Append(panel.labelStatus, false)
	panel.inlineStatusBox.Append(newLogsLinkLabel(), true)
	panel.Box.Append(panel.inlineStatusBox, false)

	return panel
}

func newLogsLinkLabel() *ui.Area {
	return newLinkLabel("Show logs...", ui.DrawTextAlignRight, showLogFolder)
}
