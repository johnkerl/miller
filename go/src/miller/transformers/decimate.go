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
var DecimateSetup = transforming.TransformerSetup{
	Verb:         "decimate",
	ParseCLIFunc: transformerDecimateParseCLI,
	IgnoresInput: false,
}

func transformerDecimateParseCLI(
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

	pDecimateCount := flagSet.Int64(
		"n",
		10,
		"Decimation factor.",
	)

	pAtStart := flagSet.Bool(
		"b",
		false,
		"Decimate by printing first of every n.",
	)

	pAtEnd := flagSet.Bool(
		"e",
		false,
		"Decimate by printing last of every n (default).",
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional group-by-field names for decimate counts, e.g. a,b,c",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerDecimateUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerDecimate(
		*pDecimateCount,
		*pAtStart,
		*pAtEnd,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerDecimateUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes through one of every n records, optionally by category.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		if f.Name == "g" {
			fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
		} else {
			fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
		}
	})
}

// ----------------------------------------------------------------
type TransformerDecimate struct {
	decimateCount        int64
	remainderToKeep      int64
	groupByFieldNameList []string

	countsByGroup map[string]int64
}

// ----------------------------------------------------------------
func NewTransformerDecimate(
	decimateCount int64,
	atStart bool,
	atEnd bool,
	groupByFieldNames string,
) (*TransformerDecimate, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	remainderToKeep := decimateCount - 1
	if atStart && !atEnd {
		remainderToKeep = 0
	}

	this := &TransformerDecimate{
		decimateCount:        decimateCount,
		remainderToKeep:      remainderToKeep,
		groupByFieldNameList: groupByFieldNameList,
		countsByGroup:        make(map[string]int64),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerDecimate) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
		if !ok {
			return // This particular record doesn't have the specified fields; ignore
		}

		countForGroup, ok := this.countsByGroup[groupingKey]
		if !ok {
			countForGroup = 0
			this.countsByGroup[groupingKey] = countForGroup
		}

		remainder := countForGroup % this.decimateCount
		if remainder == this.remainderToKeep {
			outputChannel <- inrecAndContext
		}

		countForGroup++
		this.countsByGroup[groupingKey] = countForGroup

	} else {
		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
