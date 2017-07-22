#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"
#include "cli/mlrcli.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmsv.h"
#include "containers/mlhmmv.h"
#include "parsing/mlr_dsl_wrapper.h"
#include "dsl/rval_evaluators.h"
#include "dsl/mlr_dsl_cst.h"
#include "mapping/mappers.h"

#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"

// ----------------------------------------------------------------
typedef struct _mapper_put_or_filter_state_t {
	char*          mlr_dsl_expression;

	mlr_dsl_ast_t* past;
	mlr_dsl_cst_t* pcst;

	int            at_begin;
	mlhmmv_root_t* poosvars;

	int            trace_execution;
	char*          oosvar_flatten_separator;
	int            flush_every_record;
	cli_writer_opts_t* pwriter_opts;

	local_stack_t* plocal_stack;
	loop_stack_t*  ploop_stack;

	int            put_output_disabled; // mlr put -q
	int            do_final_filter;     // mlr filter
	int            negate_final_filter; // mlr filter -x
} mapper_put_or_filter_state_t;

typedef struct _expression_info_t {
	char* filename;
	char* expression;
} expression_info_t;

// ----------------------------------------------------------------
static void      mapper_put_usage(FILE* o, char* argv0, char* verb);
static void   mapper_filter_usage(FILE* o, char* argv0, char* verb);
static void          shared_usage(FILE* o, char* argv0, char* verb);

static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);

static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);

static mapper_t* shared_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);

static mapper_t* mapper_put_or_filter_alloc(
	char*              mlr_dsl_expression,
	int                print_ast,
	int                trace_stack_allocation,
	int                trace_execution,
	mlr_dsl_ast_t*     past,
	int                put_output_disabled, // mlr put -q
	int                do_final_filter,     // mlr filter
	int                negate_final_filter, // mlr filter -x
	int                type_inferencing,
	char*              oosvar_flatten_separator,
	int                flush_every_record,
	cli_writer_opts_t* pwriter_opts,
	cli_writer_opts_t* pmain_writer_opts);

static void      mapper_put_or_filter_free(mapper_t* pmapper, context_t* pctx);

static sllv_t*   mapper_put_or_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_put_setup = {
	.verb = "put",
	.pusage_func = mapper_put_usage,
	.pparse_func = mapper_put_parse_cli,
	.ignores_input = FALSE,
};

mapper_setup_t mapper_filter_setup = {
	.verb = "filter",
	.pusage_func = mapper_filter_usage,
	.pparse_func = mapper_filter_parse_cli,
	.ignores_input = FALSE,
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

	shared_usage(o, argv0, verb);

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
	fprintf(o, "  %s %s -q '@sum += $x; end {emit @sum}'\n", argv0, verb);
	fprintf(o, "  %s %s -q '@sum[$a] += $x; end {emit @sum, \"a\"}'\n", argv0, verb);
	fprintf(o, "  %s %s -q '@sum[$a][$b] += $x; end {emit @sum, \"a\", \"b\"}'\n", argv0, verb);
	fprintf(o, "  %s %s -q '@min=min(@min,$x);@max=max(@max,$x); end{emitf @min, @max}'\n", argv0, verb);
	fprintf(o, "  %s %s -q 'is_null(@xmax) || $x > @xmax {@xmax=$x; @recmax=$*}; end {emit @recmax}'\n", argv0, verb);
	fprintf(o, "  %s %s '\n", argv0, verb);
	fprintf(o, "    $x = 1;\n");
	fprintf(o, "   #$y = 2;\n");
	fprintf(o, "    $z = 3\n");
	fprintf(o, "  '\n");
	fprintf(o, "\n");
	fprintf(o, "Please see also '%s -k' for examples using redirected output.\n", argv0);
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\".\n", argv0);
	fprintf(o, "Please see in particular:\n");
	fprintf(o, "  http://www.johnkerl.org/miller/doc/reference.html#put\n");
}

