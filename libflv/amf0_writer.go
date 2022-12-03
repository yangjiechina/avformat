package libflv

import (
	"encoding/binary"
	"math"
)

type writer interface {
	WriteTo(data []byte) int
}

type AMF0Number float64

type AMF0Boolean bool

type AMF0String string

func (a AMF0Number) WriteTo(data []byte) int {
	data[0] = byte(AMF0DataTypeNumber)
	binary.LittleEndian.PutUint64(data[1:], math.Float64bits(float64(a)))
	return 9
}

func (a AMF0Boolean) WriteTo(data []byte) int {
	data[0] = byte(AMF0DataTypeBoolean)
	if a {
		data[1] = 1
	} else {
		data[1] = 0
	}
	return 2
}

func (a AMF0String) WriteTo(data []byte) int {
	data[0] = byte(AMF0DataTypeString)
	binary.BigEndian.PutUint16(data[1:], uint16(len(a)))
	copy(data[2:], a)
	return 3 + len(a)
}

type AMF0Writer interface {
	writer
	AddNumber(float64)
	AddBoolean(bool)
	AddString(string)
	AddObject(*AMF0Object)
}

func NewAMF0Writer() AMF0Writer {
	return &amf0Writer{}
}

type amf0Writer struct {
	nodes []writer
}

func (w *amf0Writer) WriteTo(data []byte) int {
	var count int
	for _, node := range w.nodes {
		count += node.WriteTo(data[count:])
	}
	return count
}

func (w *amf0Writer) AddNumber(f float64) {
	w.nodes = append(w.nodes, AMF0Number(f))
}
func (w *amf0Writer) AddBoolean(b bool) {
	w.nodes = append(w.nodes, AMF0Boolean(b))
}
func (w *amf0Writer) AddString(str string) {
	w.nodes = append(w.nodes, AMF0String(str))
}

func (w *amf0Writer) AddObject(amf *AMF0Object) {
	w.nodes = append(w.nodes, amf)
}

type AMF0Object struct {
	amf0Writer
}

func (w *AMF0Object) WriteTo(data []byte) int {
	data[0] = byte(AMF0DataTypeObject)
	return 1 + w.amf0Writer.WriteTo(data)
}
