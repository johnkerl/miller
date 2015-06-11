#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_join_state_t {
	slls_t*  pleft_field_names;
	slls_t*  pright_field_names;
	slls_t*  poutput_field_names;
	int      emit_pairables;
	int      emit_left_unpairables;
	int      emit_right_unpairables;
	char*    left_file_name;
	lhmslv_t* precords_by_key_field_names; // For unsorted input
} mapper_join_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_unsorted(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) {
		return sllv_single(NULL);
	}
	//mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;
	//int num_found = 0;
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_sorted(lrec_t* pinrec, context_t* pctx, void* pvstate) {
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

// ----------------------------------------------------------------
static lhmslv_t* ingest_left_file(char* left_file_name) {
	lhmslv_t* precords_by_key_field_names = lhmslv_alloc();
	while (TRUE) {
		lrec_t* pinrec = NULL; // need reader ...
		if (pinrec == NULL)
			break;

		slls_t* pkey_field_names = mlr_keys_from_record(pinrec);
		sllv_t* plist = lhmslv_get(precords_by_key_field_names, pkey_field_names);
		if (plist == NULL) {
			plist = sllv_alloc();
			sllv_add(plist, pinrec);
			lhmslv_put(precords_by_key_field_names, slls_copy(pkey_field_names), plist);
		} else {
			sllv_add(plist, pinrec);
		}
	}
	return precords_by_key_field_names;
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_alloc(
	slls_t* pleft_field_names,
	slls_t* pright_field_names,
	slls_t* poutput_field_names,
	int      allow_unsorted_input,
	int      emit_pairables,
	int      emit_left_unpairables,
	int      emit_right_unpairables,
	char*    left_file_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_join_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_join_state_t));
	pstate->pleft_field_names      = pleft_field_names;
	pstate->pright_field_names     = pright_field_names;
	pstate->poutput_field_names    = poutput_field_names;
	pstate->emit_pairables         = emit_pairables;
	pstate->emit_left_unpairables  = emit_left_unpairables;
	pstate->emit_right_unpairables = emit_right_unpairables;

	pmapper->pvstate = (void*)pstate;
	if (allow_unsorted_input) {
		pstate->precords_by_key_field_names = ingest_left_file(left_file_name);
		pstate->left_file_name = NULL;
		pmapper->pprocess_func = mapper_join_process_unsorted;
	} else {
		pstate->precords_by_key_field_names = NULL;
		pstate->left_file_name = left_file_name;
		pmapper->pprocess_func = mapper_join_process_sorted;
	}
	pmapper->pfree_func    = mapper_join_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_join_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "xxx write me up.\n");
	fprintf(stdout, "-f {left file name}\n");
	fprintf(stdout, "-l {a,b,c}\n");
	fprintf(stdout, "-r {a,b,c}\n");
	fprintf(stdout, "-o {a,b,c}\n");
	fprintf(stdout, "--np\n");
	fprintf(stdout, "--ul\n");
	fprintf(stdout, "--ur\n");
	fprintf(stdout, "-u\n");
	fprintf(stdout, "-e EMPTY\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv) {
	char*   left_file_name          = NULL;
	slls_t* pleft_field_names       = NULL;
	slls_t* pright_field_names      = NULL;
	slls_t* poutput_field_names     = NULL;
	int     allow_unsorted_input    = FALSE;
	int     emit_pairables          = TRUE;
	int     emit_left_unpairables   = FALSE;
	int     emit_right_unpairables  = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate,      "-f",   &left_file_name);
	ap_define_string_list_flag(pstate, "-l",   &pleft_field_names);
	ap_define_string_list_flag(pstate, "-r",   &pright_field_names);
	ap_define_string_list_flag(pstate, "-o",   &poutput_field_names);
	ap_define_false_flag(pstate,       "--np", &emit_pairables);
	ap_define_int_flag(pstate,         "--ul", &emit_left_unpairables);
	ap_define_int_flag(pstate,         "--ur", &emit_left_unpairables);
	ap_define_true_flag(pstate,        "-u",   &allow_unsorted_input);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (left_file_name == NULL) {
		fprintf(stderr, "%s %s: need left file name\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}
	if (pleft_field_names == NULL) {
		fprintf(stderr, "%s %s: need left field names\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (pright_field_names == NULL)
		pright_field_names = slls_copy(pleft_field_names);
	if (poutput_field_names == NULL)
		poutput_field_names = slls_copy(pleft_field_names);

	return mapper_join_alloc(pleft_field_names, pright_field_names, poutput_field_names,
		allow_unsorted_input, emit_pairables, emit_left_unpairables, emit_right_unpairables,
		left_file_name);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_join_setup = {
	.verb = "join",
	.pusage_func = mapper_join_usage,
	.pparse_func = mapper_join_parse_cli,
};
