package lib

import (
	"math"
)

// ================================================================
// ABOUT DISPOSITION MATRICES/VECTORS
//
// Mlrvals can be of type MT_STRING, MT_INT, MT_FLOAT, MT_BOOLEAN, as well as
// MT_ABSENT, MT_VOID, and ERROR.  Thus when we do pairwise operations on them
// (for binary operators) or singly (for unary operators), what we do depends
// on the type.
//
// Mlrval type enums are 0-up integers precisely so that instead of if-elsing
// or switching on the types, we can instead define tables of function pointers
// and jump immediately to the right thing to do for a given type pairing.  For
// example: adding two ints, or an int and a float, or int and boolean (the
// latter being an error).
//
// The next-past-highest mlrval type enum is called MT_DIM and that is the
// dimension of the binary-operator disposition matrices and unary-operator
// disposition vectors.
//
// Note that not every operation uses disposition matrices. If something makes
// sense only for pairs of strings and nothing else, it makes sense for the
// implementing method to return an MT_STRING result if both arguments are
// MT_STRING, or MT_ERROR otherwise.
// ================================================================

// Function-pointer type for unary-operator disposition vectors.
type unaryFunc func(*Mlrval) Mlrval

// Function-pointer type for binary-operator disposition matrices.
type binaryFunc func(*Mlrval, *Mlrval) Mlrval

// ----------------------------------------------------------------
// The following are frequently used in disposition matrices for various
// operators and are defined here for re-use. The names are VERY short,
// and all the same length, so that the disposition matrices will look
// reasonable rectangular even after gofmt has been run.

// Return error (unary)
func _erro1(val1 *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (unary)
func _absn1(val1 *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (unary)
func _void1(val1 *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return argument (unary)
func _1u___(val1 *Mlrval) Mlrval {
	return *val1
}

// Return error (binary)
func _erro(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (binary)
func _absn(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (binary)
func _void(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromVoid()
}

// Return first argument (binary)
func _1___(val1, val2 *Mlrval) Mlrval {
	return *val1
}

// Return second argument (binary)
func _2___(val1, val2 *Mlrval) Mlrval {
	return *val2
}

// Return first argument, as string (binary)
func _s1__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String())
}

// Return second argument, as string (binary)
func _s2__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val2.String())
}

// Return integer zero (binary)
func _i0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
}

// Return float zero (binary)
func _f0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(0.0)
}

// ================================================================
// Unary plus operator

func upos_i_i(val1 *Mlrval) Mlrval {
	return *val1
}

func upos_f_f(val1 *Mlrval) Mlrval {
	return *val1
}

var upos_dispositions = [MT_DIM]unaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*EMPTY  */ _void1,
	/*STRING */ _erro1,
	/*INT    */ _1u___,
	/*FLOAT  */ _1u___,
	/*BOOL   */ _erro1,
}

func MlrvalUnaryPlus(val1 *Mlrval) Mlrval {
	return upos_dispositions[val1.mvtype](val1)
}

// ================================================================
// Unary minus operator

func uneg_i_i(val1 *Mlrval) Mlrval {
	return MlrvalFromInt64(-val1.intval)
}

func uneg_f_f(val1 *Mlrval) Mlrval {
	return MlrvalFromFloat64(-val1.floatval)
}

var uneg_dispositions = [MT_DIM]unaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*EMPTY  */ _void1,
	/*STRING */ _erro1,
	/*INT    */ uneg_i_i,
	/*FLOAT  */ uneg_f_f,
	/*BOOL   */ _erro1,
}

func MlrvalUnaryMinus(val1 *Mlrval) Mlrval {
	return uneg_dispositions[val1.mvtype](val1)
}

