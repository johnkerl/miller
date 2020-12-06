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
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		fmt.Fprintf(o, "  %-8s %s\n", stepperLookup.name, stepperLookup.desc)
	}
	fmt.Fprintf(o, "\n")

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
type TransformerStep struct {
	// INPUT
	stepperNames      []string
	valueFieldNames   []string
	groupByFieldNames []string
	allowIntFloat     bool
	stringAlphas      []string
	EWMASuffixes      []string

	// STATE
	// Scratch space used per-record
	valueFieldValues []types.Mlrval
	// Map from group-by field names to value-field names to array of
	// stepper objects.  See the Transform method below for more details.
	groups map[string]map[string]map[string]tStepper
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
		stepperNames:      stepperNames,
		valueFieldNames:   valueFieldNames,
		groupByFieldNames: groupByFieldNames,
		allowIntFloat:     allowIntFloat,
		stringAlphas:      stringAlphas,
		EWMASuffixes:      EWMASuffixes,
		groups:            make(map[string]map[string]map[string]tStepper),
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
		return
	}

	// ["s", "t"]
	groupingKey, gok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
	if !gok { // current record doesn't have fields to be stepped; pass it along
		outputChannel <- inrecAndContext
		return
	}

	// ["x", "y"]
	valueFieldValues := inrec.ReferenceSelectedValues(this.valueFieldNames)

	groupToAccField := this.groups[groupingKey]
	if groupToAccField == nil {
		// Populate the groups data structure on first reference if needed
		groupToAccField = make(map[string]map[string]tStepper)
		this.groups[groupingKey] = groupToAccField
	}

	// for x=1 and y=2:
	for i, valueFieldName := range this.valueFieldNames {
		// TODO: make it sparse in the GetSelectedValues() ... no `vok` return ...
		valueFieldValue := valueFieldValues[i]
		if valueFieldValue == nil { // not present in the current record
			continue
		}

		accFieldToAccState := groupToAccField[valueFieldName]
		if accFieldToAccState == nil {
			accFieldToAccState = make(map[string]tStepper)
			groupToAccField[valueFieldName] = accFieldToAccState
		}

		// for "delta", "rsum"
		for _, stepperName := range this.stepperNames {
			stepper, present := accFieldToAccState[stepperName]
			if !present {
				stepper = allocateStepper(
					stepperName,
					valueFieldName,
					this.allowIntFloat,
					this.stringAlphas,
					this.EWMASuffixes,
				)
				if stepper == nil {
					// TODO: parameterize verb name
					fmt.Fprintf(os.Stderr, "mlr step: stepper \"%s\" not found.\n",
						stepperName)
					os.Exit(1)
				}
				accFieldToAccState[stepperName] = stepper
			}

			// xxx
			// https://stackoverflow.com/questions/44370277/type-is-pointer-to-interface-not-interface-confusion
			//
			// A pointer to a struct and a pointer to an interface are not the
			// same.
			//
			// An interface can store either a struct directly or a pointer to
			// a struct. In the latter case, you still just use the interface
			// directly, not a pointer to the interface. For example:

			stepper.process(valueFieldValue, inrec)
		}
	}

	outputChannel <- inrecAndContext
}

// ================================================================
// File-local utility function

func makeIntOrFloat(literal int64, allowIntFloat bool) types.Mlrval {
	if allowIntFloat {
		return types.MlrvalFromInt64(literal)
	} else {
		return types.MlrvalFromFloat64(float64(literal))
	}
}

// ================================================================
// Lookups for individual steppers, like "delta" or "rsum"

type tStepperAllocator func(
	inputFieldName string,
	allowIntFloat bool,
	stringAlphas []string,
	EWMASuffixes []string,
) tStepper

type tStepper interface {
	process(valueFieldValue *types.Mlrval, inputRecord *types.Mlrmap)
}

type tStepperLookup struct {
	name             string
	stepperAllocator tStepperAllocator
	desc             string
}

var STEPPER_LOOKUP_TABLE = []tStepperLookup{
	{"delta", stepperDeltaAlloc, "Compute differences in field(s) between successive records"},
	{"shift", stepperShiftAlloc, "Include value(s) in field(s) from previous record, if any"},
	{"from-first", stepperFromFirstAlloc, "Compute differences in field(s) from first record"},
	{"ratio", stepperRatioAlloc, "Compute ratios in field(s) between successive records"},
	//{"rsum", stepperRsumAlloc, "Compute running sums of field(s) between successive records"},
	//{"counter", stepperCounterAlloc, "Count instances of field(s) between successive records"},
	//{"ewma", stepperEwmaAlloc, "Exponentially weighted moving average over successive records"},
}

func allocateStepper(
	stepperName string,
	inputFieldName string,
	allowIntFloat bool,
	stringAlphas []string,
	EWMASuffixes []string,
) tStepper {
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		if stepperLookup.name == stepperName {
			return stepperLookup.stepperAllocator(
				inputFieldName,
				allowIntFloat,
				stringAlphas,
				EWMASuffixes,
			)
		}
	}
	return nil
}

