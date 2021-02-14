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
	ParseCLIFunc: transformerCatParseCLI,
	UsageFunc:    transformerCatUsage,
	IgnoresInput: false,
}

func transformerCatParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	doCounters := false
	counterFieldName := ""
	groupByFieldNames := ""

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
			groupByFieldNames = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCatUsage(os.Stderr, true, 1)
		}
	}
	transformer, _ := NewTransformerCat(
		doCounters,
		counterFieldName,
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerCatUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCat)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	fmt.Fprintf(o, "-n         Prepend field \"n\" to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-N {name}  Prepend field {name} to each record with record-counter starting at 1.\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for counters, e.g. a,b,c\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerCat struct {
	doCounters           bool
	groupByFieldNameList []string

	counter          int
	countsByGroup    map[string]int
	counterFieldName string

	recordTransformerFunc transforming.RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerCat(
	doCounters bool,
	counterFieldName string,
	groupByFieldNames string,
) (*TransformerCat, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	if counterFieldName != "" {
		doCounters = true
	}

	this := &TransformerCat{
		doCounters:           doCounters,
		groupByFieldNameList: groupByFieldNameList,
		counter:              0,
		countsByGroup:        make(map[string]int),
		counterFieldName:     counterFieldName,
	}

	if !doCounters {
		this.recordTransformerFunc = this.simpleCat
	} else {
		if groupByFieldNames == "" {
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

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
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
