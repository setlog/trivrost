package gui

import (
	"context"
	"sync"

	"github.com/setlog/trivrost/pkg/misc"

	"github.com/andlabs/ui"
	log "github.com/sirupsen/logrus"
)

type progressState int

const (
	stateInfo   progressState = 1
	stateError  progressState = 2
	statePaused progressState = 3
)

const maxLineWidth = 110

var (
	window                                        *ui.Window
	windowCalculatedWidth, windowCalculatedHeight int
	windowTitle                                   string
	waitDialog                                    *ui.Window
	waitDialogText                                *ui.Label
	panelDownloadStatus                           *DownloadStatusPanel

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
		labelBox, _ := textBox(ui.NewVerticalBox(), message, maxLineWidth)
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

// Pause shows given message in the download status panel along with a clickable link
// which reads "Continue" and blocks until the user clicks it.
func Pause(ctx context.Context, message string) {
	var n int
	var hBox *ui.Box
	c := make(chan struct{}, 1)
	ui.QueueMain(func() {
		_, n = textBox(panelDownloadStatus.pauseStatusBox, message, maxLineWidth)
		hBox = ui.NewHorizontalBox()
		hBox.Append(newLinkLabel("Continue", ui.DrawTextAlignLeft, misc.WriteAttempter(c)), true)
		hBox.Append(ui.NewLabel(""), false) // Needed or else the box has no minimum dimensions.
		hBox.Append(newLogsLinkLabel(), true)
		panelDownloadStatus.pauseStatusBox.Append(hBox, false)
		setProgressState(statePaused)
		panelDownloadStatus.inlineStatusBox.Hide()
		panelDownloadStatus.pauseStatusBox.Show()
	})
	misc.WaitCancelable(ctx, c)
	ui.QueueMain(func() {
		panelDownloadStatus.pauseStatusBox.Hide()
		setWindowDimensions(window.Handle(), windowCalculatedWidth, windowCalculatedHeight)
		panelDownloadStatus.inlineStatusBox.Show()
		setProgressState(stateInfo)
		clearBox(panelDownloadStatus.pauseStatusBox, n+1)
		hBox.Destroy()
	})
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

		panelDownloadStatus = newDownloadStatusPanel()
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
			windowCalculatedWidth, windowCalculatedHeight = getWindowDimensions(window.Handle())
			centerWindow(window.Handle())
		}

		go updateProgressPeriodically(ctx)

		guiInitWaitGroup.Done()
	})
}
