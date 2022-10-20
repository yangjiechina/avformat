package libmpeg

import (
	"avformat/base"
	"encoding/hex"
	"fmt"
)

const (
	PacketStartCode       = 0x000001BA
	SystemHeaderStartCode = 0x000001BB
	PSMStartCode          = 0x000001BC
	ProgramEndCode        = 0x000001B9

	trickModeControlTypeFastForward = 0x0
	trickModeControlTypeSlowMotion  = 0x1
	trickModeControlTypeFreezeFrame = 0x2
	trickModeControlTypeFastReverse = 0x3
	trickModeControlTypeSlowReverse = 0x4

	//Reference from https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/mpeg.h
	StreamTypeVideoMPEG1     = 0x01
	StreamTypeVideoMPEG2     = 0x02
	StreamTypeAudioMPEG1     = 0x03
	StreamTypeAudioMPEG2     = 0x04
	StreamTypePrivateSection = 0x05
	StreamTypePrivateData    = 0x06
	StreamTypeAudioAAC       = 0x0F
	StreamTypeVideoMpeg4     = 0x10
	StreamTypeVideoH264      = 0x1B
	StreamTypeVideoHEVC      = 0x24
	StreamTypeVideoCAVS      = 0x42
	StreamTypeAudioAC3       = 0x81

	StreamIdPrivateStream1 = 0xBD
	StreamIdPaddingStream  = 0xBE
	StreamIdPrivateStream2 = 0xBF
	StreamIdAudio          = 0xC0 //110x xxxx
	StreamIdVideo          = 0xE0 //1110 xxxx
	StreamIdH624           = 0xE2
)

var (
	streamTypes map[int]int
)

type StreamType int

func (s StreamType) isAudio() bool {
	return streamTypes[int(s)] == StreamIdAudio
}

func (s StreamType) isVideo() bool {
	return streamTypes[int(s)] == StreamIdVideo || streamTypes[int(s)] == StreamIdH624
}

func init() {
	streamTypes = map[int]int{
		StreamTypeVideoMPEG1:     StreamIdVideo,
		StreamTypeVideoMPEG2:     StreamIdVideo,
		StreamTypeAudioMPEG1:     StreamIdAudio,
		StreamTypeAudioMPEG2:     StreamIdAudio,
		StreamTypePrivateSection: StreamIdPrivateStream1,
		StreamTypePrivateData:    StreamIdPrivateStream1,
		StreamTypeAudioAAC:       StreamIdAudio,
		StreamTypeVideoMpeg4:     StreamIdVideo,
		StreamTypeVideoH264:      StreamIdVideo,
		StreamTypeVideoHEVC:      StreamIdVideo,
		StreamTypeVideoCAVS:      StreamIdVideo,
		StreamTypeAudioAC3:       StreamIdAudio,
	}
}

type PacketHeader struct {
	mpeg2                         bool
	systemClockReferenceBase      int64  //33
	systemClockReferenceExtension uint16 //9
	programMuxRate                uint32 //22

	stuffing []byte
}

func (h *PacketHeader) ToBytes(dst []byte) int {
	base.WriteDWORD(dst, PacketStartCode)
	//2bits 01
	dst[4] = 0x40
	//3bits [32..30]
	dst[4] = dst[4] | (byte(h.systemClockReferenceBase>>30) << 3)
	//1bit marker bit
	dst[4] = dst[4] | 0x4
	//15bits [29..15]
	//2bits 29 28
	dst[4] = dst[4] | byte(h.systemClockReferenceBase>>28&0x3)
	//8bits
	dst[5] = byte(h.systemClockReferenceBase >> 20)
	//5bits
	dst[6] = byte(h.systemClockReferenceBase >> 12 & 0xF8)
	dst[6] = dst[6] | 0x4
	//15bits [14:0]
	//2bits
	dst[6] = dst[6] | byte(h.systemClockReferenceBase>>13&0x3)
	dst[7] = byte(h.systemClockReferenceBase >> 5)
	//5bits
	dst[8] = byte(h.systemClockReferenceBase&0x1f) << 3
	dst[8] = dst[8] | 0x4

	dst[8] = dst[8] | byte(h.systemClockReferenceExtension>>7&0x3)
	dst[9] = byte(h.systemClockReferenceExtension) << 1
	//1bits mark bit
	dst[9] = dst[9] | 0x1

	dst[10] = byte(h.programMuxRate >> 14)
	dst[11] = byte(h.programMuxRate >> 6)
	dst[12] = byte(h.programMuxRate) << 2
	//2bits 2 mark bit
	dst[12] = dst[12] | 0x3

	//5bits reserved
	//3bits pack_stuffing_length
	dst[13] = 0xF8
	offset := 14
	if h.stuffing != nil {
		length := len(h.stuffing)
		dst[13] = dst[13] | byte(length)
		copy(dst[offset:], h.stuffing)
		offset += length
	}

	return offset
}

