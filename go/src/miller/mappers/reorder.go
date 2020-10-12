package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var ReorderSetup = mapping.MapperSetup{
	Verb:         "reorder",
	ParseCLIFunc: mapperReorderParseCLI,
	IgnoresInput: false,
}

func mapperReorderParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

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
		mapperReorderUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if *pFieldNames == "" {
		mapperReorderUsage(os.Stderr, args[0], verb, flagSet)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperReorder(
		*pFieldNames,
		*pPutAtEnd,
	)

	*pargi = argi
	return mapper
}

func mapperReorderUsage(
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
	fmt.Fprintf(o, "Examples:\n");
	fmt.Fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n",
		argv0, verb);
	fmt.Fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n",
		argv0, verb);
}

// ----------------------------------------------------------------
type MapperReorder struct {
	// input
	fieldNameList []string

	// state
	recordMapperFunc mapping.RecordMapperFunc
}

func NewMapperReorder(
	fieldNames string,
	putAtEnd bool,
) (*MapperReorder, error) {

	fieldNameList := lib.SplitString(fieldNames, ",")
	if !putAtEnd {
		lib.ReverseStringList(fieldNameList)
	}

	this := &MapperReorder{
		fieldNameList: fieldNameList,
	}

	if !putAtEnd {
		this.recordMapperFunc = this.reorderToStart
	} else {
		this.recordMapperFunc = this.reorderToEnd
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperReorder) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordMapperFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *MapperReorder) reorderToStart(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for _, fieldName := range(this.fieldNameList) {
			inrec.MoveToHead(&fieldName)
		}
		outputChannel <- inrecAndContext

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *MapperReorder) reorderToEnd(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		for _, fieldName := range(this.fieldNameList) {
			inrec.MoveToTail(&fieldName)
		}
		outputChannel <- inrecAndContext
	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
