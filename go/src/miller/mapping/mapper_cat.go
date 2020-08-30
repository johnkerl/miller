package mapping

import (
	"miller/containers"
	"miller/runtime"
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

func (this *MapperCat) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	outrecs <- inrec
}
