package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameReorder = "reorder"

var ReorderSetup = transforming.TransformerSetup{
	Verb:         verbNameReorder,
	ParseCLIFunc: transformerReorderParseCLI,
	UsageFunc:    transformerReorderUsage,
	IgnoresInput: false,
}

func transformerReorderParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var fieldNames []string = nil
	putAtEnd := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerReorderUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-f" {
			fieldNames = clitypes.VerbGetStringArrayArgOrDie(verb, args, &argi, argc)

		} else if args[argi] == "-e" {
			putAtEnd = true
			argi++

		} else {
			transformerReorderUsage(os.Stderr, true, 1)
			os.Exit(1)
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

func transformerReorderUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", os.Args[0], verbNameReorder)
	fmt.Fprint(o,
		`Moves specified names to start of record, or end of record.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-e Put specified field names at record end: default is to put them at record start.\n")
	fmt.Fprintf(o, "-f {a,b,c} Field names to reorder.\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(
		o,
		"%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n",
		os.Args[0], verbNameReorder,
	)
	fmt.Fprintf(
		o,
		"%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n",
		os.Args[0], verbNameReorder,
	)

	if doExit {
		os.Exit(exitCode)
	}
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
