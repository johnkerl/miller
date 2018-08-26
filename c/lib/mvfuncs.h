#ifndef MVFUNCS_H
#define MVFUNCS_H

// ================================================================
// Functions involving mlrvals: primarily for the DSL but also for
// stats1/stats2, etc.
// ================================================================

#include "../lib/mlrmath.h"
#include "../lib/mlrutil.h"
#include "../lib/mlrdatetime.h"
#include "../lib/mtrand.h"
#include "../lib/string_builder.h"
#include "../lib/string_array.h"
#include "../lib/mlrval.h"

#define MV_SB_ALLOC_LENGTH 32

#define ISO8601_TIME_FORMAT   "%Y-%m-%dT%H:%M:%SZ"
#define ISO8601_TIME_FORMAT_1 "%Y-%m-%dT%H:%M:%1SZ"
#define ISO8601_TIME_FORMAT_2 "%Y-%m-%dT%H:%M:%2SZ"
#define ISO8601_TIME_FORMAT_3 "%Y-%m-%dT%H:%M:%3SZ"
#define ISO8601_TIME_FORMAT_4 "%Y-%m-%dT%H:%M:%4SZ"
#define ISO8601_TIME_FORMAT_5 "%Y-%m-%dT%H:%M:%5SZ"
#define ISO8601_TIME_FORMAT_6 "%Y-%m-%dT%H:%M:%6SZ"
#define ISO8601_TIME_FORMAT_7 "%Y-%m-%dT%H:%M:%7SZ"
#define ISO8601_TIME_FORMAT_8 "%Y-%m-%dT%H:%M:%8SZ"
#define ISO8601_TIME_FORMAT_9 "%Y-%m-%dT%H:%M:%9SZ"
#define ISO8601_DATE_FORMAT   "%Y-%m-%d"

#define ISO8601_LOCAL_TIME_FORMAT   "%Y-%m-%d %H:%M:%S"
#define ISO8601_LOCAL_TIME_FORMAT_1 "%Y-%m-%d %H:%M:%1S"
#define ISO8601_LOCAL_TIME_FORMAT_2 "%Y-%m-%d %H:%M:%2S"
#define ISO8601_LOCAL_TIME_FORMAT_3 "%Y-%m-%d %H:%M:%3S"
#define ISO8601_LOCAL_TIME_FORMAT_4 "%Y-%m-%d %H:%M:%4S"
#define ISO8601_LOCAL_TIME_FORMAT_5 "%Y-%m-%d %H:%M:%5S"
#define ISO8601_LOCAL_TIME_FORMAT_6 "%Y-%m-%d %H:%M:%6S"
#define ISO8601_LOCAL_TIME_FORMAT_7 "%Y-%m-%d %H:%M:%7S"
#define ISO8601_LOCAL_TIME_FORMAT_8 "%Y-%m-%d %H:%M:%8S"
#define ISO8601_LOCAL_TIME_FORMAT_9 "%Y-%m-%d %H:%M:%9S"

// ----------------------------------------------------------------
typedef mv_t mv_variadic_func_t(mv_t* pvals, int nvals);
typedef mv_t mv_zary_func_t();
typedef mv_t mv_unary_func_t(mv_t* pval1);
typedef mv_t mv_binary_func_t(mv_t* pval1, mv_t* pval2);
typedef mv_t mv_binary_arg3_capture_func_t(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
typedef mv_t mv_binary_arg2_regex_func_t(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);
typedef mv_t mv_binary_arg2_regex_extract_func_t(mv_t* pval1, regex_t* pregex);
typedef mv_t mv_ternary_func_t(mv_t* pval1, mv_t* pval2, mv_t* pval3);
typedef mv_t mv_ternary_arg2_regex_func_t(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);

// ----------------------------------------------------------------
static inline mv_t b_b_not_func(mv_t* pval1) {
	return mv_from_bool(!pval1->u.boolv);
}

static inline mv_t b_bb_or_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv || pval2->u.boolv);
}
static inline mv_t b_bb_and_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv && pval2->u.boolv);
}
static inline mv_t b_bb_xor_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv ^ pval2->u.boolv);
}

// ----------------------------------------------------------------
static inline mv_t f_z_urand_func() {
	return mv_from_float(get_mtrand_double()); // mtrand.h
}
static inline mv_t i_z_urand32_func() {
	return mv_from_float(get_mtrand_int32()); // mtrand.h
}
static inline mv_t f_z_systime_func() {
	return mv_from_float(get_systime()); // mlrutil.h
}