func readPackHeader(header *PacketHeader, src []byte) int {
	length := len(src)
	if length < 14 {
		return 0
	}
	header.mpeg2 = src[4]&0xC0 == 0
	//mpeg1 版本占用4bits 没有clockExtension reserved stuffingLength
	header.systemClockReferenceBase = int64(src[4]&0x38)<<27 | (int64(src[4]&0x3) << 28) | (int64(src[5]) << 20) | (int64(src[6]&0xF8) << 12) | (int64(src[6]&0x3) << 13) | (int64(src[7]) << 5) | (int64(src[8] & 0xF8 >> 3))

	header.systemClockReferenceExtension = uint16(src[8]&0x3) << 7
	header.systemClockReferenceExtension = header.systemClockReferenceExtension | uint16(src[9]>>1)

	header.programMuxRate = uint32(src[10]) << 14
	header.programMuxRate = header.programMuxRate | uint32(src[11])<<6
	header.programMuxRate = header.programMuxRate | uint32(src[12]>>2)

	l := 14 + int(src[13]&0x7)
	if l > length {
		return 0
	}

	header.stuffing = src[14:l]
	return l
}

//func (h *PacketHeader) SetStuffing(stuffing []byte) {
//	if len(stuffing) > 7 {
//		panic("Stuffing length is only 3 bits")
//	}
//	h.stuffing = stuffing
//}

func (h *PacketHeader) ToString() string {
	if h.stuffing == nil {
		return fmt.Sprintf("systemClockReferenceBase=%d\r\nsystemClockReferenceExtension=%d\r\nprogramMuxRate=%d\r\n", h.systemClockReferenceBase,
			h.systemClockReferenceExtension, h.programMuxRate)
	} else {
		return fmt.Sprintf("systemClockReferenceBase=%d\r\nsystemClockReferenceExtension=%d\r\nprogramMuxRate=%d\r\nstuffingLength=%d\r\nstuffing=%s\r\n", h.systemClockReferenceBase,
			h.systemClockReferenceExtension, h.programMuxRate, len(h.stuffing), hex.EncodeToString(h.stuffing))
	}
}

// streamHeader 3bytes.
type streamHeader struct {
	streamId byte
	//'11'
	bufferBoundScale byte   //1
	bufferSizeBound  uint16 //13
}

func (h *streamHeader) ToString() string {
	return fmt.Sprintf("streamId=%x\r\nbufferBoundScale=%d\r\nbufferSizeBound=%d\r\n", h.streamId, h.bufferBoundScale, h.bufferSizeBound)
}

func (h *streamHeader) ToBytes(data []byte) {
	data[0] = h.streamId
	data[1] = 0xc0
	data[1] = data[1] | (h.bufferBoundScale << 5)
	data[1] = data[1] | byte(h.bufferSizeBound&0x1F00>>8)
	data[2] = byte(h.bufferSizeBound & 0xFF)
}

type SystemHeader struct {
	//6 bytes
	rateBound                 uint32 //22
	audioBound                byte   //6 [0,32]
	fixedFlag                 byte   //1
	cspsFlag                  byte   //1
	systemAudioLockFlag       byte   //1
	systemVideoLockFlag       byte   //1
	videoBound                byte   //5 [0,16]
	packetRateRestrictionFlag byte   //1

	streams []streamHeader
}

func (h *SystemHeader) findStream(id byte) (streamHeader, bool) {
	if h.streams == nil {
		return streamHeader{}, false
	}
	for _, s := range h.streams {
		if s.streamId == id {
			return s, true
		}
	}

	return streamHeader{}, false
}

