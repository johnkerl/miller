package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameGroupLike = "group-like"

var GroupLikeSetup = TransformerSetup{
	Verb:         verbNameGroupLike,
	UsageFunc:    transformerGroupLikeUsage,
	ParseCLIFunc: transformerGroupLikeParseCLI,
	IgnoresInput: false,
}

func transformerGroupLikeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameGroupLike)
	fmt.Fprintln(o, "Outputs records in batches having identical field names.")
	fmt.Fprintln(o, "Options:")
	fmt.Fprintln(o, "-h|--help Show this message.")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerGroupLikeParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerGroupLikeUsage(os.Stdout, true, 0)

		} else {
			transformerGroupLikeUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerGroupLike()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerGroupLike struct {
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerGroupLike() (*TransformerGroupLike, error) {

	tr := &TransformerGroupLike{
		recordListsByGroup: lib.NewOrderedMap(),
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerGroupLike) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey := inrec.GetKeysJoined()

		recordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			recordListForGroup = list.New()
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputChannel <- inner.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
