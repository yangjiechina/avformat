package libmp4

import "avformat/utils"

type parser func(data []byte) (box, int, error)

const (
	containersBoxConsumeCount = 4
)

var (
	parsers map[string]parser
)

func init() {
	parsers = map[string]parser{
		"ftyp": parseFileTypeBox,
		"free": parseFreeBox,
		"mdat": parseMediaDataBox,
		"moov": parseMovieBox,
		"mvhd": parseMovieHeaderBox,
		"udta": parseUserDataBox,
		"trak": parseTrackBox,
		"tkhd": parseTrackHeaderBox,
		"tref": parseTrackReferenceBox,
		"hint": parseTrackReferenceTypeBox,
		"cdsc": parseTrackReferenceTypeBox,
		"font": parseTrackReferenceTypeBox,
		"hind": parseTrackReferenceTypeBox,
		"vdep": parseTrackReferenceTypeBox,
		"vplx": parseTrackReferenceTypeBox,
		"subt": parseTrackReferenceTypeBox,
		"trgr": parseTrackGroupBox,
		"msrc": parseTrackGroupTypeBox,
		"edts": parseEditBox,
		"elst": parseEditListBox,
		"mdia": parseMediaBox,
		"mdhd": parseMediaHeaderBox,
		"hdlr": parseHandlerReferenceBox,
		"elng": parseExtendedLanguageBox,
		"minf": parseMediaInformationBox,
		"vmhd": parseVideoMediaHeaderBox,
		"smhd": parseSoundMediaHeaderBox,
		"hmhd": parseHintMediaHeaderBox,
		"sthd": parseSubtitleMediaHeaderBox,
		"nmhd": parseNullMediaHeaderBox,
		"dinf": parseDataInformationBox,
		"dref": parseDataReferenceBox,
		"url":  parseDataEntryUrlBox,
		"urn":  parseDataEntryUrnBox,
		"stbl": parseSampleTableBox,
	}
}

func parseSampleDescriptionBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsd := sampleDescriptionBox{fullBox: fullBox{version: version, flags: flags}}
	stsd.entryCount = buffer.ReadUInt32()

	return &stsd, nil
}

func parseDecodingTimeToSampleBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stts := decodingTimeToSampleBox{fullBox: fullBox{version: version, flags: flags}}
	stts.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stts.entryCount); i++ {
		sampleCount := buffer.ReadUInt32()
		sampleDelta := buffer.ReadUInt32()
		println(sampleCount)
		println(sampleDelta)
	}

	return &stts, nil
}

func parseCompositionTimeToSampleBox(data []byte) (box, error) {
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
	return &ctts, nil
}

func parseCompositionToDecodeBox(data []byte) (box, error) {
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

	return &cslg, nil
}

func parseSyncSampleBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stss := syncSampleBox{fullBox: fullBox{version: version, flags: flags}}

	stss.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stss.entryCount); i++ {
		sampleCount := buffer.ReadUInt32()
		println(sampleCount)
	}

	return &stss, nil
}

func parseShadowSyncSampleBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsh := syncSampleBox{fullBox: fullBox{version: version, flags: flags}}

	stsh.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stsh.entryCount); i++ {
		shadowedSampleNumber := buffer.ReadUInt32()
		syncSampleNumber := buffer.ReadUInt32()
		println(shadowedSampleNumber)
		println(syncSampleNumber)
	}

	return &stsh, nil
}

func parseIndependentAndDisposableSamplesBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	sdtp := syncSampleBox{fullBox: fullBox{version: version, flags: flags}}

	//	for (i=0; i < sample_count; i++){
	//	unsigned int(2) is_leading;
	//	unsigned int(2) sample_depends_on;
	//	unsigned int(2) sample_is_depended_on;
	//	unsigned int(2) sample_has_redundancy;
	//}

	return &sdtp, nil
}

func parseSampleSizeBoxes(data []byte) (box, error) {
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

	return &stsz, nil
}

func parseCompactSampleSizeBox(data []byte) (box, error) {
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

	return &stz2, nil
}
func parseSampleToChunkBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stsc := sampleToChunkBox{fullBox: fullBox{version: version, flags: flags}}
	stsc.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stsc.entryCount); i++ {
		stsc.firstChunk = append(stsc.firstChunk, buffer.ReadUInt32())
		stsc.samplesPerChunk = append(stsc.samplesPerChunk, buffer.ReadUInt32())
		stsc.sampleDescriptionIndex = append(stsc.sampleDescriptionIndex, buffer.ReadUInt32())
	}

	return &stsc, nil
}

func parseChunkOffsetBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	stco := chunkOffsetBox{fullBox: fullBox{version: version, flags: flags}}
	stco.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(stco.entryCount); i++ {
		stco.chunkOffset = append(stco.chunkOffset, buffer.ReadUInt32())
	}

	return &stco, nil
}

func parseChunkLargeOffsetBox(data []byte) (box, error) {
	buffer := utils.NewByteBuffer(data)
	version := buffer.ReadUInt8()
	flags := buffer.ReadUInt24()
	co64 := chunkLargeOffsetBox{fullBox: fullBox{version: version, flags: flags}}
	co64.entryCount = buffer.ReadUInt32()
	for i := 0; i < int(co64.entryCount); i++ {
		co64.chunkOffset = append(co64.chunkOffset, buffer.ReadUInt64())
	}

	return &co64, nil
}

func parsePaddingBitsBox(data []byte) (box, error) {
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

	return &padb, nil
}

func parseSubSampleInformationBox(data []byte) (box, error) {
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

	return &subs, nil
}
