package transformers

import (
	//"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var FillDownSetup = transforming.TransformerSetup{
	Verb:         "fill-down",
	ParseCLIFunc: transformerFillDownParseCLI,
	IgnoresInput: false,
}

func transformerFillDownParseCLI(
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

	pFillDownFieldNames := flagSet.String(
		"f",
		"",
		`Field names for fill-down`,
	)

	pOnlyIfAbsentShort := flagSet.Bool(
		"a",
		false,
		`If a given record has a missing value for a given field, fill that from
the corresponding value from a previous record, if any.
By default, a 'missing' field either is absent, or has the empty-string value.
With -a, a field is 'missing' only if it is absent.`,
	)

	pOnlyIfAbsentLong := flagSet.Bool(
		"only-if-absent",
		false,
		`Synonym for -a`,
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerFillDownUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerFillDown(
		*pFillDownFieldNames,
		*pOnlyIfAbsentShort || *pOnlyIfAbsentLong,
	)

	*pargi = argi
	return transformer
}

func transformerFillDownUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Passes through the last n records, optionally by category.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerFillDown struct {
	// input
	fillDownFieldNames []string
	onlyIfAbsent       bool

	// state
	lastNonNullValues map[string]*types.Mlrval
}

func NewTransformerFillDown(
	fillDownFieldNames string,
	onlyIfAbsent bool,
) (*TransformerFillDown, error) {
	this := &TransformerFillDown{
		fillDownFieldNames: lib.SplitString(fillDownFieldNames, ","),
		onlyIfAbsent:       onlyIfAbsent,
		lastNonNullValues:  make(map[string]*types.Mlrval),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerFillDown) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fillDownFieldName := range this.fillDownFieldNames {
			present := false
			value := inrec.Get(fillDownFieldName)
			if this.onlyIfAbsent {
				present = value != nil
			} else {
				present = value != nil && !value.IsEmpty()
			}

			if present {
				// Remember it for a subsequent record lacking this field
				this.lastNonNullValues[fillDownFieldName] = value.Copy()
			} else {
				// Reuse previously seen value, if any
				prev, ok := this.lastNonNullValues[fillDownFieldName]
				if ok {
					inrec.PutCopy(fillDownFieldName, prev)
				}
			}
		}

		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
