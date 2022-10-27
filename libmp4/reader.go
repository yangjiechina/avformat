package libmp4

import (
	"avformat/utils"
	"encoding/binary"
)

type reader struct {
	data   []byte
	offset int
}

func newReader(data []byte) *reader {
	return &reader{data: data, offset: 0}
}

//func (r reader) hasNext() bool {
//Return value minus size of box header.
func (r *reader) nextSize() int64 {
	remain := len(r.data) - r.offset
	if remain < 4 {
		return -1
	}
	remain -= 4

	var isLargeSize bool
	size := int64(utils.BytesToUInt32(r.data[r.offset], r.data[r.offset+1], r.data[r.offset+2], r.data[r.offset+3]))
	if size == 0 {
		return 0
	} else if size == 1 {
		if remain < 8 {
			return -1
		}

		isLargeSize = true
		//r.offset += 4
		size = int64(binary.BigEndian.Uint64(r.data[r.offset+4:]))
		remain -= 8
		size -= 8
	} else {
		size -= 4
	}

	if size <= int64(remain) {
		if isLargeSize {
			r.offset += 12
		} else {
			r.offset += 4
		}
		return size
	} else {
		return -1
	}
}

func (r *reader) next(size int64) (string, []byte) {
	temp := r.offset + 4
	r.offset += int(size)
	return string(r.data[temp-4 : temp]), r.data[temp:r.offset]
}
