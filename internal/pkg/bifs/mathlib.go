// ================================================================
// Go math-library functions
// ================================================================

package bifs

import (
	"math"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ----------------------------------------------------------------
// Return error (unary math-library func)
func _math_unary_erro1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.ERROR
}

// Return absent (unary math-library func)
func _math_unary_absn1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.ABSENT
}

// Return null (unary math-library func)
func _math_unary_null1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.NULL
}

// Return void (unary math-library func)
func _math_unary_void1(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.VOID
}

// ----------------------------------------------------------------
func math_unary_f_i(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.FromFloat(f(float64(input1.AcquireIntValue())))
}
func math_unary_i_i(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.FromInt(int(f(float64(input1.AcquireIntValue()))))
}
func math_unary_f_f(input1 *mlrval.Mlrval, f mathLibUnaryFunc) *mlrval.Mlrval {
	return mlrval.FromFloat(f(input1.AcquireFloatValue()))
}

// Disposition vector for unary mathlib functions
var mudispo = [mlrval.MT_DIM]mathLibUnaryFuncWrapper{
	/*INT    */ math_unary_f_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*VOID   */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
	/*FUNC   */ _math_unary_erro1,
	/*ERROR  */ _math_unary_erro1,
	/*NULL   */ _math_unary_null1,
	/*ABSENT */ _math_unary_absn1,
}

func BIF_acos(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Acos) }
func BIF_acosh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Acosh)
}
func BIF_asin(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Asin) }
func BIF_asinh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Asinh)
}
func BIF_atan(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Atan) }
func BIF_atanh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Atanh)
}
func BIF_cbrt(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Cbrt) }
func BIF_cos(input1 *mlrval.Mlrval) *mlrval.Mlrval  { return mudispo[input1.Type()](input1, math.Cos) }
func BIF_cosh(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Cosh) }
func BIF_erf(input1 *mlrval.Mlrval) *mlrval.Mlrval  { return mudispo[input1.Type()](input1, math.Erf) }
func BIF_erfc(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Erfc) }
func BIF_exp(input1 *mlrval.Mlrval) *mlrval.Mlrval  { return mudispo[input1.Type()](input1, math.Exp) }
func BIF_expm1(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Expm1)
}
func BIF_invqnorm(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, lib.Invqnorm)
}
func BIF_log(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Log) }
func BIF_log10(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Log10)
}
func BIF_log1p(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Log1p)
}
func BIF_qnorm(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, lib.Qnorm)
}
func BIF_sin(input1 *mlrval.Mlrval) *mlrval.Mlrval  { return mudispo[input1.Type()](input1, math.Sin) }
func BIF_sinh(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Sinh) }
func BIF_sqrt(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Sqrt) }
func BIF_tan(input1 *mlrval.Mlrval) *mlrval.Mlrval  { return mudispo[input1.Type()](input1, math.Tan) }
func BIF_tanh(input1 *mlrval.Mlrval) *mlrval.Mlrval { return mudispo[input1.Type()](input1, math.Tanh) }

// Disposition vector for unary mathlib functions which are int-preserving
var imudispo = [mlrval.MT_DIM]mathLibUnaryFuncWrapper{
	/*INT    */ math_unary_i_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*VOID   */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
	/*FUNC   */ _math_unary_erro1,
	/*ERROR  */ _math_unary_erro1,
	/*NULL   */ _math_unary_null1,
	/*ABSENT */ _math_unary_absn1,
}

// Int-preserving
func BIF_abs(input1 *mlrval.Mlrval) *mlrval.Mlrval { return imudispo[input1.Type()](input1, math.Abs) } // xxx
func BIF_ceil(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Ceil)
} // xxx
func BIF_floor(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Floor)
} // xxx
func BIF_round(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Round)
}                                                  // xxx
func BIF_sgn(input1 *mlrval.Mlrval) *mlrval.Mlrval { return imudispo[input1.Type()](input1, lib.Sgn) } // xxx

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func pow_f_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	foutput := math.Pow(float64(input1.AcquireIntValue()), float64(input2.AcquireIntValue()))
	ioutput := int(foutput)
	// Int raised to int power should be float if it can be (i.e. unless overflow)
	if float64(ioutput) == foutput {
		return mlrval.FromInt(ioutput)
	} else {
		return mlrval.FromFloat(foutput)
	}
}
func pow_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Pow(float64(input1.AcquireIntValue()), input2.AcquireFloatValue()))
}
func pow_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Pow(input1.AcquireFloatValue(), float64(input2.AcquireIntValue())))
}
func pow_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Pow(input1.AcquireFloatValue(), input2.AcquireFloatValue()))
}

var pow_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {pow_f_ii, pow_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {pow_f_fi, pow_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_pow(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return pow_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
func atan2_f_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Atan2(float64(input1.AcquireIntValue()), float64(input2.AcquireIntValue())))
}
func atan2_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Atan2(float64(input1.AcquireIntValue()), input2.AcquireFloatValue()))
}
func atan2_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Atan2(input1.AcquireFloatValue(), float64(input2.AcquireIntValue())))
}
func atan2_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(math.Atan2(input1.AcquireFloatValue(), input2.AcquireFloatValue()))
}

var atan2_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT         FLOAT       BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {atan2_f_ii, atan2_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {atan2_f_fi, atan2_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_atan2(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return atan2_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int(mlr_roundm(float64(input1.AcquireIntValue()), float64(input2.AcquireIntValue()))))
}
func roundm_f_if(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(mlr_roundm(float64(input1.AcquireIntValue()), input2.AcquireFloatValue()))
}
func roundm_f_fi(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(mlr_roundm(input1.AcquireFloatValue(), float64(input2.AcquireIntValue())))
}
func roundm_f_ff(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromFloat(mlr_roundm(input1.AcquireFloatValue(), input2.AcquireFloatValue()))
}

var roundm_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT        BOOL   VOID   STRING ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {roundm_f_ii, roundm_f_if, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*FLOAT  */ {roundm_f_fi, roundm_f_ff, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _1___},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*VOID   */ {_void, _void, _erro, _void, _erro, _absn, _absn, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _absn, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_i0__, _f0__, _erro, _absn, _erro, _absn, _absn, _erro, _erro, _absn, _absn},
}

func BIF_roundm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return roundm_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
func BIF_logifit(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input3.IsLegit() {
		return input3
	}

	// int/float OK; rest not
	x, xok := input1.GetNumericToFloatValue()
	if !xok {
		return mlrval.ERROR
	}
	m, mok := input2.GetNumericToFloatValue()
	if !mok {
		return mlrval.ERROR
	}
	b, bok := input3.GetNumericToFloatValue()
	if !bok {
		return mlrval.ERROR
	}

	return mlrval.FromFloat(1.0 / (1.0 + math.Exp(-m*x-b)))
}
