package libmp4

import "avformat/utils"

type MetaData interface {
	MediaType() utils.AVMediaType
	CodeId() utils.AVCodecID
	setMediaType(mediaType utils.AVMediaType)
	setCodeId(id utils.AVCodecID)
}

type VideoMetaData struct {
	mediaType utils.AVMediaType
	codecId   utils.AVCodecID
	width     int
	height    int
}

func (v *VideoMetaData) MediaType() utils.AVMediaType {
	return v.mediaType
}

func (v *VideoMetaData) CodeId() utils.AVCodecID {
	return v.codecId
}

func (v *VideoMetaData) setMediaType(mediaType utils.AVMediaType) {
	v.mediaType = mediaType
}

func (v *VideoMetaData) setCodeId(id utils.AVCodecID) {
	v.codecId = id
}

func (v *VideoMetaData) Width() int {
	return v.width
}

func (v *VideoMetaData) Height() int {
	return v.height
}

type AudioMetaData struct {
	mediaType utils.AVMediaType
	codecId   utils.AVCodecID

	sampleRate   int
	sampleBit    int
	channelCount int
}

func (a *AudioMetaData) MediaType() utils.AVMediaType {
	return a.mediaType
}

func (a *AudioMetaData) CodeId() utils.AVCodecID {
	return a.codecId
}

func (a *AudioMetaData) setMediaType(mediaType utils.AVMediaType) {
	a.mediaType = mediaType
}

func (a *AudioMetaData) setCodeId(id utils.AVCodecID) {
	a.codecId = id
}

func (a *AudioMetaData) SampleRate() int {
	return a.sampleRate
}

func (a *AudioMetaData) SampleBit() int {
	return a.sampleBit
}

func (a *AudioMetaData) ChannelCount() int {
	return a.channelCount
}

type SubTitleMetaData struct {
	mediaType utils.AVMediaType
	codecId   utils.AVCodecID
}

func (s *SubTitleMetaData) MediaType() utils.AVMediaType {
	return s.mediaType
}

func (s *SubTitleMetaData) CodeId() utils.AVCodecID {
	return s.codecId
}

func (s *SubTitleMetaData) setMediaType(mediaType utils.AVMediaType) {
	s.mediaType = mediaType
}

func (s *SubTitleMetaData) setCodeId(id utils.AVCodecID) {
	s.codecId = id
}
