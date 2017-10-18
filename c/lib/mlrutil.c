#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include <sys/stat.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/free_flags.h"

// ----------------------------------------------------------------
void mlr_internal_coding_error(char* file, int line) {
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.bargv0, file, line);
	exit(1);
}

void mlr_internal_coding_error_if(int v, char* file, int line) {
	if (v) {
		mlr_internal_coding_error(file, line);
	}
}

void mlr_internal_coding_error_unless(int v, char* file, int line) {
	if (!v) {
		mlr_internal_coding_error(file, line);
	}
}

// ----------------------------------------------------------------
char* mlr_strmsep(char **pstring, const char *sep, int seplen) {
	char* string = *pstring;
	if (string == NULL) {
		return NULL;
	}
	char* pnext = strstr(string, sep);
	if (pnext == NULL) {
		*pstring = NULL;
		return string;
	} else {
		*pnext = 0;
		*pstring = pnext + seplen;
		return string;
	}
}

// ----------------------------------------------------------------
int mlr_bsearch_double_for_insert(double* array, int size, double value) {
	int lo = 0;
	int hi = size-1;
	int mid = (hi+lo)/2;
	int newmid;

	if (size == 0)
		return 0;
	if (value > array[0])
		return 0;
	if (value < array[hi])
		return size;

	while (lo < hi) {
		double a = array[mid];
		if (value == a) {
			return mid;
		}
		else if (value > a) {
			hi = mid;
			newmid = (hi+lo)/2;
		}
		else {
			lo = mid;
			newmid = (hi+lo)/2;
		}
		if (mid == newmid) {
			if (value >= array[lo])
				return lo;
			else if (value >= array[hi])
				return hi;
			else
				return hi+1;
		}
		mid = newmid;
	}

	return lo;
}

// ----------------------------------------------------------------
void* mlr_malloc_or_die(size_t size) {
	void* p = malloc(size);
	if (p == NULL) {
		fprintf(stderr, "malloc(%llu) failed.\n", (unsigned long long)size);
		exit(1);
	}
#ifdef MLR_MALLOC_TRACE
	fprintf(stderr, "MALLOC size=%llu,p=%p\n", (unsigned long long)size, p);
#endif
	return p;
}

// ----------------------------------------------------------------
void* mlr_realloc_or_die(void *optr, size_t size) {
	void* nptr = realloc(optr, size);
	if (nptr == NULL) {
		fprintf(stderr, "realloc(%llu) failed.\n", (unsigned long long)size);
		exit(1);
	}
#ifdef MLR_MALLOC_TRACE
	fprintf(stderr, "REALLOC size=%llu,p=%p\n", (unsigned long long)size, nptr);
#endif
	return nptr;
}

// ----------------------------------------------------------------
char * mlr_strdup_quoted_or_die(const char *s1) {
	int len = strlen(s1);
	char* s2 = mlr_malloc_or_die(len+3);
	s2[0] = '"';
	strcpy(&s2[1], s1);
	s2[len+1] = '"';
	s2[len+2] = 0;
	return s2;
}

// ----------------------------------------------------------------
// The caller should free the return value from each of these.

char* mlr_alloc_string_from_double(double value, char* fmt) {
	int n = snprintf(NULL, 0, fmt, value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, fmt, value);
	return string;
}

char* mlr_alloc_string_from_ull(unsigned long  long value) {
	int n = snprintf(NULL, 0, "%llu", value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "%llu", value);
	return string;
}

char* mlr_alloc_string_from_ll(long  long value) {
	int n = snprintf(NULL, 0, "%lld", value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "%lld", value);
	return string;
}

char* mlr_alloc_string_from_ll_and_format(long long value, char* fmt) {
	int n = snprintf(NULL, 0, fmt, value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, fmt, value);
	return string;
}

char* mlr_alloc_string_from_int(int value) {
	int n = snprintf(NULL, 0, "%d", value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "%d", value);
	return string;
}

char* mlr_alloc_string_from_char_range(char* start, int num_bytes) {
	char* string = mlr_malloc_or_die(num_bytes+1);
	memcpy(string, start, num_bytes);
	string[num_bytes] = 0;
	return string;
}

char* mlr_alloc_hexfmt_from_ll(long  long value) {
	int n = snprintf(NULL, 0, "0x%llx", (unsigned long long)value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "0x%llx", value);
	return string;
}

