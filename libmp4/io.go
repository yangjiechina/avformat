package libmp4

import (
	"bufio"
	"os"
)

type fileReader struct {
	path   string
	offset int64
	handle *os.File
	reader *bufio.Reader
}

func (f *fileReader) open(path string) error {
	openFile, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	f.handle = openFile
	f.reader = bufio.NewReader(f.handle)
	return nil
}

func (f *fileReader) seek(offset int64) error {
	if f.offset == offset {
		return nil
	}
	offset, err := f.handle.Seek(offset, 0)
	if err != nil {
		return err
	}

	f.offset = offset
	return err
}

func (f *fileReader) read(dst []byte) (int64, error) {
	n, err := f.reader.Read(dst)
	f.offset += int64(n)
	return int64(n), err
}

func (f *fileReader) close() error {
	return f.handle.Close()
}
