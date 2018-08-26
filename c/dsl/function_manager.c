#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "dsl/function_manager.h"
#include "dsl/context_flags.h"
#include "dsl/rval_evaluators.h"
#include "dsl/rxval_evaluators.h"

// ----------------------------------------------------------------
typedef enum _func_class_t {
	FUNC_CLASS_ARITHMETIC,
	FUNC_CLASS_MATH,
	FUNC_CLASS_BOOLEAN,
	FUNC_CLASS_STRING,
	FUNC_CLASS_CONVERSION,
	FUNC_CLASS_TYPING,
	FUNC_CLASS_MAPS,
	FUNC_CLASS_TIME
} func_class_t;

typedef enum _arity_check_t {
	ARITY_CHECK_PASS,
	ARITY_CHECK_FAIL,
	ARITY_CHECK_NO_SUCH
} arity_check_t;

typedef struct _function_lookup_t {
	func_class_t function_class;
	char*        function_name;
	int          arity; // for variadic, this is minimum arity
	int          variadic;
	char*        usage_string;
} function_lookup_t;

// This is shared between all instances
static function_lookup_t FUNCTION_LOOKUP_TABLE[];

// ----------------------------------------------------------------
// See also comments in rval_evaluators.h

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
static void fmgr_check_arity_with_report(fmgr_t* pfmgr, char* function_name,
	int user_provided_arity, int* pvariadic);

static rval_evaluator_t* fmgr_alloc_evaluator_from_variadic_func_name(
	char* function_name, rval_evaluator_t** pargs, int nargs);

static rval_evaluator_t* fmgr_alloc_evaluator_from_zary_func_name(
	char* function_name);

static rval_evaluator_t* fmgr_alloc_evaluator_from_unary_func_name(
	char* function_name, rval_evaluator_t* parg1);

static rval_evaluator_t* fmgr_alloc_evaluator_from_binary_func_name(
	char* function_name,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);

static rval_evaluator_t* fmgr_alloc_evaluator_from_binary_regex_arg2_func_name(
	char* function_name,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case);

static rval_evaluator_t* fmgr_alloc_evaluator_from_ternary_func_name(
	char* function_name,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);

static rval_evaluator_t* fmgr_alloc_evaluator_from_ternary_regex_arg2_func_name(
	char* function_name,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3);

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// For rval functions, we pass rval_evaluator_t* (CST); for rxval functions, we pass
// mlr_dsl_ast_node_t* (AST). It's easy to construct the former from the latter, of
// course. The difference is that we look up map-enabled functions by name first,
// then non-map-enabled functions by name second.
//
// * AST nodes are passed to try to look up a map-enabled function given a function name.
// * If those exist, they construct CST structures and return.
// * But if not, we look up a non-map-enabled function for the same function name.
// * If that doesn't exist either, then it's a fatal error. So we go ahead and
//   construct an rval_evaluator_t* CST structure from the AST node simply to
//   save keystrokes, passing that to the function-lookup routines.
//
// It would simpler to always construct CST structures before looking up
// function names, but the only problem is that it's hard to unconstruct CST
// structures in case the name lookup fails. (The function-manager
// as-yet-unresolved-name list points into them, whenever function arguments
// themselves include function calls). Namely, the following scenario is to be
// avoided:
//
// * Construct rxval_evaluator_t* CST structure.
// * Look up map-enabled function with a given name.
// * That doesn't exist.
// * Now the rxval_evaluator_t* can't be torn down since the fmgr points into it.

static rxval_evaluator_t* fmgr_alloc_xevaluator_from_variadic_func_name(
	char* function_name, sllv_t* parg_nodes,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/);

static rxval_evaluator_t* fmgr_alloc_xevaluator_from_unary_func_name(
	char* function_name,
	mlr_dsl_ast_node_t* parg1,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/);

static rxval_evaluator_t* fmgr_alloc_xevaluator_from_binary_func_name(
	char* function_name,
	mlr_dsl_ast_node_t* parg1, mlr_dsl_ast_node_t* pargs2,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/);

static rxval_evaluator_t* fmgr_alloc_xevaluator_from_ternary_func_name(
	char* function_name,
	mlr_dsl_ast_node_t* parg1, mlr_dsl_ast_node_t* pargs2, mlr_dsl_ast_node_t* pargs3,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/);

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
static void  resolve_func_callsite(fmgr_t* pfmgr, rval_evaluator_t*  pev);
static void resolve_func_xcallsite(fmgr_t* pfmgr, rxval_evaluator_t* pxev);
static rxval_evaluator_t* fmgr_alloc_xeval_wrapping_eval(rval_evaluator_t* pevaluator);
static rval_evaluator_t* fmgr_alloc_eval_wrapping_xeval(rxval_evaluator_t* pxevaluator);

// ----------------------------------------------------------------
fmgr_t* fmgr_alloc() {
	fmgr_t* pfmgr = mlr_malloc_or_die(sizeof(fmgr_t));

	pfmgr->function_lookup_table = &FUNCTION_LOOKUP_TABLE[0];

	pfmgr->built_in_function_names = hss_alloc();
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &pfmgr->function_lookup_table[i];
		char* fname = plookup->function_name;
		if (fname == NULL)
			break;
		hss_add(pfmgr->built_in_function_names, fname);
	}

	pfmgr->pudf_names_to_defsite_states = lhmsv_alloc();

	pfmgr->pfunc_callsite_evaluators_to_resolve  = sllv_alloc();
	pfmgr->pfunc_callsite_xevaluators_to_resolve = sllv_alloc();

	return pfmgr;
}

// ----------------------------------------------------------------
void fmgr_free(fmgr_t* pfmgr, context_t* pctx) {
	if (pfmgr == NULL)
		return;

	for (lhmsve_t* pe = pfmgr->pudf_names_to_defsite_states->phead; pe != NULL; pe = pe->pnext) {
		udf_defsite_state_t * pdefsite_state = pe->pvvalue;
		free(pdefsite_state->name);
		pdefsite_state->pfree_func(pdefsite_state->pvstate, pctx);
		free(pdefsite_state);
	}
	lhmsv_free(pfmgr->pudf_names_to_defsite_states);
	sllv_free(pfmgr->pfunc_callsite_evaluators_to_resolve);
	sllv_free(pfmgr->pfunc_callsite_xevaluators_to_resolve);
	hss_free(pfmgr->built_in_function_names);
	free(pfmgr);
}

// ----------------------------------------------------------------
void fmgr_install_udf(fmgr_t* pfmgr, udf_defsite_state_t* pdefsite_state) {
	if (hss_has(pfmgr->built_in_function_names, pdefsite_state->name)) {
		fprintf(stderr, "%s: function named \"%s\" must not override a built-in function of the same name.\n",
			MLR_GLOBALS.bargv0, pdefsite_state->name);
		exit(1);
	}
	if (lhmsv_get(pfmgr->pudf_names_to_defsite_states, pdefsite_state->name)) {
		fprintf(stderr, "%s: function named \"%s\" has already been defined.\n",
			MLR_GLOBALS.bargv0, pdefsite_state->name);
		exit(1);
	}
	lhmsv_put(pfmgr->pudf_names_to_defsite_states, mlr_strdup_or_die(pdefsite_state->name), pdefsite_state,
		FREE_ENTRY_KEY);
}

