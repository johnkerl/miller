#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"

int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

static mv_t* smv(char* strv) {
	mv_t* pmv = mlr_malloc_or_die(sizeof(mv_t));
	*pmv = mv_from_string(strv, NO_FREE);
	return pmv;
}
static mv_t* imv(long long intv) {
	mv_t* pmv = mlr_malloc_or_die(sizeof(mv_t));
	*pmv = mv_from_int(intv);
	return pmv;
}

// ----------------------------------------------------------------
static char* test_no_overlap() {
	mlhmmv_t* pmap = mlhmmv_alloc();

	printf("----------------------------------------------------------------\n");
	printf("empty map:\n");
	mlhmmv_print(pmap);

	sllmv_t* pmvkeys1 = sllmv_single(imv(3));
	mv_t value1 = mv_from_int(4LL);
	printf("\n");
	printf("keys1:  ");
	sllmv_print(pmvkeys1);
	printf("value1: %s\n", mv_alloc_format_val(&value1));
	mlhmmv_put(pmap, pmvkeys1, &value1);
	printf("map:\n");
	mlhmmv_print(pmap);

	mu_assert_lf(mlhmmv_has_keys(pmap, pmvkeys1) == FALSE); // xxx stub
	mu_assert_lf(mlhmmv_get(pmap, pmvkeys1) == NULL); // xxx stub


	sllmv_t* pmvkeys2 = sllmv_double(smv("abcde"), imv(-6));
	mv_t value2 = mv_from_int(7);
	printf("\n");
	printf("keys2:  ");
	sllmv_print(pmvkeys2);
	printf("value2: %s\n", mv_alloc_format_val(&value2));
	mlhmmv_put(pmap, pmvkeys2, &value2);
	printf("map:\n");
	mlhmmv_print(pmap);

	mu_assert_lf(mlhmmv_has_keys(pmap, pmvkeys2) == FALSE); // xxx stub
	mu_assert_lf(mlhmmv_get(pmap, pmvkeys2) == NULL); // xxx stub


	sllmv_t* pmvkeys3 = sllmv_triple(imv(0), smv("fghij"), imv(0));
	mv_t value3 = mv_from_int(0LL);
	printf("\n");
	printf("keys3:  ");
	sllmv_print(pmvkeys3);
	printf("value3: %s\n", mv_alloc_format_val(&value3));
	mlhmmv_put(pmap, pmvkeys3, &value3);
	printf("map:\n");
	mlhmmv_print(pmap);

	mu_assert_lf(mlhmmv_has_keys(pmap, pmvkeys3) == FALSE); // xxx stub
	mu_assert_lf(mlhmmv_get(pmap, pmvkeys3) == NULL); // xxx stub


	sllmv_free(pmvkeys1);
	mlhmmv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_overlap() {
	mlhmmv_t* pmap = mlhmmv_alloc();

	printf("----------------------------------------------------------------\n");
	mlhmmv_put(pmap, sllmv_single(imv(3)), imv(4));
	mlhmmv_print(pmap);
	mlhmmv_put(pmap, sllmv_single(imv(3)), imv(5));
	mlhmmv_print(pmap);

	mlhmmv_put(pmap, sllmv_double(imv(3), smv("x")), imv(6));
	mlhmmv_print(pmap);
	mlhmmv_put(pmap, sllmv_double(imv(3), smv("x")), imv(7));
	mlhmmv_print(pmap);

	mlhmmv_put(pmap, sllmv_triple(imv(3), imv(9), smv("y")), smv("z"));
	mlhmmv_print(pmap);
	mlhmmv_put(pmap, sllmv_triple(imv(3), imv(9), smv("z")), smv("y"));
	mlhmmv_print(pmap);

	mlhmmv_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_resize() {
	mlhmmv_t* pmap = mlhmmv_alloc();

	printf("----------------------------------------------------------------\n");
	for (int i = 0; i < 32; i++)
		mlhmmv_put(pmap, sllmv_single(imv(i)), imv(-i));
	mlhmmv_print(pmap);

	for (int i = 0; i < 32; i++)
		mlhmmv_put(pmap, sllmv_double(smv("a"), imv(i)), imv(-i));
	mlhmmv_print(pmap);

	for (int i = 0; i < 32; i++)
		mlhmmv_put(pmap, sllmv_triple(imv(i*100), imv(i % 4), smv("b")), smv("term"));
	mlhmmv_print(pmap);

	mlhmmv_free(pmap);
	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_no_overlap);
	mu_run_test(test_overlap);
	mu_run_test(test_resize);
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
