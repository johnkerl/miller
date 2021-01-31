package transformers

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameRemoveEmptyColumns = "remove-empty-columns"

var RemoveEmptyColumnsSetup = transforming.TransformerSetup{
	Verb:         verbNameRemoveEmptyColumns,
	ParseCLIFunc: transformerRemoveEmptyColumnsParseCLI,
	UsageFunc:    transformerRemoveEmptyColumnsUsage,
	IgnoresInput: false,
}

func transformerRemoveEmptyColumnsParseCLI(
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
			transformerRemoveEmptyColumnsUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerRemoveEmptyColumnsUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerRemoveEmptyColumns()

	*pargi = argi
	return transformer
}

func transformerRemoveEmptyColumnsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameRemoveEmptyColumns)
	fmt.Fprintf(o, "Omits fields which are empty on every input row. Non-streaming.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerRemoveEmptyColumns struct {
	recordsAndContexts      *list.List
	namesWithNonEmptyValues map[string]bool
}

func NewTransformerRemoveEmptyColumns() (*TransformerRemoveEmptyColumns, error) {
	this := &TransformerRemoveEmptyColumns{
		recordsAndContexts:      list.New(),
		namesWithNonEmptyValues: make(map[string]bool),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerRemoveEmptyColumns) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		this.recordsAndContexts.PushBack(inrecAndContext)

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !pe.Value.IsEmpty() {
				this.namesWithNonEmptyValues[pe.Key] = true
			}
		}

	} else { // end of record stream

		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
			outrec := outrecAndContext.Record

			newrec := types.NewMlrmapAsRecord()

			for pe := outrec.Head; pe != nil; pe = pe.Next {
				_, ok := this.namesWithNonEmptyValues[pe.Key]
				if ok {
					// Transferring ownership from old record to new record; no copy needed
					newrec.PutReference(pe.Key, pe.Value)
				}
			}

			outputChannel <- types.NewRecordAndContext(newrec, &outrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
