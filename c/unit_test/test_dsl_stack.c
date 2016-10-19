#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "dsls/mlr_dsl_wrapper.h"
#include "mapping/mlr_dsl_ast.h"
#include "mapping/mlr_dsl_cst.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static char* sep = "================================================================";

// ================================================================
// LOCAL SUPPORTING FUNCTIONS

// ----------------------------------------------------------------
static mlr_dsl_ast_t* my_mlr_dsl_parse(char* string) {
	mlr_dsl_ast_t* past = mlr_dsl_parse(string, FALSE);
	if (past == NULL) {
		fprintf(stderr, "Syntax error on DSL parse of '%s'\n", string);
		exit(1);
	}
	return past;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_t* my_mlr_dsl_cst_alloc(mlr_dsl_ast_t* past) {
	int print_ast = TRUE;
	int trace_stack_allocation = TRUE;

	mlr_dsl_cst_t* pcst = mlr_dsl_cst_alloc(past, print_ast, trace_stack_allocation,
		TYPE_INFER_STRING_FLOAT_INT, FALSE, FALSE, FALSE);

	return pcst;
}

// ----------------------------------------------------------------
static mv_t* imv(long long intv) {
	mv_t* pmv = mlr_malloc_or_die(sizeof(mv_t));
	*pmv = mv_from_int(intv);
	return pmv;
}

// ----------------------------------------------------------------
static variables_t* my_make_variables() {
	variables_t* pvariables = mlr_malloc_or_die(sizeof(variables_t));

	lrec_t* pinrec = lrec_unbacked_alloc();
	lhmsmv_t* ptyped_overlay = lhmsmv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL; // May be set to non-null on evaluation
	context_t* pctx = mlr_malloc_or_die(sizeof(context_t));
	*pctx = (context_t) { .nr = 1, .fnr = 1, .filenum = 0, .filename = NULL, .force_eof = FALSE };
	local_stack_t* plocal_stack = local_stack_alloc();
	loop_stack_t* ploop_stack = loop_stack_alloc();

	*pvariables = (variables_t) {
		.pinrec           = pinrec,
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.plocal_stack     = plocal_stack,
		.ploop_stack      = ploop_stack,
		.return_state = {
			.returned = FALSE,
			.retval = mv_absent(),
		}
	};

	return pvariables;
}

// ----------------------------------------------------------------
static cst_outputs_t* my_make_cst_outputs() {
	int should_emit_rec = TRUE;
	sllv_t* poutrecs = sllv_alloc();
	char* oosvar_flatten_separator = ":";
	cli_writer_opts_t* pwriter_opts = NULL;

	cst_outputs_t* pcst_outputs = mlr_malloc_or_die(sizeof(cst_outputs_t));
	*pcst_outputs = (cst_outputs_t) {
		.pshould_emit_rec         = &should_emit_rec,
		.poutrecs                 = poutrecs,
		.oosvar_flatten_separator = oosvar_flatten_separator,
		.pwriter_opts             = pwriter_opts,
	};
	return pcst_outputs;
}

// ================================================================
// TEST METHODS

// ----------------------------------------------------------------
static char* test_empty() {
	printf("%s\n", sep);
	char* mlr_dsl_expression = "";

	mlr_dsl_ast_t* past = my_mlr_dsl_parse(mlr_dsl_expression);
	mlr_dsl_cst_t* pcst = my_mlr_dsl_cst_alloc(past);

	// xxx comment
	mu_assert_lf(pcst != NULL);
	mu_assert_lf(pcst->pmain_block->pframe != NULL);
	mu_assert_lf(pcst->pmain_block->max_var_depth == 1);
	mu_assert_lf(mv_is_absent(&pcst->pmain_block->pframe->pvars[0]));

	printf("%s\n", sep);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_top_level_locals() {
	printf("%s\n", sep);
	char* mlr_dsl_expression = ""
		"local a = 1;"
		"local b = 2;"
		"local c = 3;"
		;

	mlr_dsl_ast_t* past         = my_mlr_dsl_parse(mlr_dsl_expression);
	mlr_dsl_cst_t* pcst         = my_mlr_dsl_cst_alloc(past);
	variables_t*   pvariables   = my_make_variables();
	cst_outputs_t* pcst_outputs = my_make_cst_outputs();

	// xxx comment
	mu_assert_lf(pcst != NULL);
	mu_assert_lf(pcst->pmain_block->pframe != NULL);
	mu_assert_lf(pcst->pmain_block->max_var_depth == 4);
	mv_t* pvars = pcst->pmain_block->pframe->pvars;
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));
	mu_assert_lf(mv_is_absent(&pvars[2]));
	mu_assert_lf(mv_is_absent(&pvars[3]));

	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);

	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_equals_si(&pvars[1], imv(1)));
	mu_assert_lf(mv_equals_si(&pvars[2], imv(2)));
	mu_assert_lf(mv_equals_si(&pvars[3], imv(3)));

	printf("%s\n", sep);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_top_level_clears() {
	printf("%s\n", sep);
	char* mlr_dsl_expression = ""
		"do {"
		"    if (1==1) {"
		"        local a = nonesuch;"
		"        print;"
		"        print \"NR=\".NR;"
		"        print \"oa=\".ispresent(a);"
		"        print \"oa=\".a;"
		"        if (NR == 1) {"
		"          print \"oa=\".ispresent(a);"
		"          print \"oa=\".a;"
		"        } elif (NR == 2) {"
		"          a = 999;"
		"          print \"oa=\".ispresent(a);"
		"          print \"oa=\".a;"
		"        } elif (NR == 3) {"
		"          print \"oa=\".ispresent(a);"
		"          print \"oa=\".a;"
		"        }"
		"    }"
		"} while (1==0);"
		;

	mlr_dsl_ast_t* past         = my_mlr_dsl_parse(mlr_dsl_expression);
	mlr_dsl_cst_t* pcst         = my_mlr_dsl_cst_alloc(past);
	variables_t*   pvariables   = my_make_variables();
	cst_outputs_t* pcst_outputs = my_make_cst_outputs();

	mu_assert_lf(pcst != NULL);
	mu_assert_lf(pcst->pmain_block->pframe != NULL);
	mu_assert_lf(pcst->pmain_block->max_var_depth == 2);
	mv_t* pvars = pcst->pmain_block->pframe->pvars;
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));

	pvariables->pctx->nr = 1;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));

	pvariables->pctx->nr = 2;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_equals_si(&pvars[1], imv(999)));

	pvariables->pctx->nr = 3;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));

	printf("%s\n", sep);
	return NULL;
}

