package libmp4

import "avformat/utils"

const (
	markSampleDescription = 1 << 31
	markTimeToSample      = 1 << 30
	markSampleToChunk     = 1 << 29
	markChunkOffset       = 1 << 28
	markSampleSize        = 1 << 27
	//chunkOffset64     = 1 << 28
	//CompactSampleSize = 1 << 27
)

type track struct {
	// mark the required box
	mark uint32

	//"vmhd""smhd""hmhd""sthd""nmhd"
	mediaHandlerType string
	codecType        utils.AVMediaType
	codecId          utils.AVCodecID
	width            int
	height           int

	stsd *sampleDescriptionBox
	stts *decodingTimeToSampleBox
	stsc *sampleToChunkBox
	stco *chunkOffsetBox
	stsz *sampleSizeBox

	stss *syncSampleBox
}
