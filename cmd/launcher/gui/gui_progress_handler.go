package gui

import (
	"fmt"
	"sync"

	"net/http"

	log "github.com/sirupsen/logrus"
)

// Implements fetching.DownloadProgressHandler
type GuiDownloadProgressHandler struct {
	progressMutex          *sync.RWMutex
	progressAccumulator    uint64
	ongoingProgressBuckets []uint64
	problemUrl             string
}

func NewGuiDownloadProgressHandler(bucketCount int) *GuiDownloadProgressHandler {
	return &GuiDownloadProgressHandler{progressMutex: &sync.RWMutex{}, ongoingProgressBuckets: make([]uint64, bucketCount)}
}

func (handler *GuiDownloadProgressHandler) ResetProgress() {
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.progressAccumulator = 0
	for i := range handler.ongoingProgressBuckets {
		handler.ongoingProgressBuckets[i] = 0
	}
}

func (handler *GuiDownloadProgressHandler) GetProgress() uint64 {
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	currentTotal := handler.progressAccumulator
	for _, v := range handler.ongoingProgressBuckets {
		currentTotal += v
	}
	return currentTotal
}

func (handler *GuiDownloadProgressHandler) HandleProgress(fromURL string, workerId int, receivedBytes uint64) {
	handler.progressMutex.RLock()
	defer handler.progressMutex.RUnlock()
	handler.ongoingProgressBuckets[workerId] = receivedBytes
	if fromURL == handler.problemUrl {
		ClearProblem()
	}
}

func (handler *GuiDownloadProgressHandler) HandleStartDownload(fromURL string, workerId int) {
}

func (handler *GuiDownloadProgressHandler) HandleFinishDownload(fromURL string, workerId int) {
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.progressAccumulator += handler.ongoingProgressBuckets[workerId]
	handler.ongoingProgressBuckets[workerId] = 0
}

func (handler *GuiDownloadProgressHandler) HandleFailDownload(fromURL string, workerId int, err error) {
	log.Errorf("GET %s failed: %v", fromURL, err)
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.problemUrl = fromURL
	NotifyProblem("Security error", false)
}

func (handler *GuiDownloadProgressHandler) HandleHttpGetError(fromURL string, err error) {
	log.Warnf("GET %s could not start: %v", fromURL, err)
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.problemUrl = fromURL
	NotifyProblem("Cannot reach server", false)
}

func (handler *GuiDownloadProgressHandler) HandleBadHttpResponse(fromURL string, code int) {
	log.Warnf("GET %s yielded bad HTTP response: %s (Code was %d)", fromURL, http.StatusText(code), code)
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.problemUrl = fromURL
	NotifyProblem(fmt.Sprintf("HTTP Status %d", code), false)
}

func (handler *GuiDownloadProgressHandler) HandleReadError(fromURL string, err error, firstByteIndex int64) {
	log.Warnf("GET %s interrupted: %v. Continuing from byte %d", fromURL, err, firstByteIndex)
	handler.progressMutex.Lock()
	defer handler.progressMutex.Unlock()
	handler.problemUrl = fromURL
	NotifyProblem("Connection unstable", false)
}
