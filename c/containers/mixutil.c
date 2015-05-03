#include "lib/mlrutil.h"
#include "containers/mixutil.h"

// ----------------------------------------------------------------
// xxx freeing contract
// xxx return an lrec?
slls_t* mlr_keys_from_record(lrec_t* prec) {
	slls_t* plist = slls_alloc();
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		slls_add_no_free(plist, pe->key);
	}
	return plist;
}

// ----------------------------------------------------------------
// xxx freeing contract.
// xxx behavior on missing. doc, or make a second boolean flag.
slls_t* mlr_selected_values_from_record(lrec_t* prec, slls_t* pselected_field_names) {
	slls_t* pvalue_list = slls_alloc();
	for (sllse_t* pe = pselected_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* selected_field_name = pe->value;
		char* value = lrec_get(prec, selected_field_name);
		if (value == NULL) {
			// xxx have stashed argv0 for error message.
			// xxx better to have filename + linenumber somehow.
			//fprintf(stderr, "Couldn't find field named \"%s\"\n", selected_field_name);
			//exit(1);
		} else {
			slls_add_no_free(pvalue_list, value);
		}
	}
	return pvalue_list;
}
