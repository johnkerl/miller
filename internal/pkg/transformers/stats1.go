package transformers

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/transformers/utils"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameStats1 = "stats1"

var Stats1Setup = TransformerSetup{
	Verb:         verbNameStats1,
	UsageFunc:    transformerStats1Usage,
	ParseCLIFunc: transformerStats1ParseCLI,
	IgnoresInput: false,
}

func transformerStats1Usage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameStats1)
	fmt.Fprint(o,
		`Computes univariate statistics for one or more given fields, accumulated across
the input record stream.
Options:
-a {sum,count,...} Names of accumulators: one or more of:
  median   This is the same as p50
  p10 p25.2 p50 p98 p100 etc.
  TODO: flags for interpolated percentiles
`)
	utils.ListStats1Accumulators(o)
	fmt.Fprint(o, `
-f {a,b,c}     Value-field names on which to compute statistics
--fr {regex}   Regex for value-field names on which to compute statistics
               (compute statistics on values in all field names matching regex
--fx {regex}   Inverted regex for value-field names on which to compute statistics
               (compute statistics on values in all field names not matching regex)

-g {d,e,f}     Optional group-by-field names
--gr {regex}   Regex for optional group-by-field names
               (group by values in field names matching regex)
--gx {regex}   Inverted regex for optional group-by-field names
               (group by values in field names not matching regex)

--grfx {regex} Shorthand for --gr {regex} --fx {that same regex}

-i             Use interpolated percentiles, like R's type=7; default like type=1.
               Not sensical for string-valued fields.\n");
-s             Print iterative stats. Useful in tail -f contexts (in which
               case please avoid pprint-format output since end of input
               stream will never be seen).
-h|--help      Show this message.
[TODO: more]
`)

	fmt.Fprintln(o,
		"Example: mlr stats1 -a min,p10,p50,p90,max -f value -g size,shape\n", "mlr", verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size\n", "mlr", verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size -g shape\n", "mlr", verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode --fr '^[a-h].*$' -gr '^k.*$'\n", "mlr", verbNameStats1)
	fmt.Fprintln(o,
		`        This computes count and mode statistics on all field names beginning
         with a through h, grouped by all field names starting with k.`)
	fmt.Println()
	fmt.Fprint(o,
		`Notes:
* p50 and median are synonymous.
* min and max output the same results as p0 and p100, respectively, but use
  less memory.
* String-valued data make sense unless arithmetic on them is required,
  e.g. for sum, mean, interpolated percentiles, etc. In case of mixed data,
  numbers are less than strings.
* count and mode allow text input; the rest require numeric input.
  In particular, 1 and 1.0 are distinct text for count and mode.
* When there are mode ties, the first-encountered datum wins.
`)

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerStats1ParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	accumulatorNameList := make([]string, 0)
	valueFieldNameList := make([]string, 0)
	groupByFieldNameList := make([]string, 0)

	doRegexValueFieldNames := false
	doRegexGroupByFieldNames := false
	invertRegexValueFieldNames := false
	invertRegexGroupByFieldNames := false

	doInterpolatedPercentiles := false
	doIterativeStats := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerStats1Usage(os.Stdout, true, 0)

		} else if opt == "-a" {
			accumulatorNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--fr" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexValueFieldNames = true

		} else if opt == "--fx" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexValueFieldNames = true
			invertRegexValueFieldNames = true
		} else if opt == "--gr" {
			groupByFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexGroupByFieldNames = true
		} else if opt == "--gx" {
			groupByFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			doRegexGroupByFieldNames = true
			invertRegexGroupByFieldNames = true

		} else if opt == "--grfx" {
			doRegexValueFieldNames = true
			doRegexGroupByFieldNames = true
			invertRegexValueFieldNames = true
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)
			groupByFieldNameList = lib.CopyStringArray(valueFieldNameList)

		} else if opt == "-i" {
			doInterpolatedPercentiles = true

		} else if opt == "-s" {
			doIterativeStats = true

		} else if opt == "-S" {
			// No-op pass-through for backward compatibility with Miller 5

		} else if opt == "-F" {
			// No-op pass-through for backward compatibility with Miller 5

		} else {
			transformerStats1Usage(os.Stderr, true, 1)
		}
	}

	// TODO: libify for use across verbs.
	if len(accumulatorNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -a option is required.\n", "mlr", verbNameStats1)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", "mlr", verbNameStats1)
		os.Exit(1)
	}
	if len(valueFieldNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -f option is required.\n", "mlr", verbNameStats1)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", "mlr", verbNameStats1)
		os.Exit(1)
	}

	transformer, err := NewTransformerStats1(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,

		doRegexValueFieldNames,
		doRegexGroupByFieldNames,
		invertRegexValueFieldNames,
		invertRegexGroupByFieldNames,

		doInterpolatedPercentiles,
		doIterativeStats,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerStats1 struct {
	// Input:
	accumulatorNameList  []string
	valueFieldNameList   []string
	groupByFieldNameList []string

	// If the group-by field names are non-regexed, these are just the names in
	// the groupByFieldNameList. If the group-by field names are regexed, this
	// is the union of all the group-by field names encountered in the input,
	// over all records.
	groupByFieldNamesForOutput *lib.OrderedMap

	valueFieldRegexes   []*regexp.Regexp
	groupByFieldRegexes []*regexp.Regexp

	doRegexValueFieldNames   bool
	doRegexGroupByFieldNames bool

	invertRegexValueFieldNames   bool
	invertRegexGroupByFieldNames bool

	doInterpolatedPercentiles bool
	doIterativeStats          bool

	// State:
	accumulatorFactory *utils.Stats1AccumulatorFactory

	// Accumulators are indexed by
	//   groupByFieldName -> valueFieldName -> accumulatorName -> accumulator object
	// This would be
	//   namedAccumulators map[string]map[string]map[string]Stats1NamedAccumulator
	// except we need maps that preserve insertion order.
	namedAccumulators *lib.OrderedMap

	// map[string]OrderedMap[string]*types.Mlrval
	groupingKeysToGroupByFieldValues map[string]*lib.OrderedMap
}

// Given: accumulate count,sum on values x,y group by a,b.
//
// Example input:       Example output:
//   a b x y            a b x_count x_sum y_count y_sum
//   s t 1 2            s t 2       6     2       8
//   u v 3 4            u v 1       3     1       4
//   s t 5 6            u w 1       7     1       9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   "s,t" : {                <--- group-by field names
//     "x" : {                  <--- value field name
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
//   "u,v" : {
//     "x" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
//   "u,w" : {
//     "x" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//     "y" : {
//       "count" : Stats1CountAccumulator object,
//       "sum"   : Stats1SumAccumulator  object
//     },
//   },
// }

func NewTransformerStats1(
	accumulatorNameList []string,
	valueFieldNameList []string,
	groupByFieldNameList []string,

	doRegexValueFieldNames bool,
	doRegexGroupByFieldNames bool,
	invertRegexValueFieldNames bool,
	invertRegexGroupByFieldNames bool,

	doInterpolatedPercentiles bool,
	doIterativeStats bool,
) (*TransformerStats1, error) {
	for _, name := range accumulatorNameList {
		if !utils.ValidateStats1AccumulatorName(name) {
			return nil, errors.New(
				fmt.Sprintf(
					"%s stats1: accumulator \"%s\" not found.\n",
					"mlr", name,
				),
			)
		}
	}

	tr := &TransformerStats1{
		accumulatorNameList:        accumulatorNameList,
		valueFieldNameList:         valueFieldNameList,
		groupByFieldNameList:       groupByFieldNameList,
		groupByFieldNamesForOutput: lib.NewOrderedMap(),

		doRegexValueFieldNames:       doRegexValueFieldNames,
		doRegexGroupByFieldNames:     doRegexGroupByFieldNames,
		invertRegexValueFieldNames:   invertRegexValueFieldNames,
		invertRegexGroupByFieldNames: invertRegexGroupByFieldNames,

		doInterpolatedPercentiles:        doInterpolatedPercentiles,
		doIterativeStats:                 doIterativeStats,
		accumulatorFactory:               utils.NewStats1AccumulatorFactory(),
		namedAccumulators:                lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: make(map[string]*lib.OrderedMap),
	}

	if doRegexGroupByFieldNames {
		tr.groupByFieldRegexes = lib.CompileMillerRegexesOrDie(groupByFieldNameList)
	} else {
		for _, groupByFieldName := range groupByFieldNameList {
			tr.groupByFieldNamesForOutput.Put(groupByFieldName, true)
		}
	}

	if doRegexValueFieldNames {
		tr.valueFieldRegexes = lib.CompileMillerRegexesOrDie(valueFieldNameList)
	}

	return tr, nil
}

// Transform is the function executed for every input record, as well as for
// the end-of-stream marker.
func (tr *TransformerStats1) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		tr.handleInputRecord(inrecAndContext, outputChannel)
	} else {
		tr.handleEndOfRecordStream(inrecAndContext, outputChannel)
	}
}

