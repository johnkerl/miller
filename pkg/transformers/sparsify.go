package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSparsify = "sparsify"

var SparsifySetup = TransformerSetup{
	Verb:         verbNameSparsify,
	UsageFunc:    transformerSparsifyUsage,
	ParseCLIFunc: transformerSparsifyParseCLI,
	IgnoresInput: false,
}

func transformerSparsifyUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSparsify)
	fmt.Fprint(o,
		`Unsets fields for which the key is the empty string (or, optionally, another
specified value). Only makes sense with output format not being CSV or TSV.
`)

	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-s {filler string} What values to remove. Defaults to the empty string.\n")
	fmt.Fprintf(o, "-f {a,b,c} Specify field names to be operated on; any other fields won't be\n")
	fmt.Fprintf(o, "           modified. The default is to modify all fields.\n")
	fmt.Fprintf(o, "-h|--help  Show this message.\n")

	fmt.Fprint(o,
		`Example: if input is a=1,b=,c=3 then output is a=1,c=3.
`)
}

func transformerSparsifyParseCLI(
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

	fillerString := ""
	var specifiedFieldNames []string = nil

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
			transformerSparsifyUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-s" {
			fillerString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			specifiedFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerSparsifyUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerSparsify(
		fillerString,
		specifiedFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerSparsify struct {
	fillerString          string
	fieldNamesSet         map[string]bool
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerSparsify(
	fillerString string,
	specifiedFieldNames []string,
) (*TransformerSparsify, error) {

	tr := &TransformerSparsify{
		fillerString:  fillerString,
		fieldNamesSet: lib.StringListToSet(specifiedFieldNames),
	}
	if specifiedFieldNames == nil {
		tr.recordTransformerFunc = tr.transformAll
	} else {
		tr.recordTransformerFunc = tr.transformSome
	}

	return tr, nil
}

func (tr *TransformerSparsify) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		tr.recordTransformerFunc(
			inrecAndContext,
			outputRecordsAndContexts,
			inputDownstreamDoneChannel,
			outputDownstreamDoneChannel,
		)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

func (tr *TransformerSparsify) transformAll(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	inrec := inrecAndContext.Record
	outrec := mlrval.NewMlrmapAsRecord()

	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if pe.Value.String() != tr.fillerString {
			// Reference OK because ownership transfer
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
	outputRecordsAndContexts.PushBack(outrecAndContext)
}

// ----------------------------------------------------------------
func (tr *TransformerSparsify) transformSome(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	inrec := inrecAndContext.Record
	outrec := mlrval.NewMlrmapAsRecord()

	for pe := inrec.Head; pe != nil; pe = pe.Next {
		if tr.fieldNamesSet[pe.Key] {
			if pe.Value.String() != tr.fillerString {
				// Reference OK because ownership transfer
				outrec.PutReference(pe.Key, pe.Value)
			}
		} else {
			outrec.PutReference(pe.Key, pe.Value)
		}
	}

	outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
	outputRecordsAndContexts.PushBack(outrecAndContext)
}
