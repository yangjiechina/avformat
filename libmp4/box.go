package libmp4

//MP4 file format versions
//Version	Release date	Standard	Description
//MP4 file format version 1	2001	ISO/IEC 14496-1:2001	MPEG-4 Part 1 (Systems), First edition
//MP4 file format version 2	2003	ISO/IEC 14496-14:2003	MPEG-4 Part 14 (MP4 file format), Second edition

const (
	BoxTypeFileType    = "ftyp" //only one.Forbidden to be included by other boxes.
	BoxTypeFree        = "free"
	BoxTypeMovie       = "moov"
	BoxTypeMovieHeader = "mvhd"
	BoxTypeUUID        = "uuid" // extension type
	BoxTypeTrack       = "track"

	//track type
	mediaHandlerTypeVideo    = "vmhd"
	mediaHandlerTypeAudio    = "smhd"
	mediaHandlerTypeHint     = "hmhd"
	mediaHandlerTypeSubTitle = "sthd"
	mediaHandlerTypeNull     = "nmhd"
)

//ISO/IEC 14496-12:2015
//Table 1 — Box types, structure, and cross-reference	(Informative)
//Box types, structure, and cross-reference (Informative)
//ftyp * 4.3 file type and compatibility
//pdin 8.1.3 progressive download information
//moov * 8.2.1 container for all the metadata
//mvhd * 8.2.2 movie header, overall declarations
//meta 8.11.1 metadata
//trak * 8.3.1 container for an individual track or stream
//tkhd * 8.3.2 track header, overall information about the track
//tref 8.3.3 track reference container
//trgr 8.3.4 track grouping indication
//edts 8.6.4 edit list container
//elst 8.6.6 an edit list
//meta 8.11.1 metadata
//mdia * 8.4 container for the media information in a track
//mdhd * 8.4.2 media header, overall information about the media
//hdlr * 8.4.3 handler, declares the media (handler) type
//	elng 8.4.6 extended language tag
//minf * 8.4.4 media information container
//vmhd 12.1.2 video media header, overall information (video
//track only)
//smhd 12.2.2 sound media header, overall information (sound
//track only)
//hmhd 12.4.2 hint media header, overall information (hint track
//only)
//sthd 12.6.2 subtitle media header, overall information (subtitle
//track only)
//nmhd 8.4.5.2 Null media header, overall information (some
//tracks only)
//dinf * 8.7.1 data information box, container
//ISO/IEC 14496-12:2015(E)
//16 ©	ISO/IEC	2015	–	All	rights	reserved
//Box types, structure, and cross-reference (Informative)
//dref * 8.7.2 data reference box, declares source(s) of media
//data in track
//stbl * 8.5.1 sample table box, container for the time/space
//map
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
//udta 8.10.1 user-data
//mvex 8.8.1 movie extends box
//mehd 8.8.2 movie extends header box
//trex * 8.8.3 track extends defaults
//leva 8.8.13 level assignment
//moof 8.8.4 movie fragment
//mfhd * 8.8.5 movie fragment header
//meta 8.11.1 metadata
//traf 8.8.6 track fragment
//tfhd * 8.8.7 track fragment header
//trun 8.8.8 track fragment run
//sbgp 8.9.2 sample-to-group
//sgpd 8.9.3 sample group description
//subs 8.7.7 sub-sample information
//saiz 8.7.8 sample auxiliary information sizes
//saio 8.7.9 sample auxiliary information offsets
//tfdt 8.8.12 track fragment decode time
//meta 8.11.1 metadata
//mfra 8.8.9 movie fragment random access
//tfra 8.8.10 track fragment random access
//mfro * 8.8.11 movie fragment random access offset
//mdat 8.2.2 media data container
//free 8.1.2 free space
//skip 8.1.2 free space
//udta 8.10.1 user-data
//cprt 8.10.2 copyright etc.
//tsel 8.10.3 track selection box
//strk 8.14.3 sub track box
//stri 8.14.4 sub track information box
//strd 8.14.5 sub track definition box
//meta 8.11.1 metadata
//hdlr * 8.4.3 handler, declares the metadata (handler) type
//	dinf 8.7.1 data information box, container
//dref 8.7.2 data reference box, declares source(s) of
//metadata items
//ISO/IEC 14496-12:2015(E)
//©	ISO/IEC	2015	–	All	rights	reserved 17
//Box types, structure, and cross-reference (Informative)
//iloc 8.11.3 item location
//ipro 8.11.5 item protection
//sinf 8.12.1 protection scheme information box
//frma 8.12.2 original format box
//schm 8.12.5 scheme type box
//schi 8.12.6 scheme information box
//iinf 8.11.6 item information
//xml 8.11.2 XML container
//bxml 8.11.2 binary XML container
//pitm 8.11.4 primary item reference
//fiin 8.13.2 file delivery item information
//paen 8.13.2 partition entry
//fire 8.13.7 file reservoir
//fpar 8.13.3 file partition
//fecr 8.13.4 FEC reservoir
//segr 8.13.5 file delivery session group
//gitn 8.13.6 group id to name
//idat 8.11.11 item data
//iref 8.11.12 item reference
//meco 8.11.7 additional metadata container
//mere 8.11.8 metabox relation
//meta 8.11.1 metadata
//styp 8.16.2 segment type
//	sidx 8.16.3 segment index
//ssix 8.16.4 subsegment index
//prft 8.16.5 producer reference time

