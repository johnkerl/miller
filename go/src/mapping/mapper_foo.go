package mapping

import (
	"containers"
)

func MapperFoo(lrec *containers.Lrec, dest chan<- *containers.Lrec) {
	ka := "a"
	kb := "b"
	kab := "ab"
	va := lrec.Get(&ka)
	vb := lrec.Get(&kb)
	if va != nil && vb != nil {
		vab := *va + ":" + *vb
		// To-do: put-by-value variant
		lrec.Put(&kab, &vab)
	}
	dest <- lrec
}
