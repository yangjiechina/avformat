package libflv

import (
	"avformat/utils"
	"io/ioutil"
	"os"
	"testing"
)

func TestDeMuxer(t *testing.T) {
	path := "../test.flv"
	//path := "../85003520_27.flv"

	all, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

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

	muxer := NewDeMuxer(func(mediaType utils.AVMediaType, id utils.AVCodecID, data utils.ByteBuffer, pts, dts int64) {
		switch id {
		case utils.AVCodecIdH264:
			data.ReadTo(func(bytes []byte) {
				h264File.Write(bytes)
			})
			break
		case utils.AVCodecIdAAC:
			data.ReadTo(func(bytes []byte) {
				aacFile.Write(bytes)
			})
			break
		}
	})

	err = muxer.Read(all)
	if err != nil {
		panic(err)
	}
}
