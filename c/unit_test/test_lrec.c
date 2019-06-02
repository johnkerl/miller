#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_lrec_unbacked_api() {
	lrec_t* prec = lrec_unbacked_alloc();
	mu_assert_lf(prec->field_count == 0);

	lrec_put(prec, "x", "3", NO_FREE);
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));

	lrec_put(prec, "y", "4", NO_FREE);
	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));

	lrec_put(prec, "x", "5", NO_FREE);
	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "x"), "5"));
	mu_assert_lf(streq(lrec_get(prec, "y"), "4"));

	lrec_remove(prec, "x");
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(lrec_get(prec, "x") == NULL);

	// Non-replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z", FALSE);
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 1);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	// Replacing-rename case
	prec = lrec_unbacked_alloc();

	lrec_put(prec, "x", "3", NO_FREE);
	lrec_put(prec, "y", "4", NO_FREE);
	lrec_put(prec, "z", "5", NO_FREE);
	mu_assert_lf(prec->field_count == 3);

	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z", FALSE);
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
	char* line = mlr_strdup_or_die("w=2,x=3,y=4,z=5");

	lrec_t* prec = lrec_parse_stdio_dkvp_single_sep(line, ',', '=', FALSE);
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
	lrec_rename(prec, "x", "u", FALSE);
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z", FALSE);
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
	char* line = mlr_strdup_or_die("a,b,c,d");
	lrec_t* prec = lrec_parse_stdio_nidx_single_sep(line, ',', FALSE);
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
	lrec_rename(prec, "2", "u", FALSE);
	lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "2") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "b"));

	// Replacing-rename case
	lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "3", "4", FALSE);
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
	char* hdr_line = mlr_strdup_or_die("w,x,y,z");
	slls_t* hdr_fields = split_csvlite_header_line_single_ifs(hdr_line, ',', FALSE);
	header_keeper_t* pheader_keeper = header_keeper_alloc(hdr_line, hdr_fields);

	char* data_line_1 = mlr_strdup_or_die("2,3,4,5");
	lrec_t* prec_1 = lrec_parse_stdio_csvlite_data_line_single_ifs(pheader_keeper, "test-file", 999,
		data_line_1, ',', FALSE, FALSE);

	char* data_line_2 = mlr_strdup_or_die("6,7,8,9");
	lrec_t* prec_2 = lrec_parse_stdio_csvlite_data_line_single_ifs(pheader_keeper, "test-file", 999,
		data_line_2, ',', FALSE, FALSE);

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
	lrec_rename(prec_1, "x", "u", FALSE);
	//lrec_dump_titled("After rename", prec_1);
	mu_assert_lf(prec_1->field_count == 3);
	mu_assert_lf(lrec_get(prec_1, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec_1, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec_2);
	lrec_rename(prec_2, "y", "z", FALSE);
	//lrec_dump_titled("After rename", prec_2);

	mu_assert_lf(prec_2->field_count == 3);
	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));
	mu_assert_lf(streq(lrec_get(prec_2, "x"), "7"));
	mu_assert_lf(lrec_get(prec_2, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec_2, "z"), "8"));

	lrec_free(prec_1);
	lrec_free(prec_2);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_csv_api_disjoint_allocs() {
	char* hdr_line = mlr_strdup_or_die("w,x,y,z");
	slls_t* hdr_fields = split_csvlite_header_line_single_ifs(hdr_line, ',', FALSE);
	header_keeper_t* pheader_keeper = header_keeper_alloc(hdr_line, hdr_fields);


	char* data_line_1 = mlr_strdup_or_die("2,3,4,5");
	lrec_t* prec_1 = lrec_parse_stdio_csvlite_data_line_single_ifs(pheader_keeper, "test-file", 999,
		data_line_1, ',', FALSE, FALSE);

	mu_assert_lf(prec_1->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec_1, "w"), "2"));
	mu_assert_lf(streq(lrec_get(prec_1, "x"), "3"));
	mu_assert_lf(streq(lrec_get(prec_1, "y"), "4"));
	mu_assert_lf(streq(lrec_get(prec_1, "z"), "5"));

	lrec_remove(prec_1, "w");
	mu_assert_lf(prec_1->field_count == 3);
	mu_assert_lf(lrec_get(prec_1, "w") == NULL);

	lrec_rename(prec_1, "x", "u", FALSE);
	mu_assert_lf(prec_1->field_count == 3);
	mu_assert_lf(lrec_get(prec_1, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec_1, "u"), "3"));

	lrec_free(prec_1);


	char* data_line_2 = mlr_strdup_or_die("6,7,8,9");
	lrec_t* prec_2 = lrec_parse_stdio_csvlite_data_line_single_ifs(pheader_keeper, "test-file", 999,
		data_line_2, ',', FALSE, FALSE);

	mu_assert_lf(prec_2->field_count == 4);

	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));
	mu_assert_lf(streq(lrec_get(prec_2, "x"), "7"));
	mu_assert_lf(streq(lrec_get(prec_2, "y"), "8"));
	mu_assert_lf(streq(lrec_get(prec_2, "z"), "9"));

	mu_assert_lf(prec_2->field_count == 4);
	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));

	lrec_rename(prec_2, "y", "z", FALSE);

	mu_assert_lf(prec_2->field_count == 3);
	mu_assert_lf(streq(lrec_get(prec_2, "w"), "6"));
	mu_assert_lf(streq(lrec_get(prec_2, "x"), "7"));
	mu_assert_lf(lrec_get(prec_2, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec_2, "z"), "8"));

	lrec_free(prec_2);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_xtab_api() {
	char* line_1 = mlr_strdup_or_die("w 2");
	char* line_2 = mlr_strdup_or_die("x    3");
	char* line_3 = mlr_strdup_or_die("y 4");
	char* line_4 = mlr_strdup_or_die("z  5");
	slls_t* pxtab_lines = slls_alloc();
	slls_append_with_free(pxtab_lines, line_1);
	slls_append_with_free(pxtab_lines, line_2);
	slls_append_with_free(pxtab_lines, line_3);
	slls_append_with_free(pxtab_lines, line_4);

	lrec_t* prec = lrec_parse_stdio_xtab_single_ips(pxtab_lines, ' ', TRUE);
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
	lrec_rename(prec, "x", "u", FALSE);
	//lrec_dump_titled("After rename", prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(lrec_get(prec, "x") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));

	// Replacing-rename case
	//lrec_dump_titled("Before rename", prec);
	lrec_rename(prec, "y", "z", FALSE);
	//lrec_dump_titled("After rename", prec);

	mu_assert_lf(prec->field_count == 2);
	mu_assert_lf(streq(lrec_get(prec, "u"), "3"));
	mu_assert_lf(lrec_get(prec, "y") == NULL);
	mu_assert_lf(streq(lrec_get(prec, "z"), "4"));

	lrec_free(prec);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lrec_put_after() {
	printf("TEST_LREC_PUT_AFTER ENTER\n");

	lrec_t* prec = lrec_literal_1("a", "1");
	lrece_t* pe = NULL;

	printf("BEFORE: "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 1);
	char* value = lrec_get_ext(prec, "nosuch", &pe);
	mu_assert_lf(value == NULL);
	mu_assert_lf(pe == NULL);

	value = lrec_get_ext(prec, "a", &pe);
	mu_assert_lf(value != NULL);
	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(value, "1"));
	lrec_put_after(prec, pe, "a", "2", NO_FREE);
	printf("AFTER : "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 1);
	value = lrec_get_ext(prec, "a", &pe);
	mu_assert_lf(value != NULL);
	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(value, "2"));
	lrec_free(prec);

	prec = lrec_literal_1("a", "1");
	printf("BEFORE: "); lrec_print(prec);
	value = lrec_get_ext(prec, "a", &pe);
	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(value, "1"));

	lrec_put_after(prec, pe, "b", "3", NO_FREE);
	printf("AFTER : "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 2);
	value = lrec_get(prec, "a");
	mu_assert_lf(value != NULL);
	mu_assert_lf(streq(value, "1"));
	value = lrec_get(prec, "b");
	mu_assert_lf(value != NULL);
	mu_assert_lf(streq(value, "3"));
	mu_assert_lf(streq(prec->phead->key, "a"));
	mu_assert_lf(streq(prec->phead->pnext->key, "b"));
	mu_assert_lf(streq(prec->phead->value, "1"));
	mu_assert_lf(streq(prec->phead->pnext->value, "3"));
	mu_assert_lf(prec->phead->pnext->pnext == NULL);
	lrec_free(prec);


	prec = lrec_literal_2("a", "1", "b", "2");
	printf("BEFORE: "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 2);
	value = lrec_get_ext(prec, "a", &pe);
	mu_assert_lf(value != NULL);
	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(value, "1"));

	lrec_put_after(prec, pe, "z", "9", NO_FREE);
	printf("AFTER : "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(streq(prec->phead->key, "a"));
	mu_assert_lf(streq(prec->phead->pnext->key, "z"));
	mu_assert_lf(streq(prec->phead->pnext->pnext->key, "b"));
	mu_assert_lf(streq(prec->phead->value, "1"));
	mu_assert_lf(streq(prec->phead->pnext->value, "9"));
	mu_assert_lf(streq(prec->phead->pnext->pnext->value, "2"));
	mu_assert_lf(prec->phead->pnext->pnext->pnext == NULL);
	lrec_free(prec);


	prec = lrec_literal_2("a", "1", "b", "2");
	printf("BEFORE: "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 2);
	value = lrec_get_ext(prec, "b", &pe);
	mu_assert_lf(value != NULL);
	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(value, "2"));

	lrec_put_after(prec, pe, "z", "9", NO_FREE);
	printf("AFTER : "); lrec_print(prec);
	mu_assert_lf(prec->field_count == 3);
	mu_assert_lf(streq(prec->phead->key, "a"));
	mu_assert_lf(streq(prec->phead->pnext->key, "b"));
	mu_assert_lf(streq(prec->phead->pnext->pnext->key, "z"));
	mu_assert_lf(streq(prec->phead->value, "1"));
	mu_assert_lf(streq(prec->phead->pnext->value, "2"));
	mu_assert_lf(streq(prec->phead->pnext->pnext->value, "9"));
	mu_assert_lf(prec->phead->pnext->pnext->pnext == NULL);
	lrec_free(prec);

	printf("TEST_LREC_PUT_AFTER EXIT\n");
	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_lrec_unbacked_api);
	mu_run_test(test_lrec_dkvp_api);
	mu_run_test(test_lrec_nidx_api);
	mu_run_test(test_lrec_csv_api);
	mu_run_test(test_lrec_csv_api_disjoint_allocs);
	mu_run_test(test_lrec_xtab_api);
	mu_run_test(test_lrec_put_after);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
	printf("TEST_LREC ENTER\n");
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
