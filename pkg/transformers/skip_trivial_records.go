package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameSkipTrivialRecords = "skip-trivial-records"

var SkipTrivialRecordsSetup = TransformerSetup{
	Verb:         verbNameSkipTrivialRecords,
	UsageFunc:    transformerSkipTrivialRecordsUsage,
	ParseCLIFunc: transformerSkipTrivialRecordsParseCLI,
	IgnoresInput: false,
}

func transformerSkipTrivialRecordsUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSkipTrivialRecords)
	fmt.Fprintf(o, "Passes through all records except those with zero fields,\n")
	fmt.Fprintf(o, "or those for which all fields have empty value.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerSkipTrivialRecordsParseCLI(
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
			transformerSkipTrivialRecordsUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		}
		return nil, cli.VerbErrorf(verbNameSkipTrivialRecords, "option \"%s\" not recognized", opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSkipTrivialRecords()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerSkipTrivialRecords struct {
}

func NewTransformerSkipTrivialRecords() (*TransformerSkipTrivialRecords, error) {
	tr := &TransformerSkipTrivialRecords{}
	return tr, nil
}

func (tr *TransformerSkipTrivialRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		hasAny := false
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Value.String() != "" {
				hasAny = true
				break
			}
		}

		if hasAny {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		}

	} else {
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
	}
}
