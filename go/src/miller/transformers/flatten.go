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

	// flatsep: defaults to global flatsep
	// field names: default all

	pFieldNames := flagSet.String(
		"f",
		"",
		"Field names to flatten",
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
	// TODO: type this up.
	fmt.Fprint(o,
		`TO DO: type this up.
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
	fieldNameList []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerFlatten(
	fieldNames string,
) (*TransformerFlatten, error) {

	fieldNameList := lib.SplitString(fieldNames, ",")

	this := &TransformerFlatten{
		fieldNameList: fieldNameList,
	}

	if len(fieldNameList) == 0 {
		this.recordTransformerFunc = this.flattenAll
	} else {
		this.recordTransformerFunc = this.flattenSome
	}

	return this, nil
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
		// TODO
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
		for _, fieldName := range this.fieldNameList {
			inrec.MoveToTail(&fieldName)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
