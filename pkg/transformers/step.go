// Options for the step verb are mostly simple operations involving the previous record and the
// current record, optionally grouped by one or more group-by field names. For example, with input
// data
//
//   $ cat sample.csv
//   shape,count
//   square,10
//   circle,20
//   square,11
//   circle,23
//
// using the delta stepper we have
//
//   $ mlr --csv --from sample.csv step -a delta -f count
//   shape,count,count_delta
//   square,10,0
//   circle,20,10
//   square,11,-9
//   circle,23,12
//
// whereas if we group by shape when we have
//
//   $ mlr --csv --from sample.csv step -a delta -f count -g shape
//   shape,count,count_delta
//   square,10,0
//   circle,20,0
//   square,11,1
//   circle,23,3
//
// This is (rather, was) straightforward until we added the ability to do *forward* operations such
// as shift_lead. Namely:
//
// * If the stepper is shift_lead then output lags input by one, e.g.  we emit the 10th record only
//   after seeing the 11th. Likewise, for sliding-window average with look-forward of 4, we emit the
//   10th record only after seeing the 14th. More generally, if there are multiple steppers
//   specified with -a, then the delay is the max of each stepper's look-forward.
//
// * Then we need to produce output at the end of the record stream -- e.g.  if there are only 20
//   records and we're doing shift_lead, then we'd normally emit the 20th record only when the 21st
//   is received -- but there isn't one.  And we can't use a simple next-is-nil rule for the last
//   record received in the group-by case. For example, if a given record has shape=square and we're
//   grouping by shape, we don't know a priori where in the record stream the next record with
//   shape=square will be -- or if there will be one at all.
//
// * If we keep a simple hashmap from grouping key to delayed records and process that at end of
//   record stream, since Go hashmaps don't preserve insertion order, we'd have non-deterministic
//   output ordering which would frustrate users and would also break automated regression tests.
//   For example, doing shift_lead with the above sample data, the last square and circle record
//   could appear in either order.
//
// * For these reasons we have an ordered hashmap -- basically a mashup of hashmap and doubly linked
//   list -- of all "window" objects per grouping-key.
//
// * The window object is just the current record along with previous/next records as required by a
//   given stepper. The shift_lag stepper keeps the previous and current record; when the 10th
//   record is ingested, the previous is the 9th, and it emits the 10th record with a value from the
//   9th.  The shift_lead stepper has a current and next. When the 11th record is ingested, the
//   'current' is the 10th record and the 'next' is the 11th, and it emits the 10th record with a
//   value from the 11th.
//
// * The ordered hashmap is called a "stepper log" and it has -- in order -- records pointing to the
//   window object for their grouping key.  We don't know a priori when the end of the record stream
//   is so we keep the last n records for each grouping key.  At end of the record stream we process
//   these.

package transformers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/bifs"
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/transformers/utils"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// For EWMA
const DEFAULT_STRING_ALPHA = "0.5"

const verbNameStep = "step"

var stepOptions = []OptionSpec{
	{Flag: "-a", Arg: "{delta,rsum,...}", Type: "enum", Desc: "Names of steppers: comma-separated, one or more of the listed values.", Values: []string{"counter", "delta", "ewma", "from-first", "ratio", "rprod", "rsum", "shift", "shift_lag", "shift_lead", "slwin"}},
	{Flag: "-f", Arg: "{a,b,c}", Type: "csv-list", Desc: "Value-field names on which to compute statistics."},
	{Flag: "-g", Arg: "{d,e,f}", Type: "csv-list", Desc: "Optional group-by-field names."},
	{Flag: "-F", Type: "bool", Desc: "Computes integerable things (e.g. counter) in floating point. As of Miller 6 this happens automatically, but the flag is accepted as a no-op for backward compatibility with Miller 5 and below."},
	{Flag: "-d", Arg: "{x,y,z}", Type: "csv-list", Desc: "Weights for EWMA. 1 means current sample gets all weight (no smoothing), near under 1 is light smoothing, near over 0 is heavy smoothing. Multiple weights may be specified, e.g. \"mlr step -a ewma -f sys_load -d 0.01,0.1,0.9\". Default if omitted is \"-d " + DEFAULT_STRING_ALPHA + "\"."},
	{Flag: "-o", Arg: "{a,b,c}", Type: "csv-list", Desc: "Custom suffixes for EWMA output fields. If omitted, these default to the -d values. If supplied, the number of -o values must be the same as the number of -d values."},
}

var StepSetup = TransformerSetup{
	Verb:         verbNameStep,
	UsageFunc:    transformerStepUsage,
	ParseCLIFunc: transformerStepParseCLI,
	IgnoresInput: false,
	Options:      stepOptions,
}

func transformerStepUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: mlr %s [options]\n", verbNameStep)
	fmt.Fprintf(o, "Computes values dependent on earlier/later records, optionally grouped by category.\n")
	WriteVerbOptions(o, stepOptions)

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Names of steppers for -a, comma-separated, one or more of:\n")
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		fmt.Fprintf(o, "  %-10s %s\n", stepperLookup.name, stepperLookup.desc)
	}

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Examples:\n")
	fmt.Fprintf(o, "  mlr %s -a rsum -f request_size\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a delta -f request_size -g hostname\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a ewma -d 0.1,0.9 -f x,y\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a slwin_9_0,slwin_0_9 -f x\n", verbNameStep)
	fmt.Fprintf(o, "  mlr %s -a shift_lag_12 -f sales\n", verbNameStep)
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "The shift, shift_lag, shift_lead, delta, and ratio steppers accept an\n")
	fmt.Fprintf(o, "optional trailing count: shift_lag_{n} refers n records back, and\n")
	fmt.Fprintf(o, "shift_lead_{n} refers n records forward. The plain forms are equivalent\n")
	fmt.Fprintf(o, "to a count of 1: e.g. shift_lag is the same as shift_lag_1.\n")

	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#filter or\n")
	fmt.Fprintf(o, "https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average\n")
	fmt.Fprintf(o, "for more information on EWMA.\n")
}

func transformerStepParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var stepperInputs []*tStepperInput = nil
	var valueFieldNames []string = nil
	var groupByFieldNames []string = nil
	var stringAlphas []string = nil
	var ewmaSuffixes []string = nil

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		switch opt {
		case "-h", "--help":
			transformerStepUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		case "-a":
			// Let them do '-a delta -a rsum' or '-a delta,rsum'
			stepperNames, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

			for _, stepperName := range stepperNames {
				stepperInput := stepperInputFromName(stepperName)
				if stepperInput == nil {
					if stepperNameHasBadCount(stepperName) {
						return nil, cli.VerbErrorf(
							verb, "stepper \"%s\": count must be a positive integer", stepperName,
						)
					}
					if stepperNameHasNegativeSlwin(stepperName) {
						return nil, cli.VerbErrorf(
							verb, "stepper needed non-negative num-backward & num-forward in %s.", stepperName,
						)
					}
					return nil, cli.VerbErrorf(verb, "stepper \"%s\" not found", stepperName)
				}
				stepperInputs = append(stepperInputs, stepperInput)
			}

		case "-f":
			// Let them do '-f x -f y' or '-f x,y'
			arr, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			valueFieldNames = append(valueFieldNames, arr...)

		case "-g":
			// Let them do '-g a -g b' or '-g a,b'
			arr, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			groupByFieldNames = append(groupByFieldNames, arr...)

		case "-d":
			// Let them do '-d 0.8 -d 0.9' or '-d 0.8,0.9'
			arr, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			stringAlphas = append(stringAlphas, arr...)

		case "-o":
			// Let them do '-o fast -o slow' or '-o fast,slow'
			arr, err := cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}
			ewmaSuffixes = append(ewmaSuffixes, arr...)

		case "-F":
			// As of Miller 6 this happens automatically, but the flag is accepted
			// as a no-op for backward compatibility with Miller 5 and below.

		default:
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerStep(
		stepperInputs,
		valueFieldNames,
		groupByFieldNames,
		stringAlphas,
		ewmaSuffixes,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

// This is the "stepper log" referred to in comments at the top of this file.
type tStepLogEntry struct {
	recordAndContext *types.RecordAndContext
	windowKeeper     *utils.TWindowKeeper
	// Map from value field name to stepper name to stepper.  E.g. with 'mlr step -g a,b -f x,y -a
	// shift_lag,shift_lead', value field names are 'x' and 'y', and stepper names are 'shift_lag'
	// and 'shift_lead'.
	steppers map[string]map[string]tStepper
}

type TransformerStep struct {
	// INPUT

	stepperInputs     []*tStepperInput
	valueFieldNames   []string
	groupByFieldNames []string
	stringAlphas      []string
	ewmaSuffixes      []string

	maxNumRecordsBackward int
	maxNumRecordsForward  int

	// STATE

	// Scratch space used per-record
	// Map from group-by field names to value-field names to stepper name to stepper object.  See
	// the Transform method below for more details.
	groups map[string]map[string]map[string]tStepper
	// Map from group-by field names to window-keeper object.  These keep rows before and after a
	// 'current' row center point for lag/lead computations, etc.
	windowKeepers map[string]*utils.TWindowKeeper

	// Ordered map from stringified pointer to recordAndContext, to *tStepLogEntry,
	// as described in comments at the top of this file.
	log *lib.OrderedMap[*tStepLogEntry]
}

func NewTransformerStep(
	stepperInputs []*tStepperInput,
	valueFieldNames []string,
	groupByFieldNames []string,
	stringAlphas []string,
	ewmaSuffixes []string,
) (*TransformerStep, error) {

	if len(stepperInputs) == 0 || len(valueFieldNames) == 0 {
		return nil, fmt.Errorf("mlr %s: -a and -f are both required arguments", verbNameStep)
	}
	if len(stringAlphas) != 0 && len(ewmaSuffixes) != 0 {
		if len(ewmaSuffixes) != len(stringAlphas) {
			return nil, fmt.Errorf(
				"mlr %s: If -d and -o are provided, their values must have the same length", verbNameStep,
			)
		}
	}

	maxNumRecordsBackward := 0
	maxNumRecordsForward := 0
	for _, stepperInput := range stepperInputs {
		if maxNumRecordsBackward < stepperInput.numRecordsBackward {
			maxNumRecordsBackward = stepperInput.numRecordsBackward
		}
		if maxNumRecordsForward < stepperInput.numRecordsForward {
			maxNumRecordsForward = stepperInput.numRecordsForward
		}
	}

	tr := &TransformerStep{
		stepperInputs:         stepperInputs,
		valueFieldNames:       valueFieldNames,
		groupByFieldNames:     groupByFieldNames,
		stringAlphas:          stringAlphas,
		ewmaSuffixes:          ewmaSuffixes,
		maxNumRecordsBackward: maxNumRecordsBackward,
		maxNumRecordsForward:  maxNumRecordsForward,
		groups:                make(map[string]map[string]map[string]tStepper),
		windowKeepers:         make(map[string]*utils.TWindowKeeper),
		log:                   lib.NewOrderedMap[*tStepLogEntry](),
	}

	return tr, nil
}

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
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) error {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)

	if !inrecAndContext.EndOfStream {
		if err := tr.handleRecord(inrecAndContext, outputRecordsAndContexts); err != nil {
			return err
		}

	} else {
		// As described in comments at the top of this file: process through all delayed-input
		// records for shift_lead, forward-sliding-window, etc. steppers.
		for pe := tr.log.Head; pe != nil; pe = pe.Next {
			logEntry := pe.Value
			// Shift by one -- if 'current' is the 9th record and 'next' is 10th, and there's no
			// 11th, 'current' becomes the 10th and the 'next' becomes nil.
			logEntry.windowKeeper.Ingest(nil)
			tr.handleDrainRecord(logEntry, outputRecordsAndContexts)
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		return nil
	}
	return nil
}

