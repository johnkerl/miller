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
		recordMapper.Map(lrec, outrecs)
		if lrec == nil { // end of stream
			break
		}
	}
}
