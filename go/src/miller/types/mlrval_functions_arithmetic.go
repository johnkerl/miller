package types

import (
	"math"
)

// ================================================================
// Unary plus operator

var upos_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*VOID   */ _void1,
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
	/*VOID   */ _void1,
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR ABSENT  VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT   BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR ABSENT  VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR ABSENT  VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
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

// ================================================================
// Pythonic
func mlrmod(a, m int64) int64 {
	retval := a % m
	if retval < 0 {
		retval += m
	}
	return retval
}

type i_iii_func func(a, b, m int64) int64

func imodadd(a, b, m int64) int64 {
	return mlrmod(a+b, m)
}
func imodsub(a, b, m int64) int64 {
	return mlrmod(a-b, m)
}
func imodmul(a, b, m int64) int64 {
	return mlrmod(a*b, m)
}
func imodexp(a, e, m int64) int64 {
	if e == 0 {
		return 1
	}
	if e == 1 {
		return a
	}

	// Repeated-squaring algorithm.
	// We assume our caller has verified the exponent is not negative.
	apower := a
	c := int64(1)
	u := uint64(e)

	for u != 0 {
		if (u & 1) == 1 {
			c = mlrmod(c*apower, m)
		}
		u >>= 1
		apower = mlrmod(apower*apower, m)
	}
	return c
}

func imodop(ma, mb, mc *Mlrval, iop i_iii_func) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	if !mc.IsLegit() {
		return *mc
	}
	if !ma.IsInt() {
		return MlrvalFromError()
	}
	if !mb.IsInt() {
		return MlrvalFromError()
	}
	if !mc.IsInt() {
		return MlrvalFromError()
	}

	return MlrvalFromInt64(iop(ma.intval, mb.intval, mc.intval))
}

func MlrvalModAdd(ma, mb, mc *Mlrval) Mlrval {
	return imodop(ma, mb, mc, imodadd)
}

func MlrvalModSub(ma, mb, mc *Mlrval) Mlrval {
	return imodop(ma, mb, mc, imodsub)
}

func MlrvalModMul(ma, mb, mc *Mlrval) Mlrval {
	return imodop(ma, mb, mc, imodmul)
}

func MlrvalModExp(ma, mb, mc *Mlrval) Mlrval {
	// Pre-check for negative exponent
	if mb.mvtype == MT_INT && mb.intval < 0 {
		return MlrvalFromError()
	}
	return imodop(ma, mb, mc, imodexp)
}

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
	//       .  ERROR   ABSENT VOID   STRING    INT       FLOAT     BOOL      ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _2___, _2___, _2___, _2___, _2___, _absn, _absn},
	/*VOID   */ {_erro, _1___, _void, _void, _2___, _2___, _2___, _absn, _absn},
	/*STRING */ {_erro, _1___, _void, min_s_ss, _2___, _2___, _2___, _absn, _absn},
	/*INT    */ {_erro, _1___, _1___, _1___, min_i_ii, min_f_if, _1___, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _1___, _1___, min_f_fi, min_f_ff, _1___, _absn, _absn},
	/*BOOL   */ {_erro, _1___, _1___, _1___, _2___, _2___, min_b_bb, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryMin(ma, mb *Mlrval) Mlrval {
	return (min_dispositions[ma.mvtype][mb.mvtype])(ma, mb)
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
	//       .  ERROR   ABSENT VOID   STRING    INT       FLOAT     BOOL      ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _2___, _2___, _2___, _2___, _2___, _absn, _absn},
	/*VOID   */ {_erro, _1___, _void, _2___, _1___, _1___, _1___, _absn, _absn},
	/*STRING */ {_erro, _1___, _1___, max_s_ss, _1___, _1___, _1___, _absn, _absn},
	/*INT    */ {_erro, _1___, _2___, _2___, max_i_ii, max_f_if, _2___, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _2___, _2___, max_f_fi, max_f_ff, _2___, _absn, _absn},
	/*BOOL   */ {_erro, _1___, _2___, _2___, _1___, _1___, max_b_bb, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalBinaryMax(ma, mb *Mlrval) Mlrval {
	return (max_dispositions[ma.mvtype][mb.mvtype])(ma, mb)
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
