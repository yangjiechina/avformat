package librtmp

import (
	"avformat/libflv"
	"avformat/utils"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	url2 "net/url"
	"strconv"
	"strings"
	"time"
)

type HandshakeState byte
type ParserState byte

const (
	HandshakeStateUninitialized = HandshakeState(0) //after the client sends C0
	HandshakeStateVersionSent   = HandshakeState(1) //client waiting for S1
	HandshakeStateAckSent       = HandshakeState(2) //client waiting for S2
	HandshakeStateDone          = HandshakeState(3) //client receives S2

	ParserStateInit              = ParserState(0)
	ParserStateBasicHeader       = ParserState(1)
	ParserStateTimestamp         = ParserState(2)
	ParserStateMessageLength     = ParserState(3)
	ParserStateStreamType        = ParserState(4)
	ParserStateStreamId          = ParserState(5)
	ParserStateExtendedTimestamp = ParserState(6)
	ParserStatePayload           = ParserState(7)
)

var (
	headerSize map[ChunkType]int
)

func init() {
	rand.Seed(time.Now().UnixNano())

	headerSize = map[ChunkType]int{
		ChunkType0: 11,
		ChunkType1: 7,
		ChunkType2: 3,
		ChunkType3: 0,
	}
}

type Message struct {
	ChunkHeader
	payload []byte
	length  int
}

type Parser struct {
	state             ParserState
	chunkType         ChunkType
	chunkStreamId     ChunkStreamID
	chunkStreamIdSize int
	headerSize        int
	offset            int
	extendedTimestamp bool
	msg               *Message
}

type OnVideo func(data []byte, ts int)
type OnAudio func(data []byte, ts int)

type Puller struct {
	client         utils.Transport
	handshakeState HandshakeState
	protocol       string
	url            string
	host           string
	port           int
	app            string
	streamName     string

	commandBuffer []byte
	chunkSize     int
	windowSize    int
	bandwidth     int

	messages []*Message
	parser   *Parser
	onVideo  OnVideo
	onAudio  OnAudio
}

func NewPuller(v OnVideo, a OnAudio) *Puller {
	return &Puller{commandBuffer: make([]byte, 1024*4), parser: &Parser{}, onVideo: v, onAudio: a, chunkSize: DefaultChunkSize}
}

func (p *Puller) findMessage(csid ChunkStreamID) *Message {
	for _, message := range p.messages {
		if message.chunkStreamId == csid {
			return message
		}
	}

	return nil
}
func (p *Puller) onPacket(conn net.Conn, data []byte) {
	length, i := len(data), 0
	for i < length {
		switch p.handshakeState {
		case HandshakeStateUninitialized:
			p.handshakeState = HandshakeStateVersionSent
			if data[i] < VERSION {
				fmt.Printf("unkonw rtmp version:%d", data[i])
			}
			i++
			break
		case HandshakeStateVersionSent:
			//5.2.3 The C1 and S1 packets are 1536 octets long.
			if length-i < HandshakePacketSize {
				fmt.Printf("the S1 length is less than 1536. current:%d", length-i)
			} else {
				//time
				_ = binary.BigEndian.Uint32(data[i:])
				//zero
				_ = binary.BigEndian.Uint32(data[i+4:])
				//random bytes
				i += HandshakePacketSize
				bytes := data[i-HandshakePacketSize : i]
				binary.BigEndian.PutUint32(bytes[4:], 0)
				p.client.Write(bytes)
				//send c2
				p.handshakeState = HandshakeStateAckSent
			}
			break
		case HandshakeStateAckSent:
			p.handshakeState = HandshakeStateDone
			p.connect()
			return
		case HandshakeStateDone:
			//chunks
			_ = p.processChunk(data)
			return
		}
	}
}

func (p *Puller) onDisconnected(conn net.Conn, err error) {

}

func (p *Puller) parseUrl(addr string) error {
	parse, err := url2.Parse(addr)
	if err != nil {
		return err
	}

	if "rtmp" != parse.Scheme {
		return fmt.Errorf("unknow protocol:%s", parse.Scheme)
	}

	var port int
	if p := parse.Port(); "" != p {
		if port, err = strconv.Atoi(p); err != nil {
			return err
		}
	} else {
		port = DefaultPort
	}
	p.protocol = parse.Scheme
	p.host = parse.Hostname()
	p.port = port

	split := strings.Split(parse.Path, "/")
	if len(split) > 1 {
		p.app = strings.Split(parse.Path, "/")[1]
	}
	if len(split) > 2 {
		p.streamName = strings.Split(parse.Path, "/")[2]
	}

	return nil
}

