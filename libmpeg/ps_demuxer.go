package libmpeg

import (
	"avformat/libavc"
	"avformat/utils"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type deHandler func(buffer *utils.ByteBuffer, keyFrame bool, streamType int, pts, dts int64)

type DeMuxer struct {
	handler          deHandler
	packetHeader     *PacketHeader
	systemHeader     *SystemHeader
	programStreamMap *ProgramStreamMap
	lastPesPacket    *PESPacket
	currentPesPacket *PESPacket

	packet     *utils.Packet
	streamType byte
}

func (d *DeMuxer) Close() {
	//回调最后一帧
	if d.packet.Data().Size() > 0 {
		d.callback()
	}
}

func (d *DeMuxer) callback() {
	var keyFrame bool
	switch d.lastPesPacket.streamId {
	case StreamIdAudio:
		keyFrame = true
		break
	case StreamIdVideo, StreamIdH624:
		keyFrame = libavc.IsKeyFrameFromBuffer(d.packet.Data())
		break
	}

	d.handler(d.packet.Data(), keyFrame, int(d.streamType), d.packet.Pts(), d.packet.Dts())
	d.packet.Release()
}

// Input Reference from https://github.com/ireader/media-server/blob/master/libmpeg/source/mpeg-ps-dec.c
func (d *DeMuxer) Input(data []byte) int {
	n, i, consume := 0, 0, 0
	//保存第一个pes的开始位置
	//每次Input如果没有读取到完整的一帧，回退到第一个pes的位置
	//内部不做内存拷贝，ByteBuffer只是浅拷贝
	var firstPesIndex int
	length := len(data)
	d.packet.Release()

	for i = libavc.FindStartCode(data, 0); i >= 0 && i < length; i = libavc.FindStartCode(data, i) {
		i -= 3
		switch data[i+3] {
		case 0xBA:
			n = readPackHeader(d.packetHeader, data[i:])
			break
		case 0xBB:
			n = readSystemHeader(d.systemHeader, data[i:])
			break
		case 0xBC:
			n, _ = readProgramStreamMap(d.programStreamMap, data[i:])
			break
		case 0xB9: //end code
			break
		default:
			var esPacket []byte
			if firstPesIndex == 0 {
				firstPesIndex = i
			}

			esPacket, n = readPESPacket(d.currentPesPacket, data[i:])
			if n == 0 {
				goto END
			}

			element, ok := d.programStreamMap.findElementaryStream(data[i+3])
			if !ok {
				println(fmt.Sprintf("unKnow stream:%x", data[i+3]))
				break
			}

			if d.lastPesPacket == nil {
				pesPacket := *d.currentPesPacket
				d.lastPesPacket = &pesPacket
			}

			//读到下一包，才回调前一包
			//上一包和当前包的pts/streamId不一样,才回调
			if d.currentPesPacket.streamId != d.lastPesPacket.streamId || d.currentPesPacket.pts != d.lastPesPacket.pts {
				d.callback()
				*d.lastPesPacket = *d.currentPesPacket
				firstPesIndex = i
			}

			d.streamType = element.streamType
			if d.currentPesPacket.ptsDtsFlags&0x3 != 0 {
				d.packet.SetPts(d.currentPesPacket.pts)
				d.packet.SetDts(d.currentPesPacket.dts)
			}
			d.packet.Write(esPacket)
			d.currentPesPacket.Reset()
		}

		i += n
		consume = i
	}

END:
	if firstPesIndex != 0 {
		return firstPesIndex
	} else {
		return consume
	}
}

// Open 解复用本地文件
// @readCount 每次读取多少字节. <= 0 一次性读取完
func (d *DeMuxer) Open(path string, readCount int) error {
	if readCount <= 0 {
		all, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		d.Input(all)
		return nil
	} else {
		fi, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			fi.Close()
		}()

		reader := bufio.NewReader(fi)
		offset := 0
		buffer := make([]byte, readCount)

		for {
			r, err := reader.Read(buffer[offset:])
			if err != nil {
				if err == io.EOF {
					return nil
				} else {
					return err
				}
			}

			length := offset + r
			consume := d.Input(buffer[:length])
			offset = length - consume
			copy(buffer, buffer[consume:length])
		}
	}
}

func NewDeMuxer(handler deHandler) *DeMuxer {
	return &DeMuxer{
		handler:          handler,
		packetHeader:     &PacketHeader{},
		systemHeader:     &SystemHeader{},
		programStreamMap: &ProgramStreamMap{},
		currentPesPacket: NewPESPacket(),
		packet:           utils.NewPacket(),
	}
}
