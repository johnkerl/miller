#include <stdio.h>
#include <string.h>
#include "lib/minunit.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/rslls.h"
#include "containers/sllv.h"
#include "lib/string_array.h"
#include "containers/hss.h"
#include "containers/lhmsi.h"
#include "containers/lhmsll.h"
#include "containers/lhmss.h"
#include "containers/lhmsv.h"
#include "containers/lhms2v.h"
#include "containers/lhmslv.h"
#include "containers/lhmsmv.h"
#include "containers/percentile_keeper.h"
#include "containers/top_keeper.h"
#include "containers/dheap.h"
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
static char* test_slls() {
	slls_t* plist = slls_from_line(mlr_strdup_or_die(""), ',', FALSE);
	mu_assert_lf(plist->length == 0);

	plist = slls_from_line(mlr_strdup_or_die("a"), ',', FALSE);

	mu_assert_lf(plist->length == 1);
	plist = slls_from_line(mlr_strdup_or_die("c,d,a,e,b"), ',', FALSE);
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


	plist = slls_from_line(mlr_strdup_or_die(","), ',', FALSE);
	slls_print_quoted(plist);printf("\n");
	mu_assert_lf(plist->length == 2);
	pe = plist->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "")); pe = pe->pnext;

	plist = slls_from_line(mlr_strdup_or_die("a,b,c,"), ',', FALSE);
	slls_print_quoted(plist);printf("\n");
	mu_assert_lf(plist->length == 4);
	pe = plist->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, ""));  pe = pe->pnext;

	plist = slls_from_line(mlr_strdup_or_die("a,,c,d"), ',', FALSE);
	slls_print_quoted(plist);printf("\n");
	mu_assert_lf(plist->length == 4);
	pe = plist->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, ""));  pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;

	return NULL;
}

