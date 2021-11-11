package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
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
	doExit bool,
	exitCode int,
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

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerFlattenParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
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
		argi++

		if opt == "-h" || opt == "--help" {
			transformerFlattenUsage(os.Stdout, true, 0)

		} else if opt == "-s" {
			oFlatSep = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerFlattenUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerFlatten(
		oFlatSep,
		fieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerFlatten struct {
	// input
	oFlatSep     string
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerFlatten(
	oFlatSep string,
	fieldNames []string,
) (*TransformerFlatten, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
	}

	retval := &TransformerFlatten{
		oFlatSep:     oFlatSep,
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
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerFlatten) flattenAll(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := tr.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.FLATSEP
		}
		inrec.Flatten(oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerFlatten) flattenSome(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := tr.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.FLATSEP
		}
		inrec.FlattenFields(tr.fieldNameSet, oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
