#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>

// For some Linux distros, in spite of including time.h:
char *strptime(const char *s, const char *format, struct tm *ptm);

// ----------------------------------------------------------------
char * mlr_strptime(const char *buf, const char *format, struct tm *timeptr);
size_t mlr_strftime(char *restrict s, size_t maxsize, const char *restrict format,
	const struct tm *restrict timeptr);

// ----------------------------------------------------------------
#define ISO8601_TIME_FORMAT "%Y-%m-%dT%H:%M:%SZ"

// ----------------------------------------------------------------
// seconds since the epoch
double get_systime() {
	struct timeval tv = { .tv_sec = 0, .tv_usec = 0 };
	(void)gettimeofday(&tv, NULL);
	return (double)tv.tv_sec + (double)tv.tv_usec * 1e-6;
}

// ----------------------------------------------------------------
// See the GNU timegm manpage -- this is what it does.
time_t mlr_timegm(struct tm* ptm) {
	time_t ret;
	char* tz;

	tz = getenv("TZ");
	setenv("TZ", "GMT0", 1);
	tzset();
	ret = mktime(ptm);
	if (tz) {
		setenv("TZ", tz, 1);
	} else {
		unsetenv("TZ");
	}
	tzset();
	return ret;
}

// ----------------------------------------------------------------
#define NZBUFLEN 63
char* mlr_alloc_time_string_from_seconds(time_t seconds, char* format, int use_new) {
	struct tm tm;
	struct tm *ptm = gmtime_r(&seconds, &tm);
	char* string = malloc(NZBUFLEN + 1);
	int written_length = use_new
		? mlr_strftime(string, NZBUFLEN, format, ptm)
		: strftime(string, NZBUFLEN, format, ptm);
	if (written_length > NZBUFLEN || written_length == 0) {
		fprintf(stderr, "(new=%d) Could not strftime(\"%s\", \"%s\").\n", use_new, string, format);
	}

	return string;
}

// ----------------------------------------------------------------
time_t mlr_seconds_from_time_string(char* string, char* format, int use_new) {
	struct tm tm;
	memset(&tm, 0, sizeof(tm));
	char* retval = use_new
		? mlr_strptime(string, format, &tm)
		: strptime(string, format, &tm);
	if (retval == NULL) {
		fprintf(stderr, "(new=%d) Could not strptime(\"%s\", \"%s\").\n", use_new, string, format);
	}
	return mlr_timegm(&tm);
}

// ----------------------------------------------------------------
int main(int argc, char** argv) {
	time_t t = 1491399761;
	char* format = ISO8601_TIME_FORMAT;
	if (argc == 3) {
		int i;
		if (sscanf(argv[1], "%d", &i) != 1) {
			fprintf(stderr, "b04k!\n");
			exit(1);
		}
		t = (time_t) i;
		format = argv[2];
	}
	char*  f0 = mlr_alloc_time_string_from_seconds(t, format, 0);
	char*  f1 = mlr_alloc_time_string_from_seconds(t, format, 1);
	time_t p0 = mlr_seconds_from_time_string(f0, format, 0);
	time_t p1 = mlr_seconds_from_time_string(f0, format, 1);
	printf("t  = %d\n", (int)t);
	printf("f0 = %s\n", f0);
	printf("f1 = %s\n", f1);
	printf("p0 = %d\n", (int)p0);
	printf("p1 = %d\n", (int)p1);
	return 0;
}
