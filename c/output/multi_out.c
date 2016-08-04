#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "multi_out.h"

// ----------------------------------------------------------------
multi_out_t* multi_out_alloc() {
	multi_out_t* pmo = mlr_malloc_or_die(sizeof(multi_out_t));
	pmo->pnames_to_fps = lhmsv_alloc();
	return pmo;
}

// ----------------------------------------------------------------
void multi_out_close(multi_out_t* pmo) {
	for (lhmsve_t* pe = pmo->pnames_to_fps->phead; pe != NULL; pe = pe->pnext) {
		fp_and_flag_t* pstate = pe->pvvalue;
		if (pstate->is_popen) {
			pclose(pstate->output_stream);
		} else {
			fclose(pstate->output_stream);
		}
	}
}

// ----------------------------------------------------------------
void multi_out_free(multi_out_t* pmo) {
	if (pmo == NULL)
		return;
	lhmsv_free(pmo->pnames_to_fps);
	free(pmo);
}

// ----------------------------------------------------------------
FILE* multi_out_get(multi_out_t* pmo, char* filename_or_command, file_output_mode_t file_output_mode) {
	fp_and_flag_t* pstate = lhmsv_get(pmo->pnames_to_fps, filename_or_command);
	if (pstate == NULL) {
		pstate = mlr_malloc_or_die(sizeof(fp_and_flag_t));
		char* mode_string = get_mode_string(file_output_mode);
		char* mode_desc = get_mode_desc(file_output_mode);
		if (file_output_mode == MODE_PIPE) {
			pstate->is_popen = TRUE;
			pstate->output_stream = popen(filename_or_command, mode_string);
			if (pstate->output_stream == NULL) {
				perror("popen");
				fprintf(stderr, "%s: failed popen for %s of \"%s\".\n",
					MLR_GLOBALS.bargv0, mode_desc, filename_or_command);
				exit(1);
			}
		} else {
			pstate->is_popen = FALSE;
			pstate->output_stream = fopen(filename_or_command, mode_string);
			if (pstate->output_stream == NULL) {
				perror("fopen");
				fprintf(stderr, "%s: failed fopen for %s of \"%s\".\n",
					MLR_GLOBALS.bargv0, mode_desc, filename_or_command);
				exit(1);
			}
		}
		lhmsv_put(pmo->pnames_to_fps, mlr_strdup_or_die(filename_or_command), pstate, FREE_ENTRY_KEY);
	}
	return pstate->output_stream;
}
