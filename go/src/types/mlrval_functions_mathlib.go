// ================================================================
// Go math-library functions
// ================================================================

package types

import (
	"fmt"
	"math"
	"os"
)

// ----------------------------------------------------------------
// Some wrappers around things which aren't one-liners from math.*.

func mlrSgn(a float64) float64 {
	if a > 0 {
		return 1.0
	} else if a < 0 {
		return -1.0
	} else if a == 0 {
		return 0.0
	} else {
		return math.NaN()
	}
}

// Normal cumulative distribution function, expressed in terms of erfc library
// function (which is awkward, but exists).
func mlrQnorm(x float64) float64 {
	return 0.5 * math.Erfc(-x/math.Sqrt2)
}

// This is a tangent-following method not unlike Newton-Raphson:
// * We can compute qnorm(y) = integral from -infinity to y of (1/sqrt(2pi)) exp(-t^2/2) dt.
// * We can compute derivative of qnorm(y) = (1/sqrt(2pi)) exp(-y^2/2).
// * We cannot explicitly compute invqnorm(y).
// * If dx/dy = (1/sqrt(2pi)) exp(-y^2/2) then dy/dx = sqrt(2pi) exp(y^2/2).
//
// This means we *can* compute the derivative of invqnorm even though we
// can't compute the function itself. So the essence of the method is to
// follow the tangent line to form successive approximations: we have known function input x
// and unknown function output y and initial guess y0.  At each step we find the intersection
// of the tangent line at y_n with the vertical line at x, to find y_{n+1}. Specificall:
//
// * Even though we can't compute y = q^-1(x) we can compute x = q(y).
// * Start with initial guess for y (y0 = 0.0 or y0 = x both are OK).
// * Find x = q(y). Since q (and therefore q^-1) are 1-1, we're done if qnorm(invqnorm(x)) is small.
// * Else iterate: using point-slope form, (y_{n+1} - y_n) / (x_{n+1} - x_n) = m = sqrt(2pi) exp(y_n^2/2).
//   Here x_2 = x (the input) and x_1 = q(y_1).
// * Solve for y_{n+1} and repeat.

const INVQNORM_TOL float64 = 1e-9
const INVQNORM_MAXITER int = 30

func mlrInvqnorm(x float64) float64 {
	// Initial approximation is linear. Starting with y0 = 0.0 works just as well.
	y0 := x - 0.5
	if x <= 0.0 {
		return 0.0
	}
	if x >= 1.0 {
		return 0.0
	}

	y := y0
	niter := 0

	for {

		backx := mlrQnorm(y)
		err := math.Abs(x - backx)
		if err < INVQNORM_TOL {
			break
		}
		if niter > INVQNORM_MAXITER {
			fmt.Fprintf(os.Stderr,
				"Miller: internal coding error: max iterations %d exceeded in invqnorm.\n",
				INVQNORM_MAXITER,
			)
			os.Exit(1)
		}
		m := math.Sqrt2 * math.SqrtPi * math.Exp(y*y/2.0)
		delta_y := m * (x - backx)
		y += delta_y
		niter++
	}

	return y
}

// ----------------------------------------------------------------
func math_unary_f_i(output, input1 *Mlrval, f mathLibUnaryFunc) {
	output.SetFromFloat64(f(float64(input1.intval)))
}
func math_unary_f_f(output, input1 *Mlrval, f mathLibUnaryFunc) {
	output.SetFromFloat64(f(input1.floatval))
}

// Disposition vector for unary mathlib functions
var mudispo = [MT_DIM]mathLibUnaryFuncWrapper{
	/*ERROR  */ _math_unary_erro1,
	/*ABSENT */ _math_unary_absn1,
	/*VOID   */ _math_unary_void1,
	/*STRING */ _math_unary_erro1,
	/*INT    */ math_unary_f_i,
	/*FLOAT  */ math_unary_f_f,
	/*BOOL   */ _math_unary_erro1,
	/*ARRAY  */ _math_unary_absn1,
	/*MAP    */ _math_unary_absn1,
}

