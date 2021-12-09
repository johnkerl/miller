package types

import (
	"math"
)

// ================================================================
// Unary plus operator

var upos_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _null1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*INT    */ _1u___,
	/*FLOAT  */ _1u___,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
}

func BIF_plus_unary(input1 *Mlrval) *Mlrval {
	return upos_dispositions[input1.mvtype](input1)
}

// ================================================================
// Unary minus operator

func uneg_i_i(input1 *Mlrval) *Mlrval {
	return MlrvalFromInt(-input1.intval)
}

func uneg_f_f(input1 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(-input1.floatval)
}

var uneg_dispositions = [MT_DIM]UnaryFunc{
	/*ERROR  */ _erro1,
	/*ABSENT */ _absn1,
	/*NULL   */ _null1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*INT    */ uneg_i_i,
	/*FLOAT  */ uneg_f_f,
	/*BOOL   */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
}

func BIF_minus_unary(input1 *Mlrval) *Mlrval {
	return uneg_dispositions[input1.mvtype](input1)
}

// ================================================================
// Logical NOT operator

func BIF_logicalnot(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_BOOL {
		return MlrvalFromBool(!input1.boolval)
	} else {
		return MLRVAL_ERROR
	}
}

// ================================================================
// Addition with auto-overflow from int to float when necessary.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

// Auto-overflows up to float.  Additions & subtractions overflow by at most
// one bit so it suffices to check sign-changes.
func plus_n_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval
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
		return MlrvalFromInt(c)
	}
}

func plus_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) + input2.floatval)
}
func plus_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval + float64(input2.intval))
}
func plus_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval + input2.floatval)
}

var plus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT        FLOAT      BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _null, _erro, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, plus_n_ii, plus_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, plus_f_fi, plus_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_plus_binary(input1, input2 *Mlrval) *Mlrval {
	return plus_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Subtraction with auto-overflow from int to float when necessary.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
func minus_n_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval
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
		return MlrvalFromInt(c)
	}
}

func minus_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) - input2.floatval)
}
func minus_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval - float64(input2.intval))
}
func minus_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval - input2.floatval)
}

var minus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT         FLOAT       BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _null, _erro, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, minus_n_ii, minus_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, minus_f_fi, minus_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_minus_binary(input1, input2 *Mlrval) *Mlrval {
	return minus_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Multiplication with auto-overflow from int to float when necessary.  See
// https://johnkerl.org/miller6/reference-main-arithmetic.html

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

func times_n_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval
	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt(a * b)
	}
}

func times_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) * input2.floatval)
}
func times_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval * float64(input2.intval))
}
func times_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval * input2.floatval)
}

var times_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT         FLOAT       BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _null, _erro, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, times_n_ii, times_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, times_f_fi, times_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_times(input1, input2 *Mlrval) *Mlrval {
	return times_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Pythonic division.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html
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

func divide_n_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	if a%b == 0 {
		return MlrvalFromInt(a / b)
	} else {
		return MlrvalFromFloat64(float64(a) / float64(b))
	}

	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return MlrvalFromFloat64(c)
	} else {
		return MlrvalFromInt(a * b)
	}
}

func divide_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) / input2.floatval)
}
func divide_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval / float64(input2.intval))
}
func divide_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval / input2.floatval)
}

var divide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT          FLOAT        BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, divide_n_ii, divide_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, divide_f_fi, divide_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_divide(input1, input2 *Mlrval) *Mlrval {
	return divide_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Integer division: DSL operator '//' as in Python.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func int_divide_n_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval

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
	return MlrvalFromInt(q)
}

func int_divide_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(input1.intval) / input2.floatval))
}
func int_divide_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(input1.floatval / float64(input2.intval)))
}
func int_divide_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(input1.floatval / input2.floatval))
}

var int_divide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT              FLOAT            BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, int_divide_n_ii, int_divide_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, int_divide_f_fi, int_divide_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_int_divide(input1, input2 *Mlrval) *Mlrval {
	return int_divide_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Non-auto-overflowing addition: DSL operator '.+'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func dotplus_i_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromInt(input1.intval + input2.intval)
}
func dotplus_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) + input2.floatval)
}
func dotplus_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval + float64(input2.intval))
}
func dotplus_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval + input2.floatval)
}

var dot_plus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT           FLOAT         BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _null, _erro, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, dotplus_i_ii, dotplus_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, dotplus_f_fi, dotplus_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_dot_plus(input1, input2 *Mlrval) *Mlrval {
	return dot_plus_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Non-auto-overflowing subtraction: DSL operator '.-'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func dotminus_i_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromInt(input1.intval - input2.intval)
}
func dotminus_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) - input2.floatval)
}
func dotminus_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval - float64(input2.intval))
}
func dotminus_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval - input2.floatval)
}