func readSystemHeader(header *SystemHeader, src []byte) int {
	length := len(src)
	if length < 6 {
		return 0
	}

	totalLength := 6 + (int(src[4])<<8 | int(src[5]))
	if totalLength > length {
		return 0
	}

	header.rateBound = uint32(src[6]) & 0x7E << 15
	header.rateBound = header.rateBound | uint32(src[7])<<7
	header.rateBound = header.rateBound | uint32(src[8]>>1)

	header.audioBound = src[9] >> 2
	header.fixedFlag = src[9] >> 1 & 0x1
	header.cspsFlag = src[9] & 0x1

	header.systemAudioLockFlag = src[10] >> 7
	header.systemVideoLockFlag = src[10] >> 6 & 0x1
	header.videoBound = src[10] & 0x1F
	header.packetRateRestrictionFlag = src[11] >> 7

	offset := 12
	for ; offset <= totalLength && (src[offset]&0x80) == 0x80 && (totalLength-offset)%3 == 0; offset += 3 {
		if _, ok := header.findStream(src[offset]); ok {
			continue
		}

		streamHeader := streamHeader{}
		streamHeader.streamId = src[offset]
		streamHeader.bufferBoundScale = src[offset+1] >> 5 & 0x1
		streamHeader.bufferSizeBound = uint16(src[offset+1]&0x1F) << 8
		streamHeader.bufferSizeBound = streamHeader.bufferSizeBound | uint16(src[offset+2])
		header.streams = append(header.streams, streamHeader)
	}

	return totalLength
}

func (h *SystemHeader) ToBytes(dst []byte) int {
	base.WriteDWORD(dst, SystemHeaderStartCode)
	dst[6] = 0x80
	dst[6] = dst[6] | byte(h.rateBound>>15)
	dst[7] = byte(h.rateBound >> 7)
	dst[8] = byte(h.rateBound) << 1
	//mark bit
	dst[8] = dst[8] | 0x1
	dst[9] = h.audioBound << 2
	dst[9] = dst[9] | (h.fixedFlag << 1)
	dst[9] = dst[9] | h.cspsFlag

	dst[10] = h.systemAudioLockFlag << 7
	dst[10] = dst[10] | (h.systemVideoLockFlag << 6)
	dst[10] = dst[10] | 0x20
	dst[10] = dst[10] | h.videoBound

	dst[11] = h.packetRateRestrictionFlag << 7
	dst[11] = dst[11] | 0x7F

	offset := 12
	for i, s := range h.streams {
		s.ToBytes(dst[offset:])
		offset += (i + 1) * 3
	}

	base.WriteWORD(dst[4:], uint16(offset-6))
	return offset
}

func (h *SystemHeader) ToString() string {
	sprintf := fmt.Sprintf(
		"rateBound=%d\r\n"+
			"audioBound=%d\r\n"+
			"fixedFlag=%d\r\n"+
			"cspsFlag=%d\r\n"+
			"systemAudioLockFlag=%d\r\n"+
			"systemVideoLockFlag=%d\r\n"+
			"videoBound=%d\r\n"+
			"packetRateRestrictionFlag=%d\r\n",
		h.rateBound,
		h.audioBound,
		h.fixedFlag,
		h.cspsFlag,
		h.systemAudioLockFlag,
		h.systemVideoLockFlag,
		h.videoBound,
		h.packetRateRestrictionFlag,
	)

	sprintf += "streams=[\r\n"
	for _, stream := range h.streams {
		sprintf += stream.ToString()
	}
	sprintf += "]\r\n"

	return sprintf
}

type ElementaryStream struct {
	streamType byte //2-34. 0x5 disable
	streamId   byte
	info       []byte
}

func (e ElementaryStream) ToString() string {
	if e.info == nil {
		return fmt.Sprintf("StreamType=%x\r\nstreamId=%x\r\n", e.streamType, e.streamId)
	} else {
		return fmt.Sprintf("StreamType=%x\r\nstreamId=%x\r\ninfo=%s\r\n", e.streamType, e.streamId, hex.EncodeToString(e.info))
	}
}

type ProgramStreamMap struct {
	streamId             byte
	currentNextIndicator byte //1 bit
	version              byte //5 bits
	info                 []byte
	elementaryStreams    []ElementaryStream
	crc32                uint32
}

func (h *ProgramStreamMap) findElementaryStream(streamId byte) (ElementaryStream, bool) {
	if h.elementaryStreams == nil {
		return ElementaryStream{}, false
	}

	for _, element := range h.elementaryStreams {
		if element.streamId == streamId {
			return element, true
		}
	}

	return ElementaryStream{}, false
}

