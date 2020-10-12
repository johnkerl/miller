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
var SkipTrivialRecordsSetup = mapping.MapperSetup{
	Verb:         "skip-trivial-records",
	ParseCLIFunc: mapperSkipTrivialRecordsParseCLI,
	IgnoresInput: false,
}

func mapperSkipTrivialRecordsParseCLI(
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
		mapperSkipTrivialRecordsUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperSkipTrivialRecords()

	*pargi = argi
	return mapper
}

func mapperSkipTrivialRecordsUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s {no options}\n", argv0, verb)
	fmt.Fprintf(o, "Passes through all records except those with zero fields,\n")
	fmt.Fprintf(o, "or those for which all fields have empty value.\n")

	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperSkipTrivialRecords struct {
}

func NewMapperSkipTrivialRecords() (*MapperSkipTrivialRecords, error) {
	this := &MapperSkipTrivialRecords{}
	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperSkipTrivialRecords) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		hasAny := false
		for pe := inrec.Head; pe != nil; pe = pe.Next {
			if pe.Value.String() != "" {
				hasAny = true
				break
			}
		}

		if hasAny {
			outputChannel <- inrecAndContext
		}

	} else {
		outputChannel <- inrecAndContext
	}
}
