#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/lrec_evaluators.h"

// ================================================================
typedef struct _lrec_evaluator_b_b_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_b_b_state_t;

mv_t lrec_evaluator_b_b_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_b_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_b_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_b_b_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_b_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_b_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_bb_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_bb_state_t;

mv_t lrec_evaluator_b_bb_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_bb_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);

	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_bb_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_b_bb_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_bb_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_bb_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_z_state_t {
	mv_zary_func_t* pfunc;
} lrec_evaluator_f_z_state_t;

mv_t lrec_evaluator_f_z_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_z_state_t* pstate = pvstate;

	return pstate->pfunc();
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_z_func(mv_zary_func_t* pfunc) {
	lrec_evaluator_f_z_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_z_state_t));
	pstate->pfunc = pfunc;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_z_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_f_state_t {
	mv_unary_func_t* pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_f_state_t;

mv_t lrec_evaluator_f_f_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	mt_get_double_strict(&val1);
	if (val1.type != MT_DOUBLE)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_f_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_ff_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_f_ff_state_t;

mv_t lrec_evaluator_f_ff_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mt_get_double_strict(&val1);
	if (val1.type != MT_DOUBLE)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
NULL_OR_ERROR_OUT(val2);
	mt_get_double_strict(&val2);
	if (val2.type != MT_DOUBLE)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_ff_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_f_ff_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_ff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_ff_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_s_state_t;

mv_t lrec_evaluator_s_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING) // xxx conversions?
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_f_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_f_state_t;

mv_t lrec_evaluator_s_f_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	// xxx decide & document whether to do the typing here or in the pfunc
	//if (val1.type != MT_STRING) // xxx conversions?
		//return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_f_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_s_state_t;

mv_t lrec_evaluator_f_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	// xxx decide & document whether to do the typing here or in the pfunc
	if (val1.type != MT_STRING) // xxx conversions?
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_i_s_state_t;

mv_t lrec_evaluator_i_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_i_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	// xxx decide & document whether to do the typing here or in the pfunc
	if (val1.type != MT_STRING) // xxx conversions?
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_i_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_i_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_xx_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_xx_state_t;

mv_t lrec_evaluator_b_xx_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_xx_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_xx_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_b_xx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_xx_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_xx_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_ss_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_s_ss_state_t;

mv_t lrec_evaluator_s_ss_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_ss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING) // xxx conversions?
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_ss_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_s_ss_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_ss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_ss_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_sss_state_t {
	mv_ternary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_s_sss_state_t;

mv_t lrec_evaluator_s_sss_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_sss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING) // xxx conversions?
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	mv_t val3 = pstate->parg3->pevaluator_func(prec, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT(val3);
	if (val3.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2, &val3);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_sss_func(mv_ternary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	lrec_evaluator_s_sss_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_sss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_sss_func;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_field_name_state_t {
	char* field_name;
} lrec_evaluator_field_name_state_t;

mv_t lrec_evaluator_field_name_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_field_name_state_t* pstate = pvstate;
	char* string = lrec_get(prec, pstate->field_name);
	if (string == NULL) {
		return (mv_t) {.type = MT_NULL, .u.intv = 0};
	} else {
		double dblv;
		if (mlr_try_double_from_string(string, &dblv)) {
			return (mv_t) {.type = MT_DOUBLE, .u.dblv = dblv};
		} else {
			return (mv_t) {.type = MT_STRING, .u.strv = strdup(string)};
		}
	}
}

lrec_evaluator_t* lrec_evaluator_alloc_from_field_name(char* field_name) {
	lrec_evaluator_field_name_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_field_name_state_t));
	pstate->field_name = strdup(field_name);

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_field_name_func;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_literal_state_t {
	mv_t literal;
} lrec_evaluator_literal_state_t;

mv_t lrec_evaluator_double_literal_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_literal_state_t* pstate = pvstate;
	return pstate->literal;
}
mv_t lrec_evaluator_string_literal_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_literal_state_t* pstate = pvstate;
	// xxx cmt strdup semantics :(
	return (mv_t) {.type = MT_STRING, .u.strv = strdup(pstate->literal.u.strv)};
}

lrec_evaluator_t* lrec_evaluator_alloc_from_literal(char* string) {
	lrec_evaluator_literal_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_literal_state_t));
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));

	double dblv;
	if (mlr_try_double_from_string(string, &dblv)) {
		pstate->literal = (mv_t) {.type = MT_DOUBLE, .u.dblv = dblv};
		pevaluator->pevaluator_func = lrec_evaluator_double_literal_func;
	} else {
		pstate->literal = (mv_t) {.type = MT_STRING, .u.strv = strdup(string)};
		pevaluator->pevaluator_func = lrec_evaluator_string_literal_func;
	}
	pevaluator->pvstate = pstate;

	return pevaluator;
}

