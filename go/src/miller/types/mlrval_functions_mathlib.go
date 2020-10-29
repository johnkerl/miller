package types

import (
	"fmt"
	"math"
	"os"
)

// ================================================================
// Go math-library functions

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
func math_unary_f_i(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromFloat64(f(float64(ma.intval)))
}
func math_unary_f_f(ma *Mlrval, f mathLibUnaryFunc) Mlrval {
	return MlrvalFromFloat64(f(ma.floatval))
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

func MlrvalAbs(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Abs) }
func MlrvalAcos(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Acos) }
func MlrvalAcosh(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Acosh) }
func MlrvalAsin(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Asin) }
func MlrvalAsinh(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Asinh) }
func MlrvalAtan(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Atan) }
func MlrvalAtanh(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Atanh) }
func MlrvalCbrt(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Cbrt) }
func MlrvalCeil(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Ceil) }
func MlrvalCos(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Cos) }
func MlrvalCosh(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Cosh) }
func MlrvalErf(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Erf) }
func MlrvalErfc(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Erfc) }
func MlrvalExp(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Exp) }
func MlrvalExpm1(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Expm1) }
func MlrvalFloor(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Floor) }
func MlrvalInvqnorm(ma *Mlrval) Mlrval { return mudispo[ma.mvtype](ma, mlrInvqnorm) }
func MlrvalLog(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Log) }
func MlrvalLog10(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Log10) }
func MlrvalLog1p(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Log1p) }
func MlrvalQnorm(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, mlrQnorm) }
func MlrvalRound(ma *Mlrval) Mlrval    { return mudispo[ma.mvtype](ma, math.Round) }
func MlrvalSgn(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, mlrSgn) }
func MlrvalSin(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Sin) }
func MlrvalSinh(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Sinh) }
func MlrvalSqrt(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Sqrt) }
func MlrvalTan(ma *Mlrval) Mlrval      { return mudispo[ma.mvtype](ma, math.Tan) }
func MlrvalTanh(ma *Mlrval) Mlrval     { return mudispo[ma.mvtype](ma, math.Tanh) }

// ================================================================
// Exponentiation: DSL operator '**'.  See also
// http://johnkerl.org/miller/doc/reference.html#Arithmetic.

func pow_f_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(float64(ma.intval), float64(mb.intval)))
}
func pow_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(float64(ma.intval), mb.floatval))
}
func pow_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(ma.floatval, float64(mb.intval)))
}
func pow_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Pow(ma.floatval, mb.floatval))
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

func MlrvalPow(ma, mb *Mlrval) Mlrval {
	return pow_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
func atan2_f_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Atan2(float64(ma.intval), float64(mb.intval)))
}
func atan2_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Atan2(float64(ma.intval), mb.floatval))
}
func atan2_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Atan2(ma.floatval, float64(mb.intval)))
}
func atan2_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(math.Atan2(ma.floatval, mb.floatval))
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

func MlrvalAtan2(ma, mb *Mlrval) Mlrval {
	return atan2_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
func mlr_roundm(x, m float64) float64 {
	return math.Round(x/m) * m
}

func roundm_f_ii(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(mlr_roundm(float64(ma.intval), float64(mb.intval)))
}
func roundm_f_if(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(mlr_roundm(float64(ma.intval), mb.floatval))
}
func roundm_f_fi(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(mlr_roundm(ma.floatval, float64(mb.intval)))
}
func roundm_f_ff(ma, mb *Mlrval) Mlrval {
	return MlrvalFromFloat64(mlr_roundm(ma.floatval, mb.floatval))
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

func MlrvalRoundm(ma, mb *Mlrval) Mlrval {
	return roundm_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
func MlrvalLogifit(ma, mb, mc *Mlrval) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	if !mc.IsLegit() {
		return *mc
	}

	// int/float OK; rest not
	x, xok := ma.GetFloatValue()
	if !xok {
		return MlrvalFromError()
	}
	m, mok := mb.GetFloatValue()
	if !mok {
		return MlrvalFromError()
	}
	b, bok := mc.GetFloatValue()
	if !bok {
		return MlrvalFromError()
	}

	return MlrvalFromFloat64(1.0 / (1.0 + math.Exp(-m*x-b)))
}
