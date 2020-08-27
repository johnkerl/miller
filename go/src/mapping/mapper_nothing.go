package mapping

import (
	"containers"
)

type MapperNothing struct {
	// stateless
}

func NewMapperNothing() *MapperNothing {
	return &MapperNothing {
	}
}

func (this *MapperNothing) Name() string {
	return "nothing"
}

func (this *MapperNothing) Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec) {
	if inrec == nil { // end of stream
		outrecs <- inrec
	}
}
