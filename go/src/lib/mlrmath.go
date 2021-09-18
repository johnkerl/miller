// ================================================================
// Non-mlrval math routines
// ================================================================

package lib

import (
	"fmt"
	"math"
	"os"
)

// ----------------------------------------------------------------
// Some wrappers around things which aren't one-liners from math.*.

func Sgn(a float64) float64 {
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
func Qnorm(x float64) float64 {
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

func Invqnorm(x float64) float64 {
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

		backx := Qnorm(y)
		err := math.Abs(x - backx)
		if err < INVQNORM_TOL {
			break
		}
		if niter > INVQNORM_MAXITER {
			fmt.Fprintf(os.Stderr,
				"mlr: internal coding error: max iterations %d exceeded in invqnorm.\n",
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

const JACOBI_TOLERANCE = 1e-12
const JACOBI_MAXITER = 20

// ----------------------------------------------------------------
// Jacobi real-symmetric eigensolver. Loosely adapted from Numerical Recipes.
//
// Note: this is coded for n=2 (to implement PCA linear regression on 2
// variables) but the algorithm is quite general. Changing from 2 to n is a
// matter of updating the top and bottom of the function: function signature to
// take double** matrix, double* eigenvector_1, double* eigenvector_2, and n;
// create copy-matrix and make-identity matrix functions; free temp matrices at
// the end; etc.

func GetRealSymmetricEigensystem(
	matrix [2][2]float64,
) (
	eigenvalue1 float64, // Output: dominant eigenvalue
	eigenvalue2 float64, // Output: less-dominant eigenvalue
	eigenvector1 [2]float64, // Output: corresponding to dominant eigenvalue
	eigenvector2 [2]float64, // Output: corresponding to less-dominant eigenvalue
) {
	L := [2][2]float64{
		{matrix[0][0], matrix[0][1]},
		{matrix[1][0], matrix[1][1]},
	}
	V := [2][2]float64{
		{1.0, 0.0},
		{0.0, 1.0},
	}
	var P, PT_A [2][2]float64
	n := 2

	found := false
	for iter := 0; iter < JACOBI_MAXITER; iter++ {
		sum := 0.0
		for i := 1; i < n; i++ {
			for j := 0; j < i; j++ {
				sum += math.Abs(L[i][j])
			}
		}
		if math.Abs(sum*sum) < JACOBI_TOLERANCE {
			found = true
			break
		}

		for p := 0; p < n; p++ {
			for q := p + 1; q < n; q++ {
				numer := L[p][p] - L[q][q]
				denom := L[p][q] + L[q][p]
				if math.Abs(denom) < JACOBI_TOLERANCE {
					continue
				}
				theta := numer / denom
				signTheta := 1.0
				if theta < 0 {
					signTheta = -1.0
				}
				t := signTheta / (math.Abs(theta) + math.Sqrt(theta*theta+1))
				c := 1.0 / math.Sqrt(t*t+1)
				s := t * c

				for pi := 0; pi < n; pi++ {
					for pj := 0; pj < n; pj++ {
						if pi == pj {
							P[pi][pj] = 1.0
						} else {
							P[pi][pj] = 0.0
						}
					}
				}
				P[p][p] = c
				P[p][q] = -s
				P[q][p] = s
				P[q][q] = c

				// L = P.transpose() * L * P
				// V = V * P
				matmul2t(&PT_A, &P, &L)
				matmul2(&L, &PT_A, &P)
				matmul2(&V, &V, &P)
			}
		}
	}

	if !found {
		fmt.Fprintf(os.Stderr,
			"%s: Jacobi eigensolver: max iterations (%d) exceeded.  Non-symmetric input?\n",
			"mlr",
			JACOBI_MAXITER,
		)
		os.Exit(1)
	}

	eigenvalue1 = L[0][0]
	eigenvalue2 = L[1][1]
	abs1 := math.Abs(eigenvalue1)
	abs2 := math.Abs(eigenvalue2)
	if abs1 > abs2 {
		eigenvector1[0] = V[0][0] // Column 0 of V
		eigenvector1[1] = V[1][0]
		eigenvector2[0] = V[0][1] // Column 1 of V
		eigenvector2[1] = V[1][1]
	} else {
		eigenvalue1, eigenvalue2 = eigenvalue2, eigenvalue1
		eigenvector1[0] = V[0][1]
		eigenvector1[1] = V[1][1]
		eigenvector2[0] = V[0][0]
		eigenvector2[1] = V[1][0]
	}

	return eigenvalue1, eigenvalue2, eigenvector1, eigenvector2
}

// C = A * B
func matmul2(
	C *[2][2]float64, // Output
	A *[2][2]float64, // Input
	B *[2][2]float64, // Input
) {
	var T [2][2]float64
	n := 2
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += A[i][k] * B[k][j]
			}
			T[i][j] = sum
		}
	}
	// Needs copy in case C's memory is the same as A and/or B
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			C[i][j] = T[i][j]
		}
	}
}

