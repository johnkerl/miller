#include <stdlib.h>
#include "lib/mlrutil.h"
#include "string_array.h"

string_array_t* string_array_alloc(int length) {
	string_array_t* parray = mlr_malloc_or_die(sizeof(string_array_t));
	parray->length = length;
	parray->strings_need_freeing = FALSE;
	parray->strings = mlr_malloc_or_die(length * sizeof(char**));
	for (int i = 0; i < length; i++)
		parray->strings[i] = NULL;
	return parray;
}

void string_array_free(string_array_t* parray) {
	if (parray == NULL)
		return;
	if (parray->strings_need_freeing) {
		for (int i = 0; i < parray->length; i++)
			free(parray->strings[i]);
	}
	free(parray->strings);
	free(parray);
}
