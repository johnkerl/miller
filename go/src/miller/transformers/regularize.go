package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var RegularizeSetup = transforming.TransformerSetup{
	Verb:         "regularize",
	ParseCLIFunc: transformerRegularizeParseCLI,
	IgnoresInput: false,
}

func transformerRegularizeParseCLI(
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

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerRegularizeUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerRegularize()

	*pargi = argi
	return transformer
}

func transformerRegularizeUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Outputs records sorted lexically ascending by keys.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerRegularize struct {
	// map from string to []string
	sortedToOriginal map[string][]string
}

func NewTransformerRegularize() (*TransformerRegularize, error) {
	this := &TransformerRegularize{
		make(map[string][]string),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerRegularize) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		currentFieldNames := inrec.GetKeys()
		currentSortedFieldNames := lib.SortedStrings(currentFieldNames)
		currentSortedFieldNamesJoined := strings.Join(currentSortedFieldNames, ",")
		previousSortedFieldNames := this.sortedToOriginal[currentSortedFieldNamesJoined]
		if previousSortedFieldNames == nil {
			this.sortedToOriginal[currentSortedFieldNamesJoined] = currentFieldNames
			outputChannel <- inrecAndContext
		} else {
			outrec := types.NewMlrmapAsRecord()
			for _, fieldName := range previousSortedFieldNames {
				outrec.PutReference(&fieldName, inrec.Get(&fieldName)) // inrec will be GC'ed
			}
			outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			outputChannel <- outrecAndContext
		}
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
