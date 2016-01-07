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
	char*      ors;
	char       ofs;
} lrec_writer_pprint_state_t;

static void lrec_writer_pprint_free(lrec_writer_t* pwriter);
static void lrec_writer_pprint_process(FILE* output_stream, lrec_t* prec, void* pvstate);
static void print_and_free_record_list(sllv_t* precords, FILE* output_stream, char* ors, char ofs, int left_align);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_pprint_alloc(char* ors, char ofs, int left_align) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_pprint_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_pprint_state_t));
	pstate->precords           = sllv_alloc();
	pstate->pprev_keys         = NULL;
	pstate->ors                = ors;
	pstate->ofs                = ofs;
	pstate->left_align         = left_align;
	pstate->num_blocks_written = 0LL;

	plrec_writer->pvstate       = pstate;
	plrec_writer->pprocess_func = lrec_writer_pprint_process;
	plrec_writer->pfree_func    = lrec_writer_pprint_free;

	return plrec_writer;
}

static void lrec_writer_pprint_free(lrec_writer_t* pwriter) {
	lrec_writer_pprint_state_t* pstate = pwriter->pvstate;
	if (pstate->precords != NULL) {
		sllv_free(pstate->precords);
		pstate->precords = NULL;
	}
	if (pstate->pprev_keys != NULL) {
		slls_free(pstate->pprev_keys);
		pstate->pprev_keys = NULL;
	}
	free(pstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_pprint_process(FILE* output_stream, lrec_t* prec, void* pvstate) {
	lrec_writer_pprint_state_t* pstate = pvstate;

	int drain = FALSE;

	if (prec == NULL) {
		drain = TRUE;
	} else {
		if (pstate->pprev_keys != NULL && !lrec_keys_equal_list(prec, pstate->pprev_keys)) {
			drain = TRUE;
		}
	}

	if (drain) {
		if (pstate->num_blocks_written > 0LL) // separate blocks with empty line
			fputs(pstate->ors, output_stream);
		print_and_free_record_list(pstate->precords, output_stream, pstate->ors, pstate->ofs, pstate->left_align);
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
			pstate->pprev_keys = mlr_copy_keys_from_record(prec);
	}
}

// ----------------------------------------------------------------
static void print_and_free_record_list(sllv_t* precords, FILE* output_stream, char* ors, char ofs, int left_align) {
	if (precords->length == 0) {
		sllv_free(precords);
		return;
	}
	lrec_t* prec1 = precords->phead->pvvalue;

	int* max_widths = mlr_malloc_or_die(sizeof(int) * prec1->field_count);
	int j = 0;
	for (lrece_t* pe = prec1->phead; pe != NULL; pe = pe->pnext, j++) {
		max_widths[j] = strlen_for_utf8_display(pe->key);
	}
	for (sllve_t* pnode = precords->phead; pnode != NULL; pnode = pnode->pnext) {
		lrec_t* prec = pnode->pvvalue;
		j = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
			int width = strlen_for_utf8_display(pe->value);
			if (width > max_widths[j])
				max_widths[j] = width;
		}
	}

	int onr = 0;
	for (sllve_t* pnode = precords->phead; pnode != NULL; pnode = pnode->pnext, onr++) {
		lrec_t* prec = pnode->pvvalue;

		if (onr == 0) {
			j = 0;
			for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
				if (j > 0) {
					fputc(ofs, output_stream);
				}
				if (left_align) {
					if (pe->pnext == NULL) {
						fprintf(output_stream, "%s", pe->key);
					} else {
						// "%-*s" fprintf format isn't correct for non-ASCII UTF-8
						fprintf(output_stream, "%s", pe->key);
						int d = max_widths[j] - strlen_for_utf8_display(pe->key);
						for (int i = 0; i < d; i++)
							fputc(ofs, output_stream);
					}
				} else {
					int d = max_widths[j] - strlen_for_utf8_display(pe->key);
					for (int i = 0; i < d; i++)
						fputc(ofs, output_stream);
					fprintf(output_stream, "%s", pe->key);
				}
			}
			fputs(ors, output_stream);
		}

		j = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext, j++) {
			if (j > 0) {
				fputc(ofs, output_stream);
			}
			char* value = pe->value;
			if (*value == 0) // empty string
				value = "-";
			if (left_align) {
				if (pe->pnext == NULL) {
					fprintf(output_stream, "%s", value);
				} else {
					fprintf(output_stream, "%s", value);
					int d = max_widths[j] - strlen_for_utf8_display(value);
					for (int i = 0; i < d; i++)
						fputc(ofs, output_stream);
				}
			} else {
				int d = max_widths[j] - strlen_for_utf8_display(value);
				for (int i = 0; i < d; i++)
					fputc(ofs, output_stream);
				fprintf(output_stream, "%s", value);
			}
		}
		fputs(ors, output_stream);

		lrec_free(prec); // end of baton-pass
	}

	free(max_widths);
	sllv_free(precords);
}
