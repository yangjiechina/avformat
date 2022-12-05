package librtmp

import (
	"avformat/libflv"
	"avformat/utils"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HandshakeState byte

const (
	HandshakeStateUninitialized = HandshakeState(0) //after the client sends C0
	HandshakeStateVersionSent   = HandshakeState(1) //client waiting for S1
	HandshakeStateAckSent       = HandshakeState(2) //client waiting for S2
	HandshakeStateDone          = HandshakeState(3) //client receives S2

)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Puller struct {
	client         *utils.TCPClient
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
			p.processChunk(data[i:])
			return
		}
	}
}

func (p *Puller) onDisconnected(conn net.Conn, err error) {

}

func (p *Puller) parseUrl(addr string) error {
	if index := strings.Index(addr, ":"); index == -1 {
		return fmt.Errorf("the format of URL is invalid")
	} else {
		protocol := addr[:index]
		if "rtmp" != protocol {
			return fmt.Errorf("unknow protocol:%s", protocol)
		}
		p.protocol = protocol
	}

	parse, err := url.Parse(addr)
	if err != nil {
		return err
	}

	var port int
	if p := parse.Port(); "" != p {
		if port, err = strconv.Atoi(p); err != nil {
			return err
		}
	} else {
		port = DefaultPort
	}

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
	length := writer.ToBytes(p.commandBuffer[12:])

	chunk := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdSystem,
		timestamp:       0,
		MessageLength:   length,
		messageTypeId:   MessageTypeIDCommandAMF0,
		messageStreamId: 0,
	}

	chunk.ToBytes(p.commandBuffer)
	total := 12 + utils.MinInt(length, p.chunkSize)
	for i := length - p.chunkSize; i > 0; {
		minInt := utils.MinInt(i, p.chunkSize)
		chunk.chunkType = ChunkType3
		copy(p.commandBuffer[total+1:], p.commandBuffer[total:])
		chunk.ToBytes(p.commandBuffer[total:])
		i -= minInt
		total++
		total += minInt
	}
	_, _ = p.client.Write(p.commandBuffer[:total])
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
	length := writer.ToBytes(p.commandBuffer[12:])

	chunk := ChunkHeader{
		chunkType:       ChunkType0,
		chunkStreamId:   ChunkStreamIdSystem,
		timestamp:       0,
		MessageLength:   length,
		messageTypeId:   MessageTypeIDCommandAMF0,
		messageStreamId: int(streamId),
	}

	chunk.ToBytes(p.commandBuffer)
	total := 12 + utils.MinInt(length, p.chunkSize)
	for i := length - p.chunkSize; i > 0; {
		minInt := utils.MinInt(i, p.chunkSize)
		chunk.chunkType = ChunkType3
		copy(p.commandBuffer[total+1:], p.commandBuffer[total:])
		chunk.ToBytes(p.commandBuffer[total:])
		i -= minInt
		total++
		total += minInt
	}
	_, _ = p.client.Write(p.commandBuffer[:total])
}

func (p *Puller) processChunk(data []byte) error {
	//fmt.Printf("chunk data:%s\r\n", hex.EncodeToString(data))
	length, i := len(data), 0
	header := &ChunkHeader{}
	for i < length {
		n, err := readChunkHeader(data[i:], header)
		if err != nil {
			return err
		}

		i += n
		if header.chunkType != ChunkType3 && length-i < header.MessageLength {
			return fmt.Errorf("invalid data")
		}

		switch header.messageTypeId {
		case MessageTypeIDSetChunkSize:
			p.chunkSize = utils.BytesToInt(data[i : i+header.MessageLength])
			break
		case MessageTypeIDAbortMessage:
			break
		case MessageTypeIDAcknowledgement:
			break
		case MessageTypeIDUserControlMessage:
			if header.MessageLength < 6 {
				return fmt.Errorf("invalid data")
			}
			event := binary.BigEndian.Uint16(data[i:])
			value := binary.BigEndian.Uint32(data[i+2:])
			p.processUserControlMessage(UserControlMessageEvent(event), value)
			break
		case MessageTypeIDWindowAcknowledgementSize:
			p.windowSize = utils.BytesToInt(data[i : i+header.MessageLength])
			break
		case MessageTypeIDSetPeerBandWith:
			if header.MessageLength < 5 {
				return fmt.Errorf("invalid data")
			}
			p.bandwidth = int(binary.BigEndian.Uint32(data[i:]))
			//limit type 0-hard/1-soft/2-dynamic
			_ = data[i+4]
			p.sendWindowAcknowledgementSize()
			break
		case MessageTypeIDAudio:
			break
		case MessageTypeIDVideo:
			break
		case MessageTypeIDDataAMF0:
			break
		case MessageTypeIDDataAMF3:
			break
		case MessageTypeIDCommandAMF0:
			if amf0, err := libflv.DoReadAFM0(data[i:]); err != nil {
				return err
			} else {
				l := len(amf0)
				var command string
				if l == 0 {
					return fmt.Errorf("invalid data")
				}
				command, _ = amf0[0].(string)
				if "_result" == command {
					transactionId := amf0[1].(float64)
					if TransactionIDCreateStream == TransactionID(transactionId) {
						streamId := amf0[3].(float64)
						p.play(streamId)
					}
				}

			}
			break
		case MessageTypeIDCommandAMF3:
			break
		case MessageTypeIDSharedObjectAMF0:
			break
		case MessageTypeIDSharedObjectAMF3:
			break
		case MessageTypeIDAggregateMessage:
			break
		}

		i += header.MessageLength

	}
	//read chunk
	return nil
}

func (p *Puller) processUserControlMessage(event UserControlMessageEvent, value uint32) {
	switch event {
	case UserControlMessageEventStreamBegin:
		p.createStream()
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
