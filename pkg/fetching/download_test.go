package fetching_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/setlog/trivrost/pkg/dummy"
	"github.com/setlog/trivrost/pkg/fetching"
)

func TestDownloadCompletes(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 5000, -1)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownload(t, "http://example.com")
}

func TestDownloadWithRangeRequestsCompletes(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 1000, 0)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownload(t, "http://example.com/a")
	de = fetching.CreateDummyEnvironment(t, 10000, 600)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownload(t, "http://example.com/b")
}

func TestDownloadCancels(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 1000, -1)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownloadCancel(t, "http://example.com/a", 1000)
	de = fetching.CreateDummyEnvironment(t, 10000, -1)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownloadCancel(t, "http://example.com/b", 1000)
}

func TestRequestCreationFailure(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 1000, -1)
	fetching.DoForClientFunc = de.DoForClientFunc
	de.TestDownloadFailure(t, "://badurl")
}

func TestHttpClientDoFuncFailure(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 1000, -1)
	fetching.DoForClientFunc = func(client *http.Client, req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("some error")
	}
	de.TestDownloadRetries(t, "http://example.com")
}

func TestBadHttpResponse(t *testing.T) {
	de := fetching.CreateDummyEnvironment(t, 1000, -1)
	fetching.DoForClientFunc = func(client *http.Client, req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: &dummy.ReadCloser{}}, nil
	}
	de.TestDownloadRetries(t, "http://example.com")
}
