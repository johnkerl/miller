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

	// As of Miller 6 this happens automatically, but the flag is accepted
	// as a no-op for backward compatibility with Miller 5 and below.
	_ = flagSet.Bool(
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
	fmt.Fprintf(o, "           As of Miller 6 this happens automatically, but the flag is accepted\n")
	fmt.Fprintf(o, "           as a no-op for backward compatibility with Miller 5 and below.\n")

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
				"mlr step: If -d and -o are provided, their values must have the same length.",
			)
		}
	}

	this := &TransformerStep{
		stepperNames:      stepperNames,
		valueFieldNames:   valueFieldNames,
		groupByFieldNames: groupByFieldNames,
		stringAlphas:      stringAlphas,
		EWMASuffixes:      EWMASuffixes,
		groups:            make(map[string]map[string]map[string]tStepper),
	}

	return this, nil
}

// ----------------------------------------------------------------
// Multilevel hashmap structure for the `groups` field example:
//
// * Group-by field names = ["a", "b"]
// * Value field names = ["x", "y"]
// * Steppers ["rsum", "delta"]
//
// {
//   "s,t" : {        <-- for records where 'a=s,b=t'
//     "x": {
//       "rsum": rsum stepper object,
//       "delta": delta stepper object,
//     },
//     "y": {
//       "rsum": rsum stepper object,
//       "delta": delta stepper object,
//     }
//   },
//   "u,v" : {        <-- for records where 'a=u,b=v'
//     "x": {
//       "rsum": rsum stepper object,
//       "delta": delta stepper object,
//     },
//     "y": {
//       "rsum": rsum stepper object,
//       "delta": delta stepper object,
//     }
//   }
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

	// Group-by field names are ["a", "b"]
	// Input data {"a": "s", "b": "t", "x": 3.4, "y": 5.6}
	// Grouping key is "s,t"
	groupingKey, gok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
	if !gok { // current record doesn't have fields to be stepped; pass it along
		outputChannel <- inrecAndContext
		return
	}

	// Create the data structure on first reference
	groupToAccField := this.groups[groupingKey]
	if groupToAccField == nil {
		// Populate the groups data structure on first reference if needed
		groupToAccField = make(map[string]map[string]tStepper)
		this.groups[groupingKey] = groupToAccField
	}

	// [3.4, 5.6]
	valueFieldValues, _ := inrec.ReferenceSelectedValues(this.valueFieldNames)

	// for x=3.4 and y=5.6:
	for i, valueFieldName := range this.valueFieldNames {
		valueFieldValue := valueFieldValues[i]
		if valueFieldValue == nil { // not present in the current record
			continue
		}

		accFieldToAccState := groupToAccField[valueFieldName]
		if accFieldToAccState == nil {
			accFieldToAccState = make(map[string]tStepper)
			groupToAccField[valueFieldName] = accFieldToAccState
		}

		// for "delta", "rsum":
		for _, stepperName := range this.stepperNames {
			stepper, present := accFieldToAccState[stepperName]
			if !present {
				stepper = allocateStepper(
					stepperName,
					valueFieldName,
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
			stepper.process(valueFieldValue, inrec)
		}
	}

	outputChannel <- inrecAndContext
}

// ================================================================
// Lookups for individual steppers, like "delta" or "rsum"

type tStepperAllocator func(
	inputFieldName string,
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
	{"rsum", stepperRsumAlloc, "Compute running sums of field(s) between successive records"},
	{"counter", stepperCounterAlloc, "Count instances of field(s) between successive records"},
	{"ewma", stepperEWMAAlloc, "Exponentially weighted moving average over successive records"},
}

func allocateStepper(
	stepperName string,
	inputFieldName string,
	stringAlphas []string,
	EWMASuffixes []string,
) tStepper {
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		if stepperLookup.name == stepperName {
			return stepperLookup.stepperAllocator(
				inputFieldName,
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
}

func stepperDeltaAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperDelta{
		previous:        nil,
		outputFieldName: inputFieldName + "_delta",
	}
}

func (this *tStepperDelta) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var delta types.Mlrval
	if this.previous == nil {
		delta = types.MlrvalFromInt64(0)
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
}

func stepperShiftAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperShift{
		previous:        nil,
		outputFieldName: inputFieldName + "_shift",
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
}

func stepperFromFirstAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperFromFirst{
		first:           nil,
		outputFieldName: inputFieldName + "_from_first",
	}
}

func (this *tStepperFromFirst) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var from_first types.Mlrval
	if this.first == nil {
		from_first = types.MlrvalFromInt64(0)
		this.first = valueFieldValue.Copy()
	} else {
		from_first = types.MlrvalBinaryMinus(valueFieldValue, this.first)
	}
	inrec.PutCopy(&this.outputFieldName, &from_first)
}

