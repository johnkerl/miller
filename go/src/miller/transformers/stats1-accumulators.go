// ================================================================
// For stats1 as well as merge-fields
// ================================================================

package transformers

import (
	"fmt"
	"os"

	"miller/types"
)

// ----------------------------------------------------------------
type IStats1Accumulator interface {
	Ingest(value *types.Mlrval)
	Emit() types.Mlrval
}

// ----------------------------------------------------------------
type newStats1AccumulatorFunc func() IStats1Accumulator

type stats1AccumulatorInfo struct {
	name        string
	description string
	constructor newStats1AccumulatorFunc
}

var stats1AccumulatorInfos []stats1AccumulatorInfo = []stats1AccumulatorInfo{

	{
		"count",
		"Count instances of fields",
		NewStats1CountAccumulator,
	},

	{
		"mode",
		"Find most-frequently-occurring values for fields; first-found wins tie",
		NewStats1ModeAccumulator,
	},

	{
		"antimode",
		"Find least-frequently-occurring values for fields; first-found wins tie",
		NewStats1AntimodeAccumulator,
	},

	{
		"sum",
		"Compute sums of specified fields",
		NewStats1SumAccumulator,
	},
	{
		"mean",
		"Compute averages (sample means) of specified fields",
		NewStats1MeanAccumulator,
	},

	//   stddev    Compute sample standard deviation of specified fields
	//   var       Compute sample variance of specified fields
	//   meaneb    Estimate error bars for averages (assuming no sample autocorrelation)
	//   skewness  Compute sample skewness of specified fields
	//   kurtosis  Compute sample kurtosis of specified fields

	{
		"min",
		"Compute minimum values of specified fields",
		NewStats1MinAccumulator,
	},
	{
		"max",
		"Compute maximum values of specified fields",
		NewStats1MaxAccumulator,
	},
}

// ----------------------------------------------------------------
func listStats1Accumulators(o *os.File) {
	for _, info := range stats1AccumulatorInfos {
		fmt.Fprintf(o, "  %-8s %s\n", info.name, info.description)
	}
}

// ----------------------------------------------------------------
func makeStats1Accumulator(name string) IStats1Accumulator {
	for _, info := range stats1AccumulatorInfos {
		if info.name == name {
			return info.constructor()
		}
	}
	return nil
}

// ----------------------------------------------------------------
type Stats1CountAccumulator struct {
	count types.Mlrval
}

func NewStats1CountAccumulator() IStats1Accumulator {
	return &Stats1CountAccumulator{
		count: types.MlrvalFromInt64(0),
	}
}
func (this *Stats1CountAccumulator) Ingest(value *types.Mlrval) {
	this.count.Increment()
}
func (this *Stats1CountAccumulator) Emit() types.Mlrval {
	return *this.count.Copy()
}

// ----------------------------------------------------------------
type Stats1ModeAccumulator struct {
	countsByValue map[string]int64
}

func NewStats1ModeAccumulator() IStats1Accumulator {
	return &Stats1ModeAccumulator{
		countsByValue: make(map[string]int64),
	}
}
func (this *Stats1ModeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	this.countsByValue[key] = this.countsByValue[key] + 1
}
func (this *Stats1ModeAccumulator) Emit() types.Mlrval {
	if len(this.countsByValue) == 0 {
		return types.MlrvalFromError()
	}
	maxValue := ""
	var maxCount = int64(0)
	for value, count := range this.countsByValue {
		if maxValue == "" || count > maxCount {
			maxValue = value
			maxCount = count
		}
	}
	return types.MlrvalFromString(maxValue)
}

// ----------------------------------------------------------------
type Stats1AntimodeAccumulator struct {
	countsByValue map[string]int64
}