func (p *Puller) Open(addr string) error {
	if err := p.parseUrl(addr); err != nil {
		return err
	}

	client, err := utils.NewTCPClient(nil, p.host, p.port)
	if err != nil {
		return err
	}

	p.client = client
	p.client.SetOnPacketHandler(p.onPacket)
	p.client.SetOnDisconnectedHandler(p.onDisconnected)
	p.client.Read()
	p.chunkSize = DefaultChunkSize
	p.commandBuffer = make([]byte, 1024*4)

	return p.sendHandshake()
}

func (p *Puller) sendHandshake() error {
	bytes := make([]byte, HandshakePacketSize+1)
	bytes[0] = VERSION
	//ffmpeg后面写flash client version 有的写C1。
	//gen random bytes
	length := len(bytes)
	for i := 9; i < length; i++ {
		bytes[i] = byte(rand.Intn(255))
	}

	_, err := p.client.Write(bytes)
	if err != nil {
		return err
	}

	//waiting for s1
	p.handshakeState = HandshakeStateUninitialized
	return nil
}

/*
|----------- Command Message(connect) ------->|
| |
|<------- Window Acknowledgement Size --------|
| |
|<----------- Set Peer Bandwidth -------------|
| |
|-------- Window Acknowledgement Size ------->|
| |
|<------ User Control Message(StreamBegin) ---|
| |
|<------------ Command Message ---------------|
| (_result- connect response) |
| |
*/

func (p *Puller) connect() {
	//command message {name,transactionID,object}
	writer := libflv.NewAMF0Writer()
	writer.AddString("connect")
	writer.AddNumber(float64(TransactionIDConnect)) //transaction ID. Always set to 1. 对应_result中的number
	object := libflv.AMF0Object{}
	object.AddStringProperty("app", p.app)
	object.AddStringProperty("flashVer", "LNX 9,0,124,2")
	object.AddStringProperty("tcUrl", fmt.Sprintf("%s://%s:%d/%s", p.protocol, p.host, p.port, p.app))
	object.AddBooleanProperty("fpad", false)
	object.AddNumberProperty("capabilities", 15)
	object.AddNumberProperty("audioCodecs", 0x0FFF)   //client supports. 0x0FFF supports all audio codes
	object.AddNumberProperty("videoCodecs", 0x00FF)   //client supports. 0x00FF supports all video codes
	object.AddNumberProperty("videoFunction", 0x0001) //Indicates what special video  functions are supported. 0x0001 unused.
	writer.AddObject(&object)

	bytes := make([]byte, 256)
	length := writer.ToBytes(bytes)

	chunk := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdSystem,
		timestamp:       0,
		MessageLength:   length,
		messageTypeId:   MessageTypeIDCommandAMF0,
		messageStreamId: 0,
	}

	p.sendMessage(chunk, bytes[:length])
}

func (p *Puller) sendWindowAcknowledgementSize() {
	header := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdNetwork,
		timestamp:       0,
		MessageLength:   4,
		messageTypeId:   MessageTypeIDWindowAcknowledgementSize,
		messageStreamId: 0,
	}

	bytes := header.ToBytes(p.commandBuffer)
	binary.BigEndian.PutUint32(p.commandBuffer[bytes:], uint32(p.bandwidth))
	_, _ = p.client.Write(p.commandBuffer[:4+bytes])
}

func (p *Puller) createStream() {
	writer := libflv.NewAMF0Writer()
	writer.AddString("createStream")
	writer.AddNumber(float64(TransactionIDCreateStream)) //transaction ID. Always set to 1. 对应_result中的number
	writer.AddNull()                                     //
	length := writer.ToBytes(p.commandBuffer[12:])

	header := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdNetwork,
		timestamp:       0,
		MessageLength:   length,
		messageTypeId:   MessageTypeIDCommandAMF0,
		messageStreamId: 0,
	}

	header.ToBytes(p.commandBuffer)
	length += 12
	_, _ = p.client.Write(p.commandBuffer[:length])
}