// handleRecord processes records received before the end of the record stream is seen.
// The records emitted here are the ones we can emit now. For example, with shift_lead, if the most
// recent input record is the 11th, then here we're emitting the 10th.  At EOS, we'll drain any
// delayed-input records in the order in which they were received.
func (tr *TransformerStep) handleRecord(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) error {
	inrec := inrecAndContext.Record

	// Group-by field names are ["a", "b"]
	// Input data {"a": "s", "b": "t", "x": 3.4, "y": 5.6}
	// Grouping key is "s,t"
	groupingKey, gok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
	if !gok { // current record doesn't have fields to be stepped; pass it along
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext)
		return nil
	}

	// Create the data structure on first reference
	groupToAccField := tr.groups[groupingKey]
	if groupToAccField == nil {
		// Populate the groups data structure on first reference if needed
		groupToAccField = make(map[string]map[string]tStepper)
		tr.groups[groupingKey] = groupToAccField
	}

	windowKeeper := tr.windowKeepers[groupingKey]
	if windowKeeper == nil {
		windowKeeper = utils.NewWindowKeeper(
			tr.maxNumRecordsBackward,
			tr.maxNumRecordsForward,
		)
		tr.windowKeepers[groupingKey] = windowKeeper
	}
	windowKeeper.Ingest(inrecAndContext)

	// Keep a log of delayed-input records, which we'll drain at end of record stream.
	tr.insertToLog(inrecAndContext, windowKeeper, groupToAccField)

	// E.g. if x=3.4 and y=5.6 then this is [3.4, 5.6]
	valueFieldValues, _ := inrec.ReferenceSelectedValues(tr.valueFieldNames)

	// For x=3.4 and y=5.6:
	for i, valueFieldName := range tr.valueFieldNames {
		valueFieldValue := valueFieldValues[i]
		if valueFieldValue == nil { // not present in the current record
			// Notify any steppers that cache prev values so they don't carry stale state.
			accFieldToAccState := groupToAccField[valueFieldName]
			if accFieldToAccState != nil {
				for _, stepperInput := range tr.stepperInputs {
					if stepper, present := accFieldToAccState[stepperInput.name]; present {
						if sc, ok := stepper.(tStepperPrevCacheClearer); ok {
							sc.clearPrevValue()
						}
					}
				}
			}
			continue
		}

		accFieldToAccState := groupToAccField[valueFieldName]
		if accFieldToAccState == nil {
			accFieldToAccState = make(map[string]tStepper)
			groupToAccField[valueFieldName] = accFieldToAccState
		}

		// for "delta", "rsum":
		for _, stepperInput := range tr.stepperInputs {
			stepper, present := accFieldToAccState[stepperInput.name]
			if !present {
				var err error
				stepper, err = allocateStepper(
					stepperInput,
					valueFieldName,
					tr.stringAlphas,
					tr.ewmaSuffixes,
				)
				if err != nil {
					return err
				}
				if stepper == nil {
					return cli.VerbErrorf(verbNameStep, "stepper allocation failed")
				}
				accFieldToAccState[stepperInput.name] = stepper
			}

			stepper.process(windowKeeper)
		}
	}

	if windowKeeper.Get(0) != nil {
		outrecAndContext := windowKeeper.Get(0).(*types.RecordAndContext)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)
		tr.removeFromLog(outrecAndContext)
	}
	return nil
}

// handleDrainRecord processes records received after the end of the record stream is seen.  The
// records emitted here are the ones we couldn't emit before. For example, with shift_lead, if the
// most recent input record is the 11th, then before EOS we emitted the 10th. Here, we'll drain any
// delayed-input records in the order in which they were received.
func (tr *TransformerStep) handleDrainRecord(
	logEntry *tStepLogEntry,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
) {
	inrecAndContext := logEntry.recordAndContext
	inrec := inrecAndContext.Record
	windowKeeper := logEntry.windowKeeper
	steppers := logEntry.steppers

	// [3.4, 5.6]
	valueFieldValues, _ := inrec.ReferenceSelectedValues(tr.valueFieldNames)

	// for x=3.4 and y=5.6:
	for i, valueFieldName := range tr.valueFieldNames {
		valueFieldValue := valueFieldValues[i]
		if valueFieldValue == nil { // not present in the current record
			// Notify any steppers that cache prev values so they don't carry stale state.
			accFieldToAccState := steppers[valueFieldName]
			if accFieldToAccState != nil {
				for _, stepperInput := range tr.stepperInputs {
					if stepper, present := accFieldToAccState[stepperInput.name]; present {
						if sc, ok := stepper.(tStepperPrevCacheClearer); ok {
							sc.clearPrevValue()
						}
					}
				}
			}
			continue
		}

		accFieldToAccState := steppers[valueFieldName]
		lib.InternalCodingErrorIf(accFieldToAccState == nil)

		// for "delta", "rsum":
		for _, stepperInput := range tr.stepperInputs {
			stepper, present := accFieldToAccState[stepperInput.name]
			lib.InternalCodingErrorIf(!present)
			stepper.process(windowKeeper)
		}
	}

	if windowKeeper.Get(0) != nil {
		outrecAndContext := windowKeeper.Get(0).(*types.RecordAndContext)
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, outrecAndContext)
	}
}