func (tr *TransformerStats1) handleInputRecord(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
	// then groupingKey is the string "circle,blue".
	var groupingKey string
	var groupByFieldValues *lib.OrderedMap // OrderedMap[string]*types.Mlrval
	var ok bool
	if tr.doRegexGroupByFieldNames {
		groupingKey, groupByFieldValues, ok = tr.getGroupByFieldNamesWithRegexes(inrec)
	} else {
		groupingKey, groupByFieldValues, ok = tr.getGroupByFieldNamesWithoutRegexes(inrec)
	}
	if !ok {
		return
	}

	level2 := tr.namedAccumulators.Get(groupingKey)
	if level2 == nil {
		level2 = lib.NewOrderedMap()
		tr.namedAccumulators.Put(groupingKey, level2)
		// E.g. if grouping by "color" and "shape", and the current record has
		// color=blue, shape=circle, then groupByFieldValues is the map
		// {"color": "blue", "shape": "circle"}.
		tr.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	}

	if tr.doRegexValueFieldNames {
		tr.ingestWithValueFieldRegexes(inrec, groupingKey, level2.(*lib.OrderedMap))
	} else {
		tr.ingestWithoutValueFieldRegexes(inrec, groupingKey, level2.(*lib.OrderedMap))
	}

	if tr.doIterativeStats {
		tr.emitIntoOutputRecord(
			inrecAndContext.Record,
			groupByFieldValues,
			level2.(*lib.OrderedMap),
			inrec,
		)
		outputChannel <- inrecAndContext
	}
}

