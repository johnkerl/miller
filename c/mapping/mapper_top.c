#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/top_keeper.h"
#include "containers/mixutil.h"
#include "lib/mvfuncs.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DEFAULT_OUTPUT_FIELD_NAME "top_idx"

typedef mv_t maybe_sign_flipper_t(mv_t* pval1);

typedef struct _mapper_top_state_t {
	ap_state_t* pargp;
	slls_t* pvalue_field_names;
	slls_t* pgroup_by_field_names;
	int top_count;
	int show_full_records;
	int allow_int_float;
	maybe_sign_flipper_t* pmaybe_sign_flipper;
	lhmslv_t* groups;
	char* output_field_name;
} mapper_top_state_t;

static void      mapper_top_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_top_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_top_alloc(ap_state_t* pargp, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names,
	int top_count, int do_max, int show_full_records, int allow_int_float, char* output_field_name);
static void      mapper_top_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_top_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_top_ingest(lrec_t* pinrec, mapper_top_state_t* pstate);
static sllv_t*   mapper_top_emit(mapper_top_state_t* pstate, context_t* pctx);

// ----------------------------------------------------------------
mapper_setup_t mapper_top_setup = {
	.verb = "top",
	.pusage_func = mapper_top_usage,
	.pparse_func = mapper_top_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_top_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}    Value-field names for top counts.\n");
	fprintf(o, "-g {d,e,f}    Optional group-by-field names for top counts.\n");
	fprintf(o, "-n {count}    How many records to print per category; default 1.\n");
	fprintf(o, "-a            Print all fields for top-value records; default is\n");
	fprintf(o, "              to print only value and group-by fields. Requires a single\n");
	fprintf(o, "              value-field name only.\n");
	fprintf(o, "--min         Print top smallest values; default is top largest values.\n");
	fprintf(o, "-F            Keep top values as floats even if they look like integers.\n");
	fprintf(o, "-o {name}     Field name for output indices. Default \"%s\".\n", DEFAULT_OUTPUT_FIELD_NAME);

	fprintf(o, "Prints the n records with smallest/largest values at specified fields,\n");
	fprintf(o, "optionally by category.\n");
}

static mapper_t* mapper_top_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int     top_count             = 1;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();
	int     show_full_records     = FALSE;
	int     do_max                = TRUE;
	int     allow_int_float       = TRUE;
	char*   output_field_name     = DEFAULT_OUTPUT_FIELD_NAME;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate,         "-n",    &top_count);
	ap_define_string_list_flag(pstate, "-f",    &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g",    &pgroup_by_field_names);
	ap_define_true_flag(pstate,        "-a",    &show_full_records);
	ap_define_true_flag(pstate,        "--max", &do_max);
	ap_define_false_flag(pstate,       "--min", &do_max);
	ap_define_false_flag(pstate,       "-F",    &allow_int_float);
	ap_define_string_flag(pstate,      "-o",    &output_field_name);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_top_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (pvalue_field_names == NULL) {
		mapper_top_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (pvalue_field_names->length > 1 && show_full_records) {
		mapper_top_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_top_alloc(pstate, pvalue_field_names, pgroup_by_field_names,
		top_count, do_max, show_full_records, allow_int_float, output_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_top_alloc(ap_state_t* pargp, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names,
	int top_count, int do_max, int show_full_records, int allow_int_float, char* output_field_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_top_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_top_state_t));

	pstate->pargp                 = pargp;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->show_full_records     = show_full_records;
	pstate->allow_int_float       = allow_int_float;
	pstate->top_count             = top_count;
	pstate->pmaybe_sign_flipper   = do_max ? x_x_upos_func : x_x_uneg_func;
	pstate->groups                = lhmslv_alloc();
	pstate->output_field_name     = output_field_name;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_top_process;
	pmapper->pfree_func    = mapper_top_free;

	return pmapper;
}

static void mapper_top_free(mapper_t* pmapper, context_t* _) {
	mapper_top_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pvalue_field_names);
	slls_free(pstate->pgroup_by_field_names);

	// Free the hashmap pvvalues; the lhm free methods will free the hashmap keys.
	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
		lhmsv_t* pgroup = pa->pvvalue;
		for (lhmsve_t* pb = pgroup->phead; pb != NULL; pb = pb->pnext) {
			top_keeper_t* ptop_keeper_for_group = pb->pvvalue;
			top_keeper_free(ptop_keeper_for_group);
		}
		lhmsv_free(pgroup);
	}

	lhmslv_free(pstate->groups);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_top_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_top_state_t* pstate = pvstate;

	if (pinrec != NULL) {
		mapper_top_ingest(pinrec, pstate);
		return NULL;
	} else {
		return mapper_top_emit(pstate, pctx);
	}
}

