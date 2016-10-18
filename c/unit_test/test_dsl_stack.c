#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "dsls/mlr_dsl_wrapper.h"
#include "mapping/mlr_dsl_ast.h"
//#include "containers/lrec.h"
//#include "containers/sllv.h"
//#include "input/lrec_readers.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_empty() {
	int trace_parse = FALSE;
	char* mlr_dsl_expression = "";

	mlr_dsl_ast_t* past = mlr_dsl_parse(mlr_dsl_expression, trace_parse);
	// xxx funcify
	if (past == NULL) {
		fprintf(stderr, "Syntax error on DSL parse of '%s'\n",
			mlr_dsl_expression);
		return NULL;
	}
	mlr_dsl_ast_print(past);

//	lrec_t* prec = lrec_unbacked_alloc();
//	mu_assert_lf(prec->field_count == 0);
//
//	lrec_put(prec, "x", "3", NO_FREE);
//	mu_assert_lf(prec->field_count == 1);
//	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
//
//	lrec_put(prec, "y", "4", NO_FREE);
//	mu_assert_lf(prec->field_count == 2);
//	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
//	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));
//
//	lrec_put(prec, "x", "5", NO_FREE);
//	mu_assert_lf(prec->field_count == 2);
//	mu_assert_lf(streq(lrec_get(prec, "x"), "5"));
//	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));
//
//	lrec_remove(prec, "x");
//	mu_assert_lf(prec->field_count == 1);
//	mu_assert_lf(lrec_get(prec, "x") == NULL);
//
//	// Non-replacing-rename case
//	//lrec_dump_titled("Before rename", prec);
//	lrec_rename(prec, "y", "z", FALSE);
//	//lrec_dump_titled("After rename", prec);
//	mu_assert_lf(prec->field_count == 1);
//	mu_assert_lf(lrec_get(prec, "x") == NULL);
//	mu_assert_lf(lrec_get(prec, "y") == NULL);
//	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));
//
//	lrec_free(prec);
//
//	// Replacing-rename case
//	prec = lrec_unbacked_alloc();
//
//	lrec_put(prec, "x", "3", NO_FREE);
//	lrec_put(prec, "y", "4", NO_FREE);
//	lrec_put(prec, "z", "5", NO_FREE);
//	mu_assert_lf(prec->field_count == 3);
//
//	//lrec_dump_titled("Before rename", prec);
//	lrec_rename(prec, "y", "z", FALSE);
//	//lrec_dump_titled("After rename", prec);
//
//	mu_assert_lf(prec->field_count == 2);
//	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
//	mu_assert_lf(lrec_get(prec, "y") == NULL);
//	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));
//
//	lrec_free(prec);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_empty);
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
