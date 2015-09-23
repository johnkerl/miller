#include <stdlib.h>
#include <string.h>
#include "string_builder.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"

// ----------------------------------------------------------------
void sb_init(string_builder_t* psb, int alloc_length) {
	if (alloc_length < 1) {
		fprintf(stderr, "%s: string_builder alloc_length must be >= 1; got %d.\n",
			MLR_GLOBALS.argv0, alloc_length);
		exit(1);
	}
	psb->used_length = 0;
	psb->alloc_length = alloc_length;
	psb->buffer = mlr_malloc_or_die(alloc_length); // xxx malloc ...
}

// ----------------------------------------------------------------
void sb_append_string(string_builder_t* psb, char* s) {
	for (char* p = s; *p; p++)
		sb_append_char(psb, *p);
}

// ----------------------------------------------------------------
int sb_is_empty(string_builder_t* psb) {
	return psb->used_length == 0;
}

// ----------------------------------------------------------------
// Keep and reuse the internal buffer. Allocate to the caller only
// the size needed.
char* sb_finish(string_builder_t* psb) {
	sb_append_char(psb, '\0');
	char* rv = mlr_malloc_or_die(psb->used_length);
	memcpy(rv, psb->buffer, psb->used_length);
	psb->used_length  = 0;
	return rv;
}

// ----------------------------------------------------------------
void _sb_enlarge(string_builder_t* psb) {
	int new_alloc_length = psb->alloc_length * 2;
	char* new_buffer = mlr_malloc_or_die(new_alloc_length);
	memcpy(new_buffer, psb->buffer, psb->used_length);
	psb->alloc_length = new_alloc_length;
	psb->buffer = new_buffer;
}
