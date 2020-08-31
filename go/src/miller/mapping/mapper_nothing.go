package mapping

import (
	"miller/containers"
	"miller/runtime"
)

type MapperNothing struct {
	// stateless
}

func NewMapperNothing() (*MapperNothing, error) {
	return &MapperNothing{}, nil
}

func (this *MapperNothing) Name() string {
	return "nothing"
}

func (this *MapperNothing) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	if inrec == nil { // end of stream
		outrecs <- inrec
	}
}