// ================================================================
mv_t lrec_evaluator_NF_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = prec->field_count};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NF() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_NF_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_NR_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->nr};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_NR_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FNR_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->fnr};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FNR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FNR_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENAME_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_STRING, .u.strv = strdup(pctx->filename)};
}

lrec_evaluator_t* lrec_evaluator_alloc_from_FILENAME() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FILENAME_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENUM_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->filenum};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FILENUM() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FILENUM_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_PI_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_DOUBLE, .u.dblv = M_PI};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_PI() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_PI_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_E_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_DOUBLE, .u.dblv = M_E};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_E() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_E_func;
	return pevaluator;
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_context_variable(char* variable_name) {
	if        (streq(variable_name, "NF"))       { return lrec_evaluator_alloc_from_NF();
    } else if (streq(variable_name, "NR"))       { return lrec_evaluator_alloc_from_NR();
    } else if (streq(variable_name, "FNR"))      { return lrec_evaluator_alloc_from_FNR();
    } else if (streq(variable_name, "FILENAME")) { return lrec_evaluator_alloc_from_FILENAME();
    } else if (streq(variable_name, "FILENUM"))  { return lrec_evaluator_alloc_from_FILENUM();
    } else if (streq(variable_name, "PI"))       { return lrec_evaluator_alloc_from_PI();
    } else if (streq(variable_name, "E"))        { return lrec_evaluator_alloc_from_E();

	} else  { return NULL; // xxx handle me better
	}
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_zary_func_name(char* function_name) {
	if        (streq(function_name, "urand")) {
		return lrec_evaluator_alloc_from_f_z_func(f_z_urand_func);
	} else if (streq(function_name, "systime")) {
		return lrec_evaluator_alloc_from_f_z_func(f_z_systime_func);
	} else  {
		return NULL; // xxx handle me better
	}
}

// ================================================================
typedef struct _function_lookup_t {
	char* function_name;
	int   arity;
	char* usage_string;
} function_lookup_t;

static function_lookup_t FUNCTION_LOOKUP_TABLE[] = {
	{  "abs",     1 , "Absolute value."},
	{  "acos",    1 , "Inverse trigonometric cosine."},
	{  "acosh",   1 , "Inverse hyperbolic cosine."},
	{  "asin",    1 , "Inverse trigonometric sine."},
	{  "asinh",   1 , "Inverse hyperbolic sine."},
	{  "atan",    1 , "One-argument arctangent."},
	{  "atan2",   2 , "Two-argument arctangent."},
	{  "atanh",   1 , "Inverse hyperbolic tangent."},
	{  "cbrt",    1 , "Cube root."},
	{  "ceil",    1 , "Ceiling: nearest integer at or above."},
	{  "cos",     1 , "Trigonometric cosine."},
	{  "cosh",    1 , "Hyperbolic cosine."},
	{  "erf",     1 , "Error function."},
	{  "erfc",    1 , "Complementary error function."},
	{  "exp",     1 , "Exponential function e**x."},
	{  "expm1",   1 , "e**x - 1."},
	{  "floor",   1 , "Floor: nearest integer at or below."},
	{  "gmt2sec", 1 , "Parses GMT timestamp as integer seconds since epoch."},
	{  "log",     1 , "Natural (base-e) logarithm."},
	{  "log10",   1 , "Base-10 logarithm."},
	{  "log1p",   1 , "log(1-x)."},
	{  "pow",     2 , "Exponentiation; same as **."},
	{  "qnorm",   1 , "Normal cumulative distribution function."},
	{  "round",   1 , "Nearest integer."},
	{  "sec2gmt", 1 , "Formats seconds since epoch (integer part only) as GMT timestamp."},
	{  "sin",     1 , "Trigonometric sine."},
	{  "sinh",    1 , "Hyperbolic sine."},
	{  "sqrt",    1 , "Square root."},
	{  "strlen",  1 , "String length."},
	{  "sub",     3 , "Example: '$name=sub($name, \"old\", \"new\")'. Regexes not supported."},
	{  "systime", 0 , "Floating-point seconds since the epoch." },
	{  "tan",     1 , "Hyperbolic tangent."},
	{  "tan",     1 , "Trigonometric tangent."},
	{  "tanh",    1 , "Inverse trigonometric tangent."},
	{  "tolower", 1 , "Convert string to lowercase."},
	{  "toupper", 1 , "Convert string to uppercase."},
	{  "urand",   0 , "Floating-point numbers on the unit interval. Int-valued example: '$n=floor(20+urand()*11)'." },

	{  "==",      2 , "String/numeric equality. Mixing number and string results in string compare."},
	{  "!=",      2 , "String/numeric inequality. Mixing number and string results in string compare."},
	{  ">",       2 , "String/numeric greater-than. Mixing number and string results in string compare."},
	{  ">=",      2 , "String/numeric greater-than-or-equals. Mixing number and string results in string compare."},
	{  "<",       2 , "String/numeric less-than. Mixing number and string results in string compare."},
	{  "<=",      2 , "String/numeric less-than-or-equals. Mixing number and string results in string compare."},
	{  "&&",      2 , "Logical AND."},
	{  "||",      2 , "Logical OR."},
	{  "!",       1 , "Logical negation."},
	{  "+",       2 , "Addition."},
	{  "-",       1 , "Unary minus."},
	{  "-",       2 , "Subtraction."},
	{  "*",       2 , "Multiplication."},
	{  "/",       2 , "Division."},
	{  "%",       2 , "Remainder; never negative-valued."},
	{  "**",      2 , "Exponentiation; same as pow."},
	{  ".",       2 , "String concatenation."},
	{  NULL,      -1 , NULL}, // table terminator
};

#define ARITY_CHECK_PASS    0xbb
#define ARITY_CHECK_FAIL    0xbc
#define ARITY_CHECK_NO_SUCH 0xbd

static int check_arity(function_lookup_t lookup_table[], char* function_name, int user_provided_arity, int *parity) {
	*parity = -1;
	int found_function_name = FALSE;
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &lookup_table[i];
		if (plookup->function_name == NULL)
			break;
		if (streq(function_name, plookup->function_name)) {
			found_function_name = TRUE;
			*parity = plookup->arity;
			if (user_provided_arity == plookup->arity) {
				return ARITY_CHECK_PASS;
			}
		}
	}
	if (found_function_name) {
		return ARITY_CHECK_FAIL;
	} else {
		return ARITY_CHECK_NO_SUCH;
	}
}

