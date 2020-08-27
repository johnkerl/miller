package mapping

import (
	"containers"
)

type MapperCat struct {
	// stateless
}

func NewMapperCat() *MapperCat {
	return &MapperCat {
	}
}

func (this *MapperCat) Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec) {
	outrecs <- inrec
}
