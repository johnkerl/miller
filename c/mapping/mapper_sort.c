#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

#define SORT_NUMERIC    0x40
#define SORT_DESCENDING 0x80

typedef struct _mapper_sort_state_t {
	slls_t* pkey_field_names;
	int*    psort_params;
	// map from list of string to list of record
	lhmslv_t* precords_by_key_field_names;
	int do_sort;
} mapper_sort_state_t;

static int cmp_for_qsort(const void* pva, const void* pvb);

// state for bucket choices:
// o slls of key-field values
// n slls/array of key_field_names
// config for the sort:
// * compare two arrays of u.string/double by array of num/lex asc/dec

// need num/lex & +/- flags for the sort -- static context for the non-reentrant qsort
// 'v' payload should be a pair of:
// * union of double & char*: here do the sscanf of the double if doing numeric sort on that field || abend
// * record-list as now
// at the end make an array of the pairs and qsort them on the union part
// then emit the record-lists as now
//
// before doing that, do the num/lex & the +/- at the CLI interface & handle (inefficiently)
// via repeated sscanf in the comparator:
// -rn a,b -fn c,d,e -r f,g -f h,i,j,k -- OR--
// -nr a,b -nf c,d,e -r f,g -f h,i,j,k
//
// data structure down there: sllv of sort info which is triple of char* field_name, char do_num, char do_rev.

// xxx cmt about re-entrancy
static int* pcmp_sort_params = NULL;
static int cmp_for_qsort(const void* pva, const void* pvb) {
	const lhmslve_t* pea = pva;
	const lhmslve_t* peb = pvb;
	slls_t* pa = pea->key;
	slls_t* pb = peb->key;
	if (pa->length != pb->length)
		return pa->length - pb->length;
	sllse_t* pe = pa->phead;
	sllse_t* pf = pb->phead;
	for (int i = 0; pe != NULL && pf != NULL; pe = pe->pnext, pf = pf->pnext, i++) {
		if (pcmp_sort_params[i] & SORT_NUMERIC) {
			double e, f;
			if (sscanf(pe->value, "%lf", &e) != 1) {
				fprintf(stderr, "xxx b04k!\n");
				exit(1);
			}
			if (sscanf(pf->value, "%lf", &f) != 1) {
				fprintf(stderr, "xxx b04k!\n");
				exit(1);
			}
			double d = e - f;
			int s = (d < 0) ? -1 : 1;
			if (s != 0)
				return (pcmp_sort_params[i] & SORT_DESCENDING) ? -s : s;
		} else {
			int s = strcmp(pe->value, pf->value);
			if (s != 0)
				return (pcmp_sort_params[i] & SORT_DESCENDING) ? -s : s;
		}
	}
	return 0;
}

// ----------------------------------------------------------------
sllv_t* mapper_sort_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_sort_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pkey_field_values = mlr_selected_values_from_record(pinrec, pstate->pkey_field_names);
		sllv_t* plist = lhmslv_get(pstate->precords_by_key_field_names, pkey_field_values);
		if (plist == NULL) {
			plist = sllv_alloc();
			sllv_add(plist, pinrec);
			lhmslv_put(pstate->precords_by_key_field_names, slls_copy(pkey_field_values), plist);
		} else {
			sllv_add(plist, pinrec);
		}
		return NULL;
	}
	else if (!pstate->do_sort) {
		sllv_t* poutput = sllv_alloc();
		for (lhmslve_t* pe = pstate->precords_by_key_field_names->phead; pe != NULL; pe = pe->pnext) {
			sllv_t* plist = pe->value;
			for (sllve_t* pf = plist->phead; pf != NULL; pf = pf->pnext) {
				sllv_add(poutput, pf->pvdata);
			}
		}
		sllv_add(poutput, NULL);
		return poutput;
	} else {
		int num_lists = pstate->precords_by_key_field_names->num_occupied;
		lhmslve_t* pairs = mlr_malloc_or_die(num_lists * sizeof(lhmslve_t));

		int i = 0;
		for (lhmslve_t* pe = pstate->precords_by_key_field_names->phead; pe != NULL; pe = pe->pnext, i++) {
			pairs[i].key   = pe->key;
			pairs[i].value = pe->value;
		}

		pcmp_sort_params = pstate->psort_params;
		qsort(pairs, num_lists, sizeof(pairs[0]), cmp_for_qsort);
		pcmp_sort_params = NULL;

		sllv_t* poutput = sllv_alloc();
		for (i = 0; i < num_lists; i++) {
			sllv_t* plist = pairs[i].value;
			for (sllve_t* pf = plist->phead; pf != NULL; pf = pf->pnext) {
				sllv_add(poutput, pf->pvdata);
			}
		}
		free(pairs);
		sllv_add(poutput, NULL);
		return poutput;
	}
}