// insertToLog remembers a delayed-input record so we can process it in the order it was received,
// perhaps only after the end of the record stream has been seen.
func (tr *TransformerStep) insertToLog(
	recordAndContext *types.RecordAndContext,
	windowKeeper *utils.TWindowKeeper,
	steppers map[string]map[string]tStepper,
) {
	key := tr.makeLogKey(recordAndContext)
	ientry := tr.log.Get(key)
	lib.InternalCodingErrorIf(ientry != nil)
	tr.log.Put(key, &tStepLogEntry{
		recordAndContext: recordAndContext,
		windowKeeper:     windowKeeper,
		steppers:         steppers,
	})
}

// removeFromLog shifts records out of the log. For example, with shift_lead, we only have
// look-forward of 1, so the log will only have one record per grouping key.
func (tr *TransformerStep) removeFromLog(
	recordAndContext *types.RecordAndContext,
) {
	key := tr.makeLogKey(recordAndContext)
	ientry := tr.log.Get(key)
	lib.InternalCodingErrorIf(ientry == nil)
	tr.log.Remove(key)
}

// makeLogKey stringifies record-and-context pointer for use as a map key for the stepper log.
func (tr *TransformerStep) makeLogKey(
	inrecAndContext *types.RecordAndContext,
) string {
	return fmt.Sprintf("%p", inrecAndContext)
}

// Lookups for individual steppers, like "delta" or "rsum"

type tStepperInputFromName func(
	stepperName string,
) *tStepperInput

type tOwnsPrefix func(
	stepperName string,
) bool

type tStepperAllocator func(
	stepperInput *tStepperInput,
	inputFieldName string,
	stringAlphas []string,
	ewmaSuffixes []string,
) (tStepper, error)

type tStepperInput struct {
	name               string
	numRecordsBackward int
	numRecordsForward  int
}

type tStepper interface {
	process(windowKeeper *utils.TWindowKeeper)
}

// tStepperPrevCacheClearer is an optional interface for steppers that cache the
// previous record's field value. When a record does not contain the target field,
// handleRecord calls clearPrevValue() so that the cached state reflects "no previous
// value", preventing stale values from leaking across absent-field records. This
// also allows these steppers to set numRecordsBackward=0, eliminating the race
// condition that arises when a downstream transformer (running in a separate goroutine)
// modifies a record that the upstream transformer's window keeper still references.
type tStepperPrevCacheClearer interface {
	clearPrevValue()
}

type tStepperLookup struct {
	name                 string
	nameIsVariable       bool
	ownsPrefix           tOwnsPrefix
	stepperInputFromName tStepperInputFromName
	stepperAllocator     tStepperAllocator
	desc                 string
}

var STEPPER_LOOKUP_TABLE = []tStepperLookup{
	{
		name:                 "counter",
		stepperInputFromName: stepperCounterInputFromName,
		stepperAllocator:     stepperCounterAlloc,
		desc:                 "Count instances of field(s) between successive records",
	},
	{
		name:                 "delta",
		nameIsVariable:       true,
		ownsPrefix:           stepperDeltaOwnsName,
		stepperInputFromName: stepperDeltaInputFromName,
		stepperAllocator:     stepperDeltaAlloc,
		desc:                 "Compute differences in field(s) between successive records. Use delta or equivalently delta_1 for the previous record, or delta_{n} for n records back.",
	},
	{
		name:                 "ewma",
		stepperInputFromName: stepperEWMAInputFromName,
		stepperAllocator:     stepperEWMAAlloc,
		desc:                 "Exponentially weighted moving average over successive records",
	},
	{
		name:                 "from-first",
		stepperInputFromName: stepperFromFirstInputFromName,
		stepperAllocator:     stepperFromFirstAlloc,
		desc:                 "Compute differences in field(s) from first record",
	},
	{
		name:                 "ratio",
		nameIsVariable:       true,
		ownsPrefix:           stepperRatioOwnsName,
		stepperInputFromName: stepperRatioInputFromName,
		stepperAllocator:     stepperRatioAlloc,
		desc:                 "Compute ratios in field(s) between successive records. Use ratio or equivalently ratio_1 for the previous record, or ratio_{n} for n records back.",
	},
	{
		name:                 "rprod",
		stepperInputFromName: stepperRprodInputFromName,
		stepperAllocator:     stepperRprodAlloc,
		desc:                 "Compute running products of field(s) between successive records",
	},
	{
		name:                 "rsum",
		stepperInputFromName: stepperRsumInputFromName,
		stepperAllocator:     stepperRsumAlloc,
		desc:                 "Compute running sums of field(s) between successive records",
	},
	{
		name:                 "shift",
		nameIsVariable:       true,
		ownsPrefix:           stepperShiftOwnsName,
		stepperInputFromName: stepperShiftInputFromName,
		stepperAllocator:     stepperShiftAlloc,
		desc:                 "Alias for shift_lag. Use shift or equivalently shift_1 for the previous record, or shift_{n} for n records back.",
	},
	{
		name:                 "shift_lag",
		nameIsVariable:       true,
		ownsPrefix:           stepperShiftLagOwnsName,
		stepperInputFromName: stepperShiftLagInputFromName,
		stepperAllocator:     stepperShiftLagAlloc,
		desc:                 "Include value(s) in field(s) from the previous record, if any. Use shift_lag or equivalently shift_lag_1 for the previous record, or shift_lag_{n} for n records back.",
	},
	{
		name:                 "shift_lead",
		nameIsVariable:       true,
		ownsPrefix:           stepperShiftLeadOwnsName,
		stepperInputFromName: stepperShiftLeadInputFromName,
		stepperAllocator:     stepperShiftLeadAlloc,
		desc:                 "Include value(s) in field(s) from the next record, if any. Use shift_lead or equivalently shift_lead_1 for the next record, or shift_lead_{n} for n records forward.",
	},
	{
		name:                 "slwin",
		nameIsVariable:       true,
		ownsPrefix:           stepperSlwintOwnsPrefix,
		stepperInputFromName: stepperSlwinInputFromName,
		stepperAllocator:     stepperSlwinAlloc,
		desc:                 "Sliding-window averages over m records back and n forward. E.g. slwin_7_2 for 7 back and 2 forward.",
	},
}

