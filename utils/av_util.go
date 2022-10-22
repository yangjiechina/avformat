package utils

type AVMediaType int

const (
	AVMediaTypeUnknown    = AVMediaType(-1) ///< Usually treated as AVMediaTypeData
	AVMediaTypeVideo      = AVMediaType(0)
	AVMediaTypeAudio      = AVMediaType(1)
	AVMediaTypeData       = AVMediaType(2) ///< Opaque data information usually continuous
	AVMediaTypeSubtitle   = AVMediaType(3)
	AVMediaTypeAttachment = AVMediaType(4) ///< Opaque data information usually sparse
	AVMediaTypeN          = AVMediaType(5)
)
