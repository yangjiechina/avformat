package libmp4

const (
	markSampleDescription = 1 << 31
	markTimeToSample      = 1 << 30
	markSampleToChunk     = 1 << 29
	markChunkOffset       = 1 << 28
	markSampleSize        = 1 << 27
	markMediaHeader       = 1 << 26
	markEditLit           = 1 << 25
	markSyncSample        = 1 << 24
	//chunkOffset64     = 1 << 28
	//CompactSampleSize = 1 << 27
)

type sampleIndexEntry struct {
	pos       int64 // the position in the file.
	timestamp int64
	size      uint32
	keyFrame  bool
}

type Track struct {
	// mark the required box
	mark uint32

	//"vmhd""smhd""hmhd""sthd""nmhd"
	mediaHandlerType string
	metaData         MetaData

	sampleCount        uint32
	chunkCount         uint32
	currentSample      uint32
	sampleIndexEntries []*sampleIndexEntry

	mdhd *mediaHeaderBox
	stsd *sampleDescriptionBox
	stts *decodingTimeToSampleBox
	stsc *sampleToChunkBox
	stco *chunkOffsetBox
	stsz *sampleSizeBox
	stss *syncSampleBox
	elst *editListBox
}

func (t *Track) MetaData() MetaData {
	return t.metaData
}
