package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
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
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerRegularizeUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerRegularizeUsage(os.Stderr, true, 1)
			os.Exit(1)
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
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameRegularize)
	fmt.Fprint(o,
		`Outputs records sorted lexically ascending by keys.
`)

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
