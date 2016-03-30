#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "cli/argparse.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmsv.h"
#include "containers/mlhmmv.h"
#include "mapping/mappers.h"
#include "mapping/rval_evaluators.h"
#include "dsls/mlr_dsl_wrapper.h"
#include "mlr_dsl_cst.h"

#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"

typedef struct _mapper_put_state_t {
	ap_state_t*    pargp;
	char*          mlr_dsl_expression;
	mlr_dsl_ast_t* past;
	mlr_dsl_cst_t* pcst;
	int            at_begin;
	mlhmmv_t*      poosvars;
	char*          oosvar_flatten_separator;
	int            outer_filter;
} mapper_put_state_t;

static void      mapper_put_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_put_alloc(ap_state_t* pargp, char* mlr_dsl_expression, mlr_dsl_ast_t* past,
	int outer_filter, int type_inferencing, char* oosvar_flatten_separator);
static void      mapper_put_free(mapper_t* pmapper);
static sllv_t*   mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

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
	fprintf(o, "either be assignments, or evaluate to boolean.  Booleans with following\n");
	fprintf(o, "statements in curly braces control whether those statements are executed;\n");
	fprintf(o, "booleans without following curly braces do nothing except side effects (e.g.\n");
	fprintf(o, "regex-captures into \\1, \\2, etc.).\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-v: First prints the AST (abstract syntax tree) for the expression, which gives\n");
	fprintf(o, "    full transparency on the precedence and associativity rules of Miller's\n");
	fprintf(o, "    grammar.\n");
	fprintf(o, "-q: Does not include the modified record in the output stream. Useful for when\n");
	fprintf(o, "    all desired output is in begin and/or end blocks.\n");
	fprintf(o, "-S: Keeps field values, or literals in the expression, as strings with no type \n");
	fprintf(o, "    inference to int or float.\n");
	fprintf(o, "-F: Keeps field values, or literals in the expression, as strings or floats\n");
	fprintf(o, "    with no inference to int.\n");
	fprintf(o, "\n");
	fprintf(o, "Please use a dollar sign for field names and double-quotes for string\n");
	fprintf(o, "literals. If field names have special characters such as \".\" then you might\n");
	fprintf(o, "use braces, e.g. '${field.name}'. Miller built-in variables are\n");
	fprintf(o, "NF NR FNR FILENUM FILENAME PI E, and ENV[\"namegoeshere\"] to access environment\n");
	fprintf(o, "variables. The environment-variable name may be an expression, e.g. a field value.\n");
	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s '$y = log10($x); $z = sqrt($y)'\n", argv0, verb);
	fprintf(o, "  %s %s '$x>0.0 { $y=log10($x); $z=sqrt($y) }' # does {...} only if $x > 0.0\n", argv0, verb);
	fprintf(o, "  %s %s '$x>0.0;  $y=log10($x); $z=sqrt($y)'   # does all three statements\n", argv0, verb);
	fprintf(o, "  %s %s '$a =~ \"([a-z]+)_([0-9]+);  $b = \"left_\\1\"; $c = \"right_\\2\"'\n", argv0, verb);
	fprintf(o, "  %s %s '$a =~ \"([a-z]+)_([0-9]+) { $b = \"left_\\1\"; $c = \"right_\\2\" }'\n", argv0, verb);
	fprintf(o, "  %s %s '$filename = FILENAME'\n", argv0, verb);
	fprintf(o, "  %s %s '$colored_shape = $color . \"_\" . $shape'\n", argv0, verb);
	fprintf(o, "  %s %s '$y = cos($theta); $z = atan2($y, $x)'\n", argv0, verb);
	fprintf(o, "  %s %s '$name = sub($name, \"http.*com\"i, \"\")'\n", argv0, verb);
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\".\n", argv0);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv) {
	char* verb                = argv[(*pargi)++];
	char* mlr_dsl_expression  = NULL;
	char* expression_filename = NULL;
	int   outer_filter        = TRUE;
	int   type_inferencing    = TYPE_INFER_STRING_FLOAT_INT;
	int   print_ast           = FALSE;
	char* oosvar_flatten_separator = DEFAULT_OOSVAR_FLATTEN_SEPARATOR;

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate,    "-f", &expression_filename);
	ap_define_true_flag(pstate,      "-v", &print_ast);
	ap_define_false_flag(pstate,     "-q", &outer_filter);
	ap_define_int_value_flag(pstate, "-S", TYPE_INFER_STRING_ONLY,  &type_inferencing);
	ap_define_int_value_flag(pstate, "-F", TYPE_INFER_STRING_FLOAT, &type_inferencing);
	// xxx to online help
	ap_define_string_flag(pstate,    "--oflatsep", &oosvar_flatten_separator);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_put_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (expression_filename == NULL) {
		if ((argc - *pargi) < 1) {
			mapper_put_usage(stderr, argv[0], verb);
			return NULL;
		}
		mlr_dsl_expression = mlr_strdup_or_die(argv[(*pargi)++]);
	} else {
		mlr_dsl_expression = read_file_into_memory(expression_filename, NULL);
	}

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

	return mapper_put_alloc(pstate, mlr_dsl_expression, past, outer_filter, type_inferencing,
		oosvar_flatten_separator);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_alloc(ap_state_t* pargp, char* mlr_dsl_expression, mlr_dsl_ast_t* past,
	int outer_filter, int type_inferencing, char* oosvar_flatten_separator)
{
	mapper_put_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_put_state_t));
	pstate->pargp        = pargp;
	// Retain the string contents along with any in-pointers from the AST/CST
	pstate->mlr_dsl_expression = mlr_dsl_expression;
	pstate->past         = past;
	pstate->pcst         = mlr_dsl_cst_alloc(past, type_inferencing);
	pstate->at_begin     = TRUE;
	pstate->outer_filter = outer_filter;
	pstate->poosvars     = mlhmmv_alloc();
	pstate->oosvar_flatten_separator = oosvar_flatten_separator;

	mapper_t* pmapper      = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_put_process;
	pmapper->pfree_func    = mapper_put_free;

	return pmapper;
}

