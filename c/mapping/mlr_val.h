#ifndef MLR_VAL_H
#define MLR_VAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
#include <regex.h>
#include "lib/mlrmath.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"

// ================================================================
// MT for Miller type -- highly abbreviated here since these are
// spelled out a lot in lrec_evaluators.c.
// ================================================================

// Among other things, these defines are used in mlr_val.c to index disposition matrices.
// So if the numeric values are changed, the matrices must be as well.
#define MT_NULL   0 // E.g. field name not present in input record -- not a problem.
#define MT_ERROR  1 // E.g. error encountered in one eval & it propagates up the AST.
#define MT_BOOL   2
#define MT_FLOAT  3
#define MT_INT    4
#define MT_STRING 5
#define MT_MAX    6

typedef struct _mlr_val_t {
	union {
		int        boolv;
		double     fltv;
		long long  intv;
		char*      strv;
	} u;
	unsigned char type;
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
char* mt_describe_type(int type);

char* mt_describe_type(int type);
char* mt_format_val(mv_t* pval); // For debug only; the caller should free the return value
char* mt_describe_val(mv_t val);
void mt_get_boolean_strict(mv_t* pval);
void mt_get_double_strict(mv_t* pval);
void mt_get_double_nullable(mv_t* pval);

// ----------------------------------------------------------------
typedef mv_t mv_zary_func_t();
typedef mv_t mv_unary_func_t(mv_t* pval1);
typedef mv_t mv_binary_func_t(mv_t* pval1, mv_t* pval2);
typedef mv_t mv_binary_right_regex_func_t(mv_t* pval1, regex_t* pregex);
typedef mv_t mv_ternary_func_t(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
static inline mv_t b_b_not_func(mv_t* pval1) {
	mv_t rv = {.type = MT_BOOL, .u.boolv = !pval1->u.boolv};
	return rv;
}

static inline mv_t b_bb_or_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_BOOL, .u.boolv = pval1->u.boolv || pval2->u.boolv};
	return rv;
}
static inline mv_t b_bb_and_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_BOOL, .u.boolv = pval1->u.boolv && pval2->u.boolv};
	return rv;
}

// ----------------------------------------------------------------
static inline mv_t f_z_urand_func() {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = get_mtrand_double()}; // mtrand.h
	return rv;
}
static inline mv_t f_z_systime_func() {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = get_systime()}; // mlrutil.h
	return rv;
}

