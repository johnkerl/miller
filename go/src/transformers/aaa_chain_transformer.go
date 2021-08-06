package transformers

import (
	"mlr/src/types"
)

func ChainTransformer(
	inputChannel <-chan *types.RecordAndContext,
	recordTransformers []IRecordTransformer, // not *recordTransformer since this is an interface
	outputChannel chan<- *types.RecordAndContext,
) {
	i := 0
	n := len(recordTransformers)

	intermediateChannels := make([]chan *types.RecordAndContext, n-1)
	for i = 0; i < n-1; i++ {
		intermediateChannels[i] = make(chan *types.RecordAndContext, 1)
	}

	// r M0 w
	// r M0 i0 M1 w
	// r M0 i0 M1 i1 M2 w
	// r M0 i0 M1 i1 M2 i2 M3 w

	for i, recordTransformer := range recordTransformers {
		ichan := inputChannel
		ochan := outputChannel

		if i > 0 {
			ichan = intermediateChannels[i-1]
		}
		if i < n-1 {
			ochan = intermediateChannels[i]
		}

		go runSingleTransformer(
			i == 0,
			ichan,
			recordTransformer,
			ochan,
		)
	}
}

func runSingleTransformer(
	isFirst bool,
	inputChannel <-chan *types.RecordAndContext,
	recordTransformer IRecordTransformer,
	outputChannel chan<- *types.RecordAndContext,
) {
	for {
		recordAndContext := <-inputChannel

		// Three things can come through:
		//
		// * End-of-stream marker
		// * Non-nil records to be printed
		// * Strings to be printed from put/filter DSL print/dump/etc
		//   statements. They are handled here rather than fmt.Println directly
		//   in the put/filter handlers since we want all print statements and
		//   record-output to be in the same goroutine, for deterministic
		//   output ordering.
		//
		// The first two are passed to the transformer. The third we send along
		// the output channel without involving the record-transformer, since
		// there is no record to be transformed.

		if recordAndContext.EndOfStream == true || recordAndContext.Record != nil {
			recordTransformer.Transform(recordAndContext, outputChannel)
			// TODO: nr progress mod
		} else {
			outputChannel <- recordAndContext
		}

		if recordAndContext.EndOfStream {
			break
		}
	}
}
