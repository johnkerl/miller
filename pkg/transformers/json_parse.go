package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/types"
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
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameJSONParse)
	fmt.Fprintln(
		o,
		`Tries to convert string field values to parsed JSON, e.g. "[1,2,3]" -> [1,2,3].`,
	)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {...} Comma-separated list of field names to json-parse (default all).\n")
	fmt.Fprintf(o, "-k       If supplied, then on parse fail for any cell, keep the (unparsable)\n")
	fmt.Fprintf(o, "         input value for the cell.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerJSONParseParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++
	keepFailed := false

	var fieldNames []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerJSONParseUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-k" {
			keepFailed = true

		} else {
			transformerJSONParseUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerJSONParse(
		fieldNames,
		keepFailed,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerJSONParse struct {
	// input
	fieldNameSet map[string]bool
	keepFailed   bool

	// state
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerJSONParse(
	fieldNames []string,
	keepFailed bool,
) (*TransformerJSONParse, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
	}

	retval := &TransformerJSONParse{
		fieldNameSet: fieldNameSet,
		keepFailed:   keepFailed,
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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerJSONParse) jsonParseAll(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.keepFailed {
				pe.JSONTryParseInPlace()
			} else {
				pe.JSONParseInPlace()
			}
		}
		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerJSONParse) jsonParseSome(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.fieldNameSet[pe.Key] {
				if tr.keepFailed {
					pe.JSONTryParseInPlace()
				} else {
					pe.JSONParseInPlace()
				}
			}
		}
		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
