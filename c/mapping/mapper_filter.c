#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "mapping/rval_evaluators.h"
#include "mapping/mappers.h"
#include "dsls/mlr_dsl_wrapper.h"
#include "cli/argparse.h"

typedef struct _mapper_filter_state_t {
	ap_state_t* pargp;
	char* mlr_dsl_expression;
	mlr_dsl_ast_node_t* past;
	mlhmmv_t* poosvars;
	rval_evaluator_t* pevaluator;
	int do_exclude;
} mapper_filter_state_t;

static void      mapper_filter_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_filter_alloc(ap_state_t* pargp, char* mlr_dsl_expression, mlr_dsl_ast_node_t* past,
	int type_inferencing, int do_exclude);
static void      mapper_filter_free(mapper_t* pmapper);
static sllv_t*   mapper_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

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
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-v: First prints the AST (abstract syntax tree) for the expression, which gives\n");
	fprintf(o, "    full transparency on the precedence and associativity rules of Miller's\n");
	fprintf(o, "    grammar.\n");
	fprintf(o, "-S: Keeps field values, or literals in the expression, as strings with no type \n");
	fprintf(o, "    inference to int or float.\n");
	fprintf(o, "-F: Keeps field values, or literals in the expression, as strings or floats\n");
	fprintf(o, "    with no inference to int.\n");
	fprintf(o, "-x: Prints records for which {expression} evaluates to false.\n");
	fprintf(o, "\n");
	fprintf(o, "Please use a dollar sign for field names and double-quotes for string\n");
	fprintf(o, "literals. If field names have special characters such as \".\" then you might\n");
	fprintf(o, "use braces, e.g. '${field.name}'. Miller built-in variables are\n");
	fprintf(o, "NF NR FNR FILENUM FILENAME PI E, and ENV[\"namegoeshere\"] to access environment\n");
	fprintf(o, "variables. The environment-variable name may be an expression, e.g. a field value.\n");
	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s 'log10($count) > 4.0'\n", argv0, verb);
	fprintf(o, "  %s %s 'FNR == 2          (second record in each file)'\n", argv0, verb);
	fprintf(o, "  %s %s 'urand() < 0.001'  (subsampling)\n", argv0, verb);
	fprintf(o, "  %s %s '$color != \"blue\" && $value > 4.2'\n", argv0, verb);
	fprintf(o, "  %s %s '($x<.5 && $y<.5) || ($x>.5 && $y>.5)'\n", argv0, verb);
	fprintf(o, "  %s %s '($name =~ \"^sys.*east$\") || ($name =~ \"^dev.[0-9]+\"i)'\n", argv0, verb);
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\". Please also also \"%s grep\" which is\n", argv0, argv0);
	fprintf(o, "useful when you don't yet know which field name(s) you're looking for.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv) {
	char* verb = argv[(*pargi)++];
	char* mlr_dsl_expression = NULL;
	char* expression_filename = NULL;
	int   print_ast = FALSE;
	int   type_inferencing = TYPE_INFER_STRING_FLOAT_INT;
	int   do_exclude = FALSE;

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate,    "-f", &expression_filename);
	ap_define_true_flag(pstate,      "-v", &print_ast);
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

	if (expression_filename == NULL) {
		if ((argc - *pargi) < 1) {
			mapper_filter_usage(stderr, argv[0], verb);
			return NULL;
		}
		mlr_dsl_expression = mlr_strdup_or_die(argv[(*pargi)++]);
	} else {
		mlr_dsl_expression = read_file_into_memory(expression_filename, NULL);
	}

	mlr_dsl_ast_t* past = mlr_dsl_parse(mlr_dsl_expression);
	if (past == NULL) {
		fprintf(stderr, "%s %s: syntax error on DSL parse of '%s'\n",
			argv[0], verb, mlr_dsl_expression);
		return NULL;
	}

	// For just dev-testing the parser, you can do
	//   mlr filter -v 'expression goes here' /dev/null
	if (print_ast) {
		mlr_dsl_ast_print(past);
	}

	if (past->pbegin_statements->length != 0) {
		fprintf(stderr, "%s %s: begin-statements are unsupported. Please use filter inside put.\n", argv[0], verb);
		return NULL;
	}
	if (past->pmain_statements->length != 1) {
		fprintf(stderr, "%s %s: multiple expressions are unsupported.\n", argv[0], verb);
		return NULL;
	}
	if (past->pend_statements->length != 0) {
		fprintf(stderr, "%s %s: end-statements are unsupported. Please use filter inside put.\n", argv[0], verb);
		return NULL;
	}
	mlr_dsl_ast_node_t* psubtree = sllv_pop(past->pmain_statements);
	mlr_dsl_ast_free(past);

	return mapper_filter_alloc(pstate, mlr_dsl_expression, psubtree, type_inferencing, do_exclude);
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_alloc(ap_state_t* pargp, char* mlr_dsl_expression, mlr_dsl_ast_node_t* past,
	int type_inferencing, int do_exclude)
{
	mapper_filter_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_filter_state_t));

	pstate->pargp      = pargp;
	// Retain the string contents along with any in-pointers from the AST/CST
	pstate->mlr_dsl_expression = mlr_dsl_expression;
	pstate->past       = past;
	pstate->pevaluator = rval_evaluator_alloc_from_ast(past, type_inferencing);
	pstate->poosvars   = mlhmmv_alloc();
	pstate->do_exclude = do_exclude;

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_filter_process;
	pmapper->pfree_func    = mapper_filter_free;

	return pmapper;
}

static void mapper_filter_free(mapper_t* pmapper) {
	mapper_filter_state_t* pstate = pmapper->pvstate;
	pstate->pevaluator->pfree_func(pstate->pevaluator);
	ap_free(pstate->pargp);
	mlr_dsl_ast_node_free(pstate->past);
	mlhmmv_free(pstate->poosvars);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);

	mapper_filter_state_t* pstate = pvstate;
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	sllv_t* rv = NULL;

	mv_t val = pstate->pevaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars,
		NULL, pctx, pstate->pevaluator->pvstate);
	if (mv_is_null(&val)) {
		lrec_free(pinrec);
	} else {
		mv_set_boolean_strict(&val);
		if (val.u.boolv ^ pstate->do_exclude) {
			rv = sllv_single(pinrec);
		} else {
			lrec_free(pinrec);
		}
	}

	for (lhmsve_t* pe = ptyped_overlay->phead; pe != NULL; pe = pe->pnext) {
		mv_t* pmv = pe->pvvalue;
		mv_free(pmv);
	}
	lhmsv_free(ptyped_overlay);
	return rv;
}
