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
const verbNameReorder = "reorder"

var ReorderSetup = transforming.TransformerSetup{
	Verb:         verbNameReorder,
	UsageFunc:    transformerReorderUsage,
	ParseCLIFunc: transformerReorderParseCLI,
	IgnoresInput: false,
}

func transformerReorderUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameReorder)
	fmt.Fprint(o,
		`Moves specified names to start of record, or end of record.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-e Put specified field names at record end: default is to put them at record start.\n")
	fmt.Fprintf(o, "-f {a,b,c} Field names to reorder.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(
		o,
		"%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n",
		lib.MlrExeName(), verbNameReorder,
	)
	fmt.Fprintf(
		o,
		"%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n",
		lib.MlrExeName(), verbNameReorder,
	)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerReorderParseCLI(
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

	var fieldNames []string = nil
	putAtEnd := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerReorderUsage(os.Stdout, true, 0)

		} else if opt == "-f" {
			fieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-e" {
			putAtEnd = true

		} else {
			transformerReorderUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerReorderUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerReorder(
		fieldNames,
		putAtEnd,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerReorder struct {
	// input
	fieldNames []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerReorder(
	fieldNames []string,
	putAtEnd bool,
) (*TransformerReorder, error) {

	if !putAtEnd {
		lib.ReverseStringList(fieldNames)
	}

	this := &TransformerReorder{
		fieldNames: fieldNames,
	}

	if !putAtEnd {
		this.recordTransformerFunc = this.reorderToStart
	} else {
		this.recordTransformerFunc = this.reorderToEnd
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerReorder) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerReorder) reorderToStart(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range this.fieldNames {
			inrec.MoveToHead(fieldName)
		}
		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerReorder) reorderToEnd(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		for _, fieldName := range this.fieldNames {
			inrec.MoveToTail(fieldName)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
