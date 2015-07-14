#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "containers/join_bucket_keeper.h"

#ifdef __TEST_JOIN_BUCKET_KEEPER_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

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

	*ppleft_field_names = pleft_field_names;
	*ppreader = preader;
}

// ----------------------------------------------------------------
// xxx cases:
//
// * empty left
// * single-key left, right >
// * single-key left, right == then >
// * single-key left, right < then >
// * single-key left, right < then ==

// * double-key left, right >
// * double-key left, right == last then >
// * double-key left, right between then >
// * double-key left, right between then == then >




// ----------------------------------------------------------------
static char* test1() {
	printf("test1 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(&pleft_field_names, &preader);

	void* pvhandle = NULL;  // xxx move these into the jbk obj?
	context_t* pctx = NULL; // xxx revisit

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, pvhandle, pctx,
		pleft_field_names);

	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	char* rval = "0";
	slls_t* pright_field_values = slls_alloc();
	slls_add_no_free(pright_field_values, rval);
	join_bucket_keeper_emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	printf("match %s:\n", rval);
	mu_assert_lf(pbucket_paired == NULL);
	mu_assert_lf(pbucket_left_unpaired == NULL);

	join_bucket_keeper_emit(pkeeper, NULL, &pbucket_paired, &pbucket_left_unpaired);
	printf("unpaired:\n");
	mu_assert_lf(pbucket_paired == NULL);
	mu_assert_lf(pbucket_left_unpaired != NULL);
	printf("#lunp=%d\n", pbucket_left_unpaired->length);
	mu_assert_lf(pbucket_left_unpaired->length == 6);
	for (sllve_t* pe = pbucket_left_unpaired->phead; pe != NULL; pe = pe->pnext) {
		lrec_t* prec = pe->pvdata;
		lrec_print(prec);
	}

	printf("test1 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test2() {
	printf("test2 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(&pleft_field_names, &preader);

	void* pvhandle = NULL;  // xxx move these into the jbk obj?
	context_t* pctx = NULL; // xxx revisit

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, pvhandle, pctx,
		pleft_field_names);

	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	char* rval = "6";
	slls_t* pright_field_values = slls_alloc();
	slls_add_no_free(pright_field_values, rval);
	join_bucket_keeper_emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	printf("match %s:\n", rval);
	mu_assert_lf(pbucket_paired == NULL);
	mu_assert_lf(pbucket_left_unpaired != NULL);
	printf("#lunp=%d\n", pbucket_left_unpaired->length);
	mu_assert_lf(pbucket_left_unpaired->length == 6);
	mu_assert_lf(pbucket_left_unpaired == NULL);
	for (sllve_t* pe = pbucket_left_unpaired->phead; pe != NULL; pe = pe->pnext) {
		lrec_t* prec = pe->pvdata;
		lrec_print(prec);
	}

	printf("test2 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test3() {
	printf("test3 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(&pleft_field_names, &preader);

	void* pvhandle = NULL;  // xxx move these into the jbk obj?
	context_t* pctx = NULL; // xxx revisit

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, pvhandle, pctx,
		pleft_field_names);

	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_alloc();
	slls_add_no_free(pright_field_values, "0");
	join_bucket_keeper_emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	printf("match 0:\n");
	mu_assert_lf(pbucket_paired == NULL);
	mu_assert_lf(pbucket_left_unpaired == NULL);

	pright_field_values = slls_alloc();
	slls_add_no_free(pright_field_values, "1");
	join_bucket_keeper_emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	printf("match 2:\n");
	mu_assert_lf(pbucket_paired != NULL);
	mu_assert_lf(pbucket_paired->length == 2);
	mu_assert_lf(pbucket_left_unpaired == NULL);
	for (sllve_t* pe = pbucket_paired->phead; pe != NULL; pe = pe->pnext) {
		lrec_t* prec = pe->pvdata;
		lrec_print(prec);
	}

	printf("test3 exit\n");
	printf("\n");
	return 0;
}

// ----------------------------------------------------------------
static char* test4() {
	printf("test4 enter\n");

	slls_t* pleft_field_names;
	lrec_reader_t* preader;
	set_up(&pleft_field_names, &preader);

	void* pvhandle = NULL;  // xxx move these into the jbk obj?
	context_t* pctx = NULL; // xxx revisit

	join_bucket_keeper_t* pkeeper = join_bucket_keeper_alloc_from_reader(preader, pvhandle, pctx,
		pleft_field_names);

	sllv_t* pbucket_paired;
	sllv_t* pbucket_left_unpaired;

	slls_t* pright_field_values = slls_alloc();
	slls_add_no_free(pright_field_values, "2");

	join_bucket_keeper_emit(pkeeper, pright_field_values, &pbucket_paired, &pbucket_left_unpaired);
	mu_assert_lf(pbucket_paired == NULL);
	mu_assert_lf(pbucket_left_unpaired != NULL);
	mu_assert_lf(pbucket_left_unpaired->length == 2);

	printf("test4 exit\n");
	printf("\n");
	return 0;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test1);
	mu_run_test(test2);
	mu_run_test(test3);
	mu_run_test(test4);
	return 0;
}

int main(int argc, char **argv) {
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
