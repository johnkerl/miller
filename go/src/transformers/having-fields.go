package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

type tHavingFieldsCountSource int

const (
	havingFieldsCountSourceUnspecified tHavingFieldsCountSource = iota
	havingFieldsCountFromInt
	havingFieldsCountFromFieldName
)

//typedef enum _criterion_t {
//	HAVING_FIELDS_UNSPECIFIED,
//	HAVING_FIELDS_AT_LEAST,
//	HAVING_FIELDS_WHICH_ARE,
//	HAVING_FIELDS_AT_MOST,
//	HAVING_ALL_FIELDS_MATCHING,
//	HAVING_ANY_FIELDS_MATCHING,
//	HAVING_NO_FIELDS_MATCHING
//} criterion_t;

// ----------------------------------------------------------------
const verbNameHavingFields = "having-fields"

var HavingFieldsSetup = transforming.TransformerSetup{
	Verb:         verbNameHavingFields,
	UsageFunc:    transformerHavingFieldsUsage,
	ParseCLIFunc: transformerHavingFieldsParseCLI,

	IgnoresInput: false,
}

//mapper_setup_t mapper_having_fields_setup = {
//	.verb = "having-fields",
//	.pusage_func = mapper_having_fields_usage,
//	.pparse_func = mapper_having_fields_parse_cli,
//	.ignores_input = FALSE,
//};

func transformerHavingFieldsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameHavingFields)
	fmt.Fprintf(o, "Copies input records to output records multiple times.\n")
	fmt.Fprintf(o, "Options must be exactly one of the following:\n")
	fmt.Fprintf(o, "-n {havingFields count}  HavingFields each input record this many times.\n")
	fmt.Fprintf(o, "-f {field name}    Same, but take the havingFields count from the specified\n")
	fmt.Fprintf(o, "                   field name of each input record.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo x=0 | %s %s -n 4 then put '$x=urand()'\n", lib.MlrExeName(), verbNameHavingFields)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, " x=0.488189\n")
	fmt.Fprintf(o, " x=0.484973\n")
	fmt.Fprintf(o, " x=0.704983\n")
	fmt.Fprintf(o, " x=0.147311\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f b\n", lib.MlrExeName(), verbNameHavingFields)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "Example:\n")
	fmt.Fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f c\n", lib.MlrExeName(), verbNameHavingFields)
	fmt.Fprintf(o, "produces:\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	fmt.Fprintf(o, "  a=1,b=2,c=3\n")
	if doExit {
		os.Exit(exitCode)
	}
}

//// ----------------------------------------------------------------
//static void mapper_having_fields_usage(FILE* o, char* argv0, char* verb) {
//	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
//	fprintf(o, "Conditionally passes through records depending on each record's field names.\n");
//	fprintf(o, "Options:\n");
//	fprintf(o, "  --at-least      {comma-separated names}\n");
//	fprintf(o, "  --which-are     {comma-separated names}\n");
//	fprintf(o, "  --at-most       {comma-separated names}\n");
//	fprintf(o, "  --all-matching  {regular expression}\n");
//	fprintf(o, "  --any-matching  {regular expression}\n");
//	fprintf(o, "  --none-matching {regular expression}\n");
//	fprintf(o, "Examples:\n");
//	fprintf(o, "  %s %s --which-are amount,status,owner\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching 'sda[0-9]'\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"'\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"i' (this is case-insensitive)\n", argv0, verb);
//}

func transformerHavingFieldsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	havingFieldsCountSource := havingFieldsCountSourceUnspecified
	havingFieldsCount := 0
	havingFieldsCountFieldName := ""

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerHavingFieldsUsage(os.Stdout, true, 0)

		} else if opt == "-n" {
			havingFieldsCount = cliutil.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)
			havingFieldsCountSource = havingFieldsCountFromInt

		} else if opt == "-f" {
			havingFieldsCountFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)
			havingFieldsCountSource = havingFieldsCountFromFieldName

		} else {
			transformerHavingFieldsUsage(os.Stderr, true, 1)
		}
	}

	if havingFieldsCountSource == havingFieldsCountSourceUnspecified {
		transformerHavingFieldsUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerHavingFields(
		havingFieldsCountSource,
		havingFieldsCount,
		havingFieldsCountFieldName,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerHavingFields struct {
	havingFieldsCount           int
	havingFieldsCountFieldName  string
	recordTransformerFunc transforming.RecordTransformerFunc
}

//typedef struct _mapper_having_fields_state_t {
//	slls_t* pfield_names;
//	hss_t*  pfield_name_set;
//	regex_t regex;
//} mapper_having_fields_state_t;

// ----------------------------------------------------------------
func NewTransformerHavingFields(
	havingFieldsCountSource tHavingFieldsCountSource,
	havingFieldsCount int,
	havingFieldsCountFieldName string,
) (*TransformerHavingFields, error) {

	this := &TransformerHavingFields{
		havingFieldsCount:          havingFieldsCount,
		havingFieldsCountFieldName: havingFieldsCountFieldName,
	}

	if havingFieldsCountSource == havingFieldsCountFromInt {
		this.recordTransformerFunc = this.havingFieldsByCount
	} else {
		this.recordTransformerFunc = this.havingFieldsByFieldName
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerHavingFields) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerHavingFields) havingFieldsByCount(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		for i := 0; i < this.havingFieldsCount; i++ {
			outputChannel <- types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			)
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

// ----------------------------------------------------------------
func (this *TransformerHavingFields) havingFieldsByFieldName(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		fieldValue := inrecAndContext.Record.Get(this.havingFieldsCountFieldName)
		if fieldValue == nil {
			return
		}
		havingFieldsCount, ok := fieldValue.GetIntValue()
		if !ok {
			return
		}
		for i := 0; i < int(havingFieldsCount); i++ {
			outputChannel <- types.NewRecordAndContext(
				inrecAndContext.Record.Copy(),
				&inrecAndContext.Context,
			)
		}

	} else {
		outputChannel <- inrecAndContext
	}
}





//// ----------------------------------------------------------------
//static mapper_t* mapper_having_fields_parse_cli(int* pargi, int argc, char** argv,
//	cli_reader_opts_t* _, cli_writer_opts_t* __)
//{
//	slls_t*     pfield_names = NULL;
//	char*       regex_string = NULL;
//	criterion_t criterion    = HAVING_FIELDS_UNSPECIFIED;
//
//	char* verb = argv[(*pargi)++];
//
//	int argi = *pargi;
//	while (argi < argc && argv[argi][0] == '-') {
//
//		if (streq(argv[argi], "--at-least")) {
//			criterion = HAVING_FIELDS_AT_LEAST;
//			if (pfield_names != NULL)
//				slls_free(pfield_names);
//			pfield_names = slls_from_line(argv[argi+1], ',', FALSE);
//			regex_string = NULL;
//		} else if (streq(argv[argi], "--which-are")) {
//			criterion = HAVING_FIELDS_WHICH_ARE;
//			if (pfield_names != NULL)
//				slls_free(pfield_names);
//			pfield_names = slls_from_line(argv[argi+1], ',', FALSE);
//			regex_string = NULL;
//		} else if (streq(argv[argi], "--at-most")) {
//			criterion = HAVING_FIELDS_AT_MOST;
//			if (pfield_names != NULL)
//				slls_free(pfield_names);
//			pfield_names = slls_from_line(argv[argi+1], ',', FALSE);
//			regex_string = NULL;
//
//		} else if (streq(argv[argi], "--all-matching")) {
//			criterion = HAVING_ALL_FIELDS_MATCHING;
//			if (pfield_names != NULL) {
//				slls_free(pfield_names);
//				pfield_names = NULL;
//			}
//			regex_string = argv[argi+1];
//		} else if (streq(argv[argi], "--any-matching")) {
//			criterion = HAVING_ANY_FIELDS_MATCHING;
//			if (pfield_names != NULL) {
//				slls_free(pfield_names);
//				pfield_names = NULL;
//			}
//			regex_string = argv[argi+1];
//		} else if (streq(argv[argi], "--none-matching")) {
//			criterion = HAVING_NO_FIELDS_MATCHING;
//			if (pfield_names != NULL) {
//				slls_free(pfield_names);
//				pfield_names = NULL;
//			}
//			regex_string = argv[argi+1];
//		} else {
//			mapper_having_fields_usage(stderr, argv[0], verb);
//			return NULL;
//		}
//
//		if (argc - argi < 2) {
//			mapper_having_fields_usage(stderr, argv[0], verb);
//			return NULL;
//		}
//		argi += 2;
//	}
//
//	if (pfield_names == NULL && regex_string == NULL) {
//		mapper_having_fields_usage(stderr, argv[0], verb);
//		return NULL;
//	}
//	if (criterion == FALSE) {
//		mapper_having_fields_usage(stderr, argv[0], verb);
//		return NULL;
//	}
//
//	*pargi = argi;
//	return mapper_having_fields_alloc(pfield_names, regex_string, criterion);
//}

//// ----------------------------------------------------------------
//static void mapper_having_fields_usage(FILE* o, char* argv0, char* verb) {
//	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
//	fprintf(o, "Conditionally passes through records depending on each record's field names.\n");
//	fprintf(o, "Options:\n");
//	fprintf(o, "  --at-least      {comma-separated names}\n");
//	fprintf(o, "  --which-are     {comma-separated names}\n");
//	fprintf(o, "  --at-most       {comma-separated names}\n");
//	fprintf(o, "  --all-matching  {regular expression}\n");
//	fprintf(o, "  --any-matching  {regular expression}\n");
//	fprintf(o, "  --none-matching {regular expression}\n");
//	fprintf(o, "Examples:\n");
//	fprintf(o, "  %s %s --which-are amount,status,owner\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching 'sda[0-9]'\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"'\n", argv0, verb);
//	fprintf(o, "  %s %s --any-matching '\"sda[0-9]\"i' (this is case-insensitive)\n", argv0, verb);
//}

//// ----------------------------------------------------------------
//static mapper_t* mapper_having_fields_alloc(slls_t* pfield_names, char* regex_string, criterion_t criterion) {
//	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
//
//	mapper_having_fields_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_having_fields_state_t));
//
//	pmapper->pvstate = (void*)pstate;
//
//	if (regex_string != NULL) {
//		pstate->pfield_names    = NULL;
//		pstate->pfield_name_set = hss_alloc();
//
//		// Let them type in a.*b if they want, or "a.*b", or "a.*b"i.
//		// Strip off the leading " and trailing " or "i.
//		regcomp_or_die_quoted(&pstate->regex, regex_string, REG_NOSUB);
//
//		if (criterion == HAVING_ALL_FIELDS_MATCHING)
//			pmapper->pprocess_func = mapper_having_all_fields_matching_process;
//		else if (criterion == HAVING_ANY_FIELDS_MATCHING)
//			pmapper->pprocess_func = mapper_having_any_fields_matching_process;
//		else if (criterion == HAVING_NO_FIELDS_MATCHING)
//			pmapper->pprocess_func = mapper_having_no_fields_matching_process;
//		pmapper->pfree_func = mapper_having_fields_free;
//
//	} else {
//		pstate->pfield_names    = pfield_names;
//		pstate->pfield_name_set = hss_alloc();
//		regcomp_or_die(&pstate->regex, ".", 0);
//		for (sllse_t* pe = pfield_names->phead; pe != NULL; pe = pe->pnext)
//			hss_add(pstate->pfield_name_set, pe->value);
//
//		if (criterion == HAVING_FIELDS_AT_LEAST)
//			pmapper->pprocess_func = mapper_having_fields_at_least_process;
//		else if (criterion == HAVING_FIELDS_WHICH_ARE)
//			pmapper->pprocess_func = mapper_having_fields_which_are_process;
//		else if (criterion == HAVING_FIELDS_AT_MOST)
//			pmapper->pprocess_func = mapper_having_fields_at_most_process;
//		pmapper->pfree_func = mapper_having_fields_free;
//	}
//
//	return pmapper;
//}

//// ----------------------------------------------------------------
//static sllv_t* mapper_having_fields_at_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	int num_found = 0;
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (hss_has(pstate->pfield_name_set, pe->key)) {
//			num_found++;
//			if (num_found == pstate->pfield_name_set->num_occupied)
//				return sllv_single(pinrec);
//		}
//	}
//	lrec_free(pinrec);
//	return NULL;
//}
//
//static sllv_t* mapper_having_fields_which_are_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	if (pinrec->field_count != pstate->pfield_name_set->num_occupied) {
//		lrec_free(pinrec);
//		return NULL;
//	}
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (!hss_has(pstate->pfield_name_set, pe->key)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}
//
//static sllv_t* mapper_having_fields_at_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (!hss_has(pstate->pfield_name_set, pe->key)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}

//// ----------------------------------------------------------------
//static sllv_t* mapper_having_all_fields_matching_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (!regmatch_or_die(&pstate->regex, pe->key, 0, NULL)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}
//
//static sllv_t* mapper_having_any_fields_matching_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (regmatch_or_die(&pstate->regex, pe->key, 0, NULL)) {
//			return sllv_single(pinrec);
//		}
//	}
//	lrec_free(pinrec);
//	return NULL;
//}
//
//static sllv_t* mapper_having_no_fields_matching_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (regmatch_or_die(&pstate->regex, pe->key, 0, NULL)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}
