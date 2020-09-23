package types

import (
	"math"
	"regexp"
	"strings"
	"time"

	"miller/lib"
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
//
// Naming conventions: since these functions fit into disposition matrices, the
// names are kept quite short. Many are of the form 'plus_f_fi', 'eq_b_xs',
// etc. The conventions are:
//
// * The 'plus_', 'eq_', etc is for the name of the operator.
//
// * For binary operators, things like _f_fi indicate the type of the return
//   value (e.g. 'f') and the types of the two arguments (e.g. 'fi').
//
// * For unary operators, things like _i_i indicate the type of the return
//   value and the type of the argument.
//
// * Type names:
//   's' for string
//   'i' for int64
//   'f' for float64
//   'n' for number return types -- e.g. the auto-overflow from
//       int to float plus_n_ii returns MT_INT if integer-additio overflow
//       didn't happen, or MT_FLOAT if it did.
//   'b' for boolean
//   'x' for don't-care slots, e.g. eq_b_sx for comparing MT_STRING
//       ('s') to anything else ('x').
// ================================================================

// Function-pointer type for zary functions.
type ZaryFunc func() Mlrval

// Function-pointer type for unary-operator disposition vectors.
type UnaryFunc func(*Mlrval) Mlrval

// Helps keystroke-saving for wrapping Go math-library functions
// Examples: cos, sin, etc.
type mathLibUnaryFunc func(float64) float64
type mathLibUnaryFuncWrapper func(*Mlrval, mathLibUnaryFunc) Mlrval

// Function-pointer type for binary-operator disposition matrices.
type BinaryFunc func(*Mlrval, *Mlrval) Mlrval

// Function-pointer type for ternary functions
type TernaryFunc func(*Mlrval, *Mlrval, *Mlrval) Mlrval

// Function-pointer type for variadic functions.
type VariadicFunc func([]*Mlrval) Mlrval

// Function-pointer type for sorting. Returns < 0 if a < b, 0 if a == b, > 0 if a > b.
type ComparatorFunc func(*Mlrval, *Mlrval) int

// ================================================================
// The following are frequently used in disposition matrices for various
// operators and are defined here for re-use. The names are VERY short,
// and all the same length, so that the disposition matrices will look
// reasonable rectangular even after gofmt has been run.

// ----------------------------------------------------------------
// Return error (unary)
func _erro1(ma *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (unary)
func _absn1(ma *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (unary)
func _void1(ma *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return argument (unary)
func _1u___(ma *Mlrval) Mlrval {
	return *ma
}

// ----------------------------------------------------------------
// Return error (unary math-library func)
func _math_unary_erro1(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromError()
}

// Return absent (unary math-library func)
func _math_unary_absn1(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (unary math-library func)
func _math_unary_void1(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromAbsent()
}

// ----------------------------------------------------------------
// Return error (binary)
func _erro(ma, mb *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent (binary)
func _absn(ma, mb *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void (binary)
func _void(ma, mb *Mlrval) Mlrval {
	return MlrvalFromVoid()
}

// Return first argument (binary)
func _1___(ma, mb *Mlrval) Mlrval {
	return *ma
}

// Return second argument (binary)
func _2___(ma, mb *Mlrval) Mlrval {
	return *mb
}

// Return first argument, as string (binary)
func _s1__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(ma.String())
}

// Return second argument, as string (binary)
func _s2__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(mb.String())
}

// Return integer zero (binary)
func _i0__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
}

// Return float zero (binary)
func _f0__(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(0.0)
}

// ================================================================
// ZARY FUNCTIONS
// ================================================================

func MlrvalSystime() Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func MlrvalSystimeInt() Mlrval {
	return MlrvalFromInt64(time.Now().Unix())
}

func MlrvalUrand() Mlrval {
	return MlrvalFromFloat64(
		lib.RandFloat64(),
	)
}

func MlrvalUrand32() Mlrval {
	return MlrvalFromInt64(
		int64(
			lib.RandUint32(),
		),
	)
}

// ================================================================
// UNARY FUNCTIONS
// ================================================================

// ================================================================
// Unary plus operator

var upos_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*EMPTY  */ _void1,
	/*STRING */ _erro1,
	/*INT    */ _1u___,
	/*FLOAT  */ _1u___,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
}

func MlrvalUnaryPlus(ma *Mlrval) Mlrval {
	return upos_dispositions[ma.mvtype](ma)
}

// ================================================================
// Unary minus operator

func uneg_i_i(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(-ma.intval)
}

func uneg_f_f(ma *Mlrval) Mlrval {
	return MlrvalFromFloat64(-ma.floatval)
}

var uneg_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*EMPTY  */ _void1,
	/*STRING */ _erro1,
	/*INT    */ uneg_i_i,
	/*FLOAT  */ uneg_f_f,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
}

func MlrvalUnaryMinus(ma *Mlrval) Mlrval {
	return uneg_dispositions[ma.mvtype](ma)
}

// ================================================================
// Logical NOT operator

func MlrvalLogicalNOT(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL {
		return MlrvalFromBool(!ma.boolval)
	} else {
		return MlrvalFromError()
	}
}

// ================================================================
// Bitwise NOT

func bitwise_not_i_i(ma *Mlrval) Mlrval {
	return MlrvalFromInt64(^ma.intval)
}

var bitwise_not_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*EMPTY  */ _void1,
	/*STRING */ _erro1,
	/*INT    */ bitwise_not_i_i,
	/*FLOAT  */ _erro1,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
}

func MlrvalBitwiseNOT(ma *Mlrval) Mlrval {
	return bitwise_not_dispositions[ma.mvtype](ma)
}

// ================================================================
// Go math-library functions

func math_unary_f_i(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromFloat64(f(float64(ma.intval)))
}
func math_unary_f_f(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromFloat64(f(ma.floatval))
}

var mudispo = [MT_DIM]mathLibUnaryFuncWrapper{
	/*ERROR  */ _math_unary_erro1,
	/*ABSENT */ _math_unary_absn1,
	/*EMPTY  */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*INT    */ math_unary_f_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
}

func MlrvalAbs(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Abs) }
func MlrvalAcos(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Acos) }
func MlrvalAcosh(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Acosh) }
func MlrvalAsin(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Asin) }
func MlrvalAsinh(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Asinh) }
func MlrvalAtan(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Atan) }
func MlrvalAtanh(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Atanh) }
func MlrvalCbrt(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Cbrt) }
func MlrvalCeil(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Ceil) }
func MlrvalCos(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Cos) }
func MlrvalCosh(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Cosh) }
func MlrvalErf(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Erf) }
func MlrvalErfc(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Erfc) }
func MlrvalExp(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Exp) }
func MlrvalExpm1(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Expm1) }
func MlrvalFloor(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Floor) }
func MlrvalLog(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Log) }
func MlrvalLog10(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Log10) }
func MlrvalLog1p(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Log1p) }
func MlrvalRound(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Round) }
func MlrvalSin(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Sin) }
func MlrvalSinh(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Sinh) }
func MlrvalSqrt(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Sqrt) }
func MlrvalTan(ma *Mlrval) Mlrval   { return mudispo[ma.mvtype](ma, math.Tan) }
func MlrvalTanh(ma *Mlrval) Mlrval  { return mudispo[ma.mvtype](ma, math.Tanh) }

