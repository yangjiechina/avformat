package libmpeg

import (
	"avformat/utils"
	"io/ioutil"
	"os"
	"testing"
)

func TestTSDeMuxer(t *testing.T) {
	path := "../sample_1280x720.ts"
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	h264File, _ := os.OpenFile(path+".h264", os.O_WRONLY|os.O_CREATE, 132)
	defer func() {
		h264File.Close()
	}()

	muxer := NewTSDeMuxer(func(buffer utils.ByteBuffer, keyFrame bool, streamType int, pts, dts int64) {
		buffer.ReadTo(func(bytes []byte) {
			h264File.Write(bytes)
		})
	})

	length := len(file)
	for count := 0; length >= 188; count++ {
		i := len(file) - length
		err := muxer.doRead(file[i : i+188])
		if err != nil {
			panic(err)
		}
		length -= 188
	}

}
