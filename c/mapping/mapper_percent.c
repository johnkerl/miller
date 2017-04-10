#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsmv.h"
#include "containers/mixutil.h"
#include "containers/mvfuncs.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// pass 1:
// * append records to single record list
// * accumulate sums of percent-field names grouped by group-by field names
//   - two-level map: group-by field values -> percent field name -> mv_t sum
//     lhmslv_t -> lhmsmv
//   - how to handle het case: do sum only if record has all group-by field names
//
// pass 2:
// * iterate over record list
// * get group-by field values
// * if all are present:
//   for each group-by field name
//     get value
//     get sum
//     maybe * 100
//     maybe cumu
//     lrec_put w/ _percent or _fraction

typedef struct _mapper_percent_state_t {
	ap_state_t* pargp;
	slls_t* ppercent_field_names;
	slls_t* pgroup_by_field_names;
	sllv_t* precords;
	lhmslv_t* psums;
} mapper_percent_state_t;

static void      mapper_percent_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_percent_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_percent_alloc(ap_state_t* pargp, slls_t* ppercent_field_names, slls_t* pgroup_by_field_names);
static void      mapper_percent_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_percent_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_percent_setup = {
	.verb = "percent",
	.pusage_func = mapper_percent_usage,
	.pparse_func = mapper_percent_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_percent_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	// xxx narrate here
	// xxx streaming/two-pass disclaimer
	// -1 flag
	// --cumu flag
	fprintf(o, "-- UNDER CONSTRUCTION --\n");
	fprintf(o, "-f {a,b,c}    Field names for percent calculation\n");
	fprintf(o, "-g {d,e,f}    Optional group-by-field names for percent counts\n");
	fprintf(o, "Passes through the last n records, optionally by category.\n");
}

static mapper_t* mapper_percent_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* ppercent_field_names = slls_alloc();
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &ppercent_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_percent_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (ppercent_field_names->length == 0) {
		mapper_percent_usage(stderr, argv[0], verb);
		return NULL;
	}

	// xxx abend on empty -f

	return mapper_percent_alloc(pstate, ppercent_field_names, pgroup_by_field_names);
}

// ----------------------------------------------------------------
static mapper_t* mapper_percent_alloc(ap_state_t* pargp, slls_t* ppercent_field_names, slls_t* pgroup_by_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_percent_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_percent_state_t));

	pstate->pargp                  = pargp;
	pstate->ppercent_field_names   = ppercent_field_names;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->precords               = sllv_alloc();
	pstate->psums                  = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_percent_process;
	pmapper->pfree_func    = mapper_percent_free;

	return pmapper;
}

static void mapper_percent_free(mapper_t* pmapper, context_t* _) {
	mapper_percent_state_t* pstate = pmapper->pvstate;
	if (pstate->ppercent_field_names != NULL)
		slls_free(pstate->ppercent_field_names);
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);

	// The process method will have emptied out the list. We just need to free the container.
	sllv_free(pstate->precords);

	// xxx update
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->psums->phead; pa != NULL; pa = pa->pnext) {
//		sllv_t* precord_list_for_group = pa->pvvalue;
//		// outrecs were freed by caller of mapper_percent_process. Here, just free
//		// the sllv container itself.
//		sllv_free(precord_list_for_group);
	}
	lhmslv_free(pstate->psums);

	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_percent_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_percent_state_t* pstate = pvstate;
	if (pinrec != NULL) { // Not end of stream; pass 1

		sllv_append(pstate->precords, pinrec);

		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);

		if (pgroup_by_field_values != NULL) {
			lhmsmv_t* psums_for_group = lhmslv_get(pstate->psums, pgroup_by_field_values);
			if (psums_for_group == NULL) {
				psums_for_group = lhmsmv_alloc();
				lhmslv_put(pstate->psums, slls_copy(pgroup_by_field_values),
					psums_for_group, FREE_ENTRY_KEY);
			}

			for (sllse_t* pf = pstate->ppercent_field_names->phead; pf != NULL; pf = pf->pnext) {
				char* percent_field_name = pf->value;
				char* lrec_string_value = lrec_get(pinrec, percent_field_name);
				if (lrec_string_value != NULL) {
					mv_t lrec_num_value = mv_scan_number_or_die(lrec_string_value);
					mv_t* psum = lhmsmv_get(psums_for_group, percent_field_name);
					if (psum == NULL) { // First value for group
						lhmsmv_put(psums_for_group, percent_field_name, &lrec_num_value, FREE_ENTRY_VALUE);
					} else {
						*psum = x_xx_plus_func(psum, &lrec_num_value);
					}
				}
			}

			slls_free(pgroup_by_field_values);
		}

		return NULL;

	} else { // End of stream; pass 2
		sllv_t* poutrecs = sllv_alloc();

		while (pstate->precords->phead != NULL) {
			lrec_t* poutrec = sllv_pop(pstate->precords);

			slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(poutrec,
				pstate->pgroup_by_field_names);
			if (pgroup_by_field_values != NULL) {
				lhmsmv_t* psums_for_group = lhmslv_get(pstate->psums, pgroup_by_field_values);
				MLR_INTERNAL_CODING_ERROR_IF(psums_for_group == NULL);
				for (sllse_t* pf = pstate->ppercent_field_names->phead; pf != NULL; pf = pf->pnext) {
					char* percent_field_name = pf->value;
					mv_t* psum = lhmsmv_get(psums_for_group, percent_field_name);
					char* output_value = mv_alloc_format_val(psum); // xxx temp

					// xxx
					// maybe * 100
					// maybe cumu
					// lrec_put w/ _percent or _fraction

					char* output_field_name = mlr_paste_2_strings(percent_field_name, "_percent");
					lrec_put(poutrec, output_field_name, output_value, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
				}
				slls_free(pgroup_by_field_values);
			}

			sllv_append(poutrecs, poutrec);
		}

		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}