// ================================================================
static function_lookup_t FUNCTION_LOOKUP_TABLE[] = {

	{FUNC_CLASS_ARITHMETIC, "+",  2,0, "Addition."},
	{FUNC_CLASS_ARITHMETIC, "+",  1,0, "Unary plus."},
	{FUNC_CLASS_ARITHMETIC, "-",  2,0, "Subtraction."},
	{FUNC_CLASS_ARITHMETIC, "-",  1,0, "Unary minus."},
	{FUNC_CLASS_ARITHMETIC, "*",  2,0, "Multiplication."},
	{FUNC_CLASS_ARITHMETIC, "/",  2,0, "Division."},
	{FUNC_CLASS_ARITHMETIC, "//", 2,0, "Integer division: rounds to negative (pythonic)."},

	{FUNC_CLASS_ARITHMETIC, ".+",  2,0, "Addition, with integer-to-integer overflow"},
	{FUNC_CLASS_ARITHMETIC, ".+",  1,0, "Unary plus, with integer-to-integer overflow."},
	{FUNC_CLASS_ARITHMETIC, ".-",  2,0, "Subtraction, with integer-to-integer overflow."},
	{FUNC_CLASS_ARITHMETIC, ".-",  1,0, "Unary minus, with integer-to-integer overflow."},
	{FUNC_CLASS_ARITHMETIC, ".*",  2,0, "Multiplication, with integer-to-integer overflow."},
	{FUNC_CLASS_ARITHMETIC, "./",  2,0, "Division, with integer-to-integer overflow."},
	{FUNC_CLASS_ARITHMETIC, ".//", 2,0, "Integer division: rounds to negative (pythonic), with integer-to-integer overflow."},

	{FUNC_CLASS_ARITHMETIC, "%",  2,0, "Remainder; never negative-valued (pythonic)."},
	{FUNC_CLASS_ARITHMETIC, "**", 2,0, "Exponentiation; same as pow, but as an infix\noperator."},
	{FUNC_CLASS_ARITHMETIC, "|",  2,0, "Bitwise OR."},
	{FUNC_CLASS_ARITHMETIC, "^",  2,0, "Bitwise XOR."},
	{FUNC_CLASS_ARITHMETIC, "&",  2,0, "Bitwise AND."},
	{FUNC_CLASS_ARITHMETIC, "~",  1,0,
		"Bitwise NOT. Beware '$y=~$x' since =~ is the\nregex-match operator: try '$y = ~$x'."},
	{FUNC_CLASS_ARITHMETIC, "<<", 2,0, "Bitwise left-shift."},
	{FUNC_CLASS_ARITHMETIC, ">>", 2,0, "Bitwise right-shift."},
	{FUNC_CLASS_ARITHMETIC, "bitcount",  1,0, "Count of 1-bits"},

	{FUNC_CLASS_BOOLEAN, "==",  2,0, "String/numeric equality. Mixing number and string\nresults in string compare."},
	{FUNC_CLASS_BOOLEAN, "!=",  2,0, "String/numeric inequality. Mixing number and string\nresults in string compare."},
	{FUNC_CLASS_BOOLEAN, "=~",  2,0,
		"String (left-hand side) matches regex (right-hand\n"
		"side), e.g. '$name =~ \"^a.*b$\"'."},
	{FUNC_CLASS_BOOLEAN, "!=~", 2,0,
		"String (left-hand side) does not match regex\n"
		"(right-hand side), e.g. '$name !=~ \"^a.*b$\"'."},
	{FUNC_CLASS_BOOLEAN, ">",   2,0,
		"String/numeric greater-than. Mixing number and string\n"
		"results in string compare."},
	{FUNC_CLASS_BOOLEAN, ">=",  2,0,
		"String/numeric greater-than-or-equals. Mixing number\n"
		"and string results in string compare."},
	{FUNC_CLASS_BOOLEAN, "<",   2,0,
		"String/numeric less-than. Mixing number and string\n"
		"results in string compare."},
	{FUNC_CLASS_BOOLEAN, "<=",  2,0,
		"String/numeric less-than-or-equals. Mixing number\n"
		"and string results in string compare."},
	{FUNC_CLASS_BOOLEAN, "&&",  2,0, "Logical AND."},
	{FUNC_CLASS_BOOLEAN, "||",  2,0, "Logical OR."},
	{FUNC_CLASS_BOOLEAN, "^^",  2,0, "Logical XOR."},
	{FUNC_CLASS_BOOLEAN, "!",   1,0, "Logical negation."},
	{FUNC_CLASS_BOOLEAN, "? :", 3,0, "Ternary operator."},

	{FUNC_CLASS_STRING, ".",        2,0, "String concatenation."},
	{FUNC_CLASS_STRING, "gsub",     3,0, "Example: '$name=gsub($name, \"old\", \"new\")'\n(replace all)."},
	{FUNC_CLASS_STRING, "regex_extract", 2,0, "Example: '$name=regex_extract($name, \"[A-Z]{3}[0-9]{2}\")'\n."},
	{FUNC_CLASS_STRING, "strlen",   1,0, "String length."},
	{FUNC_CLASS_STRING, "sub",      3,0, "Example: '$name=sub($name, \"old\", \"new\")'\n(replace once)."},
	{FUNC_CLASS_STRING, "ssub",     3,0, "Like sub but does no regexing. No characters are special."},
	{FUNC_CLASS_STRING, "substr",   3,0,
		"substr(s,m,n) gives substring of s from 0-up position m to n \n"
		"inclusive. Negative indices -len .. -1 alias to 0 .. len-1."},
	{FUNC_CLASS_STRING, "tolower",  1,0, "Convert string to lowercase."},
	{FUNC_CLASS_STRING, "toupper",  1,0, "Convert string to uppercase."},

	{FUNC_CLASS_MATH, "abs",      1,0, "Absolute value."},
	{FUNC_CLASS_MATH, "acos",     1,0, "Inverse trigonometric cosine."},
	{FUNC_CLASS_MATH, "acosh",    1,0, "Inverse hyperbolic cosine."},
	{FUNC_CLASS_MATH, "asin",     1,0, "Inverse trigonometric sine."},
	{FUNC_CLASS_MATH, "asinh",    1,0, "Inverse hyperbolic sine."},
	{FUNC_CLASS_MATH, "atan",     1,0, "One-argument arctangent."},
	{FUNC_CLASS_MATH, "atan2",    2,0, "Two-argument arctangent."},
	{FUNC_CLASS_MATH, "atanh",    1,0, "Inverse hyperbolic tangent."},
	{FUNC_CLASS_MATH, "cbrt",     1,0, "Cube root."},
	{FUNC_CLASS_MATH, "ceil",     1,0, "Ceiling: nearest integer at or above."},
	{FUNC_CLASS_MATH, "cos",      1,0, "Trigonometric cosine."},
	{FUNC_CLASS_MATH, "cosh",     1,0, "Hyperbolic cosine."},
	{FUNC_CLASS_MATH, "erf",      1,0, "Error function."},
	{FUNC_CLASS_MATH, "erfc",     1,0, "Complementary error function."},
	{FUNC_CLASS_MATH, "exp",      1,0, "Exponential function e**x."},
	{FUNC_CLASS_MATH, "expm1",    1,0, "e**x - 1."},
	{FUNC_CLASS_MATH, "floor",    1,0, "Floor: nearest integer at or below."},
	// See also http://johnkerl.org/doc/randuv.pdf for more about urand() -> other distributions
	{FUNC_CLASS_MATH, "invqnorm", 1,0,
		"Inverse of normal cumulative distribution\n"
		"function. Note that invqorm(urand()) is normally distributed."},
	{FUNC_CLASS_MATH, "log",      1,0, "Natural (base-e) logarithm."},
	{FUNC_CLASS_MATH, "log10",    1,0, "Base-10 logarithm."},
	{FUNC_CLASS_MATH, "log1p",    1,0, "log(1-x)."},
	{FUNC_CLASS_MATH, "logifit",  3,0, "Given m and b from logistic regression, compute\nfit: $yhat=logifit($x,$m,$b)."},
	{FUNC_CLASS_MATH, "madd",     3,0, "a + b mod m (integers)"},
	{FUNC_CLASS_MATH, "max",      0,1, "max of n numbers; null loses"},
	{FUNC_CLASS_MATH, "mexp",     3,0, "a ** b mod m (integers)"},
	{FUNC_CLASS_MATH, "min",      0,1, "Min of n numbers; null loses"},
	{FUNC_CLASS_MATH, "mmul",     3,0, "a * b mod m (integers)"},
	{FUNC_CLASS_MATH, "msub",     3,0, "a - b mod m (integers)"},
	{FUNC_CLASS_MATH, "pow",      2,0, "Exponentiation; same as **."},
	{FUNC_CLASS_MATH, "qnorm",    1,0, "Normal cumulative distribution function."},
	{FUNC_CLASS_MATH, "round",    1,0, "Round to nearest integer."},
	{FUNC_CLASS_MATH, "roundm",   2,0, "Round to nearest multiple of m: roundm($x,$m) is\nthe same as round($x/$m)*$m"},
	{FUNC_CLASS_MATH, "sgn",      1,0, "+1 for positive input, 0 for zero input, -1 for\nnegative input."},
	{FUNC_CLASS_MATH, "sin",      1,0, "Trigonometric sine."},
	{FUNC_CLASS_MATH, "sinh",     1,0, "Hyperbolic sine."},
	{FUNC_CLASS_MATH, "sqrt",     1,0, "Square root."},
	{FUNC_CLASS_MATH, "tan",      1,0, "Trigonometric tangent."},
	{FUNC_CLASS_MATH, "tanh",     1,0, "Hyperbolic tangent."},
	{FUNC_CLASS_MATH, "urand",    0,0,
		"Floating-point numbers on the unit interval.\n"
		"Int-valued example: '$n=floor(20+urand()*11)'." },
	{FUNC_CLASS_MATH, "urand32",  0,0, "Integer uniformly distributed 0 and 2**32-1\n"
	"inclusive." },
	{FUNC_CLASS_MATH, "urandint", 2,0, "Integer uniformly distributed between inclusive\ninteger endpoints." },

	{FUNC_CLASS_TIME, "dhms2fsec", 1,0,
		"Recovers floating-point seconds as in\n"
		"dhms2fsec(\"5d18h53m20.250000s\") = 500000.250000"},
	{FUNC_CLASS_TIME, "dhms2sec",  1,0, "Recovers integer seconds as in\ndhms2sec(\"5d18h53m20s\") = 500000"},
	{FUNC_CLASS_TIME, "fsec2dhms", 1,0,
		"Formats floating-point seconds as in\nfsec2dhms(500000.25) = \"5d18h53m20.250000s\""},
	{FUNC_CLASS_TIME, "fsec2hms",  1,0,
		"Formats floating-point seconds as in\nfsec2hms(5000.25) = \"01:23:20.250000\""},

	{FUNC_CLASS_TIME, "gmt2sec",   1,0, "Parses GMT timestamp as integer seconds since\nthe epoch."},
	{FUNC_CLASS_TIME, "localtime2sec", 1,0, "Parses local timestamp as integer seconds since\n"
		"the epoch. Consults $TZ environment variable."},

	{FUNC_CLASS_TIME, "hms2fsec",  1,0,
		"Recovers floating-point seconds as in\nhms2fsec(\"01:23:20.250000\") = 5000.250000"},
	{FUNC_CLASS_TIME, "hms2sec",   1,0, "Recovers integer seconds as in\nhms2sec(\"01:23:20\") = 5000"},
	{FUNC_CLASS_TIME, "sec2dhms",  1,0, "Formats integer seconds as in sec2dhms(500000)\n= \"5d18h53m20s\""},

	{FUNC_CLASS_TIME, "sec2gmt",   1,0,
		"Formats seconds since epoch (integer part)\n"
		"as GMT timestamp, e.g. sec2gmt(1440768801.7) = \"2015-08-28T13:33:21Z\".\n"
		"Leaves non-numbers as-is."},
	{FUNC_CLASS_TIME, "sec2gmt",   2,0,
		"Formats seconds since epoch as GMT timestamp with n\n"
		"decimal places for seconds, e.g. sec2gmt(1440768801.7,1) = \"2015-08-28T13:33:21.7Z\".\n"
		"Leaves non-numbers as-is."},
	{FUNC_CLASS_TIME, "sec2gmtdate", 1,0,
		"Formats seconds since epoch (integer part)\n"
		"as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = \"2015-08-28\".\n"
		"Leaves non-numbers as-is."},

	{FUNC_CLASS_TIME, "sec2localtime",   1,0, "Formats seconds since epoch (integer part)\n"
		"as local timestamp, e.g. sec2localtime(1440768801.7) = \"2015-08-28T13:33:21Z\".\n"
		"Consults $TZ environment variable. Leaves non-numbers as-is."},
	{FUNC_CLASS_TIME, "sec2localtime",   2,0,
		"Formats seconds since epoch as local timestamp with n\n"
		"decimal places for seconds, e.g. sec2localtime(1440768801.7,1) = \"2015-08-28T13:33:21.7Z\".\n"
		"Consults $TZ environment variable. Leaves non-numbers as-is."},
	{FUNC_CLASS_TIME, "sec2localdate", 1,0,
		"Formats seconds since epoch (integer part)\n"
		"as local timestamp with year-month-date, e.g. sec2localdate(1440768801.7) = \"2015-08-28\".\n"
		"Consults $TZ environment variable. Leaves non-numbers as-is."},

	{FUNC_CLASS_TIME, "sec2hms",   1,0,
		"Formats integer seconds as in\n"
		"sec2hms(5000) = \"01:23:20\""},
	{FUNC_CLASS_TIME, "strftime",  2,0,
		"Formats seconds since the epoch as timestamp, e.g.\n"
		"strftime(1440768801.7,\"%Y-%m-%dT%H:%M:%SZ\") = \"2015-08-28T13:33:21Z\", and\n"
		"strftime(1440768801.7,\"%Y-%m-%dT%H:%M:%3SZ\") = \"2015-08-28T13:33:21.700Z\".\n"
		"Format strings are as in the C library (please see \"man strftime\" on your system),\n"
		"with the Miller-specific addition of \"%1S\" through \"%9S\" which format the seconds\n"
		"with 1 through 9 decimal places, respectively. (\"%S\" uses no decimal places.)\n"
		"See also strftime_local."},
	{FUNC_CLASS_TIME, "strftime_local",  2,0,
		"Like strftime but consults the $TZ environment variable to get local time zone."},
	{FUNC_CLASS_TIME, "strptime",  2,0,
		"Parses timestamp as floating-point seconds since the epoch,\n"
		"e.g. strptime(\"2015-08-28T13:33:21Z\",\"%Y-%m-%dT%H:%M:%SZ\") = 1440768801.000000,\n"
		"and  strptime(\"2015-08-28T13:33:21.345Z\",\"%Y-%m-%dT%H:%M:%SZ\") = 1440768801.345000.\n"
		"See also strptime_local."},
	{FUNC_CLASS_TIME, "strptime_local",  2,0,
		"Like strptime, but consults $TZ environment variable to find and use local timezone."},
	{FUNC_CLASS_TIME, "systime",   0,0,
		"Floating-point seconds since the epoch,\n"
		"e.g. 1440768801.748936." },

	{FUNC_CLASS_TYPING, "is_absent",      1,0, "False if field is present in input, false otherwise"},
	{FUNC_CLASS_TYPING, "is_bool",        1,0, "True if field is present with boolean value. Synonymous with is_boolean."},
	{FUNC_CLASS_TYPING, "is_boolean",     1,0, "True if field is present with boolean value. Synonymous with is_bool."},
	{FUNC_CLASS_TYPING, "is_empty",       1,0, "True if field is present in input with empty string value, false otherwise."},
	{FUNC_CLASS_TYPING, "is_empty_map",    1,0, "True if argument is a map which is empty."},
	{FUNC_CLASS_TYPING, "is_float",       1,0, "True if field is present with value inferred to be float"},
	{FUNC_CLASS_TYPING, "is_int",         1,0, "True if field is present with value inferred to be int "},
	{FUNC_CLASS_TYPING, "is_map",         1,0, "True if argument is a map."},
	{FUNC_CLASS_TYPING, "is_nonempty_map", 1,0, "True if argument is a map which is non-empty."},
	{FUNC_CLASS_TYPING, "is_not_empty",    1,0, "False if field is present in input with empty value, false otherwise"},
	{FUNC_CLASS_TYPING, "is_not_map",      1,0, "True if argument is not a map."},
	{FUNC_CLASS_TYPING, "is_not_null",     1,0, "False if argument is null (empty or absent), true otherwise."},
	{FUNC_CLASS_TYPING, "is_null",        1,0, "True if argument is null (empty or absent), false otherwise."},
	{FUNC_CLASS_TYPING, "is_numeric",     1,0, "True if field is present with value inferred to be int or float"},
	{FUNC_CLASS_TYPING, "is_present",     1,0, "True if field is present in input, false otherwise."},
	{FUNC_CLASS_TYPING, "is_string",      1,0, "True if field is present with string (including empty-string) value"},

	{FUNC_CLASS_TYPING, "asserting_absent",      1,0, "Returns argument if it is absent in the input data, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_bool",        1,0, "Returns argument if it is present with boolean value, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_boolean",     1,0, "Returns argument if it is present with boolean value, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_empty",       1,0, "Returns argument if it is present in input with empty value,\n"
		"else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_empty_map",    1,0, "Returns argument if it is a map with empty value, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_float",       1,0, "Returns argument if it is present with float value, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_int",         1,0, "Returns argument if it is present with int value, else\n"
		"throws an error."},
	{FUNC_CLASS_TYPING, "asserting_map",         1,0, "Returns argument if it is a map, else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_nonempty_map", 1,0, "Returns argument if it is a non-empty map, else throws\n"
		"an error."},
	{FUNC_CLASS_TYPING, "asserting_not_empty",    1,0, "Returns argument if it is present in input with non-empty\n"
		"value, else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_not_map",      1,0, "Returns argument if it is not a map, else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_not_null",     1,0, "Returns argument if it is non-null (non-empty and non-absent),\n"
		"else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_null",        1,0, "Returns argument if it is null (empty or absent), else throws\n"
		"an error."},
	{FUNC_CLASS_TYPING, "asserting_numeric",     1,0, "Returns argument if it is present with int or float value,\n"
		"else throws an error."},
	{FUNC_CLASS_TYPING, "asserting_present",     1,0, "Returns argument if it is present in input, else throws\n"
		"an error."},
	{FUNC_CLASS_TYPING, "asserting_string",      1,0, "Returns argument if it is present with string (including\n"
		"empty-string) value, else throws an error."},

	{FUNC_CLASS_CONVERSION, "boolean",     1,0, "Convert int/float/bool/string to boolean."},
	{FUNC_CLASS_CONVERSION, "float",       1,0, "Convert int/float/bool/string to float."},
	{FUNC_CLASS_CONVERSION, "fmtnum",    2,0,
		"Convert int/float/bool to string using\n"
		"printf-style format string, e.g. '$s = fmtnum($n, \"%06lld\")'. WARNING: Miller numbers\n"
		"are all long long or double. If you use formats like %d or %f, behavior is undefined."},
	{FUNC_CLASS_CONVERSION, "hexfmt",    1,0, "Convert int to string, e.g. 255 to \"0xff\"."},
	{FUNC_CLASS_CONVERSION, "int",       1,0, "Convert int/float/bool/string to int."},
	{FUNC_CLASS_CONVERSION, "string",    1,0, "Convert int/float/bool/string to string."},
	{FUNC_CLASS_CONVERSION, "typeof",    1,0,
		"Convert argument to type of argument (e.g.\n"
		"MT_STRING). For debug."},

	{FUNC_CLASS_MAPS, "depth",         1,0, "Prints maximum depth of hashmap: ''. Scalars have depth 0."},
	{FUNC_CLASS_MAPS, "haskey",        2,0, "True/false if map has/hasn't key, e.g. 'haskey($*, \"a\")' or\n"
		"'haskey(mymap, mykey)'. Error if 1st argument is not a map."},
	{FUNC_CLASS_MAPS, "joink",         2,0, "Makes string from map keys. E.g. 'joink($*, \",\")'."},
	{FUNC_CLASS_MAPS, "joinkv",        3,0, "Makes string from map key-value pairs. E.g. 'joinkv(@v[2], \"=\", \",\")'"},
	{FUNC_CLASS_MAPS, "joinv",         2,0, "Makes string from map keys. E.g. 'joinv(mymap, \",\")'."},
	{FUNC_CLASS_MAPS, "leafcount",     1,0, "Counts total number of terminal values in hashmap. For single-level maps,\n"
		"same as length."},
	{FUNC_CLASS_MAPS, "length",        1,0, "Counts number of top-level entries in hashmap. Scalars have length 1."},
	{FUNC_CLASS_MAPS, "mapdiff",       0,1, "With 0 args, returns empty map. With 1 arg, returns copy of arg.\n"
		"With 2 or more, returns copy of arg 1 with all keys from any of remaining argument maps removed."},
	{FUNC_CLASS_MAPS, "mapexcept",     1,1, "Returns a map with keys from remaining arguments, if any, unset.\n"
		"E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'."},
	{FUNC_CLASS_MAPS, "mapselect",       1,1, "Returns a map with only keys from remaining arguments set.\n"
		"E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'."},
	{FUNC_CLASS_MAPS, "mapsum",        0,1, "With 0 args, returns empty map. With >= 1 arg, returns a map with\n"
		"key-value pairs from all arguments. Rightmost collisions win, e.g. 'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'."},
	{FUNC_CLASS_MAPS, "splitkv",       3,0, "Splits string by separators into map with type inference.\n"
		"E.g. 'splitkv(\"a=1,b=2,c=3\", \"=\", \",\")' gives '{\"a\" : 1, \"b\" : 2, \"c\" : 3}'."},
	{FUNC_CLASS_MAPS, "splitkvx",      3,0, "Splits string by separators into map without type inference (keys and\n"
		"values are strings). E.g. 'splitkv(\"a=1,b=2,c=3\", \"=\", \",\")' gives\n"
			"'{\"a\" : \"1\", \"b\" : \"2\", \"c\" : \"3\"}'."},
	{FUNC_CLASS_MAPS, "splitnv",       2,0, "Splits string by separator into integer-indexed map with type inference.\n"
		"E.g. 'splitnv(\"a,b,c\" , \",\")' gives '{1 : \"a\", 2 : \"b\", 3 : \"c\"}'."},
	{FUNC_CLASS_MAPS, "splitnvx",      2,0, "Splits string by separator into integer-indexed map without type\n"
		"inference (values are strings). E.g. 'splitnv(\"4,5,6\" , \",\")' gives '{1 : \"4\", 2 : \"5\", 3 : \"6\"}'."},

	{0, NULL, -1 , -1, NULL}, // table terminator
};

