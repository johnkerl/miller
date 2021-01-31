package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/cliutil"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameUnflatten = "unflatten"

var UnflattenSetup = transforming.TransformerSetup{
	Verb:         verbNameUnflatten,
	ParseCLIFunc: transformerUnflattenParseCLI,
	UsageFunc:    transformerUnflattenUsage,
	IgnoresInput: false,
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

	iFlatSep := "" // means take it from the record context
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
			iFlatSep = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerUnflattenUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerUnflatten(
		iFlatSep,
		fieldNames,
	)

	*pargi = argi
	return transformer
}

func transformerUnflattenUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameUnflatten)
	fmt.Fprint(o,
		`Reverses flatten. Example: field with name 'a:b:c' and value 4
becomes name 'a' and value '{"b": { "c": 4 }}'.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c} Comma-separated list of field names to unflatten (default all).\n")
	fmt.Fprintf(o, "-s {string} Separator, defaulting to %s --jflatsep value.\n", lib.MlrExeName())
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerUnflatten struct {
	// input
	iFlatSep     string
	fieldNameSet map[string]bool

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerUnflatten(
	iFlatSep string,
	fieldNames []string,
) (*TransformerUnflatten, error) {
	var fieldNameSet map[string]bool = nil
	if fieldNames != nil {
		fieldNameSet = lib.StringListToSet(fieldNames)
	}

	retval := &TransformerUnflatten{
		iFlatSep:     iFlatSep,
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
		iFlatSep := this.iFlatSep
		if iFlatSep == "" {
			iFlatSep = inrecAndContext.Context.IFLATSEP
		}
		inrec.Unflatten(iFlatSep)
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
		iFlatSep := this.iFlatSep
		if iFlatSep == "" {
			iFlatSep = inrecAndContext.Context.IFLATSEP
		}
		inrec.UnflattenFields(this.fieldNameSet, iFlatSep)
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
