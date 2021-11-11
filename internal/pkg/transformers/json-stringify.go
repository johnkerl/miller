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
const verbNameJSONStringify = "json-stringify"

var JSONStringifySetup = TransformerSetup{
	Verb:         verbNameJSONStringify,
	UsageFunc:    transformerJSONStringifyUsage,
	ParseCLIFunc: transformerJSONStringifyParseCLI,
	IgnoresInput: false,
}

func transformerJSONStringifyUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameJSONStringify)
	fmt.Fprint(o,
		`Produces string field values from field-value data, e.g. [1,2,3] -> "[1,2,3]".
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {...} Comma-separated list of field names to json-parse (default all).\n")
	fmt.Fprintf(o, "--jvstack Produce multi-line JSON output.\n")
	fmt.Fprintf(o, "--no-jvstack Produce single-line JSON output per record (default).\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerJSONStringifyParseCLI(
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
	jvStack := false // TODO: ??

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerJSONStringifyUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--jvstack" {
			jvStack = true

		} else if opt == "--no-jvstack" {
			jvStack = false

		} else {
			transformerJSONStringifyUsage(os.Stderr, true, 1)
		}
	}

	var jsonFormatting types.TJSONFormatting = types.JSON_SINGLE_LINE
	if jvStack {
		jsonFormatting = types.JSON_MULTILINE
	} else {
		jsonFormatting = types.JSON_SINGLE_LINE
	}

	transformer, err := NewTransformerJSONStringify(
		jsonFormatting,
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
type TransformerJSONStringify struct {
	// input
	jsonFormatting types.TJSONFormatting
	fieldNameSet   map[string]bool

	// state
	recordTransformerFunc RecordTransformerFunc
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

func (tr *TransformerJSONStringify) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerJSONStringify) jsonStringifyAll(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			pe.JSONStringifyInPlace(tr.jsonFormatting)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerJSONStringify) jsonStringifySome(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if tr.fieldNameSet[pe.Key] {
				pe.JSONStringifyInPlace(tr.jsonFormatting)
			}
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
