package libmp4

import (
	"io/ioutil"
)

type DeMuxer struct {
}

func (d *DeMuxer) recursive(parent box, data []byte) (bool, error) {
	r := newReader(data)
	var size int64
	for size = r.nextSize(); size > 0; {
		t, _ := r.next(size)
		println(t)
		parse, ok := parsers[t]
		if !ok {
			println("unKnow box type " + t)
		}
		b, consume, err := parse(data[r.offset-int(size) : r.offset])
		if err != nil {
			panic(err)
		}

		parent.addChild(b)
		if b.hasContainer() {
			d.recursive(b, data[r.offset-int(size)+consume:r.offset])
		}
		size = r.nextSize()
	}

	//Not the last box. need more...
	if size != 0 {
		return false, nil
	}

	return true, nil
}

func (d *DeMuxer) Read(path string) {
	all, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	root := file{}
	recursive, err := d.recursive(&root, all)
	if recursive {
		//解析子box
	}

}
