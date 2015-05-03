#ifndef MLRSTAT_H
#define MLRSTAT_H

double mlr_get_stddev(unsigned long long count, double sum, double sum2);
double mlr_get_cov(unsigned long long count, double sumx, double sumy, double sumxy);

#endif // MLRSTAT_H
