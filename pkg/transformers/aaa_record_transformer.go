package transformers

import (
	"container/list"
	"os"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/types"
)

// IRecordTransformer is the interface satisfied by all transformers, i.e.,
// Miller verbs. See stream.go for context on the channels, as well as for
// context on end-of-record-stream signaling.
type IRecordTransformer interface {
	Transform(
		inrecAndContext *types.RecordAndContext,
		outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
		inputDownstreamDoneChannel <-chan bool,
		outputDownstreamDoneChannel chan<- bool,
	)
}

type RecordTransformerFunc func(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
)

// Used within some verbs
type RecordTransformerHelperFunc func(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
)

type TransformerUsageFunc func(
	ostream *os.File,
)

type TransformerParseCLIFunc func(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
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
