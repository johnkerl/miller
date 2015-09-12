#include <stdio.h>
#include "lib/mlrutil.h"
#include "cli/argparse.h"

// ================================================================
#define AP_INT_VALUE_FLAG   0xf6
#define AP_INT_FLAG         0xe7
#define AP_CHAR_FLAG        0xd8
#define AP_DOUBLE_FLAG      0xc9
#define AP_STRING_FLAG      0xba
#define AP_STRING_LIST_FLAG 0xab

typedef struct _ap_flag_def_t {
	char* flag_name;
	int   type;
	int   intval;
	void* pval;
	int   count; // 1 for bool flags; 2 for the rest
} ap_flag_def_t;

static ap_flag_def_t* ap_find(ap_state_t* pstate, char* flag_name) {
	for (sllve_t* pe = pstate->pflag_defs->phead; pe != NULL; pe = pe->pnext) {
		ap_flag_def_t* pdef = pe->pvdata;
		if (streq(pdef->flag_name, flag_name))
			return pdef;
	}
	return NULL;
}

static ap_flag_def_t* ap_flag_def_alloc(char* flag_name, int ap_type, int intval, void* pval, int count) {
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
		ap_flag_def_t* pdef = pe->pvdata;
		if (pdef->type == AP_STRING_LIST_FLAG && pdef->pval != NULL) {
			slls_t** pplist = pdef->pval;
			if (*pplist != NULL)
				slls_free(*pplist);
		}
		free(pdef);
	}
	sllv_free(pstate->pflag_defs);

	free(pstate);
}

// ----------------------------------------------------------------
void ap_define_true_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, TRUE, pintval, 1));
}

void ap_define_false_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, FALSE, pintval, 1));
}

void ap_define_int_value_flag(ap_state_t* pstate, char* flag_name, int intval, int* pintval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_VALUE_FLAG, intval, pintval, 1));
}

void ap_define_char_flag(ap_state_t* pstate, char* flag_name, char* pcharval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_CHAR_FLAG, 0, pcharval, 2));
}

void ap_define_int_flag(ap_state_t* pstate, char* flag_name, int* pintval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_INT_FLAG, 0, pintval, 2));
}

void ap_define_double_flag(ap_state_t* pstate, char* flag_name, double* pdoubleval) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_DOUBLE_FLAG, 0, pdoubleval, 2));
}

void ap_define_string_flag(ap_state_t* pstate, char* flag_name, char** pstring) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_FLAG, 0, pstring, 2));
}

void ap_define_string_list_flag(ap_state_t* pstate, char* flag_name, slls_t** pplist) {
	sllv_add(pstate->pflag_defs, ap_flag_def_alloc(flag_name, AP_STRING_LIST_FLAG, 0, pplist, 2));
}

// xxx make common w/ mlrcli.c: libify
static int try_sep_from_arg(char* arg, char* pout) {
	if (streq(arg, "tab"))
		*pout = '\t';
	else if (streq(arg, "space"))
		*pout = ' ';
	else if (streq(arg, "newline"))
		*pout = '\n';
	else if (streq(arg, "pipe"))
		*pout = '|';
	else if (streq(arg, "semicolon"))
		*pout = '|';
	else if (strlen(arg) != 1)
		return FALSE;
	*pout = arg[0];
	return TRUE;
}

// ----------------------------------------------------------------
int ap_parse(ap_state_t* pstate, char* verb, int* pargi, int argc, char** argv) {

	int argi = *pargi;
	int ok = TRUE;

	while (argi < argc) {
		if (argv[argi][0] != '-') {
			break;
		}
		if (streq(argv[argi], "-h") || streq(argv[argi], "--help")) {
			ok = FALSE;
			break;
		}

		ap_flag_def_t* pdef = ap_find(pstate, argv[argi]);
		if (pdef == NULL) {
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
		} else if (pdef->type == AP_CHAR_FLAG) {
			if (!try_sep_from_arg(argv[argi+1], (char *)pdef->pval)) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as character.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}

		} else if (pdef->type == AP_INT_FLAG) {
			if (sscanf(argv[argi+1], "%d", (int *)pdef->pval) != 1) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as integer.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}

		} else if (pdef->type == AP_DOUBLE_FLAG) {
			if (!mlr_try_double_from_string(argv[argi+1], (double *)pdef->pval)) {
				fprintf(stderr, "%s %s: couldn't parse \"%s\" after \"%s\" as double.\n",
					argv[0], verb, argv[argi+1], argv[argi]);
				fprintf(stderr, "\n");
			}
		} else if (pdef->type == AP_STRING_FLAG) {
			char** pstring = pdef->pval;
			*pstring = argv[argi+1];
			pdef->pval = pstring;
		} else if (pdef->type == AP_STRING_LIST_FLAG) {
			slls_t** pplist = pdef->pval;

			if (*pplist != NULL)
				slls_free(*pplist);
			*pplist = slls_from_line(argv[argi+1], ',', FALSE);

			pdef->pval = pplist;
		} else {
			ok = FALSE;
			fprintf(stderr, "argparse.c: coding error: flag-def type %x not recognized.\n", pdef->type);
			fprintf(stderr, "\n");
			break;
		}

		argi += pdef->count;
	}

	*pargi = argi;
	return ok;
}

// ================================================================
#ifdef __AP_MAIN__
int main(int argc, char** argv) {
	int     bflag  = TRUE;
	int     intv   = 0;
	double  dblv   = 0.0;
	char*   string = NULL;
	slls_t* plist  = NULL;
	ap_state_t* pstate = ap_alloc();

	ap_define_true_flag(pstate,        "-t",   &bflag);
	ap_define_false_flag(pstate,       "-f",   &bflag);
	ap_define_int_value_flag(pstate,   "-100", 100,      &intv);
	ap_define_int_value_flag(pstate,   "-200", 200,      &intv);
	ap_define_int_flag(pstate,         "-i",   &intv);
	ap_define_double_flag(pstate,      "-d",   &dblv);
	ap_define_string_flag(pstate,      "-s",   &string);
	ap_define_string_list_flag(pstate, "-S",   &plist);

	char* verb = "stub";
	int argi = 1;
	if (ap_parse(pstate, verb, &argi, argc, argv) == TRUE) {
		printf("OK\n");
	} else {
		printf("Usage!\n");
	}

	printf("argi  is %d\n", argi);
	printf("argc  is %d\n", argc);
	printf("rest  is");
	for (; argi < argc; argi++)
		printf(" %s", argv[argi]);
	printf("\n");
	printf("bflag is %d\n", bflag);
	printf("intv  is %d\n", intv);
	printf("dblv  is %g\n", dblv);

	if (string == NULL) {
		printf("string  is null\n");
	} else {
		printf("string  is \"%s\"\n", string);
	}

	if (plist == NULL) {
		printf("list  is null\n");
	} else {
		char* out = slls_join(plist, ",");
		printf("list  is %s\n", out);
		free(out);
	}

	ap_free(pstate);

	return 0;
}
#endif
