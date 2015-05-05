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
	double* pq00, double* pq01, double* pq10, double* pq11)
{
	double denominator = n - 1;
	*pq00 = (sumx2 - sumx*sumx/n) / denominator;
	*pq01 = (sumxy - sumx*sumy/n) / denominator;
	*pq10 = *pq01;
	*pq11 = (sumy2 - sumy*sumy/n) / denominator;
}

