package transformers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transformers/utils"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameStats1 = "stats1"

var Stats1Setup = transforming.TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameStats1)
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
-f {a,b,c}   Value-field names on which to compute statistics
-g {d,e,f}   Optional group-by-field names

-i           Use interpolated percentiles, like R's type=7; default like type=1.\n");
             Not sensical for string-valued fields.\n");
-s           Print iterative stats. Useful in tail -f contexts (in which
             case please avoid pprint-format output since end of input
             stream will never be seen).
-h|--help    Show this message.
[TODO: more]
`)

	fmt.Fprintln(o,
		"Example: mlr stats1 -a min,p10,p50,p90,max -f value -g size,shape\n", lib.MlrExeName(), verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size\n", lib.MlrExeName(), verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size -g shape\n", lib.MlrExeName(), verbNameStats1)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode --fr '^[a-h].*$' -gr '^k.*$'\n", lib.MlrExeName(), verbNameStats1)
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
	_ *cliutil.TOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	accumulatorNameList := make([]string, 0)
	valueFieldNameList := make([]string, 0)
	groupByFieldNameList := make([]string, 0)
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
			accumulatorNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "--fr" {
			// TODO: port field-name regexing from C to Go
			valueFieldNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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
		fmt.Fprintf(os.Stderr, "%s %s: -a option is required.\n", lib.MlrExeName(), verbNameStats1)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", lib.MlrExeName(), verbNameStats1)
		os.Exit(1)
	}
	if len(valueFieldNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -f option is required.\n", lib.MlrExeName(), verbNameStats1)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", lib.MlrExeName(), verbNameStats1)
		os.Exit(1)
	}

	transformer, err := NewTransformerStats1(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,
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
	accumulatorNameList       []string
	valueFieldNameList        []string
	groupByFieldNameList      []string
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

	groupingKeysToGroupByFieldValues map[string][]*types.Mlrval
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
	doInterpolatedPercentiles bool,
	doIterativeStats bool,
) (*TransformerStats1, error) {
	for _, name := range accumulatorNameList {
		if !utils.ValidateStats1AccumulatorName(name) {
			return nil, errors.New(
				fmt.Sprintf(
					"%s stats1: accumulator \"%s\" not found.\n",
					lib.MlrExeName(), name,
				),
			)
		}
	}

	tr := &TransformerStats1{
		accumulatorNameList:              accumulatorNameList,
		valueFieldNameList:               valueFieldNameList,
		groupByFieldNameList:             groupByFieldNameList,
		doInterpolatedPercentiles:        doInterpolatedPercentiles,
		doIterativeStats:                 doIterativeStats,
		accumulatorFactory:               utils.NewStats1AccumulatorFactory(),
		namedAccumulators:                lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: make(map[string][]*types.Mlrval),
	}
	return tr, nil
}

// ----------------------------------------------------------------
func (tr *TransformerStats1) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.handleInputRecord(inrecAndContext, outputChannel)
	} else {
		tr.handleEndOfRecordStream(inrecAndContext, outputChannel)
	}
}

// ----------------------------------------------------------------
func (tr *TransformerStats1) handleInputRecord(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// TODO: make a function-pointer variant for non-iterative which doesn't get the
	// unnecessary groupByFieldValues.

	// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
	// then groupingKey is the string "circle,blue".
	groupingKey, groupByFieldValues, ok := inrec.GetSelectedValuesAndJoined(
		tr.groupByFieldNameList,
	)
	if !ok {
		return
	}

	level2 := tr.namedAccumulators.Get(groupingKey)
	if level2 == nil {
		level2 = lib.NewOrderedMap()
		tr.namedAccumulators.Put(groupingKey, level2)
		// E.g. if grouping by "color" and "shape", and the current record has
		// color=blue, shape=circle, then groupByFieldValues is the array
		// ["blue", "circle"].
		tr.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	}
	for _, valueFieldName := range tr.valueFieldNameList {
		valueFieldValue := inrec.Get(valueFieldName)
		if valueFieldValue == nil {
			continue
		}
		level3 := level2.(*lib.OrderedMap).Get(valueFieldName)
		if level3 == nil {
			level3 = lib.NewOrderedMap()
			level2.(*lib.OrderedMap).Put(valueFieldName, level3)
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
			namedAccumulator.(*utils.Stats1NamedAccumulator).Ingest(valueFieldValue)
		}
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

// ----------------------------------------------------------------
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

// ----------------------------------------------------------------
// TODO: comment
func (tr *TransformerStats1) emitIntoOutputRecord(
	inrec *types.Mlrmap,
	groupByFieldValues []*types.Mlrval,
	level2accumulators *lib.OrderedMap,
	outrec *types.Mlrmap,
) {
	for i, groupByFieldName := range tr.groupByFieldNameList {
		outrec.PutCopy(groupByFieldName, groupByFieldValues[i])
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