// C = A^t * B
func matmul2t(
	C *[2][2]float64, // Output
	A *[2][2]float64, // Input
	B *[2][2]float64, // Input
) {
	var T [2][2]float64
	n := 2
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for k := 0; k < n; k++ {
				sum += A[k][i] * B[k][j]
			}
			T[i][j] = sum
		}
	}
	// Needs copy in case C's memory is the same as A and/or B
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			C[i][j] = T[i][j]
		}
	}
}

// ================================================================
// Logisitic regression
//
// Real-valued x_0 .. x_{N-1}
// 0/1-valued  y_0 .. y_{N-1}
// Model p(x_i == 1)  as
//   p(x, m, b) = 1 / (1 + exp(-m*x-b)
// which is the same as
//   log(p/(1-p)) = m*x + b
// then
//   p(x, m, b) = 1 / (1 + exp(-m*x-b)
//              = exp(m*x+b) / (1 + exp(m*x+b)
// and
//   1-p        = exp(-m*x-b) / (1 + exp(-m*x-b)
//              = 1 / (1 + exp(m*x+b)
// Note for reference just below that
//   dp/dm      = -1 / [1 + exp(-m*x-b)]**2 * (-x) * exp(-m*x-b)
//              = [x exp(-m*x-b)) ] / [1 + exp(-m*x-b)]**2
//              = x * p * (1-p)
// and
//   dp/db      = -1 / [1 + exp(-m*x-b)]**2 * (-1) * exp(-m*x-b)
//              = [exp(-m*x-b)) ] / [1 + exp(-m*x-b)]**2
//              = p * (1-p)
// Write p_i for p(x_i, m, b)
//
// Maximum-likelihood equation:
//   L(m, b)    = prod_{i=0}^{N-1} [ p_i**y_i * (1-p_i)**(1-y_i) ]
//
// Log-likelihood equation:
//   ell(m, b)  = sum{i=0}^{N-1} [ y_i log(p_i) + (1-y_i) log(1-p_i) ]
//              = sum{i=0}^{N-1} [ log(1-p_i) + y_i log(p_i/(1-p_i)) ]
//              = sum{i=0}^{N-1} [ log(1-p_i) + y_i*(m*x_i+b) ]
// Differentiate with respect to parameters:
//
//   d ell/dm   = sum{i=0}^{N-1} [ -1/(1-p_i) dp_i/dm + x_i*y_i ]
//              = sum{i=0}^{N-1} [ -1/(1-p_i) x_i*p_i*(1-p_i) + x_i*y_i ]
//              = sum{i=0}^{N-1} [ x_i(y_i-p_i) ]
//
//   d ell/db   = sum{i=0}^{N-1} [ -1/(1-p_i) dp_i/db + y_i ]
//              = sum{i=0}^{N-1} [ -1/(1-p_i) p_i*(1-p_i) + y_i ]
//              = sum{i=0}^{N-1} [ y_i - p_i ]
//
//
//   d2ell/dm2  = sum{i=0}^{N-1} [ -x_i dp_i/dm ]
//              = sum{i=0}^{N-1} [ -x_i**2 * p_i * (1-p_i) ]
//
//   d2ell/dmdb = sum{i=0}^{N-1} [ -x_i dp_i/db ]
//              = sum{i=0}^{N-1} [ -x_i * p_i * (1-p_i) ]
//
//   d2ell/dbdm = sum{i=0}^{N-1} [ -dp_i/dm ]
//              = sum{i=0}^{N-1} [ -x_i * p_i * (1-p_i) ]
//
//   d2ell/db2  = sum{i=0}^{N-1} [ -dp_i/db ]
//              = sum{i=0}^{N-1} [ -p_i * (1-p_i) ]
//
// Newton-Raphson to minimize ell(m, b):
// * Pick m0, b0
// * [m_{j+1], b_{j+1}] = H^{-1} grad ell(m_j, b_j)
// * grad ell =
//   [ d ell/dm ]
//   [ d ell/db ]
// * H = Hessian of ell = Jacobian of grad ell =
//   [ d2ell/dm2  d2ell/dmdb ]
//   [ d2ell/dmdb d2ell/db2  ]

