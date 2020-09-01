package mapping

import (
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/runtime"
)

var MapperCatSetup = MapperSetup{
	Verb:         "cat",
	ParseCLIFunc: mapperCatParseCLIFunc,
	UsageFunc:    mapperCatUsageFunc,
	IgnoresInput: false,
}

func mapperCatParseCLIFunc(
	pargi *int,
	argc int,
	args []string,
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) IRecordMapper {
	//	char* default_counter_field_name = DEFAULT_COUNTER_FIELD_NAME;
	//	char* counter_field_name = NULL;
	//	int   do_counters = FALSE;
	//	int   verbose = FALSE;
	//	slls_t* pgroup_by_field_names = slls_alloc();
	//
	//	if ((argc - *pargi) < 1) {
	//		mapper_cat_usage(stderr, argv[0], argv[*pargi]);
	//		return NULL;
	//	}
	//	char* verb = argv[*pargi];
	*pargi += 1
	//
	//	ap_state_t* pstate = ap_alloc();
	//	ap_define_true_flag(pstate, "-n",   &do_counters);
	//	ap_define_true_flag(pstate, "-v",   &verbose);
	//	ap_define_string_flag(pstate, "-N", &counter_field_name);
	//	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	//
	//	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
	//		mapper_cat_usage(stderr, argv[0], verb);
	//		return NULL;
	//	}
	//
	//	if (counter_field_name != NULL) {
	//		do_counters = TRUE;
	//	} else if (do_counters) {
	//		counter_field_name = default_counter_field_name;
	//	}
	//
	//	mapper_t* pmapper = mapper_cat_alloc(pstate, do_counters, verbose, counter_field_name, pgroup_by_field_names);
	//	return pmapper;

	// xxx temp err keep or no
	mapper, _ := NewMapperCat()
	return mapper
}

func mapperCatUsageFunc(
	o *os.File,
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	// xxx to do
	//	fmt.Fprintf(o, "Options:\n");
	//	fmt.Fprintf(o, "-n        Prepend field \"%s\" to each record with record-counter starting at 1\n",
	//		DEFAULT_COUNTER_FIELD_NAME);
	//	fmt.Fprintf(o, "-g {comma-separated field name(s)} When used with -n/-N, writes record-counters\n");
	//	fmt.Fprintf(o, "          keyed by specified field name(s).\n");
	//	fmt.Fprintf(o, "-v        Write a low-level record-structure dump to stderr.\n");
	//	fmt.Fprintf(o, "-N {name} Prepend field {name} to each record with record-counter starting at 1\n");

}

// ----------------------------------------------------------------
type MapperCat struct {
	// stateless
}

func NewMapperCat() (*MapperCat, error) {
	return &MapperCat{}, nil
}

func (this *MapperCat) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	outrecs <- inrec
}
