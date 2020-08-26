package mapping

import (
	"containers"
)

func MapperFoo(lrec *containers.Lrec, dest chan<- *containers.Lrec) {
	k := "foo"
	v := "bar"
	// To-do: put-by-value variant
	lrec.Put(&k, &v)
	//dest <- lrec
	dest <- lrec
}