static void check_arity_with_report(function_lookup_t function_lookup_table[], char* function_name,
	int user_provided_arity)
{
	int arity = -1;
	int result = check_arity(function_lookup_table, function_name, user_provided_arity, &arity);
	if (result == ARITY_CHECK_NO_SUCH) {
		fprintf(stderr, "Function name \"%s\" not found.\n", function_name);
		exit(1);
	}
	if (result == ARITY_CHECK_FAIL) {
		// More flexibly, I'd have a list of arities supported by each
		// function. But this is overkill: there are unary and binary minus,
		// and everything else has a single arity.
		if (streq(function_name, "-")) {
			fprintf(stderr, "Function named \"%s\" takes one argument or two; got %d.\n",
				function_name, user_provided_arity);
		} else {
		}
			fprintf(stderr, "Function named \"%s\" takes %d argument%s; got %d.\n",
				function_name, arity, (arity == 1) ? "" : "s", user_provided_arity);
		exit(1);
	}
}

void lrec_evaluator_list_functions(FILE* output_stream) {
	fprintf(output_stream, "Functions for filter and put:\n");

	int linelen = 0;
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL)
			break;
		linelen += 1 + strlen(FUNCTION_LOOKUP_TABLE[i].function_name);
		if (linelen > 80) {
			fprintf(output_stream, "\n");
			linelen = 0;
		}
		if ((i > 0) && (linelen > 0))
			fprintf(output_stream, " ");
		else
			fprintf(output_stream, "   ");
		fprintf(output_stream, "%s", FUNCTION_LOOKUP_TABLE[i].function_name);
	}
	fprintf(output_stream, "\n");
}

