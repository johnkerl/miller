// ================================================================
// Adding a new builtin function:
// * New entry in makeBuiltinFunctionLookupTable
// * Implement the function in mlrval_functions.go
// ================================================================

package cst

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"mlr/src/colorizer"
	"mlr/src/lib"
	"mlr/src/types"
)

type TFunctionClass string

const (
	FUNC_CLASS_ARITHMETIC  TFunctionClass = "arithmetic"
	FUNC_CLASS_MATH                       = "math"
	FUNC_CLASS_BOOLEAN                    = "boolean"
	FUNC_CLASS_STRING                     = "string"
	FUNC_CLASS_HASHING                    = "hashing"
	FUNC_CLASS_CONVERSION                 = "conversion"
	FUNC_CLASS_TYPING                     = "typing"
	FUNC_CLASS_COLLECTIONS                = "collections"
	FUNC_CLASS_SYSTEM                     = "system"
	FUNC_CLASS_TIME                       = "time"
)

// ================================================================
type BuiltinFunctionInfo struct {
	name                   string
	class                  TFunctionClass
	help                   string
	hasMultipleArities     bool
	minimumVariadicArity   int
	maximumVariadicArity   int // 0 means no max
	zaryFunc               types.ZaryFunc
	unaryFunc              types.UnaryFunc
	unaryFuncWithContext   types.UnaryFuncWithContext   // asserting_{typename}
	regexCaptureBinaryFunc types.RegexCaptureBinaryFunc // =~ and !=~
	binaryFuncWithState    BinaryFuncWithState          // sortaf and sortmf
	binaryFunc             types.BinaryFunc
	ternaryFunc            types.TernaryFunc
	variadicFunc           types.VariadicFunc
}

// ================================================================

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func startsWithLetter(s string) bool {
	if len(s) < 1 {
		return false
	} else {
		return isLetter(s[0])
	}
}

