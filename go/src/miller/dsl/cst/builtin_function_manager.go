// ================================================================
// Adding a new builtin function:
// * New entry in BUILTIN_FUNCTION_LOOKUP_TABLE
// * Implement the function in mlrval_functions.go
// ================================================================

package cst

import (
	"fmt"
	"os"
	"strings"

	"miller/lib"
	"miller/types"
)

type TFunctionClass string

const (
	FUNC_CLASS_ARITHMETIC  TFunctionClass = "arithmetic"
	FUNC_CLASS_MATH                       = "math"
	FUNC_CLASS_BOOLEAN                    = "boolean"
	FUNC_CLASS_STRING                     = "string"
	FUNC_CLASS_CONVERSION                 = "conversion"
	FUNC_CLASS_TYPING                     = "typing"
	FUNC_CLASS_COLLECTIONS                = "maps/arrays"
	FUNC_CLASS_TIME                       = "time"
)

// ================================================================
type BuiltinFunctionInfo struct {
	name                 string
	class                TFunctionClass
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

// ================================================================
var BUILTIN_FUNCTION_LOOKUP_TABLE = []BuiltinFunctionInfo{

	// ----------------------------------------------------------------
	// FUNC_CLASS_ARITHMETIC
	{
		name:               "+",
		class:              FUNC_CLASS_ARITHMETIC,
		unaryFunc:          types.MlrvalUnaryPlus,
		binaryFunc:         types.MlrvalBinaryPlus,
		hasMultipleArities: true,
	},

	{
		name:               "-",
		class:              FUNC_CLASS_ARITHMETIC,
		unaryFunc:          types.MlrvalUnaryMinus,
		binaryFunc:         types.MlrvalBinaryMinus,
		hasMultipleArities: true,
	},

	{
		name:       "*",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalTimes,
	},

	{
		name:       "/",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalDivide,
	},

	{
		name:       "//",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalIntDivide,
	},

	{
		name:       "**",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalPow,
	},

	{
		name:       "pow",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalPow,
	},

	{
		name:       ".+",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalDotPlus,
	},

	{
		name:       ".-",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalDotMinus,
	},

	{
		name:       ".*",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalDotTimes,
	},

	{
		name:       "./",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalDotDivide,
	},

	{
		name:       "%",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalModulus,
	},

	{
		name:      "~",
		class:     FUNC_CLASS_ARITHMETIC,
		unaryFunc: types.MlrvalBitwiseNOT,
	},

	{
		name:       "&",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalBitwiseAND,
	},

	{
		name:       "|",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalBitwiseOR,
	},

	{
		name:       "^",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalBitwiseXOR,
	},

	{
		name:       "<<",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalLeftShift,
	},

	{
		name:       ">>",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalSignedRightShift,
	},

	{
		name:       ">>>",
		class:      FUNC_CLASS_ARITHMETIC,
		binaryFunc: types.MlrvalUnsignedRightShift,
	},

	{
		name:      "bitcount",
		class:     FUNC_CLASS_ARITHMETIC,
		help:      "Count of 1-bits.",
		unaryFunc: types.MlrvalBitCount,
	},

	{
		name:        "madd",
		class:       FUNC_CLASS_ARITHMETIC,
		help:        `a + b mod m (integers)`,
		ternaryFunc: types.MlrvalModAdd,
	},

	{
		name:        "msub",
		class:       FUNC_CLASS_ARITHMETIC,
		help:        `a - b mod m (integers)`,
		ternaryFunc: types.MlrvalModSub,
	},

	{
		name:        "mmul",
		class:       FUNC_CLASS_ARITHMETIC,
		help:        `a * b mod m (integers)`,
		ternaryFunc: types.MlrvalModMul,
	},

	{
		name:        "mexp",
		class:       FUNC_CLASS_ARITHMETIC,
		help:        `a ** b mod m (integers)`,
		ternaryFunc: types.MlrvalModExp,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_BOOLEAN
	{
		name:      "!",
		class:     FUNC_CLASS_BOOLEAN,
		unaryFunc: types.MlrvalLogicalNOT,
	},

	{
		name:       "==",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalEquals,
	},

	{
		name:       "!=",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalNotEquals,
	},

	{
		name:       ">",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalGreaterThan,
	},

	{
		name:       ">=",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalGreaterThanOrEquals,
	},

	{
		name:       "<",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalLessThan,
	},

	{
		name:       "<=",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalLessThanOrEquals,
	},

	{
		name:  "=~",
		class: FUNC_CLASS_BOOLEAN,
		help: `String (left-hand side) matches regex (right-hand
side), e.g. '$name =~ "^a.*b$"'.`,
		binaryFunc: types.MlrvalStringMatchesRegexp,
	},

	{
		name:  "!=~",
		class: FUNC_CLASS_BOOLEAN,
		help: `String (left-hand side) does not match regex (right-hand
side), e.g. '$name !=~ "^a.*b$"'.`,
		binaryFunc: types.MlrvalStringDoesNotMatchRegexp,
	},

	{
		name:       "&&",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: BinaryShortCircuitPlaceholder,
	},

	{
		name:       "||",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: BinaryShortCircuitPlaceholder,
	},

	{
		name:       "??",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: BinaryShortCircuitPlaceholder,
	},

	{
		name:       "???",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: BinaryShortCircuitPlaceholder,
	},

	{
		name:       "^^",
		class:      FUNC_CLASS_BOOLEAN,
		binaryFunc: types.MlrvalLogicalXOR,
	},

	{
		name:        "?:",
		class:       FUNC_CLASS_BOOLEAN,
		ternaryFunc: TernaryShortCircuitPlaceholder,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_STRING

	{
		name:       ".",
		class:      FUNC_CLASS_STRING,
		binaryFunc: types.MlrvalDot,
	},

	{
		name:      "capitalize",
		class:     FUNC_CLASS_STRING,
		help:      "Convert string's first character to uppercase.",
		unaryFunc: types.MlrvalCapitalize,
	},

	{
		name:      "clean_whitespace",
		class:     FUNC_CLASS_STRING,
		help:      "Same as collapse_whitespace and strip.",
		unaryFunc: types.MlrvalCleanWhitespace,
	},

	{
		name:      "collapse_whitespace",
		class:     FUNC_CLASS_STRING,
		help:      "Strip repeated whitespace from string.",
		unaryFunc: types.MlrvalCollapseWhitespace,
	},

	{
		name:        "gsub",
		class:       FUNC_CLASS_STRING,
		ternaryFunc: types.MlrvalGsub,
	},

	{
		name:      "lstrip",
		class:     FUNC_CLASS_STRING,
		help:      "Strip leading whitespace from string.",
		unaryFunc: types.MlrvalLStrip,
	},

	{
		name:      "rstrip",
		class:     FUNC_CLASS_STRING,
		help:      "Strip trailing whitespace from string.",
		unaryFunc: types.MlrvalRStrip,
	},

	{
		name:      "strip",
		class:     FUNC_CLASS_STRING,
		help:      "Strip leading and trailing whitespace from string.",
		unaryFunc: types.MlrvalStrip,
	},

	{
		name:      "strlen",
		class:     FUNC_CLASS_STRING,
		help:      "String length.",
		unaryFunc: types.MlrvalStrlen,
	},

	{
		name:        "ssub",
		class:       FUNC_CLASS_STRING,
		ternaryFunc: types.MlrvalSsub,
	},

	{
		name:        "sub",
		class:       FUNC_CLASS_STRING,
		ternaryFunc: types.MlrvalSub,
	},

	{
		name:  "substr",
		class: FUNC_CLASS_STRING,
		help: `substr(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len.`,
		ternaryFunc: types.MlrvalSubstr,
	},

	{
		name:      "tolower",
		class:     FUNC_CLASS_STRING,
		help:      "Convert string to lowercase.",
		unaryFunc: types.MlrvalToLower,
	},

	{
		name:      "toupper",
		class:     FUNC_CLASS_STRING,
		help:      "Convert string to uppercase.",
		unaryFunc: types.MlrvalToUpper,
	},

	{
		name:       "truncate",
		class:      FUNC_CLASS_STRING,
		binaryFunc: types.MlrvalTruncate,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_MATH

	{
		name:      "abs",
		class:     FUNC_CLASS_MATH,
		help:      "Absolute value.",
		unaryFunc: types.MlrvalAbs,
	},

	{
		name:      "acos",
		class:     FUNC_CLASS_MATH,
		help:      "Inverse trigonometric cosine.",
		unaryFunc: types.MlrvalAcos,
	},

	{
		name:      "acosh",
		class:     FUNC_CLASS_MATH,
		help:      "Inverse hyperbolic cosine.",
		unaryFunc: types.MlrvalAcosh,
	},

	{
		name:      "asin",
		class:     FUNC_CLASS_MATH,
		help:      "Inverse trigonometric sine.",
		unaryFunc: types.MlrvalAsin,
	},

	{
		name:      "asinh",
		class:     FUNC_CLASS_MATH,
		help:      "Inverse hyperbolic sine.",
		unaryFunc: types.MlrvalAsinh,
	},

	{
		name:      "atan",
		class:     FUNC_CLASS_MATH,
		help:      "One-argument arctangent.",
		unaryFunc: types.MlrvalAtan,
	},

	{
		name:      "atanh",
		class:     FUNC_CLASS_MATH,
		help:      "Inverse hyperbolic tangent.",
		unaryFunc: types.MlrvalAtanh,
	},

	{
		name:       "atan2",
		class:      FUNC_CLASS_MATH,
		help:       "Two-argument arctangent.",
		binaryFunc: types.MlrvalAtan2,
	},

	{
		name:      "cbrt",
		class:     FUNC_CLASS_MATH,
		help:      "Cube root.",
		unaryFunc: types.MlrvalCbrt,
	},

	{
		name:      "ceil",
		class:     FUNC_CLASS_MATH,
		help:      "Ceiling: nearest integer at or above.",
		unaryFunc: types.MlrvalCeil,
	},

	{
		name:      "cos",
		class:     FUNC_CLASS_MATH,
		help:      "Trigonometric cosine.",
		unaryFunc: types.MlrvalCos,
	},

	{
		name:      "cosh",
		class:     FUNC_CLASS_MATH,
		help:      "Hyperbolic cosine.",
		unaryFunc: types.MlrvalCosh,
	},

	{
		name:      "erf",
		class:     FUNC_CLASS_MATH,
		help:      "Error function.",
		unaryFunc: types.MlrvalErf,
	},

	{
		name:      "erfc",
		class:     FUNC_CLASS_MATH,
		help:      "Complementary error function.",
		unaryFunc: types.MlrvalErfc,
	},

	{
		name:      "exp",
		class:     FUNC_CLASS_MATH,
		help:      "Exponential function e**x.",
		unaryFunc: types.MlrvalExp,
	},

	{
		name:      "expm1",
		class:     FUNC_CLASS_MATH,
		help:      "e**x - 1.",
		unaryFunc: types.MlrvalExpm1,
	},

	{
		name:      "floor",
		class:     FUNC_CLASS_MATH,
		help:      "Floor: nearest integer at or below.",
		unaryFunc: types.MlrvalFloor,
	},

	{
		name:  "invqnorm",
		class: FUNC_CLASS_MATH,
		help: `Inverse of normal cumulative distribution function.
Note that invqorm(urand()) is normally distributed.`,
		unaryFunc: types.MlrvalInvqnorm,
	},

	{
		name:      "log",
		class:     FUNC_CLASS_MATH,
		help:      "Natural (base-e) logarithm.",
		unaryFunc: types.MlrvalLog,
	},

	{
		name:      "log10",
		class:     FUNC_CLASS_MATH,
		help:      "Base-10 logarithm.",
		unaryFunc: types.MlrvalLog10,
	},

	{
		name:      "log1p",
		class:     FUNC_CLASS_MATH,
		help:      "log(1-x).",
		unaryFunc: types.MlrvalLog1p,
	},

	{
		name:  "logifit",
		class: FUNC_CLASS_MATH,
		help: ` Given m and b from logistic regression, compute fit:
$yhat=logifit($x,$m,$b).`,
		ternaryFunc: types.MlrvalLogifit,
	},

	{
		name:         "max",
		class:        FUNC_CLASS_MATH,
		variadicFunc: types.MlrvalVariadicMax,
	},

	{
		name:         "min",
		class:        FUNC_CLASS_MATH,
		variadicFunc: types.MlrvalVariadicMin,
	},

	{
		name:      "qnorm",
		class:     FUNC_CLASS_MATH,
		help:      `Normal cumulative distribution function.`,
		unaryFunc: types.MlrvalQnorm,
	},

	{
		name:      "round",
		class:     FUNC_CLASS_MATH,
		help:      "Round to nearest integer.",
		unaryFunc: types.MlrvalRound,
	},

	{
		name:      "sgn",
		class:     FUNC_CLASS_MATH,
		help:      ` +1, 0, -1 for positive, zero, negative input respectively.`,
		unaryFunc: types.MlrvalSgn,
	},

	{
		name:      "sin",
		class:     FUNC_CLASS_MATH,
		help:      "Trigonometric sine.",
		unaryFunc: types.MlrvalSin,
	},

	{
		name:      "sinh",
		class:     FUNC_CLASS_MATH,
		help:      "Hyperbolic sine.",
		unaryFunc: types.MlrvalSinh,
	},

	{
		name:      "sqrt",
		class:     FUNC_CLASS_MATH,
		help:      "Square root.",
		unaryFunc: types.MlrvalSqrt,
	},

	{
		name:      "tan",
		class:     FUNC_CLASS_MATH,
		help:      "Trigonometric tangent.",
		unaryFunc: types.MlrvalTan,
	},

	{
		name:      "tanh",
		class:     FUNC_CLASS_MATH,
		help:      "Hyperbolic tangent.",
		unaryFunc: types.MlrvalTanh,
	},

	{
		name:  "roundm",
		class: FUNC_CLASS_MATH,
		help: `Round to nearest multiple of m: roundm($x,$m) is
the same as round($x/$m)*$m.`,
		binaryFunc: types.MlrvalRoundm,
	},

	{
		name:     "urand",
		class:    FUNC_CLASS_MATH,
		zaryFunc: types.MlrvalUrand,
	},

	{
		name:  "urandint",
		class: FUNC_CLASS_MATH,
		help: `Integer uniformly distributed between inclusive
integer endpoints.`,
		binaryFunc: types.MlrvalUrandInt,
	},

	{
		name:       "urandrange",
		class:      FUNC_CLASS_MATH,
		help:       `Floating-point numbers uniformly distributed on the interval [a, b).`,
		binaryFunc: types.MlrvalUrandRange,
	},

	{
		name:     "urand32",
		class:    FUNC_CLASS_MATH,
		zaryFunc: types.MlrvalUrand32,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_TIME

	{
		name:  "sec2gmt",
		class: FUNC_CLASS_TIME,
		help: `Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is.`,
		unaryFunc:          types.MlrvalSec2GMTUnary,
		binaryFunc:         types.MlrvalSec2GMTBinary,
		hasMultipleArities: true,
	},

	{
		name:     "systime",
		class:    FUNC_CLASS_TIME,
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystime,
	},

	{
		name:     "systimeint",
		class:    FUNC_CLASS_TIME,
		help:     "help string will go here",
		zaryFunc: types.MlrvalSystimeInt,
	},

	{
		name:     "uptime",
		class:    FUNC_CLASS_TIME,
		help:     "help string will go here",
		zaryFunc: types.MlrvalUptime,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_TYPING

	{
		name:      "is_absent",
		class:     FUNC_CLASS_TYPING,
		help:      "False if field is present in input, true otherwise",
		unaryFunc: types.MlrvalIsAbsent,
	},

	{
		name:      "is_array",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is an array.",
		unaryFunc: types.MlrvalIsArray,
	},

	{
		name:      "is_bool",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with boolean value. Synonymous with is_boolean.",
		unaryFunc: types.MlrvalIsBool,
	},

	{
		name:      "is_boolean",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with boolean value. Synonymous with is_bool.",
		unaryFunc: types.MlrvalIsBoolean,
	},

	{
		name:      "is_empty",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present in input with empty string value, false otherwise.",
		unaryFunc: types.MlrvalIsEmpty,
	},

	{
		name:      "is_empty_map",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is a map which is empty.",
		unaryFunc: types.MlrvalIsEmptyMap,
	},

	{
		name:      "is_error",
		class:     FUNC_CLASS_TYPING,
		help:      "True if if argument is an error, such as taking string length of an integer.",
		unaryFunc: types.MlrvalIsError,
	},

	{
		name:      "is_float",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with value inferred to be float",
		unaryFunc: types.MlrvalIsFloat,
	},

	{
		name:      "is_int",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with value inferred to be int",
		unaryFunc: types.MlrvalIsInt,
	},

	{
		name:      "is_map",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is a map.",
		unaryFunc: types.MlrvalIsMap,
	},

	{
		name:      "is_nonempty_map",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is a map which is non-empty.",
		unaryFunc: types.MlrvalIsNonEmptyMap,
	},

	{
		name:      "is_not_empty",
		class:     FUNC_CLASS_TYPING,
		help:      "False if field is present in input with empty value, true otherwise",
		unaryFunc: types.MlrvalIsNotEmpty,
	},

	{
		name:      "is_not_map",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is not a map.",
		unaryFunc: types.MlrvalIsNotMap,
	},

	{
		name:      "is_not_array",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is not an array.",
		unaryFunc: types.MlrvalIsNotArray,
	},

	{
		name:      "is_not_null",
		class:     FUNC_CLASS_TYPING,
		help:      "False if argument is null (empty or absent), true otherwise.",
		unaryFunc: types.MlrvalIsNotNull,
	},

	{
		name:      "is_null",
		class:     FUNC_CLASS_TYPING,
		help:      "True if argument is null (empty or absent), false otherwise.",
		unaryFunc: types.MlrvalIsNull,
	},

	{
		name:      "is_numeric",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with value inferred to be int or float",
		unaryFunc: types.MlrvalIsNumeric,
	},

	{
		name:      "is_present",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present in input, false otherwise.",
		unaryFunc: types.MlrvalIsPresent,
	},

	{
		name:      "is_string",
		class:     FUNC_CLASS_TYPING,
		help:      "True if field is present with string (including empty-string) value",
		unaryFunc: types.MlrvalIsString,
	},

	{
		name:  "asserting_absent",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_absent on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingAbsent,
	},

	{
		name:  "asserting_array",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_array on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingArray,
	},

	{
		name:  "asserting_bool",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_bool on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingBool,
	},

	{
		name:  "asserting_boolean",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_boolean on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingBoolean,
	},

	{
		name:  "asserting_error",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_error on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingError,
	},

	{
		name:  "asserting_empty",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_empty on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingEmpty,
	},

	{
		name:  "asserting_empty_map",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingEmptyMap,
	},

	{
		name:  "asserting_float",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_float on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingFloat,
	},

	{
		name:  "asserting_int",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_int on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingInt,
	},

	{
		name:  "asserting_map",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingMap,
	},

	{
		name:  "asserting_nonempty_map",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNonEmptyMap,
	},

	{
		name:  "asserting_not_empty",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotEmpty,
	},

	{
		name:  "asserting_not_map",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_not_map on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotMap,
	},

	{
		name:  "asserting_not_array",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_not_array on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotArray,
	},

	{
		name:  "asserting_not_null",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_not_null on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNotNull,
	},

	{
		name:  "asserting_null",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_null on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNull,
	},

	{
		name:  "asserting_numeric",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_numeric on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingNumeric,
	},

	{
		name:  "asserting_present",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_present on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingPresent,
	},

	{
		name:  "asserting_string",
		class: FUNC_CLASS_TYPING,
		help: `Aborts with an error if is_string on the argument returns false,
else returns its argument.`,
		contextualUnaryFunc: types.MlrvalAssertingString,
	},

	{
		name:      "typeof",
		class:     FUNC_CLASS_TYPING,
		help:      "Convert argument to type of argument (e.g. \"str\"). For debug.",
		unaryFunc: types.MlrvalTypeof,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_CONVERSION

	{
		name:      "boolean",
		class:     FUNC_CLASS_CONVERSION,
		help:      "Convert int/float/bool/string to boolean.",
		unaryFunc: types.MlrvalToBoolean,
	},

	{
		name:      "float",
		class:     FUNC_CLASS_CONVERSION,
		help:      "Convert int/float/bool/string to float.",
		unaryFunc: types.MlrvalToFloat,
	},

	{
		name:  "fmtnum",
		class: FUNC_CLASS_CONVERSION,
		help: `Convert int/float/bool to string using
printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.`,
		binaryFunc: types.MlrvalFmtNum,
	},

	{
		name:      "hexfmt",
		class:     FUNC_CLASS_CONVERSION,
		help:      `Convert int to hex string, e.g. 255 to "0xff".`,
		unaryFunc: types.MlrvalHexfmt,
	},

	{
		name:      "int",
		class:     FUNC_CLASS_CONVERSION,
		help:      "Convert int/float/bool/string to int.",
		unaryFunc: types.MlrvalToInt,
	},

	{
		name:  "joink",
		class: FUNC_CLASS_CONVERSION,
		help: `Makes string from map/array keys. Examples:
joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
joink([1,2,3], ",") = "1,2,3".`,
		binaryFunc: types.MlrvalJoinK,
	},

	{
		name:  "joinv",
		class: FUNC_CLASS_CONVERSION,
		help: `Makes string from map/array values.
joinv([3,4,5], ",") = "3,4,5"
joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"`,
		binaryFunc: types.MlrvalJoinV,
	},

	{
		name:  "joinkv",
		class: FUNC_CLASS_CONVERSION,
		help: `Makes string from map/array key-value pairs. Examples:
joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"`,
		ternaryFunc: types.MlrvalJoinKV,
	},

	{
		name:  "splita",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string into array with type inference. Example:
splita("3,4,5", ",") = [3,4,5]`,
		binaryFunc: types.MlrvalSplitA,
	},

	{
		name:  "splitax",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string into array without type inference. Example:
splita("3,4,5", ",") = ["3","4","5"]`,
		binaryFunc: types.MlrvalSplitAX,
	},

	{
		name:  "splitkv",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string by separators into map with type inference. Example:
splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}`,
		ternaryFunc: types.MlrvalSplitKV,
	},

	{
		name:  "splitkvx",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string by separators into map without type inference (keys and
values are strings). Example:
splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}`,
		ternaryFunc: types.MlrvalSplitKVX,
	},

	{
		name:  "splitnv",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string by separator into integer-indexed map with type inference. Example:
splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}`,
		binaryFunc: types.MlrvalSplitNV,
	},

	{
		name:  "splitnvx",
		class: FUNC_CLASS_CONVERSION,
		help: `Splits string by separator into integer-indexed map without type
inference (values are strings). Example:
splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}`,
		binaryFunc: types.MlrvalSplitNVX,
	},

	{
		name:      "string",
		class:     FUNC_CLASS_CONVERSION,
		help:      "Convert int/float/bool/string/array/map to string.",
		unaryFunc: types.MlrvalToString,
	},

	// ----------------------------------------------------------------
	// FUNC_CLASS_COLLECTIONS

	{
		name:       "append",
		class:      FUNC_CLASS_COLLECTIONS,
		help:       "Appends second argument to end of first argument, which must be an array.",
		binaryFunc: types.MlrvalAppend,
	},

	{
		name:      "depth",
		class:     FUNC_CLASS_COLLECTIONS,
		help:      "Prints maximum depth of map/array. Scalars have depth 0.",
		unaryFunc: types.MlrvalDepth,
	},

	{
		name:      "get_keys",
		class:     FUNC_CLASS_COLLECTIONS,
		help:      "Returns array of keys of map or array",
		unaryFunc: types.MlrvalGetKeys,
	},

	{
		name:      "get_values",
		class:     FUNC_CLASS_COLLECTIONS,
		help:      "Returns array of keys of map or array -- in the latter case, returns a copy of the array",
		unaryFunc: types.MlrvalGetValues,
	},

	{
		name:  "haskey",
		class: FUNC_CLASS_COLLECTIONS,
		help: `True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.`,
		binaryFunc: types.MlrvalHasKey,
	},

	{
		name:  "leafcount",
		class: FUNC_CLASS_COLLECTIONS,
		help: `Counts total number of terminal values in map/array. For single-level
map/array, same as length.`,
		unaryFunc: types.MlrvalLeafCount,
	},

	{
		name:      "length",
		class:     FUNC_CLASS_COLLECTIONS,
		help:      "Counts number of top-level entries in array/map. Scalars have length 1.",
		unaryFunc: types.MlrvalLength,
	},

	{
		name:  "mapdiff",
		class: FUNC_CLASS_COLLECTIONS,
		help: `With 0 args, returns empty map. With 1 arg, returns copy of arg.
With 2 or more, returns copy of arg 1 with all keys from any of remaining
argument maps removed.`,
		variadicFunc: types.MlrvalMapDiff,
	},

	{
		name:  "mapexcept",
		class: FUNC_CLASS_COLLECTIONS,
		help: `Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.
E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.`,
		variadicFunc:         types.MlrvalMapExcept,
		minimumVariadicArity: 1,
	},

	{
		name:  "mapselect",
		class: FUNC_CLASS_COLLECTIONS,
		help: `Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.
E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.`,
		variadicFunc:         types.MlrvalMapSelect,
		minimumVariadicArity: 1,
	},

	{
		name:  "mapsum",
		class: FUNC_CLASS_COLLECTIONS,
		help: `With 0 args, returns empty map. With >= 1 arg, returns a map with
key-value pairs from all arguments. Rightmost collisions win, e.g.
'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.`,
		variadicFunc: types.MlrvalMapSum,
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
		fmt.Fprintf(o, "%-s  (class=%s #args=%s) %s\n",
			builtinFunctionInfo.name,
			builtinFunctionInfo.class,
			describeNargs(builtinFunctionInfo),
			builtinFunctionInfo.help,
		)
	}
}

func describeNargs(info BuiltinFunctionInfo) string {
	if info.hasMultipleArities {
		pieces := make([]string, 0)
		if info.zaryFunc != nil {
			pieces = append(pieces, "0")
		}
		if info.unaryFunc != nil {
			pieces = append(pieces, "1")
		}
		if info.contextualUnaryFunc != nil {
			pieces = append(pieces, "1")
		}
		if info.binaryFunc != nil {
			pieces = append(pieces, "2")
		}
		if info.ternaryFunc != nil {
			pieces = append(pieces, "3")
			return "3"
		}
		return strings.Join(pieces, ",")

	} else {

		if info.zaryFunc != nil {
			return "0"
		}
		if info.unaryFunc != nil {
			return "1"
		}
		if info.contextualUnaryFunc != nil {
			return "1"
		}
		if info.binaryFunc != nil {
			return "2"
		}
		if info.ternaryFunc != nil {
			return "3"
		}
		if info.variadicFunc != nil {
			return "variadic"
		}
	}
	lib.InternalCodingErrorIf(true)
	return "(error)" // solely to appease the Go compiler; not reached
}

// ================================================================
// This is a singleton so the online-help functions can query it for listings,
// online help, etc.
var BuiltinFunctionManagerInstance *BuiltinFunctionManager = NewBuiltinFunctionManager()