// Pass function_name == NULL to get usage for all functions.
void lrec_evaluator_function_usage(FILE* output_stream, char* function_name) {
	int found = FALSE;
	char* fmt = (function_name == NULL)
		? "%-10s (#args=%d): %s\n"
		: "%s (#args=%d): %s\n";

	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL)
			break;
		if (function_name == NULL || streq(function_name, plookup->function_name)) {
			fprintf(output_stream, fmt, plookup->function_name, plookup->arity, plookup->usage_string);
			found = TRUE;
		}
	}
	if (!found)
		fprintf(output_stream, "%s: no such function.\n", function_name);
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_unary_func_name(char* fnnm, lrec_evaluator_t* parg1) {
	if        (streq(fnnm, "!"))       { return lrec_evaluator_alloc_from_b_b_func(b_b_not_func,     parg1);
	} else if (streq(fnnm, "-"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_uneg_func,    parg1);
	} else if (streq(fnnm, "abs"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_abs_func,     parg1);
	} else if (streq(fnnm, "acos"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_acos_func,    parg1);
	} else if (streq(fnnm, "acosh"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_acosh_func,   parg1);
	} else if (streq(fnnm, "asin"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_asin_func,    parg1);
	} else if (streq(fnnm, "asinh"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_asinh_func,   parg1);
	} else if (streq(fnnm, "atan"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_atan_func,    parg1);
	} else if (streq(fnnm, "atanh"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_atanh_func,   parg1);
	} else if (streq(fnnm, "cbrt"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_cbrt_func,    parg1);
	} else if (streq(fnnm, "ceil"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_ceil_func,    parg1);
	} else if (streq(fnnm, "cos"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_cos_func,     parg1);
	} else if (streq(fnnm, "cosh"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_cosh_func,    parg1);
	} else if (streq(fnnm, "erf"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_erf_func,     parg1);
	} else if (streq(fnnm, "erfc"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_erfc_func,    parg1);
	} else if (streq(fnnm, "exp"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_exp_func,     parg1);
	} else if (streq(fnnm, "floor"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_floor_func,   parg1);
	} else if (streq(fnnm, "gmt2sec")) { return lrec_evaluator_alloc_from_i_s_func(i_s_gmt2sec_func, parg1);
	} else if (streq(fnnm, "log"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_log_func,     parg1);
	} else if (streq(fnnm, "log10"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_log10_func,   parg1);
	} else if (streq(fnnm, "qnorm"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_qnorm_func,   parg1);
	} else if (streq(fnnm, "round"))   { return lrec_evaluator_alloc_from_f_f_func(f_f_round_func,   parg1);
	} else if (streq(fnnm, "sec2gmt")) { return lrec_evaluator_alloc_from_s_f_func(s_f_sec2gmt_func, parg1);
	} else if (streq(fnnm, "sin"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_sin_func,     parg1);
	} else if (streq(fnnm, "sinh"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_sinh_func,    parg1);
	} else if (streq(fnnm, "sqrt"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_sqrt_func,    parg1);
	} else if (streq(fnnm, "strlen"))  { return lrec_evaluator_alloc_from_i_s_func(i_s_strlen_func,  parg1);
	} else if (streq(fnnm, "tan"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_tan_func,     parg1);
	} else if (streq(fnnm, "tanh"))    { return lrec_evaluator_alloc_from_f_f_func(f_f_tanh_func,    parg1);
	} else if (streq(fnnm, "tolower")) { return lrec_evaluator_alloc_from_s_s_func(s_s_tolower_func, parg1);
	} else if (streq(fnnm, "toupper")) { return lrec_evaluator_alloc_from_s_s_func(s_s_toupper_func, parg1);

	} else return NULL; // xxx handle me better
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_binary_func_name(char* fnnm,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	if        (streq(fnnm, "&&"))    { return lrec_evaluator_alloc_from_b_bb_func(b_bb_and_func,    parg1, parg2);
	} else if (streq(fnnm, "||"))    { return lrec_evaluator_alloc_from_b_bb_func(b_bb_or_func,     parg1, parg2);
	} else if (streq(fnnm, "=="))    { return lrec_evaluator_alloc_from_b_xx_func(eq_op_func,       parg1, parg2);
	} else if (streq(fnnm, "!="))    { return lrec_evaluator_alloc_from_b_xx_func(ne_op_func,       parg1, parg2);
	} else if (streq(fnnm, ">"))     { return lrec_evaluator_alloc_from_b_xx_func(gt_op_func,       parg1, parg2);
	} else if (streq(fnnm, ">="))    { return lrec_evaluator_alloc_from_b_xx_func(ge_op_func,       parg1, parg2);
	} else if (streq(fnnm, "<"))     { return lrec_evaluator_alloc_from_b_xx_func(lt_op_func,       parg1, parg2);
	} else if (streq(fnnm, "<="))    { return lrec_evaluator_alloc_from_b_xx_func(le_op_func,       parg1, parg2);
	} else if (streq(fnnm, "."))     { return lrec_evaluator_alloc_from_s_ss_func(s_ss_dot_func,    parg1, parg2);
	} else if (streq(fnnm, "pow"))   { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,    parg1, parg2);
	} else if (streq(fnnm, "+"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_plus_func,   parg1, parg2);
	} else if (streq(fnnm, "-"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_minus_func,  parg1, parg2);
	} else if (streq(fnnm, "*"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func,  parg1, parg2);
	} else if (streq(fnnm, "/"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_divide_func, parg1, parg2);
	} else if (streq(fnnm, "**"))    { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,    parg1, parg2);
	} else if (streq(fnnm, "%"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_mod_func,    parg1, parg2);
	} else if (streq(fnnm, "atan2")) { return lrec_evaluator_alloc_from_f_ff_func(f_ff_atan2_func,  parg1, parg2);
	} else  { return NULL; /* xxx handle me better */ }
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_ternary_func_name(char* fnnm,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	if (streq(fnnm, "sub")) { return lrec_evaluator_alloc_from_s_sss_func(s_sss_sub_func,   parg1, parg2, parg3);
	} else  { return NULL; /* xxx handle me better */ }
}

