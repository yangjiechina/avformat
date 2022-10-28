package libmp4

import "avformat/utils"

//stsd * 8.5.2 sample descriptions (codec types, initialization
//etc.)
//stts * 8.6.1.2 (decoding) time-to-sample
//ctts 8.6.1.3 (composition) time to sample
//cslg 8.6.1.4 composition to decode timeline mapping
//stsc * 8.7.4 sample-to-chunk, partial data-offset information
//stsz 8.7.3.2 sample sizes (framing)
//stz2 8.7.3.3 compact sample sizes (framing)
//stco * 8.7.5 chunk offset, partial data-offset information
//co64 8.7.5 64-bit chunk offset
//stss 8.6.2 sync sample table
//stsh 8.6.3 shadow sync sample table
//padb 8.7.6 sample padding bits
//stdp 8.7.6 sample degradation priority
//sdtp 8.6.4 independent and disposable samples
//sbgp 8.9.2 sample-to-group
//sgpd 8.9.3 sample group description
//subs 8.7.7 sub-sample information
//saiz 8.7.8 sample auxiliary information sizes
//saio 8.7.9 sample auxiliary information offsets

/**
Box	Types:	 ‘stsd’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type sampleDescriptionBox struct {
	fullBox
	containerBox
	entryCount uint32
}

/**
Box	Type:	 ‘stdp’
Container:	 Sample	Table	Box	(‘stbl’).
Mandatory:	No.
Quantity:	 Zero	or	one.
*/
type degradationPriorityBox struct {
	fullBox
	finalBox
	//sampleCount from 'stsz'
	//int i;
	//for (i=0; i < sample_count; i++) {
	//unsigned int(16) priority;
}

/**
Box	Type:	 ‘stts’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type decodingTimeToSampleBox struct {
	fullBox
	finalBox
	entryCount  uint32
	sampleCount []uint32
	sampleDelta []uint32
}

/**
Box	Type:	 ‘ctts’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type compositionTimeToSampleBox struct {
	fullBox
	finalBox
	entryCount uint32
}

/**
Box	Type:	 ‘cslg’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Extension	Properties	Box	(‘trep’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type compositionToDecodeBox struct {
	fullBox
	finalBox
	compositionToDTSShift        int64
	leastDecodeToDisplayDelta    int64
	greatestDecodeToDisplayDelta int64
	compositionStartTime         int64
	compositionEndTime           int64
}

/**
Box	Type:	 ‘stss’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type syncSampleBox struct {
	fullBox
	finalBox
	entryCount uint32
}

/**
Box	Type:	 ‘stsh’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type shadowSyncSampleBox struct {
	fullBox
	finalBox
	entryCount uint32
}

/**
Box	Types:	 ‘sdtp’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type independentAndDisposableSamplesBox struct {
	fullBox
	finalBox
	//sample_count from 'stsz' or ‘stz2’
}

/**
Box	Type:	 ‘stsz’,	‘stz2’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	Yes
Quantity:	 Exactly	one	variant	must	be	present
*/
type sampleSizeBox struct {
	fullBox
	finalBox
	sampleSize  uint32
	sampleCount uint32
	entrySize   []uint32
}

// stz2
type compactSampleSizeBox struct {
	fullBox
	finalBox
	fieldSize   uint8
	sampleCount uint32
	entrySize   []int64
	//	for (i=1; i <= sample_count; i++) {
	//	unsigned int(field_size) entry_size;
	//}
}

/**
Box	Type:	 ‘stsc’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type sampleToChunkBox struct {
	fullBox
	finalBox
	entryCount             uint32
	firstChunk             []uint32
	samplesPerChunk        []uint32
	sampleDescriptionIndex []uint32
}

/**
Box	Type:	 ‘stco’,	‘co64’
Container:	 Sample	Table	Box	(‘stbl’)
Mandatory:	Yes
Quantity:	 Exactly	one	variant	must	be	present
*/
type chunkOffsetBox struct {
	fullBox
	finalBox
	entryCount  uint32
	chunkOffset []uint32
}

//‘co64’
type chunkLargeOffsetBox struct {
	fullBox
	finalBox
	entryCount  uint32
	chunkOffset []uint64
}

