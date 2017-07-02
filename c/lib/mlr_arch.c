#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mlr_globals.h"
#include "mlr_arch.h"
#include "mlrutil.h"
#include "netbsd_strptime.h"
#include "nlnet_timegm.h"

// For some Linux distros, in spite of including time.h:
char *strptime(const char *s, const char *format, struct tm *ptm);

// ----------------------------------------------------------------
int mlr_arch_setenv(const char *name, const char *value) {
#ifdef MLR_ON_MSYS2
	fprintf(stderr, "%s: setenv is not supported on this architecture.\n", MLR_GLOBALS.bargv0);
	exit(1);
#else
	return setenv(name, value, 1 /*overwrite*/);
#endif
}

// ----------------------------------------------------------------
int mlr_arch_unsetenv(const char *name) {
#ifdef MLR_ON_MSYS2
	fprintf(stderr, "%s: unsetenv is not supported on this architecture.\n", MLR_GLOBALS.bargv0);
	exit(1);
#else
	return unsetenv(name);
#endif
}

// ----------------------------------------------------------------
char *mlr_arch_strptime(const char *s, const char *format, struct tm *ptm) {
#ifdef MLR_ON_MSYS2
	return netbsd_strptime(s, format, ptm);
#else
	return strptime(s, format, ptm);
#endif
}

// ----------------------------------------------------------------
// See the GNU timegm manpage -- this is what it does.
time_t mlr_arch_timegm(struct tm* ptm) {
#ifdef MLR_ON_MSYS2
	return nlnet_timegm(ptm);
#else
	time_t ret;
	char* tz = getenv("TZ");
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
#endif
}
