#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"
#include "lib/mvfuncs.h"

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
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();
	int error = 0;

	printf("================================================================\n");
	printf("empty map:\n");
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);

	sllmv_t* pmvkeys1 = sllmv_single_with_free(imv(3));
	mv_t value1 = mv_from_int(4LL);
	printf("\n");
	printf("keys1:  ");
	sllmv_print(pmvkeys1);
	printf("value1: %s\n", mv_alloc_format_val(&value1));
	mlhmmv_root_put_terminal(pmap, pmvkeys1, &value1);
	printf("map:\n");
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys1, &error), &value1));

	sllmv_t* pmvkeys2 = sllmv_double_with_free(smv("abcde"), imv(-6));
	mv_t value2 = mv_from_int(7);
	printf("\n");
	printf("keys2:  ");
	sllmv_print(pmvkeys2);
	printf("value2: %s\n", mv_alloc_format_val(&value2));
	mlhmmv_root_put_terminal(pmap, pmvkeys2, &value2);
	printf("map:\n");
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys2, &error), &value2));

	sllmv_t* pmvkeys3 = sllmv_triple_with_free(imv(0), smv("fghij"), imv(0));
	mv_t value3 = mv_from_int(0LL);
	printf("\n");
	printf("keys3:  ");
	sllmv_print(pmvkeys3);
	printf("value3: %s\n", mv_alloc_format_val(&value3));
	mlhmmv_root_put_terminal(pmap, pmvkeys3, &value3);
	printf("map:\n");
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys3, &error), &value3));

	sllmv_free(pmvkeys1);
	mlhmmv_root_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_overlap() {
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();
	int error = 0;

	printf("================================================================\n");
	sllmv_t* pmvkeys = sllmv_single_with_free(imv(3));
	mv_t* ptermval = imv(4);
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	ptermval = imv(5);
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double_with_free(imv(3), smv("x"));
	ptermval = imv(6);
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	ptermval = imv(7);
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple_with_free(imv(3), imv(9), smv("y"));
	ptermval = smv("z");
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple_with_free(imv(3), imv(9), smv("z"));
	ptermval = smv("y");
	mlhmmv_root_put_terminal(pmap, pmvkeys, ptermval);
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	mlhmmv_root_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_resize() {
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();
	int error;

	printf("================================================================\n");
	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_root_put_terminal(pmap, sllmv_single_with_free(imv(i)), imv(-i));
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	printf("\n");

	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_root_put_terminal(pmap, sllmv_double_with_free(smv("a"), imv(i)), imv(-i));
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	printf("\n");

	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(imv(i*100), imv(i % 4), smv("b")), smv("term"));
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);

	sllmv_t* pmvkeys = sllmv_single_with_free(imv(2));
	mv_t* ptermval = imv(-2);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double_with_free(smv("a"), imv(9));
	ptermval = imv(-9);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double_with_free(smv("a"), imv(31));
	ptermval = imv(-31);
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple_with_free(imv(0), imv(0), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple_with_free(imv(100), imv(1), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple_with_free(imv(1700), imv(1), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_root_look_up_and_ref_terminal(pmap, pmvkeys, &error), ptermval));

	mlhmmv_root_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_depth_errors() {
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();
	int error;

	printf("================================================================\n");
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(imv(1), imv(2), imv(3)), imv(4));

	mu_assert_lf(NULL != mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_triple_with_free(imv(1), imv(2), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_triple_with_free(imv(0), imv(2), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_triple_with_free(imv(1), imv(0), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_triple_with_free(imv(1), imv(2), imv(0)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_quadruple_with_free(imv(1), imv(2), imv(3), imv(4)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_KEYLIST_TOO_DEEP);

	mu_assert_lf(NULL == mlhmmv_root_look_up_and_ref_terminal(pmap, sllmv_double_with_free(imv(1), imv(2)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_KEYLIST_TOO_SHALLOW);

	mlhmmv_root_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_mlhmmv_to_lrecs() {
	printf("================================================================\n");
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();

	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("s"), smv("x")), imv(1));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("s"), smv("y")), imv(2));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("t"), smv("x")), imv(3));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("t"), smv("y")), imv(4));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("s"), smv("x")), imv(5));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("s"), smv("y")), imv(6));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("t"), smv("x")), imv(7));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("t"), smv("y")), imv(8));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("u"), smv("x")), imv(9));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("triple"), smv("u"), smv("y")), imv(10));
	mlhmmv_root_put_terminal(pmap, sllmv_triple_with_free(smv("not"), smv("u"), smv("y")), imv(11));
	mlhmmv_root_put_terminal(pmap, sllmv_single_with_free(smv("single")), imv(99));
	mlhmmv_root_put_terminal(pmap, sllmv_double_with_free(smv("double"), smv("i")), imv(111));
	mlhmmv_root_put_terminal(pmap, sllmv_double_with_free(smv("double"), smv("j")), imv(112));
	mlhmmv_root_put_terminal(pmap, sllmv_double_with_free(smv("double"), smv("k")), imv(113));

	printf("full map:\n");
	mlhmmv_root_print_json_stacked(pmap, TRUE, FALSE, "", "\n", stdout);
	printf("\n");

	sllv_t* poutrecs = NULL;

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("single")),
		sllmv_alloc(), poutrecs, TRUE, ":");
	printf("single-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("single")),
		sllmv_single_with_free(smv("first")), poutrecs, TRUE, ":");
	printf("single-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("double")),
		sllmv_alloc(), poutrecs, TRUE, ":");
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("double")),
		sllmv_single_with_free(smv("first")), poutrecs, TRUE, ":");
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("double")),
		sllmv_double_with_free(smv("first"), smv("second")), poutrecs, TRUE, ":");
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("triple")),
		sllmv_alloc(), poutrecs, TRUE, ":");
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("triple")),
		sllmv_single_with_free(smv("first")), poutrecs, TRUE, ":");
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("triple")),
		sllmv_double_with_free(smv("first"), smv("second")), poutrecs, TRUE, ":");
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_root_partial_to_lrecs(pmap, sllmv_single_with_free(smv("triple")),
		sllmv_triple_with_free(smv("first"), smv("second"), smv("third")), poutrecs, TRUE, ":");
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	mlhmmv_root_free(pmap);
	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_no_overlap);
	mu_run_test(test_overlap);
	mu_run_test(test_resize);
	mu_run_test(test_depth_errors);
	mu_run_test(test_mlhmmv_to_lrecs);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);
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