// ================================================================
// Dot operator, with loose typecasting.
//
// For most operations, I don't like loose typecasting -- for example, in PHP
// "10" + 2 is the number 12 and in JavaScript it's the string "102", and I
// find both of those horrid and error-prone. In Miller, "10"+2 is MT_ERROR, by
// design, unless intentional casting is done like '$x=int("10")+2'.
//
// However, for dotting, in practice I tipped over and allowed dotting of
// strings and ints: so while "10" + 2 is an error in Miller, '"10". 2' is
// "102". Unlike with "+", with "." there is no ambiguity about what the output
// should be: always the string concatenation of the string representations of
// the two arguments. So, we do the string-cast for the user.

func dot_s_xx(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String() + val2.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT       FLOAT     BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__},
	/*EMPTY  */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
}

func MlrvalDot(val1, val2 *Mlrval) Mlrval {
	return dot_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Addition with auto-overflow from int to float when necessary.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func plus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + float64(val2.intval))
}
func plus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) + val2.floatval)
}
func plus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + val2.floatval)
}

// Auto-overflows up to float.  Additions & subtractions overflow by at most
// one bit so it suffices to check sign-changes.
func plus_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := a + b

	overflowed := false
	if a > 0 {
		if b > 0 && c < 0 {
			overflowed = true
		}
	} else if a < 0 {
		if b < 0 && c > 0 {
			overflowed = true
		}
	}

	if overflowed {
		return MlrvalFromFloat64(float64(a) + float64(b))
	} else {
		return MlrvalFromInt64(c)
	}
}

var plus_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, plus_n_ii, plus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, plus_f_fi, plus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalPlus(val1, val2 *Mlrval) Mlrval {
	return plus_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Subtraction with auto-overflow from int to float when necessary.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func minus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - val2.floatval)
}
func minus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - float64(val2.intval))
}
func minus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) - val2.floatval)
}

// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
func minus_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := a - b

	overflowed := false
	if a > 0 {
		if b < 0 && c < 0 {
			overflowed = true
		}
	} else if a < 0 {
		if b > 0 && c > 0 {
			overflowed = true
		}
	}

	if overflowed {
		return MlrvalFromFloat64(float64(a) - float64(b))
	} else {
		return MlrvalFromInt64(c)
	}
}

var minus_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, minus_n_ii, minus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, minus_f_fi, minus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalMinus(val1, val2 *Mlrval) Mlrval {
	return minus_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Multiplication with auto-overflow from int to float when necessary.  See
// also http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func times_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * float64(val2.intval))
}
func times_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) * val2.floatval)
}
func times_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * val2.floatval)
}

// Auto-overflows up to float.
//
// Unlike adds & subtracts which overflow by at most one bit, multiplies can
// overflow by a word size. Thus detecting sign-changes does not suffice to
// detect overflow. Instead we test whether the floating-point product exceeds
// the representable integer range. Now 64-bit integers have 64-bit precision
// while IEEE-doubles have only 52-bit mantissas -- so, 53 bits including
// implicit leading one.
//
// The following experiment explicitly demonstrates the resolution at this range:
//
//    64-bit integer     64-bit integer     Casted to double           Back to 64-bit
//        in hex           in decimal                                    integer
// 0x7ffffffffffff9ff 9223372036854774271 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffa00 9223372036854774272 9223372036854773760.000000 0x7ffffffffffff800
// 0x7ffffffffffffbff 9223372036854774783 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffc00 9223372036854774784 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffdff 9223372036854775295 9223372036854774784.000000 0x7ffffffffffffc00
// 0x7ffffffffffffe00 9223372036854775296 9223372036854775808.000000 0x8000000000000000
// 0x7ffffffffffffffe 9223372036854775806 9223372036854775808.000000 0x8000000000000000
// 0x7fffffffffffffff 9223372036854775807 9223372036854775808.000000 0x8000000000000000
//
// That is, we cannot check an integer product to see if it is greater than
// 2**63-1 (or is less than -2**63) using integer arithmetic (it may have
// already overflowed) *or* using double-precision (granularity). Instead we
// check if the absolute value of the product exceeds the largest representable
// double less than 2**63. (An alterative would be to do all integer multiplies
// using handcrafted multi-word 128-bit arithmetic).

func times_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval
	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt64(a * b)
	}
}

