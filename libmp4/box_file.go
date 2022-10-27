package libmp4

import "avformat/utils"

/**
root
*/
type file struct {
	containerBox
}

type fileTypeBox struct {
	finalBox
	majorBrand       uint32
	minorVersion     uint32
	compatibleBrands []uint32
}

/**
Box	Types: ‘pdin’
Container: File
Mandatory: No
Quantity: Zero or One
*/
type progressiveDownloadInformationBox struct {
	rate         []uint32 //download rate expressed	in bytes/second
	initialDelay []uint32 //playing suggested delay
}

type freeBox struct {
	finalBox
}

/**
Box Type: ‘mdat’
Container: File
Mandatory: No
Quantity: Zero or more
*/
type mediaDataBox struct {
	finalBox
}

/**
Box Type: ‘moov’
Container: File
Mandatory: Yes
Quantity: Exactly one
*/
type movieBox struct {
	fullBox
	containerBox
}

func parseFileTypeBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	ftyp := fileTypeBox{}
	ftyp.majorBrand = buffer.ReadUInt32()
	ftyp.minorVersion = buffer.ReadUInt32()
	length := len(data)
	for i := 8; i <= length && length%4 == 0; i++ {
		ftyp.compatibleBrands = append(ftyp.compatibleBrands, buffer.ReadUInt32())
		i += 4
	}
	return &ftyp, len(data), nil
}

func parseFreeBox(data []byte) (box, int, error) {
	return &freeBox{}, len(data), nil
}

func parseMediaDataBox(data []byte) (box, int, error) {
	return &mediaDataBox{}, len(data), nil
}

func parseMovieBox(data []byte) (box, int, error) {

	return &movieBox{}, containersBoxConsumeCount, nil
}
