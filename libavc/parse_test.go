package libavc

import (
	"io/ioutil"
	"testing"
)

func TestParseH264(t *testing.T) {
	file, err := ioutil.ReadFile("../001.h264")
	if err != nil {
		panic(file)
	}
	//units := ParseNalUnits(file)
	frame := IsKeyFrame(file)
	//println(units)
	println(frame)
}
