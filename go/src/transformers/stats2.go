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
const verbNameStats2 = "stats2"

var Stats2Setup = transforming.TransformerSetup{
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
	argv0 := lib.MlrExeName()
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
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	argv0 := lib.MlrExeName()

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
			accumulatorNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-f" {
			valueFieldNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-g" {
			groupByFieldNameList = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

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

	transformer, _ := NewTransformerStats2(
		accumulatorNameList,
		valueFieldNameList,
		groupByFieldNameList,
		doVerbose,
		doIterativeStats,
		doHoldAndFit,
	)

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
	//   groupByFieldName . value1FieldName+value2FieldName . accumulatorName . accumulator object
	// This would be
	//   namedAccumulators map[string]map[string]map[string]Stats2NamedAccumulator
	// except we need maps that preserve insertion order.
	namedAccumulators *lib.OrderedMap

	groupingKeysToGroupByFieldValues map[string][]*types.Mlrval
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
					lib.MlrExeName(), name,
				),
			)
		}
	}

	this := &TransformerStats2{
		accumulatorNameList:              accumulatorNameList,
		valueFieldNameList:               valueFieldNameList,
		groupByFieldNameList:             groupByFieldNameList,
		doVerbose:                        doVerbose,
		doIterativeStats:                 doIterativeStats,
		doHoldAndFit:                     doHoldAndFit,
		accumulatorFactory:               utils.NewStats2AccumulatorFactory(),
		namedAccumulators:                lib.NewOrderedMap(),
		groupingKeysToGroupByFieldValues: make(map[string][]*types.Mlrval),
	}
	return this, nil
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
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
// }
// ================================================================
//
// In the iterative case, add to the current record its current group's stats fields.
// In the non-iterative case, produce output only at the end of the input stream.

// ----------------------------------------------------------------
func (this *TransformerStats2) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.handleInputRecord(inrecAndContext, outputChannel)
	} else {
		this.handleEndOfRecordStream(inrecAndContext, outputChannel)
	}
}

// ----------------------------------------------------------------
func (this *TransformerStats2) handleInputRecord(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record

// ----------------------------------------------------------------
//		mapper_stats2_ingest(pinrec, pctx, this)

//		if (this.doIterativeStats) {
//			// The input record is modified in this case, with new fields appended
//			return sllv_single(pinrec)
//		} else if (this.doHoldAndFit) {
//			// The input record is held by the ingestor
//			return nil
//		} else {
//			lrec_free(pinrec)
//			return nil
//		}
// ----------------------------------------------------------------

	// E.g. if grouping by "a" and "b", and the current record has a=circle, b=blue,
	// then groupingKey is the string "circle,blue".
	groupingKey, groupByFieldValues, ok := inrec.GetSelectedValuesAndJoined(
		this.groupByFieldNameList,
	)
	if !ok {
		return
	}

	level2 := this.namedAccumulators.Get(groupingKey)
	if level2 == nil {
		level2 = lib.NewOrderedMap()
		this.namedAccumulators.Put(groupingKey, level2)
		// E.g. if grouping by "color" and "shape", and the current record has
		// color=blue, shape=circle, then groupByFieldValues is the array
		// ["blue", "circle"].
		this.groupingKeysToGroupByFieldValues[groupingKey] = groupByFieldValues
	}
	for _, valueFieldName := range this.valueFieldNameList {
		valueFieldValue := inrec.Get(valueFieldName)
		if valueFieldValue == nil {
			continue
		}
		level3 := level2.(*lib.OrderedMap).Get(valueFieldName)
		if level3 == nil {
			level3 = lib.NewOrderedMap()
			level2.(*lib.OrderedMap).Put(valueFieldName, level3)
		}
		for _, accumulatorName := range this.accumulatorNameList {
			namedAccumulator := level3.(*lib.OrderedMap).Get(accumulatorName)
			if namedAccumulator == nil {
				namedAccumulator = this.accumulatorFactory.MakeNamedAccumulator(
					accumulatorName,
					groupingKey,
					valueFieldName, // TODO
					valueFieldName,
				)
				level3.(*lib.OrderedMap).Put(accumulatorName, namedAccumulator)
			}
			// TODO
			namedAccumulator.(*utils.Stats2NamedAccumulator).Ingest(valueFieldValue, valueFieldValue)
		}
	}

	if this.doIterativeStats {
		this.emitIntoOutputRecord(
			inrecAndContext.Record,
			groupByFieldValues,
			level2.(*lib.OrderedMap),
			inrec,
		)
		outputChannel <- inrecAndContext
	}
}

