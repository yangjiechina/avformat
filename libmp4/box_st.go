package libmp4

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
Box	Type:	 ‘stbl’
Container:	 Media	Information	Box	(‘minf’)
Mandatory:	Yes
Quantity:	 Exactly	one
*/
type sampleTableBox struct {
}

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
	entryCount uint32
	//unsigned int(32) entry_count;
	//int i;
	//for (i=0; i < entry_count; i++) {
	//unsigned int(32) sample_count;
	//unsigned int(32) sample_delta;
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

//Box	Type:	 ‘saiz’
//Box	Type:	 ‘saio’
//Box	Type:	 ‘sbgp’
//Box	Type:	 ‘sgpd’
