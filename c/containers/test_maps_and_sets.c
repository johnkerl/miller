#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/sllv.h"

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

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "c"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "d"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "a"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "e"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "b"));
	pe = pe->pnext;

	mu_assert_lf(pe == NULL);

	slls_sort(plist);

	mu_assert_lf(plist->length == 5);
	pe = plist->phead;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "a"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "b"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "c"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "d"));
	pe = pe->pnext;

	mu_assert_lf(pe != NULL);
	mu_assert_lf(streq(pe->value, "e"));
	pe = pe->pnext;

	mu_assert_lf(pe == NULL);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_sllv_append() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//static void print_list(sllv_t* pa, char* desc) {
//	printf("%s [%d]:\n", desc, pa->length);
//	for (sllve_t* pe = pa->phead; pe != NULL; pe = pe->pnext) {
//		printf("  %s\n", (char*)pe->pvdata);
//	}
//}

//	sllv_t* pa = sllv_alloc();
//	sllv_add(pa, "a");
//	sllv_add(pa, "b");
//	sllv_add(pa, "c");
//
//	sllv_t* pb = sllv_alloc();
//	sllv_add(pb, "d");
//	sllv_add(pb, "e");
//
//	print_list(pa, "A");
//	print_list(pb, "B");
//
//	pa = sllv_append(pa, pb);
//	print_list(pa, "A+B");

// ----------------------------------------------------------------
static char* test_hss() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	hss_t *pset = hss_alloc();
//	hss_add(pset, "x");
//	hss_add(pset, "y");
//	hss_add(pset, "x");
//	hss_add(pset, "z");
//	hss_remove(pset, "y");
//	printf("set size = %d\n", hss_size(pset));
//	hss_dump(pset);
//	hss_check_counts(pset);
//	hss_free(pset);

// ----------------------------------------------------------------
static char* test_lhms2v() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	lhms2v_t *pmap = lhms2v_alloc();
//	lhms2v_put(pmap, "a", "x", "3");
//	lhms2v_put(pmap, "a", "y", "5");
//	lhms2v_put(pmap, "a", "x", "4");
//	lhms2v_put(pmap, "b", "z", "7");
//	lhms2v_remove(pmap, "a", "y");
//	printf("map size = %d\n", lhms2v_size(pmap));
//	lhms2v_dump(pmap);
//	lhms2v_check_counts(pmap);
//	lhms2v_free(pmap);

// ----------------------------------------------------------------
static char* test_lhmsi() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	lhmsi_t *pmap = lhmsi_alloc();
//	lhmsi_put(pmap, "x", 3);
//	lhmsi_put(pmap, "y", 5);
//	lhmsi_put(pmap, "x", 4);
//	lhmsi_put(pmap, "z", 7);
//	lhmsi_remove(pmap, "y");
//	printf("map size = %d\n", pmap->num_occupied);
//	lhmsi_dump(pmap);
//	printf("map has(\"w\") = %d\n", lhmsi_has_key(pmap, "w"));
//	printf("map has(\"x\") = %d\n", lhmsi_has_key(pmap, "x"));
//	printf("map has(\"y\") = %d\n", lhmsi_has_key(pmap, "y"));
//	printf("map has(\"z\") = %d\n", lhmsi_has_key(pmap, "z"));
//	lhmsi_check_counts(pmap);
//	lhmsi_free(pmap);

// ----------------------------------------------------------------
static char* test_lhmslv() {
	mu_assert_lf(0 == 0);

	return NULL;
}

//	slls_t* ax = slls_alloc();
//	slls_add_no_free(ax, "a");
//	slls_add_no_free(ax, "x");
//
//	slls_t* ay = slls_alloc();
//	slls_add_no_free(ay, "a");
//	slls_add_no_free(ay, "y");
//
//	slls_t* bz = slls_alloc();
//	slls_add_no_free(bz, "b");
//	slls_add_no_free(bz, "z");
//
//	lhmslv_t *pmap = lhmslv_alloc();
//	lhmslv_put(pmap, ax, "3");
//	lhmslv_put(pmap, ay, "5");
//	lhmslv_put(pmap, ax, "4");
//	lhmslv_put(pmap, bz, "7");
//	lhmslv_remove(pmap, ay);
//	printf("map size = %d\n", lhmslv_size(pmap));
//	lhmslv_dump(pmap);
//	lhmslv_check_counts(pmap);
//	lhmslv_free(pmap);

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
	mu_run_test(test_lhms2v);
	mu_run_test(test_lhmsi);
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
