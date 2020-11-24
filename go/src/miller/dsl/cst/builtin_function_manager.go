// ================================================================
// Adding a new builtin function:
// * New entry in BUILTIN_FUNCTION_LOOKUP_TABLE
// * Implement the function in mlrval_functions.go
// ================================================================

package cst

import (
	"fmt"
	"miller/types"
	"os"
)

// ================================================================
type BuiltinFunctionInfo struct {
	name string
	// class      string -- "math", "time", "typing", "maps", etc
	help                 string
	hasMultipleArities   bool
	minimumVariadicArity int
	zaryFunc             types.ZaryFunc
	unaryFunc            types.UnaryFunc
	contextualUnaryFunc  types.ContextualUnaryFunc // asserting_{typename}
	binaryFunc           types.BinaryFunc
	ternaryFunc          types.TernaryFunc
	variadicFunc         types.VariadicFunc
}

//// ----------------------------------------------------------------
//typedef enum _func_class_t {
//	FUNC_CLASS_ARITHMETIC,
//	FUNC_CLASS_MATH,
//	FUNC_CLASS_BOOLEAN,
//	FUNC_CLASS_STRING,
//	FUNC_CLASS_CONVERSION,
//	FUNC_CLASS_TYPING,
//	FUNC_CLASS_MAPS,
//	FUNC_CLASS_TIME
//} func_class_t;

