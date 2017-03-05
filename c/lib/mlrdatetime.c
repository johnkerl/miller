#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/time.h>
#include <sys/stat.h>
#include "lib/mlr_arch.h"
#include "lib/mlrdatetime.h"

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
