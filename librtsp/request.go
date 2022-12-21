package librtsp

import (
	"fmt"
)

const (
	DefaultPort = 554
)

var (
	delimiter []byte
)

func init() {
	delimiter = make([]byte, 2)
	delimiter[0] = 0x0D
	delimiter[1] = 0x0A
}

type Request struct {
	//request line
	method  string
	url     string
	version string
	header  map[string]string
	body    string
}

func (r Request) toBytes(data []byte) int {
	var n int
	line := fmt.Sprintf("%s %s RTSP/%s\r\n", r.method, r.url, r.version)
	copy(data, line)
	n += len(line)
	for k, v := range r.header {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		copy(data[n:], header)
		n += len(header)

	}

	copy(data[n:], "\r\n")
	n += 2
	if r.body != "" {
		copy(data[n:], r.body)
		n += len(r.body)

	}
	return n
}
