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
var JSONParseSetup = transforming.TransformerSetup{
	Verb:         "json-parse",
	ParseCLIFunc: transformerJSONParseParseCLI,
	IgnoresInput: false,
}

func transformerJSONParseParseCLI(
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

	pFieldNames := flagSet.String(
		"f",
		"",
		"Comma-separated list of field names to json-parse (default all).",
	)
	// TODO: single-line / multiline

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerJSONParseUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerJSONParse(
		*pFieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerJSONParseUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Tries to convert string field values to parsed JSON, e.g. "[1,2,3]" -> [1,2,3].
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type TransformerJSONParse struct {
	// input
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerJSONParse(
	fieldNames string,
) (*TransformerJSONParse, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != "" {
		fieldNameSet = lib.StringListToSet(
			lib.SplitString(fieldNames, ","),
		)
	}

	retval := &TransformerJSONParse{
		fieldNameSet: fieldNameSet,
	}

	retval.recordTransformerFunc = retval.jsonParseAll
	if fieldNameSet != nil {
		retval.recordTransformerFunc = retval.jsonParseSome
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (this *TransformerJSONParse) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerJSONParse) jsonParseAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			pe.JSONParseInPlace()
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerJSONParse) jsonParseSome(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if this.fieldNameSet[*pe.Key] {
				pe.JSONParseInPlace()
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
