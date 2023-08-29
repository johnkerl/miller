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
func _math_unary_erro1(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.FromTypeErrorUnary(fname, input1)
}

// Return absent (unary math-library func)
func _math_unary_absn1(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.ABSENT
}

// Return null (unary math-library func)
func _math_unary_null1(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.NULL
}

// Return void (unary math-library func)
func _math_unary_void1(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.VOID
}

// ----------------------------------------------------------------
func math_unary_f_i(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.FromFloat(f(float64(input1.AcquireIntValue())))
}
func math_unary_i_i(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
	return mlrval.FromInt(int64(f(float64(input1.AcquireIntValue()))))
}
func math_unary_f_f(input1 *mlrval.Mlrval, f mathLibUnaryFunc, fname string) *mlrval.Mlrval {
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

func BIF_acos(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Acos, "acos")
}
func BIF_acosh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Acosh, "acosh")
}
func BIF_asin(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Asin, "asin")
}
func BIF_asinh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Asinh, "asinh")
}
func BIF_atan(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Atan, "atan")
}
func BIF_atanh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Atanh, "atanh")
}
func BIF_cbrt(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Cbrt, "atan")
}
func BIF_cos(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Cos, "cos")
}
func BIF_cosh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Cosh, "cosh")
}
func BIF_erf(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Erf, "erf")
}
func BIF_erfc(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Erfc, "erfc")
}
func BIF_exp(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Exp, "exp")
}
func BIF_expm1(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Expm1, "expm1")
}
func BIF_invqnorm(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, lib.Invqnorm, "invqnorm")
}
func BIF_log(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Log, "log")
}
func BIF_log10(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Log10, "log10")
}
func BIF_log1p(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Log1p, "log1p")
}
func BIF_qnorm(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, lib.Qnorm, "qnorm")
}
func BIF_sin(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Sin, "sin")
}
func BIF_sinh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Sinh, "sinh")
}
func BIF_sqrt(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Sqrt, "sqrt")
}
func BIF_tan(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Tan, "tan")
}
func BIF_tanh(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mudispo[input1.Type()](input1, math.Tanh, "tanh")
}

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
func BIF_abs(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Abs, "abs")
} // xxx
func BIF_ceil(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Ceil, "ceil")
} // xxx
func BIF_floor(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Floor, "floor")
} // xxx
func BIF_round(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, math.Round, "round")
} // xxx
func BIF_sgn(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return imudispo[input1.Type()](input1, lib.Sgn, "sgn")
} // xxx

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// https://johnkerl.org/miller6/reference-main-arithmetic.html

func pow_f_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	foutput := math.Pow(float64(input1.AcquireIntValue()), float64(input2.AcquireIntValue()))
	ioutput := int64(foutput)
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

func powte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("**", input1, input2)
}

var pow_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT        FLOAT     BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {pow_f_ii, pow_f_if, powte, _void, powte, powte, powte, powte, powte, powte, _1___},
	/*FLOAT  */ {pow_f_fi, pow_f_ff, powte, _void, powte, powte, powte, powte, powte, powte, _1___},
	/*BOOL   */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*VOID   */ {_void, _void, powte, _void, powte, powte, powte, powte, powte, powte, _absn},
	/*STRING */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*ARRAY  */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*MAP    */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*FUNC   */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*ERROR  */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*NULL   */ {powte, powte, powte, powte, powte, powte, powte, powte, powte, powte, _absn},
	/*ABSENT */ {_i0__, _f0__, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
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

func atan2te(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("atan2", input1, input2)
}

var atan2_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT          FLOAT       BOOL     VOID     STRING   ARRAY    MAP      FUNC     ERROR    NULL     ABSENT
	/*INT    */ {atan2_f_ii, atan2_f_if, atan2te, _void, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _1___},
	/*FLOAT  */ {atan2_f_fi, atan2_f_ff, atan2te, _void, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _1___},
	/*BOOL   */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*VOID   */ {_void, _void, atan2te, _void, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*STRING */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*ARRAY  */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*MAP    */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*FUNC   */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*ERROR  */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*NULL   */ {atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, atan2te, _absn},
	/*ABSENT */ {_i0__, _f0__, atan2te, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func BIF_atan2(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return atan2_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromInt(int64(mlr_roundm(float64(input1.AcquireIntValue()), float64(input2.AcquireIntValue()))))
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

func rdmte(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("roundm", input1, input2)
}

var roundm_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT           FLOAT        BOOL   VOID   STRING ARRAY  MAP    FUNC   ERROR  NULL   ABSENT
	/*INT    */ {roundm_f_ii, roundm_f_if, rdmte, _void, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _1___},
	/*FLOAT  */ {roundm_f_fi, roundm_f_ff, rdmte, _void, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _1___},
	/*BOOL   */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*VOID   */ {_void, _void, rdmte, _void, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*STRING */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*ARRAY  */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*MAP    */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*FUNC   */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*ERROR  */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*NULL   */ {rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, rdmte, _absn},
	/*ABSENT */ {_i0__, _f0__, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func BIF_roundm(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return roundm_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
func logifit_te(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromTypeErrorBinary("logifit", input1, input2)
}

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
		return logifit_te(input1, input2)
	}
	m, mok := input2.GetNumericToFloatValue()
	if !mok {
		return logifit_te(input1, input2)
	}
	b, bok := input3.GetNumericToFloatValue()
	if !bok {
		return logifit_te(input1, input2)
	}

	return mlrval.FromFloat(1.0 / (1.0 + math.Exp(-m*x-b)))
}