// ----------------------------------------------------------------
static inline mv_t f_f_acos_func(mv_t*     pval1) {return mv_from_float( acos     (pval1->u.fltv));}
static inline mv_t f_f_acosh_func(mv_t*    pval1) {return mv_from_float( acosh    (pval1->u.fltv));}
static inline mv_t f_f_asin_func(mv_t*     pval1) {return mv_from_float( asin     (pval1->u.fltv));}
static inline mv_t f_f_asinh_func(mv_t*    pval1) {return mv_from_float( asinh    (pval1->u.fltv));}
static inline mv_t f_f_atan_func(mv_t*     pval1) {return mv_from_float( atan     (pval1->u.fltv));}
static inline mv_t f_f_atanh_func(mv_t*    pval1) {return mv_from_float( atanh    (pval1->u.fltv));}
static inline mv_t f_f_cbrt_func(mv_t*     pval1) {return mv_from_float( cbrt     (pval1->u.fltv));}
static inline mv_t f_f_cos_func(mv_t*      pval1) {return mv_from_float( cos      (pval1->u.fltv));}
static inline mv_t f_f_cosh_func(mv_t*     pval1) {return mv_from_float( cosh     (pval1->u.fltv));}
static inline mv_t f_f_erf_func(mv_t*      pval1) {return mv_from_float( erf      (pval1->u.fltv));}
static inline mv_t f_f_erfc_func(mv_t*     pval1) {return mv_from_float( erfc     (pval1->u.fltv));}
static inline mv_t f_f_exp_func(mv_t*      pval1) {return mv_from_float( exp      (pval1->u.fltv));}
static inline mv_t f_f_expm1_func(mv_t*    pval1) {return mv_from_float( expm1    (pval1->u.fltv));}
static inline mv_t f_f_invqnorm_func(mv_t* pval1) {return mv_from_float( invqnorm (pval1->u.fltv));}
static inline mv_t f_f_log10_func(mv_t*    pval1) {return mv_from_float( log10    (pval1->u.fltv));}
static inline mv_t f_f_log1p_func(mv_t*    pval1) {return mv_from_float( log1p    (pval1->u.fltv));}
static inline mv_t f_f_log_func(mv_t*      pval1) {return mv_from_float( log      (pval1->u.fltv));}
static inline mv_t f_f_qnorm_func(mv_t*    pval1) {return mv_from_float( qnorm    (pval1->u.fltv));}
static inline mv_t f_f_sin_func(mv_t*      pval1) {return mv_from_float( sin      (pval1->u.fltv));}
static inline mv_t f_f_sinh_func(mv_t*     pval1) {return mv_from_float( sinh     (pval1->u.fltv));}
static inline mv_t f_f_sqrt_func(mv_t*     pval1) {return mv_from_float( sqrt     (pval1->u.fltv));}
static inline mv_t f_f_tan_func(mv_t*      pval1) {return mv_from_float( tan      (pval1->u.fltv));}
static inline mv_t f_f_tanh_func(mv_t*     pval1) {return mv_from_float( tanh     (pval1->u.fltv));}

static inline mv_t f_ff_pow_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_float(pow(pval1->u.fltv, pval2->u.fltv));
}

// These four overflow from 64-bit ints to double. This is for general use.
mv_t x_xx_plus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_minus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_times_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_divide_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_int_divide_func(mv_t* pval1, mv_t* pval2);

// These four intentionally overflow 64-bit ints. This is for use-cases where
// people want that, e.g. 64-bit integer math.
mv_t x_xx_oplus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_ominus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_otimes_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_odivide_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_int_odivide_func(mv_t* pval1, mv_t* pval2);

mv_t x_xx_mod_func(mv_t* pval1, mv_t* pval2);
mv_t x_x_upos_func(mv_t* pval1);
mv_t x_x_uneg_func(mv_t* pval1);

// Bitwise
mv_t x_xx_bxor_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_band_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_bor_func(mv_t* pval1, mv_t* pval2);

mv_t x_x_abs_func(mv_t* pval1);
mv_t x_x_ceil_func(mv_t* pval1);
mv_t x_x_floor_func(mv_t* pval1);
mv_t x_x_round_func(mv_t* pval1);
mv_t x_x_sgn_func(mv_t* pval1);

mv_t variadic_min_func(mv_t* pvals, int nvals);
mv_t variadic_max_func(mv_t* pvals, int nvals);

mv_t x_xx_min_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_max_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_roundm_func(mv_t* pval1, mv_t* pval2);

mv_t i_x_int_func(mv_t* pval1);
mv_t f_x_float_func(mv_t* pval1);
mv_t b_x_boolean_func(mv_t* pval1);
mv_t s_x_string_func(mv_t* pval1);
mv_t s_sii_substr_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t s_x_hexfmt_func(mv_t* pval1);
mv_t s_xs_fmtnum_func(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
static inline mv_t f_ff_atan2_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_float(atan2(pval1->u.fltv, pval2->u.fltv));
}

static inline mv_t f_fff_logifit_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	double x = pval1->u.fltv;
	double m = pval2->u.fltv;
	double b = pval3->u.fltv;
	return mv_from_float(1.0 / (1.0 + exp(-m*x-b)));
}

static inline mv_t i_ii_urandint_func(mv_t* pval1, mv_t* pval2) {
	long long a = pval1->u.intv;
	long long b = pval2->u.intv;
	long long lo, hi;
	if (a <= b) {
		lo = a;
		hi = b + 1;
	} else {
		lo = b;
		hi = a + 1;
	}
	long long u  = lo + (hi - lo) * get_mtrand_double();
	return mv_from_int(u);
}

