// ================================================================
// Argument-parsing library, with non-getopt semantics.
// ================================================================

#ifndef ARGPARSE_H
#define ARGPARSE_H
#include "containers/slls.h"
#include "containers/sllv.h"

typedef struct _ap_state_t {
	sllv_t* pflag_defs;
} ap_state_t;

ap_state_t* ap_alloc();
void ap_free(ap_state_t* pstate);

void        ap_define_true_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void       ap_define_false_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void   ap_define_int_value_flag(ap_state_t* pstate, char* flag_name, int value, int* pintval);
void        ap_define_char_flag(ap_state_t* pstate, char* flag_name, char* pcharval);
void         ap_define_int_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void      ap_define_double_flag(ap_state_t* pstate, char* flag_name, double* pdoubleval);
void      ap_define_string_flag(ap_state_t* pstate, char* flag_name, char** pstring);
void ap_define_string_list_flag(ap_state_t* pstate, char* flag_name, slls_t** pplist);

int ap_parse(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv);

#endif // ARGPARSE_H