func stepperInputFromName(
	name string,
) *tStepperInput {
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		if stepperLookup.nameIsVariable {
			stepperInput := stepperLookup.stepperInputFromName(name)
			if stepperInput != nil {
				return stepperInput
			}
		} else {
			if stepperLookup.name == name {
				return stepperLookup.stepperInputFromName(name)
			}
		}
	}
	return nil
}

func allocateStepper(
	stepperInput *tStepperInput,
	inputFieldName string,
	stringAlphas []string,
	ewmaSuffixes []string,
) (tStepper, error) {
	for _, stepperLookup := range STEPPER_LOOKUP_TABLE {
		if stepperLookup.nameIsVariable {
			if stepperLookup.ownsPrefix(stepperInput.name) {
				return stepperLookup.stepperAllocator(
					stepperInput,
					inputFieldName,
					stringAlphas,
					ewmaSuffixes,
				)
			}
		} else {
			if stepperLookup.name == stepperInput.name {
				return stepperLookup.stepperAllocator(
					stepperInput,
					inputFieldName,
					stringAlphas,
					ewmaSuffixes,
				)
			}
		}
	}
	return nil, nil
}

// parseStepperCount parses stepper names which accept an optional trailing
// number-of-records-back/forward, like "shift_lag" (count 1) or "shift_lag_12"
// (count 12). It returns ok=false if the name is neither the base name nor the
// base name followed by an underscore and a positive integer.
func parseStepperCount(
	stepperName string,
	baseName string,
) (count int, ok bool) {
	if stepperName == baseName {
		return 1, true
	}
	prefix := baseName + "_"
	if !strings.HasPrefix(stepperName, prefix) {
		return 0, false
	}
	suffix := stepperName[len(prefix):]
	if suffix == "" {
		return 0, false
	}
	for _, c := range suffix {
		if c < '0' || c > '9' {
			return 0, false
		}
	}
	n, err := strconv.Atoi(suffix)
	if err != nil || n < 1 {
		return 0, false
	}
	return n, true
}

// stepperCountSuffixBases lists the steppers which accept an optional
// trailing count, e.g. "shift_lag_12". Longer names come first so that e.g.
// "shift_lag_0" reports against base "shift_lag" rather than "shift".
var stepperCountSuffixBases = []string{"shift_lead", "shift_lag", "shift", "delta", "ratio"}

// stepperNameHasBadCount detects stepper names which look like the
// count-suffixed form of a known stepper but whose suffix is not a positive
// integer -- e.g. "delta_0", "delta_-1", "delta_x" -- so the CLI parser can
// give a more specific error than "stepper not found".
func stepperNameHasBadCount(
	stepperName string,
) bool {
	for _, baseName := range stepperCountSuffixBases {
		if strings.HasPrefix(stepperName, baseName+"_") {
			if _, ok := parseStepperCount(stepperName, baseName); !ok {
				return true
			}
		}
	}
	return false
}

// tValueRing is a fixed-size ring buffer of cached field values, used by the
// steppers which reference a value from n records back -- e.g. shift_lag_7,
// delta_7, ratio_7 -- as well as their unparameterized forms (n=1). Steppers
// cache copies of previous values in their own state, rather than reading
// from the window keeper's stored (already-emitted) records, which avoids a
// race condition when downstream transformers run concurrently.
type tValueRing struct {
	values []*mlrval.Mlrval
	next   int
	filled int
}

func newValueRing(n int) *tValueRing {
	return &tValueRing{
		values: make([]*mlrval.Mlrval, n),
		next:   0,
		filled: 0,
	}
}

// push inserts the newest value into the ring, returning the value from n
// records back -- along with whether n records have been seen yet. A nil
// value means the record n back did not have the field in question.
func (ring *tValueRing) push(value *mlrval.Mlrval) (nBack *mlrval.Mlrval, has bool) {
	nBack = ring.values[ring.next]
	has = ring.filled >= len(ring.values)
	ring.values[ring.next] = value
	ring.next = (ring.next + 1) % len(ring.values)
	if !has {
		ring.filled++
	}
	return nBack, has
}

// Implementations of individual steppers, like "delta" or "rsum"

