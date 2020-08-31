package mapping

import (
	"miller/containers"
	"miller/runtime"
)

func ChannelMapper(
	inrecs <-chan *containers.Lrec,
	context *runtime.Context,
	recordMapper IRecordMapper, // not *recordMapper since this is an interface
	outrecs chan<- *containers.Lrec,
) {
	for {
		lrec := <-inrecs

		context.UpdateForInputRecord(lrec)

		recordMapper.Map(lrec, context, outrecs)
		if lrec == nil { // end of stream
			break
		}
	}
}
