package bifs

import (
	"math"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
// Unary plus operator

func upos_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_unary("+", input1)
}

var upos_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ _1u___,
	/*FLOAT  */ _1u___,
	/*BOOL   */ upos_te,
	/*VOID   */ _zero1,
	/*STRING */ upos_te,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ upos_te,
	/*ERROR  */ upos_te,
	/*NULL   */ _null1,
	/*ABSENT */ _absn1,
}

func BIF_plus_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return upos_dispositions[input1.Type()](input1)
}

// ================================================================
// Unary minus operator

func uneg_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_unary("-", input1)
}

func uneg_i_i(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(-input1.AcquireIntValue())
}

func uneg_f_f(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(-input1.AcquireFloatValue())
}

var uneg_dispositions = [mlrval.MT_DIM]UnaryFunc{
	/*INT    */ uneg_i_i,
	/*FLOAT  */ uneg_f_f,
	/*BOOL   */ uneg_te,
	/*VOID   */ _zero1,
	/*STRING */ uneg_te,
	/*ARRAY  */ _absn1,
	/*MAP    */ _absn1,
	/*FUNC   */ uneg_te,
	/*ERROR  */ uneg_te,
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

func plste(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("+", input1, input2)
}

var plus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT         FLOAT      BOOL     VOID     STRING   ARRAY    MAP      FUNC     ERROR    NULL     ABSENT
	/*INT    */ {plus_n_ii, plus_f_if, plste, _1___, plste, _absn, _absn, plste, plste, _1___, _1___},
	/*FLOAT  */ {plus_f_fi, plus_f_ff, plste, _1___, plste, _absn, _absn, plste, plste, _1___, _1___},
	/*BOOL   */ {plste, plste, plste, plste, plste, _absn, _absn, plste, plste, plste, plste},
	/*VOID   */ {_2___, _2___, plste, _void, plste, _absn, _absn, plste, plste, plste, _absn},
	/*STRING */ {plste, plste, plste, plste, plste, _absn, _absn, plste, plste, plste, plste},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, plste, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, plste, _absn, _absn, _absn},
	/*FUNC   */ {plste, plste, plste, plste, plste, plste, plste, plste, plste, plste, plste},
	/*ERROR  */ {plste, plste, plste, plste, plste, _absn, _absn, plste, plste, plste, plste},
	/*NULL   */ {_2___, _2___, plste, plste, plste, _absn, _absn, plste, plste, _null, _absn},
	/*ABSENT */ {_2___, _2___, plste, _absn, plste, _absn, _absn, plste, plste, _absn, _absn},
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

func mnste(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("-", input1, input2)
}

var minus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT       BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {minus_n_ii, minus_f_if, mnste, _1___, mnste, _absn, _absn, mnste, mnste, _1___, _1___},
	/*FLOAT  */ {minus_f_fi, minus_f_ff, mnste, _1___, mnste, _absn, _absn, mnste, mnste, _1___, _1___},
	/*BOOL   */ {mnste, mnste, mnste, mnste, mnste, _absn, _absn, mnste, mnste, mnste, mnste},
	/*VOID   */ {_n2__, _n2__, mnste, _void, mnste, _absn, _absn, mnste, mnste, mnste, _absn},
	/*STRING */ {mnste, mnste, mnste, mnste, mnste, _absn, _absn, mnste, mnste, mnste, mnste},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, mnste, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, mnste, _absn, _absn, _absn},
	/*FUNC   */ {mnste, mnste, mnste, mnste, mnste, mnste, mnste, mnste, mnste, mnste, mnste},
	/*ERROR  */ {mnste, mnste, mnste, mnste, mnste, _absn, _absn, mnste, mnste, mnste, mnste},
	/*NULL   */ {_2___, _2___, mnste, mnste, mnste, _absn, _absn, mnste, mnste, _null, _absn},
	/*ABSENT */ {_2___, _2___, mnste, _absn, mnste, _absn, _absn, mnste, mnste, _absn, _absn},
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

func tmste(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("*", input1, input2)
}

