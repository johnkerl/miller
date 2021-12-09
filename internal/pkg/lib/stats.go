// ================================================================
// These are intended for streaming (i.e. single-pass) applications. Otherwise
// the formulas look different (and are more intuitive).
// ================================================================

package lib

import (
	"math"
)

// ----------------------------------------------------------------
// Univariate linear regression
// ----------------------------------------------------------------
// There are N (xi, yi) pairs.
//
// minimize E = sum (yi - m xi - b)^2
//
// Set the two partial derivatives to zero and solve for m and b:
//
// DE/Dm = sum 2 (yi - m xi - b) (-xi) = 0
// DE/Db = sum 2 (yi - m xi - b) (-1)  = 0
//
// sum (yi - m xi - b) (xi) = 0
// sum (yi - m xi - b)      = 0
//
// sum (xi yi - m xi^2 - b xi) = 0
// sum (yi - m xi - b)         = 0
//
// m sum(xi^2) + b sum(xi) = sum(xi yi)
// m sum(xi)   + b N       = sum(yi)
//
// [ sum(xi^2)   sum(xi) ] [ m ] = [ sum(xi yi) ]
// [ sum(xi)     N       ] [ b ] = [ sum(yi)    ]
//
// [ m ] = [ sum(xi^2) sum(xi) ]^-1  [ sum(xi yi) ]
// [ b ]   [ sum(xi)   N       ]     [ sum(yi)    ]
//
//       = [ N         -sum(xi)  ]  [ sum(xi yi) ] * 1/D
//         [ -sum(xi)   sum(xi^2)]  [ sum(yi)    ]
//
// where
//
//   D = N sum(xi^2) - sum(xi)^2.
//
// So
//
//      N sum(xi yi) - sum(xi) sum(yi)
// m = --------------------------------
//                   D
//
//      -sum(xi)sum(xi yi) + sum(xi^2) sum(yi)
// b = ----------------------------------------
//                   D
//
// ----------------------------------------------------------------

func GetLinearRegressionOLS(
	nint int,
	sumx float64,
	sumx2 float64,
	sumxy float64,
	sumy float64,
) (m, b float64) {

	n := float64(nint)
	D := n*sumx2 - sumx*sumx
	m = (n*sumxy - sumx*sumy) / D
	b = (-sumx*sumxy + sumx2*sumy) / D
	return m, b
}

// We would need a second pass through the data to compute the error-bars given
// the data and the m and the b.
//
//	# Young 1962, pp. 122-124.  Compute sample variance of linear
//	# approximations, then variances of m and b.
//	var_z = 0.0
//	for i in range(0, N):
//		var_z += (m * xs[i] + b - ys[i])**2
//	var_z /= N
//
//	var_m = (N * var_z) / D
//	var_b = (var_z * sumx2) / D
//
//	output = [m, b, math.sqrt(var_m), math.sqrt(var_b)]

// ----------------------------------------------------------------

// GetVar is the finalizing function for computing variance from streamed
// accumulator values.
func GetVar(
	nint int,
	sumx float64,
	sumx2 float64,
) float64 {

	n := float64(nint)
	mean := sumx / n
	numerator := sumx2 - mean*(2.0*sumx-n*mean)
	if numerator < 0.0 { // round-off error
		numerator = 0.0
	}
	denominator := n - 1.0
	return numerator / denominator
}

// ----------------------------------------------------------------
// Unbiased estimator:
//    (1/n)   sum{(xi-mean)**3}
//  -----------------------------
// [(1/(n-1)) sum{(xi-mean)**2}]**1.5

// mean = sumx / n; n mean = sumx

// sum{(xi-mean)^3}
//   = sum{xi^3 - 3 mean xi^2 + 3 mean^2 xi - mean^3}
//   = sum{xi^3} - 3 mean sum{xi^2} + 3 mean^2 sum{xi} - n mean^3
//   = sumx3 - 3 mean sumx2 + 3 mean^2 sumx - n mean^3
//   = sumx3 - 3 mean sumx2 + 3n mean^3 - n mean^3
//   = sumx3 - 3 mean sumx2 + 2n mean^3
//   = sumx3 - mean*(3 sumx2 + 2n mean^2)

// sum{(xi-mean)^2}
//   = sum{xi^2 - 2 mean xi + mean^2}
//   = sum{xi^2} - 2 mean sum{xi} + n mean^2
//   = sumx2 - 2 mean sumx + n mean^2
//   = sumx2 - 2 n mean^2 + n mean^2
//   = sumx2 - n mean^2

// ----------------------------------------------------------------