//// ----------------------------------------------------------------
//static void mapper_stats2_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats2_state_t* this) {
//	// ["s", "t"]
//	slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec, this.pgroup_by_field_names)
//	if (pgroup_by_field_values == nil) {
//		return
//	}
//
//	lhms2v_t* pgroup_to_acc_field = lhmslv_get(this.acc_groups, pgroup_by_field_values)
//	if (pgroup_to_acc_field == nil) {
//		pgroup_to_acc_field = lhms2v_alloc()
//		lhmslv_put(this.acc_groups, slls_copy(pgroup_by_field_values), pgroup_to_acc_field, FREE_ENTRY_KEY)
//	}
//
//	if (this.doHoldAndFit) { // Retain the input record in memory, for fitting and delivery at end of stream
//		sllv_t* group_to_records = lhmslv_get(this.record_groups, pgroup_by_field_values)
//		if (group_to_records == nil) {
//			group_to_records = sllv_alloc()
//			lhmslv_put(this.record_groups, slls_copy(pgroup_by_field_values), group_to_records, FREE_ENTRY_KEY)
//		}
//		sllv_append(group_to_records, pinrec)
//	}
//
//	// for [["x","y"]]
//	int n = this.pvalue_field_name_pairs.length
//	for (int i = 0; i < n; i += 2) {
//		char* value_field_name_1 = this.pvalue_field_name_pairs.strings[i]
//		char* value_field_name_2 = this.pvalue_field_name_pairs.strings[i+1]
//
//		lhmsv_t* pacc_fields_to_acc_state = lhms2v_get(pgroup_to_acc_field, value_field_name_1, value_field_name_2)
//		if (pacc_fields_to_acc_state == nil) {
//			pacc_fields_to_acc_state = lhmsv_alloc()
//			lhms2v_put(pgroup_to_acc_field, value_field_name_1, value_field_name_2, pacc_fields_to_acc_state, NO_FREE)
//		}
//
//		char* sval1 = lrec_get(pinrec, value_field_name_1)
//		char* sval2 = lrec_get(pinrec, value_field_name_2)
//		if (sval1 == nil) // Key not present
//			continue
//		if (*sval1 == 0) // Key present with null value
//			continue
//		if (sval2 == nil) // Key not present
//			continue
//		if (*sval2 == 0) // Key present with null value
//			continue
//
//		// for ["corr", "cov"]
//		sllse_t* pc = this.paccumulator_names.phead
//		for ( ; pc != nil; pc = pc.pnext) {
//			char* stats2_acc_name = pc.value
//			stats2_acc_t* pstats2_acc = lhmsv_get(pacc_fields_to_acc_state, stats2_acc_name)
//			if (pstats2_acc == nil) {
//				pstats2_acc = make_stats2(value_field_name_1, value_field_name_2, stats2_acc_name, this.do_verbose)
//				if (pstats2_acc == nil) {
//					fprintf(stderr, "mlr stats2: accumulator \"%s\" not found.\n",
//						stats2_acc_name)
//					exit(1)
//				}
//				lhmsv_put(pacc_fields_to_acc_state, stats2_acc_name, pstats2_acc, NO_FREE)
//			}
//			if (sval1 == nil || sval2 == nil)
//				continue
//
//			double dval1 = mlr_double_from_string_or_die(sval1)
//			double dval2 = mlr_double_from_string_or_die(sval2)
//			pstats2_acc.pingest_func(pstats2_acc.pvstate, dval1, dval2)
//		}
//		if (this.doIterativeStats) {
//			mapper_stats2_emit(this, pinrec, value_field_name_1, value_field_name_2,
//				pacc_fields_to_acc_state)
//		}
//	}
//}

// ----------------------------------------------------------------
func (this *TransformerStats2) handleEndOfRecordStream(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if this.doIterativeStats {
		outputChannel <- inrecAndContext // end-of-stream marker
		return
	}

	for pa := this.namedAccumulators.Head; pa != nil; pa = pa.Next {
		groupingKey := pa.Key
		level2 := pa.Value.(*lib.OrderedMap)
		groupByFieldValues := this.groupingKeysToGroupByFieldValues[groupingKey]

		newrec := types.NewMlrmapAsRecord()

		this.emitIntoOutputRecord(
			inrecAndContext.Record,
			groupByFieldValues,
			level2,
			newrec,
		)

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)
	}

	outputChannel <- inrecAndContext // end-of-stream marker
}

//	if (!this.doIterativeStats) {
//		if (!this.doHoldAndFit) {
//			return mapper_stats2_emit_all(this)
//		} else {
//			return mapper_stats2_fit_all(this)
//		}
//	} else {
//		return nil
//	}