// ================================================================
// Implementations of individual steppers, like "delta" or "rsum"

// ================================================================
type tStepperDelta struct {
	previous        *types.Mlrval
	outputFieldName string
	allowIntFloat   bool
}

func stepperDeltaAlloc(
	inputFieldName string,
	allowIntFloat bool,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperDelta{
		previous:        nil,
		outputFieldName: inputFieldName + "_delta",
		allowIntFloat:   allowIntFloat,
	}
}

func (this *tStepperDelta) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var delta types.Mlrval
	if this.previous == nil {
		delta = makeIntOrFloat(0, this.allowIntFloat)
	} else {
		delta = types.MlrvalBinaryMinus(valueFieldValue, this.previous)
	}
	inrec.PutCopy(&this.outputFieldName, &delta)

	this.previous = valueFieldValue.Copy()

	// TODO: from C impl: if input is empty:
	// lrec_put(prec, this.output_field_name, "", NO_FREE);
}

// ================================================================
type tStepperShift struct {
	previous        *types.Mlrval
	outputFieldName string
	allowIntFloat   bool
}

func stepperShiftAlloc(
	inputFieldName string,
	allowIntFloat bool,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperShift{
		previous:        nil,
		outputFieldName: inputFieldName + "_shift",
		allowIntFloat:   allowIntFloat,
	}
}

func (this *tStepperShift) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if this.previous == nil {
		shift := types.MlrvalFromVoid()
		inrec.PutCopy(&this.outputFieldName, &shift)
	} else {
		inrec.PutCopy(&this.outputFieldName, this.previous)
		this.previous = valueFieldValue.Copy()
	}
	this.previous = valueFieldValue.Copy()
}

// ================================================================
type tStepperFromFirst struct {
	first           *types.Mlrval
	outputFieldName string
	allowIntFloat   bool
}

func stepperFromFirstAlloc(
	inputFieldName string,
	allowIntFloat bool,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperFromFirst{
		first:           nil,
		outputFieldName: inputFieldName + "_from_first",
		allowIntFloat:   allowIntFloat,
	}
}

func (this *tStepperFromFirst) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var from_first types.Mlrval
	if this.first == nil {
		from_first = makeIntOrFloat(0, this.allowIntFloat)
		this.first = valueFieldValue.Copy()
	} else {
		from_first = types.MlrvalBinaryMinus(valueFieldValue, this.first)
	}
	inrec.PutCopy(&this.outputFieldName, &from_first)

	// TODO: from C impl: if input is empty:
	// lrec_put(prec, this.output_field_name, "", NO_FREE);
}

// ================================================================
type tStepperRatio struct {
	previous        *types.Mlrval
	outputFieldName string
	allowIntFloat   bool
}

func stepperRatioAlloc(
	inputFieldName string,
	allowIntFloat bool,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperRatio{
		previous:        nil,
		outputFieldName: inputFieldName + "_ratio",
		allowIntFloat:   allowIntFloat,
	}
}

func (this *tStepperRatio) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var ratio types.Mlrval
	if this.previous == nil {
		ratio = makeIntOrFloat(1, this.allowIntFloat)
	} else {
		ratio = types.MlrvalDivide(valueFieldValue, this.previous)
	}
	inrec.PutCopy(&this.outputFieldName, &ratio)

	this.previous = valueFieldValue.Copy()

	// TODO: from C impl: if input is empty:
	// lrec_put(prec, this.output_field_name, "", NO_FREE);
}

//// ----------------------------------------------------------------
//typedef struct _step_ratio_state_t {
//	double prev;
//	int    have_prev;
//	char*  output_field_name;
//} step_ratio_state_t;
//static void step_ratio_dprocess(double fltv, lrec_t* prec) {
//	double ratio = 1.0;
//	if (this.have_prev) {
//		ratio = fltv / this.prev;
//	} else {
//		this.have_prev = TRUE;
//	}
//	lrec_put(prec, this.output_field_name, mlr_alloc_string_from_double(ratio, MLR_GLOBALS.ofmt),
//		FREE_ENTRY_VALUE);
//	this.prev = fltv;
//}
//static void step_ratio_zprocess(lrec_t* prec) {
//	lrec_put(prec, this.output_field_name, "", NO_FREE);
//}
//static tStepper* step_ratio_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_ratio_state_t* pstate = mlr_malloc_or_die(sizeof(step_ratio_state_t));
//	this.prev          = -999.0;
//	this.have_prev     = FALSE;
//	this.output_field_name = mlr_paste_2_strings(input_field_name, "_ratio");
//
//	stepper->pdprocess_func = step_ratio_dprocess;
//	stepper->pnprocess_func = nil;
//	stepper->psprocess_func = nil;
//	stepper->pzprocess_func = step_ratio_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_rsum_state_t {
//	mv_t   rsum;
//	char*  output_field_name;
//	int    allow_int_float;
//} step_rsum_state_t;
//static void step_rsum_nprocess(mv_t* pnumv, lrec_t* prec) {
//	this.rsum = x_xx_plus_func(&this.rsum, pnumv);
//	lrec_put(prec, this.output_field_name, mv_alloc_format_val(&this.rsum),
//		FREE_ENTRY_VALUE);
//}
//static void step_rsum_zprocess(lrec_t* prec) {
//	lrec_put(prec, this.output_field_name, "", NO_FREE);
//}
//static tStepper* step_rsum_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_rsum_state_t* pstate = mlr_malloc_or_die(sizeof(step_rsum_state_t));
//	this.allow_int_float = allow_int_float;
//	this.rsum = this.allow_int_float ? types.MlrvalFromInt64(0) : MlrvalFromFloat64(0.0);
//	this.output_field_name = mlr_paste_2_strings(input_field_name, "_rsum");
//	stepper->pdprocess_func = nil;
//	stepper->pnprocess_func = step_rsum_nprocess;
//	stepper->psprocess_func = nil;
//	stepper->pzprocess_func = step_rsum_zprocess;
//	return stepper;
//}

