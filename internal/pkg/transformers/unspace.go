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
const verbNameUnspace = "unspace"

var UnspaceSetup = TransformerSetup{
	Verb:         verbNameUnspace,
	UsageFunc:    transformerUnspaceUsage,
	ParseCLIFunc: transformerUnspaceParseCLI,
	IgnoresInput: false,
}

func transformerUnspaceUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameUnspace)
	fmt.Fprintf(o, "Replaces spaces in record keys and/or values with _. This is helpful for PPRINT output.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {x}    Replace spaces with specified filler character.\n")
	fmt.Fprintf(o, "-k        Unspace only keys, not keys and values.\n")
	fmt.Fprintf(o, "-v        Unspace only values, not keys and values.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerUnspaceParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	filler := "_"
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
			transformerUnspaceUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-f" {
			filler = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-k" {
			which = "keys_only"

		} else if opt == "-v" {
			which = "values_only"

		} else {
			transformerUnspaceUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerUnspace(filler, which)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerUnspace struct {
	filler                string
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerUnspace(
	filler string,
	which string,
) (*TransformerUnspace, error) {
	tr := &TransformerUnspace{filler: filler}
	if which == "keys_only" {
		tr.recordTransformerFunc = tr.transformKeysOnly
	} else if which == "values_only" {
		tr.recordTransformerFunc = tr.transformValuesOnly
	} else {
		tr.recordTransformerFunc = tr.transformKeysAndValues
	}
	return tr, nil
}

func (tr *TransformerUnspace) Transform(
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

func (tr *TransformerUnspace) transformKeysOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		newkey := tr.unspace(pe.Key)
		// Reference not copy since this is ownership transer of the value from the now-abandoned inrec
		newrec.PutReference(newkey, pe.Value)
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}

func (tr *TransformerUnspace) transformValuesOnly(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		stringval, ok := pe.Value.GetStringValue()
		if ok {
			pe.Value = mlrval.FromString(tr.unspace(stringval))
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(inrec, &inrecAndContext.Context))
}

func (tr *TransformerUnspace) transformKeysAndValues(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	_ <-chan bool,
	__ chan<- bool,
) {
	inrec := inrecAndContext.Record
	newrec := mlrval.NewMlrmapAsRecord()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		newkey := tr.unspace(pe.Key)
		stringval, ok := pe.Value.GetStringValue()
		if ok {
			stringval = tr.unspace(stringval)
			newrec.PutReference(newkey, mlrval.FromString(stringval))
		} else {
			newrec.PutReference(newkey, pe.Value)
		}
	}
	outputRecordsAndContexts.PushBack(types.NewRecordAndContext(newrec, &inrecAndContext.Context))
}

func (tr *TransformerUnspace) unspace(input string) string {
	return strings.ReplaceAll(input, " ", tr.filler)
}