/**
Box	Type:	 ‘padb’
Container:	 Sample	Table	(‘stbl’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type paddingBitsBox struct {
	fullBox
	finalBox
	sampleCount uint32
}

/**
Box	Type:	 ‘subs’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	(‘traf’)
Mandatory:	No
Quantity:	 Zero	or	more
*/
type subSampleInformationBox struct {
	fullBox
	finalBox
	entryCount uint32
}

/**
Box	Type:	 ‘saiz’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	('traf')
Mandatory:	No
Quantity:	 Zero	or	More
*/
type sampleAuxiliaryInformationSizesBox struct {
	fullBox
	finalBox
	auxInfoType           uint32
	auxInfoTypeParameter  uint32
	defaultSampleInfoSize uint8
	sampleCount           uint32
	sampleInfoSize        []uint8
}

/**
Box	Type:	 ‘saio’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	('traf')
Mandatory:	No
Quantity:	 Zero	or	More
*/
type sampleAuxiliaryInformationOffsetsBox struct {
	fullBox
	finalBox
	auxInfoType          uint32
	auxInfoTypeParameter uint32
	entryCount           uint32
	offset               []uint8
}

/**
Box	Type:	 ‘sbgp’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	(‘traf’)
Mandatory:	No
Quantity:	 Zero	or	more
*/
type sampleToGroupBox struct {
	fullBox
	finalBox

	groupingType          uint32
	groupingTypeParameter uint32
	entryCount            uint32
	sampleCount           []uint32
	groupDescriptionIndex []uint32
}

type sampleGroupDescriptionEntry interface {
}

/**
Box	Type:	 ‘sgpd’
Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	(‘traf’)
Mandatory:	No
Quantity:	 Zero	or	more,	with	one	for	each	Sample	to	Group	Box.
*/
type sampleGroupDescriptionBox struct {
	fullBox
	finalBox
	groupingType                  uint32
	defaultLength                 uint32
	defaultSampleDescriptionIndex uint32
	entryCount                    uint32
	descriptionLength             []uint32
	sampleEntries                 []sampleGroupDescriptionEntry
}

func parseSampleDescriptionBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsd := sampleDescriptionBox{fullBox: fullBox{version: version, flags: flags}}
	stsd.entryCount = buffer.ReadUInt32()

	return &stsd, 8, nil
}

func parseDecodingTimeToSampleBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stts := decodingTimeToSampleBox{fullBox: fullBox{version: version, flags: flags}}
	stts.entryCount = buffer.ReadUInt32()
	stts.sampleCount = make([]uint32, 0, stts.entryCount)
	stts.sampleDelta = make([]uint32, 0, stts.entryCount)
	for i := 0; i < int(stts.entryCount); i++ {
		stts.sampleCount = append(stts.sampleCount, buffer.ReadUInt32())
		stts.sampleDelta = append(stts.sampleDelta, buffer.ReadUInt32())
	}

	return &stts, len(data), nil
}

func parseCompositionTimeToSampleBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	ctts := compositionTimeToSampleBox{fullBox: fullBox{version: version, flags: flags}}

	ctts.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(ctts.entryCount); i++ {
		sampleCount := buffer.ReadUInt32()
		println(sampleCount)
		if version == 0 {
			sampleOffset := buffer.ReadUInt32()
			println(sampleOffset)
		} else if version == 1 {
			sampleOffset := buffer.ReadInt32()
			println(sampleOffset)
		}
	}
	return &ctts, len(data), nil
}

func parseCompositionToDecodeBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	cslg := compositionToDecodeBox{fullBox: fullBox{version: version, flags: flags}}
	if version == 0 {
		cslg.compositionToDTSShift = int64(buffer.ReadUInt32())
		cslg.leastDecodeToDisplayDelta = int64(buffer.ReadUInt32())
		cslg.greatestDecodeToDisplayDelta = int64(buffer.ReadUInt32())
		cslg.compositionStartTime = int64(buffer.ReadUInt32())
		cslg.compositionEndTime = int64(buffer.ReadUInt32())
	} else if version == 1 {
		cslg.compositionToDTSShift = buffer.ReadInt64()
		cslg.leastDecodeToDisplayDelta = buffer.ReadInt64()
		cslg.greatestDecodeToDisplayDelta = buffer.ReadInt64()
		cslg.compositionStartTime = buffer.ReadInt64()
		cslg.compositionEndTime = buffer.ReadInt64()
	}

	return &cslg, len(data), nil
}

func parseSampleToChunkBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsc := sampleToChunkBox{fullBox: fullBox{version: version, flags: flags}}
	stsc.entryCount = buffer.ReadUInt32()

	stsc.firstChunk = make([]uint32, stsc.entryCount)
	stsc.samplesPerChunk = make([]uint32, stsc.entryCount)
	stsc.sampleDescriptionIndex = make([]uint32, 0, stsc.entryCount)
	for i := 0; i < int(stsc.entryCount); i++ {
		stsc.firstChunk = append(stsc.firstChunk, buffer.ReadUInt32())
		stsc.samplesPerChunk = append(stsc.samplesPerChunk, buffer.ReadUInt32())
		stsc.sampleDescriptionIndex = append(stsc.sampleDescriptionIndex, buffer.ReadUInt32())
	}

	return &stsc, len(data), nil
}

func parseSampleSizeBoxes(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsz := sampleSizeBox{fullBox: fullBox{version: version, flags: flags}}
	stsz.sampleSize = buffer.ReadUInt32()
	stsz.sampleCount = buffer.ReadUInt32()
	if stsz.sampleSize == 0 {
		for i := 0; i < int(stsz.sampleCount); i++ {
			stsz.entrySize = append(stsz.entrySize, buffer.ReadUInt32())
		}
	}

	return &stsz, len(data), nil
}

func parseCompactSampleSizeBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stz2 := compactSampleSizeBox{fullBox: fullBox{version: version, flags: flags}}
	stz2.fieldSize = buffer.ReadUInt8()
	stz2.sampleCount = buffer.ReadUInt32()
	for i := 0; i < int(stz2.sampleCount); i++ {
		//unsigned int(field_size) entry_size
		switch stz2.fieldSize {
		case 8:
			buffer.ReadUInt8()
			break
		case 16:
			buffer.ReadUInt16()
			break
		case 32:
			buffer.ReadUInt32()
			break
		case 64:
			buffer.ReadUInt64()
			break
		}
	}

	return &stz2, len(data), nil
}

func parseChunkOffsetBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stco := chunkOffsetBox{fullBox: fullBox{version: version, flags: flags}}
	stco.entryCount = buffer.ReadUInt32()
	stco.chunkOffset = make([]uint32, 0, stco.entryCount)
	for i := 0; i < int(stco.entryCount); i++ {
		stco.chunkOffset = append(stco.chunkOffset, buffer.ReadUInt32())
	}

	return &stco, len(data), nil
}

func parseChunkLargeOffsetBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	co64 := chunkLargeOffsetBox{fullBox: fullBox{version: version, flags: flags}}
	co64.entryCount = buffer.ReadUInt32()

	co64.chunkOffset = make([]uint64, co64.entryCount)
	for i := 0; i < int(co64.entryCount); i++ {
		co64.chunkOffset = append(co64.chunkOffset, buffer.ReadUInt64())
	}

	return &co64, len(data), nil
}

func parseSyncSampleBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stss := syncSampleBox{fullBox: fullBox{version: version, flags: flags}}

	stss.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stss.entryCount); i++ {
		sampleCount := buffer.ReadUInt32()
		println(sampleCount)
	}

	return &stss, len(data), nil
}

func parseShadowSyncSampleBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsh := shadowSyncSampleBox{fullBox: fullBox{version: version, flags: flags}}

	stsh.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stsh.entryCount); i++ {
		shadowedSampleNumber := buffer.ReadUInt32()
		syncSampleNumber := buffer.ReadUInt32()
		println(shadowedSampleNumber)
		println(syncSampleNumber)
	}

	return &stsh, len(data), nil
}

func parsePaddingBitsBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	padb := paddingBitsBox{fullBox: fullBox{version: version, flags: flags}}
	padb.sampleCount = buffer.ReadUInt32()
	for i := 0; i < int(padb.sampleCount+1)/2; i++ {
		//bit(1) reserved = 0;
		//bit(3) pad1;
		//bit(1) reserved = 0;
		//bit(3) pad2;
	}

	return &padb, len(data), nil
}

func parseDegradationPriorityBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stdp := degradationPriorityBox{fullBox: fullBox{version: version, flags: flags}}
	return &stdp, len(data), nil
}

func parseIndependentAndDisposableSamplesBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	sdtp := independentAndDisposableSamplesBox{fullBox: fullBox{version: version, flags: flags}}

	//	for (i=0; i < sample_count; i++){
	//	unsigned int(2) is_leading;
	//	unsigned int(2) sample_depends_on;
	//	unsigned int(2) sample_is_depended_on;
	//	unsigned int(2) sample_has_redundancy;
	//}

	return &sdtp, len(data), nil
}

func parseSubSampleInformationBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	subs := subSampleInformationBox{fullBox: fullBox{version: version, flags: flags}}
	subs.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(subs.entryCount); i++ {
		sample_delta := buffer.ReadUInt32()
		println(sample_delta)
		subsample_count := buffer.ReadUInt16()
		if subsample_count > 0 {
			for j := 0; j < int(subsample_count); j++ {
				if version == 1 {
					subsample_size := buffer.ReadUInt32()
					println(subsample_size)
				} else {
					subsample_size := buffer.ReadUInt16()
					println(subsample_size)
				}
				subsample_priority := buffer.ReadUInt8()
				discardable := buffer.ReadUInt8()
				codec_specific_parameters := buffer.ReadUInt32()
				println(subsample_priority)
				println(discardable)
				println(codec_specific_parameters)

			}
		}
	}

	return &subs, len(data), nil
}

func parseSampleAuxiliaryInformationSizesBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	saiz := sampleAuxiliaryInformationSizesBox{fullBox: fullBox{version: version, flags: flags}}
	if saiz.flags&0x1 != 0 {
		saiz.auxInfoType = buffer.ReadUInt32()
		saiz.auxInfoTypeParameter = buffer.ReadUInt32()
	}
	saiz.defaultSampleInfoSize = buffer.ReadUInt8()
	if saiz.defaultSampleInfoSize == 0 {
		saiz.sampleInfoSize = data[len(data)-8:]
	}

	return &saiz, len(data), nil
}

func parseSampleAuxiliaryInformationOffsetsBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	saio := sampleAuxiliaryInformationOffsetsBox{fullBox: fullBox{version: version, flags: flags}}
	if saio.flags&0x1 != 0 {
		saio.auxInfoType = buffer.ReadUInt32()
		saio.auxInfoTypeParameter = buffer.ReadUInt32()
	}
	saio.entryCount = buffer.ReadUInt32()

	return &saio, len(data), nil
}

func parseSampleToGroupBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	sbgp := sampleToGroupBox{fullBox: fullBox{version: version, flags: flags}}
	sbgp.groupingType = buffer.ReadUInt32()
	if sbgp.version == 1 {
		sbgp.groupingTypeParameter = buffer.ReadUInt32()
	}
	sbgp.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(sbgp.entryCount); i++ {
		sbgp.sampleCount = append(sbgp.sampleCount, buffer.ReadUInt32())
		sbgp.groupDescriptionIndex = append(sbgp.groupDescriptionIndex, buffer.ReadUInt32())
	}

	return &sbgp, len(data), nil
}

func parseSampleGroupDescriptionBox(data []byte) (box, int, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	sgpd := sampleGroupDescriptionBox{fullBox: fullBox{version: version, flags: flags}}
	sgpd.groupingType = buffer.ReadUInt32()
	if sgpd.version == 1 {
		sgpd.defaultLength = buffer.ReadUInt32()
	} else if version >= 2 {
		sgpd.defaultSampleDescriptionIndex = buffer.ReadUInt32()
	}

	sgpd.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(sgpd.entryCount); i++ {
		if sgpd.version == 1 && sgpd.defaultLength == 0 {
			sgpd.descriptionLength = append(sgpd.descriptionLength, buffer.ReadUInt32())
		}
		//SampleGroupEntry (grouping_type);
		// an instance of a class derived from SampleGroupEntry
		// that is appropriate and permitted for the media type
	}

	return &sgpd, len(data), nil
}
