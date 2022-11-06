package libmp4

import (
	"avformat/libavc"
	"avformat/utils"
	"os"
	"testing"
)

func TestMp4DeMuxer(t *testing.T) {
	path := "../232937384-1-208_baseline.mp4"
	h264File, err := os.OpenFile(path+".h264", os.O_WRONLY|os.O_CREATE, 132)
	if err != nil {
		panic(err)
	}
	defer func() {
		h264File.Close()
	}()

	aacFile, err := os.OpenFile(path+".aac", os.O_WRONLY|os.O_CREATE, 132)
	if err != nil {
		panic(err)
	}
	defer func() {
		aacFile.Close()
	}()

	convertBuffer := utils.NewByteBuffer()
	var videoTrack *Track
	muxer := NewDeMuxer(func(data []byte, pts, dts int64, mediaType utils.AVMediaType, id utils.AVCodecID) {
		switch id {
		case utils.AVCodecIdH264:
			libavc.Mp4ToAnnexB(convertBuffer, data, videoTrack.MetaData().extraData())
			convertBuffer.ReadTo(func(bytes []byte) {
				h264File.Write(bytes)
			})
			break
		case utils.AVCodecIdH265:
			break
		case utils.AVCodecIdAAC:
			aacFile.Write(data)
			break
		}

		convertBuffer.Release()
	})

	if err := muxer.Open(path); err != nil {
		panic(err)
	}

	videoTrack = muxer.FindTrack(utils.AVMediaTypeVideo)
	if videoTrack == nil {
		panic("Not find for video track.")
	}

	for err = muxer.Read(); err == nil; err = muxer.Read() {

	}
	//muxer.Read("../LB1l2iXISzqK1RjSZFjXXblCFXa.mp4")
}
