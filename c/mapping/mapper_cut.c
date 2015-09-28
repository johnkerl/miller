#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_cut_state_t {
	slls_t* pfield_name_list;
	hss_t*  pfield_name_set;
	int     do_arg_order;
	int     do_complement;
} mapper_cut_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_cut_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_cut_state_t* pstate = (mapper_cut_state_t*)pvstate;
		if (!pstate->do_complement) {
			// Loop over the record and free the fields not in the
			// to-be-retained set, being careful about the fact that we're
			// modifying what we're looping over.
			for (lrece_t* pe = pinrec->phead; pe != NULL; /* next in loop */) {
				if (!hss_has(pstate->pfield_name_set, pe->key)) {
					lrece_t* pf = pe->pnext;
					lrec_remove(pinrec, pe->key);
					pe = pf;
				} else {
					pe = pe->pnext;
				}
			}
			if (pstate->do_arg_order) {
				// OK since the field-name list was reversed at construction time.
				for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
					char* field_name = pe->value;
					lrec_move_to_head(pinrec, field_name);
				}
			}
			return sllv_single(pinrec);
		} else {
			for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
				char* field_name = pe->value;
				lrec_remove(pinrec, field_name);
			}
			return sllv_single(pinrec);
		}
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static void mapper_cut_free(void* pvstate) {
	mapper_cut_state_t* pstate = (mapper_cut_state_t*)pvstate;
	if (pstate->pfield_name_list != NULL)
		slls_free(pstate->pfield_name_list);
	hss_free(pstate->pfield_name_set);
}

static mapper_t* mapper_cut_alloc(slls_t* pfield_name_list, int do_arg_order, int do_complement) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_cut_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_cut_state_t));
	pstate->pfield_name_list   = pfield_name_list;
	slls_reverse(pstate->pfield_name_list);
	pstate->pfield_name_set    = hss_from_slls(pfield_name_list);
	pstate->do_arg_order       = do_arg_order;
	pstate->do_complement      = do_complement;

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_cut_process;
	pmapper->pfree_func    = mapper_cut_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_cut_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}       Field names to include for cut.\n");
	fprintf(o, "-o               Retain fields in the order specified here in the argument list.\n");
	fprintf(o, "                 Default is to retain them in the order found in the input data.\n");
	fprintf(o, "-x|--complement  Exclude, rather that include, field names specified by -f.\n");
	fprintf(o, "Passes through input records with specified fields included/excluded.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_cut_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_name_list  = NULL;
	int     do_arg_order  = FALSE;
	int     do_complement = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_name_list);
	ap_define_true_flag(pstate, "-o",           &do_arg_order);
	ap_define_true_flag(pstate, "-x",           &do_complement);
	ap_define_true_flag(pstate, "--complement", &do_complement);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_cut_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_name_list == NULL) {
		mapper_cut_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_cut_alloc(pfield_name_list, do_arg_order, do_complement);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_cut_setup = {
	.verb = "cut",
	.pusage_func = mapper_cut_usage,
	.pparse_func = mapper_cut_parse_cli,
};
