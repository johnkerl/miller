// ================================================================
// Functions involving more than one container type
// ================================================================

#ifndef MIXUTIL_H
#define MIXUTIL_H
#include "containers/lrec.h"
#include "containers/slls.h"
#include "containers/hss.h"
#include "containers/lhmss.h"
#include "lib/string_array.h"
#include "lib/mlrregex.h"

// Makes a list with values pointing to the lrec's keys. slls_free() will respect that and not corrupt the lrec.
// However, the slls values will be invalid after the lrec is freed.
slls_t* mlr_reference_keys_from_record(lrec_t* prec);
// Makes a list with values pointing to the lrec's values. slls_free() will respect that and not corrupt the lrec.
// However, the slls values will be invalid after the lrec is freed.
slls_t* mlr_reference_values_from_record(lrec_t* prec);

slls_t* mlr_reference_keys_from_record_except(lrec_t* prec, lrece_t* px);
slls_t* mlr_reference_values_from_record_except(lrec_t* prec, lrece_t* px);

// Copies data; no referencing concerns.
slls_t* mlr_copy_keys_from_record(lrec_t* prec);

// Makes a list with values pointing into the lrec's values. slls_free() will
// respect that and not corrupt the lrec. However, the slls values will be
// invalid after the lrec is freed.
slls_t* mlr_reference_selected_values_from_record(lrec_t* prec, slls_t* pselected_field_names);
void mlr_reference_values_from_record_into_string_array(lrec_t* prec, string_array_t* pselected_field_names,
	string_array_t* pvalues);
int record_has_all_keys(lrec_t* prec, slls_t* pselected_field_names);

lhmss_t* mlr_reference_key_value_pairs_from_regex_names(lrec_t* prec, regex_t* pregexes, int num_regexes,
	int invert_matches);

// Copies data; no referencing concerns.
hss_t*  hss_from_slls(slls_t* plist);

// Prints a list of lrecs using lrec_print.
void lrec_print_list(sllv_t* plist);
void lrec_print_list_with_prefix(sllv_t* plist, char* prefix);

// Same as
//   slls_t* prec_values = mlr_reference_selected_values_from_record(prec, pkeys);
//   return slls_compare_lexically(plist, prec_values);
// but without the unnecessary copy.
int slls_lrec_compare_lexically(
	slls_t* plist,
	lrec_t* prec,
	slls_t* pkeys);
int lrec_slls_compare_lexically(
	lrec_t* prec,
	slls_t* pkeys,
	slls_t* plist);

int lrec_keys_equal_list(
	lrec_t* prec,
	slls_t* plist);

#endif // MIXUTIL_H
