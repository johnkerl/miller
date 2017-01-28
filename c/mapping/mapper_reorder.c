#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_reorder_state_t {
	ap_state_t* pargp;
	slls_t* pfield_name_list;
	int     put_at_end;
} mapper_reorder_state_t;

static void      mapper_reorder_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_reorder_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_reorder_alloc(ap_state_t* pargp, slls_t* pfield_name_list, int put_at_end);
static void      mapper_reorder_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_reorder_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

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
	fprintf(o, "Examples:\n");
	fprintf(o, "%s %s    -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"a=1,b=2,d=4,c=3\".\n",
		argv0, verb);
	fprintf(o, "%s %s -e -f a,b sends input record \"d=4,b=2,a=1,c=3\" to \"d=4,c=3,a=1,b=2\".\n",
		argv0, verb);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reorder_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pfield_name_list = NULL;
	int     put_at_end       = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_name_list);
	ap_define_true_flag(pstate, "-e", &put_at_end);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_reorder_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_name_list == NULL) {
		mapper_reorder_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_reorder_alloc(pstate, pfield_name_list, put_at_end);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reorder_alloc(ap_state_t* pargp, slls_t* pfield_name_list, int put_at_end) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_reorder_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_reorder_state_t));
	pstate->pargp = pargp;
	pstate->pfield_name_list = pfield_name_list;
	pstate->put_at_end = put_at_end;
	if (!put_at_end)
		slls_reverse(pstate->pfield_name_list);

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_reorder_process;
	pmapper->pfree_func    = mapper_reorder_free;

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
