#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "cli/argparse.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmsv.h"
#include "containers/mlhmmv.h"
#include "mapping/mappers.h"
#include "mapping/lrec_evaluators.h"
#include "dsls/mlr_dsl_wrapper.h"
#include "mlr_dsl_cst.h"

typedef struct _mapper_put_state_t {
	ap_state_t*    pargp;
	mlr_dsl_ast_t* past;
	mlr_dsl_cst_t* pcst;
	int            at_begin;
	lhmsv_t*       poosvars;
	mlhmmv_t*      pmoosvars;
} mapper_put_state_t;

static mapper_t* mapper_put_alloc(ap_state_t* pargp, mlr_dsl_ast_t* past, int type_inferencing);
static void      mapper_put_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv);
static sllv_t*   mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_put_free(mapper_t* pmapper);

static void evaluate_statements(
	lrec_t*             pinrec,
	lhmsv_t*            ptyped_overlay,
	string_array_t*     pregex_captures,
	mapper_put_state_t* pstate,
	context_t*          pctx,
	sllv_t*             pcst_statements,
	int*                pemit_rec,
	sllv_t*             poutrecs
);

// ----------------------------------------------------------------
mapper_setup_t mapper_put_setup = {
	.verb = "put",
	.pusage_func = mapper_put_usage,
	.pparse_func = mapper_put_parse_cli
};

// ----------------------------------------------------------------
static void mapper_put_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {expression}\n", argv0, verb);
	fprintf(o, "Adds/updates specified field(s). Expressions are semicolon-separated and must\n");
	fprintf(o, "either be assignments, or evaluate to boolean.  Each expression is evaluated in\n");
	fprintf(o, "turn from left to right. Assignment expressions are applied to the current\n");
	fprintf(o, "record; once a boolean expression evaluates to false, the record is emitted\n");
	fprintf(o, "with all changes up to that point and remaining expressions to the right are\n");
	fprintf(o, "not evaluated.\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-v: First prints the AST (abstract syntax tree) for the expression, which gives\n");
	fprintf(o, "    full transparency on the precedence and associativity rules of Miller's\n");
	fprintf(o, "    grammar.\n");
	fprintf(o, "-S: Keeps field values, or literals in the expression, as strings with no type \n");
	fprintf(o, "    inference to int or float.\n");
	fprintf(o, "-F: Keeps field values, or literals in the expression, as strings or floats\n");
	fprintf(o, "    with no inference to int.\n");
	fprintf(o, "\n");
	fprintf(o, "Please use a dollar sign for field names and double-quotes for string\n");
	fprintf(o, "literals. If field names have special characters such as \".\" then you might\n");
	fprintf(o, "use braces, e.g. '${field.name}'. Miller built-in variables are\n");
	fprintf(o, "NF NR FNR FILENUM FILENAME PI E.\n");
	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  Assignment only:\n");
	fprintf(o, "  %s %s '$y = log10($x); $z = sqrt($y)'\n", argv0, verb);
	fprintf(o, "  %s %s '$filename = FILENAME'\n", argv0, verb);
	fprintf(o, "  %s %s '$colored_shape = $color . \"_\" . $shape'\n", argv0, verb);
	fprintf(o, "  %s %s '$y = cos($theta); $z = atan2($y, $x)'\n", argv0, verb);
	fprintf(o, "  %s %s '$name = sub($name, \"http.*com\"i, \"\")'\n", argv0, verb);
	fprintf(o, "  Mixed assignment/boolean:\n");
	fprintf(o, "  %s %s '$x > 0.0; $y = log10($x); $z = sqrt($y)'\n", argv0, verb);
	fprintf(o, "  %s %s '$y = log10($x); 1.1 < $y && $y < 7.0; $z = sqrt($y)'\n", argv0, verb);
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\".\n", argv0);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv) {
	char* verb = argv[(*pargi)++];
	char* mlr_dsl_expression = NULL;
	int   type_inferencing = TYPE_INFER_STRING_FLOAT_INT;
	int   print_ast = FALSE;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate,      "-v", &print_ast);
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
	mlr_dsl_ast_t* past = mlr_dsl_parse(mlr_dsl_expression);
	if (past == NULL) {
		fprintf(stderr, "%s %s: syntax error on DSL parse of '%s'\n",
			argv[0], verb, mlr_dsl_expression);
		return NULL;
	}

	// For just dev-testing the parser, you can do
	//   mlr put -v 'expression goes here' /dev/null
	if (print_ast)
		mlr_dsl_ast_print(past);

	return mapper_put_alloc(pstate, past, type_inferencing);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_alloc(ap_state_t* pargp, mlr_dsl_ast_t* past, int type_inferencing) {
	mapper_put_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_put_state_t));
	pstate->pargp     = pargp;
	pstate->past      = past;
	pstate->pcst      = mlr_dsl_cst_alloc(past, type_inferencing);
	pstate->at_begin  = TRUE;
	pstate->poosvars  = lhmsv_alloc();
	pstate->pmoosvars = mlhmmv_alloc();

	mapper_t* pmapper      = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_put_process;
	pmapper->pfree_func    = mapper_put_free;

	return pmapper;
}

