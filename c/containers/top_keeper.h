#ifndef TOP_KEEPER_H
#define TOP_KEEPER_H
#include "containers/lrec.h"

// xxx option to keep entire records ...
typedef struct _top_keeper_t {
	double*  top_values;
	lrec_t** top_precords;
	int     capacity;
	int     size;
} top_keeper_t;

top_keeper_t* top_keeper_alloc(int capacity);
void top_keeper_free(top_keeper_t* ptop_keeper);
void top_keeper_add(top_keeper_t* ptop_keeper, double value, lrec_t* prec);

#endif // TOP_KEEPER_H
