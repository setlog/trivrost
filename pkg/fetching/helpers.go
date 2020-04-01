package fetching

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-ieproxy"
	log "github.com/sirupsen/logrus"
)

func MakeClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: defaultTimeout,
		}).Dial,
		TLSHandshakeTimeout:   defaultTimeout,
		IdleConnTimeout:       defaultTimeout,
		ResponseHeaderTimeout: defaultTimeout,
		ExpectContinueTimeout: defaultTimeout,
		Proxy:                 GetProxyLoggingFunc(),
	}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	return &http.Client{
		Timeout:   time.Hour * 12, // <- This is not a timeout. This is a deadline.
		Transport: transport,
	}
}

func GetProxyLoggingFunc() func(req *http.Request) (*url.URL, error) {
	proxyFunc := ieproxy.GetProxyFunc()
	return func(req *http.Request) (*url.URL, error) {
		proxyURL, err := proxyFunc(req)
		if err != nil {
			log.Warnf("Getting proxy for URL %v failed: %v", req.URL, err)
		} else {
			if proxyURL == nil {
				log.Infof("GET %v (direct).", req.URL)
			} else {
				log.Infof("GET %v (with proxy: %v).", req.URL, proxyURL)
			}
		}
		return proxyURL, err
	}
}

func newRequestWithCancel(ctx context.Context, fromUrl string) (*http.Request, context.CancelFunc) {
	req, err := http.NewRequest("GET", fromUrl, nil)
	if err != nil {
		panic(DownloadError(err.Error()))
	}
	ctx, cancelFunc := context.WithCancel(ctx)

	return req.WithContext(ctx), cancelFunc
}

func newRangeRequestWithCancel(ctx context.Context, fromUrl string, firstByte int64, lastByte int64) (*http.Request, context.CancelFunc) {
	req, cancelFunc := newRequestWithCancel(ctx, fromUrl)
	if lastByte >= 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", firstByte, lastByte))
	} else {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", firstByte))
	}
	return req, cancelFunc
}

func ParseRange(rangeHeader string, endMax int64) (rangeStart, rangeEnd int64, err error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, fmt.Errorf("Range header does not start with 'bytes='")
	}
	parts := strings.Split(rangeHeader[6:], ",")
	if len(parts) != 1 {
		return 0, 0, fmt.Errorf("Range header did not contain exactly 1 range")
	}
	subParts := strings.Split(parts[0], "-")
	if len(subParts) == 0 || len(subParts) > 2 {
		return 0, 0, fmt.Errorf("Range header contains malformed range")
	}
	rangeStart, err = strconv.ParseInt(strings.Trim(subParts[0], " "), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	if len(subParts) == 2 {
		rangeEndStr := strings.Trim(subParts[1], " ")
		if rangeEndStr == "" {
			rangeEnd = endMax
		} else {
			rangeEnd, err = strconv.ParseInt(rangeEndStr, 10, 64)
		}
	} else {
		err = fmt.Errorf("Range header is missing '-' to indicate open-ended range")
	}
	return rangeStart, rangeEnd, err
}

func ReadExactly(r io.Reader, byteCount int) (data []byte, err error) {
	p := make([]byte, 1024)
	totalRead := 0
	for len(data) < byteCount {
		if len(p) > byteCount-totalRead {
			p = make([]byte, byteCount-totalRead)
		}
		n, err := r.Read(p)
		totalRead += n
		data = append(data, p[0:n]...)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

// BitRateToByteDuration returns the duration between bytes for a given bit rate in bits per second,
// rounded down to nanoseconds.
func BitRateToByteDuration(bitRate int64) time.Duration {
	return time.Nanosecond * time.Duration(8*1000000000/bitRate)
}

func createWorkerIdChannel(maxWorkers int) chan int {
	workerIds := make(chan int, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workerIds <- i
	}
	return workerIds
}

func parseTotalLengthFromContentRangeHeader(contentRange string) int64 {
	slashIndex := strings.LastIndex(contentRange, "/")
	if slashIndex == -1 {
		return -1
	}
	i, err := strconv.ParseInt(contentRange[slashIndex+1:], 10, 64)
	if err != nil {
		return -1
	}
	return i
}

func stringStringMapKeys(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	return s
}
