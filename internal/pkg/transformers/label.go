package transformers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameLabel = "label"

var LabelSetup = TransformerSetup{
	Verb:         verbNameLabel,
	UsageFunc:    transformerLabelUsage,
	ParseCLIFunc: transformerLabelParseCLI,
	IgnoresInput: false,
}

func transformerLabelUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {new1,new2,new3,...}\n", "mlr", verbNameLabel)
	fmt.Fprintf(o, "Given n comma-separated names, renames the first n fields of each record to\n")
	fmt.Fprintf(o, "have the respective name. (Fields past the nth are left with their original\n")
	fmt.Fprintf(o, "names.) Particularly useful with --inidx or --implicit-csv-header, to give\n")
	fmt.Fprintf(o, "useful names to otherwise integer-indexed fields.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerLabelParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerLabelUsage(os.Stdout, true, 0)

		} else {
			transformerLabelUsage(os.Stderr, true, 1)
		}
	}

	// Get the label field names from the command line
	if argi >= argc {
		transformerLabelUsage(os.Stderr, true, 1)
	}
	newNames := lib.SplitString(args[argi], ",")
	argi++

	transformer, err := NewTransformerLabel(
		newNames,
	)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
		// TODO: return nil to caller and have it exit, maybe
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerLabel struct {
	newNames []string
}

func NewTransformerLabel(
	newNames []string,
) (*TransformerLabel, error) {
	// TODO: make this a library function.
	uniquenessChecker := make(map[string]bool)
	for _, newName := range newNames {
		_, ok := uniquenessChecker[newName]
		if ok {
			return nil, errors.New(
				fmt.Sprintf(
					"mlr label: labels must be unique; got duplicate \"%s\"\n",
					newName,
				),
			)
		}
		uniquenessChecker[newName] = true
	}

	tr := &TransformerLabel{
		newNames: newNames,
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerLabel) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.Label(tr.newNames)
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}
