package mapping

import (
	"miller/containers"
)

type MapperCat struct {
	// stateless
}

func NewMapperCat() *MapperCat {
	return &MapperCat{}
}

func (this *MapperCat) Name() string {
	return "cat"
}

func (this *MapperCat) Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec) {
	outrecs <- inrec
}