// ----------------------------------------------------------------
static arity_check_t check_arity(function_lookup_t lookup_table[], char* function_name,
	int user_provided_arity, int *parity, int* pvariadic)
{
	*parity = -1;
	*pvariadic = FALSE;
	int found_function_name = FALSE;
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &lookup_table[i];
		if (plookup->function_name == NULL)
			break;
		if (streq(function_name, plookup->function_name)) {
			found_function_name = TRUE;
			*parity = plookup->arity;
			if (plookup->variadic) {
				*pvariadic = TRUE;
				if (user_provided_arity < plookup->arity) {
					return ARITY_CHECK_FAIL;
				}
				return ARITY_CHECK_PASS;
			}
			if (user_provided_arity == plookup->arity) {
				return ARITY_CHECK_PASS;
			}
		}
	}
	if (found_function_name) {
		return ARITY_CHECK_FAIL;
	} else {
		return ARITY_CHECK_NO_SUCH;
	}
}

static void fmgr_check_arity_with_report(fmgr_t* pfmgr, char* function_name,
	int user_provided_arity, int* pvariadic)
{
	int arity = -1;
	arity_check_t result = check_arity(pfmgr->function_lookup_table, function_name, user_provided_arity,
		&arity, pvariadic);
	if (result == ARITY_CHECK_NO_SUCH) {
		fprintf(stderr, "%s: Function name \"%s\" not found.\n", MLR_GLOBALS.bargv0, function_name);
		exit(1);
	}
	if (result == ARITY_CHECK_FAIL) {
		// More flexibly, I'd have a list of arities supported by each
		// function. But this is overkill: there are unary and binary minus and sec2gmt,
		// and everything else has a single arity.
		if (streq(function_name, "-") || streq(function_name, "sec2gmt") || streq(function_name, "sec2localtime")) {
			fprintf(stderr, "%s: Function named \"%s\" takes one argument or two; got %d.\n",
				MLR_GLOBALS.bargv0, function_name, user_provided_arity);
		} else if (*pvariadic) {
			fprintf(stderr, "%s: Function named \"%s\" takes at least %d argument%s; got %d.\n",
				MLR_GLOBALS.bargv0, function_name, arity, (arity == 1) ? "" : "s", user_provided_arity);
		} else {
			fprintf(stderr, "%s: Function named \"%s\" takes %d argument%s; got %d.\n",
				MLR_GLOBALS.bargv0, function_name, arity, (arity == 1) ? "" : "s", user_provided_arity);
		}
		exit(1);
	}
}

