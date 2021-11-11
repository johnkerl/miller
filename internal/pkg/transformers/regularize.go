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
const verbNameRegularize = "regularize"

var RegularizeSetup = TransformerSetup{
	Verb:         verbNameRegularize,
	UsageFunc:    transformerRegularizeUsage,
	ParseCLIFunc: transformerRegularizeParseCLI,
	IgnoresInput: false,
}

func transformerRegularizeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameRegularize)
	fmt.Fprintf(o, "Outputs records sorted lexically ascending by keys.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerRegularizeParseCLI(
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
			transformerRegularizeUsage(os.Stdout, true, 0)

		} else {
			transformerRegularizeUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerRegularize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerRegularize struct {
	// map from string to []string
	sortedToOriginal map[string][]string
}

func NewTransformerRegularize() (*TransformerRegularize, error) {
	tr := &TransformerRegularize{
		make(map[string][]string),
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerRegularize) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		currentFieldNames := inrec.GetKeys()
		currentSortedFieldNames := lib.SortedStrings(currentFieldNames)
		currentSortedFieldNamesJoined := strings.Join(currentSortedFieldNames, ",")
		previousSortedFieldNames := tr.sortedToOriginal[currentSortedFieldNamesJoined]
		if previousSortedFieldNames == nil {
			tr.sortedToOriginal[currentSortedFieldNamesJoined] = currentFieldNames
			outputChannel <- inrecAndContext
		} else {
			outrec := types.NewMlrmapAsRecord()
			for _, fieldName := range previousSortedFieldNames {
				outrec.PutReference(fieldName, inrec.Get(fieldName)) // inrec will be GC'ed
			}
			outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			outputChannel <- outrecAndContext
		}
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
