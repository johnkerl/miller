#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "mapping/lrec_evaluators.h"
#include "mapping/mappers.h"
#include "dsls/filter_dsl_wrapper.h"
#include "cli/argparse.h"

typedef struct _mapper_filter_state_t {
	ap_state_t* pargp;
	lrec_evaluator_t* pevaluator;
	int do_exclude;
} mapper_filter_state_t;

static sllv_t*   mapper_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_filter_free(mapper_t* pmapper);
static mapper_t* mapper_filter_alloc(ap_state_t* pargp, mlr_dsl_ast_node_t* past,
	int type_inferencing, int do_exclude);
static void      mapper_filter_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_filter_setup = {
	.verb = "filter",
	.pusage_func = mapper_filter_usage,
	.pparse_func = mapper_filter_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_filter_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {expression}\n", argv0, verb);
	fprintf(o, "Prints records for which {expression} evaluates to true.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-x: Prints records for which {expression} evaluates to false.\n");
	fprintf(o, "-v: First prints the AST (abstract syntax tree) for the expression, which gives\n");
	fprintf(o, "    full transparency on the precedence and associativity rules of Miller's\n");
	fprintf(o, "    grammar.\n");
	fprintf(o, "-S: Keeps field values, or literals in the expression, as strings with no type \n");
	fprintf(o, "    inference to int or float.\n");
	fprintf(o, "-F: Keeps field values, or literals in the expression, as strings or floats\n");
	fprintf(o, "    with no inference to int.\n");
	fprintf(o, "Please use a dollar sign for field names and double-quotes for string\n");
	fprintf(o, "literals. If field names have special characters such as \".\" then you might\n");
	fprintf(o, "use braces, e.g. '${field.name}'. Miller built-in variables are\n");
	fprintf(o, "NF NR FNR FILENUM FILENAME PI E.\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s 'log10($count) > 4.0'\n", argv0, verb);
	fprintf(o, "  %s %s 'FNR == 2          (second record in each file)'\n", argv0, verb);
	fprintf(o, "  %s %s 'urand() < 0.001'  (subsampling)\n", argv0, verb);
	fprintf(o, "  %s %s '$color != \"blue\" && $value > 4.2'\n", argv0, verb);
	fprintf(o, "  %s %s '($x<.5 && $y<.5) || ($x>.5 && $y>.5)'\n", argv0, verb);
	fprintf(o, "  %s %s '($name =~ \"^sys.*east$\") || ($name =~ \"^dev.[0-9]+\"i)'\n", argv0, verb);
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\". Please also also \"%s grep\" which is\n", argv0, argv0);
	fprintf(o, "useful when you don't yet know which field name(s) you're looking for.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv) {
	char* verb = argv[(*pargi)++];
	char* mlr_dsl_expression = NULL;
	int   print_asts = FALSE;
	int   type_inferencing = TYPE_INFER_STRING_FLOAT_INT;
	int   do_exclude = FALSE;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate,      "-v", &print_asts);
	ap_define_int_value_flag(pstate, "-S", TYPE_INFER_STRING_ONLY,  &type_inferencing);
	ap_define_int_value_flag(pstate, "-F", TYPE_INFER_STRING_FLOAT, &type_inferencing);
	ap_define_true_flag(pstate,      "-x", &do_exclude);

	// Pass error_on_unrecognized == FALSE to ap_parse so expressions starting
	// with a minus sign aren't treated as errors. Example: "mlr filter '-$x ==
	// $y'".
	if (!ap_parse_aux(pstate, verb, pargi, argc, argv, FALSE)) {
		mapper_filter_usage(stderr, argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_filter_usage(stderr, argv[0], verb);
		return NULL;
	}
	mlr_dsl_expression = argv[(*pargi)++];

	mlr_dsl_ast_node_holder_t* past = filter_dsl_parse(mlr_dsl_expression);
	if (past == NULL) {
		mapper_filter_usage(stderr, argv[0], verb);
		return NULL;
	}

	// For just dev-testing the parser, you can do
	//   mlr filter -v 'expression goes here' /dev/null
	if (print_asts) {
		mlr_dsl_ast_node_print(past->proot);
	}

	return mapper_filter_alloc(pstate, past->proot, type_inferencing, do_exclude);
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_alloc(ap_state_t* pargp, mlr_dsl_ast_node_t* past,
	int type_inferencing, int do_exclude)
{
	mapper_filter_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_filter_state_t));

	pstate->pargp      = pargp;
	pstate->pevaluator = lrec_evaluator_alloc_from_ast(past, type_inferencing);
	pstate->do_exclude = do_exclude;

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_filter_process;
	pmapper->pfree_func    = mapper_filter_free;

	return pmapper;
}

static void mapper_filter_free(mapper_t* pmapper) {
	mapper_filter_state_t* pstate = pmapper->pvstate;
	pstate->pevaluator->pfree_func(pstate->pevaluator->pvstate);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_filter_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mv_t val = pstate->pevaluator->pprocess_func(pinrec, pctx, pstate->pevaluator->pvstate);
		if (val.type == MT_NULL) {
			lrec_free(pinrec);
			return NULL;
		} else {
			mv_set_boolean_strict(&val);
			if (val.u.boolv ^ pstate->do_exclude) {
				return sllv_single(pinrec);
			} else {
				lrec_free(pinrec);
				return NULL;
			}
		}
	} else {
		return sllv_single(NULL);
	}
}
