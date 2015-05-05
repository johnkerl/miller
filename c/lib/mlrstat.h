#ifndef MLRSTAT_H
#define MLRSTAT_H

double mlr_get_stddev(unsigned long long count, double sum, double sum2);
double mlr_get_cov(unsigned long long count, double sumx, double sumy, double sumxy);
void mlr_get_cov_matrix(unsigned long long n,
	double sumx, double sumx2, double sumy, double sumy2, double sumxy,
	double* pq00, double* pq01, double* pq10, double* pq11);

#endif // MLRSTAT_H