func (p *Puller) play(streamId float64) {
	writer := libflv.NewAMF0Writer()
	writer.AddString("play")
	writer.AddNumber(float64(TransactionIDPlay)) //transaction ID. Always set to 1. 对应_result中的number
	writer.AddNull()
	writer.AddString(p.streamName)
	//start duration reset
	writer.AddNumber(-2)    //default
	writer.AddNumber(-1)    //default
	writer.AddBoolean(true) //flush any previous playlist

	bytes := make([]byte, 256)
	length := writer.ToBytes(bytes)

	chunk := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdSystem,
		timestamp:       0,
		MessageLength:   length,
		messageTypeId:   MessageTypeIDCommandAMF0,
		messageStreamId: int(streamId),
	}

	p.sendMessage(chunk, bytes[:length])
}

func (p *Puller) processChunk(data []byte) error {
	length, i := len(data), 0
	for i < length {
		switch p.parser.state {

		case ParserStateInit:
			*p.parser = Parser{}

			t := ChunkType(data[i] >> 6)
			if t > ChunkType3 {
				return fmt.Errorf("unknow chunk type:%d", t)
			}

			if data[i]&0x3F == 0 {
				p.parser.chunkStreamIdSize = 1
			} else if data[i]&0x3F == 1 {
				p.parser.chunkStreamIdSize = 2
			} else {
				p.parser.chunkStreamIdSize = 0
				p.parser.chunkStreamId = ChunkStreamID(data[i] & 0x3F)
			}

			p.parser.chunkType = t
			p.parser.headerSize = headerSize[p.parser.chunkType]
			p.parser.state = ParserStateBasicHeader
			i++
			break

		case ParserStateBasicHeader:
			for p.parser.chunkStreamIdSize > 0 {
				p.parser.chunkStreamId <<= 8
				p.parser.chunkStreamId |= ChunkStreamID(data[i])
				p.parser.chunkStreamIdSize--
				i++
			}

			if p.parser.chunkStreamIdSize == 0 {
				message := p.findMessage(p.parser.chunkStreamId)
				if message == nil {
					message = &Message{ChunkHeader{chunkType: p.parser.chunkType, chunkStreamId: p.parser.chunkStreamId}, nil, 0}
				}
				p.messages = append(p.messages, message)
				p.parser.msg = message

				if p.parser.chunkType < ChunkType3 {
					p.parser.state = ParserStateTimestamp
				} else {
					p.parser.state = ParserStatePayload
				}
			}
			break

		case ParserStateTimestamp:
			for p.parser.offset < 3 && i < length {
				p.parser.msg.timestamp <<= 8
				p.parser.msg.timestamp |= int(data[i])
				p.parser.offset++
				i++
			}

			if p.parser.offset == 3 {
				p.parser.extendedTimestamp = p.parser.msg.timestamp == 0xFFFFFF
				if p.parser.chunkType < ChunkType2 {
					p.parser.state = ParserStateMessageLength
				} else if p.parser.extendedTimestamp {
					p.parser.state = ParserStateExtendedTimestamp
				} else {
					p.parser.state = ParserStatePayload
				}
			}
			break

		case ParserStateMessageLength:
			for p.parser.offset < 6 && i < length {
				p.parser.msg.MessageLength <<= 8
				p.parser.msg.MessageLength |= int(data[i])
				p.parser.offset++
				i++
			}

			if p.parser.offset == 6 {
				p.parser.state = ParserStateStreamType
			}
			break

		case ParserStateStreamType:
			p.parser.msg.messageTypeId = MessageTypeID(data[i])
			i++
			p.parser.offset++
			if p.parser.chunkType == ChunkType0 {
				p.parser.state = ParserStateStreamId
			} else if p.parser.extendedTimestamp {
				p.parser.state = ParserStateExtendedTimestamp
			} else {
				p.parser.state = ParserStatePayload
			}
			break

		case ParserStateStreamId:
			for p.parser.offset < 11 && i < length {
				p.parser.msg.messageStreamId <<= 8
				p.parser.msg.messageStreamId |= int(data[i])
				p.parser.offset++
				i++
			}

			if p.parser.offset == 11 {
				if p.parser.extendedTimestamp {
					p.parser.state = ParserStateExtendedTimestamp
				} else {
					p.parser.state = ParserStatePayload
				}
			}
			break

		case ParserStateExtendedTimestamp:
			for p.parser.offset < 15 && i < length {
				p.parser.msg.timestamp <<= 8
				p.parser.msg.timestamp |= int(data[i])
				p.parser.offset++
				i++
			}

			if p.parser.offset == 15 {
				if p.parser.extendedTimestamp {
					p.parser.state = ParserStateExtendedTimestamp
				} else {
					p.parser.state = ParserStatePayload
				}
			}
			break

		case ParserStatePayload:
			remain := length - i
			need := p.parser.msg.MessageLength - p.parser.msg.length
			consume := utils.MinInt(need, p.chunkSize-(p.parser.msg.length%p.chunkSize))
			consume = utils.MinInt(consume, remain)
			if len(p.parser.msg.payload) < p.parser.msg.MessageLength {
				bytes := make([]byte, p.parser.msg.MessageLength+1024)
				copy(bytes, p.parser.msg.payload)
				p.parser.msg.payload = bytes
			}

			copy(p.parser.msg.payload[p.parser.msg.length:], data[i:i+consume])
			p.parser.msg.length += consume

			if p.parser.msg.length >= p.parser.msg.MessageLength {
				if p.parser.msg.length != 0 {
					err := p.processMessage(p.parser.msg.messageTypeId, p.parser.msg.payload[:p.parser.msg.length], p.parser.msg.timestamp)
					if err != nil {
						return err
					}
				}

				*p.parser.msg = Message{}
				p.parser.state = ParserStateInit
			} else if p.parser.msg.length%p.chunkSize == 0 {
				p.parser.state = ParserStateInit
			}

			i += consume
			break
		}
	}

	return nil
}

