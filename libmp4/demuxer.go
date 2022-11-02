package libmp4

import (
	"fmt"
	"io/ioutil"
)

type DeMuxer struct {
	ctx deMuxContext
}

type deMuxContext struct {
	root   *file
	tracks []*track
}

func (d *DeMuxer) recursive(ctx *deMuxContext, parent box, data []byte) (bool, error) {
	r := newReader(data)
	var size int64
	for size = r.nextSize(); size > 0; size = r.nextSize() {
		name, n := r.next(size)
		fmt.Printf("size:%d name:%s\r\n", size, name)

		if parse, ok := parsers[name]; !ok {
			return false, fmt.Errorf("unKnow box type:%s", name)
		} else {
			b, consume, err := parse(ctx, data[r.offset-n:r.offset])
			if err != nil {
				return false, err
			}

			parent.addChild(b)
			if b.hasContainer() {
				_, e := d.recursive(ctx, b, data[r.offset-n+int64(consume):r.offset])
				if e != nil {
					return false, e
				}
			}
		}
	}

	//Not the last box. need more...
	if size != 0 {
		return false, nil
	}

	return true, nil
}

func (d *deMuxContext) findNextTrack() {
	/*for _, t := range d.tracks {

	}*/
}

func findNextSample(t *track) {
	if t.currentSample+1 >= t.sampleCount {
		return
	}

}

func buildIndex(ctx *deMuxContext) error {
	length := len(ctx.tracks)
	if length == 0 {
		return fmt.Errorf("uninvalid data")
	}

	for _, t := range ctx.tracks {
		if t.mark>>26 != 0x3F {
			return fmt.Errorf("uninvalid data")
		}

		t.sampleCount = t.stsz.sampleCount
		t.chunkCount = t.stco.entryCount
		t.sampleIndexEntries = make([]*sampleIndexEntry, t.sampleCount)
		var index uint32
		var duration uint32
		var dts int64
		addSampleIndex := func(chunkOffsetIndex, size int) {
			chunkOffset := t.stco.chunkOffset[chunkOffsetIndex]
			var sampleOffset uint32
			for n := 0; n < size; n++ {
				entry := sampleIndexEntry{}
				entry.pos = int64(chunkOffset + sampleOffset)
				entry.size = t.stsz.entrySize[index]
				if t.stss != nil {
					_, ok := t.stss.sampleNumber[index+1]
					entry.keyFrame = ok
				}

				tempIndex := index
				for i := 0; i < len(t.stts.sampleCount); i++ {
					if tempIndex < t.stts.sampleCount[i] {
						duration = t.stts.sampleDelta[i]
						break
					} else {
						tempIndex -= t.stts.sampleCount[i]
					}
				}

				entry.timestamp = dts
				dts += int64(duration)

				sampleOffset += entry.size
				t.sampleIndexEntries[index] = &entry
				index++
			}
		}

		for i := 0; i < len(t.stsc.firstChunk); i++ {
			chunk := t.stsc.firstChunk[i]
			size := t.stsc.samplesPerChunk[i]
			// All subsequent chunks size
			if i+1 == len(t.stsc.firstChunk) {
				for ; chunk <= t.chunkCount; chunk++ {
					addSampleIndex(int(chunk-1), int(size))
				}
			} else {
				nextChunk := t.stsc.firstChunk[i+1]
				for ; chunk < nextChunk; chunk++ {
					addSampleIndex(int(chunk-1), int(size))
				}
			}

		}
	}

	return nil
}

func (d *DeMuxer) Read(path string) {
	all, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	context := &deMuxContext{}
	context.root = &file{}
	end, err := d.recursive(context, context.root, all)
	if err != nil {
		panic(end)
	}

	buildIndex(context)
}
