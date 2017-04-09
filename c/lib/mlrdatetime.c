#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include "lib/mlr_globals.h"
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlrdatetime.h"

// ----------------------------------------------------------------
// seconds since the epoch
double get_systime() {
	struct timeval tv = { .tv_sec = 0, .tv_usec = 0 };
	(void)gettimeofday(&tv, NULL);
	return (double)tv.tv_sec + (double)tv.tv_usec * 1e-6;
}

// ----------------------------------------------------------------
#define NZBUFLEN 63
char* mlr_alloc_time_string_from_seconds(double seconds_since_the_epoch, char* format) {
	time_t xxx_temp = (time_t)seconds_since_the_epoch;
	struct tm tm = *gmtime(&xxx_temp); // No gmtime_r on windows
	char* string = mlr_malloc_or_die(NZBUFLEN + 1);
	int written_length = strftime(string, NZBUFLEN, format, &tm);
	if (written_length > NZBUFLEN || written_length == 0) {
		fprintf(stderr, "%s: could not strftime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, string, format, MLR_GLOBALS.bargv0);
		exit(1);
	}

	return string;
}

// ----------------------------------------------------------------
double mlr_seconds_from_time_string(char* string, char* format) {
	struct tm tm;
	memset(&tm, 0, sizeof(tm));
	char* retval = mlr_arch_strptime(string, format, &tm);
	if (retval == NULL) {
		fprintf(stderr, "%s: could not strptime(\"%s\", \"%s\"). See \"%s --help-function strptime\".\n",
			MLR_GLOBALS.bargv0, string, format, MLR_GLOBALS.bargv0);
		exit(1);
	}
	MLR_INTERNAL_CODING_ERROR_IF(*retval != 0); // Parseable input followed by non-parseable
	time_t iseconds = mlr_arch_timegm(&tm);
	return (double)iseconds;
}
