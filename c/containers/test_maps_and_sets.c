#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "containers/lhmsi.h"
#include "containers/lhms2v.h"
#include "containers/lhmslv.h"

#ifdef __TEST_MAPS_AND_SETS_MAIN__
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char* test_slls() {
	slls_t* plist = slls_from_line(strdup(""), ',', FALSE);
	mu_assert_lf(plist->length == 0);

	plist = slls_from_line(strdup("a"), ',', FALSE);
	mu_assert_lf(plist->length == 1);

	plist = slls_from_line(strdup("c,d,a,e,b"), ',', FALSE);
	mu_assert_lf(plist->length == 5);

	sllse_t* pe = plist->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "e")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "b")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	slls_sort(plist);

	mu_assert_lf(plist->length == 5);
	pe = plist->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_sllv_append() {
	mu_assert_lf(0 == 0);

	sllv_t* pa = sllv_alloc();
	sllv_add(pa, "a");
	sllv_add(pa, "b");
	sllv_add(pa, "c");
	mu_assert_lf(pa->length == 3);

	sllve_t* pe = pa->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "c")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	sllv_t* pb = sllv_alloc();
	sllv_add(pb, "d");
	sllv_add(pb, "e");
	mu_assert_lf(pb->length == 2);

	pe = pb->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	pa = sllv_append(pa, pb);

	mu_assert_lf(pa->length == 5);
	mu_assert_lf(pb->length == 2);

	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	pe = pb->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvdata, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_hss() {

	hss_t *pset = hss_alloc();
	mu_assert_lf(pset->num_occupied == 0);

	hss_add(pset, "x");
	mu_assert_lf(pset->num_occupied == 1);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(hss_has(pset, "x"));
	mu_assert_lf(!hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "y");
	mu_assert_lf(pset->num_occupied == 2);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(hss_has(pset, "x"));
	mu_assert_lf(hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "x");
	mu_assert_lf(pset->num_occupied == 2);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(hss_has(pset, "x"));
	mu_assert_lf(hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "z");
	mu_assert_lf(pset->num_occupied == 3);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(hss_has(pset, "x"));
	mu_assert_lf(hss_has(pset, "y"));
	mu_assert_lf(hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_remove(pset, "y");
	mu_assert_lf(pset->num_occupied == 2);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(hss_has(pset, "x"));
	mu_assert_lf(!hss_has(pset, "y"));
	mu_assert_lf(hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_clear(pset);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf(!hss_has(pset, "x"));
	mu_assert_lf(!hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_free(pset);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsi() {
	mu_assert_lf(0 == 0);

	lhmsi_t *pmap = lhmsi_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(!lhmsi_has_key(pmap, "x"));
	mu_assert_lf(!lhmsi_has_key(pmap, "y"));
	mu_assert_lf(!lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 3);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(lhmsi_has_key(pmap, "x"));
	mu_assert_lf(!lhmsi_has_key(pmap, "y"));
	mu_assert_lf(!lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "y", 5);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(lhmsi_has_key(pmap, "x"));
	mu_assert_lf(lhmsi_has_key(pmap, "y"));
	mu_assert_lf(!lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 4);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(lhmsi_has_key(pmap, "x"));
	mu_assert_lf(lhmsi_has_key(pmap, "y"));
	mu_assert_lf(!lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "z", 7);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(lhmsi_has_key(pmap, "x"));
	mu_assert_lf(lhmsi_has_key(pmap, "y"));
	mu_assert_lf(lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_remove(pmap, "y");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(lhmsi_has_key(pmap, "x"));
	mu_assert_lf(!lhmsi_has_key(pmap, "y"));
	mu_assert_lf(lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_clear(pmap);
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsi_has_key(pmap, "w"));
	mu_assert_lf(!lhmsi_has_key(pmap, "x"));
	mu_assert_lf(!lhmsi_has_key(pmap, "y"));
	mu_assert_lf(!lhmsi_has_key(pmap, "z"));
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhms2v() {
	mu_assert_lf(0 == 0);

	lhms2v_t *pmap = lhms2v_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "3");
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "y", "5");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "4");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "b", "z", "7");
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_remove(pmap, "a", "y");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_clear(pmap);
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmslv() {
	mu_assert_lf(0 == 0);

	slls_t* ax = slls_alloc();
	slls_add_no_free(ax, "a");
	slls_add_no_free(ax, "x");
	// xxx assertions here

	slls_t* ay = slls_alloc();
	slls_add_no_free(ay, "a");
	slls_add_no_free(ay, "y");

	slls_t* bz = slls_alloc();
	slls_add_no_free(bz, "b");
	slls_add_no_free(bz, "z");

	lhmslv_t *pmap = lhmslv_alloc();
	lhmslv_put(pmap, ax, "3");
	lhmslv_put(pmap, ay, "5");
	lhmslv_put(pmap, ax, "4");
	lhmslv_put(pmap, bz, "7");
	lhmslv_remove(pmap, ay);

	lhmslv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmss() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	lhmss_t *pmap = lhmss_alloc();
//	lhmss_put(pmap, "x", "3");
//	lhmss_put(pmap, "y", "5");
//	lhmss_put(pmap, "x", "4");
//	lhmss_put(pmap, "z", "7");
//	lhmss_remove(pmap, "y");
//	printf("map size = %d\n", pmap->num_occupied);
//	lhmss_dump(pmap);
//	lhmss_check_counts(pmap);
//	lhmss_free(pmap);

// ----------------------------------------------------------------
static char* test_lhmsv() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	int x3 = 3;
//	int x5 = 5;
//	int x4 = 4;
//	int x7 = 7;
//	lhmsv_t *pmap = lhmsv_alloc();
//	lhmsv_put(pmap, "x", &x3);
//	lhmsv_put(pmap, "y", &x5);
//	lhmsv_put(pmap, "x", &x4);
//	lhmsv_put(pmap, "z", &x7);
//	lhmsv_remove(pmap, "y");
//	printf("map size = %d\n", pmap->num_occupied);
//	lhmsv_dump(pmap);
//	lhmsv_check_counts(pmap);
//	lhmsv_free(pmap);

// ----------------------------------------------------------------
static char* test_percentile_keeper() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//void percentile_keeper_dump(percentile_keeper_t* ppercentile_keeper) {
//	for (int i = 0; i < ppercentile_keeper->size; i++)
//		printf("[%02d] %.8lf\n", i, ppercentile_keeper->data[i]);
//}

//	char buffer[1024];
//	percentile_keeper_t* ppercentile_keeper = percentile_keeper_alloc();
//	char* line;
//	while ((line = fgets(buffer, sizeof(buffer), stdin)) != NULL) {
//		int len = strlen(line);
//		if (len >= 1) // xxx write and use a chomp()
//			if (line[len-1] == '\n')
//				line[len-1] = 0;
//		double v;
//		if (!mlr_try_double_from_string(line, &v)) {
//			percentile_keeper_ingest(ppercentile_keeper, v);
//		} else {
//			printf("meh? >>%s<<\n", line);
//		}
//	}
//	percentile_keeper_dump(ppercentile_keeper);
//	printf("\n");
//	double p;
//	p = 0.10; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
//	p = 0.50; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
//	p = 0.90; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
//	printf("\n");
//	percentile_keeper_dump(ppercentile_keeper);

// ----------------------------------------------------------------
static char* test_top_keeper() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//void top_keeper_dump(top_keeper_t* ptop_keeper) {
//	for (int i = 0; i < ptop_keeper->size; i++)
//		printf("[%02d] %.8lf\n", i, ptop_keeper->top_values[i]);
//	for (int i = ptop_keeper->size; i < ptop_keeper->capacity; i++)
//		printf("[%02d] ---\n", i);
//}

//	int capacity = 5;
//	char buffer[1024];
//	if (argc == 2)
//		(void)sscanf(argv[1], "%d", &capacity);
//	top_keeper_t* ptop_keeper = top_keeper_alloc(capacity);
//	char* line;
//	while ((line = fgets(buffer, sizeof(buffer), stdin)) != NULL) {
//		int len = strlen(line);
//		if (len >= 1) // xxx write and use a chomp()
//			if (line[len-1] == '\n')
//				line[len-1] = 0;
//		if (streq(line, "")) {
//			//top_keeper_dump(ptop_keeper);
//			printf("\n");
//		} else {
//			double v;
//			if (!mlr_try_double_from_string(line, &v)) {
//				top_keeper_add(ptop_keeper, v, NULL);
//				top_keeper_dump(ptop_keeper);
//				printf("\n");
//			} else {
//				printf("meh? >>%s<<\n", line);
//			}
//		}
//	}

// ----------------------------------------------------------------
static char* test_dheap() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	dheap_t *pdheap = dheap_alloc();
//	dheap_check(pdheap, __FILE__,  __LINE__);
//	dheap_add(pdheap, 4.1);
//	dheap_add(pdheap, 3.1);
//	dheap_add(pdheap, 2.1);
//	dheap_add(pdheap, 6.1);
//	dheap_add(pdheap, 5.1);
//	dheap_add(pdheap, 8.1);
//	dheap_add(pdheap, 7.1);
//	dheap_print(pdheap);
//	dheap_check(pdheap, __FILE__,  __LINE__);
//
//	printf("\n");
//	printf("remove %lf\n", dheap_remove(pdheap));
//	printf("remove %lf\n", dheap_remove(pdheap));
//	printf("remove %lf\n", dheap_remove(pdheap));
//	printf("remove %lf\n", dheap_remove(pdheap));
//	printf("\n");
//
//	dheap_print(pdheap);
//	dheap_check(pdheap, __FILE__,  __LINE__);
//
//	dheap_free(pdheap);

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_slls);
	mu_run_test(test_sllv_append);
	mu_run_test(test_hss);
	mu_run_test(test_lhmsi);
	mu_run_test(test_lhms2v);
	mu_run_test(test_lhmslv);
	mu_run_test(test_lhmss);
	mu_run_test(test_lhmsv);
	mu_run_test(test_percentile_keeper);
	mu_run_test(test_top_keeper);
	mu_run_test(test_dheap);
	return 0;
}

int main(int argc, char **argv) {
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MAPS_AND_SETS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
#endif // __TEST_MAPS_AND_SETS_MAIN__
