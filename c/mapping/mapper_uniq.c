#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_uniq_state_t {
	slls_t* pgroup_by_field_names;
	int show_counts;
	lhmslv_t* pcounts_by_group;
} mapper_uniq_state_t;

static sllv_t*   mapper_uniq_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_uniq_free(void* pvstate);
static mapper_t* mapper_uniq_alloc(slls_t* pgroup_by_field_names, int show_counts);
static void      mapper_uniq_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_uniq_parse_cli(int* pargi, int argc, char** argv);
static void      mapper_count_distinct_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_count_distinct_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_count_distinct_setup = {
	.verb = "count-distinct",
	.pusage_func = mapper_count_distinct_usage,
	.pparse_func = mapper_count_distinct_parse_cli,
};

mapper_setup_t mapper_uniq_setup = {
	.verb = "uniq",
	.pusage_func = mapper_uniq_usage,
	.pparse_func = mapper_uniq_parse_cli
};

// ----------------------------------------------------------------
static void mapper_count_distinct_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}   Field names for distinct count.\n");
	fprintf(o, "Prints number of records having distinct values for specified field names.\n");
	fprintf(o, "Same as uniq -c.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_count_distinct_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_names = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_count_distinct_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_names == NULL) {
		mapper_count_distinct_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_uniq_alloc(pfield_names, TRUE);
}

// ----------------------------------------------------------------
static void mapper_uniq_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-g {d,e,f}    Group-by-field names for uniq counts\n");
	fprintf(o, "-c            Show repeat counts in addition to unique values\n");
	fprintf(o, "Prints distinct values for specified field names. With -c, same as\n");
	fprintf(o, "count-distinct.\n");
}

static mapper_t* mapper_uniq_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pgroup_by_field_names = NULL;
	int     show_counts = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate,        "-c", &show_counts);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_uniq_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pgroup_by_field_names == NULL) {
		mapper_uniq_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_uniq_alloc(pgroup_by_field_names, show_counts);
}

// ----------------------------------------------------------------
static mapper_t* mapper_uniq_alloc(slls_t* pgroup_by_field_names, int show_counts) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_uniq_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_uniq_state_t));

	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->show_counts           = show_counts;
	pstate->pcounts_by_group      = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_uniq_process;
	pmapper->pfree_func    = mapper_uniq_free;

	return pmapper;
}

static void mapper_uniq_free(void* pvstate) {
	mapper_uniq_state_t* pstate = pvstate;
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the void-star payload
	lhmslv_free(pstate->pcounts_by_group);
	pstate->pgroup_by_field_names = NULL;
	pstate->pcounts_by_group = NULL;
}

// ----------------------------------------------------------------
static sllv_t* mapper_uniq_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_uniq_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record_or_die(pinrec, pstate->pgroup_by_field_names);

		unsigned long long* pcount = lhmslv_get(pstate->pcounts_by_group, pgroup_by_field_values);
		if (pcount == NULL) {
			pcount = mlr_malloc_or_die(sizeof(unsigned long long));
			*pcount = 1LL;
			lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values), pcount);
		} else {
			(*pcount)++;
		}

		lrec_free(pinrec);
		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
			lrec_t* poutrec = lrec_unbacked_alloc();

			slls_t* pgroup_by_field_values = pa->key;

			sllse_t* pb = pstate->pgroup_by_field_names->phead;
			sllse_t* pc =         pgroup_by_field_values->phead;
			for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
				lrec_put(poutrec, pb->value, pc->value, 0);
			}

			if (pstate->show_counts) {
				unsigned long long* pcount = pa->pvvalue;
				lrec_put(poutrec, "count", mlr_alloc_string_from_ull(*pcount), LREC_FREE_ENTRY_VALUE);
			}

			sllv_add(poutrecs, poutrec);
		}
		sllv_add(poutrecs, NULL);
		return poutrecs;
	}
}
