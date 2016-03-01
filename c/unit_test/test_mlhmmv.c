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
	int error = 0;

	printf("================================================================\n");
	printf("empty map:\n");
	mlhmmv_print_json_stacked(pmap, FALSE);

	sllmv_t* pmvkeys1 = sllmv_single(imv(3));
	mv_t value1 = mv_from_int(4LL);
	printf("\n");
	printf("keys1:  ");
	sllmv_print(pmvkeys1);
	printf("value1: %s\n", mv_alloc_format_val(&value1));
	mlhmmv_put(pmap, pmvkeys1, &value1);
	printf("map:\n");
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys1, &error), &value1));

	sllmv_t* pmvkeys2 = sllmv_double(smv("abcde"), imv(-6));
	mv_t value2 = mv_from_int(7);
	printf("\n");
	printf("keys2:  ");
	sllmv_print(pmvkeys2);
	printf("value2: %s\n", mv_alloc_format_val(&value2));
	mlhmmv_put(pmap, pmvkeys2, &value2);
	printf("map:\n");
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys2, &error), &value2));

	sllmv_t* pmvkeys3 = sllmv_triple(imv(0), smv("fghij"), imv(0));
	mv_t value3 = mv_from_int(0LL);
	printf("\n");
	printf("keys3:  ");
	sllmv_print(pmvkeys3);
	printf("value3: %s\n", mv_alloc_format_val(&value3));
	mlhmmv_put(pmap, pmvkeys3, &value3);
	printf("map:\n");
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys3, &error), &value3));

	sllmv_free(pmvkeys1);
	mlhmmv_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_overlap() {
	mlhmmv_t* pmap = mlhmmv_alloc();
	int error = 0;

	printf("================================================================\n");
	sllmv_t* pmvkeys = sllmv_single(imv(3));
	mv_t* ptermval = imv(4);
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	ptermval = imv(5);
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double(imv(3), smv("x"));
	ptermval = imv(6);
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	ptermval = imv(7);
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple(imv(3), imv(9), smv("y"));
	ptermval = smv("z");
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple(imv(3), imv(9), smv("z"));
	ptermval = smv("y");
	mlhmmv_put(pmap, pmvkeys, ptermval);
	mlhmmv_print_json_stacked(pmap, FALSE);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	mlhmmv_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_resize() {
	mlhmmv_t* pmap = mlhmmv_alloc();
	int error;

	printf("================================================================\n");
	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_put(pmap, sllmv_single(imv(i)), imv(-i));
	mlhmmv_print_json_stacked(pmap, FALSE);
	printf("\n");

	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_put(pmap, sllmv_double(smv("a"), imv(i)), imv(-i));
	mlhmmv_print_json_stacked(pmap, FALSE);
	printf("\n");

	for (int i = 0; i < 2*MLHMMV_INITIAL_ARRAY_LENGTH; i++)
		mlhmmv_put(pmap, sllmv_triple(imv(i*100), imv(i % 4), smv("b")), smv("term"));
	mlhmmv_print_json_stacked(pmap, FALSE);

	sllmv_t* pmvkeys = sllmv_single(imv(2));
	mv_t* ptermval = imv(-2);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double(smv("a"), imv(9));
	ptermval = imv(-9);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_double(smv("a"), imv(31));
	ptermval = imv(-31);
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple(imv(0), imv(0), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple(imv(100), imv(1), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	pmvkeys = sllmv_triple(imv(1700), imv(1), smv("b"));
	ptermval = smv("term");
	mu_assert_lf(mv_equals_si(mlhmmv_get(pmap, pmvkeys, &error), ptermval));

	mlhmmv_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_depth_errors() {
	mlhmmv_t* pmap = mlhmmv_alloc();
	int error;

	printf("================================================================\n");
	mlhmmv_put(pmap, sllmv_triple(imv(1), imv(2), imv(3)), imv(4));

	mu_assert_lf(NULL != mlhmmv_get(pmap, sllmv_triple(imv(1), imv(2), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_get(pmap, sllmv_triple(imv(0), imv(2), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_get(pmap, sllmv_triple(imv(1), imv(0), imv(3)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_get(pmap, sllmv_triple(imv(1), imv(2), imv(0)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_NONE);

	mu_assert_lf(NULL == mlhmmv_get(pmap, sllmv_quadruple(imv(1), imv(2), imv(3), imv(4)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_KEYLIST_TOO_DEEP);

	mu_assert_lf(NULL == mlhmmv_get(pmap, sllmv_double(imv(1), imv(2)), &error));
	mu_assert_lf(error == MLHMMV_ERROR_KEYLIST_TOO_SHALLOW);

	mlhmmv_free(pmap);
	return NULL;
}

// ----------------------------------------------------------------
// $ mlr put -q '@sum[$a][$b][$i]=$x; end{dump}' ../data/small | mlr --ijson --oxtab cat
// sum:pan:pan:1=0.346790
// sum:pan:wye:10=0.502626
// sum:eks:pan:2=0.758680
// sum:eks:wye:4=0.381399
// sum:eks:zee:7=0.611784
// sum:wye:wye:3=0.204603
// sum:wye:pan:5=0.573289
// sum:zee:pan:6=0.527126
// sum:zee:wye:8=0.598554
// sum:hat:wye:9=0.031442


// $ mlr put -q '@sum[$a][$b][$i]=$x; end{emit @sum}' ../data/small
// sum:pan:pan:1=0.346790
// sum:pan:wye:10=0.502626
// sum:eks:pan:2=0.758680
// sum:eks:wye:4=0.381399
// sum:eks:zee:7=0.611784
// sum:wye:wye:3=0.204603
// sum:wye:pan:5=0.573289
// sum:zee:pan:6=0.527126
// sum:zee:wye:8=0.598554
// sum:hat:wye:9=0.031442

// $ mlr put -q '@sum[$a][$b][$i]=$x; end{emit @sum, "a"}' ../data/small
// a=pan,sum:pan:1=0.346790
// a=pan,sum:wye:10=0.502626
// a=eks,sum:pan:2=0.758680
// a=eks,sum:wye:4=0.381399
// a=eks,sum:zee:7=0.611784
// a=wye,sum:wye:3=0.204603
// a=wye,sum:pan:5=0.573289
// a=zee,sum:pan:6=0.527126
// a=zee,sum:wye:8=0.598554
// a=hat,sum:wye:9=0.031442

// $ mlr put -q '@sum[$a][$b][$i]=$x; end{emit @sum, "a","b"}' ../data/small
// a=pan,b=pan,sum:1=0.346790
// a=pan,b=wye,sum:10=0.502626
// a=eks,b=pan,sum:2=0.758680
// a=eks,b=wye,sum:4=0.381399
// a=eks,b=zee,sum:7=0.611784
// a=wye,b=wye,sum:3=0.204603
// a=wye,b=pan,sum:5=0.573289
// a=zee,b=pan,sum:6=0.527126
// a=zee,b=wye,sum:8=0.598554
// a=hat,b=wye,sum:9=0.031442

// $ mlr put -q '@sum[$a][$b][$i]=$x; end{emit @sum, "a","b","i"}' ../data/small
// a=pan,b=pan,i=1,sum=0.346790
// a=pan,b=wye,i=10,sum=0.502626
// a=eks,b=pan,i=2,sum=0.758680
// a=eks,b=wye,i=4,sum=0.381399
// a=eks,b=zee,i=7,sum=0.611784
// a=wye,b=wye,i=3,sum=0.204603
// a=wye,b=pan,i=5,sum=0.573289
// a=zee,b=pan,i=6,sum=0.527126
// a=zee,b=wye,i=8,sum=0.598554
// a=hat,b=wye,i=9,sum=0.031442

static char* test_mlhmmv_to_lrecs() {
	printf("================================================================\n");
	mlhmmv_t* pmap = mlhmmv_alloc();

	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("s"), smv("x")), imv(1));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("s"), smv("y")), imv(2));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("t"), smv("x")), imv(3));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("t"), smv("y")), imv(4));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("s"), smv("x")), imv(5));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("s"), smv("y")), imv(6));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("t"), smv("x")), imv(7));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("t"), smv("y")), imv(8));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("u"), smv("x")), imv(9));
	mlhmmv_put(pmap, sllmv_triple(smv("triple"), smv("u"), smv("y")), imv(10));
	mlhmmv_put(pmap, sllmv_triple(smv("not"), smv("u"), smv("y")), imv(11));
	mlhmmv_put(pmap, sllmv_single(smv("single")), imv(99));
	mlhmmv_put(pmap, sllmv_double(smv("double"), smv("i")), imv(111));
	mlhmmv_put(pmap, sllmv_double(smv("double"), smv("j")), imv(112));
	mlhmmv_put(pmap, sllmv_double(smv("double"), smv("k")), imv(113));

	printf("full map:\n");
	mlhmmv_print_json_stacked(pmap, FALSE);
	printf("\n");

	sllv_t* poutrecs = NULL;

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_single(smv("single")), poutrecs);
	printf("single-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_double(smv("single"), smv("first")), poutrecs);
	printf("single-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_single(smv("double")), poutrecs);
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_double(smv("double"), smv("first")), poutrecs);
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_triple(smv("double"), smv("first"), smv("second")), poutrecs);
	printf("double-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	printf("----------------------------------------------------------------\n");
	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_single(smv("triple")), poutrecs);
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_double(smv("triple"), smv("first")), poutrecs);
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_triple(smv("triple"), smv("first"), smv("second")), poutrecs);
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	poutrecs = sllv_alloc();
	mlhmmv_to_lrecs(pmap, sllmv_quadruple(smv("triple"), smv("first"), smv("second"), smv("third")), poutrecs);
	printf("triple-level outrecs (%lld):\n", poutrecs->length);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_print(pe->pvvalue);
	for (sllve_t* pe = poutrecs->phead; pe != NULL; pe = pe->pnext)
		lrec_free(pe->pvvalue);
	printf("\n");
	sllv_free(poutrecs);

	mlhmmv_free(pmap);
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