static void mapper_put_free(mapper_t* pmapper) {
	mapper_put_state_t* pstate = pmapper->pvstate;

	for (lhmsve_t* pe = pstate->poosvars->phead; pe != NULL; pe = pe->pnext)
		mv_free(pe->pvvalue);
	lhmsv_free(pstate->poosvars);

	mlhmmv_free(pstate->pmoosvars);

	mlr_dsl_cst_free(pstate->pcst);
	mlr_dsl_ast_free(pstate->past);

	ap_free(pstate->pargp);

	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
// The typed-overlay holds intermediate values such as in
//
//   echo x=1 | mlr put '$y = string($x); $z = $y . $y'
//
// because otherwise
// * lrecs store insertion-ordered maps of string to string (since this is ultimately file I/O);
// * types are inferred at entry to put;
// * x=1 would be inferred to int; string($x) would be string; written back to the
//   lrec, y would be "1" which would be re-inferred to int.
//
// So the typed overlay allows us to remember that y is string "1" not integer 1.
//
// But this raises the question: why stop here? Why not have lrecs be insertion-ordered maps from
// string to mlrval? Then we could preserve types for the duration of each lrec, not just for
// the duration of the put operation. Reasons:
// * The compare_lexically operation would suffer a performance regression;
// * Worse, all lhmslv group-by operations (used by many Miller verbs) would likewise suffer a performance regression.
//
// ----------------------------------------------------------------
// The regex-capture string-array holds copies of regex matches, e.g. in
//
//   echo name=abc_def | mlr put '$name =~ "(.*)_(.*)"; $left = "\1"; $right = "\2"'
//
// produces a record with left=abc and right=def.
//
// There is an important trick here with the length of the string-array:
// * It is allocated here with length 0.
// * It is passed by reference to the lrec-evaluator tree. In particular, the matches and does-not-match functions
//   (which implement the =~ and !=~ operators) resize it and populate it.
// * For simplicity, it is a 1-up array: so \1, \2, \3 are at array indices 1, 2, 3.
// * If the matches/does-not-match functions are entered, even with no matches, the regex-captures string-array
//   will be resized to have length at least 1: length 1 for 0 matches, length 2 for 1 match, etc. since
//   the array is indexed 1-up.
// * When the lrec-evaluator's from-literal function is invoked, the interpolate_regex_captures function can quickly
//   check to see if the regex-captures array has length 0 and thereby know that a time-consuming scan for \1, \2, \3,
//   etc. does not need to be done.

static sllv_t* mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_put_state_t* pstate = (mapper_put_state_t*)pvstate;

	string_array_t* pregex_captures = string_array_alloc(0);
	sllv_t* poutrecs = sllv_alloc();
	int emit_rec = TRUE;

	if (pstate->at_begin) {
		evaluate_statements(NULL, NULL, pregex_captures, pstate, pctx,
			pstate->pcst->pbegin_statements, &emit_rec, poutrecs);
		pstate->at_begin = FALSE;
	}

	if (pinrec == NULL) { // End of input stream
		evaluate_statements(NULL, NULL, pregex_captures, pstate, pctx,
			pstate->pcst->pend_statements, &emit_rec, poutrecs);

		string_array_free(pregex_captures);
		sllv_add(poutrecs, NULL);
		return poutrecs;
	}

	lhmsv_t* ptyped_overlay = lhmsv_alloc();

	evaluate_statements(pinrec, ptyped_overlay, pregex_captures, pstate, pctx,
		pstate->pcst->pmain_statements, &emit_rec, poutrecs);

	if (emit_rec) {
		// Write the output fields from the typed overlay back to the lrec.
		for (lhmsve_t* pe = ptyped_overlay->phead; pe != NULL; pe = pe->pnext) {
			char* output_field_name = pe->key;
			mv_t* pval = pe->pvvalue;

			if (pval->type == MT_STRING) {
				// Ownership transfer from mv_t to lrec.
				lrec_put(pinrec, output_field_name, pval->u.strv, pval->free_flags);
			} else {
				char free_flags = NO_FREE;
				char* string = mv_format_val(pval, &free_flags);
				lrec_put(pinrec, output_field_name, string, free_flags);
			}
			free(pval);

		}
	}
	lhmsv_free(ptyped_overlay);
	string_array_free(pregex_captures);

	if (emit_rec) {
		sllv_add(poutrecs, pinrec);
	} else {
		lrec_free(pinrec);
	}
	return poutrecs;
}

