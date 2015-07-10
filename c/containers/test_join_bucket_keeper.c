#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"

#ifdef __TEST_JOIN_BUCKET_KEEPER_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_foo() {
	sllv_t* precords = sllv_alloc();
	sllv_add(precords, lrec_literal_2("a","1", "b","10"));
	sllv_add(precords, lrec_literal_2("a","1", "b","11"));
	sllv_add(precords, lrec_literal_2("a","2", "b","12"));
	sllv_add(precords, lrec_literal_2("a","2", "b","13"));
	sllv_add(precords, lrec_literal_2("a","3", "b","14"));
	sllv_add(precords, lrec_literal_2("a","3", "b","15"));
	lrec_reader_t* preader = lrec_reader_in_memory_alloc(precords);
	printf("#=%d\n", precords->length);
	mu_assert_lf(precords->length == 6);
	while (TRUE) {
		lrec_t* precord = preader->pprocess_func(NULL, preader->pvstate, NULL);
		if (precord == NULL)
			break;
		lrec_print(precord);
	}
	printf("#=%d\n", precords->length);
	mu_assert_lf(precords->length == 0);

	return 0;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_fotest_foo);
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
