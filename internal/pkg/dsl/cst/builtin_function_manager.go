// ================================================================
// Adding a new builtin function:
// * New entry in makeBuiltinFunctionLookupTable
// * Implement the function in mlrval_functions.go
//
// Note: Miller-DSL functions, e.g. sec2gmt, are implemented by Go functions
// with names like BIF_sec2gmt not Sec2GMT. This is an intentional departure
// from Go naming conventions: it makes it easier to mentally pair up
// Miller-DSL functions with their Go implementations. Please preserve this
// naming convention.
// ================================================================

package cst

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
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
	FUNC_CLASS_HOFS                       = "higher-order-functions"
	FUNC_CLASS_SYSTEM                     = "system"
	FUNC_CLASS_TIME                       = "time"
)

// ================================================================
type BuiltinFunctionInfo struct {
	name  string
	class TFunctionClass
	// For source-code storage, these have newlines in them. For any presentation to the user, they must be
	// formatted using the JoinHelp() method which joins newlines. This is crucial for rendering of
	// help-strings for manual page, webdocs, etc wherein we must let the user's resizing of the terminal
	// window or browser determine -- at their choosing -- where lines wrap.
	help                   string
	examples               []string
	hasMultipleArities     bool
	minimumVariadicArity   int
	maximumVariadicArity   int // 0 means no max
	zaryFunc               types.ZaryFunc
	unaryFunc              types.UnaryFunc
	binaryFunc             types.BinaryFunc
	ternaryFunc            types.TernaryFunc
	variadicFunc           types.VariadicFunc
	unaryFuncWithContext   types.UnaryFuncWithContext   // asserting_{typename}
	regexCaptureBinaryFunc types.RegexCaptureBinaryFunc // =~ and !=~
	binaryFuncWithState    BinaryFuncWithState          // select, apply, reduce
	ternaryFuncWithState   TernaryFuncWithState         // fold
	variadicFuncWithState  VariadicFuncWithState        // sort
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
			unaryFunc:          types.BIF_plus_unary,
			binaryFunc:         types.BIF_plus_binary,
			hasMultipleArities: true,
		},

		{
			name:               "-",
			class:              FUNC_CLASS_ARITHMETIC,
			help:               `Subtraction as binary operator; unary negation operator.`,
			unaryFunc:          types.BIF_minus_unary,
			binaryFunc:         types.BIF_minus_binary,
			hasMultipleArities: true,
		},

		{
			name:       "*",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Multiplication, with integer*integer overflow to float.`,
			binaryFunc: types.BIF_times,
		},

		{
			name:       "/",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Division. Integer / integer is floating-point.`,
			binaryFunc: types.BIF_divide,
		},

		{
			name:       "//",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Pythonic integer division, rounding toward negative.`,
			binaryFunc: types.BIF_int_divide,
		},

		{
			name:       "**",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Exponentiation. Same as pow, but as an infix operator.`,
			binaryFunc: types.BIF_pow,
		},

		{
			name:       "pow",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Exponentiation. Same as **, but as a function.`,
			binaryFunc: types.BIF_pow,
		},

		{
			name:       ".+",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Addition, with integer-to-integer overflow.`,
			binaryFunc: types.BIF_dot_plus,
		},

		{
			name:       ".-",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Subtraction, with integer-to-integer overflow.`,
			binaryFunc: types.BIF_dot_minus,
		},

		{
			name:       ".*",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Multiplication, with integer-to-integer overflow.`,
			binaryFunc: types.BIF_dot_times,
		},

		{
			name:       "./",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Integer division; not pythonic.`,
			binaryFunc: types.BIF_dot_divide,
		},

		{
			name:       "%",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Remainder; never negative-valued (pythonic).`,
			binaryFunc: types.BIF_modulus,
		},

		{
			name:      "~",
			class:     FUNC_CLASS_ARITHMETIC,
			help:      `Bitwise NOT. Beware '$y=~$x' since =~ is the regex-match operator: try '$y = ~$x'.`,
			unaryFunc: types.BIF_bitwise_not,
		},

		{
			name:       "&",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise AND.`,
			binaryFunc: types.BIF_bitwise_and,
		},

		{
			name:       "|",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise OR.`,
			binaryFunc: types.BIF_bitwise_or,
		},

		{
			name:       "^",
			help:       `Bitwise XOR.`,
			class:      FUNC_CLASS_ARITHMETIC,
			binaryFunc: types.BIF_bitwise_xor,
		},

		{
			name:       "<<",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise left-shift.`,
			binaryFunc: types.BIF_left_shift,
		},

		{
			name:       ">>",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise signed right-shift.`,
			binaryFunc: types.BIF_signed_right_shift,
		},

		{
			name:       ">>>",
			class:      FUNC_CLASS_ARITHMETIC,
			help:       `Bitwise unsigned right-shift.`,
			binaryFunc: types.BIF_unsigned_right_shift,
		},

		{
			name:      "bitcount",
			class:     FUNC_CLASS_ARITHMETIC,
			help:      "Count of 1-bits.",
			unaryFunc: types.BIF_bitcount,
		},

		{
			name:        "madd",
			class:       FUNC_CLASS_ARITHMETIC,
			help:        `a + b mod m (integers)`,
			ternaryFunc: types.BIF_mod_add,
		},

		{
			name:        "msub",
			class:       FUNC_CLASS_ARITHMETIC,
			help:        `a - b mod m (integers)`,
			ternaryFunc: types.BIF_mod_sub,
		},

		{
			name:        "mmul",
			class:       FUNC_CLASS_ARITHMETIC,
			help:        `a * b mod m (integers)`,
			ternaryFunc: types.BIF_mod_mul,
		},

		{
			name:        "mexp",
			class:       FUNC_CLASS_ARITHMETIC,
			help:        `a ** b mod m (integers)`,
			ternaryFunc: types.BIF_mod_exp,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_BOOLEAN

		{
			name:      "!",
			class:     FUNC_CLASS_BOOLEAN,
			help:      `Logical negation.`,
			unaryFunc: types.BIF_logicalnot,
		},

		{
			name:  "==",
			class: FUNC_CLASS_BOOLEAN,

			help:       `String/numeric equality. Mixing number and string results in string compare.`,
			binaryFunc: types.BIF_equals,
		},

		{
			name:       "!=",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric inequality. Mixing number and string results in string compare.`,
			binaryFunc: types.BIF_not_equals,
		},

		{
			name:       ">",
			help:       `String/numeric greater-than. Mixing number and string results in string compare.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.BIF_greater_than,
		},

		{
			name:       ">=",
			help:       `String/numeric greater-than-or-equals. Mixing number and string results in string compare.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.BIF_greater_than_or_equals,
		},

		{
			name:       "<=>",
			help:       `Comparator, nominally for sorting. Given a <=> b, returns <0, 0, >0 as a < b, a == b, or a > b, respectively.`,
			class:      FUNC_CLASS_BOOLEAN,
			binaryFunc: types.BIF_cmp,
		},

		{
			name:       "<",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric less-than. Mixing number and string results in string compare.`,
			binaryFunc: types.BIF_less_than,
		},

		{
			name:       "<=",
			class:      FUNC_CLASS_BOOLEAN,
			help:       `String/numeric less-than-or-equals. Mixing number and string results in string compare.`,
			binaryFunc: types.BIF_less_than_or_equals,
		},

		{
			name:                   "=~",
			class:                  FUNC_CLASS_BOOLEAN,
			help:                   `String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.`,
			regexCaptureBinaryFunc: types.BIF_string_matches_regexp,
		},

		{
			name:                   "!=~",
			class:                  FUNC_CLASS_BOOLEAN,
			help:                   `String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.`,
			regexCaptureBinaryFunc: types.BIF_string_does_not_match_regexp,
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
			binaryFunc: types.BIF_logicalxor,
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
			binaryFunc: types.BIF_dot,
		},

		{
			name:      "capitalize",
			class:     FUNC_CLASS_STRING,
			help:      "Convert string's first character to uppercase.",
			unaryFunc: types.BIF_capitalize,
		},

		{
			name:      "clean_whitespace",
			class:     FUNC_CLASS_STRING,
			help:      "Same as collapse_whitespace and strip.",
			unaryFunc: types.BIF_clean_whitespace,
		},

		{
			name:      "collapse_whitespace",
			class:     FUNC_CLASS_STRING,
			help:      "Strip repeated whitespace from string.",
			unaryFunc: types.BIF_collapse_whitespace,
		},

		{
			name:        "gsub",
			class:       FUNC_CLASS_STRING,
			help:        `'$name=gsub($name, "old", "new")' (replace all).`,
			ternaryFunc: types.BIF_gsub,
		},

		{
			name:      "lstrip",
			class:     FUNC_CLASS_STRING,
			help:      "Strip leading whitespace from string.",
			unaryFunc: types.BIF_lstrip,
		},

		{
			name:       "regextract",
			class:      FUNC_CLASS_STRING,
			help:       `'$name=regextract($name, "[A-Z]{3}[0-9]{2}")'`,
			binaryFunc: types.BIF_regextract,
		},

		{
			name:        "regextract_or_else",
			class:       FUNC_CLASS_STRING,
			help:        `'$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'`,
			ternaryFunc: types.BIF_regextract_or_else,
		},

		{
			name:      "rstrip",
			class:     FUNC_CLASS_STRING,
			help:      "Strip trailing whitespace from string.",
			unaryFunc: types.BIF_rstrip,
		},

		{
			name:      "strip",
			class:     FUNC_CLASS_STRING,
			help:      "Strip leading and trailing whitespace from string.",
			unaryFunc: types.BIF_strip,
		},

		{
			name:      "strlen",
			class:     FUNC_CLASS_STRING,
			help:      "String length.",
			unaryFunc: types.BIF_strlen,
		},

		{
			name:        "ssub",
			class:       FUNC_CLASS_STRING,
			help:        `Like sub but does no regexing. No characters are special.`,
			ternaryFunc: types.BIF_ssub,
		},

		{
			name:        "sub",
			class:       FUNC_CLASS_STRING,
			help:        `'$name=sub($name, "old", "new")' (replace once).`,
			ternaryFunc: types.BIF_sub,
		},

		{
			name:  "substr0",
			class: FUNC_CLASS_STRING,
			help: `substr0(s,m,n) gives substring of s from 0-up position m to n inclusive.
Negative indices -len .. -1 alias to 0 .. len-1. See also substr and substr1.`,
			ternaryFunc: types.BIF_substr_0_up,
		},
		{
			name:  "substr1",
			class: FUNC_CLASS_STRING,
			help: `substr1(s,m,n) gives substring of s from 1-up position m to n inclusive.
Negative indices -len .. -1 alias to 1 .. len. See also substr and substr0.`,
			ternaryFunc: types.BIF_substr_1_up,
		},
		{
			name:  "substr",
			class: FUNC_CLASS_STRING,
			help: `substr is an alias for substr0. See also substr1. Miller is generally 1-up with all
array and string indices, but, this is a backward-compatibility issue with Miller 5 and below.
Arrays are new in Miller 6; the substr function is older.`,
			ternaryFunc: types.BIF_substr_0_up,
		},

		{
			name:      "tolower",
			class:     FUNC_CLASS_STRING,
			help:      "Convert string to lowercase.",
			unaryFunc: types.BIF_tolower,
		},

		{
			name:      "toupper",
			class:     FUNC_CLASS_STRING,
			help:      "Convert string to uppercase.",
			unaryFunc: types.BIF_toupper,
		},

		{
			name:       "truncate",
			class:      FUNC_CLASS_STRING,
			help:       `Truncates string first argument to max length of int second argument.`,
			binaryFunc: types.BIF_truncate,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_HASHING

		{
			name:      "md5",
			class:     FUNC_CLASS_HASHING,
			help:      `MD5 hash.`,
			unaryFunc: types.BIF_md5,
		},
		{
			name:      "sha1",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA1 hash.`,
			unaryFunc: types.BIF_sha1,
		},
		{
			name:      "sha256",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA256 hash.`,
			unaryFunc: types.BIF_sha256,
		},
		{
			name:      "sha512",
			class:     FUNC_CLASS_HASHING,
			help:      `SHA512 hash.`,
			unaryFunc: types.BIF_sha512,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_MATH

		{
			name:      "abs",
			class:     FUNC_CLASS_MATH,
			help:      "Absolute value.",
			unaryFunc: types.BIF_abs,
		},

		{
			name:      "acos",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse trigonometric cosine.",
			unaryFunc: types.BIF_acos,
		},

		{
			name:      "acosh",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse hyperbolic cosine.",
			unaryFunc: types.BIF_acosh,
		},

		{
			name:      "asin",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse trigonometric sine.",
			unaryFunc: types.BIF_asin,
		},

		{
			name:      "asinh",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse hyperbolic sine.",
			unaryFunc: types.BIF_asinh,
		},

		{
			name:      "atan",
			class:     FUNC_CLASS_MATH,
			help:      "One-argument arctangent.",
			unaryFunc: types.BIF_atan,
		},

		{
			name:       "atan2",
			class:      FUNC_CLASS_MATH,
			help:       "Two-argument arctangent.",
			binaryFunc: types.BIF_atan2,
		},

		{
			name:      "atanh",
			class:     FUNC_CLASS_MATH,
			help:      "Inverse hyperbolic tangent.",
			unaryFunc: types.BIF_atanh,
		},

		{
			name:      "cbrt",
			class:     FUNC_CLASS_MATH,
			help:      "Cube root.",
			unaryFunc: types.BIF_cbrt,
		},

		{
			name:      "ceil",
			class:     FUNC_CLASS_MATH,
			help:      "Ceiling: nearest integer at or above.",
			unaryFunc: types.BIF_ceil,
		},

		{
			name:      "cos",
			class:     FUNC_CLASS_MATH,
			help:      "Trigonometric cosine.",
			unaryFunc: types.BIF_cos,
		},

		{
			name:      "cosh",
			class:     FUNC_CLASS_MATH,
			help:      "Hyperbolic cosine.",
			unaryFunc: types.BIF_cosh,
		},

		{
			name:      "erf",
			class:     FUNC_CLASS_MATH,
			help:      "Error function.",
			unaryFunc: types.BIF_erf,
		},

		{
			name:      "erfc",
			class:     FUNC_CLASS_MATH,
			help:      "Complementary error function.",
			unaryFunc: types.BIF_erfc,
		},

		{
			name:      "exp",
			class:     FUNC_CLASS_MATH,
			help:      "Exponential function e**x.",
			unaryFunc: types.BIF_exp,
		},

		{
			name:      "expm1",
			class:     FUNC_CLASS_MATH,
			help:      "e**x - 1.",
			unaryFunc: types.BIF_expm1,
		},

		{
			name:      "floor",
			class:     FUNC_CLASS_MATH,
			help:      "Floor: nearest integer at or below.",
			unaryFunc: types.BIF_floor,
		},

		{
			name:  "invqnorm",
			class: FUNC_CLASS_MATH,
			help: `Inverse of normal cumulative distribution function.  Note that invqorm(urand())
is normally distributed.`,
			unaryFunc: types.BIF_invqnorm,
		},

		{
			name:      "log",
			class:     FUNC_CLASS_MATH,
			help:      "Natural (base-e) logarithm.",
			unaryFunc: types.BIF_log,
		},

		{
			name:      "log10",
			class:     FUNC_CLASS_MATH,
			help:      "Base-10 logarithm.",
			unaryFunc: types.BIF_log10,
		},

		{
			name:      "log1p",
			class:     FUNC_CLASS_MATH,
			help:      "log(1-x).",
			unaryFunc: types.BIF_log1p,
		},

		{
			name:        "logifit",
			class:       FUNC_CLASS_MATH,
			help:        `Given m and b from logistic regression, compute fit: $yhat=logifit($x,$m,$b).`,
			ternaryFunc: types.BIF_logifit,
		},

		{
			name:         "max",
			class:        FUNC_CLASS_MATH,
			help:         `Max of n numbers; null loses.`,
			variadicFunc: types.BIF_max_variadic,
		},

		{
			name:         "min",
			class:        FUNC_CLASS_MATH,
			help:         `Min of n numbers; null loses.`,
			variadicFunc: types.BIF_min_variadic,
		},

		{
			name:      "qnorm",
			class:     FUNC_CLASS_MATH,
			help:      `Normal cumulative distribution function.`,
			unaryFunc: types.BIF_qnorm,
		},

		{
			name:      "round",
			class:     FUNC_CLASS_MATH,
			help:      "Round to nearest integer.",
			unaryFunc: types.BIF_round,
		},

		{
			name:      "sgn",
			class:     FUNC_CLASS_MATH,
			help:      `+1, 0, -1 for positive, zero, negative input respectively.`,
			unaryFunc: types.BIF_sgn,
		},

		{
			name:      "sin",
			class:     FUNC_CLASS_MATH,
			help:      "Trigonometric sine.",
			unaryFunc: types.BIF_sin,
		},

		{
			name:      "sinh",
			class:     FUNC_CLASS_MATH,
			help:      "Hyperbolic sine.",
			unaryFunc: types.BIF_sinh,
		},

		{
			name:      "sqrt",
			class:     FUNC_CLASS_MATH,
			help:      "Square root.",
			unaryFunc: types.BIF_sqrt,
		},

		{
			name:      "tan",
			class:     FUNC_CLASS_MATH,
			help:      "Trigonometric tangent.",
			unaryFunc: types.BIF_tan,
		},

		{
			name:      "tanh",
			class:     FUNC_CLASS_MATH,
			help:      "Hyperbolic tangent.",
			unaryFunc: types.BIF_tanh,
		},

		{
			name:       "roundm",
			class:      FUNC_CLASS_MATH,
			help:       `Round to nearest multiple of m: roundm($x,$m) is the same as round($x/$m)*$m.`,
			binaryFunc: types.BIF_roundm,
		},

		{
			name:  "urand",
			class: FUNC_CLASS_MATH,
			help:  `Floating-point numbers uniformly distributed on the unit interval.`,
			examples: []string{
				"Int-valued example: '$n=floor(20+urand()*11)'.",
			},
			zaryFunc: types.BIF_urand,
		},

		{
			name:       "urandint",
			class:      FUNC_CLASS_MATH,
			help:       `Integer uniformly distributed between inclusive integer endpoints.`,
			binaryFunc: types.BIF_urandint,
		},

		{
			name:       "urandrange",
			class:      FUNC_CLASS_MATH,
			help:       `Floating-point numbers uniformly distributed on the interval [a, b).`,
			binaryFunc: types.BIF_urandrange,
		},

		{
			name:     "urand32",
			class:    FUNC_CLASS_MATH,
			help:     `Integer uniformly distributed 0 and 2**32-1 inclusive.`,
			zaryFunc: types.BIF_urand32,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_TIME

		{
			name:  "gmt2sec",
			class: FUNC_CLASS_TIME,
			help:  `Parses GMT timestamp as integer seconds since the epoch.`,
			examples: []string{
				`gmt2sec("2001-02-03T04:05:06Z") = 981173106`,
			},
			unaryFunc: types.BIF_gmt2sec,
		},

		{
			name:  "localtime2sec",
			class: FUNC_CLASS_TIME,
			help: `Parses local timestamp as integer seconds since the epoch. Consults $TZ environment variable,
unless second argument is supplied.`,
			examples: []string{
				`localtime2sec("2001-02-03 04:05:06") = 981165906 with TZ="Asia/Istanbul"`,
				`localtime2sec("2001-02-03 04:05:06", "Asia/Istanbul") = 981165906"`,
			},
			// TODO: help-string
			unaryFunc:          types.BIF_localtime2sec_unary,
			binaryFunc:         types.BIF_localtime2sec_binary,
			hasMultipleArities: true,
		},

		{
			name:  "sec2gmt",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch as GMT timestamp. Leaves non-numbers as-is. With second integer
argument n, includes n decimal places for the seconds part.`,
			examples: []string{
				`sec2gmt(1234567890)           = "2009-02-13T23:31:30Z"`,
				`sec2gmt(1234567890.123456)    = "2009-02-13T23:31:30Z"`,
				`sec2gmt(1234567890.123456, 6) = "2009-02-13T23:31:30.123456Z"`,
			},
			unaryFunc:          types.BIF_sec2gmt_unary,
			binaryFunc:         types.BIF_sec2gmt_binary,
			hasMultipleArities: true,
		},

		{
			name:  "sec2localtime",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch (integer part) as local timestamp.  Consults $TZ
environment variable unless third argument is supplied. Leaves non-numbers as-is. With second integer argument n,
includes n decimal places for the seconds part`,
			examples: []string{
				`sec2localtime(1234567890)           = "2009-02-14 01:31:30"        with TZ="Asia/Istanbul"`,
				`sec2localtime(1234567890.123456)    = "2009-02-14 01:31:30"        with TZ="Asia/Istanbul"`,
				`sec2localtime(1234567890.123456, 6) = "2009-02-14 01:31:30.123456" with TZ="Asia/Istanbul"`,
				`sec2localtime(1234567890.123456, 6, "Asia/Istanbul") = "2009-02-14 01:31:30.123456"`,
			},
			unaryFunc:          types.BIF_sec2localtime_unary,
			binaryFunc:         types.BIF_sec2localtime_binary,
			ternaryFunc:        types.BIF_sec2localtime_ternary,
			hasMultipleArities: true,
		},

		{
			name:  "sec2gmtdate",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch (integer part) as GMT timestamp with year-month-date.
Leaves non-numbers as-is.`,
			examples: []string{
				`sec2gmtdate(1440768801.7) = "2015-08-28".`,
			},
			unaryFunc: types.BIF_sec2gmtdate,
		},

		{
			name:  "sec2localdate",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since epoch (integer part) as local timestamp with year-month-date.
Leaves non-numbers as-is. Consults $TZ environment variable unless second argument is supplied.`,
			examples: []string{
				`sec2localdate(1440768801.7) = "2015-08-28" with TZ="Asia/Istanbul"`,
				`sec2localdate(1440768801.7, "Asia/Istanbul") = "2015-08-28"`,
			},
			unaryFunc:          types.BIF_sec2localdate_unary,
			binaryFunc:         types.BIF_sec2localdate_binary,
			hasMultipleArities: true,
		},

		{
			name:  "localtime2gmt",
			class: FUNC_CLASS_TIME,
			help: `Convert from a local-time string to a GMT-time string. Consults $TZ unless second argument
is supplied.`,
			examples: []string{
				`localtime2gmt("2000-01-01 00:00:00") = "1999-12-31T22:00:00Z" with TZ="Asia/Istanbul"`,
				`localtime2gmt("2000-01-01 00:00:00", "Asia/Istanbul") = "1999-12-31T22:00:00Z"`,
			},
			unaryFunc:          types.BIF_localtime2gmt_unary,
			binaryFunc:         types.BIF_localtime2gmt_binary,
			hasMultipleArities: true,
		},

		{
			name:  "gmt2localtime",
			class: FUNC_CLASS_TIME,
			help: `Convert from a GMT-time string to a local-time string. Consulting $TZ unless second argument
is supplied.`,
			examples: []string{
				`gmt2localtime("1999-12-31T22:00:00Z") = "2000-01-01 00:00:00" with TZ="Asia/Istanbul"`,
				`gmt2localtime("1999-12-31T22:00:00Z", "Asia/Istanbul") = "2000-01-01 00:00:00"`,
			},
			unaryFunc:          types.BIF_gmt2localtime_unary,
			binaryFunc:         types.BIF_gmt2localtime_binary,
			hasMultipleArities: true,
		},

		{
			name:  "strftime",
			class: FUNC_CLASS_TIME,
			help: `Formats seconds since the epoch as timestamp. Format strings are as in the C library
(please see "man strftime" on your system), with the Miller-specific addition of "%1S" through "%9S" which
format the seconds with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.) See also
strftime_local.`,
			examples: []string{
				`strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ")  = "2015-08-28T13:33:21Z"`,
				`strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z"`,
			},
			binaryFunc: types.BIF_strftime,
		},

		{
			name:  "strptime",
			class: FUNC_CLASS_TIME,
			help:  `strptime: Parses timestamp as floating-point seconds since the epoch. See also strptime_local.`,
			examples: []string{
				`strptime("2015-08-28T13:33:21Z",      "%Y-%m-%dT%H:%M:%SZ")   = 1440768801.000000`,
				`strptime("2015-08-28T13:33:21.345Z",  "%Y-%m-%dT%H:%M:%SZ")   = 1440768801.345000`,
				`strptime("1970-01-01 00:00:00 -0400", "%Y-%m-%d %H:%M:%S %z") = 14400`,
				`strptime("1970-01-01 00:00:00 EET",   "%Y-%m-%d %H:%M:%S %Z") = -7200`,
			},

			binaryFunc: types.BIF_strptime,
		},

		{
			name:  "strftime_local",
			class: FUNC_CLASS_TIME,
			help:  `Like strftime but consults the $TZ environment variable to get local time zone.`,
			examples: []string{
				`strftime_local(1440768801.7, "%Y-%m-%d %H:%M:%S %z")  = "2015-08-28 16:33:21 +0300" with TZ="Asia/Istanbul"`,
				`strftime_local(1440768801.7, "%Y-%m-%d %H:%M:%3S %z") = "2015-08-28 16:33:21.700 +0300" with TZ="Asia/Istanbul"`,
				`strftime_local(1440768801.7, "%Y-%m-%d %H:%M:%3S %z", "Asia/Istanbul") = "2015-08-28 16:33:21.700 +0300"`,
			},
			binaryFunc:         types.BIF_strftime_local_binary,
			ternaryFunc:        types.BIF_strftime_local_ternary,
			hasMultipleArities: true,
		},

		{
			name:  "strptime_local",
			class: FUNC_CLASS_TIME,
			help:  `Like stpftime but consults the $TZ environment variable to get local time zone.`,
			examples: []string{
				`strptime_local("2015-08-28T13:33:21Z",    "%Y-%m-%dT%H:%M:%SZ") = 1440758001     with TZ="Asia/Istanbul"`,
				`strptime_local("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440758001.345 with TZ="Asia/Istanbul"`,
				`strptime_local("2015-08-28 13:33:21",     "%Y-%m-%d %H:%M:%S")  = 1440758001     with TZ="Asia/Istanbul"`,
				`strptime_local("2015-08-28 13:33:21",     "%Y-%m-%d %H:%M:%S", "Asia/Istanbul") = 1440758001`,
				// TODO: fix parse error on decimal part
				//`strptime_local("2015-08-28 13:33:21.345","%Y-%m-%d %H:%M:%S") = 1440758001.345`,
			},
			binaryFunc:         types.BIF_strptime_local_binary,
			ternaryFunc:        types.BIF_strptime_local_ternary,
			hasMultipleArities: true,
		},

		{
			name:      "dhms2fsec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000`,
			unaryFunc: types.BIF_dhms2fsec,
		},

		{
			name:      "dhms2sec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000`,
			unaryFunc: types.BIF_dhms2sec,
		},

		{
			name:      "fsec2dhms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"`,
			unaryFunc: types.BIF_fsec2dhms,
		},

		{
			name:      "fsec2hms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"`,
			unaryFunc: types.BIF_fsec2hms,
		},

		{
			name:      "hms2fsec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000`,
			unaryFunc: types.BIF_hms2fsec,
		},

		{
			name:      "hms2sec",
			class:     FUNC_CLASS_TIME,
			help:      `Recovers integer seconds as in hms2sec("01:23:20") = 5000`,
			unaryFunc: types.BIF_hms2sec,
		},

		{
			name:      "sec2dhms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"`,
			unaryFunc: types.BIF_sec2dhms,
		},

		{
			name:      "sec2hms",
			class:     FUNC_CLASS_TIME,
			help:      `Formats integer seconds as in sec2hms(5000) = "01:23:20"`,
			unaryFunc: types.BIF_sec2hms,
		},

		{
			name:     "systime",
			class:    FUNC_CLASS_TIME,
			help:     "help string will go here",
			zaryFunc: types.BIF_systime,
		},

		{
			name:     "systimeint",
			class:    FUNC_CLASS_TIME,
			help:     "help string will go here",
			zaryFunc: types.BIF_systimeint,
		},

		{
			name:     "uptime",
			class:    FUNC_CLASS_TIME,
			help:     "help string will go here",
			zaryFunc: types.BIF_uptime,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_TYPING

		{
			name:      "is_absent",
			class:     FUNC_CLASS_TYPING,
			help:      "False if field is present in input, true otherwise",
			unaryFunc: types.BIF_is_absent,
		},

		{
			name:      "is_array",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is an array.",
			unaryFunc: types.BIF_is_array,
		},

		{
			name:      "is_bool",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with boolean value. Synonymous with is_boolean.",
			unaryFunc: types.BIF_is_bool,
		},

		{
			name:      "is_boolean",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with boolean value. Synonymous with is_bool.",
			unaryFunc: types.BIF_is_boolean,
		},

		{
			name:      "is_empty",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present in input with empty string value, false otherwise.",
			unaryFunc: types.BIF_is_empty,
		},

		{
			name:      "is_empty_map",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is a map which is empty.",
			unaryFunc: types.BIF_is_emptymap,
		},

		{
			name:      "is_error",
			class:     FUNC_CLASS_TYPING,
			help:      "True if if argument is an error, such as taking string length of an integer.",
			unaryFunc: types.BIF_is_error,
		},

		{
			name:      "is_float",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with value inferred to be float",
			unaryFunc: types.BIF_is_float,
		},

		{
			name:      "is_int",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with value inferred to be int",
			unaryFunc: types.BIF_is_int,
		},

		{
			name:      "is_map",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is a map.",
			unaryFunc: types.BIF_is_map,
		},

		{
			name:      "is_nonempty_map",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is a map which is non-empty.",
			unaryFunc: types.BIF_is_nonemptymap,
		},

		{
			name:      "is_not_empty",
			class:     FUNC_CLASS_TYPING,
			help:      "False if field is present in input with empty value, true otherwise",
			unaryFunc: types.BIF_is_notempty,
		},

		{
			name:      "is_not_map",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is not a map.",
			unaryFunc: types.BIF_is_notmap,
		},

		{
			name:      "is_not_array",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is not an array.",
			unaryFunc: types.BIF_is_notarray,
		},

		{
			name:      "is_not_null",
			class:     FUNC_CLASS_TYPING,
			help:      "False if argument is null (empty or absent), true otherwise.",
			unaryFunc: types.BIF_is_notnull,
		},

		{
			name:      "is_null",
			class:     FUNC_CLASS_TYPING,
			help:      "True if argument is null (empty or absent), false otherwise.",
			unaryFunc: types.BIF_is_null,
		},

		{
			name:      "is_numeric",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with value inferred to be int or float",
			unaryFunc: types.BIF_is_numeric,
		},

		{
			name:      "is_present",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present in input, false otherwise.",
			unaryFunc: types.BIF_is_present,
		},

		{
			name:      "is_string",
			class:     FUNC_CLASS_TYPING,
			help:      "True if field is present with string (including empty-string) value",
			unaryFunc: types.BIF_is_string,
		},

		{
			name:                 "asserting_absent",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_absent on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_absent,
		},

		{
			name:                 "asserting_array",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_array on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_array,
		},

		{
			name:                 "asserting_bool",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_bool on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_bool,
		},

		{
			name:                 "asserting_boolean",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_boolean on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_boolean,
		},

		{
			name:                 "asserting_error",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_error on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_error,
		},

		{
			name:                 "asserting_empty",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_empty on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_empty,
		},

		{
			name:                 "asserting_empty_map",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_empty_map on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_emptyMap,
		},

		{
			name:                 "asserting_float",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_float on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_float,
		},

		{
			name:                 "asserting_int",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_int on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_int,
		},

		{
			name:                 "asserting_map",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_map on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_map,
		},

		{
			name:                 "asserting_nonempty_map",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_nonempty_map on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_nonempty_map,
		},

		{
			name:                 "asserting_not_empty",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_not_empty on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_not_empty,
		},

		{
			name:                 "asserting_not_map",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_not_map on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_not_map,
		},

		{
			name:                 "asserting_not_array",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_not_array on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_not_array,
		},

		{
			name:                 "asserting_not_null",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_not_null on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_not_null,
		},

		{
			name:                 "asserting_null",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_null on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_null,
		},

		{
			name:                 "asserting_numeric",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_numeric on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_numeric,
		},

		{
			name:                 "asserting_present",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_present on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_present,
		},

		{
			name:                 "asserting_string",
			class:                FUNC_CLASS_TYPING,
			help:                 `Aborts with an error if is_string on the argument returns false, else returns its argument.`,
			unaryFuncWithContext: types.BIF_asserting_string,
		},

		{
			name:      "typeof",
			class:     FUNC_CLASS_TYPING,
			help:      "Convert argument to type of argument (e.g. \"str\"). For debug.",
			unaryFunc: types.BIF_typeof,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_CONVERSION

		{
			name:      "boolean",
			class:     FUNC_CLASS_CONVERSION,
			help:      "Convert int/float/bool/string to boolean.",
			unaryFunc: types.BIF_boolean,
		},

		{
			name:      "float",
			class:     FUNC_CLASS_CONVERSION,
			help:      "Convert int/float/bool/string to float.",
			unaryFunc: types.BIF_float,
		},

		{
			name:  "fmtnum",
			class: FUNC_CLASS_CONVERSION,
			help: `Convert int/float/bool to string using printf-style format string, e.g.
'$s = fmtnum($n, "%08d")' or '$t = fmtnum($n, "%.6e")'.`,
			binaryFunc: types.BIF_fmtnum,
		},

		{
			name:      "hexfmt",
			class:     FUNC_CLASS_CONVERSION,
			help:      `Convert int to hex string, e.g. 255 to "0xff".`,
			unaryFunc: types.BIF_hexfmt,
		},

		{
			name:      "int",
			class:     FUNC_CLASS_CONVERSION,
			help:      "Convert int/float/bool/string to int.",
			unaryFunc: types.BIF_int,
		},

		{
			name:  "joink",
			class: FUNC_CLASS_CONVERSION,
			help:  `Makes string from map/array keys.`,
			examples: []string{
				`joink({"a":3,"b":4,"c":5}, ",") = "a,b,c".`,
				`joink([1,2,3], ",") = "1,2,3".`,
			},
			binaryFunc: types.BIF_joink,
		},

		{
			name:  "joinv",
			class: FUNC_CLASS_CONVERSION,
			help:  `Makes string from map/array values.`,
			examples: []string{
				`joinv([3,4,5], ",") = "3,4,5"`,
				`joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"`,
			},
			binaryFunc: types.BIF_joinv,
		},

		{
			name:  "joinkv",
			class: FUNC_CLASS_CONVERSION,
			help:  `Makes string from map/array key-value pairs.`,
			examples: []string{
				`joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"`,
				`joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"`,
			},
			ternaryFunc: types.BIF_joinkv,
		},

		{
			name:  "splita",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string into array with type inference.`,
			examples: []string{
				`splita("3,4,5", ",") = [3,4,5]`,
			},
			binaryFunc: types.BIF_splita,
		},

		{
			name:  "splitax",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string into array without type inference.`,
			examples: []string{
				`splita("3,4,5", ",") = ["3","4","5"]`,
			},
			binaryFunc: types.BIF_splitax,
		},

		{
			name:  "splitkv",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string by separators into map with type inference.`,
			examples: []string{
				`splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}`,
			},
			ternaryFunc: types.BIF_splitkv,
		},

		{
			name:  "splitkvx",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string by separators into map without type inference (keys and values are strings).`,
			examples: []string{
				`splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}`,
			},
			ternaryFunc: types.BIF_splitkvx,
		},

		{
			name:  "splitnv",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string by separator into integer-indexed map with type inference.`,
			examples: []string{
				`splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}`,
			},
			binaryFunc: types.BIF_splitnv,
		},

		{
			name:  "splitnvx",
			class: FUNC_CLASS_CONVERSION,
			help:  `Splits string by separator into integer-indexed map without type inference (values are strings).`,
			examples: []string{
				`splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}`,
			},
			binaryFunc: types.BIF_splitnvx,
		},

		{
			name:      "string",
			class:     FUNC_CLASS_CONVERSION,
			help:      "Convert int/float/bool/string/array/map to string.",
			unaryFunc: types.BIF_string,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_COLLECTIONS

		{
			name:       "append",
			class:      FUNC_CLASS_COLLECTIONS,
			help:       "Appends second argument to end of first argument, which must be an array.",
			binaryFunc: types.BIF_append,
		},

		{
			name:  "arrayify",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Walks through a nested map/array, converting any map with consecutive keys
"1", "2", ... into an array. Useful to wrap the output of unflatten.`,
			unaryFunc: types.BIF_arrayify,
		},

		{
			name:      "depth",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      "Prints maximum depth of map/array. Scalars have depth 0.",
			unaryFunc: types.BIF_depth,
		},

		{
			name:  "flatten",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Flattens multi-level maps to single-level ones. Useful for nested JSON-like structures
for non-JSON file formats like CSV.`,
			examples: []string{
				`flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.`,
				`flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.`,
				`Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).`,
			},
			binaryFunc:         types.BIF_flatten_binary,
			ternaryFunc:        types.BIF_flatten,
			hasMultipleArities: true,
		},

		{
			name:      "get_keys",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      "Returns array of keys of map or array",
			unaryFunc: types.BIF_get_keys,
		},

		{
			name:      "get_values",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      "Returns array of keys of map or array -- in the latter case, returns a copy of the array",
			unaryFunc: types.BIF_get_values,
		},

		{
			name:  "haskey",
			class: FUNC_CLASS_COLLECTIONS,
			help: `True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or 'haskey(mymap, mykey)',
