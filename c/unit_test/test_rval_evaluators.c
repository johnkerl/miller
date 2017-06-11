#include "lib/mlr_globals.h"
#include "lib/minunit.h"
#include "dsl/rval_evaluators.h"
#include "dsl/function_manager.h"

// ----------------------------------------------------------------
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test_caps() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test_caps ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	rval_evaluator_t* pnr       = rval_evaluator_alloc_from_NR();
	rval_evaluator_t* pfnr      = rval_evaluator_alloc_from_FNR();
	rval_evaluator_t* pfilename = rval_evaluator_alloc_from_FILENAME();
	rval_evaluator_t* pfilenum  = rval_evaluator_alloc_from_FILENUM();

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsmv_t* ptyped_overlay = lhmsmv_alloc();
	mlhmmv_root_t* poosvars = mlhmmv_root_alloc();
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t val = pnr->pprocess_func(pnr->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 888);

	val = pfnr->pprocess_func(pfnr->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 999);

	val = pfilename->pprocess_func(pfilename->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "filename-goes-here"));

	val = pfilenum->pprocess_func(pfilenum->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 123);

	return 0;
}

// ----------------------------------------------------------------
static char * test_strings() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test_strings ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	rval_evaluator_t* ps       = rval_evaluator_alloc_from_field_name("s", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* pdef     = rval_evaluator_alloc_from_string_literal("def");
	rval_evaluator_t* pdot     = rval_evaluator_alloc_from_x_ss_func(s_xx_dot_func, ps, pdef);
	rval_evaluator_t* ptolower = rval_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	rval_evaluator_t* ptoupper = rval_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsmv_t* ptyped_overlay = lhmsmv_alloc();
	mlhmmv_root_t* poosvars = mlhmmv_root_alloc();
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	lrec_put(prec, "s", "abc", NO_FREE);
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t val = ps->pprocess_func(ps->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abc"));

	val = pdef->pprocess_func(pdef->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "def"));

	val = pdot->pprocess_func(pdot->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptolower->pprocess_func(ptolower->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptoupper->pprocess_func(ptoupper->pvstate, &variables);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "ABCDEF"));

	return 0;
}

// ----------------------------------------------------------------
static char * test_numbers() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test_numbers ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	rval_evaluator_t* p2     = rval_evaluator_alloc_from_numeric_literal("2.0");
	rval_evaluator_t* px     = rval_evaluator_alloc_from_field_name("x", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* plogx  = rval_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	rval_evaluator_t* p2logx = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, plogx);
	rval_evaluator_t* px2    = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, px, px);
	rval_evaluator_t* p4     = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MD_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MD_AST_NODE_TYPE_FUNCTION_CALLSITE);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	fmgr_t* pfmgr = fmgr_alloc();
	rval_evaluator_t*  pastr = rval_evaluator_alloc_from_ast(p2logxnode, pfmgr, TYPE_INFER_STRING_FLOAT_INT, 0);
	fmgr_resolve_func_callsites(pfmgr);
	fmgr_free(pfmgr, &ctx);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsmv_t* ptyped_overlay = lhmsmv_alloc();
	mlhmmv_root_t* poosvars = mlhmmv_root_alloc();
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	lrec_put(prec, "x", "4.5", NO_FREE);

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t valp2     = p2->pprocess_func(p2->pvstate, &variables);
	mv_t valp4     = p4->pprocess_func(p4->pvstate, &variables);
	mv_t valpx     = px->pprocess_func(px->pvstate, &variables);
	mv_t valpx2    = px2->pprocess_func(px2->pvstate, &variables);
	mv_t valplogx  = plogx->pprocess_func(plogx->pvstate, &variables);
	mv_t valp2logx = p2logx->pprocess_func(p2logx->pvstate, &variables);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

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
		mv_describe_val(pastr->pprocess_func(pastr->pvstate, &variables)));
	printf("\n");

	lrec_rename(prec, "x", "y", FALSE);

	valp2     = p2->pprocess_func(p2->pvstate, &variables);
	valp4     = p4->pprocess_func(p4->pvstate, &variables);
	valpx     = px->pprocess_func(px->pvstate, &variables);
	valpx2    = px2->pprocess_func(px2->pvstate, &variables);
	valplogx  = plogx->pprocess_func(plogx->pvstate, &variables);
	valp2logx = p2logx->pprocess_func(p2logx->pvstate, &variables);

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
	mu_assert_lf(valpx2.type    == MT_ABSENT);
	mu_assert_lf(valplogx.type  == MT_ABSENT);
	mu_assert_lf(valp2logx.type == MT_FLOAT);

	mu_assert_lf(valp2.u.fltv     == 2.0);
	mu_assert_lf(valp4.u.fltv     == 4.0);

	return 0;
}