// TODO: port from C
//func MlrvalInvqnorm(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, math.Invqnorm) }
//func MlrvalMax(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Max) }
//func MlrvalMin(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Min) }
//func MlrvalQnorm(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Qnorm) }
//func MlrvalSgn(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Sgn) }

// ================================================================
// BINARY FUNCTIONS
// ================================================================

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

func dot_s_xx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(ma.String() + mb.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*EMPTY  */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDot(ma, mb *Mlrval) Mlrval {
	return dot_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Addition with auto-overflow from int to float when necessary.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

// Auto-overflows up to float.  Additions & subtractions overflow by at most
// one bit so it suffices to check sign-changes.
func plus_n_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval
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

func plus_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) + mb.floatval)
}
func plus_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval + float64(mb.intval))
}
func plus_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval + mb.floatval)
}

var plus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, plus_n_ii, plus_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, plus_f_fi, plus_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryPlus(ma, mb *Mlrval) Mlrval {
	return plus_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Subtraction with auto-overflow from int to float when necessary.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
func minus_n_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval
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

func minus_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) - mb.floatval)
}
func minus_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval - float64(mb.intval))
}
func minus_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval - mb.floatval)
}

var minus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, minus_n_ii, minus_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, minus_f_fi, minus_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryMinus(ma, mb *Mlrval) Mlrval {
	return minus_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Multiplication with auto-overflow from int to float when necessary.  See
// also http://johnkerl.org/miller/doc/reference.html#Arithmetic.

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

func times_n_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval
	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt64(a * b)
	}
}

