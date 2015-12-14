#ifndef MLR_VAL_H
#define MLR_VAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
#include <regex.h>
#include "lib/mlrmath.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "lib/string_builder.h"
#include "containers/free_flags.h"

// ================================================================
// MT for Miller type -- highly abbreviated here since these are
// spelled out a lot in lrec_evaluators.c.
// ================================================================

// Among other things, these defines are used in mlrval.c to index disposition matrices.
// So if the numeric values are changed, the matrices must be as well.
#define MT_NULL   0 // E.g. field name not present in input record -- not a problem.
#define MT_ERROR  1 // E.g. error encountered in one eval & it propagates up the AST.
#define MT_BOOL   2
#define MT_FLOAT  3
#define MT_INT    4
#define MT_STRING 5
#define MT_MAX    6

#define MV_SB_ALLOC_LENGTH 32

#define ISO8601_TIME_FORMAT "%Y-%m-%dT%H:%M:%SZ"

typedef struct _mv_t {
	union {
		int        boolv;
		double     fltv;
		long long  intv;
		char*      strv;
	} u;
	unsigned char type;
	unsigned char free_flags;
} mv_t;

// ----------------------------------------------------------------
extern mv_t MV_NULL;
extern mv_t MV_ERROR;

#define NULL_OR_ERROR_OUT(val) { \
	if ((val).type == MT_ERROR) \
		return MV_ERROR; \
	if ((val).type == MT_NULL) \
		return MV_NULL; \
}

#define NULL_OUT(val) { \
	if ((val).type == MT_NULL) \
		return MV_NULL; \
}
#define ERROR_OUT(val) { \
	if ((val).type == MT_ERROR) \
		return MV_ERROR; \
}

// ----------------------------------------------------------------
static inline mv_t mv_from_float(double d) {
	return (mv_t) {.type = MT_FLOAT, .free_flags = NO_FREE, .u.fltv = d};
}

static inline mv_t mv_from_int(long long i) {
	return (mv_t) {.type = MT_INT, .free_flags = NO_FREE, .u.intv = i};
}

static inline mv_t mv_from_bool(int b) {
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = b};
}
static inline mv_t mv_from_true() {
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = TRUE};
}
static inline mv_t mv_from_false() {
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = FALSE};
}

static inline mv_t mv_from_null() {
	return (mv_t) {.type = MT_NULL, .free_flags = NO_FREE, .u.intv = 0LL};
}

static inline mv_t mv_from_string_with_free(char* s) {
	return (mv_t) {.type = MT_STRING, .free_flags = FREE_ENTRY_KEY, .u.strv = s};
}

static inline mv_t mv_from_string_no_free(char* s) {
	return (mv_t) {.type = MT_STRING, .free_flags = NO_FREE, .u.strv = s};
}

// ----------------------------------------------------------------
static inline int mv_is_numeric(mv_t* pval) {
	return pval->type == MT_INT || pval->type == MT_FLOAT;
}

static inline int mv_is_null(mv_t* pval) {
	return pval->type == MT_NULL;
}

// ----------------------------------------------------------------
char* mt_describe_type(int type);

char* mt_describe_type(int type);
char* mv_format_val(mv_t* pval); // For debug only; the caller should free the return value
char* mv_describe_val(mv_t val);
void mv_get_boolean_strict(mv_t* pval);
void mv_get_float_strict(mv_t* pval);
void mv_get_float_nullable(mv_t* pval);
void mv_get_int_nullable(mv_t* pval);

// int or float:
void mv_get_number_nullable(mv_t* pval);
mv_t mv_scan_number_nullable(char* string);
mv_t mv_scan_number_or_die(char* string);

// ----------------------------------------------------------------
typedef mv_t mv_zary_func_t();
typedef mv_t mv_unary_func_t(mv_t* pval1);
typedef mv_t mv_binary_func_t(mv_t* pval1, mv_t* pval2);
typedef mv_t mv_binary_arg2_regex_func_t(mv_t* pval1, regex_t* pregex, string_builder_t* psb);
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