// ----------------------------------------------------------------
static char * test_logical_and() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test4 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	lrec_t* prec = NULL;
	lhmsmv_t* ptyped_overlay = NULL;
	mlhmmv_root_t* poosvars = NULL;
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t t = mv_from_bool(TRUE);
	mv_t f = mv_from_bool(FALSE);
	mv_t a = mv_absent();

	rval_evaluator_t* pt = rval_evaluator_alloc_from_mlrval(&t);
	rval_evaluator_t* pf = rval_evaluator_alloc_from_mlrval(&f);
	rval_evaluator_t* pa = rval_evaluator_alloc_from_mlrval(&a);

	rval_evaluator_t* ptt = rval_evaluator_alloc_from_b_bb_and_func(pt, pt);
	rval_evaluator_t* ptf = rval_evaluator_alloc_from_b_bb_and_func(pt, pf);
	rval_evaluator_t* pta = rval_evaluator_alloc_from_b_bb_and_func(pt, pa);
	rval_evaluator_t* pft = rval_evaluator_alloc_from_b_bb_and_func(pf, pt);
	rval_evaluator_t* pff = rval_evaluator_alloc_from_b_bb_and_func(pf, pf);
	rval_evaluator_t* pfa = rval_evaluator_alloc_from_b_bb_and_func(pf, pa);
	rval_evaluator_t* pat = rval_evaluator_alloc_from_b_bb_and_func(pa, pt);
	rval_evaluator_t* paf = rval_evaluator_alloc_from_b_bb_and_func(pa, pf);
	rval_evaluator_t* paa = rval_evaluator_alloc_from_b_bb_and_func(pa, pa);

	mv_t ott = ptt->pprocess_func(ptt->pvstate, &variables);
	mv_t otf = ptf->pprocess_func(ptf->pvstate, &variables);
	mv_t oft = pft->pprocess_func(pft->pvstate, &variables);
	mv_t off = pff->pprocess_func(pff->pvstate, &variables);
	mv_t oat = pat->pprocess_func(pat->pvstate, &variables);
	mv_t oaf = paf->pprocess_func(paf->pvstate, &variables);
	mv_t ota = pta->pprocess_func(pta->pvstate, &variables);
	mv_t ofa = pfa->pprocess_func(pfa->pvstate, &variables);
	mv_t oaa = paa->pprocess_func(paa->pvstate, &variables);

	mu_assert_lf(ott.type == MT_BOOLEAN); mu_assert_lf(ott.u.boolv == TRUE);
	mu_assert_lf(otf.type == MT_BOOLEAN); mu_assert_lf(otf.u.boolv == FALSE);
	mu_assert_lf(oft.type == MT_BOOLEAN); mu_assert_lf(oft.u.boolv == FALSE);
	mu_assert_lf(off.type == MT_BOOLEAN); mu_assert_lf(off.u.boolv == FALSE);
	mu_assert_lf(oat.type == MT_BOOLEAN); mu_assert_lf(oat.u.boolv == TRUE);
	mu_assert_lf(oaf.type == MT_BOOLEAN); mu_assert_lf(oaf.u.boolv == FALSE);
	mu_assert_lf(ota.type == MT_BOOLEAN); mu_assert_lf(ota.u.boolv == TRUE);
	mu_assert_lf(ofa.type == MT_BOOLEAN); mu_assert_lf(ofa.u.boolv == FALSE);
	mu_assert_lf(oaa.type == MT_ABSENT);

	return 0;
}

