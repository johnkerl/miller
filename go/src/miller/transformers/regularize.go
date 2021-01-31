package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/cliutil"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameRegularize = "regularize"

var RegularizeSetup = transforming.TransformerSetup{
	Verb:         verbNameRegularize,
	ParseCLIFunc: transformerRegularizeParseCLI,
	UsageFunc:    transformerRegularizeUsage,
	IgnoresInput: false,
}

func transformerRegularizeParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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

	transformer, _ := NewTransformerRegularize()

	*pargi = argi
	return transformer
}

func transformerRegularizeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameRegularize)
	fmt.Fprint(o, "Outputs records sorted lexically ascending by keys.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerRegularize struct {
	// map from string to []string
	sortedToOriginal map[string][]string
}

func NewTransformerRegularize() (*TransformerRegularize, error) {
	this := &TransformerRegularize{
		make(map[string][]string),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerRegularize) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		currentFieldNames := inrec.GetKeys()
		currentSortedFieldNames := lib.SortedStrings(currentFieldNames)
		currentSortedFieldNamesJoined := strings.Join(currentSortedFieldNames, ",")
		previousSortedFieldNames := this.sortedToOriginal[currentSortedFieldNamesJoined]
		if previousSortedFieldNames == nil {
			this.sortedToOriginal[currentSortedFieldNamesJoined] = currentFieldNames
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
