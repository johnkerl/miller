#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include <sys/time.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"

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
#ifdef MLR_MALLOC_TRACE
	fprintf(stderr, "MALLOC size=%d,p=%p\n", (int)size, p);
#endif
	return p;
}

// ----------------------------------------------------------------
void* mlr_realloc_or_die(void *optr, size_t size) {
	void* nptr = realloc(optr, size);
	if (nptr == NULL) {
		fprintf(stderr, "realloc(%lu) failed.\n", (unsigned long)size);
		exit(1);
	}
#ifdef MLR_MALLOC_TRACE
	fprintf(stderr, "REALLOC size=%d,p=%p\n", (int)size, nptr);
#endif
	return nptr;
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
			MLR_GLOBALS.argv0, string);
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
// See the GNU timegm manpage -- this is what it does.
time_t mlr_timegm(struct tm* tm) {
	time_t ret;
	char* tz;

	tz = getenv("TZ");
	setenv("TZ", "GMT0", 1);
	tzset();
	ret = mktime(tm);
	if (tz) {
		setenv("TZ", tz, 1);
	} else {
		unsetenv("TZ");
	}
	tzset();
	return ret;
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
// This is inefficient. It's quite fine for call-once, small-n use.

int power_of_two_ceil(int n) {
	while (n&(n-1))
		n++;
	return n;
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

char* mlr_unbackslash(char* input) {
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

// ----------------------------------------------------------------
// Succeeds or aborts the process. cflag REG_EXTENDED is already included.
regex_t* regcomp_or_die(regex_t* pregex, char* regex_string, int cflags) {
	cflags |= REG_EXTENDED;
	int rc = regcomp(pregex, regex_string, cflags);
	if (rc != 0) {
		size_t nbytes = regerror(rc, pregex, NULL, 0);
		char* errbuf = malloc(nbytes);
		(void)regerror(rc, pregex, errbuf, nbytes);
		fprintf(stderr, "%s: could not compile regex \"%s\" : %s\n",
			MLR_GLOBALS.argv0, regex_string, errbuf);
		exit(1);
	}
	return pregex;
}

// Always uses cflags with REG_EXTENDED.
// If the regex_string is of the form a.*b, compiles it using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b", compiles a.*b using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b"i, compiles a.*b using cflags with REG_ICASE.
regex_t* regcomp_or_die_quoted(regex_t* pregex, char* orig_regex_string, int cflags) {
	cflags |= REG_EXTENDED;
	if (string_starts_with(orig_regex_string, "\"")) {
		char* regex_string = mlr_strdup_or_die(orig_regex_string);
		int len = 0;
		if (string_ends_with(regex_string, "\"", &len)) {
			regex_string[len-1] = 0;
		} else if (string_ends_with(regex_string, "\"i", &len)) {
			regex_string[len-2] = 0;
			cflags |= REG_ICASE;
		} else {
			fprintf(stderr, "%s: imbalanced double-quote in regex [%s].\n",
				MLR_GLOBALS.argv0, regex_string);
			exit(1);
		}
		regcomp_or_die(pregex, regex_string+1, cflags);
		free(regex_string);
	} else {
		regcomp_or_die(pregex, orig_regex_string, cflags);
	}
	return pregex;
}

// Returns TRUE for match, FALSE for no match, and aborts the process if
// regexec returns anything else.
int regmatch_or_die(const regex_t* pregex, const char* restrict match_string,
	size_t nmatch, regmatch_t pmatch[restrict])
{
	int rc = regexec(pregex, match_string, nmatch, pmatch, 0);
	if (rc == 0) {
		return TRUE;
	} else if (rc == REG_NOMATCH) {
		return FALSE;
	} else {
		size_t nbytes = regerror(rc, pregex, NULL, 0);
		char* errbuf = malloc(nbytes);
		(void)regerror(rc, pregex, errbuf, nbytes);
		printf("regexec failure: %s\n", errbuf);
		exit(1);
	}
}

// If there is a match, the return value is dynamically allocated.  If not, the
// input is returned.
//
// Capture-group example:
// sed: $ echo '<<abcdefg>>'|sed 's/ab\(.\)d\(..\)g/AYEBEE\1DEE\2GEE/' gives <<AYEBEEcDEEefGEE>>
// mlr: echo 'x=<<abcdefg>>' | mlr put '$x = sub($x, "ab(.)d(..)g", "AYEBEE\1DEE\2GEE")' x=<<AYEBEEcDEEefGEE>>

char* regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement, int* pmatched, int *pall_captured) {
	const size_t nmatch = 10; // Capture-groups \1 through \9 supported, along with entire-string match
	regmatch_t matches[nmatch];
	*pall_captured = TRUE;

	*pmatched = regmatch_or_die(pregex, input, nmatch, matches);
	if (!*pmatched) {
		return input;
	} else {
		sb_append_chars(psb, input, 0, matches[0].rm_so-1);
		char* p = replacement;
		while (*p) {
			if (p[0] == '\\' && isdigit(p[1])) {
				int idx = p[1] - '0';
				regmatch_t* pmatch = &matches[idx];
				if (pmatch->rm_so == -1) {
					*pall_captured = FALSE;
					sb_append_chars(psb, p, 0, 1);
				} else {
					sb_append_chars(psb, input, matches[idx].rm_so, matches[idx].rm_eo-1);
				}
				p += 2;
			} else {
				sb_append_char(psb, *p);
				p++;
			}
		}
		sb_append_chars(psb, input, matches[0].rm_eo, strlen(input));

		return sb_finish(psb);
	}
}

char* regex_gsub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement, int *pmatched, int* pall_captured) {
	const size_t nmatch = 10;
	regmatch_t matches[nmatch];
	*pmatched = FALSE;
	*pall_captured = TRUE;

	int   match_start = 0;
	char* current_input = input;

	while (TRUE) {
		int matched = regmatch_or_die(pregex, &current_input[match_start], nmatch, matches);
		if (!matched) {
			return current_input;
		}
		*pmatched = TRUE;

		sb_append_chars(psb, current_input, 0, match_start + matches[0].rm_so-1);

		char* p = replacement;
		int len1 = psb->used_length;
		while (*p) {
			if (p[0] == '\\' && isdigit(p[1])) {
				int idx = p[1] - '0';
				regmatch_t* pmatch = &matches[idx];
				if (pmatch->rm_so == -1) {
					*pall_captured = FALSE;
					sb_append_chars(psb, p, 0, 1);
				} else {
					sb_append_chars(psb, &current_input[match_start], matches[idx].rm_so, matches[idx].rm_eo-1);
				}
				p += 2;
			} else {
				sb_append_char(psb, *p);
				p++;
			}
		}

		int replen = psb->used_length - len1;
		sb_append_chars(psb, current_input, match_start + matches[0].rm_eo, strlen(current_input));

		current_input = sb_finish(psb);

		match_start += matches[0].rm_so + replen;
	}
}
