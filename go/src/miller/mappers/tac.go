package mappers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var TacSetup = mapping.MapperSetup{
	Verb:         "tac",
	ParseCLIFunc: mapperTacParseCLI,
	IgnoresInput: false,
}

func mapperTacParseCLI(
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
		mapperTacUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperTac()

	*pargi = argi
	return mapper
}

func mapperTacUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s\n", argv0, verb)
	fmt.Fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperTac struct {
	recordsAndContexts *list.List
}

func NewMapperTac() (*MapperTac, error) {
	return &MapperTac{
		recordsAndContexts: list.New(),
	}, nil
}

func (this *MapperTac) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.Record != nil {
		this.recordsAndContexts.PushFront(inrecAndContext)
	} else {
		// end of stream
		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outputChannel <- e.Value.(*types.RecordAndContext)
		}
		outputChannel <- types.NewRecordAndContext(
			nil, // signals end of input record stream
			&inrecAndContext.Context,
		)
	}
}
