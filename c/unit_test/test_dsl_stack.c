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

// ----------------------------------------------------------------
// xxx being tested:
// * placements
// * clears
// * stack-frame over/underflow
// * correct bindings including ispresent

// xxx:
// ? begin top-level
// ? main  top-level
// ? end   top-level
// ? func  top-level
// ? subr  top-level

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

static mv_t* smv(char* strv) {
	mv_t* pmv = mlr_malloc_or_die(sizeof(mv_t));
	*pmv = mv_from_string(strv, NO_FREE);
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
static char* test_placements() {
	printf("%s\n", sep);
	char* mlr_dsl_expression = ""
		"                                                      # absent-null is at local 0\n"
		"@v[NR] = NR*1000;\n"
		"true {\n"
		"    local val = nonesuch;                             # local 1\n"
		"    while (1 == 1) {\n"
		"        if (NR == 1) {\n"
		"            for (k, v in $*) {                        # locals 2 and 3\n"
		"                val = v;\n"
		"            }\n"
		"        } elif (NR == 2) {\n"
		"            for (k in @v) {                           # local 2\n"
		"            }\n"
		"        } elif (NR == 3) {\n"
		"            for (k, v in @v) {                        # locals 2 and 3\n"
		"            }\n"
		"        } elif (NR == 4) {\n"
		"            local sum = 0;                            # local 2\n"
		"            for (local i = 1; i <= 10; i += 1) {      # local 3\n"
		"                for (local j = 1; j <= 10; j += 1) {  # local 4\n"
		"                    sum += i*j;\n"
		"                }\n"
		"            }\n"
		"            val = sum;\n"
		"        }\n"
		"        break;\n"
		"    }\n"
		"}\n"
		;

	mlr_dsl_ast_t* past         = my_mlr_dsl_parse(mlr_dsl_expression);
	mlr_dsl_cst_t* pcst         = my_mlr_dsl_cst_alloc(past);
	variables_t*   pvariables   = my_make_variables();
	cst_outputs_t* pcst_outputs = my_make_cst_outputs();
	lrec_put(pvariables->pinrec, "xyz", "ijkl", NO_FREE);

	printf("Expression:\n");
	printf("%s\n", mlr_dsl_expression);

	mu_assert_lf(pcst != NULL);
	mu_assert_lf(pcst->pmain_block->pframe != NULL);
	mu_assert_lf(pcst->pmain_block->max_var_depth == 5);
	mv_t* pvars = pcst->pmain_block->pframe->pvars;

	// Before any records are read
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));
	mu_assert_lf(mv_is_absent(&pvars[2]));
	mu_assert_lf(mv_is_absent(&pvars[3]));
	mu_assert_lf(mv_is_absent(&pvars[4]));

	// Push record 1
	pvariables->pctx->nr = 1;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	printf("\n");
	for (int i = 1; i <= 4; i++) { printf("1:%d:%s\n", i, mv_alloc_format_val(&pvars[i])); }
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_equals_si(&pvars[1], smv("ijkl")));
	mu_assert_lf(mv_equals_si(&pvars[2], smv("xyz")));
	mu_assert_lf(mv_equals_si(&pvars[3], smv("ijkl")));
	mu_assert_lf(mv_is_absent(&pvars[4]));

	// Push record 2
	pvariables->pctx->nr = 2;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	printf("\n");
	for (int i = 1; i <= 4; i++) { printf("2:%d:%s\n", i, mv_alloc_format_val(&pvars[i])); }
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));
	mu_assert_lf(mv_equals_si(&pvars[2], imv(2)));
	mu_assert_lf(mv_equals_si(&pvars[3], smv("ijkl")));
	mu_assert_lf(mv_is_absent(&pvars[4]));

	// Push record 3
	pvariables->pctx->nr = 3;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	printf("\n");
	for (int i = 1; i <= 4; i++) { printf("3:%d:%s\n", i, mv_alloc_format_val(&pvars[i])); }
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_is_absent(&pvars[1]));
	mu_assert_lf(mv_equals_si(&pvars[2], imv(3)));
	mu_assert_lf(mv_equals_si(&pvars[3], imv(3000)));
	mu_assert_lf(mv_is_absent(&pvars[4]));

	// Push record 4
	pvariables->pctx->nr = 4;
	mlr_dsl_cst_handle_top_level_statement_block(pcst->pmain_block, pvariables, pcst_outputs);
	printf("\n");
	for (int i = 1; i <= 4; i++) { printf("4:%d:%s\n", i, mv_alloc_format_val(&pvars[i])); }
	mu_assert_lf(mv_is_absent(&pvars[0]));
	mu_assert_lf(mv_equals_si(&pvars[1], imv(3025)));
	mu_assert_lf(mv_equals_si(&pvars[2], imv(3025)));
	mu_assert_lf(mv_equals_si(&pvars[3], imv(11)));
	mu_assert_lf(mv_equals_si(&pvars[3], imv(11)));

	printf("%s\n", sep);
	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_empty);
	mu_run_test(test_top_level_locals);
	mu_run_test(test_top_level_clears);
	mu_run_test(test_placements);
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