// ----------------------------------------------------------------
static inline mv_t f_f_abs_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=fabs(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_acos_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=acos(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_acosh_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=acosh(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_asin_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=asin(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_asinh_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=asinh(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_atan_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=atan(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_atanh_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=atanh(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_cbrt_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=cbrt(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_ceil_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=ceil(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_cos_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=cos(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_cosh_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=cosh(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_erf_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=erf(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_erfc_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=erfc(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_exp_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=exp(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_expm1_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=expm1(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_floor_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=floor(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_invqnorm_func(mv_t* pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=invqnorm(pval1->u.fltv)}; return rv;}
static inline mv_t f_f_log10_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=log10(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_log1p_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=log1p(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_log_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=log(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_qnorm_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=qnorm(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_round_func(mv_t*    pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=round(pval1->u.fltv)};    return rv;}
static inline mv_t f_f_sin_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=sin(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_sinh_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=sinh(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_sqrt_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=sqrt(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_tan_func(mv_t*      pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=tan(pval1->u.fltv)};      return rv;}
static inline mv_t f_f_tanh_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=tanh(pval1->u.fltv)};     return rv;}
static inline mv_t f_f_uneg_func(mv_t*     pval1){mv_t rv={.type=MT_FLOAT,.u.fltv=-pval1->u.fltv};          return rv;}

mv_t i_x_int_func(mv_t* pval1);
mv_t f_x_float_func(mv_t* pval1);
mv_t b_x_boolean_func(mv_t* pval1);
mv_t s_x_string_func(mv_t* pval1);
mv_t s_x_hexfmt_func(mv_t* pval1);
mv_t s_xs_fmtnum_func(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
static inline mv_t f_ff_plus_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = pval1->u.fltv + pval2->u.fltv};
	return rv;
}
static inline mv_t f_ff_minus_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = pval1->u.fltv - pval2->u.fltv};
	return rv;
}
static inline mv_t f_ff_times_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = pval1->u.fltv * pval2->u.fltv};
	return rv;
}
static inline mv_t f_ff_divide_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = pval1->u.fltv / pval2->u.fltv};
	return rv;
}
static inline mv_t f_ff_max_func(mv_t* pval1, mv_t* pval2) {
	if (pval1->type == MT_NULL) {
		return *pval2;
	} else if (pval2->type == MT_NULL) {
		return *pval1;
	} else {
		mv_t rv = {.type = MT_FLOAT, .u.fltv = fmax(pval1->u.fltv, pval2->u.fltv)};
		return rv;
	}
}
static inline mv_t f_ff_min_func(mv_t* pval1, mv_t* pval2) {
	if (pval1->type == MT_NULL) {
		return *pval2;
	} else if (pval2->type == MT_NULL) {
		return *pval1;
	} else {
		mv_t rv = {.type = MT_FLOAT, .u.fltv = fmin(pval1->u.fltv, pval2->u.fltv)};
		return rv;
	}
}
static inline mv_t f_ff_pow_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = pow(pval1->u.fltv, pval2->u.fltv)};
	return rv;
}
static inline mv_t f_ff_mod_func(mv_t* pval1, mv_t* pval2) {
	long long i1 = (long long)pval1->u.fltv;
	long long i2 = (long long)pval2->u.fltv;
	long long i3 = i1 % i2;
	if (i3 < 0)
		i3 += i2; // C mod is insane
	mv_t rv = {.type = MT_FLOAT, .u.fltv = (double)i3};
	return rv;
}

static inline mv_t f_ff_atan2_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_FLOAT, .u.fltv = atan2(pval1->u.fltv, pval2->u.fltv)};
	return rv;
}

static inline mv_t f_ff_roundm_func(mv_t* pval1, mv_t* pval2) {
	double x = pval1->u.fltv;
	double m = pval2->u.fltv;
	mv_t rv = {.type = MT_FLOAT, .u.fltv = round(x / m) * m};
	return rv;
}

// ----------------------------------------------------------------
mv_t s_s_tolower_func(mv_t* pval1);
mv_t s_s_toupper_func(mv_t* pval1);

mv_t s_ss_dot_func(mv_t* pval1, mv_t* pval2);

mv_t s_sss_sub_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_f_sec2gmt_func(mv_t* pval1);
mv_t i_s_gmt2sec_func(mv_t* pval1);

mv_t s_i_sec2hms_func(mv_t* pval1);
mv_t s_f_fsec2hms_func(mv_t* pval1);
mv_t s_i_sec2dhms_func(mv_t* pval1);
mv_t s_f_fsec2dhms_func(mv_t* pval1);
mv_t i_s_hms2sec_func(mv_t* pval1);
mv_t f_s_hms2fsec_func(mv_t* pval1);
mv_t i_s_dhms2sec_func(mv_t* pval1);
mv_t f_s_dhms2fsec_func(mv_t* pval1);

mv_t i_s_strlen_func(mv_t* pval1);

// ----------------------------------------------------------------
// arg2 evaluates to string via compound expression; regexes compiled on each call
mv_t matches_no_precomp_func(mv_t* pval1, mv_t* pval2);
mv_t does_not_match_no_precomp_func(mv_t* pval1, mv_t* pval2);
// arg2 is a string, compiled to regex only once at alloc time
mv_t matches_precomp_func(mv_t* pval1, regex_t* pregex);
mv_t does_not_match_precomp_func(mv_t* pval1, regex_t* pregex);

mv_t eq_op_func(mv_t* pval1, mv_t* pval2);
mv_t ne_op_func(mv_t* pval1, mv_t* pval2);
mv_t gt_op_func(mv_t* pval1, mv_t* pval2);
mv_t ge_op_func(mv_t* pval1, mv_t* pval2);
mv_t lt_op_func(mv_t* pval1, mv_t* pval2);
mv_t le_op_func(mv_t* pval1, mv_t* pval2);

#endif // MLR_VAL_H
