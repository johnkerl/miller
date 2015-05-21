#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/mixutil.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_pprint_state_t {
	sllv_t*    precords;
	slls_t*    pprev_keys;
	int        left_align;
	long long  num_blocks_written;
} lrec_writer_pprint_state_t;

static void print_and_free_record_list(sllv_t* precords, FILE* output_stream, int left_align);

// ----------------------------------------------------------------
static void lrec_writer_pprint_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	lrec_writer_pprint_state_t* pstate = pvstate;

	int drain = FALSE;
	slls_t* pcurr_keys = NULL;

	if (prec == NULL) {
		drain = TRUE;
	}
	else {
		// xxx make a fcn which does the cmp of lrec & slls w/o the copy
		pcurr_keys = mlr_keys_from_record(prec);
		if (pstate->pprev_keys != NULL && !slls_equals(pcurr_keys, pstate->pprev_keys)) {
			drain = TRUE;
		}
	}

	if (drain) {
		if (pstate->num_blocks_written > 0LL) // xxx cmt
			fputc('\n', output_stream);
		print_and_free_record_list(pstate->precords, output_stream, pstate->left_align);
		if (pstate->pprev_keys != NULL) {
			slls_free(pstate->pprev_keys);
			pstate->pprev_keys = NULL;
		}
		pstate->precords = sllv_alloc();
		pstate->num_blocks_written++;
	}
	if (prec != NULL) {
		sllv_add(pstate->precords, prec);
		if (pstate->pprev_keys == NULL)
			pstate->pprev_keys = slls_copy(pcurr_keys);
	}
}

// ----------------------------------------------------------------
static void print_and_free_record_list(sllv_t* precords, FILE* output_stream, int left_align) {
	if (precords->length == 0)
		return;
	lrec_t* prec1 = precords->phead->pvdata;

	int* max_widths = mlr_malloc_or_die(sizeof(int) * prec1->field_count);
	int j = 0;
	for (lrece_t* pe = prec1->phead; pe != NULL; pe = pe->pnext, j++) {
		max_widths[j] = strlen(pe->key);
	}
	for (sllve_t* pnode = precords->phead; pnode != NULL; pnode = pnode->pnext) {
		lrec_t* prec = pnode->pvdata;
		j = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
			int width = strlen(pe->value);
			if (width > max_widths[j])
				max_widths[j] = width;
		}
	}

	int onr = 0;
	for (sllve_t* pnode = precords->phead; pnode != NULL; pnode = pnode->pnext, onr++) {
		lrec_t* prec = pnode->pvdata;

		if (onr == 0) {
			j = 0;
			for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
				if (j > 0) {
					fputc(' ', output_stream);
				}
				if (left_align) {
					if (pe->pnext == NULL)
						fprintf(output_stream, "%s", pe->key);
					else
						fprintf(output_stream, "%-*s", max_widths[j], pe->key);
				} else {
					fprintf(output_stream, "%*s", max_widths[j], pe->key);
				}
			}
			fputc('\n', output_stream);
		}

		j = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
			if (j > 0) {
				fputc(' ', output_stream);
			}
			char* value = pe->value;
			if (*value == 0) // empty string
				value = "-";
			if (left_align) {
				if (pe->pnext == NULL)
					fprintf(output_stream, "%s", value);
				else
					fprintf(output_stream, "%-*s", max_widths[j], value);
			} else {
				fprintf(output_stream, "%*s", max_widths[j], value);
			}
		}
		fputc('\n', output_stream);

		lrec_free(prec); // xxx cmt mem-mgmt
	}

	free(max_widths);
	sllv_free(precords);
}

static void lrec_writer_pprint_free(void* pvstate) {
	lrec_writer_pprint_state_t* pstate = pvstate;
	if (pstate->precords != NULL) {
		sllv_free(pstate->precords);
		pstate->precords = NULL;
	}
	if (pstate->pprev_keys != NULL) {
		slls_free(pstate->pprev_keys);
		pstate->pprev_keys = NULL;
	}
}

lrec_writer_t* lrec_writer_pprint_alloc(int left_align) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_pprint_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_pprint_state_t));
	pstate->precords           = sllv_alloc();
	pstate->pprev_keys         = NULL;
	pstate->left_align         = left_align;
	pstate->num_blocks_written = 0LL;
	plrec_writer->pvstate           = pstate;

	plrec_writer->plrec_writer_func = &lrec_writer_pprint_func;
	plrec_writer->pfree_func   = &lrec_writer_pprint_free;

	return plrec_writer;
}
