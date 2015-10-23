#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/mixutil.h"

// ----------------------------------------------------------------
// Makes a list with values pointing to the lrec's keys. slls_free() will
// respect that and not corrupt the lrec. However, the slls values will be
// invalid after the lrec is freed.

slls_t* mlr_reference_keys_from_record(lrec_t* prec) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		slls_add_no_free(plist, pe->key);
	}
	return plist;
}

slls_t* mlr_copy_keys_from_record(lrec_t* prec) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		slls_add_with_free(plist, mlr_strdup_or_die(pe->key));
	}
	return plist;
}

// ----------------------------------------------------------------
// Makes a list with values pointing into the lrec's values. slls_free() will
// respect that and not corrupt the lrec. However, the slls values will be
// invalid after the lrec is freed.

slls_t* mlr_selected_values_from_record(lrec_t* prec, slls_t* pselected_field_names) {
	slls_t* pvalue_list = slls_alloc();
	for (sllse_t* pe = pselected_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* selected_field_name = pe->value;
		char* value = lrec_get(prec, selected_field_name);
		if (value == NULL) {
			slls_free(pvalue_list);
			return NULL;
		} else {
			slls_add_no_free(pvalue_list, value);
		}
	}
	return pvalue_list;
}

// ----------------------------------------------------------------
hss_t* hss_from_slls(slls_t* plist) {
	hss_t* pset = hss_alloc();
	for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext)
		hss_add(pset, pe->value);
	return pset;
}

// ----------------------------------------------------------------
void lrec_print_list(sllv_t* plist) {
	for (sllve_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
		lrec_print(pe->pvdata);
	}
}

void lrec_print_list_with_prefix(sllv_t* plist, char* prefix) {
	for (sllve_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
		printf("%s", prefix);
		lrec_print(pe->pvdata);
	}
}

// ----------------------------------------------------------------
int slls_lrec_compare_lexically(
	slls_t* plist,
	lrec_t* prec,
	slls_t* pkeys)
{
	sllse_t* pe = plist->phead;
	sllse_t* pf = pkeys->phead;
	while (TRUE) {
		if (pe == NULL && pf == NULL)
			return 0;
		if (pe == NULL)
			return 1;
		if (pf == NULL)
			return -1;

		char* precval = lrec_get(prec, pf->value);
		if (precval == NULL) {
			return -1;
		} else {
			int rc = strcmp(pe->value, precval);
			if (rc != 0)
				return rc;
		}

		pe = pe->pnext;
		pf = pf->pnext;
	}
}

// ----------------------------------------------------------------
int lrec_slls_compare_lexically(
	lrec_t* prec,
	slls_t* pkeys,
	slls_t* plist)
{
	return -slls_lrec_compare_lexically(plist, prec, pkeys);
}

// ----------------------------------------------------------------
int lrec_keys_equal_list(
	lrec_t* prec,
	slls_t* plist)
{
	lrece_t* pe = prec->phead;
	sllse_t* pf = plist->phead;
	while (TRUE) {
		if (pe == NULL && pf == NULL)
			return TRUE;
		if (pe == NULL || pf == NULL)
			return FALSE;
		if (!streq(pe->key, pf->value))
			return FALSE;
		pe = pe->pnext;
		pf = pf->pnext;
	}
}
