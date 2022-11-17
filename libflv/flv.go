package libflv

//5bits reserved. must be 0
//1bit audio tags are present
//1bit reserved.
//1bit video tags are present
type typeFlag byte

func (t typeFlag) ExistAudio() bool {
	return t>>2&0x1 == 1
}
func (t typeFlag) ExistVideo() bool {
	return t&0x1 == 1
}

type header struct {
	signature  string //always `0x46 0x4c 0x56(flv)`
	version    byte
	flags      typeFlag
	dataOffset uint32 // value of 9 for version 1.
}

type body struct {
	preTagSize uint32
}

type tag struct {
	tagType           byte   //8-audio/9-video/18-script
	dataSize          uint32 //3bytes
	timestamp         uint32 //3bytes
	timestampExtended byte
	streamId          uint32 //3bytes

	data interface{}
}

type audioData struct {
	desc byte
	data []byte
}

func (a audioData) soundFormat() int {
	return int(a.desc >> 4 & 0xF)
}
func (a audioData) soundRate() int {
	return int(a.desc >> 2 & 0x3)
}
func (a audioData) soundSize() int {
	return int(a.desc >> 1 & 0x1)
}
func (a audioData) soundType() int {
	return int(a.desc & 0x1)
}

type videoData struct {
	frameType byte //4bits
	codeId    byte //4bits
	data      []byte
}