func times_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) * mb.floatval)
}
func times_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval * float64(mb.intval))
}
func times_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval * mb.floatval)
}

var times_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, times_n_ii, times_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, times_f_fi, times_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalTimes(ma, mb *Mlrval) Mlrval {
	return times_dispositions[ma.mvtype][mb.mvtype](ma, mb)
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

func divide_n_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval

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

func divide_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) / mb.floatval)
}
func divide_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval / float64(mb.intval))
}
func divide_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval / mb.floatval)
}

var divide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, divide_n_ii, divide_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, divide_f_fi, divide_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDivide(ma, mb *Mlrval) Mlrval {
	return divide_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Integer division: DSL operator '//' as in Python.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func int_divide_n_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval

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

func int_divide_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(ma.intval) / mb.floatval))
}
func int_divide_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(ma.floatval / float64(mb.intval)))
}
func int_divide_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(ma.floatval / mb.floatval))
}

var int_divide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, int_divide_n_ii, int_divide_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, int_divide_f_fi, int_divide_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalIntDivide(ma, mb *Mlrval) Mlrval {
	return int_divide_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func pow_f_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(float64(ma.intval), float64(mb.intval)))
}
func pow_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(float64(ma.intval), mb.floatval))
}
func pow_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(ma.floatval, float64(mb.intval)))
}
func pow_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(ma.floatval, mb.floatval))
}

var pow_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, pow_f_ii, pow_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, pow_f_fi, pow_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalPow(ma, mb *Mlrval) Mlrval {
	return pow_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Non-auto-overflowing addition: DSL operator '.+'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotplus_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval + mb.intval)
}
func dotplus_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) + mb.floatval)
}
func dotplus_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval + float64(mb.intval))
}
func dotplus_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval + mb.floatval)
}

var dot_plus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR ABSENT  EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, dotplus_i_ii, dotplus_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotplus_f_fi, dotplus_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDotPlus(ma, mb *Mlrval) Mlrval {
	return dot_plus_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// Non-auto-overflowing subtraction: DSL operator '.-'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotminus_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval - mb.intval)
}
func dotminus_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) - mb.floatval)
}
func dotminus_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval - float64(mb.intval))
}
func dotminus_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval - mb.floatval)
}

var dotminus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, dotminus_i_ii, dotminus_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotminus_f_fi, dotminus_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDotMinus(ma, mb *Mlrval) Mlrval {
	return dotminus_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Non-auto-overflowing multiplication: DSL operator '.*'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dottimes_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval * mb.intval)
}
func dottimes_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) * mb.floatval)
}
func dottimes_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval * float64(mb.intval))
}
func dottimes_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval * mb.floatval)
}

var dottimes_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT   BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, dottimes_i_ii, dottimes_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dottimes_f_fi, dottimes_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDotTimes(ma, mb *Mlrval) Mlrval {
	return dottimes_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotdivide_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval / mb.intval)
}
func dotdivide_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(float64(ma.intval) / mb.floatval)
}
func dotdivide_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval / float64(mb.intval))
}
func dotdivide_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(ma.floatval / mb.floatval)
}

