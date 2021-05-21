package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameCat = "cat"

var CatSetup = transforming.TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCat)
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
	_ *cliutil.TOptions,
) transforming.IRecordTransformer {

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
			counterFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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

	recordTransformerFunc transforming.RecordTransformerFunc
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

	this := &TransformerCat{
		doCounters:        doCounters,
		groupByFieldNames: groupByFieldNames,
		counter:           0,
		countsByGroup:     make(map[string]int),
		counterFieldName:  counterFieldName,
	}

	if !doCounters {
		this.recordTransformerFunc = this.simpleCat
	} else {
		if groupByFieldNames == nil {
			this.recordTransformerFunc = this.countersUngrouped
		} else {
			this.recordTransformerFunc = this.countersGrouped
		}
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerCat) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerCat) simpleCat(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (this *TransformerCat) countersUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		this.counter++
		key := this.counterFieldName
		value := types.MlrvalFromInt(this.counter)
		inrec.PrependCopy(key, &value)
	}
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (this *TransformerCat) countersGrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		var counter int = 0
		if !ok {
			// Treat as unkeyed
			this.counter++
			counter = this.counter
		} else {
			counter, ok = this.countsByGroup[groupingKey]
			if ok {
				counter++
			} else {
				counter = 1
			}
			this.countsByGroup[groupingKey] = counter
		}

		key := this.counterFieldName
		value := types.MlrvalFromInt(counter)
		inrec.PrependCopy(key, &value)
	}
	outputChannel <- inrecAndContext
}
