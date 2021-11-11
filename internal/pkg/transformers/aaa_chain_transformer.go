package transformers

import (
	"fmt"
	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
	"os"
)

// ================================================================
// ChainTransformer is a refinement of Miller's high-level sketch in stream.go.
// As far as stream.go is concerned, the transformer-chain is a box which reads
// from an input-record channel (from the record-reader object) and writes to
// an output-record channel (to the record-writer object). Inside that box is a
// bit more complexity, including channels between transformers in the chain.
//
// ----------------------------------------------------------------
// Channel structure from outside the box:
// * readerInputRecordChannel "ichan" passes records from the record-reader
//   to the transformer chain.
// * writerOutputRecordChannel "ochan" passes records from the transformer chain
//   to the record-writer.
// * readerDownstreamDoneChannel "dchan" signifies back to the record-reader
//   that it can stop reading files, since all downstream consumers aren't
//   interested in any more data.
//
//   Record-reader
//     | ichan  ^ dchan
//     v        |
// +---------------+
// | Transformer 0 |
// | Transformer 1 |
// | Transformer 2 |
// | Transformer 3 |
// +---------------+
//     | ochan
//     v
//   Record-writer
//
// ----------------------------------------------------------------
// Channel structure from inside the box:
//
//   Record-reader
//     | ichan ^ dchan
//     v       |
// +---------------+
// | Transformer 0 |
// |   |       ^   |
// |   v       |   |
// | Transformer 1 |
// |   |       ^   |
// |   v       |   |
// | Transformer 2 |
// |   |       ^   |
// |   v       |   |
// | Transformer 3 |
// +---------------+
//     | ochan
//     v
//   Record-writer
//
// ----------------------------------------------------------------
// Each transformer has four channels from its point of view:
//
// | irchan  odchan |
// |   |       ^    |
// |   v       |    |
// | Transformer i  |
// |   |       ^    |
// |   v       |    |
// | orchan  idchan |
//
// * inputRecordChannel "irchan" is where it reads records from.
//   o If the transformer is the first in the chain, this is the
//     readerInputRecordChannel passed into ChainTransformer.
//   o Otherwise it's an intermediary channel from transformer i-1's output,
//     created and connected by ChainTransformer.
//
// * outputRecordChannel "orchan" is where it writes records to.
//   o If the transformer is the last in the chain, this is the
//     writerOutputRecordChannel passed into ChainTransformer.
//   o Otherwise it's an intermediary channel to transformer i+1's input,
//     created and connected by ChainTransformer.
//
// * inputDownstreamDoneChannel "idchan" is where it reads a downstream-done flag from.
//   o These are all created internally by ChainTransformer.
//   o If the transform is the last in the chain, nothing writes to its idchan.
//   o Otherwise transformer i's idchan is connected to transformer i+1's odchan,
//     so transformer i can accept a downstream-done flag.
//
// * outputDownstreamDoneChannel "odchan" is where it writes a downstream-done flag to.
//   o If the transformer is the first in the chain, this is the
//     readerDownstreamDoneChannel passed into ChainTransformer.
//   o Otherwise it's an intermediary channel connected to transformer i-1's idchan,
//     so transformer i can produce a downstream-done flag.
//
// ----------------------------------------------------------------
// Handling in practice:
//
// * Most verbs pass a downstrem-done flag on their idchan upstream to their odchan.
//   The HandleDefaultDownstreamDone function is for this purpose: most verbs use it.
//
// * The exceptional verbs are HEAD, TEE, and SEQGEN.
//
// * mlr head is the reason this exists. The problem to solve is that people
//   want to do 'mlr head -n 10 myhugefile.dat' and have mlr exit as soon as those
//   10 records have been read.
//
// * However, if someone does 'mlr cut -f foo then tee bar.dat then head -n 10',
//   they want bar.dat to have all the records seen by tee; head -n 10 should produce
//   only the first 10 but bar.dat should have them all.
//
// * Likewise, 'mlr seqgen --stop 1000000000 then head -n 10' should result
//   in seqgen breaking out of its data-production loop.
//
// * In head.go, tee.go, and seqgen.go you will see specific handling of
//   reading idchan and writing odchan.
//
// ----------------------------------------------------------------
// TESTING
//
// * This is a bit awkward with regard to regression-test -- we don't want a
//   multi-GB data file in our repo for the continuous integration job to check
//   that the processing finishes quickly.
//
// * Nonetheless: all these should finish quickly/
//
//   mlr head -n 10 ~/tmp/huge
//   mlr cat then head -n 10 ~/tmp/huge
//   mlr head -n 100 then tee foo.txt then head -n 10 ~/tmp/huge
//     check `wc -l foo.txt` is 100
//   mlr head -n 100 then tee foo.txt then head -n 10 then tee bar.txt ~/tmp/huge
//     check `wc -l foo.txt` is 100 and `wc -l bar.txt` is 10
//   mlr seqgen --stop 100000000 then head -n 10
//
// ================================================================

