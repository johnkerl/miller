#include <stdio.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "lib/minunit.h"
#include "input/byte_readers.h"

#ifdef __TEST_STRING_BYTE_READER_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_it() {
	byte_reader_t* pbr = string_byte_reader_alloc();

	int ok = pbr->popen_func(pbr, "");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

	ok = pbr->popen_func(pbr, "a");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

	ok = pbr->popen_func(pbr, "abc");
	mu_assert_lf(ok == TRUE);
	mu_assert_lf(pbr->pread_func(pbr) == 'a');
	mu_assert_lf(pbr->pread_func(pbr) == 'b');
	mu_assert_lf(pbr->pread_func(pbr) == 'c');
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	mu_assert_lf(pbr->pread_func(pbr) == EOF);
	pbr->pclose_func(pbr);

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