// The stepper name may be "delta" for difference from the previous record, or
// e.g. "delta_7" for difference from 7 records back.
type tStepperDelta struct {
	inputFieldName  string
	outputFieldName string
	prevValues      *tValueRing
}

func stepperDeltaOwnsName(
	stepperName string,
) bool {
	_, ok := parseStepperCount(stepperName, "delta")
	return ok
}

func stepperDeltaInputFromName(
	stepperName string,
) *tStepperInput {
	if _, ok := parseStepperCount(stepperName, "delta"); !ok {
		return nil
	}
	return &tStepperInput{
		name: stepperName,
		// numRecordsBackward is 0: delta caches previous values in its own struct
		// rather than reading from the window keeper's stored (already-emitted) records,
		// which avoids a race condition when downstream transformers run concurrently.
		numRecordsBackward: 0,
		numRecordsForward:  0,
	}
}

func stepperDeltaAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	count, _ := parseStepperCount(stepperInput.name, "delta")
	return &tStepperDelta{
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_" + stepperInput.name,
		prevValues:      newValueRing(count),
	}, nil
}

func (stepper *tStepperDelta) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		stepper.prevValues.push(nil)
		return
	}

	if currval.IsVoid() {
		stepper.prevValues.push(nil)
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
		return
	}

	prevValue, hasPrev := stepper.prevValues.push(currval.Copy())
	delta := mlrval.FromInt(0)
	if hasPrev && prevValue != nil {
		delta = bifs.BIF_minus_binary(currval, prevValue)
	}
	currec.PutCopy(stepper.outputFieldName, delta.Copy())
}

func (stepper *tStepperDelta) clearPrevValue() {
	stepper.prevValues.push(nil)
}

// shift is an alias for shift_lag. The stepper name may be "shift_lag" for
// the value from the previous record, or e.g. "shift_lag_7" for the value
// from 7 records back; likewise "shift" and "shift_7".
type tStepperShiftLag struct {
	inputFieldName  string
	outputFieldName string
	prevValues      *tValueRing
}

func stepperShiftOwnsName(
	stepperName string,
) bool {
	_, ok := parseStepperCount(stepperName, "shift")
	return ok
}

func stepperShiftLagOwnsName(
	stepperName string,
) bool {
	_, ok := parseStepperCount(stepperName, "shift_lag")
	return ok
}

func stepperShiftInputFromName(
	stepperName string,
) *tStepperInput {
	if _, ok := parseStepperCount(stepperName, "shift"); !ok {
		return nil
	}
	return &tStepperInput{
		name: stepperName,
		// numRecordsBackward is 0: shift caches previous values in its own struct
		// rather than reading from the window keeper's stored (already-emitted) records,
		// which avoids a race condition when downstream transformers run concurrently.
		numRecordsBackward: 0,
		numRecordsForward:  0,
	}
}

func stepperShiftLagInputFromName(
	stepperName string,
) *tStepperInput {
	if _, ok := parseStepperCount(stepperName, "shift_lag"); !ok {
		return nil
	}
	return &tStepperInput{
		name: stepperName,
		// numRecordsBackward is 0: shift_lag caches previous values in its own struct
		// rather than reading from the window keeper's stored (already-emitted) records,
		// which avoids a race condition when downstream transformers run concurrently.
		numRecordsBackward: 0,
		numRecordsForward:  0,
	}
}

func stepperShiftAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	count, _ := parseStepperCount(stepperInput.name, "shift")
	return &tStepperShiftLag{
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_" + stepperInput.name,
		prevValues:      newValueRing(count),
	}, nil
}

func stepperShiftLagAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	count, _ := parseStepperCount(stepperInput.name, "shift_lag")
	return &tStepperShiftLag{
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_" + stepperInput.name,
		prevValues:      newValueRing(count),
	}, nil
}

func (stepper *tStepperShiftLag) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record

	curval := currec.Get(stepper.inputFieldName)

	var newValue *mlrval.Mlrval = nil
	if curval != nil {
		newValue = curval.Copy()
	}
	prevValue, hasPrev := stepper.prevValues.push(newValue)

	if !hasPrev || prevValue == nil {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
	} else {
		currec.PutCopy(stepper.outputFieldName, prevValue.Copy())
	}
}

func (stepper *tStepperShiftLag) clearPrevValue() {
	stepper.prevValues.push(nil)
}

// The stepper name may be "shift_lead" for the value from the next record, or
// e.g. "shift_lead_7" for the value from 7 records forward.
type tStepperShiftLead struct {
	inputFieldName    string
	outputFieldName   string
	numRecordsForward int
}

func stepperShiftLeadOwnsName(
	stepperName string,
) bool {
	_, ok := parseStepperCount(stepperName, "shift_lead")
	return ok
}

func stepperShiftLeadInputFromName(
	stepperName string,
) *tStepperInput {
	count, ok := parseStepperCount(stepperName, "shift_lead")
	if !ok {
		return nil
	}
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0,
		numRecordsForward:  count,
	}
}

func stepperShiftLeadAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	return &tStepperShiftLead{
		inputFieldName:    inputFieldName,
		outputFieldName:   inputFieldName + "_" + stepperInput.name,
		numRecordsForward: stepperInput.numRecordsForward,
	}, nil
}

func (stepper *tStepperShiftLead) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record

	inextrec := windowKeeper.Get(stepper.numRecordsForward)
	if inextrec == nil {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
		return
	}
	nextrec := inextrec.(*types.RecordAndContext).Record
	nextval := nextrec.Get(stepper.inputFieldName)

	if nextval != nil {
		currec.PutCopy(stepper.outputFieldName, nextval.Copy())
	}
}

