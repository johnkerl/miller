package transformers

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/transformers/utils"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameStats2 = "stats2"

// For joining "x" and "y" into "x...y" for map keys. "," is another natural choice but would break
// if we were ever asked to process field names with commas in them.
const stats2KeySeparator = "\001"

var Stats2Setup = TransformerSetup{
	Verb:         verbNameStats2,
	UsageFunc:    transformerStats2Usage,
	ParseCLIFunc: transformerStats2ParseCLI,
	IgnoresInput: false,
}

func transformerStats2Usage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameStats2

	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Computes bivariate statistics for one or more given field-name pairs,\n")
	fmt.Fprintf(o, "accumulated across the input record stream.\n")
	fmt.Fprintf(o, "-a {linreg-ols,corr,...}  Names of accumulators: one or more of:\n")

	utils.ListStats2Accumulators(o)

	fmt.Fprintf(o, "-f {a,b,c,d}   Value-field name-pairs on which to compute statistics.\n")
	fmt.Fprintf(o, "               There must be an even number of names.\n")
	fmt.Fprintf(o, "-g {e,f,g}     Optional group-by-field names.\n")
	fmt.Fprintf(o, "-v             Print additional output for linreg-pca.\n")
	fmt.Fprintf(o, "-s             Print iterative stats. Useful in tail -f contexts (in which\n")
	fmt.Fprintf(o, "               case please avoid pprint-format output since end of input\n")
	fmt.Fprintf(o, "               stream will never be seen).\n")
	fmt.Fprintf(o, "--fit          Rather than printing regression parameters, applies them to\n")
	fmt.Fprintf(o, "               the input data to compute new fit fields. All input records are\n")
	fmt.Fprintf(o, "               held in memory until end of input stream. Has effect only for\n")
	fmt.Fprintf(o, "               linreg-ols, linreg-pca, and logireg.\n")
	fmt.Fprintf(o, "Only one of -s or --fit may be used.\n")
	fmt.Fprintf(o, "Example: %s %s -a linreg-pca -f x,y\n", argv0, verb)
	fmt.Fprintf(o, "Example: %s %s -a linreg-ols,r2 -f x,y -g size,shape\n", argv0, verb)
	fmt.Fprintf(o, "Example: %s %s -a corr -f x,y\n", argv0, verb)

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
func transformerStats2ParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	argv0 := "mlr"

	var accumulatorNameList []string = nil
	var valueFieldNameList []string = nil
	groupByFieldNameList := make([]string, 0)
	doVerbose := false
	doIterativeStats := false
	doHoldAndFit := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerStats2Usage(os.Stdout, true, 0)

		} else if opt == "-a" {
			accumulatorNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNameList = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-v" {
			doVerbose = true

		} else if opt == "-s" {
			doIterativeStats = true

		} else if opt == "--fit" {
			doHoldAndFit = true

		} else if opt == "-S" {
			// No-op pass-through for backward compatibility with Miller 5

		} else if opt == "-F" {
			// The -F flag isn't used for stats2: all arithmetic here is
			// floating-point. Yet it is supported for step and stats1 for all
			// applicable stats1/step accumulators, so we accept here as well
			// for all applicable stats2 accumulators (i.e. none of them).

		} else {
			transformerStats2Usage(os.Stderr, true, 1)
		}
	}

	if doIterativeStats && doHoldAndFit {
		transformerStats2Usage(os.Stderr, true, 1)
	}
	if accumulatorNameList == nil {
		fmt.Fprintf(os.Stderr, "%s %s: -a option is required.\n", argv0, verb)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", argv0, verb)
		os.Exit(1)
	}
	if valueFieldNameList == nil {
		fmt.Fprintf(os.Stderr, "%s %s: -f option is required.\n", argv0, verb)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", argv0, verb)
		os.Exit(1)
	}
	if len(valueFieldNameList)%2 != 0 {
		fmt.Fprintf(os.Stderr, "%s %s: argument to -f must have even number of fields.\n", argv0, verb)
		fmt.Fprintf(os.Stderr, "Please see %s %s --help for more information.\n", argv0, verb)
		os.Exit(1)
	}

	transformer, err := NewTransformerStats2(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,
		doVerbose,
		doIterativeStats,
		doHoldAndFit,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerStats2 struct {
	// Input:
	accumulatorNameList  []string
	valueFieldNameList   []string
	groupByFieldNameList []string

	doVerbose        bool
	doIterativeStats bool
	doHoldAndFit     bool

	// State:
	accumulatorFactory *utils.Stats2AccumulatorFactory

	// Accumulators are indexed by
	//   groupByFieldName . value1FieldName+sep+value2FieldName . accumulatorName . accumulator object
	// This would be
	//   namedAccumulators map[string]map[string]map[string]IStats2Accumulator
	// except we need maps that preserve insertion order.
	namedAccumulators *lib.OrderedMap

	groupingKeysToGroupByFieldValues *lib.OrderedMap

	// For hold-and-fit:
	// ordered map from grouping-key to list of RecordAndContext
	recordGroups *lib.OrderedMap
}

func NewTransformerStats2(
	accumulatorNameList []string,
	valueFieldNameList []string,
	groupByFieldNameList []string,
	doVerbose bool,
	doIterativeStats bool,
	doHoldAndFit bool,
) (*TransformerStats2, error) {
	for _, name := range accumulatorNameList {
		if !utils.ValidateStats2AccumulatorName(name) {
			return nil, errors.New(
				fmt.Sprintf(
					"%s stats2: accumulator \"%s\" not found.\n",
					"mlr", name,
				),
			)
		}
	}

	tr := &TransformerStats2{
		accumulatorNameList:              accumulatorNameList,
		valueFieldNameList:               valueFieldNameList,
		groupByFieldNameList:             groupByFieldNameList,
		doVerbose:                        doVerbose,
		doIterativeStats:                 doIterativeStats,
		doHoldAndFit:                     doHoldAndFit,
		accumulatorFactory:               utils.NewStats2AccumulatorFactory(),
		namedAccumulators:                lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: lib.NewOrderedMap(),
		recordGroups:                     lib.NewOrderedMap(),
	}
	return tr, nil
}

// ================================================================
// Given: accumulate corr,cov on values x,y group by a,b.
// Example input:       Example output:
//   a b x y            a b x_corr x_cov y_corr y_cov
//   s t 1 2            s t 2       6    2      8
//   u v 3 4            u v 1       3    1      4
//   s t 5 6            u w 1       7    1      9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   ["s","t"] : {                    <--- group-by field names
//     ["x","y"] : {                  <--- value field names
//       "corr" : stats2_corr object,
//       "cov"  : stats2_cov  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr object,
//       "cov"  : stats2_cov  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr object,
//       "cov"  : stats2_cov  object
//     }
//   },
// }
//
// In the iterative case, add to the current record its current group's stats fields.
// In the non-iterative case, produce output only at the end of the input stream.
// ================================================================

// ----------------------------------------------------------------

func (tr *TransformerStats2) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {

		tr.ingest(inrecAndContext)

		if tr.doIterativeStats {
			// The input record is modified in this case, with new fields appended
			outputChannel <- inrecAndContext
		}
		// if tr.doHoldAndFit, the input record is held by the ingestor

	} else { // end of record stream
		if !tr.doIterativeStats { // in the iterative case, already emitted per-record
			if tr.doHoldAndFit {
				tr.fit(outputChannel)
			} else {
				tr.emit(outputChannel, &inrecAndContext.Context)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerStats2) ingest(
	inrecAndContext *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

	// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
	// then groupingKey is the string "circle,blue".
	groupingKey, groupByFieldValues, ok := inrec.GetSelectedValuesAndJoined(tr.groupByFieldNameList)
	if !ok {
		return
	}

	tr.groupingKeysToGroupByFieldValues.Put(groupingKey, groupByFieldValues)

	groupToValueFields := tr.namedAccumulators.Get(groupingKey)
	if groupToValueFields == nil {
		groupToValueFields = lib.NewOrderedMap()
		tr.namedAccumulators.Put(groupingKey, groupToValueFields)
	}

	if tr.doHoldAndFit { // Retain the input record in memory, for fitting and delivery at end of stream
		groupToRecords := tr.recordGroups.Get(groupingKey)
		if groupToRecords == nil {
			groupToRecords = list.New()
			tr.recordGroups.Put(groupingKey, groupToRecords)
		}
		groupToRecords.(*list.List).PushBack(inrecAndContext)
	}

	// for [["x","y"]]
	n := len(tr.valueFieldNameList)
	for i := 0; i < n; i += 2 {
		valueFieldName1 := tr.valueFieldNameList[i]
		valueFieldName2 := tr.valueFieldNameList[i+1]

		key := valueFieldName1 + stats2KeySeparator + valueFieldName2

		valueFieldsToAccumulator := groupToValueFields.(*lib.OrderedMap).Get(key)
		if valueFieldsToAccumulator == nil {
			valueFieldsToAccumulator = lib.NewOrderedMap()
			groupToValueFields.(*lib.OrderedMap).Put(key, valueFieldsToAccumulator)
		}

		mval1 := inrec.Get(valueFieldName1)
		mval2 := inrec.Get(valueFieldName2)
		if mval1 == nil || mval2 == nil { // Key absent in current record
			continue
		}
		if mval1.IsVoid() || mval2.IsVoid() { // Key present in current record but with empty value
			continue
		}

		// for ["corr", "cov"]
		for _, accumulatorName := range tr.accumulatorNameList {
			accumulator := valueFieldsToAccumulator.(*lib.OrderedMap).Get(accumulatorName)
			if accumulator == nil {
				accumulator = tr.accumulatorFactory.Make(
					valueFieldName1,
					valueFieldName2,
					accumulatorName,
					tr.doVerbose,
				)
				if accumulator == nil {
					fmt.Fprintf(os.Stderr, "%s %s: accumulator \"%s\" not found.\n",
						"mlr", verbNameStats2, accumulatorName,
					)
					os.Exit(1)
				}
				valueFieldsToAccumulator.(*lib.OrderedMap).Put(accumulatorName, accumulator)
			}
			accumulator.(utils.IStats2Accumulator).Ingest(
				mval1.GetNumericToFloatValueOrDie(),
				mval2.GetNumericToFloatValueOrDie(),
			)
		}

		if tr.doIterativeStats {
			tr.populateRecord(
				inrecAndContext.Record,
				valueFieldName1,
				valueFieldName2,
				valueFieldsToAccumulator.(*lib.OrderedMap),
			)
		}
	}
}

// ----------------------------------------------------------------
func (tr *TransformerStats2) emit(
	outputChannel chan<- *types.RecordAndContext,
	context *types.Context,
) {
	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		outrec := types.NewMlrmapAsRecord()

		// Add in a=s,b=t fields:
		groupingKey := pa.Key
		groupByFieldValues := tr.groupingKeysToGroupByFieldValues.Get(groupingKey).([]*types.Mlrval)
		for i, groupByFieldName := range tr.groupByFieldNameList {
			outrec.PutReference(groupByFieldName, groupByFieldValues[i].Copy())
		}

		// Add in fields such as x_y_corr, etc.
		groupToValueFields := tr.namedAccumulators.Get(groupingKey).(*lib.OrderedMap)

		// For "x","y"
		for pc := groupToValueFields.Head; pc != nil; pc = pc.Next {

			pairs := strings.Split(pc.Key, stats2KeySeparator)
			valueFieldName1 := pairs[0]
			valueFieldName2 := pairs[1]
			valueFieldsToAccumulator := pc.Value.(*lib.OrderedMap)

			tr.populateRecord(outrec, valueFieldName1, valueFieldName2, valueFieldsToAccumulator)

			// For "corr", "linreg"
			for pd := valueFieldsToAccumulator.Head; pd != nil; pd = pd.Next {
				accumulator := pd.Value.(utils.IStats2Accumulator)
				accumulator.Populate(valueFieldName1, valueFieldName2, outrec)
			}
		}

		outputChannel <- types.NewRecordAndContext(outrec, context)
	}
}

func (tr *TransformerStats2) populateRecord(
	outrec *types.Mlrmap,
	valueFieldName1 string,
	valueFieldName2 string,
	valueFieldsToAccumulator *lib.OrderedMap,
) {
	// For "corr", "linreg"
	for pe := valueFieldsToAccumulator.Head; pe != nil; pe = pe.Next {
		accumulator := pe.Value.(utils.IStats2Accumulator)
		accumulator.Populate(valueFieldName1, valueFieldName2, outrec)
	}
}

func (tr *TransformerStats2) fit(
	outputChannel chan<- *types.RecordAndContext,
) {
	for pa := tr.namedAccumulators.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		groupToValueFields := pa.Value.(*lib.OrderedMap)
		recordsAndContexts := tr.recordGroups.Get(groupingKey).(*list.List)

		for recordsAndContexts.Front() != nil {
			recordAndContext := recordsAndContexts.Remove(recordsAndContexts.Front()).(*types.RecordAndContext)
			record := recordAndContext.Record

			// For "x","y"
			for pb := groupToValueFields.Head; pb != nil; pb = pb.Next {
				pairs := strings.Split(pb.Key, stats2KeySeparator)
				valueFieldName1 := pairs[0]
				valueFieldName2 := pairs[1]
				valueFieldsToAccumulator := pb.Value.(*lib.OrderedMap)

				// For "linreg-ols", "logireg"
				for pc := valueFieldsToAccumulator.Head; pc != nil; pc = pc.Next {
					accumulator := pc.Value.(utils.IStats2Accumulator)

					// Note R2, cov, corr, etc have no non-trivial fit-function
					mval1 := record.Get(valueFieldName1)
					mval2 := record.Get(valueFieldName2)
					if mval1 != nil && mval2 != nil {
						accumulator.Fit(
							mval1.GetNumericToFloatValueOrDie(),
							mval2.GetNumericToFloatValueOrDie(),
							record,
						)
					}
				}
			}

			outputChannel <- recordAndContext
		}
	}
}
