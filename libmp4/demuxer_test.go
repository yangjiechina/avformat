package libmp4

import (
	"avformat/utils"
	"testing"
)

func TestMp4DeMuxer(t *testing.T) {
	muxer := NewDeMuxer(func(data []byte, pts, dts int64, mediaType utils.AVMediaType, id utils.AVCodecID) {
		switch id {
		case utils.AVCodecIdH264:
			break
		case utils.AVCodecIdH265:
			break
		}
	})
	err := muxer.Open("../232937384-1-208_baseline.mp4")
	if err != nil {
		panic(err)
	}

	for err = muxer.Read(); err != nil; err = muxer.Read() {

	}
	//muxer.Read("../LB1l2iXISzqK1RjSZFjXXblCFXa.mp4")
}
