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
static inline FILE* multi_out_get(multi_out_t* pmo, char* filename,
	char* mode, char* mode_desc)
{
	FILE* outfp = lhmsv_get(pmo->pnames_to_fps, filename);
	if (outfp == NULL) {
		outfp = fopen(filename, mode);
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

// ----------------------------------------------------------------
FILE* multi_out_get_for_write(multi_out_t* pmo, char* filename) {
	return multi_out_get(pmo, filename, "w", "write");
}

// ----------------------------------------------------------------
FILE* multi_out_get_for_append(multi_out_t* pmo, char* filename) {
	return multi_out_get(pmo, filename, "a", "append");
}
