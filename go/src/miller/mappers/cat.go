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
var CatSetup = mapping.MapperSetup{
	Verb:         "cat",
	ParseCLIFunc: mapperCatParseCLI,
	IgnoresInput: false,
}

func mapperCatParseCLI(
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

	pDoCounters := flagSet.Bool(
		"n",
		false,
		"Prepend field \"n\" to each record with record-counter starting at 1",
	)

	pCounterFieldName := flagSet.String(
		"N",
		"",
		"Prepend field {name} to each record with record-counter starting at 1",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperCatUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// xxx to port:
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	//	fmt.Fprintf(o, "Options:\n");
	//	fmt.Fprintf(o, "-n        Prepend field \"%s\" to each record with record-counter starting at 1\n",
	//		DEFAULT_COUNTER_FIELD_NAME);
	//	fmt.Fprintf(o, "-g {comma-separated field name(s)} When used with -n/-N, writes record-counters\n");
	//	fmt.Fprintf(o, "          keyed by specified field name(s).\n");
	//	fmt.Fprintf(o, "-v        Write a low-level record-structure dump to stderr.\n");
	//	fmt.Fprintf(o, "-N {name} Prepend field {name} to each record with record-counter starting at 1\n");
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperCat(*pDoCounters, pCounterFieldName)

	*pargi = argi
	return mapper
}

func mapperCatUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperCat struct {
	doCounters bool

	counter          int64
	counterFieldName string
}

func NewMapperCat(
	doCounters bool,
	pCounterFieldName *string,
) (*MapperCat, error) {

	counterFieldName := "n"
	if *pCounterFieldName != "" {
		counterFieldName = *pCounterFieldName
		doCounters = true
	}

	return &MapperCat{
		doCounters:       doCounters,
		counter:          0,
		counterFieldName: counterFieldName,
	}, nil
}

func (this *MapperCat) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	record := inrecAndContext.Record
	if record != nil { // not end of record stream
		if this.doCounters {
			this.counter++
			key := this.counterFieldName
			value := types.MlrvalFromInt64(this.counter)
			record.PrependCopy(&key, &value)
		}
	}
	outputChannel <- inrecAndContext
}
