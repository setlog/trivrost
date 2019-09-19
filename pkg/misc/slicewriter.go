package misc

type ByteSliceWriter struct {
	Data []byte
}

func (dw *ByteSliceWriter) Write(p []byte) (n int, err error) {
	dw.Data = append(dw.Data, p...)
	return len(p), nil
}
