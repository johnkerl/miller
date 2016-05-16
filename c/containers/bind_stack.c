#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/bind_stack.h"

bind_stack_t* bind_stack_alloc() {
	return NULL; // xxx stub
}

void bind_stack_free(bind_stack_t* pstack) {
}

void bind_stack_push(bind_stack_t* pstack, mlhmmv_t* pframe) {
}

void bind_stack_pop(bind_stack_t* pstack) {
}

mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key) {
	return NULL; // xxx stub
}