var times_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT       BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {times_n_ii, times_f_if, tmste, _1___, tmste, _absn, _absn, tmste, tmste, _1___, _1___},
	/*FLOAT  */ {times_f_fi, times_f_ff, tmste, _1___, tmste, _absn, _absn, tmste, tmste, _1___, _1___},
	/*BOOL   */ {tmste, tmste, tmste, tmste, tmste, _absn, _absn, tmste, tmste, tmste, tmste},
	/*VOID   */ {_2___, _2___, tmste, _void, tmste, _absn, _absn, tmste, tmste, tmste, _absn},
	/*STRING */ {tmste, tmste, tmste, tmste, tmste, _absn, _absn, tmste, tmste, tmste, tmste},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, tmste, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, tmste, _absn, _absn, _absn},
	/*FUNC   */ {tmste, tmste, tmste, tmste, tmste, tmste, tmste, tmste, tmste, tmste, tmste},
	/*ERROR  */ {tmste, tmste, tmste, tmste, tmste, _absn, _absn, tmste, tmste, tmste, tmste},
	/*NULL   */ {_2___, _2___, tmste, tmste, tmste, _absn, _absn, tmste, tmste, _null, _absn},
	/*ABSENT */ {_2___, _2___, tmste, _absn, tmste, _absn, _absn, tmste, tmste, _absn, _absn},
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

func dvdte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("/", input1, input2)
}

var divide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT        BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {divide_n_ii, divide_f_if, dvdte, _void, dvdte, _absn, _absn, dvdte, dvdte, _1___, _1___},
	/*FLOAT  */ {divide_f_fi, divide_f_ff, dvdte, _void, dvdte, _absn, _absn, dvdte, dvdte, _1___, _1___},
	/*BOOL   */ {dvdte, dvdte, dvdte, dvdte, dvdte, _absn, _absn, dvdte, dvdte, dvdte, dvdte},
	/*VOID   */ {_void, _void, dvdte, _void, dvdte, _absn, _absn, dvdte, dvdte, dvdte, _absn},
	/*STRING */ {dvdte, dvdte, dvdte, dvdte, dvdte, _absn, _absn, dvdte, dvdte, dvdte, dvdte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dvdte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dvdte, _absn, _absn, _absn},
	/*FUNC   */ {dvdte, dvdte, dvdte, dvdte, dvdte, dvdte, dvdte, dvdte, dvdte, dvdte, dvdte},
	/*ERROR  */ {dvdte, dvdte, dvdte, dvdte, dvdte, _absn, _absn, dvdte, dvdte, dvdte, dvdte},
	/*NULL   */ {_i0__, _f0__, dvdte, dvdte, dvdte, _absn, _absn, dvdte, dvdte, dvdte, _absn},
	/*ABSENT */ {_i0__, _f0__, dvdte, _absn, dvdte, _absn, _absn, dvdte, dvdte, _absn, _absn},
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

func idvte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("//", input1, input2)
}

var int_divide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT            BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {int_divide_n_ii, int_divide_f_if, idvte, _void, idvte, _absn, _absn, idvte, idvte, idvte, _1___},
	/*FLOAT  */ {int_divide_f_fi, int_divide_f_ff, idvte, _void, idvte, _absn, _absn, idvte, idvte, idvte, _1___},
	/*BOOL   */ {idvte, idvte, idvte, idvte, idvte, _absn, _absn, idvte, idvte, idvte, idvte},
	/*VOID   */ {_void, _void, idvte, _void, idvte, _absn, _absn, idvte, idvte, idvte, _absn},
	/*STRING */ {idvte, idvte, idvte, idvte, idvte, _absn, _absn, idvte, idvte, idvte, idvte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, idvte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, idvte, _absn, _absn, _absn},
	/*FUNC   */ {idvte, idvte, idvte, idvte, idvte, idvte, idvte, idvte, idvte, idvte, idvte},
	/*ERROR  */ {idvte, idvte, idvte, idvte, idvte, _absn, _absn, idvte, idvte, idvte, idvte},
	/*NULL   */ {idvte, idvte, idvte, idvte, idvte, _absn, _absn, idvte, idvte, idvte, _absn},
	/*ABSENT */ {_i0__, _f0__, idvte, _absn, idvte, _absn, _absn, idvte, idvte, _absn, _absn},
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

func dplte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary(".+", input1, input2)
}

