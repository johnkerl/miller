// ================================================================
// For mlr stats1 percentiles
// ================================================================

#ifndef PERCENTILE_KEEPER_H
#define PERCENTILE_KEEPER_H
#include "containers/mlrval.h"

typedef struct _percentile_keeper_t {
	mv_t* data;
	int   size;
	int   capacity;
	int   sorted;
} percentile_keeper_t;

percentile_keeper_t* percentile_keeper_alloc();
void percentile_keeper_free(percentile_keeper_t* ppercentile_keeper);
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, mv_t value);
mv_t percentile_keeper_emit(percentile_keeper_t* ppercentile_keeper, double percentile);

// For debug/test
void percentile_keeper_print(percentile_keeper_t* ppercentile_keeper);

#endif // PERCENTILE_KEEPER_H
