package librtsp

type Writer struct {
	Buffer []byte
	Length int
}

func NewWriter(size int) *Writer {
	return &Writer{make([]byte, size), 0}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	copy(w.Buffer[w.Length:], p)
	w.Length += len(p)
	return len(p), err
}

func (w *Writer) Reset() {
	w.Length = 0
}
