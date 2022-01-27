package bifs

import (
	"math"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
// Unary plus operator

var upos_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ _1u___,
	/*FLOAT  */ _1u___,
	/*BOOL   */ _erro1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_plus_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return upos_dispositions[input1.Type()](input1)
}

// ================================================================
// Unary minus operator

func uneg_i_i(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(-input1.AcquireIntValue())
}

func uneg_f_f(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(-input1.AcquireFloatValue())
}

var uneg_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ uneg_i_i,
	/*FLOAT  */ uneg_f_f,
	/*BOOL   */ _erro1,
	/*VOID   */ _void1,
	/*STRING */ _erro1,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ _erro1,
	/*ERROR  */ _erro1,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_minus_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return uneg_dispositions[input1.Type()](input1)
}

// ================================================================
// Addition with auto-overflow from int to float when necessary.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

// Auto-overflows up to float.  Additions & subtractions overflow by at most
// one bit so it suffices to check sign-changes.
func plus_n_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()
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
		return mlrval.FromFloat(float64(a) + float64(b))
	} else {
		return mlrval.FromInt(c)
	}
}

func plus_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) + input2.AcquireFloatValue())
}
func plus_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() + float64(input2.AcquireIntValue()))
}
func plus_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() + input2.AcquireFloatValue())
}

var plus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT      BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {plus_n_ii, plus_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {plus_f_fi, plus_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_plus_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return plus_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Subtraction with auto-overflow from int to float when necessary.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

// Adds & subtracts overflow by at most one bit so it suffices to check
// sign-changes.
func minus_n_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()
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
		return mlrval.FromFloat(float64(a) - float64(b))
	} else {
		return mlrval.FromInt(c)
	}
}

func minus_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) - input2.AcquireFloatValue())
}
func minus_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() - float64(input2.AcquireIntValue()))
}
func minus_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() - input2.AcquireFloatValue())
}

var minus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT         FLOAT       BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {minus_n_ii, minus_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {minus_f_fi, minus_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_minus_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return minus_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Multiplication with auto-overflow from int to float when necessary.  See
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

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
// double less than 2**63. (An alternative would be to do all integer multiplies
// using handcrafted multi-word 128-bit arithmetic).

func times_n_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()
	c := float64(a) * float64(b)

	if math.Abs(c) > 9223372036854774784.0 {
		return mlrval.FromFloat(c)
	} else {
		return mlrval.FromInt(a * b)
	}
}

func times_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) * input2.AcquireFloatValue())
}
func times_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() * float64(input2.AcquireIntValue()))
}
func times_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() * input2.AcquireFloatValue())
}

var times_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT         FLOAT       BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {times_n_ii, times_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {times_f_fi, times_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_times(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return times_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Pythonic division.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic
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

func divide_n_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return mlrval.FromFloat(float64(a) / float64(b))
	}

	// Pythonic division, not C division.
	if a%b == 0 {
		return mlrval.FromInt(a / b)
	} else {
		return mlrval.FromFloat(float64(a) / float64(b))
	}
}

func divide_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) / input2.AcquireFloatValue())
}
func divide_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() / float64(input2.AcquireIntValue()))
}
func divide_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() / input2.AcquireFloatValue())
}

var divide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT        BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {divide_n_ii, divide_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {divide_f_fi, divide_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_i0__, _f0__, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return divide_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Integer division: DSL operator '//' as in Python.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func int_divide_n_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return mlrval.FromFloat(float64(a) / float64(b))
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
	return mlrval.FromInt(q)
}

func int_divide_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(float64(input1.AcquireIntValue()) / input2.AcquireFloatValue()))
}
func int_divide_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(input1.AcquireFloatValue() / float64(input2.AcquireIntValue())))
}
func int_divide_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(input1.AcquireFloatValue() / input2.AcquireFloatValue()))
}

var int_divide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT            BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {int_divide_n_ii, int_divide_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {int_divide_f_fi, int_divide_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_int_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return int_divide_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Non-auto-overflowing addition: DSL operator '.+'.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func dotplus_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() + input2.AcquireIntValue())
}
func dotplus_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) + input2.AcquireFloatValue())
}
func dotplus_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() + float64(input2.AcquireIntValue()))
}
func dotplus_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() + input2.AcquireFloatValue())
}

