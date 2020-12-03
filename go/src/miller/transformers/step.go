package transformers

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

const DEFAULT_STRING_ALPHA = "0.5"

// ----------------------------------------------------------------
var StepSetup = transforming.TransformerSetup{
	Verb:         "step",
	ParseCLIFunc: transformerStepParseCLI,
	IgnoresInput: false,
}

func transformerStepParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pStepperNamesString := flagSet.String(
		"a",
		"",
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	pValueFieldNamesString := flagSet.String(
		"f",
		"",
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	pGroupByFieldNamesString := flagSet.String(
		"g",
		"",
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	pAllowIntFloat := flagSet.Bool(
		"F",
		false,
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	pAlphasString := flagSet.String(
		"d",
		DEFAULT_STRING_ALPHA,
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	pEWMASuffixesString := flagSet.String(
		"o",
		"",
		`It's a coding error if you see this.`, // Handled in transformerStepUsage() below
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerStepUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	stepperNames := lib.SplitString(*pStepperNamesString, ",")
	valueFieldNames := lib.SplitString(*pValueFieldNamesString, ",")
	groupByFieldNames := lib.SplitString(*pGroupByFieldNamesString, ",")
	stringAlphas := lib.SplitString(*pAlphasString, ",")
	EWMASuffixes := lib.SplitString(*pEWMASuffixesString, ",")

	transformer, err := NewTransformerStep(
		stepperNames,
		valueFieldNames,
		groupByFieldNames,
		stringAlphas,
		EWMASuffixes,
		*pAllowIntFloat,
	)
	// TODO: put error return into this API
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	*pargi = argi
	return transformer
}

func transformerStepUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Computes values dependent on the previous record, optionally grouped by category.\n")
	// For this transformer, we do NOT use flagSet.VisitAll -- we use our own print statements.

	fmt.Fprintf(o, "-a {delta,rsum,...}   Names of steppers: comma-separated, one or more of:\n")
	// TODO
	//	for (int i = 0; i < step_lookup_table_length; i++) {
	//		fprintf(o, "  %-8s %s\n", step_lookup_table[i].name, step_lookup_table[i].desc);
	//	}

	fmt.Fprintf(o, "-f {a,b,c} Value-field names on which to compute statistics\n")

	fmt.Fprintf(o, "-g {d,e,f} Optional group-by-field names\n")

	fmt.Fprintf(o, "-F         Computes integerable things (e.g. counter) in floating point.\n")

	fmt.Fprintf(o, "-d {x,y,z} Weights for ewma. 1 means current sample gets all weight (no\n")
	fmt.Fprintf(o, "           smoothing), near under under 1 is light smoothing, near over 0 is\n")
	fmt.Fprintf(o, "           heavy smoothing. Multiple weights may be specified, e.g.\n")
	fmt.Fprintf(o, "           \"%s %s -a ewma -f sys_load -d 0.01,0.1,0.9\". Default if omitted\n", argv0, verb)
	fmt.Fprintf(o, "           is \"-d %s\".\n", DEFAULT_STRING_ALPHA)

	fmt.Fprintf(o, "-o {a,b,c} Custom suffixes for EWMA output fields. If omitted, these default to\n")
	fmt.Fprintf(o, "           the -d values. If supplied, the number of -o values must be the same\n")
	fmt.Fprintf(o, "           as the number of -d values.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s -a rsum -f request_size\n", argv0, verb)
	fmt.Fprintf(o, "  %s %s -a delta -f request_size -g hostname\n", argv0, verb)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -f x,y\n", argv0, verb)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y\n", argv0, verb)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name\n", argv0, verb)

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#filter or\n")
	fmt.Fprintf(o, "https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average\n")
	fmt.Fprintf(o, "for more information on EWMA.\n")
}

// ----------------------------------------------------------------
type tStepper struct {
}

type TransformerStep struct {
	// Input:
	stepperNames    []string
	valueFieldNames []string
	groupByFieldNames []string

	// State:

	valueFieldValues []types.Mlrval // scratch space used per-record

	// Map from group-by field names to value-field names to array of stepper objects.
	// See the Transform method below for more details.
	groups map[string]map[string][]tStepper
}

func NewTransformerStep(
	stepperNames []string,
	valueFieldNames []string,
	groupByFieldNames []string,
	stringAlphas []string,
	EWMASuffixes []string,
	allowIntFloat bool,
) (*TransformerStep, error) {

	if len(stepperNames) == 0 || len(valueFieldNames) == 0 {
		return nil, errors.New(
			// TODO: parameterize verb here somehow
			"mlr step: -a and -f are both required arguments.",
		)
	}
	if len(stringAlphas) != 0 && len(EWMASuffixes) != 0 {
		if len(EWMASuffixes) != len(stringAlphas) {
			return nil, errors.New(
				// TODO: parameterize verb here somehow
				"mlr step: If -d and -o are provied, their values must have the same length.",
			)
		}
	}

	// TODO: flesh out
	this := &TransformerStep{
		stepperNames:    stepperNames,
		valueFieldNames: valueFieldNames,
		groupByFieldNames: groupByFieldNames,
	}

	return this, nil
}

// ----------------------------------------------------------------
// Multilevel hashmap structure:
// {
//   ["s","t"] : {              <--- group-by field names
//     ["x","y"] : {            <--- value field names
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
// }

func (this *TransformerStep) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	if inrec == nil { // end of record stream
		outputChannel <- inrecAndContext
	}

	// ["s", "t"]
	groupingKey, groupByFieldValues, gok := inrec.GetSelectedValuesAndJoined(this.groupByFieldNames)
	if !gok { // current record doesn't have fields to be stepped; pass it along
		outputChannel <- inrecAndContext
		return
	}

	// ["x", "y"]
	valueFieldValues, vok := inrec.GetSelectedValues(this.valueFieldNames)
	if !vok { // current record doesn't have fields to be stepped; pass it along
		outputChannel <- inrecAndContext
		return
	}

	groupToAccField := this.groups[groupingKey]
	if groupToAccField == nil {
		// Populate the groups data structure on first reference if needed
		groupToAccField := make(map[string][]tStepper)
		this.groups[groupingKey] = groupToAccField
	}

	// for x=1 and y=2:
	for i, valueFieldName := range this.valueFieldNames {
		// TODO: make it sparse in the GetSelectedValues() ... no `vok` return ...
		valueFieldValue := valueFieldValues[i]
		//		char* value_field_sval = pstate->pvalue_field_values->strings[i];
		//		if (value_field_sval == NULL) // Key not present
		//			continue;

		//		int have_dval = FALSE;
		//		int have_nval = FALSE;
		//		double value_field_dval = -999.0;
		//		mv_t   value_field_nval = mv_absent();
		//
		//		lhmsv_t* pacc_field_to_acc_state = lhmsv_get(pgroup_to_acc_field, value_field_name);
		//		if (pacc_field_to_acc_state == NULL) {
		//			pacc_field_to_acc_state = lhmsv_alloc();
		//			lhmsv_put(pgroup_to_acc_field, value_field_name, pacc_field_to_acc_state, NO_FREE);
		//		}

		// for "delta", "rsum"
		for _, stepperName := range this.stepperNames {
			//			tStepper* stepper = lhmsv_get(pacc_field_to_acc_state, stepperName);
			//			if (stepper == NULL) {
			//				stepper = make_step(stepperName, value_field_name, pstate->allow_int_float,
			//					pstate->pstring_alphas, pstate->pewma_suffixes);
			//				if (stepper == NULL) {
			//					fprintf(stderr, "mlr step: stepper \"%s\" not found.\n",
			//						stepperName);
			//					exit(1);
			//				}
			//				lhmsv_put(pacc_field_to_acc_state, stepperName, stepper, NO_FREE);
			//			}
			//
			//			if (*value_field_sval == 0) { // Key present with null value
			//				if (stepper->pzprocess_func != NULL) {
			//					stepper->pzprocess_func(stepper->pvstate, pinrec);
			//				}
			//			} else {
			//
			//				if (stepper->pdprocess_func != NULL) {
			//					if (!have_dval) {
			//						value_field_dval = mlr_double_from_string_or_die(value_field_sval);
			//						have_dval = TRUE;
			//					}
			//					stepper->pdprocess_func(stepper->pvstate, value_field_dval, pinrec);
			//				}
			//
			//				if (stepper->pnprocess_func != NULL) {
			//					if (!have_nval) {
			//						value_field_nval = pstate->allow_int_float
			//							? mv_scan_number_or_die(value_field_sval)
			//							: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
			//						have_nval = TRUE;
			//					}
			//					stepper->pnprocess_func(stepper->pvstate, &value_field_nval, pinrec);
			//				}
			//
			//				if (stepper->psprocess_func != NULL) {
			//					stepper->psprocess_func(stepper->pvstate, value_field_sval, pinrec);
			//				}
			//			}
		}
	}

	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
//typedef struct _step_lookup_t {
//	name string
//	allocFunc stepAllocFunc
//	desc string;
//} step_lookup_t;

//static step_lookup_t step_lookup_table[] = {
//	{"delta",      step_delta_alloc,      "Compute differences in field(s) between successive records"},
//	{"shift",      step_shift_alloc,      "Include value(s) in field(s) from previous record, if any"},
//	{"from-first", step_from_first_alloc, "Compute differences in field(s) from first record"},
//	{"ratio",      step_ratio_alloc,      "Compute ratios in field(s) between successive records"},
//	{"rsum",       step_rsum_alloc,       "Compute running sums of field(s) between successive records"},
//	{"counter",    step_counter_alloc,    "Count instances of field(s) between successive records"},
//	{"ewma",       step_ewma_alloc,       "Exponentially weighted moving average over successive records"},
//};
// ----------------------------------------------------------------

//static tStepper* make_step(char* step_name, char* input_field_name, int allow_int_float,
//	slls_t* pstring_alphas, slls_t* pewma_suffixes)
//{
//	for (int i = 0; i < step_lookup_table_length; i++)
//		if (streq(step_name, step_lookup_table[i].name))
//			return step_lookup_table[i].palloc_func(input_field_name, allow_int_float,
//				pstring_alphas, pewma_suffixes);
//	return NULL;
//}

//// ----------------------------------------------------------------
//typedef struct _step_delta_state_t {
//	mv_t  prev;
//	char* output_field_name;
//	int   allow_int_float;
//} step_delta_state_t;
//static void step_delta_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
//	step_delta_state_t* pstate = pvstate;
//	mv_t delta;
//	if (mv_is_null(&pstate->prev)) {
//		delta = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
//	} else {
//		delta = x_xx_minus_func(pnumv, &pstate->prev);
//	}
//	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&delta), FREE_ENTRY_VALUE);
//	pstate->prev = *pnumv;
//}
//static void step_delta_zprocess(void* pvstate, lrec_t* prec) {
//	step_delta_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_delta_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_delta_state_t* pstate = mlr_malloc_or_die(sizeof(step_delta_state_t));
//	pstate->prev = mv_absent();
//	pstate->allow_int_float = allow_int_float;
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_delta");
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = NULL;
//	stepper->pnprocess_func = step_delta_nprocess;
//	stepper->psprocess_func = NULL;
//	stepper->pzprocess_func = step_delta_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_shift_state_t {
//	char* prev;
//	char* output_field_name;
//	int   allow_int_float;
//} step_shift_state_t;
//static void step_shift_sprocess(void* pvstate, char* strv, lrec_t* prec) {
//	step_shift_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, pstate->prev, FREE_ENTRY_VALUE);
//	pstate->prev = mlr_strdup_or_die(strv);
//}
//static void step_shift_zprocess(void* pvstate, lrec_t* prec) {
//	step_shift_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_shift_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_shift_state_t* pstate = mlr_malloc_or_die(sizeof(step_shift_state_t));
//	pstate->prev = mlr_strdup_or_die("");
//	pstate->allow_int_float = allow_int_float;
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_shift");
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = NULL;
//	stepper->pnprocess_func = NULL;
//	stepper->psprocess_func = step_shift_sprocess;
//	stepper->pzprocess_func = step_shift_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_from_first_state_t {
//	mv_t  first;
//	char* output_field_name;
//	int   allow_int_float;
//} step_from_first_state_t;
//static void step_from_first_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
//	step_from_first_state_t* pstate = pvstate;
//	mv_t from_first;
//	if (mv_is_null(&pstate->first)) {
//		from_first = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
//		pstate->first = *pnumv;
//	} else {
//		from_first = x_xx_minus_func(pnumv, &pstate->first);
//	}
//	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&from_first), FREE_ENTRY_VALUE);
//}
//static void step_from_first_zprocess(void* pvstate, lrec_t* prec) {
//	step_from_first_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_from_first_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_from_first_state_t* pstate = mlr_malloc_or_die(sizeof(step_from_first_state_t));
//	pstate->first = mv_absent();
//	pstate->allow_int_float = allow_int_float;
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_from_first");
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = NULL;
//	stepper->pnprocess_func = step_from_first_nprocess;
//	stepper->psprocess_func = NULL;
//	stepper->pzprocess_func = step_from_first_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_ratio_state_t {
//	double prev;
//	int    have_prev;
//	char*  output_field_name;
//} step_ratio_state_t;
//static void step_ratio_dprocess(void* pvstate, double fltv, lrec_t* prec) {
//	step_ratio_state_t* pstate = pvstate;
//	double ratio = 1.0;
//	if (pstate->have_prev) {
//		ratio = fltv / pstate->prev;
//	} else {
//		pstate->have_prev = TRUE;
//	}
//	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_double(ratio, MLR_GLOBALS.ofmt),
//		FREE_ENTRY_VALUE);
//	pstate->prev = fltv;
//}
//static void step_ratio_zprocess(void* pvstate, lrec_t* prec) {
//	step_ratio_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_ratio_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_ratio_state_t* pstate = mlr_malloc_or_die(sizeof(step_ratio_state_t));
//	pstate->prev          = -999.0;
//	pstate->have_prev     = FALSE;
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_ratio");
//
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = step_ratio_dprocess;
//	stepper->pnprocess_func = NULL;
//	stepper->psprocess_func = NULL;
//	stepper->pzprocess_func = step_ratio_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_rsum_state_t {
//	mv_t   rsum;
//	char*  output_field_name;
//	int    allow_int_float;
//} step_rsum_state_t;
//static void step_rsum_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
//	step_rsum_state_t* pstate = pvstate;
//	pstate->rsum = x_xx_plus_func(&pstate->rsum, pnumv);
//	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&pstate->rsum),
//		FREE_ENTRY_VALUE);
//}
//static void step_rsum_zprocess(void* pvstate, lrec_t* prec) {
//	step_rsum_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_rsum_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_rsum_state_t* pstate = mlr_malloc_or_die(sizeof(step_rsum_state_t));
//	pstate->allow_int_float = allow_int_float;
//	pstate->rsum = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_rsum");
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = NULL;
//	stepper->pnprocess_func = step_rsum_nprocess;
//	stepper->psprocess_func = NULL;
//	stepper->pzprocess_func = step_rsum_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_counter_state_t {
//	mv_t counter;
//	mv_t one;
//	char*  output_field_name;
//} step_counter_state_t;
//static void step_counter_sprocess(void* pvstate, char* strv, lrec_t* prec) {
//	step_counter_state_t* pstate = pvstate;
//	pstate->counter = x_xx_plus_func(&pstate->counter, &pstate->one);
//	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&pstate->counter),
//		FREE_ENTRY_VALUE);
//}
//static void step_counter_zprocess(void* pvstate, lrec_t* prec) {
//	step_counter_state_t* pstate = pvstate;
//	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
//}
//static tStepper* step_counter_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_counter_state_t* pstate = mlr_malloc_or_die(sizeof(step_counter_state_t));
//	pstate->counter = allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
//	pstate->one     = allow_int_float ? mv_from_int(1LL) : mv_from_float(1.0);
//	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_counter");
//
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = NULL;
//	stepper->pnprocess_func = NULL;
//	stepper->psprocess_func = step_counter_sprocess;
//	stepper->pzprocess_func = step_counter_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
//typedef struct _step_ewma_state_t {
//	int     num_alphas;
//	double* alphas;
//	double* alphacompls;
//	double* prevs;
//	int     have_prevs;
//	char**  output_field_names;
//} step_ewma_state_t;
//static void step_ewma_dprocess(void* pvstate, double fltv, lrec_t* prec) {
//	step_ewma_state_t* pstate = pvstate;
//	if (!pstate->have_prevs) {
//		for (int i = 0; i < pstate->num_alphas; i++) {
//			lrec_put(prec, pstate->output_field_names[i], mlr_alloc_string_from_double(fltv, MLR_GLOBALS.ofmt),
//				FREE_ENTRY_VALUE);
//			pstate->prevs[i] = fltv;
//		}
//		pstate->have_prevs = TRUE;
//	} else {
//		for (int i = 0; i < pstate->num_alphas; i++) {
//			double curr = fltv;
//			curr = pstate->alphas[i] * curr + pstate->alphacompls[i] * pstate->prevs[i];
//			lrec_put(prec, pstate->output_field_names[i], mlr_alloc_string_from_double(curr, MLR_GLOBALS.ofmt),
//				FREE_ENTRY_VALUE);
//			pstate->prevs[i] = curr;
//		}
//	}
//}
//static void step_ewma_zprocess(void* pvstate, lrec_t* prec) {
//	step_ewma_state_t* pstate = pvstate;
//	for (int i = 0; i < pstate->num_alphas; i++)
//		lrec_put(prec, pstate->output_field_names[i], "", NO_FREE);
//}

//static tStepper* step_ewma_alloc(char* input_field_name, int unused, slls_t* pstring_alphas, slls_t* pewma_suffixes) {
//	tStepper* stepper              = mlr_malloc_or_die(sizeof(tStepper));
//
//	step_ewma_state_t* pstate  = mlr_malloc_or_die(sizeof(step_ewma_state_t));
//	int n                      = pstring_alphas->length;
//	pstate->num_alphas         = n;
//	pstate->alphas             = mlr_malloc_or_die(n * sizeof(double));
//	pstate->alphacompls        = mlr_malloc_or_die(n * sizeof(double));
//	pstate->prevs              = mlr_malloc_or_die(n * sizeof(double));
//	pstate->have_prevs         = FALSE;
//	pstate->output_field_names = mlr_malloc_or_die(n * sizeof(char*));
//	slls_t* psuffixes = (pewma_suffixes == NULL) ? pstring_alphas : pewma_suffixes;
//	sllse_t* pe = pstring_alphas->phead;
//	sllse_t* pf = psuffixes->phead;
//	for (int i = 0; i < n; i++, pe = pe->pnext, pf = pf->pnext) {
//		char* string_alpha     = pe->value;
//		char* suffix           = pf->value;
//		pstate->alphas[i]      = mlr_double_from_string_or_die(string_alpha);
//		pstate->alphacompls[i] = 1.0 - pstate->alphas[i];
//		pstate->prevs[i]       = 0.0;
//		pstate->output_field_names[i] = mlr_paste_3_strings(input_field_name, "_ewma_", suffix);
//	}
//	pstate->have_prevs = FALSE;
//
//	stepper->pvstate        = (void*)pstate;
//	stepper->pdprocess_func = step_ewma_dprocess;
//	stepper->pnprocess_func = NULL;
//	stepper->psprocess_func = NULL;
//	stepper->pzprocess_func = step_ewma_zprocess;
//	return stepper;
//}
