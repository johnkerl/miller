#ifndef BIND_STACK_H
#define BIND_STACK_H

#include "containers/mlrval.h"
#include "containers/lhmsmv.h"

// Bound & scoped variables for use in for-loops.
typedef struct _bind_stack_t {
	int num_used;
	int num_allocated;
	lhmsmv_t** pframes;
} bind_stack_t;

bind_stack_t* bind_stack_alloc();
void bind_stack_free(bind_stack_t* pstack);

void bind_stack_push(bind_stack_t* pstack, lhmsmv_t* pframe);
void bind_stack_pop(bind_stack_t* pstack);
mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key);

void bind_stack_print(bind_stack_t* pstack);

#endif // BIND_STACK_H
