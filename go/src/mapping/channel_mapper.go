package mapping

import (
	"containers"
)

func ChannelMapper(
	inrecs <-chan *containers.Lrec,
	recordMapper RecordMapper, // not *recordMapper since this is an interface
	outrecs chan<- *containers.Lrec,
) {
	for {
		lrec := <-inrecs
		if lrec == nil {
			outrecs <- nil
			break
		}
		recordMapper.Map(lrec, outrecs)
	}
}
