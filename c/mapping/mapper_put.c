#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"
#include "mapping/lrec_evaluators.h"
#include "dsls/put_dsl_wrapper.h"
#include "cli/argparse.h"

typedef struct _mapper_put_state_t {
	ap_state_t* pargp;
	int num_evaluators;
	char** output_field_names;
	lrec_evaluator_t** pevaluators;
} mapper_put_state_t;

static sllv_t*   mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_put_free(void* pvstate);
static mapper_t* mapper_put_alloc(ap_state_t* pargp, sllv_t* pasts, int type_inferencing);
static void      mapper_put_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_put_setup = {
	.verb = "put",
	.pusage_func = mapper_put_usage,
	.pparse_func = mapper_put_parse_cli
};

// ----------------------------------------------------------------
static void mapper_put_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {expression}\n", argv0, verb);
	fprintf(o, "Adds/updates specified field(s).\n");
	fprintf(o, "Options:\n");
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
	fprintf(o, "  %s %s '$y = log10($x); $z = sqrt($y)'\n", argv0, verb);
	fprintf(o, "  %s %s '$filename = FILENAME'\n", argv0, verb);
	fprintf(o, "  %s %s '$colored_shape = $color . \"_\" . $shape'\n", argv0, verb);
	fprintf(o, "  %s %s '$y = cos($theta); $z = atan2($y, $x)'\n", argv0, verb);
	fprintf(o, "  %s %s '$name = sub($name, \"http.*com\"i, \"\")'\n", argv0, verb);
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv) {
	char* verb = argv[(*pargi)++];
	char* mlr_dsl_expression = NULL;
	int   type_inferencing = TYPE_INFER_STRING_FLOAT_INT;
	int   print_asts = FALSE;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-v", &print_asts);
	ap_define_int_value_flag(pstate, "-S", TYPE_INFER_STRING_ONLY,  &type_inferencing);
	ap_define_int_value_flag(pstate, "-F", TYPE_INFER_STRING_FLOAT, &type_inferencing);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_put_usage(stderr, argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_put_usage(stderr, argv[0], verb);
		return NULL;
	}
	mlr_dsl_expression = argv[(*pargi)++];

	// Linked list of mlr_dsl_ast_node_t*.
	sllv_t* pasts = put_dsl_parse(mlr_dsl_expression);
	if (pasts == NULL) {
		mapper_put_usage(stderr, argv[0], verb);
		return NULL;
	}

	// For just dev-testing the parser, you can do
	//   mlr put -v 'expression goes here' /dev/null
	if (print_asts) {
		for (sllve_t* pe = pasts->phead; pe != NULL; pe = pe->pnext)
			mlr_dsl_ast_node_print(pe->pvdata);
	}

	return mapper_put_alloc(pstate, pasts, type_inferencing);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_alloc(ap_state_t* pargp, sllv_t* pasts, int type_inferencing) {
	mapper_put_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_put_state_t));
	pstate->pargp = pargp;
	pstate->num_evaluators = pasts->length;
	pstate->output_field_names = mlr_malloc_or_die(pasts->length * sizeof(char*));
	pstate->pevaluators = mlr_malloc_or_die(pasts->length * sizeof(lrec_evaluator_t*));

	int i = 0;
	for (sllve_t* pe = pasts->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* past = pe->pvdata;

		if ((past->type != MLR_DSL_AST_NODE_TYPE_OPERATOR) || !streq(past->text, "=")) {
			fprintf(stderr,
				"%s: expected assignment-rooted AST; got operator \"%s\" with node type %s.\n",
					MLR_GLOBALS.argv0, past->text, mlr_dsl_ast_node_describe_type(past->type));
			return NULL;
		} else if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvdata;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvdata;

		if (pleft->type != MLR_DSL_AST_NODE_TYPE_FIELD_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		char* output_field_name = pleft->text;
		lrec_evaluator_t* pevaluator = lrec_evaluator_alloc_from_ast(pright, type_inferencing);

		pstate->pevaluators[i] = pevaluator;
		pstate->output_field_names[i] = output_field_name;
	}

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_put_process;
	pmapper->pfree_func    = mapper_put_free;

	return pmapper;
}

static void mapper_put_free(void* pvstate) {
	mapper_put_state_t* pstate = (mapper_put_state_t*)pvstate;
	free(pstate->output_field_names);
	for (int i = 0; i < pstate->num_evaluators; i++) {
		lrec_evaluator_t* pevaluator = pstate->pevaluators[i];
		pevaluator->pfree_func(pevaluator->pvstate);
	}
	free(pstate->pevaluators);
	ap_free(pstate->pargp);
}

// ----------------------------------------------------------------
static sllv_t* mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_put_state_t* pstate = (mapper_put_state_t*)pvstate;
		for (int i = 0; i < pstate->num_evaluators; i++) {
			mv_t val = pstate->pevaluators[i]->pprocess_func(pinrec,
				pctx, pstate->pevaluators[i]->pvstate);
			char free_flags;
			if (val.type == MT_STRING) {
				lrec_put(pinrec, pstate->output_field_names[i], val.u.strv, val.free_flags);
			} else {
				char* string = mv_format_val(&val, &free_flags);
				lrec_put(pinrec, pstate->output_field_names[i], string, free_flags);
			}
		}
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}
