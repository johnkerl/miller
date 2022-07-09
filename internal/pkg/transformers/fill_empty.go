package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameFillEmpty = "fill-empty"
const defaultFillEmptyString = "N/A"

var FillEmptySetup = TransformerSetup{
	Verb:         verbNameFillEmpty,
	UsageFunc:    transformerFillEmptyUsage,
	ParseCLIFunc: transformerFillEmptyParseCLI,
	IgnoresInput: false,
}

func transformerFillEmptyUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameFillEmpty)
	fmt.Fprintf(o, "Fills empty-string fields with specified fill-value.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-v {string} Fill-value: defaults to \"%s\"\n", defaultFillEmptyString)
	fmt.Fprintf(o, "-S          Don't infer type -- so '-v 0' would fill string 0 not int 0.\n")
}

func transformerFillEmptyParseCLI(
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

	fillString := defaultFillEmptyString
	inferType := true

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
			transformerFillEmptyUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-v" {
			fillString = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-S" {
			inferType = false

		} else {
			transformerFillEmptyUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerFillEmpty(fillString, inferType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerFillEmpty struct {
	fillValue *mlrval.Mlrval
}

func NewTransformerFillEmpty(
	fillString string,
	inferType bool,
) (*TransformerFillEmpty, error) {
	tr := &TransformerFillEmpty{}
	if inferType {
		tr.fillValue = mlrval.FromInferredType(fillString)
	} else {
		tr.fillValue = mlrval.FromString(fillString)
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerFillEmpty) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Value.IsVoid() {
				pe.Value = tr.fillValue
			}
		}

		outputRecordsAndContexts.PushBack(inrecAndContext)

	} else { // end of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext) // emit end-of-stream marker
	}
}
