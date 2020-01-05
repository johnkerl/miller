#include <stdio.h>
#include "lib/mlrutil.h"
#include "cli/argparse.h"

// ================================================================
typedef enum _ap_flag_t {
	AP_INT_VALUE_FLAG,
	AP_INT_FLAG,
	AP_LONG_LONG_FLAG,
	AP_DOUBLE_FLAG,
	AP_STRING_FLAG,
	AP_STRING_BUILD_LIST_FLAG,
	AP_STRING_LIST_FLAG,
	AP_STRING_ARRAY_FLAG
} ap_flag_t;

typedef struct _ap_flag_def_t {
	char*     flag_name;
	ap_flag_t type;
	int       intval;
	void*     pval;
	int       count; // 1 for bool flags; 2 for the rest
} ap_flag_def_t;

static ap_flag_def_t* ap_find(ap_state_t* pstate, char* flag_name) {
	for (sllve_t* pe = pstate->pflag_defs->phead; pe != NULL; pe = pe->pnext) {
		ap_flag_def_t* pdef = pe->pvvalue;
		if (streq(pdef->flag_name, flag_name))
			return pdef;
	}
	return NULL;
}

static ap_flag_def_t* ap_flag_def_alloc(char* flag_name, ap_flag_t ap_type, int intval, void* pval, int count) {
	ap_flag_def_t* pdef = mlr_malloc_or_die(sizeof(ap_flag_def_t));
	pdef->flag_name = flag_name;
	pdef->type      = ap_type;
	pdef->intval    = intval;
	pdef->pval      = pval;
	pdef->count     = count;
	return pdef;
}

// ================================================================
ap_state_t* ap_alloc() {
	ap_state_t* pstate = mlr_malloc_or_die(sizeof(ap_state_t));
	pstate->pflag_defs = sllv_alloc();

	return pstate;
}

void ap_free(ap_state_t* pstate) {
	if (pstate == NULL)
		return;

	for (sllve_t* pe = pstate->pflag_defs->phead; pe != NULL; pe = pe->pnext) {
		ap_flag_def_t* pdef = pe->pvvalue;

		// Linked-lists are pointed to by mappers and freed by their free
		// methods.  If any mappers miss on that contract, we can find out by
		// using valgrind --leak-check=full (e.g. reg_test/run --valgrind).
		//
		//if (pdef->type == AP_STRING_LIST_FLAG && pdef->pval != NULL) {
		//	slls_t** pplist = pdef->pval;
		//	slls_free(*pplist);
		//}

		free(pdef);
	}
	sllv_free(pstate->pflag_defs);

	free(pstate);
}

// ----------------------------------------------------------------
void ap_define_true_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, TRUE, pintval, 1));
}

void ap_define_false_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, FALSE, pintval, 1));
}

void ap_define_int_value_flag(ap_state_t* pstate, char* flag_name, int intval, int* pintval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, intval, pintval, 1));
}

void ap_define_int_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_FLAG, 0, pintval, 2));
}
void ap_define_long_long_flag(ap_state_t* pstate, char* flag_name, long long* pintval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_LONG_LONG_FLAG, 0, pintval, 2));
}

void ap_define_float_flag(ap_state_t* pstate, char* flag_name, double* pdoubleval) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_DOUBLE_FLAG, 0, pdoubleval, 2));
}

void ap_define_string_flag(ap_state_t* pstate, char* flag_name, char** pstring) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_FLAG, 0, pstring, 2));
}

void ap_define_string_build_list_flag(ap_state_t* pstate, char* flag_name, slls_t** pplist) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_BUILD_LIST_FLAG, 0, pplist, 2));
}

void ap_define_string_list_flag(ap_state_t* pstate, char* flag_name, slls_t** pplist) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_LIST_FLAG, 0, pplist, 2));
}

void ap_define_string_array_flag(ap_state_t* pstate, char* flag_name, string_array_t** pparray) {
	sllv_append(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_ARRAY_FLAG, 0, pparray, 2));
}

// ----------------------------------------------------------------
int ap_parse(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv) {
	return ap_parse_aux(pstate, verb, pargi, argc, argv, TRUE);
}

int ap_parse_aux(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv,
	int error_on_unrecognized)
{
	int argi = *pargi;
	int ok = TRUE;

	while (argi < argc) {
		if (argv[argi][0] != '-' && argv[argi][0] != '+') {
			break;
		}
		if (streq(argv[argi], "-h") || streq(argv[argi], "--help")) {
			ok = FALSE;
			break;
		}

		ap_flag_def_t* pdef = ap_find(pstate, argv[argi]);
		if (pdef == NULL) {
			if (error_on_unrecognized)
				ok = FALSE;
			break;
		}

		if ((argc-argi) < pdef->count) {
			fprintf(stderr, "%s %s: option %s requires an argument.\n",
				argv[0], verb, argv[argi]);
			fprintf(stderr, "\n");
			ok = FALSE;
			break;
		}

		if (pdef->type == AP_INT_VALUE_FLAG) {
			*(int *)pdef->pval = pdef->intval;

		} else if (pdef->type == AP_INT_FLAG) {
			if (sscanf(argv[argi+1], "%d", (int *)pdef->pval) != 1) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as integer.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}

		} else if (pdef->type == AP_LONG_LONG_FLAG) {
			if (sscanf(argv[argi+1], "%lld", (long long *)pdef->pval) != 1) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as integer.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}

		} else if (pdef->type == AP_DOUBLE_FLAG) {
			if (!mlr_try_float_from_string(argv[argi+1], (double *)pdef->pval)) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as double.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}

		} else if (pdef->type == AP_STRING_FLAG) {
			char** pstring = pdef->pval;
			*pstring = argv[argi+1];
			pdef->pval = pstring;

		} else if (pdef->type == AP_STRING_BUILD_LIST_FLAG) {
			slls_t** pplist = pdef->pval;
			if (*pplist == NULL) {
				*pplist = slls_alloc();
			}
			slls_append_no_free(*pplist, argv[argi+1]);
			pdef->pval = pplist;

		} else if (pdef->type == AP_STRING_LIST_FLAG) {
			slls_t** pplist = pdef->pval;
			if (*pplist != NULL)
				slls_free(*pplist);
			*pplist = slls_from_line(argv[argi+1], ',', FALSE);
			pdef->pval = pplist;

		} else if (pdef->type == AP_STRING_ARRAY_FLAG) {
			string_array_t** pparray = pdef->pval;
			if (*pparray != NULL)
				string_array_free(*pparray);
			*pparray = string_array_from_line(argv[argi+1], ',');
			pdef->pval = pparray;

		} else {
			ok = FALSE;
			fprintf(stderr, "argparse.c: internal coding error: flag-def type %x not recognized.\n", pdef->type);
			fprintf(stderr, "\n");
			break;
		}

		argi += pdef->count;
	}

	*pargi = argi;
	return ok;
}
