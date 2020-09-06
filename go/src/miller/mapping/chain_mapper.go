package mapping

import (
	"miller/lib"
)

func ChainMapper(
	inputChannel <-chan *lib.RecordAndContext,
	recordMappers []IRecordMapper, // not *recordMapper since this is an interface
	outputChannel chan<- *lib.RecordAndContext,
) {
	i := 0
	n := len(recordMappers)

	intermediateChannels := make([]chan *lib.RecordAndContext, n-1)
	for i = 0; i < n-1; i++ {
		intermediateChannels[i] = make(chan *lib.RecordAndContext, 1)
	}

	// r M0 w
	// r M0 i0 M1 w
	// r M0 i0 M1 i1 M2 w
	// r M0 i0 M1 i1 M2 i2 M3 w

	for i, recordMapper := range recordMappers {
		ichan := inputChannel
		ochan := outputChannel

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
	inputChannel <-chan *lib.RecordAndContext,
	recordMapper IRecordMapper,
	outputChannel chan<- *lib.RecordAndContext,
) {
	for {
		recordAndContext := <-inputChannel
		recordMapper.Map(recordAndContext, outputChannel)
		if recordAndContext.Record == nil { // end of stream
			break
		}
	}
}