// ================================================================
var BUILTIN_FUNCTION_LOOKUP_TABLE = []BuiltinFunctionInfo{

	// ----------------------------------------------------------------
	// Zary built-in functions
	{
		name:     "systime",
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystime,
	},
	{
		name:     "uptime",
		help:     "help string will go here",
		zaryFunc: types.MlrvalUptime,
	},
	{
		name:     "systimeint",
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystimeInt,
	},
	{
		name:     "urand",
		zaryFunc: types.MlrvalUrand,
	},
	{
		name:     "urand32",
		zaryFunc: types.MlrvalUrand32,
	},

	// ----------------------------------------------------------------
	// Multiple-arity built-in functions
	{
		name:               "+",
		unaryFunc:          types.MlrvalUnaryPlus,
		binaryFunc:         types.MlrvalBinaryPlus,
		hasMultipleArities: true,
	},
	{
		name:               "-",
		unaryFunc:          types.MlrvalUnaryMinus,
		binaryFunc:         types.MlrvalBinaryMinus,
		hasMultipleArities: true,
	},
	{
		name: "sec2gmt",
		help: `Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is.`,
		unaryFunc:          types.MlrvalSec2GMTUnary,
		binaryFunc:         types.MlrvalSec2GMTBinary,
		hasMultipleArities: true,
	},

	// ----------------------------------------------------------------
	// Unary built-in functions
	{
		name:      "~",
		unaryFunc: types.MlrvalBitwiseNOT,
	},
	{
		name:      "!",
		unaryFunc: types.MlrvalLogicalNOT,
	},

	{
		name:      "abs",
		help:      "Absolute value.",
		unaryFunc: types.MlrvalAbs,
	},
	{
		name:      "acos",
		help:      "Inverse trigonometric cosine.",
		unaryFunc: types.MlrvalAcos,
	},
	{
		name:      "acosh",
		help:      "Inverse hyperbolic cosine.",
		unaryFunc: types.MlrvalAcosh,
	},
	{
		name:      "asin",
		help:      "Inverse trigonometric sine.",
		unaryFunc: types.MlrvalAsin,
	},
	{
		name:      "asinh",
		help:      "Inverse hyperbolic sine.",
		unaryFunc: types.MlrvalAsinh,
	},
	{
		name:      "atan",
		help:      "One-argument arctangent.",
		unaryFunc: types.MlrvalAtan,
	},
	{
		name:      "atanh",
		help:      "Inverse hyperbolic tangent.",
		unaryFunc: types.MlrvalAtanh,
	},
	{
		name:      "cbrt",
		help:      "Cube root.",
		unaryFunc: types.MlrvalCbrt,
	},
	{
		name:      "ceil",
		help:      "Ceiling: nearest integer at or above.",
		unaryFunc: types.MlrvalCeil,
	},
	{
		name:      "cos",
		help:      "Trigonometric cosine.",
		unaryFunc: types.MlrvalCos,
	},
	{
		name:      "cosh",
		help:      "Hyperbolic cosine.",
		unaryFunc: types.MlrvalCosh,
	},
	{
		name:      "erf",
		help:      "Error function.",
		unaryFunc: types.MlrvalErf,
	},
	{
		name:      "erfc",
		help:      "Complementary error function.",
		unaryFunc: types.MlrvalErfc,
	},
	{
		name:      "exp",
		help:      "Exponential function e**x.",
		unaryFunc: types.MlrvalExp,
	},
	{
		name:      "expm1",
		help:      "e**x - 1.",
		unaryFunc: types.MlrvalExpm1,
	},
	{
		name:      "floor",
		help:      "Floor: nearest integer at or below.",
		unaryFunc: types.MlrvalFloor,
	},
	{
		name: "invqnorm",
		help: `Inverse of normal cumulative distribution function.
Note that invqorm(urand()) is normally distributed.`,
		unaryFunc: types.MlrvalInvqnorm,
	},
	{
		name:      "log",
		help:      "Natural (base-e) logarithm.",
		unaryFunc: types.MlrvalLog,
	},
	{
		name:      "log10",
		help:      "Base-10 logarithm.",
		unaryFunc: types.MlrvalLog10,
	},
	{
		name:      "log1p",
		help:      "log(1-x).",
		unaryFunc: types.MlrvalLog1p,
	},
	{
		name:      "qnorm",
		help:      `Normal cumulative distribution function.`,
		unaryFunc: types.MlrvalQnorm,
	},
	{
		name:      "round",
		help:      "Round to nearest integer.",
		unaryFunc: types.MlrvalRound,
	},
	{
		name:      "sgn",
		help:      ` +1, 0, -1 for positive, zero, negative input respectively.`,
		unaryFunc: types.MlrvalSgn,
	},
	{
		name:      "sin",
		help:      "Trigonometric sine.",
		unaryFunc: types.MlrvalSin,
	},
	{
		name:      "sinh",
		help:      "Hyperbolic sine.",
		unaryFunc: types.MlrvalSinh,
	},
	{
		name:      "sqrt",
		help:      "Square root.",
		unaryFunc: types.MlrvalSqrt,
	},
	{
		name:      "tan",
		help:      "Trigonometric tangent.",
		unaryFunc: types.MlrvalTan,
	},
	{
		name:      "tanh",
		help:      "Hyperbolic tangent.",
		unaryFunc: types.MlrvalTanh,
	},
	{
		name:      "bitcount",
		help:      "Count of 1-bits.",
		unaryFunc: types.MlrvalBitCount,
	},
	{
		name:      "clean_whitespace",
		help:      "Same as collapse_whitespace and strip.",
		unaryFunc: types.MlrvalCleanWhitespace,
	},
	{
		name:      "collapse_whitespace",
		help:      "Strip repeated whitespace from string.",
		unaryFunc: types.MlrvalCollapseWhitespace,
	},
	{
		name:      "length",
		help:      "Counts number of top-level entries in array/map. Scalars have length 1.",
		unaryFunc: types.MlrvalLength,
	},
	{
		name:      "lstrip",
		help:      "Strip leading whitespace from string.",
		unaryFunc: types.MlrvalLStrip,
	},
	{
		name:      "rstrip",
		help:      "Strip trailing whitespace from string.",
		unaryFunc: types.MlrvalRStrip,
	},

	{
		name:      "string",
		help:      "Convert int/float/bool/string/array/map to string.",
		unaryFunc: types.MlrvalToString,
	},
	{
		name:      "int",
		help:      "Convert int/float/bool/string to int.",
		unaryFunc: types.MlrvalToInt,
	},
	{
		name:      "float",
		help:      "Convert int/float/bool/string to float.",
		unaryFunc: types.MlrvalToFloat,
	},
	{
		name:      "boolean",
		help:      "Convert int/float/bool/string to boolean.",
		unaryFunc: types.MlrvalToBoolean,
	},
	{
		name:      "hexfmt",
		help:      `Convert int to hex string, e.g. 255 to "0xff".`,
		unaryFunc: types.MlrvalHexfmt,
	},

	{
		name:      "strip",
		help:      "Strip leading and trailing whitespace from string.",
		unaryFunc: types.MlrvalStrip,
	},
	{
		name:      "strlen",
		help:      "String length.",
		unaryFunc: types.MlrvalStrlen,
	},
	{
		name:      "tolower",
		help:      "Convert string to lowercase.",
		unaryFunc: types.MlrvalToLower,
	},
	{
		name:      "toupper",
		help:      "Convert string to uppercase.",
		unaryFunc: types.MlrvalToUpper,
	},
	{
		name:      "capitalize",
		help:      "Convert string's first character to uppercase.",
		unaryFunc: types.MlrvalCapitalize,
	},

	{
		name:      "is_absent",
		help:      "False if field is present in input, true otherwise",
		unaryFunc: types.MlrvalIsAbsent,
	},
	{
		name:      "is_error",
		help:      "True if if argument is an error, such as taking string length of an integer.",
		unaryFunc: types.MlrvalIsError,
	},
	{
		name:      "is_bool",
		help:      "True if field is present with boolean value. Synonymous with is_boolean.",
		unaryFunc: types.MlrvalIsBool,
	},
	{
		name:      "is_boolean",
		help:      "True if field is present with boolean value. Synonymous with is_bool.",
		unaryFunc: types.MlrvalIsBoolean,
	},
	{
		name:      "is_empty",
		help:      "True if field is present in input with empty string value, false otherwise.",
		unaryFunc: types.MlrvalIsEmpty,
	},
	{
		name:      "is_empty_map",
		help:      "True if argument is a map which is empty.",
		unaryFunc: types.MlrvalIsEmptyMap,
	},
	{
		name:      "is_float",
		help:      "True if field is present with value inferred to be float",
		unaryFunc: types.MlrvalIsFloat,
	},
	{
		name:      "is_int",
		help:      "True if field is present with value inferred to be int",
		unaryFunc: types.MlrvalIsInt,
	},
	{
		name:      "is_map",
		help:      "True if argument is a map.",
		unaryFunc: types.MlrvalIsMap,
	},
	{
		name:      "is_array",
		help:      "True if argument is an array.",
		unaryFunc: types.MlrvalIsArray,
	},
	{
		name:      "is_nonempty_map",
		help:      "True if argument is a map which is non-empty.",
		unaryFunc: types.MlrvalIsNonEmptyMap,
	},
	{
		name:      "is_not_empty",
		help:      "False if field is present in input with empty value, true otherwise",
		unaryFunc: types.MlrvalIsNotEmpty,
	},
	{
		name:      "is_not_map",
		help:      "True if argument is not a map.",
		unaryFunc: types.MlrvalIsNotMap,
	},
	{
		name:      "is_not_array",
		help:      "True if argument is not an array.",
		unaryFunc: types.MlrvalIsNotArray,
	},
	{
		name:      "is_not_null",
		help:      "False if argument is null (empty or absent), true otherwise.",
		unaryFunc: types.MlrvalIsNotNull,
	},
	{
		name:      "is_null",
		help:      "True if argument is null (empty or absent), false otherwise.",
		unaryFunc: types.MlrvalIsNull,
	},
	{
		name:      "is_numeric",
		help:      "True if field is present with value inferred to be int or float",
		unaryFunc: types.MlrvalIsNumeric,
	},
	{
		name:      "is_present",
		help:      "True if field is present in input, false otherwise.",
		unaryFunc: types.MlrvalIsPresent,
	},
	{
		name:      "is_string",
		help:      "True if field is present with string (including empty-string) value",
		unaryFunc: types.MlrvalIsString,
	},

	{
		name: "asserting_absent",
		help: `Aborts with an error if is_absent on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingAbsent,
	},
	{
		name: "asserting_error",
		help: `Aborts with an error if is_error on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingError,
	},
	{
		name: "asserting_bool",
		help: `Aborts with an error if is_bool on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingBool,
	},
	{
		name: "asserting_boolean",
		help: `Aborts with an error if is_boolean on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingBoolean,
	},
	{
		name: "asserting_empty",
		help: `Aborts with an error if is_empty on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingEmpty,
	},
	{
		name: "asserting_empty_map",
		help: `Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingEmptyMap,
	},
	{
		name: "asserting_float",
		help: `Aborts with an error if is_float on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingFloat,
	},
	{
		name: "asserting_int",
		help: `Aborts with an error if is_int on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingInt,
	},
	{
		name: "asserting_map",
		help: `Aborts with an error if is_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingMap,
	},
	{
		name: "asserting_array",
		help: `Aborts with an error if is_array on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingArray,
	},
	{
		name: "asserting_nonempty_map",
		help: `Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNonEmptyMap,
	},
	{
		name: "asserting_not_empty",
		help: `Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotEmpty,
	},
	{
		name: "asserting_not_map",
		help: `Aborts with an error if is_not_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotMap,
	},
	{
		name: "asserting_not_array",
		help: `Aborts with an error if is_not_array on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotArray,
	},
	{
		name: "asserting_not_null",
		help: `Aborts with an error if is_not_null on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotNull,
	},
	{
		name: "asserting_null",
		help: `Aborts with an error if is_null on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNull,
	},
	{
		name: "asserting_numeric",
		help: `Aborts with an error if is_numeric on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNumeric,
	},
	{
		name: "asserting_present",
		help: `Aborts with an error if is_present on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingPresent,
	},
	{
		name: "asserting_string",
		help: `Aborts with an error if is_string on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingString,
	},

	{
		name:      "typeof",
		help:      "Convert argument to type of argument (e.g. \"str\"). For debug.",
		unaryFunc: types.MlrvalTypeof,
	},
	{
		name:      "depth",
		help:      "Prints maximum depth of map/array. Scalars have depth 0.",
		unaryFunc: types.MlrvalDepth,
	},
	{
		name: "leafcount",
		help: `Counts total number of terminal values in map/array. For single-level
map/array, same as length.`,
		unaryFunc: types.MlrvalLeafCount,
	},
	{
		name:      "keys",
		help:      "Returns array of keys of map or array",
		unaryFunc: types.MlrvalKeys,
	},
	{
		name:      "values",
		help:      "Returns array of keys of map or array -- in the latter case, returns a copy of the array",
		unaryFunc: types.MlrvalValues,
	},

	// ----------------------------------------------------------------
	// Binary built-in functions
	{
		name:       ".",
		binaryFunc: types.MlrvalDot,
	},
	{
		name:       "*",
		binaryFunc: types.MlrvalTimes,
	},
	{
		name:       "/",
		binaryFunc: types.MlrvalDivide,
	},
	{
		name:       "//",
		binaryFunc: types.MlrvalIntDivide,
	},
	{
		name:       "**",
		binaryFunc: types.MlrvalPow,
	},
	{
		name:       "pow",
		binaryFunc: types.MlrvalPow,
	},
	{
		name:       ".+",
		binaryFunc: types.MlrvalDotPlus,
	},
	{
		name:       ".-",
		binaryFunc: types.MlrvalDotMinus,
	},
	{
		name:       ".*",
		binaryFunc: types.MlrvalDotTimes,
	},
	{
		name:       "./",
		binaryFunc: types.MlrvalDotDivide,
	},
	{
		name:       "%",
		binaryFunc: types.MlrvalModulus,
	},

	{
		name:       "==",
		binaryFunc: types.MlrvalEquals,
	},
	{
		name:       "!=",
		binaryFunc: types.MlrvalNotEquals,
	},
	{
		name:       ">",
		binaryFunc: types.MlrvalGreaterThan,
	},
	{
		name:       ">=",
		binaryFunc: types.MlrvalGreaterThanOrEquals,
	},
	{
		name:       "<",
		binaryFunc: types.MlrvalLessThan,
	},
	{
		name:       "<=",
		binaryFunc: types.MlrvalLessThanOrEquals,
	},
	{
		name: "=~",
		help: `String (left-hand side) matches regex (right-hand
side), e.g. '$name =~ "^a.*b$"'.`,
		binaryFunc: types.MlrvalStringMatchesRegexp,
	},
	{
		name: "!=~",
		help: `String (left-hand side) does not match regex (right-hand
side), e.g. '$name !=~ "^a.*b$"'.`,
		binaryFunc: types.MlrvalStringDoesNotMatchRegexp,
	},

	{
		name:       "&&",
		binaryFunc: BinaryShortCircuitPlaceholder,
	},
	{
		name:       "||",
		binaryFunc: BinaryShortCircuitPlaceholder,
	},
	{
		name:       "??",
		binaryFunc: BinaryShortCircuitPlaceholder,
	},
	{
		name:       "^^",
		binaryFunc: types.MlrvalLogicalXOR,
	},
	{
		name:       "&",
		binaryFunc: types.MlrvalBitwiseAND,
	},
	{
		name:       "|",
		binaryFunc: types.MlrvalBitwiseOR,
	},
	{
		name:       "^",
		binaryFunc: types.MlrvalBitwiseXOR,
	},
	{
		name:       "<<",
		binaryFunc: types.MlrvalLeftShift,
	},
	{
		name:       ">>",
		binaryFunc: types.MlrvalSignedRightShift,
	},
	{
		name:       ">>>",
		binaryFunc: types.MlrvalUnsignedRightShift,
	},
	{
		name: "roundm",
		help: `Round to nearest multiple of m: roundm($x,$m) is
the same as round($x/$m)*$m.`,
		binaryFunc: types.MlrvalRoundm,
	},
	{
		name: "fmtnum",
		help: `Convert int/float/bool to string using
printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.`,
		binaryFunc: types.MlrvalFmtNum,
	},

	{
		name: "urandint",
		help: `Integer uniformly distributed between inclusive
integer endpoints.`,
		binaryFunc: types.MlrvalUrandInt,
	},
	{
		name:       "urandrange",
		help:       `Floating-point numbers uniformly distributed on the interval [a, b).`,
		binaryFunc: types.MlrvalUrandRange,
	},

	{
		name:       "atan2",
		help:       "Two-argument arctangent.",
		binaryFunc: types.MlrvalAtan2,
	},
	{
		name:       "truncate",
		binaryFunc: types.MlrvalTruncate,
	},
	{
		name: "haskey",
		help: `True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.`,
		binaryFunc: types.MlrvalHasKey,
	},

	{
		name: "joink",
		help: `Makes string from map/array keys. Examples:
joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
joink([1,2,3], ",") = "1,2,3".`,
		binaryFunc: types.MlrvalJoinK,
	},
	{
		name: "joinv",
		help: `Makes string from map/array values.
joinv([3,4,5], ",") = "3,4,5"
joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"`,
		binaryFunc: types.MlrvalJoinV,
	},
	{
		name: "splita",
		help: `Splits string into array with type inference. Example:
splita("3,4,5", ",") = [3,4,5]`,
		binaryFunc: types.MlrvalSplitA,
	},
	{
		name: "splitax",
		help: `Splits string into array without type inference. Example:
splita("3,4,5", ",") = ["3","4","5"]`,
		binaryFunc: types.MlrvalSplitAX,
	},
	{
		name:       "append",
		help:       "Appends second argument to end of first argument, which must be an array.",
		binaryFunc: types.MlrvalAppend,
	},

	//pow (class=math #args=2): Exponentiation; same as **.
	//roundm (class=math #args=2): Round to nearest multiple of m: roundm($x,$m) is
	//urandrange (class=math #args=2): Floating-point numbers uniformly distributed on the interval [a, b).
	//urandint (class=math #args=2): Integer uniformly distributed between inclusive
	//atan2 (class=math #args=2): Two-argument arctangent.

	// Ternary built-in functions
	//logifit (class=math #args=3): Given m and b from logistic regression, compute
	//madd (class=math #args=3): a + b mod m (integers)
	//mexp (class=math #args=3): a ** b mod m (integers)
	//mmul (class=math #args=3): a * b mod m (integers)
	//msub (class=math #args=3): a - b mod m (integers)
	{
		name:        "?:",
		ternaryFunc: TernaryShortCircuitPlaceholder,
	},
	{
		name:        "ssub",
		ternaryFunc: types.MlrvalSsub,
	},
	{
		name:        "sub",
		ternaryFunc: types.MlrvalSub,
	},
	{
		name:        "gsub",
		ternaryFunc: types.MlrvalGsub,
	},
	{
		name: "substr",
		help: `substr(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len.`,
		ternaryFunc: types.MlrvalSubstr,
	},
	{
		name:        "madd",
		help:        `a + b mod m (integers)`,
		ternaryFunc: types.MlrvalModAdd,
	},
	{
		name:        "msub",
		help:        `a - b mod m (integers)`,
		ternaryFunc: types.MlrvalModSub,
	},
	{
		name:        "mmul",
		help:        `a * b mod m (integers)`,
		ternaryFunc: types.MlrvalModMul,
	},
	{
		name:        "mexp",
		help:        `a ** b mod m (integers)`,
		ternaryFunc: types.MlrvalModExp,
	},
	{
		name: "logifit",
		help: ` Given m and b from logistic regression, compute fit:
$yhat=logifit($x,$m,$b).`,
		ternaryFunc: types.MlrvalLogifit,
	},
	{
		name: "joinkv",
		help: `Makes string from map/array key-value pairs. Examples:
joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"`,
		ternaryFunc: types.MlrvalJoinKV,
	},

	{
		name: "splitkv",
		help: `Splits string by separators into map with type inference. Example:
splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}`,
		ternaryFunc: types.MlrvalSplitKV,
	},
	{
		name: "splitkvx",
		help: `Splits string by separators into map without type inference (keys and
values are strings). Example:
splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}`,
		ternaryFunc: types.MlrvalSplitKVX,
	},
	{
		name: "splitnv",
		help: `Splits string by separator into integer-indexed map with type inference. Example:
splitnv("a=3,b=4,c=5", "=", ",") = {"1":3,"2":4,"3":5}`,
		ternaryFunc: types.MlrvalSplitNV,
	},
	{
		name: "splitnvx",
		help: `Splits string by separator into integer-indexed map without type
inference (values are strings). Example:
splitnvx("a=3,b=4,c=5", "=", ",") = {"1":"3","2":"4","3":"5"}`,
		ternaryFunc: types.MlrvalSplitNVX,
	},
	{
		name: "splitak",
		help: `Splits keys out of string representation of map into array. Example:
splitak("a=3,b=4,c=5", "=", ",") = ["a","b","c"]`,
		ternaryFunc: types.MlrvalSplitAK,
	},
	{
		name: "splitav",
		help: `Splits type-inferred values out of string representation of map into array. Example:
splitav("a=3,b=4,c=5", "=", ",") = [3,4,5]`,
		ternaryFunc: types.MlrvalSplitAV,
	},
	{
		name: "splitavx",
		help: `Splits non-type-inferred values out of string representation of map into array. Example:
splitav("a=3,b=4,c=5", "=", ",") = ["3","4","5"]`,
		ternaryFunc: types.MlrvalSplitAVX,
	},

	// Variadic built-in functions
	{
		name:         "max",
		variadicFunc: types.MlrvalVariadicMax,
	},
	{
		name:         "min",
		variadicFunc: types.MlrvalVariadicMin,
	},

	{
		name: "mapselect",
		help: `Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.
E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.`,
		variadicFunc:         types.MlrvalMapSelect,
		minimumVariadicArity: 1,
	},
	{
		name: "mapexcept",
		help: `Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.
E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.`,
		variadicFunc:         types.MlrvalMapExcept,
		minimumVariadicArity: 1,
	},
	{
		name: "mapsum",
		help: `With 0 args, returns empty map. With >= 1 arg, returns a map with
key-value pairs from all arguments. Rightmost collisions win, e.g.
'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.`,
		variadicFunc: types.MlrvalMapSum,
	},
	{
		name: "mapdiff",
		help: `With 0 args, returns empty map. With 1 arg, returns copy of arg.
With 2 or more, returns copy of arg 1 with all keys from any of remaining
argument maps removed.`,
		variadicFunc: types.MlrvalMapDiff,
	},
}

