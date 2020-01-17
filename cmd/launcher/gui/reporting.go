package gui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/setlog/trivrost/pkg/misc"

	"github.com/setlog/trivrost/cmd/launcher/places"
	"github.com/setlog/trivrost/pkg/logging"

	"github.com/setlog/trivrost/pkg/system"

	log "github.com/sirupsen/logrus"

	"github.com/setlog/trivrost/cmd/launcher/flags"
)

func HandlePanic(launcherFlags *flags.LauncherFlags) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok && errors.Is(err, context.Canceled) {
			log.Infof("Quitting: %v", err)
		} else {
			PanicInformatively(r, launcherFlags)
		}
	}
}

func PanicInformatively(r interface{}, launcherFlags *flags.LauncherFlags) {
	defer presentError(getPanicMessage(r), launcherFlags.DismissGuiPrompts)
	misc.LogRecoveredValue(r)
}

func getPanicMessage(r interface{}) string {
	message := "Something went wrong. The program will now close."

	userError, ok := r.(misc.IUserError)
	if ok && !misc.IsNil(userError) {
		message = userError.UserError()
	}

	fileSystemError, ok := r.(*system.FileSystemError)
	if ok && fileSystemError != nil {
		if os.IsPermission(fileSystemError.Unwrap()) {
			message = "Error: Insufficient permissions to write files in your own user directory. " +
				"Please contact your system administrator and verify that you have full access to your user directory."
		} else {
			message = fmt.Sprintf("Error: Your machine's file system denied a required operation. The error received was: %v", fileSystemError.Unwrap())
		}
	}

	if !strings.HasSuffix(message, ".") && !strings.HasSuffix(message, "!") && !strings.HasSuffix(message, "?") {
		message += "."
	}

	return message
}

func presentError(message string, dismissGuiPrompts bool) {
	if BlockingDialog("Error", fmt.Sprintf("%s\n\nYou can find technical information in the log files under\n%s\n",
		message, places.GetAppLogFolderPath()), []string{"Open log folder and close", "Close"}, 1, dismissGuiPrompts) == 0 {
		showLogFolder()
	}
}

func showLogFolder() {
	log.Infof("Showing file \"%s\" in file manager.", logging.GetLogFilePath())
	err := system.ShowLocalFileInFileManager(logging.GetLogFilePath())
	if err != nil {
		log.Errorf("Error showing file \"%s\" in file manager: %v", logging.GetLogFilePath(), err)
	}
}
