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
const verbNameJSONStringify = "json-stringify"

var JSONStringifySetup = transforming.TransformerSetup{
	Verb:         verbNameJSONStringify,
	ParseCLIFunc: transformerJSONStringifyParseCLI,
	IgnoresInput: false,
}

func transformerJSONStringifyParseCLI(
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

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pFieldNames := flagSet.String(
		"f",
		"",
		"Comma-separated list of field names to json-stringify (default all).",
	)
	pJVStack := flagSet.Bool(
		"jvstack",
		false,
		"Produce multi-line output",
	)
	pNoJVStack := flagSet.Bool(
		"no-jvstack",
		false,
		"Produce single-line output",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerJSONStringifyUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if *pJVStack {
		jsonFormatting = types.JSON_MULTILINE
	}
	if *pNoJVStack {
		jsonFormatting = types.JSON_SINGLE_LINE
	}

	transformer, _ := NewTransformerJSONStringify(
		jsonFormatting,
		*pFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerJSONStringifyUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Produces string field values from field-value data, e.g. [1,2,3] -> "[1,2,3]".
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerJSONStringify struct {
	// input
	jsonFormatting types.TJSONFormatting
	fieldNameSet   map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerJSONStringify(
	jsonFormatting types.TJSONFormatting,
	fieldNames string,
) (*TransformerJSONStringify, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != "" {
		fieldNameSet = lib.StringListToSet(
			lib.SplitString(fieldNames, ","),
		)
	}

	retval := &TransformerJSONStringify{
		jsonFormatting: jsonFormatting,
		fieldNameSet:   fieldNameSet,
	}

	retval.recordTransformerFunc = retval.jsonStringifyAll
	if fieldNameSet != nil {
		retval.recordTransformerFunc = retval.jsonStringifySome
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (this *TransformerJSONStringify) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerJSONStringify) jsonStringifyAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			pe.JSONStringifyInPlace(this.jsonFormatting)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerJSONStringify) jsonStringifySome(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if this.fieldNameSet[pe.Key] {
				pe.JSONStringifyInPlace(this.jsonFormatting)
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
