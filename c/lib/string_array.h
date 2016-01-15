#ifndef STRING_ARRAY_H
#define STRING_ARRAY_H

// Container class for keeping an array of strings, some of which may be null.
typedef struct _string_array_t {
	int length;
	int strings_need_freeing;
	char** strings;
} string_array_t;

string_array_t* string_array_alloc(int length);
void string_array_free(string_array_t* parray);
void string_array_realloc(string_array_t* parray, int new_length);
string_array_t* string_array_from_line(char* line, char ifs);

#endif // STRING_ARRAY_H
