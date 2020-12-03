package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var ReorderSetup = transforming.TransformerSetup{
	Verb:         "reorder",
	ParseCLIFunc: transformerReorderParseCLI,
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

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pFieldNames := flagSet.String(
		"f",
		"",
		"Field names to reorder",
	)

	pPutAtEnd := flagSet.Bool(
		"e",
		false,
		"Put specified field names at record end: default is to put them at record start",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerReorderUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if *pFieldNames == "" {
		transformerReorderUsage(os.Stderr, args[0], verb, flagSet)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerReorder(
		*pFieldNames,
		*pPutAtEnd,
	)

	*pargi = argi
	return transformer
}

func transformerReorderUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Moves specified names to start of record, or end of record.
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
	})
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n",
		argv0, verb)
	fmt.Fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n",
		argv0, verb)
}

// ----------------------------------------------------------------
type TransformerReorder struct {
	// input
	fieldNameList []string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
}

func NewTransformerReorder(
	fieldNames string,
	putAtEnd bool,
) (*TransformerReorder, error) {

	fieldNameList := lib.SplitString(fieldNames, ",")
	if !putAtEnd {
		lib.ReverseStringList(fieldNameList)
	}

	this := &TransformerReorder{
		fieldNameList: fieldNameList,
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
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for _, fieldName := range this.fieldNameList {
			inrec.MoveToHead(&fieldName)
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
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for _, fieldName := range this.fieldNameList {
			inrec.MoveToTail(&fieldName)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
