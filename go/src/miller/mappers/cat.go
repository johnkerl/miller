package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/lib"
	"miller/mapping"
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
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// xxx to be ported:
	//
	//	char* default_counter_field_name = DEFAULT_COUNTER_FIELD_NAME;
	//	char* counter_field_name = nil;
	//	int   do_counters = FALSE;
	//	int   verbose = FALSE;
	//	slls_t* pgroup_by_field_names = slls_alloc();
	//
	//	ap_state_t* pstate = ap_alloc();
	//	ap_define_true_flag(pstate, "-n",   &do_counters);
	//	ap_define_true_flag(pstate, "-v",   &verbose);
	//	ap_define_string_flag(pstate, "-N", &counter_field_name);
	//	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	//
	//	if (!ap_parse(pstate, verb, pargi, argc, args)) {
	//		mapper_cat_usage(stderr, args[0], verb);
	//		return nil;
	//	}
	//
	//	if (counter_field_name != nil) {
	//		do_counters = TRUE;
	//	} else if (do_counters) {
	//		counter_field_name = default_counter_field_name;
	//	}
	//
	//	mapper_t* pmapper = mapper_cat_alloc(pstate, do_counters, verbose, counter_field_name, pgroup_by_field_names);
	//	return pmapper;
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

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
	inrecAndContext *containers.LrecAndContext,
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	lrec := inrecAndContext.Lrec
	if lrec != nil { // not end of record stream
		if this.doCounters {
			this.counter++
			key := this.counterFieldName
			value := lib.MlrvalFromInt64(this.counter)
			lrec.Prepend(&key, &value)
		}
	}
	outrecsAndContexts <- inrecAndContext
}
