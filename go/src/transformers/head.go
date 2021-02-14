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
const verbNameHead = "head"

var HeadSetup = transforming.TransformerSetup{
	Verb:         verbNameHead,
	UsageFunc:    transformerHeadUsage,
	ParseCLIFunc: transformerHeadParseCLI,
	IgnoresInput: false,
}

func transformerHeadUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameHead)
	fmt.Fprintf(o, "Passes through the first n records, optionally by category.\n")

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for head counts, e.g. a,b,c.\n")
	fmt.Fprintf(o, "-n {n} Head-count to print. Default 10.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	// TODO: work on this, keeping in mind https://github.com/johnkerl/miller/issues/291
	//	fmt.Fprint(o,
	//		`Without -g, ceases consuming more input (i.e. is fast) when n records
	//have been read.
	//`)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerHeadParseCLI(
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

	headCount := 10
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerHeadUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			headCount = cliutil.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerHeadUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerHead(
		headCount,
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerHead struct {
	// input
	headCount         int
	groupByFieldNames []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
	unkeyedRecordCount    int
	keyedRecordCounts     map[string]int
}

func NewTransformerHead(
	headCount int,
	groupByFieldNames []string,
) (*TransformerHead, error) {

	this := &TransformerHead{
		headCount:         headCount,
		groupByFieldNames: groupByFieldNames,

		unkeyedRecordCount: 0,
		keyedRecordCounts:  make(map[string]int),
	}

	if groupByFieldNames == nil {
		this.recordTransformerFunc = this.mapUnkeyed
	} else {
		this.recordTransformerFunc = this.mapKeyed
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerHead) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

func (this *TransformerHead) mapUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.unkeyedRecordCount++
		if this.unkeyedRecordCount <= this.headCount {
			outputChannel <- inrecAndContext
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

func (this *TransformerHead) mapKeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok {
			return
		}

		count, present := this.keyedRecordCounts[groupingKey]
		if !present { // first time
			this.keyedRecordCounts[groupingKey] = 1
			count = 1
		} else {
			this.keyedRecordCounts[groupingKey] += 1
			count += 1
		}

		if count <= this.headCount {
			outputChannel <- inrecAndContext
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