double mlr_double_from_string_or_die(char* string) {
	double d;
	if (!mlr_try_float_from_string(string, &d)) {
		fprintf(stderr, "%s: couldn't parse \"%s\" as number.\n",
			MLR_GLOBALS.bargv0, string);
		exit(1);
	}
	return d;
}

// E.g. "300" is a number; "300ms" is not.
int mlr_try_float_from_string(char* string, double* pval) {
	int num_bytes_scanned;
	int rc = sscanf(string, "%lf%n", pval, &num_bytes_scanned);
	if (rc != 1)
		return 0;
	if (string[num_bytes_scanned] != 0) // scanned to end of string?
		return 0;
	return 1;
}

long long mlr_int_from_string_or_die(char* string) {
	long long i;
	if (!mlr_try_int_from_string(string, &i)) {
		fprintf(stderr, "Couldn't parse \"%s\" as number.\n", string);
		exit(1);
	}
	return i;
}

// E.g. "300" is a number; "300ms" is not.
int mlr_try_int_from_string(char* string, long long* pval) {
	int num_bytes_scanned, rc;
	// sscanf with %li / %lli doesn't scan correctly when the high bit is set
	// on hex input; it just returns max signed. So we need to special-case hex
	// input.
	if (string[0] == '0' && (string[1] == 'x' || string[1] == 'X')) {
		rc = sscanf(string, "%llx%n", pval, &num_bytes_scanned);
	} else {
		rc = sscanf(string, "%lli%n", pval, &num_bytes_scanned);
	}
	if (rc != 1)
		return 0;
	if (string[num_bytes_scanned] != 0) // scanned to end of string?
		return 0;
	return 1;
}

// ----------------------------------------------------------------
static char* low_int_to_string_data[] = {
	"0",   "1",  "2",  "3",  "4",  "5",  "6",  "7",  "8",  "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "100"
};

char* low_int_to_string(int idx, char* pfree_flags) {
	if ((0 <= idx) && (idx <= 100)) {
		*pfree_flags = 0;
		return low_int_to_string_data[idx];
	} else {
		char buf[32];
		sprintf(buf, "%d", idx);
		*pfree_flags = FREE_ENTRY_KEY;
		return mlr_strdup_or_die(buf);
	}
}

// ----------------------------------------------------------------
char* mlr_paste_2_strings(char* s1, char* s2) {
	int n1 = strlen(s1);
	int n2 = strlen(s2);
	char* s = mlr_malloc_or_die(n1+n2+1);
	strcpy(s, s1);
	strcat(s, s2);
	return s;
}

char* mlr_paste_3_strings(char* s1, char* s2, char* s3) {
	int n1 = strlen(s1);
	int n2 = strlen(s2);
	int n3 = strlen(s3);
	char* s = mlr_malloc_or_die(n1+n2+n3+1);
	strcpy(s, s1);
	strcat(s, s2);
	strcat(s, s3);
	return s;
}

char* mlr_paste_4_strings(char* s1, char* s2, char* s3, char* s4) {
	int n1 = strlen(s1);
	int n2 = strlen(s2);
	int n3 = strlen(s3);
	int n4 = strlen(s4);
	char* s = mlr_malloc_or_die(n1+n2+n3+n4+1);
	strcpy(s, s1);
	strcat(s, s2);
	strcat(s, s3);
	strcat(s, s4);
	return s;
}

char* mlr_paste_5_strings(char* s1, char* s2, char* s3, char* s4, char* s5) {
	int n1 = strlen(s1);
	int n2 = strlen(s2);
	int n3 = strlen(s3);
	int n4 = strlen(s4);
	int n5 = strlen(s5);
	char* s = mlr_malloc_or_die(n1+n2+n3+n4+n5+1);
	strcpy(s, s1);
	strcat(s, s2);
	strcat(s, s3);
	strcat(s, s4);
	strcat(s, s5);
	return s;
}

// ----------------------------------------------------------------
// This is djb2.
int mlr_string_hash_func(char *str) {
	unsigned long hash = 5381;
	int c;

	while ((c = *str++) != 0)
		hash = ((hash << 5) + hash) + c; /* hash * 33 + c */

	return (int)hash;
}

