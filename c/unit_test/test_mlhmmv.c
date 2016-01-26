#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_stub() {
	mlhmmv_t* pmap = mlhmmv_alloc();

	mv_t key = mv_from_int(3LL);
	sllmv_t* pmkeys = sllmv_single(&key);
	mv_t value = mv_from_int(4LL);

	mlhmmv_put(pmap, pmkeys, &value);

	int ret = mlhmmv_has_keys(pmap, pmkeys);
	mu_assert_lf(ret == FALSE); // xxx stub

	mv_t* pback = mlhmmv_get(pmap, pmkeys);
	mu_assert_lf(pback == NULL); // xxx stub

	sllmv_free(pmkeys);
	mlhmmv_free(pmap);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_stub);
	return 0;
}

int main(int argc, char **argv) {
	printf("TEST_MLHMMV ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MLHMMV: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