var dotdivide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR ABSENT  EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, dotdivide_i_ii, dotdivide_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotdivide_f_fi, dotdivide_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDotDivide(ma, mb *Mlrval) Mlrval {
	return dotdivide_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './/'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func dotidivide_i_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval

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

func dotidivide_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(ma.intval) / mb.floatval))
}
func dotidivide_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(ma.floatval / float64(mb.intval)))
}
func dotidivide_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Floor(ma.floatval / mb.floatval))
}

var dotidivide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR ABSENT  EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, dotidivide_i_ii, dotidivide_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, dotidivide_f_fi, dotidivide_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalDotIntDivide(ma, mb *Mlrval) Mlrval {
	return dotidivide_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Modulus

func modulus_i_ii(ma, mb *Mlrval) Mlrval {
	a := ma.intval
	b := mb.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	m := a % b
	if a >= 0 {
		if b < 0 {
			m += b
		}
	} else {
		if b >= 0 {
			m += b
		}
	}

	return MlrvalFromInt64(m)
}

func modulus_f_fi(ma, mb *Mlrval) Mlrval {
	a := ma.floatval
	b := float64(mb.intval)
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

func modulus_f_if(ma, mb *Mlrval) Mlrval {
	a := float64(ma.intval)
	b := mb.floatval
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

func modulus_f_ff(ma, mb *Mlrval) Mlrval {
	a := ma.floatval
	b := mb.floatval
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

var modulus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, modulus_i_ii, modulus_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, modulus_f_fi, modulus_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalModulus(ma, mb *Mlrval) Mlrval {
	return modulus_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Bitwise AND

func bitwise_and_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval & mb.intval)
}

var bitwise_and_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_and_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseAND(ma, mb *Mlrval) Mlrval {
	return bitwise_and_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Bitwise OR

func bitwise_or_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval | mb.intval)
}

var bitwise_or_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_or_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseOR(ma, mb *Mlrval) Mlrval {
	return bitwise_or_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Bitwise XOR

func bitwise_xor_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval ^ mb.intval)
}

var bitwise_xor_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, bitwise_xor_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBitwiseXOR(ma, mb *Mlrval) Mlrval {
	return bitwise_xor_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Left shift

func lsh_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval << uint64(mb.intval))
}

var left_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, lsh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalLeftShift(ma, mb *Mlrval) Mlrval {
	return left_shift_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Signed right shift

func srsh_i_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromInt64(ma.intval >> uint64(mb.intval))
}