static char* function_class_to_desc(func_class_t function_class) {
	switch(function_class) {
	case FUNC_CLASS_ARITHMETIC: return "arithmetic"; break;
	case FUNC_CLASS_MATH:       return "math";       break;
	case FUNC_CLASS_BOOLEAN:    return "boolean";    break;
	case FUNC_CLASS_STRING:     return "string";     break;
	case FUNC_CLASS_CONVERSION: return "conversion"; break;
	case FUNC_CLASS_TYPING:     return "typing";     break;
	case FUNC_CLASS_MAPS:       return "maps";       break;
	case FUNC_CLASS_TIME:       return "time";       break;
	default:                    return "???";        break;
	}
}

void fmgr_list_functions(fmgr_t* pfmgr, FILE* output_stream, char* leader) {
	char* separator = " ";
	int leaderlen = strlen(leader);
	int separatorlen = strlen(separator);
	int linelen = leaderlen;
	int j = 0;

	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		char* fname = plookup->function_name;
		if (fname == NULL)
			break;
		int fnamelen = strlen(fname);
		linelen += separatorlen + fnamelen;
		if (linelen >= 80) {
			fprintf(output_stream, "\n");
			linelen = 0;
			linelen = leaderlen + separatorlen + fnamelen;
			j = 0;
		}
		if (j == 0)
			fprintf(output_stream, "%s", leader);
		fprintf(output_stream, "%s%s", separator, fname);
		j++;
	}
	fprintf(output_stream, "\n");
}

// Pass function_name == NULL to get usage for all functions.
void fmgr_function_usage(fmgr_t* pfmgr, FILE* output_stream, char* function_name) {
	int found = FALSE;
	char* nfmt = "%s (class=%s #args=%d): %s\n";
	char* vfmt = "%s (class=%s variadic): %s\n";

	int num_printed = 0; // > 1 matches e.g. for - and sec2gmt
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL) // end of table
			break;
		if (function_name == NULL || streq(function_name, plookup->function_name)) {
			if (++num_printed > 1)
				fprintf(output_stream, "\n");
			if (plookup->variadic) {
				fprintf(output_stream, vfmt, plookup->function_name,
					function_class_to_desc(plookup->function_class),
					plookup->usage_string);
			} else {
				fprintf(output_stream, nfmt, plookup->function_name,
					function_class_to_desc(plookup->function_class),
					plookup->arity, plookup->usage_string);
			}
			found = TRUE;
		}
		if (function_name == NULL)
			fprintf(output_stream, "\n");
	}
	if (!found)
		fprintf(output_stream, "%s: no such function.\n", function_name);
	if (function_name == NULL) {
		fprintf(output_stream, "To set the seed for urand, you may specify decimal or hexadecimal 32-bit\n");
		fprintf(output_stream, "numbers of the form \"%s --seed 123456789\" or \"%s --seed 0xcafefeed\".\n",
			MLR_GLOBALS.bargv0, MLR_GLOBALS.bargv0);
		fprintf(output_stream, "Miller's built-in variables are NF, NR, FNR, FILENUM, and FILENAME (awk-like)\n");
		fprintf(output_stream, "along with the mathematical constants M_PI and M_E.\n");
	}
}

void fmgr_list_all_functions_raw(fmgr_t* pfmgr, FILE* output_stream) {
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL) // end of table
			break;
		printf("%s\n", plookup->function_name);
	}
}

// ================================================================
typedef struct _udf_callsite_state_t {
	int arity;
	rxval_evaluator_t** pevals;
	boxed_xval_t* args;
	udf_defsite_state_t* pdefsite_state;
} udf_callsite_state_t;

// ----------------------------------------------------------------
static udf_callsite_state_t* udf_callsite_state_alloc(
	fmgr_t*              pfmgr,
	udf_defsite_state_t* pdefsite_state,
	mlr_dsl_ast_node_t*  pnode,
	int                  arity,
	int                  type_inferencing,
	int                  context_flags)
{
	udf_callsite_state_t* pstate = mlr_malloc_or_die(sizeof(udf_callsite_state_t));

	pstate->arity = pnode->pchildren->length;

	pstate->pevals = mlr_malloc_or_die(pstate->arity * sizeof(rxval_evaluator_t*));
	int i = 0;
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* parg_node = pe->pvvalue;
		pstate->pevals[i] = rxval_evaluator_alloc_from_ast(parg_node,
			pfmgr, type_inferencing, context_flags);
	}

	pstate->args = mlr_malloc_or_die(pstate->arity * sizeof(boxed_xval_t));
	for (i = 0; i < pstate->arity; i++) {
		// Ownership will be transferred to local-stack which will be responsible for freeing.
		pstate->args[i] = box_ephemeral_val(mv_absent());
	}

	pstate->pdefsite_state = pdefsite_state;

	return pstate;
}

// ----------------------------------------------------------------
static void udf_callsite_state_eval_args(udf_callsite_state_t* pstate, variables_t* pvars) {
	for (int i = 0; i < pstate->arity; i++) {
		pstate->args[i] = pstate->pevals[i]->pprocess_func(pstate->pevals[i]->pvstate, pvars);
	}
}

// ----------------------------------------------------------------
static void udf_callsite_state_free(udf_callsite_state_t* pstate) {
	for (int i = 0; i < pstate->arity; i++) {
		rxval_evaluator_t* pxev = pstate->pevals[i];
		pxev->pfree_func(pxev);
	}
	free(pstate->pevals);
	free(pstate->args);
	free(pstate);
}

// ----------------------------------------------------------------
static mv_t rval_evaluator_udf_callsite_process(void* pvstate, variables_t* pvars) {
	udf_callsite_state_t* pstate = pvstate;

	udf_callsite_state_eval_args(pstate, pvars);

	// Functions returning map values in a scalar context get their return values treated as
	// absent-null. (E.g. f() returns a map and g() returns an int and the statement is '$x
	// = f() + g()'.) Non-scalar-context return values are handled separately (not here).
	boxed_xval_t retval = pstate->pdefsite_state->pprocess_func(
		pstate->pdefsite_state->pvstate, pstate->arity, pstate->args, pvars);

	if (retval.xval.is_terminal) {
		return retval.xval.terminal_mlrval;
	} else {
		if (retval.is_ephemeral) {
			mlhmmv_xvalue_free(&retval.xval);
		}
		return mv_absent();
	}
}

static boxed_xval_t rxval_evaluator_udf_xcallsite_process(void* pvstate, variables_t* pvars) {
	udf_callsite_state_t* pstate = pvstate;
	udf_callsite_state_eval_args(pstate, pvars);
	return pstate->pdefsite_state->pprocess_func(
		pstate->pdefsite_state->pvstate, pstate->arity, pstate->args, pvars);
}

static void rval_evaluator_udf_callsite_free(rval_evaluator_t* pevaluator) {
	udf_callsite_state_t* pstate = pevaluator->pvstate;
	udf_callsite_state_free(pstate);
	free(pevaluator);
}

static void rxval_evaluator_udf_xcallsite_free(rxval_evaluator_t* pxevaluator) {
	udf_callsite_state_t* pstate = pxevaluator->pvstate;
	udf_callsite_state_free(pstate);
	free(pxevaluator);
}

static rval_evaluator_t* fmgr_alloc_from_udf_callsite(fmgr_t* pfmgr, udf_defsite_state_t* pdefsite_state,
	mlr_dsl_ast_node_t* pnode, char* function_name, int arity, int type_inferencing, int context_flags)
{
	rval_evaluator_t* pudf_callsite_evaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	udf_callsite_state_t* pstate = udf_callsite_state_alloc(pfmgr, pdefsite_state, pnode,
		arity, type_inferencing, context_flags);

	pudf_callsite_evaluator->pvstate = pstate;
	pudf_callsite_evaluator->pprocess_func = rval_evaluator_udf_callsite_process;
	pudf_callsite_evaluator->pfree_func = rval_evaluator_udf_callsite_free;

	return pudf_callsite_evaluator;
}

static rxval_evaluator_t* fmgr_alloc_from_udf_xcallsite(fmgr_t* pfmgr, udf_defsite_state_t* pdefsite_state,
	mlr_dsl_ast_node_t* pnode, char* function_name, int arity, int type_inferencing, int context_flags)
{
	rxval_evaluator_t* pudf_xcallsite_evaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	udf_callsite_state_t* pstate = udf_callsite_state_alloc(pfmgr, pdefsite_state, pnode,
		arity, type_inferencing, context_flags);

	pudf_xcallsite_evaluator->pvstate = pstate;
	pudf_xcallsite_evaluator->pprocess_func = rxval_evaluator_udf_xcallsite_process;
	pudf_xcallsite_evaluator->pfree_func = rxval_evaluator_udf_xcallsite_free;

	return pudf_xcallsite_evaluator;
}

// ================================================================
typedef struct _unresolved_func_callsite_state_t {
	char* function_name;
	int arity;
	int type_inferencing;
	int context_flags;
	mlr_dsl_ast_node_t* pnode;
} unresolved_func_callsite_state_t;

static unresolved_func_callsite_state_t* unresolved_callsite_alloc(char* function_name, int arity,
	int type_inferencing, int context_flags, mlr_dsl_ast_node_t* pnode)
{
	unresolved_func_callsite_state_t* pstate = mlr_malloc_or_die(sizeof(unresolved_func_callsite_state_t));
	pstate->function_name    = mlr_strdup_or_die(function_name);
	pstate->arity            = arity;
	pstate->type_inferencing = type_inferencing;
	pstate->context_flags    = context_flags;
	pstate->pnode            = pnode;
	return pstate;
}

