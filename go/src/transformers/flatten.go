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
const verbNameFlatten = "flatten"

var FlattenSetup = transforming.TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameFlatten)
	fmt.Fprint(o,
		`Flattens multi-level maps to single-level ones. Example: field with name 'a'
and value '{"b": { "c": 4 }}' becomes name 'a.b.c' and value 4.
`)
	fmt.Fprint(o, "Options:\n")
	fmt.Fprint(o, "-f Comma-separated list of field names to flatten (default all).\n")
	fmt.Fprintf(o, "-s Separator, defaulting to %s --oflatsep value.\n", lib.MlrExeName())
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerFlattenParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TOptions,
) transforming.IRecordTransformer {

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
			oFlatSep = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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
	recordTransformerFunc transforming.RecordTransformerFunc
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
func (this *TransformerFlatten) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerFlatten) flattenAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.Flatten(oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerFlatten) flattenSome(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.FlattenFields(this.fieldNameSet, oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