// ================================================================
static lrec_evaluator_t* lrec_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	function_lookup_t* function_lookup_table)
{
	if (pnode->pchildren == NULL) { // leaf node
		if (pnode->type == MLR_DSL_AST_NODE_TYPE_FIELD_NAME) {
			return lrec_evaluator_alloc_from_field_name(pnode->text);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_LITERAL) {
			return lrec_evaluator_alloc_from_literal(pnode->text);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE) {
			return lrec_evaluator_alloc_from_context_variable(pnode->text);
		} else {
			fprintf(stderr, "xxx write this error message please.\n");
			return NULL;
		}
	} else { // operator/function
		if ((pnode->type != MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME)
		&& (pnode->type != MLR_DSL_AST_NODE_TYPE_OPERATOR)) {
			fprintf(stderr, "yyy write this error message please: %04x.\n", pnode->type);
			return NULL;
		}
		char* func_name = pnode->text;

		int user_provided_arity = pnode->pchildren->length;

		check_arity_with_report(function_lookup_table, func_name, user_provided_arity);

		lrec_evaluator_t* pevaluator = NULL;
		if (user_provided_arity == 0) {
			pevaluator = lrec_evaluator_alloc_from_zary_func_name(func_name);
		} else if (user_provided_arity == 1) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_unary_func_name(func_name, parg1);
		} else if (user_provided_arity == 2) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_binary_func_name(func_name, parg1, parg2);
		} else if (user_provided_arity == 3) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvdata;
			mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, function_lookup_table);
			lrec_evaluator_t* parg3 = lrec_evaluator_alloc_from_ast_aux(parg3_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_ternary_func_name(func_name, parg1, parg2, parg3);
		} else {
			fprintf(stderr, "Internal coding error:  arity for function name \"%s\" misdetected.\n",
				func_name);
			exit(1);
		}
		if (pevaluator == NULL) {
			fprintf(stderr, "Unrecognized function name \"%s\".\n", func_name);
			exit(1);
		}
		return pevaluator;
	}
}

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode) {
	lrec_evaluator_t* pevaluator = lrec_evaluator_alloc_from_ast_aux(pnode, FUNCTION_LOOKUP_TABLE);
	return pevaluator;
}

