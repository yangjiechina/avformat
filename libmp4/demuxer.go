package libmp4

import (
	"fmt"
	"io/ioutil"
)

type DeMuxer struct {
	ctx DeMuxContext
}

type DeMuxContext struct {
	root   *file
	tracks []*track
}

func (d *DeMuxer) recursive(ctx *DeMuxContext, parent box, data []byte) (bool, error) {
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

func buildIndex(ctx *DeMuxContext) error {
	length := len(ctx.tracks)
	if length == 0 {
		return fmt.Errorf("uninvalid data")
	}

	for _, t := range ctx.tracks {
		if t.mark>>27 != 0x1F {
			return fmt.Errorf("uninvalid data")
		}
	}

	return nil
}

func (d *DeMuxer) Read(path string) {
	all, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	context := &DeMuxContext{}
	context.root = &file{}
	end, err := d.recursive(context, context.root, all)
	if err != nil {
		panic(end)
	}

	buildIndex(context)
}
