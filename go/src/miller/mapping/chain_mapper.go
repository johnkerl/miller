package mapping

import (
	"miller/containers"
)

func ChainMapper(
	inrecsAndContexts <-chan *containers.LrecAndContext,
	recordMappers []IRecordMapper, // not *recordMapper since this is an interface
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	recordMapper := recordMappers[0] // xxx temp

	for {
		lrecAndContext := <-inrecsAndContexts

		recordMapper.Map(lrecAndContext, outrecsAndContexts)

		if lrecAndContext.Lrec == nil { // end of stream
			break
		}
	}
}
