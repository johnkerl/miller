package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameRemoveEmptyColumns = "remove-empty-columns"

var RemoveEmptyColumnsSetup = transforming.TransformerSetup{
	Verb:         verbNameRemoveEmptyColumns,
	UsageFunc:    transformerRemoveEmptyColumnsUsage,
	ParseCLIFunc: transformerRemoveEmptyColumnsParseCLI,
	IgnoresInput: false,
}

func transformerRemoveEmptyColumnsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameRemoveEmptyColumns)
	fmt.Fprintf(o, "Omits fields which are empty on every input row. Non-streaming.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerRemoveEmptyColumnsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TOptions,
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
			transformerRemoveEmptyColumnsUsage(os.Stdout, true, 0)

		} else {
			transformerRemoveEmptyColumnsUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerRemoveEmptyColumns()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
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