int mlr_string_pair_hash_func(char* str1, char* str2) {
	unsigned long hash = 5381;
	int c;

	while ((c = *str1++) != 0)
		hash = ((hash << 5) + hash) + c; /* hash * 33 + c */
	while ((c = *str2++) != 0)
		hash = ((hash << 5) + hash) + c; /* hash * 33 + c */

	return (int)hash;
}

// ----------------------------------------------------------------
// 0x00-0x7f (MSB is 0) are ASCII and printable.
// 0x80-0xbf (MSBs are 10) are continuation characters and don't add to printable length.
// 0xc0-0xfe (MSBs are 11) are leading characters and do add to printable length.
// (0xff, incidentally, is never a valid UTF-8 byte).
int strlen_for_utf8_display(char* str) {
	int len = 0;
	for (char* p = str; *p; p++) {
		if ((*p & 0xc0) != 0x80)
			len++;
	}
	return len;
}

// ----------------------------------------------------------------
// These are for low-volume, call-at-startup applications. If they get used
// record-by-record they should be replaced with pointer-walking logic which
// avoids the unnecessary expense of calling strlen.

int string_starts_with(char* string, char* prefix) {
	int prefixlen = strlen(prefix);
	return !strncmp(string, prefix, prefixlen);
}

int string_ends_with(char* string, char* suffix, int* pstringlen) {
	int stringlen = strlen(string);
	int suffixlen = strlen(suffix);
	if (pstringlen != NULL)
		*pstringlen = stringlen;
	if (stringlen < suffixlen)
		return FALSE;
	return !strcmp(&string[stringlen-suffixlen], suffix);
}

// ----------------------------------------------------------------
int mlr_imax2(int a, int b) {
	if (a >= b)
		return a;
	else
		return b;
}

// ----------------------------------------------------------------
int power_of_two_above(int n) {
	n |= (n >> 1);
	n |= (n >> 2);
	n |= (n >> 4);
	n |= (n >> 8);
	n |= (n >> 16);
	return(n+1);
}

// ----------------------------------------------------------------
static int is_backslash_octal(char* input, int* pcode) {
	if (strlen(input) < 4)
		return FALSE;
	if (input[0] != '\\')
		return FALSE;
	if (input[1] < '0' || input[1] > '7')
		return FALSE;
	if (input[2] < '0' || input[2] > '7')
		return FALSE;
	if (input[3] < '0' || input[3] > '7')
		return FALSE;
	*pcode = (input[1] - '0') * 64
		+ (input[2] - '0') * 8
		+ (input[3] - '0');
	return TRUE;
}

static int is_backslash_hex(char* input, int* pcode) {
	if (strlen(input) < 4)
		return FALSE;
	if (input[0] != '\\')
		return FALSE;
	if (input[1] != 'x')
		return FALSE;
	if (!isxdigit(input[2]))
		return FALSE;
	if (!isxdigit(input[3]))
		return FALSE;

	char buf[3];
	buf[0] = input[2];
	buf[1] = input[3];
	buf[2] = 0;
	if (sscanf(buf, "%x", pcode) != 1) {
		fprintf(stderr, "Miller: internal coding error detected in file %s at line %d.\n",
			__FILE__, __LINE__);
		exit(1);
	}
	return TRUE;
}

char* mlr_alloc_unbackslash(char* input) {
	// Do the strdup even if there's nothing to expand, so the caller can unconditionally
	// free what we return.
	char* output = mlr_strdup_or_die(input);
	char* pi = input;
	char* po = output;
	int code = 0;
	while (*pi) {
		// https://en.wikipedia.org/wiki/Escape_sequences_in_C
		if (streqn(pi, "\\a", 2)) {
			pi += 2;
			*(po++) = '\a';
		} else if (streqn(pi, "\\b", 2)) {
			pi += 2;
			*(po++) = '\b';
		} else if (streqn(pi, "\\f", 2)) {
			pi += 2;
			*(po++) = '\f';
		} else if (streqn(pi, "\\n", 2)) {
			pi += 2;
			*(po++) = '\n';
		} else if (streqn(pi, "\\r", 2)) {
			pi += 2;
			*(po++) = '\r';
		} else if (streqn(pi, "\\t", 2)) {
			pi += 2;
			*(po++) = '\t';
		} else if (streqn(pi, "\\v", 2)) {
			pi += 2;
			*(po++) = '\v';
		} else if (streqn(pi, "\\\\", 2)) {
			pi += 2;
			*(po++) = '\\';
		} else if (streqn(pi, "\\'", 2)) {
			pi += 2;
			*(po++) = '\'';
		} else if (streqn(pi, "\\\"", 2)) {
			pi += 2;
			*(po++) = '"';
		} else if (streqn(pi, "\\?", 2)) {
			pi += 2;
			*(po++) = '?';
		} else if (is_backslash_octal(pi, &code)) {
			pi += 4;
			*(po++) = code;
		} else if (is_backslash_hex(pi, &code)) {
			pi += 4;
			*(po++) = code;
		} else {
			*po = *pi;
			pi++;
			po++;
		}
	}
	*po = 0;

	return output;
}

