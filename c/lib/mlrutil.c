#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include "lib/mlrutil.h"

//// ----------------------------------------------------------------
// inlined in mrlutil.h
//int mlr_canonical_mod(int a, int n) {
//	int r = a % n;
//	if (r >= 0)
//		return r;
//	else
//		return r+n;
//}

// ----------------------------------------------------------------
// xxx cmt top insert ...
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
// seconds since the epoch
double get_systime() {
	struct timeval tv = { .tv_sec = 0, .tv_usec = 0 };
	(void)gettimeofday(&tv, NULL);
	return (double)tv.tv_sec + (double)tv.tv_usec * 1e-6;
}

// ----------------------------------------------------------------
void* mlr_malloc_or_die(size_t size) {
	void* p = malloc(size);
	if (p == NULL) {
		fprintf(stderr, "malloc(%lu) failed.\n", (unsigned long)size);
		exit(1);
	}

	//memset(p, 0, size); // xxx temp memset

	return p;
}

// ----------------------------------------------------------------
// xxx cmt mem mgt
// xxx use stack buf & avoid double calls to the formatter
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
	int n = snprintf(NULL, 0, "%lli", value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "%lli", value);
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

char* mlr_alloc_hexfmt_from_ll(long  long value) {
	int n = snprintf(NULL, 0, "0x%llx", (unsigned long long)value);
	char* string = mlr_malloc_or_die(n+1);
	sprintf(string, "0x%llx", value);
	return string;
}

double mlr_double_from_string_or_die(char* string) {
	double d;
	if (!mlr_try_double_from_string(string, &d)) {
		fprintf(stderr, "Couldn't parse \"%s\" as number.\n", string);
		exit(1);
	}
	return d;
}

// E.g. "300" is a number; "300ms" is not.
int mlr_try_double_from_string(char* string, double* pval) {
	int num_bytes_scanned;
	int rc = sscanf(string, "%lf%n", pval, &num_bytes_scanned);
	if (rc != 1)
		return 0;
	if (string[num_bytes_scanned] != 0) // scanned to end of string?
		return 0;
	return 1;
}

// E.g. "300" is a number; "300ms" is not.
int mlr_try_int_from_string(char* string, long long* pval) {
	int num_bytes_scanned;
	int rc = sscanf(string, "%lli%n", pval, &num_bytes_scanned);
	if (rc != 1)
		return 0;
	if (string[num_bytes_scanned] != 0) // scanned to end of string?
		return 0;
	return 1;
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
// Found on the web.
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
char* mlr_get_line(FILE* input_stream, char rs) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = getdelim(&line, &linecap, rs, input_stream);
	if (linelen <= 0) {
		return NULL;
	}
	if (line[linelen-1] == '\n') { // chomp
		line[linelen-1] = 0;
		linelen--;
	}

	return line;
}

// ----------------------------------------------------------------
// See the GNU timegm manpage -- this is what it does.
time_t mlr_timegm (struct tm* tm) {
	time_t ret;
	char* tz;

	tz = getenv("TZ");
	if (tz) {
		tz = strdup(tz);
	}
	setenv("TZ", "GMT0", 1);
	tzset();
	ret = mktime(tm);
	if (tz) {
		setenv("TZ", tz, 1);
		free(tz);
	} else {
		unsetenv("TZ");
	}
	tzset();
	return ret;
}