func (h *ProgramStreamMap) ToString() string {
	var info string
	if h.info != nil {
		info = hex.EncodeToString(h.info)
	}

	var elements string
	if h.elementaryStreams != nil {
		for _, element := range h.elementaryStreams {
			elements += element.ToString()
		}
	}

	return fmt.Sprintf("streamId=%x\r\ncurrentNextIndicator=%d\r\nversion=%d\r\ninfo=%s\r\nelements=[\r\n%s]\r\ncrc32=%d\r\n",
		h.streamId, h.currentNextIndicator, h.version, info, elements, h.crc32)
}

func readProgramStreamMap(header *ProgramStreamMap, src []byte) (int, error) {
	length := len(src)
	if length < 16 {
		return 0, nil
	}
	totalLength := 6 + base.BytesToInt(src[4], src[5])
	if totalLength > length {
		return 0, nil
	}

	header.streamId = src[3]
	header.currentNextIndicator = src[6] >> 7
	header.version = src[6] & 0x1F

	infoLength := base.BytesToInt(src[8], src[9])
	offset := 10
	if infoLength > 0 {
		// +2 reserved elementary_stream_map_length
		if 10+2+infoLength > totalLength-4 {
			return 0, fmt.Errorf("bad bytes:%s", hex.EncodeToString(src))
		}

		offset += infoLength
		header.info = src[10:offset]
	}

	elementaryLength := base.BytesToInt(src[offset], src[offset+1])
	offset += 2
	if offset+elementaryLength > totalLength-4 {
		return 0, fmt.Errorf("bad bytes:%s", hex.EncodeToString(src))
	}

	for i := offset; i < offset+elementaryLength; i += 4 {
		eInfoLength := base.BytesToInt(src[i+2], src[i+3])

		if _, ok := header.findElementaryStream(src[i+1]); !ok {
			element := ElementaryStream{}
			element.streamType = src[i]
			element.streamId = src[i+1]

			if eInfoLength > 0 {
				//if i+4+eInfoLength > offset+elementaryLength {
				if i+4+eInfoLength > totalLength-4 {
					return 0, fmt.Errorf("bad bytes:%s", hex.EncodeToString(src))
				}
				element.info = src[i+4 : i+4+eInfoLength]
			}

			header.elementaryStreams = append(header.elementaryStreams, element)
		}

		i += eInfoLength
	}

	header.crc32 = base.BytesToUInt32(src[totalLength-4], src[totalLength-3], src[totalLength-2], src[totalLength-1])

	return totalLength, nil
}

func (h *ProgramStreamMap) ToBytes(dst []byte) int {
	base.WriteDWORD(dst, PSMStartCode)
	//current_next_indicator
	dst[6] = 0x80
	//reserved
	dst[6] = dst[6] | (0x3 << 5)
	//program_stream_map_version
	dst[6] = dst[6] | 0x1
	//reserved
	dst[7] = 0xFE
	//mark bit
	dst[7] = dst[7] | 0x1

	offset := 10
	if h.info != nil {
		length := len(h.info)
		copy(dst[offset:], h.info)
		base.WriteWORD(dst[8:], uint16(length))
		offset += length
	} else {
		base.WriteWORD(dst[8:], 0)
	}
	//elementary length
	offset += 2
	temp := offset
	for _, elementaryStream := range h.elementaryStreams {
		dst[offset] = elementaryStream.streamType
		offset++
		dst[offset] = elementaryStream.streamId
		offset += 3
		if elementaryStream.info != nil {
			length := len(elementaryStream.info)
			copy(dst[offset:], elementaryStream.info)
			base.WriteWORD(dst[offset-2:], uint16(length))
			offset += length
		} else {
			base.WriteWORD(dst[offset-2:], 0)
		}
	}

	elementaryLength := offset - temp
	base.WriteWORD(dst[temp-2:], uint16(elementaryLength))

	crc32 := base.CalculateCrcMpeg2(dst[:offset])
	base.WriteDWORD(dst[offset:], crc32)

	offset += 4
	base.WriteWORD(dst[4:], uint16(offset-6))

	return offset
}