// ----------------------------------------------------------------
static void mapper_filter_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {expression}\n", argv0, verb);
	fprintf(o, "Prints records for which {expression} evaluates to true.\n");
	fprintf(o, "If there are multiple semicolon-delimited expressions, all of them are\n");
	fprintf(o, "evaluated and the last one is used as the filter criterion.\n");
	fprintf(o, "\n");

	shared_usage(o, argv0, verb);

	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s 'log10($count) > 4.0'\n", argv0, verb);
	fprintf(o, "  %s %s 'FNR == 2          (second record in each file)'\n", argv0, verb);
	fprintf(o, "  %s %s 'urand() < 0.001'  (subsampling)\n", argv0, verb);
	fprintf(o, "  %s %s '$color != \"blue\" && $value > 4.2'\n", argv0, verb);
	fprintf(o, "  %s %s '($x<.5 && $y<.5) || ($x>.5 && $y>.5)'\n", argv0, verb);
	fprintf(o, "  %s %s '($name =~ \"^sys.*east$\") || ($name =~ \"^dev.[0-9]+\"i)'\n", argv0, verb);
	fprintf(o, "  %s %s '$ab = $a+$b; $cd = $c+$d; $ab != $cd'\n", argv0, verb);
	fprintf(o, "  %s %s '\n", argv0, verb);
	fprintf(o, "    NR == 1 ||\n");
	fprintf(o, "   #NR == 2 ||\n");
	fprintf(o, "    NR == 3\n");
	fprintf(o, "  '\n");
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including function list. Or \"%s -f\". Please also also \"%s grep\" which is\n", argv0, argv0);
	fprintf(o, "useful when you don't yet know which field name(s) you're looking for.\n");
}

static void shared_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "All field values are type-inferred to int/float/string unless this behavior is\n");
	fprintf(o, "suppressed with -S or -F (details below).\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-v: Prints the expressions's AST (abstract syntax tree), which gives\n");
	fprintf(o, "    full transparency on the precedence and associativity rules of\n");
	fprintf(o, "    Miller's grammar, to stdout.\n");
	fprintf(o, "-a: Prints a low-level stack-allocation trace to stdout.\n");
	fprintf(o, "-t: Prints a low-level parser trace to stderr.\n");
	fprintf(o, "-T: Prints a every statement to stderr as it is executed.\n");
	if (streq(verb, "put")) {
		fprintf(o, "-q: Does not include the modified record in the output stream. Useful for when\n");
		fprintf(o, "    all desired output is in begin and/or end blocks.\n");
	}
	if (streq(verb, "filter")) {
		fprintf(o, "-x: Prints records for which {expression} evaluates to false.\n");
	}
	fprintf(o, "-S: Keeps field values as strings with no type inference to int or float.\n");
	fprintf(o, "-F: Keeps field values as strings or floats with no inference to int.\n");
	fprintf(o, "--oflatsep {string}: Separator to use when flattening multi-level @-variables\n");
	fprintf(o, "    to output records for emit. Default \"%s\".\n", DEFAULT_OOSVAR_FLATTEN_SEPARATOR);
	fprintf(o, "--jknquoteint: For dump output (JSON-formatted), do not quote map keys if non-string.\n");
	fprintf(o, "--jvquoteall: For dump output (JSON-formatted), quote map values even if non-string.\n");
	fprintf(o, "-f {filename}: the DSL expression is taken from the specified file rather\n");
	fprintf(o, "    than from the command line. Outer single quotes wrapping the expression\n");
	fprintf(o, "    should not be placed in the file. If -f is specified more than once,\n");
	fprintf(o, "    all input files specified using -f are concatenated to produce the expression.\n");
	fprintf(o, "    (For example, you can define functions in one file and call them from another.)\n");
	fprintf(o, "-e {expression}: You can use this after -f to add an expression. Example use\n");
	fprintf(o, "    case: define functions/subroutines in a file you specify with -f, then call\n");
	fprintf(o, "    them with an expression you specify with -e.\n");
	fprintf(o, "(If you mix -e and -f then the expressions are evaluated in the order encountered.\n");
	fprintf(o, "Since the expression pieces are simply concatenated, please be sure to use intervening\n");
	fprintf(o, "semicolons to separate expressions.)\n");
	fprintf(o, "--no-fflush: for emit, tee, print, and dump, don't call fflush() after every\n");
	fprintf(o, "    record.\n");
	fprintf(o, "Any of the output-format command-line flags (see %s -h). Example: using\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "  %s --icsv --opprint ... then put --ojson 'tee > \"mytap-\".$a.\".dat\", $*' then ...\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "the input is CSV, the output is pretty-print tabular, but the tee-file output\n");
	fprintf(o, "is written in JSON format.\n");
	fprintf(o, "\n");
	fprintf(o, "Please use a dollar sign for field names and double-quotes for string\n");
	fprintf(o, "literals. If field names have special characters such as \".\" then you might\n");
	fprintf(o, "use braces, e.g. '${field.name}'. Miller built-in variables are\n");
	fprintf(o, "NF NR FNR FILENUM FILENAME PI E, and ENV[\"namegoeshere\"] to access environment\n");
	fprintf(o, "variables. The environment-variable name may be an expression, e.g. a field\n");
	fprintf(o, "value.\n");
	fprintf(o, "\n");
	fprintf(o, "Use # to comment to end of line.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_filter_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* pmain_writer_opts)
{
	return shared_parse_cli(pargi, argc, argv, _, pmain_writer_opts);
}

static mapper_t* mapper_put_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* pmain_writer_opts)
{
	return shared_parse_cli(pargi, argc, argv, _, pmain_writer_opts);
}

