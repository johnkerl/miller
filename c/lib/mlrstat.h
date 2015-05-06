#ifndef MLRSTAT_H
#define MLRSTAT_H

double mlr_get_stddev(unsigned long long count, double sum, double sum2);

double mlr_get_cov(unsigned long long count, double sumx, double sumy, double sumxy);

void mlr_get_cov_matrix(unsigned long long n,
	double sumx, double sumx2, double sumy, double sumy2, double sumxy, double Q[2][2]);

void mlr_get_linear_regression_pca(
	// Inputs:
	double eigenvalue_1,
	double eigenvalue_2,
	double eigenvector_1[2],
	double eigenvector_2[2],
	double x_mean, double y_mean,
	// Outputs, with quality 1 being a tight fit and quality 0 being a loose one.
	double* pm, double* pb, double* pquality);

#endif // MLRSTAT_H