var signed_right_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, srsh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalSignedRightShift(ma, mb *Mlrval) Mlrval {
	return signed_right_shift_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Unsigned right shift

func ursh_i_ii(ma, mb *Mlrval) Mlrval {
	var ua uint64 = uint64(ma.intval)
	var ub uint64 = uint64(mb.intval)
	var uc = ua >> ub
	return MlrvalFromInt64(int64(uc))
}

var unsigned_right_shift_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _erro, _erro, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, ursh_i_ii, _erro, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _erro, _void, _erro, _erro, _erro, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalUnsignedRightShift(ma, mb *Mlrval) Mlrval {
	return unsigned_right_shift_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
// Boolean expressions for ==, !=, >, >=, <, <=

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ss(ma *Mlrval, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep == mb.printrep)
}
func ne_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep != mb.printrep)
}
func gt_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep > mb.printrep)
}
func ge_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep >= mb.printrep)
}
func lt_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep < mb.printrep)
}
func le_b_ss(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep <= mb.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() == mb.printrep)
}
func ne_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() != mb.printrep)
}
func gt_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() > mb.printrep)
}
func ge_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() >= mb.printrep)
}
func lt_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() < mb.printrep)
}
func le_b_xs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.String() <= mb.printrep)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep == mb.String())
}
func ne_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep != mb.String())
}
func gt_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep > mb.String())
}
func ge_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep >= mb.String())
}
func lt_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep < mb.String())
}
func le_b_sx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.printrep <= mb.String())
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval == mb.intval)
}
func ne_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval != mb.intval)
}
func gt_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval > mb.intval)
}
func ge_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval >= mb.intval)
}
func lt_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval < mb.intval)
}
func le_b_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.intval <= mb.intval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) == mb.floatval)
}
func ne_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) != mb.floatval)
}
func gt_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) > mb.floatval)
}
func ge_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) >= mb.floatval)
}
func lt_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) < mb.floatval)
}
func le_b_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(float64(ma.intval) <= mb.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval == float64(mb.intval))
}
func ne_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != float64(mb.intval))
}
func gt_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval > float64(mb.intval))
}
func ge_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval >= float64(mb.intval))
}
func lt_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval < float64(mb.intval))
}
func le_b_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval <= float64(mb.intval))
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
func eq_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval == mb.floatval)
}
func ne_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval != mb.floatval)
}
func gt_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval > mb.floatval)
}
func ge_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval >= mb.floatval)
}
func lt_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval < mb.floatval)
}
func le_b_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromBool(ma.floatval <= mb.floatval)
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
var eq_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, eq_b_ss, eq_b_ss, eq_b_sx, eq_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_ii, eq_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, eq_b_xs, eq_b_xs, eq_b_fi, eq_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var ne_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, ne_b_ss, ne_b_ss, ne_b_sx, ne_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_ii, ne_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, ne_b_xs, ne_b_xs, ne_b_fi, ne_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var gt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, gt_b_ss, gt_b_ss, gt_b_sx, gt_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_ii, gt_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, gt_b_xs, gt_b_xs, gt_b_fi, gt_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var ge_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, ge_b_ss, ge_b_ss, ge_b_sx, ge_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_ii, ge_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, ge_b_xs, ge_b_xs, ge_b_fi, ge_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var lt_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, lt_b_ss, lt_b_ss, lt_b_sx, lt_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_ii, lt_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, lt_b_xs, lt_b_xs, lt_b_fi, lt_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

