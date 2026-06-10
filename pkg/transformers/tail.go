package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for tail counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Tail-count to print. Default 10.\n")
	fmt.Fprintf(o, "           A leading '+' means start at the nth record rather than print\n")
	fmt.Fprintf(o, "           the last n: e.g. -n +3 passes through all but the first 2\n")
	fmt.Fprintf(o, "           records, optionally by category.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerTailParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	tailCount := int64(10)
	fromStart := false
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
			return nil, cli.ErrHelpRequested

		} else if opt == "-n" {
			if argi < argc && strings.HasPrefix(args[argi], "+") {
				fromStart = true
			}
			n, err := cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			tailCount = n

		} else if opt == "-g" {
			names, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			groupByFieldNames = names

		} else {
			transformerTailUsage(os.Stderr)
			return nil, fmt.Errorf("%s %s: option \"%s\" not recognized", "mlr", verb, opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerTail(
		tailCount,
		fromStart,
		groupByFieldNames,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerTail struct {
	// input
	count             int64 // last-n count, or skip-count in from-start mode
	groupByFieldNames []string

	// state
	recordTransformerFunc RecordTransformerFunc
	// map from string to record slices
	recordListsByGroup *lib.OrderedMap[*[]*types.RecordAndContext]
	// for the from-start mode: per-group counts of records seen so far
	countsByGroup map[string]int64
}

func NewTransformerTail(
	tailCount int64,
	fromStart bool,
	groupByFieldNames []string,
) (*TransformerTail, error) {

	tr := &TransformerTail{
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap[*[]*types.RecordAndContext](),
		countsByGroup:      make(map[string]int64),
	}

	if fromStart {
		// '-n +N' is a 1-based record index, i.e. skip the first n-1.
		tr.count = tailCount - 1
		if tr.count < 0 {
			tr.count = 0
		}
		tr.recordTransformerFunc = tr.transformFromStart
	} else {
		if tailCount < 0 {
			tailCount = -tailCount
		}
		tr.count = tailCount
		tr.recordTransformerFunc = tr.transformLastN
	}

	return tr, nil
}

func (tr *TransformerTail) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerTail) transformLastN(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		recordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			records := []*types.RecordAndContext{}
			recordListForGroup = &records
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		*recordListForGroup = append(*recordListForGroup, inrecAndContext)
		for int64(len(*recordListForGroup)) > tr.count {
			(*recordListForGroup)[0] = nil // release the backing-array slot's reference
			*recordListForGroup = (*recordListForGroup)[1:]
		}

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, *outer.Value...)
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}

func (tr *TransformerTail) transformFromStart(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		tr.countsByGroup[groupingKey]++
		if tr.countsByGroup[groupingKey] > tr.count {
			// Emit records now that we skipped the requested number of them.
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