static void unresolved_callsite_free(unresolved_func_callsite_state_t* pstate) {
	if (pstate == NULL)
		return;
	free(pstate->function_name);
	free(pstate);
}

// ----------------------------------------------------------------
static mv_t provisional_call_func(void* pvstate, variables_t* pvars) {
	unresolved_func_callsite_state_t* pstate = pvstate;
	fprintf(stderr,
		"%s: internal coding error: unresolved scalar-return-value callsite \"%s\".\n",
		MLR_GLOBALS.bargv0, pstate->function_name);
	exit(1);
}

static void provisional_call_free(rval_evaluator_t* pevaluator) {
	unresolved_func_callsite_state_t* pstate = pevaluator->pvstate;
	unresolved_callsite_free(pstate);
	free(pevaluator);
}

rval_evaluator_t* fmgr_alloc_provisional_from_operator_or_function_call(fmgr_t* pfmgr, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	char* function_name = pnode->text;
	int user_provided_arity = pnode->pchildren->length;

	unresolved_func_callsite_state_t* pstate = unresolved_callsite_alloc(function_name, user_provided_arity,
		type_inferencing, context_flags, pnode);

	rval_evaluator_t* pev = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pev->pvstate       = pstate;
	pev->pprocess_func = provisional_call_func;
	pev->pfree_func    = provisional_call_free;

	// Remember this callsite to a function which may or may not have been defined yet.
	// Then later we can resolve them to point to UDF bodies which have been defined.
	fmgr_mark_callsite_to_resolve(pfmgr, pev);

	return pev;
}

// ----------------------------------------------------------------
static boxed_xval_t provisional_xcall_func(void* pvstate, variables_t* pvars) {
	unresolved_func_callsite_state_t* pstate = pvstate;
	fprintf(stderr,
		"%s: internal coding error: unresolved map-return-value callsite \"%s\".\n",
		MLR_GLOBALS.bargv0, pstate->function_name);
	exit(1);
}