static mapper_t* shared_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* pmain_writer_opts)
{
	sllv_t* expression_infos         = sllv_alloc();
	char*   mlr_dsl_expression       = NULL;
	int     put_output_disabled      = FALSE;
	int     do_final_filter          = FALSE;
	int     negate_final_filter      = FALSE;
	int     type_inferencing         = TYPE_INFER_STRING_FLOAT_INT;
	int     print_ast                = FALSE;
	int     trace_stack_allocation   = FALSE;
	int     trace_parse              = FALSE;
	int     trace_execution          = FALSE;
	char*   oosvar_flatten_separator = DEFAULT_OOSVAR_FLATTEN_SEPARATOR;
	int     flush_every_record       = TRUE;

	cli_writer_opts_t* pwriter_opts = mlr_malloc_or_die(sizeof(cli_writer_opts_t));
	cli_writer_opts_init(pwriter_opts);

	int argi = *pargi;
	if ((argc - argi) < 1) {
		mapper_put_usage(stderr, argv[0], argv[argi]);
		return NULL;
	}
	char* verb = argv[argi++];

	if (streq(verb, "filter"))
		do_final_filter = TRUE;

	cli_writer_opts_init(pwriter_opts);
	for (; argi < argc; /* variable increment: 1 or 2 depending on flag */) {

		if (argv[argi][0] != '-') {
			break; // No more flag options to process

		} else if (cli_handle_writer_options(argv, argc, &argi, pwriter_opts)) {
			// handled

		} else if (streq(argv[argi], "-f")) {
			if ((argc - argi) < 2) {
				mapper_put_usage(stderr, argv[0], verb);
				return NULL;
			}
			expression_info_t* pexpression_info = mlr_malloc_or_die(sizeof(expression_info_t));
			pexpression_info->filename = argv[argi+1];
			pexpression_info->expression = NULL;
			sllv_append(expression_infos, pexpression_info);
			argi += 2;

		} else if (streq(argv[argi], "-e")) {
			if ((argc - argi) < 2) {
				mapper_put_usage(stderr, argv[0], verb);
				return NULL;
			}
			expression_info_t* pexpression_info = mlr_malloc_or_die(sizeof(expression_info_t));
			pexpression_info->filename = NULL;
			pexpression_info->expression = argv[argi+1];
			sllv_append(expression_infos, pexpression_info);
			argi += 2;

		} else if (streq(argv[argi], "-v")) {
			print_ast = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-a")) {
			trace_stack_allocation = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-t")) {
			trace_parse = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-T")) {
			trace_execution = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-q") && streq(verb, "put")) {

			put_output_disabled = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-x") && streq(verb, "filter")) {
			negate_final_filter = TRUE;
			argi += 1;

		} else if (streq(argv[argi], "-S")) {
			type_inferencing = TYPE_INFER_STRING_ONLY;
			argi += 1;
		} else if (streq(argv[argi], "-F")) {
			type_inferencing = TYPE_INFER_STRING_FLOAT;
			argi += 1;
		} else if (streq(argv[argi], "--oflatsep")) {
			if ((argc - argi) < 2) {
				mapper_put_usage(stderr, argv[0], verb);
				return NULL;
			}
			oosvar_flatten_separator = argv[argi+1];
			argi += 2;
		} else if (streq(argv[argi], "--no-fflush") || streq(argv[argi], "--no-flush")) {
			flush_every_record = FALSE;
			argi += 1;

		} else {
			mapper_put_usage(stderr, argv[0], verb);
			return NULL;
		}
	}

	if (expression_infos->length == 0) {
		if ((argc - argi) < 1) {
			mapper_put_usage(stderr, argv[0], verb);
			return NULL;
		}
		mlr_dsl_expression = mlr_strdup_or_die(argv[argi++]);
	} else {
		string_builder_t *psb = sb_alloc(1024);

		for (sllve_t* pe = expression_infos->phead; pe != NULL; pe = pe->pnext) {
			expression_info_t* pexpression_info = pe->pvvalue;
			if (pexpression_info->filename != NULL) {
				char* mlr_dsl_expression_piece = read_file_into_memory(pexpression_info->filename, NULL);
				sb_append_string(psb, mlr_dsl_expression_piece);
				free(mlr_dsl_expression_piece);
			} else {
				sb_append_string(psb, pexpression_info->expression);
			}
			free(pexpression_info);
		}

		mlr_dsl_expression = sb_finish(psb);
		sb_free(psb);
	}
	sllv_free(expression_infos);

	// On Linux, '$x=4' is correct; single-quotes inhibit shell-variable expansion which is desired
	// since Miller variables also start with $. In either case, the shell removes the outermost
	// quotes before passing the string to us.  On Windows, though, double-quotes are the right
	// thing and are removed for us, while single-quotes are passed through to us.
	if (mlr_dsl_expression[0] == '\'') {
		int len = strlen(mlr_dsl_expression);
		if (len >= 2 && mlr_dsl_expression[len-1] == '\'') {
			mlr_dsl_expression[0] = ' ';
			mlr_dsl_expression[len-1] = 0;
		}
	}

	mlr_dsl_ast_t* past = mlr_dsl_parse(mlr_dsl_expression, trace_parse);
	if (past == NULL) {
		fprintf(stderr, "%s %s: syntax error on DSL parse of '%s'\n",
			MLR_GLOBALS.bargv0, verb, mlr_dsl_expression);
		return NULL;
	}

	// For just dev-testing the parser, you can do
	//   mlr put -v 'expression goes here' /dev/null
	if (print_ast) {
		printf("RAW AST:\n");
		printf("\n");
		mlr_dsl_ast_print(past);
	}

	*pargi = argi;
	return mapper_put_or_filter_alloc(mlr_dsl_expression, print_ast, trace_stack_allocation, trace_execution,
		past, put_output_disabled, do_final_filter, negate_final_filter, type_inferencing, oosvar_flatten_separator,
			flush_every_record, pwriter_opts, pmain_writer_opts);
}

