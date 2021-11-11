package transformers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

const DEFAULT_STRING_ALPHA = "0.5"

// ----------------------------------------------------------------
const verbNameStep = "step"

var StepSetup = TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameStep)
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
	fmt.Fprintf(o, "           \"%s %s -a ewma -f sys_load -d 0.01,0.1,0.9\". Default if omitted\n", "mlr", verbNameStep)
	fmt.Fprintf(o, "           is \"-d %s\".\n", DEFAULT_STRING_ALPHA)

	fmt.Fprintf(o, "-o {a,b,c} Custom suffixes for EWMA output fields. If omitted, these default to\n")
	fmt.Fprintf(o, "           the -d values. If supplied, the number of -o values must be the same\n")
	fmt.Fprintf(o, "           as the number of -d values.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  %s %s -a rsum -f request_size\n", "mlr", verbNameStep)
	fmt.Fprintf(o, "  %s %s -a delta -f request_size -g hostname\n", "mlr", verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -f x,y\n", "mlr", verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y\n", "mlr", verbNameStep)
	fmt.Fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name\n", "mlr", verbNameStep)

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
	_ *cli.TOptions,
) IRecordTransformer {

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
			stepperNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-d" {
			stringAlphas = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			ewmaSuffixes = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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

	tr := &TransformerStep{
		stepperNames:      stepperNames,
		valueFieldNames:   valueFieldNames,
		groupByFieldNames: groupByFieldNames,
		stringAlphas:      stringAlphas,
		ewmaSuffixes:      ewmaSuffixes,
		groups:            make(map[string]map[string]map[string]tStepper),
	}

	return tr, nil
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

func (tr *TransformerStep) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext
		return
	}

	inrec := inrecAndContext.Record

	// Group-by field names are ["a", "b"]
	// Input data {"a": "s", "b": "t", "x": 3.4, "y": 5.6}
	// Grouping key is "s,t"
	groupingKey, gok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
	if !gok { // current record doesn't have fields to be stepped; pass it along
		outputChannel <- inrecAndContext
		return
	}

	// Create the data structure on first reference
	groupToAccField := tr.groups[groupingKey]
	if groupToAccField == nil {
		// Populate the groups data structure on first reference if needed
		groupToAccField = make(map[string]map[string]tStepper)
		tr.groups[groupingKey] = groupToAccField
	}

	// [3.4, 5.6]
	valueFieldValues, _ := inrec.ReferenceSelectedValues(tr.valueFieldNames)

	// for x=3.4 and y=5.6:
	for i, valueFieldName := range tr.valueFieldNames {
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
		for _, stepperName := range tr.stepperNames {
			stepper, present := accFieldToAccState[stepperName]
			if !present {
				stepper = allocateStepper(
					stepperName,
					valueFieldName,
					tr.stringAlphas,
					tr.ewmaSuffixes,
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

func (stepper *tStepperDelta) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if valueFieldValue.IsEmpty() {
		inrec.PutCopy(stepper.outputFieldName, types.MLRVAL_VOID)
		return
	}

	delta := types.MlrvalFromInt(0)
	if stepper.previous != nil {
		delta = types.BIF_minus_binary(valueFieldValue, stepper.previous)
	}
	inrec.PutCopy(stepper.outputFieldName, delta)

	stepper.previous = valueFieldValue.Copy()
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

func (stepper *tStepperShift) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if stepper.previous == nil {
		shift := types.MLRVAL_VOID
		inrec.PutCopy(stepper.outputFieldName, shift)
	} else {
		inrec.PutCopy(stepper.outputFieldName, stepper.previous)
		stepper.previous = valueFieldValue.Copy()
	}
	stepper.previous = valueFieldValue.Copy()
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

func (stepper *tStepperFromFirst) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	fromFirst := types.MlrvalFromInt(0)
	if stepper.first == nil {
		stepper.first = valueFieldValue.Copy()
	} else {
		fromFirst = types.BIF_minus_binary(valueFieldValue, stepper.first)
	}
	inrec.PutCopy(stepper.outputFieldName, fromFirst)
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

func (stepper *tStepperRatio) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if valueFieldValue.IsEmpty() {
		inrec.PutCopy(stepper.outputFieldName, types.MLRVAL_VOID)
		return
	}

	ratio := types.MlrvalFromInt(1)
	if stepper.previous != nil {
		ratio = types.BIF_divide(valueFieldValue, stepper.previous)
	}
	inrec.PutCopy(stepper.outputFieldName, ratio)

	stepper.previous = valueFieldValue.Copy()
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
		rsum:            types.MlrvalFromInt(0),
		outputFieldName: inputFieldName + "_rsum",
	}
}

func (stepper *tStepperRsum) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if valueFieldValue.IsEmpty() {
		inrec.PutCopy(stepper.outputFieldName, types.MLRVAL_VOID)
	} else {
		stepper.rsum = types.BIF_plus_binary(valueFieldValue, stepper.rsum)
		inrec.PutCopy(stepper.outputFieldName, stepper.rsum)
	}
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
		counter:         types.MlrvalFromInt(0),
		one:             types.MlrvalFromInt(1),
		outputFieldName: inputFieldName + "_counter",
	}
}

func (stepper *tStepperCounter) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if valueFieldValue.IsEmpty() {
		inrec.PutCopy(stepper.outputFieldName, types.MLRVAL_VOID)
	} else {
		stepper.counter = types.BIF_plus_binary(stepper.counter, stepper.one)
		inrec.PutCopy(stepper.outputFieldName, stepper.counter)
	}
}

// ----------------------------------------------------------------
// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average

// ================================================================
type tStepperEWMA struct {
	alphas           []*types.Mlrval
	oneMinusAlphas   []*types.Mlrval
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

	alphas := make([]*types.Mlrval, n)
	oneMinusAlphas := make([]*types.Mlrval, n)
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

func (stepper *tStepperEWMA) process(
	valueFieldValue *types.Mlrval,
	inrec *types.Mlrmap,
) {
	if !stepper.havePrevs {
		for i := range stepper.alphas {
			inrec.PutCopy(stepper.outputFieldNames[i], valueFieldValue)
			stepper.prevs[i] = valueFieldValue.Copy()
		}
		stepper.havePrevs = true
	} else {
		for i := range stepper.alphas {
			curr := valueFieldValue.Copy()
			// xxx pending pointer-output refactor
			product1 := types.BIF_times(curr, stepper.alphas[i])
			product2 := types.BIF_times(stepper.prevs[i], stepper.oneMinusAlphas[i])
			next := types.BIF_plus_binary(product1, product2)
			inrec.PutCopy(stepper.outputFieldNames[i], next)
			stepper.prevs[i] = next
		}
	}
}
