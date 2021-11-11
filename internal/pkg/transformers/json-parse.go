package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameJSONParse = "json-parse"

var JSONParseSetup = TransformerSetup{
	Verb:         verbNameJSONParse,
	UsageFunc:    transformerJSONParseUsage,
	ParseCLIFunc: transformerJSONParseParseCLI,
	IgnoresInput: false,
}

func transformerJSONParseUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameJSONParse)
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

func transformerJSONParseParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

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
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerJSONParseUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerJSONParse(
		fieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerJSONParse struct {
	// input
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc RecordTransformerFunc
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

func (tr *TransformerJSONParse) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerJSONParse) jsonParseAll(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
func (tr *TransformerJSONParse) jsonParseSome(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.fieldNameSet[pe.Key] {
				pe.JSONParseInPlace()
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
