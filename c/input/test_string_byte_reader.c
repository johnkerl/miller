#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"

#ifdef __TEST_STRING_BYTE_READER_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_it() {
	mu_assert_lf(0 == 0);
	mu_assert_lf(1 == 1);
	mu_assert_lf(2 == 2);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_it);
	return 0;
}

int main(int argc, char **argv) {
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_STRING_BYTE_READER: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_STRING_BYTE_READER_MAIN__
