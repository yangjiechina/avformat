package libflv

import (
	"encoding/binary"
	"math"
)

type writer interface {
	ToBytes(data []byte) int
}

type amf0Number float64

type amf0Boolean bool

type amf0String string

type AMF0Object struct{ amf0Writer }

type afm0ObjectProperty [2]writer

type afm0Null byte

func (a amf0Number) ToBytes(data []byte) int {
	data[0] = byte(AMF0DataTypeNumber)
	binary.BigEndian.PutUint64(data[1:], math.Float64bits(float64(a)))
	return 9
}

func (a amf0Boolean) ToBytes(data []byte) int {
	data[0] = byte(AMF0DataTypeBoolean)
	if a {
		data[1] = 1
	} else {
		data[1] = 0
	}
	return 2
}

func (a amf0String) ToBytes(data []byte) int {
	data[0] = byte(AMF0DataTypeString)
	binary.BigEndian.PutUint16(data[1:], uint16(len(a)))
	copy(data[3:], a)
	return 3 + len(a)
}

func (a afm0Null) ToBytes(data []byte) int {
	data[0] = byte(AMF0DataTypeNull)
	return 1
}

type AMF0Writer interface {
	writer
	AddNumber(float64)
	AddBoolean(bool)
	AddString(string)
	AddNull()
	AddObject(*AMF0Object)
}

func NewAMF0Writer() AMF0Writer {
	return &amf0Writer{}
}

type amf0Writer struct {
	nodes []writer
}

func (w *amf0Writer) AddNull() {
	w.nodes = append(w.nodes, afm0Null(0))
}

func (w *amf0Writer) ToBytes(data []byte) int {
	var count int
	for _, node := range w.nodes {
		count += node.ToBytes(data[count:])
	}

	return count
}

func (w *amf0Writer) AddNumber(f float64) {
	w.nodes = append(w.nodes, amf0Number(f))
}
func (w *amf0Writer) AddBoolean(b bool) {
	w.nodes = append(w.nodes, amf0Boolean(b))
}
func (w *amf0Writer) AddString(str string) {
	w.nodes = append(w.nodes, amf0String(str))
}

func (w *amf0Writer) AddObject(amf *AMF0Object) {
	w.nodes = append(w.nodes, amf)
}

func (a afm0ObjectProperty) ToBytes(data []byte) int {
	length := uint16(len(a[0].(amf0String)))
	binary.BigEndian.PutUint16(data, length)
	copy(data[2:], a[0].(amf0String))
	length += 2
	return a[1].ToBytes(data[length:]) + int(length)
}

func (w *AMF0Object) ToBytes(data []byte) int {
	data[0] = byte(AMF0DataTypeObject)
	i := 1 + w.amf0Writer.ToBytes(data[1:])
	i += 3
	data[i-3] = 0x0
	data[i-2] = 0x0
	data[i-1] = byte(AMF0DataTyeObjectEnd)
	return i
}

func (w *AMF0Object) AddStringProperty(name, value string) {
	w.nodes = append(w.nodes, afm0ObjectProperty([2]writer{amf0String(name), amf0String(value)}))
}

func (w *AMF0Object) AddBooleanProperty(name string, b bool) {
	w.nodes = append(w.nodes, afm0ObjectProperty([2]writer{amf0String(name), amf0Boolean(b)}))
}

func (w *AMF0Object) AddNumberProperty(name string, f float64) {
	w.nodes = append(w.nodes, afm0ObjectProperty([2]writer{amf0String(name), amf0Number(f)}))
}