static inline mv_t i_ii_bitwise_lsh_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv << pval2->u.intv);
}
static inline mv_t i_ii_bitwise_rsh_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv >> pval2->u.intv);
}
static inline mv_t i_i_bitwise_not_func(mv_t* pval1) {
	return mv_from_int(~pval1->u.intv);
}
mv_t i_i_bitcount_func(mv_t* pval1);

mv_t i_iii_modadd_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modsub_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modmul_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modexp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_s_tolower_func(mv_t* pval1);
mv_t s_s_toupper_func(mv_t* pval1);
mv_t i_s_strlen_func(mv_t* pval1);
mv_t s_x_typeof_func(mv_t* pval1);

mv_t s_xx_dot_func(mv_t* pval1, mv_t* pval2);

mv_t sub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t sub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);
mv_t gsub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t gsub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);
mv_t regex_extract_no_precomp_func(mv_t* pval1, mv_t* pval2);
mv_t regex_extract_precomp_func(mv_t* pval1, regex_t* pregex);
// String-substitution with no regexes or special characters.
mv_t s_sss_ssub_func(mv_t* pstring, mv_t* pold, mv_t* pnew);

// ----------------------------------------------------------------
mv_t s_x_sec2gmt_func(mv_t* pval1);
mv_t s_xi_sec2gmt_func(mv_t* pval1, mv_t* pval2);
mv_t s_x_sec2gmtdate_func(mv_t* pval1);

mv_t s_x_sec2localtime_func(mv_t* pval1);
mv_t s_xi_sec2localtime_func(mv_t* pval1, mv_t* pval2);
mv_t s_x_sec2localdate_func(mv_t* pval1);

mv_t i_s_gmt2sec_func(mv_t* pval1);
mv_t i_s_localtime2sec_func(mv_t* pval1);

mv_t s_ns_strftime_func(mv_t* pval1, mv_t* pval2);
mv_t s_ns_strftime_local_func(mv_t* pval1, mv_t* pval2);

mv_t i_ss_strptime_func(mv_t* pval1, mv_t* pval2);
mv_t i_ss_strptime_local_func(mv_t* pval1, mv_t* pval2);

mv_t s_i_sec2hms_func(mv_t* pval1);
mv_t s_f_fsec2hms_func(mv_t* pval1);
mv_t s_i_sec2dhms_func(mv_t* pval1);
mv_t s_f_fsec2dhms_func(mv_t* pval1);
mv_t i_s_hms2sec_func(mv_t* pval1);
mv_t f_s_hms2fsec_func(mv_t* pval1);
mv_t i_s_dhms2sec_func(mv_t* pval1);
mv_t f_s_dhms2fsec_func(mv_t* pval1);

mv_t time_string_from_seconds(mv_t* psec, char* format,
	timezone_handling_t timezone_handling);

// ----------------------------------------------------------------
// arg2 evaluates to string via compound expression; regexes compiled on each call
mv_t matches_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
mv_t does_not_match_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
// arg2 is a string, compiled to regex only once at alloc time
mv_t matches_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);
mv_t does_not_match_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);

// For filter/put DSL:
mv_t eq_op_func(mv_t* pval1, mv_t* pval2);
mv_t ne_op_func(mv_t* pval1, mv_t* pval2);
mv_t gt_op_func(mv_t* pval1, mv_t* pval2);
mv_t ge_op_func(mv_t* pval1, mv_t* pval2);
mv_t lt_op_func(mv_t* pval1, mv_t* pval2);
mv_t le_op_func(mv_t* pval1, mv_t* pval2);

// Assumes inputs are MT_STRING or MT_INT. Nominally intended for mlhmmv which uses only string/int mlrvals.
int mv_equals_si(mv_t* pa, mv_t* pb);

// For non-DSL comparison of mlrvals:
int mv_i_nn_eq(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ne(mv_t* pval1, mv_t* pval2);
int mv_i_nn_gt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ge(mv_t* pval1, mv_t* pval2);
int mv_i_nn_lt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_le(mv_t* pval1, mv_t* pval2);

// For unit-test keystroke-saving:
int mveq(mv_t* pval1, mv_t* pval2);
int mvne(mv_t* pval1, mv_t* pval2);
int mveqcopy(mv_t* pval1, mv_t* pval2);
int mvnecopy(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
// For qsort of numeric mlrvals.
int mv_nn_comparator(const void* pva, const void* pvb);

// For qsort of arbitrary mlrvals. Sort rules:
// * Across types:
//   NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * Within types:
//   o numeric compares on numbers
//   o false < true
//   o string compares on strings
//   o error == error (this is a singleton type)
//   o absent == absent (this is a singleton type)
int mv_xx_comparator(const void* pva, const void* pvb);

int mlr_bsearch_mv_n_for_insert(mv_t* array, int size, mv_t* pvalue);

#endif // MVFUNCS_H
