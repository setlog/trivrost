package dummy

import (
	"bytes"
)

type ByteReadCloser struct {
	*bytes.Buffer
}

func (rc *ByteReadCloser) Close() error { return nil }
