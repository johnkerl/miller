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

// Function-pointer type for binary-operator disposition matrices.
type binaryFunc func(*Mlrval, *Mlrval) Mlrval

// ----------------------------------------------------------------
// The following are frequently used in disposition matrices for various
// operators and are defined here for re-use. The names are VERY short,
// and all the same length, so that the disposition matrices will look
// reasonable rectangular even after gofmt has been run.

// Return error
func _erro(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromError()
}

// Return absent
func _absn(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromAbsent()
}

// Return void
func _void(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromVoid()
}

// Return first argument
func _1___(val1, val2 *Mlrval) Mlrval {
	return *val1
}

// Return second argument
func _2___(val1, val2 *Mlrval) Mlrval {
	return *val2
}

// Return first argument, as string
func _s1__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String())
}

// Return second argument, as string
func _s2__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val2.String())
}

// Return integer zero
func _i0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromInt64(0)
}

// Return float zero
func _f0__(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromFloat64(0.0)
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
// "102".

func dot_s_xx(val1, val2 *Mlrval) Mlrval {
	return MlrvalFromString(val1.String() + val2.String())
}

var dotDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return dotDispositions[val1.mvtype][val2.mvtype](val1, val2)
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

var plusDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return plusDispositions[val1.mvtype][val2.mvtype](val1, val2)
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

var minusDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return minusDispositions[val1.mvtype][val2.mvtype](val1, val2)
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

var timesDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return timesDispositions[val1.mvtype][val2.mvtype](val1, val2)
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

var divideDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return divideDispositions[val1.mvtype][val2.mvtype](val1, val2)
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

var int_divideDispositions = [MT_DIM][MT_DIM]binaryFunc{
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
	return int_divideDispositions[val1.mvtype][val2.mvtype](val1, val2)
}

// ================================================================
// Non-auto-overflowing addition: DSL operator '.+'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

//static mv_t oplus_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a + b);
//}
//static mv_t oplus_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	long long c = a + b;
//	return mv_from_int(c);
//}
//
//static mv_binary_func_t* oplus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT         FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,         _2,         _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,       _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  oplus_n_ii, oplus_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  oplus_f_fi, oplus_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,       _erro,       _erro},
//};
//
//mv_t x_xx_oplus_func(mv_t* pval1, mv_t* pval2) { return (oplus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ================================================================
// Non-auto-overflowing subtraction: DSL operator '.-'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

//static mv_t ominus_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a - b);
//}
//static mv_t ominus_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	long long c = a - b;
//	return mv_from_int(c);
//}
//
//static mv_binary_func_t* ominus_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT          FLOAT        BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,          _2,          _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,        _void,        _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  ominus_n_ii, ominus_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  ominus_f_fi, ominus_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,        _erro,        _erro},
//};
//
//mv_t x_xx_ominus_func(mv_t* pval1, mv_t* pval2) { return (ominus_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
// Non-auto-overflowing multiplication: DSL operator '.*'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

//static mv_t otimes_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a * b);
//}
//static mv_t otimes_n_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	return mv_from_int(a * b);
//}
//
//static mv_binary_func_t* otimes_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT          FLOAT       BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _2,          _2,         _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,        _void,       _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  otimes_n_ii, otimes_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  otimes_f_fi, otimes_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,        _erro,       _erro},
//};
//
//mv_t x_xx_otimes_func(mv_t* pval1, mv_t* pval2) { return (otimes_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

//static mv_t odivide_f_ff(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = pb->u.fltv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_f_fi(mv_t* pa, mv_t* pb) {
//	double a = pa->u.fltv;
//	double b = (double)pb->u.intv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_f_if(mv_t* pa, mv_t* pb) {
//	double a = (double)pa->u.intv;
//	double b = pb->u.fltv;
//	return mv_from_float(a / b);
//}
//static mv_t odivide_i_ii(mv_t* pa, mv_t* pb) {
//	long long a = pa->u.intv;
//	long long b = pb->u.intv;
//	return mv_from_int(a / b);
//}
//
//static mv_binary_func_t* odivide_dispositions[MT_DIM][MT_DIM] = {
//	//         ERROR  ABSENT EMPTY STRING INT           FLOAT         BOOL
//	/*ERROR*/  {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//	/*ABSENT*/ {_erro, _a,    _a,   _erro,  _i0,          _f0,          _erro},
//	/*EMPTY*/  {_erro, _a,    _void, _erro,  _void,         _void,         _erro},
//	/*STRING*/ {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//	/*INT*/    {_erro, _1,    _void, _erro,  odivide_i_ii, odivide_f_if, _erro},
//	/*FLOAT*/  {_erro, _1,    _void, _erro,  odivide_f_fi, odivide_f_ff, _erro},
//	/*BOOL*/   {_erro, _erro,  _erro, _erro,  _erro,         _erro,         _erro},
//};
//
//mv_t x_xx_odivide_func(mv_t* pval1, mv_t* pval2) { return (odivide_dispositions[pval1->type][pval2->type])(pval1,pval2); }

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './/'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

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