// ----------------------------------------------------------------
static char * test_logical_or() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test4 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	lrec_t* prec = NULL;
	lhmsmv_t* ptyped_overlay = NULL;
	mlhmmv_root_t* poosvars = NULL;
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t t = mv_from_bool(TRUE);
	mv_t f = mv_from_bool(FALSE);
	mv_t a = mv_absent();

	rval_evaluator_t* pt = rval_evaluator_alloc_from_mlrval(&t);
	rval_evaluator_t* pf = rval_evaluator_alloc_from_mlrval(&f);
	rval_evaluator_t* pa = rval_evaluator_alloc_from_mlrval(&a);

	rval_evaluator_t* ptt = rval_evaluator_alloc_from_b_bb_or_func(pt, pt);
	rval_evaluator_t* ptf = rval_evaluator_alloc_from_b_bb_or_func(pt, pf);
	rval_evaluator_t* pta = rval_evaluator_alloc_from_b_bb_or_func(pt, pa);
	rval_evaluator_t* pft = rval_evaluator_alloc_from_b_bb_or_func(pf, pt);
	rval_evaluator_t* pff = rval_evaluator_alloc_from_b_bb_or_func(pf, pf);
	rval_evaluator_t* pfa = rval_evaluator_alloc_from_b_bb_or_func(pf, pa);
	rval_evaluator_t* pat = rval_evaluator_alloc_from_b_bb_or_func(pa, pt);
	rval_evaluator_t* paf = rval_evaluator_alloc_from_b_bb_or_func(pa, pf);
	rval_evaluator_t* paa = rval_evaluator_alloc_from_b_bb_or_func(pa, pa);

	mv_t ott = ptt->pprocess_func(ptt->pvstate, &variables);
	mv_t otf = ptf->pprocess_func(ptf->pvstate, &variables);
	mv_t ota = pta->pprocess_func(pta->pvstate, &variables);
	mv_t oft = pft->pprocess_func(pft->pvstate, &variables);
	mv_t off = pff->pprocess_func(pff->pvstate, &variables);
	mv_t ofa = pfa->pprocess_func(pfa->pvstate, &variables);
	mv_t oat = pat->pprocess_func(pat->pvstate, &variables);
	mv_t oaf = paf->pprocess_func(paf->pvstate, &variables);
	mv_t oaa = paa->pprocess_func(paa->pvstate, &variables);

	mu_assert_lf(ott.type == MT_BOOLEAN); mu_assert_lf(ott.u.boolv == TRUE);
	mu_assert_lf(otf.type == MT_BOOLEAN); mu_assert_lf(otf.u.boolv == TRUE);
	mu_assert_lf(ota.type == MT_BOOLEAN); mu_assert_lf(ota.u.boolv == TRUE);
	mu_assert_lf(oft.type == MT_BOOLEAN); mu_assert_lf(oft.u.boolv == TRUE);
	mu_assert_lf(off.type == MT_BOOLEAN); mu_assert_lf(off.u.boolv == FALSE);
	mu_assert_lf(ofa.type == MT_BOOLEAN); mu_assert_lf(ofa.u.boolv == FALSE);
	mu_assert_lf(oat.type == MT_BOOLEAN); mu_assert_lf(oat.u.boolv == TRUE);
	mu_assert_lf(oaf.type == MT_BOOLEAN); mu_assert_lf(oaf.u.boolv == FALSE);
	mu_assert_lf(oaa.type == MT_ABSENT);

	return 0;
}

