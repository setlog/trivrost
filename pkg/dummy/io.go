package dummy

import "io"

type ReadCloser struct {
	Data []byte
	pos  int
}

func (rc *ReadCloser) Read(p []byte) (n int, err error) {
	if rc.Data == nil {
		return 0, nil
	}
	var bytesWritten int
	if len(p) < len(rc.Data)-rc.pos {
		bytesWritten = rc.write(p, len(p))
	} else {
		bytesWritten = rc.write(p, len(rc.Data)-rc.pos)
	}
	if rc.pos == len(rc.Data) {
		return bytesWritten, io.EOF
	}
	return bytesWritten, nil
}

func (rc *ReadCloser) write(p []byte, count int) int {
	x := 0
	for i := rc.pos; i < rc.pos+count; i++ {
		p[x] = rc.Data[i]
		x++
	}
	rc.pos += x
	return x
}

func (rc *ReadCloser) Close() error {
	return nil
}
