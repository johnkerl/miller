// ================================================================
// For mlr stats1 percentiles
// ================================================================

#ifndef PERCENTILE_KEEPER_H
#define PERCENTILE_KEEPER_H
#include "containers/lrec.h"

typedef struct _percentile_keeper_t {
	double* data;
	int     size;
	int     capacity;
	int     sorted;
} percentile_keeper_t;

percentile_keeper_t* percentile_keeper_alloc();
void percentile_keeper_free(percentile_keeper_t* ppercentile_keeper);
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, double value);
double percentile_keeper_emit(percentile_keeper_t* ppercentile_keeper, double percentile);

#endif // PERCENTILE_KEEPER_H