// ----------------------------------------------------------------
static char * test_logical_xor() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test4 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here", .force_eof = FALSE,
		.ips = "=", .ifs = ",", .irs = "\n", .ops = "=", .ofs = ",", .ors = "\n", .auto_line_term = "\n"
	};
	context_t* pctx = &ctx;

	lrec_t* prec = NULL;
	lhmsmv_t* ptyped_overlay = NULL;
	mlhmmv_root_t* poosvars = NULL;
	string_array_t* pregex_captures = NULL;
	loop_stack_t* ploop_stack = loop_stack_alloc();

	variables_t variables = (variables_t) {
		.pinrec           = prec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.ploop_stack      = ploop_stack,
	};

	mv_t t = mv_from_bool(TRUE);
	mv_t f = mv_from_bool(FALSE);
	mv_t a = mv_absent();

	rval_evaluator_t* pt = rval_evaluator_alloc_from_mlrval(&t);
	rval_evaluator_t* pf = rval_evaluator_alloc_from_mlrval(&f);
	rval_evaluator_t* pa = rval_evaluator_alloc_from_mlrval(&a);

	rval_evaluator_t* ptt = rval_evaluator_alloc_from_b_bb_xor_func(pt, pt);
	rval_evaluator_t* ptf = rval_evaluator_alloc_from_b_bb_xor_func(pt, pf);
	rval_evaluator_t* pft = rval_evaluator_alloc_from_b_bb_xor_func(pf, pt);
	rval_evaluator_t* pff = rval_evaluator_alloc_from_b_bb_xor_func(pf, pf);
	rval_evaluator_t* pat = rval_evaluator_alloc_from_b_bb_xor_func(pa, pt);
	rval_evaluator_t* paf = rval_evaluator_alloc_from_b_bb_xor_func(pa, pf);
	rval_evaluator_t* pta = rval_evaluator_alloc_from_b_bb_xor_func(pt, pa);
	rval_evaluator_t* pfa = rval_evaluator_alloc_from_b_bb_xor_func(pf, pa);
	rval_evaluator_t* paa = rval_evaluator_alloc_from_b_bb_xor_func(pa, pa);

	mv_t ott = ptt->pprocess_func(ptt->pvstate, &variables);
	mv_t otf = ptf->pprocess_func(ptf->pvstate, &variables);
	mv_t oft = pft->pprocess_func(pft->pvstate, &variables);
	mv_t off = pff->pprocess_func(pff->pvstate, &variables);
	mv_t oat = pat->pprocess_func(pat->pvstate, &variables);
	mv_t oaf = paf->pprocess_func(paf->pvstate, &variables);
	mv_t ota = pta->pprocess_func(pta->pvstate, &variables);
	mv_t ofa = pfa->pprocess_func(pfa->pvstate, &variables);
	mv_t oaa = paa->pprocess_func(paa->pvstate, &variables);

	mu_assert_lf(ott.type == MT_BOOLEAN); mu_assert_lf(ott.u.boolv == FALSE);
	mu_assert_lf(otf.type == MT_BOOLEAN); mu_assert_lf(otf.u.boolv == TRUE);
	mu_assert_lf(oft.type == MT_BOOLEAN); mu_assert_lf(oft.u.boolv == TRUE);
	mu_assert_lf(off.type == MT_BOOLEAN); mu_assert_lf(off.u.boolv == FALSE);
	mu_assert_lf(oat.type == MT_BOOLEAN); mu_assert_lf(oat.u.boolv == TRUE);
	mu_assert_lf(oaf.type == MT_BOOLEAN); mu_assert_lf(oaf.u.boolv == FALSE);
	mu_assert_lf(ota.type == MT_BOOLEAN); mu_assert_lf(ota.u.boolv == TRUE);
	mu_assert_lf(ofa.type == MT_BOOLEAN); mu_assert_lf(ofa.u.boolv == FALSE);
	mu_assert_lf(oaa.type == MT_ABSENT);

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test_caps);
	mu_run_test(test_strings);
	mu_run_test(test_numbers);
	mu_run_test(test_logical_and);
	mu_run_test(test_logical_or);
	mu_run_test(test_logical_xor);
	// There is more operator testing in reg_test/run
	return 0;
}

// ----------------------------------------------------------------
int main(int argc, char **argv) {
	mlr_global_init(argv[0], "%lf");

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