mv_t n_nn_plus_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_minus_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_times_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_divide_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_int_divide_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_mod_func(mv_t* pval1, mv_t* pval2);
mv_t n_n_upos_func(mv_t* pval1);
mv_t n_n_uneg_func(mv_t* pval1);

mv_t n_n_abs_func(mv_t* pval1);
mv_t n_n_ceil_func(mv_t* pval1);
mv_t n_n_floor_func(mv_t* pval1);
mv_t n_n_round_func(mv_t* pval1);
mv_t n_n_sgn_func(mv_t* pval1);

mv_t n_nn_min_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_max_func(mv_t* pval1, mv_t* pval2);
mv_t n_nn_roundm_func(mv_t* pval1, mv_t* pval2);

mv_t i_x_int_func(mv_t* pval1);
mv_t f_x_float_func(mv_t* pval1);
mv_t b_x_boolean_func(mv_t* pval1);
mv_t s_x_string_func(mv_t* pval1);
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

static inline mv_t i_ii_bitwise_or_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv | pval2->u.intv);
}
static inline mv_t i_ii_bitwise_xor_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv ^ pval2->u.intv);
}
static inline mv_t i_ii_bitwise_and_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv & pval2->u.intv);
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

mv_t i_iii_modadd_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modsub_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modmul_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modexp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_s_tolower_func(mv_t* pval1);
mv_t s_s_toupper_func(mv_t* pval1);

mv_t s_ss_dot_func(mv_t* pval1, mv_t* pval2);

mv_t sub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t sub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);
mv_t gsub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t gsub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_n_sec2gmt_func(mv_t* pval1);
mv_t i_s_gmt2sec_func(mv_t* pval1);
mv_t s_ns_strftime_func(mv_t* pval1, mv_t* pval2);
mv_t i_ss_strptime_func(mv_t* pval1, mv_t* pval2);

mv_t s_i_sec2hms_func(mv_t* pval1);
mv_t s_f_fsec2hms_func(mv_t* pval1);
mv_t s_i_sec2dhms_func(mv_t* pval1);
mv_t s_f_fsec2dhms_func(mv_t* pval1);
mv_t i_s_hms2sec_func(mv_t* pval1);
mv_t f_s_hms2fsec_func(mv_t* pval1);
mv_t i_s_dhms2sec_func(mv_t* pval1);
mv_t f_s_dhms2fsec_func(mv_t* pval1);

mv_t time_string_from_seconds(mv_t* psec, char* format);

mv_t i_s_strlen_func(mv_t* pval1);

// ----------------------------------------------------------------
// arg2 evaluates to string via compound expression; regexes compiled on each call
mv_t matches_no_precomp_func(mv_t* pval1, mv_t* pval2);
mv_t does_not_match_no_precomp_func(mv_t* pval1, mv_t* pval2);
// arg2 is a string, compiled to regex only once at alloc time
mv_t matches_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb);
mv_t does_not_match_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb);

// For filter/put DSLs:
mv_t eq_op_func(mv_t* pval1, mv_t* pval2);
mv_t ne_op_func(mv_t* pval1, mv_t* pval2);
mv_t gt_op_func(mv_t* pval1, mv_t* pval2);
mv_t ge_op_func(mv_t* pval1, mv_t* pval2);
mv_t lt_op_func(mv_t* pval1, mv_t* pval2);
mv_t le_op_func(mv_t* pval1, mv_t* pval2);

// For non-DSL comparison of mlrvals:
int mv_i_nn_eq(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ne(mv_t* pval1, mv_t* pval2);
int mv_i_nn_gt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ge(mv_t* pval1, mv_t* pval2);
int mv_i_nn_lt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_le(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
// For qsort of numeric mlrvals.
int mv_nn_comparator(const void* pva, const void* pvb);

int mlr_bsearch_mv_n_for_insert(mv_t* array, int size, mv_t* pvalue);

#endif // MLR_VAL_H
