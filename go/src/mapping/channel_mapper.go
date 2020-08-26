package mapping

import (
	"containers"
)

func ChannelMapper(
	inrecs <-chan *containers.Lrec,
	outrecs chan<- *containers.Lrec,
) {
	for {
		lrec := <-inrecs
		if lrec == nil {
			outrecs <- nil
			break
		}
		MapperFoo(lrec, outrecs)
	}
}
