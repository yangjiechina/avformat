package libflv

import (
	"avformat/libavc"
	"avformat/utils"
	"fmt"
)

type TagType byte
type VideoCodecId byte

const (
	TagTypeAudioData        = TagType(8)
	TagTypeVideoData        = TagType(9)
	TagTypeScriptDataObject = TagType(18) //metadata value https://en.wikipedia.org/wiki/Flash_Video

	VideoCodeIdH263     = VideoCodecId(2)
	VideoCodeIdSCREEN   = VideoCodecId(3)
	VideoCodeIdVP6      = VideoCodecId(4)
	VideoCodeIdVP6Alpha = VideoCodecId(5)
	VideoCodeIdScreenV2 = VideoCodecId(6)
	VideoCodeIdH264     = VideoCodecId(7)
)

type Handler func(mediaType utils.AVMediaType, id utils.AVCodecID, data utils.ByteBuffer, pts, dts int64)

type DeMuxer struct {
	videoExtraData []byte
	audioConfig    *utils.MPEG4AudioConfig
	aacADtsHeader  []byte
	handler        Handler

	/**
	duration: DOUBLE
	width: DOUBLE
	height: DOUBLE
	videodatarate: DOUBLE
	framerate: DOUBLE
	videocodecid: DOUBLE
	audiosamplerate: DOUBLE
	audiosamplesize: DOUBLE
	stereo: BOOL
	audiocodecid: DOUBLE
	filesize: DOUBLE
	*/
	metaData []interface{}
}

func NewDeMuxer(handler Handler) *DeMuxer {
	return &DeMuxer{handler: handler}
}

func (d *DeMuxer) readAudioTag(data []byte, dst utils.ByteBuffer) (utils.AVCodecID, error) {
	soundFormat := data[0] >> 4
	//soundRate := data[0] >> 2 & 3
	//soundSize := data[0] >> 1 & 0x1
	//soundType := data[0] & 0x1
	soundData := data[1:]

	//aac audio data
	if soundFormat == 10 {
		//audio specificConfig
		if soundData[0] == 0x0 {
			config, err := utils.ParseMpeg4AudioConfig(data[2:])
			if err != nil {
				return utils.AVCodecIdAAC, err
			}
			d.audioConfig = config
			d.aacADtsHeader = make([]byte, 7)
		} else if soundData[0] == 0x1 {
			utils.SetADtsHeader(d.aacADtsHeader, 0, d.audioConfig.ObjectType-1, d.audioConfig.SamplingIndex, d.audioConfig.ChanConfig, 7+len(soundData[1:]))
			dst.Write(d.aacADtsHeader)
			dst.Write(soundData[1:])
			return utils.AVCodecIdAAC, nil
		}
	}

	return utils.AVCodecIdNONE, nil
}

func (d *DeMuxer) readVideoTag(data []byte, dst utils.ByteBuffer) (utils.AVCodecID, int, error) {
	//frameType := data[0] >> 4 & 0xF
	codecId := VideoCodecId(data[0] & 0xF)
	if codecId == VideoCodeIdH264 {
		pktType := data[1]
		ct := (int(data[2]) << 16) | (int(data[3]) << 8) | int(data[4])
		if pktType == 0 {
			b, err := libavc.ExtraDataToAnnexB(data[5:])
			if err != nil {
				return utils.AVCodecIdNONE, 0, err
			}
			d.videoExtraData = b
		} else if pktType == 1 {
			libavc.Mp4ToAnnexB(dst, data[5:], d.videoExtraData)
			return utils.AVCodecIdH264, ct, nil
		} else if pktType == 2 {
			//empty
		}

	} else {

	}

	return utils.AVCodecIdNONE, 0, nil
}

func (d *DeMuxer) readScriptDataObject(data []byte) error {
	buffer := utils.NewByteBuffer(data)

	if err := buffer.PeekCount(1); err != nil {
		return err
	}

	metaData, err := DoReadAFM0FromBuffer(buffer)
	if err != nil {
		return err
	}
	if len(metaData) <= 0 {
		return fmt.Errorf("invalid data")
	}
	if s, ok := metaData[0].(string); s == "" || !ok {
		return fmt.Errorf("not find the ONMETADATA of AMF0")
	}

	d.metaData = metaData
	return nil
}

func (d *DeMuxer) Read(data []byte) error {
	buffer := utils.NewByteBuffer(data)
	if err := buffer.PeekCount(9); err != nil {
		return err
	}

	if data[0] != 0x46 || data[1] != 0x4C || data[2] != 0x56 {
		return fmt.Errorf("invalid data")
	}

	buffer.Skip(3)
	h := header{}
	h.version = buffer.ReadUInt8()
	h.flags = typeFlag(buffer.ReadUInt8())
	h.dataOffset = buffer.ReadUInt32()

	if h.version == 1 && h.dataOffset != 9 {
		return fmt.Errorf("invalid data")
	}

	if !h.flags.ExistAudio() && !h.flags.ExistAudio() {
		return fmt.Errorf("invalid data")
	}

	callBackBuffer := utils.NewByteBuffer()
	//pre size length + tag header size
	for buffer.ReadableBytes() > 15 {
		//preSize
		_ = buffer.ReadUInt32()
		tagType := buffer.ReadUInt8()
		dataSize := int(buffer.ReadUInt24())
		timestamp := int(buffer.ReadUInt24())
		timestamp |= int(buffer.ReadUInt8()) << 24

		// streamId always 0.
		_ = buffer.ReadUInt24()
		if buffer.ReadableBytes() < dataSize {
			break
		}

		offset := buffer.ReadOffset()
		dataBuffer := data[offset : offset+dataSize]
		buffer.Skip(dataSize)

		callBackBuffer.Clear()
		//data
		if TagTypeAudioData == TagType(tagType) {
			codeId, err := d.readAudioTag(dataBuffer[:dataSize], callBackBuffer)
			if err != nil {
				return err
			}
			if d.handler != nil && codeId != utils.AVCodecIdNONE {
				d.handler(utils.AVMediaTypeAudio, codeId, callBackBuffer, int64(timestamp), int64(timestamp))
			}
		} else if TagTypeVideoData == TagType(tagType) {
			codeId, ct, err := d.readVideoTag(dataBuffer[:dataSize], callBackBuffer)
			if err != nil {
				return err
			}
			if d.handler != nil && codeId != utils.AVCodecIdNONE {
				d.handler(utils.AVMediaTypeAudio, codeId, callBackBuffer, int64(timestamp+ct), int64(timestamp))
			}

		} else if TagTypeScriptDataObject == TagType(tagType) {
			if err := d.readScriptDataObject(dataBuffer[:dataSize]); err != nil {
				return err
			}
		}
	}

	return nil
}