// ChainTransformer is a refinement of Miller's high-level sketch in stream.go.
// While stream.go sees goroutines for record reader, transformer chain, and
// record writer, with input channel from record-reader to transformer chain
// and output channel from transformer chain to record-writer, ChainTransformer
// subdivides goroutines for each transformer in the chain, with intermediary
// channels between them.
func ChainTransformer(
	readerInputRecordChannel <-chan *types.RecordAndContext,
	readerDownstreamDoneChannel chan<- bool, // for mlr head -- see also stream.go
	recordTransformers []IRecordTransformer, // not *recordTransformer since this is an interface
	writerOutputRecordChannel chan<- *types.RecordAndContext,
	options *cli.TOptions,
) {
	i := 0
	n := len(recordTransformers)

	intermediateRecordChannels := make([]chan *types.RecordAndContext, n-1)
	for i = 0; i < n-1; i++ {
		intermediateRecordChannels[i] = make(chan *types.RecordAndContext, 1)
	}

	intermediateDownstreamDoneChannels := make([]chan bool, n)
	for i = 0; i < n; i++ {
		intermediateDownstreamDoneChannels[i] = make(chan bool, 1)
	}

	for i, recordTransformer := range recordTransformers {
		// Downstream flow: channel a given transformer reads records from
		irchan := readerInputRecordChannel
		// Downstream flow: channel a given transformer writes transformed
		// records to
		orchan := writerOutputRecordChannel
		// Upstream signaling: channel a given transformer reads to see if
		// downstream transformers are done (e.g. mlr head)
		idchan := intermediateDownstreamDoneChannels[i]
		// Upstream signaling: channel a given transformer (e.g. mlr head)
		// writes to signal to upstream transformers that it will ignore
		// further input.
		odchan := readerDownstreamDoneChannel

		if i > 0 {
			irchan = intermediateRecordChannels[i-1]
			odchan = intermediateDownstreamDoneChannels[i-1]
		}
		if i < n-1 {
			orchan = intermediateRecordChannels[i]
		}

		go runSingleTransformer(
			recordTransformer,
			i == 0,
			irchan,
			orchan,
			idchan,
			odchan,
			options,
		)
	}
}

func runSingleTransformer(
	recordTransformer IRecordTransformer,
	isFirst bool,
	inputRecordChannel <-chan *types.RecordAndContext,
	outputRecordChannel chan<- *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	options *cli.TOptions,
) {

	for {
		recordAndContext := <-inputRecordChannel

		// --nr-progress-mod
		// TODO: function-pointer this away to reduce instruction count in the
		// normal case which it isn't used at all. No need to test if {static thing} != 0
		// on every record.
		if options.NRProgressMod != 0 {
			if isFirst && recordAndContext.Record != nil {
				context := &recordAndContext.Context
				if context.NR%options.NRProgressMod == 0 {
					fmt.Fprintf(os.Stderr, "NR=%d FNR=%d FILENAME=%s\n", context.NR, context.FNR, context.FILENAME)
				}
			}
		}

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
			recordTransformer.Transform(
				recordAndContext,
				inputDownstreamDoneChannel,
				outputDownstreamDoneChannel,
				outputRecordChannel,
			)
		} else {
			outputRecordChannel <- recordAndContext
		}

		if recordAndContext.EndOfStream {
			break
		}
	}
}
