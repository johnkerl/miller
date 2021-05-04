// ================================================================
// For stats2
// ================================================================

package utils

import (
	"fmt"
	"os"

	"miller/src/lib"
	"miller/src/types"
)

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
// ================================================================

// ----------------------------------------------------------------
type IStats2Accumulator interface {
	Ingest(value1, value2 *types.Mlrval)
	Emit() types.Mlrval
}

// ----------------------------------------------------------------
type newStats2AccumulatorFunc func() IStats2Accumulator

type stats2AccumulatorInfo struct {
	name        string
	description string
	constructor newStats2AccumulatorFunc
}

var stats2AccumulatorInfos []stats2AccumulatorInfo = []stats2AccumulatorInfo{
	{
		"linreg-pca",
		"Linear regression using principal component analysis",
		NewStats2LinRegPCAAccumulator,
	},
	{
		"linreg-ols",
		"Linear regression using ordinary least squares",
		NewStats2LinRegOLSAccumulator,
	},
	{
		"r2",
		"Quality metric for linreg-ols (linreg-pca emits its own)",
		NewStats2R2Accumulator,
	},
	{
		"logireg",
		"Logistic regression",
		NewStats2LogiRegAccumulator,
	},
	{
		"corr",
		"Sample correlation",
		NewStats2CorrAccumulator,
	},
	{
		"cov",
		"Sample covariance",
		NewStats2CovAccumulator,
	},
	{
		"covx",
		"Sample-covariance matrix",
		NewStats2CovXAccumulator,
	},
}

// ================================================================
type Stats2NamedAccumulator struct {
	value1FieldName string
	value2FieldName string
	accumulatorName string
	accumulator     IStats2Accumulator
	outputFieldName string
}

func NewStats2NamedAccumulator(
	value1FieldName string,
	value2FieldName string,
	accumulatorName string,
	accumulator IStats2Accumulator,
) *Stats2NamedAccumulator {
	return &Stats2NamedAccumulator{
		value1FieldName: value1FieldName,
		value2FieldName: value2FieldName,
		accumulatorName: accumulatorName,
		accumulator:     accumulator,
		outputFieldName: value1FieldName + "_" + value2FieldName + "_" + accumulatorName,
	}
}

func (this *Stats2NamedAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.accumulator.Ingest(value1, value2)
}

func (this *Stats2NamedAccumulator) Emit() (key string, value types.Mlrval) {
	return this.outputFieldName, this.accumulator.Emit()
}

// ----------------------------------------------------------------
type Stats2AccumulatorFactory struct {
}

func NewStats2AccumulatorFactory() *Stats2AccumulatorFactory {
	return &Stats2AccumulatorFactory{}
}

// ----------------------------------------------------------------
func ListStats2Accumulators(o *os.File) {
	for _, info := range stats2AccumulatorInfos {
		fmt.Fprintf(o, "  %-8s %s\n", info.name, info.description)
	}
}