var dot_plus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT            FLOAT         BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {dotplus_i_ii, dotplus_f_if, dplte, _1___, dplte, _absn, _absn, dplte, dplte, _1___, _1___},
	/*FLOAT  */ {dotplus_f_fi, dotplus_f_ff, dplte, _1___, dplte, _absn, _absn, dplte, dplte, _1___, _1___},
	/*BOOL   */ {dplte, dplte, dplte, dplte, dplte, _absn, _absn, dplte, dplte, dplte, dplte},
	/*VOID   */ {_2___, _2___, dplte, _void, dplte, _absn, _absn, dplte, dplte, dplte, _absn},
	/*STRING */ {dplte, dplte, dplte, dplte, dplte, _absn, _absn, dplte, dplte, dplte, dplte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dplte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dplte, _absn, _absn, _absn},
	/*FUNC   */ {dplte, dplte, dplte, dplte, dplte, dplte, dplte, dplte, dplte, dplte, dplte},
	/*ERROR  */ {dplte, dplte, dplte, dplte, dplte, _absn, _absn, dplte, dplte, dplte, dplte},
	/*NULL   */ {_2___, _2___, dplte, dplte, dplte, _absn, _absn, dplte, dplte, _null, _absn},
	/*ABSENT */ {_2___, _2___, dplte, _absn, dplte, _absn, _absn, dplte, dplte, _absn, _absn},
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

func dmnte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary(".-", input1, input2)
}

var dotminus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT             FLOAT          BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {dotminus_i_ii, dotminus_f_if, dmnte, _1___, dmnte, _absn, _absn, dmnte, dmnte, _1___, _1___},
	/*FLOAT  */ {dotminus_f_fi, dotminus_f_ff, dmnte, _1___, dmnte, _absn, _absn, dmnte, dmnte, _1___, _1___},
	/*BOOL   */ {dmnte, dmnte, dmnte, dmnte, dmnte, _absn, _absn, dmnte, dmnte, dmnte, dmnte},
	/*VOID   */ {_n2__, _n2__, dmnte, _void, dmnte, _absn, _absn, dmnte, dmnte, dmnte, _absn},
	/*STRING */ {dmnte, dmnte, dmnte, dmnte, dmnte, _absn, _absn, dmnte, dmnte, dmnte, dmnte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dmnte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dmnte, _absn, _absn, _absn},
	/*FUNC   */ {dmnte, dmnte, dmnte, dmnte, dmnte, dmnte, dmnte, dmnte, dmnte, dmnte, dmnte},
	/*ERROR  */ {dmnte, dmnte, dmnte, dmnte, dmnte, _absn, _absn, dmnte, dmnte, dmnte, dmnte},
	/*NULL   */ {_n2__, _n2__, dmnte, dmnte, dmnte, _absn, _absn, dmnte, dmnte, _null, _absn},
	/*ABSENT */ {_n2__, _n2__, dmnte, _absn, dmnte, _absn, _absn, dmnte, dmnte, _absn, _absn},
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

func dttte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary(".*", input1, input2)
}

var dottimes_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT             FLOAT          BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {dottimes_i_ii, dottimes_f_if, dttte, _1___, dttte, _absn, _absn, dttte, dttte, _1___, _1___},
	/*FLOAT  */ {dottimes_f_fi, dottimes_f_ff, dttte, _1___, dttte, _absn, _absn, dttte, dttte, _1___, _1___},
	/*BOOL   */ {dttte, dttte, dttte, dttte, dttte, _absn, _absn, dttte, dttte, dttte, dttte},
	/*VOID   */ {_n2__, _n2__, dttte, _void, dttte, _absn, _absn, dttte, dttte, dttte, _absn},
	/*STRING */ {dttte, dttte, dttte, dttte, dttte, _absn, _absn, dttte, dttte, dttte, dttte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dttte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, dttte, _absn, _absn, _absn},
	/*FUNC   */ {dttte, dttte, dttte, dttte, dttte, dttte, dttte, dttte, dttte, dttte, dttte},
	/*ERROR  */ {dttte, dttte, dttte, dttte, dttte, _absn, _absn, dttte, dttte, dttte, dttte},
	/*NULL   */ {_2___, _2___, dttte, dttte, dttte, _absn, _absn, dttte, dttte, dttte, _absn},
	/*ABSENT */ {_2___, _2___, dttte, _absn, dttte, _absn, _absn, dttte, dttte, _absn, _absn},
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

func ddvte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("./", input1, input2)
}

