#ifndef MLR_VAL_H
#define MLR_VAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
#include "lib/mlrutil.h"
#include "lib/mtrand.h"

// MT for Miller type -- highly abbreviated here since these are
// spelled out a lot in lrec_evaluators.c.

// Among other things, these are used in mlr_val.c to index disposition matrices.
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
		int        bool_val;
		double     double_val;
		long long  int_val;
		char*      string_val;
	} u;
	unsigned char type;
} mlr_val_t;

// ----------------------------------------------------------------
extern mlr_val_t MV_NULL;
extern mlr_val_t MV_ERROR;

#define NULL_OR_ERROR_OUT(val) { \
	if (val.type == MT_ERROR) \
		return MV_ERROR; \
	if (val.type == MT_NULL) \
		return MV_NULL; \
}

// ----------------------------------------------------------------
char* mt_describe_type(int type);

// xxx cmt mem-mgt
char* mt_format_val(mlr_val_t* pval);

char* mt_describe_type(int type);
char* mt_format_val(mlr_val_t* pval);
char* mt_describe_val(mlr_val_t val);
int mt_get_boolean_strict(mlr_val_t* pval);
void mt_get_double_strict(mlr_val_t* pval);
int mt_get_boolean_strict(mlr_val_t* pval);

// ----------------------------------------------------------------
typedef mlr_val_t mv_zary_func_t();
typedef mlr_val_t mv_unary_func_t(mlr_val_t* pval1);
typedef mlr_val_t mv_binary_func_t(mlr_val_t* pval1, mlr_val_t* pval2);

// ----------------------------------------------------------------
static inline mlr_val_t b_b_not_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = !pval1->u.bool_val};
	return rv;
}

static inline mlr_val_t b_bb_or_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = pval1->u.bool_val || pval2->u.bool_val};
	return rv;
}
static inline mlr_val_t b_bb_and_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_BOOL, .u.bool_val = pval1->u.bool_val && pval2->u.bool_val};
	return rv;
}

// ----------------------------------------------------------------
static inline mlr_val_t f_z_urand_func() {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = get_mtrand_double()}; // mtrand.h
	return rv;
}
static inline mlr_val_t f_z_systime_func() {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = get_systime()}; // mlrutil.h
	return rv;
}

// ----------------------------------------------------------------
static inline mlr_val_t f_f_uneg_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = -pval1->u.double_val};
	return rv;
}
static inline mlr_val_t f_f_abs_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = fabs(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_log_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = log(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_log10_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = log10(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_exp_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = exp(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_sin_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = sin(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_cos_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = cos(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_tan_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = tan(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_sqrt_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = sqrt(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_round_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = round(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_floor_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = floor(pval1->u.double_val)};
	return rv;
}
static inline mlr_val_t f_f_ceil_func(mlr_val_t* pval1) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = ceil(pval1->u.double_val)};
	return rv;
}

// ----------------------------------------------------------------
static inline mlr_val_t f_ff_plus_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = pval1->u.double_val + pval2->u.double_val};
	return rv;
}
static inline mlr_val_t f_ff_minus_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = pval1->u.double_val - pval2->u.double_val};
	return rv;
}
static inline mlr_val_t f_ff_times_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = pval1->u.double_val * pval2->u.double_val};
	return rv;
}
static inline mlr_val_t f_ff_divide_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = pval1->u.double_val / pval2->u.double_val};
	return rv;
}
static inline mlr_val_t f_ff_pow_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = pow(pval1->u.double_val, pval2->u.double_val)};
	return rv;
}
static inline mlr_val_t f_ff_atan2_func(mlr_val_t* pval1, mlr_val_t* pval2) {
	mlr_val_t rv = {.type = MT_DOUBLE, .u.double_val = atan2(pval1->u.double_val, pval2->u.double_val)};
	return rv;
}

// ----------------------------------------------------------------
mlr_val_t s_s_tolower_func(mlr_val_t* pval1);
mlr_val_t s_s_toupper_func(mlr_val_t* pval1);

mlr_val_t s_ss_dot_func(mlr_val_t* pval1, mlr_val_t* pval2);

// ----------------------------------------------------------------
mlr_val_t eq_op_func(mlr_val_t* pval1, mlr_val_t* pval2);
mlr_val_t ne_op_func(mlr_val_t* pval1, mlr_val_t* pval2);
mlr_val_t gt_op_func(mlr_val_t* pval1, mlr_val_t* pval2);
mlr_val_t ge_op_func(mlr_val_t* pval1, mlr_val_t* pval2);
mlr_val_t lt_op_func(mlr_val_t* pval1, mlr_val_t* pval2);
mlr_val_t le_op_func(mlr_val_t* pval1, mlr_val_t* pval2);

#endif // MLR_VAL_H