static void mapper_put_free(mapper_t* pmapper) {
	mapper_put_state_t* pstate = pmapper->pvstate;

	free(pstate->mlr_dsl_expression);
	mlhmmv_free(pstate->poosvars);
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
//
// * It is set here to null.
//
// * It is passed by reference to the lrec-evaluator tree. In particular, the matches and does-not-match functions
//   (which implement the =~ and !=~ operators) allocate it (or resize it, as necessary) and populate it.
//
// * If the matches/does-not-match functions are entered, even with no matches, the regex-captures string-array
//   will be resized to have length 0.
//
// * When the lrec-evaluator's from-literal function is invoked, the interpolate_regex_captures function can quickly
//   check to see if the regex-captures array is null and thereby know that a time-consuming scan for \1, \2, \3, etc.
//   does not need to be done. On the other hand, if the regex-captures array is non-null and has length
//   zero, then \0 .. \9 should all be replaced with the empty string.
//
// ----------------------------------------------------------------
// The oosvars multi-level hashmap contains out-of-stream variables which can be written/read/output
// in begin{}/end{} blocks, and/or per record.
//
// ----------------------------------------------------------------
// The context_t contains information about record-number, file-number, file-name, etc. from
// which the current stream-record was obtained.
// ----------------------------------------------------------------

static sllv_t* mapper_put_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_put_state_t* pstate = (mapper_put_state_t*)pvstate;

	string_array_t* pregex_captures = NULL; // May be set to non-null on evaluation
	sllv_t* poutrecs = sllv_alloc();
	int emit_rec = TRUE;

	if (pstate->at_begin) {
		mlr_dsl_cst_evaluate(pstate->pcst->pbegin_statements,
			pstate->poosvars, NULL, NULL, &pregex_captures, pctx, &emit_rec, poutrecs,
				pstate->oosvar_flatten_separator);
		pstate->at_begin = FALSE;
	}

	if (pinrec == NULL) { // End of input stream
		mlr_dsl_cst_evaluate(pstate->pcst->pend_statements,
			pstate->poosvars, NULL, NULL, &pregex_captures, pctx, &emit_rec, poutrecs,
				pstate->oosvar_flatten_separator);

		string_array_free(pregex_captures);
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}

	lhmsv_t* ptyped_overlay = lhmsv_alloc();

	emit_rec = TRUE;
	mlr_dsl_cst_evaluate(pstate->pcst->pmain_statements,
		pstate->poosvars, pinrec, ptyped_overlay, &pregex_captures, pctx, &emit_rec, poutrecs,
			pstate->oosvar_flatten_separator);

	if (emit_rec && pstate->outer_filter) {
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

	if (emit_rec && pstate->outer_filter) {
		sllv_append(poutrecs, pinrec);
	} else {
		lrec_free(pinrec);
	}
	return poutrecs;
}