// ----------------------------------------------------------------
static void evaluate_statements(
	lrec_t*             pinrec,
	lhmsv_t*            ptyped_overlay,
	string_array_t*     pregex_captures,
	mapper_put_state_t* pstate,
	context_t*          pctx,
	sllv_t*             pcst_statements,
	int*                pemit_rec,
	sllv_t*             poutrecs
) {

	// xxx move some/all of this into mlr_dsl_cst.c -- ?

	// Do the evaluations, writing typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	*pemit_rec = TRUE;

	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;

		int node_type = pstatement->ast_node_type;

		if (node_type == MD_AST_NODE_TYPE_SREC_ASSIGNMENT || node_type == MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT) {
			mlr_dsl_cst_statement_item_t* pitem = pstatement->pitems->phead->pvvalue;
			int lhs_type = pitem->lhs_type;
			char* output_field_name = pitem->output_field_name;
			lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

			mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars, pstate->pmoosvars,
				pregex_captures, pctx, prhs_evaluator->pvstate);
			mv_t* pval = mlr_malloc_or_die(sizeof(mv_t));
			*pval = val;

			if (lhs_type == MLR_DSL_CST_LHS_TYPE_OOSVAR) {
				lhmsv_put(pstate->poosvars, output_field_name, pval, NO_FREE);
			} else {
				// The lrec_evaluator reads the overlay in preference to the lrec. E.g. if the input had
				// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
				// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
				// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
				// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
				// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
				// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
				// here. And putting something statically allocated minimizes copying/freeing.
				lhmsv_put(ptyped_overlay, output_field_name, pval, NO_FREE);
				lrec_put(pinrec, output_field_name, "bug", NO_FREE);
			}

		} else if (node_type == MD_AST_NODE_TYPE_MOOSVAR_ASSIGNMENT) {
			mlr_dsl_cst_statement_item_t* pitem = pstatement->pitems->phead->pvvalue;

			lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;
			mv_t rhs_value = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars,
				pstate->pmoosvars, pregex_captures, pctx, prhs_evaluator->pvstate);

			sllmv_t* pmvkeys = sllmv_alloc();
			int ok = TRUE;
			for (sllve_t* pe = pitem->pmoosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
				lrec_evaluator_t* pmvkey_evaluator = pe->pvvalue;
				mv_t mvkey = pmvkey_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars,
					pstate->pmoosvars, pregex_captures, pctx, pmvkey_evaluator->pvstate);
				if (mv_is_null(&mvkey)) {
					ok = FALSE;
					printf("xxx temp stub NOT OK\n");
					break;
				}
				// xxx make this no-copy, or a no-copy variant ... some such.
				sllmv_add(pmvkeys, &mvkey);
				mv_free(&mvkey);
			}

			// xxx move this null-check into mlhmmv.c?
			if (ok)
				mlhmmv_put(pstate->pmoosvars, pmvkeys, &rhs_value);

			sllmv_free(pmvkeys);

		} else if (node_type == MD_AST_NODE_TYPE_EMIT) {
			lrec_t* prec_to_emit = lrec_unbacked_alloc();
			for (sllve_t* pf = pstatement->pitems->phead; pf != NULL; pf = pf->pnext) {
				mlr_dsl_cst_statement_item_t* pitem = pf->pvvalue;
				char* output_field_name = pitem->output_field_name;
				lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

				// xxx this is overkill ... the grammar allows only for oosvar names as args to emit.  so we could
				// bypass that and just hashmap-get keyed by output_field_name here.
				mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars, pstate->pmoosvars,
					pregex_captures, pctx, prhs_evaluator->pvstate);

				if (val.type == MT_STRING) {
					// Ownership transfer from (newly created) mlrval to (newly created) lrec.
					lrec_put(prec_to_emit, output_field_name, val.u.strv, val.free_flags);
				} else {
					char free_flags = NO_FREE;
					char* string = mv_format_val(&val, &free_flags);
					lrec_put(prec_to_emit, output_field_name, string, free_flags);
				}
			}
			sllv_add(poutrecs, prec_to_emit);

		} else if (node_type == MD_AST_NODE_TYPE_DUMP) {
			mlhmmv_print(pstate->pmoosvars);

		} else if (node_type == MD_AST_NODE_TYPE_FILTER) {
			mlr_dsl_cst_statement_item_t* pitem = pstatement->pitems->phead->pvvalue;
			lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

			mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars, pstate->pmoosvars,
				pregex_captures, pctx, prhs_evaluator->pvstate);
			if (val.type != MT_NULL) {
				mv_set_boolean_strict(&val);
				if (!val.u.boolv) {
					*pemit_rec = FALSE;
					break;
				}
			}

		} else if (node_type == MD_AST_NODE_TYPE_GATE) {
			mlr_dsl_cst_statement_item_t* pitem = pstatement->pitems->phead->pvvalue;
			lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

			mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars, pstate->pmoosvars,
				pregex_captures, pctx, prhs_evaluator->pvstate);
			if (val.type == MT_NULL)
				break;
			mv_set_boolean_strict(&val);
			if (!val.u.boolv) {
				break;
			}

		} else { // Bare-boolean statement
			mlr_dsl_cst_statement_item_t* pitem = pstatement->pitems->phead->pvvalue;
			lrec_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

			mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, pstate->poosvars, pstate->pmoosvars,
				pregex_captures, pctx, prhs_evaluator->pvstate);
			if (val.type != MT_NULL)
				mv_set_boolean_strict(&val);
		}
	}
}
