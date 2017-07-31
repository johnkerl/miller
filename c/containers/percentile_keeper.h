// ================================================================
// For mlr stats1 percentiles
// ================================================================

#ifndef PERCENTILE_KEEPER_H
#define PERCENTILE_KEEPER_H
#include "lib/mlrval.h"

typedef struct _percentile_keeper_t {
	mv_t* data;
	unsigned long long size;
	unsigned long long capacity;
	int   sorted;
} percentile_keeper_t;

percentile_keeper_t* percentile_keeper_alloc();
void percentile_keeper_free(percentile_keeper_t* ppercentile_keeper);
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, mv_t value);

typedef mv_t percentile_keeper_emitter_t(percentile_keeper_t* ppercentile_keeper, double percentile);
mv_t percentile_keeper_emit_non_interpolated(percentile_keeper_t* ppercentile_keeper, double percentile);
mv_t percentile_keeper_emit_linearly_interpolated(percentile_keeper_t* ppercentile_keeper, double percentile);

// For debug/test
void percentile_keeper_print(percentile_keeper_t* ppercentile_keeper);

#endif // PERCENTILE_KEEPER_H
