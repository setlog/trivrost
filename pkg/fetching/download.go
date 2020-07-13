package fetching

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Changed only during tests.
var DoForClientFunc = DoForClient

func DoForClient(client *http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req)
}

type DownloadError string

func (err DownloadError) Error() string {
	return string(err)
}

type writeCounter struct {
	counted  uint64
	url      string
	workerId int
	handler  DownloadProgressHandler
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	delta := len(p)
	wc.counted += uint64(delta)
	if delta > 0 {
		wc.handler.HandleProgress(wc.url, wc.workerId, wc.counted)
	}
	return delta, nil
}

// Download wraps the retrieval of a resource at a URL using HTTP GET requests and exposes
// the received data through its implementation of the io.Reader interface.
// The Read() method will only return a non-nil error other than io.EOF when the
// Content-Length of the requested resource changes during Download's attempts
// to retrieve it. See Download.handler for Download's behavior in other error scenarios.
type Download struct {
	url string // The URL of the resource to download.

	ctx context.Context

	// Calls to methods of handler occur during Download.Read(). The calls are followed by these behaviors:
	// HandleStartDownload(): Read() will start its first (ideally only) HTTP request.
	// HandleHttpGetError() and HandleBadHttpResponse(): Read() will not return. The request will be retried after one second.
	// HandleReadError(): If at least 1 byte has been read, Read() will return with a nil error, having accepted all data which
	//                    has been received so far. Either way, Download will continue with a range-request immediately.
	// HandleFailDownload(): Read() will return with a non-nil, non-io.EOF error.
	// HandleFinishDownload(): Read() will return with io.EOF.
	handler DownloadProgressHandler

	// An arbitrary integer id or value which should be contained in applicable calls on DownloadProgressHandler. Useful for concurrency.
	workerId int

	isDownloadStarted     bool
	gotValidFirstResponse bool
	firstByteIndex        int64
	lastByteIndex         int64

	client         *http.Client
	request        *http.Request
	cancelRequest  context.CancelFunc
	cooldownTime   time.Time
	cooldownStacks int

	response       *http.Response
	responseReader io.Reader
}

func NewDownload(ctx context.Context, resourceUrl string) *Download {
	return &Download{url: resourceUrl, client: MakeClient(), ctx: ctx, handler: &EmptyHandler{}}
}

func NewHandledDownload(ctx context.Context, resourceUrl string, handler DownloadProgressHandler) *Download {
	return &Download{url: resourceUrl, client: MakeClient(), ctx: ctx, handler: handler}
}

func NewDownloadForConcurrentUse(ctx context.Context, resourceUrl string, client *http.Client, handler DownloadProgressHandler, workerId int) *Download {
	return &Download{
		url:      resourceUrl,
		ctx:      ctx,
		client:   client,
		handler:  handler,
		workerId: workerId,
	}
}

func (dl *Download) URL() string {
	return dl.url
}

// Read reads some data of the requested resource into p.
// Calling Read() again after it returned a non-nil error results in undefined behaviour.
func (dl *Download) Read(p []byte) (n int, err error) {
	defer dl.handlePanic(&err)
	for n == 0 && err == nil {
		if len(p) == 0 || dl.ctx.Err() != nil {
			return 0, dl.ctx.Err()
		}
		dl.waitCooldown()
		if !dl.isDownloadStarted {
			dl.handler.HandleStartDownload(dl.url, dl.workerId)
			dl.isDownloadStarted = true
		}
		n, err = dl.readDownload(p)
	}
	if err == io.EOF {
		dl.handler.HandleFinishDownload(dl.url, dl.workerId)
	} else if err != nil {
		dl.handler.HandleFailDownload(dl.url, dl.workerId, err)
	}
	return n, err
}

// Close closes the underlying HTTP response body of the Download if one currently exists.
// If you Read() until you receive a non-nil error, you do not have to call Close(); otherwise you must.
func (dl *Download) Close() error {
	if dl.response != nil {
		return dl.response.Body.Close()
	}
	return nil
}

func (dl *Download) readDownload(p []byte) (bytesReadCount int, err error) {
	if dl.response == nil {
		dl.request, dl.cancelRequest = dl.createRequest()
		dl.response = dl.sendRequest(dl.request)
		if dl.response == nil {
			return 0, nil
		}
		dl.processResponse()
	}
	if dl.response != nil {
		dl.resetCooldown()
		return dl.readFromResponse(p)
	}
	return 0, nil
}

func (dl *Download) createRequest() (*http.Request, context.CancelFunc) {
	if !dl.gotValidFirstResponse {
		return newRequestWithCancel(dl.ctx, dl.url)
	}
	return newRangeRequestWithCancel(dl.ctx, dl.url, dl.firstByteIndex, dl.lastByteIndex)
}