type tStepperFromFirst struct {
	first           *mlrval.Mlrval
	inputFieldName  string
	outputFieldName string
}

func stepperFromFirstInputFromName(
	stepperName string,
) *tStepperInput {
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0, // doesn't use record-windowing; retains its own pointer
		numRecordsForward:  0,
	}
}

func stepperFromFirstAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	return &tStepperFromFirst{
		first:           nil,
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_from_first",
	}, nil
}

func (stepper *tStepperFromFirst) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		return
	}

	fromFirst := mlrval.FromInt(0)
	if stepper.first == nil {
		stepper.first = currval.Copy()
	} else {
		fromFirst = bifs.BIF_minus_binary(currval, stepper.first)
	}
	currec.PutCopy(stepper.outputFieldName, fromFirst)
}

// The stepper name may be "ratio" for ratio against the previous record, or
// e.g. "ratio_7" for ratio against the record 7 records back.
type tStepperRatio struct {
	inputFieldName  string
	outputFieldName string
	prevValues      *tValueRing
}

func stepperRatioOwnsName(
	stepperName string,
) bool {
	_, ok := parseStepperCount(stepperName, "ratio")
	return ok
}

func stepperRatioInputFromName(
	stepperName string,
) *tStepperInput {
	if _, ok := parseStepperCount(stepperName, "ratio"); !ok {
		return nil
	}
	return &tStepperInput{
		name: stepperName,
		// numRecordsBackward is 0: ratio caches previous values in its own struct
		// rather than reading from the window keeper's stored (already-emitted) records,
		// which avoids a race condition when downstream transformers run concurrently.
		numRecordsBackward: 0,
		numRecordsForward:  0,
	}
}

func stepperRatioAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	count, _ := parseStepperCount(stepperInput.name, "ratio")
	return &tStepperRatio{
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_" + stepperInput.name,
		prevValues:      newValueRing(count),
	}, nil
}

func (stepper *tStepperRatio) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		stepper.prevValues.push(nil)
		return
	}

	if currval.IsVoid() {
		stepper.prevValues.push(nil)
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
		return
	}

	prevValue, hasPrev := stepper.prevValues.push(currval.Copy())
	ratio := mlrval.FromInt(1)
	if hasPrev && prevValue != nil {
		ratio = bifs.BIF_divide(currval, prevValue)
	}
	currec.PutCopy(stepper.outputFieldName, ratio.Copy())
}

func (stepper *tStepperRatio) clearPrevValue() {
	stepper.prevValues.push(nil)
}

type tStepperRprod struct {
	rprod           *mlrval.Mlrval
	inputFieldName  string
	outputFieldName string
}

func stepperRprodInputFromName(
	stepperName string,
) *tStepperInput {
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0, // doesn't use record-windowing; retains its own pointer
		numRecordsForward:  0,
	}
}

func stepperRprodAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	return &tStepperRprod{
		rprod:           mlrval.FromInt(1),
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_rprod",
	}, nil
}

func (stepper *tStepperRprod) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		return
	}

	if currval.IsVoid() {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
	} else {
		stepper.rprod = bifs.BIF_times(currval, stepper.rprod)
		currec.PutCopy(stepper.outputFieldName, stepper.rprod)
	}
}

type tStepperRsum struct {
	rsum            *mlrval.Mlrval
	inputFieldName  string
	outputFieldName string
}

func stepperRsumInputFromName(
	stepperName string,
) *tStepperInput {
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0, // doesn't use record-windowing; retains its own pointer
		numRecordsForward:  0,
	}
}

func stepperRsumAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	return &tStepperRsum{
		rsum:            mlrval.FromInt(0),
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_rsum",
	}, nil
}

func (stepper *tStepperRsum) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		return
	}

	if currval.IsVoid() {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
	} else {
		stepper.rsum = bifs.BIF_plus_binary(currval, stepper.rsum)
		currec.PutCopy(stepper.outputFieldName, stepper.rsum)
	}
}

type tStepperCounter struct {
	counter         *mlrval.Mlrval
	inputFieldName  string
	outputFieldName string
}

func stepperCounterInputFromName(
	stepperName string,
) *tStepperInput {
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0, // doesn't use record-windowing; retains its own pointer
		numRecordsForward:  0,
	}
}

func stepperCounterAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	return &tStepperCounter{
		counter:         mlrval.FromInt(0),
		inputFieldName:  inputFieldName,
		outputFieldName: inputFieldName + "_counter",
	}, nil
}

func (stepper *tStepperCounter) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		return
	}

	if currval.IsVoid() {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
	} else {
		stepper.counter = bifs.BIF_plus_binary(stepper.counter, mlrval.ONE)
		currec.PutCopy(stepper.outputFieldName, stepper.counter)
	}
}

// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average

type tStepperEWMA struct {
	alphas           []*mlrval.Mlrval
	oneMinusAlphas   []*mlrval.Mlrval
	prevs            []*mlrval.Mlrval
	inputFieldName   string
	outputFieldNames []string
	havePrevs        bool
}

func stepperEWMAInputFromName(
	stepperName string,
) *tStepperInput {
	return &tStepperInput{
		name:               stepperName,
		numRecordsBackward: 0, // doesn't use record-windowing; retains its own accumulators
		numRecordsForward:  0,
	}
}

func stepperEWMAAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	stringAlphas []string,
	ewmaSuffixes []string,
) (tStepper, error) {

	// We trust our caller has already checked len(stringAlphas) == len(ewmaSuffixes) in the CLI
	// parser.
	n := len(stringAlphas)

	alphas := make([]*mlrval.Mlrval, n)
	oneMinusAlphas := make([]*mlrval.Mlrval, n)
	prevs := make([]*mlrval.Mlrval, n)
	outputFieldNames := make([]string, n)

	suffixes := stringAlphas
	if len(ewmaSuffixes) != 0 {
		suffixes = ewmaSuffixes
	}

	for i, stringAlpha := range stringAlphas {
		suffix := suffixes[i]

		dalpha, ok := lib.TryFloatFromString(stringAlpha)
		if !ok {
			return nil, cli.VerbErrorf(verbNameStep,
				"could not parse \"%s\" as floating-point EWMA coefficient.", stringAlpha)
		}
		alphas[i] = mlrval.FromFloat(dalpha)
		oneMinusAlphas[i] = mlrval.FromFloat(1.0 - dalpha)
		prevs[i] = mlrval.FromFloat(0.0)
		outputFieldNames[i] = inputFieldName + "_ewma_" + suffix
	}

	return &tStepperEWMA{
		alphas:           alphas,
		oneMinusAlphas:   oneMinusAlphas,
		prevs:            prevs,
		inputFieldName:   inputFieldName,
		outputFieldNames: outputFieldNames,
		havePrevs:        false,
	}, nil
}

func (stepper *tStepperEWMA) process(
	windowKeeper *utils.TWindowKeeper,
) {
	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currecAndContext := icur.(*types.RecordAndContext)
	currec := currecAndContext.Record
	currval := currec.Get(stepper.inputFieldName)

	// The record at the window center may lack the field, e.g. when this
	// stepper is combined with a forward-window stepper such as shift_lead,
	// over heterogeneous data.
	if currval == nil {
		return
	}

	if !stepper.havePrevs {
		for i := range stepper.alphas {
			currec.PutCopy(stepper.outputFieldNames[i], currval)
			stepper.prevs[i] = currval.Copy()
		}
		stepper.havePrevs = true
	} else {
		for i := range stepper.alphas {
			curr := currval.Copy()
			next := bifs.BIF_plus_binary(
				bifs.BIF_times(curr, stepper.alphas[i]),
				bifs.BIF_times(stepper.prevs[i], stepper.oneMinusAlphas[i]),
			)
			currec.PutCopy(stepper.outputFieldNames[i], next)
			stepper.prevs[i] = next
		}
	}
}

type tStepperSlwin struct {
	inputFieldName     string
	numRecordsBackward int
	numRecordsForward  int
	outputFieldName    string
}

func stepperSlwintOwnsPrefix(
	stepperName string,
) bool {
	return strings.HasPrefix(stepperName, "slwin")
}

// stepperNameHasNegativeSlwin detects slwin stepper names whose
// num-backward/num-forward parsed but were negative, so the CLI parser can
// give a more specific error than "stepper not found".
func stepperNameHasNegativeSlwin(stepperName string) bool {
	var numRecordsBackward, numRecordsForward int
	n, err := fmt.Sscanf(stepperName, "slwin_%d_%d", &numRecordsBackward, &numRecordsForward)
	return n == 2 && err == nil && (numRecordsBackward < 0 || numRecordsForward < 0)
}

func stepperSlwinInputFromName(
	stepperName string,
) *tStepperInput {
	var numRecordsBackward, numRecordsForward int
	n, err := fmt.Sscanf(stepperName, "slwin_%d_%d", &numRecordsBackward, &numRecordsForward)
	if n == 2 && err == nil {
		if numRecordsBackward < 0 || numRecordsForward < 0 {
			// The CLI parser reports the specific negative-parameter error;
			// see stepperNameHasNegativeSlwin.
			return nil
		}
		return &tStepperInput{
			name:               stepperName,
			numRecordsBackward: numRecordsBackward, // doesn't use record-windowing; retains its own accumulators
			numRecordsForward:  numRecordsForward,
		}
	} else {
		return nil
	}
}

func stepperSlwinAlloc(
	stepperInput *tStepperInput,
	inputFieldName string,
	_unused1 []string,
	_unused2 []string,
) (tStepper, error) {
	nb := stepperInput.numRecordsBackward
	nf := stepperInput.numRecordsForward
	return &tStepperSlwin{
		inputFieldName:     inputFieldName,
		outputFieldName:    fmt.Sprintf("%s_%d_%d", inputFieldName, nb, nf),
		numRecordsBackward: nb,
		numRecordsForward:  nf,
	}, nil
}

func (stepper *tStepperSlwin) process(
	windowKeeper *utils.TWindowKeeper,
) {
	count := int64(0)
	sum := mlrval.FromFloat(0.0)
	for i := -stepper.numRecordsBackward; i <= stepper.numRecordsForward; i++ {
		irac := windowKeeper.Get(i)
		if irac == nil {
			continue
		}
		rac := irac.(*types.RecordAndContext)
		rec := rac.Record
		val := rec.Get(stepper.inputFieldName)
		if val == nil || val.IsVoid() {
			continue
		}
		sum = bifs.BIF_plus_binary(sum, val)
		count++
	}

	icur := windowKeeper.Get(0)
	if icur == nil {
		return
	}
	currac := icur.(*types.RecordAndContext)
	currec := currac.Record

	if count == 0 {
		currec.PutCopy(stepper.outputFieldName, mlrval.VOID)
	} else {
		currec.PutReference(
			stepper.outputFieldName,
			bifs.BIF_divide(sum, mlrval.FromInt(count)),
		)
	}
}
