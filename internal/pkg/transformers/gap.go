package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameGap = "gap"

var GapSetup = TransformerSetup{
	Verb:         verbNameGap,
	UsageFunc:    transformerGapUsage,
	ParseCLIFunc: transformerGapParseCLI,
	IgnoresInput: false,
}

func transformerGapUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameGap)
	fmt.Fprint(o, "Emits an empty record every n records, or when certain values change.\n")
	fmt.Fprintf(o, "Options:\n")

	fmt.Fprintf(o, "Emits an empty record every n records, or when certain values change.\n")
	fmt.Fprintf(o, "-g {a,b,c} Print a gap whenever values of these fields (e.g. a,b,c) changes.\n")
	fmt.Fprintf(o, "-n {n} Print a gap every n records.\n")
	fmt.Fprintf(o, "One of -f or -g is required.\n")
	fmt.Fprintf(o, "-n is ignored if -g is present.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerGapParseCLI(
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

	gapCount := -1
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
			transformerGapUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			gapCount = cli.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerGapUsage(os.Stderr, true, 1)
		}
	}

	if gapCount == -1 && groupByFieldNames == nil {
		transformerGapUsage(os.Stderr, true, 1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerGap(
		gapCount,
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerGap struct {
	// input
	gapCount          int
	groupByFieldNames []string

	// state
	recordTransformerFunc RecordTransformerFunc
	recordCount           int
	previousGroupingKey   string
}

func NewTransformerGap(
	gapCount int,
	groupByFieldNames []string,
) (*TransformerGap, error) {

	tr := &TransformerGap{
		gapCount:          gapCount,
		groupByFieldNames: groupByFieldNames,

		recordCount:         0,
		previousGroupingKey: "",
	}

	if groupByFieldNames == nil {
		tr.recordTransformerFunc = tr.transformUnkeyed
	} else {
		tr.recordTransformerFunc = tr.transformKeyed
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerGap) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerGap) transformUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		if tr.recordCount > 0 && tr.recordCount%tr.gapCount == 0 {
			newrec := mlrval.NewMlrmapAsRecord()
			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}
		outputRecordsAndContexts.PushBack(inrecAndContext)

		tr.recordCount++

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

func (tr *TransformerGap) transformKeyed(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			groupingKey = ""
		}

		if groupingKey != tr.previousGroupingKey && tr.recordCount > 0 {
			newrec := mlrval.NewMlrmapAsRecord()
			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)

		tr.previousGroupingKey = groupingKey
		tr.recordCount++

	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}