var dotdivide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT             FLOAT           BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotdivide_i_ii, dotdivide_f_if, ddvte, _void, ddvte, _absn, _absn, ddvte, ddvte, ddvte, _1___},
	/*FLOAT  */ {dotdivide_f_fi, dotdivide_f_ff, ddvte, _void, ddvte, _absn, _absn, ddvte, ddvte, ddvte, _1___},
	/*BOOL   */ {ddvte, ddvte, ddvte, ddvte, ddvte, _absn, _absn, ddvte, ddvte, ddvte, ddvte},
	/*VOID   */ {_void, _void, ddvte, _void, ddvte, _absn, _absn, ddvte, ddvte, ddvte, _absn},
	/*STRING */ {ddvte, ddvte, ddvte, ddvte, ddvte, _absn, _absn, ddvte, ddvte, ddvte, ddvte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, ddvte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, ddvte, _absn, _absn, _absn},
	/*FUNC   */ {ddvte, ddvte, ddvte, ddvte, ddvte, ddvte, ddvte, ddvte, ddvte, ddvte, ddvte},
	/*ERROR  */ {ddvte, ddvte, ddvte, ddvte, ddvte, _absn, _absn, ddvte, ddvte, ddvte, ddvte},
	/*NULL   */ {ddvte, ddvte, ddvte, ddvte, ddvte, _absn, _absn, ddvte, ddvte, ddvte, _absn},
	/*ABSENT */ {_2___, _2___, ddvte, _absn, ddvte, _absn, _absn, ddvte, ddvte, _absn, _absn},
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

func didte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary(".//", input1, input2)
}

var dotidivide_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT              FLOAT            BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dotidivide_i_ii, dotidivide_f_if, didte, _void, didte, _absn, _absn, didte, didte, didte, _1___},
	/*FLOAT  */ {dotidivide_f_fi, dotidivide_f_ff, didte, _void, didte, _absn, _absn, didte, didte, didte, _1___},
	/*BOOL   */ {didte, didte, didte, didte, didte, _absn, _absn, didte, didte, didte, didte},
	/*VOID   */ {_void, _void, didte, _void, didte, _absn, _absn, didte, didte, didte, _absn},
	/*STRING */ {didte, didte, didte, didte, didte, _absn, _absn, didte, didte, didte, didte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, didte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, didte, _absn, _absn, _absn},
	/*FUNC   */ {didte, didte, didte, didte, didte, didte, didte, didte, didte, didte, didte},
	/*ERROR  */ {didte, didte, didte, didte, didte, _absn, _absn, didte, didte, didte, didte},
	/*NULL   */ {didte, didte, didte, didte, didte, _absn, _absn, didte, didte, didte, _absn},
	/*ABSENT */ {_2___, _2___, didte, _absn, didte, _absn, _absn, didte, didte, didte, _absn},
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

func modte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("%", input1, input2)
}

var modulus_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT           FLOAT         BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {modulus_i_ii, modulus_f_if, modte, _void, modte, _absn, _absn, modte, modte, modte, _1___},
	/*FLOAT  */ {modulus_f_fi, modulus_f_ff, modte, _void, modte, _absn, _absn, modte, modte, modte, _1___},
	/*BOOL   */ {modte, modte, modte, modte, modte, _absn, _absn, modte, modte, modte, modte},
	/*VOID   */ {_void, _void, modte, _void, modte, _absn, _absn, modte, modte, modte, _absn},
	/*STRING */ {modte, modte, modte, modte, modte, _absn, _absn, modte, modte, modte, modte},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, modte, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, modte, _absn, _absn, _absn},
	/*FUNC   */ {modte, modte, modte, modte, modte, modte, modte, modte, modte, modte, modte},
	/*ERROR  */ {modte, modte, modte, modte, modte, _absn, _absn, modte, modte, modte, modte},
	/*NULL   */ {modte, modte, modte, modte, modte, _absn, _absn, modte, modte, modte, _absn},
	/*ABSENT */ {_i0__, _f0__, modte, _absn, modte, _absn, _absn, modte, modte, _absn, _absn},
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

func imodop(input1, input2, input3 *mlrval.Mlrval, iop i_iii_func, funcname string) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input3.IsLegit() {
		return input3
	}
	if !input1.IsInt() || !input2.IsInt() || !input3.IsInt() {
		return type_error_ternary(funcname, input1, input2, input3)
	}

	return mlrval.FromInt(
		iop(
			input1.AcquireIntValue(),
			input2.AcquireIntValue(),
			input3.AcquireIntValue(),
		),
	)
}

func BIF_mod_add(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodadd, "madd")
}

func BIF_mod_sub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodsub, "msub")
}

func BIF_mod_mul(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return imodop(input1, input2, input3, imodmul, "mmul")
}