func makeBuiltinFunctionLookupTable() []BuiltinFunctionInfo {
	lookupTable := []BuiltinFunctionInfo{

		// ----------------------------------------------------------------
		// FUNC_CLASS_ARITHMETIC
		{
			name:               "+",
			class:              FUNC_CLASS_ARITHMETIC,
			help:               `Addition as binary operator; unary plus operator.`,
			unaryFunc:          types.MlrvalUnaryPlus,
			binaryFunc:         types.MlrvalBinaryPlus,
			hasMultipleArities: true,
		},

		{
			name:               "-",
			class:              FUNC_CLASS_ARITHMETIC,
			help:               `Subtraction as binary operator; unary negation operator.`,
			unaryFunc:          types.MlrvalUnaryMinus,
			binaryFunc:         types.MlrvalBinaryMinus,
			hasMultipleArities: true,
		},

		{
			name:       "*",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Multiplication, with integer*integer overflow to float.`,
			binaryFunc: types.MlrvalTimes,
		},

		{
			name:       "/",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Division. Integer / integer is floating-point.`,
			binaryFunc: types.MlrvalDivide,
		},

		{
			name:       "//",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Pythonic integer division, rounding toward negative.`,
			binaryFunc: types.MlrvalIntDivide,
		},

		{
			name:       "**",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Exponentiation. Same as pow, but as an infix operator.`,
			binaryFunc: types.MlrvalPow,
		},

		{
			name:       "pow",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Exponentiation. Same as **, but as a function.`,
			binaryFunc: types.MlrvalPow,
		},

		{
			name:       ".+",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Addition, with integer-to-integer overflow.`,
			binaryFunc: types.MlrvalDotPlus,
		},

		{
			name:       ".-",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Subtraction, with integer-to-integer overflow.`,
			binaryFunc: types.MlrvalDotMinus,
		},

		{
			name:       ".*",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Multiplication, with integer-to-integer overflow.`,
			binaryFunc: types.MlrvalDotTimes,
		},

		{
			name:       "./",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Integer division; not pythonic.`,
			binaryFunc: types.MlrvalDotDivide,
		},

		{
			name:       "%",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Remainder; never negative-valued (pythonic).`,
			binaryFunc: types.MlrvalModulus,
		},

		{
			name:  "~",
			class: FUNC_CLASS_ARITHMETIC,
			help: `Bitwise NOT. Beware '$y=~$x' since =~ is the
regex-match operator: try '$y = ~$x'.`,
			unaryFunc: types.MlrvalBitwiseNOT,
		},

		{
			name:       "&",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise AND.`,
			binaryFunc: types.MlrvalBitwiseAND,
		},

		{
			name:       "|",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise OR.`,
			binaryFunc: types.MlrvalBitwiseOR,
		},

		{
			name:       "^",
			help:       `Bitwise XOR.`,
			class:      FUNC_CLASS_ARITHMETIC,
			binaryFunc: types.MlrvalBitwiseXOR,
		},

		{
			name:       "<<",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise left-shift.`,
			binaryFunc: types.MlrvalLeftShift,
		},

		{
			name:       ">>",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise signed right-shift.`,
			binaryFunc: types.MlrvalSignedRightShift,
		},

		{
			name:       ">>>",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise unsigned right-shift.`,
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
			help:      `Logical negation.`,
			unaryFunc: types.MlrvalLogicalNOT,
		},

		{
			name:  "==",
			class: FUNC_CLASS_BOOLEAN,

			help:       `String/numeric equality. Mixing number and string results in string compare.`,
			binaryFunc: types.MlrvalEquals,
		},

		{
			name:       "!=",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric inequality. Mixing number and string results in string compare.`,
			binaryFunc: types.MlrvalNotEquals,
		},

		{
			name:       ">",
			help:       `String/numeric greater-than. Mixing number and string results in string compare.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.MlrvalGreaterThan,
		},

		{
			name:       ">=",
			help:       `String/numeric greater-than-or-equals. Mixing number and string results in string compare.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.MlrvalGreaterThanOrEquals,
		},

		{
			name:       "<=>",
			help:       `Comparator, nominally for sorting. Given a <=> b, returns <0, 0, >0 as a < b, a == b, or a > b, respectively.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.MlrvalCmp,
		},

		{
			name:       "<",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric less-than. Mixing number and string results in string compare.`,
			binaryFunc: types.MlrvalLessThan,
		},

		{
			name:       "<=",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric less-than-or-equals. Mixing number and string results in string compare.`,
			binaryFunc: types.MlrvalLessThanOrEquals,
		},

		{
			name:                   "=~",
			class:                  FUNC_CLASS_BOOLEAN,
			help:                   `String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.`,
			regexCaptureBinaryFunc: types.MlrvalStringMatchesRegexp,
		},

		{
			name:                   "!=~",
			class:                  FUNC_CLASS_BOOLEAN,
			help:                   `String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.`,
			regexCaptureBinaryFunc: types.MlrvalStringDoesNotMatchRegexp,
		},

		{
			name:       "&&",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `Logical AND.`,
			binaryFunc: BinaryShortCircuitPlaceholder,
		},

		{
			name:       "||",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `Logical OR.`,
			binaryFunc: BinaryShortCircuitPlaceholder,
		},

		{
			name:       "^^",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `Logical XOR.`,
			binaryFunc: types.MlrvalLogicalXOR,
		},

		{
			name:       "??",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.`,
			binaryFunc: BinaryShortCircuitPlaceholder,
		},

		{
			name:       "???",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.`,
			binaryFunc: BinaryShortCircuitPlaceholder,
		},

		{
			name:        "?:",
			class:       FUNC_CLASS_BOOLEAN,
			help:        `Standard ternary operator.`,
			ternaryFunc: TernaryShortCircuitPlaceholder,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_STRING

		{
			name:       ".",
			class:      FUNC_CLASS_STRING,
			help:       `String concatenation.`,
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
			help:        `Example: '$name=gsub($name, "old", "new")' (replace all).`,
			ternaryFunc: types.MlrvalGsub,
		},

		{
			name:      "lstrip",
			class:     FUNC_CLASS_STRING,
			help:      "Strip leading whitespace from string.",
			unaryFunc: types.MlrvalLStrip,
		},

		{
			name:       "regextract",
			class:      FUNC_CLASS_STRING,
			help:       `Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'`,
			binaryFunc: types.MlrvalRegextract,
		},

		{
			name:        "regextract_or_else",
			class:       FUNC_CLASS_STRING,
			help:        `Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'`,
			ternaryFunc: types.MlrvalRegextractOrElse,
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
			help:        `Like sub but does no regexing. No characters are special.`,
			ternaryFunc: types.MlrvalSsub,
		},

		{
			name:        "sub",
			class:       FUNC_CLASS_STRING,
			help:        `Example: '$name=sub($name, "old", "new")' (replace once).`,
			ternaryFunc: types.MlrvalSub,
		},

		{
			name:  "substr0",
			class: FUNC_CLASS_STRING,
			help: `substr0(s,m,n) gives substring of s from 0-up position m to n
inclusive. Negative indices -len .. -1 alias to 0 .. len-1. See also substr and substr1.`,
			ternaryFunc: types.MlrvalSubstr0Up,
		},
		{
			name:  "substr1",
			class: FUNC_CLASS_STRING,
			help: `substr1(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len. See also substr and substr0.`,
			ternaryFunc: types.MlrvalSubstr1Up,
		},
		{
			name:  "substr",
			class: FUNC_CLASS_STRING,
			help: `substr is an alias for substr0. See also substr1. Miller is generally 1-up
with all array and string indices, but, this is a backward-compatibility issue with Miller 5
and below. Arrays are new in Miller 6; the substr function is older.`,
			ternaryFunc: types.MlrvalSubstr0Up,
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
			help:       `Truncates string first argument to max length of int second argument.`,
			binaryFunc: types.MlrvalTruncate,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_HASHING

		{
			name:      "md5",
			class:     FUNC_CLASS_HASHING,
			help:      `MD5 hash.`,
			unaryFunc: types.MlrvalMD5,
		},
		{
			name:      "sha1",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA1 hash.`,
			unaryFunc: types.MlrvalSHA1,
		},
		{
			name:      "sha256",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA256 hash.`,
			unaryFunc: types.MlrvalSHA256,
		},
		{
			name:      "sha512",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA512 hash.`,
			unaryFunc: types.MlrvalSHA512,
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
			name:       "atan2",
			class:      FUNC_CLASS_MATH,
			help:       "Two-argument arctangent.",
			binaryFunc: types.MlrvalAtan2,
		},

		{
			name:      "atanh",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse hyperbolic tangent.",
			unaryFunc: types.MlrvalAtanh,
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
			help:         `Max of n numbers; null loses.`,
			variadicFunc: types.MlrvalVariadicMax,
		},

		{
			name:         "min",
			class:        FUNC_CLASS_MATH,
			help:         `Min of n numbers; null loses.`,
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
			name:  "urand",
			class: FUNC_CLASS_MATH,
			help: `Floating-point numbers uniformly distributed on the unit interval.
Int-valued example: '$n=floor(20+urand()*11)'.`,
			zaryFunc: types.MlrvalUrand,
		},

		{
			name:       "urandint",
			class:      FUNC_CLASS_MATH,
			help:       `Integer uniformly distributed between inclusive integer endpoints.`,
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
			help:     `Integer uniformly distributed 0 and 2**32-1 inclusive.`,
			zaryFunc: types.MlrvalUrand32,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_TIME

		{
			name:      "gmt2sec",
			class:     FUNC_CLASS_TIME,
			help:      `Parses GMT timestamp as integer seconds since the epoch.`,
			unaryFunc: types.MlrvalGMT2Sec,
		},

		{
			name:  "sec2gmt",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is. With second integer argument n, includes n decimal places
for the seconds part`,
			unaryFunc:          types.MlrvalSec2GMTUnary,
			binaryFunc:         types.MlrvalSec2GMTBinary,
			hasMultipleArities: true,
		},

		{
			name:  "sec2gmtdate",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch (integer part)
as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
Leaves non-numbers as-is.`,
			unaryFunc: types.MlrvalSec2GMTDate,
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

		{
			name:  "strftime",
			class: FUNC_CLASS_TIME,
			help: ` Formats seconds since the epoch as timestamp, e.g.
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
	Format strings are as in the C library (please see "man strftime" on your system),
	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
	See also strftime_local.`,
			binaryFunc: types.MlrvalStrftime,
		},

		{
			name:  "strptime",
			class: FUNC_CLASS_TIME,
			help: `strptime: Parses timestamp as floating-point seconds since the epoch,
	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
	See also strptime_local.`,
			binaryFunc: types.MlrvalStrptime,
		},

		// TODO:

		// strftime_local (class=time #args=2): Like strftime but consults the $TZ environment variable to get local time zone.

		// strptime_local (class=time #args=2): Like strptime, but consults $TZ environment variable to find and use local timezone.

		{
			name:      "dhms2fsec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000`,
			unaryFunc: types.MlrvalDHMS2FSec,
		},

		{
			name:      "dhms2sec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000`,
			unaryFunc: types.MlrvalDHMS2Sec,
		},

		{
			name:      "fsec2dhms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"`,
			unaryFunc: types.MlrvalFSec2DHMS,
		},

		{
			name:      "fsec2hms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"`,
			unaryFunc: types.MlrvalFSec2HMS,
		},

		{
			name:      "hms2fsec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000`,
			unaryFunc: types.MlrvalHMS2FSec,
		},

		{
			name:      "hms2sec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers integer seconds as in hms2sec("01:23:20") = 5000`,
			unaryFunc: types.MlrvalHMS2Sec,
		},

		{
			name:      "sec2dhms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"`,
			unaryFunc: types.MlrvalSec2DHMS,
		},

		{
			name:      "sec2hms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats integer seconds as in sec2hms(5000) = "01:23:20"`,
			unaryFunc: types.MlrvalSec2HMS,
		},

		// localtime2sec (class=time #args=1): Parses local timestamp as integer seconds since
		// the epoch. Consults $TZ environment variable.

		// sec2localtime (class=time #args=1): Formats seconds since epoch (integer part)
		// as local timestamp, e.g. sec2localtime(1440768801.7) = "2015-08-28T13:33:21Z".
		// Consults $TZ environment variable. Leaves non-numbers as-is.
		//
		// sec2localtime (class=time #args=2): Formats seconds since epoch as local timestamp with n
		// decimal places for seconds, e.g. sec2localtime(1440768801.7,1) = "2015-08-28T13:33:21.7Z".
		// Consults $TZ environment variable. Leaves non-numbers as-is.
		//
		// sec2localdate (class=time #args=1): Formats seconds since epoch (integer part)
		// as local timestamp with year-month-date, e.g. sec2localdate(1440768801.7) = "2015-08-28".
		// Consults $TZ environment variable. Leaves non-numbers as-is.

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
			unaryFuncWithContext: types.MlrvalAssertingAbsent,
		},

		{
			name:  "asserting_array",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_array on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingArray,
		},

		{
			name:  "asserting_bool",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_bool on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingBool,
		},

		{
			name:  "asserting_boolean",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_boolean on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingBoolean,
		},

		{
			name:  "asserting_error",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_error on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingError,
		},

		{
			name:  "asserting_empty",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_empty on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingEmpty,
		},

		{
			name:  "asserting_empty_map",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingEmptyMap,
		},

		{
			name:  "asserting_float",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_float on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingFloat,
		},

		{
			name:  "asserting_int",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_int on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingInt,
		},

		{
			name:  "asserting_map",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_map on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingMap,
		},

		{
			name:  "asserting_nonempty_map",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNonEmptyMap,
		},

		{
			name:  "asserting_not_empty",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNotEmpty,
		},

		{
			name:  "asserting_not_map",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_not_map on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNotMap,
		},

		{
			name:  "asserting_not_array",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_not_array on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNotArray,
		},

		{
			name:  "asserting_not_null",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_not_null on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNotNull,
		},

		{
			name:  "asserting_null",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_null on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNull,
		},

		{
			name:  "asserting_numeric",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_numeric on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingNumeric,
		},

		{
			name:  "asserting_present",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_present on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingPresent,
		},

		{
			name:  "asserting_string",
			class: FUNC_CLASS_TYPING,
			help: `Aborts with an error if is_string on the argument returns false,
else returns its argument.`,
			unaryFuncWithContext: types.MlrvalAssertingString,
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
			name:  "arrayify",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Walks through a nested map/array, converting any map with consecutive keys
"1", "2", ... into an array. Useful to wrap the output of unflatten.`,
			unaryFunc: types.MlrvalArrayify,
		},

		{
			name:      "depth",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      "Prints maximum depth of map/array. Scalars have depth 0.",
			unaryFunc: types.MlrvalDepth,
		},

		{
			name:  "flatten",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Flattens multi-level maps to single-level ones. Examples:
flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
Useful for nested JSON-like structures for non-JSON file formats like CSV.`,
			binaryFunc:         types.MlrvalFlattenBinary,
			ternaryFunc:        types.MlrvalFlatten,
			hasMultipleArities: true,
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
			name:      "json_parse",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      `Converts value from JSON-formatted string.`,
			unaryFunc: types.MlrvalJSONParse,
		},
		{
			name:  "json_stringify",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Converts value to JSON-formatted string. Default output is single-line.
With optional second boolean argument set to true, produces multiline output.`,
			unaryFunc:          types.MlrvalJSONStringifyUnary,
			binaryFunc:         types.MlrvalJSONStringifyBinary,
			hasMultipleArities: true,
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

		{
			name:  "unflatten",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Reverses flatten. Example:
unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.
Useful for nested JSON-like structures for non-JSON file formats like CSV.
See also arrayify.`,
			binaryFunc: types.MlrvalUnflatten,
		},

		{
			name:                 "sorta",
			class:                FUNC_CLASS_COLLECTIONS,
			help:                 "Returns a copy of an array, sorted ascending. Coming soon: other sort options.",
			variadicFunc:         types.MlrvalSortA,
			minimumVariadicArity: 1,
			maximumVariadicArity: 2,
		},

		{
			name:                 "sortmk",
			class:                FUNC_CLASS_COLLECTIONS,
			help:                 "Returns a copy of a map, sorted ascending by map key. Coming soon: other sort options.",
			variadicFunc:         types.MlrvalSortMK,
			minimumVariadicArity: 1,
			maximumVariadicArity: 2,
		},

		{
			name:  "sortaf",
			class: FUNC_CLASS_COLLECTIONS,
			help:  "Sorts an array (1st argument) using a comparator function you specify by name (2nd argument).  Example: 'sortaf([5,2,3,1,4], \"f\")'.  Forward sort: 'func f(a,b) {return a <=> b}'.  Reverse sort: 'func f(a,b) {return b <=> a}'.  And so on -- you can implement logic you choose. Your function should return a number <0, ==0, >0 as a<b, a==b, or a>b respectively.",
			// Not in the types package since it uses UDFs which are in the dsl/cst
			// package, which would make a circular reference.
			binaryFuncWithState: SortAF,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_SYSTEM

		{
			name:     "hostname",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the hostname as a string.`,
			zaryFunc: types.MlrvalHostname,
		},

		{
			name:     "os",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the operating-system name as a string.`,
			zaryFunc: types.MlrvalOS,
		},

		{
			name:      "system",
			class:     FUNC_CLASS_SYSTEM,
			help:      `Run command string, yielding its stdout minus final carriage return.`,
			unaryFunc: types.MlrvalSystem,
		},

		{
			name:     "version",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the Miller version as a string.`,
			zaryFunc: types.MlrvalVersion,
		},
	}

	// Sort the function table.  Useful for online help and autogenned docs / manpage.

	// Go sort API: for ascending sort, return true if element i < element j.
	// Put symbols like '||' after text like 'gsub', not before.
	sort.Slice(lookupTable, func(i, j int) bool {
		namei := lookupTable[i].name
		namej := lookupTable[j].name

		si := startsWithLetter(namei)
		sj := startsWithLetter(namej)

		if si && !sj {
			return true
		} else if !si && sj {
			return false
		} else {
			if namei < namej {
				return true
			} else {
				return false
			}
		}
	})

	return lookupTable
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
	lookupTable := makeBuiltinFunctionLookupTable()
	hashTable := hashifyLookupTable(&lookupTable)
	return &BuiltinFunctionManager{
		lookupTable: &lookupTable,
		hashTable:   hashTable,
	}
}

func (manager *BuiltinFunctionManager) LookUp(functionName string) *BuiltinFunctionInfo {
	return manager.hashTable[functionName]
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
func (manager *BuiltinFunctionManager) ListBuiltinFunctionClasses() {
	classesList := manager.getBuiltinFunctionClasses()
	for _, class := range classesList {
		fmt.Println(class)
	}
}

func (manager *BuiltinFunctionManager) getBuiltinFunctionClasses() []string {
	classesSeen := make(map[string]bool)
	classesList := make([]string, 0)
	for _, builtinFunctionInfo := range *manager.lookupTable {
		class := string(builtinFunctionInfo.class)
		if classesSeen[class] == false {
			classesList = append(classesList, class)
			classesSeen[class] = true
		}
	}
	sort.Strings(classesList)
	return classesList
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionsInClass(class string) {
	for _, builtinFunctionInfo := range *manager.lookupTable {
		if string(builtinFunctionInfo.class) == class {
			fmt.Println(builtinFunctionInfo.name)
		}
	}
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionNamesVertically() {
	for _, builtinFunctionInfo := range *manager.lookupTable {
		fmt.Println(builtinFunctionInfo.name)
	}
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionNamesAsParagraph() {
	functionNames := make([]string, len(*manager.lookupTable))
	for i, builtinFunctionInfo := range *manager.lookupTable {
		functionNames[i] = builtinFunctionInfo.name
	}
	lib.PrintWordsAsParagraph(functionNames)
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionsAsTable() {
	fmt.Printf("%-30s %-12s %s\n", "Name", "Class", "Args")
	for _, builtinFunctionInfo := range *manager.lookupTable {
		fmt.Printf("%-30s %-12s %s\n",
			builtinFunctionInfo.name,
			builtinFunctionInfo.class,
			describeNargs(&builtinFunctionInfo),
		)
	}
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionUsages() {
	for i, builtinFunctionInfo := range *manager.lookupTable {
		if i > 0 {
			fmt.Println()
		}
		lib.InternalCodingErrorIf(builtinFunctionInfo.help == "")
		fmt.Print(colorizer.MaybeColorizeHelp(builtinFunctionInfo.name, true))
		fmt.Printf("  (class=%s #args=%s) %s\n",
			builtinFunctionInfo.class,
			describeNargs(&builtinFunctionInfo),
			builtinFunctionInfo.help,
		)
	}
}

// ----------------------------------------------------------------
func (manager *BuiltinFunctionManager) ListBuiltinFunctionUsagesByClass() {
	classesList := manager.getBuiltinFunctionClasses()

	for _, class := range classesList {
		fmt.Println()
		fmt.Println(colorizer.MaybeColorizeHelp(strings.ToUpper(class), true))
		fmt.Println()
		for _, builtinFunctionInfo := range *manager.lookupTable {
			if string(builtinFunctionInfo.class) != class {
				continue
			}
			lib.InternalCodingErrorIf(builtinFunctionInfo.help == "")
			fmt.Print(colorizer.MaybeColorizeHelp(builtinFunctionInfo.name, true))
			fmt.Printf("  (class=%s #args=%s) %s\n",
				builtinFunctionInfo.class,
				describeNargs(&builtinFunctionInfo),
				builtinFunctionInfo.help,
			)
			fmt.Println()
		}
	}
}

func (manager *BuiltinFunctionManager) ListBuiltinFunctionUsage(functionName string) {
	if !manager.TryListBuiltinFunctionUsage(functionName) {
		fmt.Fprintf(os.Stderr, "Function \"%s\" not found.\n", functionName)
	}
}

func (manager *BuiltinFunctionManager) TryListBuiltinFunctionUsage(functionName string) bool {
	builtinFunctionInfo := manager.LookUp(functionName)
	if builtinFunctionInfo == nil {
		manager.listBuiltinFunctionUsageApproximate(functionName)
		return false
	}
	manager.listBuiltinFunctionUsageExact(builtinFunctionInfo)
	return true
}

func (manager *BuiltinFunctionManager) listBuiltinFunctionUsageExact(
	builtinFunctionInfo *BuiltinFunctionInfo,
) {
	lib.InternalCodingErrorIf(builtinFunctionInfo.help == "")
	fmt.Printf("%-s  (class=%s #args=%s) %s\n",
		colorizer.MaybeColorizeHelp(builtinFunctionInfo.name, true),
		builtinFunctionInfo.class,
		describeNargs(builtinFunctionInfo),
		builtinFunctionInfo.help,
	)
}

func (manager *BuiltinFunctionManager) listBuiltinFunctionUsageApproximate(
	text string,
) {
	fmt.Printf("No exact match for \"%s\". Inexact matches:\n", text)
	found := false
	for _, builtinFunctionInfo := range *manager.lookupTable {
		if strings.Contains(builtinFunctionInfo.name, text) {
			fmt.Printf("  %s\n", builtinFunctionInfo.name)
			found = true
		}
	}
	if !found {
		fmt.Println("None found.")
	}
}

func describeNargs(info *BuiltinFunctionInfo) string {
	if info.hasMultipleArities {
		pieces := make([]string, 0)
		if info.zaryFunc != nil {
			pieces = append(pieces, "0")
		}
		if info.unaryFunc != nil {
			pieces = append(pieces, "1")
		}
		if info.unaryFuncWithContext != nil {
			pieces = append(pieces, "1")
		}
		if info.binaryFunc != nil {
			pieces = append(pieces, "2")
		}
		if info.regexCaptureBinaryFunc != nil {
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
		if info.unaryFuncWithContext != nil {
			return "1"
		}
		if info.binaryFunc != nil {
			return "2"
		}
		if info.binaryFuncWithState != nil {
			return "2"
		}
		if info.regexCaptureBinaryFunc != nil {
			return "2"
		}
		if info.ternaryFunc != nil {
			return "3"
		}
		if info.variadicFunc != nil {
			if info.maximumVariadicArity != 0 {
				return fmt.Sprintf("%d-%d", info.minimumVariadicArity, info.maximumVariadicArity)
			} else {
				return "variadic"
			}
		}
	}
	lib.InternalCodingErrorIf(true)
	return "(error)" // solely to appease the Go compiler; not reached
}

// ================================================================
// This is a singleton so the online-help functions can query it for listings,
// online help, etc.
var BuiltinFunctionManagerInstance *BuiltinFunctionManager = NewBuiltinFunctionManager()
