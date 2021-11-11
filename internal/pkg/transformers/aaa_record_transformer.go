package transformers

import (
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// IRecordTransformer is the interface satisfied by all transformers, i.e.,
// Miller verbs. See stream.go for context on the channels, as well as for
// context on end-of-record-stream signaling.
type IRecordTransformer interface {
	Transform(
		inrecAndContext *types.RecordAndContext,
		inputDownstreamDoneChannel <-chan bool,
		outputDownstreamDoneChannel chan<- bool,
		outputChannel chan<- *types.RecordAndContext,
	)
}

type RecordTransformerFunc func(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
)

type TransformerUsageFunc func(
	ostream *os.File,
	doExit bool,
	exitCode int,
)

type TransformerParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
) IRecordTransformer

type TransformerSetup struct {
	Verb         string
	UsageFunc    TransformerUsageFunc
	ParseCLIFunc TransformerParseCLIFunc

	// For seqgen only -- all other transformers process records sourced by the
	// record-reader.  The seqgen verb, by contrast, is a record-source of its
	// own. (The seqgen verb probably should have been designed as a zero-file
	// record "reader" object, rather than a verb, alas.)
	IgnoresInput bool
}

// HandleDefaultDownstreamDone is a utility function for most verbs other than
// head, tee, and seqgen to use for passing downstream-done flags back
// upstream.  See ChainTransformer for context.
func HandleDefaultDownstreamDone(
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	select {
	case b := <-inputDownstreamDoneChannel:
		outputDownstreamDoneChannel <- b
		break
	default:
		break
	}
}
