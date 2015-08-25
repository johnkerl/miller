#include <stdlib.h>
#include "cli/mlrcli.h" // xxx move QUOTE_* to another header file
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/mixutil.h"
#include "output/lrec_writers.h"

// See https://github.com/johnkerl/miller/issues/4
// Temporary status:
// * --csv     from the command line maps into the (existing) csvlite I/O
// * --csvlite from the command line maps into the (existing) csvlite I/O
// * --csvex   from the command line maps into the (new & experimental & unadvertised) rfc-csv I/O
// Ultimate status:
// * --csvlite from the command line will maps into the csvlite I/O
// * --csv     from the command line will maps into the rfc-csv I/O

typedef void     quoted_output_func_t(FILE* fp, char* string, char ors, char ofs);
static void     quote_all_output_func(FILE* fp, char* string, char ors, char ofs);
static void    quote_none_output_func(FILE* fp, char* string, char ors, char ofs);
static void quote_minimal_output_func(FILE* fp, char* string, char ors, char ofs);
static void quote_numeric_output_func(FILE* fp, char* string, char ors, char ofs);

typedef struct _lrec_writer_csv_state_t {
	int  onr;
	char ors; // xxx char -> char*
	char ofs; // xxx char -> char*
	quoted_output_func_t* pquoted_output_func;
	long long num_header_lines_output;
	slls_t* plast_header_output;
} lrec_writer_csv_state_t;

// ----------------------------------------------------------------
// xxx cmt mem-mgmt

static void lrec_writer_csv_process(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_csv_state_t* pstate = pvstate;
	char ors = pstate->ors;
	char ofs = pstate->ofs;

	if (pstate->plast_header_output != NULL) {
		// xxx make a fcn to compare these w/o copy: put it in mixutil.
		if (!lrec_keys_equal_list(prec, pstate->plast_header_output)) {
			slls_free(pstate->plast_header_output);
			pstate->plast_header_output = NULL;
			if (pstate->num_header_lines_output > 0LL)
				fputc(ors, output_stream);
		}
	}

	if (pstate->plast_header_output == NULL) {
		int nf = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			if (nf > 0)
				fputc(ofs, output_stream);
			pstate->pquoted_output_func(output_stream, pe->key, pstate->ors, pstate->ofs);
			nf++;
		}
		fputc(ors, output_stream);
		pstate->plast_header_output = mlr_copy_keys_from_record(prec);
		pstate->num_header_lines_output++;
	}

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(ofs, output_stream);
		pstate->pquoted_output_func(output_stream, pe->value, pstate->ors, pstate->ofs);
		nf++;
	}
	fputc(ors, output_stream);
	pstate->onr++;

	lrec_free(prec); // xxx cmt mem-mgmt
}

static void lrec_writer_csv_free(void* pvstate) {
	lrec_writer_csv_state_t* pstate = pvstate;
	if (pstate->plast_header_output != NULL) {
		slls_free(pstate->plast_header_output);
		pstate->plast_header_output = NULL;
	}
}

lrec_writer_t* lrec_writer_csv_alloc(char ors, char ofs, int oquoting) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_csv_state_t));
	pstate->onr                     = 0;
	pstate->ors                     = ors;
	pstate->ofs                     = ofs;
	switch(oquoting) {
	case QUOTE_ALL:     pstate->pquoted_output_func = quote_all_output_func;     break;
	case QUOTE_NONE:    pstate->pquoted_output_func = quote_none_output_func;    break;
	case QUOTE_MINIMAL: pstate->pquoted_output_func = quote_minimal_output_func; break;
	case QUOTE_NUMERIC: pstate->pquoted_output_func = quote_numeric_output_func; break;
	default:
		fprintf(stderr, "%s: internal coding error: output-quoting style 0x%x unrecognized.\n",
			MLR_GLOBALS.argv0, oquoting);
		exit(1);
	}

	pstate->num_header_lines_output = 0LL;
	pstate->plast_header_output     = NULL;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = lrec_writer_csv_process;
	plrec_writer->pfree_func    = lrec_writer_csv_free;

	return plrec_writer;
}

static void quote_all_output_func(FILE* fp, char* string, char ors, char ofs) {
	fputc('"', fp);
	fputs(string, fp);
	fputc('"', fp);
}

static void quote_none_output_func(FILE* fp, char* string, char ors, char ofs) {
	fputs(string, fp);
}

static void quote_minimal_output_func(FILE* fp, char* string, char ors, char ofs) {
	int output_quotes = FALSE;
	for (char* p = string; *p; p++) {
		if (*p == ors || *p == ofs) {
			output_quotes = TRUE;
			break;
		}
	}
	if (output_quotes) {
		fputc('"', fp);
		fputs(string, fp);
		fputc('"', fp);
	} else {
		fputs(string, fp);
	}
}

static void quote_numeric_output_func(FILE* fp, char* string, char ors, char ofs) {
	double temp;
	if (mlr_try_double_from_string(string, &temp)) {
		fputc('"', fp);
		fputs(string, fp);
		fputc('"', fp);
	} else {
		fputs(string, fp);
	}
}
