#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include "lib/mlr_globals.h"
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlrdatetime.h"

// For some Linux distros, in spite of including time.h:
char *strptime(const char *s, const char *format, struct tm *ptm);

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
	mlr_arch_setenv("TZ", "GMT0");
	tzset();
	ret = mktime(ptm);
	if (tz) {
		mlr_arch_setenv("TZ", tz);
	} else {
		mlr_arch_unsetenv("TZ");
	}
	tzset();
	return ret;
}

// ----------------------------------------------------------------
#define NZBUFLEN 63
char* mlr_alloc_time_string_from_seconds(time_t seconds, char* format) {
	struct tm tm;
	struct tm *ptm = gmtime_r(&seconds, &tm);
	MLR_INTERNAL_CODING_ERROR_IF(ptm == NULL);
	char* string = mlr_malloc_or_die(NZBUFLEN + 1);
	int written_length = strftime(string, NZBUFLEN, format, ptm);
	if (written_length > NZBUFLEN || written_length == 0) {
		fprintf(stderr, "%s: could not strftime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, string, format, MLR_GLOBALS.bargv0);
		exit(1);
	}

	return string;
}

// ----------------------------------------------------------------
time_t mlr_seconds_from_time_string(char* string, char* format) {
	struct tm tm;
	memset(&tm, 0, sizeof(tm));
	char* retval = strptime(string, format, &tm);
	if (retval == NULL) {
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, string, format, MLR_GLOBALS.bargv0);
		exit(1);
	}
	MLR_INTERNAL_CODING_ERROR_IF(*retval != 0); // Parseable input followed by non-parseable
	return mlr_timegm(&tm);
}
