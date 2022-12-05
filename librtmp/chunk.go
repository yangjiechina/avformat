package librtmp

import (
	"avformat/utils"
	"encoding/binary"
	"fmt"
)

//https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol
//https://rtmp.veriskope.com/pdf/rtmp_specification_1.0.pdf
type ChunkType byte
type ChunkStreamID int
type MessageTypeID int
type MessageStreamID int
type UserControlMessageEvent uint16
type TransactionID int

/*
ChunkHeader Format
Each chunk consists of a header and data. The header itself has
three parts:
+--------------+----------------+--------------------+--------------+
| Basic Header | Message Header | Extended Timestamp | ChunkHeader Data |
+--------------+----------------+--------------------+--------------+
| |
|<------------------- ChunkHeader Header ----------------->|
*/

/**
type 0
0 1 2 3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| timestamp |message length |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| message length (cont) |message type id| msg stream id |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| message stream id (cont) |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

const (
	ChunkType0 = ChunkType(0x00)
	ChunkType1 = ChunkType(0x01)
	ChunkType2 = ChunkType(0x02)
	ChunkType3 = ChunkType(0x03)

	ChunkStreamIdNetwork = ChunkStreamID(2)
	ChunkStreamIdSystem  = ChunkStreamID(3)
	ChunkStreamIdAudio   = ChunkStreamID(4)
	ChunkStreamIdVideo   = ChunkStreamID(6)
	ChunkStreamIdSource  = ChunkStreamID(8)

	MessageTypeIDSetChunkSize              = MessageTypeID(1)
	MessageTypeIDAbortMessage              = MessageTypeID(2)
	MessageTypeIDAcknowledgement           = MessageTypeID(3)
	MessageTypeIDUserControlMessage        = MessageTypeID(4)
	MessageTypeIDWindowAcknowledgementSize = MessageTypeID(5)
	MessageTypeIDSetPeerBandWith           = MessageTypeID(6)
	MessageTypeIDAudio                     = MessageTypeID(8)
	MessageTypeIDVideo                     = MessageTypeID(9)

	// MessageTypeIDDataAMF0 MessageTypeIDDataAMF3 metadata:creation time, duration, theme...
	MessageTypeIDDataAMF0 = MessageTypeID(18)
	MessageTypeIDDataAMF3 = MessageTypeID(15)
	// MessageTypeIDCommandAMF0 MessageTypeIDCommandAMF3  connect, createStream, publish, play, pause
	MessageTypeIDCommandAMF0 = MessageTypeID(20)
	MessageTypeIDCommandAMF3 = MessageTypeID(17)

	MessageTypeIDSharedObjectAMF0 = MessageTypeID(19)
	MessageTypeIDSharedObjectAMF3 = MessageTypeID(16)
	MessageTypeIDAggregateMessage = MessageTypeID(22)

	UserControlMessageEventStreamBegin      = UserControlMessageEvent(0x00)
	UserControlMessageEventStreamEOF        = UserControlMessageEvent(0x01)
	UserControlMessageEventStreamDry        = UserControlMessageEvent(0x02)
	UserControlMessageEventSetBufferLength  = UserControlMessageEvent(0x03)
	UserControlMessageEventStreamIsRecorded = UserControlMessageEvent(0x04)
	UserControlMessageEventPingRequest      = UserControlMessageEvent(0x06)
	UserControlMessageEventPingResponse     = UserControlMessageEvent(0x07)

	TransactionIDConnect      = TransactionID(1)
	TransactionIDCreateStream = TransactionID(2)
	TransactionIDPlay         = TransactionID(0)
	DefaultChunkSize          = 128
)

type ChunkHeader struct {
	//basic header
	chunkType     ChunkType     //1-3bytes.低6位等于0,2字节;低6位等于1,3字节
	chunkStreamId ChunkStreamID //customized by users

	timestamp       int
	MessageLength   int
	messageTypeId   MessageTypeID
	messageStreamId int //customized by users. LittleEndian
}

func (h ChunkHeader) ToBytes(dst []byte) int {
	var index int
	index++
	dst[0] = byte(h.chunkType) << 6
	if h.chunkStreamId <= 63 {
		dst[0] = dst[0] | 0x3
	} else if h.chunkStreamId <= 0xFF {
		dst[0] = dst[0] & 0xC0
		dst[1] = byte(h.chunkStreamId)
		index++
	} else if h.chunkStreamId <= 0xFFFF {
		dst[0] = dst[0] & 0xC0
		dst[0] = dst[0] | 0x1
		binary.BigEndian.PutUint16(dst[1:], uint16(h.chunkStreamId))
		index += 2
	}

	if h.chunkType != ChunkType3 {
		for i := 0; i < 3; i++ {
			if i == 0 {
				if h.timestamp >= 0xFFFFFF {
					utils.WriteUInt24(dst[index:], 0xFFFFFF)
					index += 3
				} else {
					utils.WriteUInt24(dst[index:], uint32(h.timestamp))
					index += 3
				}
				if h.chunkType == ChunkType2 {
					break
				}
			} else if i == 1 {
				utils.WriteUInt24(dst[index:], uint32(h.MessageLength))
				index += 4
				dst[index-1] = byte(h.messageTypeId)
				if h.chunkType == ChunkType1 {
					break
				}
			} else {
				binary.LittleEndian.PutUint32(dst[index:], uint32(h.messageStreamId))
				index += 4
			}
		}

		if h.timestamp >= 0xFFFFFF {
			binary.BigEndian.PutUint32(dst[index:], uint32(h.timestamp))
			index += 4
		}
	}

	return index
}

func readChunkHeader(src []byte, header *ChunkHeader) (int, error) {
	var i int
	header.chunkType = ChunkType(src[0] >> 6)
	if header.chunkType > 0x3 {
		return -1, fmt.Errorf("unknow chunk type:%d", header.chunkType)
	}

	i++
	switch src[0] & 0x3F {
	case 0:
		//64-(64+255)
		header.chunkStreamId = ChunkStreamID(64 + int(src[i]))
		break
	case 1:
		//64-(65535+64)
		header.chunkStreamId = ChunkStreamID(64 + int(binary.BigEndian.Uint16(src[i:])))
		i += 2
		break
	case 2:
	default:
		//1bytes
		header.chunkStreamId = ChunkStreamID(src[0] & 0x3F)
		break
	}

	if header.chunkType != ChunkType3 {
		for j := 0; j < 3; j++ {
			if j == 0 {
				i += 3
				header.timestamp = int(utils.BytesToUInt24(src[i-3], src[i-2], src[i-1]))
				if header.chunkType == ChunkType2 {
					break
				}
			} else if j == 1 {
				i += 4
				header.MessageLength = int(utils.BytesToUInt24(src[i-4], src[i-3], src[i-2]))
				header.messageTypeId = MessageTypeID(src[i-1])
				if header.chunkType == ChunkType1 {
					break
				}
			} else {
				binary.LittleEndian.Uint32(src[i:])
				i += 4
			}
		}

		if header.timestamp == 0xFFFFFF {
			header.timestamp = int(binary.BigEndian.Uint32(src[i:]))
			i += 4
		}
	}

	return i, nil
}
