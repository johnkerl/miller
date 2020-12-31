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
var FlattenSetup = transforming.TransformerSetup{
	Verb:         "flatten",
	ParseCLIFunc: transformerFlattenParseCLI,
	IgnoresInput: false,
}

func transformerFlattenParseCLI(
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

	pOFlatSep := flagSet.String(
		"s",
		"",
		"Separator, defaulting to mlr --jflatsep value",
	)

	pFieldNames := flagSet.String(
		"f",
		"",
		"Comma-separated list of field names to flatten (default all).",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerFlattenUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerFlatten(
		*pOFlatSep,
		*pFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerFlattenUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Flattens multi-level maps to single-level ones. Example: field with name 'a'
and value '{"b": { "c": 4 }}' becomes name 'a:b:c' and value 4.
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerFlatten struct {
	// input
	oFlatSep     string
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerFlatten(
	oFlatSep string,
	fieldNames string,
) (*TransformerFlatten, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != "" {
		fieldNameSet = lib.StringListToSet(
			lib.SplitString(fieldNames, ","),
		)
	}

	retval := &TransformerFlatten{
		oFlatSep:     oFlatSep,
		fieldNameSet: fieldNameSet,
	}

	retval.recordTransformerFunc = retval.flattenAll
	if fieldNameSet != nil {
		retval.recordTransformerFunc = retval.flattenSome
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (this *TransformerFlatten) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerFlatten) flattenAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.Flatten(oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerFlatten) flattenSome(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.FlattenFields(this.fieldNameSet, oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
