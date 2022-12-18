package librtmp

import (
	"avformat/libavc"
	"avformat/libflv"
	"avformat/utils"
	"os"
	"testing"
)

func TestRTMPPuller(t *testing.T) {
	url := "rtmp://ns8.indexforce.com/home/mystream"
	h264File, err := os.OpenFile("../rtmp.h264", os.O_WRONLY|os.O_CREATE, 132)
	if err != nil {
		panic(err)
	}
	defer func() {
		h264File.Close()
	}()

	videoCallbackBuffer := make([]byte, 1024*1204)
	videoBufferLength := 0
	videoLastPts := 0
	var extraData []byte
	puller := NewPuller(func(data []byte, ts int) {
		if ts != 0 {
			//payload
			codecId := libflv.VideoCodecId(videoCallbackBuffer[0] & 0xF)
			if codecId == libflv.VideoCodeIdH264 {
				pktType := videoCallbackBuffer[1]
				ct := (int(videoCallbackBuffer[2]) << 16) | (int(videoCallbackBuffer[3]) << 8) | int(videoCallbackBuffer[4])

				if pktType == 0 {
					b, err := libavc.ExtraDataToAnnexB(videoCallbackBuffer[5:])
					//if err != nil {
					//	return utils.AVCodecIdNONE, 0, err
					//}
					//d.videoExtraData = b
					println(b)
					println(err)
					println(ct)
					extraData = b
				} else if pktType == 1 {
					buffer := utils.NewByteBuffer()
					libavc.Mp4ToAnnexB(buffer, videoCallbackBuffer[5:], extraData)
					buffer.ReadTo(func(bytes []byte) {
						h264File.Write(bytes)
					})
					//return utils.AVCodecIdH264, ct, nil
				} else if pktType == 2 {
					//empty
				}

			} else {

			}
			videoBufferLength = 0
			videoLastPts += ts
		}

		copy(videoCallbackBuffer[videoBufferLength:], data)
		videoBufferLength += len(data)

	}, func(data []byte, ts int) {

	})

	err = puller.Open(url)
	if err != nil {
		panic(err)
	}

	select {}
}
