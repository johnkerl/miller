package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameTac = "tac"

var TacSetup = TransformerSetup{
	Verb:         verbNameTac,
	UsageFunc:    transformerTacUsage,
	ParseCLIFunc: transformerTacParseCLI,
	IgnoresInput: false,
}

func transformerTacUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameTac)
	fmt.Fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerTacParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

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
			transformerTacUsage(os.Stdout, true, 0)

		} else {
			transformerTacUsage(os.Stderr, true, 1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerTac()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerTac struct {
	recordsAndContexts *list.List
}

func NewTransformerTac() (*TransformerTac, error) {
	return &TransformerTac{
		recordsAndContexts: list.New(),
	}, nil
}

func (tr *TransformerTac) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.recordsAndContexts.PushFront(inrecAndContext)
	} else {
		// end of stream
		for e := tr.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outputRecordsAndContexts.PushBack(e.Value.(*types.RecordAndContext))
		}
		outputRecordsAndContexts.PushBack(types.NewEndOfStreamMarker(&inrecAndContext.Context))
	}
}