func (dl *Download) sendRequest(req *http.Request) *http.Response {
	resp, err := DoForClientFunc(dl.client, req)
	if err != nil {
		dl.cleanUp()
		dl.handler.HandleHttpGetError(dl.url, err)
		dl.inscribeCooldown()
	} else {
		counter := &writeCounter{counted: uint64(dl.firstByteIndex), url: dl.url, workerId: dl.workerId, handler: dl.handler}
		timeoutingBodyReader := &TimeoutingReader{Reader: resp.Body, Timeout: defaultTimeout * 30}
		dl.responseReader = io.TeeReader(timeoutingBodyReader, counter)
	}
	return resp
}

func (dl *Download) processResponse() {
	if !isRangeRequest(dl.request) && dl.response.StatusCode == http.StatusOK {
		dl.acceptFirstResponseHeader(dl.response.Header)
	} else if !(isRangeRequest(dl.request) && dl.response.StatusCode == http.StatusPartialContent) {
		dl.cleanUp()
		if dl.response.StatusCode == http.StatusRequestedRangeNotSatisfiable {
			panic(DownloadError("remote file changed during download"))
		}
		dl.handler.HandleBadHttpResponse(dl.url, dl.response.StatusCode)
		dl.response = nil
		dl.inscribeCooldown()
	}
}

func (dl *Download) acceptFirstResponseHeader(header http.Header) {
	contentLength, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Printf("Assuming remote file won't change: Could not get Content-Length header from \"%s\": %v.", dl.url, err)
		dl.lastByteIndex = -1
	} else {
		dl.lastByteIndex = contentLength - 1
	}
	dl.gotValidFirstResponse = true
}

func (dl *Download) readFromResponse(p []byte) (n int, err error) {
	n, err = dl.responseReader.Read(p)
	dl.firstByteIndex += int64(n)
	if (dl.lastByteIndex >= 0) && (dl.firstByteIndex > dl.lastByteIndex+1) {
		panic(fmt.Errorf("read more bytes than expected"))
	}
	if err != nil {
		dl.cleanUp() // Note: https://github.com/golang/go/issues/26095#issuecomment-400903313
		dl.response = nil
		if err != io.EOF { // Network failures are temporary. Keep trying until it works.
			dl.handler.HandleReadError(dl.url, err, dl.firstByteIndex)
			return n, nil
		}
		if (dl.lastByteIndex >= 0) && (dl.firstByteIndex < dl.lastByteIndex+1) {
			return n, fmt.Errorf("read less bytes than expected")
		}
	}
	return n, err
}

func (dl *Download) handlePanic(errPtr *error) {
	if r := recover(); r != nil {
		dl.cleanUp()
		if !setError(r, errPtr) {
			panic(r)
		}
	}
}

func setError(panicObject interface{}, errPtr *error) bool {
	if panicErr, ok := panicObject.(DownloadError); ok {
		*errPtr = panicErr
		return true
	} else if panicObject == context.Canceled {
		*errPtr = context.Canceled
		return true
	}
	return false
}

func (dl *Download) cleanUp() {
	if dl.response != nil && dl.response.Body != nil {
		dl.response.Body.Close()
	}
	if dl.cancelRequest != nil {
		dl.cancelRequest()
	}
}

func (dl *Download) waitCooldown() {
	now := time.Now()
	if dl.cooldownStacks > 0 && dl.cooldownTime.After(now) {
		select {
		case <-time.NewTimer(dl.cooldownTime.Sub(now)).C:
		case <-dl.ctx.Done():
			panic(dl.ctx.Err())
		}
	}
}

func (dl *Download) inscribeCooldown() {
	cooldownIntervalOptions := []time.Duration{1, 1, 2, 3, 5, 8, 13}
	cooldownOptionIndex := intMin(dl.cooldownStacks, len(cooldownIntervalOptions)-1)
	cooldownDuration := time.Second * cooldownIntervalOptions[cooldownOptionIndex]

	p := make([]byte, 1)
	_, err := rand.Read(p)
	if err != nil {
		log.Printf("Could not crypto/rand.Read(): %v\n", err)
	} else {
		rand := ((time.Duration(p[0]) - 127) * cooldownDuration) / (10 * 128) // +/- 10%
		cooldownDuration += rand
	}

	dl.cooldownTime = time.Now().Add(cooldownDuration)
	dl.cooldownStacks++
}

func (dl *Download) resetCooldown() {
	dl.cooldownStacks = 0
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isRangeRequest(req *http.Request) bool {
	return strings.HasPrefix(req.Header.Get("Range"), "bytes=")
}
