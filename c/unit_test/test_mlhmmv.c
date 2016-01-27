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

	printf("----------------------------------------------------------------\n");
	printf("empty map:\n");
	mlhmmv_print(pmap);

	mv_t key1 = mv_from_int(3LL);
	sllmv_t* pmvkeys1 = sllmv_single(&key1);
	mv_t value1 = mv_from_int(4LL);

	mlhmmv_put(pmap, pmvkeys1, &value1);
	printf("map:\n");
	mlhmmv_print(pmap);

	int ret = mlhmmv_has_keys(pmap, pmvkeys1);
	mu_assert_lf(ret == FALSE); // xxx stub

	mv_t* pback = mlhmmv_get(pmap, pmvkeys1);
	mu_assert_lf(pback == NULL); // xxx stub


	mv_t key2a = mv_from_string("abcde", NO_FREE);
	mv_t key2b = mv_from_int(-6LL);
	sllmv_t* pmvkeys2 = sllmv_double(&key2a, &key2b);
	mv_t value2 = mv_from_int(7LL);

	mlhmmv_put(pmap, pmvkeys2, &value2);
	printf("map:\n");
	mlhmmv_print(pmap);

	ret = mlhmmv_has_keys(pmap, pmvkeys2);
	mu_assert_lf(ret == FALSE); // xxx stub

	pback = mlhmmv_get(pmap, pmvkeys2);
	mu_assert_lf(pback == NULL); // xxx stub


	mv_t key3a = mv_from_int(0LL);
	mv_t key3b = mv_from_string("fghij", NO_FREE);
	mv_t key3c = mv_from_int(0LL);
	sllmv_t* pmvkeys3 = sllmv_triple(&key3a, &key3b, &key3c);
	mv_t value3 = mv_from_int(17LL);

	mlhmmv_put(pmap, pmvkeys3, &value3);
	printf("map:\n");
	mlhmmv_print(pmap);

	ret = mlhmmv_has_keys(pmap, pmvkeys3);
	mu_assert_lf(ret == FALSE); // xxx stub

	pback = mlhmmv_get(pmap, pmvkeys3);
	mu_assert_lf(pback == NULL); // xxx stub


	sllmv_free(pmvkeys1);
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
