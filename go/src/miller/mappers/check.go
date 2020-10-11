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
var CheckSetup = mapping.MapperSetup{
	Verb:         "check",
	ParseCLIFunc: mapperCheckParseCLI,
	IgnoresInput: false,
}

func mapperCheckParseCLI(
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
		mapperCheckUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperCheck()

	*pargi = argi
	return mapper
}

func mapperCheckUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s\n", argv0, verb)
	fmt.Fprintf(o, "Consumes records without printing any output.\n")
	fmt.Fprintf(o, "Useful for doing a well-formatted check on input data.\n")

	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperCheck struct {
	// stateless
}

func NewMapperCheck() (*MapperCheck, error) {
	return &MapperCheck{}, nil
}

func (this *MapperCheck) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.Record == nil { // end of stream
		outputChannel <- inrecAndContext
	}
}
