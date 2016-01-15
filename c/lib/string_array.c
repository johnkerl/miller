#include <stdlib.h>
#include "lib/mlrutil.h"
#include "string_array.h"

string_array_t* string_array_alloc(int length) {
	string_array_t* parray = mlr_malloc_or_die(sizeof(string_array_t));
	parray->length = length;
	parray->strings_need_freeing = FALSE;
	parray->strings = mlr_malloc_or_die(length * sizeof(char*));
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

void string_array_realloc(string_array_t* parray, int new_length) {
	if (parray->strings_need_freeing)
		for (int i = 0; i < parray->length; i++)
			free(parray->strings[i]);
	if (new_length > parray->length) // else re-use
		parray->strings = mlr_realloc_or_die(parray->strings, new_length * sizeof(char*));
	parray->length = new_length;
	for (int i = 0; i < parray->length; i++)
		parray->strings[i] = NULL;
	parray->strings_need_freeing = FALSE;
}

string_array_t* string_array_from_line(char* line, char ifs) {
	if (*line == 0) // empty string splits to empty array
		return string_array_alloc(0);

	int num_commas = 0;

	for (char* p = line; *p; p++)
		if (*p == ifs)
			num_commas++;

	string_array_t* parray = string_array_alloc(num_commas + 1);

	char* start = line;
	int i = 0;
	for (char* p = line; *p; p++) {
		if (*p == ifs) {
			*p = 0;
			p++;
			parray->strings[i++] = start;
			start = p;
		}
	}
	parray->strings[i++] = start;

	return parray;
}
