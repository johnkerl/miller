package mapping

import (
	"miller/containers"
	"miller/runtime"
)

type MapperFoo struct {
	// stateless
}

func NewMapperFoo() *MapperFoo {
	return &MapperFoo{}
}

func (this *MapperFoo) Name() string {
	return "foo"
}

func (this *MapperFoo) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	ka := "a"
	kb := "b"
	kab := "ab"
	va := inrec.Get(&ka)
	vb := inrec.Get(&kb)
	if va != nil && vb != nil {
		vab := *va + ":" + *vb
		// To-do: put-by-value variant
		inrec.Put(&kab, &vab)
	}
	outrecs <- inrec
}