static void provisional_xcall_free(rxval_evaluator_t* pxevaluator) {
	unresolved_func_callsite_state_t* pstate = pxevaluator->pvstate;
	unresolved_callsite_free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* fmgr_xalloc_provisional_from_operator_or_function_call(fmgr_t* pfmgr, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	char* function_name = pnode->text;
	int user_provided_arity = pnode->pchildren->length;

	unresolved_func_callsite_state_t* pstate = unresolved_callsite_alloc(function_name, user_provided_arity,
		type_inferencing, context_flags, pnode);

	rxval_evaluator_t* pxev = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxev->pvstate       = pstate;
	pxev->pprocess_func = provisional_xcall_func;
	pxev->pfree_func    = provisional_xcall_free;

	// Remember this callsite to a function which may or may not have been defined yet.
	// Then later we can resolve them to point to UDF bodies which have been defined.
	fmgr_mark_xcallsite_to_resolve(pfmgr, pxev);

	return pxev;
}

// ----------------------------------------------------------------
void fmgr_mark_callsite_to_resolve(fmgr_t* pfmgr, rval_evaluator_t* pev) {
	sllv_append(pfmgr->pfunc_callsite_evaluators_to_resolve, pev);
}

void fmgr_mark_xcallsite_to_resolve(fmgr_t* pfmgr, rxval_evaluator_t* pxev) {
	sllv_append(pfmgr->pfunc_callsite_xevaluators_to_resolve, pxev);
}

// ----------------------------------------------------------------
// Resolving a callsite involves treewalking the AST which may find more callsites to
// resolve. E.g. in '$y = f(g($x))', f is initially unresolved (f and/or g perhaps as yet
// undefined as of when the callsite is parsed), then at resolution time for f, its
// argument 'g($x)' is encountered, initially unresolved, then resolved.
// Hence the outer loop.
void fmgr_resolve_func_callsites(fmgr_t* pfmgr) {
	while (TRUE) {
		int did = FALSE;
		while (pfmgr->pfunc_callsite_xevaluators_to_resolve->phead != NULL) {
			did = TRUE;
			rxval_evaluator_t* pxev = sllv_pop(pfmgr->pfunc_callsite_xevaluators_to_resolve);
			unresolved_func_callsite_state_t* ptemp_state = pxev->pvstate;
			resolve_func_xcallsite(pfmgr, pxev);
			unresolved_callsite_free(ptemp_state);
		}

		while (pfmgr->pfunc_callsite_evaluators_to_resolve->phead != NULL) {
			did = TRUE;
			rval_evaluator_t* pev = sllv_pop(pfmgr->pfunc_callsite_evaluators_to_resolve);
			unresolved_func_callsite_state_t* ptemp_state = pev->pvstate;
			resolve_func_callsite(pfmgr, pev);
			unresolved_callsite_free(ptemp_state);
		}
		if (!did) {
			break;
		}
	}
}

// ----------------------------------------------------------------
static rval_evaluator_t* construct_udf_callsite_evaluator(
	fmgr_t* pfmgr,
	unresolved_func_callsite_state_t* pcallsite)
{
	char* function_name       = pcallsite->function_name;
	int   user_provided_arity = pcallsite->arity;
	int   type_inferencing    = pcallsite->type_inferencing;
	int   context_flags       = pcallsite->context_flags;
	mlr_dsl_ast_node_t* pnode = pcallsite->pnode;

	udf_defsite_state_t* pudf_defsite_state = lhmsv_get(pfmgr->pudf_names_to_defsite_states,
		pcallsite->function_name);

	if (pudf_defsite_state != NULL) {
		int udf_arity = pudf_defsite_state->arity;
		if (user_provided_arity != udf_arity) {
			fprintf(stderr, "Function named \"%s\" takes %d argument%s; got %d.\n",
				function_name, udf_arity, (udf_arity == 1) ? "" : "s", user_provided_arity);
			exit(1);
		}

		return fmgr_alloc_from_udf_callsite(pfmgr, pudf_defsite_state,
			pnode, function_name, user_provided_arity, type_inferencing, context_flags);
	} else {
		return NULL;
	}
}

static rxval_evaluator_t* construct_udf_defsite_xevaluator(
	fmgr_t* pfmgr,
	unresolved_func_callsite_state_t* pcallsite)
{
	char* function_name       = pcallsite->function_name;
	int   user_provided_arity = pcallsite->arity;
	int   type_inferencing    = pcallsite->type_inferencing;
	int   context_flags       = pcallsite->context_flags;
	mlr_dsl_ast_node_t* pnode = pcallsite->pnode;

	udf_defsite_state_t* pudf_defsite_state = lhmsv_get(pfmgr->pudf_names_to_defsite_states,
		pcallsite->function_name);

	if (pudf_defsite_state != NULL) {
		int udf_arity = pudf_defsite_state->arity;
		if (user_provided_arity != udf_arity) {
			fprintf(stderr, "Function named \"%s\" takes %d argument%s; got %d.\n",
				function_name, udf_arity, (udf_arity == 1) ? "" : "s", user_provided_arity);
			exit(1);
		}

		return fmgr_alloc_from_udf_xcallsite(pfmgr, pudf_defsite_state,
			pnode, function_name, user_provided_arity, type_inferencing, context_flags);
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static rval_evaluator_t* construct_builtin_function_callsite_evaluator(
	fmgr_t* pfmgr,
	unresolved_func_callsite_state_t* pcallsite)
{
	char* function_name       = pcallsite->function_name;
	int   user_provided_arity = pcallsite->arity;
	int   type_inferencing    = pcallsite->type_inferencing;
	int   context_flags       = pcallsite->context_flags;
	mlr_dsl_ast_node_t* pnode = pcallsite->pnode;

	int variadic = FALSE;
	fmgr_check_arity_with_report(pfmgr, function_name, user_provided_arity, &variadic);

	rval_evaluator_t* pevaluator = NULL;
	if (variadic) {
		int nargs = pnode->pchildren->length;
		rval_evaluator_t** pargs = mlr_malloc_or_die(nargs * sizeof(rval_evaluator_t*));
		int i = 0;
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pargs[i] = rval_evaluator_alloc_from_ast(pchild, pfmgr, type_inferencing, context_flags);
		}
		pevaluator = fmgr_alloc_evaluator_from_variadic_func_name(function_name, pargs, nargs);

	} else if (user_provided_arity == 0) {
		pevaluator = fmgr_alloc_evaluator_from_zary_func_name(function_name);
	} else if (user_provided_arity == 1) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
		pevaluator = fmgr_alloc_evaluator_from_unary_func_name(function_name, parg1);
	} else if (user_provided_arity == 2) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
		int type2 = parg2_node->type;

		int is_regexy =
			streq(function_name, "=~") ||
			streq(function_name, "!=~") ||
			streq(function_name, "regex_extract");

		if (is_regexy && type2 == MD_AST_NODE_TYPE_STRING_LITERAL) {
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_binary_regex_arg2_func_name(function_name,
				parg1, parg2_node->text, FALSE);
		} else if (is_regexy && type2 == MD_AST_NODE_TYPE_REGEXI) {
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_binary_regex_arg2_func_name(function_name, parg1, parg2_node->text,
				TYPE_INFER_STRING_FLOAT_INT);
		} else {
			// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
			// the regexes will be compiled record-by-record rather than once at alloc time, which will
			// be slower.
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast(parg2_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_binary_func_name(function_name, parg1, parg2);
		}

	} else if (user_provided_arity == 3) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
		mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvvalue;
		int type2 = parg2_node->type;

		int is_regexy =
			streq(function_name, "sub") ||
			streq(function_name, "gsub");

		if (is_regexy && type2 == MD_AST_NODE_TYPE_STRING_LITERAL) {
			// sub/gsub-regex special case:
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast(parg3_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_ternary_regex_arg2_func_name(function_name, parg1, parg2_node->text,
				FALSE, parg3);

		} else if (is_regexy && type2 == MD_AST_NODE_TYPE_REGEXI) {
			// sub/gsub-regex special case:
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast(parg3_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_ternary_regex_arg2_func_name(function_name, parg1, parg2_node->text,
				TYPE_INFER_STRING_FLOAT_INT, parg3);

		} else {
			// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
			// the regexes will be compiled record-by-record rather than once at alloc time, which will
			// be slower.
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast(parg1_node, pfmgr, type_inferencing, context_flags);
			rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast(parg2_node, pfmgr, type_inferencing, context_flags);
			rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast(parg3_node, pfmgr, type_inferencing, context_flags);
			pevaluator = fmgr_alloc_evaluator_from_ternary_func_name(function_name, parg1, parg2, parg3);
		}

	} else {
		fprintf(stderr, "Miller: internal coding error:  arity for function name \"%s\" misdetected.\n",
			function_name);
		exit(1);
	}

	return pevaluator;
}

// ----------------------------------------------------------------
// At callsites, arguments can be scalars or maps; return values can be scalars
// or maps.  At the user level, a function take map input and produce scalar
// output or vice versa. As of this writing, though, *internally* functions
// go from scalars to scalar or maps to map. This wrapper wraps scalar input
// to functions which know about maps.

typedef struct _xeval_wrapping_eval_state_t {
	rval_evaluator_t* pevaluator;
} xeval_wrapping_eval_state_t;

static boxed_xval_t xeval_wrapping_eval_func(void* pvstate, variables_t* pvars) {
	xeval_wrapping_eval_state_t* pstate = pvstate;
	rval_evaluator_t* pevaluator = pstate->pevaluator;
	mv_t val = pevaluator->pprocess_func(pevaluator->pvstate, pvars);
	return (boxed_xval_t) {
		.xval = mlhmmv_xvalue_wrap_terminal(val),
		.is_ephemeral = TRUE, // xxx verify reference semantics for RHS evaluators!
	};
}

static void xeval_wrapping_eval_free(rxval_evaluator_t* pxevaluator) {
	xeval_wrapping_eval_state_t* pstate = pxevaluator->pvstate;
	pstate->pevaluator->pfree_func(pstate->pevaluator);
	free(pstate);
	free(pxevaluator);
}

static rxval_evaluator_t* fmgr_alloc_xeval_wrapping_eval(rval_evaluator_t* pevaluator) {
	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));

	xeval_wrapping_eval_state_t* pstate = mlr_malloc_or_die(sizeof(xeval_wrapping_eval_state_t));
	pstate->pevaluator = pevaluator;

	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = xeval_wrapping_eval_func;
	pxevaluator->pfree_func    = xeval_wrapping_eval_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
// At callsites, arguments can be scalars or maps; return values can be scalars
// or maps.  At the user level, a function take map input and produce scalar
// output or vice versa. As of this writing, though, *internally* functions go
// from scalars to scalar or maps to map. This wrapper wraps maybe-map input to
// functions which do not know about maps.

typedef struct _eval_wrapping_xeval_state_t {
	rxval_evaluator_t* pxevaluator;
} eval_wrapping_xeval_state_t;

static mv_t eval_wrapping_xeval_func(void* pvstate, variables_t* pvars) {
	eval_wrapping_xeval_state_t* pstate = pvstate;
	rxval_evaluator_t* pxevaluator = pstate->pxevaluator;
	boxed_xval_t bxval = pxevaluator->pprocess_func(pxevaluator->pvstate, pvars);

	if (bxval.xval.is_terminal) {
		if (bxval.is_ephemeral) {
			return bxval.xval.terminal_mlrval;
		} else {
			return mv_copy(&bxval.xval.terminal_mlrval);
		}

	} else {
		if (bxval.is_ephemeral) {
			mlhmmv_xvalue_free(&bxval.xval);
		}
		return mv_error();
	}

}

static void eval_wrapping_xeval_free(rval_evaluator_t* pevaluator) {
	eval_wrapping_xeval_state_t* pstate = pevaluator->pvstate;
	pstate->pxevaluator->pfree_func(pstate->pxevaluator);
	free(pstate);
	free(pevaluator);
}

static rval_evaluator_t* fmgr_alloc_eval_wrapping_xeval(rxval_evaluator_t* pxevaluator) {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	eval_wrapping_xeval_state_t* pstate = mlr_malloc_or_die(sizeof(eval_wrapping_xeval_state_t));
	pstate->pxevaluator = pxevaluator;

	pevaluator->pvstate       = pstate;
	pevaluator->pprocess_func = eval_wrapping_xeval_func;
	pevaluator->pfree_func    = eval_wrapping_xeval_free;

	return pevaluator;
}

// ================================================================
static rval_evaluator_t* fmgr_alloc_evaluator_from_variadic_func_name(char* fnnm, rval_evaluator_t** pargs, int nargs) {
	if        (streq(fnnm, "min")) { return rval_evaluator_alloc_from_variadic_func(variadic_min_func, pargs, nargs);
	} else if (streq(fnnm, "max")) { return rval_evaluator_alloc_from_variadic_func(variadic_max_func, pargs, nargs);
	} else return NULL;
}

// ================================================================
static rval_evaluator_t* fmgr_alloc_evaluator_from_zary_func_name(char* function_name) {
	if        (streq(function_name, "urand")) {
		return rval_evaluator_alloc_from_x_z_func(f_z_urand_func);
	} else if (streq(function_name, "urand32")) {
		return rval_evaluator_alloc_from_x_z_func(i_z_urand32_func);
	} else if (streq(function_name, "systime")) {
		return rval_evaluator_alloc_from_x_z_func(f_z_systime_func);
	} else  {
		return NULL;
	}
}

// ================================================================
static rval_evaluator_t* fmgr_alloc_evaluator_from_unary_func_name(char* fnnm, rval_evaluator_t* parg1)  {
	if        (streq(fnnm, "!"))               { return rval_evaluator_alloc_from_b_b_func(b_b_not_func,         parg1);
	} else if (streq(fnnm, "+"))               { return rval_evaluator_alloc_from_x_x_func(x_x_upos_func,        parg1);
	} else if (streq(fnnm, "-"))               { return rval_evaluator_alloc_from_x_x_func(x_x_uneg_func,        parg1);
	} else if (streq(fnnm, ".+"))              { return rval_evaluator_alloc_from_x_x_func(x_x_upos_func,        parg1);
	} else if (streq(fnnm, ".-"))              { return rval_evaluator_alloc_from_x_x_func(x_x_uneg_func,        parg1);
	} else if (streq(fnnm, "abs"))             { return rval_evaluator_alloc_from_x_x_func(x_x_abs_func,         parg1);
	} else if (streq(fnnm, "acos"))            { return rval_evaluator_alloc_from_f_f_func(f_f_acos_func,        parg1);
	} else if (streq(fnnm, "acosh"))           { return rval_evaluator_alloc_from_f_f_func(f_f_acosh_func,       parg1);
	} else if (streq(fnnm, "asin"))            { return rval_evaluator_alloc_from_f_f_func(f_f_asin_func,        parg1);
	} else if (streq(fnnm, "asinh"))           { return rval_evaluator_alloc_from_f_f_func(f_f_asinh_func,       parg1);
	} else if (streq(fnnm, "atan"))            { return rval_evaluator_alloc_from_f_f_func(f_f_atan_func,        parg1);
	} else if (streq(fnnm, "atanh"))           { return rval_evaluator_alloc_from_f_f_func(f_f_atanh_func,       parg1);
	} else if (streq(fnnm, "bitcount"))        { return rval_evaluator_alloc_from_i_i_func(i_i_bitcount_func,    parg1);
	} else if (streq(fnnm, "boolean"))         { return rval_evaluator_alloc_from_x_x_func(b_x_boolean_func,     parg1);
	} else if (streq(fnnm, "cbrt"))            { return rval_evaluator_alloc_from_f_f_func(f_f_cbrt_func,        parg1);
	} else if (streq(fnnm, "ceil"))            { return rval_evaluator_alloc_from_x_x_func(x_x_ceil_func,        parg1);
	} else if (streq(fnnm, "cos"))             { return rval_evaluator_alloc_from_f_f_func(f_f_cos_func,         parg1);
	} else if (streq(fnnm, "cosh"))            { return rval_evaluator_alloc_from_f_f_func(f_f_cosh_func,        parg1);
	} else if (streq(fnnm, "dhms2fsec"))       { return rval_evaluator_alloc_from_f_s_func(f_s_dhms2fsec_func,   parg1);
	} else if (streq(fnnm, "dhms2sec"))        { return rval_evaluator_alloc_from_f_s_func(i_s_dhms2sec_func,    parg1);
	} else if (streq(fnnm, "erf"))             { return rval_evaluator_alloc_from_f_f_func(f_f_erf_func,         parg1);
	} else if (streq(fnnm, "erfc"))            { return rval_evaluator_alloc_from_f_f_func(f_f_erfc_func,        parg1);
	} else if (streq(fnnm, "exp"))             { return rval_evaluator_alloc_from_f_f_func(f_f_exp_func,         parg1);
	} else if (streq(fnnm, "expm1"))           { return rval_evaluator_alloc_from_f_f_func(f_f_expm1_func,       parg1);
	} else if (streq(fnnm, "float"))           { return rval_evaluator_alloc_from_x_x_func(f_x_float_func,       parg1);
	} else if (streq(fnnm, "floor"))           { return rval_evaluator_alloc_from_x_x_func(x_x_floor_func,       parg1);
	} else if (streq(fnnm, "fsec2dhms"))       { return rval_evaluator_alloc_from_s_f_func(s_f_fsec2dhms_func,   parg1);
	} else if (streq(fnnm, "fsec2hms"))        { return rval_evaluator_alloc_from_s_f_func(s_f_fsec2hms_func,    parg1);
	} else if (streq(fnnm, "gmt2sec"))         { return rval_evaluator_alloc_from_i_s_func(i_s_gmt2sec_func,     parg1);
	} else if (streq(fnnm, "localtime2sec"))   { return rval_evaluator_alloc_from_i_s_func(i_s_localtime2sec_func, parg1);
	} else if (streq(fnnm, "hexfmt"))          { return rval_evaluator_alloc_from_x_x_func(s_x_hexfmt_func,      parg1);
	} else if (streq(fnnm, "hms2fsec"))        { return rval_evaluator_alloc_from_f_s_func(f_s_hms2fsec_func,    parg1);
	} else if (streq(fnnm, "hms2sec"))         { return rval_evaluator_alloc_from_f_s_func(i_s_hms2sec_func,     parg1);
	} else if (streq(fnnm, "int"))             { return rval_evaluator_alloc_from_x_x_func(i_x_int_func,         parg1);
	} else if (streq(fnnm, "invqnorm"))        { return rval_evaluator_alloc_from_f_f_func(f_f_invqnorm_func,    parg1);
	} else if (streq(fnnm, "log"))             { return rval_evaluator_alloc_from_f_f_func(f_f_log_func,         parg1);
	} else if (streq(fnnm, "log10"))           { return rval_evaluator_alloc_from_f_f_func(f_f_log10_func,       parg1);
	} else if (streq(fnnm, "log1p"))           { return rval_evaluator_alloc_from_f_f_func(f_f_log1p_func,       parg1);
	} else if (streq(fnnm, "qnorm"))           { return rval_evaluator_alloc_from_f_f_func(f_f_qnorm_func,       parg1);
	} else if (streq(fnnm, "round"))           { return rval_evaluator_alloc_from_x_x_func(x_x_round_func,       parg1);
	} else if (streq(fnnm, "sec2dhms"))        { return rval_evaluator_alloc_from_s_i_func(s_i_sec2dhms_func,    parg1);
	} else if (streq(fnnm, "sec2gmt"))         { return rval_evaluator_alloc_from_x_x_func(s_x_sec2gmt_func,     parg1);
	} else if (streq(fnnm, "sec2gmtdate"))     { return rval_evaluator_alloc_from_x_x_func(s_x_sec2gmtdate_func, parg1);
	} else if (streq(fnnm, "sec2localtime"))   { return rval_evaluator_alloc_from_x_x_func(s_x_sec2localtime_func, parg1);
	} else if (streq(fnnm, "sec2localdate"))   { return rval_evaluator_alloc_from_x_x_func(s_x_sec2localdate_func, parg1);
	} else if (streq(fnnm, "sec2hms"))         { return rval_evaluator_alloc_from_s_i_func(s_i_sec2hms_func,     parg1);
	} else if (streq(fnnm, "sgn"))             { return rval_evaluator_alloc_from_x_x_func(x_x_sgn_func,         parg1);
	} else if (streq(fnnm, "sin"))             { return rval_evaluator_alloc_from_f_f_func(f_f_sin_func,         parg1);
	} else if (streq(fnnm, "sinh"))            { return rval_evaluator_alloc_from_f_f_func(f_f_sinh_func,        parg1);
	} else if (streq(fnnm, "sqrt"))            { return rval_evaluator_alloc_from_f_f_func(f_f_sqrt_func,        parg1);
	} else if (streq(fnnm, "string"))          { return rval_evaluator_alloc_from_x_x_func(s_x_string_func,      parg1);
	} else if (streq(fnnm, "strlen"))          { return rval_evaluator_alloc_from_i_s_func(i_s_strlen_func,      parg1);
	} else if (streq(fnnm, "tan"))             { return rval_evaluator_alloc_from_f_f_func(f_f_tan_func,         parg1);
	} else if (streq(fnnm, "tanh"))            { return rval_evaluator_alloc_from_f_f_func(f_f_tanh_func,        parg1);
	} else if (streq(fnnm, "tolower"))         { return rval_evaluator_alloc_from_s_s_func(s_s_tolower_func,     parg1);
	} else if (streq(fnnm, "toupper"))         { return rval_evaluator_alloc_from_s_s_func(s_s_toupper_func,     parg1);
	} else if (streq(fnnm, "~"))               { return rval_evaluator_alloc_from_i_i_func(i_i_bitwise_not_func, parg1);

	} else return NULL;
}

// ================================================================
static rval_evaluator_t* fmgr_alloc_evaluator_from_binary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	if        (streq(fnnm, "&&"))   { return rval_evaluator_alloc_from_b_bb_and_func(parg1, parg2);
	} else if (streq(fnnm, "||"))   { return rval_evaluator_alloc_from_b_bb_or_func (parg1, parg2);
	} else if (streq(fnnm, "^^"))   { return rval_evaluator_alloc_from_b_bb_xor_func(parg1, parg2);
	} else if (streq(fnnm, "=~"))   { return rval_evaluator_alloc_from_x_ssc_func(
		matches_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "regex_extract"))   { return rval_evaluator_alloc_from_x_ss_func(
		regex_extract_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "!=~"))  { return rval_evaluator_alloc_from_x_ssc_func(does_not_match_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "=="))   { return rval_evaluator_alloc_from_x_xx_func(eq_op_func,             parg1, parg2);
	} else if (streq(fnnm, "!="))   { return rval_evaluator_alloc_from_x_xx_func(ne_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">"))    { return rval_evaluator_alloc_from_x_xx_func(gt_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">="))   { return rval_evaluator_alloc_from_x_xx_func(ge_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<"))    { return rval_evaluator_alloc_from_x_xx_func(lt_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<="))   { return rval_evaluator_alloc_from_x_xx_func(le_op_func,             parg1, parg2);
	} else if (streq(fnnm, "."))    { return rval_evaluator_alloc_from_x_xx_func(s_xx_dot_func,          parg1, parg2);

	} else if (streq(fnnm, "+"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_plus_func,         parg1, parg2);
	} else if (streq(fnnm, "-"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_minus_func,        parg1, parg2);
	} else if (streq(fnnm, "*"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_times_func,        parg1, parg2);
	} else if (streq(fnnm, "/"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_divide_func,       parg1, parg2);
	} else if (streq(fnnm, "//"))   { return rval_evaluator_alloc_from_x_xx_func(x_xx_int_divide_func,   parg1, parg2);

	} else if (streq(fnnm, ".+"))   { return rval_evaluator_alloc_from_x_xx_func(x_xx_oplus_func,        parg1, parg2);
	} else if (streq(fnnm, ".-"))   { return rval_evaluator_alloc_from_x_xx_func(x_xx_ominus_func,       parg1, parg2);
	} else if (streq(fnnm, ".*"))   { return rval_evaluator_alloc_from_x_xx_func(x_xx_otimes_func,       parg1, parg2);
	} else if (streq(fnnm, "./"))   { return rval_evaluator_alloc_from_x_xx_func(x_xx_odivide_func,      parg1, parg2);
	} else if (streq(fnnm, ".//"))  { return rval_evaluator_alloc_from_x_xx_func(x_xx_int_odivide_func,  parg1, parg2);

	} else if (streq(fnnm, "%"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_mod_func,          parg1, parg2);
	} else if (streq(fnnm, "**"))   { return rval_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "pow"))  { return rval_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "atan2")){ return rval_evaluator_alloc_from_f_ff_func(f_ff_atan2_func,        parg1, parg2);
	} else if (streq(fnnm, "roundm")) { return rval_evaluator_alloc_from_x_xx_func(x_xx_roundm_func,     parg1, parg2);
	} else if (streq(fnnm, "fmtnum")) { return rval_evaluator_alloc_from_s_xs_func(s_xs_fmtnum_func,     parg1, parg2);
	} else if (streq(fnnm, "urandint")) { return rval_evaluator_alloc_from_i_ii_func(i_ii_urandint_func, parg1, parg2);
	} else if (streq(fnnm, "sec2gmt"))  { return rval_evaluator_alloc_from_x_xi_func(s_xi_sec2gmt_func,  parg1, parg2);
	} else if (streq(fnnm, "sec2localtime")) { return rval_evaluator_alloc_from_x_xi_func(s_xi_sec2localtime_func, parg1, parg2);
	} else if (streq(fnnm, "&"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_band_func,         parg1, parg2);
	} else if (streq(fnnm, "|"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_bor_func,          parg1, parg2);
	} else if (streq(fnnm, "^"))    { return rval_evaluator_alloc_from_x_xx_func(x_xx_bxor_func,         parg1, parg2);
	} else if (streq(fnnm, "<<"))   { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_lsh_func,  parg1, parg2);
	} else if (streq(fnnm, ">>"))   { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_rsh_func,  parg1, parg2);
	} else if (streq(fnnm, "strftime")) { return rval_evaluator_alloc_from_x_ns_func(s_ns_strftime_func, parg1, parg2);
	} else if (streq(fnnm, "strftime_local")) { return rval_evaluator_alloc_from_x_ns_func(s_ns_strftime_local_func, parg1, parg2);
	} else if (streq(fnnm, "strptime")) { return rval_evaluator_alloc_from_x_ss_func(i_ss_strptime_func, parg1, parg2);
	} else if (streq(fnnm, "strptime_local")) { return rval_evaluator_alloc_from_x_ss_func(i_ss_strptime_local_func, parg1, parg2);
	} else  { return NULL; }
}