or true/false if array index is in bounds / out of bounds.  Error if 1st argument is not a map or array. Note
-n..-1 alias to 1..n in Miller arrays.`,
			binaryFunc: types.BIF_haskey,
		},

		{
			name:      "json_parse",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      `Converts value from JSON-formatted string.`,
			unaryFunc: types.BIF_json_parse,
		},
		{
			name:  "json_stringify",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Converts value to JSON-formatted string. Default output is single-line.
With optional second boolean argument set to true, produces multiline output.`,
			unaryFunc:          types.BIF_json_stringify_unary,
			binaryFunc:         types.BIF_json_stringify_binary,
			hasMultipleArities: true,
		},

		{
			name:      "leafcount",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      `Counts total number of terminal values in map/array. For single-level map/array, same as length.`,
			unaryFunc: types.BIF_leafcount,
		},

		{
			name:      "length",
			class:     FUNC_CLASS_COLLECTIONS,
			help:      "Counts number of top-level entries in array/map. Scalars have length 1.",
			unaryFunc: types.BIF_length,
		},

		{
			name:  "mapdiff",
			class: FUNC_CLASS_COLLECTIONS,
			help: `With 0 args, returns empty map. With 1 arg, returns copy of arg.  With 2 or more,
returns copy of arg 1 with all keys from any of remaining argument maps removed.`,
			variadicFunc: types.BIF_mapdiff,
		},

		{
			name:  "mapexcept",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.  E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.`,
			variadicFunc:         types.BIF_mapexcept,
			minimumVariadicArity: 1,
		},

		{
			name:  "mapselect",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.  E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is
'{1:2,5:6}' and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.`,
			variadicFunc:         types.BIF_mapselect,
			minimumVariadicArity: 1,
		},

		{
			name:  "mapsum",
			class: FUNC_CLASS_COLLECTIONS,
			help: `With 0 args, returns empty map. With >= 1 arg, returns a map with key-value pairs
from all arguments. Rightmost collisions win, e.g.  'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.`,
			variadicFunc: types.BIF_mapsum,
		},

		{
			name:  "unflatten",
			class: FUNC_CLASS_COLLECTIONS,
			help: `Reverses flatten. Useful for nested JSON-like structures for non-JSON file formats like CSV.
See also arrayify.`,
			examples: []string{
				`unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.`,
			},
			binaryFunc: types.BIF_unflatten,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_HOFS

		// Note: most UDFs are in the types package, but these use UDFs which are in the dsl/cst.

		{
			name:  "select",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, includes each
input element in the output if the function returns true. For arrays, the function should take one argument,
for array element; for maps, it should take two, for map-element key and value. In either case it should
return a boolean.`,
			examples: []string{
				`Array example: select([1,2,3,4,5], func(e) {return e >= 3}) returns [3, 4, 5].`,
				`Map example: select({"a":1, "b":3, "c":5}, func(k,v) {return v >= 3}) returns {"b":3, "c": 5}.`,
			},
			binaryFuncWithState: SelectHOF,
		},

		{
			name:  "apply",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, applies the
function to each element of the array/map.  For arrays, the function should take one argument, for array
element; it should return a new element. For maps, it should take two arguments, for map-element key and
value; it should return a new key-value pair (i.e. a single-entry map).`,
			examples: []string{
				`Array example: apply([1,2,3,4,5], func(e) {return e ** 3}) returns [1, 8, 27, 64, 125].`,
				`Map example: apply({"a":1, "b":3, "c":5}, func(k,v) {return {toupper(k): v ** 2}}) returns {"A": 1, "B":9, "C": 25}",`,
			},
			binaryFuncWithState: ApplyHOF,
		},

		{
			name:  "reduce",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, accumulates entries
into a final output -- for example, sum or product. For arrays, the function should take two arguments, for
accumulated value and array element, and return the accumulated element. For maps, it should take four
arguments, for accumulated key and value, and map-element key and value; it should return the updated
accumulator as a new key-value pair (i.e. a single-entry map). The start value for the accumulator is the
first element for arrays, or the first element's key-value pair for maps.`,
			examples: []string{
				`Array example: reduce([1,2,3,4,5], func(acc,e) {return acc + e**3}) returns 225.`,
				`Map example: reduce({"a":1, "b":3, "c": 5}, func(acck,accv,ek,ev) {return {"sum_of_squares": accv + ev**2}}) returns {"sum_of_squares": 35}.`,
			},
			binaryFuncWithState: ReduceHOF,
		},

		{
			name:  "fold",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, accumulates
entries into a final output -- for example, sum or product. For arrays, the function should take two
arguments, for accumulated value and array element. For maps, it should take four arguments, for accumulated
key and value, and map-element key and value; it should return the updated accumulator as a new key-value pair
(i.e. a single-entry map). The start value for the accumulator is taken from the third argument.`,
			examples: []string{
				`Array example: fold([1,2,3,4,5], func(acc,e) {return acc + e**3}, 10000) returns 10225.`,
				`Map example: fold({"a":1, "b":3, "c": 5}, func(acck,accv,ek,ev) {return {"sum": accv+ev**2}}, {"sum":10000}) returns 10035.`,
			},
			ternaryFuncWithState: FoldHOF,
		},

		{
			name:  "sort",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and string flags or function as optional second argument,
returns a sorted copy of the input. With one argument, sorts array elements naturally, and maps naturally by
map keys.  If the second argument is a string, it can contain any of "f" for lexical (default "n" for
natural/numeric), "), "c" for case-folded lexical, and "r" for reversed/descending sort. If the second
argument is a function, then for arrays it should take two arguments a and b, returning < 0, 0, or > 0 as a <
b, a == b, or a > b respectively; for maps the function should take four arguments ak, av, bk, and bv, again
returning < 0, 0, or > 0, using a and b's keys and values.`,
			examples: []string{
				`Array example: sort([5,2,3,1,4], func(a,b) {return b <=> a}) returns [5,4,3,2,1].`,
				`Map example: sort({"c":2,"a":3,"b":1}, func(ak,av,bk,bv) {return bv <=> av}) returns {"a":3,"c":2,"b":1}.`,
			},
			variadicFuncWithState: SortHOF,
			minimumVariadicArity:  1,
			maximumVariadicArity:  2,
		},

		{
			name:  "any",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, yields a boolean true
if the argument function returns true for any array/map element, false otherwise.  For arrays, the function
should take one argument, for array element; for maps, it should take two, for map-element key and value. In
either case it should return a boolean.`,
			examples: []string{
				`Array example: any([10,20,30], func(e) {return $index == e})`,
				`Map example: any({"a": "foo", "b": "bar"}, func(k,v) {return $[k] == v})`,
			},
			binaryFuncWithState: AnyHOF,
		},

		{
			name:  "every",
			class: FUNC_CLASS_HOFS,
			help: `Given a map or array as first argument and a function as second argument, yields a boolean true
if the argument function returns true for every array/map element, false otherwise.  For arrays, the function
should take one argument, for array element; for maps, it should take two, for map-element key and value. In
either case it should return a boolean.`,
			examples: []string{
				`Array example: every(["a", "b", "c"], func(e) {return $[e] >= 0})`,
				`Map example: every({"a": "foo", "b": "bar"}, func(k,v) {return $[k] == v})`,
			},
			binaryFuncWithState: EveryHOF,
		},

		// ----------------------------------------------------------------
		// FUNC_CLASS_SYSTEM

		{
			name:     "hostname",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the hostname as a string.`,
			zaryFunc: types.BIF_hostname,
		},

		{
			name:     "os",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the operating-system name as a string.`,
			zaryFunc: types.BIF_os,
		},

		{
			name:      "system",
			class:     FUNC_CLASS_SYSTEM,
			help:      `Run command string, yielding its stdout minus final carriage return.`,
			unaryFunc: types.BIF_system,
		},

		{
			name:     "version",
			class:    FUNC_CLASS_SYSTEM,
			help:     `Returns the Miller version as a string.`,
			zaryFunc: types.BIF_version,
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

func (manager *BuiltinFunctionManager) ListBuiltinFunctionNamesVertically() {
	for _, builtinFunctionInfo := range *manager.lookupTable {
		fmt.Println(builtinFunctionInfo.name)
	}
}

func (manager *BuiltinFunctionManager) ListBuiltinFunctionNamesAsParagraph() {
	functionNames := make([]string, len(*manager.lookupTable))
	for i, builtinFunctionInfo := range *manager.lookupTable {
		functionNames[i] = builtinFunctionInfo.name
	}
	lib.PrintWordsAsParagraph(functionNames)
}

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
		manager.showSingleUsage(&builtinFunctionInfo)
	}
}

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
			manager.showSingleUsage(&builtinFunctionInfo)
		}
	}
}

func (manager *BuiltinFunctionManager) ListBuiltinFunctionUsage(functionName string) {
	// xxx remove bool-arg code smell
	if !manager.TryListBuiltinFunctionUsage(functionName, true) {
		fmt.Fprintf(os.Stderr, "Function \"%s\" not found.\n", functionName)
	}
}

func (manager *BuiltinFunctionManager) TryListBuiltinFunctionUsage(
	functionName string,
	showApproximate bool,
) bool {
	builtinFunctionInfo := manager.LookUp(functionName)
	if builtinFunctionInfo == nil {
		if showApproximate {
			manager.listBuiltinFunctionUsageApproximate(functionName)
		}
		return false
	}
	manager.listBuiltinFunctionUsageExact(builtinFunctionInfo)
	return true
}

func (manager *BuiltinFunctionManager) listBuiltinFunctionUsageExact(
	builtinFunctionInfo *BuiltinFunctionInfo,
) {
	manager.showSingleUsage(builtinFunctionInfo)
}

func (manager *BuiltinFunctionManager) showSingleUsage(
	builtinFunctionInfo *BuiltinFunctionInfo,
) {
	lib.InternalCodingErrorIf(builtinFunctionInfo.help == "")

	fmt.Printf("%s  (class=%s #args=%s) %s\n",
		colorizer.MaybeColorizeHelp(builtinFunctionInfo.name, true),
		builtinFunctionInfo.class,
		describeNargs(builtinFunctionInfo),
		builtinFunctionInfo.JoinHelp(),
	)
	if len(builtinFunctionInfo.examples) == 1 {
		fmt.Println("Example:")
	}
	if len(builtinFunctionInfo.examples) > 1 {
		fmt.Println("Examples:")
	}
	for _, example := range builtinFunctionInfo.examples {
		fmt.Println(example)
	}
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
		if info.ternaryFuncWithState != nil {
			return "3"
		}
		if info.variadicFunc != nil || info.variadicFuncWithState != nil {
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

var multiSpaceRegex = regexp.MustCompile(`\s+`)

// JoinHelp must be used to format any function help-strings for output. for
// source-code storage, these have newlines in them. For any presentation to
// the user, they must be formatted using the JoinHelp() method which joins
// newlines. This is crucial for rendering of help-strings for manual page,
// webdocs, etc wherein we must let the user's resizing of the terminal window
// or browser determine -- at their choosing -- where lines wrap.
func (info *BuiltinFunctionInfo) JoinHelp() string {
	return multiSpaceRegex.ReplaceAllString(strings.ReplaceAll(info.help, "\n", " "), " ")
}

// ================================================================
// This is a singleton so the online-help functions can query it for listings,
// online help, etc.
var BuiltinFunctionManagerInstance *BuiltinFunctionManager = NewBuiltinFunctionManager()