// ----------------------------------------------------------------
// being tested:
// * placements
// * clears
// * stack-frame over/underflow
// * correct bindings including ispresent

// xxx assert placements & remaining absents
// xxx drive another lrec (w/ if:NR etc.)
// xxx assert clears

// xxx x all control structures (w/ nesting)
// xxx from begin/end/main/func/subr
// xxx also regtestrun:repro w/ all such cases


// ? begin top-level
// ? main  top-level
// ? end   top-level
// ? func  top-level
// ? subr  top-level

// xxx:
// cond true
// if true elif true else
// while true w/ break
// for-srec w/ non-empty lrec & break
// for full-oosvar w/ non-empty oos & break
// for full-oosvar-key-only w/ non-empty oos & break
// for oosvar w/ non-empty oos & break
// for oosvar-key-only w/ non-empty oos & break
// triple-for w/ break

// ? md_cond_block
// ? md_while_block
// ? md_for_loop_full_srec
// ? md_for_loop_full_oosvar
// ? md_for_loop_full_oosvar_key_only
// ? md_for_loop_oosvar
// ? md_for_loop_oosvar_key_only
// ? md_triple_for
// ? md_if_chain


// ================================================================
static char * run_all_tests() {
	mu_run_test(test_empty);
	mu_run_test(test_top_level_locals);
	mu_run_test(test_top_level_clears);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_DSL_STACK ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_DSL_STACK: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
