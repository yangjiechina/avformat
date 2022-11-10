package libavc

import (
	"avformat/utils"
	"encoding/binary"
)

var (
	StartCode3 = []byte{0x00, 0x00, 0x01}
	StartCode4 = []byte{0x00, 0x00, 0x00, 0x01}
)

func FindStartCode(p []byte, offset int) int {
	length := len(p)
	i := offset + 2

	for i < length {
		if p[i] > 1 {
			i += 3
		} else if p[i-1] != 0 {
			i += 2
		} else if (p[i-2] | (p[i] - 1)) != 0 {
			i++
		} else {
			i++
			break
		}
	}

	if i < length {
		return i
	} else {
		return -1
	}
}

func FindStartCodeFromBuffer(buffer *utils.ByteBuffer, offset int) int {
	length := buffer.Size()
	i := offset + 2

	for i < length {
		if buffer.At(i) > 1 {
			i += 3
		} else if buffer.At(i-1) != 0 {
			i += 2
		} else if (buffer.At(i-2) | (buffer.At(i) - 1)) != 0 {
			i++
		} else {
			i++
			break
		}
	}

	if i < length {
		return i
	} else {
		return -1
	}
}

func IsKeyFrame(p []byte) bool {
	index := 0
	for {
		index = FindStartCode(p, index)
		if index < 0 {
			return false
		}
		state := p[index]
		switch state & 0x1F {
		case H264NalSPS:
			break
		case H264NalPPS:
			break
		case H264NalSEI:
			break
		case H264NalIDRSlice:
			return true
		case H264NalSlice:
			return false
		default:
			return false
		}
	}
}

func IsKeyFrameFromBuffer(buffer *utils.ByteBuffer) bool {
	index := 0
	for {
		index = FindStartCodeFromBuffer(buffer, index)
		if index < 0 {
			return false
		}
		state := buffer.At(index)
		switch state & 0x1F {
		case H264NalSPS:
			break
		case H264NalPPS:
			break
		case H264NalSEI:
			break
		case H264NalIDRSlice:
			return true
		case H264NalSlice:
			return false
		default:
			return false
		}
	}
}

func ParseNalUnits(p []byte) int {
	for {
		index := FindStartCode(p, 0)
		state := p[index]
		switch state & 0x1F {
		case H264NalSlice:
		case H264NalIDRSlice:
			break
		}
	}
}

func Mp4ToAnnexB(buffer *utils.ByteBuffer, data, extra []byte) {
	length := len(data)
	outSize, spsSeen, ppsSeen := 0, false, false
	for index := 4; index < length; index += 4 {
		size := int(binary.BigEndian.Uint32(data[index-4:]))
		if size == 0 || length-index < size {
			break
		}
		unitType := data[index] & 0x1F
		switch unitType {
		case H264NalSPS:
			spsSeen = true
			break
		case H264NalPPS:
			ppsSeen = true
			break
		case H264NalIDRSlice:
			if !spsSeen && !ppsSeen {
				outSize += copyNalU(buffer, extra, outSize, false)
			}
			break
		}

		bytes := data[index : index+size]
		outSize += copyNalU(buffer, bytes, outSize, true)
		index += size
	}
}

func copyNalU(buffer *utils.ByteBuffer, data []byte, outSize int, append bool) int {
	var startCodeSize int

	if append {
		if outSize == 0 {
			startCodeSize = 4
		} else {
			startCodeSize = 3
		}

		if startCodeSize == 4 {
			buffer.Write(StartCode4)
		} else if startCodeSize != 0 {
			buffer.Write(StartCode3)
		}
	}

	buffer.Write(data)

	return startCodeSize + len(data)
}

func ExtraDataToAnnexB(src []byte) []byte {
	buffer := utils.NewByteBuffer(src)
	//unsigned int(8) configurationVersion = 1;
	//unsigned int(8) AVCProfileIndication;
	//unsigned int(8) profile_compatibility;
	//unsigned int(8) AVCLevelIndication;
	buffer.Skip(4)
	_ = buffer.ReadUInt8()&0x3 + 1
	unitNb := buffer.ReadUInt8() & 0x1f
	if unitNb == 0 {
		return nil
	}

	dstBuffer := utils.NewByteBuffer()
	spsDone := 0
	for unitNb != 0 {
		unitNb--
		size := int(buffer.ReadUInt16())
		dstBuffer.Write(StartCode4)
		readOffset := buffer.ReadOffset()
		dstBuffer.Write(src[readOffset : readOffset+size])
		buffer.Skip(size)

		bytes := buffer.ReadableBytes()
		spsDone++
		if bytes > 2 && unitNb == 0 && spsDone == 1 {
			unitNb = buffer.ReadUInt8()
		}
	}

	return dstBuffer.ToBytes()
}