func (p *Puller) processUserControlMessage(event UserControlMessageEvent, value uint32) {
	switch event {
	case UserControlMessageEventStreamBegin:
		break
	case UserControlMessageEventStreamEOF:
		break
	case UserControlMessageEventStreamDry:
		break
	case UserControlMessageEventSetBufferLength:
		break
	case UserControlMessageEventStreamIsRecorded:
		break
	case UserControlMessageEventPingRequest:
		break
	case UserControlMessageEventPingResponse:
		break
	default:
		fmt.Printf("unkonw control event:%d", event)
		break
	}
}

func (p *Puller) processMessage(typeId MessageTypeID, data []byte, timestamp int) error {
	switch typeId {
	case MessageTypeIDSetChunkSize:
		p.chunkSize = utils.BytesToInt(data)
		break
	case MessageTypeIDAbortMessage:
		break
	case MessageTypeIDAcknowledgement:
		break
	case MessageTypeIDUserControlMessage:
		event := binary.BigEndian.Uint16(data)
		value := binary.BigEndian.Uint32(data[2:])
		p.processUserControlMessage(UserControlMessageEvent(event), value)
		break
	case MessageTypeIDWindowAcknowledgementSize:
		p.windowSize = utils.BytesToInt(data)
		break
	case MessageTypeIDSetPeerBandWith:
		p.bandwidth = int(binary.BigEndian.Uint32(data))
		//limit type 0-hard/1-soft/2-dynamic
		_ = data[4:]
		p.sendWindowAcknowledgementSize()
		break
	case MessageTypeIDAudio:
		p.onAudio(data, timestamp)
		break
	case MessageTypeIDVideo:
		p.onVideo(data, timestamp)
		break
	//case MessageTypeIDDataAMF0:
	//	break
	case MessageTypeIDDataAMF3:
		break
	case MessageTypeIDDataAMF0, MessageTypeIDCommandAMF0, MessageTypeIDSharedObjectAMF0:
		if amf0, err := libflv.DoReadAFM0(data); err != nil {
			return err
		} else {
			l := len(amf0)
			var command string
			if l == 0 {
				return fmt.Errorf("invalid data")
			}

			command, _ = amf0[0].(string)
			if "_result" == command || "_error" == command {
				transactionId := amf0[1].(float64)
				if TransactionIDConnect == TransactionID(transactionId) {
					p.createStream()
				} else if TransactionIDCreateStream == TransactionID(transactionId) {
					streamId := amf0[3].(float64)
					p.play(streamId)
				}
			}
		}
		break
	case MessageTypeIDCommandAMF3:
		break
	//case MessageTypeIDSharedObjectAMF0:
	//	break
	case MessageTypeIDSharedObjectAMF3:
		break
	case MessageTypeIDAggregateMessage:
		//unsupported
		break
	}

	return nil
}

func (p *Puller) sendMessage(header ChunkHeader, payload []byte) {
	length, index := len(payload), 0
	for length > 0 {
		minInt := utils.MinInt(p.chunkSize, length)
		if length != len(payload) {
			header.chunkType = ChunkType3
		}

		index += header.ToBytes(p.commandBuffer[index:])
		copy(p.commandBuffer[index:], payload[len(payload)-length:len(payload)-length+minInt])
		length -= minInt
		index += minInt
	}

	_, _ = p.client.Write(p.commandBuffer[:index])
}