// ================================================================
type BuiltinFunctionManager struct {
	// We need both the array and the hashmap since Go maps are not
	// insertion-order-preserving: to produce a sensical help-all-functions
	// list, etc., we want the original ordering.
	lookupTable *[]BuiltinFunctionInfo
	hashTable   map[string]*BuiltinFunctionInfo
}

func NewBuiltinFunctionManager() *BuiltinFunctionManager {
	// TODO: temp -- one big one -- pending UDFs
	lookupTable := &BUILTIN_FUNCTION_LOOKUP_TABLE
	hashTable := hashifyLookupTable(lookupTable)
	return &BuiltinFunctionManager{
		lookupTable: lookupTable,
		hashTable:   hashTable,
	}
}

func (this *BuiltinFunctionManager) LookUp(functionName string) *BuiltinFunctionInfo {
	return this.hashTable[functionName]
}

func hashifyLookupTable(lookupTable *[]BuiltinFunctionInfo) map[string]*BuiltinFunctionInfo {
	hashTable := make(map[string]*BuiltinFunctionInfo)
	for _, builtinFunctionInfo := range *lookupTable {
		// Each function name should appear only once in the table.  If it has
		// multiple arities (e.g. unary and binary "-") there should be
		// multiple function-pointers in a single row.
		if hashTable[builtinFunctionInfo.name] != nil {
			fmt.Fprintf(
				os.Stderr,
				"Internal coding error: function name \"%s\" is non-unique",
				builtinFunctionInfo.name,
			)
			os.Exit(1)
		}
		clone := builtinFunctionInfo
		hashTable[builtinFunctionInfo.name] = &clone
	}
	return hashTable
}

// ----------------------------------------------------------------
func (this *BuiltinFunctionManager) ListBuiltinFunctionsRaw(o *os.File) {
	for _, builtinFunctionInfo := range *this.lookupTable {
		fmt.Fprintln(o, builtinFunctionInfo.name)
	}
}

// ----------------------------------------------------------------
func (this *BuiltinFunctionManager) ListBuiltinFunctionUsages(o *os.File) {
	for i, builtinFunctionInfo := range *this.lookupTable {
		if i > 0 {
			fmt.Fprintln(o)
		}
		fmt.Fprintf(o, "%-20s  %s\n", builtinFunctionInfo.name, builtinFunctionInfo.help)
	}
}

// ================================================================
// This is a singleton so the online-help functions can query it for listings,
// online help, etc.
var BuiltinFunctionManagerInstance *BuiltinFunctionManager = NewBuiltinFunctionManager()
