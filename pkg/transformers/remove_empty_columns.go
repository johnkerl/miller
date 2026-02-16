package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

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
) (RecordTransformer, error) {

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
			return nil, cli.ErrHelpRequested

		}
		return nil, cli.VerbErrorf(verbNameRemoveEmptyColumns, "option \"%s\" not recognized", opt)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerRemoveEmptyColumns()
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerRemoveEmptyColumns struct {
	recordsAndContexts      []*types.RecordAndContext
	namesWithNonEmptyValues map[string]bool
}

func NewTransformerRemoveEmptyColumns() (*TransformerRemoveEmptyColumns, error) {
	tr := &TransformerRemoveEmptyColumns{
		recordsAndContexts:      []*types.RecordAndContext{},
		namesWithNonEmptyValues: make(map[string]bool),
	}
	return tr, nil
}

func (tr *TransformerRemoveEmptyColumns) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		tr.recordsAndContexts = append(tr.recordsAndContexts, inrecAndContext)

		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if !pe.Value.IsVoid() {
				tr.namesWithNonEmptyValues[pe.Key] = true
			}
		}

	} else { // end of record stream

		for _, outrecAndContext := range tr.recordsAndContexts {
			outrec := outrecAndContext.Record

			newrec := mlrval.NewMlrmapAsRecord()

			for pe := outrec.Head; pe != nil; pe = pe.Next {
				_, ok := tr.namesWithNonEmptyValues[pe.Key]
				if ok {
					// Transferring ownership from old record to new record; no copy needed
					newrec.PutReference(pe.Key, pe.Value)
				}
			}

			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(newrec, &outrecAndContext.Context))
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // Emit the stream-terminating null record
	}
}
