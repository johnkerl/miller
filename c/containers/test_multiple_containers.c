#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "containers/lhmsi.h"
#include "containers/lhmss.h"
#include "containers/lhmsv.h"
#include "containers/lhms2v.h"
#include "containers/lhmslv.h"
#include "containers/percentile_keeper.h"
#include "containers/top_keeper.h"
#include "containers/dheap.h"

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
static char* test_sllv() {

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

	lhmsi_t *pmap = lhmsi_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 3);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "x"));  mu_assert_lf(lhmsi_get(pmap, "x") == 3);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "y", 5);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "x"));  mu_assert_lf(lhmsi_get(pmap, "x") == 3);
	mu_assert_lf(lhmsi_has_key(pmap, "y"));  mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 4);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "x"));  mu_assert_lf(lhmsi_get(pmap, "x") == 4);
	mu_assert_lf(lhmsi_has_key(pmap, "y"));  mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "z", 7);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "x"));  mu_assert_lf(lhmsi_get(pmap, "x") == 4);
	mu_assert_lf(lhmsi_has_key(pmap, "y"));  mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(lhmsi_has_key(pmap, "z"));  mu_assert_lf(lhmsi_get(pmap, "z") == 7);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_remove(pmap, "y");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "x"));  mu_assert_lf(lhmsi_get(pmap, "x") == 4);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == -999);
	mu_assert_lf(lhmsi_has_key(pmap, "z"));  mu_assert_lf(lhmsi_get(pmap, "z") == 7);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_clear(pmap);
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "x") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmss() {

	lhmss_t *pmap = lhmss_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(!lhmss_has_key(pmap, "x")); mu_assert_lf(lhmss_get(pmap, "x") == NULL);
	mu_assert_lf(!lhmss_has_key(pmap, "y")); mu_assert_lf(lhmss_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));
	lhmss_put(pmap, "x", "3");
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "x"));  mu_assert_lf(streq(lhmss_get(pmap, "x"), "3"));
	mu_assert_lf(!lhmss_has_key(pmap, "y")); mu_assert_lf(lhmss_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "y", "5");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "x"));  mu_assert_lf(streq(lhmss_get(pmap, "x"), "3"));
	mu_assert_lf(lhmss_has_key(pmap, "y"));  mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "x", "4");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "x"));  mu_assert_lf(streq(lhmss_get(pmap, "x"), "4"));
	mu_assert_lf(lhmss_has_key(pmap, "y"));  mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "z", "7");
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "x"));  mu_assert_lf(streq(lhmss_get(pmap, "x"), "4"));
	mu_assert_lf(lhmss_has_key(pmap, "y"));  mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf(lhmss_has_key(pmap, "z"));  mu_assert_lf(streq(lhmss_get(pmap, "z"), "7"));
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_remove(pmap, "y");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "x"));  mu_assert_lf(streq(lhmss_get(pmap, "x"), "4"));
	mu_assert_lf(!lhmss_has_key(pmap, "y")); mu_assert_lf(lhmss_get(pmap, "y") == NULL);
	mu_assert_lf(lhmss_has_key(pmap, "z"));  mu_assert_lf(streq(lhmss_get(pmap, "z"), "7"));
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsv() {

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
//	lhmsv_print(pmap);
//	lhmsv_check_counts(pmap);
//	lhmsv_free(pmap);

// ----------------------------------------------------------------
static char* test_lhms2v() {

	// xxx more assertions here
	lhms2v_t *pmap = lhms2v_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(lhms2v_get(pmap, "a", "x") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "3");
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "a", "x"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "3"));
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "y", "5");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "a", "x"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "3"));
	mu_assert_lf(lhms2v_has_key(pmap, "a", "y"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "4");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "a", "x"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "4"));
	mu_assert_lf(lhms2v_has_key(pmap, "a", "y"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "b", "z", "7");
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "a", "x"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "4"));
	mu_assert_lf(lhms2v_has_key(pmap, "a", "y"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf(lhms2v_has_key(pmap, "b", "z"));  mu_assert_lf(streq(lhms2v_get(pmap, "b", "z"), "7"));
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_remove(pmap, "a", "y");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "a", "x"));  mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "4"));
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(lhms2v_has_key(pmap, "b", "z"));  mu_assert_lf(streq(lhms2v_get(pmap, "b", "z"), "7"));
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_clear(pmap);
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(lhms2v_get(pmap, "a", "x") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmslv() {

	slls_t* aw = slls_alloc(); slls_add_no_free(aw, "a"); slls_add_no_free(aw, "w");
	slls_t* ax = slls_alloc(); slls_add_no_free(ax, "a"); slls_add_no_free(ax, "x");
	slls_t* ay = slls_alloc(); slls_add_no_free(ay, "a"); slls_add_no_free(ay, "y");
	slls_t* bz = slls_alloc(); slls_add_no_free(bz, "b"); slls_add_no_free(bz, "z");

	lhmslv_t *pmap = lhmslv_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ax)); mu_assert_lf(lhmslv_get(pmap, ax) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ax, "3");
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "3"));
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ay, "5");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "3"));
	mu_assert_lf(lhmslv_has_key(pmap,  ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ax, "4");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "4"));
	mu_assert_lf(lhmslv_has_key(pmap,  ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, bz, "7");
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "4"));
	mu_assert_lf(lhmslv_has_key(pmap,  ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf(lhmslv_has_key(pmap,  bz)); mu_assert_lf(streq(lhmslv_get(pmap, bz), "7"));
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_remove(pmap, ay);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "4"));
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(lhmslv_has_key(pmap,  bz)); mu_assert_lf(streq(lhmslv_get(pmap, bz), "7"));
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_clear(pmap);
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ax)); mu_assert_lf(lhmslv_get(pmap, ax) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_percentile_keeper() {

	percentile_keeper_t* ppercentile_keeper = percentile_keeper_alloc();
	percentile_keeper_ingest(ppercentile_keeper, 1.0);
	percentile_keeper_ingest(ppercentile_keeper, 2.0);
	percentile_keeper_ingest(ppercentile_keeper, 3.0);
	percentile_keeper_ingest(ppercentile_keeper, 4.0);
	percentile_keeper_ingest(ppercentile_keeper, 5.0);
	percentile_keeper_print(ppercentile_keeper);

	double p, q;
	p = 0.0;
	q = percentile_keeper_emit(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q);
	mu_assert_lf(q == 1.0);

	p = 10.0;
	q = percentile_keeper_emit(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q);
	mu_assert_lf(q == 1.0);

	p = 50.0;
	q = percentile_keeper_emit(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q);
	mu_assert_lf(q == 3.0);

	p = 90.0;
	q = percentile_keeper_emit(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q);
	mu_assert_lf(q == 5.0);

	p = 100.0;
	q = percentile_keeper_emit(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q);
	mu_assert_lf(q == 5.0);

	percentile_keeper_free(ppercentile_keeper);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_top_keeper() {
	int capacity = 3;

	top_keeper_t* ptop_keeper = top_keeper_alloc(capacity);
	mu_assert_lf(ptop_keeper->size == 0);

	top_keeper_add(ptop_keeper, 5.0, NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 1);
	mu_assert_lf(ptop_keeper->top_values[0] == 5.0);

	top_keeper_add(ptop_keeper, 6.0, NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 2);
	mu_assert_lf(ptop_keeper->top_values[0] == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1] == 5.0);

	top_keeper_add(ptop_keeper, 4.0, NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0] == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1] == 5.0);
	mu_assert_lf(ptop_keeper->top_values[2] == 4.0);

	top_keeper_add(ptop_keeper, 2.0, NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0] == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1] == 5.0);
	mu_assert_lf(ptop_keeper->top_values[2] == 4.0);

	top_keeper_add(ptop_keeper, 7.0, NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0] == 7.0);
	mu_assert_lf(ptop_keeper->top_values[1] == 6.0);
	mu_assert_lf(ptop_keeper->top_values[2] == 5.0);

	top_keeper_free(ptop_keeper);
	return NULL;
}

