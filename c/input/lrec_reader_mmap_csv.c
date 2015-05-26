#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

// Idea of pheader_keepers: each header_keeper object retains the input-line backing
// and the slls_t for a CSV header line which is used by one or more CSV data
// lines.  Meanwhile some mappers retain input records from the entire data
// stream, including header-schema changes in the input stream. This means we
// need to keep headers intact as long as any lrecs are pointing to them.  One
// option is reference-counting which I experimented with; it was messy and
// error-prone. The approach used here is to keep a hash map from header-schema
// to header_keeper object. The current pheader_keeper is a pointer into one of
// those.  Then when the reader is freed, all the header-keepers are freed.

typedef struct _lrec_reader_mmap_csv_state_t {
	long long  ifnr; // xxx cmt w/r/t pctx
	long long  ilno; // xxx cmt w/r/t pctx
	char irs;
	char ifs;
	int  allow_repeat_ifs;

	int  expect_header_line_next;
	header_keeper_t* pheader_keeper;
	lhmslv_t*     pheader_keepers;
} lrec_reader_mmap_csv_state_t;

// Cases:
//
// a,a        a,b        c          d
// -- FILE1:  -- FILE1:  -- FILE1:  -- FILE1:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// -- FILE2:  -- FILE2:
// a,b,c      d,e,f,g    a,b,c      d,e,f
// 7,8,9      3,4,5,6    7,8,9      3,4,5
// --OUTPUT:  --OUTPUT:  --OUTPUT:  --OUTPUT:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// 7,8,9                 7,8,9
//            d,e,f,g               d,e,f
//            3,4,5,6               3,4,5

// ----------------------------------------------------------------
// xxx needs abend on null lhs.
//
// etc.

static lrec_t* lrec_reader_mmap_csv_process(file_reader_mmap_state_t* phandle, void* pvstate, context_t* pctx) {
	return NULL; // xxx stub
//	lrec_reader_mmap_csv_state_t* pstate = pvstate;
//
//	while (TRUE) {
//		if (pstate->expect_header_line_next) {
//			// xxx cmt
//			while (TRUE) {

//				xxx need a func to take mmap handle & return slls_t* pheader_fields

//				char* hline = mlr_get_line(input_stream, pstate->irs);
//				if (hline == NULL) // EOF
//					return NULL;
//				pstate->ilno++;
//
//				slls_t* pheader_fields = split_csv_header_line(hline, pstate->ifs, pstate->allow_repeat_ifs);
//				if (pheader_fields->length == 0) {
//					pstate->expect_header_line_next = TRUE;
//					if (pstate->pheader_keeper != NULL) {
//						pstate->pheader_keeper = NULL;
//					}
//				} else {
//					pstate->expect_header_line_next = FALSE;
//
//					pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
//					if (pstate->pheader_keeper == NULL) {
//						pstate->pheader_keeper = header_keeper_alloc(hline, pheader_fields);
//						lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper);
//					} else { // Re-use the header-keeper in the header cache
//						slls_free(pheader_fields);
//					}
//					break;
//				}
//			}
//		}
//

//		xxx func not unlike the mmapnidx parser

//		char* line = mlr_get_line(input_stream, pstate->irs);
//		if (line == NULL) // EOF
//			return NULL;
//
//		// xxx empty-line check ... make a lib func is_empty_modulo_whitespace().
//		if (!*line) {
//			if (pstate->pheader_keeper != NULL) {
//				pstate->pheader_keeper = NULL;
//				pstate->expect_header_line_next = TRUE;
//				continue;
//			}
//		} else {
//			pstate->ifnr++;
//			return lrec_parse_mmap_csv(pstate->pheader_keeper, phandle, pstate->irs, pstate->ifs, pstate->allow_repeat_ifs);
//		}
//	}
}

// ----------------------------------------------------------------
static void lrec_reader_mmap_csv_sof(void* pvstate) {
	lrec_reader_mmap_csv_state_t* pstate = pvstate;
	pstate->ifnr = 0LL;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = TRUE;
}

// xxx restore free func for header_keepers ... ?

// ----------------------------------------------------------------
lrec_reader_mmap_t* lrec_reader_mmap_csv_alloc(char irs, char ifs, int allow_repeat_ifs) {
	lrec_reader_mmap_t* plrec_reader_mmap = mlr_malloc_or_die(sizeof(lrec_reader_mmap_t));

	lrec_reader_mmap_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_csv_state_t));
	pstate->ifnr                     = 0LL;
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->allow_repeat_ifs         = allow_repeat_ifs;
	pstate->expect_header_line_next  = TRUE;
	pstate->pheader_keeper           = NULL;
	pstate->pheader_keepers          = lhmslv_alloc();

	plrec_reader_mmap->pvstate       = (void*)pstate;
	plrec_reader_mmap->pprocess_func = &lrec_reader_mmap_csv_process;
	plrec_reader_mmap->psof_func     = &lrec_reader_mmap_csv_sof;

	return plrec_reader_mmap;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_csv(header_keeper_t* pheader_keeper, file_reader_mmap_state_t* phandle, char irs, char ifs, int allow_repeat_ifs) {
	return NULL; // xxx stub
}
