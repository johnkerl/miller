package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
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
			transformerRegularizeUsage(os.Stdout, true, 0)

		} else {
			transformerRegularizeUsage(os.Stderr, true, 1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerRegularize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
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
			outputRecordsAndContexts.PushBack(inrecAndContext)
		} else {
			outrec := mlrval.NewMlrmapAsRecord()
			for _, fieldName := range previousSortedFieldNames {
				outrec.PutReference(fieldName, inrec.Get(fieldName)) // inrec will be GC'ed
			}
			outrecAndContext := types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			outputRecordsAndContexts.PushBack(outrecAndContext)
		}
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
