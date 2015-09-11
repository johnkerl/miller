#include <stdio.h>
#include <string.h>
#ifdef MLR_USE_MCHECK
#include <mcheck.h>
#endif // MLR_USE_MCHECK
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"

#ifdef __TEST_LREC_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_lrec_unbacked_api() {
	lrec_t* prec = lrec_unbacked_alloc();
	mu_assert_lf(prec->field_count == 0);

	lrec_put_no_free(prec, "x", "3");
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));

	lrec_put_no_free(prec, "y", "4");
	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));

	lrec_put_no_free(prec, "x", "5");
	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "x"), "5"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));

	lrec_remove(prec, "x");
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(lrec_get(prec, "x") == NULL);

	// Non-replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z");
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	// Replacing-rename case
	prec = lrec_unbacked_alloc();

	lrec_put_no_free(prec, "x", "3");
	lrec_put_no_free(prec, "y", "4");
	lrec_put_no_free(prec, "z", "5");
	mu_assert_lf(prec->field_count == 3);

	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z");
	//lrec_dump_titled("After rename", prec);

	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_dkvp_api() {
	char* line = strdup("w=2,x=3,y=4,z=5");
	lrec_t* prec = lrec_parse_stdio_dkvp(line, ',', '=', FALSE);
	mu_assert_lf(prec->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec, "w"), "2"));
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));
	mu_assert_lf(streq(lrec_get(prec, "z"), "5"));

	lrec_remove(prec, "w");
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "w") == NULL);

	// Non-replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "x", "u");
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z");
	//lrec_dump_titled("After rename", prec);

	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_nidx_api() {
	char* line = strdup("a,b,c,d");
	lrec_t* prec = lrec_parse_stdio_nidx(line, ',', FALSE);
	mu_assert_lf(prec->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec, "1"), "a"));
	mu_assert_lf(streq(lrec_get(prec, "2"), "b"));
	mu_assert_lf(streq(lrec_get(prec, "3"), "c"));
	mu_assert_lf(streq(lrec_get(prec, "4"), "d"));

	lrec_remove(prec, "1");
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "1") == NULL);

	// Non-replacing-rename case
	lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "2", "u");
	lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "2") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "b"));

	// Replacing-rename case
	lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "3", "4");
	lrec_dump_titled("After rename", prec);

	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "u"), "b"));
	mu_assert_lf(lrec_get(prec, "3") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "4"), "c"));

	lrec_free(prec);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_csv_api() {
	char* hdr_line = strdup("w,x,y,z");
	slls_t* hdr_fields = split_csvlite_header_line(hdr_line, ',', FALSE);
	header_keeper_t* pheader_keeper = header_keeper_alloc(hdr_line, hdr_fields);

	char* data_line_1 = strdup("2,3,4,5");
	lrec_t* prec_1 = lrec_parse_stdio_csvlite_data_line(pheader_keeper, data_line_1, ',', FALSE);

	char* data_line_2 = strdup("6,7,8,9");
	lrec_t* prec_2 = lrec_parse_stdio_csvlite_data_line(pheader_keeper, data_line_2, ',', FALSE);

	mu_assert_lf(prec_1->field_count == 4);
	mu_assert_lf(prec_2->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec_1, "w"), "2"));
	mu_assert_lf(streq(lrec_get(prec_1, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec_1, "y"), "4"));
	mu_assert_lf(streq(lrec_get(prec_1, "z"), "5"));

	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));
	mu_assert_lf(streq(lrec_get(prec_2, "x"), "7"));
	mu_assert_lf(streq(lrec_get(prec_2, "y"), "8"));
	mu_assert_lf(streq(lrec_get(prec_2, "z"), "9"));

	lrec_remove(prec_1, "w");
	mu_assert_lf(prec_1->field_count == 3);
	mu_assert_lf(prec_2->field_count == 4);
	mu_assert_lf(lrec_get(prec_1, "w") == NULL);
	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));

	// Non-replacing-rename case
	//lrec_dump_titled("Before rename", prec_1);
	lrec_rename(prec_1, "x", "u");
	//lrec_dump_titled("After rename", prec_1);
	mu_assert_lf(prec_1->field_count == 3);
	mu_assert_lf(lrec_get(prec_1, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec_1, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec_2);
	lrec_rename(prec_2, "y", "z");
	//lrec_dump_titled("After rename", prec_2);

	mu_assert_lf(prec_2->field_count == 3);
	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));
	mu_assert_lf(streq(lrec_get(prec_2, "x"), "7"));
	mu_assert_lf(lrec_get(prec_2, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec_2, "z"), "8"));

	lrec_free(prec_1);
	lrec_free(prec_2);

	// xxx need a test case for alloc1,free1,alloc2,free2 w/ same hdr.
	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_xtab_api() {
	char* line_1 = strdup("w 2");
	char* line_2 = strdup("x    3");
	char* line_3 = strdup("y 4");
	char* line_4 = strdup("z  5");
	slls_t* pxtab_lines = slls_alloc();
	slls_add_with_free(pxtab_lines, line_1);
	slls_add_with_free(pxtab_lines, line_2);
	slls_add_with_free(pxtab_lines, line_3);
	slls_add_with_free(pxtab_lines, line_4);

	lrec_t* prec = lrec_parse_stdio_xtab(pxtab_lines, ' ', TRUE);
	mu_assert_lf(prec->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec, "w"), "2"));
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));
	mu_assert_lf(streq(lrec_get(prec, "z"), "5"));

	lrec_remove(prec, "w");
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "w") == NULL);

	// Non-replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "x", "u");
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z");
	//lrec_dump_titled("After rename", prec);

	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_lrec_unbacked_api);
	mu_run_test(test_lrec_dkvp_api);
	mu_run_test(test_lrec_nidx_api);
	mu_run_test(test_lrec_csv_api);
	mu_run_test(test_lrec_xtab_api);
	return 0;
}

int main(int argc, char **argv) {
#ifdef MLR_USE_MCHECK
	if (mcheck(NULL) != 0) {
		printf("Could not set up mcheck\n");
		exit(1);
	}
	printf("Set up mcheck\n");
#endif // MLR_USE_MCHECK

	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_LREC: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_LREC_MAIN__
