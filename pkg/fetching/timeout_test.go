package fetching_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/setlog/trivrost/pkg/fetching"
)

type IdleReader struct {
}

func (ir *IdleReader) Read(p []byte) (n int, err error) {
	time.Sleep(time.Minute)
	return 0, fmt.Errorf("Stopped idling")
}

func TestTimeoutingReaderTimeout(t *testing.T) {
	trc := &fetching.TimeoutingReader{Reader: &IdleReader{}, Timeout: time.Nanosecond}
	_, err := trc.Read(nil)
	if err == nil {
		t.Fatalf("Expected read to fail with timeout, but got no error.")
	}
	if !strings.Contains(err.Error(), "Read() timed out: failed to write") {
		t.Fatalf("Expected read to fail with timeout, but got the following error instead: %v", err)
	}
}

func TestTimeoutingReaderRead(t *testing.T) {
	trc := &fetching.TimeoutingReader{Reader: bytes.NewReader([]byte{42}), Timeout: time.Minute}
	b := make([]byte, 1)
	n, err := trc.Read(b)
	if err != nil {
		t.Fatalf("Expected read to succeed, but got error: %v", err)
	}
	if n != 1 {
		t.Fatalf("Read %d instead of 1 byte.", n)
	}
}