static rval_evaluator_t* fmgr_alloc_evaluator_from_binary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	if        (streq(fnnm, "=~"))  {
		return rval_evaluator_alloc_from_x_sr_func(matches_precomp_func,        parg1, regex_string, ignore_case);
	} else if (streq(fnnm, "!=~")) {
		return rval_evaluator_alloc_from_x_sr_func(does_not_match_precomp_func, parg1, regex_string, ignore_case);
	} else if (streq(fnnm, "regex_extract")) {
		return rval_evaluator_alloc_from_x_se_func(regex_extract_precomp_func, parg1, regex_string, ignore_case);
	} else  { return NULL; }
}

// ================================================================
static rval_evaluator_t* fmgr_alloc_evaluator_from_ternary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	if (streq(fnnm, "sub")) {
		return rval_evaluator_alloc_from_s_sss_func(sub_no_precomp_func,  parg1, parg2, parg3);
	} else if (streq(fnnm, "gsub")) {
		return rval_evaluator_alloc_from_s_sss_func(gsub_no_precomp_func, parg1, parg2, parg3);
	} else if (streq(fnnm, "ssub")) {
		return rval_evaluator_alloc_from_s_sss_func(s_sss_ssub_func,      parg1, parg2, parg3);
	} else if (streq(fnnm, "logifit")) {
		return rval_evaluator_alloc_from_f_fff_func(f_fff_logifit_func,   parg1, parg2, parg3);
	} else if (streq(fnnm, "madd")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modadd_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "msub")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modsub_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mmul")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modmul_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mexp")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modexp_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "substr")) {
		return rval_evaluator_alloc_from_s_sii_func(s_sii_substr_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "? :")) {
		return rval_evaluator_alloc_from_ternop(parg1, parg2, parg3);
	} else  { return NULL; }
}

static rval_evaluator_t* fmgr_alloc_evaluator_from_ternary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3)
{
	if (streq(fnnm, "sub"))  {
		return rval_evaluator_alloc_from_x_srs_func(sub_precomp_func,  parg1, regex_string, ignore_case, parg3);
	} else if (streq(fnnm, "gsub"))  {
		return rval_evaluator_alloc_from_x_srs_func(gsub_precomp_func, parg1, regex_string, ignore_case, parg3);
	} else  { return NULL; }
}

