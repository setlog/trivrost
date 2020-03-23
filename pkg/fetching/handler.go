package fetching

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// DownloadProgressHandler is an interface which defines callbacks for typical events which
// will or may occur during a resource download via an HTTP GET request, such as receiving bytes,
// the connection being interrupted or a bad HTTP response code being received.
type DownloadProgressHandler interface {
	HandleStartDownload(fromURL string, workerId int)                  // The download of resource at given URL will be started.
	HandleProgress(fromURL string, workerId int, receivedBytes uint64) // The download of resource at given URL made some progress.
	HandleFinishDownload(fromURL string, workerId int)                 // The download of resource at given URL finished successfully.
	HandleFailDownload(fromURL string, workerId int, err error)        // The download of resource at given URL failed irrecoverably.

	HandleHttpGetError(fromURL string, err error)   // The HTTP GET request to download the resource did not receive an HTTP response.
	HandleBadHttpResponse(fromURL string, code int) // A bad HTTP response code was received.
	HandleReadError(fromURL string, err error, firstByteIndex int64)      // An error occurred while reading the response body (data) of the resource at the given URL.
}

type ConsoleDownloadProgressHandler struct {
}

func (handler *ConsoleDownloadProgressHandler) HandleProgress(fromURL string, workerId int, receivedBytes uint64) {
}

func (handler *ConsoleDownloadProgressHandler) HandleStartDownload(fromURL string, workerId int) {
	log.Infof("Downloading %s", fromURL)
}

func (handler *ConsoleDownloadProgressHandler) HandleFinishDownload(fromURL string, workerId int) {
}

func (handler *ConsoleDownloadProgressHandler) HandleFailDownload(fromURL string, workerId int, err error) {
}

func (handler *ConsoleDownloadProgressHandler) HandleHttpGetError(fromURL string, err error) {
	log.Errorf("Error downloading %s: %v.", fromURL, err)
	os.Exit(1)
}

func (handler *ConsoleDownloadProgressHandler) HandleBadHttpResponse(fromURL string, code int) {
	log.Errorf("GET %s yielded bad HTTP response: %v (Code was %d)", fromURL, http.StatusText(code), code)
	os.Exit(1)
}

func (handler *ConsoleDownloadProgressHandler) HandleReadError(fromURL string, err error, firstByteIndex int64) {
	log.Errorf("Could not copy bytes from \"%s\": %v", fromURL, err)
	os.Exit(1)
}

type EmptyHandler struct {
}

func (handler *EmptyHandler) HandleProgress(fromURL string, workerId int, receivedBytes uint64) {
}

func (handler *EmptyHandler) HandleStartDownload(fromURL string, workerId int) {
}

func (handler *EmptyHandler) HandleFinishDownload(fromURL string, workerId int) {
}

func (handler *EmptyHandler) HandleFailDownload(fromURL string, workerId int, err error) {
}

func (handler *EmptyHandler) HandleHttpGetError(fromURL string, err error) {
}

func (handler *EmptyHandler) HandleBadHttpResponse(fromURL string, code int) {
}

func (handler *EmptyHandler) HandleReadError(fromURL string, err error, firstByteIndex int64) {
}
