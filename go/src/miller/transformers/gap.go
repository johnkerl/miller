package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var GapSetup = transforming.TransformerSetup{
	Verb:         "gap",
	ParseCLIFunc: transformerGapParseCLI,
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

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pGapCount := flagSet.Int64(
		"n",
		-1,
		`Print a gap every n records.`,
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Print a gap whenever values of these fields (e.g. a,b,c) changes",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerGapUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}
	if *pGapCount == -1 && *pGroupByFieldNames == "" {
		transformerGapUsage(os.Stderr, args[0], verb, flagSet)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerGap(
		*pGapCount,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerGapUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o, "Emits an empty record every n records, or when certain values change.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
	fmt.Fprint(o, "One of -f or -g is required.\n")
	fmt.Fprint(o, "-n is ignored if -g is present.\n")
}

// ----------------------------------------------------------------
type TransformerGap struct {
	// input
	gapCount             int64
	groupByFieldNameList []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
	recordCount           int64
	previousGroupingKey   string
}

func NewTransformerGap(
	gapCount int64,
	groupByFieldNames string,
) (*TransformerGap, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &TransformerGap{
		gapCount:             gapCount,
		groupByFieldNameList: groupByFieldNameList,

		recordCount:         0,
		previousGroupingKey: "",
	}

	if len(groupByFieldNameList) == 0 {
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
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

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
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
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