var le_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*EMPTY  */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _erro, _absn, _absn},
	/*STRING */ {_erro, _absn, le_b_ss, le_b_ss, le_b_sx, le_b_sx, _erro, _absn, _absn},
	/*INT    */ {_erro, _absn, le_b_xs, le_b_xs, le_b_ii, le_b_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _absn, le_b_xs, le_b_xs, le_b_fi, le_b_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalEquals(ma, mb *Mlrval) Mlrval {
	return eq_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalNotEquals(ma, mb *Mlrval) Mlrval {
	return ne_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalGreaterThan(ma, mb *Mlrval) Mlrval {
	return gt_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalGreaterThanOrEquals(ma, mb *Mlrval) Mlrval {
	return ge_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalLessThan(ma, mb *Mlrval) Mlrval {
	return lt_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func MlrvalLessThanOrEquals(ma, mb *Mlrval) Mlrval {
	return le_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ----------------------------------------------------------------
func MlrvalLogicalAND(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval && mb.boolval)
	} else {
		return MlrvalFromError()
	}
}

func MlrvalLogicalOR(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval || mb.boolval)
	} else {
		return MlrvalFromError()
	}
}

func MlrvalLogicalXOR(ma, mb *Mlrval) Mlrval {
	if ma.mvtype == MT_BOOL && mb.mvtype == MT_BOOL {
		return MlrvalFromBool(ma.boolval != mb.boolval)
	} else {
		return MlrvalFromError()
	}
}

// ================================================================
// VARIADIC FUNCTIONS
// ================================================================

// ================================================================
// MIN AND MAX

// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true
// Exceptions for min & max:
// * absent-null always loses
// * empty-null always loses against numbers

// ----------------------------------------------------------------
func min_f_ff(ma, mb *Mlrval) Mlrval {
	var a float64 = ma.floatval
	var b float64 = mb.floatval
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_f_fi(ma, mb *Mlrval) Mlrval {
	var a float64 = ma.floatval
	var b float64 = float64(mb.intval)
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_f_if(ma, mb *Mlrval) Mlrval {
	var a float64 = float64(ma.intval)
	var b float64 = mb.floatval
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_i_ii(ma, mb *Mlrval) Mlrval {
	var a int64 = ma.intval
	var b int64 = mb.intval
	if a < b {
		return *ma
	} else {
		return *mb
	}
}

// min | b=F   b=T
// --- + ----- -----
// a=F | min=a min=a
// a=T | min=b min=b
func min_b_bb(ma, mb *Mlrval) Mlrval {
	if ma.boolval == false {
		return *ma
	} else {
		return *mb
	}
}

func min_s_ss(ma, mb *Mlrval) Mlrval {
	var a string = ma.printrep
	var b string = mb.printrep
	if a < b {
		return *ma
	} else {
		return *mb
	}
}

var min_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//          ERROR   ABSENT EMPTY  STRING  INT  FLOAT   BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _2___, _2___, _2___, _2___, _2___, _absn, _absn},
	/*EMPTY  */ {_erro, _1___, _void, _void, _2___, _2___, _2___, _absn, _absn},
	/*STRING */ {_erro, _1___, _void, min_s_ss, _2___, _2___, _2___, _absn, _absn},
	/*INT    */ {_erro, _1___, _1___, _1___, min_i_ii, min_f_if, _1___, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _1___, _1___, min_f_fi, min_f_ff, _1___, _absn, _absn},
	/*BOOL   */ {_erro, _1___, _1___, _1___, _2___, _2___, min_b_bb, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryMin(ma, mb *Mlrval) Mlrval {
	return (min_dispositions[ma.mvtype][ma.mvtype])(ma, mb)
}

func MlrvalVariadicMin(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalFromVoid()
	} else {
		retval := *mlrvals[0]
		for _, mlrval := range mlrvals[1:] {
			retval = MlrvalBinaryMin(&retval, mlrval)
		}
		return retval
	}
}

// ----------------------------------------------------------------
func max_f_ff(ma, mb *Mlrval) Mlrval {
	var a float64 = ma.floatval
	var b float64 = mb.floatval
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_f_fi(ma, mb *Mlrval) Mlrval {
	var a float64 = ma.floatval
	var b float64 = float64(mb.intval)
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_f_if(ma, mb *Mlrval) Mlrval {
	var a float64 = float64(ma.intval)
	var b float64 = mb.floatval
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_i_ii(ma, mb *Mlrval) Mlrval {
	var a int64 = ma.intval
	var b int64 = mb.intval
	if a > b {
		return *ma
	} else {
		return *mb
	}
}

// max | b=F   b=T
// --- + ----- -----
// a=F | max=a max=b
// a=T | max=a max=b
func max_b_bb(ma, mb *Mlrval) Mlrval {
	if mb.boolval == false {
		return *ma
	} else {
		return *mb
	}
}

func max_s_ss(ma, mb *Mlrval) Mlrval {
	var a string = ma.printrep
	var b string = mb.printrep
	if a > b {
		return *ma
	} else {
		return *mb
	}
}

var max_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT EMPTY  STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _2___, _2___, _2___, _2___, _2___, _absn, _absn},
	/*EMPTY  */ {_erro, _1___, _void, _void, _2___, _2___, _2___, _absn, _absn},
	/*STRING */ {_erro, _1___, _void, max_s_ss, _2___, _2___, _2___, _absn, _absn},
	/*INT    */ {_erro, _1___, _1___, _1___, max_i_ii, max_f_if, _1___, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _1___, _1___, max_f_fi, max_f_ff, _1___, _absn, _absn},
	/*BOOL   */ {_erro, _1___, _1___, _1___, _2___, _2___, max_b_bb, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryMax(ma, mb *Mlrval) Mlrval {
	return (max_dispositions[ma.mvtype][ma.mvtype])(ma, mb)
}

func MlrvalVariadicMax(mlrvals []*Mlrval) Mlrval {
	if len(mlrvals) == 0 {
		return MlrvalFromVoid()
	} else {
		retval := *mlrvals[0]
		for _, mlrval := range mlrvals[1:] {
			retval = MlrvalBinaryMax(&retval, mlrval)
		}
		return retval
	}
}

// ================================================================
// For sorting

// Lexical sort: just stringify everything.
func LexicalAscendingComparator(ma *Mlrval, mb *Mlrval) int {
	sa := ma.String()
	sb := mb.String()
	if sa < sb {
		return -1
	} else if sa > sb {
		return 1
	} else {
		return 0
	}
}
func LexicalDescendingComparator(ma *Mlrval, mb *Mlrval) int {
	return LexicalAscendingComparator(mb, ma)
}

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

func _neg1(ma, mb *Mlrval) int {
	return -1
}
func _zero(ma, mb *Mlrval) int {
	return 0
}
func _pos1(ma, mb *Mlrval) int {
	return 1
}

func _scmp(ma, mb *Mlrval) int {
	if ma.printrep < mb.printrep {
		return -1
	} else if ma.printrep > mb.printrep {
		return 1
	} else {
		return 0
	}
}

func iicmp(ma, mb *Mlrval) int {
	ca := ma.intval
	cb := mb.intval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ifcmp(ma, mb *Mlrval) int {
	ca := float64(ma.intval)
	cb := mb.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ficmp(ma, mb *Mlrval) int {
	ca := ma.floatval
	cb := float64(mb.intval)
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}
func ffcmp(ma, mb *Mlrval) int {
	ca := ma.floatval
	cb := mb.floatval
	if ca < cb {
		return -1
	} else if ca > cb {
		return 1
	} else {
		return 0
	}
}

func bbcmp(ma, mb *Mlrval) int {
	a := ma.boolval
	b := mb.boolval
	if a == false {
		if b == false {
			return 0
		} else {
			return -1
		}
	} else {
		if b == false {
			return 1
		} else {
			return 0
		}
	}
}

// ----------------------------------------------------------------
// Sort rules (same for min, max, and comparator):
// * NUMERICS < BOOL < STRINGS < ERROR < ABSENT
// * error == error (singleton type)
// * absent == absent (singleton type)
// * string compares on strings
// * numeric compares on numbers
// * false < true

var num_cmp_dispositions = [MT_DIM][MT_DIM]ComparatorFunc{
	//       .  ERROR   ABSENT EMPTY  STRING INT    FLOAT  BOOL    ARRAY MAP
	/*ERROR  */ {_zero, _neg1, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*ABSENT */ {_pos1, _zero, _pos1, _pos1, _pos1, _pos1, _pos1, _zero, _zero},
	/*EMPTY  */ {_neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*STRING */ {_neg1, _neg1, _scmp, _scmp, _pos1, _pos1, _pos1, _zero, _zero},
	/*INT    */ {_neg1, _neg1, _neg1, _neg1, iicmp, ifcmp, _neg1, _zero, _zero},
	/*FLOAT  */ {_neg1, _neg1, _neg1, _neg1, ficmp, ffcmp, _neg1, _zero, _zero},
	/*BOOL   */ {_neg1, _neg1, _neg1, _neg1, _pos1, _pos1, bbcmp, _zero, _zero},
	/*ARRAY  */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
	/*MAP    */ {_zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero, _zero},
}

func NumericAscendingComparator(ma *Mlrval, mb *Mlrval) int {
	return num_cmp_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
func NumericDescendingComparator(ma *Mlrval, mb *Mlrval) int {
	return NumericAscendingComparator(mb, ma)
}

// ================================================================
func MlrvalSsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mb.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		strings.Replace(ma.printrep, mb.printrep, mc.printrep, 1),
	)
}

// ================================================================
func MlrvalGsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mb.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := regexp.MustCompile(mb.printrep)
	return MlrvalFromString(
		re.ReplaceAllString(ma.printrep, mc.printrep),
	)
}
