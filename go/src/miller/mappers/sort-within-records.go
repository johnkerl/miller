package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var SortWithinRecordsSetup = mapping.MapperSetup{
	Verb:         "sort-within-records",
	ParseCLIFunc: mapperSortWithinRecordsParseCLI,
	IgnoresInput: false,
}

func mapperSortWithinRecordsParseCLI(
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

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperSortWithinRecordsUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperSortWithinRecords()

	*pargi = argi
	return mapper
}

func mapperSortWithinRecordsUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Outputs records sorted lexically ascending by keys.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperSortWithinRecords struct {
}

func NewMapperSortWithinRecords() (*MapperSortWithinRecords, error) {

	this := &MapperSortWithinRecords{
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperSortWithinRecords) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		inrec.SortByKey()

	} else {
	}
	outputChannel <- inrecAndContext // end-of-stream marker
}