func MlrvalAbs(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Abs) }
func MlrvalAcos(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Acos) }
func MlrvalAcosh(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Acosh) }
func MlrvalAsin(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Asin) }
func MlrvalAsinh(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Asinh) }
func MlrvalAtan(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Atan) }
func MlrvalAtanh(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Atanh) }
func MlrvalCbrt(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Cbrt) }
func MlrvalCeil(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Ceil) }
func MlrvalCos(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Cos) }
func MlrvalCosh(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Cosh) }
func MlrvalErf(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Erf) }
func MlrvalErfc(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Erfc) }
func MlrvalExp(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Exp) }
func MlrvalExpm1(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Expm1) }
func MlrvalFloor(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Floor) }
func MlrvalInvqnorm(output, input1 *Mlrval) {
	mudispo[input1.mvtype](output, input1, mlrInvqnorm)
}
func MlrvalLog(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Log) }
func MlrvalLog10(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Log10) }
func MlrvalLog1p(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Log1p) }
func MlrvalQnorm(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, mlrQnorm) }
func MlrvalRound(output, input1 *Mlrval) { mudispo[input1.mvtype](output, input1, math.Round) }
func MlrvalSgn(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, mlrSgn) }
func MlrvalSin(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Sin) }
func MlrvalSinh(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Sinh) }
func MlrvalSqrt(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Sqrt) }
func MlrvalTan(output, input1 *Mlrval)   { mudispo[input1.mvtype](output, input1, math.Tan) }
func MlrvalTanh(output, input1 *Mlrval)  { mudispo[input1.mvtype](output, input1, math.Tanh) }

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func pow_f_ii(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Pow(float64(input1.intval), float64(input2.intval)))
}
func pow_f_if(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Pow(float64(input1.intval), input2.floatval))
}
func pow_f_fi(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Pow(input1.floatval, float64(input2.intval)))
}
func pow_f_ff(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Pow(input1.floatval, input2.floatval))
}

var pow_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, pow_f_ii, pow_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, pow_f_fi, pow_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalPow(output, input1, input2 *Mlrval) {
	pow_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ================================================================
func atan2_f_ii(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Atan2(float64(input1.intval), float64(input2.intval)))
}
func atan2_f_if(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Atan2(float64(input1.intval), input2.floatval))
}
func atan2_f_fi(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Atan2(input1.floatval, float64(input2.intval)))
}
func atan2_f_ff(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(math.Atan2(input1.floatval, input2.floatval))
}

var atan2_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, atan2_f_ii, atan2_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, atan2_f_fi, atan2_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalAtan2(output, input1, input2 *Mlrval) {
	atan2_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(mlr_roundm(float64(input1.intval), float64(input2.intval)))
}
func roundm_f_if(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(mlr_roundm(float64(input1.intval), input2.floatval))
}
func roundm_f_fi(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(mlr_roundm(input1.floatval, float64(input2.intval)))
}
func roundm_f_ff(output, input1, input2 *Mlrval) {
	output.SetFromFloat64(mlr_roundm(input1.floatval, input2.floatval))
}

var roundm_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//           ERROR  ABSENT VOID   STRING INT    FLOAT  BOOL ARRAY MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _absn, _erro, _i0__, _f0__, _erro, _absn, _absn},
	/*VOID   */ {_erro, _absn, _void, _erro, _void, _void, _erro, _absn, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*INT    */ {_erro, _1___, _void, _erro, roundm_f_ii, roundm_f_if, _erro, _absn, _absn},
	/*FLOAT  */ {_erro, _1___, _void, _erro, roundm_f_fi, roundm_f_ff, _erro, _absn, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ARRAY  */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
	/*MAP    */ {_absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn, _absn},
}

func MlrvalRoundm(output, input1, input2 *Mlrval) {
	roundm_dispositions[input1.mvtype][input2.mvtype](output, input1, input2)
}

// ================================================================
func MlrvalLogifit(output, input1, input2, input3 *Mlrval) {
	if !input1.IsLegit() {
		output.CopyFrom(input1)
		return
	}
	if !input2.IsLegit() {
		output.CopyFrom(input2)
		return
	}
	if !input3.IsLegit() {
		output.CopyFrom(input3)
		return
	}

	// int/float OK; rest not
	x, xok := input1.GetNumericToFloatValue()
	if !xok {
		output.SetFromError()
		return
	}
	m, mok := input2.GetNumericToFloatValue()
	if !mok {
		output.SetFromError()
		return
	}
	b, bok := input3.GetNumericToFloatValue()
	if !bok {
		output.SetFromError()
		return
	}

	output.SetFromFloat64(1.0 / (1.0 + math.Exp(-m*x-b)))
}
