package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/mapping"
)

// ----------------------------------------------------------------
var NothingSetup = mapping.MapperSetup{
	Verb:         "nothing",
	ParseCLIFunc: mapperNothingParseCLI,
	IgnoresInput: false,
}

func mapperNothingParseCLI(
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
		mapperNothingUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperNothing()

	*pargi = argi
	return mapper
}

func mapperNothingUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s\n", argv0, verb)
	fmt.Fprintf(o, "Drops all input records. Useful for testing, or after tee/print/etc. have\n")
	fmt.Fprintf(o, "produced other output.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperNothing struct {
	// stateless
}

func NewMapperNothing() (*MapperNothing, error) {
	return &MapperNothing{}, nil
}

func (this *MapperNothing) Map(
	inrecAndContext *containers.LrecAndContext,
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	if inrecAndContext.Lrec == nil { // end of stream
		outrecsAndContexts <- inrecAndContext
	}
}
