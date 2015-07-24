// Functions involving more than one container type
#ifndef MIXUTIL_H
#define MIXUTIL_H
#include "containers/lrec.h"
#include "containers/slls.h"
#include "containers/hss.h"
slls_t* mlr_keys_from_record(lrec_t* prec);
slls_t* mlr_selected_values_from_record(lrec_t* prec, slls_t* pselected_field_names);
hss_t*  hss_from_slls(slls_t* plist);
void lrec_print_list(sllv_t* plist);
void lrec_print_list_with_prefix(sllv_t* plist, char* prefix);

// Same as
//   slls_t* prec_values = mlr_selected_values_from_record(prec, pkeys);
//   return slls_compare_lexically(plist, prec_values);
// but without the unnecessary copy.
// xxx context here ...
int slls_lrec_compare_lexically(
	slls_t* plist,
	lrec_t* prec,
	slls_t* pkeys);

#endif // MIXUTIL_H