type PESPacket struct {
	streamId     byte
	packetLength uint16

	//'10' 2 bslbf
	pesScramblingControl   byte //2
	pesPriority            byte //1
	dataAlignmentIndicator byte //1
	copyright              byte //1
	originalOrCopy         byte //1

	ptsDtsFlags            byte //2
	escrFlag               byte //1
	esRateFlag             byte //1
	dsmTrickModeFlag       byte //1
	additionalCopyInfoFlag byte //1
	pesCrcFlag             byte //1
	pesExtensionFlag       byte //1
	pesHeaderDataLength    byte //8

	escrBase      uint64
	escrExtension uint16 //9 bits
	esRate        uint32 //22 bits

	pts int64
	dts int64
}

func NewPESPacket() *PESPacket {
	return &PESPacket{
		pts: -1,
		dts: -1,
	}
}

func (p *PESPacket) Reset() {
	//p.streamId = 0
	p.packetLength = 0
	p.pesScramblingControl = 0
	p.pesPriority = 0
	//p.dataAlignmentIndicator = 0
	p.copyright = 0
	p.originalOrCopy = 0
	p.ptsDtsFlags = 0
	p.escrFlag = 0
	p.esRateFlag = 0
	p.dsmTrickModeFlag = 0
	p.additionalCopyInfoFlag = 0
	p.pesCrcFlag = 0
	p.pesExtensionFlag = 0
	p.pesHeaderDataLength = 0
	p.escrBase = 0
	p.escrExtension = 0
	p.esRate = 0
	//p.pts = -1
	//p.dts = -1
}

func (p *PESPacket) ToBytes(dst []byte) int {
	dst[0] = 0x00
	dst[1] = 0x00
	dst[2] = 0x01
	dst[3] = p.streamId

	dst[6] = 0x80
	dst[6] = dst[6] | p.pesScramblingControl<<4
	dst[6] = dst[6] | p.pesPriority<<3
	dst[6] = dst[6] | p.dataAlignmentIndicator<<2
	dst[6] = dst[6] | p.copyright<<1
	dst[6] = dst[6] | p.originalOrCopy

	dst[7] = p.ptsDtsFlags << 6
	dst[7] = dst[7] | p.escrFlag<<5
	dst[7] = dst[7] | p.esRateFlag<<4
	dst[7] = dst[7] | p.dsmTrickModeFlag<<3
	dst[7] = dst[7] | p.additionalCopyInfoFlag<<2
	dst[7] = dst[7] | p.pesCrcFlag<<1
	dst[7] = dst[7] | p.pesExtensionFlag

	//dst[8] = p.pesHeaderDataLength

	offset, temp := 9, 9
	if p.ptsDtsFlags&0x2 == 0x2 {
		//4bits
		dst[offset] = 0x20
		//PTS [32..30]
		dst[offset] = dst[offset] | (byte(p.pts>>30) << 1)
		//mark bit
		dst[offset] = dst[offset] | 0x1
		offset++
		dst[offset] = byte(p.pts >> 22)
		offset++
		dst[offset] = byte(p.pts >> 14)
		dst[offset] = dst[offset] | 0x1
		offset++
		dst[offset] = byte(p.pts >> 7)
		offset++
		dst[offset] = byte(p.pts) << 1
		dst[offset] = dst[offset] | 0x1

		offset++
	}

	if p.ptsDtsFlags&0x1 == 0x1 {
		dst[temp] = dst[temp] | 0x30

		//4bits `0001`
		dst[offset] = 0x10
		//PTS [32..30]
		dst[offset] = dst[offset] | (byte(p.dts>>30) << 1)
		//mark bit
		dst[offset] = dst[offset] | 0x1
		offset++
		dst[offset] = byte(p.dts >> 22)
		offset++
		dst[offset] = byte(p.dts >> 14)
		dst[offset] = dst[offset] | 0x1
		offset++
		dst[offset] = byte(p.dts >> 7)
		offset++
		dst[offset] = byte(p.dts) << 1
		dst[offset] = dst[offset] | 0x1

		offset++
	}

	p.pesHeaderDataLength = byte(offset - temp)
	dst[8] = p.pesHeaderDataLength
	return offset
}

