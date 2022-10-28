package libmp4

type parser func(data []byte) (box, int, error)

const (
	containersBoxConsumeCount = 0
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
		"stsd": parseSampleDescriptionBox,
		"stts": parseDecodingTimeToSampleBox,
		"ctts": parseCompositionTimeToSampleBox,
		"cslg": parseCompositionToDecodeBox,
		"stsc": parseSampleToChunkBox,
		"stsz": parseSampleSizeBoxes,
		"stz2": parseCompactSampleSizeBox,
		"stco": parseChunkOffsetBox,
		"co64": parseChunkLargeOffsetBox,
		"stss": parseSyncSampleBox,
		"stsh": parseShadowSyncSampleBox,
		"padb": parsePaddingBitsBox,
		"stdp": parseDegradationPriorityBox,
		"sdtp": parseIndependentAndDisposableSamplesBox,
		"sbgp": parseSampleToGroupBox,
		"sgpd": parseSampleGroupDescriptionBox,
		"subs": parseSubSampleInformationBox,
		"saiz": parseSampleAuxiliaryInformationSizesBox,
		"saio": parseSampleAuxiliaryInformationOffsetsBox,
	}
}
