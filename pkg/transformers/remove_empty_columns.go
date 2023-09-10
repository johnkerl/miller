package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/types"
)

// ----------------------------------------------------------------
const verbNameRemoveEmptyColumns = "remove-empty-columns"

var RemoveEmptyColumnsSetup = TransformerSetup{
	Verb:         verbNameRemoveEmptyColumns,
	UsageFunc:    transformerRemoveEmptyColumnsUsage,
	ParseCLIFunc: transformerRemoveEmptyColumnsParseCLI,
	IgnoresInput: false,
}

func transformerRemoveEmptyColumnsUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameRemoveEmptyColumns)
	fmt.Fprintf(o, "Omits fields which are empty on every input row. Non-streaming.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerRemoveEmptyColumnsParseCLI(
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
			transformerRemoveEmptyColumnsUsage(os.Stdout)
			os.Exit(0)

		} else {
			transformerRemoveEmptyColumnsUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerRemoveEmptyColumns()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerRemoveEmptyColumns struct {
	recordsAndContexts      *list.List
	namesWithNonEmptyValues map[string]bool
}

func NewTransformerRemoveEmptyColumns() (*TransformerRemoveEmptyColumns, error) {
	tr := &TransformerRemoveEmptyColumns{
		recordsAndContexts:      list.New(),
		namesWithNonEmptyValues: make(map[string]bool),
	}
	return tr, nil
}

// ---------------------------------------------------------------

func (tr *TransformerRemoveEmptyColumns) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		tr.recordsAndContexts.PushBack(inrecAndContext)

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !pe.Value.IsVoid() {
				tr.namesWithNonEmptyValues[pe.Key] = true
			}
		}

	} else { // end of record stream

		for e := tr.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outrecAndContext := e.Value.(*types.RecordAndContext)
			outrec := outrecAndContext.Record

			newrec := mlrval.NewMlrmapAsRecord()

			for pe := outrec.Head; pe != nil; pe = pe.Next {
				_, ok := tr.namesWithNonEmptyValues[pe.Key]
				if ok {
					// Transferring ownership from old record to new record; no copy needed
					newrec.PutReference(pe.Key, pe.Value)
				}
			}

			outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &outrecAndContext.Context))
		}

		outputRecordsAndContexts.PushBack(inrecAndContext) // Emit the stream-terminating null record
	}
}
