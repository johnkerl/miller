package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameUpcase = "upcase"

var UpcaseSetup = TransformerSetup{
	Verb:         verbNameUpcase,
	UsageFunc:    transformerUpcaseUsage,
	ParseCLIFunc: transformerUpcaseParseCLI,
	IgnoresInput: false,
}

func transformerUpcaseUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameUpcase)
	fmt.Fprintf(o, "Uppercases strings in record keys and/or values.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-k        Upcase only keys, not keys and values.\n")
	fmt.Fprintf(o, "-v        Upcase only values, not keys and values.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerUpcaseParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	which := "keys_and_values"

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
			transformerUpcaseUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-k" {
			which = "keys_only"

		} else if opt == "-v" {
			which = "values_only"

		} else {
			transformerUpcaseUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerUpcase(which)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerUpcase struct {
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerUpcase(
	which string,
) (*TransformerUpcase, error) {
	tr := &TransformerUpcase{}
	if which == "keys_only" {
		tr.recordTransformerFunc = tr.transformKeysOnly
	} else if which == "values_only" {
		tr.recordTransformerFunc = tr.transformValuesOnly
	} else {
		tr.recordTransformerFunc = tr.transformKeysAndValues
	}
	return tr, nil
}

func (tr *TransformerUpcase) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.recordTransformerFunc(
			inrecAndContext,
			outputRecordsAndContexts,
			inputDownstreamDoneChannel,
			outputDownstreamDoneChannel,
		)
	} else { // end of record stream
		outputRecordsAndContexts.PushBack(inrecAndContext)
	}
}

func (tr *TransformerUpcase) transformKeysOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		newkey := strings.ToUpper(pe.Key)
		// Reference not copy since this is ownership transfer of the value from the now-abandoned inrec
		newrec.PutReference(newkey, pe.Value)
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}

func (tr *TransformerUpcase) transformValuesOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		stringval, ok := pe.Value.GetStringValue()
		if ok {
			pe.Value = mlrval.FromString(strings.ToUpper(stringval))
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(inrec, &inrecAndContext.Context))
}

func (tr *TransformerUpcase) transformKeysAndValues(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		newkey := strings.ToUpper(pe.Key)
		stringval, ok := pe.Value.GetStringValue()
		if ok {
			stringval = strings.ToUpper(stringval)
			newrec.PutReference(newkey, mlrval.FromString(stringval))
		} else {
			newrec.PutReference(newkey, pe.Value)
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}