func ValidateStats2AccumulatorName(
	accumulatorName string,
) bool {
	for _, info := range stats2AccumulatorInfos {
		if info.name == accumulatorName {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
func (this *Stats2AccumulatorFactory) MakeNamedAccumulator(
	accumulatorName string,
	groupingKey string,
	value1FieldName string,
	value2FieldName string,
) *Stats2NamedAccumulator {

	accumulator := this.MakeAccumulator(
		accumulatorName,
		groupingKey,
		value1FieldName,
		value2FieldName,
	)
	// We don't return errors.New here. The nominal case is that the stats2
	// verb has already pre-validated accumulator names, and this is just a
	// fallback. The accumulators are instantiated for every unique combination
	// of group-by field values in the record stream, only as those values are
	// encountered: for example, with 'mlr stats2 -a count,sum -f x,y -g
	// color,shape', we make a new accumulator the first time we find a record
	// with 'color=blue,shape=square' and another the first time we find a
	// record with 'color=red,shape=circle', and so on. The right thing is to
	// pre-validate names once when the stats2 transformer is being
	// instantiated.
	lib.InternalCodingErrorIf(accumulator == nil)

	return NewStats2NamedAccumulator(
		value1FieldName,
		value2FieldName,
		accumulatorName,
		accumulator,
	)
}

func (this *Stats2AccumulatorFactory) MakeAccumulator(
	accumulatorName string,
	groupingKey string,
	value1FieldName string,
	value2FieldName string,
) IStats2Accumulator {
	for _, info := range stats2AccumulatorInfos {
		if info.name == accumulatorName {
			return info.constructor()
		}
	}
	return nil
}

// ================================================================
type Stats2LinRegPCAAccumulator struct {
	count int
}

func NewStats2LinRegPCAAccumulator() IStats2Accumulator {
	return &Stats2LinRegPCAAccumulator{
		count: 0,
	}
}
func (this *Stats2LinRegPCAAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2LinRegPCAAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

// ================================================================
type Stats2LinRegOLSAccumulator struct {
	count int
	sumx  *types.Mlrval
	sumx2 *types.Mlrval
	sumy  *types.Mlrval
	sumxy *types.Mlrval

// TODO
//	char*  m_output_field_name
//	char*  b_output_field_name
//	char*  n_output_field_name
//
//	char*  fit_output_field_name
//	int    fit_ready
//	double m
//	double b

}

func NewStats2LinRegOLSAccumulator() IStats2Accumulator {
	return &Stats2LinRegOLSAccumulator{
		count: 0,
		sumx:   types.MlrvalPointerFromInt(0),
		sumx2:  types.MlrvalPointerFromInt(0),
		sumy:   types.MlrvalPointerFromInt(0),
		sumxy:  types.MlrvalPointerFromInt(0),

//static stats2_acc_t* stats2_linreg_ols_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	this.m_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_m")
//	this.b_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_b")
//	this.n_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_n")
//	this.fit_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_fit")
//	this.fit_ready = FALSE
//	this.m         = -999.0
//	this.b         = -999.0
//
//	pstats2_acc.pingest_func = stats2_linreg_ols_ingest
//	pstats2_acc.pemit_func   = stats2_linreg_ols_emit
//	pstats2_acc.pfit_func    = stats2_linreg_ols_fit

	}
}

func (this *Stats2LinRegOLSAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
	x2 := types.MlrvalTimes(value1, value1)
	xy := types.MlrvalTimes(value1, value2)
	this.sumx = types.MlrvalBinaryPlus(this.sumx, value1)
	this.sumy = types.MlrvalBinaryPlus(this.sumx, value2)
	this.sumx2 = types.MlrvalBinaryPlus(this.sumx2, x2)
	this.sumxy = types.MlrvalBinaryPlus(this.sumx2, xy)
}

func (this *Stats2LinRegOLSAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

//static void stats2_linreg_ols_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
//	stats2_linreg_ols_state_t* this = pvstate
//
//	if (this.count < 2) {
//		lrec_put(poutrec, this.m_output_field_name, "", NO_FREE)
//		lrec_put(poutrec, this.b_output_field_name, "", NO_FREE)
//	} else {
//		double m, b
//		mlr_get_linear_regression_ols(this.count, this.sumx, this.sumx2, this.sumxy, this.sumy, &m, &b)
//		char* mval = mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt)
//		char* bval = mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt)
//
//		lrec_put(poutrec, this.m_output_field_name, mval, FREE_ENTRY_VALUE)
//		lrec_put(poutrec, this.b_output_field_name, bval, FREE_ENTRY_VALUE)
//	}
//
//	char* nval = mlr_alloc_string_from_ll(this.count)
//	lrec_put(poutrec, this.n_output_field_name, nval, FREE_ENTRY_VALUE)
//}

//static void stats2_linreg_ols_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
//	stats2_linreg_ols_state_t* this = pvstate
//
//	if (!this.fit_ready) {
//		mlr_get_linear_regression_ols(this.count, this.sumx, this.sumx2, this.sumxy, this.sumy,
//			&this.m, &this.b)
//		this.fit_ready = TRUE
//	}
//
//	if (this.count < 2) {
//		lrec_put(poutrec, this.fit_output_field_name, "", NO_FREE)
//	} else {
//		double yfit = this.m * x + this.b
//		char* sfit = mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt)
//		lrec_put(poutrec, this.fit_output_field_name, sfit, FREE_ENTRY_VALUE)
//	}
//}

// ================================================================
type Stats2R2Accumulator struct {
	count int
}

func NewStats2R2Accumulator() IStats2Accumulator {
	return &Stats2R2Accumulator{
		count: 0,
	}
}
func (this *Stats2R2Accumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2R2Accumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

//// ----------------------------------------------------------------
//// http://en.wikipedia.org/wiki/Pearson_product-moment_correlation_coefficient
//// Alternatively, just use sqrt(corr) as defined above.
//
//typedef struct _stats2_r2_state_t {
//	unsigned long long count
//	double sumx
//	double sumy
//	double sumx2
//	double sumxy
//	double sumy2
//	char*  r2_output_field_name
//} stats2_r2_state_t
//static void stats2_r2_ingest(void* pvstate, double x, double y) {
//	stats2_r2_state_t* this = pvstate
//	this.count++
//	this.sumx  += x
//	this.sumy  += y
//	this.sumx2 += x*x
//	this.sumxy += x*y
//	this.sumy2 += y*y
//}
//static void stats2_r2_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
//	stats2_r2_state_t* this = pvstate
//	if (this.count < 2LL) {
//		lrec_put(poutrec, this.r2_output_field_name, "", NO_FREE)
//	} else {
//		unsigned long long n = this.count
//		double sumx  = this.sumx
//		double sumy  = this.sumy
//		double sumx2 = this.sumx2
//		double sumy2 = this.sumy2
//		double sumxy = this.sumxy
//		double numerator = n*sumxy - sumx*sumy
//		numerator = numerator * numerator
//		double denominator = (n*sumx2 - sumx*sumx) * (n*sumy2 - sumy*sumy)
//		double output = numerator/denominator
//		char* val = mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt)
//		lrec_put(poutrec, this.r2_output_field_name, val, FREE_ENTRY_VALUE)
//	}
//}
//
//static stats2_acc_t* stats2_r2_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	stats2_acc_t* pstats2_acc = mlr_malloc_or_die(sizeof(stats2_acc_t))
//	stats2_r2_state_t* this = mlr_malloc_or_die(sizeof(stats2_r2_state_t))
//	this.count     = 0LL
//	this.sumx      = 0.0
//	this.sumy      = 0.0
//	this.sumx2     = 0.0
//	this.sumxy     = 0.0
//	this.sumy2     = 0.0
//	this.r2_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_r2")
//
//	pstats2_acc.pvstate      = (void*)this
//	pstats2_acc.pingest_func = stats2_r2_ingest
//	pstats2_acc.pemit_func   = stats2_r2_emit
//	pstats2_acc.pfit_func    = NULL
//
//	return pstats2_acc
//}

// ================================================================
type Stats2LogiRegAccumulator struct {
	count int
}

func NewStats2LogiRegAccumulator() IStats2Accumulator {
	return &Stats2LogiRegAccumulator{
		count: 0,
	}
}
func (this *Stats2LogiRegAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2LogiRegAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

//// ----------------------------------------------------------------
//#define LOGIREG_DVECTOR_INITIAL_SIZE 1024
//typedef struct _stats2_logireg_state_t {
//	dvector_t* pxs
//	dvector_t* pys
//	char*  m_output_field_name
//	char*  b_output_field_name
//	char*  n_output_field_name
//	char*  fit_output_field_name
//	int    fit_ready
//	double m
//	double b
//} stats2_logireg_state_t
//static void stats2_logireg_ingest(void* pvstate, double x, double y) {
//	stats2_logireg_state_t* this = pvstate
//	dvector_append(this.pxs, x)
//	dvector_append(this.pys, y)
//}
//
//static void stats2_logireg_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
//	stats2_logireg_state_t* this = pvstate
//
//	if (this.pxs.size < 2) {
//		lrec_put(poutrec, this.m_output_field_name, "", NO_FREE)
//		lrec_put(poutrec, this.b_output_field_name, "", NO_FREE)
//	} else {
//		double m, b
//		mlr_logistic_regression(this.pxs.data, this.pys.data, this.pxs.size, &m, &b)
//		char* mval = mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt)
//		char* bval = mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt)
//
//		lrec_put(poutrec, this.m_output_field_name, mval, FREE_ENTRY_VALUE)
//		lrec_put(poutrec, this.b_output_field_name, bval, FREE_ENTRY_VALUE)
//	}
//
//	char* nval = mlr_alloc_string_from_ll(this.pxs.size)
//	lrec_put(poutrec, this.n_output_field_name, nval, FREE_ENTRY_VALUE)
//}
//
//static void stats2_logireg_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
//	stats2_logireg_state_t* this = pvstate
//
//	if (!this.fit_ready) {
//		mlr_logistic_regression(this.pxs.data, this.pys.data, this.pxs.size, &this.m, &this.b)
//		this.fit_ready = TRUE
//	}
//
//	if (this.pxs.size < 2) {
//		lrec_put(poutrec, this.fit_output_field_name, "", NO_FREE)
//	} else {
//		double yfit = 1.0 / (1.0 + exp(-this.m*x - this.b))
//		char* fitval = mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt)
//		lrec_put(poutrec, this.fit_output_field_name, fitval, FREE_ENTRY_VALUE)
//	}
//}
//
//static stats2_acc_t* stats2_logireg_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	stats2_acc_t* pstats2_acc = mlr_malloc_or_die(sizeof(stats2_acc_t))
//	stats2_logireg_state_t* this = mlr_malloc_or_die(sizeof(stats2_logireg_state_t))
//	this.pxs = dvector_alloc(LOGIREG_DVECTOR_INITIAL_SIZE)
//	this.pys = dvector_alloc(LOGIREG_DVECTOR_INITIAL_SIZE)
//	this.m_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_m")
//	this.b_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_b")
//	this.n_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_n")
//	this.fit_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_fit")
//	this.fit_ready = FALSE
//	this.m         = -999.0
//	this.b         = -999.0
//
//	pstats2_acc.pvstate = (void*)this
//	pstats2_acc.pingest_func = stats2_logireg_ingest
//	pstats2_acc.pemit_func   = stats2_logireg_emit
//	pstats2_acc.pfit_func    = stats2_logireg_fit
//	return pstats2_acc
//}

// ================================================================
type Stats2CorrAccumulator struct {
	count int
}

func NewStats2CorrAccumulator() IStats2Accumulator {
	return &Stats2CorrAccumulator{
		count: 0,
	}
}
func (this *Stats2CorrAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2CorrAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

// ================================================================
type Stats2CovAccumulator struct {
	count int
}

func NewStats2CovAccumulator() IStats2Accumulator {
	return &Stats2CovAccumulator{
		count: 0,
	}
}
func (this *Stats2CovAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2CovAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

// ================================================================
type Stats2CovXAccumulator struct {
	count int
}

func NewStats2CovXAccumulator() IStats2Accumulator {
	return &Stats2CovXAccumulator{
		count: 0,
	}
}
func (this *Stats2CovXAccumulator) Ingest(value1, value2 *types.Mlrval) {
	this.count++
}
func (this *Stats2CovXAccumulator) Emit() types.Mlrval {
	return types.MlrvalFromInt(this.count)
}

// ================================================================
// ================================================================

//// ----------------------------------------------------------------
//// Corr(X,Y) = Cov(X,Y) / sigma_X sigma_Y.
//typedef struct _stats2_corr_cov_state_t {
//	unsigned long long count
//	double sumx
//	double sumy
//	double sumx2
//	double sumxy
//	double sumy2
//	bivar_measure_t do_which
//	int    do_verbose
//
//	char*  covx_00_output_field_name
//	char*  covx_01_output_field_name
//	char*  covx_10_output_field_name
//	char*  covx_11_output_field_name
//
//	char*  pca_m_output_field_name
//	char*  pca_b_output_field_name
//	char*  pca_n_output_field_name
//	char*  pca_q_output_field_name
//	char* pca_l1_output_field_name
//	char* pca_l2_output_field_name
//	char* pca_v11_output_field_name
//	char* pca_v12_output_field_name
//	char* pca_v21_output_field_name
//	char* pca_v22_output_field_name
//	char* pca_fit_output_field_name
//	int   fit_ready
//	double m
//	double b
//	double q
//
//	char*  corr_output_field_name
//	char*   cov_output_field_name
//
//} stats2_corr_cov_state_t
//static void stats2_corr_cov_ingest(void* pvstate, double x, double y) {
//	stats2_corr_cov_state_t* this = pvstate
//	this.count++
//	this.sumx  += x
//	this.sumy  += y
//	this.sumx2 += x*x
//	this.sumxy += x*y
//	this.sumy2 += y*y
//}
//
//static void stats2_corr_cov_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
//	stats2_corr_cov_state_t* this = pvstate
//	if (this.do_which == DO_COVX) {
//		char* key00 = this.covx_00_output_field_name
//		char* key01 = this.covx_01_output_field_name
//		char* key10 = this.covx_10_output_field_name
//		char* key11 = this.covx_11_output_field_name
//		if (this.count < 2LL) {
//			lrec_put(poutrec, key00, "", NO_FREE)
//			lrec_put(poutrec, key01, "", NO_FREE)
//			lrec_put(poutrec, key10, "", NO_FREE)
//			lrec_put(poutrec, key11, "", NO_FREE)
//		} else {
//			double Q[2][2]
//			mlr_get_cov_matrix(this.count,
//				this.sumx, this.sumx2, this.sumy, this.sumy2, this.sumxy, Q)
//			char* val00 = mlr_alloc_string_from_double(Q[0][0], MLR_GLOBALS.ofmt)
//			char* val01 = mlr_alloc_string_from_double(Q[0][1], MLR_GLOBALS.ofmt)
//			char* val10 = mlr_alloc_string_from_double(Q[1][0], MLR_GLOBALS.ofmt)
//			char* val11 = mlr_alloc_string_from_double(Q[1][1], MLR_GLOBALS.ofmt)
//			lrec_put(poutrec, key00, val00, FREE_ENTRY_VALUE)
//			lrec_put(poutrec, key01, val01, FREE_ENTRY_VALUE)
//			lrec_put(poutrec, key10, val10, FREE_ENTRY_VALUE)
//			lrec_put(poutrec, key11, val11, FREE_ENTRY_VALUE)
//		}
//
//	} else if (this.do_which == DO_LINREG_PCA) {
//		char* keym   = this.pca_m_output_field_name
//		char* keyb   = this.pca_b_output_field_name
//		char* keyn   = this.pca_n_output_field_name
//		char* keyq   = this.pca_q_output_field_name
//		char* keyl1  = this.pca_l1_output_field_name
//		char* keyl2  = this.pca_l2_output_field_name
//		char* keyv11 = this.pca_v11_output_field_name
//		char* keyv12 = this.pca_v12_output_field_name
//		char* keyv21 = this.pca_v21_output_field_name
//		char* keyv22 = this.pca_v22_output_field_name
//		if (this.count < 2LL) {
//			lrec_put(poutrec, keym,   "", NO_FREE)
//			lrec_put(poutrec, keyb,   "", NO_FREE)
//			lrec_put(poutrec, keyn,   "", NO_FREE)
//			lrec_put(poutrec, keyq,   "", NO_FREE)
//			if (this.do_verbose) {
//				lrec_put(poutrec, keyl1,  "", NO_FREE)
//				lrec_put(poutrec, keyl2,  "", NO_FREE)
//				lrec_put(poutrec, keyv11, "", NO_FREE)
//				lrec_put(poutrec, keyv12, "", NO_FREE)
//				lrec_put(poutrec, keyv21, "", NO_FREE)
//				lrec_put(poutrec, keyv22, "", NO_FREE)
//			}
//		} else {
//			double Q[2][2]
//			mlr_get_cov_matrix(this.count,
//				this.sumx, this.sumx2, this.sumy, this.sumy2, this.sumxy, Q)
//
//			double l1, l2;       // Eigenvalues
//			double v1[2], v2[2]; // Eigenvectors
//			mlr_get_real_symmetric_eigensystem(Q, &l1, &l2, v1, v2)
//
//			double x_mean = this.sumx / this.count
//			double y_mean = this.sumy / this.count
//			double m, b, q
//			mlr_get_linear_regression_pca(l1, l2, v1, v2, x_mean, y_mean, &m, &b, &q)
//
//			lrec_put(poutrec, keym, mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//			lrec_put(poutrec, keyb, mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//			lrec_put(poutrec, keyn, mlr_alloc_string_from_ll(this.count),           FREE_ENTRY_VALUE)
//			lrec_put(poutrec, keyq, mlr_alloc_string_from_double(q, MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//			if (this.do_verbose) {
//				lrec_put(poutrec, keyl1,  mlr_alloc_string_from_double(l1,    MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//				lrec_put(poutrec, keyl2,  mlr_alloc_string_from_double(l2,    MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//				lrec_put(poutrec, keyv11, mlr_alloc_string_from_double(v1[0], MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//				lrec_put(poutrec, keyv12, mlr_alloc_string_from_double(v1[1], MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//				lrec_put(poutrec, keyv21, mlr_alloc_string_from_double(v2[0], MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//				lrec_put(poutrec, keyv22, mlr_alloc_string_from_double(v2[1], MLR_GLOBALS.ofmt), FREE_ENTRY_VALUE)
//			}
//		}
//	} else {
//		char* key = (this.do_which == DO_CORR) ? this.corr_output_field_name : this.cov_output_field_name
//		if (this.count < 2LL) {
//			lrec_put(poutrec, key, "", NO_FREE)
//		} else {
//			double output = mlr_get_cov(this.count, this.sumx, this.sumy, this.sumxy)
//			if (this.do_which == DO_CORR) {
//				double sigmax = sqrt(mlr_get_var(this.count, this.sumx, this.sumx2))
//				double sigmay = sqrt(mlr_get_var(this.count, this.sumy, this.sumy2))
//				output = output / sigmax / sigmay
//			}
//			char* val = mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt)
//			lrec_put(poutrec, key, val, FREE_ENTRY_VALUE)
//		}
//	}
//}
//
//static void linreg_pca_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
//	stats2_corr_cov_state_t* this = pvstate
//
//	if (!this.fit_ready) {
//		double Q[2][2]
//		mlr_get_cov_matrix(this.count,
//			this.sumx, this.sumx2, this.sumy, this.sumy2, this.sumxy, Q)
//
//		double l1, l2;       // Eigenvalues
//		double v1[2], v2[2]; // Eigenvectors
//		mlr_get_real_symmetric_eigensystem(Q, &l1, &l2, v1, v2)
//
//		double x_mean = this.sumx / this.count
//		double y_mean = this.sumy / this.count
//		mlr_get_linear_regression_pca(l1, l2, v1, v2, x_mean, y_mean, &this.m, &this.b, &this.q)
//
//		this.fit_ready = TRUE
//	}
//	if (this.count < 2LL) {
//		lrec_put(poutrec, this.pca_fit_output_field_name, "", NO_FREE)
//	} else {
//		double yfit = this.m * x + this.b
//		lrec_put(poutrec, this.pca_fit_output_field_name, mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt),
//			FREE_ENTRY_VALUE)
//	}
//}
//
//static stats2_acc_t* stats2_corr_cov_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name,
//	bivar_measure_t do_which, int do_verbose)
//{
//	stats2_acc_t* pstats2_acc = mlr_malloc_or_die(sizeof(stats2_acc_t))
//	stats2_corr_cov_state_t* this = mlr_malloc_or_die(sizeof(stats2_corr_cov_state_t))
//	this.count      = 0LL
//	this.sumx       = 0.0
//	this.sumy       = 0.0
//	this.sumx2      = 0.0
//	this.sumxy      = 0.0
//	this.sumy2      = 0.0
//	this.do_which   = do_which
//	this.do_verbose = do_verbose
//
//	char* name1 = value_field_name_1
//	char* name2 = value_field_name_2
//
//	this.covx_00_output_field_name = mlr_paste_4_strings(name1, "_", name1, "_covx")
//	this.covx_01_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_covx")
//	this.covx_10_output_field_name = mlr_paste_4_strings(name2, "_", name1, "_covx")
//	this.covx_11_output_field_name = mlr_paste_4_strings(name2, "_", name2, "_covx")
//
//	this.pca_m_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_m")
//	this.pca_b_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_b")
//	this.pca_n_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_n")
//	this.pca_q_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_quality")
//	this.pca_l1_output_field_name  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival1")
//	this.pca_l2_output_field_name  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival2")
//	this.pca_v11_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec11")
//	this.pca_v12_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec12")
//	this.pca_v21_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec21")
//	this.pca_v22_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec22")
//	this.pca_fit_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_fit")
//	this.fit_ready = FALSE
//	this.m         = -999.0
//	this.b         = -999.0
//
//	this.corr_output_field_name    = mlr_paste_4_strings(name1, "_", name2, "_corr")
//	this.cov_output_field_name     = mlr_paste_4_strings(name1, "_", name2, "_cov")
//
//	pstats2_acc.pvstate      = (void*)this
//	pstats2_acc.pingest_func = stats2_corr_cov_ingest
//	pstats2_acc.pemit_func   = stats2_corr_cov_emit
//	if (do_which == DO_LINREG_PCA)
//		pstats2_acc.pfit_func = linreg_pca_fit
//	else
//		pstats2_acc.pfit_func = NULL
//
//	return pstats2_acc
//}

//static stats2_acc_t* stats2_corr_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_acc_name, DO_CORR, do_verbose)
//}
//static stats2_acc_t* stats2_cov_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_acc_name, DO_COV, do_verbose)
//}
//static stats2_acc_t* stats2_covx_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_acc_name, DO_COVX, do_verbose)
//}
//static stats2_acc_t* stats2_linreg_pca_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_acc_name, int do_verbose) {
//	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_acc_name, DO_LINREG_PCA, do_verbose)
//}
