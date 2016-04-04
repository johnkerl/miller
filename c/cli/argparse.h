// ================================================================
// Argument-parsing library, with non-getopt semantics.
// ================================================================

#ifndef ARGPARSE_H
#define ARGPARSE_H
#include "containers/slls.h"
#include "containers/sllv.h"
#include "lib/string_array.h"

typedef struct _ap_state_t {
	sllv_t* pflag_defs;
} ap_state_t;

ap_state_t* ap_alloc();
void ap_free(ap_state_t* pstate);

void         ap_define_true_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void        ap_define_false_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void    ap_define_int_value_flag(ap_state_t* pstate, char* flag_name, int value, int* pintval);
void          ap_define_int_flag(ap_state_t* pstate, char* flag_name, int* pintval);
void    ap_define_long_long_flag(ap_state_t* pstate, char* flag_name, long long* pintval);
void       ap_define_float_flag(ap_state_t* pstate, char* flag_name, double* pdoubleval);
void       ap_define_string_flag(ap_state_t* pstate, char* flag_name, char** pstring);
void  ap_define_string_list_flag(ap_state_t* pstate, char* flag_name, slls_t** pplist);
void ap_define_string_array_flag(ap_state_t* pstate, char* flag_name, string_array_t** pparray);

int ap_parse(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv);
int ap_parse_aux(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv,
	int error_on_unrecognized);

#endif // ARGPARSE_H