// GetSkewness is the finalizing function for computing skewness from streamed
// accumulator values.
func GetSkewness(
	nint int,
	sumx float64,
	sumx2 float64,
	sumx3 float64,
) float64 {

	n := float64(nint)
	mean := sumx / n
	numerator := sumx3 - mean*(3*sumx2-2*n*mean*mean)
	numerator = numerator / n
	denominator := (sumx2 - n*mean*mean) / (n - 1)
	denominator = math.Pow(denominator, 1.5)
	return numerator / denominator
}

// ----------------------------------------------------------------
// Unbiased:
//  (1/n) sum{(x-mean)**4}
//  ----------------------- - 3
// [(1/n) sum{(x-mean)**2}]**2

// sum{(xi-mean)^4}
//   = sum{xi^4 - 4 mean xi^3 + 6 mean^2 xi^2 - 4 mean^3 xi + mean^4}
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 4 mean^3 sum{xi} + n mean^4
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 4 n mean^4 + n mean^4
//   = sum{xi^4} - 4 mean sum{xi^3} + 6 mean^2 sum{xi^2} - 3 n mean^4
//   = sum{xi^4} - mean*(4 sum{xi^3} - 6 mean sum{xi^2} + 3 n mean^3)
//   = sumx4 - mean*(4 sumx3 - 6 mean sumx2 + 3 n mean^3)
//   = sumx4 - mean*(4 sumx3 - mean*(6 sumx2 - 3 n mean^2))

func GetKurtosis(
	nint int,
	sumx float64,
	sumx2 float64,
	sumx3 float64,
	sumx4 float64,
) float64 {

	n := float64(nint)
	mean := sumx / n
	numerator := sumx4 - mean*(4*sumx3-mean*(6*sumx2-3*n*mean*mean))
	numerator = numerator / n
	denominator := (sumx2 - n*mean*mean) / n
	denominator = denominator * denominator
	return numerator/denominator - 3.0
}

// ----------------------------------------------------------------
// Non-streaming implementation:
//
// def find_sample_covariance(xs, ys):
//      n = len(xs)
//      mean_x = find_mean(xs)
//      mean_y = find_mean(ys)
//
//      sum = 0.0
//      for k in range(0, n):
//              sum += (xs[k] - mean_x) * (ys[k] - mean_y)
//
//      return sum / (n-1.0)

func GetCov(
	nint int,
	sumx float64,
	sumy float64,
	sumxy float64,
) float64 {

	n := float64(nint)
	meanx := sumx / n
	meany := sumy / n
	numerator := sumxy - meanx*sumy - meany*sumx + n*meanx*meany
	denominator := n - 1
	return numerator / denominator
}

// ----------------------------------------------------------------
func GetCovMatrix(
	nint int,
	sumx float64,
	sumx2 float64,
	sumy float64,
	sumy2 float64,
	sumxy float64,
) (Q [2][2]float64) {

	n := float64(nint)
	denominator := n - 1

	Q[0][0] = (sumx2 - sumx*sumx/n) / denominator
	Q[0][1] = (sumxy - sumx*sumy/n) / denominator
	Q[1][0] = Q[0][1]
	Q[1][1] = (sumy2 - sumy*sumy/n) / denominator

	return Q
}

// ----------------------------------------------------------------
// Principal component analysis can be used for linear regression:
//
// * Compute the covariance matrix for the x's and y's.
//
// * Find its eigenvalues and eigenvectors of the cov. (This is real-symmetric
//   so Jacobi iteration is simple and fine.)
//
// * The principal eigenvector points in the direction of the fit.
//
// * The covariance matrix is computed on zero-mean data so the intercept
//   is zero. The fit equation is of the form (y - nu) = m*(x - mu) where mu
//   and nu are x and y means, respectively.
//
// * If the fit is perfect then the 2nd eigenvalue will be zero; if the fit is
//   good then the 2nd eigenvalue will be smaller; if the fit is bad then
//   they'll be about the same. I use 1 - |lambda2|/|lambda1| as an indication
//   of quality of the fit.
//
// Standard ("ordinary least-squares") linear regression is appropriate when
// the errors are thought to be all in the y's. PCA ("total least-squares") is
// appropriate when the x's and the y's are thought to both have errors.

func GetLinearRegressionPCA(
	eigenvalue_1 float64,
	eigenvalue_2 float64,
	eigenvector_1 [2]float64,
	eigenvector_2 [2]float64,
	x_mean float64,
	y_mean float64,
) (m, b, quality float64) {

	abs_1 := math.Abs(eigenvalue_1)
	abs_2 := math.Abs(eigenvalue_2)
	quality = 1.0
	if abs_1 == 0.0 {
		quality = 0.0
	} else if abs_2 > 0.0 {
		quality = 1.0 - abs_2/abs_1
	}
	a0 := eigenvector_1[0]
	a1 := eigenvector_1[1]
	m = a1 / a0
	b = y_mean - m*x_mean
	return m, b, quality
}
