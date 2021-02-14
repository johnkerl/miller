package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameJSONParse = "json-parse"

var JSONParseSetup = transforming.TransformerSetup{
	Verb:         verbNameJSONParse,
	ParseCLIFunc: transformerJSONParseParseCLI,
	UsageFunc:    transformerJSONParseUsage,
	IgnoresInput: false,
}

func transformerJSONParseParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var fieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerJSONParseUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerJSONParseUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerJSONParse(
		fieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerJSONParseUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameJSONParse)
	fmt.Fprintln(
		o,
		`Tries to convert string field values to parsed JSON, e.g. "[1,2,3]" -> [1,2,3].`,
	)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {...} Comma-separated list of field names to json-parse (default all).\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerJSONParse struct {
	// input
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerJSONParse(
	fieldNames []string,
) (*TransformerJSONParse, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
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
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
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
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if this.fieldNameSet[pe.Key] {
				pe.JSONParseInPlace()
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
