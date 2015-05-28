#ifndef MLR_VAL_H
#define MLR_VAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
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
#define MT_DOUBLE 3
#define MT_INT    4
#define MT_STRING 5
#define MT_MAX    6

typedef struct _mlr_val_t {
	union {
		int        bval;
		double     dval;
		long long  ival;
		char*      string_val;
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

// ----------------------------------------------------------------
char* mt_describe_type(int type);

char* mt_describe_type(int type);
char* mt_format_val(mv_t* pval); // xxx cmt mem-mgt
char* mt_describe_val(mv_t val);
// xxx explain why one is void & the other isn't
int  mt_get_boolean_strict(mv_t* pval);
void mt_get_double_strict(mv_t* pval);

// ----------------------------------------------------------------
typedef mv_t mv_zary_func_t();
typedef mv_t mv_unary_func_t(mv_t* pval1);
typedef mv_t mv_binary_func_t(mv_t* pval1, mv_t* pval2);
typedef mv_t mv_ternary_func_t(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
static inline mv_t b_b_not_func(mv_t* pval1) {
	mv_t rv = {.type = MT_BOOL, .u.bval = !pval1->u.bval};
	return rv;
}

static inline mv_t b_bb_or_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_BOOL, .u.bval = pval1->u.bval || pval2->u.bval};
	return rv;
}
static inline mv_t b_bb_and_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_BOOL, .u.bval = pval1->u.bval && pval2->u.bval};
	return rv;
}

// ----------------------------------------------------------------
static inline mv_t f_z_urand_func() {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = get_mtrand_double()}; // mtrand.h
	return rv;
}
static inline mv_t f_z_systime_func() {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = get_systime()}; // mlrutil.h
	return rv;
}

// ----------------------------------------------------------------
static inline mv_t f_f_uneg_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = -pval1->u.dval};
	return rv;
}
static inline mv_t f_f_abs_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = fabs(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_log_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = log(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_log10_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = log10(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_exp_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = exp(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_sin_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = sin(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_cos_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = cos(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_tan_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = tan(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_sqrt_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = sqrt(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_round_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = round(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_floor_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = floor(pval1->u.dval)};
	return rv;
}
static inline mv_t f_f_ceil_func(mv_t* pval1) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = ceil(pval1->u.dval)};
	return rv;
}

// ----------------------------------------------------------------
static inline mv_t f_ff_plus_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = pval1->u.dval + pval2->u.dval};
	return rv;
}
static inline mv_t f_ff_minus_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = pval1->u.dval - pval2->u.dval};
	return rv;
}
static inline mv_t f_ff_times_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = pval1->u.dval * pval2->u.dval};
	return rv;
}
static inline mv_t f_ff_divide_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = pval1->u.dval / pval2->u.dval};
	return rv;
}
static inline mv_t f_ff_pow_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = pow(pval1->u.dval, pval2->u.dval)};
	return rv;
}
static inline mv_t f_ff_mod_func(mv_t* pval1, mv_t* pval2) {
	long long i1 = (long long)pval1->u.dval;
	long long i2 = (long long)pval2->u.dval;
	long long i3 = i1 % i2;
	if (i3 < 0)
		i3 += i2; // C mod is insane
	mv_t rv = {.type = MT_DOUBLE, .u.dval = (double)i3};
	return rv;
}
static inline mv_t f_ff_atan2_func(mv_t* pval1, mv_t* pval2) {
	mv_t rv = {.type = MT_DOUBLE, .u.dval = atan2(pval1->u.dval, pval2->u.dval)};
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
mv_t i_s_strlen_func(mv_t* pval1);

// ----------------------------------------------------------------
mv_t eq_op_func(mv_t* pval1, mv_t* pval2);
mv_t ne_op_func(mv_t* pval1, mv_t* pval2);
mv_t gt_op_func(mv_t* pval1, mv_t* pval2);
mv_t ge_op_func(mv_t* pval1, mv_t* pval2);
mv_t lt_op_func(mv_t* pval1, mv_t* pval2);
mv_t le_op_func(mv_t* pval1, mv_t* pval2);

#endif // MLR_VAL_H
