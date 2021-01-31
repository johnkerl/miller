package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameSkipTrivialRecords = "skip-trivial-records"

var SkipTrivialRecordsSetup = transforming.TransformerSetup{
	Verb:         verbNameSkipTrivialRecords,
	ParseCLIFunc: transformerSkipTrivialRecordsParseCLI,
	UsageFunc:    transformerSkipTrivialRecordsUsage,
	IgnoresInput: false,
}

func transformerSkipTrivialRecordsParseCLI(
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
			transformerSkipTrivialRecordsUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerSkipTrivialRecordsUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerSkipTrivialRecords()

	*pargi = argi
	return transformer
}

func transformerSkipTrivialRecordsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameSkipTrivialRecords)
	fmt.Fprintf(o, "Passes through all records except those with zero fields,\n")
	fmt.Fprintf(o, "or those for which all fields have empty value.\n")

	if doExit {
		os.Exit(exitCode)
	}
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
