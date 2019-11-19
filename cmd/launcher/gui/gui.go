package gui

import (
	"context"
	"strings"
	"sync"

	"github.com/andlabs/ui"
	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/pkg/misc"
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

var (
	window              *ui.Window
	windowTitle         string
	waitDialog          *ui.Window
	waitDialogText      *ui.Label
	panelDownloadStatus *DownloadStatusPanel

	guiInitWaitGroup = &sync.WaitGroup{}
	didQuit          bool

	uiShutdownMutex *sync.Mutex
)

func init() {
	guiInitWaitGroup.Add(1)
	uiShutdownMutex = &sync.Mutex{}
}

func Quit() {
	uiShutdownMutex.Lock()
	defer uiShutdownMutex.Unlock()
	if didQuit {
		panic("Called gui.Quit() more than once.")
	}
	didQuit = true
	ui.QueueMain(func() {
		window.Destroy()
		ui.Quit()
	})
}

func WaitUntilReady() {
	guiInitWaitGroup.Wait()
}

func BlockingDialog(title, message string, options []string, defaultOption int, dismissGuiPrompts bool) int {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	chosenOption := defaultOption
	message = misc.WordWrap(message, 120) // The ui library itself wraps nothing, not even spaces, resulting in very wide windows (>10000 pixels) without this.
	var waitGroupDoneTrigger sync.Once
	ui.QueueMain(func() {
		dialogWindow := ui.NewWindow(title, 600, 90, false)
		applyIconToWindow(dialogWindow.Handle())
		applyWindowStyle(dialogWindow.Handle())
		dialogWindow.OnClosing(func(*ui.Window) bool {
			waitGroupDoneTrigger.Do(waitGroup.Done)
			return true
		})

		mainBox := ui.NewVerticalBox()
		mainBox.SetPadded(true)

		labelBox := ui.NewVerticalBox()

		// HACK: When the GUI lib calculates string width, it ignores newlines,
		// resulting in too large results, and thus too wide windows.
		lines := strings.Split(message, "\n")
		for _, line := range lines {
			lineLabel := ui.NewLabel(line)
			labelBox.Append(lineLabel, false)
		}

		mainBox.Append(labelBox, true)

		if len(options) > 0 {
			buttonBox := ui.NewHorizontalBox()
			buttonBox.SetPadded(true)
			for i, option := range options {
				optionButton := ui.NewButton(option)
				val := i
				optionButton.OnClicked(func(*ui.Button) {
					waitGroupDoneTrigger.Do(func() {
						chosenOption = val
						waitGroup.Done()
						dialogWindow.Destroy()
					})
				})
				buttonBox.Append(optionButton, true)
			}
			mainBox.Append(buttonBox, false)
		}

		dialogWindow.SetChild(mainBox)
		dialogWindow.SetMargined(true)
		centerWindow(dialogWindow.Handle())
		dialogWindow.Show()
		centerWindow(dialogWindow.Handle())

		if dismissGuiPrompts {
			log.Infof("Automatically dismissing dialog \"%s\" with default option %d.", title, defaultOption)
			chosenOption = defaultOption
			waitGroup.Done()
			dialogWindow.Destroy()
		}
	})
	waitGroup.Wait()
	return chosenOption
}

func ShowWaitDialog(title, text string) {
	ui.QueueMain(func() {
		if waitDialog == nil {
			waitDialog = ui.NewWindow(title, 300, 90, false)
			applyIconToWindow(waitDialog.Handle())
			applyWindowStyle(waitDialog.Handle())
			waitDialog.OnClosing(func(*ui.Window) bool {
				return false
			})
			mainBox := ui.NewVerticalBox()
			mainBox.SetPadded(true)
			waitDialogText = ui.NewLabel(text)
			mainBox.Append(waitDialogText, false)

			waitDialog.SetMargined(true)
			waitDialog.SetChild(mainBox)
		} else {
			waitDialog.SetTitle(title)
			waitDialogText.SetText(text)
		}

		centerWindow(waitDialog.Handle())
		waitDialog.Show()
		centerWindow(waitDialog.Handle())
	})
}

func HideWaitDialog() {
	ui.QueueMain(func() {
		if waitDialog != nil {
			waitDialog.Hide()
		}
	})
}

func Pause(message string) {
	// TODO: Pause
}

// Main hands control over to ui.Main() to initialize and manage the GUI. It blocks until gui.Quit() is called.
func Main(ctx context.Context, cancelFunc func(), title string, showMainWindow bool) error {
	log.WithFields(log.Fields{"title": title, "showMainWindow": showMainWindow}).Info("Initializing GUI.")
	// Note: ui.Main() calls any functions queued with ui.QueueMain() before the one we provide via parameter.
	return ui.Main(func() {
		windowTitle = title
		window = ui.NewWindow(windowTitle, 600, 50, false)
		applyIconToWindow(window.Handle())
		applyWindowStyle(window.Handle())

		window.OnClosing(func(*ui.Window) bool {
			log.Info("User tries to close the window.")
			cancelFunc()
			return false
		})

		panelDownloadStatus = makeContent()
		window.SetChild(panelDownloadStatus)
		window.SetMargined(true)

		ui.OnShouldQuit(func() bool {
			log.Info("OnShouldQuit().")
			cancelFunc()
			return false
		})

		if showMainWindow {
			centerWindow(window.Handle())
			window.Show()
			centerWindow(window.Handle())
		}

		go updateProgressPeriodically(ctx)

		guiInitWaitGroup.Done()
	})
}

func makeContent() *DownloadStatusPanel {
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
