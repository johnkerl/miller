package mapping

import (
	"miller/containers"
)

func ChainMapper(
	inrecsAndContexts <-chan *containers.LrecAndContext,
	recordMappers []IRecordMapper, // not *recordMapper since this is an interface
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	i := 0
	n := len(recordMappers)

	intermediateChannels := make([]chan *containers.LrecAndContext, n-1)
	for i = 0; i < n-1; i++ {
		intermediateChannels[i] = make(chan *containers.LrecAndContext, 1)
	}

	// r M0 w
	// r M0 i0 M1 w
	// r M0 i0 M1 i1 M2 w
	// r M0 i0 M1 i1 M2 i2 M3 w

	for i, recordMapper := range recordMappers {
		ichan := inrecsAndContexts
		ochan := outrecsAndContexts

		if i > 0 {
			ichan = intermediateChannels[i-1]
		}
		if i < n-1 {
			ochan = intermediateChannels[i]
		}

		go runSingleMapper(
			ichan,
			recordMapper,
			ochan,
		)
	}
}

func runSingleMapper(
	inrecsAndContexts <-chan *containers.LrecAndContext,
	recordMapper IRecordMapper,
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	for {
		lrecAndContext := <-inrecsAndContexts
		recordMapper.Map(lrecAndContext, outrecsAndContexts)
		if lrecAndContext.Lrec == nil { // end of stream
			break
		}
	}
}
