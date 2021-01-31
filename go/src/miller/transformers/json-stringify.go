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
const verbNameJSONStringify = "json-stringify"

var JSONStringifySetup = transforming.TransformerSetup{
	Verb:         verbNameJSONStringify,
	ParseCLIFunc: transformerJSONStringifyParseCLI,
	UsageFunc:    transformerJSONStringifyUsage,
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

	var fieldNames []string = nil
	jvStack := false // TODO: ??

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerJSONStringifyUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-f" {
			fieldNames = clitypes.VerbGetStringArrayArgOrDie(verb, args, &argi, argc)

		} else if args[argi] == "--jvstack" {
			jvStack = true
			argi++

		} else if args[argi] == "--no-jvstack" {
			jvStack = false
			argi++

		} else {
			transformerJSONStringifyUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if jvStack {
		jsonFormatting = types.JSON_MULTILINE
	} else {
		jsonFormatting = types.JSON_SINGLE_LINE
	}

	transformer, _ := NewTransformerJSONStringify(
		jsonFormatting,
		fieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerJSONStringifyUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameJSONStringify)
	fmt.Fprint(o,
		`Produces string field values from field-value data, e.g. [1,2,3] -> "[1,2,3]".
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {...} Comma-separated list of field names to json-parse (default all).\n")
	fmt.Fprintf(o, "--jvstack Produce multi-line JSON output.\n")
	fmt.Fprintf(o, "--no-jvstack Produce single-line JSON output per record (default).\n")

	if doExit {
		os.Exit(exitCode)
	}
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
	fieldNames []string,
) (*TransformerJSONStringify, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
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
