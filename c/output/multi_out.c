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
void multi_out_free(multi_out_t* pmo) {
	if (pmo == NULL)
		return;

	lhmsv_free(pmo->pnames_to_fps);
	free(pmo);
}

// ----------------------------------------------------------------
FILE* multi_out_get(multi_out_t* pmo, char* filename, file_output_mode_t file_output_mode) {
	FILE* outfp = lhmsv_get(pmo->pnames_to_fps, filename);
	if (outfp == NULL) {
		char* mode_string = get_mode_string(file_output_mode);
		char* mode_desc = get_mode_desc(file_output_mode);
		outfp = fopen(filename, mode_string);
		if (outfp == NULL) {
			perror("fopen");
			fprintf(stderr, "%s: failed fopen for %s of \"%s\".\n",
				MLR_GLOBALS.bargv0, mode_desc, filename);
			exit(1);
		}
		lhmsv_put(pmo->pnames_to_fps, mlr_strdup_or_die(filename), outfp, FREE_ENTRY_KEY);
	}
	return outfp;
}