func NewStats1AntimodeAccumulator() IStats1Accumulator {
	return &Stats1AntimodeAccumulator{
		countsByValue: make(map[string]int64),
	}
}
func (this *Stats1AntimodeAccumulator) Ingest(value *types.Mlrval) {
	key := value.String() // 1, 1.0, and 1.000 are distinct
	this.countsByValue[key] = this.countsByValue[key] + 1
}
func (this *Stats1AntimodeAccumulator) Emit() types.Mlrval {
	if len(this.countsByValue) == 0 {
		return types.MlrvalFromError()
	}
	maxValue := ""
	var maxCount = int64(0)
	for value, count := range this.countsByValue {
		if maxValue == "" || count < maxCount {
			maxValue = value
			maxCount = count
		}
	}
	return types.MlrvalFromString(maxValue)
}

// ----------------------------------------------------------------
type Stats1SumAccumulator struct {
	sum types.Mlrval
}

func NewStats1SumAccumulator() IStats1Accumulator {
	return &Stats1SumAccumulator{
		sum: types.MlrvalFromInt64(0),
	}
}
func (this *Stats1SumAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
}
func (this *Stats1SumAccumulator) Emit() types.Mlrval {
	return *this.sum.Copy()
}

// ----------------------------------------------------------------
type Stats1MeanAccumulator struct {
	sum   types.Mlrval
	count types.Mlrval
}

func NewStats1MeanAccumulator() IStats1Accumulator {
	return &Stats1MeanAccumulator{
		sum:   types.MlrvalFromInt64(0),
		count: types.MlrvalFromInt64(0),
	}
}
func (this *Stats1MeanAccumulator) Ingest(value *types.Mlrval) {
	this.sum = types.MlrvalBinaryPlus(&this.sum, value)
	this.count.Increment()
}
func (this *Stats1MeanAccumulator) Emit() types.Mlrval {
	if this.count.IsIntZero() {
		return types.MlrvalFromVoid()
	} else {
		return types.MlrvalDivide(&this.sum, &this.count)
	}
}

// ----------------------------------------------------------------
type Stats1MinAccumulator struct {
	min types.Mlrval
}

func NewStats1MinAccumulator() IStats1Accumulator {
	return &Stats1MinAccumulator{
		min: types.MlrvalFromAbsent(),
	}
}
func (this *Stats1MinAccumulator) Ingest(value *types.Mlrval) {
	this.min = types.MlrvalBinaryMin(&this.min, value)
}
func (this *Stats1MinAccumulator) Emit() types.Mlrval {
	return *this.min.Copy()
}

// ----------------------------------------------------------------
type Stats1MaxAccumulator struct {
	max types.Mlrval
}

func NewStats1MaxAccumulator() IStats1Accumulator {
	return &Stats1MaxAccumulator{
		max: types.MlrvalFromAbsent(),
	}
}
func (this *Stats1MaxAccumulator) Ingest(value *types.Mlrval) {
	this.max = types.MlrvalBinaryMax(&this.max, value)
}
func (this *Stats1MaxAccumulator) Emit() types.Mlrval {
	return *this.max.Copy()
}

//// ----------------------------------------------------------------
//typedef struct _stats1_stddev_var_meaneb_state_t {
//	unsigned long long count;
//	double sumx;
//	double sumx2;
//	cumulant2o_t  do_which;
//	char* output_field_name;
//} stats1_stddev_var_meaneb_state_t;
//static void stats1_stddev_var_meaneb_dingest(void* pvstate, double val) {
//	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
//	pstate->count++;
//	pstate->sumx  += val;
//	pstate->sumx2 += val*val;
//}
//
//static void stats1_stddev_var_meaneb_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
//	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
//	if (pstate->count < 2LL) {
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
//		else
//			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
//	} else {
//		double output = mlr_get_var(pstate->count, pstate->sumx, pstate->sumx2);
//		if (pstate->do_which == DO_STDDEV)
//			output = sqrt(output);
//		else if (pstate->do_which == DO_MEANEB)
//			output = sqrt(output / pstate->count);
//		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//		else
//			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
//	}
//}
//
//stats1_acc_t* stats1_stddev_var_meaneb_alloc(char* value_field_name, char* stats1_acc_name, cumulant2o_t do_which) {
//	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
//	stats1_stddev_var_meaneb_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_stddev_var_meaneb_state_t));
//	pstate->count              = 0LL;
//	pstate->sumx               = 0.0;
//	pstate->sumx2              = 0.0;
//	pstate->do_which           = do_which;
//	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
//
//	pstats1_acc->pvstate       = (void*)pstate;
//	pstats1_acc->pdingest_func = stats1_stddev_var_meaneb_dingest;
//	pstats1_acc->pningest_func = NULL;
//	pstats1_acc->psingest_func = NULL;
//	pstats1_acc->pemit_func    = stats1_stddev_var_meaneb_emit;
//	return pstats1_acc;
//}
//stats1_acc_t* stats1_stddev_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_STDDEV);
//}
//stats1_acc_t* stats1_var_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_VAR);
//}
//stats1_acc_t* stats1_meaneb_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_MEANEB);
//}

