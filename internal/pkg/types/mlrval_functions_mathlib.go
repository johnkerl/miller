// ================================================================
// Go math-library functions
// ================================================================

package types

import (
	"math"

	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
func math_unary_f_i(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MlrvalFromFloat64(f(float64(input1.intval)))
}
func math_unary_i_i(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MlrvalFromInt(int(f(float64(input1.intval))))
}
func math_unary_f_f(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MlrvalFromFloat64(f(input1.floatval))
}

// Disposition vector for unary mathlib functions
var mudispo = [MT_DIM]mathLibUnaryFuncWrapper{
	/*ERROR  */ _math_unary_erro1,
	/*ABSENT */ _math_unary_absn1,
	/*NULL   */ _math_unary_null1,
	/*VOID   */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*INT    */ math_unary_f_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
	/*FUNC   */ _math_unary_erro1,
}

func BIF_acos(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Acos) }
func BIF_acosh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Acosh) }
func BIF_asin(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Asin) }
func BIF_asinh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Asinh) }
func BIF_atan(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Atan) }
func BIF_atanh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Atanh) }
func BIF_cbrt(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Cbrt) }
func BIF_cos(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Cos) }
func BIF_cosh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Cosh) }
func BIF_erf(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Erf) }
func BIF_erfc(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Erfc) }
func BIF_exp(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Exp) }
func BIF_expm1(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Expm1) }
func BIF_invqnorm(input1 *Mlrval) *Mlrval { return mudispo[input1.mvtype](input1, lib.Invqnorm) }
func BIF_log(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Log) }
func BIF_log10(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Log10) }
func BIF_log1p(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Log1p) }
func BIF_qnorm(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, lib.Qnorm) }
func BIF_sin(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Sin) }
func BIF_sinh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Sinh) }
func BIF_sqrt(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Sqrt) }
func BIF_tan(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Tan) }
func BIF_tanh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Tanh) }

// Disposition vector for unary mathlib functions which are int-preserving
var imudispo = [MT_DIM]mathLibUnaryFuncWrapper{
	/*ERROR  */ _math_unary_erro1,
	/*ABSENT */ _math_unary_absn1,
	/*NULL   */ _math_unary_null1,
	/*VOID   */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*INT    */ math_unary_i_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
	/*FUNC   */ _math_unary_erro1,
}

// Int-preserving
func BIF_abs(input1 *Mlrval) *Mlrval   { return imudispo[input1.mvtype](input1, math.Abs) }   // xxx
func BIF_ceil(input1 *Mlrval) *Mlrval  { return imudispo[input1.mvtype](input1, math.Ceil) }  // xxx
func BIF_floor(input1 *Mlrval) *Mlrval { return imudispo[input1.mvtype](input1, math.Floor) } // xxx
func BIF_round(input1 *Mlrval) *Mlrval { return imudispo[input1.mvtype](input1, math.Round) } // xxx
func BIF_sgn(input1 *Mlrval) *Mlrval   { return imudispo[input1.mvtype](input1, lib.Sgn) }    // xxx

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func pow_f_ii(input1, input2 *Mlrval) *Mlrval {
	foutput := math.Pow(float64(input1.intval), float64(input2.intval))
	ioutput := int(foutput)
	// Int raised to int power should be float if it can be (i.e. unless overflow)
	if float64(ioutput) == foutput {
		return MlrvalFromInt(ioutput)
	} else {
		return MlrvalFromFloat64(foutput)
	}
}
func pow_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Pow(float64(input1.intval), input2.floatval))
}
func pow_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Pow(input1.floatval, float64(input2.intval)))
}
func pow_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Pow(input1.floatval, input2.floatval))
}

var pow_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT       FLOAT     BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, pow_f_ii, pow_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, pow_f_fi, pow_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_pow(input1, input2 *Mlrval) *Mlrval {
	return pow_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func atan2_f_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Atan2(float64(input1.intval), float64(input2.intval)))
}
func atan2_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Atan2(float64(input1.intval), input2.floatval))
}
func atan2_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Atan2(input1.floatval, float64(input2.intval)))
}
func atan2_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(math.Atan2(input1.floatval, input2.floatval))
}

var atan2_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT         FLOAT       BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, atan2_f_ii, atan2_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, atan2_f_fi, atan2_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_atan2(input1, input2 *Mlrval) *Mlrval {
	return atan2_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromInt(int(mlr_roundm(float64(input1.intval), float64(input2.intval))))
}
func roundm_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(mlr_roundm(float64(input1.intval), input2.floatval))
}
func roundm_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(mlr_roundm(input1.floatval, float64(input2.intval)))
}
func roundm_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromFloat64(mlr_roundm(input1.floatval, input2.floatval))
}

var roundm_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT          FLOAT        BOOL   ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn, _erro},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, roundm_f_ii, roundm_f_if, _erro, _absn, _absn, _erro},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, roundm_f_fi, roundm_f_ff, _erro, _absn, _absn, _erro},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_roundm(input1, input2 *Mlrval) *Mlrval {
	return roundm_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func BIF_logifit(input1, input2, input3 *Mlrval) *Mlrval {
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
		return MLRVAL_ERROR
	}
	m, mok := input2.GetNumericToFloatValue()
	if !mok {
		return MLRVAL_ERROR
	}
	b, bok := input3.GetNumericToFloatValue()
	if !bok {
		return MLRVAL_ERROR
	}

	return MlrvalFromFloat64(1.0 / (1.0 + math.Exp(-m*x-b)))
}
