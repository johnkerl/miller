#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "containers/join_bucket_keeper.h"
#include "containers/mixutil.h"

#ifdef __TEST_JOIN_BUCKET_KEEPER_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static int tjbk_verbose = FALSE;

// ----------------------------------------------------------------
static sllv_t* make_records_empty() {
	sllv_t* precords = sllv_alloc();
	return precords;
}

// ----------------------------------------------------------------
static sllv_t* make_records_113335() {
	sllv_t* precords = sllv_alloc();
	sllv_add(precords, lrec_literal_2("l","1", "b","10"));
	sllv_add(precords, lrec_literal_2("l","1", "b","11"));
	sllv_add(precords, lrec_literal_2("l","3", "b","12"));
	sllv_add(precords, lrec_literal_2("l","3", "b","13"));
	sllv_add(precords, lrec_literal_2("l","3", "b","14"));
	sllv_add(precords, lrec_literal_2("l","5", "b","15"));
	return precords;
}

// ----------------------------------------------------------------
static void set_up(
	sllv_t* precords,
	slls_t** ppleft_field_names,
	lrec_reader_t** ppreader)
{
	slls_t* pleft_field_names = slls_alloc();
	slls_add_no_free(pleft_field_names, "l");

	lrec_reader_t* preader = lrec_reader_in_memory_alloc(precords);
	printf("left records:\n");
	lrec_print_list_with_prefix(precords, "  ");
	printf("\n");

	*ppleft_field_names = pleft_field_names;
	*ppreader = preader;
}

static void emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	if (tjbk_verbose) {
		printf("BEFORE EMIT\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}

	join_bucket_keeper_emit(pkeeper, pright_field_values, ppbucket_paired, ppbucket_left_unpaired);

	if (tjbk_verbose) {
		printf("AFTER EMIT\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}
}

static int list_is_null(sllv_t* plist, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%-8s is null with rval=\"%s\"; ok\n", desc, rval);
		return TRUE;
	} else {
		printf("%-8s should be null with rval=\"%s\" and is not (length=%d):\n",
			desc, rval, plist->length);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

static int list_has_length(sllv_t* plist, int length, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%-8s is null with rval=\"%s\" and should not be\n", desc, rval);
		return FALSE;
	}

	if (plist->length == length) {
		printf("%-8s length=%d with rval=\"%s\"; ok:\n", desc, length, rval);
		lrec_print_list_with_prefix(plist, "  ");
		return TRUE;
	} else {
		printf("%-8s length=%d with rval=\"%s\" but should be %d:\n", desc, plist->length, rval, length);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

// ----------------------------------------------------------------
static char* test00() {
	printf("----------------------------------------------------------------\n");
	printf("test00 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_empty(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test00 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test01() {
	printf("----------------------------------------------------------------\n");
	printf("test01 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_empty(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test01 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test0() {
	printf("----------------------------------------------------------------\n");
	printf("test0 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 6, "unpaired", "(eof)"));
	printf("\n");

	printf("test0 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test1() {
	printf("----------------------------------------------------------------\n");
	printf("test1 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 6, "unpaired", "(eof)"));
	printf("\n");

	printf("test1 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test2() {
	printf("----------------------------------------------------------------\n");
	printf("test2 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 6, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 0, "unpaired", "(eof)"));
	printf("\n");

	printf("test2 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test3() {
	printf("----------------------------------------------------------------\n");
	printf("test3 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 4, "unpaired", "(eof)"));
	printf("\n");

	printf("test3 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test4() {
	printf("----------------------------------------------------------------\n");
	printf("test4 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 2, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 1, "unpaired", "(eof)"));
	printf("\n");

	printf("test4 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test5() {
	printf("----------------------------------------------------------------\n");
	printf("test5 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(make_records_113335(), &pleft_field_names, &preader);
	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("0");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("1");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 2, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("2");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("3");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 3, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("4");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("4");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("5");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 1, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("5");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_has_length(pbucket_paired, 1, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_is_null(pbucket_left_unpaired, "unpaired", pright_field_values->phead->value));
	printf("\n");

	emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", "(eof)"));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 0, "unpaired", pright_field_values->phead->value));
	printf("\n");

	printf("test5 exit\n");
	printf("\n");
	return 0;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test00);
	mu_run_test(test01);
	mu_run_test(test0);
	mu_run_test(test1);
	mu_run_test(test2);
	mu_run_test(test3);
	mu_run_test(test4);
	mu_run_test(test5);
	return 0;
}

// xxx make a -v flag with conditional bucket-dumps :)
int main(int argc, char **argv) {
	if ((argc == 2) && streq(argv[1], "-v"))
		tjbk_verbose = TRUE;

	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_JOIN_BUCKET_KEEPER: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_JOIN_BUCKET_KEEPER_MAIN__