type box interface {
	//version() byte
	//flags() int
	//Does not contain the length of box header.
	//fixedSize() int
	hasContainer() bool
	getChildren() []box
	addChild(b box)
	toBytes(dst []byte) int
}

type boxHeader struct {
	/**
	Most boxes will use the compact (32‐bit) size. Typically only the Media Data Box(es) need the 64‐bit size.
	if size is 1 then the actual size is in the field largesize; if size is 0, then this box is the last one in the file.
	*/
	boxSize uint64
	boxType uint32
}

type fullBox struct {
	version byte
	flags   uint32 //24 bits
}

type containerBox struct {
	children []box
}

func (c *containerBox) hasContainer() bool {
	return true
}

func (c *containerBox) getChildren() []box {
	return c.children
}

func (c *containerBox) addChild(b box) {
	c.children = append(c.children, b)
}

func (c *containerBox) toBytes(dst []byte) int {
	//offset := 0
	//for _, child := range c.children {
	//
	//}
	return 0
}

type finalBox struct {
}

func (c *finalBox) hasContainer() bool {
	return false
}

func (c *finalBox) getChildren() []box {
	return nil
}

func (c *finalBox) addChild(b box) {

}

func (c *finalBox) toBytes(dst []byte) int {

	return 0
}

/**
Box	Type:	 ‘mdia’
Container:	 Track	Box	(‘trak’)
Mandatory:	Yes
Quantity:	 Exactly	one*
*/
type mediaBox struct {
}

/**
Box	Type:	 ‘mdhd’
Container:	 Media	Box	(‘mdia’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type mediaHeaderBox struct {
	fullBox
	containerBox
	creationTime     uint64
	modificationTime uint64
	timescale        uint32
	duration         uint64
	pad              bool
	language         [3]byte
	preDefined       uint16

	//bit(1) pad = 0;
	//unsigned int(5)[3] language; // ISO-639-2/T language code
	//unsigned int(16) pre_defined = 0;
}

/**
Box	Type:	 ‘hdlr’
Container:	 Media	Box	(‘mdia’)	or	Meta	Box	(‘meta’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type handlerReferenceBox struct {
	fullBox
	finalBox

	preDefined  uint32
	handlerType uint32
	//const unsigned int(32)[3] reserved = 0
	name string
}

/**
Box	Type:	 ‘minf’
Container:	 Media	Box	(‘mdia’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type mediaInformationBox struct {
}

/**
Box	Types:	 	‘nmhd’
Container:	 Media	Information	Box	(‘minf’)
Mandatory:	Yes
Quantity:	 Exactly	one	specific	media	header	shall	be	present
*/
type nullMediaHeaderBox struct {
	fullBox
}

/**
Box	Type:	 ‘elng’
Container:	 Media	Box	(‘mdia’)
Mandatory:	No
Quantity:	 Zero	or	one
*/
type extendedLanguageTag struct {
	fullBox
	extendedLanguage string
}
