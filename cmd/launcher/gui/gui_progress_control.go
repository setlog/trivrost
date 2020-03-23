package gui

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/andlabs/ui"
	log "github.com/sirupsen/logrus"
)

// ProgressFunc should be set to a function which reports the progress of the given stage.
var ProgressFunc = func(s Stage) uint64 {
	return 0
}

// SetStage sets up the GUI to determine the progress bar value based on the progress
// interval of the given stage. When progressTotal is >0, you can set gui.ProgressFunc
// to a function which reports the current progress.
func SetStage(s Stage, progressTarget uint64) {
	log.Debugf("Changing stage to %v with total %d.\n", s, progressTarget)
	ui.QueueMain(func() {
		isStateChange := panelDownloadStatus.stage.IsWaitingStage() != s.IsWaitingStage()
		panelDownloadStatus.stage = s
		panelDownloadStatus.progressCurrent = 0
		panelDownloadStatus.progressPrevious = 0
		panelDownloadStatus.progressTarget = progressTarget
		panelDownloadStatus.labelStage.SetText(s.getText())
		barProgress, percentage := calculateProgress(panelDownloadStatus.stage, panelDownloadStatus.progressCurrent, panelDownloadStatus.progressTarget)
		window.SetTitle(fmt.Sprintf("[%.1f%%] %s", percentage, windowTitle))
		panelDownloadStatus.currentProblemMessage = ""
		panelDownloadStatus.labelStatus.SetText("")
		if isStateChange {
			if s.IsWaitingStage() {
				setProgressState(statePaused)
			} else {
				setProgressState(stateInfo)
			}
		}
		setBarProgress(panelDownloadStatus.barTotalProgress, barProgress)
	})
}

// Lerp within the interval of the given stage using current/limit.
func calculateProgress(s Stage, current, total uint64) (barProgress int, percentage float64) {
	lowerEnd, upperEnd := s.getProgressInterval()
	if total == 0 || lowerEnd == upperEnd {
		return lowerEnd, float64(lowerEnd)
	}
	if current > total {
		current = total
	}
	percentage = float64(lowerEnd) + (float64(upperEnd)-float64(lowerEnd))*float64(current)/float64(total)
	if percentage > float64(upperEnd) {
		percentage = float64(upperEnd)
	} else if percentage <= float64(lowerEnd) {
		percentage = float64(lowerEnd)
	}
	barProgress = int(math.Round(percentage))
	if barProgress > upperEnd {
		barProgress = upperEnd
	} else if barProgress < lowerEnd {
		barProgress = lowerEnd
	}
	return barProgress, percentage
}

func updateProgressPeriodically(ctx context.Context) {
	const barUpdateInterval = time.Millisecond * 100
	const labelUpdateInterval = time.Second
	const titleUpdateInterval = time.Millisecond * 500
	barTimer := time.NewTimer(barUpdateInterval)
	labelTimer := time.NewTimer(labelUpdateInterval)
	titleTimer := time.NewTimer(titleUpdateInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-barTimer.C:
			{
				updateProgressBar()
				barTimer.Reset(barUpdateInterval)
			}
		case <-labelTimer.C:
			{
				updateProgressLabel()
				labelTimer.Reset(labelUpdateInterval)
			}
		case <-titleTimer.C:
			{
				updateWindowTitle()
				titleTimer.Reset(titleUpdateInterval)
			}
		}
	}
}

func NotifyProblem(problemMessage string, requiresUserAction bool) {
	uiShutdownMutex.Lock()
	defer uiShutdownMutex.Unlock()
	if !didQuit {
		ui.QueueMain(func() {
			if problemMessage == "" {
				panelDownloadStatus.currentProblemMessage = ""
			} else if requiresUserAction {
				panelDownloadStatus.currentProblemMessage = "Cannot continue: " + problemMessage
			} else {
				panelDownloadStatus.currentProblemMessage = "Taking longer than usual: " + problemMessage
			}
		})
	}
}

func ClearProblem() {
	NotifyProblem("", false)
}

func updateProgressBar() {
	uiShutdownMutex.Lock()
	defer uiShutdownMutex.Unlock()
	if !didQuit {
		ui.QueueMain(func() {
			barProgress, _ := calculateProgress(panelDownloadStatus.stage, ProgressFunc(panelDownloadStatus.stage), panelDownloadStatus.progressTarget)
			setBarProgress(panelDownloadStatus.barTotalProgress, barProgress)
		})
	}
}

func updateProgressLabel() {
	uiShutdownMutex.Lock()
	defer uiShutdownMutex.Unlock()
	if !didQuit {
		ui.QueueMain(func() {
			panelDownloadStatus.progressPrevious = panelDownloadStatus.progressCurrent
			panelDownloadStatus.progressCurrent = ProgressFunc(panelDownloadStatus.stage)

			delta := panelDownloadStatus.progressCurrent - panelDownloadStatus.progressPrevious
			var message string
			if panelDownloadStatus.stage.IsDownloadStage() {
				message = fmt.Sprintf("Downloading at %s. ", rateString(delta))
			}
			if panelDownloadStatus.currentProblemMessage != "" {
				message += fmt.Sprintf("(%s)", panelDownloadStatus.currentProblemMessage)
			}
			panelDownloadStatus.labelStatus.SetText(message)
		})
	}
}

func updateWindowTitle() {
	uiShutdownMutex.Lock()
	defer uiShutdownMutex.Unlock()
	if !didQuit {
		ui.QueueMain(func() {
			_, percentage := calculateProgress(panelDownloadStatus.stage, ProgressFunc(panelDownloadStatus.stage), panelDownloadStatus.progressTarget)

			// This should not be called too frequently; we observed Kubuntu's UI hanging for long durations (>5 seconds) already at 10 calls per second.
			window.SetTitle(fmt.Sprintf("[%.1f%%] %s", percentage, windowTitle))
		})
	}
}

func rateString(rate uint64) string {
	if rate < 1000 {
		return fmt.Sprintf("%d B/s", rate)
	} else if rate < 1024*10 {
		return fmt.Sprintf("%.2f KiB/s", float64(rate)/1024)
	} else if rate < 1024*100 {
		return fmt.Sprintf("%.1f KiB/s", float64(rate)/1024)
	} else if rate < 1024*1000 {
		return fmt.Sprintf("%d KiB/s", rate/1024)
	} else if rate < 1024*1024*10 {
		return fmt.Sprintf("%.2f MiB/s", float64(rate)/(1024*1024))
	}
	return fmt.Sprintf("%.1f MiB/s", float64(rate)/(1024*1024))
}
