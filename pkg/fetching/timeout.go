package fetching

import (
	"fmt"
	"io"
	"time"
)

// TimeoutingReader wraps an io.Reader with a configurable timeout.
type TimeoutingReader struct {
	Reader  io.Reader
	Timeout time.Duration
}

type readResult struct {
	n   int
	err error
}

func (trc *TimeoutingReader) Read(p []byte) (n int, err error) {
	t := time.NewTimer(trc.Timeout)
	c := make(chan readResult, 1)
	go trc.read(p, c)
	select {
	case res := <-c:
		{
			t.Stop()
			return res.n, res.err
		}
	case <-t.C:
		{
			return 0, fmt.Errorf("Read() timed out: failed to write max %d bytes in %v", len(p), trc.Timeout)
		}
	}
}

func (trc *TimeoutingReader) read(p []byte, c chan readResult) {
	n, err := trc.Reader.Read(p)
	c <- readResult{n, err}
}
