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
#endif // MIXUTIL_H
