package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameTac = "tac"

var TacSetup = TransformerSetup{
	Verb:         verbNameTac,
	UsageFunc:    transformerTacUsage,
	ParseCLIFunc: transformerTacParseCLI,
	IgnoresInput: false,
}

func transformerTacUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameTac)
	fmt.Fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerTacParseCLI(
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
			transformerTacUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		}
		return nil, cli.VerbErrorf(verbNameTac, "option \"%s\" not recognized", opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerTac()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerTac struct {
	recordsAndContexts []*types.RecordAndContext
}

func NewTransformerTac() (*TransformerTac, error) {
	return &TransformerTac{
		recordsAndContexts: []*types.RecordAndContext{},
	}, nil
}

func (tr *TransformerTac) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.recordsAndContexts = append(tr.recordsAndContexts, inrecAndContext)
	} else {
		// end of stream
		for i := len(tr.recordsAndContexts) - 1; i >= 0; i-- {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, tr.recordsAndContexts[i])
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewEndOfStreamMarker(&inrecAndContext.Context))
	}
}
