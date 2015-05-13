#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_reorder_state_t {
	slls_t* pfield_name_list;
	int     put_at_end;
} mapper_reorder_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_reorder_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pinrec != NULL) {
		if (!pstate->put_at_end) {
			// OK since the field-name list was reversed at construction time.
			for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext)
				lrec_move_to_head(pinrec, pe->value);
		} else {
			for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext)
				lrec_move_to_tail(pinrec, pe->value);
		}
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static void mapper_reorder_free(void* pvstate) {
	mapper_reorder_state_t* pstate = (mapper_reorder_state_t*)pvstate;
	if (pstate->pfield_name_list != NULL)
		slls_free(pstate->pfield_name_list);
}

static mapper_t* mapper_reorder_alloc(slls_t* pfield_name_list, int put_at_end) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_reorder_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_reorder_state_t));
	pstate->pfield_name_list = pfield_name_list;
	pstate->put_at_end = put_at_end;
	if (!put_at_end)
		slls_reverse(pstate->pfield_name_list);

	pmapper->pvstate              = (void*)pstate;
	pmapper->pmapper_process_func = mapper_reorder_process;
	pmapper->pmapper_free_func    = mapper_reorder_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_reorder_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-f {a,b,c}       Field names to reorder.\n");
	fprintf(stdout, "-e           Put specified field names at record end: default is to put at record start.\n");
	fprintf(stdout, "Example: %s %s    -f a,b sends input record d=4,b=2,a=1,c=3 to a=1,b=2,d=4,c=3.\n", argv0, verb);
	fprintf(stdout, "Example: %s %s -e -f a,b sends input record d=4,b=2,a=1,c=3 to d=4,c=3,a=1,b=2.\n", argv0, verb);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reorder_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_name_list = NULL;
	int     put_at_end       = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_name_list);
	ap_define_true_flag(pstate, "-e", &put_at_end);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_reorder_usage(argv[0], verb);
		return NULL;
	}

	if (pfield_name_list == NULL) {
		mapper_reorder_usage(argv[0], verb);
		return NULL;
	}

	return mapper_reorder_alloc(pfield_name_list, put_at_end);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_reorder_setup = {
	.verb = "reorder",
	.pusage_func = mapper_reorder_usage,
	.pparse_func = mapper_reorder_parse_cli,
};