var times_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, times_n_ii, times_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, times_f_fi, times_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalTimes(val1, val2 *Mlrval) Mlrval {
	return times_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Pythonic division.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.
//
// Int/int pairings don't produce overflow.
//
// IEEE-754 handles float overflow/underflow:
//
//   $ echo 'x=1e-300,y=1e300' | mlr put '$z=$x*$y'
//   x=1e-300,y=1e300,z=1
//
//   $ echo 'x=1e-300,y=1e300' | mlr put '$z=$x/$y'
//   x=1e-300,y=1e300,z=0
//
//   $ echo 'x=1e-300,y=1e300' | mlr put '$z=$y/$x'
//   x=1e-300,y=1e300,z=+Inf

func divide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / float64(val2.intval))
}
func divide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) / val2.floatval)
}
func divide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / val2.floatval)
}

func divide_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	if a%b == 0 {
		return MlrvalFromInt64(a / b)
	} else {
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt64(a * b)
	}
}

var divide_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, divide_n_ii, divide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, divide_f_fi, divide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDivide(val1, val2 *Mlrval) Mlrval {
	return divide_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Integer division: DSL operator '//' as in Python.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func int_divide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(val1.floatval / float64(val2.intval)))
}
func int_divide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(val1.intval) / val2.floatval))
}
func int_divide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(val1.floatval / val2.floatval))
}

func int_divide_n_ii(val1, val2 *Mlrval) Mlrval {
	a := val1.intval
	b := val2.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	q := a / b
	r := a % b
	if a < 0 {
		if b > 0 {
			if r != 0 {
				q--
			}
		}
	} else {
		if b < 0 {
			if r != 0 {
				q--
			}
		}
	}
	return MlrvalFromInt64(q)
}

var int_divide_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, int_divide_n_ii, int_divide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, int_divide_f_fi, int_divide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalIntDivide(val1, val2 *Mlrval) Mlrval {
	return int_divide_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Non-auto-overflowing addition: DSL operator '.+'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotplus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + val2.floatval)
}

func dotplus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval + float64(val2.intval))
}

func dotplus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) + val2.floatval)
}

func dotplus_i_ii(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(val1.intval + val2.intval)
}

var dot_plus_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT    FLOAT         BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, dotplus_i_ii, dotplus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotplus_f_fi, dotplus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotPlus(val1, val2 *Mlrval) Mlrval {
	return dot_plus_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Non-auto-overflowing subtraction: DSL operator '.-'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotminus_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - val2.floatval)
}

func dotminus_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval - float64(val2.intval))
}

func dotminus_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) - val2.floatval)
}

func dotminus_i_ii(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(val1.intval - val2.intval)
}

var dotminus_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT    FLOAT         BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, dotminus_i_ii, dotminus_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotminus_f_fi, dotminus_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotMinus(val1, val2 *Mlrval) Mlrval {
	return dotminus_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ----------------------------------------------------------------
// Non-auto-overflowing multiplication: DSL operator '.*'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dottimes_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * val2.floatval)
}

func dottimes_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval * float64(val2.intval))
}

func dottimes_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) * val2.floatval)
}

func dottimes_i_ii(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(val1.intval * val2.intval)
}

var dottimes_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT    FLOAT         BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, dottimes_i_ii, dottimes_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dottimes_f_fi, dottimes_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotTimes(val1, val2 *Mlrval) Mlrval {
	return dottimes_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotdivide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / val2.floatval)
}

func dotdivide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / float64(val2.intval))
}

func dotdivide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) / val2.floatval)
}

func dotdivide_i_ii(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(val1.intval / val2.intval)
}

var dotdivide_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT    FLOAT         BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, dotdivide_i_ii, dotdivide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotdivide_f_fi, dotdivide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotDivide(val1, val2 *Mlrval) Mlrval {
	return dotdivide_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './/'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotidivide_f_ff(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / val2.floatval)
}

func dotidivide_f_fi(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(val1.floatval / float64(val2.intval))
}

