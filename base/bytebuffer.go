package base

// ByteBuffer unsafe thread
type ByteBuffer struct {
	data     [][]byte
	itemSize []int
	size     int
}

func (b *ByteBuffer) Write(data []byte, position, length int) {
	b.data = append(b.data, data[position:position+length])
	b.size += length
	b.itemSize = append(b.itemSize, b.size)
}

func (b *ByteBuffer) Size() int {
	return b.size
}

func (b *ByteBuffer) Release() {
	b.data = nil
	b.size = 0
	b.itemSize = nil
}

func (b *ByteBuffer) ToBytes() []byte {
	if b.data == nil {
		return nil
	}

	dst := make([]byte, b.size)
	offset := 0
	for _, bytes := range b.data {
		copy(dst[offset:], bytes)
		offset += len(bytes)
	}

	b.Release()

	return dst
}

func (b *ByteBuffer) PeekTo(handle func([]byte)) {
	for _, bytes := range b.data {
		handle(bytes)
	}
}

func (b *ByteBuffer) ReadTo(handle func([]byte)) {
	b.PeekTo(handle)
	b.Release()
}

func (b *ByteBuffer) At(index int) byte {
	for i, v := range b.itemSize {
		if index < v {
			if i > 0 {
				return b.data[i][index-b.itemSize[i-1]]
			} else {
				return b.data[i][index]
			}
		}
	}

	panic("slice index out of range")
}

func (b *ByteBuffer) ForEach(start int, handle func(i int, v byte) (bool, int)) {
	index := 0
	offset := start
	if start >= b.size {
		panic("slice index out of range")
	}

	for i := 0; i < len(b.data); i++ {
		bytes := b.data[i]
		length := len(bytes)
		total := index + length

		if offset > length {
			offset -= length
			index = total
			continue
		}
		for j := offset; j < length; {
			if broken, next := handle(index, bytes[i]); broken {
				return
			} else {
				if next >= b.size {
					panic("slice index out of range")
				}
				if next < total {
					j = next
				} else {
					offset = length - j
				}
			}
		}

		index = total
	}
}