// ----------------------------------------------------------------
static char* test_rslls() {

	rslls_t* pa = rslls_alloc();
	rslls_append(pa, "a", NO_FREE, 0);
	rslls_append(pa, "b", NO_FREE, 0);
	rslls_append(pa, "c", NO_FREE, 0);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 3);
	rsllse_t* pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "c")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_reset(pa);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 0);
	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_append(pa, "d", NO_FREE, 0);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 1);
	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_append(pa, "e", NO_FREE, 0);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 2);
	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "e")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(pe->value == NULL); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_append(pa, "f", NO_FREE, 0);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 3);
	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "e")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "f")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_append(pa, "g", NO_FREE, 0);

	rslls_print(pa); printf("\n");
	mu_assert_lf(pa->length == 4);
	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "e")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "f")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->value, "g")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	rslls_free(pa);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_sllv() {

	sllv_t* pa = sllv_alloc();
	sllv_append(pa, "a");
	sllv_append(pa, "b");
	sllv_append(pa, "c");
	mu_assert_lf(pa->length == 3);

	sllve_t* pe = pa->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "c")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	sllv_t* pb = sllv_alloc();
	sllv_append(pb, "d");
	sllv_append(pb, "e");
	mu_assert_lf(pb->length == 2);

	pe = pb->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	sllv_transfer(pa, pb);

	mu_assert_lf(pa->length == 5);
	mu_assert_lf(pb->length == 0);
	mu_assert_lf(pb->phead == NULL);
	mu_assert_lf(pb->ptail == NULL);

	pe = pa->phead;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "a")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "d")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "e")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	sllv_free(pb);

	pa = sllv_alloc();
	sllv_push(pa, "a");
	sllv_push(pa, "b");
	sllv_push(pa, "c");
	mu_assert_lf(pa->length == 3);

	pe = pa->phead;

	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "c")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "b")); pe = pe->pnext;
	mu_assert_lf(pe != NULL); mu_assert_lf(streq(pe->pvvalue, "a")); pe = pe->pnext;
	mu_assert_lf(pe == NULL);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_string_array() {
	string_array_t* parray = string_array_from_line(mlr_strdup_or_die(""), ',');
	mu_assert_lf(parray->length == 0);
	string_array_free(parray);

	parray = string_array_from_line(mlr_strdup_or_die("a"), ',');
	mu_assert_lf(parray->length == 1);
	mu_assert_lf(streq(parray->strings[0], "a"));
	string_array_free(parray);

	parray = string_array_from_line(mlr_strdup_or_die("c,d,a,e,b"), ',');
	mu_assert_lf(parray->length == 5);
	mu_assert_lf(streq(parray->strings[0], "c"));
	mu_assert_lf(streq(parray->strings[1], "d"));
	mu_assert_lf(streq(parray->strings[2], "a"));
	mu_assert_lf(streq(parray->strings[3], "e"));
	mu_assert_lf(streq(parray->strings[4], "b"));
	string_array_free(parray);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_hss() {

	hss_t *pset = hss_alloc();
	mu_assert_lf(pset->num_occupied == 0);

	hss_add(pset, "x");
	mu_assert_lf(pset->num_occupied == 1);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf( hss_has(pset, "x"));
	mu_assert_lf(!hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "y");
	mu_assert_lf(pset->num_occupied == 2);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf( hss_has(pset, "x"));
	mu_assert_lf( hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "x");
	mu_assert_lf(pset->num_occupied == 2);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf( hss_has(pset, "x"));
	mu_assert_lf( hss_has(pset, "y"));
	mu_assert_lf(!hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_add(pset, "z");
	mu_assert_lf(pset->num_occupied == 3);
	mu_assert_lf(!hss_has(pset, "w"));
	mu_assert_lf( hss_has(pset, "x"));
	mu_assert_lf( hss_has(pset, "y"));
	mu_assert_lf(hss_has(pset, "z"));
	mu_assert_lf(hss_check_counts(pset));

	hss_free(pset);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsi() {

	lhmsi_t *pmap = lhmsi_alloc();
	int val = -123;
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf(lhmsi_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "x", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "y", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 3, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf( lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "x") == 3);
	mu_assert_lf(!lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == -999);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 3);
	mu_assert_lf(lhmsi_test_and_get(pmap, "y", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "y", 5, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf( lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "x") == 3);
	mu_assert_lf( lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 3);
	mu_assert_lf(lhmsi_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsi_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "x", 4, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf( lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "x") == 4);
	mu_assert_lf( lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsi_has_key(pmap, "z")); mu_assert_lf(lhmsi_get(pmap, "z") == -999);
	mu_assert_lf(lhmsi_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 4);
	mu_assert_lf(lhmsi_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsi_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_put(pmap, "z", 7, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsi_has_key(pmap, "w")); mu_assert_lf(lhmsi_get(pmap, "w") == -999);
	mu_assert_lf( lhmsi_has_key(pmap, "x")); mu_assert_lf(lhmsi_get(pmap, "x") == 4);
	mu_assert_lf( lhmsi_has_key(pmap, "y")); mu_assert_lf(lhmsi_get(pmap, "y") == 5);
	mu_assert_lf(lhmsi_has_key(pmap, "z"));  mu_assert_lf(lhmsi_get(pmap, "z") == 7);
	mu_assert_lf(lhmsi_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsi_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 4);
	mu_assert_lf(lhmsi_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsi_test_and_get(pmap, "z", &val) == TRUE); mu_assert_lf(val == 7);
	mu_assert_lf(lhmsi_check_counts(pmap));

	lhmsi_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsll() {

	lhmsll_t *pmap = lhmsll_alloc();
	long long val = -123;
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsll_has_key(pmap, "w")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsll_has_key(pmap, "x")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsll_has_key(pmap, "y")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf(!lhmsll_has_key(pmap, "z")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf(lhmsll_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "x", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "y", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsll_check_counts(pmap));

	lhmsll_put(pmap, "x", 3, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsll_has_key(pmap, "w")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf( lhmsll_has_key(pmap, "x")); mu_assert_lf(lhmsll_get(pmap, "x") == 3);
	mu_assert_lf(!lhmsll_has_key(pmap, "y")); mu_assert_lf(lhmsll_get(pmap, "y") == -999);
	mu_assert_lf(!lhmsll_has_key(pmap, "z")); mu_assert_lf(lhmsll_get(pmap, "z") == -999);
	mu_assert_lf(lhmsll_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 3);
	mu_assert_lf(lhmsll_test_and_get(pmap, "y", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsll_check_counts(pmap));

	lhmsll_put(pmap, "y", 5, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsll_has_key(pmap, "w")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf( lhmsll_has_key(pmap, "x")); mu_assert_lf(lhmsll_get(pmap, "x") == 3);
	mu_assert_lf( lhmsll_has_key(pmap, "y")); mu_assert_lf(lhmsll_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsll_has_key(pmap, "z")); mu_assert_lf(lhmsll_get(pmap, "z") == -999);
	mu_assert_lf(lhmsll_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 3);
	mu_assert_lf(lhmsll_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsll_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsll_check_counts(pmap));

	lhmsll_put(pmap, "x", 4, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsll_has_key(pmap, "w")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf( lhmsll_has_key(pmap, "x")); mu_assert_lf(lhmsll_get(pmap, "x") == 4);
	mu_assert_lf( lhmsll_has_key(pmap, "y")); mu_assert_lf(lhmsll_get(pmap, "y") == 5);
	mu_assert_lf(!lhmsll_has_key(pmap, "z")); mu_assert_lf(lhmsll_get(pmap, "z") == -999);
	mu_assert_lf(lhmsll_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 4);
	mu_assert_lf(lhmsll_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsll_test_and_get(pmap, "z", &val) == FALSE);
	mu_assert_lf(lhmsll_check_counts(pmap));

	lhmsll_put(pmap, "z", 7, NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsll_has_key(pmap, "w")); mu_assert_lf(lhmsll_get(pmap, "w") == -999);
	mu_assert_lf( lhmsll_has_key(pmap, "x")); mu_assert_lf(lhmsll_get(pmap, "x") == 4);
	mu_assert_lf( lhmsll_has_key(pmap, "y")); mu_assert_lf(lhmsll_get(pmap, "y") == 5);
	mu_assert_lf(lhmsll_has_key(pmap, "z"));  mu_assert_lf(lhmsll_get(pmap, "z") == 7);
	mu_assert_lf(lhmsll_test_and_get(pmap, "w", &val) == FALSE);
	mu_assert_lf(lhmsll_test_and_get(pmap, "x", &val) == TRUE); mu_assert_lf(val == 4);
	mu_assert_lf(lhmsll_test_and_get(pmap, "y", &val) == TRUE); mu_assert_lf(val == 5);
	mu_assert_lf(lhmsll_test_and_get(pmap, "z", &val) == TRUE); mu_assert_lf(val == 7);
	mu_assert_lf(lhmsll_check_counts(pmap));

	lhmsll_free(pmap);

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

	lhmss_put(pmap, "x", "3", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf( lhmss_has_key(pmap, "x")); mu_assert_lf(streq(lhmss_get(pmap, "x"), "3"));
	mu_assert_lf(!lhmss_has_key(pmap, "y")); mu_assert_lf(lhmss_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "y", "5", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf( lhmss_has_key(pmap, "x")); mu_assert_lf(streq(lhmss_get(pmap, "x"), "3"));
	mu_assert_lf( lhmss_has_key(pmap, "y")); mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "x", "4", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf( lhmss_has_key(pmap, "x")); mu_assert_lf(streq(lhmss_get(pmap, "x"), "4"));
	mu_assert_lf( lhmss_has_key(pmap, "y")); mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmss_has_key(pmap, "z")); mu_assert_lf(lhmss_get(pmap, "z") == NULL);
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_put(pmap, "z", "7", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmss_has_key(pmap, "w")); mu_assert_lf(lhmss_get(pmap, "w") == NULL);
	mu_assert_lf( lhmss_has_key(pmap, "x")); mu_assert_lf(streq(lhmss_get(pmap, "x"), "4"));
	mu_assert_lf( lhmss_has_key(pmap, "y")); mu_assert_lf(streq(lhmss_get(pmap, "y"), "5"));
	mu_assert_lf( lhmss_has_key(pmap, "z")); mu_assert_lf(streq(lhmss_get(pmap, "z"), "7"));
	mu_assert_lf(lhmss_check_counts(pmap));

	lhmss_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsv() {

	lhmsv_t *pmap = lhmsv_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsv_has_key(pmap, "w")); mu_assert_lf(lhmsv_get(pmap, "w") == NULL);
	mu_assert_lf(!lhmsv_has_key(pmap, "x")); mu_assert_lf(lhmsv_get(pmap, "x") == NULL);
	mu_assert_lf(!lhmsv_has_key(pmap, "y")); mu_assert_lf(lhmsv_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmsv_has_key(pmap, "z")); mu_assert_lf(lhmsv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsv_check_counts(pmap));

	lhmsv_put(pmap, "x", "3", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsv_has_key(pmap, "w")); mu_assert_lf(lhmsv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsv_has_key(pmap, "x")); mu_assert_lf(streq(lhmsv_get(pmap, "x"), "3"));
	mu_assert_lf(!lhmsv_has_key(pmap, "y")); mu_assert_lf(lhmsv_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmsv_has_key(pmap, "z")); mu_assert_lf(lhmsv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsv_check_counts(pmap));

	lhmsv_put(pmap, "y", "5", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsv_has_key(pmap, "w")); mu_assert_lf(lhmsv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsv_has_key(pmap, "x")); mu_assert_lf(streq(lhmsv_get(pmap, "x"), "3"));
	mu_assert_lf( lhmsv_has_key(pmap, "y")); mu_assert_lf(streq(lhmsv_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmsv_has_key(pmap, "z")); mu_assert_lf(lhmsv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsv_check_counts(pmap));

	lhmsv_put(pmap, "x", "4", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsv_has_key(pmap, "w")); mu_assert_lf(lhmsv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsv_has_key(pmap, "x")); mu_assert_lf(streq(lhmsv_get(pmap, "x"), "4"));
	mu_assert_lf( lhmsv_has_key(pmap, "y")); mu_assert_lf(streq(lhmsv_get(pmap, "y"), "5"));
	mu_assert_lf(!lhmsv_has_key(pmap, "z")); mu_assert_lf(lhmsv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsv_check_counts(pmap));

	lhmsv_put(pmap, "z", "7", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsv_has_key(pmap, "w")); mu_assert_lf(lhmsv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsv_has_key(pmap, "x")); mu_assert_lf(streq(lhmsv_get(pmap, "x"), "4"));
	mu_assert_lf( lhmsv_has_key(pmap, "y")); mu_assert_lf(streq(lhmsv_get(pmap, "y"), "5"));
	mu_assert_lf( lhmsv_has_key(pmap, "z")); mu_assert_lf(streq(lhmsv_get(pmap, "z"), "7"));
	mu_assert_lf(lhmsv_check_counts(pmap));

	lhmsv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhms2v() {

	lhms2v_t *pmap = lhms2v_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(lhms2v_get(pmap, "a", "x") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "3", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf( lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "3"));
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(lhms2v_get(pmap, "a", "y") == NULL);
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "y", "5", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf( lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "3"));
	mu_assert_lf( lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "a", "x", "4", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf( lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "4"));
	mu_assert_lf( lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf(!lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(lhms2v_get(pmap, "b", "z") == NULL);
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_put(pmap, "b", "z", "7", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhms2v_has_key(pmap, "a", "w")); mu_assert_lf(lhms2v_get(pmap, "a", "w") == NULL);
	mu_assert_lf( lhms2v_has_key(pmap, "a", "x")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "x"), "4"));
	mu_assert_lf( lhms2v_has_key(pmap, "a", "y")); mu_assert_lf(streq(lhms2v_get(pmap, "a", "y"), "5"));
	mu_assert_lf( lhms2v_has_key(pmap, "b", "z")); mu_assert_lf(streq(lhms2v_get(pmap, "b", "z"), "7"));
	mu_assert_lf(lhms2v_check_counts(pmap));

	lhms2v_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmslv() {

	slls_t* aw = slls_alloc(); slls_append_no_free(aw, "a"); slls_append_no_free(aw, "w");
	slls_t* ax = slls_alloc(); slls_append_no_free(ax, "a"); slls_append_no_free(ax, "x");
	slls_t* ay = slls_alloc(); slls_append_no_free(ay, "a"); slls_append_no_free(ay, "y");
	slls_t* bz = slls_alloc(); slls_append_no_free(bz, "b"); slls_append_no_free(bz, "z");

	lhmslv_t *pmap = lhmslv_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ax)); mu_assert_lf(lhmslv_get(pmap, ax) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ax, "3", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf( lhmslv_has_key(pmap, ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "3"));
	mu_assert_lf(!lhmslv_has_key(pmap, ay)); mu_assert_lf(lhmslv_get(pmap, ay) == NULL);
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ay, "5", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf( lhmslv_has_key(pmap, ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "3"));
	mu_assert_lf( lhmslv_has_key(pmap, ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, ax, "4", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf( lhmslv_has_key(pmap, ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "4"));
	mu_assert_lf( lhmslv_has_key(pmap, ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf(!lhmslv_has_key(pmap, bz)); mu_assert_lf(lhmslv_get(pmap, bz) == NULL);
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_put(pmap, bz, "7", NO_FREE);
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmslv_has_key(pmap, aw)); mu_assert_lf(lhmslv_get(pmap, aw) == NULL);
	mu_assert_lf( lhmslv_has_key(pmap, ax)); mu_assert_lf(streq(lhmslv_get(pmap, ax), "4"));
	mu_assert_lf( lhmslv_has_key(pmap, ay)); mu_assert_lf(streq(lhmslv_get(pmap, ay), "5"));
	mu_assert_lf( lhmslv_has_key(pmap, bz)); mu_assert_lf(streq(lhmslv_get(pmap, bz), "7"));
	mu_assert_lf(lhmslv_check_counts(pmap));

	lhmslv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_lhmsmv() {
	printf("\n");

	lhmsmv_t *pmap = lhmsmv_alloc();
	mu_assert_lf(pmap->num_occupied == 0);
	mu_assert_lf(!lhmsmv_has_key(pmap, "w")); mu_assert_lf(lhmsmv_get(pmap, "w") == NULL);
	mu_assert_lf(!lhmsmv_has_key(pmap, "x")); mu_assert_lf(lhmsmv_get(pmap, "x") == NULL);
	mu_assert_lf(!lhmsmv_has_key(pmap, "y")); mu_assert_lf(lhmsmv_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmsmv_has_key(pmap, "z")); mu_assert_lf(lhmsmv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsmv_check_counts(pmap));

	lhmsmv_put(pmap, "x", imv(3), NO_FREE);
	lhmsmv_dump(pmap);
	printf("\n");
	mu_assert_lf(pmap->num_occupied == 1);
	mu_assert_lf(!lhmsmv_has_key(pmap, "w")); mu_assert_lf(lhmsmv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsmv_has_key(pmap, "x")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "x"), imv(3)));
	mu_assert_lf( lhmsmv_has_key(pmap, "x")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "x"), smv("3")));
	mu_assert_lf(!lhmsmv_has_key(pmap, "y")); mu_assert_lf(lhmsmv_get(pmap, "y") == NULL);
	mu_assert_lf(!lhmsmv_has_key(pmap, "z")); mu_assert_lf(lhmsmv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsmv_check_counts(pmap));

	lhmsmv_put(pmap, "y", smv("5"), NO_FREE);
	lhmsmv_dump(pmap);
	printf("\n");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsmv_has_key(pmap, "w")); mu_assert_lf(lhmsmv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsmv_has_key(pmap, "x")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "x"), smv("3")));
	mu_assert_lf( lhmsmv_has_key(pmap, "y")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "y"), smv("5")));
	mu_assert_lf( lhmsmv_has_key(pmap, "y")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "y"), imv(5)));
	mu_assert_lf(!lhmsmv_has_key(pmap, "z")); mu_assert_lf(lhmsmv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsmv_check_counts(pmap));

	lhmsmv_put(pmap, "x", smv("f"), NO_FREE);
	lhmsmv_dump(pmap);
	printf("\n");
	mu_assert_lf(pmap->num_occupied == 2);
	mu_assert_lf(!lhmsmv_has_key(pmap, "w")); mu_assert_lf(lhmsmv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsmv_has_key(pmap, "x")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "x"), smv("f")));
	mu_assert_lf( lhmsmv_has_key(pmap, "y")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "y"), smv("5")));
	mu_assert_lf(!lhmsmv_has_key(pmap, "z")); mu_assert_lf(lhmsmv_get(pmap, "z") == NULL);
	mu_assert_lf(lhmsmv_check_counts(pmap));

	lhmsmv_put(pmap, "z", imv(7), NO_FREE);
	lhmsmv_dump(pmap);
	printf("\n");
	mu_assert_lf(pmap->num_occupied == 3);
	mu_assert_lf(!lhmsmv_has_key(pmap, "w")); mu_assert_lf(lhmsmv_get(pmap, "w") == NULL);
	mu_assert_lf( lhmsmv_has_key(pmap, "x")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "x"), smv("f")));
	mu_assert_lf( lhmsmv_has_key(pmap, "y")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "y"), smv("5")));
	mu_assert_lf( lhmsmv_has_key(pmap, "y")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "y"), imv(5)));
	mu_assert_lf( lhmsmv_has_key(pmap, "z")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "z"), smv("7")));
	mu_assert_lf( lhmsmv_has_key(pmap, "z")); mu_assert_lf(mveqcopy(lhmsmv_get(pmap, "z"), imv(7)));
	mu_assert_lf(lhmsmv_check_counts(pmap));

	lhmsmv_free(pmap);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_percentile_keeper() {

	percentile_keeper_t* ppercentile_keeper = percentile_keeper_alloc();
	percentile_keeper_ingest(ppercentile_keeper, mv_from_float(1.0));
	percentile_keeper_ingest(ppercentile_keeper, mv_from_int(2));
	percentile_keeper_ingest(ppercentile_keeper, mv_from_int(3));
	percentile_keeper_ingest(ppercentile_keeper, mv_from_float(4.0));
	percentile_keeper_ingest(ppercentile_keeper, mv_from_float(5.0));
	percentile_keeper_print(ppercentile_keeper);

	double p;
	mv_t q;
	p = 0.0;
	q = percentile_keeper_emit_non_interpolated(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q.u.fltv);
	mu_assert_lf(q.type == MT_FLOAT);
	mu_assert_lf(q.u.fltv == 1.0);

	p = 10.0;
	q = percentile_keeper_emit_non_interpolated(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q.u.fltv);
	mu_assert_lf(q.type == MT_FLOAT);
	mu_assert_lf(q.u.fltv == 1.0);

	p = 50.0;
	q = percentile_keeper_emit_non_interpolated(ppercentile_keeper, p);
	printf("%4.2lf -> %lld\n", p, q.u.intv);
	mu_assert_lf(q.type == MT_INT);
	mu_assert_lf(q.u.intv == 3LL);

	p = 90.0;
	q = percentile_keeper_emit_non_interpolated(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q.u.fltv);
	mu_assert_lf(q.type == MT_FLOAT);
	mu_assert_lf(q.u.fltv == 5.0);

	p = 100.0;
	q = percentile_keeper_emit_non_interpolated(ppercentile_keeper, p);
	printf("%4.2lf -> %7.4lf\n", p, q.u.fltv);
	mu_assert_lf(q.type == MT_FLOAT);
	mu_assert_lf(q.u.fltv == 5.0);

	percentile_keeper_free(ppercentile_keeper);

	return NULL;
}