// ----------------------------------------------------------------
static char* test_dheap() {

	dheap_t *pdheap = dheap_alloc();
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 0);

	dheap_add(pdheap, 4.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 1);

	dheap_add(pdheap, 3.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 2);

	dheap_add(pdheap, 2.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 3);

	dheap_add(pdheap, 6.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 4);

	dheap_add(pdheap, 5.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 5);

	dheap_add(pdheap, 8.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 6);

	dheap_add(pdheap, 7.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 7);

	dheap_print(pdheap);

	mu_assert_lf(dheap_remove(pdheap) == 8.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 6);

	mu_assert_lf(dheap_remove(pdheap) == 7.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 5);

	mu_assert_lf(dheap_remove(pdheap) == 6.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 4);

	mu_assert_lf(dheap_remove(pdheap) == 5.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 3);

	mu_assert_lf(dheap_remove(pdheap) == 4.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 2);

	mu_assert_lf(dheap_remove(pdheap) == 3.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 1);

	mu_assert_lf(dheap_remove(pdheap) == 2.25);
	mu_assert_lf(dheap_check(pdheap, __FILE__,  __LINE__));
	mu_assert_lf(pdheap->n == 0);

	dheap_free(pdheap);

	return NULL;
}

// ================================================================
static char * run_all_tests() {
	mu_run_test(test_slls);
	mu_run_test(test_sllv);
	mu_run_test(test_hss);
	mu_run_test(test_lhmsi);
	mu_run_test(test_lhmss);
	mu_run_test(test_lhmsv);
	mu_run_test(test_lhms2v);
	mu_run_test(test_lhmslv);
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