var dot_plus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT           FLOAT         BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotplus_i_ii, dotplus_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {dotplus_f_fi, dotplus_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_dot_plus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dot_plus_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// Non-auto-overflowing subtraction: DSL operator '.-'.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func dotminus_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() - input2.AcquireIntValue())
}
func dotminus_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) - input2.AcquireFloatValue())
}
func dotminus_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() - float64(input2.AcquireIntValue()))
}
func dotminus_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() - input2.AcquireFloatValue())
}

var dotminus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT            FLOAT          BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotminus_i_ii, dotminus_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {dotminus_f_fi, dotminus_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_n2__, _n2__, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _absn},
	/*ABSENT */ {_n2__, _n2__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_dot_minus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dotminus_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Non-auto-overflowing multiplication: DSL operator '.*'.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func dottimes_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() * input2.AcquireIntValue())
}
func dottimes_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) * input2.AcquireFloatValue())
}
func dottimes_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() * float64(input2.AcquireIntValue()))
}
func dottimes_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() * input2.AcquireFloatValue())
}

var dottimes_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT            FLOAT          BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dottimes_i_ii, dottimes_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {dottimes_f_fi, dottimes_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_dot_times(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dottimes_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './'.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func dotdivide_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(input1.AcquireIntValue() / input2.AcquireIntValue())
}
func dotdivide_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(float64(input1.AcquireIntValue()) / input2.AcquireFloatValue())
}
func dotdivide_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() / float64(input2.AcquireIntValue()))
}
func dotdivide_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(input1.AcquireFloatValue() / input2.AcquireFloatValue())
}

var dotdivide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT             FLOAT           BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotdivide_i_ii, dotdivide_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {dotdivide_f_fi, dotdivide_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_dot_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dotdivide_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// 64-bit integer division: DSL operator './/'.  See also
// https://miller.readthedocs.io/en/latest/reference-main-arithmetic

func dotidivide_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return mlrval.FromFloat(float64(a) / float64(b))
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
	return mlrval.FromInt(q)
}

func dotidivide_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(float64(input1.AcquireIntValue()) / input2.AcquireFloatValue()))
}
func dotidivide_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(input1.AcquireFloatValue() / float64(input2.AcquireIntValue())))
}
func dotidivide_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Floor(input1.AcquireFloatValue() / input2.AcquireFloatValue()))
}

var dotidivide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT            BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotidivide_i_ii, dotidivide_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {dotidivide_f_fi, dotidivide_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_2___, _2___, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
}

func BIF_dot_int_divide(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dotidivide_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ----------------------------------------------------------------
// Modulus

func modulus_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireIntValue()
	b := input2.AcquireIntValue()

	if b == 0 {
		// Compute inf/nan as with floats rather than fatal runtime FPE on integer divide by zero
		return mlrval.FromFloat(float64(a) / float64(b))
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

	return mlrval.FromInt(m)
}

func modulus_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireFloatValue()
	b := float64(input2.AcquireIntValue())
	return mlrval.FromFloat(a - b*math.Floor(a/b))
}

func modulus_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := float64(input1.AcquireIntValue())
	b := input2.AcquireFloatValue()
	return mlrval.FromFloat(a - b*math.Floor(a/b))
}

func modulus_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	a := input1.AcquireFloatValue()
	b := input2.AcquireFloatValue()
	return mlrval.FromFloat(a - b*math.Floor(a/b))
}

var modulus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT           FLOAT         BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {modulus_i_ii, modulus_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {modulus_f_fi, modulus_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_modulus(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return modulus_dispositions[input1.Type()][input2.Type()](input1, input2)
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

func imodop(input1, input2, input3 *mlrval.Mlrval, iop i_iii_func) *mlrval.Mlrval {
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
		return mlrval.ERROR
	}
	if !input2.IsInt() {
		return mlrval.ERROR
	}
	if !input3.IsInt() {
		return mlrval.ERROR
	}

	return mlrval.FromInt(iop(input1.AcquireIntValue(), input2.AcquireIntValue(), input3.AcquireIntValue()))
}

func BIF_mod_add(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodadd)
}

func BIF_mod_sub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodsub)
}

func BIF_mod_mul(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodmul)
}