// ----------------------------------------------------------------
static mapper_t* mapper_put_or_filter_alloc(
	char*              mlr_dsl_expression,
	int                print_ast,
	int                trace_stack_allocation,
	int                trace_execution,
	mlr_dsl_ast_t*     past,
	int                put_output_disabled, // mlr put -q
	int                do_final_filter,     // mlr filter
	int                negate_final_filter, // mlr filter -x
	int                type_inferencing,
	char*              oosvar_flatten_separator,
	int                flush_every_record,
	cli_writer_opts_t* pwriter_opts,
	cli_writer_opts_t* pmain_writer_opts)
{
	mapper_put_or_filter_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_put_or_filter_state_t));
	// Retain the string contents along with any in-pointers from the AST/CST
	pstate->mlr_dsl_expression = mlr_dsl_expression;
	pstate->past                     = past;
	pstate->pcst                     = mlr_dsl_cst_alloc(past, print_ast, trace_stack_allocation,
		type_inferencing, flush_every_record, do_final_filter, negate_final_filter);
	pstate->at_begin                     = TRUE;
	pstate->put_output_disabled          = put_output_disabled;
	pstate->poosvars                     = mlhmmv_root_alloc();
	pstate->trace_execution              = trace_execution;
	pstate->oosvar_flatten_separator     = oosvar_flatten_separator;
	pstate->flush_every_record           = flush_every_record;
	pstate->plocal_stack                 = local_stack_alloc();
	pstate->ploop_stack                  = loop_stack_alloc();
	pstate->pwriter_opts                 = pwriter_opts;

	cli_merge_writer_opts(pstate->pwriter_opts, pmain_writer_opts);

	mapper_t* pmapper      = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_put_or_filter_process;
	pmapper->pfree_func    = mapper_put_or_filter_free;

	return pmapper;
}