func BIF_mod_exp(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	// Pre-check for negative exponent
	if input2.IsInt() && input2.AcquireIntValue() < 0 {
		return mlrval.ERROR
	}
	return imodop(input1, input2, input3, imodexp, "mexp")
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

func min_te(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("min", input1, input2)
}

var min_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT     BOOL      VOID   STRING    ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {min_i_ii, min_f_if, _1___, _1___, _1___, _absn, _absn, min_te, min_te, _1___, _1___},
	/*FLOAT  */ {min_f_fi, min_f_ff, _1___, _1___, _1___, _absn, _absn, min_te, min_te, _1___, _1___},
	/*BOOL   */ {_2___, _2___, min_b_bb, _1___, _1___, _absn, _absn, min_te, min_te, _1___, _1___},
	/*VOID   */ {_2___, _2___, _2___, _void, _void, _absn, _absn, min_te, min_te, _1___, _1___},
	/*STRING */ {_2___, _2___, _2___, _void, min_s_ss, _absn, _absn, min_te, min_te, _1___, _1___},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, min_te, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, min_te, _absn, _absn, _absn},
	/*FUNC   */ {min_te, min_te, min_te, min_te, min_te, min_te, min_te, min_te, min_te, min_te, min_te},
	/*ERROR  */ {min_te, min_te, min_te, min_te, min_te, _absn, _absn, min_te, min_te, min_te, min_te},
	/*NULL   */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, min_te, min_te, _null, _null},
	/*ABSENT */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, min_te, min_te, _null, _absn},
}

// BIF_min_binary is not a direct DSL function. It's a helper here,
// and is also exposed publicly for use by the stats1 verb.
func BIF_min_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return (min_dispositions[input1.Type()][input2.Type()])(input1, input2)
}

func BIF_min_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	}
	return mlrval.ArrayFold(
		mlrvals,
		bif_min_unary(mlrvals[0]),
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_min_binary(bif_min_unary(a), bif_min_unary(b))
		},
	)
}

func BIF_min_within_map_values(m *mlrval.Mlrmap) *mlrval.Mlrval {
	if m.Head == nil {
		return mlrval.VOID
	}
	return mlrval.MapFold(
		m,
		m.Head.Value,
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_min_binary(a, b)
		},
	)
}

// bif_min_unary allows recursion into arguments, so users can do either
// min(1,2,3) or min([1,2,3]).
func bif_min_unary_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_min_variadic(input1.AcquireArrayValue())
}
func bif_min_unary_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_min_within_map_values(input1.AcquireMapValue())
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var min_unary_dispositions = [mlrval.MT_DIM]UnaryFunc{}

func min_unary_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_unary("min", input1)
}

func init() {
	min_unary_dispositions = [mlrval.MT_DIM]UnaryFunc{
		/*INT    */ _1u___,
		/*FLOAT  */ _1u___,
		/*BOOL   */ _1u___,
		/*VOID   */ _1u___,
		/*STRING */ _1u___,
		/*ARRAY  */ bif_min_unary_array,
		/*MAP    */ bif_min_unary_map,
		/*FUNC   */ min_unary_te,
		/*ERROR  */ min_unary_te,
		/*NULL   */ _null1,
		/*ABSENT */ _absn1,
	}
}

func bif_min_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return min_unary_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func BIF_minlen_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	}
	// Do the bulk arithmetic on native ints not Mlrvals, to avoid unnecessary allocation.
	retval := lib.UTF8Strlen(mlrvals[0].OriginalString())
	for i, _ := range mlrvals {
		clen := lib.UTF8Strlen(mlrvals[i].OriginalString())
		if clen < retval {
			retval = clen
		}
	}
	return mlrval.FromInt(retval)
}