// ----------------------------------------------------------------
static void mapper_sort_free(void* pvstate) {
	mapper_sort_state_t* pstate = pvstate;
	if (pstate->pkey_field_names != NULL)
		slls_free(pstate->pkey_field_names);
	if (pstate->precords_by_key_field_names != NULL)
		// xxx free void-star payloads 1st
		lhmslv_free(pstate->precords_by_key_field_names);
	if (pstate->psort_params != NULL)
		free(pstate->psort_params);
}

mapper_t* mapper_sort_alloc(slls_t* pkey_field_names, int* psort_params, int do_sort) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sort_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sort_state_t));

	pstate->pkey_field_names            = pkey_field_names;
	pstate->psort_params                = psort_params;
	pstate->precords_by_key_field_names = lhmslv_alloc();
	pstate->do_sort                     = do_sort;

	// xxx put this line by the next x all mapper_*_allocs
	pmapper->pvstate = pstate;
	pmapper->pmapper_process_func = mapper_sort_func;
	pmapper->pmapper_free_func    = mapper_sort_free;

	return pmapper;
}

// ----------------------------------------------------------------
void mapper_group_by_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s {comma-separated field names}\n", argv0, verb);
}
mapper_t* mapper_group_by_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 2) {
		mapper_group_by_usage(argv[0], argv[*pargi]);
		return NULL;
	}

	slls_t* pnames = slls_from_line(argv[*pargi+1], ',', FALSE);

	*pargi += 2;
	return mapper_sort_alloc(pnames, NULL, FALSE);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_group_by_setup = {
	.verb = "group-by",
	.pusage_func = mapper_group_by_usage,
	.pparse_func = mapper_group_by_parse_cli
};

// ----------------------------------------------------------------
void mapper_sort_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s {xxx}\n", argv0, verb);
	fprintf(stdout, "       xxx info goes here\n");
}
mapper_t* mapper_sort_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 3) {
		mapper_sort_usage(argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;
	slls_t* pnames = slls_alloc();
	slls_t* pflags = slls_alloc();

	while ((argc - *pargi) >= 1 && argv[*pargi][0] == '-') {
		if ((argc - *pargi) < 2)
			mapper_sort_usage(argv[0], verb);
		char* flag  = argv[*pargi];
		char* value = argv[*pargi+1];
		*pargi += 2;

		if (streq(flag, "-f")) {
		} else if (streq(flag, "-n")) {
		} else if (streq(flag, "-nf")) {
		} else if (streq(flag, "-r")) {
		} else if (streq(flag, "-nr")) {
		} else {
			mapper_sort_usage(argv[0], verb);
		}
		slls_t* pnames_for_flag = slls_from_line(value, ',', FALSE);
		for (sllse_t* pe = pnames_for_flag->phead; pe != NULL; pe = pe->pnext) {
			slls_add_no_free(pnames, pe->value);
			slls_add_no_free(pflags, flag);
		}
		slls_free(pnames_for_flag);
	}

	if (pnames->length < 1)
		mapper_sort_usage(argv[0], verb);

	int* opt_array  = mlr_malloc_or_die(pnames->length * sizeof(int));
	sllse_t* pe;
	int di;
	for (pe = pflags->phead, di = 0; pe != NULL; pe = pe->pnext, di++) {
		char* flag = pe->value;
		int opt =
			streq(flag, "-nf") ? SORT_NUMERIC :
			streq(flag, "-n")  ? SORT_NUMERIC :
			streq(flag, "-r")  ? SORT_DESCENDING :
			streq(flag, "-nr") ? SORT_NUMERIC|SORT_DESCENDING :
			0;
		opt_array[di] =opt;
	}
	slls_free(pflags);

	return mapper_sort_alloc(pnames, opt_array, TRUE);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_sort_setup = {
	.verb = "sort",
	.pusage_func = mapper_sort_usage,
	.pparse_func = mapper_sort_parse_cli
};
