#include "lib/mlr_globals.h"
#include "lib/minunit.h"
#include "mapping/rval_evaluators.h"

// ----------------------------------------------------------------
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test1() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test1 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* pnr       = rval_evaluator_alloc_from_NR();
	rval_evaluator_t* pfnr      = rval_evaluator_alloc_from_FNR();
	rval_evaluator_t* pfilename = rval_evaluator_alloc_from_FILENAME();
	rval_evaluator_t* pfilenum  = rval_evaluator_alloc_from_FILENUM();

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;

	mv_t val = pnr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 888);

	val = pfnr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 999);

	val = pfilename->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfilename->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "filename-goes-here"));

	val = pfilenum->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfilenum->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 123);

	return 0;
}

// ----------------------------------------------------------------
static char * test2() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test2 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* ps       = rval_evaluator_alloc_from_field_name("s", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* pdef     = rval_evaluator_alloc_from_strnum_literal("def", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* pdot     = rval_evaluator_alloc_from_x_ss_func(s_ss_dot_func, ps, pdef);
	rval_evaluator_t* ptolower = rval_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	rval_evaluator_t* ptoupper = rval_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "s", "abc", NO_FREE);
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	mv_t val = ps->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ps->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abc"));

	val = pdef->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pdef->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "def"));

	val = pdot->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pdot->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptolower->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ptolower->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptoupper->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ptoupper->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "ABCDEF"));

	return 0;
}

// ----------------------------------------------------------------
static char * test3() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test3 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* p2     = rval_evaluator_alloc_from_strnum_literal("2.0", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* px     = rval_evaluator_alloc_from_field_name("x", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* plogx  = rval_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	rval_evaluator_t* p2logx = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, plogx);
	rval_evaluator_t* px2    = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, px, px);
	rval_evaluator_t* p4     = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MD_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MD_AST_NODE_TYPE_FUNCTION_NAME);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MD_AST_NODE_TYPE_STRNUM_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	rval_evaluator_t*  pastr = rval_evaluator_alloc_from_ast(p2logxnode, TYPE_INFER_STRING_FLOAT_INT);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "x", "4.5", NO_FREE);

	mv_t valp2     = p2->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p2->pvstate);
	mv_t valp4     = p4->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p4->pvstate);
	mv_t valpx     = px->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, px->pvstate);
	mv_t valpx2    = px2->pprocess_func(prec,    ptyped_overlay, poosvars, &pregex_captures, pctx, px2->pvstate);
	mv_t valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, poosvars, &pregex_captures, pctx, plogx->pvstate);
	mv_t valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, p2logx->pvstate);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

	printf("XXX %s\n", mt_describe_type(valp2.type));
	mu_assert_lf(valp2.type     == MT_FLOAT);
	mu_assert_lf(valp4.type     == MT_FLOAT);
	mu_assert_lf(valpx.type     == MT_FLOAT);
	mu_assert_lf(valpx2.type    == MT_FLOAT);
	mu_assert_lf(valplogx.type  == MT_FLOAT);
	mu_assert_lf(valp2logx.type == MT_FLOAT);

	mu_assert_lf(valp2.u.fltv     == 2.0);
	mu_assert_lf(valp4.u.fltv     == 4.0);
	mu_assert_lf(valpx.u.fltv     == 4.5);
	mu_assert_lf(valpx2.u.fltv    == 20.25);
	mu_assert_lf(fabs(valplogx.u.fltv  - 0.653213) < 1e-5);
	mu_assert_lf(fabs(valp2logx.u.fltv - 1.306425) < 1e-5);

	mlr_dsl_ast_node_print(p2logxnode);
	printf("newval AST      = %s\n",
		mv_describe_val(pastr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pastr->pvstate)));
	printf("\n");

	lrec_rename(prec, "x", "y", FALSE);

	valp2     = p2->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p2->pvstate);
	valp4     = p4->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p4->pvstate);
	valpx     = px->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, px->pvstate);
	valpx2    = px2->pprocess_func(prec,    ptyped_overlay, poosvars, &pregex_captures, pctx, px2->pvstate);
	valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, poosvars, &pregex_captures, pctx, plogx->pvstate);
	valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, p2logx->pvstate);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

	mu_assert_lf(valp2.type     == MT_FLOAT);
	mu_assert_lf(valp4.type     == MT_FLOAT);
	mu_assert_lf(valpx.type     == MT_ABSENT);
	mu_assert_lf(valpx2.type    == MT_VOID);
	mu_assert_lf(valplogx.type  == MT_ABSENT);
	mu_assert_lf(valp2logx.type == MT_FLOAT);

	mu_assert_lf(valp2.u.fltv     == 2.0);
	mu_assert_lf(valp4.u.fltv     == 4.0);

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test1);
	mu_run_test(test2);
	mu_run_test(test3);
	return 0;
}

// ----------------------------------------------------------------
int main(int argc, char **argv) {
	mlr_global_init(argv[0], "%lf", NULL);

	printf("TEST_RVAL_EVALUATORS ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_RVAL_EVALUATORS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
