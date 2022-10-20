package base

/**
全都大端子序
*/

func WriteWORD(dst []byte, src uint16) {
	dst[0] = byte(src >> 8)
	dst[1] = byte(src)
}

func WriteDWORD(dst []byte, src uint32) {
	dst[0] = byte(src >> 24)
	dst[1] = byte(src >> 16)
	dst[2] = byte(src >> 8)
	dst[3] = byte(src)
}

func WriteInt(dst []byte, value, bytes int) {
	for i := 0; i < bytes; i++ {
		dst[i] = byte(value >> ((bytes - i - 1) * 8))
	}
}

func BytesToInt(b1 byte, b2 byte) int {
	return int(b1)<<8 | int(b2)
}

func BytesToUInt32(b1 byte, b2 byte, b3 byte, b4 byte) uint32 {
	return (uint32(b1) << 24) | (uint32(b2) << 16) | (uint32(b3) << 8) | uint32(b4)
}

func MinInt(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

// CalculateCrcMpeg2 from https://github.com/onurkaraoglan/go-crc32-mpeg2/blob/master/crcmpeg2.go
func CalculateCrcMpeg2(data []byte) uint32 {
	crc := uint32(0xffffffff)

	for _, v := range data {
		crc ^= uint32(v) << 24
		for i := 0; i < 8; i++ {
			if (crc & 0x80000000) != 0 {
				crc = (crc << 1) ^ 0x04C11DB7
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

func ReadBits() {

}

func WriteBits() {

}