func dotidivide_f_if(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(val1.intval) / val2.floatval)
}

func dotidivide_i_ii(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(val1.intval / val2.intval)
}

var dotidivide_dispositions = [MT_DIM][MT_DIM]binaryFunc{
	//       ERROR ABSENT  EMPTY  STRING INT    FLOAT         BOOL
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _1___, _void, _erro, dotidivide_i_ii, dotidivide_f_if, _erro},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotidivide_f_fi, dotidivide_f_ff, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotIntDivide(val1, val2 *Mlrval) Mlrval {
	return dotidivide_dispositions[val1.mvtype][val2.mvtype](val1, val2)
}




//static mv_t oidiv_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(floor(a / b));
//}
//static mv_t oidiv_i_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//
//	// Pythonic division, not C division.
//	long long q = a / b;
//	long long r = a % b;
//	if (a < 0) {
//		if (b > 0) {
//			if (r != 0)
//				q--;
//		}
//	} else {
//		if (b < 0) {
//			if (r != 0)
//				q--;
//		}
//	}
//	return mv_from_int(q);
//}
//
//static mv_binary_func_t* oidiv_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _i0,        _f0,        _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,       _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  oidiv_i_ii, oidiv_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  oidiv_f_fi, oidiv_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//};
//
//mv_t x_xx_int_odivide_func(mv_t* pval1, mv_t* pval2) {
//	return (oidiv_dispositions[pval1->type][pval2->type])(pval1,pval2);
//}

//
//// ================================================================
//static mv_t mod_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b * floor(a / b));
//}
//static mv_t mod_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a - b * floor(a / b));
//}
//static mv_t mod_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b * floor(a / b));
//}
//static mv_t mod_i_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	if (b == 0LL) { // Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
//		return mv_from_float((double)a / (double)b);
//	}
//	long long u = a % b;
//	// Pythonic division, not C division.
//	if (a >= 0LL) {
//		if (b < 0LL) {
//			u += b;
//		}
//	} else {
//		if (b >= 0LL) {
//			u += b;
//		}
//	}
//	return mv_from_int(u);
//}
//
//static mv_binary_func_t* mod_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT       FLOAT     BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err,     _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _i0,      _f0,      _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,     _emt,     _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err,     _err},
//	/*INT*/    {_err, _1,    _emt, _err,  mod_i_ii, mod_f_if, _err},
//	/*FLOAT*/  {_err, _1,    _emt, _err,  mod_f_fi, mod_f_ff, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err,     _err},
//};
//
//mv_t x_xx_mod_func(mv_t* pval1, mv_t* pval2) {
//	return (mod_dispositions[pval1->type][pval2->type])(pval1,pval2);
//}

