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
) (RecordTransformer, error) {

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
			return nil, cli.ErrHelpRequested

		} else if opt == "-n" {
			n, err := cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			tailCount = n

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
		groupByFieldNames,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerTail struct {
	// input
	tailCount         int64
	groupByFieldNames []string

	// state
	// map from string to record slices
	recordListsByGroup *lib.OrderedMap[*[]*types.RecordAndContext]
}

func NewTransformerTail(
	tailCount int64,
	groupByFieldNames []string,
) (*TransformerTail, error) {

	tr := &TransformerTail{
		tailCount:         tailCount,
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap[*[]*types.RecordAndContext](),
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
		for int64(len(*recordListForGroup)) > tr.tailCount {
			*recordListForGroup = (*recordListForGroup)[1:]
		}

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, *outer.Value...)
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