// ----------------------------------------------------------------
static char* test_top_keeper() {
	int capacity = 3;

	top_keeper_t* ptop_keeper = top_keeper_alloc(capacity);
	mu_assert_lf(ptop_keeper->size == 0);

	top_keeper_add(ptop_keeper, mv_from_float(5.0), NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 1);
	mu_assert_lf(ptop_keeper->top_values[0].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[0].u.fltv == 5.0);

	top_keeper_add(ptop_keeper, mv_from_float(6.0), NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 2);
	mu_assert_lf(ptop_keeper->top_values[0].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[0].u.fltv == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[1].u.fltv == 5.0);

	top_keeper_add(ptop_keeper, mv_from_int(4), NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[0].u.fltv == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[1].u.fltv == 5.0);
	mu_assert_lf(ptop_keeper->top_values[2].type == MT_INT);
	mu_assert_lf(ptop_keeper->top_values[2].u.intv == 4.0);

	top_keeper_add(ptop_keeper, mv_from_int(2), NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[0].u.fltv == 6.0);
	mu_assert_lf(ptop_keeper->top_values[1].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[1].u.fltv == 5.0);
	mu_assert_lf(ptop_keeper->top_values[2].type == MT_INT);
	mu_assert_lf(ptop_keeper->top_values[2].u.intv == 4.0);

	top_keeper_add(ptop_keeper, mv_from_int(7), NULL);
	top_keeper_print(ptop_keeper);
	mu_assert_lf(ptop_keeper->size == 3);
	mu_assert_lf(ptop_keeper->top_values[0].type == MT_INT);
	mu_assert_lf(ptop_keeper->top_values[0].u.intv == 7);
	mu_assert_lf(ptop_keeper->top_values[1].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[1].u.fltv == 6.0);
	mu_assert_lf(ptop_keeper->top_values[2].type == MT_FLOAT);
	mu_assert_lf(ptop_keeper->top_values[2].u.fltv == 5.0);

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
	mu_run_test(test_rslls);
	mu_run_test(test_sllv);
	mu_run_test(test_string_array);
	mu_run_test(test_hss);
	mu_run_test(test_lhmsi);
	mu_run_test(test_lhmsll);
	mu_run_test(test_lhmss);
	mu_run_test(test_lhmsv);
	mu_run_test(test_lhms2v);
	mu_run_test(test_lhmslv);
	mu_run_test(test_lhmsmv);
	mu_run_test(test_percentile_keeper);
	mu_run_test(test_top_keeper);
	mu_run_test(test_dheap);
	return 0;
}

int main(int argc, char **argv) {
	mlr_global_init(argv[0], NULL);

	printf("TEST_MULTIPLE_CONTAINERS ENTER\n");
	char *result = run_all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_MULTIPLE_CONTAINERS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