//// ----------------------------------------------------------------
//typedef struct _stats1_skewness_state_t {
//	unsigned long long count;
//	double sumx;
//	double sumx2;
//	double sumx3;
//	char* output_field_name;
//} stats1_skewness_state_t;
//static void stats1_skewness_dingest(void* pvstate, double val) {
//	stats1_skewness_state_t* pstate = pvstate;
//	pstate->count++;
//	pstate->sumx  += val;
//	pstate->sumx2 += val*val;
//	pstate->sumx3 += val*val*val;
//}
//
//static void stats1_skewness_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
//	stats1_skewness_state_t* pstate = pvstate;
//	if (pstate->count < 2LL) {
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
//		else
//			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
//	} else {
//		double output = mlr_get_skewness(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3);
//		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//		else
//			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
//	}
//}
//
//stats1_acc_t* stats1_skewness_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
//	stats1_skewness_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_skewness_state_t));
//	pstate->count              = 0LL;
//	pstate->sumx               = 0.0;
//	pstate->sumx2              = 0.0;
//	pstate->sumx3              = 0.0;
//	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
//
//	pstats1_acc->pvstate       = (void*)pstate;
//	pstats1_acc->pdingest_func = stats1_skewness_dingest;
//	pstats1_acc->pningest_func = NULL;
//	pstats1_acc->psingest_func = NULL;
//	pstats1_acc->pemit_func    = stats1_skewness_emit;
//	return pstats1_acc;
//}

//// ----------------------------------------------------------------
//typedef struct _stats1_kurtosis_state_t {
//	unsigned long long count;
//	double sumx;
//	double sumx2;
//	double sumx3;
//	double sumx4;
//	char* output_field_name;
//} stats1_kurtosis_state_t;
//static void stats1_kurtosis_dingest(void* pvstate, double val) {
//	stats1_kurtosis_state_t* pstate = pvstate;
//	pstate->count++;
//	pstate->sumx  += val;
//	pstate->sumx2 += val*val;
//	pstate->sumx3 += val*val*val;
//	pstate->sumx4 += val*val*val*val;
//}
//
//static void stats1_kurtosis_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
//	stats1_kurtosis_state_t* pstate = pvstate;
//	if (pstate->count < 2LL) {
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
//		else
//			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
//	} else {
//		double output = mlr_get_kurtosis(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3, pstate->sumx4);
//		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//		else
//			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
//	}
//}
//stats1_acc_t* stats1_kurtosis_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
//	stats1_kurtosis_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_kurtosis_state_t));
//	pstate->count              = 0LL;
//	pstate->sumx               = 0.0;
//	pstate->sumx2              = 0.0;
//	pstate->sumx3              = 0.0;
//	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
//
//	pstats1_acc->pvstate       = (void*)pstate;
//	pstats1_acc->pdingest_func = stats1_kurtosis_dingest;
//	pstats1_acc->pningest_func = NULL;
//	pstats1_acc->psingest_func = NULL;
//	pstats1_acc->pemit_func    = stats1_kurtosis_emit;
//	return pstats1_acc;
//}

