package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCat = "cat"

var CatSetup = TransformerSetup{
	Verb:         verbNameCat,
	UsageFunc:    transformerCatUsage,
	ParseCLIFunc: transformerCatParseCLI,
	IgnoresInput: false,
}

func transformerCatUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCat)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-n         Prepend field \"n\" to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-N {name}  Prepend field {name} to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for counters, e.g. a,b,c\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCatParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	doCounters := false
	counterFieldName := ""
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCatUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			counterFieldName = "n"

		} else if opt == "-N" {
			counterFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCatUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerCat(
		doCounters,
		counterFieldName,
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
type TransformerCat struct {
	doCounters        bool
	groupByFieldNames []string

	counter          int
	countsByGroup    map[string]int
	counterFieldName string

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerCat(
	doCounters bool,
	counterFieldName string,
	groupByFieldNames []string,
) (*TransformerCat, error) {

	if counterFieldName != "" {
		doCounters = true
	}

	tr := &TransformerCat{
		doCounters:        doCounters,
		groupByFieldNames: groupByFieldNames,
		counter:           0,
		countsByGroup:     make(map[string]int),
		counterFieldName:  counterFieldName,
	}

	if !doCounters {
		tr.recordTransformerFunc = tr.simpleCat
	} else {
		if groupByFieldNames == nil {
			tr.recordTransformerFunc = tr.countersUngrouped
		} else {
			tr.recordTransformerFunc = tr.countersGrouped
		}
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCat) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerCat) simpleCat(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (tr *TransformerCat) countersUngrouped(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		tr.counter++
		key := tr.counterFieldName
		inrec.PrependCopy(key, types.MlrvalFromInt(tr.counter))
	}
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (tr *TransformerCat) countersGrouped(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		var counter int = 0
		if !ok {
			// Treat as unkeyed
			tr.counter++
			counter = tr.counter
		} else {
			counter, ok = tr.countsByGroup[groupingKey]
			if ok {
				counter++
			} else {
				counter = 1
			}
			tr.countsByGroup[groupingKey] = counter
		}

		key := tr.counterFieldName
		inrec.PrependCopy(key, types.MlrvalFromInt(counter))
	}
	outputChannel <- inrecAndContext
}