// ----------------------------------------------------------------
// TODO: comment
func (this *TransformerStats2) emitIntoOutputRecord(
	inrec *types.Mlrmap,
	groupByFieldValues []*types.Mlrval,
	level2accumulators *lib.OrderedMap,
	outrec *types.Mlrmap,
) {
	for i, groupByFieldName := range this.groupByFieldNameList {
		outrec.PutCopy(groupByFieldName, groupByFieldValues[i])
	}

	for pb := level2accumulators.Head; pb != nil; pb = pb.Next {
		level3 := pb.Value.(*lib.OrderedMap)
		for pc := level3.Head; pc != nil; pc = pc.Next {
			namedAccumulator := pc.Value.(*utils.Stats2NamedAccumulator)
			key, value := namedAccumulator.Emit()
			outrec.PutCopy(key, &value)
		}
	}
}

//// ----------------------------------------------------------------
//static sllv_t* mapper_stats2_emit_all(mapper_stats2_state_t* this) {
//	sllv_t* poutrecs = sllv_alloc()
//
//	for (lhmslve_t* pa = this.acc_groups.phead; pa != nil; pa = pa.pnext) {
//		lrec_t* poutrec = lrec_unbacked_alloc()
//
//		// Add in a=s,b=t fields:
//		slls_t* pgroup_by_field_values = pa.key
//		sllse_t* pb = this.pgroup_by_field_names.phead
//		sllse_t* pc =         pgroup_by_field_values.phead
//		for ( ; pb != nil && pc != nil; pb = pb.pnext, pc = pc.pnext) {
//			lrec_put(poutrec, pb.value, pc.value, 0)
//		}
//
//		// Add in fields such as x_y_corr, etc.
//		lhms2v_t* pgroup_to_acc_field = pa.pvvalue
//
//		// For "x","y"
//		for (lhms2ve_t* pd = pgroup_to_acc_field.phead; pd != nil; pd = pd.pnext) {
//			char*    value_field_name_1 = pd.key1
//			char*    value_field_name_2 = pd.key2
//			lhmsv_t* pacc_fields_to_acc_state = pd.pvvalue
//
//			mapper_stats2_emit(this, poutrec, value_field_name_1, value_field_name_2,
//				pacc_fields_to_acc_state)
//
//			// For "corr", "linreg"
//			for (lhmsve_t* pe = pacc_fields_to_acc_state.phead; pe != nil; pe = pe.pnext) {
//				stats2_acc_t* pstats2_acc = pe.pvvalue
//				pstats2_acc.pemit_func(pstats2_acc.pvstate, value_field_name_1, value_field_name_2, poutrec)
//			}
//		}
//
//		sllv_append(poutrecs, poutrec)
//	}
//	sllv_append(poutrecs, nil)
//	return poutrecs
//}

//static void mapper_stats2_emit(mapper_stats2_state_t* this, lrec_t* poutrec,
//	char* value_field_name_1, char* value_field_name_2, lhmsv_t* pacc_fields_to_acc_state)
//{
//	// For "corr", "linreg"
//	for (lhmsve_t* pe = pacc_fields_to_acc_state.phead; pe != nil; pe = pe.pnext) {
//		stats2_acc_t* pstats2_acc = pe.pvvalue
//		pstats2_acc.pemit_func(pstats2_acc.pvstate, value_field_name_1, value_field_name_2, poutrec)
//	}
//}

//// ----------------------------------------------------------------
//static sllv_t* mapper_stats2_fit_all(mapper_stats2_state_t* this) {
//	sllv_t* poutrecs = sllv_alloc()
//
//	for (lhmslve_t* pa = this.acc_groups.phead; pa != nil; pa = pa.pnext) {
//		slls_t* pgroup_by_field_values = pa.key
//		sllv_t* precords = lhmslv_get(this.record_groups, pgroup_by_field_values)
//
//		while (precords.phead) {
//			lrec_t* prec = sllv_pop(precords)
//
//			lhms2v_t* pgroup_to_acc_field = pa.pvvalue
//
//			// For "x","y"
//			for (lhms2ve_t* pd = pgroup_to_acc_field.phead; pd != nil; pd = pd.pnext) {
//				char*    value_field_name_1 = pd.key1
//				char*    value_field_name_2 = pd.key2
//				lhmsv_t* pacc_fields_to_acc_state = pd.pvvalue
//
//				// For "linreg-ols", "logireg"
//				for (lhmsve_t* pe = pacc_fields_to_acc_state.phead; pe != nil; pe = pe.pnext) {
//					stats2_acc_t* pstats2_acc = pe.pvvalue
//					if (pstats2_acc.pfit_func != nil) {
//						char* sx = lrec_get(prec, value_field_name_1)
//						char* sy = lrec_get(prec, value_field_name_2)
//						if (sx != nil && sy != nil) {
//							double x = mlr_double_from_string_or_die(sx)
//							double y = mlr_double_from_string_or_die(sy)
//							pstats2_acc.pfit_func(pstats2_acc.pvstate, x, y, prec)
//						}
//					}
//				}
//			}
//
//			sllv_append(poutrecs, prec)
//		}
//	}
//	sllv_append(poutrecs, nil)
//	return poutrecs
//}