func readPESPacket(p *PESPacket, src []byte) ([]byte, int) {
	length := len(src)
	if length < 9 {
		return nil, 0
	}

	p.streamId = src[3]
	packetLength := base.BytesToInt(src[4], src[5])
	totalLength := 6 + packetLength
	if totalLength > length {
		return nil, 0
	}
	//1011 1100 1 program_stream_map
	//1011 1101 2 private_stream_1
	//1011 1110 padding_stream
	//1011 1111 3 private_stream_2
	//110x xxxx ISO/IEC 13818-3 or ISO/IEC 11172-3 or ISO/IEC 13818-7 or ISO/IEC 14496-3 audio stream number x xxxx
	//1110 xxxx ITU-T Rec. H.262 | ISO/IEC 13818-2 or ISO/IEC 11172-2 or ISO/IEC 14496-2 video stream number xxxx
	//1111 0000 3 ECM_stream
	//1111 0001 3 EMM_stream
	//1111 0010 5 ITU-T Rec. H.222.0 | ISO/IEC 13818-1 Annex A or ISO/IEC 13818- 6_DSMCC_stream
	//1111 0011 2 ISO/IEC_13522_stream
	//1111 0100 6 ITU-T Rec. H.222.1 type A
	//1111 0101 6 ITU-T Rec. H.222.1 type B
	//1111 0110 6 ITU-T Rec. H.222.1 type C
	//1111 0111 6 ITU-T Rec. H.222.1 type D
	//1111 1000 6 ITU-T Rec. H.222.1 type E
	//1111 1001 7 ancillary_stream
	//1111 1010 ISO/IEC14496-1_SL-packetized_stream
	//1111 1011 ISO/IEC14496-1_FlexMux_stream
	//1111 1100 … 1111 1110 reserved data stream
	//1111 1111 4 program_stream_directory

	//if (stream_id != program_stream_map
	//&& stream_id != padding_stream
	//&& stream_id != private_stream_2
	//&& stream_id != ECM
	//&& stream_id != EMM
	//&& stream_id != program_stream_directory
	//&& stream_id != DSMCC_stream
	//&& stream_id != ITU-T Rec. H.222.1 type E stream)

	if p.streamId != 0xBC && p.streamId != 0xBE && p.streamId != 0xBF && p.streamId != 0xF0 && p.streamId != 0xF1 && p.streamId != 0xff && p.streamId != 0xF2 && p.streamId != 0xF8 {

	} else {
		panic("Other unfinished")
	}
	p.pesScramblingControl = src[6] >> 4 & 0x3
	p.pesPriority = src[6] >> 3 & 0x1
	p.dataAlignmentIndicator = src[6] >> 2 & 0x1
	p.copyright = src[6] >> 1 & 0x1
	p.originalOrCopy = src[6] & 0x1
	p.ptsDtsFlags = src[7] >> 6 & 0x3
	p.escrFlag = src[7] >> 5 & 0x1
	p.esRateFlag = src[7] >> 4 & 0x1
	p.dsmTrickModeFlag = src[7] >> 3 & 0x1
	p.additionalCopyInfoFlag = src[7] >> 2 & 0x1
	p.pesCrcFlag = src[7] >> 1 & 0x1
	p.pesExtensionFlag = src[7] & 0x1
	p.pesHeaderDataLength = src[8]

	offset := 9
	if p.ptsDtsFlags&0x2 == 0x2 {
		p.pts = int64(src[offset]&0xE)<<29 | (int64(src[offset+1]) << 22) | (int64(src[offset+2]&0xFE) << 14) | (int64(src[offset+3]) << 7) | int64(src[offset+4]>>1)
		offset += 5
	}

	if p.ptsDtsFlags&0x1 == 0x1 {
		p.dts = int64(src[offset]&0xE)<<29 | (int64(src[offset+1]) << 22) | (int64(src[offset+2]&0xFE) << 14) | (int64(src[offset+3]) << 7) | int64(src[offset+4]>>1)
		offset += 5
	}

	if p.escrFlag == 0x1 {
		p.escrBase = (uint64(src[offset]&0x38) << 27) | (uint64(src[offset]&0x3) << 28) | (uint64(src[offset+1]) << 20) | (uint64(src[offset+2]&0xF8) << 12) | (uint64(src[offset+2]&0x3) << 13) | (uint64(src[offset+3]) << 5) | (uint64(src[offset+4] >> 3))
		p.escrExtension = uint16(src[offset+4]&0x3<<6) | uint16(src[offset+5]>>1)
		offset += 6
	}

	if p.esRateFlag == 0x1 {
		p.esRate = (uint32(src[offset]&0x7F) << 15) | (uint32(src[offset+1]) << 7) | uint32(src[offset+2]>>1)
		offset += 3
	}

	return src[offset:totalLength], totalLength
}
