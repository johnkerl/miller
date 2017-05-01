#include "lib/mlr_globals.h"
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
		slls_append_no_free(plist, pe->key);
	}
	return plist;
}

slls_t* mlr_copy_keys_from_record(lrec_t* prec) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		slls_append_with_free(plist, mlr_strdup_or_die(pe->key));
	}
	return plist;
}

slls_t* mlr_reference_values_from_record(lrec_t* prec) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		slls_append_no_free(plist, pe->value);
	}
	return plist;
}

slls_t* mlr_reference_keys_from_record_except(lrec_t* prec, lrece_t* px) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (pe != px)
			slls_append_no_free(plist, pe->key);
	}
	return plist;
}

slls_t* mlr_reference_values_from_record_except(lrec_t* prec, lrece_t* px) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (pe != px)
			slls_append_no_free(plist, pe->value);
	}
	return plist;
}

// ----------------------------------------------------------------
// Makes a list with values pointing into the lrec's values. slls_free() will
// respect that and not corrupt the lrec. However, the slls values will be
// invalid after the lrec is freed.

slls_t* mlr_reference_selected_values_from_record(lrec_t* prec, slls_t* pselected_field_names) {
	slls_t* pvalue_list = slls_alloc();
	for (sllse_t* pe = pselected_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* selected_field_name = pe->value;
		char* value = lrec_get(prec, selected_field_name);
		if (value == NULL) {
			slls_free(pvalue_list);
			return NULL;
		} else {
			slls_append_no_free(pvalue_list, value);
		}
	}
	return pvalue_list;
}

// Makes an array with values pointing into the lrec's values.
// string_array_free() will respect that and not corrupt the lrec. However,
// the array's values will be invalid after the lrec is freed.

void mlr_reference_values_from_record_into_string_array(lrec_t* prec, string_array_t* pselected_field_names,
	string_array_t* pvalues)
{
	MLR_INTERNAL_CODING_ERROR_IF(pselected_field_names->length != pvalues->length);
	pvalues->strings_need_freeing = FALSE;
	for (int i = 0; i < pselected_field_names->length; i++) {
		char* selected_field_name = pselected_field_names->strings[i];
		if (selected_field_name == NULL) {
			pvalues->strings[i] = NULL;
		} else {
			pvalues->strings[i] = lrec_get(prec, selected_field_name);
		}

	}
}

int record_has_all_keys(lrec_t* prec, slls_t* pselected_field_names) {
	for (sllse_t* pe = pselected_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* selected_field_name = pe->value;
		char* value = lrec_get(prec, selected_field_name);
		if (value == NULL)
			return FALSE;
	}
	return TRUE;
}

// ----------------------------------------------------------------
lhmss_t* mlr_reference_key_value_pairs_from_regex_names(lrec_t* prec, regex_t* pregexes, int num_regexes,
	int invert_matches)
{
	lhmss_t* pmap = lhmss_alloc();

	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		int matches_any = FALSE;
		for (int i = 0; i < num_regexes; i++) {
			regex_t* pregex = &pregexes[i];
			if (regmatch_or_die(pregex, pe->key, 0, NULL)) {
				matches_any = TRUE;
				break;
			}
		}
		if (matches_any) {
			lhmss_put(pmap, pe->key, pe->value, NO_FREE);
		}
	}

	return pmap;
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
		lrec_print(pe->pvvalue);
	}
}

void lrec_print_list_with_prefix(sllv_t* plist, char* prefix) {
	if (plist == NULL) {
		printf("%s NULL", prefix);
	} else {
		for (sllve_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
			printf("%s", prefix);
			lrec_print(pe->pvvalue);
		}
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