// Does a strdup even if there's nothing to expand, so the caller can unconditionally
// free what we return.
char* mlr_alloc_double_backslash(char* input) {
	char *p, *q;
	int input_length = 0;
	int num_backslashes = 0;
	for (p = input; *p; p++) {
		input_length++;
		if (*p == '\\') {
			if (p[1] != '.') {
				num_backslashes++;
			}
		}
	}
	char* output = mlr_malloc_or_die(input_length + num_backslashes + 1);
	for (p = input, q = output; *p; p++) {
		if (*p == '\\') {
			if (p[1] != '.') {
				*(q++) = *p;
			}
			*(q++) = *p;
		} else {
			*(q++) = *p;
		}
	}
	*q = 0;

	return output;
}

// ----------------------------------------------------------------
char* read_file_into_memory(char* filename, size_t* psize) {
	struct stat statbuf;
	if (stat(filename, &statbuf) < 0) {
		perror("stat");
		fprintf(stderr, "%s: could not stat \"%s\"\n", MLR_GLOBALS.bargv0, filename);
		exit(1);
	}
	char* buffer = mlr_malloc_or_die(statbuf.st_size + 1);

	FILE* fp = fopen(filename, "r");
	if (fp == NULL) {
		perror("fopen");
		fprintf(stderr, "%s: could not fopen \"%s\"\n", MLR_GLOBALS.bargv0, filename);
		free(buffer);
		return NULL;
	}

	int rc = fread(buffer, statbuf.st_size, 1, fp);
	if (rc != 1) {
		fprintf(stderr, "Unable t read content of %s\n", filename);
		perror("fread");
		fprintf(stderr, "%s: could not fread \"%s\"\n", MLR_GLOBALS.bargv0, filename);
		fclose(fp);
		free(buffer);
		return NULL;
	}
	fclose(fp);
	buffer[statbuf.st_size] = 0;
	if (psize)
		*psize = statbuf.st_size;
	return buffer;
}

// ----------------------------------------------------------------
#define INITIAL_ALLOC_SIZE 16384
#define BLOCK_SIZE 16384
char* read_fp_into_memory(FILE* fp, size_t* psize) {
	size_t file_size = 0;
	size_t alloc_size = INITIAL_ALLOC_SIZE;
	char* buffer = mlr_malloc_or_die(alloc_size);

	while (TRUE) {
		if (file_size + BLOCK_SIZE > alloc_size) {
			alloc_size *= 2;
			buffer = mlr_realloc_or_die(buffer, alloc_size);
		}
		size_t block_num_bytes_read = fread(&buffer[file_size], 1, BLOCK_SIZE, fp);
		if (block_num_bytes_read == 0) {
			if (feof(fp))
				break;
			perror("fread");
			fprintf(stderr, "%s: stdio/popen fread failed\n", MLR_GLOBALS.bargv0);
			free(buffer);
			*psize = 0;
			return NULL;
		}
		file_size += block_num_bytes_read;
	}

	*psize = file_size;
	return buffer;
}

// ----------------------------------------------------------------
char* alloc_suffixed_temp_file_name(char* filename) {
	const int suffix_length = 6;
	static char bag[] = "abcdefghijklmnopqrstuvwxyz" "ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123456789";
	const static int bag_length = sizeof(bag) - 1;

	char* output = mlr_malloc_or_die(strlen(filename) + 2 + suffix_length);

	int rand_start_index = sprintf(output, "%s.", filename);
	char* rand_start_ptr = &output[rand_start_index];

	int i = 0;
	for ( ; i < suffix_length; i++) {
		rand_start_ptr[i] = bag[get_mtrand_int32() % bag_length];
	}
	rand_start_ptr[i] = 0;

	return output;
}
