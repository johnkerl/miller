package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameSkipTrivialRecords = "skip-trivial-records"

var SkipTrivialRecordsSetup = transforming.TransformerSetup{
	Verb:         verbNameSkipTrivialRecords,
	UsageFunc:    transformerSkipTrivialRecordsUsage,
	ParseCLIFunc: transformerSkipTrivialRecordsParseCLI,
	IgnoresInput: false,
}

func transformerSkipTrivialRecordsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameSkipTrivialRecords)
	fmt.Fprintf(o, "Passes through all records except those with zero fields,\n")
	fmt.Fprintf(o, "or those for which all fields have empty value.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSkipTrivialRecordsParseCLI(
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
			transformerSkipTrivialRecordsUsage(os.Stdout, true, 0)

		} else {
			transformerSkipTrivialRecordsUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerSkipTrivialRecords()

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerSkipTrivialRecords struct {
}

func NewTransformerSkipTrivialRecords() (*TransformerSkipTrivialRecords, error) {
	this := &TransformerSkipTrivialRecords{}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerSkipTrivialRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
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
			outputChannel <- inrecAndContext
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