// E.g. if grouping by "a" and "b", and the current record has a=circle,
// b=blue, then groupingKey is the string "circle,blue".  For grouping without
// regexed group-by field names, the group-by field names/values are the same
// on every record.
func (tr *TransformerStats1) getGroupByFieldNamesWithoutRegexes(
	inrec *types.Mlrmap,
) (
	groupingKey string,
	groupByFieldValues *lib.OrderedMap, // OrderedMap[string]*types.Mlrval,
	ok bool,
) {
	var groupByFieldValuesArray []*types.Mlrval
	groupingKey, groupByFieldValuesArray, ok = inrec.GetSelectedValuesAndJoined(tr.groupByFieldNameList)
	if !ok {
		return groupingKey, nil, false
	}
	groupByFieldValues = lib.NewOrderedMap()
	for i, groupByFieldValue := range groupByFieldValuesArray {
		groupByFieldValues.Put(tr.groupByFieldNameList[i], groupByFieldValue)
	}
	return groupingKey, groupByFieldValues, ok
}

// E.g. if grouping by "a" and "b", and the current record has a=circle,
// b=blue, then groupingKey is the string "circle,blue".  For grouping with
// regexed group-by field names, the group-by field names/values may or may not
// be the same on every record.
func (tr *TransformerStats1) getGroupByFieldNamesWithRegexes(
	inrec *types.Mlrmap,
) (
	groupingKey string,
	groupByFieldValues *lib.OrderedMap, // OrderedMap[string]*types.Mlrval,
	ok bool,
) {

	var buffer bytes.Buffer
	groupByFieldValues = lib.NewOrderedMap()
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		groupByFieldName := pe.Key
		if !tr.matchGroupByFieldName(groupByFieldName) {
			continue
		}

		// Remember the union of all encountered group-by field names
		// for output at the end of the record stream.
		tr.groupByFieldNamesForOutput.Put(groupByFieldName, true)

		groupByFieldValue := pe.Value.Copy()
		if !groupByFieldValues.IsEmpty() {
			buffer.WriteString(",")
		}
		buffer.WriteString(groupByFieldValue.String())
		groupByFieldValues.Put(groupByFieldName, groupByFieldValue)
	}
	groupingKey = buffer.String()

	return groupingKey, groupByFieldValues, true
}

