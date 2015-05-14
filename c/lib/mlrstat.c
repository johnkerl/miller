#include <math.h>
#include "lib/mlrstat.h"

// ================================================================
// These are intended for streaming (i.e. single-pass) applications. Otherwise
// the formulas look different (and are more intuitive).
// ================================================================

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

void mlr_get_linear_regression_ols(unsigned long long n, double sumx, double sumx2, double sumxy, double sumy,
	double* pm, double* pb)
{
	double D =  n * sumx2 - sumx*sumx;
	double m = (n * sumxy - sumx * sumy) / D;
	double b = (-sumx * sumxy + sumx2 * sumy) / D;

	*pm = m;
	*pb = b;
}

// xxx gah ... need a 2nd pass through the data to get the error-bars.
// xxx make a 2nd filter to compute the error-bars given the data & the m & the b?
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
double mlr_get_var(unsigned long long n, double sum, double sum2) {
	double mean = sum / n;
	double numerator = sum2 - 2.0*mean*sum + n*mean*mean;
	if (numerator < 0.0) // round-off error
		numerator = 0.0;
	double denominator = n - 1LL;
	return numerator / denominator;
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

double mlr_get_cov(unsigned long long n, double sumx, double sumy, double sumxy) {
	double meanx = sumx / n;
	double meany = sumy / n;
	double numerator = sumxy - meanx*sumy - meany*sumx + n*meanx*meany;
	double denominator = n - 1;
	return numerator / denominator;
}

// ----------------------------------------------------------------
void mlr_get_cov_matrix(unsigned long long n,
	double sumx, double sumx2, double sumy, double sumy2, double sumxy,
	double Q[2][2])
{
	double denominator = n - 1;
	Q[0][0] = (sumx2 - sumx*sumx/n) / denominator;
	Q[0][1] = (sumxy - sumx*sumy/n) / denominator;
	Q[1][0] = Q[0][1];
	Q[1][1] = (sumy2 - sumy*sumy/n) / denominator;
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

void mlr_get_linear_regression_pca(
	// Inputs:
	double eigenvalue_1,
	double eigenvalue_2,
	double eigenvector_1[2],
	double eigenvector_2[2],
	double x_mean, double y_mean,
	// Outputs:
	double* pm, double* pb, double* pquality)
{
	double abs_1 = fabs(eigenvalue_1);
	double abs_2 = fabs(eigenvalue_2);
	double quality = 1.0;
	if (abs_1 == 0.0)
		quality = 0.0;
	else if (abs_2 > 0.0)
		quality = 1.0 - abs_2 / abs_1;

	double a0 = eigenvector_1[0];
	double a1 = eigenvector_1[1];
	double m = a1 / a0;
	double b = y_mean - m * x_mean;

	*pm = m;
	*pb = b;
	*pquality = quality;
}
