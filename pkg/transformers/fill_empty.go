package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	fillString := defaultFillEmptyString
	inferType := true

	var err error
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
			return nil, cli.ErrHelpRequested

		} else if opt == "-v" {
			fillString, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else if opt == "-S" {
			inferType = false

		} else {
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerFillEmpty(fillString, inferType)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

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

func (tr *TransformerFillEmpty) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)

	} else { // end of record stream
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // emit end-of-stream marker
	}
}