func BIF_mod_exp(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	// Pre-check for negative exponent
	if input2.IsInt() && input2.AcquireIntValue() < 0 {
		return mlrval.ERROR
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
func min_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = input1.AcquireFloatValue()
	var b float64 = input2.AcquireFloatValue()
	return mlrval.FromFloat(math.Min(a, b))
}

func min_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = input1.AcquireFloatValue()
	var b float64 = float64(input2.AcquireIntValue())
	return mlrval.FromFloat(math.Min(a, b))
}

func min_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = float64(input1.AcquireIntValue())
	var b float64 = input2.AcquireFloatValue()
	return mlrval.FromFloat(math.Min(a, b))
}

func min_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a int64 = input1.AcquireIntValue()
	var b int64 = input2.AcquireIntValue()
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
func min_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.AcquireBoolValue() == false {
		return input1
	} else {
		return input2
	}
}

func min_s_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a string = input1.AcquireStringValue()
	var b string = input2.AcquireStringValue()
	if a < b {
		return input1
	} else {
		return input2
	}
}

var min_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL      VOID   STRING    ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {min_i_ii, min_f_if, _1___, _1___, _1___, _absn, _absn, _erro, _erro, _1___, _1___},
	/*FLOAT  */ {min_f_fi, min_f_ff, _1___, _1___, _1___, _absn, _absn, _erro, _erro, _1___, _1___},
	/*BOOL   */ {_2___, _2___, min_b_bb, _1___, _1___, _absn, _absn, _erro, _erro, _1___, _1___},
	/*VOID   */ {_2___, _2___, _2___, _void, _void, _absn, _absn, _erro, _erro, _1___, _1___},
	/*STRING */ {_2___, _2___, _2___, _void, min_s_ss, _absn, _absn, _erro, _erro, _1___, _1___},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro, _erro, _null, _null},
	/*ABSENT */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro, _erro, _null, _absn},
}

func BIF_min_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return (min_dispositions[input1.Type()][input2.Type()])(input1, input2)
}

func BIF_min_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	} else {
		retval := mlrvals[0]
		for i := range mlrvals {
			if i > 0 {
				retval = BIF_min_binary(retval, mlrvals[i])
			}
		}
		return retval
	}
}

// ----------------------------------------------------------------
func max_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = input1.AcquireFloatValue()
	var b float64 = input2.AcquireFloatValue()
	return mlrval.FromFloat(math.Max(a, b))
}

func max_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = input1.AcquireFloatValue()
	var b float64 = float64(input2.AcquireIntValue())
	return mlrval.FromFloat(math.Max(a, b))
}

func max_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a float64 = float64(input1.AcquireIntValue())
	var b float64 = input2.AcquireFloatValue()
	return mlrval.FromFloat(math.Max(a, b))
}

func max_i_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a int64 = input1.AcquireIntValue()
	var b int64 = input2.AcquireIntValue()
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
func max_b_bb(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input2.AcquireBoolValue() == false {
		return input1
	} else {
		return input2
	}
}

func max_s_ss(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	var a string = input1.AcquireStringValue()
	var b string = input2.AcquireStringValue()
	if a > b {
		return input1
	} else {
		return input2
	}
}

var max_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL      VOID   STRING    ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {max_i_ii, max_f_if, _2___, _2___, _2___, _absn, _absn, _erro, _erro, _null, _1___},
	/*FLOAT  */ {max_f_fi, max_f_ff, _2___, _2___, _2___, _absn, _absn, _erro, _erro, _null, _1___},
	/*BOOL   */ {_1___, _1___, max_b_bb, _2___, _2___, _absn, _absn, _erro, _erro, _null, _1___},
	/*VOID   */ {_1___, _1___, _1___, _void, _2___, _absn, _absn, _erro, _erro, _null, _1___},
	/*STRING */ {_1___, _1___, _1___, _1___, max_s_ss, _absn, _absn, _erro, _erro, _null, _1___},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _absn, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _null, _erro},
	/*NULL   */ {_null, _null, _null, _null, _null, _absn, _absn, _erro, _null, _null, _absn},
	/*ABSENT */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_max_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return (max_dispositions[input1.Type()][input2.Type()])(input1, input2)
}

func BIF_max_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	} else {
		retval := mlrvals[0]
		for i := range mlrvals {
			if i > 0 {
				retval = BIF_max_binary(retval, mlrvals[i])
			}
		}
		return retval
	}
}
