package mapping

import (
	// System:
	"container/list"
	// Miller:
	"containers"
)

type MapperTac struct {
	lrecs *list.List
}

func NewMapperTac() *MapperTac {
	return &MapperTac {
		list.New(),
	}
}

func (this *MapperTac) Name() string {
	return "tac"
}

func (this *MapperTac) Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec) {
	if inrec != nil {
		this.lrecs.PushFront(inrec)
	} else {
		// end of stream
		for e := this.lrecs.Front(); e != nil; e = e.Next() {
			outrecs <- e.Value.(*containers.Lrec)
		}
		outrecs <- nil
	}
}
