#include <math.h>
#include "lib/mlrstat.h"

// xxx cmt intended for streaming applications. otherwise the formulas are
// different (and more intuitive).
double mlr_get_stddev(unsigned long long count, double sum, double sum2) {
	double mean = sum / count;
	double numerator = sum2 - 2.0*mean*sum + count*mean*mean;
	if (numerator < 0.0) // round-off error
		numerator = 0.0;
	double denominator = count - 1LL;
	return sqrt(numerator / denominator);
}

double mlr_get_cov(unsigned long long count, double sumx, double sumy, double sumxy) {
	double meanx = sumx / count;
	double meany = sumy / count;
	double numerator = sumxy - meanx * sumy - meany * sumx + count * meanx * meany;
	double denominator = count - 1;
	return numerator / denominator;
}
