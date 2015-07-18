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
static void set_up(
	slls_t** ppleft_field_names,
	lrec_reader_t** ppreader)
{
	slls_t* pleft_field_names = slls_alloc();
	slls_add_no_free(pleft_field_names, "l");

	sllv_t* precords = sllv_alloc();
	sllv_add(precords, lrec_literal_2("l","1", "b","10"));
	sllv_add(precords, lrec_literal_2("l","1", "b","11"));
	sllv_add(precords, lrec_literal_2("l","3", "b","12"));
	sllv_add(precords, lrec_literal_2("l","3", "b","13"));
	sllv_add(precords, lrec_literal_2("l","3", "b","14"));
	sllv_add(precords, lrec_literal_2("l","5", "b","15"));
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
		printf("BEFORE\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}
	join_bucket_keeper_emit(pkeeper, pright_field_values, ppbucket_paired, ppbucket_left_unpaired);
	if (tjbk_verbose) {
		printf("AFTER\n");
		join_bucket_keeper_print(pkeeper);
		printf("\n");
	}
}

static int list_is_null(sllv_t* plist, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%s is null with rval=\"%s\"; ok\n", desc, rval);
		return TRUE;
	} else {
		printf("%s should be null with rval=\"%s\" and is not:\n", desc, rval);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

static int list_has_length(sllv_t* plist, int length, char* desc, char* rval) {
	if (plist == NULL) {
		printf("%s is null with rval=\"%s\" and should not be\n", desc, rval);
		return FALSE;
	}

	if (plist->length == length) {
		printf("%s length is %d with rval=\"%s\"; ok:\n", desc, length, rval);
		lrec_print_list_with_prefix(plist, "  ");
		return TRUE;
	} else {
		printf("%s length is %d with rval=\"%s\" but should be %d:\n", desc, plist->length, rval, length);
		lrec_print_list_with_prefix(plist, "  ");
		return FALSE;
	}
}

// ----------------------------------------------------------------
static char* test1() {
	printf("----------------------------------------------------------------\n");
	printf("test1 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(&pleft_field_names, &preader);

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
	set_up(&pleft_field_names, &preader);

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);
	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("6");
	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 6, "unpaired", pright_field_values->phead->value));

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
	set_up(&pleft_field_names, &preader);

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
	set_up(&pleft_field_names, &preader);

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, NULL, pleft_field_names);

	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_single_no_free("2");

	emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);

	mu_assert_lf(list_is_null(pbucket_paired, "paired", pright_field_values->phead->value));
	mu_assert_lf(list_has_length(pbucket_left_unpaired, 2, "unpaired", pright_field_values->phead->value));

	printf("\n");
	printf("test4 exit\n");
	printf("\n");
	return 0;
}

// test4 enter
// left records:
//   l=1,b=10
//   l=1,b=11
//   l=3,b=12
//   l=3,b=13
//   l=3,b=14
//   l=5,b=15
//
// paired should be null with rval="2" and is not:
//   l=3,b=12
//   l=3,b=13
//   l=3,b=14

// typedef struct _join_bucket_keeper_t {
//  lrec_reader_t* plrec_reader;
//  void*          pvhandle;
//  context_t*     pctx;
//
//  slls_t*        pleft_field_names;
//
//  join_bucket_t* pbucket;
//
//  lrec_t*        prec_peek;
//  int            leof;
//  int            state;
// } join_bucket_keeper_t;

// ================================================================
static char * run_all_tests() {
	mu_run_test(test1);
	mu_run_test(test2);
	mu_run_test(test3);
	mu_run_test(test4);
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
