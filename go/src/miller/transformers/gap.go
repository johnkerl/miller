package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameGap = "gap"

var GapSetup = transforming.TransformerSetup{
	Verb:         verbNameGap,
	ParseCLIFunc: transformerGapParseCLI,
	UsageFunc:    transformerGapUsage,
	IgnoresInput: false,
}

func transformerGapParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	gapCount := -1
	var groupByFieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerGapUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-n" {
			gapCount = clitypes.VerbGetIntArgOrDie(verb, args, &argi, argc)

		} else if args[argi] == "-g" {
			groupByFieldNames = clitypes.VerbGetStringArrayArgOrDie(verb, args, &argi, argc)

		} else {
			transformerGapUsage(os.Stderr, true, 1)
		}
	}

	if gapCount == -1 && groupByFieldNames == nil {
		transformerGapUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerGap(
		gapCount,
		groupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerGapUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameGap)
	fmt.Fprint(o, "Emits an empty record every n records, or when certain values change.\n")
	fmt.Fprintf(o, "Options:\n")

	fmt.Fprintf(o, "Emits an empty record every n records, or when certain values change.\n")
	fmt.Fprintf(o, "-g {comma-separated field names} Print a gap whenever values of these fields (e.g. a,b,c) changes.\n")
	fmt.Fprintf(o, "-n {n} Print a gap every n records.\n")
	fmt.Fprintf(o, "One of -f or -g is required.\n")
	fmt.Fprintf(o, "-n is ignored if -g is present.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerGap struct {
	// input
	gapCount          int
	groupByFieldNames []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
	recordCount           int
	previousGroupingKey   string
}

func NewTransformerGap(
	gapCount int,
	groupByFieldNames []string,
) (*TransformerGap, error) {

	this := &TransformerGap{
		gapCount:          gapCount,
		groupByFieldNames: groupByFieldNames,

		recordCount:         0,
		previousGroupingKey: "",
	}

	if groupByFieldNames == nil {
		this.recordTransformerFunc = this.mapUnkeyed
	} else {
		this.recordTransformerFunc = this.mapKeyed
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerGap) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

func (this *TransformerGap) mapUnkeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		if this.recordCount > 0 && this.recordCount%this.gapCount == 0 {
			newrec := types.NewMlrmapAsRecord()
			outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
		}
		outputChannel <- inrecAndContext

		this.recordCount++

	} else {
		outputChannel <- inrecAndContext
	}
}

func (this *TransformerGap) mapKeyed(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok {
			groupingKey = ""
		}

		if groupingKey != this.previousGroupingKey && this.recordCount > 0 {
			newrec := types.NewMlrmapAsRecord()
			outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
		}

		outputChannel <- inrecAndContext

		this.previousGroupingKey = groupingKey
		this.recordCount++

	} else {
		outputChannel <- inrecAndContext
	}
}
