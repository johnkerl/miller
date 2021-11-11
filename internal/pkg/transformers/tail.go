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
const verbNameTail = "tail"

var TailSetup = TransformerSetup{
	Verb:         verbNameTail,
	UsageFunc:    transformerTailUsage,
	ParseCLIFunc: transformerTailParseCLI,
	IgnoresInput: false,
}

func transformerTailUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameTail)
	fmt.Fprintln(o, "Passes through the last n records, optionally by category.")

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for head counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Head-count to print. Default 10.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerTailParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	tailCount := 10
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerTailUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			tailCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerTailUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerTail(
		tailCount,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerTail struct {
	// input
	tailCount         int
	groupByFieldNames []string

	// state
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerTail(
	tailCount int,
	groupByFieldNames []string,
) (*TransformerTail, error) {

	tr := &TransformerTail{
		tailCount:         tailCount,
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap(),
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerTail) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		irecordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if irecordListForGroup == nil { // first time
			irecordListForGroup = list.New()
			tr.recordListsByGroup.Put(groupingKey, irecordListForGroup)
		}
		recordListForGroup := irecordListForGroup.(*list.List)

		recordListForGroup.PushBack(inrecAndContext)
		for recordListForGroup.Len() > tr.tailCount {
			recordListForGroup.Remove(recordListForGroup.Front())
		}

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
