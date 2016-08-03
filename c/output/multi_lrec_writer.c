#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "cli/mlrcli.h"
#include "output/multi_lrec_writer.h"

// ----------------------------------------------------------------
multi_lrec_writer_t* multi_lrec_writer_alloc() {
	multi_lrec_writer_t* pmlw = mlr_malloc_or_die(sizeof(multi_lrec_writer_t));
	pmlw->pnames_to_lrec_writers_and_fps = lhmsv_alloc();
	return pmlw;
}

// ----------------------------------------------------------------
void multi_lrec_writer_free(multi_lrec_writer_t* pmlw) {
	if (pmlw == NULL)
		return;

	// xxx in loop: free writers and filenames
	for (lhmsve_t* pe = pmlw->pnames_to_lrec_writers_and_fps->phead; pe != NULL; pe = pe->pnext) {
		lrec_writer_and_fp_t* pstate = pe->pvvalue;
		pstate->plrec_writer->pfree_func(pstate->plrec_writer);
		free(pstate->filename);
	}

	lhmsv_free(pmlw->pnames_to_lrec_writers_and_fps);
	free(pmlw);
}

// ----------------------------------------------------------------
void multi_lrec_writer_output_srec(multi_lrec_writer_t* pmlw, lrec_t* poutrec, char* filename,
	file_output_mode_t file_output_mode, int flush_every_record)
{
	lrec_writer_and_fp_t* pstate = lhmsv_get(pmlw->pnames_to_lrec_writers_and_fps, filename);
	if (pstate == NULL) {
		pstate = mlr_malloc_or_die(sizeof(lrec_writer_and_fp_t)); // xxx free at destruct
		cli_opts_t* popts = MLR_GLOBALS.popts;
		pstate->plrec_writer = lrec_writer_alloc(popts->ofile_fmt, popts->ors, popts->ofs, popts->ops,
			popts->headerless_csv_output, popts->oquoting, popts->left_align_pprint,
			popts->right_justify_xtab_value, popts->json_flatten_separator, popts->quote_json_values_always,
			popts->stack_json_output_vertically, popts->wrap_json_output_in_outer_list);
		if (pstate->plrec_writer == NULL) {
			fprintf(stderr, "%s: internal coding error detected in file \"%s\" at line %d.\n",
				MLR_GLOBALS.bargv0, __FILE__, __LINE__);
			exit(1);
		}
		pstate->filename = mlr_strdup_or_die(filename);
		char* mode_string = get_mode_string(file_output_mode);
		char* mode_desc = get_mode_desc(file_output_mode);
		pstate->output_stream = fopen(filename, mode_string);
		if (pstate->output_stream == NULL) {
			perror("fopen");
			fprintf(stderr, "%s: failed fopen for %s on \"%s\".\n", MLR_GLOBALS.bargv0, mode_desc, filename);
			exit(1);
		}

		lhmsv_put(pmlw->pnames_to_lrec_writers_and_fps, mlr_strdup_or_die(filename), pstate, FREE_ENTRY_KEY);
	}

	pstate->plrec_writer->pprocess_func(pstate->output_stream, poutrec, pstate->plrec_writer->pvstate);

	if (poutrec != NULL) {
		if (flush_every_record)
			fflush(pstate->output_stream);
	} else {
		if (fclose(pstate->output_stream) != 0) {
			perror("fclose");
			fprintf(stderr, "%s: fclose error on \"%s\".\n", MLR_GLOBALS.bargv0, filename);
			exit(1);
		}
		pstate->output_stream = NULL;
	}
}

void multi_lrec_writer_output_list(multi_lrec_writer_t* pmlw, sllv_t* poutrecs, char* filename,
	file_output_mode_t file_output_mode, int flush_every_record)
{
	if (poutrecs == NULL) // synonym for empty record-list
		return;

	while (poutrecs->phead) {
		lrec_t* poutrec = sllv_pop(poutrecs);
		multi_lrec_writer_output_srec(pmlw, poutrec, filename, file_output_mode, flush_every_record);
	}
}

void multi_lrec_writer_drain(multi_lrec_writer_t* pmlw) {
	for (lhmsve_t* pe = pmlw->pnames_to_lrec_writers_and_fps->phead; pe != NULL; pe = pe->pnext) {
		lrec_writer_and_fp_t* pstate = pe->pvvalue;
		pstate->plrec_writer->pprocess_func(pstate->output_stream, NULL, pstate->plrec_writer->pvstate);
		fflush(pstate->output_stream);
		if (fclose(pstate->output_stream) != 0) {
			perror("fclose");
			fprintf(stderr, "%s: fclose error on \"%s\".\n", MLR_GLOBALS.bargv0, pstate->filename);
			exit(1);
		}
	}
}
