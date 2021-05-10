package transformers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

const DEFAULT_STRING_ALPHA = "0.5"

// ----------------------------------------------------------------
const verbNameStep = "step"

var StepSetup = transforming.TransformerSetup{
	Verb:         verbNameStep,
	UsageFunc:    transformerStepUsage,
	ParseCLIFunc: transformerStepParseCLI,
	IgnoresInput: false,
}

func transformerStepUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "Computes values dependent on the previous record, optionally grouped by category.\n")
	fmt.Fprintf(o, "Options:\n")

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
	fmt.Fprintf(o, "           \"%s %s -a ewma -f sys_load -d 0.01,0.1,0.9\". Default if omitted\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "           is \"-d %s\".\n", DEFAULT_STRING_ALPHA)

	fmt.Fprintf(o, "-o {a,b,c} Custom suffixes for EWMA output fields. If omitted, these default to\n")
	fmt.Fprintf(o, "           the -d values. If supplied, the number of -o values must be the same\n")
	fmt.Fprintf(o, "           as the number of -d values.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s -a rsum -f request_size\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "  %s %s -a delta -f request_size -g hostname\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -f x,y\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y\n", lib.MlrExeName(), verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name\n", lib.MlrExeName(), verbNameStep)

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#filter or\n")
	fmt.Fprintf(o, "https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average\n")
	fmt.Fprintf(o, "for more information on EWMA.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerStepParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var stepperNames []string = nil
	var valueFieldNames []string = nil
	var groupByFieldNames []string = nil
	var stringAlphas []string = nil
	var ewmaSuffixes []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerStepUsage(os.Stdout, true, 0)

		} else if opt == "-a" {
			stepperNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-d" {
			stringAlphas = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			ewmaSuffixes = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-F" {
			// As of Miller 6 this happens automatically, but the flag is accepted
			// as a no-op for backward compatibility with Miller 5 and below.

		} else {
			transformerStepUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerStep(
		stepperNames,
		valueFieldNames,
		groupByFieldNames,
		stringAlphas,
		ewmaSuffixes,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerStep struct {
	// INPUT
	stepperNames      []string
	valueFieldNames   []string
	groupByFieldNames []string
	stringAlphas      []string
	ewmaSuffixes      []string

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
	ewmaSuffixes []string,
) (*TransformerStep, error) {

	if len(stepperNames) == 0 || len(valueFieldNames) == 0 {
		return nil, errors.New(
			// TODO: parameterize verb here somehow
			"mlr step: -a and -f are both required arguments.",
		)
	}
	if len(stringAlphas) != 0 && len(ewmaSuffixes) != 0 {
		if len(ewmaSuffixes) != len(stringAlphas) {
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
		ewmaSuffixes:      ewmaSuffixes,
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
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext
		return
	}

	inrec := inrecAndContext.Record

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
					this.ewmaSuffixes,
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
	ewmaSuffixes []string,
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
	ewmaSuffixes []string,
) tStepper {
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		if stepperLookup.name == stepperName {
			return stepperLookup.stepperAllocator(
				inputFieldName,
				stringAlphas,
				ewmaSuffixes,
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
	delta := types.MlrvalPointerFromInt(0)
	if this.previous != nil {
		delta = types.MlrvalBinaryMinus(valueFieldValue, this.previous)
	}
	inrec.PutCopy(this.outputFieldName, delta)

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
		shift := types.MLRVAL_VOID
		inrec.PutCopy(this.outputFieldName, shift)
	} else {
		inrec.PutCopy(this.outputFieldName, this.previous)
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
	fromFirst := types.MlrvalPointerFromInt(0)
	if this.first == nil {
		this.first = valueFieldValue.Copy()
	} else {
		fromFirst = types.MlrvalBinaryMinus(valueFieldValue, this.first)
	}
	inrec.PutCopy(this.outputFieldName, fromFirst)
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
	ratio := types.MlrvalPointerFromInt(1)
	if this.previous != nil {
		ratio = types.MlrvalDivide(valueFieldValue, this.previous)
	}
	inrec.PutCopy(this.outputFieldName, ratio)

	this.previous = valueFieldValue.Copy()
}

// ================================================================
type tStepperRsum struct {
	rsum            *types.Mlrval
	outputFieldName string
}

func stepperRsumAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperRsum{
		rsum:            types.MlrvalPointerFromInt(0),
		outputFieldName: inputFieldName + "_rsum",
	}
}

func (this *tStepperRsum) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	this.rsum = types.MlrvalBinaryPlus(valueFieldValue, this.rsum)
	inrec.PutCopy(this.outputFieldName, this.rsum)
}

// ================================================================
type tStepperCounter struct {
	counter         *types.Mlrval
	one             *types.Mlrval
	outputFieldName string
}

func stepperCounterAlloc(
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) tStepper {
	return &tStepperCounter{
		counter:         types.MlrvalPointerFromInt(0),
		one:             types.MlrvalPointerFromInt(1),
		outputFieldName: inputFieldName + "_counter",
	}
}

func (this *tStepperCounter) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	this.counter = types.MlrvalBinaryPlus(this.counter, this.one)
	inrec.PutCopy(this.outputFieldName, this.counter)
}

// ----------------------------------------------------------------
// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average

// ================================================================
type tStepperEWMA struct {
	alphas           []types.Mlrval
	oneMinusAlphas   []types.Mlrval
	prevs            []*types.Mlrval
	outputFieldNames []string
	havePrevs        bool
}

func stepperEWMAAlloc(
	inputFieldName string,
	stringAlphas []string,
	ewmaSuffixes []string,
) tStepper {

	// We trust our caller has already checked len(stringAlphas) ==
	// len(ewmaSuffixes) in the CLI parser.
	n := len(stringAlphas)

	alphas := make([]types.Mlrval, n)
	oneMinusAlphas := make([]types.Mlrval, n)
	prevs := make([]*types.Mlrval, n)
	outputFieldNames := make([]string, n)

	suffixes := stringAlphas
	if len(ewmaSuffixes) != 0 {
		suffixes = ewmaSuffixes
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
		prevs[i] = types.MlrvalPointerFromFloat64(0.0)
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
			inrec.PutCopy(this.outputFieldNames[i], valueFieldValue)
			this.prevs[i] = valueFieldValue.Copy()
		}
		this.havePrevs = true
	} else {
		for i, _ := range this.alphas {
			curr := valueFieldValue.Copy()
			// xxx pending pointer-output refactor
			product1 := types.MlrvalTimes(curr, &this.alphas[i])
			product2 := types.MlrvalTimes(this.prevs[i], &this.oneMinusAlphas[i])
			next := types.MlrvalBinaryPlus(product1, product2)
			inrec.PutCopy(this.outputFieldNames[i], next)
			this.prevs[i] = next
		}
	}
}
