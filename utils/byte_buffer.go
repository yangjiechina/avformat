package utils

// ByteBuffer unsafe thread
type ByteBuffer struct {
	data       [][]byte
	itemSize   []int
	size       int
	readOffset int
}

func NewByteBuffer(data ...[]byte) *ByteBuffer {
	buffer := &ByteBuffer{}
	for _, datum := range data {
		buffer.Write(datum, 0, len(datum))
	}
	return buffer
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
	b.readOffset = 0
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
	i1, i2 := b.offset()
	for i, bytes := range b.data[i1:] {
		if i == 0 {
			handle(bytes[i2:])
		} else {
			handle(bytes)
		}
	}
}

func (b *ByteBuffer) ReadTo(handle func([]byte)) {
	b.PeekTo(handle)
	b.Release()
}

func (b *ByteBuffer) offset() (int, int) {
	if len(b.itemSize) == 1 {
		return 0, b.readOffset
	}

	for i, v := range b.itemSize {
		if b.readOffset < v {
			return i, b.readOffset - b.itemSize[i-1]
		}
	}

	panic("slice index out of range")
}

func (b *ByteBuffer) At(index int) byte {
	if len(b.itemSize) == 1 {
		return b.data[0][index]
	}

	for i, v := range b.itemSize {
		if index < v {
			return b.data[i][index-b.itemSize[i-1]]
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

func (b *ByteBuffer) ReadByte() byte {
	i := b.PeekByte()
	b.readOffset++
	return i
}

func (b *ByteBuffer) ReadUInt16() uint16 {
	i := b.PeekUInt16()
	b.readOffset += 2
	return i
}

func (b *ByteBuffer) ReadUInt24() uint32 {
	i := b.PeekUInt24()
	b.readOffset += 3
	return i
}

func (b *ByteBuffer) ReadUInt32() uint32 {
	i := b.PeekUInt32()
	b.readOffset += 4
	return i
}

func (b *ByteBuffer) ReadUInt64() uint64 {
	i := b.PeekUInt64()
	b.readOffset += 8
	return i
}

func (b *ByteBuffer) PeekByte() byte {
	return b.At(b.readOffset)
}

func (b *ByteBuffer) PeekUInt16() uint16 {
	return BytesToUInt16(b.At(b.readOffset), b.At(b.readOffset+1))
}

func (b *ByteBuffer) PeekUInt24() uint32 {
	return BytesToUInt24(b.At(b.readOffset), b.At(b.readOffset+1), b.At(b.readOffset+2))
}

func (b *ByteBuffer) PeekUInt32() uint32 {
	return BytesToUInt32(b.At(b.readOffset), b.At(b.readOffset+1), b.At(b.readOffset+2), b.At(b.readOffset+3))
}

func (b *ByteBuffer) PeekUInt64() uint64 {
	return BytesToUInt64(b.At(b.readOffset), b.At(b.readOffset+1), b.At(b.readOffset+2), b.At(b.readOffset+3), b.At(b.readOffset+4), b.At(b.readOffset+5), b.At(b.readOffset+6), b.At(b.readOffset+7))
}

func (b *ByteBuffer) Skip(count int) {
	b.readOffset += count
	if b.readOffset >= b.size {
		panic("slice index out of range")
	}
}