// ================================================================
static rxval_evaluator_t* construct_builtin_function_callsite_xevaluator(
	fmgr_t* pfmgr,
	unresolved_func_callsite_state_t* pcallsite)
{
	char* function_name       = pcallsite->function_name;
	int   user_provided_arity = pcallsite->arity;
	int   type_inferencing    = pcallsite->type_inferencing;
	int   context_flags       = pcallsite->context_flags;
	mlr_dsl_ast_node_t* pnode = pcallsite->pnode;

	int variadic = FALSE;
	fmgr_check_arity_with_report(pfmgr, function_name, user_provided_arity, &variadic);

	rxval_evaluator_t* pxevaluator = NULL;
	if (variadic) {
		pxevaluator = fmgr_alloc_xevaluator_from_variadic_func_name(function_name, pnode->pchildren,
			pfmgr, type_inferencing, context_flags);

	} else if (user_provided_arity == 1) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		pxevaluator = fmgr_alloc_xevaluator_from_unary_func_name(function_name, parg1_node,
			pfmgr, type_inferencing, context_flags);

	} else if (user_provided_arity == 2) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
		pxevaluator = fmgr_alloc_xevaluator_from_binary_func_name(function_name, parg1_node, parg2_node,
			pfmgr, type_inferencing, context_flags);

	} else if (user_provided_arity == 3) {
		mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
		mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvvalue;
		pxevaluator = fmgr_alloc_xevaluator_from_ternary_func_name(function_name, parg1_node, parg2_node, parg3_node,
			pfmgr, type_inferencing, context_flags);
	}

	return pxevaluator;
}

// ----------------------------------------------------------------
static rxval_evaluator_t* fmgr_alloc_xevaluator_from_variadic_func_name(
	char*               function_name,
	sllv_t*             parg_nodes,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags)
{
	if (streq(function_name, "mapsum")) {
		return rxval_evaluator_alloc_from_variadic_func(variadic_mapsum_xfunc, parg_nodes,
			pfmgr, type_inferencing, context_flags);
	} else if (streq(function_name, "mapdiff")) {
		return rxval_evaluator_alloc_from_variadic_func(variadic_mapdiff_xfunc, parg_nodes,
			pfmgr, type_inferencing, context_flags);
	} else if (streq(function_name, "mapexcept")) {
		return rxval_evaluator_alloc_from_variadic_func(variadic_mapexcept_xfunc, parg_nodes,
			pfmgr, type_inferencing, context_flags);
	} else if (streq(function_name, "mapselect")) {
		return rxval_evaluator_alloc_from_variadic_func(variadic_mapselect_xfunc, parg_nodes,
			pfmgr, type_inferencing, context_flags);
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static rxval_evaluator_t* fmgr_alloc_xevaluator_from_unary_func_name(char* fnnm,
	mlr_dsl_ast_node_t* parg1,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/)
{

	if (streq(fnnm, "asserting_absent")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_absent_no_free_xfunc, parg1, pf, ti, cf, "absent");
	} else if (streq(fnnm, "asserting_bool")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_boolean_no_free_xfunc, parg1, pf, ti, cf, "boolean");
	} else if (streq(fnnm, "asserting_boolean")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_boolean_no_free_xfunc, parg1, pf, ti, cf, "boolean");
	} else if (streq(fnnm, "asserting_empty")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_empty_no_free_xfunc, parg1, pf, ti, cf, "empty");
	} else if (streq(fnnm, "asserting_empty_map")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_empty_map_no_free_xfunc, parg1, pf, ti, cf, "empty_map");
	} else if (streq(fnnm, "asserting_float")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_float_no_free_xfunc, parg1, pf, ti, cf, "float");
	} else if (streq(fnnm, "asserting_int")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_int_no_free_xfunc, parg1, pf, ti, cf, "int");
	} else if (streq(fnnm, "asserting_map")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_map_no_free_xfunc, parg1, pf, ti, cf, "map");
	} else if (streq(fnnm, "asserting_nonempty_map")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_nonempty_map_no_free_xfunc, parg1, pf, ti, cf,
			"nonempty_map");
	} else if (streq(fnnm, "asserting_not_empty")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_not_empty_no_free_xfunc, parg1, pf, ti, cf, "not_empty");
	} else if (streq(fnnm, "asserting_not_map")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_not_map_no_free_xfunc, parg1, pf, ti, cf, "not_map");
	} else if (streq(fnnm, "asserting_not_null")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_not_null_no_free_xfunc, parg1, pf, ti, cf, "not_null");
	} else if (streq(fnnm, "asserting_null")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_null_no_free_xfunc, parg1, pf, ti, cf, "null");
	} else if (streq(fnnm, "asserting_numeric")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_numeric_no_free_xfunc, parg1, pf, ti, cf, "numeric");
	} else if (streq(fnnm, "asserting_present")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_present_no_free_xfunc, parg1, pf, ti, cf, "present");
	} else if (streq(fnnm, "asserting_string")) {
		return rxval_evaluator_alloc_from_A_x_func(b_x_is_string_no_free_xfunc, parg1, pf, ti, cf, "string");

	} else if (streq(fnnm, "is_absent")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_absent_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_bool")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_boolean_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_boolean")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_boolean_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_empty")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_empty_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_empty_map")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_empty_map_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_float")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_float_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_int")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_int_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_map")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_map_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_nonempty_map")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_nonempty_map_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_not_empty")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_not_empty_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_not_map")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_not_map_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_not_null")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_not_null_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_null")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_null_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_numeric")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_numeric_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_present")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_present_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "is_string")) {
		return rxval_evaluator_alloc_from_x_x_func(b_x_is_string_xfunc, parg1, pf, ti, cf);

	} else if (streq(fnnm, "typeof")) {
		return rxval_evaluator_alloc_from_x_x_func(s_x_typeof_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "length")) {
		return rxval_evaluator_alloc_from_x_x_func(i_x_length_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "depth")) {
		return rxval_evaluator_alloc_from_x_x_func(i_x_depth_xfunc, parg1, pf, ti, cf);
	} else if (streq(fnnm, "leafcount")) {
		return rxval_evaluator_alloc_from_x_x_func(i_x_leafcount_xfunc, parg1, pf, ti, cf);
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static rxval_evaluator_t* fmgr_alloc_xevaluator_from_binary_func_name(char* fnnm,
	mlr_dsl_ast_node_t* parg1, mlr_dsl_ast_node_t* parg2,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/)
{
	if (streq(fnnm, "haskey")) {
		return rxval_evaluator_alloc_from_x_mx_func(b_xx_haskey_xfunc, parg1, parg2, pf, ti, cf);
	} else if (streq(fnnm, "splitnv")) {
		return rxval_evaluator_alloc_from_x_ss_func(m_ss_splitnv_xfunc, parg1, parg2, pf, ti, cf);
	} else if (streq(fnnm, "splitnvx")) {
		return rxval_evaluator_alloc_from_x_ss_func(m_ss_splitnvx_xfunc, parg1, parg2, pf, ti, cf);
	} else if (streq(fnnm, "joink")) {
		return rxval_evaluator_alloc_from_x_ms_func(s_ms_joink_xfunc, parg1, parg2, pf, ti, cf);
	} else if (streq(fnnm, "joinv")) {
		return rxval_evaluator_alloc_from_x_ms_func(s_ms_joinv_xfunc, parg1, parg2, pf, ti, cf);
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static rxval_evaluator_t* fmgr_alloc_xevaluator_from_ternary_func_name(char* fnnm,
	mlr_dsl_ast_node_t* parg1, mlr_dsl_ast_node_t* parg2, mlr_dsl_ast_node_t* parg3,
	fmgr_t* pf, int ti /*type_inferencing*/, int cf /*context_flags*/)
{
	if (streq(fnnm, "joinkv")) {
		return rxval_evaluator_alloc_from_x_mss_func(s_mss_joinkv_xfunc, parg1, parg2, parg3, pf, ti, cf);
	} else if (streq(fnnm, "splitkv")) {
		return rxval_evaluator_alloc_from_x_sss_func(m_sss_splitkv_xfunc, parg1, parg2, parg3, pf, ti, cf);
	} else if (streq(fnnm, "splitkvx")) {
		return rxval_evaluator_alloc_from_x_sss_func(m_sss_splitkvx_xfunc, parg1, parg2, parg3, pf, ti, cf);
	} else {
		return NULL;
	}
}

// ================================================================
// Return value is in scalar context.
static void resolve_func_callsite(fmgr_t* pfmgr, rval_evaluator_t* pev) {
	unresolved_func_callsite_state_t* pcallsite = pev->pvstate;

	rval_evaluator_t* pevaluator = construct_udf_callsite_evaluator(pfmgr, pcallsite);
	if (pevaluator != NULL) {
		// Struct assignment into the callsite space
		*pev = *pevaluator;
		free(pevaluator);
		return;
	}

	// Really there are map-in,map-out, map-in,scalar-out, and
	// scalar-in,scalar-out: and actually even more subtle, e.g. the join
	// functions take a mix of map and string arguments.  What we have
	// internally are builtin function evaluators (scalars only) and builtin
	// function xevaluators (at least one argument, and/or retval, is a map).
	rxval_evaluator_t* pxevaluator = construct_builtin_function_callsite_xevaluator(pfmgr, pcallsite);
	if (pxevaluator != NULL) {
		pevaluator = fmgr_alloc_eval_wrapping_xeval(pxevaluator);
		*pev = *pevaluator;
		free(pevaluator);
		return;
	}

	pevaluator = construct_builtin_function_callsite_evaluator(pfmgr, pcallsite);
	if (pevaluator != NULL) {
		*pev = *pevaluator;
		free(pevaluator);
		return;
	}

	fprintf(stderr, "Miller: unrecognized function name \"%s\".\n", pcallsite->function_name);
	exit(1);
}

// ----------------------------------------------------------------
// Return value is in map context.
static void resolve_func_xcallsite(fmgr_t* pfmgr, rxval_evaluator_t* pxev) {
	unresolved_func_callsite_state_t* pcallsite = pxev->pvstate;

	rxval_evaluator_t* pxevaluator = construct_udf_defsite_xevaluator(pfmgr, pcallsite);
	if (pxevaluator != NULL) {
		// Struct assignment into the callsite space
		*pxev = *pxevaluator;
		free(pxevaluator);
		return;
	}

	pxevaluator = construct_builtin_function_callsite_xevaluator(pfmgr, pcallsite);
	if (pxevaluator != NULL) {
		*pxev = *pxevaluator;
		free(pxevaluator);
		return;
	}

	rval_evaluator_t* pevaluator = construct_builtin_function_callsite_evaluator(pfmgr, pcallsite);
	pxevaluator = fmgr_alloc_xeval_wrapping_eval(pevaluator);
	if (pxevaluator != NULL) {
		*pxev = *pxevaluator;
		free(pxevaluator);
		return;
	}

	fprintf(stderr, "Miller: unrecognized function name \"%s\".\n", pcallsite->function_name);
	exit(1);
}