static void mapper_put_or_filter_free(mapper_t* pmapper, context_t* pctx) {
	mapper_put_or_filter_state_t* pstate = pmapper->pvstate;

	free(pstate->mlr_dsl_expression);
	mlhmmv_root_free(pstate->poosvars);
	local_stack_free(pstate->plocal_stack);
	loop_stack_free(pstate->ploop_stack);
	mlr_dsl_cst_free(pstate->pcst, pctx);
	// Free what's left of the stripped AST after the CST reorganized it.
	mlr_dsl_ast_free(pstate->past);

	free(pstate->pwriter_opts);
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

static sllv_t* mapper_put_or_filter_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_put_or_filter_state_t* pstate = (mapper_put_or_filter_state_t*)pvstate;

	sllv_t* poutrecs = sllv_alloc();
	int should_emit_rec = TRUE;

	if (pstate->at_begin) {
		string_array_t* pregex_captures = NULL; // May be set to non-null on evaluation

		variables_t variables = (variables_t) {
			.pinrec           = NULL,
			.ptyped_overlay   = NULL,
			.poosvars         = pstate->poosvars,
			.ppregex_captures = &pregex_captures,
			.pctx             = pctx,
			.plocal_stack     = pstate->plocal_stack,
			.ploop_stack      = pstate->ploop_stack,
			.return_state = {
				.returned = FALSE,
				.retval = box_ephemeral_val(mv_absent()),
			},
			.trace_execution              = pstate->trace_execution,
			.json_quote_int_keys          = pstate->pwriter_opts->json_quote_int_keys,
			.json_quote_non_string_values = pstate->pwriter_opts->json_quote_non_string_values,
		};
		cst_outputs_t cst_outputs = (cst_outputs_t) {
			.pshould_emit_rec             = &should_emit_rec,
			.poutrecs                     = poutrecs,
			.oosvar_flatten_separator     = pstate->oosvar_flatten_separator,
			.pwriter_opts                 = pstate->pwriter_opts,
		};

		string_array_free(pregex_captures);
		mlr_dsl_cst_handle_top_level_statement_blocks(pstate->pcst->pbegin_blocks, &variables, &cst_outputs);
		pstate->at_begin = FALSE;
	}

	if (pinrec == NULL) { // End of input stream
		string_array_t* pregex_captures = NULL; // May be set to non-null on evaluation

		variables_t variables = (variables_t) {
			.pinrec           = NULL,
			.ptyped_overlay   = NULL,
			.poosvars         = pstate->poosvars,
			.ppregex_captures = &pregex_captures,
			.pctx             = pctx,
			.plocal_stack     = pstate->plocal_stack,
			.ploop_stack      = pstate->ploop_stack,
			.return_state = {
				.returned = FALSE,
				.retval = box_ephemeral_val(mv_absent()),
			},
			.trace_execution              = pstate->trace_execution,
			.json_quote_int_keys          = pstate->pwriter_opts->json_quote_int_keys,
			.json_quote_non_string_values = pstate->pwriter_opts->json_quote_non_string_values,
		};
		cst_outputs_t cst_outputs = (cst_outputs_t) {
			.pshould_emit_rec             = &should_emit_rec,
			.poutrecs                     = poutrecs,
			.oosvar_flatten_separator     = pstate->oosvar_flatten_separator,
			.pwriter_opts                 = pstate->pwriter_opts,
		};

		mlr_dsl_cst_handle_top_level_statement_blocks(pstate->pcst->pend_blocks, &variables, &cst_outputs);

		string_array_free(pregex_captures);
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}

	lhmsmv_t* ptyped_overlay = lhmsmv_alloc();
	string_array_t* pregex_captures = NULL; // May be set to non-null on evaluation

	should_emit_rec = TRUE;

	variables_t variables = (variables_t) {
		.pinrec           = pinrec, // Note variables.pinrec pointer can update on '$* = ...'
		.ptyped_overlay   = ptyped_overlay,
		.poosvars         = pstate->poosvars,
		.ppregex_captures = &pregex_captures,
		.pctx             = pctx,
		.plocal_stack     = pstate->plocal_stack,
		.ploop_stack      = pstate->ploop_stack,
		.return_state = {
			.returned = FALSE,
			.retval = box_ephemeral_val(mv_absent()),
		},
		.trace_execution              = pstate->trace_execution,
		.json_quote_int_keys          = pstate->pwriter_opts->json_quote_int_keys,
		.json_quote_non_string_values = pstate->pwriter_opts->json_quote_non_string_values,
	};
	cst_outputs_t cst_outputs = (cst_outputs_t) {
		.pshould_emit_rec             = &should_emit_rec,
		.poutrecs                     = poutrecs,
		.oosvar_flatten_separator     = pstate->oosvar_flatten_separator,
		.pwriter_opts                 = pstate->pwriter_opts,
	};

	mlr_dsl_cst_handle_top_level_statement_block(pstate->pcst->pmain_block, &variables, &cst_outputs);

	if (should_emit_rec && !pstate->put_output_disabled) {
		// Write the output fields from the typed overlay back to the lrec.
		for (lhmsmve_t* pe = variables.ptyped_overlay->phead; pe != NULL; pe = pe->pnext) {
			char* output_field_name = pe->key;
			mv_t* pval = &pe->value;

			// Ownership transfer from mv_t to lrec.
			if (pval->type == MT_STRING) {
				lrec_put(variables.pinrec, output_field_name, pval->u.strv, pval->free_flags);
			} else {
				char free_flags = NO_FREE;
				char* string = mv_format_val(pval, &free_flags);
				lrec_put(variables.pinrec, output_field_name, string, pval->free_flags | free_flags);
			}
			pval->free_flags = NO_FREE;
		}
	}
	lhmsmv_free(variables.ptyped_overlay);
	string_array_free(pregex_captures);

	// Note variables.pinrec pointer can update on '$* = ...'
	if (should_emit_rec && !pstate->put_output_disabled) {
		sllv_append(poutrecs, variables.pinrec);
	} else {
		lrec_free(variables.pinrec);
	}
	return poutrecs;
}