var dotminus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT            FLOAT          BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _n2__, _n2__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _null, _erro, _erro, _n2__, _n2__, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, dotminus_i_ii, dotminus_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, dotminus_f_fi, dotminus_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_dot_minus(input1, input2 *Mlrval) *Mlrval {
	return dotminus_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ----------------------------------------------------------------
// Non-auto-overflowing multiplication: DSL operator '.*'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func dottimes_i_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromInt(input1.intval * input2.intval)
}
func dottimes_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) * input2.floatval)
}
func dottimes_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval * float64(input2.intval))
}
func dottimes_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval * input2.floatval)
}

var dottimes_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT            FLOAT          BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _void, _erro, dottimes_i_ii, dottimes_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _void, _erro, dottimes_f_fi, dottimes_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_dot_times(input1, input2 *Mlrval) *Mlrval {
	return dottimes_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func dotdivide_i_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromInt(input1.intval / input2.intval)
}
func dotdivide_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(float64(input1.intval) / input2.floatval)
}
func dotdivide_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval / float64(input2.intval))
}
func dotdivide_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(input1.floatval / input2.floatval)
}

var dotdivide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT             FLOAT           BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, dotdivide_i_ii, dotdivide_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, dotdivide_f_fi, dotdivide_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_dot_divide(input1, input2 *Mlrval) *Mlrval {
	return dotdivide_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './/'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func dotidivide_i_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval

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
	return MlrvalFromInt(q)
}

func dotidivide_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(float64(input1.intval) / input2.floatval))
}
func dotidivide_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(input1.floatval / float64(input2.intval)))
}
func dotidivide_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Floor(input1.floatval / input2.floatval))
}

