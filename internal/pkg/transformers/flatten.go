package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameFlatten = "flatten"

var FlattenSetup = TransformerSetup{
	Verb:         verbNameFlatten,
	UsageFunc:    transformerFlattenUsage,
	ParseCLIFunc: transformerFlattenParseCLI,
	IgnoresInput: false,
}

func transformerFlattenUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameFlatten)
	fmt.Fprint(o,
		`Flattens multi-level maps to single-level ones. Example: field with name 'a'
and value '{"b": { "c": 4 }}' becomes name 'a.b.c' and value 4.
`)
	fmt.Fprint(o, "Options:\n")
	fmt.Fprint(o, "-f Comma-separated list of field names to flatten (default all).\n")
	fmt.Fprintf(o, "-s Separator, defaulting to %s --flatsep value.\n", "mlr")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerFlattenParseCLI(
	pargi *int,
	argc int,
	args []string,
	options *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	oFlatSep := "" // means take it from the record context
	var fieldNames []string = nil

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
			transformerFlattenUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-s" {
			oFlatSep = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerFlattenUsage(os.Stderr)
			os.Exit(1)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerFlatten(
		oFlatSep,
		options,
		fieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerFlatten struct {
	// input
	oFlatSep     string
	options      *cli.TOptions
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerFlatten(
	oFlatSep string,
	options *cli.TOptions,
	fieldNames []string,
) (*TransformerFlatten, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
	}

	retval := &TransformerFlatten{
		oFlatSep:     oFlatSep,
		options:      options,
		fieldNameSet: fieldNameSet,
	}

	retval.recordTransformerFunc = retval.flattenAll
	if fieldNameSet != nil {
		retval.recordTransformerFunc = retval.flattenSome
	}

	return retval, nil
}

// ----------------------------------------------------------------

func (tr *TransformerFlatten) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerFlatten) flattenAll(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := tr.oFlatSep
		if oFlatSep == "" {
			oFlatSep = tr.options.WriterOptions.FLATSEP
		}
		inrec.Flatten(oFlatSep)
		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerFlatten) flattenSome(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := tr.oFlatSep
		if oFlatSep == "" {
			oFlatSep = tr.options.WriterOptions.FLATSEP
		}
		inrec.FlattenFields(tr.fieldNameSet, oFlatSep)
		outputRecordsAndContexts.PushBack(inrecAndContext)
	} else {
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
