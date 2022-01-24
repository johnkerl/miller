package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
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
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerGroupLikeUsage(os.Stdout, true, 0)

		} else {
			transformerGroupLikeUsage(os.Stderr, true, 1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerGroupLike()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
				outputRecordsAndContexts.PushBack(inner.Value.(*types.RecordAndContext))
			}
		}
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
