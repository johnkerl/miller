#include <stdlib.h>
#include <string.h>
#include "string_builder.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"

// To avoid heap fragmentation, round alloc lengths up to a round number.
#define BLOCK_SIZE 32
#define BLOCK_LENGTH_MASK  (BLOCK_SIZE-1)
#define BLOCK_LENGTH_NMASK (~(BLOCK_SIZE-1))

// ----------------------------------------------------------------
string_builder_t* sb_alloc(int alloc_length) {
	string_builder_t* psb = mlr_malloc_or_die(sizeof(string_builder_t));
	sb_init(psb, alloc_length);
	return psb;
}

// ----------------------------------------------------------------
void  sb_free(string_builder_t* psb) {
	if (psb == NULL)
		return;
	free(psb->buffer);
	free(psb);
}

// ----------------------------------------------------------------
void sb_init(string_builder_t* psb, int alloc_length) {
	if (alloc_length < 1) {
		fprintf(stderr, "%s: string_builder alloc_length must be >= 1; got %d.\n",
			MLR_GLOBALS.bargv0, alloc_length);
		exit(1);
	}
	psb->used_length = 0;
	psb->alloc_length = alloc_length;
	psb->buffer = mlr_malloc_or_die(alloc_length);
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
	int alloc_length = (psb->used_length + BLOCK_LENGTH_MASK) & BLOCK_LENGTH_NMASK;
	char* rv = mlr_malloc_or_die(alloc_length);
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