// ================================================================
#ifdef __LREC_EVALUATORS_MAIN__
int main(int argc, char** argv) {
	mtrand_init_default();

	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	// ----------------------------------------------------------------
	lrec_evaluator_t* pnr       = lrec_evaluator_alloc_from_NR();
	lrec_evaluator_t* pfnr      = lrec_evaluator_alloc_from_FNR();
	lrec_evaluator_t* pfilename = lrec_evaluator_alloc_from_FILENAME();
	lrec_evaluator_t* pfilenum  = lrec_evaluator_alloc_from_FILENUM();

	lrec_t* prec = lrec_alloc();

	mv_t val = pnr->pevaluator_func(prec, pctx, pnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	val = pfnr->pevaluator_func(prec, pctx, pfnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	val = pfilename->pevaluator_func(prec, pctx, pfilename->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	val = pfilenum->pevaluator_func(prec, pctx, pfilenum->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	// ----------------------------------------------------------------
	// $s + "def"

	lrec_evaluator_t* ps       = lrec_evaluator_alloc_from_field_name("s");
	lrec_evaluator_t* pdef     = lrec_evaluator_alloc_from_literal("def");
	lrec_evaluator_t* pdot     = lrec_evaluator_alloc_from_s_ss_func(s_ss_dot_func, ps, pdef);
	lrec_evaluator_t* ptolower = lrec_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	lrec_evaluator_t* ptoupper = lrec_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	prec = lrec_alloc();
	lrec_put(prec, "s", "abc");
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	val = ps->pevaluator_func(prec, pctx, ps->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	val = pdef->pevaluator_func(prec, pctx, pdef->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	val = pdot->pevaluator_func(prec, pctx, pdot->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	val = ptolower->pevaluator_func(prec, pctx, ptolower->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	val = ptoupper->pevaluator_func(prec, pctx, ptoupper->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));

	// ----------------------------------------------------------------
	// 2.0 * log($x) + rand()

	lrec_evaluator_t* p2     = lrec_evaluator_alloc_from_literal("2.0");
	lrec_evaluator_t* px     = lrec_evaluator_alloc_from_field_name("x");
	lrec_evaluator_t* plogx  = lrec_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	lrec_evaluator_t* p2logx = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, p2, plogx);
	lrec_evaluator_t* prand  = lrec_evaluator_alloc_from_f_z_func(f_z_urand_func);
	lrec_evaluator_t* psum   = lrec_evaluator_alloc_from_f_ff_func(f_ff_plus_func, p2logx, prand);
	lrec_evaluator_t* px2    = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, px, px);
	lrec_evaluator_t* p4     = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MLR_DSL_AST_NODE_TYPE_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MLR_DSL_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	lrec_evaluator_t*  pastr = lrec_evaluator_alloc_from_ast(p2logxnode);

	prec = lrec_alloc();
	lrec_put(prec, "x", "4.5");

    printf("lrec   x        = %s\n", lrec_get(prec, "x"));
    printf("newval 2        = %s\n", mt_describe_val(p2->pevaluator_func(prec,     pctx,  p2->pvstate)));
    printf("newval 4        = %s\n", mt_describe_val(p4->pevaluator_func(prec,     pctx,  p4->pvstate)));
    printf("newval x        = %s\n", mt_describe_val(px->pevaluator_func(prec,     pctx,  px->pvstate)));
    printf("newval x^2      = %s\n", mt_describe_val(px2->pevaluator_func(prec,    pctx,  px2->pvstate)));
    printf("newval log(x)   = %s\n", mt_describe_val(plogx->pevaluator_func(prec,  pctx,  plogx->pvstate)));
    printf("newval 2*log(x) = %s\n", mt_describe_val(p2logx->pevaluator_func(prec, pctx,  p2logx->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));

	printf("newval sum      = %s\n",  mt_describe_val(psum->pevaluator_func(prec, pctx, psum->pvstate)));

	mlr_dsl_ast_node_print(p2logxnode);
	printf("newval AST      = %s\n",  mt_describe_val(pastr->pevaluator_func(prec, pctx, pastr->pvstate)));
	printf("\n");

	lrec_rename(prec, "x", "y");

    printf("lrec   x        = %s\n", lrec_get(prec, "x"));
    printf("newval 2        = %s\n", mt_describe_val(p2->pevaluator_func(prec,     pctx,  p2->pvstate)));
    printf("newval 4        = %s\n", mt_describe_val(p4->pevaluator_func(prec,     pctx,  p4->pvstate)));
    printf("newval x        = %s\n", mt_describe_val(px->pevaluator_func(prec,     pctx,  px->pvstate)));
    printf("newval x^2      = %s\n", mt_describe_val(px2->pevaluator_func(prec,    pctx,  px2->pvstate)));
    printf("newval log(x)   = %s\n", mt_describe_val(plogx->pevaluator_func(prec,  pctx,  plogx->pvstate)));
    printf("newval 2*log(x) = %s\n", mt_describe_val(p2logx->pevaluator_func(prec, pctx,  p2logx->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));
    printf("newval urand    = %s\n", mt_describe_val(prand->pevaluator_func(prec,  pctx,  prand->pvstate)));
    printf("newval sum      = %s\n", mt_describe_val(psum->pevaluator_func(prec,   pctx,  psum->pvstate)));

	mlr_dsl_ast_node_print(p2logxnode);
	printf("newval AST      = %s\n",  mt_describe_val(pastr->pevaluator_func(prec, pctx, pastr->pvstate)));
	printf("\n");

	return 0;
}
#endif // __LREC_EVALUATORS_MAIN__