// ----------------------------------------------------------------
static void mapper_top_ingest(lrec_t* pinrec, mapper_top_state_t* pstate) {
	// ["s", "t"]
	slls_t* pvalue_field_values    = mlr_reference_selected_values_from_record(pinrec, pstate->pvalue_field_names);
	slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

	// Heterogeneous-data case -- not all sought fields were present in record
	if (pvalue_field_values == NULL || pgroup_by_field_values == NULL) {
		slls_free(pvalue_field_values);
		slls_free(pgroup_by_field_values);
		lrec_free(pinrec);
		return;
	}
	if (pgroup_by_field_values->length != pstate->pgroup_by_field_names->length) {
		slls_free(pvalue_field_values);
		slls_free(pgroup_by_field_values);
		lrec_free(pinrec);
		return;
	}

	lhmsv_t* group_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (group_to_acc_field == NULL) {
		group_to_acc_field = lhmsv_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), group_to_acc_field, FREE_ENTRY_KEY);
	}
	slls_free(pgroup_by_field_values);

	sllse_t* pa = pstate->pvalue_field_names->phead;
	sllse_t* pb =         pvalue_field_values->phead;
	// for "x", "y" and "1", "2"
	for ( ; pa != NULL && pb != NULL; pa = pa->pnext, pb = pb->pnext) {
		char*  value_field_name = pa->value;
		char*  value_field_sval = pb->value;
		if (value_field_sval == NULL) { // Key not present
			if (pstate->show_full_records)
				lrec_free(pinrec);
			continue;
		}

		top_keeper_t* ptop_keeper_for_group = lhmsv_get(group_to_acc_field, value_field_name);
		if (ptop_keeper_for_group == NULL) {
			ptop_keeper_for_group = top_keeper_alloc(pstate->top_count);
			lhmsv_put(group_to_acc_field, value_field_name, ptop_keeper_for_group, NO_FREE);
		}

		if (*value_field_sval == 0) { // Key present with null value
			if (pstate->show_full_records)
				lrec_free(pinrec);
			continue;
		}

		mv_t value_field_nval = pstate->allow_int_float
			? mv_scan_number_or_die(value_field_sval)
			: mv_from_float(mlr_double_from_string_or_die(value_field_sval));

		// The top-keeper object will free the record if it isn't retained, or
		// keep it if it is.
		top_keeper_add(ptop_keeper_for_group, pstate->pmaybe_sign_flipper(&value_field_nval),
			(pstate->show_full_records) ? pinrec : NULL);
	}
	if (!pstate->show_full_records)
		lrec_free(pinrec);
	slls_free(pvalue_field_values);
}

// ----------------------------------------------------------------
static sllv_t* mapper_top_emit(mapper_top_state_t* pstate, context_t* pctx) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {

		// Above we required that there was only one value field in the
		// show-full-records case. That's for two reasons: (1) here, we print
		// each record at most once, which would need a change in the format
		// presented as output; (2) there would be double-frees in our
		// ingester.
		if (pstate->show_full_records) {
			lhmsv_t* group_to_acc_field = pa->pvvalue;
			for (lhmsve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
				top_keeper_t* ptop_keeper_for_group = pd->pvvalue;
				for (int i = 0;  i < ptop_keeper_for_group->size; i++) {
					sllv_append(poutrecs, ptop_keeper_for_group->top_precords[i]);
					ptop_keeper_for_group->top_precords[i] = NULL;
				}
			}
		}

		else {
			slls_t* pgroup_by_field_values = pa->key;
			for (int i = 0; i < pstate->top_count; i++) {
				lrec_t* poutrec = lrec_unbacked_alloc();

				// Add in a=s,b=t fields:
				sllse_t* pb = pstate->pgroup_by_field_names->phead;
				sllse_t* pc =         pgroup_by_field_values->phead;
				for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
					lrec_put(poutrec, pb->value, pc->value, NO_FREE);
				}

				// Add in fields such as x_top_1=#
				lhmsv_t* group_to_acc_field = pa->pvvalue;
				// for "x", "y"
				for (lhmsve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
					char* value_field_name = pd->key;
					top_keeper_t* ptop_keeper_for_group = pd->pvvalue;

					char* key = mlr_paste_2_strings(value_field_name, "_top");
					if (i < ptop_keeper_for_group->size) {
						mv_t numv = pstate->pmaybe_sign_flipper(&ptop_keeper_for_group->top_values[i]);
						char* strv = mv_alloc_format_val(&numv);
						lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(i+1), FREE_ENTRY_VALUE);
						lrec_put(poutrec, key, strv, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
					} else {
						lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(i+1), FREE_ENTRY_VALUE);
						lrec_put(poutrec, key, "", FREE_ENTRY_KEY);
					}
				}
				sllv_append(poutrecs, poutrec);
			}
		}
	}

	sllv_append(poutrecs, NULL);
	return poutrecs;
}
