#include <stdlib.h>
#include "cli/quoting.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/mixutil.h"
#include "output/lrec_writers.h"

typedef void       quoted_output_func_t(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static  void      quote_all_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static  void     quote_none_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static  void  quote_minimal_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static  void  quote_minimal_auto_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char qf);
static  void  quote_numeric_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static  void quote_original_output_func(FILE* fp,char*s,char*ors,char*ofs, int orslen,int ofslen, char quote_flags);
static void            csv_quote_string(FILE* fp, char* string);

typedef struct _lrec_writer_csv_state_t {
	int   onr;
	char *ors;
	char *ofs;
	int   orslen;
	int   ofslen;
	quoted_output_func_t* pquoted_output_func;
	long long num_header_lines_output;
	slls_t* plast_header_output;
	int headerless_csv_output;
} lrec_writer_csv_state_t;

// ----------------------------------------------------------------
static void lrec_writer_csv_process(void* pvstate, FILE* output_stream, lrec_t* prec, char* ors);
static void lrec_writer_csv_process_auto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx);
static void lrec_writer_csv_process_nonauto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx);
static void lrec_writer_csv_free(lrec_writer_t* pwriter, context_t* pctx);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_csv_alloc(char* ors, char* ofs, quoting_t oquoting, int headerless_csv_output) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_csv_state_t));
	pstate->onr    = 0;
	pstate->ors    = ors;
	pstate->ofs    = ofs;
	pstate->orslen = strlen(pstate->ors);
	pstate->ofslen = strlen(pstate->ofs);
	pstate->headerless_csv_output = headerless_csv_output;

	switch(oquoting) {
	case QUOTE_ALL:      pstate->pquoted_output_func = quote_all_output_func;      break;
	case QUOTE_NONE:     pstate->pquoted_output_func = quote_none_output_func;     break;
	case QUOTE_MINIMAL:  pstate->pquoted_output_func = quote_minimal_output_func;  break;
	case QUOTE_NUMERIC:  pstate->pquoted_output_func = quote_numeric_output_func;  break;
	case QUOTE_ORIGINAL: pstate->pquoted_output_func = quote_original_output_func; break;
	default:
		MLR_INTERNAL_CODING_ERROR();
	}

	pstate->num_header_lines_output = 0LL;
	pstate->plast_header_output     = NULL;

	plrec_writer->pvstate = (void*)pstate;
	if (streq(ors, "auto")) {
		plrec_writer->pprocess_func = lrec_writer_csv_process_auto_ors;
		if (oquoting == QUOTE_MINIMAL) {
			pstate->pquoted_output_func = quote_minimal_auto_output_func;
		}
	} else {
		plrec_writer->pprocess_func = lrec_writer_csv_process_nonauto_ors;
	}
	plrec_writer->pfree_func = lrec_writer_csv_free;

	return plrec_writer;
}

static void lrec_writer_csv_free(lrec_writer_t* pwriter, context_t* pctx) {
	lrec_writer_csv_state_t* pstate = pwriter->pvstate;
	slls_free(pstate->plast_header_output);
	free(pstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_csv_process_auto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx) {
	lrec_writer_csv_process(pvstate, output_stream, prec, pctx->auto_line_term);
}

static void lrec_writer_csv_process_nonauto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx) {
	lrec_writer_csv_state_t* pstate = pvstate;
	lrec_writer_csv_process(pvstate, output_stream, prec, pstate->ors);
}

static void lrec_writer_csv_process(void* pvstate, FILE* output_stream, lrec_t* prec, char* ors) {
	if (prec == NULL)
		return;
	lrec_writer_csv_state_t* pstate = pvstate;
	char *ofs = pstate->ofs;
	int orslen = strlen(ors);

	if (pstate->plast_header_output != NULL) {
		if (!lrec_keys_equal_list(prec, pstate->plast_header_output)) {
			slls_free(pstate->plast_header_output);
			pstate->plast_header_output = NULL;
			if (pstate->num_header_lines_output > 0LL)
				fputs(ors, output_stream);
		}
	}

	if (pstate->plast_header_output == NULL) {
		int nf = 0;
		if (!pstate->headerless_csv_output) {
			for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
				if (nf > 0)
					fputs(ofs, output_stream);
				pstate->pquoted_output_func(output_stream, pe->key, pstate->ors, pstate->ofs,
					orslen, pstate->ofslen, 0);
				nf++;
			}
			fputs(ors, output_stream);
		}
		pstate->plast_header_output = mlr_copy_keys_from_record(prec);
		pstate->num_header_lines_output++;
	}

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputs(ofs, output_stream);
		pstate->pquoted_output_func(output_stream, pe->value, pstate->ors, pstate->ofs,
			orslen, pstate->ofslen, pe->quote_flags);
		nf++;
	}
	fputs(ors, output_stream);
	pstate->onr++;

	// See ../README.md for memory-management conventions
	lrec_free(prec);
}

// ----------------------------------------------------------------
static void quote_all_output_func(FILE* fp, char* string, char* ors, char* ofs, int orslen, int ofslen,
	char quote_flags)
{
	csv_quote_string(fp, string);
}

static void quote_none_output_func(FILE* fp, char* string, char* ors, char* ofs, int orslen, int ofslen,
	char quote_flags)
{
	fputs(string, fp);
}

static void quote_minimal_output_func(FILE* fp, char* string, char* ors, char* ofs, int orslen, int ofslen,
	char quote_flags)
{
	int output_quotes = FALSE;
	for (char* p = string; *p; p++) {
		if (streqn(p, ors, orslen) || streqn(p, ofs, ofslen)) {
			output_quotes = TRUE;
			break;
		}
		if (*p == '"') {
			output_quotes = TRUE;
			break;
		}
	}
	if (output_quotes) {
		csv_quote_string(fp, string);
	} else {
		fputs(string, fp);
	}
}

static void quote_minimal_auto_output_func(FILE* fp, char* string, char* _, char* ofs, int __, int ofslen,
	char quote_flags)
{
	int output_quotes = FALSE;
	for (char* p = string; *p; p++) {
		if (streqn(p, "\n", 1) || streqn(p, "\r\n", 2) || streqn(p, ofs, ofslen)) {
			output_quotes = TRUE;
			break;
		}
		if (*p == '"') {
			output_quotes = TRUE;
			break;
		}
	}
	if (output_quotes) {
		csv_quote_string(fp, string);
	} else {
		fputs(string, fp);
	}
}

static void quote_numeric_output_func(FILE* fp, char* string, char* ors, char* ofs, int orslen, int ofslen,
	char quote_flags)
{
	double temp;
	if (mlr_try_float_from_string(string, &temp)) {
		csv_quote_string(fp, string);
	} else {
		fputs(string, fp);
	}
}

static void quote_original_output_func(FILE* fp, char* string, char* ors, char* ofs, int orslen, int ofslen,
	char quote_flags)
{
	if (quote_flags & FIELD_QUOTED_ON_INPUT) {
		csv_quote_string(fp, string);
	} else {
		fputs(string, fp);
	}
}

// ----------------------------------------------------------------
static void csv_quote_string(FILE* fp, char* string) {
	fputc('"', fp);
	for (char* p = string; *p; p++) {
		if (*p == '"')
			fputs("\"\"", fp);
		else
			fputc(*p, fp);
	}
	fputc('"', fp);
}