//// ----------------------------------------------------------------
//typedef struct _stats1_min_state_t {
//	mv_t min;
//	char* output_field_name;
//} stats1_min_state_t;
//static void stats1_min_singest(void* pvstate, char* sval) {
//	stats1_min_state_t* pstate = pvstate;
//	mv_t val = mv_copy_type_infer_string_or_float_or_int(sval);
//	pstate->min = x_xx_min_func(&pstate->min, &val);
//}
//static void stats1_min_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
//	stats1_min_state_t* pstate = pvstate;
//	if (mv_is_null(&pstate->min)) {
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
//		else
//			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
//	} else {
//		if (copy_data)
//			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mv_alloc_format_val(&pstate->min),
//				FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//		else
//			lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->min),
//				FREE_ENTRY_VALUE);
//	}
//}
//stats1_acc_t* stats1_min_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	stats1_acc_t* pstats1_acc  = mlr_malloc_or_die(sizeof(stats1_acc_t));
//	stats1_min_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_min_state_t));
//	pstate->min                = mv_absent();
//	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
//
//	pstats1_acc->pvstate       = (void*)pstate;
//	pstats1_acc->pdingest_func = NULL;
//	pstats1_acc->pningest_func = NULL;
//	pstats1_acc->psingest_func = stats1_min_singest;
//	pstats1_acc->pemit_func    = stats1_min_emit;
//	return pstats1_acc;
//}

//// ----------------------------------------------------------------
//typedef struct _stats1_percentile_state_t {
//	percentile_keeper_t* ppercentile_keeper;
//	lhmss_t* poutput_field_names;
//	int reference_count;
//	percentile_keeper_emitter_t* ppercentile_keeper_emitter;
//} stats1_percentile_state_t;
//static void stats1_percentile_singest(void* pvstate, char* sval) {
//	stats1_percentile_state_t* pstate = pvstate;
//	mv_t val = mv_copy_type_infer_string_or_float_or_int(sval);
//	percentile_keeper_ingest(pstate->ppercentile_keeper, val);
//}
//
//static void stats1_percentile_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
//	stats1_percentile_state_t* pstate = pvstate;
//	double p;
//
//	if (stats1_acc_name[0] == 'm') { // Pre-validated to be either p{number} or median.
//		p = 50.0;
//	} else {
//		// TODO: do the sscanf once at alloc time and store the double in the state struct for a minor perf gain.
//		(void)sscanf(stats1_acc_name, "p%lf", &p); // Assuming this was range-checked earlier on to be in [0,100].
//	}
//	mv_t v = pstate->ppercentile_keeper_emitter(pstate->ppercentile_keeper, p);
//	char* s = mv_alloc_format_val(&v);
//	// For this type, one accumulator tracks many stats1_names, but a single value_field_name.
//	char* output_field_name = lhmss_get(pstate->poutput_field_names, stats1_acc_name);
//	if (output_field_name == NULL) {
//		output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
//		lhmss_put(pstate->poutput_field_names, mlr_strdup_or_die(stats1_acc_name),
//			output_field_name, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//	}
//	lrec_put(poutrec, mlr_strdup_or_die(output_field_name), s, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
//}
//
//stats1_acc_t* stats1_percentile_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
//	int do_interpolated_percentiles)
//{
//	stats1_acc_t* pstats1_acc   = mlr_malloc_or_die(sizeof(stats1_acc_t));
//	stats1_percentile_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_percentile_state_t));
//	pstate->ppercentile_keeper  = percentile_keeper_alloc();
//	pstate->poutput_field_names = lhmss_alloc();
//	pstate->reference_count     = 1;
//	pstate->ppercentile_keeper_emitter = (do_interpolated_percentiles)
//		? percentile_keeper_emit_linearly_interpolated
//		: percentile_keeper_emit_non_interpolated;
//
//	pstats1_acc->pvstate        = (void*)pstate;
//	pstats1_acc->pdingest_func  = NULL;
//	pstats1_acc->pningest_func  = NULL;
//	pstats1_acc->psingest_func  = stats1_percentile_singest;
//	pstats1_acc->pemit_func     = stats1_percentile_emit;
//	return pstats1_acc;
//}
//void stats1_percentile_reuse(stats1_acc_t* pstats1_acc) {
//	stats1_percentile_state_t* pstate = pstats1_acc->pvstate;
//	pstate->reference_count++;
//}