var dotidivide_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT              FLOAT            BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _erro, _absn, _erro, _2___, _2___, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, dotidivide_i_ii, dotidivide_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, dotidivide_f_fi, dotidivide_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalDotIntDivide(input1, input2 *Mlrval) *Mlrval {
	return dotidivide_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ----------------------------------------------------------------
// Modulus

func modulus_i_ii(input1, input2 *Mlrval) *Mlrval {
	a := input1.intval
	b := input2.intval

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

	return MlrvalFromInt(m)
}

func modulus_f_fi(input1, input2 *Mlrval) *Mlrval {
	a := input1.floatval
	b := float64(input2.intval)
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

func modulus_f_if(input1, input2 *Mlrval) *Mlrval {
	a := float64(input1.intval)
	b := input2.floatval
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

func modulus_f_ff(input1, input2 *Mlrval) *Mlrval {
	a := input1.floatval
	b := input2.floatval
	return MlrvalFromFloat64(a - b*math.Floor(a/b))
}

var modulus_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT           FLOAT         BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, modulus_i_ii, modulus_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, modulus_f_fi, modulus_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_modulus(input1, input2 *Mlrval) *Mlrval {
	return modulus_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// Pythonic
func mlrmod(a, m int) int {
	retval := a % m
	if retval < 0 {
		retval += m
	}
	return retval
}

type i_iii_func func(a, b, m int) int

func imodadd(a, b, m int) int {
	return mlrmod(a+b, m)
}
func imodsub(a, b, m int) int {
	return mlrmod(a-b, m)
}
func imodmul(a, b, m int) int {
	return mlrmod(a*b, m)
}
func imodexp(a, e, m int) int {
	if e == 0 {
		return 1
	}
	if e == 1 {
		return a
	}

	// Repeated-squaring algorithm.
	// We assume our caller has verified the exponent is not negative.
	apower := a
	c := int(1)
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

func imodop(input1, input2, input3 *Mlrval, iop i_iii_func) *Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input3.IsLegit() {
		return input3
	}
	if !input1.IsInt() {
		return MLRVAL_ERROR
	}
	if !input2.IsInt() {
		return MLRVAL_ERROR
	}
	if !input3.IsInt() {
		return MLRVAL_ERROR
	}

	return MlrvalFromInt(iop(input1.intval, input2.intval, input3.intval))
}

func BIF_mod_add(input1, input2, input3 *Mlrval) *Mlrval {
	return imodop(input1, input2, input3, imodadd)
}

func BIF_mod_sub(input1, input2, input3 *Mlrval) *Mlrval {
	return imodop(input1, input2, input3, imodsub)
}

func BIF_mod_mul(input1, input2, input3 *Mlrval) *Mlrval {
	return imodop(input1, input2, input3, imodmul)
}

func BIF_mod_exp(input1, input2, input3 *Mlrval) *Mlrval {
	// Pre-check for negative exponent
	if input2.mvtype == MT_INT && input2.intval < 0 {
		return MLRVAL_ERROR
	}
	return imodop(input1, input2, input3, imodexp)
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
func min_f_ff(input1, input2 *Mlrval) *Mlrval {
	var a float64 = input1.floatval
	var b float64 = input2.floatval
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_f_fi(input1, input2 *Mlrval) *Mlrval {
	var a float64 = input1.floatval
	var b float64 = float64(input2.intval)
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_f_if(input1, input2 *Mlrval) *Mlrval {
	var a float64 = float64(input1.intval)
	var b float64 = input2.floatval
	return MlrvalFromFloat64(math.Min(a, b))
}

func min_i_ii(input1, input2 *Mlrval) *Mlrval {
	var a int = input1.intval
	var b int = input2.intval
	if a < b {
		return input1
	} else {
		return input2
	}
}

// min | b=F   b=T
// --- + ----- -----
// a=F | min=a min=a
// a=T | min=b min=b
func min_b_bb(input1, input2 *Mlrval) *Mlrval {
	if input1.boolval == false {
		return input1
	} else {
		return input2
	}
}

func min_s_ss(input1, input2 *Mlrval) *Mlrval {
	var a string = input1.printrep
	var b string = input2.printrep
	if a < b {
		return input1
	} else {
		return input2
	}
}

var min_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING    INT       FLOAT     BOOL      ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _null, _2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro},
	/*NULL   */ {_erro, _null, _null, _2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro},
	/*VOID   */ {_erro, _1___, _1___, _void, _void, _2___, _2___, _2___, _absn, _absn, _erro},
	/*STRING */ {_erro, _1___, _1___, _void, min_s_ss, _2___, _2___, _2___, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _1___, _1___, _1___, min_i_ii, min_f_if, _1___, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _1___, _1___, _1___, min_f_fi, min_f_ff, _1___, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _1___, _1___, _1___, _1___, _2___, _2___, min_b_bb, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalBinaryMin(input1, input2 *Mlrval) *Mlrval {
	return (min_dispositions[input1.mvtype][input2.mvtype])(input1, input2)
}

func BIF_min_variadic(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) == 0 {
		return MLRVAL_VOID
	} else {
		retval := mlrvals[0]
		for i := range mlrvals {
			if i > 0 {
				retval = MlrvalBinaryMin(retval, mlrvals[i])
			}
		}
		return retval
	}
}

// ----------------------------------------------------------------
func max_f_ff(input1, input2 *Mlrval) *Mlrval {
	var a float64 = input1.floatval
	var b float64 = input2.floatval
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_f_fi(input1, input2 *Mlrval) *Mlrval {
	var a float64 = input1.floatval
	var b float64 = float64(input2.intval)
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_f_if(input1, input2 *Mlrval) *Mlrval {
	var a float64 = float64(input1.intval)
	var b float64 = input2.floatval
	return MlrvalFromFloat64(math.Max(a, b))
}

func max_i_ii(input1, input2 *Mlrval) *Mlrval {
	var a int = input1.intval
	var b int = input2.intval
	if a > b {
		return input1
	} else {
		return input2
	}
}

// max | b=F   b=T
// --- + ----- -----
// a=F | max=a max=b
// a=T | max=a max=b
func max_b_bb(input1, input2 *Mlrval) *Mlrval {
	if input2.boolval == false {
		return input1
	} else {
		return input2
	}
}

func max_s_ss(input1, input2 *Mlrval) *Mlrval {
	var a string = input1.printrep
	var b string = input2.printrep
	if a > b {
		return input1
	} else {
		return input2
	}
}

var max_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING    INT       FLOAT     BOOL      ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _null, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro},
	/*NULL   */ {_null, _absn, _null, _null, _null, _null, _null, _null, _absn, _absn, _erro},
	/*VOID   */ {_erro, _1___, _null, _void, _2___, _1___, _1___, _1___, _absn, _absn, _erro},
	/*STRING */ {_erro, _1___, _null, _1___, max_s_ss, _1___, _1___, _1___, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _null, _2___, _2___, max_i_ii, max_f_if, _2___, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _null, _2___, _2___, max_f_fi, max_f_ff, _2___, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _1___, _null, _2___, _2___, _1___, _1___, max_b_bb, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalBinaryMax(input1, input2 *Mlrval) *Mlrval {
	return (max_dispositions[input1.mvtype][input2.mvtype])(input1, input2)
}

func BIF_max_variadic(mlrvals []*Mlrval) *Mlrval {
	if len(mlrvals) == 0 {
		return MLRVAL_VOID
	} else {
		retval := mlrvals[0]
		for i := range mlrvals {
			if i > 0 {
				retval = MlrvalBinaryMax(retval, mlrvals[i])
			}
		}
		return retval
	}
}
