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
const verbNameTail = "tail"

var TailSetup = TransformerSetup{
	Verb:         verbNameTail,
	UsageFunc:    transformerTailUsage,
	ParseCLIFunc: transformerTailParseCLI,
	IgnoresInput: false,
}

func transformerTailUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameTail)
	fmt.Fprintln(o, "Passes through the last n records, optionally by category.")

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for head counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Head-count to print. Default 10.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerTailParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	tailCount := int64(10)
	var groupByFieldNames []string = nil

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
			transformerTailUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-n" {
			tailCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

			// This is a bit of a hack. In our Getoptify routine we preprocess
			// the command line sending '-xyz' to '-x -y -z', but leaving
			// '--xyz' as-is. Also, Unix-like tools often support 'head -n4'
			// and 'tail -n4' in addition to 'head -n 4' and 'tail -n 4'.  Our
			// getoptify paradigm, combined with syntax familiar to users,
			// means we get '-n -4' here. So, take the absolute value to handle this.
			if tailCount < 0 {
				tailCount = -tailCount
			}

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerTailUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerTail(
		tailCount,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerTail struct {
	// input
	tailCount         int64
	groupByFieldNames []string

	// state
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerTail(
	tailCount int64,
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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
		for int64(recordListForGroup.Len()) > tr.tailCount {
			recordListForGroup.Remove(recordListForGroup.Front())
		}

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
