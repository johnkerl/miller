#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "mapping/lrec_evaluators.h"
#include "mapping/mappers.h"
#include "dsls/filter_dsl_wrapper.h"
#include "cli/argparse.h"

typedef struct _mapper_filter_state_t {
	lrec_evaluator_t* pevaluator;
} mapper_filter_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_filter_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mv_t val = pstate->pevaluator->pevaluator_func(pinrec,
			pctx, pstate->pevaluator->pvstate);
		int bool_val = mt_get_boolean_strict(&val);
		if (bool_val) {
			return sllv_single(pinrec);
		} else {
			lrec_free(pinrec);
			return NULL;
		}
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static void mapper_filter_free(void* pvstate) {
	//mapper_filter_state_t* pstate = (mapper_filter_state_t*)pvstate;
	//xxx lrec_evaluator needs a pfree_func
	//if (pstate->pevaluator != NULL)
		//hss_free(pstate->pevaluator);
}

// xxx comment me ...
static mapper_t* mapper_filter_alloc(mlr_dsl_ast_node_t* past) {
	mapper_filter_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_filter_state_t));

	// xxx attempt to determine: does this AST evaluate to boolean? rather than
	// waiting to error out on the first record.
	pstate->pevaluator = lrec_evaluator_alloc_from_ast(past);

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_filter_process;
	pmapper->pfree_func    = mapper_filter_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_filter_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [-v] {expression}\n", argv0, verb);
	fprintf(stdout, "Prints records for which {expression} evaluates to true.\n");
	fprintf(stdout, "With -v, first prints the AST (abstract syntax tree) for the expression, which\n");
	fprintf(stdout, "gives full transparency on the precedence and associativity rules of Miller's grammar.\n");
	fprintf(stdout, "Please use a dollar sign for field names and double-quotes for string literals.\n");
	fprintf(stdout, "Miller built-in variables are NF NR FNR FILENUM FILENAME.\n");
	fprintf(stdout, "Examples:\n");
	fprintf(stdout, "  %s %s 'log10($count) > 4.0'\n", argv0, verb);
	fprintf(stdout, "  %s %s 'FNR == 2'\n", argv0, verb);
	fprintf(stdout, "  %s %s 'urand() < 0.001'\n", argv0, verb);
	fprintf(stdout, "  %s %s '$color != \"blue\" && $value > 4.2'\n", argv0, verb);
	fprintf(stdout, "  %s %s '($x<.5 && $y<.5) || ($x>.5 && $y>.5)'\n", argv0, verb);
	fprintf(stdout, "Please see http://johnkerl.org/miller/doc/reference.html for more information including function list.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv) {
	char* verb = argv[(*pargi)++];
	char* mlr_dsl_expression = NULL;
	int   print_asts = FALSE;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-v", &print_asts);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_filter_usage(argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_filter_usage(argv[0], verb);
		return NULL;
	}
	mlr_dsl_expression = argv[(*pargi)++];

	mlr_dsl_ast_node_holder_t* past = filter_dsl_parse(mlr_dsl_expression);
	if (past == NULL) {
		mapper_filter_usage(argv[0], verb);
		return NULL;
	}
	if (print_asts) {
		mlr_dsl_ast_node_print(past->proot);
	}

	return mapper_filter_alloc(past->proot);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_filter_setup = {
	.verb = "filter",
	.pusage_func = mapper_filter_usage,
	.pparse_func = mapper_filter_parse_cli,
};