//
//// ----------------------------------------------------------------
//static mv_t band_i_ii(mv_t* pa, mv_t* pb) {
//	return mv_from_int(pa->u.intv & pb->u.intv);
//}
//
//static mv_binary_func_t* band_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT        FLOAT BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err, _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _err, _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt, _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err, _err},
//	/*INT*/    {_err, _1,    _emt, _err,  band_i_ii, _err, _err},
//	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,      _err, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err, _err},
//};
//
//mv_t x_xx_band_func(mv_t* pval1, mv_t* pval2) { return (band_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t bor_i_ii(mv_t* pa, mv_t* pb) {
//	return mv_from_int(pa->u.intv | pb->u.intv);
//}
//
//static mv_binary_func_t* bor_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT       FLOAT BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,     _err, _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,       _err, _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,     _emt, _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,     _err, _err},
//	/*INT*/    {_err, _1,    _emt, _err,  bor_i_ii, _err, _err},
//	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,     _err, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,     _err, _err},
//};
//
//mv_t x_xx_bor_func(mv_t* pval1, mv_t* pval2) { return (bor_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t bxor_i_ii(mv_t* pa, mv_t* pb) {
//	return mv_from_int(pa->u.intv ^ pb->u.intv);
//}
//
//static mv_binary_func_t* bxor_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT        FLOAT BOOL
//	/*ERROR*/  {_err, _err,  _err, _err,  _err,      _err, _err},
//	/*ABSENT*/ {_err, _a,    _a,   _err,  _2,        _err, _err},
//	/*EMPTY*/  {_err, _a,    _emt, _err,  _emt,      _emt, _err},
//	/*STRING*/ {_err, _err,  _err, _err,  _err,      _err, _err},
//	/*INT*/    {_err, _1,    _emt, _err,  bxor_i_ii, _err, _err},
//	/*FLOAT*/  {_err, _err,  _emt, _err,  _err,      _err, _err},
//	/*BOOL*/   {_err, _err,  _err, _err,  _err,      _err, _err},
//};
//
//mv_t x_xx_bxor_func(mv_t* pval1, mv_t* pval2) { return (bxor_dispositions[pval1->type][pval2->type])(pval1,pval2); }
//
//// ----------------------------------------------------------------
//static mv_t boolean_b_b(mv_t* pa) { return mv_from_bool(pa->u.boolv); }
//static mv_t boolean_b_f(mv_t* pa) { return mv_from_bool((pa->u.fltv == 0.0) ? FALSE : TRUE); }
//static mv_t boolean_b_i(mv_t* pa) { return mv_from_bool((pa->u.intv == 0LL) ? FALSE : TRUE); }
//static mv_t boolean_b_s(mv_t* pa) { return mv_from_bool((streq(pa->u.strv, "true") || streq(pa->u.strv, "TRUE")) ? TRUE : FALSE);}
//
//static mv_unary_func_t* boolean_dispositions[MT_DIM] = {
//	/*ERROR*/  _err1,
//	/*ABSENT*/ _a1,
//	/*EMPTY*/  _emt1,
//	/*STRING*/ boolean_b_s,
//	/*INT*/    boolean_b_i,
//	/*FLOAT*/  boolean_b_f,
//	/*BOOL*/   boolean_b_b,
//};
//
//mv_t b_x_boolean_func(mv_t* pval1) { return (boolean_dispositions[pval1->type])(pval1); }
//
//// ----------------------------------------------------------------
//static mv_t string_s_b(mv_t* pa) { return mv_from_string_no_free(pa->u.boolv?"true":"false"); }
//static mv_t string_s_f(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_double(pa->u.fltv, MLR_GLOBALS.ofmt)); }
//static mv_t string_s_i(mv_t* pa) { return mv_from_string_with_free(mlr_alloc_string_from_ll(pa->u.intv)); }
//static mv_t string_s_s(mv_t* pa) {
//	char free_flags = pa->free_flags;
//	pa->free_flags = NO_FREE;
//	return mv_from_string(pa->u.strv, free_flags);
//}
//
//static mv_unary_func_t* string_dispositions[MT_DIM] = {
//	/*ERROR*/  _err1,
//	/*ABSENT*/ _a1,
//	/*EMPTY*/  _emt1,
//	/*STRING*/ string_s_s,
//	/*INT*/    string_s_i,
//	/*FLOAT*/  string_s_f,
//	/*BOOL*/   string_s_b,
//};
//
//mv_t s_x_string_func(mv_t* pval1) { return (string_dispositions[pval1->type])(pval1); }
//

