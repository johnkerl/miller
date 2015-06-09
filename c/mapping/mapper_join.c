#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmsi.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_join_state_t {
	slls_t*  pleft_field_names;
	slls_t*  pright_field_names;
	slls_t*  poutput_field_names;
	int      allow_unsorted_input;
	int      emit_pairables;
	int      emit_left_unpairables;
	int      emit_right_unpairables;
} mapper_join_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_join_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) {
		return sllv_single(NULL);
	}
	//mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;
	//int num_found = 0;
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static void mapper_join_free(void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;
	if (pstate->pleft_field_names != NULL)
		slls_free(pstate->pleft_field_names);
	if (pstate->pright_field_names != NULL)
		slls_free(pstate->pright_field_names);
	if (pstate->poutput_field_names != NULL)
		slls_free(pstate->poutput_field_names);
}

static mapper_t* mapper_join_alloc(
	slls_t* pleft_field_names,
	slls_t* pright_field_names,
	slls_t* poutput_field_names,
	int      allow_unsorted_input,
	int      emit_pairables,
	int      emit_left_unpairables,
	int      emit_right_unpairables)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_join_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_join_state_t));
	pstate->pleft_field_names      = pleft_field_names;
	pstate->pright_field_names     = pright_field_names;
	pstate->poutput_field_names    = poutput_field_names;
	pstate->allow_unsorted_input   = allow_unsorted_input;
	pstate->emit_pairables         = emit_pairables;
	pstate->emit_left_unpairables  = emit_left_unpairables;
	pstate->emit_right_unpairables = emit_right_unpairables;

	pmapper->pvstate = (void*)pstate;
	pmapper->pprocess_func = mapper_join_process;
	pmapper->pfree_func    = mapper_join_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_join_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-1 {a,b,c}\n");
	fprintf(stdout, "-2 {a,b,c}\n");
	fprintf(stdout, "-j {a,b,c}\n");
	fprintf(stdout, "-a 1\n");
	fprintf(stdout, "-a 2\n");
	fprintf(stdout, "-v 1\n");
	fprintf(stdout, "-v 2\n");
	fprintf(stdout, "-o {a,b,c} -- do i want this?? or use then cut -f {a,b,c}...\n");
	fprintf(stdout, "-e EMPTY\n");
	fprintf(stdout, "xxx write me up.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pleft_field_names       = NULL;
	slls_t* pright_field_names      = NULL;
	slls_t* poutput_field_names     = NULL;
	int     allow_unsorted_input    = FALSE;
	int     emit_pairables          = TRUE;
	int     emit_left_unpairables   = FALSE;
	int     emit_right_unpairables  = FALSE;

	char* verb = argv[(*pargi)++];

	int argi = *pargi;
	while (argv[argi][0] == '-') {
//		if (streq(argv[argi], "--at-least")) {
//			criterion = HAVING_FIELDS_AT_LEAST;
//		} else if (streq(argv[argi], "--which-are")) {
//			criterion = HAVING_FIELDS_WHICH_ARE;
//		} else if (streq(argv[argi], "--at-most")) {
//			criterion = HAVING_FIELDS_AT_MOST;
//		} else {
//			mapper_join_usage(argv[0], verb);
//			return NULL;
//		}

		if (argc - argi < 2) {
			return NULL;
		}
		if (pleft_field_names != NULL)
			slls_free(pleft_field_names);
		pleft_field_names = slls_from_line(argv[argi+1], ',', FALSE);
		argi += 2;
	}

	if (pleft_field_names == NULL) {
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	*pargi = argi;
	return mapper_join_alloc(pleft_field_names, pright_field_names, poutput_field_names,
		allow_unsorted_input, emit_pairables, emit_left_unpairables, emit_right_unpairables);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_join_setup = {
	.verb = "join",
	.pusage_func = mapper_join_usage,
	.pparse_func = mapper_join_parse_cli,
};
