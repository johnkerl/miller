#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_reorder_state_t {
	ap_state_t* pargp;
	slls_t* pfield_name_list;
} mapper_reorder_state_t;

static void      mapper_reorder_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_reorder_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_reorder_alloc(ap_state_t* pargp, slls_t* pfield_name_list,
	int put_at_end, char* before_field_name, char* after_field_name);
static void      mapper_reorder_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_reorder_to_start(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reorder_to_end(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reorder_before(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reorder_after(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_reorder_setup = {
	.verb = "reorder",
	.pusage_func = mapper_reorder_usage,
	.pparse_func = mapper_reorder_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_reorder_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}   Field names to reorder.\n");
	fprintf(o, "-e           Put specified field names at record end: default is to put\n");
	fprintf(o, "             them at record start.\n");
	fprintf(o, "-b {x}     Put field names specified with -f before field name specified by {x},\n");
	fprintf(o, "           if any. If {x} isn't present in a given record, the specified field names\n");
	fprintf(o, "           be moved to the start of the record.\n");
	fprintf(o, "-a {x}     Put field names specified with -f after field name specified by {x},\n");
	fprintf(o, "           if any. If {x} isn't present in a given record, the specified field names\n");
	fprintf(o, "           be moved to the end of the record.\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n", argv0, verb);
	fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n", argv0, verb);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reorder_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pfield_name_list  = NULL;
	int     put_at_end        = FALSE;
	char*   before_field_name = NULL;
	char*   after_field_name  = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_name_list);
	ap_define_true_flag(pstate, "-e", &put_at_end);
	ap_define_string_flag(pstate, "-b", &before_field_name);
	ap_define_string_flag(pstate, "-a", &after_field_name);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_reorder_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_name_list == NULL) {
		mapper_reorder_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_reorder_alloc(pstate, pfield_name_list, put_at_end, before_field_name, after_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reorder_alloc(ap_state_t* pargp, slls_t* pfield_name_list,
	int put_at_end, char* before_field_name, char* after_field_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_reorder_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_reorder_state_t));
	pstate->pargp = pargp;
	pstate->pfield_name_list = pfield_name_list;

	pmapper->pvstate = (void*)pstate;

	if (put_at_end) {
		pmapper->pprocess_func = mapper_reorder_to_end;
	} else if (before_field_name != NULL) {
		pmapper->pprocess_func = mapper_reorder_before;
	} else if (after_field_name != NULL) {
		pmapper->pprocess_func = mapper_reorder_after;
	} else {
		pmapper->pprocess_func = mapper_reorder_to_start;
		slls_reverse(pstate->pfield_name_list);
	}
	pmapper->pfree_func = mapper_reorder_free;

	return pmapper;
}

static void mapper_reorder_free(mapper_t* pmapper, context_t* _) {
	mapper_reorder_state_t* pstate = pmapper->pvstate;
	if (pstate->pfield_name_list != NULL)
		slls_free(pstate->pfield_name_list);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_reorder_to_start(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pinrec != NULL) {
		// OK since the field-name list was reversed at construction time.
		for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
			lrec_move_to_head(pinrec, pe->value);
		}
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_reorder_to_end(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pinrec != NULL) {
		for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
			lrec_move_to_tail(pinrec, pe->value);
		}
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_reorder_before(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pinrec != NULL) {
		for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
			lrec_move_to_tail(pinrec, pe->value);
		}
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_reorder_after(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pinrec != NULL) {
		for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
			lrec_move_to_tail(pinrec, pe->value);
		}
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}