//// ----------------------------------------------------------------
//static mv_t eq_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.intv); }
//static mv_t ne_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.intv); }
//static mv_t gt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.intv); }
//static mv_t ge_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.intv); }
//static mv_t lt_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.intv); }
//static mv_t le_b_ii(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.intv); }
//
//static mv_t eq_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.fltv); }
//static mv_t ne_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.fltv); }
//static mv_t gt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.fltv); }
//static mv_t ge_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.fltv); }
//static mv_t lt_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.fltv); }
//static mv_t le_b_ff(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.fltv); }
//
//static mv_t eq_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv == pb->u.intv); }
//static mv_t ne_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv != pb->u.intv); }
//static mv_t gt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >  pb->u.intv); }
//static mv_t ge_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv >= pb->u.intv); }
//static mv_t lt_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <  pb->u.intv); }
//static mv_t le_b_fi(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.fltv <= pb->u.intv); }
//
//static mv_t eq_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv == pb->u.fltv); }
//static mv_t ne_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv != pb->u.fltv); }
//static mv_t gt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >  pb->u.fltv); }
//static mv_t ge_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv >= pb->u.fltv); }
//static mv_t lt_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <  pb->u.fltv); }
//static mv_t le_b_if(mv_t* pa, mv_t* pb) { return mv_from_bool(pa->u.intv <= pb->u.fltv); }
//
//static mv_t eq_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) == 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ne_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) != 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t gt_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) > 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ge_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) >= 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t lt_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) < 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t le_b_xs(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sa = mv_format_val(pa, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(sa, pb->u.strv) <= 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sa);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//
//static mv_t eq_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) == 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ne_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) != 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t gt_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) > 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ge_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) >= 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t lt_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) < 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t le_b_sx(mv_t* pa, mv_t* pb) {
//	char free_flags;
//	char* sb = mv_format_val(pb, &free_flags);
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, sb) <= 0);
//	if (free_flags & FREE_ENTRY_VALUE)
//		free(sb);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//
//static mv_t eq_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) == 0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ne_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) != 0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t gt_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) >  0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t ge_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) >= 0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t lt_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) <  0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//static mv_t le_b_ss(mv_t*pa, mv_t*pb) {
//	mv_t rv = mv_from_bool(strcmp(pa->u.strv, pb->u.strv) <= 0);
//	mv_free(pa);
//	mv_free(pb);
//	return rv;
//}
//
//static mv_binary_func_t* eq_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _err},
//	/*STRING*/ {_err, _a,    eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _err},
//	/*INT*/    {_err, _a,    eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _err},
//	/*FLOAT*/  {_err, _a,    eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//static mv_binary_func_t* ne_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _err},
//	/*STRING*/ {_err, _a,    ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _err},
//	/*INT*/    {_err, _a,    ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _err},
//	/*FLOAT*/  {_err, _a,    ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//static mv_binary_func_t* gt_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _err},
//	/*STRING*/ {_err, _a,    gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _err},
//	/*INT*/    {_err, _a,    gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _err},
//	/*FLOAT*/  {_err, _a,    gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//static mv_binary_func_t* ge_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _err},
//	/*STRING*/ {_err, _a,    ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _err},
//	/*INT*/    {_err, _a,    ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _err},
//	/*FLOAT*/  {_err, _a,    ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//static mv_binary_func_t* lt_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _err},
//	/*STRING*/ {_err, _a,    lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _err},
//	/*INT*/    {_err, _a,    lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _err},
//	/*FLOAT*/  {_err, _a,    lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//static mv_binary_func_t* le_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY    STRING   INT      FLOAT    BOOL
//	/*ERROR*/  {_err, _err,  _err,    _err,    _err,    _err,    _err},
//	/*ABSENT*/ {_err, _a,    _a,      _a,      _a,      _a,      _a},
//	/*EMPTY*/  {_err, _a,    le_b_ss, le_b_ss, le_b_sx, le_b_sx, _err},
//	/*STRING*/ {_err, _a,    le_b_ss, le_b_ss, le_b_sx, le_b_sx, _err},
//	/*INT*/    {_err, _a,    le_b_xs, le_b_xs, le_b_ii, le_b_if, _err},
//	/*FLOAT*/  {_err, _a,    le_b_xs, le_b_xs, le_b_fi, le_b_ff, _err},
//	/*BOOL*/   {_err, _err,  _a,      _err,    _err,    _err,    _err},
//};
//
//mv_t eq_op_func(mv_t* pval1, mv_t* pval2) { return (eq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//mv_t ne_op_func(mv_t* pval1, mv_t* pval2) { return (ne_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//mv_t gt_op_func(mv_t* pval1, mv_t* pval2) { return (gt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//mv_t ge_op_func(mv_t* pval1, mv_t* pval2) { return (ge_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//mv_t lt_op_func(mv_t* pval1, mv_t* pval2) { return (lt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//mv_t le_op_func(mv_t* pval1, mv_t* pval2) { return (le_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//
//// ----------------------------------------------------------------
//int mv_equals_si(mv_t* pa, mv_t* pb) {
//	if (pa->type == MT_INT) {
//		return (pb->type == MT_INT) ? pa->u.intv == pb->u.intv : FALSE;
//	} else {
//		return (pb->type == MT_STRING) ? streq(pa->u.strv, pb->u.strv) : FALSE;
//	}
//}
//
//// ----------------------------------------------------------------
//static int eq_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv == pb->u.intv; }
//static int ne_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv != pb->u.intv; }
//static int gt_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv >  pb->u.intv; }
//static int ge_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv >= pb->u.intv; }
//static int lt_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv <  pb->u.intv; }
//static int le_i_ii(mv_t* pa, mv_t* pb) { return  pa->u.intv <= pb->u.intv; }
//
//static int eq_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv == pb->u.fltv; }
//static int ne_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv != pb->u.fltv; }
//static int gt_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv >  pb->u.fltv; }
//static int ge_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv >= pb->u.fltv; }
//static int lt_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv <  pb->u.fltv; }
//static int le_i_ff(mv_t* pa, mv_t* pb) { return  pa->u.fltv <= pb->u.fltv; }
//
//static int eq_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv == pb->u.intv; }
//static int ne_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv != pb->u.intv; }
//static int gt_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv >  pb->u.intv; }
//static int ge_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv >= pb->u.intv; }
//static int lt_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv <  pb->u.intv; }
//static int le_i_fi(mv_t* pa, mv_t* pb) { return  pa->u.fltv <= pb->u.intv; }
//
//static int eq_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv == pb->u.fltv; }
//static int ne_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv != pb->u.fltv; }
//static int gt_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv >  pb->u.fltv; }
//static int ge_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv >= pb->u.fltv; }
//static int lt_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv <  pb->u.fltv; }
//static int le_i_if(mv_t* pa, mv_t* pb) { return  pa->u.intv <= pb->u.fltv; }
//
//static mv_i_nn_comparator_func_t* ieq_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  eq_i_ii, eq_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  eq_i_fi, eq_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//static mv_i_nn_comparator_func_t* ine_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  ne_i_ii, ne_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  ne_i_fi, ne_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//static mv_i_nn_comparator_func_t* igt_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  gt_i_ii, gt_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  gt_i_fi, gt_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//static mv_i_nn_comparator_func_t* ige_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  ge_i_ii, ge_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  ge_i_fi, ge_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//static mv_i_nn_comparator_func_t* ilt_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  lt_i_ii, lt_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  lt_i_fi, lt_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//static mv_i_nn_comparator_func_t* ile_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT      FLOAT    BOOL
//	/*ERROR*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*ABSENT*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*EMPTY*/  {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*STRING*/ {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//	/*INT*/    {NULL, NULL,  NULL, NULL,  le_i_ii, le_i_if, NULL},
//	/*FLOAT*/  {NULL, NULL,  NULL, NULL,  le_i_fi, le_i_ff, NULL},
//	/*BOOL*/   {NULL, NULL,  NULL, NULL,  NULL,    NULL,    NULL},
//};
//
//int mv_i_nn_eq(mv_t* pval1, mv_t* pval2) { return (ieq_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//int mv_i_nn_ne(mv_t* pval1, mv_t* pval2) { return (ine_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//int mv_i_nn_gt(mv_t* pval1, mv_t* pval2) { return (igt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//int mv_i_nn_ge(mv_t* pval1, mv_t* pval2) { return (ige_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//int mv_i_nn_lt(mv_t* pval1, mv_t* pval2) { return (ilt_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//int mv_i_nn_le(mv_t* pval1, mv_t* pval2) { return (ile_dispositions[pval1->type][pval2->type])(pval1, pval2); }
//