func (tr *TransformerStats1) ingestWithoutValueFieldRegexes(
	inrec *types.Mlrmap,
	groupingKey string,
	level2 *lib.OrderedMap,
) {
	for _, valueFieldName := range tr.valueFieldNameList {
		valueFieldValue := inrec.Get(valueFieldName)
		if valueFieldValue == nil {
			continue
		}
		level3 := level2.Get(valueFieldName)
		if level3 == nil {
			level3 = lib.NewOrderedMap()
			level2.Put(valueFieldName, level3)
		}
		for _, accumulatorName := range tr.accumulatorNameList {
			namedAccumulator := level3.(*lib.OrderedMap).Get(accumulatorName)
			if namedAccumulator == nil {
				namedAccumulator = tr.accumulatorFactory.MakeNamedAccumulator(
					accumulatorName,
					groupingKey,
					valueFieldName,
					tr.doInterpolatedPercentiles,
				)
				level3.(*lib.OrderedMap).Put(accumulatorName, namedAccumulator)
			}
			if valueFieldValue.IsVoid() {
				// The accumulator has been initialized with default values;
				// continue here. (If we were to continue outside of this loop
				// we would be failing to construct the accumulator.)
				continue
			}
			namedAccumulator.(*utils.Stats1NamedAccumulator).Ingest(valueFieldValue)
		}
	}
}

func (tr *TransformerStats1) ingestWithValueFieldRegexes(
	inrec *types.Mlrmap,
	groupingKey string,
	level2 *lib.OrderedMap,
) {
	for pe := inrec.Head; pe != nil; pe = pe.Next {
		valueFieldName := pe.Key

		if !tr.matchValueFieldName(valueFieldName) {
			continue
		}

		valueFieldValue := inrec.Get(valueFieldName)
		if valueFieldValue == nil {
			continue
		}
		level3 := level2.Get(valueFieldName)
		if level3 == nil {
			level3 = lib.NewOrderedMap()
			level2.Put(valueFieldName, level3)
		}
		for _, accumulatorName := range tr.accumulatorNameList {
			namedAccumulator := level3.(*lib.OrderedMap).Get(accumulatorName)
			if namedAccumulator == nil {
				namedAccumulator = tr.accumulatorFactory.MakeNamedAccumulator(
					accumulatorName,
					groupingKey,
					valueFieldName,
					tr.doInterpolatedPercentiles,
				)
				level3.(*lib.OrderedMap).Put(accumulatorName, namedAccumulator)
			}
			if valueFieldValue.IsVoid() {
				// The accumulator has been initialized with default values;
				// continue here. (If we were to continue outside of this loop
				// we would be failing to construct the accumulator.)
				continue
			}
			namedAccumulator.(*utils.Stats1NamedAccumulator).Ingest(valueFieldValue)
		}
	}
}

func (tr *TransformerStats1) matchGroupByFieldName(
	groupByFieldName string,
) bool {
	matches := false
	for _, groupByFieldRegex := range tr.groupByFieldRegexes {
		if groupByFieldRegex.MatchString(groupByFieldName) {
			matches = true
			break
		}
	}
	return matches != tr.invertRegexGroupByFieldNames
}

func (tr *TransformerStats1) matchValueFieldName(
	valueFieldName string,
) bool {
	matches := false
	for _, valueFieldRegex := range tr.valueFieldRegexes {
		if valueFieldRegex.MatchString(valueFieldName) {
			matches = true
			break
		}
	}
	return matches != tr.invertRegexValueFieldNames
}

func (tr *TransformerStats1) handleEndOfRecordStream(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if tr.doIterativeStats {
		outputChannel <- inrecAndContext // end-of-stream marker
		return
	}

	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		level2 := pa.Value.(*lib.OrderedMap)
		groupByFieldValues := tr.groupingKeysToGroupByFieldValues[groupingKey]

		newrec := types.NewMlrmapAsRecord()

		tr.emitIntoOutputRecord(
			inrecAndContext.Record,
			groupByFieldValues,
			level2,
			newrec,
		)

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	}

	outputChannel <- inrecAndContext // end-of-stream marker
}

func (tr *TransformerStats1) emitIntoOutputRecord(
	inrec *types.Mlrmap,
	groupByFieldValues *lib.OrderedMap, // OrderedMap[string]*types.Mlrval,
	level2accumulators *lib.OrderedMap,
	outrec *types.Mlrmap,
) {

	for pa := tr.groupByFieldNamesForOutput.Head; pa != nil; pa = pa.Next {
		groupByFieldName := pa.Key
		iValue := groupByFieldValues.Get(groupByFieldName)
		if iValue != nil {
			outrec.PutCopy(groupByFieldName, iValue.(*types.Mlrval))
		}
	}

	for pb := level2accumulators.Head; pb != nil; pb = pb.Next {
		level3 := pb.Value.(*lib.OrderedMap)
		for pc := level3.Head; pc != nil; pc = pc.Next {
			namedAccumulator := pc.Value.(*utils.Stats1NamedAccumulator)
			key, value := namedAccumulator.Emit()
			outrec.PutCopy(key, value)
		}
	}
}
