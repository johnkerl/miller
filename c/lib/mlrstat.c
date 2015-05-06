#include <math.h>
#include "lib/mlrstat.h"

// xxx cmt intended for streaming applications. otherwise the formulas are
// different (and more intuitive).
double mlr_get_stddev(unsigned long long n, double sum, double sum2) {
	double mean = sum / n;
	double numerator = sum2 - 2.0*mean*sum + n*mean*mean;
	if (numerator < 0.0) // round-off error
		numerator = 0.0;
	double denominator = n - 1LL;
	return sqrt(numerator / denominator);
}

double mlr_get_cov(unsigned long long n, double sumx, double sumy, double sumxy) {
	double meanx = sumx / n;
	double meany = sumy / n;
	double numerator = sumxy - meanx*sumy - meany*sumx + n*meanx*meany;
	double denominator = n - 1;
	return numerator / denominator;
}

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
// * Compute the covariance matrix for the x's and y's.
// * Find its eigenvalues and eigenvectors of the cov. (This is real-symmetric
//   so Jacobi iteration is simple and fine.)
// * The principal eigenvector points in the direction of the fit.
// * The covariance matrix is computed on zero-mean data so the intercept
//   is zero, of the form (y - nu) = m*(x - mu) where mu and nu are x and y
//   means, respectively.
// * If the fit is perfect then the 2nd eigenvalue will be zero; if the fit is
//   good then the 2nd eigenvalue will be smaller; if the fit is bad then
//   they'll be about the same. I use 1 minus ratio of absolute values
//   of 2nd to 1st eigenvalues as an indication of quality of the fit.
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