// p(x,m,b) for logistic regression:
func lrp(x, m, b float64) float64 {
	return 1.0 / (1.0 + math.Exp(-m*x-b))
}

// 1 - p(x,m,b) for logistic regression:
func lrq(x, m, b float64) float64 {
	return 1.0 / (1.0 + math.Exp(m*x+b))
}

func LogisticRegression(xs, ys []float64) (m, b float64) {
	m0 := -0.001
	b0 := 0.002
	tol := 1e-9
	maxits := 100
	return logisticRegressionAux(xs, ys, m0, b0, tol, maxits)
}

// Supporting routine for mlr_logistic_regression():
func logisticRegressionAux(
	xs, ys []float64,
	m0, b0, tol float64,
	maxits int,
) (m, b float64) {

	InternalCodingErrorIf(len(xs) != len(ys))
	n := len(xs)

	its := 0
	done := false
	m = m0
	b = b0

	for !done {
		// Compute derivatives
		dldm := 0.0
		dldb := 0.0
		d2ldm2 := 0.0
		d2ldmdb := 0.0
		d2ldb2 := 0.0
		ell0 := 0.0

		for i := 0; i < n; i++ {
			xi := xs[i]
			yi := ys[i]
			pi := lrp(xi, m0, b0)
			qi := lrq(xi, m0, b0)
			dldm += xi * (yi - pi)
			dldb += yi - pi
			piqi := pi * qi
			xipiqi := xi * piqi
			xi2piqi := xi * xipiqi
			d2ldm2 -= xi2piqi
			d2ldmdb -= xipiqi
			d2ldb2 -= piqi
			ell0 += math.Log(qi) + yi*(m0*xi+b0)
		}

		// Form the Hessian
		ha := d2ldm2
		hb := d2ldmdb
		hc := d2ldmdb
		hd := d2ldb2

		// Invert the Hessian
		D := ha*hd - hb*hc
		Hinva := hd / D
		Hinvb := -hb / D
		Hinvc := -hc / D
		Hinvd := ha / D

		// Compute H^-1 times grad ell
		Hinvgradm := Hinva*dldm + Hinvb*dldb
		Hinvgradb := Hinvc*dldm + Hinvd*dldb

		// Update [m,b]
		m = m0 - Hinvgradm
		b = b0 - Hinvgradb

		ell := 0.0
		for i := 0; i < n; i++ {
			xi := xs[i]
			yi := ys[i]
			qi := lrq(xi, m, b)
			ell += math.Log(qi) + yi*(m0*xi+b0)
		}

		// Check for convergence
		dell := math.Max(ell, ell0)
		err := 0.0
		if dell != 0.0 {
			err = math.Abs(ell-ell0) / dell
		}

		if err < tol {
			done = true
		}
		its++
		if its > maxits {
			fmt.Fprintf(os.Stderr,
				"mlr_logistic_regression: Newton-Raphson convergence failed after %d iterations. m=%e, b=%e.\n",
				its, m, b)
			os.Exit(1)
		}

		m0 = m
		b0 = b
	}

	return m, b
}
