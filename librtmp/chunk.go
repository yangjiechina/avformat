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

	MessageTypeIDSetChunkSize               = MessageTypeID(1)
	MessageTypeIDAbortMessage               = MessageTypeID(2)
	MessageTypeIDAcknowledgement            = MessageTypeID(3)
	MessageTypeIDUserControlMessage         = MessageTypeID(4)
	MessageTypeIDWindowAcknowledgementSize  = MessageTypeID(5)
	MessageTypeIDSetPeerBandWith            = MessageTypeID(6)
	MessageTypeIDAudio                      = MessageTypeID(8)
	MessageTypeIDVideo                      = MessageTypeID(9)
	MessageTypeIDDataAMF0                   = MessageTypeID(18) // MessageTypeIDDataAMF0 MessageTypeIDDataAMF3 metadata:creation time, duration, theme...
	MessageTypeIDDataAMF3                   = MessageTypeID(15)
	MessageTypeIDCommandAMF0                = MessageTypeID(20) // MessageTypeIDCommandAMF0 MessageTypeIDCommandAMF3  connect, createStream, publish, play, pause
	MessageTypeIDCommandAMF3                = MessageTypeID(17)
	MessageTypeIDSharedObjectAMF0           = MessageTypeID(19)
	MessageTypeIDSharedObjectAMF3           = MessageTypeID(16)
	MessageTypeIDAggregateMessage           = MessageTypeID(22)
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

	if h.chunkType < ChunkType3 {
		if h.timestamp >= 0xFFFFFF {
			utils.WriteUInt24(dst[index:], 0xFFFFFF)
		} else {
			utils.WriteUInt24(dst[index:], uint32(h.timestamp))
		}
		index += 3
	}

	if h.chunkType < ChunkType2 {
		utils.WriteUInt24(dst[index:], uint32(h.MessageLength))
		index += 4
		dst[index-1] = byte(h.messageTypeId)
	}

	if h.chunkType < ChunkType1 {
		binary.LittleEndian.PutUint32(dst[index:], uint32(h.messageStreamId))
		index += 4
	}

	if h.timestamp >= 0xFFFFFF {
		binary.BigEndian.PutUint32(dst[index:], uint32(h.timestamp))
		index += 4
	}

	return index
}

func readBasicHeader(src []byte) (ChunkType, ChunkStreamID, int, error) {
	t := ChunkType(src[0] >> 6)
	if t > 0x3 {
		return t, 0, 0, fmt.Errorf("unknow chunk type:%d", t)
	}

	switch src[0] & 0x3F {
	case 0:
		//64-(64+255)
		return t, ChunkStreamID(64 + int(src[1])), 2, nil
	case 1:
		//64-(65535+64)
		return t, ChunkStreamID(64 + int(binary.BigEndian.Uint16(src[1:]))), 3, nil
	//case 2:
	default:
		//1bytes
		return t, ChunkStreamID(src[0] & 0x3F), 1, nil
	}
}

func readChunkHeader(src []byte) (ChunkHeader, int, error) {
	t, csid, i, err := readBasicHeader(src)
	if err != nil {
		return ChunkHeader{}, 0, err
	}

	header := ChunkHeader{
		chunkType:     t,
		chunkStreamId: csid,
	}

	if header.chunkType < ChunkType3 {
		header.timestamp = utils.BytesToInt(src[i : i+3])
		i += 3
	}

	if header.chunkType < ChunkType2 {
		i += 3
		header.MessageLength = utils.BytesToInt(src[i-3 : i])
		header.messageTypeId = MessageTypeID(src[i])
		i++
	}

	if header.chunkType < ChunkType1 {
		header.messageStreamId = int(binary.LittleEndian.Uint32(src[i:]))
		i += 4
	}

	if header.timestamp == 0xFFFFFF {
		header.timestamp = int(binary.BigEndian.Uint32(src[i:]))
		i += 4
	}

	return header, i, nil
}
