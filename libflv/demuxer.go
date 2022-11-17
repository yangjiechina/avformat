package libflv

import (
	"avformat/libavc"
	"avformat/utils"
	"fmt"
	"math"
)

type TagType byte
type VideoCodecId byte
type dataType byte

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

	dataTypeNumber     = dataType(0)
	dataTypeBoolean    = dataType(1)
	dataTypeString     = dataType(2)
	dataTypeObject     = dataType(3)
	dataTypeMovieClip  = dataType(4)
	dataTypeNull       = dataType(5)
	dataTypeUnDefined  = dataType(6)
	dataTypeReference  = dataType(7)
	dataTypeECMA       = dataType(8)
	dataTypeStrict     = dataType(10)
	dataTypeDate       = dataType(11)
	dataTypeLongString = dataType(12)

	endMark = 0x09
	//Action message Format
)

type Handler func(mediaType utils.AVMediaType, id utils.AVCodecID, data *utils.ByteBuffer, pts, dts int64)

type DeMuxer struct {
	videoExtraData []byte
	audioConfig    *utils.MPEG4AudioConfig
	aacADtsHeader  []byte
	handler        Handler

	amfObjects map[string]interface {
	}
}

func NewDeMuxer(handler Handler) *DeMuxer {
	return &DeMuxer{amfObjects: make(map[string]interface{}, 10), handler: handler}
}

func (d *DeMuxer) readAudioTag(data []byte, dst *utils.ByteBuffer) (utils.AVCodecID, error) {
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

func (d *DeMuxer) readVideoTag(data []byte, dst *utils.ByteBuffer) (utils.AVCodecID, int, error) {
	//frameType := data[0] >> 4 & 0xF
	codecId := VideoCodecId(data[0] & 0xF)
	if codecId == VideoCodeIdH264 {
		pktType := data[1]
		ct := (int(data[2]) << 16) | (int(data[3]) << 8) | int(data[4])
		if pktType == 0 {
			d.videoExtraData = libavc.ExtraDataToAnnexB(data[5:])
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

func (d *DeMuxer) readAMFString(buffer *utils.ByteBuffer, larger bool) (string, bool) {
	var length int
	if larger {
		length = int(buffer.ReadUInt32())
	} else {
		length = int(buffer.ReadUInt16())
	}

	if buffer.ReadableBytes() < length {
		return "", false
	}

	return string(buffer.ReadBytesWithShallowCopy(length)), true
}

func (d *DeMuxer) readAMFObject(buffer *utils.ByteBuffer, name string) error {

	t := buffer.ReadUInt8()
	//■ duration: DOUBLE
	//■ width: DOUBLE
	//■ height: DOUBLE
	//■ videodatarate: DOUBLE
	//■ framerate: DOUBLE
	//■ videocodecid: DOUBLE
	//■ audiosamplerate: DOUBLE
	//■ audiosamplesize: DOUBLE
	//■ stereo: BOOL
	//■ audiocodecid: DOUBLE
	//■ filesize: DOUBLE
	switch dataType(t) {
	case dataTypeNumber:
		//double
		d.amfObjects[name] = math.Float64frombits(buffer.ReadUInt64())
		break
	case dataTypeBoolean:
		d.amfObjects[name] = float64(buffer.ReadUInt8())
		break
	case dataTypeString:
		if amfString, b := d.readAMFString(buffer, false); !b {
			return fmt.Errorf("the AMF String type parsing failed")
		} else {
			d.amfObjects[name] = amfString
		}
		break
	case dataTypeObject:
		for bytes := buffer.ReadableBytes(); bytes > 6; bytes = buffer.ReadableBytes() {
			if _, b := d.readAMFString(buffer, false); b {
				if err := d.readAMFObject(buffer, ""); err != nil {
					return err
				}
			}
		}
		break
	case dataTypeMovieClip:
		//string
		break
	case dataTypeNull:
	case dataTypeUnDefined:
	case dataTypeReference:
		break
	case dataTypeECMA:
		//array
		//script data variable
		//max index
		_ = buffer.ReadUInt32()
		for bytes := buffer.ReadableBytes(); bytes > 6; bytes = buffer.ReadableBytes() {
			if amfString, b := d.readAMFString(buffer, false); amfString != "" && b {
				if err := d.readAMFObject(buffer, amfString); err != nil {
					return err
				}
			}
		}

		break
	case dataTypeStrict:
		count := int(buffer.ReadUInt32())
		for i := 0; i < count; i++ {
			if err := d.readAMFObject(buffer, ""); err != nil {
				return err
			}
		}
		break
	case dataTypeDate:
		d.amfObjects["DateTime"] = math.Float64frombits(buffer.ReadUInt64())
		d.amfObjects["LocalDateTimeOffset"] = buffer.ReadUInt16()
		break
	case dataTypeLongString:
		if amfString, b := d.readAMFString(buffer, true); !b {
			return fmt.Errorf("the AMF larger String type parsing failed")
		} else {
			d.amfObjects[name] = amfString
		}
		break
	}

	return nil
}

func (d *DeMuxer) readScriptDataObject(data []byte) error {

	buffer := utils.NewByteBuffer(data)
	t := buffer.ReadUInt8()
	if dataTypeString != dataType(t) {
		return fmt.Errorf("unknow type")
	}

	//onMetaData
	if name, _ := d.readAMFString(buffer, false); "onMetaData" != name {
		return fmt.Errorf("unknow type")
	}

	return d.readAMFObject(buffer, "onMetaData")
}

func (d *DeMuxer) Read(data []byte) error {
	if data[0] != 0x46 || data[1] != 0x4C || data[2] != 0x56 {
		return fmt.Errorf("invalid data")
	}
	buffer := utils.NewByteBuffer(data)
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
		//preSize := buffer.ReadUInt32()
		tagType := TagType(buffer.ReadUInt8())
		dataSize := int(buffer.ReadUInt24())
		timestamp := int(buffer.ReadUInt24())
		timestamp |= int(buffer.ReadUInt8()) << 24
		_ = buffer.ReadUInt24() // streamId always 0.

		if buffer.ReadableBytes() < dataSize {
			break
		}

		offset := buffer.ReadOffset()
		dataBuffer := data[offset : offset+dataSize]
		buffer.Skip(dataSize)

		callBackBuffer.Clear()
		//data
		if TagTypeAudioData == tagType {
			codeId, err := d.readAudioTag(dataBuffer[:dataSize], callBackBuffer)
			if err != nil {
				return err
			}
			if d.handler != nil && codeId != utils.AVCodecIdNONE {
				d.handler(utils.AVMediaTypeAudio, codeId, callBackBuffer, int64(timestamp), int64(timestamp))
			}
		} else if TagTypeVideoData == tagType {
			codeId, ct, err := d.readVideoTag(dataBuffer[:dataSize], callBackBuffer)
			if err != nil {
				return err
			}
			if d.handler != nil && codeId != utils.AVCodecIdNONE {
				d.handler(utils.AVMediaTypeAudio, codeId, callBackBuffer, int64(timestamp+ct), int64(timestamp))
			}

		} else if TagTypeScriptDataObject == tagType {
			if err := d.readScriptDataObject(dataBuffer[:dataSize]); err != nil {
				return err
			}
		}
	}

	return nil
}