// ================================================================
type tStepperRatio struct {
	previous        *types.Mlrval
	outputFieldName string
}

func stepperRatioAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperRatio{
		previous:        nil,
		outputFieldName: inputFieldName + "_ratio",
	}
}

func (this *tStepperRatio) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	var ratio types.Mlrval
	if this.previous == nil {
		ratio = types.MlrvalFromInt64(1)
	} else {
		ratio = types.MlrvalDivide(valueFieldValue, this.previous)
	}
	inrec.PutCopy(&this.outputFieldName, &ratio)

	this.previous = valueFieldValue.Copy()
}

// ================================================================
type tStepperRsum struct {
	rsum            types.Mlrval
	outputFieldName string
}

func stepperRsumAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperRsum{
		rsum:            types.MlrvalFromInt64(0),
		outputFieldName: inputFieldName + "_rsum",
	}
}

func (this *tStepperRsum) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	this.rsum = types.MlrvalBinaryPlus(valueFieldValue, &this.rsum)
	inrec.PutCopy(&this.outputFieldName, &this.rsum)
}

// ================================================================
type tStepperCounter struct {
	counter         types.Mlrval
	one             types.Mlrval
	outputFieldName string
}

func stepperCounterAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperCounter{
		counter:         types.MlrvalFromInt64(0),
		one:             types.MlrvalFromInt64(1),
		outputFieldName: inputFieldName + "_counter",
	}
}

func (this *tStepperCounter) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	this.counter = types.MlrvalBinaryPlus(&this.counter, &this.one)
	inrec.PutCopy(&this.outputFieldName, &this.counter)
}

// ----------------------------------------------------------------
// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average

// ================================================================
type tStepperEWMA struct {
	alphas           []types.Mlrval
	oneMinusAlphas   []types.Mlrval
	prevs            []types.Mlrval
	outputFieldNames []string
	havePrevs        bool
}

func stepperEWMAAlloc(
	inputFieldName string,
	stringAlphas []string,
	EWMASuffixes []string,
) tStepper {

	// We trust our caller has already checked len(stringAlphas) ==
	// len(EWMASuffixes) in the CLI parser.
	n := len(stringAlphas)

	alphas := make([]types.Mlrval, n)
	oneMinusAlphas := make([]types.Mlrval, n)
	prevs := make([]types.Mlrval, n)
	outputFieldNames := make([]string, n)

	suffixes := stringAlphas
	if len(EWMASuffixes) != 0 {
		suffixes = EWMASuffixes
	}

	for i, stringAlpha := range stringAlphas {
		suffix := suffixes[i]

		dalpha, ok := lib.TryFloat64FromString(stringAlpha)
		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"mlr step: could not parse \"%s\" as floating-point EWMA coefficient.\n",
				stringAlpha,
			)
			os.Exit(1)
		}
		alphas[i] = types.MlrvalFromFloat64(dalpha)
		oneMinusAlphas[i] = types.MlrvalFromFloat64(1.0 - dalpha)
		prevs[i] = types.MlrvalFromFloat64(0.0)
		outputFieldNames[i] = inputFieldName + "_ewma_" + suffix
	}

	return &tStepperEWMA{
		alphas:           alphas,
		oneMinusAlphas:   oneMinusAlphas,
		prevs:            prevs,
		outputFieldNames: outputFieldNames,
		havePrevs:        false,
	}
}

func (this *tStepperEWMA) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if !this.havePrevs {
		for i, _ := range this.alphas {
			inrec.PutCopy(&this.outputFieldNames[i], valueFieldValue)
			this.prevs[i] = *valueFieldValue.Copy()
		}
		this.havePrevs = true
	} else {
		for i, _ := range this.alphas {
			curr := valueFieldValue.Copy()
			product1 := types.MlrvalTimes(curr, &this.alphas[i])
			product2 := types.MlrvalTimes(&this.prevs[i], &this.oneMinusAlphas[i])
			next := types.MlrvalBinaryPlus(&product1, &product2)
			inrec.PutCopy(&this.outputFieldNames[i], &next)
			this.prevs[i] = next
		}
	}
}
