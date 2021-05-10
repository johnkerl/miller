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
const verbNameUnflatten = "unflatten"

var UnflattenSetup = transforming.TransformerSetup{
	Verb:         verbNameUnflatten,
	UsageFunc:    transformerUnflattenUsage,
	ParseCLIFunc: transformerUnflattenParseCLI,
	IgnoresInput: false,
}

func transformerUnflattenUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameUnflatten)
	fmt.Fprint(o,
		`Reverses flatten. Example: field with name 'a.b.c' and value 4
becomes name 'a' and value '{"b": { "c": 4 }}'.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c} Comma-separated list of field names to unflatten (default all).\n")
	fmt.Fprintf(o, "-s {string} Separator, defaulting to %s --oflatsep value.\n", lib.MlrExeName())
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerUnflattenParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
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
			transformerUnflattenUsage(os.Stdout, true, 0)

		} else if opt == "-s" {
			oFlatSep = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerUnflattenUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerUnflatten(
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
type TransformerUnflatten struct {
	// input
	oFlatSep     string
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerUnflatten(
	oFlatSep string,
	fieldNames []string,
) (*TransformerUnflatten, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
	}

	retval := &TransformerUnflatten{
		oFlatSep:     oFlatSep,
		fieldNameSet: fieldNameSet,
	}
	retval.recordTransformerFunc = retval.unflattenAll
	if fieldNameSet != nil {
		retval.recordTransformerFunc = retval.unflattenSome
	}

	return retval, nil
}

// ----------------------------------------------------------------
func (this *TransformerUnflatten) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerUnflatten) unflattenAll(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.Unflatten(oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerUnflatten) unflattenSome(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		oFlatSep := this.oFlatSep
		if oFlatSep == "" {
			oFlatSep = inrecAndContext.Context.OFLATSEP
		}
		inrec.UnflattenFields(this.fieldNameSet, oFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
