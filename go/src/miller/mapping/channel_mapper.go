package mapping

import (
	"miller/containers"
	"miller/runtime"
)

func ChannelMapper(
	inrecsAndContexts <-chan *runtime.LrecAndContext,
	recordMapper IRecordMapper, // not *recordMapper since this is an interface
	outrecs chan<- *containers.Lrec,
) {
	for {
		lrecAndContext := <-inrecsAndContexts
		lrec := lrecAndContext.Lrec
		context := lrecAndContext.Context

		recordMapper.Map(lrec, &context, outrecs)
		if lrec == nil { // end of stream
			break
		}
	}
}
