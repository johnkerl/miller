package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var UnflattenSetup = transforming.TransformerSetup{
	Verb:         "unflatten",
	ParseCLIFunc: transformerUnflattenParseCLI,
	IgnoresInput: false,
}

func transformerUnflattenParseCLI(
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

	pIFlatSep := flagSet.String(
		"s",
		"",
		"Separator, defaulting to mlr --jflatsep value",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerUnflattenUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerUnflatten(
		*pIFlatSep,
	)

	*pargi = argi
	return transformer
}

func transformerUnflattenUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Reverses flatten. Example: field with name 'a:b:c' and value 4
becomes name 'a' and value '{"b": { "c": 4 }}'.
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerUnflatten struct {
	// input
	iFlatSep string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerUnflatten(
	iFlatSep string,
) (*TransformerUnflatten, error) {
	return &TransformerUnflatten{
		iFlatSep: iFlatSep,
	}, nil
}

// ----------------------------------------------------------------
func (this *TransformerUnflatten) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		iFlatSep := this.iFlatSep
		if iFlatSep == "" {
			iFlatSep = inrecAndContext.Context.IFLATSEP
		}
		inrec.Unflatten(iFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
