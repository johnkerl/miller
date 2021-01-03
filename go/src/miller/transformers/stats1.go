package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var Stats1Setup = transforming.TransformerSetup{
	Verb:         "stats1",
	ParseCLIFunc: transformerStats1ParseCLI,
	IgnoresInput: false,
}

func transformerStats1ParseCLI(
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

	accumulatorNameList := make([]string, 0)
	valueFieldNameList := make([]string, 0)
	groupByFieldNameList := make([]string, 0)

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process
		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerStats1Usage(os.Stdout, args[0], verb, nil)
			return nil // help intentionally requested

		} else if args[argi] == "-a" {
			checkArgCountStats1(verb, args, argi, argc, 2)
			accumulatorNameList = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else if args[argi] == "-f" {
			checkArgCountStats1(verb, args, argi, argc, 2)
			valueFieldNameList = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else if args[argi] == "-g" {
			checkArgCountStats1(verb, args, argi, argc, 2)
			groupByFieldNameList = lib.SplitString(args[argi+1], ",")
			argi += 2

		} else {
			transformerStats1Usage(os.Stderr, args[0], verb, nil)
			os.Exit(1)
		}
	}

	// TODO: libify for use across verbs.
	if len(accumulatorNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -a option is required.\n", args[0], verb)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", args[0], verb)
		os.Exit(1)
	}
	if len(valueFieldNameList) == 0 {
		fmt.Fprintf(os.Stderr, "%s %s: -f option is required.\n", args[0], verb)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", args[0], verb)
		os.Exit(1)
	}

	transformer, _ := NewTransformerStats1(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,
	)

	*pargi = argi
	return transformer
}

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func checkArgCountStats1(verb string, args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s %s: option \"%s\" missing argument(s).\n", args[0], verb, args[argi])
		os.Exit(1)
	}
}

func transformerStats1Usage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Computes univariate statistics for one or more given fields, accumulated across
the input record stream.
Options:
-a {sum,count,...} Names of accumulators: one or more of:
  median   This is the same as p50
  p10 p25.2 p50 p98 p100 etc.
  TODO: flags for interpolated percentiles
`)
	listStats1Accumulators(o)
	fmt.Fprint(o, `
-f {a,b,c}   Value-field names on which to compute statistics
-g {d,e,f}   Optional group-by-field names

[TODO: more]
`)

	fmt.Fprintln(o,
		"Example: mlr stats1 -a min,p10,p50,p90,max -f value -g size,shape\n", argv0, verb)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size\n", argv0, verb)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode -f size -g shape\n", argv0, verb)
	fmt.Fprintln(o,
		"Example: mlr stats1 -a count,mode --fr '^[a-h].*$' -gr '^k.*$'\n", argv0, verb)
	fmt.Fprintln(o,
		`        This computes count and mode statistics on all field names beginning
         with a through h, grouped by all field names starting with k.
`)
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
}

// ----------------------------------------------------------------
type TransformerStats1 struct {
	// Input:
	accumulatorNameList  []string
	valueFieldNameList   []string
	groupByFieldNameList []string

	// State:
	accumulatorFactory *Stats1AccumulatorFactory

	// Accumulators are indexed by
	//   groupByFieldName -> valueFieldName -> accumulatorName -> accumulator object
	// This would be
	//   accumulators map[string]map[string]map[string]IStats1Accumulator
	// except we need maps that preserve insertion order.
	accumulators *lib.OrderedMap

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
) (*TransformerStats1, error) {
	this := &TransformerStats1{
		accumulatorNameList:              accumulatorNameList,
		valueFieldNameList:               valueFieldNameList,
		groupByFieldNameList:             groupByFieldNameList,
		accumulatorFactory:               NewStats1AccumulatorFactory(),
		accumulators:                     lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: make(map[string][]*types.Mlrval),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerStats1) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.Record != nil {
		this.handleInputRecord(inrecAndContext, outputChannel)
	} else {
		this.handleEndOfRecordStream(inrecAndContext, outputChannel)
	}
}

// ----------------------------------------------------------------
func (this *TransformerStats1) handleInputRecord(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
	// then groupingKey is the string "circle,blue".
	groupingKey, ok := inrec.GetSelectedValuesJoined(
		this.groupByFieldNameList,
	)
	if !ok {
		return
	}

	level2 := this.accumulators.Get(groupingKey)
	if level2 == nil {
		level2 = lib.NewOrderedMap()
		this.accumulators.Put(groupingKey, level2)
		// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
		// then groupByFieldValues is the array ["circle", "blue"].
		groupByFieldValues, _ := inrec.GetSelectedValues(this.groupByFieldNameList)
		this.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	}
	for _, valueFieldName := range this.valueFieldNameList {
		valueFieldValue := inrec.Get(&valueFieldName)
		if valueFieldValue == nil {
			continue
		}
		level3 := level2.(*lib.OrderedMap).Get(valueFieldName)
		if level3 == nil {
			level3 = lib.NewOrderedMap()
			level2.(*lib.OrderedMap).Put(valueFieldName, level3)
		}
		for _, accumulatorName := range this.accumulatorNameList {
			accumulator := level3.(*lib.OrderedMap).Get(accumulatorName)
			if accumulator == nil {
				// TODO: validate names at constructor time
				accumulator = this.accumulatorFactory.Make(accumulatorName, valueFieldName)
				if accumulator == nil {
					fmt.Fprintf(
						os.Stderr,
						"%s stats1: accumulator \"%s\" not found.\n",
						os.Args[0], accumulatorName,
					)
					os.Exit(1)
				}
				level3.(*lib.OrderedMap).Put(accumulatorName, accumulator)
			}
			accumulator.(IStats1Accumulator).Ingest(valueFieldValue)
		}
	}
}

// ----------------------------------------------------------------
func (this *TransformerStats1) handleEndOfRecordStream(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	for pa := this.accumulators.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		groupByFieldValues := this.groupingKeysToGroupByFieldValues[groupingKey]

		newrec := types.NewMlrmapAsRecord()
		for i, groupByFieldName := range this.groupByFieldNameList {
			newrec.PutCopy(&groupByFieldName, groupByFieldValues[i])
		}

		level2 := pa.Value.(*lib.OrderedMap)
		for pb := level2.Head; pb != nil; pb = pb.Next {
			valueFieldName := pb.Key
			level3 := pb.Value.(*lib.OrderedMap)
			for pc := level3.Head; pc != nil; pc = pc.Next {
				accumulatorName := pc.Key
				accumulator := pc.Value.(IStats1Accumulator)
				output := accumulator.Emit()
				key := valueFieldName + "_" + accumulatorName
				newrec.PutCopy(&key, &output)
			}
		}
		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	}

	outputChannel <- inrecAndContext // end-of-stream marker
}
