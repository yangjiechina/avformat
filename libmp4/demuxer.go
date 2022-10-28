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
		name, n := r.next(size)
		println(name)
		parse, ok := parsers[name]
		if !ok {
			println("unKnow box type " + name)
		} else {
			b, consume, err := parse(data[r.offset-n : r.offset])
			if err != nil {
				panic(err)
			}

			parent.addChild(b)
			if b.hasContainer() {
				d.recursive(b, data[r.offset-n+int64(consume):r.offset])
			}
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
	end, err := d.recursive(&root, all)
	if end {
		//解析子box
	}

}