func BIF_minlen_within_map_values(m *mlrval.Mlrmap) *mlrval.Mlrval {
	if m.Head == nil {
		return mlrval.VOID
	}
	// Do the bulk arithmetic on native ints not Mlrvals, to avoid unnecessary allocation.
	retval := lib.UTF8Strlen(m.Head.Value.OriginalString())
	for pe := m.Head.Next; pe != nil; pe = pe.Next {
		clen := lib.UTF8Strlen(pe.Value.OriginalString())
		if clen < retval {
			retval = clen
		}
	}
	return mlrval.FromInt(retval)
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

func max_te(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_binary("max", input1, input2)
}

var max_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL      VOID   STRING    ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {max_i_ii, max_f_if, _2___, _2___, _2___, _absn, _absn, max_te, max_te, _null, _1___},
	/*FLOAT  */ {max_f_fi, max_f_ff, _2___, _2___, _2___, _absn, _absn, max_te, max_te, _null, _1___},
	/*BOOL   */ {_1___, _1___, max_b_bb, _2___, _2___, _absn, _absn, max_te, max_te, _null, _1___},
	/*VOID   */ {_1___, _1___, _1___, _void, _2___, _absn, _absn, max_te, max_te, _null, _1___},
	/*STRING */ {_1___, _1___, _1___, _1___, max_s_ss, _absn, _absn, max_te, max_te, _null, _1___},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, max_te, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, max_te, _absn, _absn, _absn},
	/*FUNC   */ {max_te, max_te, max_te, max_te, max_te, max_te, max_te, max_te, max_te, max_te, max_te},
	/*ERROR  */ {max_te, max_te, max_te, max_te, max_te, _absn, _absn, max_te, max_te, _null, max_te},
	/*NULL   */ {_null, _null, _null, _null, _null, _absn, _absn, max_te, _null, _null, _absn},
	/*ABSENT */ {_2___, _2___, _2___, _2___, _2___, _absn, _absn, max_te, max_te, _absn, _absn},
}

// BIF_max_binary is not a direct DSL function. It's a helper here,
// and is also exposed publicly for use by the stats1 verb.
func BIF_max_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return (max_dispositions[input1.Type()][input2.Type()])(input1, input2)
}

func BIF_max_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	}
	return mlrval.ArrayFold(
		mlrvals,
		bif_max_unary(mlrvals[0]),
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_max_binary(bif_max_unary(a), bif_max_unary(b))
		},
	)
}

func BIF_max_within_map_values(m *mlrval.Mlrmap) *mlrval.Mlrval {
	if m.Head == nil {
		return mlrval.VOID
	}
	return mlrval.MapFold(
		m,
		m.Head.Value,
		func(a, b *mlrval.Mlrval) *mlrval.Mlrval {
			return BIF_max_binary(a, b)
		},
	)
}

// bif_max_unary allows recursion into arguments, so users can do either
// max(1,2,3) or max([1,2,3]).
func bif_max_unary_array(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_max_variadic(input1.AcquireArrayValue())
}
func bif_max_unary_map(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_max_within_map_values(input1.AcquireMapValue())
}

// We get a Golang "initialization loop" due to recursive depth computation
// if this is defined statically. So, we use a "package init" function.
var max_unary_dispositions = [mlrval.MT_DIM]UnaryFunc{}

func max_unary_te(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return type_error_unary("max", input1)
}

func init() {
	max_unary_dispositions = [mlrval.MT_DIM]UnaryFunc{
		/*INT    */ _1u___,
		/*FLOAT  */ _1u___,
		/*BOOL   */ _1u___,
		/*VOID   */ _1u___,
		/*STRING */ _1u___,
		/*ARRAY  */ bif_max_unary_array,
		/*MAP    */ bif_max_unary_map,
		/*FUNC   */ max_unary_te,
		/*ERROR  */ max_unary_te,
		/*NULL   */ _null1,
		/*ABSENT */ _absn1,
	}
}

func bif_max_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return max_unary_dispositions[input1.Type()](input1)
}

// ----------------------------------------------------------------
func BIF_maxlen_variadic(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	}
	// Do the bulk arithmetic on native ints not Mlrvals, to avoid unnecessary allocation.
	retval := lib.UTF8Strlen(mlrvals[0].OriginalString())
	for i, _ := range mlrvals {
		clen := lib.UTF8Strlen(mlrvals[i].OriginalString())
		if clen > retval {
			retval = clen
		}
	}
	return mlrval.FromInt(retval)
}

func BIF_maxlen_within_map_values(m *mlrval.Mlrmap) *mlrval.Mlrval {
	if m.Head == nil {
		return mlrval.VOID
	}
	// Do the bulk arithmetic on native ints not Mlrvals, to avoid unnecessary allocation.
	retval := lib.UTF8Strlen(m.Head.Value.OriginalString())
	for pe := m.Head.Next; pe != nil; pe = pe.Next {
		clen := lib.UTF8Strlen(pe.Value.OriginalString())
		if clen > retval {
			retval = clen
		}
	}
	return mlrval.FromInt(retval)
}
