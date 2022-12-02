package utils

type Packet struct {
	data ByteBuffer
	pts  int64
	dts  int64
}

func (p *Packet) Pts() int64 {
	return p.pts
}

func (p *Packet) Dts() int64 {
	return p.dts
}

func (p *Packet) SetPts(pts int64) {
	p.pts = pts
}

func (p *Packet) SetDts(dts int64) {
	p.dts = dts
}

func (p *Packet) Data() ByteBuffer {
	return p.data
}

func (p *Packet) Write(data []byte) {
	p.data.Write(data)
}
func (p *Packet) Release() {
	p.data.Clear()
	p.pts = -1
	p.dts = -1
}

func NewPacket() *Packet {
	return &Packet{
		data: NewByteBuffer(),
		pts:  -1,
		dts:  -1,
	}
}
