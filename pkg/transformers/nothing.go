package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameNothing = "nothing"

var NothingSetup = TransformerSetup{
	Verb:         verbNameNothing,
	ParseCLIFunc: transformerNothingParseCLI,
	UsageFunc:    transformerNothingUsage,
	IgnoresInput: false,
}

func transformerNothingUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameNothing)
	fmt.Fprintf(o, "Drops all input records. Useful for testing, or after tee/print/etc. have\n")
	fmt.Fprintf(o, "produced other output.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerNothingParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

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
			transformerNothingUsage(os.Stdout)
			return nil, cli.ErrHelpRequested
		}
		transformerNothingUsage(os.Stderr)
		return nil, fmt.Errorf("%s %s: option \"%s\" not recognized", "mlr", verbNameNothing, opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerNothing()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerNothing struct {
	// stateless
}

func NewTransformerNothing() (*TransformerNothing, error) {
	return &TransformerNothing{}, nil
}

func (tr *TransformerNothing) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if inrecAndContext.EndOfStream {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