//// ----------------------------------------------------------------
//typedef struct _step_counter_state_t {
//	mv_t counter;
//	mv_t one;
//	char*  output_field_name;
//} step_counter_state_t;
//static void step_counter_sprocess(char* strv, lrec_t* prec) {
//	this.counter = x_xx_plus_func(&this.counter, &this.one);
//	lrec_put(prec, this.output_field_name, mv_alloc_format_val(&this.counter),
//		FREE_ENTRY_VALUE);
//}
//static void step_counter_zprocess(lrec_t* prec) {
//	lrec_put(prec, this.output_field_name, "", NO_FREE);
//}
//static tStepper* step_counter_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
//	tStepper* stepper = mlr_malloc_or_die(sizeof(tStepper));
//	step_counter_state_t* pstate = mlr_malloc_or_die(sizeof(step_counter_state_t));
//	this.counter = allow_int_float ? types.MlrvalFromInt64(0) : MlrvalFromFloat64(0.0);
//	this.one     = allow_int_float ? types.MlrvalFromInt64(1) : MlrvalFromFloat64(1.0);
//	this.output_field_name = mlr_paste_2_strings(input_field_name, "_counter");
//
//	stepper->pdprocess_func = nil;
//	stepper->pnprocess_func = nil;
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
//static void step_ewma_dprocess(double fltv, lrec_t* prec) {
//	if (!this.have_prevs) {
//		for (int i = 0; i < this.num_alphas; i++) {
//			lrec_put(prec, this.output_field_names[i], mlr_alloc_string_from_double(fltv, MLR_GLOBALS.ofmt),
//				FREE_ENTRY_VALUE);
//			this.prevs[i] = fltv;
//		}
//		this.have_prevs = TRUE;
//	} else {
//		for (int i = 0; i < this.num_alphas; i++) {
//			double curr = fltv;
//			curr = this.alphas[i] * curr + this.alphacompls[i] * this.prevs[i];
//			lrec_put(prec, this.output_field_names[i], mlr_alloc_string_from_double(curr, MLR_GLOBALS.ofmt),
//				FREE_ENTRY_VALUE);
//			this.prevs[i] = curr;
//		}
//	}
//}
//static void step_ewma_zprocess(lrec_t* prec) {
//	for (int i = 0; i < this.num_alphas; i++)
//		lrec_put(prec, this.output_field_names[i], "", NO_FREE);
//}

//static tStepper* step_ewma_alloc(char* input_field_name, int unused, slls_t* pstring_alphas, slls_t* pewma_suffixes) {
//	tStepper* stepper              = mlr_malloc_or_die(sizeof(tStepper));
//
//	step_ewma_state_t* pstate  = mlr_malloc_or_die(sizeof(step_ewma_state_t));
//	int n                      = pstring_alphas->length;
//	this.num_alphas         = n;
//	this.alphas             = mlr_malloc_or_die(n * sizeof(double));
//	this.alphacompls        = mlr_malloc_or_die(n * sizeof(double));
//	this.prevs              = mlr_malloc_or_die(n * sizeof(double));
//	this.have_prevs         = FALSE;
//	this.output_field_names = mlr_malloc_or_die(n * sizeof(char*));
//	slls_t* psuffixes = (pewma_suffixes == nil) ? pstring_alphas : pewma_suffixes;
//	sllse_t* pe = pstring_alphas->phead;
//	sllse_t* pf = psuffixes->phead;
//	for (int i = 0; i < n; i++, pe = pe->pnext, pf = pf->pnext) {
//		char* string_alpha     = pe->value;
//		char* suffix           = pf->value;
//		this.alphas[i]      = mlr_double_from_string_or_die(string_alpha);
//		this.alphacompls[i] = 1.0 - this.alphas[i];
//		this.prevs[i]       = 0.0;
//		this.output_field_names[i] = mlr_paste_3_strings(input_field_name, "_ewma_", suffix);
//	}
//	this.have_prevs = FALSE;
//
//	stepper->pdprocess_func = step_ewma_dprocess;
//	stepper->pnprocess_func = nil;
//	stepper->psprocess_func = nil;
//	stepper->pzprocess_func = step_ewma_zprocess;
//	return stepper;
//}
