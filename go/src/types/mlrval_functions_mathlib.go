// ================================================================
// Go math-library functions
// ================================================================

package types

import (
	"math"

	"mlr/src/lib"
)

// ----------------------------------------------------------------
func math_unary_f_i(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MlrvalPointerFromFloat64(f(float64(input1.intval)))
}
func math_unary_f_f(input1 *Mlrval, f mathLibUnaryFunc) *Mlrval {
	return MlrvalPointerFromFloat64(f(input1.floatval))
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
}

func MlrvalAbs(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Abs) }
func MlrvalAcos(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Acos) }
func MlrvalAcosh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Acosh) }
func MlrvalAsin(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Asin) }
func MlrvalAsinh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Asinh) }
func MlrvalAtan(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Atan) }
func MlrvalAtanh(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Atanh) }
func MlrvalCbrt(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Cbrt) }
func MlrvalCeil(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Ceil) }
func MlrvalCos(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Cos) }
func MlrvalCosh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Cosh) }
func MlrvalErf(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Erf) }
func MlrvalErfc(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Erfc) }
func MlrvalExp(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Exp) }
func MlrvalExpm1(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Expm1) }
func MlrvalFloor(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Floor) }
func MlrvalInvqnorm(input1 *Mlrval) *Mlrval { return mudispo[input1.mvtype](input1, lib.Invqnorm) }
func MlrvalLog(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Log) }
func MlrvalLog10(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Log10) }
func MlrvalLog1p(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Log1p) }
func MlrvalQnorm(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, lib.Qnorm) }
func MlrvalRound(input1 *Mlrval) *Mlrval    { return mudispo[input1.mvtype](input1, math.Round) }
func MlrvalSgn(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, lib.Sgn) }
func MlrvalSin(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Sin) }
func MlrvalSinh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Sinh) }
func MlrvalSqrt(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Sqrt) }
func MlrvalTan(input1 *Mlrval) *Mlrval      { return mudispo[input1.mvtype](input1, math.Tan) }
func MlrvalTanh(input1 *Mlrval) *Mlrval     { return mudispo[input1.mvtype](input1, math.Tanh) }

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func pow_f_ii(input1, input2 *Mlrval) *Mlrval {
	foutput := math.Pow(float64(input1.intval), float64(input2.intval))
	ioutput := int(foutput)
	// Int raised to int power should be float if it can be (i.e. unless overflow)
	if float64(ioutput) == foutput {
		return MlrvalPointerFromInt(ioutput)
	} else {
		return MlrvalPointerFromFloat64(foutput)
	}
}
func pow_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Pow(float64(input1.intval), input2.floatval))
}
func pow_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Pow(input1.floatval, float64(input2.intval)))
}
func pow_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Pow(input1.floatval, input2.floatval))
}

var pow_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT       FLOAT     BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, pow_f_ii, pow_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, pow_f_fi, pow_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalPow(input1, input2 *Mlrval) *Mlrval {
	return pow_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func atan2_f_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Atan2(float64(input1.intval), float64(input2.intval)))
}
func atan2_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Atan2(float64(input1.intval), input2.floatval))
}
func atan2_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Atan2(input1.floatval, float64(input2.intval)))
}
func atan2_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(math.Atan2(input1.floatval, input2.floatval))
}

var atan2_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT         FLOAT       BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, atan2_f_ii, atan2_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, atan2_f_fi, atan2_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalAtan2(input1, input2 *Mlrval) *Mlrval {
	return atan2_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(mlr_roundm(float64(input1.intval), float64(input2.intval)))
}
func roundm_f_if(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(mlr_roundm(float64(input1.intval), input2.floatval))
}
func roundm_f_fi(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(mlr_roundm(input1.floatval, float64(input2.intval)))
}
func roundm_f_ff(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromFloat64(mlr_roundm(input1.floatval, input2.floatval))
}

var roundm_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING INT          FLOAT        BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _erro, _void, _erro, roundm_f_ii, roundm_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _erro, _void, _erro, roundm_f_fi, roundm_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _erro, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalRoundm(input1, input2 *Mlrval) *Mlrval {
	return roundm_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
func MlrvalLogifit(input1, input2, input3 *Mlrval) *Mlrval {
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

	return MlrvalPointerFromFloat64(1.0 / (1.0 + math.Exp(-m*x-b)))
}
