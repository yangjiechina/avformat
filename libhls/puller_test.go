package libhls

import (
	"os"
	"testing"
)

func TestPuller(t *testing.T) {
	tsOutputFile, _ := os.OpenFile("../hls.payload", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	defer func() {
		tsOutputFile.Close()
	}()

	// 请求数据
	url := "http://videocdn.renrenjiang.cn/Act-ss-m3u8-sd/1037359_1546064640169/1037359_1546064640169.m3u8"
	//url := "http://kbs-dokdo.gscdn.com/dokdo_300/_definst_/dokdo_300.stream/playlist.m3u8"
	puller := NewPuller()
	if err := puller.Open(url); err != nil {
		panic(err)
	}
	for data, err := puller.Read(); err == nil; data, err = puller.Read() {
		tsOutputFile.Write(data)
	}

	puller.Close()
}
