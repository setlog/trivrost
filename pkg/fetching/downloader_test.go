package fetching

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"

	"github.com/setlog/trivrost/pkg/launcher/config"
)

type ErrorRecordingHandler struct {
	ErrChan chan error
}

func (handler *ErrorRecordingHandler) HandleProgress(fromURL string, workerId int, receivedBytes uint64) {
}

func (handler *ErrorRecordingHandler) HandleStartDownload(fromURL string, workerId int) {
}

func (handler *ErrorRecordingHandler) HandleFinishDownload(fromURL string, workerId int) {
}

func (handler *ErrorRecordingHandler) HandleFailDownload(fromURL string, workerId int, err error) {
}

func (handler *ErrorRecordingHandler) HandleHttpGetError(fromURL string, err error) {
	handler.ErrChan <- err
}

func (handler *ErrorRecordingHandler) HandleBadHttpResponse(fromURL string, code int) {
	handler.ErrChan <- fmt.Errorf("HTTP %d: %s", code, http.StatusText(code))
}

func (handler *ErrorRecordingHandler) HandleReadError(fromURL string, err error, firstByteIndex int64) {
	handler.ErrChan <- err
}

type DummyEnvironment struct {
	DoForClientFunc   func(client *http.Client, req *http.Request) (*http.Response, error)
	Data              []byte
	OmitContentLength bool
}

type RiggedReader struct {
	*bytes.Reader
	readCount int
	failEvery int
	ctx       context.Context
}

func (rr *RiggedReader) Read(p []byte) (n int, err error) {
	if rr.ctx.Err() != nil {
		return 0, rr.ctx.Err()
	}
	n, err = rr.Reader.Read(p)
	rr.readCount += n
	if err == nil {
		if rr.failEvery >= 0 && rr.readCount >= rr.failEvery {
			rr.readCount = 0
			return n, fmt.Errorf("Rigged error")
		}
	}
	return n, err
}

func (db *RiggedReader) Close() error {
	return nil
}

func (de *DummyEnvironment) TestDownload(t *testing.T, fromUrl string) {
	dl := NewDownload(context.Background(), fromUrl)
	data, err := ioutil.ReadAll(dl)
	if err != nil {
		t.Fatalf("Download \"%s\" failed unexpectedly: %v", fromUrl, err)
	}
	if !bytes.Equal(data, de.Data) {
		t.Fatalf("Download \"%s\" yielded unexpected data", fromUrl)
	}
}

func (de *DummyEnvironment) TestDownloadCancel(t *testing.T, fromUrl string, cancelAfter int) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	dl := NewDownload(ctx, fromUrl)
	_, err := ReadExactly(dl, cancelAfter)
	if err != nil {
		t.Fatalf("Download \"%s\" could not deliver %d bytes: %v", fromUrl, cancelAfter, err)
	}
	cancelFunc()
	data, err := ioutil.ReadAll(dl)
	if err != context.Canceled {
		t.Fatalf("Download \"%s\" did not fail with context.Canceled after %d bytes. Remainder: %d. Error: %v", fromUrl, cancelAfter, len(data), err)
	}
}

func (de *DummyEnvironment) TestDownloadFailure(t *testing.T, fromUrl string) {
	dl := NewDownload(context.Background(), fromUrl)
	_, err := ioutil.ReadAll(dl)
	if err == nil {
		t.Fatalf("Download \"%s\" succeeded unexpectedly", fromUrl)
	}
}

func (de *DummyEnvironment) TestDownloadRetries(t *testing.T, fromUrl string) {
	h := &ErrorRecordingHandler{ErrChan: make(chan error, 1)}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	dl := NewHandledDownload(ctx, fromUrl, h)
	p := make([]byte, 100)
	go dl.Read(p)
	err := <-h.ErrChan
	if err == nil {
		t.Fatalf("Got nil error when expected non-nil error.")
	}
}

func CreateDummyEnvironment(t *testing.T, dataLength, failEvery int) *DummyEnvironment {
	de := &DummyEnvironment{}
	de.Data = arbitraryData(dataLength)
	de.DoForClientFunc = func(client *http.Client, req *http.Request) (response *http.Response, err error) {
		response = &http.Response{StatusCode: 200, Header: make(http.Header)}
		var rangeStart, rangeEnd int64 = 0, int64(dataLength - 1)
		requestedRange := req.Header.Get("Range")
		if requestedRange != "" {
			rangeStart, rangeEnd, err = ParseRange(requestedRange, rangeEnd)
			if err != nil {
				response.StatusCode = http.StatusBadRequest
				return response, nil
			}
			response.StatusCode = http.StatusPartialContent
			response.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, dataLength))
		}
		if !de.OmitContentLength {
			response.Header.Set("Content-Length", fmt.Sprintf("%d", rangeEnd-rangeStart+1))
		}
		response.Header.Set("ETag", hex.EncodeToString(de.Data))
		response.Body = &RiggedReader{Reader: bytes.NewReader(de.Data[rangeStart : rangeEnd+1]), failEvery: failEvery, ctx: req.Context()}
		return response, nil
	}
	return de
}

func arbitraryData(length int) []byte {
	data := make([]byte, length)
	rng := rand.New(rand.NewSource(int64(length)))
	for i := 0; i < length; i++ {
		data[i] = byte(rng.Uint32())
	}
	return data
}

func TestUpdateFile(t *testing.T) {
	d, err := ioutil.TempDir(".", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(d)
	const dataSize = 2000
	de := CreateDummyEnvironment(t, dataSize, -1)
	DoForClientFunc = de.DoForClientFunc
	x := sha256.Sum256(de.Data)
	expectedSha := hex.EncodeToString(x[:])
	dl := NewDownload(context.Background(), "http://example.com")
	di := &config.FileInfo{SHA256: expectedSha, Size: dataSize}
	err = updateFile(dl, di, "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testfile.txt")
	diskData, err := ioutil.ReadFile("testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(diskData, de.Data) {
		t.Fatalf("Data on disk mismatches data of download")
	}
}
